// Package layout translates struct offsets of the original 32-bit layout to the current architecture.
//
// Legacy code accesses many struct fields by byte offset (unsafe.Add(p, 432)). Offsets are written as in the
// original 32-bit (386) binary; Of[T]().Off(off) returns the offset of the same field (and the same position
// inside it) on the current architecture.
package layout

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Table maps 386 offsets of one struct type to offsets on the current architecture.
type Table struct {
	typ  reflect.Type
	size uintptr // size on 386
	off  []uint32
}

var tables sync.Map // reflect.Type -> *Table

// Of returns the offset table of a struct type.
func Of[T any]() *Table {
	t := reflect.TypeFor[T]()
	if v, ok := tables.Load(t); ok {
		return v.(*Table)
	}
	tb := build(t)
	v, _ := tables.LoadOrStore(t, tb)
	return v.(*Table)
}

// Off translates a 386 offset into the struct to the current architecture.
func (t *Table) Off(off uintptr) uintptr {
	if off >= uintptr(len(t.off)) {
		if off == t.size {
			return t.typ.Size() // end of the struct
		}
		panic(fmt.Errorf("layout: offset %d beyond %s (386 size %d)", off, t.typ, t.size))
	}
	v := t.off[off]
	if v == padding {
		panic(fmt.Errorf("layout: offset %d of %s is padding", off, t.typ))
	}
	return uintptr(v)
}

// padding marks bytes of the 386 layout that are not part of a field.
const padding = ^uint32(0)

func build(t reflect.Type) *Table {
	size := size386(t)
	tb := &Table{typ: t, size: size}
	tb.off = make([]uint32, size)
	for i := range tb.off {
		tb.off[i] = padding
	}
	fill(tb.off, t, 0, 0)
	return tb
}

// fill maps every byte of a value of type t at off32 (386) to off (current).
func fill(m []uint32, t reflect.Type, off32, off uintptr) {
	switch t.Kind() {
	case reflect.Struct:
		var o32 uintptr
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			o32 = align(o32, align386(f.Type))
			fill(m, f.Type, off32+o32, off+f.Offset)
			o32 += size386(f.Type)
		}
	case reflect.Array:
		e32 := size386(t.Elem())
		for i := 0; i < t.Len(); i++ {
			fill(m, t.Elem(), off32+uintptr(i)*e32, off+uintptr(i)*t.Elem().Size())
		}
	default:
		// A leaf field: bytes inside it keep their position (pointers keep their first 4 bytes).
		for i := uintptr(0); i < size386(t); i++ {
			m[off32+i] = uint32(off + i)
		}
	}
}

func align(x, a uintptr) uintptr { return (x + a - 1) &^ (a - 1) }

// model is the size and alignment rules of a target architecture (gc).
type model struct {
	ptr uintptr // pointer and int size
	i64 uintptr // alignment of 64-bit values
}

var (
	model386   = model{ptr: 4, i64: 4}
	modelAmd64 = model{ptr: 8, i64: 8}
)

func size386(t reflect.Type) uintptr  { return model386.size(t) }
func align386(t reflect.Type) uintptr { return model386.align(t) }

// size and align compute the size and alignment of a type as on the model architecture.
func (m model) size(t reflect.Type) uintptr {
	switch t.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr, reflect.Int, reflect.Uint, reflect.Map, reflect.Chan, reflect.Func:
		return m.ptr
	case reflect.String, reflect.Interface:
		return 2 * m.ptr
	case reflect.Slice:
		return 3 * m.ptr
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	case reflect.Array:
		return m.size(t.Elem()) * uintptr(t.Len())
	case reflect.Struct:
		var o, a uintptr = 0, 1
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type
			fa := m.align(ft)
			o = align(o, fa) + m.size(ft)
			a = max(a, fa)
		}
		if t.NumField() > 0 && m.size(t.Field(t.NumField()-1).Type) == 0 && o > 0 {
			o++ // a final zero-size field gets padding
		}
		return align(o, a)
	default:
		return t.Size() // bool, int8..int32, float32
	}
}

func (m model) align(t reflect.Type) uintptr {
	switch t.Kind() {
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex128:
		return m.i64
	case reflect.Complex64:
		return 4
	case reflect.Array:
		return m.align(t.Elem())
	case reflect.Struct:
		var a uintptr = 1
		for i := 0; i < t.NumField(); i++ {
			a = max(a, m.align(t.Field(i).Type))
		}
		return a
	case reflect.String, reflect.Slice, reflect.Interface:
		return m.ptr
	default:
		return min(m.size(t), m.ptr)
	}
}

// Leaf is the innermost field at an offset.
type Leaf struct {
	Path string       // field path, e.g. "Field12[3].X"
	Type reflect.Type // type of the leaf
	Off  uintptr      // offset of the leaf start
	Size uintptr
}

// leafAt finds the leaf field containing offset off of a value of type t, as on the model architecture.
func (m model) leafAt(t reflect.Type, off uintptr) (Leaf, bool) {
	var path string
	var base uintptr
	for {
		switch t.Kind() {
		case reflect.Struct:
			var o uintptr
			found := false
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				o = align(o, m.align(f.Type))
				sz := m.size(f.Type)
				if off >= o && off < o+sz {
					if path != "" {
						path += "."
					}
					path += fieldName(t, i)
					base += o
					off -= o
					t = f.Type
					found = true
					break
				}
				o += sz
			}
			if !found {
				return Leaf{}, false // padding
			}
		case reflect.Array:
			es := m.size(t.Elem())
			if es == 0 {
				return Leaf{}, false
			}
			i := off / es
			path += fmt.Sprintf("[%d]", i)
			base += i * es
			off -= i * es
			t = t.Elem()
		default:
			return Leaf{Path: path, Type: t, Off: base, Size: m.size(t)}, true
		}
	}
}

// offsetOf returns the offset of a field path (as returned by leafAt) on the model architecture.
func (m model) offsetOf(t reflect.Type, path string) uintptr {
	var off uintptr
	for path != "" {
		if path[0] == '[' {
			j := strings.IndexByte(path, ']')
			i, _ := strconv.Atoi(path[1:j])
			off += uintptr(i) * m.size(t.Elem())
			t = t.Elem()
			path = strings.TrimPrefix(path[j+1:], ".")
			continue
		}
		j := strings.IndexAny(path, ".[")
		name := path
		if j >= 0 {
			name, path = path[:j], strings.TrimPrefix(path[j:], ".")
		} else {
			path = ""
		}
		var o uintptr
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			o = align(o, m.align(f.Type))
			if fieldName(t, i) == name {
				off += o
				t = f.Type
				break
			}
			o += m.size(f.Type)
		}
	}
	return off
}

// fieldName names a field in a path; blank fields are named by index.
func fieldName(t reflect.Type, i int) string {
	if n := t.Field(i).Name; n != "_" {
		return n
	}
	return "#" + strconv.Itoa(i)
}

// Moves reports how an access of width bytes at the 386 offset off of type t changes on amd64:
// the field it reads, the amd64 offset, and if the field moves or changes size.
func Moves(t reflect.Type, off, width uintptr) (leaf Leaf, off64 uintptr, moved, resized bool) {
	l, ok := model386.leafAt(t, off)
	if !ok {
		return Leaf{Path: "<padding>"}, off, false, false
	}
	o64 := modelAmd64.offsetOf(t, l.Path) + (off - l.Off)
	sz64 := modelAmd64.size(l.Type)
	return l, o64, o64 != off, sz64 != l.Size || width > l.Size
}

// Runs calls fn for the pieces of the current architecture that hold the 386 byte range [off, off+n) of the
// struct, in order: off386 is the 386 offset of the piece, off its offset on the current architecture.
// A field that is larger on the current architecture (a pointer) only gets its first bytes; padding is skipped.
func (t *Table) Runs(off, n uintptr, fn func(off386, off, n uintptr)) {
	if n == 0 {
		return
	}
	var start386, start, prev, cnt uintptr
	for i := uintptr(0); i < n; i++ {
		if off+i >= uintptr(len(t.off)) {
			panic(fmt.Errorf("layout: offset %d beyond %s (386 size %d)", off+i, t.typ, t.size))
		}
		v := t.off[off+i]
		if v == padding {
			continue
		}
		cur := uintptr(v)
		if cnt != 0 && cur == prev+1 {
			prev = cur
			cnt++
			continue
		}
		if cnt != 0 {
			fn(start386, start, cnt)
		}
		start386, start, prev, cnt = off+i, cur, cur, 1
	}
	if cnt != 0 {
		fn(start386, start, cnt)
	}
}
