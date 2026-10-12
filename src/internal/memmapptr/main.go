// Command memmapptr lists every access to the original Nox memory blobs (common/memmap) and finds the blob
// slots that hold pointers.
//
// The blobs keep the 32-bit layout of the original binary, so a pointer stored in a blob only fits on 386.
// This tool is the first step of the 64-bit port: it shows which slots must become typed Go variables.
//
// Usage (from src/):
//
//	go run ./internal/memmapptr [-v] [-gap 128] [-json out.json] [-unref slots.txt] [packages]
//
// -v lists every access site and integer constant per table, -json writes all accesses for further tooling.
// Tables nothing refers to are most likely dead data that the Go code no longer reads.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"

	"github.com/brutusmcforce/opennox/common/memmap"
	_ "github.com/brutusmcforce/opennox/common/memmap/nox"
)

const memmapPkg = "github.com/brutusmcforce/opennox/common/memmap"

var (
	fJSON    = flag.String("json", "", "write all accesses as JSON to this file")
	fVerbose = flag.Bool("v", false, "list every access site of the pointer slots")
)

func main() {
	flag.Parse()
	pats := flag.Args()
	if len(pats) == 0 {
		pats = []string{"./..."}
	}
	if err := run(pats); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Term is a non-constant part of an offset expression: Coef * Expr.
type Term struct {
	Expr string `json:"expr"`
	Coef int64  `json:"coef"`
}

// Access is one call to a memmap accessor.
type Access struct {
	Pos   string  `json:"pos"`
	Func  string  `json:"func"`
	Blob  uintptr `json:"blob"`
	Off   int64   `json:"off"`             // constant part of the offset
	Terms []Term  `json:"terms,omitempty"` // non-constant parts (table index)
	Kind  string  `json:"kind"`            // ptr, pword, u8, i8, u16, i16, u32, i32, u64, i64, f32, f64, str, addr, overlay
	Size  int64   `json:"size"`            // access size on 386, 0 if unknown
	Type  string  `json:"type,omitempty"`  // overlay type
	// HasPtr is set for overlays that contain pointers, ArchDep when the overlay size differs on amd64.
	HasPtr  bool `json:"has_ptr,omitempty"`
	ArchDep bool `json:"arch_dep,omitempty"`
	Write   bool `json:"write,omitempty"`
	Escapes bool `json:"escapes,omitempty"` // the returned pointer is kept (not dereferenced at once)
	// PtrConv is set for an integer access whose value is converted from or to a pointer.
	PtrConv bool `json:"ptr_conv,omitempty"`
	Init    bool `json:"init,omitempty"` // static data initialization (blobdata)
}

func (a *Access) Addr() uintptr { return a.Blob + uintptr(a.Off) }

// IsPtr reports if the access treats the slot as a pointer.
func (a *Access) IsPtr() bool { return a.Kind == "ptr" || a.Kind == "pword" || a.PtrConv }

func run(pats []string) error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps,
		Env:   append(os.Environ(), "GOARCH=386", "GOOS=linux", "CGO_ENABLED=0"),
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, pats...)
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	var (
		accs []Access
		lits []Literal
	)
	seen := make(map[string]bool)
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		if p.TypesInfo == nil || seen[p.PkgPath] || p.PkgPath == memmapPkg {
			return
		}
		seen[p.PkgPath] = true
		if len(p.Errors) != 0 {
			// Packages that need cgo (OpenGL, SDL) don't type-check without it; accesses there are still found.
			fmt.Fprintf(os.Stderr, "warning: %s: %d type errors, first: %v\n", p.PkgPath, len(p.Errors), p.Errors[0])
		}
		pa, pl := scanPackage(p, wd)
		accs = append(accs, pa...)
		lits = append(lits, pl...)
	})
	sort.SliceStable(accs, func(i, j int) bool { return accs[i].Addr() < accs[j].Addr() })
	if *fJSON != "" {
		data, err := json.MarshalIndent(accs, "", "\t")
		if err != nil {
			return err
		}
		if err := os.WriteFile(*fJSON, data, 0644); err != nil {
			return err
		}
	}
	report(os.Stdout, accs, lits)
	return nil
}

var (
	sizes386   = types.SizesFor("gc", "386")
	sizesAmd64 = types.SizesFor("gc", "amd64")
)

var scalarFuncs = map[string]string{
	"PtrUint8": "u8", "PtrInt8": "i8", "PtrUint16": "u16", "PtrInt16": "i16",
	"PtrUint32": "u32", "PtrInt32": "i32", "PtrUint64": "u64", "PtrInt64": "i64",
	"PtrFloat32": "f32", "PtrFloat64": "f64",
	"Uint8": "u8", "Int8": "i8", "Uint16": "u16", "Int16": "i16",
	"Uint32": "u32", "Int32": "i32", "Uint64": "u64", "Int64": "i64",
	"Float32": "f32", "Float64": "f64",
}

var kindSize = map[string]int64{
	"u8": 1, "i8": 1, "u16": 2, "i16": 2, "u32": 4, "i32": 4, "f32": 4, "ptr": 4, "pword": 4,
	"u64": 8, "i64": 8, "f64": 8,
}

// Literal is an integer constant in the code. Offsets into blob tables that are not passed to memmap directly
// (a table walked with a loop variable, a base offset returned by a helper) are found through them.
type Literal struct {
	Pos string
	Val int64
}

func scanPackage(p *packages.Package, wd string) ([]Access, []Literal) {
	var (
		out  []Access
		lits []Literal
	)
	position := func(n ast.Node) string {
		pos := p.Fset.Position(n.Pos())
		rel, err := filepath.Rel(wd, pos.Filename)
		if err != nil {
			rel = pos.Filename
		}
		return fmt.Sprintf("%s:%d", rel, pos.Line)
	}
	insp := inspector.New(p.Syntax)
	info := p.TypesInfo
	memmapArgs := 0 // depth of memmap calls on the stack
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil), (*ast.BasicLit)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch n := n.(type) {
		case *ast.BasicLit:
			if push && memmapArgs == 0 && n.Kind == token.INT && p.Name != "blobdata" {
				if v, err := strconv.ParseInt(n.Value, 0, 64); err == nil {
					lits = append(lits, Literal{Pos: position(n), Val: v})
				}
			}
			return true
		case *ast.CallExpr:
			fn := calledMemmapFunc(info, n)
			if fn == "" {
				return true
			}
			if !push {
				memmapArgs--
				return true
			}
			memmapArgs++
			a, ok := newAccess(info, fn, n, stack)
			if !ok {
				return true
			}
			a.Pos = position(n)
			a.Init = p.Name == "blobdata"
			out = append(out, a)
		}
		return true
	})
	return out, lits
}

// calledMemmapFunc returns the name of the memmap accessor called by the expression.
func calledMemmapFunc(info *types.Info, call *ast.CallExpr) string {
	fun := ast.Unparen(call.Fun)
	if ix, ok := fun.(*ast.IndexExpr); ok { // PtrT[T]
		fun = ix.X
	}
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	obj, ok := info.Uses[sel.Sel].(*types.Func)
	if !ok || obj.Pkg() == nil || obj.Pkg().Path() != memmapPkg {
		return ""
	}
	return obj.Name()
}

func newAccess(info *types.Info, fn string, call *ast.CallExpr, stack []ast.Node) (Access, bool) {
	a := Access{Func: fn}
	switch fn {
	case "Ptr", "PtrSize":
		addr, ok := constUint(info, call.Args[0])
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: non-constant address in memmap.%s\n", fn)
			return a, false
		}
		blob, off := memmap.RelativeAddr(addr)
		a.Blob, a.Off = blob, int64(off)
	default:
		if len(call.Args) < 2 {
			return a, false // registration and query functions
		}
		blob, ok := constUint(info, call.Args[0])
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: non-constant blob in memmap.%s\n", fn)
			return a, false
		}
		a.Blob = blob
		a.Off, a.Terms = linear(info, call.Args[1])
	}
	parents := stack[:len(stack)-1]
	switch fn {
	case "PtrPtr":
		a.Kind = "ptr"
	case "LoadPtr", "StorePtr":
		a.Kind = "pword" // 32-bit pointer word, see memmap.PtrWord
		a.Write = fn == "StorePtr"
	case "String":
		a.Kind = "str"
	case "PtrT":
		inst := info.Instances[call.Fun.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel]
		setOverlay(&a, inst.TypeArgs.At(0))
	case "Ptr", "PtrSize", "PtrOff", "PtrSizeOff", "Slice":
		a.Kind = "addr"
		if len(parents) > 0 {
			if conv, ok := parents[len(parents)-1].(*ast.CallExpr); ok && isConversion(info, conv) {
				if pt, ok := info.TypeOf(conv).Underlying().(*types.Pointer); ok {
					setOverlay(&a, pt.Elem())
					parents = parents[:len(parents)-1]
				}
			}
		}
	default:
		k, ok := scalarFuncs[fn]
		if !ok {
			return a, false
		}
		a.Kind = k
	}
	if a.Size == 0 {
		a.Size = kindSize[a.Kind]
	}
	// Ptr* accessors return a pointer: find out if it is dereferenced and written.
	if strings.HasPrefix(fn, "Ptr") || fn == "Slice" {
		if len(parents) > 0 {
			if star, ok := parents[len(parents)-1].(*ast.StarExpr); ok {
				a.Write = isAssigned(star, parents[:len(parents)-1])
			} else if _, ok := parents[len(parents)-1].(*ast.IndexExpr); !ok {
				a.Escapes = true
			}
		}
	}
	if isInt(a.Kind) {
		a.PtrConv = usesPtrConv(info, call, parents)
	}
	return a, true
}

func isInt(kind string) bool {
	switch kind {
	case "u32", "i32", "u64", "i64":
		return true
	}
	return false
}

func setOverlay(a *Access, t types.Type) {
	a.Type = types.TypeString(t, func(p *types.Package) string { return p.Name() })
	switch u := t.Underlying().(type) {
	case *types.Pointer:
		a.Kind = "ptr"
		return
	case *types.Basic:
		switch u.Kind() {
		case types.UnsafePointer, types.Uintptr:
			a.Kind = "ptr"
			return
		case types.Uint8:
			a.Kind = "u8"
		case types.Int8:
			a.Kind = "i8"
		case types.Uint16:
			a.Kind = "u16"
		case types.Int16:
			a.Kind = "i16"
		case types.Uint32, types.Uint:
			a.Kind = "u32"
		case types.Int32, types.Int:
			a.Kind = "i32"
		case types.Uint64:
			a.Kind = "u64"
		case types.Int64:
			a.Kind = "i64"
		case types.Float32:
			a.Kind = "f32"
		case types.Float64:
			a.Kind = "f64"
		default:
			a.Kind = "overlay"
		}
		a.Size = sizes386.Sizeof(t)
		return
	}
	a.Kind = "overlay"
	a.Size = sizes386.Sizeof(t)
	a.HasPtr = hasPointers(t, make(map[types.Type]bool))
	a.ArchDep = sizes386.Sizeof(t) != sizesAmd64.Sizeof(t)
}

func hasPointers(t types.Type, seen map[types.Type]bool) bool {
	if seen[t] {
		return false
	}
	seen[t] = true
	switch u := t.Underlying().(type) {
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Signature, *types.Interface:
		return true
	case *types.Basic:
		return u.Kind() == types.UnsafePointer || u.Kind() == types.Uintptr || u.Kind() == types.String
	case *types.Array:
		return hasPointers(u.Elem(), seen)
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if hasPointers(u.Field(i).Type(), seen) {
				return true
			}
		}
	}
	return false
}

// isAssigned reports if the dereference is written to.
func isAssigned(star ast.Expr, parents []ast.Node) bool {
	if len(parents) == 0 {
		return false
	}
	switch p := parents[len(parents)-1].(type) {
	case *ast.AssignStmt:
		for _, l := range p.Lhs {
			if l == star {
				return true
			}
		}
	case *ast.IncDecStmt:
		return true
	case *ast.UnaryExpr:
		return p.Op == token.AND // address taken: may be written
	}
	return false
}

func isConversion(info *types.Info, call *ast.CallExpr) bool {
	tv, ok := info.Types[call.Fun]
	return ok && tv.IsType() && len(call.Args) == 1
}

func isPtrType(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Pointer:
		return true
	case *types.Basic:
		return u.Kind() == types.UnsafePointer
	}
	return false
}

// usesPtrConv reports if an integer value read from or written to the blob is converted from or to a pointer.
func usesPtrConv(info *types.Info, call *ast.CallExpr, parents []ast.Node) bool {
	// Value flows up through expressions: look for a conversion to a pointer type.
	var top ast.Node = call
	for i := len(parents) - 1; i >= 0; i-- {
		p := parents[i]
		if c, ok := p.(*ast.CallExpr); ok && isConversion(info, c) && isPtrType(info.TypeOf(c)) {
			return true
		}
		if _, ok := p.(ast.Expr); !ok {
			top = p
			break
		}
	}
	// A written value: look for a conversion from a pointer on the right side.
	if as, ok := top.(*ast.AssignStmt); ok && inLHS(as, call) {
		for _, r := range as.Rhs {
			found := false
			ast.Inspect(r, func(n ast.Node) bool {
				if c, ok := n.(*ast.CallExpr); ok && isConversion(info, c) {
					if t := info.TypeOf(c.Args[0]); t != nil && isPtrType(t) {
						found = true
					}
				}
				return !found
			})
			if found {
				return true
			}
		}
	}
	return false
}

// inLHS reports if the call is part of the assignment's left side.
func inLHS(as *ast.AssignStmt, call *ast.CallExpr) bool {
	for _, l := range as.Lhs {
		if l.Pos() <= call.Pos() && call.End() <= l.End() {
			return true
		}
	}
	return false
}

func constUint(info *types.Info, e ast.Expr) (uintptr, bool) {
	tv, ok := info.Types[e]
	if !ok || tv.Value == nil {
		return 0, false
	}
	v, ok := constant.Uint64Val(constant.ToInt(tv.Value))
	return uintptr(v), ok
}

// linear splits an offset expression into a constant and Coef*Expr terms.
func linear(info *types.Info, e ast.Expr) (int64, []Term) {
	if tv, ok := info.Types[e]; ok && tv.Value != nil {
		v, _ := constant.Int64Val(constant.ToInt(tv.Value))
		return v, nil
	}
	switch e := e.(type) {
	case *ast.ParenExpr:
		return linear(info, e.X)
	case *ast.CallExpr:
		if isConversion(info, e) {
			if b, ok := info.TypeOf(e).Underlying().(*types.Basic); ok && b.Info()&types.IsInteger != 0 {
				return linear(info, e.Args[0])
			}
		}
	case *ast.UnaryExpr:
		if e.Op == token.ADD {
			return linear(info, e.X)
		} else if e.Op == token.SUB {
			c, t := linear(info, e.X)
			return -c, scale(t, -1)
		}
	case *ast.BinaryExpr:
		switch e.Op {
		case token.ADD, token.SUB:
			c1, t1 := linear(info, e.X)
			c2, t2 := linear(info, e.Y)
			if e.Op == token.SUB {
				c2, t2 = -c2, scale(t2, -1)
			}
			return c1 + c2, append(t1, t2...)
		case token.MUL:
			c1, t1 := linear(info, e.X)
			c2, t2 := linear(info, e.Y)
			if len(t1) == 0 {
				return c1 * c2, scale(t2, c1)
			} else if len(t2) == 0 {
				return c1 * c2, scale(t1, c2)
			}
		}
	}
	return 0, []Term{{Expr: types.ExprString(e), Coef: 1}}
}

func scale(ts []Term, k int64) []Term {
	for i := range ts {
		ts[i].Coef *= k
	}
	return ts
}
