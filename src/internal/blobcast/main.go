// Command blobcast lists conversions of blob memory (the 32-bit memory images of the original binary, accessed
// through common/memmap) to Go types whose layout differs between 386 and amd64. Blobs keep the 386 layout, so
// such a type does not fit its slot on amd64: a struct with pointers, Go ints or list nodes overlaps the data
// that follows it.
//
// A blob pointer is a call to memmap.Ptr/PtrOff/PtrSizeOff/PtrUint8..., or a call to a function (in any of the
// loaded packages) whose result is one, e.g. func head() unsafe.Pointer { return memmap.PtrOff(0x5D4594, 12) }.
// The conversions are (*T)(p), (*[N]T)(p), unsafe.Slice((*T)(p), n) and layout.As[T](p).
//
// Usage (from src/):
//
//	go run ./internal/blobcast [-v] [packages]
//
// Output: one line per site, "file:line func type", sorted; -v adds the 386 and amd64 sizes.
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

	"golang.org/x/tools/go/packages"
)

var (
	sizes386   = types.SizesFor("gc", "386")
	sizesAmd64 = types.SizesFor("gc", "amd64")
)

const memmapPkg = "github.com/brutusmcforce/opennox/common/memmap"

func main() {
	verbose := flag.Bool("v", false, "print type sizes")
	flag.Parse()
	pats := flag.Args()
	if len(pats) == 0 {
		pats = []string{"./..."}
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Env:  append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0"),
	}
	pkgs, err := packages.Load(cfg, pats...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	wd, _ := os.Getwd()

	// functions whose result is a blob pointer
	blobFuncs := make(map[*types.Func]bool)
	for changed := true; changed; {
		changed = false
		for _, p := range pkgs {
			for _, f := range p.Syntax {
				for _, d := range f.Decls {
					fd, ok := d.(*ast.FuncDecl)
					if !ok || fd.Body == nil {
						continue
					}
					fn, _ := p.TypesInfo.Defs[fd.Name].(*types.Func)
					if fn == nil || blobFuncs[fn] || fn.Type().(*types.Signature).Results().Len() != 1 {
						continue
					}
					all, any := true, false
					ast.Inspect(fd.Body, func(n ast.Node) bool {
						if _, ok := n.(*ast.FuncLit); ok {
							return false
						}
						ret, ok := n.(*ast.ReturnStmt)
						if !ok || len(ret.Results) != 1 {
							return true
						}
						if isBlob(p.TypesInfo, blobFuncs, ret.Results[0]) {
							any = true
						} else {
							all = false
						}
						return true
					})
					if any && all {
						blobFuncs[fn] = true
						changed = true
					}
				}
			}
		}
	}

	var out []string
	seen := make(map[string]bool)
	for _, p := range pkgs {
		info := p.TypesInfo
		for _, f := range p.Syntax {
			var fnName string
			ast.Inspect(f, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					fnName = fd.Name.Name
				}
				call, ok := n.(*ast.CallExpr)
				if !ok || len(call.Args) == 0 {
					return true
				}
				t, src := convTarget(info, call)
				blob := t != nil && isBlob(info, blobFuncs, src)
				if !blob {
					// memmap.PtrT[T](base, off) and other memmap generics returning a typed blob pointer
					if t = memmapTypeArg(info, call); t != nil {
						blob = true
					}
				}
				if !blob || sameLayout(t) {
					return true
				}
				pos := p.Fset.Position(call.Pos())
				rel, err := filepath.Rel(wd, pos.Filename)
				if err != nil {
					rel = pos.Filename
				}
				line := fmt.Sprintf("%s:%d %s %s", rel, pos.Line, fnName, types.TypeString(t, shortQual))
				if *verbose {
					line += fmt.Sprintf(" (386 %d, amd64 %d bytes)", sizes386.Sizeof(t), sizesAmd64.Sizeof(t))
				}
				if !seen[line] {
					seen[line] = true
					out = append(out, line)
				}
				return true
			})
		}
	}
	sort.Strings(out)
	for _, l := range out {
		fmt.Println(l)
	}
}

func shortQual(p *types.Package) string { return p.Name() }

// convTarget returns the element type and the source expression of a conversion of a pointer to typed memory.
func convTarget(info *types.Info, call *ast.CallExpr) (types.Type, ast.Expr) {
	// unsafe.Slice((*T)(p), n)
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Slice" && len(call.Args) == 2 {
		if id, ok := sel.X.(*ast.Ident); ok {
			if pn, ok := info.Uses[id].(*types.PkgName); ok && pn.Imported().Path() == "unsafe" {
				if inner, ok := ast.Unparen(call.Args[0]).(*ast.CallExpr); ok {
					return convTarget(info, inner)
				}
				return nil, nil
			}
		}
	}
	// layout.As[T](p)
	if ix, ok := call.Fun.(*ast.IndexExpr); ok && len(call.Args) == 1 {
		if sel, ok := ix.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "As" {
			if tv, ok := info.Types[ix.Index]; ok && tv.IsType() {
				return tv.Type, call.Args[0]
			}
		}
		return nil, nil
	}
	// (*T)(p), (*[N]T)(p)
	if len(call.Args) != 1 {
		return nil, nil
	}
	tv, ok := info.Types[call.Fun]
	if !ok || !tv.IsType() {
		return nil, nil
	}
	pt, ok := tv.Type.Underlying().(*types.Pointer)
	if !ok {
		return nil, nil
	}
	elem := pt.Elem()
	if arr, ok := elem.Underlying().(*types.Array); ok {
		elem = arr.Elem()
	}
	return elem, call.Args[0]
}

// memmapTypeArg returns the type argument of a call to a generic memmap function (memmap.PtrT[T](...)).
func memmapTypeArg(info *types.Info, call *ast.CallExpr) types.Type {
	var id *ast.Ident
	switch f := ast.Unparen(call.Fun).(type) {
	case *ast.IndexExpr:
		if sel, ok := f.X.(*ast.SelectorExpr); ok {
			id = sel.Sel
		}
	case *ast.SelectorExpr:
		id = f.Sel // inferred type arguments
	}
	if id == nil {
		return nil
	}
	fn, ok := info.Uses[id].(*types.Func)
	if !ok || fn.Pkg() == nil || fn.Pkg().Path() != memmapPkg {
		return nil
	}
	inst, ok := info.Instances[id]
	if !ok || inst.TypeArgs.Len() != 1 {
		return nil
	}
	if _, ok := fn.Type().(*types.Signature).Results().At(0).Type().(*types.Pointer); !ok {
		return nil
	}
	t := inst.TypeArgs.At(0)
	if arr, ok := t.Underlying().(*types.Array); ok {
		t = arr.Elem()
	}
	return t
}

// isBlob reports if e is a blob pointer (see the package doc), through conversions.
func isBlob(info *types.Info, blobFuncs map[*types.Func]bool, e ast.Expr) bool {
	e = ast.Unparen(e)
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	if tv, ok := info.Types[call.Fun]; ok && tv.IsType() && len(call.Args) == 1 {
		return isBlob(info, blobFuncs, call.Args[0]) // unsafe.Pointer(p), uintptr(p)
	}
	var id *ast.Ident
	switch f := call.Fun.(type) {
	case *ast.Ident:
		id = f
	case *ast.SelectorExpr:
		id = f.Sel
	default:
		return false
	}
	fn, ok := info.Uses[id].(*types.Func)
	if !ok {
		return false
	}
	if fn.Pkg() != nil && fn.Pkg().Path() == memmapPkg && strings.HasPrefix(fn.Name(), "Ptr") && fn.Name() != "PtrWord" {
		return true
	}
	return blobFuncs[fn.Origin()]
}

// sameLayout reports if t has the same size and the same field offsets on 386 and amd64.
// Types that cannot be sized (type parameters, invalid types) are skipped.
func sameLayout(t types.Type) (same bool) {
	defer func() {
		if recover() != nil {
			same = true
		}
	}()
	return sameLayoutOf(t)
}

func sameLayoutOf(t types.Type) bool {
	if sizes386.Sizeof(t) != sizesAmd64.Sizeof(t) {
		return false
	}
	switch u := t.Underlying().(type) {
	case *types.Struct:
		fields := make([]*types.Var, u.NumFields())
		for i := range fields {
			fields[i] = u.Field(i)
		}
		o386, o64 := sizes386.Offsetsof(fields), sizesAmd64.Offsetsof(fields)
		for i := range fields {
			if o386[i] != o64[i] || !sameLayoutOf(fields[i].Type()) {
				return false
			}
		}
	case *types.Array:
		return sameLayoutOf(u.Elem())
	}
	return true
}
