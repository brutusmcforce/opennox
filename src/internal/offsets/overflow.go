package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Overflow mode (-overflow) lists reinterpretations of a value as a larger type, (*T)(unsafe.Pointer(&x)) or
// (*T)(unsafe.Add(unsafe.Pointer(&x), k)), that fit on 386 but not on amd64 (e.g. a union with pointers).

func overflows(pkgs []*packages.Package, wd string) {
	var out []string
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		info := p.TypesInfo
		for _, f := range p.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				c, ok := n.(*ast.CallExpr)
				if !ok || len(c.Args) != 1 {
					return true
				}
				tv, ok := info.Types[c.Fun]
				if !ok || !tv.IsType() {
					return true
				}
				pt, ok := tv.Type.Underlying().(*types.Pointer)
				if !ok {
					return true
				}
				arg := ast.Unparen(c.Args[0])
				var k int64
				if add, ok := arg.(*ast.CallExpr); ok && isUnsafeAdd(info, add) {
					v, ok := constInt(info, add.Args[1])
					if !ok {
						return true
					}
					k, arg = v, ast.Unparen(add.Args[0])
				}
				conv, ok := arg.(*ast.CallExpr)
				if !ok || len(conv.Args) != 1 || !isUnsafePointer(info.TypeOf(conv)) {
					return true
				}
				ue, ok := ast.Unparen(conv.Args[0]).(*ast.UnaryExpr)
				if !ok || ue.Op.String() != "&" {
					return true
				}
				ft := info.TypeOf(ue.X)
				if ft == nil {
					return true
				}
				t := pt.Elem()
				if hasTypeParam(t) || hasTypeParam(ft) {
					return true
				}
				if k+s386.Sizeof(t) <= s386.Sizeof(ft) && k+sAmd64.Sizeof(t) > sAmd64.Sizeof(ft) {
					pos := p.Fset.Position(c.Pos())
					rel, _ := filepath.Rel(wd, pos.Filename)
					out = append(out, fmt.Sprintf("%s:%d\t%s (%d bytes at +%d) over %s (%d bytes) on amd64", rel, pos.Line, t, sAmd64.Sizeof(t), k, ft, sAmd64.Sizeof(ft)))
				}
				return true
			})
		}
	}
	sort.Strings(out)
	for _, l := range out {
		fmt.Println(l)
	}
	fmt.Printf("# %d overflowing reinterpretations\n", len(out))
}

func hasTypeParam(t types.Type) (found bool) {
	defer func() {
		if recover() != nil {
			found = true
		}
	}()
	s386.Sizeof(t)
	return false
}

// allocSizes lists conversions (*T)(alloc(..., N)) of allocations with a constant size that matches T on 386
// but not on amd64 (the size was written as a number).
func allocSizes(pkgs []*packages.Package, wd string) {
	var out []string
	for _, p := range pkgs {
		if p.TypesInfo == nil {
			continue
		}
		info := p.TypesInfo
		for _, f := range p.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				c, ok := n.(*ast.CallExpr)
				if !ok || len(c.Args) != 1 {
					return true
				}
				tv, ok := info.Types[c.Fun]
				if !ok || !tv.IsType() {
					return true
				}
				pt, ok := tv.Type.Underlying().(*types.Pointer)
				if !ok {
					return true
				}
				call, ok := ast.Unparen(c.Args[0]).(*ast.CallExpr)
				if !ok || len(call.Args) == 0 {
					return true
				}
				name := types.ExprString(call.Fun)
				if !strings.Contains(strings.ToLower(name), "alloc") {
					return true
				}
				// total size: product of the constant args (count, size), or the single size
				for _, a := range call.Args {
					if strings.Contains(types.ExprString(a), "Sizeof") {
						return true // sized from the type
					}
				}
				cval := func(a ast.Expr) (int64, bool) {
					if v, ok := constInt(info, a); ok {
						return v, true
					}
					if conv, ok := ast.Unparen(a).(*ast.CallExpr); ok && len(conv.Args) == 1 {
						return constInt(info, conv.Args[0])
					}
					return 0, false
				}
				total := int64(1)
				for _, a := range call.Args {
					v, ok := cval(a)
					if !ok {
						// a variable count: check the element size (the last argument)
						v, ok = cval(call.Args[len(call.Args)-1])
						if !ok || len(call.Args) < 2 {
							return true
						}
						total = v
						break
					}
					total *= v
				}
				t := pt.Elem()
				if hasTypeParam(t) {
					return true
				}
				s32, s64 := s386.Sizeof(t), sAmd64.Sizeof(t)
				if s32 == s64 {
					return true
				}
				if total%s32 == 0 {
					pos := p.Fset.Position(c.Pos())
					rel, _ := filepath.Rel(wd, pos.Filename)
					out = append(out, fmt.Sprintf("%s:%d\t%s(%d bytes) of %s: %d bytes on 386, %d on amd64", rel, pos.Line, name, total, t, s32, s64))
				}
				return true
			})
		}
	}
	sort.Strings(out)
	for _, l := range out {
		fmt.Println(l)
	}
	fmt.Printf("# %d allocations sized for 386\n", len(out))
}
