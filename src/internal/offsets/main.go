// Command offsets finds struct fields that legacy code accesses by byte offset (unsafe.Add(base, N), or a helper
// that wraps it) and checks if the offset still matches the field on amd64.
//
// The base struct type is found from the base expression: x.C()-style methods that return unsafe.Pointer(x),
// unsafe.Pointer(x) conversions, and local variables assigned once from such expressions. Bases that are untyped
// unsafe.Pointer fields (e.g. Object.UpdateData) are reported by field name.
//
// Usage (from src/):
//
//	go run ./internal/offsets [-v] [packages]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var (
	fVerbose  = flag.Bool("v", false, "list every site")
	fFix      = flag.Bool("fix", false, "translate the offsets of sites with a known base type through common/layout")
	fTyped    = flag.String("typed", "", "rewrite offset accesses to typed field accesses; the file lists probed base types (site<TAB>type,...)")
	fTypedLog = flag.String("typedlog", "typed.log", "report of -typed")
	fDry      = flag.Bool("dry", false, "with -typed: report only, do not rewrite")
	fOverflow = flag.Bool("overflow", false, "list reinterpretations that overflow their storage on amd64")
	fProbe    = flag.Bool("probe", false, "instrument every offset access with layout.Probe (for a copy of the tree only)")
)

var (
	s386   = types.SizesFor("gc", "386")
	sAmd64 = types.SizesFor("gc", "amd64")
)

// Site is one access by offset.
type Site struct {
	Pos    string
	Base   string // struct type, "field:Name" for unsafe.Pointer fields, or "?" if unknown
	Off    int64
	Access string // the type the access reads or writes, if known
	Field  string // field path at the 386 offset
	Off64  int64  // offset of that field on amd64
	Moved  bool
	Inside bool // the offset points inside a field, not at its start
	Expr   string
	Helper string
	Width  bool // the access width differs from the field width on amd64

	file       string
	start, end int    // byte range of the offset expression
	offText    string // source of the offset expression
	offType    string // type of the offset expression
	typ        types.Type
	pkg        *types.Package
}

type pkgCtx struct {
	p       *packages.Package
	helpers map[*types.Func]helper
}

// helper is a function whose body is "return (*X)(unsafe.Add(p, off))", with p and off parameters.
type helper struct {
	baseParam, offParam int
	access              types.Type
}

func main() {
	flag.Parse()
	pats := flag.Args()
	if len(pats) == 0 {
		pats = []string{"./..."}
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Env:  append(os.Environ(), "GOARCH=386", "GOOS=linux", "CGO_ENABLED=0"),
	}
	pkgs, err := packages.Load(cfg, pats...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	wd, _ := os.Getwd()
	if *fTyped != "" {
		lf, err := os.Create(*fTypedLog)
		if err != nil {
			panic(err)
		}
		defer lf.Close()
		env := &typedEnv{dry: *fDry, byName: namedTypes(pkgs), probed: loadProbed(*fTyped), wd: wd, log: bufio.NewWriter(lf), stats: make(map[string]int)}
		env.run(pkgs)
		return
	}
	if *fOverflow {
		overflows(pkgs, wd)
		allocSizes(pkgs, wd)
		return
	}
	if *fProbe {
		edits := make(map[string][]insert)
		n := 0
		for _, p := range pkgs {
			if p.TypesInfo != nil {
				n += (&pkgCtx{p: p, helpers: findHelpers(p)}).probes(wd, edits)
			}
		}
		applyProbes(edits)
		fmt.Printf("# instrumented %d sites in %d files\n", n, len(edits))
		return
	}
	var sites []Site
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		c := &pkgCtx{p: p, helpers: findHelpers(p)}
		sites = append(sites, c.scan(wd)...)
	}
	report(sites)
	if *fFix {
		applyFixes(sites)
	}
}

func isUnsafeAdd(info *types.Info, call *ast.CallExpr) bool {
	sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Add" {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	pn, ok := info.Uses[id].(*types.PkgName)
	return ok && pn.Imported().Path() == "unsafe"
}

func paramIndex(info *types.Info, fn *types.Signature, e ast.Expr) int {
	id, ok := ast.Unparen(e).(*ast.Ident)
	if !ok {
		return -1
	}
	obj := info.Uses[id]
	for i := 0; i < fn.Params().Len(); i++ {
		if fn.Params().At(i) == obj {
			return i
		}
	}
	return -1
}

// rawAtHelper recognizes helpers whose body is "return rawAt[X](p, off)" (legacy/raw_offsets.go).
func rawAtHelper(info *types.Info, fd *ast.FuncDecl, call *ast.CallExpr) (helper, *types.Func, bool) {
	ix, ok := call.Fun.(*ast.IndexExpr)
	if !ok || len(call.Args) != 2 {
		return helper{}, nil, false
	}
	if id, ok := ix.X.(*ast.Ident); !ok || id.Name != "rawAt" {
		return helper{}, nil, false
	}
	tv, ok := info.Types[ix.Index]
	if !ok || !tv.IsType() {
		return helper{}, nil, false
	}
	fn, ok := info.Defs[fd.Name].(*types.Func)
	if !ok {
		return helper{}, nil, false
	}
	sig := fn.Type().(*types.Signature)
	bi := paramIndex(info, sig, call.Args[0])
	oi := paramIndex(info, sig, call.Args[1])
	if oi < 0 {
		if c, ok := ast.Unparen(call.Args[1]).(*ast.CallExpr); ok && len(c.Args) == 1 {
			oi = paramIndex(info, sig, c.Args[0])
		}
	}
	if bi < 0 || oi < 0 {
		return helper{}, nil, false
	}
	return helper{baseParam: bi, offParam: oi, access: tv.Type}, fn, true
}

func findHelpers(p *packages.Package) map[*types.Func]helper {
	out := make(map[*types.Func]helper)
	info := p.TypesInfo
	for _, f := range p.Syntax {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil || len(fd.Body.List) != 1 {
				continue
			}
			ret, ok := fd.Body.List[0].(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				continue
			}
			conv, ok := ast.Unparen(ret.Results[0]).(*ast.CallExpr)
			if !ok {
				continue
			}
			if h, fn, ok := rawAtHelper(info, fd, conv); ok {
				out[fn] = h
				continue
			}
			if len(conv.Args) != 1 {
				continue
			}
			tv, ok := info.Types[conv.Fun]
			if !ok || !tv.IsType() {
				continue
			}
			pt, ok := tv.Type.Underlying().(*types.Pointer)
			if !ok {
				continue
			}
			add, ok := ast.Unparen(conv.Args[0]).(*ast.CallExpr)
			if !ok || !isUnsafeAdd(info, add) {
				continue
			}
			fn, ok := info.Defs[fd.Name].(*types.Func)
			if !ok {
				continue
			}
			sig := fn.Type().(*types.Signature)
			// base: a parameter, or param.C() (e.g. dr.C())
			bi := paramIndex(info, sig, add.Args[0])
			if bi < 0 {
				if c, ok := ast.Unparen(add.Args[0]).(*ast.CallExpr); ok {
					if s, ok := c.Fun.(*ast.SelectorExpr); ok && len(c.Args) == 0 {
						bi = paramIndex(info, sig, s.X)
					}
				}
			}
			oi := paramIndex(info, sig, add.Args[1])
			if oi < 0 {
				// off may be converted: uintptr(off)
				if c, ok := ast.Unparen(add.Args[1]).(*ast.CallExpr); ok && len(c.Args) == 1 {
					oi = paramIndex(info, sig, c.Args[0])
				}
			}
			if bi < 0 || oi < 0 {
				continue
			}
			out[fn] = helper{baseParam: bi, offParam: oi, access: pt.Elem()}
		}
	}
	return out
}

func constInt(info *types.Info, e ast.Expr) (int64, bool) {
	tv, ok := info.Types[e]
	if !ok || tv.Value == nil {
		return 0, false
	}
	return constant.Int64Val(constant.ToInt(tv.Value))
}

// assignments maps local variables to their only defining expression.
func assignments(info *types.Info, body *ast.BlockStmt) map[types.Object]ast.Expr {
	defs := make(map[types.Object]ast.Expr)
	count := make(map[types.Object]int)
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.AssignStmt:
			if len(n.Lhs) != len(n.Rhs) {
				for _, l := range n.Lhs {
					if id, ok := l.(*ast.Ident); ok {
						if o := objOf(info, id); o != nil {
							count[o] += 2
						}
					}
				}
				return true
			}
			for i, l := range n.Lhs {
				if id, ok := l.(*ast.Ident); ok {
					if o := objOf(info, id); o != nil {
						count[o]++
						defs[o] = n.Rhs[i]
					}
				}
			}
		case *ast.ValueSpec:
			for i, id := range n.Names {
				if o := info.Defs[id]; o != nil && i < len(n.Values) {
					count[o]++
					defs[o] = n.Values[i]
				}
			}
		}
		return true
	})
	for o, n := range count {
		if n != 1 {
			delete(defs, o)
		}
	}
	return defs
}

func objOf(info *types.Info, id *ast.Ident) types.Object {
	if o := info.Defs[id]; o != nil {
		return o
	}
	return info.Uses[id]
}

// baseType finds the struct type that the base expression points to.
func baseType(info *types.Info, defs map[types.Object]ast.Expr, e ast.Expr, depth int) (types.Type, string) {
	if depth > 6 {
		return nil, "?"
	}
	e = ast.Unparen(e)
	if t := info.TypeOf(e); t != nil {
		if pt, ok := t.Underlying().(*types.Pointer); ok {
			if _, ok := pt.Elem().Underlying().(*types.Struct); ok {
				return pt.Elem(), ""
			}
		}
	}
	switch e := e.(type) {
	case *ast.CallExpr:
		if tv, ok := info.Types[e.Fun]; ok && tv.IsType() && len(e.Args) == 1 {
			return baseType(info, defs, e.Args[0], depth+1) // unsafe.Pointer(x)
		}
		if s, ok := e.Fun.(*ast.SelectorExpr); ok && len(e.Args) == 0 {
			// x.C(), x.CObj(), ...: methods returning the receiver as unsafe.Pointer
			if fn, ok := info.Uses[s.Sel].(*types.Func); ok {
				if sig := fn.Type().(*types.Signature); sig.Recv() != nil {
					if b := basicKind(sig.Results()); b == types.UnsafePointer {
						name := s.Sel.Name
						if name == "C" || name == "CObj" || name == "SObj" || name == "UPtr" || strings.HasPrefix(name, "C") && len(name) <= 4 {
							return baseType(info, defs, s.X, depth+1)
						}
					}
				}
			}
		}
	case *ast.SelectorExpr:
		if v, ok := info.Uses[e.Sel].(*types.Var); ok && v.IsField() {
			if b, ok := v.Type().Underlying().(*types.Basic); ok && b.Kind() == types.UnsafePointer {
				return nil, "field:" + e.Sel.Name
			}
		}
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			if t := info.TypeOf(e.X); t != nil {
				if _, ok := t.Underlying().(*types.Struct); ok {
					return t, ""
				}
			}
		}
	case *ast.Ident:
		if o := info.Uses[e]; o != nil {
			if d, ok := defs[o]; ok {
				return baseType(info, defs, d, depth+1)
			}
			if _, ok := o.(*types.Var); ok && o.Parent() != nil && o.Parent() != o.Pkg().Scope() {
				return nil, "param/var " + e.Name
			}
		}
	}
	return nil, "?"
}

func basicKind(t *types.Tuple) types.BasicKind {
	if t.Len() != 1 {
		return types.Invalid
	}
	if b, ok := t.At(0).Type().Underlying().(*types.Basic); ok {
		return b.Kind()
	}
	return types.Invalid
}

type leaf struct {
	path           string
	off32, off64   int64
	size32, size64 int64
}

func leaves(t types.Type, prefix string, o32, o64 int64, out *[]leaf) {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		var fs []*types.Var
		for i := 0; i < u.NumFields(); i++ {
			fs = append(fs, u.Field(i))
		}
		a, b := s386.Offsetsof(fs), sAmd64.Offsetsof(fs)
		for i, f := range fs {
			n := f.Name()
			if prefix != "" {
				n = prefix + "." + n
			}
			leaves(f.Type(), n, o32+a[i], o64+b[i], out)
		}
	case *types.Array:
		e32, e64 := s386.Sizeof(u.Elem()), sAmd64.Sizeof(u.Elem())
		if _, isStruct := u.Elem().Underlying().(*types.Struct); !isStruct && e32 == e64 {
			*out = append(*out, leaf{prefix, o32, o64, s386.Sizeof(t), sAmd64.Sizeof(t)})
			return
		}
		for i := int64(0); i < u.Len() && i < 4096; i++ {
			leaves(u.Elem(), fmt.Sprintf("%s[%d]", prefix, i), o32+i*e32, o64+i*e64, out)
		}
	default:
		*out = append(*out, leaf{prefix, o32, o64, s386.Sizeof(t), sAmd64.Sizeof(t)})
	}
}

var leafCache = make(map[types.Type][]leaf)

func fieldAt(t types.Type, off int64) (leaf, bool, bool) {
	ls, ok := leafCache[t]
	if !ok {
		leaves(t, "", 0, 0, &ls)
		leafCache[t] = ls
	}
	for _, l := range ls {
		if off >= l.off32 && off < l.off32+max(l.size32, 1) {
			return l, true, off != l.off32
		}
	}
	return leaf{}, false, false
}

func (c *pkgCtx) scan(wd string) []Site {
	p := c.p
	info := p.TypesInfo
	var out []Site
	insp := inspector.New(p.Syntax)
	var defs map[types.Object]ast.Expr
	var hints map[string]types.Type
	insp.WithStack([]ast.Node{(*ast.FuncDecl)(nil), (*ast.FuncLit)(nil), (*ast.CallExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Body != nil {
				defs = assignments(info, n.Body)
				hints = dataHints(info, n.Body)
			}
			return true
		case *ast.FuncLit:
			return true
		}
		call := n.(*ast.CallExpr)
		var base, offE ast.Expr
		var access types.Type
		helperName := ""
		if isUnsafeAdd(info, call) {
			base, offE = call.Args[0], call.Args[1]
			if len(stack) >= 2 {
				if conv, ok := stack[len(stack)-2].(*ast.CallExpr); ok {
					if tv, ok := info.Types[conv.Fun]; ok && tv.IsType() {
						if pt, ok := tv.Type.Underlying().(*types.Pointer); ok {
							access = pt.Elem()
						}
					}
				}
			}
		} else if id := calleeFunc(info, call); id != nil {
			h, ok := c.helpers[id]
			if !ok || len(call.Args) <= max(h.baseParam, h.offParam) {
				return true
			}
			base, offE, access = call.Args[h.baseParam], call.Args[h.offParam], h.access
			helperName = id.Name()
		} else {
			return true
		}
		off, ok := constInt(info, offE)
		if !ok {
			return true
		}
		pos := p.Fset.Position(call.Pos())
		rel, _ := filepath.Rel(wd, pos.Filename)
		s := Site{Pos: fmt.Sprintf("%s:%d", rel, pos.Line), Off: off, Helper: helperName, Expr: types.ExprString(base), pkg: p.Types}
		f := p.Fset.File(offE.Pos())
		s.file, s.start, s.end = f.Name(), f.Offset(offE.Pos()), f.Offset(offE.End())
		if src, err := os.ReadFile(s.file); err == nil {
			s.offText = string(src[s.start:s.end])
		}
		if t := info.TypeOf(offE); t != nil {
			s.offType = types.TypeString(t, func(p *types.Package) string { return p.Name() })
		}
		if access != nil {
			s.Access = types.TypeString(access, func(p *types.Package) string { return p.Name() })
		}
		t, why := baseType(info, defs, base, 0)
		if t == nil && strings.HasPrefix(why, "field:") {
			if h, ok := hints[strings.TrimPrefix(why, "field:")]; ok && h != nil {
				t = h
			}
		}
		if t == nil {
			s.Base = why
		} else {
			s.Base = types.TypeString(t, func(p *types.Package) string { return p.Name() })
			s.typ = t
			if l, ok, inside := fieldAt(t, off); ok {
				s.Field, s.Off64, s.Inside = l.path, l.off64+(off-l.off32), inside
				s.Moved = s.Off64 != off || l.size32 != l.size64
				if access != nil && l.size32 != l.size64 && sAmd64.Sizeof(access) != l.size64 {
					s.Width = true
				}
			} else {
				s.Field = "(beyond struct)"
				s.Moved = true
			}
		}
		out = append(out, s)
		return true
	})
	return out
}

func calleeFunc(info *types.Info, call *ast.CallExpr) *types.Func {
	switch f := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		fn, _ := info.Uses[f].(*types.Func)
		return fn
	case *ast.SelectorExpr:
		fn, _ := info.Uses[f.Sel].(*types.Func)
		return fn
	}
	return nil
}

func report(sites []Site) {
	sort.Slice(sites, func(i, j int) bool { return sites[i].Pos < sites[j].Pos })
	var typed, moved, unknown int
	byBase := make(map[string][2]int)
	for _, s := range sites {
		if strings.HasPrefix(s.Base, "?") || strings.HasPrefix(s.Base, "field:") || strings.HasPrefix(s.Base, "param") {
			unknown++
			key := s.Base
			if strings.HasPrefix(key, "param") {
				key = "param/var"
			}
			v := byBase[key]
			v[0]++
			byBase[key] = v
			continue
		}
		typed++
		v := byBase[s.Base]
		v[0]++
		if s.Moved {
			moved++
			v[1]++
		}
		byBase[s.Base] = v
	}
	fmt.Printf("# sites: %d (base type known %d, of those moved on amd64 %d; unknown base %d)\n", len(sites), typed, moved, unknown)
	var keys []string
	for k := range byBase {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return byBase[keys[i]][0] > byBase[keys[j]][0] })
	for _, k := range keys {
		fmt.Printf("%6d %6d  %s\n", byBase[k][0], byBase[k][1], k)
	}
	if *fVerbose {
		for _, s := range sites {
			m := ""
			if s.Moved {
				m = "MOVED"
			}
			fmt.Printf("%-50s %-28s %6d -> %-6d %-5s %-34s %-14s %s %s\n", s.Pos, s.Base, s.Off, s.Off64, m, s.Field, s.Access, s.Helper, s.Expr)
		}
	}
}

// tableVar returns the name of the package variable holding the offset table of a type.
func tableVar(t types.Type) string {
	n, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	name := n.Obj().Name()
	return "layout" + strings.ToUpper(n.Obj().Pkg().Name()[:1]) + n.Obj().Pkg().Name()[1:] + strings.ToUpper(name[:1]) + name[1:]
}

func applyFixes(sites []Site) {
	type edit struct {
		start, end int
		text       string
	}
	byFile := make(map[string][]edit)
	vars := make(map[string]map[string]string) // dir -> var -> type expression
	pkgNames := make(map[string]string)        // dir -> package name
	imports := make(map[string]map[string]bool)
	var width []Site
	n := 0
	for _, s := range sites {
		if s.typ == nil || !s.Moved || s.Field == "(beyond struct)" {
			continue
		}
		if s.Width {
			width = append(width, s)
		}
		v := tableVar(s.typ)
		if v == "" {
			continue
		}
		nt := s.typ.(*types.Named)
		dir := filepath.Dir(s.file)
		if vars[dir] == nil {
			vars[dir] = make(map[string]string)
			imports[dir] = make(map[string]bool)
		}
		pkgNames[dir] = s.pkg.Name()
		texpr := nt.Obj().Name()
		if nt.Obj().Pkg() != s.pkg {
			texpr = nt.Obj().Pkg().Name() + "." + texpr
			imports[dir][nt.Obj().Pkg().Path()] = true
		}
		vars[dir][v] = texpr
		call := v + ".Off(" + s.offText + ")"
		if s.offType != "uintptr" && s.offType != "untyped int" && s.offType != "" {
			call = s.offType + "(" + call + ")"
		}
		byFile[s.file] = append(byFile[s.file], edit{s.start, s.end, call + " /* " + s.Field + " */"})
		n++
	}
	for file, es := range byFile {
		sort.Slice(es, func(i, j int) bool { return es[i].start > es[j].start })
		data, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		src := string(data)
		last := len(src) + 1
		for _, e := range es {
			if e.end > last {
				continue
			}
			src = src[:e.start] + e.text + src[e.end:]
			last = e.start
		}
		if err := os.WriteFile(file, []byte(src), 0644); err != nil {
			panic(err)
		}
	}
	for dir, vs := range vars {
		var b strings.Builder
		b.WriteString("// Code generated by internal/offsets; DO NOT EDIT.\n\n")
		b.WriteString("package " + pkgNames[dir] + "\n\nimport (\n\t\"github.com/brutusmcforce/opennox/common/layout\"\n")
		var ips []string
		for ip := range imports[dir] {
			ips = append(ips, ip)
		}
		sort.Strings(ips)
		for _, ip := range ips {
			b.WriteString("\t\"" + ip + "\"\n")
		}
		b.WriteString(")\n\n// Offset tables for struct fields accessed by their 386 offset (see common/layout).\nvar (\n")
		var names []string
		for v := range vs {
			names = append(names, v)
		}
		sort.Strings(names)
		path := filepath.Join(dir, "layout_tables.go")
		old, _ := os.ReadFile(path)
		for _, v := range names {
			if strings.Contains(string(old), "\t"+v+" ") {
				continue
			}
			b.WriteString("\t" + v + " = layout.Of[" + vs[v] + "]()\n")
		}
		b.WriteString(")\n")
		if len(old) != 0 {
			// merge: append new variables to the existing file
			add := b.String()
			i := strings.Index(add, "var (\n")
			body := add[i+len("var (\n") : len(add)-2]
			o := string(old)
			o = strings.TrimSuffix(strings.TrimSpace(o), ")") + body + ")\n"
			for _, ip := range ips {
				if !strings.Contains(o, "\""+ip+"\"") {
					o = strings.Replace(o, "import (\n", "import (\n\t\""+ip+"\"\n", 1)
				}
			}
			os.WriteFile(path, []byte(o), 0644)
		} else {
			os.WriteFile(path, []byte(b.String()), 0644)
		}
	}
	fmt.Printf("# rewrote %d sites\n", n)
	fmt.Printf("# access width differs from the field on amd64 (review): %d\n", len(width))
	for _, s := range width {
		fmt.Printf("%-50s %-24s %6d %-30s access %s\n", s.Pos, s.Base, s.Off, s.Field, s.Access)
	}
}

// dataHints finds the typed per-class data of a function: calls of Object.UpdateDataX()-style accessors and
// conversions like (*server.XUpdateData)(u.UpdateData). A field with more than one candidate type gets nil.
func dataHints(info *types.Info, body *ast.BlockStmt) map[string]types.Type {
	out := make(map[string]types.Type)
	add := func(field string, t types.Type) {
		if old, ok := out[field]; ok && (old == nil || !types.Identical(old, t)) {
			out[field] = nil
			return
		}
		out[field] = t
	}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr); ok && len(call.Args) == 0 {
			for _, f := range []string{"UpdateData", "InitData", "CollideData"} {
				if strings.HasPrefix(sel.Sel.Name, f) && sel.Sel.Name != f {
					if pt, ok := info.TypeOf(call).(*types.Pointer); ok {
						if _, ok := pt.Elem().Underlying().(*types.Struct); ok {
							add(f, pt.Elem())
						}
					}
				}
			}
		}
		if tv, ok := info.Types[call.Fun]; ok && tv.IsType() && len(call.Args) == 1 {
			if pt, ok := tv.Type.Underlying().(*types.Pointer); ok {
				if _, ok := pt.Elem().Underlying().(*types.Struct); ok {
					if s, ok := ast.Unparen(call.Args[0]).(*ast.SelectorExpr); ok {
						switch s.Sel.Name {
						case "UpdateData", "InitData", "CollideData":
							add(s.Sel.Name, pt.Elem())
						}
					}
				}
			}
		}
		return true
	})
	return out
}
