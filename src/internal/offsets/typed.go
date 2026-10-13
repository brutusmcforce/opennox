package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// Typed mode (-typed FILE) rewrites offset accesses into typed field accesses:
//
//	*(*uint32)(unsafe.Add(ud, 188))  ->  layout.As[server.PlayerUpdateData](ud).Field47_0
//
// The base type of a site comes from, in order: the offset table of an earlier -fix rewrite (x.Off(N)),
// the static type of the base expression, the probe file (site<TAB>type lines, from layout probes), or
// another site in the same function with the same base expression (propagated).

type typedEnv struct {
	dry    bool // report only
	facts  paramFacts
	byName map[string]*types.Named // "server.PlayerUpdateData"
	probed map[string][]string     // site -> base types
	wd     string
	log    *bufio.Writer
	stats  map[string]int
}

func loadProbed(path string) map[string][]string {
	out := make(map[string][]string)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	for _, l := range strings.Split(string(data), "\n") {
		f := strings.Split(l, "\t")
		if len(f) < 2 || f[1] == "" {
			continue
		}
		out[f[0]] = strings.Split(f[1], ",")
	}
	return out
}

func namedTypes(pkgs []*packages.Package) map[string]*types.Named {
	out := make(map[string]*types.Named)
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		if p.Types == nil {
			return
		}
		sc := p.Types.Scope()
		for _, n := range sc.Names() {
			tn, ok := sc.Lookup(n).(*types.TypeName)
			if !ok {
				continue
			}
			nt, ok := tn.Type().(*types.Named)
			if !ok {
				continue
			}
			if _, ok := nt.Underlying().(*types.Struct); !ok {
				continue
			}
			out[p.Types.Name()+"."+n] = nt
		}
	})
	return out
}

// tedit replaces the source range of node with text. Inner edits are rendered into the texts of outer ones.
type tedit struct {
	start, end int
	build      func(render func(ast.Node) string) string
	text       string
}

type fileEdits struct {
	src     []byte
	fset    *token.FileSet
	edits   []*tedit
	imports map[string]string // path -> name
	pkg     *types.Package
}

func (fe *fileEdits) off(p token.Pos) int { return fe.fset.Position(p).Offset }

// render returns the source of n with all edits inside it applied.
func (fe *fileEdits) render(n ast.Node) string {
	s, e := fe.off(n.Pos()), fe.off(n.End())
	return fe.renderRange(s, e)
}

func (fe *fileEdits) renderRange(s, e int) string {
	var inner []*tedit
	for _, x := range fe.edits {
		if x.text != "" && x.start >= s && x.end <= e && !(x.start == s && x.end == e) {
			inner = append(inner, x)
		}
	}
	// keep outermost only
	sort.Slice(inner, func(i, j int) bool {
		return inner[i].start < inner[j].start || inner[i].start == inner[j].start && inner[i].end > inner[j].end
	})
	var b strings.Builder
	pos := s
	for _, x := range inner {
		if x.start < pos {
			continue // nested in an edit already rendered
		}
		b.Write(fe.src[pos:x.start])
		b.WriteString(x.text)
		pos = x.end
	}
	b.Write(fe.src[pos:e])
	return b.String()
}

// pathNode is one step of a field path.
type pathNode struct {
	path string
	typ  types.Type
	off  int64 // 386 offset from the base
}

// fieldPath lists the nested fields of t that contain the byte at off, outermost first (386 layout).
func fieldPath(t types.Type, off int64) ([]pathNode, string) {
	var out []pathNode
	path := ""
	var base int64
	for {
		switch u := t.Underlying().(type) {
		case *types.Struct:
			var fs []*types.Var
			for i := 0; i < u.NumFields(); i++ {
				fs = append(fs, u.Field(i))
			}
			offs := s386.Offsetsof(fs)
			found := false
			for i := len(fs) - 1; i >= 0; i-- {
				sz := s386.Sizeof(fs[i].Type())
				if off >= offs[i] && off < offs[i]+max(sz, 1) {
					if fs[i].Name() == "_" {
						return out, "blank field"
					}
					path += "." + fs[i].Name()
					base += offs[i]
					off -= offs[i]
					t = fs[i].Type()
					out = append(out, pathNode{path, t, base})
					found = true
					break
				}
			}
			if !found {
				return out, "padding"
			}
		case *types.Array:
			es := s386.Sizeof(u.Elem())
			if es == 0 {
				return out, "zero-size element"
			}
			i := off / es
			path += fmt.Sprintf("[%d]", i)
			base += i * es
			off -= i * es
			t = u.Elem()
			out = append(out, pathNode{path, t, base})
		default:
			return out, ""
		}
	}
}

func (env *typedEnv) note(kind, site, msg string) {
	env.stats[kind]++
	fmt.Fprintf(env.log, "%s\t%s\t%s\n", kind, site, msg)
}

type tsite struct {
	call   *ast.CallExpr
	stack  []ast.Node
	fn     ast.Node
	base   ast.Expr
	off    int64
	access types.Type // nil if unknown
	helper bool
	site   string
	static types.Type   // base type known statically
	owner  types.Object // the object whose UpdateData is the base, if any
	fnObj  *types.Func
	defs   map[types.Object]ast.Expr
	stride int64 // linear offset off + stride*index (an array element), if not 0
	index  ast.Expr
}

func (env *typedEnv) run(pkgs []*packages.Package) {
	env.facts = registrationFacts(pkgs)
	files := make(map[string]*fileEdits)
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		skip := false
		for _, s := range skipProbe {
			if strings.Contains(p.PkgPath+"/", s+"/") || strings.HasSuffix(p.PkgPath, s) {
				skip = true
			}
		}
		if skip {
			continue
		}
		env.pkg(p, files)
	}
	n := 0
	for path, fe := range files {
		if env.dry {
			n += len(fe.edits)
			continue
		}
		var top []*tedit
		sort.Slice(fe.edits, func(i, j int) bool { return fe.edits[i].end-fe.edits[i].start < fe.edits[j].end-fe.edits[j].start })
		for _, x := range fe.edits {
			x.text = x.build(fe.render)
		}
		for _, x := range fe.edits {
			outer := false
			for _, y := range fe.edits {
				if y != x && y.start <= x.start && y.end >= x.end && (y.end-y.start) > (x.end-x.start) {
					outer = true
					break
				}
			}
			if !outer {
				top = append(top, x)
			}
		}
		sort.Slice(top, func(i, j int) bool { return top[i].start > top[j].start })
		src := string(fe.src)
		for _, x := range top {
			src = src[:x.start] + x.text + src[x.end:]
		}
		for ip, name := range fe.imports {
			imp := `"` + ip + `"`
			if name != "" && name != filepath.Base(ip) {
				imp = name + " " + imp
			}
			if strings.Contains(src, `"`+ip+`"`) {
				continue
			}
			if i := strings.Index(src, "import (\n"); i >= 0 {
				src = src[:i+len("import (\n")] + "\t" + imp + "\n" + src[i+len("import (\n"):]
			} else if i := strings.Index(src, "\nimport "); i >= 0 {
				src = src[:i+1] + "import " + imp + "\n" + src[i+1:]
			} else {
				i := strings.Index(src, "\n") // after the package clause
				src = src[:i+1] + "\nimport " + imp + "\n" + src[i+1:]
			}
		}
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			panic(err)
		}
		n += len(fe.edits)
	}
	env.log.Flush()
	var ks []string
	for k := range env.stats {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		fmt.Printf("%6d %s\n", env.stats[k], k)
	}
	fmt.Printf("# %d edits in %d files\n", n, len(files))
}

// tableType returns the struct type of an offset table variable declared as layout.Of[T]().
func tableType(p *packages.Package, v *types.Var) types.Type {
	for _, f := range p.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, sp := range gd.Specs {
				vs := sp.(*ast.ValueSpec)
				for i, id := range vs.Names {
					if p.TypesInfo.Defs[id] != v || i >= len(vs.Values) {
						continue
					}
					call, ok := vs.Values[i].(*ast.CallExpr)
					if !ok {
						return nil
					}
					var fid *ast.Ident
					switch f := call.Fun.(type) {
					case *ast.IndexExpr:
						if s, ok := f.X.(*ast.SelectorExpr); ok {
							fid = s.Sel
						}
					}
					if fid == nil {
						return nil
					}
					if inst, ok := p.TypesInfo.Instances[fid]; ok && inst.TypeArgs.Len() == 1 {
						return inst.TypeArgs.At(0)
					}
				}
			}
		}
	}
	return nil
}

// offsetOf returns the constant 386 offset of an offset expression: a constant, or table.Off(constant).
func offsetOf(p *packages.Package, e ast.Expr, tables map[*types.Var]types.Type) (int64, types.Type, bool) {
	info := p.TypesInfo
	if v, ok := constInt(info, e); ok {
		return v, nil, true
	}
	e = ast.Unparen(e)
	if c, ok := e.(*ast.CallExpr); ok && len(c.Args) == 1 {
		if tv, ok := info.Types[c.Fun]; ok && tv.IsType() {
			return offsetOf(p, c.Args[0], tables) // int(x.Off(N))
		}
		if s, ok := c.Fun.(*ast.SelectorExpr); ok && s.Sel.Name == "Off" {
			if id, ok := s.X.(*ast.Ident); ok {
				if v, ok := info.Uses[id].(*types.Var); ok {
					t, seen := tables[v]
					if !seen {
						t = tableType(p, v)
						tables[v] = t
					}
					if t != nil {
						if n, ok := constInt(info, c.Args[0]); ok {
							return n, t, true
						}
					}
				}
			}
		}
	}
	return 0, nil, false
}

func (env *typedEnv) pkg(p *packages.Package, files map[string]*fileEdits) {
	info := p.TypesInfo
	helpers := findHelpers(p)
	tables := make(map[*types.Var]types.Type)
	var sites []*tsite
	insp := inspector.New(p.Syntax)
	var fn ast.Node
	var fnObj *types.Func
	var defs map[types.Object]ast.Expr
	insp.WithStack([]ast.Node{(*ast.FuncDecl)(nil), (*ast.CallExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		if fd, ok := n.(*ast.FuncDecl); ok {
			fn = fd
			fnObj, _ = info.Defs[fd.Name].(*types.Func)
			if fd.Body != nil {
				defs = assignments(info, fd.Body)
			}
			// helpers themselves are rewritten at their call sites
			if f, ok := info.Defs[fd.Name].(*types.Func); ok {
				if _, isHelper := helpers[f]; isHelper {
					return false
				}
			}
			return true
		}
		call := n.(*ast.CallExpr)
		pos := p.Fset.Position(call.Pos())
		if strings.HasSuffix(pos.Filename, "_test.go") || strings.HasSuffix(pos.Filename, "layout_tables.go") {
			return true
		}
		rel, _ := filepath.Rel(env.wd, pos.Filename)
		s := &tsite{call: call, stack: append([]ast.Node(nil), stack...), fn: fn, fnObj: fnObj, defs: defs, site: fmt.Sprintf("%s:%d:%d", rel, pos.Line, pos.Column)}
		var offE ast.Expr
		if isUnsafeAdd(info, call) {
			s.base, offE = call.Args[0], call.Args[1]
			if len(stack) >= 2 {
				if conv, ok := stack[len(stack)-2].(*ast.CallExpr); ok {
					if tv, ok := info.Types[conv.Fun]; ok && tv.IsType() {
						if pt, ok := tv.Type.Underlying().(*types.Pointer); ok {
							s.access = pt.Elem()
						}
					}
				}
			}
		} else if f := calleeFunc(info, call); f != nil {
			h, ok := helpers[f]
			if !ok || len(call.Args) <= max(h.baseParam, h.offParam) {
				return true
			}
			s.base, offE, s.access, s.helper = call.Args[h.baseParam], call.Args[h.offParam], h.access, true
		} else {
			return true
		}
		off, tt, ok := offsetOf(p, offE, tables)
		if !ok {
			c, k, idx, lok := linearOffset(info, offE)
			if !lok {
				if _, probed := env.probed[s.site]; probed {
					env.note("skip-nonconst", s.site, types.ExprString(offE))
				}
				return true
			}
			off, s.stride, s.index = c, k, idx
		}
		s.off = off
		s.static = tt
		s.owner = dataOwner(info, defs, s.base)
		if s.static == nil {
			if t, _ := baseType(info, defs, s.base, 0); t != nil {
				s.static = t
			}
		}
		sites = append(sites, s)
		return true
	})

	// probed types per (function, base expression), for propagation
	type fkey struct {
		fn   ast.Node
		base string
	}
	byFn := make(map[fkey]map[string]bool)
	for _, s := range sites {
		if ts := env.probed[s.site]; len(ts) != 0 && isStableExpr(s.base) {
			k := fkey{s.fn, types.ExprString(s.base)}
			if byFn[k] == nil {
				byFn[k] = make(map[string]bool)
			}
			for _, t := range ts {
				byFn[k][t] = true
			}
		}
	}
	for _, s := range sites {
		var t types.Type
		how := ""
		probed := env.probed[s.site]
		var ptype types.Type
		if len(probed) == 1 {
			if nt := env.byName[probed[0]]; nt != nil {
				ptype = nt
			}
		}
		switch {
		case s.static != nil:
			t, how = s.static, "static"
			if ptype != nil && !types.Identical(ptype, t) {
				env.note("conflict", s.site, fmt.Sprintf("static %s, probed %s", t, ptype))
				continue
			}
		case len(probed) > 1:
			env.note("multi-type", s.site, strings.Join(probed, ","))
			continue
		case len(probed) == 1 && s.owner != nil && classData[classGuard(info, s.stack, s.owner)] != "" && classData[classGuard(info, s.stack, s.owner)] != probed[0]:
			env.note("conflict", s.site, fmt.Sprintf("guard %s, probed %s", classData[classGuard(info, s.stack, s.owner)], probed[0]))
			continue
		case len(probed) == 1 && probed[0] == "raw" && s.access != nil && isPointerLike(s.access):
			env.note("raw-pointer", s.site, fmt.Sprintf("+%d %s", s.off, s.access))
			continue
		case len(probed) == 1:
			if ptype == nil {
				env.note("unknown-type", s.site, probed[0])
				continue
			}
			t, how = ptype, "probed"
		default:
			switch fact := env.facts.factFor(info, s.defs, s.fnObj, s.base); fact {
			case "":
			case factRaw:
				env.note("raw-registered", s.site, types.ExprString(s.base))
				continue
			case factConflict:
				env.note("registered-conflict", s.site, types.ExprString(s.base))
			default:
				if nt := env.byName[fact]; nt != nil {
					t, how = nt, "registered"
				}
			}
			if t == nil {
				if x := s.owner; x != nil {
					if name, ok := classData[classGuard(info, s.stack, x)]; ok {
						t, how = env.byName[name], "guarded"
					}
				}
			}
			if t == nil {
				env.stats["untyped"]++
				continue
			}
		}
		env.rewrite(p, s, t, how, files)
	}
}

// isStableExpr reports if e is an identifier or a selector chain on one (the same value within a function,
// unless reassigned).
func isStableExpr(e ast.Expr) bool {
	switch e := ast.Unparen(e).(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isStableExpr(e.X)
	}
	return false
}

func (env *typedEnv) rewrite(p *packages.Package, s *tsite, t types.Type, how string, files map[string]*fileEdits) {
	info := p.TypesInfo
	off := s.off
	nodes, why := fieldPath(t, off)
	if len(nodes) == 0 {
		env.note("no-field", s.site, fmt.Sprintf("%s+%d: %s", t, off, why))
		return
	}
	j0 := int64(0)
	if s.stride != 0 {
		// off + stride*index: an element of an array field with that element size
		var arr *pathNode
		for i := range nodes {
			if at, ok := nodes[i].typ.Underlying().(*types.Array); ok && s386.Sizeof(at.Elem()) == s.stride {
				arr = &nodes[i]
				break
			}
		}
		// the array node is the parent of the element node: find the element of the path
		if arr == nil {
			env.note("no-array", s.site, fmt.Sprintf("%s+%d+%d*i", t, off, s.stride))
			return
		}
		at := arr.typ.Underlying().(*types.Array)
		j0 = (off - arr.off) / s.stride
		r := (off - arr.off) % s.stride
		en, _ := fieldPath(at.Elem(), r)
		nodes = append([]pathNode{{path: arr.path + "[\x00]", typ: at.Elem(), off: 0}}, nil...)
		for _, n := range en {
			nodes = append(nodes, pathNode{path: arr.path + "[\x00]" + n.path, typ: n.typ, off: n.off})
		}
		off = r
	}
	var w int64 = 1
	if s.access != nil {
		w = s386.Sizeof(s.access)
	}
	// choose the deepest node that starts at off and has the access type; else the deepest that starts at off
	// and holds the whole access; else the leaf with a residual offset
	var pick *pathNode
	for i := len(nodes) - 1; i >= 0; i-- {
		nd := &nodes[i]
		if nd.off == off && s.access != nil && types.Identical(nd.typ, s.access) {
			pick = nd
			break
		}
	}
	if pick == nil {
		for i := len(nodes) - 1; i >= 0; i-- {
			nd := &nodes[i]
			if nd.off == off && off+w <= nd.off+s386.Sizeof(nd.typ) {
				if s.access == nil {
					// unknown access: the outermost field at off (e.g. a struct passed by pointer)
					for j := 0; j < i; j++ {
						if nodes[j].off == off {
							nd = &nodes[j]
							break
						}
					}
				}
				pick = nd
				break
			}
		}
	}
	resid := int64(0)
	if pick == nil {
		leaf := &nodes[len(nodes)-1]
		end := leaf.off + s386.Sizeof(leaf.typ)
		if len(nodes) >= 2 {
			// an access over several elements of an array of plain values is fine: they are contiguous
			if at, ok := nodes[len(nodes)-2].typ.Underlying().(*types.Array); ok && !isPointerLike(at.Elem()) {
				if _, basic := at.Elem().Underlying().(*types.Basic); basic {
					end = nodes[len(nodes)-2].off + s386.Sizeof(at)
				}
			}
		}
		if off+w > end {
			env.note("straddle", s.site, fmt.Sprintf("%s+%d width %d over %s", t, off, w, leaf.path))
			return
		}
		pick, resid = leaf, off-leaf.off
	}
	if s386.Sizeof(pick.typ) != sAmd64.Sizeof(pick.typ) && s.access != nil && !types.Identical(pick.typ, s.access) {
		env.note("width", s.site, fmt.Sprintf("%s%s is %s, accessed as %s", t, pick.path, pick.typ, s.access))
		// still rewritten: the access keeps the 386 width, the field position is right
	}
	if nt, ok := t.(*types.Named); ok && nt.Obj().Pkg() != p.Types && !nt.Obj().Exported() {
		env.note("unexported", s.site, t.String())
		return
	}
	// all fields on the path must be accessible
	if nt, ok := t.(*types.Named); ok && nt.Obj().Pkg() != p.Types {
		for _, seg := range strings.FieldsFunc(pick.path, func(r rune) bool { return r == '.' || r == '[' }) {
			if seg == "" || seg[0] >= '0' && seg[0] <= '9' || seg[0] == 0 {
				continue
			}
			if !token.IsExported(seg) {
				env.note("unexported", s.site, fmt.Sprintf("%s%s", t, pick.path))
				return
			}
		}
	}
	pos := p.Fset.Position(s.call.Pos())
	fe := files[pos.Filename]
	if fe == nil {
		src, err := os.ReadFile(pos.Filename)
		if err != nil {
			panic(err)
		}
		fe = &fileEdits{src: src, fset: p.Fset, imports: make(map[string]string), pkg: p.Types}
		files[pos.Filename] = fe
	}
	// the typed base expression
	baseExpr := func(render func(ast.Node) string) string {
		b := ast.Unparen(s.base)
		if bt := info.TypeOf(b); bt != nil {
			if pt, ok := bt.Underlying().(*types.Pointer); ok && types.Identical(pt.Elem(), t) {
				return render(b)
			}
		}
		if c, ok := b.(*ast.CallExpr); ok && len(c.Args) == 1 {
			if tv, ok := info.Types[c.Fun]; ok && tv.IsType() {
				if pt, ok := info.TypeOf(c.Args[0]).Underlying().(*types.Pointer); ok && types.Identical(pt.Elem(), t) {
					return render(c.Args[0]) // unsafe.Pointer(x)
				}
			}
		}
		if c, ok := b.(*ast.CallExpr); ok && len(c.Args) == 0 {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok {
				if pt, ok := info.TypeOf(sel.X).Underlying().(*types.Pointer); ok && types.Identical(pt.Elem(), t) {
					return render(sel.X) // x.C()
				}
			}
		}
		fe.imports["github.com/brutusmcforce/opennox/common/layout"] = ""
		return "layout.As[" + env.typeName(t, p.Types, fe) + "](" + render(b) + ")"
	}
	field := func(render func(ast.Node) string) string {
		path := pick.path
		if s.stride != 0 {
			idx := render(s.index)
			if j0 != 0 {
				idx = fmt.Sprintf("%d+%s", j0, idx)
			}
			path = strings.Replace(path, "\x00", idx, 1)
		}
		return baseExpr(render) + path
	}
	// the node to replace and its text
	var node ast.Node = s.call
	var build func(render func(ast.Node) string) string
	exact := resid == 0 && s.access != nil && types.Identical(pick.typ, s.access)
	parentConv := !s.helper && len(s.stack) >= 2 && isConvTo(info, s.stack[len(s.stack)-2], s.access)
	var star bool
	switch {
	case exact:
		if !s.helper && parentConv {
			node = s.stack[len(s.stack)-2]
		}
		if len(s.stack) >= 2 {
			k := len(s.stack) - 2
			if !s.helper && parentConv {
				k--
			}
			if k >= 0 {
				if pp, ok := parentOf(s.stack[0], node).(*ast.StarExpr); ok {
					node, star = pp, true
				}
			}
		}
		if star {
			build = field
		} else if s.helper || parentConv {
			build = func(r func(ast.Node) string) string { return "&" + field(r) }
		} else {
			build = func(r func(ast.Node) string) string { return "unsafe.Pointer(&" + field(r) + ")" }
		}
	default:
		ptr := func(r func(ast.Node) string) string {
			x := "unsafe.Pointer(&" + field(r) + ")"
			if resid != 0 {
				x = fmt.Sprintf("unsafe.Add(%s, %d)", x, resid)
			}
			return x
		}
		if s.helper {
			at := env.typeName(s.access, p.Types, fe)
			build = func(r func(ast.Node) string) string { return "(*" + at + ")(" + ptr(r) + ")" }
		} else {
			build = ptr
		}
	}
	fe.edits = append(fe.edits, &tedit{start: fe.off(node.Pos()), end: fe.off(node.End()), build: build})
	env.stats["rewritten-"+how]++
	fmt.Fprintf(env.log, "ok-%s\t%s\t%s+%d -> %s\n", how, s.site, t, off, pick.path)
}

func isConvTo(info *types.Info, n ast.Node, t types.Type) bool {
	c, ok := n.(*ast.CallExpr)
	if !ok || t == nil || len(c.Args) != 1 {
		return false
	}
	tv, ok := info.Types[c.Fun]
	if !ok || !tv.IsType() {
		return false
	}
	pt, ok := tv.Type.Underlying().(*types.Pointer)
	return ok && types.Identical(pt.Elem(), t)
}

// parentOf finds the parent of n under root, skipping parentheses.
func parentOf(root, n ast.Node) ast.Node {
	var stack []ast.Node
	var found ast.Node
	ast.Inspect(root, func(x ast.Node) bool {
		if found != nil {
			return false
		}
		if x == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		if x == n {
			for i := len(stack) - 1; i >= 0; i-- {
				if _, ok := stack[i].(*ast.ParenExpr); !ok {
					found = stack[i]
					break
				}
			}
			return false
		}
		stack = append(stack, x)
		return true
	})
	return found
}

// typeName writes a type as seen from package pkg, adding imports as needed.
func (env *typedEnv) typeName(t types.Type, pkg *types.Package, fe *fileEdits) string {
	return types.TypeString(t, func(p *types.Package) string {
		if p == pkg {
			return ""
		}
		fe.imports[p.Path()] = p.Name()
		return p.Name()
	})
}

// Class guards: an object's per-class data has a known type where the class is checked first, e.g.
//
//	if u.ObjClass&4 != 0 { ... u.UpdateData ... }        // PlayerUpdateData
//	if u.ObjClass&2 == 0 { return }; ... u.UpdateData ... // MonsterUpdateData

var classData = map[int64]string{
	4: "server.PlayerUpdateData",
	2: "server.MonsterUpdateData",
}

// dataOwner returns the object whose UpdateData the base expression is (u for u.UpdateData, or a local
// assigned once from it).
func dataOwner(info *types.Info, defs map[types.Object]ast.Expr, e ast.Expr) types.Object {
	e = ast.Unparen(e)
	if id, ok := e.(*ast.Ident); ok {
		if d, ok := defs[info.Uses[id]]; ok {
			e = ast.Unparen(d)
		}
	}
	sel, ok := e.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "UpdateData" {
		return nil
	}
	id, ok := ast.Unparen(sel.X).(*ast.Ident)
	if !ok {
		return nil
	}
	return info.Uses[id]
}

// classTest matches "x.ObjClass&K != 0" (want true) or "x.ObjClass&K == 0" (want false), and the
// x.Class().Has(object.ClassX) forms, for a single-bit K. It returns K.
func classTest(info *types.Info, e ast.Expr, x types.Object, want bool) int64 {
	e = ast.Unparen(e)
	if u, ok := e.(*ast.UnaryExpr); ok && u.Op == token.NOT {
		return classTest(info, u.X, x, !want)
	}
	if be, ok := e.(*ast.BinaryExpr); ok && (be.Op == token.NEQ || be.Op == token.EQL) {
		if v, ok := constInt(info, be.Y); !ok || v != 0 {
			return 0
		}
		if (be.Op == token.NEQ) != want {
			return 0
		}
		and, ok := ast.Unparen(be.X).(*ast.BinaryExpr)
		if !ok || and.Op != token.AND {
			return 0
		}
		sel, ok := ast.Unparen(and.X).(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "ObjClass" {
			return 0
		}
		if id, ok := ast.Unparen(sel.X).(*ast.Ident); !ok || info.Uses[id] != x {
			return 0
		}
		k, ok := constInt(info, and.Y)
		if !ok || k&(k-1) != 0 {
			return 0
		}
		return k
	}
	if c, ok := e.(*ast.CallExpr); ok && want && len(c.Args) == 1 {
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Has" && sel.Sel.Name != "HasAny" {
			return 0
		}
		inner, ok := ast.Unparen(sel.X).(*ast.CallExpr)
		if !ok {
			return 0
		}
		isel, ok := inner.Fun.(*ast.SelectorExpr)
		if !ok || isel.Sel.Name != "Class" {
			return 0
		}
		if id, ok := ast.Unparen(isel.X).(*ast.Ident); !ok || info.Uses[id] != x {
			return 0
		}
		switch types.ExprString(c.Args[0]) {
		case "object.ClassPlayer":
			return 4
		case "object.ClassMonster":
			return 2
		}
	}
	return 0
}

// condTrue returns a class bit that holds when cond is true (a conjunct).
func condTrue(info *types.Info, cond ast.Expr, x types.Object) int64 {
	cond = ast.Unparen(cond)
	if be, ok := cond.(*ast.BinaryExpr); ok && be.Op == token.LAND {
		if k := condTrue(info, be.X, x); k != 0 {
			return k
		}
		return condTrue(info, be.Y, x)
	}
	return classTest(info, cond, x, true)
}

// condFalse returns a class bit that holds when cond is false (a disjunct "class bit not set").
func condFalse(info *types.Info, cond ast.Expr, x types.Object) int64 {
	cond = ast.Unparen(cond)
	if be, ok := cond.(*ast.BinaryExpr); ok && be.Op == token.LOR {
		if k := condFalse(info, be.X, x); k != 0 {
			return k
		}
		return condFalse(info, be.Y, x)
	}
	return classTest(info, cond, x, false)
}

func endsInReturn(b *ast.BlockStmt) bool {
	if len(b.List) == 0 {
		return false
	}
	switch b.List[len(b.List)-1].(type) {
	case *ast.ReturnStmt, *ast.BranchStmt:
		return true
	}
	return false
}

// classGuard finds the class bit of object x that guards the innermost node of stack.
func classGuard(info *types.Info, stack []ast.Node, x types.Object) int64 {
	for i := len(stack) - 2; i >= 0; i-- {
		switch n := stack[i].(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			return 0
		case *ast.IfStmt:
			if stack[i+1] == n.Body {
				if k := condTrue(info, n.Cond, x); k != 0 {
					return k
				}
			}
		case *ast.BlockStmt:
			for _, st := range n.List {
				if st == stack[i+1] {
					break
				}
				if ifs, ok := st.(*ast.IfStmt); ok && ifs.Else == nil && endsInReturn(ifs.Body) {
					if k := condFalse(info, ifs.Cond, x); k != 0 {
						return k
					}
				}
			}
		}
	}
	return 0
}

// linearOffset matches c + k*i, k*i + c, k*i, i*k and c + i (with conversions), for constants c and k.
func linearOffset(info *types.Info, e ast.Expr) (c, k int64, idx ast.Expr, ok bool) {
	e = unconv(info, e)
	if be, ok := e.(*ast.BinaryExpr); ok && be.Op == token.ADD {
		if v, ok := constInt(info, be.X); ok {
			k, idx, ok := scaled(info, be.Y)
			return v, k, idx, ok
		}
		if v, ok := constInt(info, be.Y); ok {
			k, idx, ok := scaled(info, be.X)
			return v, k, idx, ok
		}
		return 0, 0, nil, false
	}
	k, idx, ok = scaled(info, e)
	if !ok || k == 1 {
		return 0, 0, nil, false // a bare variable is not an index
	}
	return 0, k, idx, true
}

// scaled matches k*i, i*k and i (k = 1).
func scaled(info *types.Info, e ast.Expr) (int64, ast.Expr, bool) {
	e = unconv(info, e)
	if be, ok := e.(*ast.BinaryExpr); ok && be.Op == token.MUL {
		if v, ok := constInt(info, be.X); ok {
			return v, be.Y, true
		}
		if v, ok := constInt(info, be.Y); ok {
			return v, be.X, true
		}
		return 0, nil, false
	}
	if _, ok := constInt(info, e); ok {
		return 0, nil, false
	}
	return 1, e, true
}

// unconv strips parentheses and integer conversions.
func unconv(info *types.Info, e ast.Expr) ast.Expr {
	for {
		e = ast.Unparen(e)
		c, ok := e.(*ast.CallExpr)
		if !ok || len(c.Args) != 1 {
			return e
		}
		tv, ok := info.Types[c.Fun]
		if !ok || !tv.IsType() {
			return e
		}
		if b, ok := tv.Type.Underlying().(*types.Basic); !ok || b.Info()&types.IsInteger == 0 {
			return e
		}
		e = c.Args[0]
	}
}

// isPointerLike reports if values of t hold pointers (or uintptr) that do not fit in 32 bits on 64-bit platforms.
func isPointerLike(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Pointer, *types.Signature, *types.Map, *types.Chan, *types.Slice, *types.Interface:
		return true
	case *types.Basic:
		return u.Kind() == types.UnsafePointer || u.Kind() == types.Uintptr || u.Kind() == types.Int || u.Kind() == types.Uint || u.Kind() == types.String
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if isPointerLike(u.Field(i).Type()) {
				return true
			}
		}
	case *types.Array:
		return isPointerLike(u.Elem())
	}
	return false
}

// Registration facts: object functions registered with server.RegisterObjectX(name, ..., data) get the object
// whose data block (XData) has the registered type. The facts flow to the functions they pass the object to.

// regField maps registration kinds to the data field of the object.
var regField = map[string]string{"Init": "InitData", "Update": "UpdateData", "Collide": "CollideData", "Use": "UseData", "Death": "DeathData"}

const factRaw = "raw" // the data block is raw (keeps the 32-bit layout)
const factConflict = "conflict"

// paramFacts: function -> parameter index -> data field -> type name (or factRaw / factConflict)
type paramFacts map[*types.Func]map[int]map[string]string

func (pf paramFacts) add(fn *types.Func, i int, field, typ string) bool {
	if pf[fn] == nil {
		pf[fn] = make(map[int]map[string]string)
	}
	if pf[fn][i] == nil {
		pf[fn][i] = make(map[string]string)
	}
	old, ok := pf[fn][i][field]
	switch {
	case !ok:
		pf[fn][i][field] = typ
		return true
	case old != typ && old != factConflict:
		pf[fn][i][field] = factConflict
		return true
	}
	return false
}

// dataDefType returns the type name of a data definition argument: DataOf[T]() -> "pkg.T", RawData(n) -> raw.
func dataDefType(info *types.Info, e ast.Expr) (string, bool) {
	e = ast.Unparen(e)
	switch e := e.(type) {
	case *ast.CallExpr:
		switch f := e.Fun.(type) {
		case *ast.IndexExpr: // DataOf[T]()
			if sel, ok := f.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "DataOf" || isIdent(f.X, "DataOf") {
				if t := info.TypeOf(f.Index); t != nil {
					return types.TypeString(t, func(p *types.Package) string { return p.Name() }), true
				}
			}
		case *ast.SelectorExpr:
			if f.Sel.Name == "RawData" {
				return factRaw, true
			}
		case *ast.Ident:
			if f.Name == "RawData" {
				return factRaw, true
			}
		}
	case *ast.SelectorExpr:
		if e.Sel.Name == "NoData" {
			return factRaw, true
		}
	case *ast.Ident:
		if e.Name == "NoData" {
			return factRaw, true
		}
	}
	return "", false
}

func isIdent(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == name
}

// registrationFacts finds registered object functions and propagates the facts through calls.
func registrationFacts(pkgs []*packages.Package) paramFacts {
	pf := make(paramFacts)
	funcs := make(map[*types.Func]*ast.FuncDecl)
	infos := make(map[*types.Func]*types.Info)
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		for _, f := range p.Syntax {
			for _, d := range f.Decls {
				if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
					if fn, ok := p.TypesInfo.Defs[fd.Name].(*types.Func); ok {
						funcs[fn] = fd
						infos[fn] = p.TypesInfo
					}
				}
			}
		}
	}
	// seeds
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		info := p.TypesInfo
		for _, f := range p.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok || len(call.Args) < 3 {
					return true
				}
				var name string
				switch fn := call.Fun.(type) {
				case *ast.SelectorExpr:
					name = fn.Sel.Name
				case *ast.Ident:
					name = fn.Name
				}
				var field string
				for kind, fld := range regField {
					if strings.HasPrefix(name, "RegisterObject"+kind) && !strings.HasSuffix(name, "Parse") {
						field = fld
					}
				}
				if field == "" {
					return true
				}
				typ, ok := dataDefType(info, call.Args[len(call.Args)-1])
				if !ok {
					return true
				}
				for _, a := range call.Args {
					lit, ok := a.(*ast.FuncLit)
					if !ok || len(lit.Type.Params.List) == 0 || len(lit.Body.List) == 0 {
						continue
					}
					var params []types.Object
					for _, fl := range lit.Type.Params.List {
						for _, nm := range fl.Names {
							params = append(params, info.Defs[nm])
						}
					}
					if len(params) == 0 {
						continue
					}
					// calls in the closure body that pass its first parameter
					ast.Inspect(lit.Body, func(n ast.Node) bool {
						c, ok := n.(*ast.CallExpr)
						if !ok {
							return true
						}
						callee := calleeFunc(info, c)
						if callee == nil {
							return true
						}
						for j, arg := range c.Args {
							if id, ok := ast.Unparen(arg).(*ast.Ident); ok && info.Uses[id] == params[0] {
								pf.add(callee, j, field, typ)
							}
						}
						return true
					})
				}
				return true
			})
		}
	}
	// propagate: a function passes its parameter on to other functions
	for changed := true; changed; {
		changed = false
		for fn, byParam := range pf {
			fd := funcs[fn]
			if fd == nil {
				continue
			}
			info := infos[fn]
			sig := fn.Type().(*types.Signature)
			for i, facts := range byParam {
				if i >= sig.Params().Len() {
					continue
				}
				param := sig.Params().At(i)
				ast.Inspect(fd.Body, func(n ast.Node) bool {
					c, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					callee := calleeFunc(info, c)
					if callee == nil || callee == fn {
						return true
					}
					for j, arg := range c.Args {
						if id, ok := ast.Unparen(arg).(*ast.Ident); ok && info.Uses[id] == param {
							for field, typ := range facts {
								if pf.add(callee, j, field, typ) {
									changed = true
								}
							}
						}
					}
					return true
				})
			}
		}
	}
	return pf
}

// factFor returns the registered type of the data block that base refers to (param.XData, or a local assigned
// once from it) in function fn.
func (pf paramFacts) factFor(info *types.Info, defs map[types.Object]ast.Expr, fn *types.Func, base ast.Expr) string {
	if fn == nil || pf[fn] == nil {
		return ""
	}
	e := ast.Unparen(base)
	if id, ok := e.(*ast.Ident); ok {
		if d, ok := defs[info.Uses[id]]; ok {
			e = ast.Unparen(d)
		}
	}
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	id, ok := ast.Unparen(sel.X).(*ast.Ident)
	if !ok {
		return ""
	}
	obj := info.Uses[id]
	sig := fn.Type().(*types.Signature)
	for i := 0; i < sig.Params().Len(); i++ {
		if sig.Params().At(i) == obj {
			return pf[fn][i][sel.Sel.Name]
		}
	}
	return ""
}
