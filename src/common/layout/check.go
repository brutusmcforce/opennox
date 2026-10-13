package layout

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Layout checks.
//
// With NOX_LAYOUTCHECK set, allocators record the type of every block they hand out (Track), and typed
// conversions of legacy pointers (As) verify that the block really holds that type at that position. This
// finds wrong type assumptions on 386, where they would otherwise go unnoticed because both layouts agree.
//
// Probe is used by instrumented builds only: it records which block types each offset access touches.

// Checking reports if layout checks are enabled.
var Checking = os.Getenv("NOX_LAYOUTCHECK") != ""

func init() {
	if Checking {
		fmt.Fprintln(os.Stderr, "layoutcheck: enabled")
	}
}

type block struct {
	start, size uintptr
	typ         reflect.Type // nil: untyped memory
	name        string       // for untyped memory
}

const pageShift = 12

var reg struct {
	sync.Mutex
	pages  map[uintptr][]*block // every page a block covers -> blocks
	starts map[uintptr]*block
}

// Track records the type of the block at p. A nil type marks untyped memory, with an optional name.
// Tracking a block at the start of a tracked block replaces it.
func Track(p unsafe.Pointer, size uintptr, t reflect.Type, name string) {
	if !Checking || p == nil || size == 0 {
		return
	}
	reg.Lock()
	defer reg.Unlock()
	if reg.pages == nil {
		reg.pages = make(map[uintptr][]*block)
		reg.starts = make(map[uintptr]*block)
	}
	untrack(uintptr(p))
	b := &block{start: uintptr(p), size: size, typ: t, name: name}
	reg.starts[b.start] = b
	for pg := b.start >> pageShift; pg <= (b.start+size-1)>>pageShift; pg++ {
		reg.pages[pg] = append(reg.pages[pg], b)
	}
}

// TrackAs records that the block at p holds a T.
func TrackAs[T any](p unsafe.Pointer) {
	if Checking {
		Track(p, reflect.TypeFor[T]().Size(), reflect.TypeFor[T](), "")
	}
}

// Untrack forgets the block starting at p.
func Untrack(p unsafe.Pointer) {
	if !Checking || p == nil {
		return
	}
	reg.Lock()
	defer reg.Unlock()
	untrack(uintptr(p))
}

func untrack(start uintptr) {
	b := reg.starts[start]
	if b == nil {
		return
	}
	delete(reg.starts, start)
	for pg := b.start >> pageShift; pg <= (b.start+b.size-1)>>pageShift; pg++ {
		l := reg.pages[pg]
		for i, x := range l {
			if x == b {
				l = append(l[:i], l[i+1:]...)
				break
			}
		}
		if len(l) == 0 {
			delete(reg.pages, pg)
		} else {
			reg.pages[pg] = l
		}
	}
}

// find returns the smallest tracked block containing p.
func find(p uintptr) *block {
	reg.Lock()
	defer reg.Unlock()
	var best *block
	for _, b := range reg.pages[p>>pageShift] {
		if p >= b.start && p < b.start+b.size && (best == nil || b.size < best.size) {
			best = b
		}
	}
	return best
}

// As converts a legacy pointer to *T. With checks enabled, it reports if the memory is not a T.
func As[T any](p unsafe.Pointer) *T {
	if Checking && p != nil {
		checkAs(p, reflect.TypeFor[T]())
	}
	return (*T)(p)
}

var unions sync.Map // union type -> map[reflect.Type]bool of its variants

// Union declares that values of type U hold one of the variants (checked casts to a variant are allowed).
func Union[U any](variants ...any) {
	m := make(map[reflect.Type]bool)
	for _, v := range variants {
		m[reflect.TypeOf(v)] = true
	}
	unions.Store(reflect.TypeFor[U](), m)
}

// typeAt reports if a value of type want lies at offset off of a value of type t.
func typeAt(t reflect.Type, off uintptr, want reflect.Type) bool {
	for {
		if off == 0 && t == want {
			return true
		}
		if off == 0 {
			if m, ok := unions.Load(t); ok && m.(map[reflect.Type]bool)[want] {
				return true
			}
		}
		switch t.Kind() {
		case reflect.Struct:
			found := false
			for i := t.NumField() - 1; i >= 0; i-- {
				f := t.Field(i)
				if off >= f.Offset && off < f.Offset+max(f.Type.Size(), 1) {
					off -= f.Offset
					t = f.Type
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case reflect.Array:
			es := t.Elem().Size()
			if es == 0 {
				return false
			}
			off %= es
			t = t.Elem()
		default:
			return false
		}
	}
}

func checkAs(p unsafe.Pointer, want reflect.Type) {
	b := find(uintptr(p))
	switch {
	case b == nil:
		report("as", "%s: untracked memory", want)
	case b.typ == nil:
		report("as", "%s: untyped block %q (%d bytes, at +%d)", want, b.name, b.size, uintptr(p)-b.start)
	case !typeAt(b.typ, uintptr(p)-b.start, want):
		report("as", "%s: WRONG TYPE: block is %s, at +%d", want, b.typ, uintptr(p)-b.start)
	}
}

var reportHook func(line string) // for tests

var reported struct {
	sync.Mutex
	seen map[string]bool
}

// report logs a check failure once per call site and message.
func report(kind, format string, args ...any) {
	site := callerSite()
	msg := fmt.Sprintf(format, args...)
	key := site + " " + msg
	reported.Lock()
	defer reported.Unlock()
	if reported.seen == nil {
		reported.seen = make(map[string]bool)
	}
	if reported.seen[key] {
		return
	}
	reported.seen[key] = true
	line := fmt.Sprintf("layoutcheck: %s %s: %s", kind, site, msg)
	if reportHook != nil {
		reportHook(line)
		return
	}
	fmt.Fprintln(os.Stderr, line)
}

// callerSite returns the first caller outside this package.
func callerSite() string {
	var pcs [8]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		if !strings.Contains(f.Function, "/common/layout.") {
			file := f.File
			if i := strings.Index(file, "/src/"); i >= 0 {
				file = file[i+len("/src/"):]
			}
			return fmt.Sprintf("%s:%d", file, f.Line)
		}
		if !more {
			return "?"
		}
	}
}

// Probes.

type probeKey struct {
	site string
	typ  string
	off  int // relative to the base pointer
}

type probeVal struct {
	count            int
	width            uintptr
	leaf             string
	off64            uintptr
	moved, resized   bool
	untyped, unknown bool
}

var probes struct {
	sync.Mutex
	m     map[probeKey]*probeVal
	dirty bool
}

// Probe returns unsafe.Add(base, off) and records the block type the access touches.
// It is inserted by instrumented builds only.
func Probe(base unsafe.Pointer, off int, width uintptr, site string) unsafe.Pointer {
	p := unsafe.Add(base, off)
	if base == nil {
		return p
	}
	b := find(uintptr(base))
	if b == nil {
		b = find(uintptr(p))
	}
	k := probeKey{site: site, off: off}
	v := probeVal{width: width}
	switch {
	case b == nil:
		k.typ, v.unknown = "?", true
	case b.typ == nil:
		k.typ, v.untyped = "raw:"+b.name, true
	default:
		bt := baseType(b.typ, uintptr(base)-b.start, off, width)
		if bt == nil {
			// the base is not at the start of a value that contains the access
			k.typ = fmt.Sprintf("block:%s+%d", b.typ, int(uintptr(base)-b.start))
			break
		}
		k.typ = bt.String()
		leaf, off64, moved, resized := Moves(bt, uintptr(off), width)
		v.leaf, v.off64, v.moved, v.resized = leaf.Path, off64, moved, resized
	}
	probes.Lock()
	defer probes.Unlock()
	if probes.m == nil {
		probes.m = make(map[probeKey]*probeVal)
		if path := os.Getenv("NOX_LAYOUTPROBE"); path != "" {
			go probeWriter(path)
		}
	}
	if x := probes.m[k]; x != nil {
		x.count++
		return p
	}
	v.count = 1
	probes.m[k] = &v
	probes.dirty = true
	return p
}

// baseType returns the outermost type of the value that starts at offset bo of a t and contains
// the access [bo+off, bo+off+width). An array is replaced by its element if the access fits in the first one.
func baseType(t reflect.Type, bo uintptr, off int, width uintptr) reflect.Type {
	for bo != 0 {
		switch t.Kind() {
		case reflect.Struct:
			found := false
			for i := t.NumField() - 1; i >= 0; i-- {
				f := t.Field(i)
				if bo >= f.Offset && bo < f.Offset+max(f.Type.Size(), 1) {
					bo -= f.Offset
					t = f.Type
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		case reflect.Array:
			es := t.Elem().Size()
			if es == 0 {
				return nil
			}
			bo %= es
			t = t.Elem()
		default:
			return nil
		}
	}
	for {
		if off < 0 || uintptr(off)+max(width, 1) > t.Size() {
			return nil
		}
		if t.Kind() == reflect.Array && uintptr(off)+max(width, 1) <= t.Elem().Size() {
			t = t.Elem()
			continue
		}
		return t
	}
}

// WriteProbes writes the probe results, one line per site, block type and offset.
func WriteProbes(path string) error {
	probes.Lock()
	type line struct {
		k probeKey
		v probeVal
	}
	var lines []line
	for k, v := range probes.m {
		lines = append(lines, line{k, *v})
	}
	probes.dirty = false
	probes.Unlock()
	sort.Slice(lines, func(i, j int) bool {
		a, b := lines[i].k, lines[j].k
		if a.site != b.site {
			return a.site < b.site
		}
		if a.typ != b.typ {
			return a.typ < b.typ
		}
		return a.off < b.off
	})
	var sb strings.Builder
	for _, l := range lines {
		flags := ""
		if l.v.moved {
			flags += " MOVED"
		}
		if l.v.resized {
			flags += " RESIZED"
		}
		fmt.Fprintf(&sb, "%s\t%s\t%d\t%d\t%s\t%d\t%d%s\n", l.k.site, l.k.typ, l.k.off, l.v.width, l.v.leaf, l.v.off64, l.v.count, flags)
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func probeWriter(path string) {
	for {
		time.Sleep(2 * time.Second)
		probes.Lock()
		dirty := probes.dirty
		probes.Unlock()
		if dirty {
			if err := WriteProbes(path); err != nil {
				fmt.Fprintln(os.Stderr, "layout probes:", err)
			}
		}
	}
}

// TrackLike records that the block at p has the type and name of the tracked block at src (a copy of it).
func TrackLike(p, src unsafe.Pointer) {
	if !Checking || p == nil || src == nil {
		return
	}
	reg.Lock()
	b := reg.starts[uintptr(src)]
	reg.Unlock()
	if b != nil {
		Track(p, b.size, b.typ, b.name)
	}
}

// ProbeBase records an offset access like Probe and returns base. It wraps the base of offset helper calls;
// P must be a pointer type.
func ProbeBase[P any](base P, off int, width uintptr, site string) P {
	Probe(*(*unsafe.Pointer)(unsafe.Pointer(&base)), off, width, site)
	return base
}

// reportOnce logs a check failure once per key.
func reportOnce(key, format string, args ...any) {
	reported.Lock()
	defer reported.Unlock()
	if reported.seen == nil {
		reported.seen = make(map[string]bool)
	}
	if reported.seen[key] {
		return
	}
	reported.seen[key] = true
	line := "layoutcheck: " + fmt.Sprintf(format, args...)
	if reportHook != nil {
		reportHook(line)
		return
	}
	fmt.Fprintln(os.Stderr, line)
}

// CheckRaw checks an access of width bytes at the 386 offset off from base, made by a raw offset helper
// (unsafe.Add with a number from the original binary). Such accesses are only valid in memory that keeps the
// 386 layout: blobs and untyped blocks. With checks enabled, it reports raw accesses into a tracked Go struct
// that land on a different field on the current architecture. skip is the number of helper frames between the
// access site and CheckRaw.
func CheckRaw(base unsafe.Pointer, off int, width uintptr, skip int) {
	if !Checking || base == nil {
		return
	}
	b := find(uintptr(base))
	if b == nil || b.typ == nil {
		return
	}
	tb := tableFor(b.typ)
	// The base may point inside the block: find its 386 offset.
	rel := uint32(uintptr(base) - b.start)
	base386 := -1
	for i, v := range tb.off {
		if v == rel {
			base386 = i
			break
		}
	}
	site := callerSiteSkip(skip)
	if base386 < 0 {
		reportOnce("raw "+site, "raw %s: base at +%d of %s is not a field start on 386", site, rel, b.typ)
		return
	}
	at := base386 + off
	if at < 0 || at+int(width) > len(tb.off) {
		reportOnce("raw "+site, "raw %s: 386 offset %d (width %d) outside %s", site, at, width, b.typ)
		return
	}
	for i := 0; i < int(width); i++ {
		if v := tb.off[at+i]; v == padding || uintptr(v) != uintptr(rel)+uintptr(off)+uintptr(i) {
			leaf, _ := model386.leafAt(b.typ, uintptr(at))
			reportOnce("raw "+site, "raw %s: 386 offset %d of %s (%s) is at a different place on this architecture", site, at, b.typ, leaf.Path)
			return
		}
	}
}

// tableFor returns the offset table of t (see Of).
func tableFor(t reflect.Type) *Table {
	if v, ok := tables.Load(t); ok {
		return v.(*Table)
	}
	v, _ := tables.LoadOrStore(t, build(t))
	return v.(*Table)
}

// callerSiteSkip returns the call site skip frames above the first caller outside this package.
func callerSiteSkip(skip int) string {
	var pcs [16]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		if !strings.Contains(f.Function, "/common/layout.") {
			if skip == 0 {
				file := f.File
				if i := strings.Index(file, "/src/"); i >= 0 {
					file = file[i+len("/src/"):]
				}
				return fmt.Sprintf("%s:%d", file, f.Line)
			}
			skip--
		}
		if !more {
			return "?"
		}
	}
}
