// Command layoutdiff compares struct layouts on 386 and amd64. For every struct type that differs it prints
// the first 386 offset that moves, and with -v the offsets of all fields on both architectures.
//
// Usage (from src/):
//
//	go run ./internal/layoutdiff [-v] [-type Name] [packages]
package main

import (
	"flag"
	"fmt"
	"go/types"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	fVerbose = flag.Bool("v", false, "print all fields")
	fType    = flag.String("type", "", "only this type (pkg.Name or Name)")
)

var (
	s386   = types.SizesFor("gc", "386")
	sAmd64 = types.SizesFor("gc", "amd64")
)

type field struct {
	path   string
	off32  int64
	off64  int64
	size32 int64
	size64 int64
	typ    string
}

// flatten lists the leaf fields of a struct with their offsets on both architectures.
func flatten(t types.Type, prefix string, off32, off64 int64, out *[]field) {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		var fs []*types.Var
		for i := 0; i < u.NumFields(); i++ {
			fs = append(fs, u.Field(i))
		}
		o32 := s386.Offsetsof(fs)
		o64 := sAmd64.Offsetsof(fs)
		for i, f := range fs {
			name := f.Name()
			if prefix != "" {
				name = prefix + "." + name
			}
			flatten(f.Type(), name, off32+o32[i], off64+o64[i], out)
		}
	case *types.Array:
		e32, e64 := s386.Sizeof(u.Elem()), sAmd64.Sizeof(u.Elem())
		if e32 == e64 && !hasArchDep(u.Elem()) {
			*out = append(*out, field{prefix, off32, off64, s386.Sizeof(t), sAmd64.Sizeof(t), types.TypeString(t, nil)})
			return
		}
		for i := int64(0); i < u.Len(); i++ {
			flatten(u.Elem(), fmt.Sprintf("%s[%d]", prefix, i), off32+i*e32, off64+i*e64, out)
		}
	default:
		*out = append(*out, field{prefix, off32, off64, s386.Sizeof(t), sAmd64.Sizeof(t), types.TypeString(t, func(p *types.Package) string { return p.Name() })})
	}
}

func hasArchDep(t types.Type) bool { return s386.Sizeof(t) != sAmd64.Sizeof(t) }

func main() {
	flag.Parse()
	pats := flag.Args()
	if len(pats) == 0 {
		pats = []string{"./server", "./client/...", "./legacy/...", "./common/..."}
	}
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes, Env: append(os.Environ(), "GOARCH=386", "CGO_ENABLED=0")}
	pkgs, err := packages.Load(cfg, pats...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	type res struct {
		name           string
		size32, size64 int64
		first          int64
		fields         []field
	}
	var all []res
	for _, p := range pkgs {
		if p.Types == nil {
			continue
		}
		sc := p.Types.Scope()
		for _, n := range sc.Names() {
			tn, ok := sc.Lookup(n).(*types.TypeName)
			if !ok {
				continue
			}
			if _, ok := tn.Type().Underlying().(*types.Struct); !ok {
				continue
			}
			full := p.Types.Name() + "." + n
			if *fType != "" && *fType != full && *fType != n {
				continue
			}
			z32, z64 := s386.Sizeof(tn.Type()), sAmd64.Sizeof(tn.Type())
			if z32 == z64 && *fType == "" {
				continue
			}
			var fs []field
			flatten(tn.Type(), "", 0, 0, &fs)
			first := int64(-1)
			for _, f := range fs {
				if f.off32 != f.off64 || f.size32 != f.size64 {
					first = f.off32
					break
				}
			}
			all = append(all, res{full, z32, z64, first, fs})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].size32 > all[j].size32 })
	for _, r := range all {
		fmt.Printf("%-40s 386=%-6d amd64=%-6d first moved 386 offset=%d\n", r.name, r.size32, r.size64, r.first)
		if *fVerbose {
			for _, f := range r.fields {
				mark := ""
				if f.off32 != f.off64 || f.size32 != f.size64 {
					mark = " *"
				}
				fmt.Printf("    %6d %6d  %-40s %s%s\n", f.off32, f.off64, f.path, strings.TrimSpace(f.typ), mark)
			}
		}
	}
}
