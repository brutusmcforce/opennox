// Command ptrint lists conversions that keep pointers in 32-bit integers, which break on 64-bit platforms:
//
//   - narrow: a pointer converted through uintptr to a 32-bit (or smaller) integer, e.g. uint32(uintptr(p));
//   - widen: a 32-bit integer converted through uintptr to a pointer, e.g. unsafe.Pointer(uintptr(v));
//   - arith: address arithmetic on a 32-bit integer, e.g. unsafe.Pointer(uintptr(v) + 120).
//
// Usage (from src/):
//
//	go run ./internal/ptrint [-v] [packages]
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var (
	fVerbose = flag.Bool("v", false, "list every site")
	fFix     = flag.Bool("fix", false, "rewrite narrow and widen sites to memmap pointer words (PtrWord/WordPtr)")
)

type site struct {
	pos, kind, fn, expr string
	fix                 *edit
}

// edit replaces a byte range of a file.
type edit struct {
	file       string
	start, end int
	text       string
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
	var sites []site
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		sites = append(sites, scan(p, wd)...)
	}
	sort.Slice(sites, func(i, j int) bool { return sites[i].pos < sites[j].pos })
	if *fFix {
		applyFixes(sites)
		return
	}
	byFile := make(map[string]int)
	byKind := make(map[string]int)
	for _, s := range sites {
		byFile[strings.SplitN(s.pos, ":", 2)[0]]++
		byKind[s.kind]++
	}
	fmt.Printf("# sites: %d (narrow %d, widen %d, arith %d) in %d files\n", len(sites), byKind["narrow"], byKind["widen"], byKind["arith"], len(byFile))
	if *fVerbose {
		for _, s := range sites {
			fmt.Printf("%-50s %-6s %-35s %s\n", s.pos, s.kind, s.fn, s.expr)
		}
		return
	}
	type fc struct {
		f string
		n int
	}
	var fl []fc
	for f, n := range byFile {
		fl = append(fl, fc{f, n})
	}
	sort.Slice(fl, func(i, j int) bool { return fl[i].n > fl[j].n || fl[i].n == fl[j].n && fl[i].f < fl[j].f })
	for _, x := range fl {
		fmt.Printf("%5d %s\n", x.n, x.f)
	}
}

func isConv(info *types.Info, c *ast.CallExpr) (types.Type, bool) {
	if len(c.Args) != 1 {
		return nil, false
	}
	tv, ok := info.Types[c.Fun]
	if !ok || !tv.IsType() {
		return nil, false
	}
	return tv.Type, true
}

func basic(t types.Type) *types.Basic {
	if t == nil {
		return nil
	}
	b, _ := t.Underlying().(*types.Basic)
	return b
}

func isPtr(t types.Type) bool {
	if t == nil {
		return false
	}
	switch u := t.Underlying().(type) {
	case *types.Pointer:
		return true
	case *types.Basic:
		return u.Kind() == types.UnsafePointer
	}
	return false
}

// narrowInt reports if t is an integer type of at most 32 bits.
func narrowInt(t types.Type) bool {
	b := basic(t)
	if b == nil || b.Info()&types.IsInteger == 0 {
		return false
	}
	switch b.Kind() {
	case types.Int8, types.Int16, types.Int32, types.Uint8, types.Uint16, types.Uint32:
		return true
	}
	return false
}

func isInt(t types.Type) bool {
	b := basic(t)
	return b != nil && (b.Kind() == types.Int || b.Kind() == types.Uint)
}

func isUintptr(t types.Type) bool {
	b := basic(t)
	return b != nil && b.Kind() == types.Uintptr
}

func scan(p *packages.Package, wd string) []site {
	var out []site
	info := p.TypesInfo
	insp := inspector.New(p.Syntax)
	fnName := ""
	insp.WithStack([]ast.Node{(*ast.FuncDecl)(nil), (*ast.CallExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			if push {
				fnName = fd.Name.Name
			}
			return true
		}
		if !push {
			return true
		}
		c := n.(*ast.CallExpr)
		to, ok := isConv(info, c)
		if !ok {
			return true
		}
		arg := ast.Unparen(c.Args[0])
		inner, ok2 := arg.(*ast.CallExpr)
		var kind string
		switch {
		case (narrowInt(to) || isInt(to)) && ok2:
			// uint32(uintptr(p)) with p a pointer; int(uintptr(p)) is 32 bits on 386 too
			if t2, ok := isConv(info, inner); ok && isUintptr(t2) && isPtr(info.TypeOf(inner.Args[0])) {
				kind = "narrow"
			}
		case isPtr(to) && ok2:
			// unsafe.Pointer(uintptr(v)) with v a 32-bit integer
			if t2, ok := isConv(info, inner); ok && isUintptr(t2) && narrowInt(info.TypeOf(inner.Args[0])) {
				kind = "widen"
			}
		case isPtr(to):
			// unsafe.Pointer(uintptr(v) + n): address arithmetic on a 32-bit integer
			if be, ok := arg.(*ast.BinaryExpr); ok {
				ast.Inspect(be, func(n ast.Node) bool {
					if c2, ok := n.(*ast.CallExpr); ok {
						if t2, ok := isConv(info, c2); ok && isUintptr(t2) && narrowInt(info.TypeOf(c2.Args[0])) {
							kind = "arith"
						}
					}
					return kind == ""
				})
			}
		}
		if kind == "" {
			return true
		}
		pos := p.Fset.Position(c.Pos())
		rel, err := filepath.Rel(wd, pos.Filename)
		if err != nil {
			rel = pos.Filename
		}
		expr := types.ExprString(c)
		if len(expr) > 90 {
			expr = expr[:90] + "..."
		}
		s := site{pos: fmt.Sprintf("%s:%d", rel, pos.Line), kind: kind, fn: fnName, expr: expr}
		if *fFix {
			s.fix = makeFix(p, info, c, inner, to, kind)
		}
		out = append(out, s)
		return true
	})
	return out
}

func srcOf(p *packages.Package, n ast.Node) string {
	f := p.Fset.File(n.Pos())
	data, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	return string(data[f.Offset(n.Pos()):f.Offset(n.End())])
}

// makeFix builds the replacement for a narrow or widen site. Other sites are left for manual review.
func makeFix(p *packages.Package, info *types.Info, c, inner *ast.CallExpr, to types.Type, kind string) *edit {
	var text string
	switch kind {
	case "narrow":
		b := basic(to)
		if b == nil || b.Kind() == types.Int || b.Kind() == types.Uint {
			return nil // int: review by hand
		}
		ptr := inner.Args[0]
		arg := srcOf(p, ptr)
		if b2 := basic(info.TypeOf(ptr)); b2 == nil || b2.Kind() != types.UnsafePointer {
			arg = "unsafe.Pointer(" + arg + ")"
		}
		text = "memmap.PtrWord(" + arg + ")"
		if b.Kind() != types.Uint32 {
			text = srcOf(p, c.Fun) + "(" + text + ")"
		}
	case "widen":
		v := inner.Args[0]
		arg := srcOf(p, v)
		if b2 := basic(info.TypeOf(v)); b2 == nil || b2.Kind() != types.Uint32 {
			arg = "uint32(" + arg + ")"
		}
		text = "memmap.WordPtr(" + arg + ")"
		if b := basic(to); b == nil || b.Kind() != types.UnsafePointer {
			text = srcOf(p, c.Fun) + "(" + text + ")"
		}
	default:
		return nil
	}
	f := p.Fset.File(c.Pos())
	return &edit{file: f.Name(), start: f.Offset(c.Pos()), end: f.Offset(c.End()), text: text}
}

const memmapImport = `"github.com/brutusmcforce/opennox/common/memmap"`

func applyFixes(sites []site) {
	byFile := make(map[string][]*edit)
	for _, s := range sites {
		if s.fix != nil && !strings.Contains(s.fix.file, "/common/memmap/") && !strings.Contains(s.fix.file, "/internal/") {
			byFile[s.fix.file] = append(byFile[s.fix.file], s.fix)
		}
	}
	total, skipped := 0, 0
	for file, edits := range byFile {
		sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
		data, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		src := string(data)
		lastStart := len(src) + 1
		for _, e := range edits {
			if e.end > lastStart {
				skipped++ // nested in an edit applied already; the next run handles it
				continue
			}
			src = src[:e.start] + e.text + src[e.end:]
			lastStart = e.start
			total++
		}
		if !strings.Contains(src, memmapImport) {
			if i := strings.Index(src, "import (\n"); i >= 0 {
				src = src[:i+len("import (\n")] + "\t" + memmapImport + "\n" + src[i+len("import (\n"):]
			} else if i := strings.Index(src, "\nimport "); i >= 0 {
				src = src[:i+1] + "import " + memmapImport + "\n" + src[i+1:]
			}
		}
		if !strings.Contains(src, `"unsafe"`) && strings.Contains(src, "unsafe.") {
			if i := strings.Index(src, "import (\n"); i >= 0 {
				src = src[:i+len("import (\n")] + "\t\"unsafe\"\n" + src[i+len("import (\n"):]
			}
		}
		if err := os.WriteFile(file, []byte(src), 0644); err != nil {
			panic(err)
		}
	}
	fmt.Printf("# rewrote %d sites in %d files, %d nested sites left for another run\n", total, len(byFile), skipped)
}
