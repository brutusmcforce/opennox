package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

// Probe mode (-probe) instruments a copy of the tree: every unsafe.Add becomes layout.Probe, and the base of
// every offset helper call is wrapped in layout.ProbeBase. With NOX_LAYOUTCHECK and NOX_LAYOUTPROBE set, the
// instrumented game then records the block type and offset that each site touches (see common/layout).

// skipProbe lists packages that must not be instrumented (the probe machinery and the allocators under it).
var skipProbe = []string{"/common/layout", "/common/memmap", "/legacy/common/alloc", "/internal/"}

type insert struct {
	at   int
	text string
	del  int // bytes to delete at at
}

func (c *pkgCtx) probes(wd string, out map[string][]insert) int {
	p := c.p
	for _, s := range skipProbe {
		if strings.Contains(p.PkgPath+"/", s+"/") || strings.HasSuffix(p.PkgPath, s) {
			return 0
		}
	}
	info := p.TypesInfo
	n := 0
	insp := inspector.New(p.Syntax)
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(nd ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		call := nd.(*ast.CallExpr)
		pos := p.Fset.Position(call.Pos())
		if strings.HasSuffix(pos.Filename, "_test.go") {
			return true
		}
		rel, _ := filepath.Rel(wd, pos.Filename)
		site := fmt.Sprintf("%s:%d:%d", rel, pos.Line, pos.Column)
		f := p.Fset.File(call.Pos())
		off := func(e ast.Node) int { return f.Offset(e.Pos()) }
		end := func(e ast.Node) int { return f.Offset(e.End()) }
		width := func(access types.Type) int64 {
			if access == nil {
				return 0
			}
			return s386.Sizeof(access)
		}
		if isUnsafeAdd(info, call) {
			var access types.Type
			if len(stack) >= 2 {
				if conv, ok := stack[len(stack)-2].(*ast.CallExpr); ok {
					if tv, ok := info.Types[conv.Fun]; ok && tv.IsType() {
						if pt, ok := tv.Type.Underlying().(*types.Pointer); ok {
							access = pt.Elem()
						}
					}
				}
			}
			offE := call.Args[1]
			fun := ast.Unparen(call.Fun)
			out[pos.Filename] = append(out[pos.Filename],
				insert{at: off(fun), del: end(fun) - off(fun), text: "layoutprobe.Probe"},
				insert{at: off(offE), text: "int("},
				insert{at: end(offE), text: fmt.Sprintf("), %d, %q", width(access), site)},
			)
			n++
			return true
		}
		if isUnsafeSlice(info, call) {
			// unsafe.Slice((*T)(p), n): an access of n*sizeof(T) bytes at p
			conv, ok := ast.Unparen(call.Args[0]).(*ast.CallExpr)
			if !ok || len(conv.Args) != 1 {
				return true
			}
			tv, ok := info.Types[conv.Fun]
			if !ok || !tv.IsType() {
				return true
			}
			pt, ok := tv.Type.Underlying().(*types.Pointer)
			if !ok {
				return true
			}
			arg := conv.Args[0]
			if at := info.TypeOf(arg); at == nil || !isUnsafePointer(at) {
				return true // a typed pointer: not a reinterpretation
			}
			if inner, ok := ast.Unparen(arg).(*ast.CallExpr); ok {
				if sel, ok := ast.Unparen(inner.Fun).(*ast.SelectorExpr); ok {
					if id, ok := sel.X.(*ast.Ident); ok && id.Name == "memmap" {
						return true // static blob data
					}
				}
			}
			var w int64
			if n, ok := constInt(info, call.Args[1]); ok {
				w = n * s386.Sizeof(pt.Elem())
			}
			out[pos.Filename] = append(out[pos.Filename],
				insert{at: off(arg), text: "layoutprobe.ProbeBase("},
				insert{at: end(arg), text: fmt.Sprintf(", 0, %d, %q)", w, site)},
			)
			n++
			return true
		}
		if id := calleeFunc(info, call); id != nil {
			h, ok := c.helpers[id]
			if !ok || len(call.Args) <= max(h.baseParam, h.offParam) {
				return true
			}
			offE := call.Args[h.offParam]
			if _, ok := constInt(info, offE); !ok {
				if tv, ok := info.Types[offE]; !ok || tv.Value == nil {
					// only constant offsets: the offset expression is evaluated twice
					if _, isCall := ast.Unparen(offE).(*ast.CallExpr); isCall || !isSimple(offE) {
						return true
					}
				}
			}
			src, err := os.ReadFile(pos.Filename)
			if err != nil {
				panic(err)
			}
			offText := string(src[off(offE):end(offE)])
			base := call.Args[h.baseParam]
			out[pos.Filename] = append(out[pos.Filename],
				insert{at: off(base), text: "layoutprobe.ProbeBase("},
				insert{at: end(base), text: fmt.Sprintf(", int(%s), %d, %q)", offText, width(h.access), site)},
			)
			n++
		}
		return true
	})
	return n
}

func isUnsafeSlice(info *types.Info, call *ast.CallExpr) bool {
	sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Slice" || len(call.Args) != 2 {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	pn, ok := info.Uses[id].(*types.PkgName)
	return ok && pn.Imported().Path() == "unsafe"
}

func isUnsafePointer(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Kind() == types.UnsafePointer
}

// isSimple reports if an expression has no side effects (identifiers, constants, selectors, arithmetic).
func isSimple(e ast.Expr) bool {
	ok := true
	ast.Inspect(e, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.CallExpr, *ast.FuncLit, *ast.UnaryExpr:
			ok = false
		}
		return ok
	})
	return ok
}

const layoutImport = `layoutprobe "github.com/brutusmcforce/opennox/common/layout"`

func applyProbes(out map[string][]insert) {
	for file, ins := range out {
		// stable order: later positions first; at equal positions, keep the insertion order reversed so that
		// text inserted earlier ends up first
		for i := range ins {
			ins[i].at = ins[i].at*1024 + i // encode the order
		}
		sort.Slice(ins, func(i, j int) bool { return ins[i].at > ins[j].at })
		data, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		src := string(data)
		for _, x := range ins {
			at := x.at / 1024
			src = src[:at] + x.text + src[at+x.del:]
		}
		if !strings.Contains(src, layoutImport) {
			if i := strings.Index(src, "import (\n"); i >= 0 {
				src = src[:i+len("import (\n")] + "\t" + layoutImport + "\n" + src[i+len("import (\n"):]
			} else if i := strings.Index(src, "\nimport "); i >= 0 {
				src = src[:i+1] + "import " + layoutImport + "\n" + src[i+1:]
			}
		}
		if strings.Contains(src, `"unsafe"`) {
			src += "\nvar _ unsafe.Pointer\n"
		}
		if err := os.WriteFile(file, []byte(src), 0644); err != nil {
			panic(err)
		}
	}
}
