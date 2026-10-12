package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/brutusmcforce/opennox/common/memmap"
)

var (
	fGap   = flag.Int("gap", 128, "max distance in bytes between pointer slots of one table")
	fUnref = flag.String("unref", "", "write the pointer slots of unreferenced tables to this file (one \"blob offset\" per line)")
)

func varName(addr uintptr) string {
	if v := memmap.VariableByAddr(addr); v != nil {
		if v.Addr == addr {
			return v.Name
		}
		return fmt.Sprintf("%s+%d", v.Name, addr-v.Addr)
	}
	blob, off := memmap.RelativeAddr(addr)
	return fmt.Sprintf("byte_%X_%d", blob, off)
}

func offExpr(a *Access) string {
	parts := []string{fmt.Sprint(a.Off)}
	for _, t := range a.Terms {
		parts = append(parts, fmt.Sprintf("%d*(%s)", t.Coef, t.Expr))
	}
	return strings.Join(parts, " + ")
}

func kindSummary(accs []*Access) string {
	kinds := make(map[string]int)
	for _, a := range accs {
		k := a.Kind
		if a.PtrConv {
			k += "(ptrconv)"
		}
		kinds[k]++
	}
	var ks []string
	for k, n := range kinds {
		ks = append(ks, fmt.Sprintf("%s:%d", k, n))
	}
	sort.Strings(ks)
	return strings.Join(ks, " ")
}

// table is a run of pointer slots that are close to each other, usually one table of the original binary.
type table struct {
	blob     uintptr
	lo, hi   uintptr // absolute addresses of the first and the last pointer slot
	slots    []uintptr
	accs     []*Access // runtime memmap accesses into the table
	lits     []Literal // integer constants pointing into the table
	relocTgt int       // pointers to this table stored by the static initialization
}

func (t *table) stride() uintptr {
	gaps := make(map[uintptr]int)
	for i := 1; i < len(t.slots); i++ {
		gaps[t.slots[i]-t.slots[i-1]]++
	}
	var best uintptr
	for g, n := range gaps {
		if n > gaps[best] || n == gaps[best] && g < best {
			best = g
		}
	}
	return best
}

func (t *table) status() string {
	var fixed, indexed, dynamic bool
	for _, a := range t.accs {
		switch {
		case len(a.Terms) != 0:
			indexed = true
		default:
			fixed = true
		}
	}
	dynamic = len(t.lits) != 0
	var s []string
	if fixed {
		s = append(s, "fixed")
	}
	if indexed {
		s = append(s, "indexed")
	}
	if dynamic {
		s = append(s, "literals")
	}
	if t.relocTgt != 0 {
		s = append(s, "reloc")
	}
	if len(s) == 0 {
		return "UNREFERENCED"
	}
	return strings.Join(s, "+")
}

func report(w io.Writer, all []Access, lits []Literal) {
	var runtime []*Access
	cnt := make(map[string]int)
	for i := range all {
		a := &all[i]
		cnt[a.Kind]++
		if !a.Init {
			runtime = append(runtime, a)
		}
	}
	var kinds []string
	for k := range cnt {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	fmt.Fprintf(w, "# memmap calls: %d (runtime %d, static initialization %d)\n", len(all), len(runtime), len(all)-len(runtime))
	for _, k := range kinds {
		fmt.Fprintf(w, "#   %-8s %d\n", k, cnt[k])
	}

	// All slots known to hold a pointer.
	ptrSet := make(map[uintptr]bool)
	relocTargets := make(map[uintptr]int)
	for i := range all {
		a := &all[i]
		if len(a.Terms) == 0 && a.IsPtr() {
			ptrSet[a.Addr()] = true
		} else if len(a.Terms) != 0 && (a.IsPtr() || a.HasPtr) {
			ptrSet[a.Addr()] = true // table base; the extent is unknown
		}
		if a.Init && a.Kind == "addr" {
			relocTargets[a.Addr()]++
		}
	}
	addrs := make([]uintptr, 0, len(ptrSet))
	for p := range ptrSet {
		addrs = append(addrs, p)
	}
	sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })

	var tables []*table
	for _, p := range addrs {
		blob, _ := memmap.RelativeAddr(p)
		if n := len(tables); n != 0 && tables[n-1].blob == blob && p-tables[n-1].hi <= uintptr(*fGap) {
			t := tables[n-1]
			t.hi = p
			t.slots = append(t.slots, p)
			continue
		}
		tables = append(tables, &table{blob: blob, lo: p, hi: p, slots: []uintptr{p}})
	}
	// Assign accesses, literals and relocation targets. A table of structs may start before its first pointer field.
	const before = 64
	for _, t := range tables {
		lo, hi := t.lo-before, t.hi+4
		for _, a := range runtime {
			if a.Addr() >= t.lo && a.Addr() < hi || len(a.Terms) != 0 && a.Addr() >= lo && a.Addr() < hi {
				t.accs = append(t.accs, a)
			}
		}
		offLo, offHi := int64(lo-t.blob), int64(hi-t.blob)
		for _, l := range lits {
			if l.Val >= 256 && l.Val >= offLo && l.Val < offHi {
				t.lits = append(t.lits, l)
			}
		}
		for addr, n := range relocTargets {
			if addr >= lo && addr < hi {
				t.relocTgt += n
			}
		}
	}

	statusCnt := make(map[string]int)
	for _, t := range tables {
		statusCnt[t.status()]++
	}
	var sts []string
	for s, n := range statusCnt {
		sts = append(sts, fmt.Sprintf("%s:%d", s, n))
	}
	sort.Strings(sts)
	fmt.Fprintf(w, "\n## pointer tables: %d pointer slots in %d tables (slots closer than %d bytes)\n", len(addrs), len(tables), *fGap)
	fmt.Fprintf(w, "# reached by: %s\n", strings.Join(sts, " "))
	fmt.Fprintln(w, "# fixed = memmap call at a constant address, indexed = memmap call with an index,")
	fmt.Fprintln(w, "# literals = integer constants in range (dynamic offsets), reloc = pointed to by other static data")
	fmt.Fprintln(w, "# blob+off          slots stride  span  reached by            memmap calls")
	for _, t := range tables {
		fmt.Fprintf(w, "%X+%-9d %5d %6d %5d  %-21s %s\n", t.blob, t.lo-t.blob, len(t.slots), t.stride(), t.hi+4-t.lo, t.status(), kindSummary(t.accs))
		if *fVerbose {
			for _, a := range t.accs {
				fmt.Fprintf(w, "        %-5s %-10s %s  %s\n", a.Kind, a.Func, offExpr(a), a.Pos)
			}
			for _, l := range t.lits {
				fmt.Fprintf(w, "        lit   %-10d %s\n", l.Val, l.Pos)
			}
		}
	}

	if *fUnref != "" {
		var buf strings.Builder
		for _, t := range tables {
			if t.status() == "UNREFERENCED" {
				for _, p := range t.slots {
					fmt.Fprintf(&buf, "0x%X %d\n", t.blob, p-t.blob)
				}
			}
		}
		if err := os.WriteFile(*fUnref, []byte(buf.String()), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// Pointer slots also accessed as plain integers: these need a look before they become typed variables.
	fmt.Fprintln(w, "\n## pointer slots with non-pointer accesses")
	for _, p := range addrs {
		var ptr, other []*Access
		for _, a := range runtime {
			if len(a.Terms) != 0 || a.Kind == "addr" || a.Size == 0 {
				continue
			}
			if a.Addr() < p+4 && a.Addr()+uintptr(a.Size) > p {
				if a.IsPtr() && a.Addr() == p {
					ptr = append(ptr, a)
				} else {
					other = append(other, a)
				}
			}
		}
		if len(other) == 0 {
			continue
		}
		fmt.Fprintf(w, "0x%07X %-28s %s\n", p, varName(p), kindSummary(append(ptr, other...)))
		if *fVerbose {
			for _, a := range append(ptr, other...) {
				fmt.Fprintf(w, "        %-5s %-10s %+d  %s\n", a.Kind, a.Func, int64(a.Addr()-p), a.Pos)
			}
		}
	}

	// Calls with a small constant part and a variable offset: the address comes from elsewhere.
	fmt.Fprintln(w, "\n## dynamic-offset calls (address computed elsewhere)")
	for _, a := range runtime {
		if len(a.Terms) != 0 && a.Off < 64 {
			fmt.Fprintf(w, "%X  %-5s %-10s %-40s %s\n", a.Blob, a.Kind, a.Func, offExpr(a), a.Pos)
		}
	}

	fmt.Fprintln(w, "\n## overlays with pointers or arch-dependent layout")
	for _, a := range runtime {
		if a.Kind == "overlay" && (a.HasPtr || a.ArchDep) {
			fmt.Fprintf(w, "0x%07X %-28s %-30s size=%d ptr=%v archdep=%v  %s\n", a.Addr(), varName(a.Addr()), a.Type, a.Size, a.HasPtr, a.ArchDep, a.Pos)
		}
	}

	esc := 0
	for _, a := range runtime {
		if a.Escapes {
			esc++
		}
	}
	fmt.Fprintf(w, "\n## accessor results kept as raw pointers (not dereferenced at once): %d (see -json, field escapes)\n", esc)
}
