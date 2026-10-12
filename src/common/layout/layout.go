// Package layout translates struct offsets of the original 32-bit layout to the current architecture.
//
// Legacy code accesses many struct fields by byte offset (unsafe.Add(p, 432)). Offsets are written as on 386;
// Of[T]().Off(off) returns the offset of the same field (and the same position inside it) on the current
// architecture. On 386 it returns off unchanged.
package layout

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// Table maps 386 offsets of one struct type to offsets on the current architecture.
type Table struct {
	typ  reflect.Type
	size uintptr // size on 386
	off  []uint32
}

var tables sync.Map // reflect.Type -> *Table

const is386 = runtime.GOARCH == "386"

// forceTables builds offset tables on 386 too (for tests).
var forceTables = false

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
	if is386 {
		return off
	}
	if off >= uintptr(len(t.off)) {
		if off == t.size {
			return t.typ.Size() // end of the struct
		}
		panic(fmt.Errorf("layout: offset %d beyond %s (386 size %d)", off, t.typ, t.size))
	}
	return uintptr(t.off[off])
}

func build(t reflect.Type) *Table {
	size := size386(t)
	tb := &Table{typ: t, size: size}
	if is386 && !forceTables {
		return tb
	}
	tb.off = make([]uint32, size)
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

// size386 and align386 compute the size and alignment of a type as on 386 (gc).
func size386(t reflect.Type) uintptr {
	switch t.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr, reflect.Int, reflect.Uint, reflect.Map, reflect.Chan, reflect.Func:
		return 4
	case reflect.String:
		return 8
	case reflect.Slice:
		return 12
	case reflect.Interface:
		return 8
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	case reflect.Array:
		return size386(t.Elem()) * uintptr(t.Len())
	case reflect.Struct:
		var o, a uintptr = 0, 1
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type
			fa := align386(ft)
			o = align(o, fa) + size386(ft)
			a = max(a, fa)
		}
		if t.NumField() > 0 && size386(t.Field(t.NumField()-1).Type) == 0 && o > 0 {
			o++ // a final zero-size field gets padding
		}
		return align(o, a)
	default:
		return t.Size() // bool, int8..int32, float32
	}
}

func align386(t reflect.Type) uintptr {
	switch t.Kind() {
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return 4
	case reflect.Array:
		return align386(t.Elem())
	case reflect.Struct:
		var a uintptr = 1
		for i := 0; i < t.NumField(); i++ {
			a = max(a, align386(t.Field(i).Type))
		}
		return a
	case reflect.String, reflect.Slice, reflect.Interface:
		return 4
	default:
		return min(size386(t), 4)
	}
}
