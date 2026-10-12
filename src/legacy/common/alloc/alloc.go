package alloc

import (
	"sync"
	"unsafe"
)

var (
	allocMu sync.Mutex
	allocs  = make(map[unsafe.Pointer]uintptr)
)

func Malloc(size uintptr) (unsafe.Pointer, func()) {
	if size == 0 {
		panic("zero alloc")
	}
	ptr := RawCalloc(1, size)
	if ptr == nil {
		panic("cannot allocate")
	}
	allocMu.Lock()
	allocs[ptr] = size
	allocMu.Unlock()
	return ptr, func() {
		FreePtr(ptr)
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

// New is an analog of builtin new that allocates on C heap.
// Since passing types as parameters is not supported, the first argument must be a value, instead of a type.
func New[T comparable](zero T) (*T, func()) {
	ptr, free := Malloc(unsafe.Sizeof(zero))
	return (*T)(ptr), free
}

// Make is an analog of builtin make that allocates on C heap.
// Since passing types as parameters is not supported, the first argument must be a slice value, instead of a type.
// If the slice value is not empty, the content of it will be copied to the beginning of a new slice.
func Make[T comparable, N number](src []T, n N) ([]T, func()) {
	if n < 0 {
		panic("negative size")
	}
	var elem T
	ptr, free := Calloc(int(n), unsafe.Sizeof(elem))
	out := unsafe.Slice((*T)(ptr), n)
	copy(out, src)
	return out, free
}

// CloneSlice copies the slice to the C heap.
func CloneSlice[T comparable](src []T) ([]T, func()) {
	return Make(src, len(src))
}

func Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	if size == 0 {
		panic("zero alloc")
	}
	old := ptr
	ptr = RawRealloc(ptr, size)
	allocMu.Lock()
	if ptr != old {
		delete(allocs, old)
	}
	allocs[ptr] = size
	allocMu.Unlock()
	return ptr
}

func Calloc(num int, size uintptr) (unsafe.Pointer, func()) {
	if uintptr(num)*size == 0 {
		panic("zero alloc")
	}
	ptr := RawCalloc(uintptr(num), size)
	allocMu.Lock()
	allocs[ptr] = uintptr(num) * size
	allocMu.Unlock()
	return ptr, func() {
		FreePtr(ptr)
	}
}

func FreePtr(ptr unsafe.Pointer) {
	allocMu.Lock()
	_, ok := allocs[ptr]
	delete(allocs, ptr)
	allocMu.Unlock()
	if !ok {
		panic("incorrect free")
	}
	RawFree(ptr)
}

func Free[T comparable](ptr *T) {
	FreePtr(unsafe.Pointer(ptr))
}

func FreeSlice[T comparable](b []T) {
	b = b[:1]
	Free(&b[0])
}

func strlen[T comparable](s *T) int {
	if s == nil {
		return 0
	}
	var zero T
	elem := unsafe.Sizeof(zero)
	n := 0
	for *s != zero {
		n++
		s = (*T)(unsafe.Add(unsafe.Pointer(s), elem))
	}
	return n
}

func strlenN[T comparable](s *T, max int) int {
	if s == nil {
		return 0
	}
	var zero T
	elem := unsafe.Sizeof(zero)
	n := 0
	for *s != zero {
		s = (*T)(unsafe.Add(unsafe.Pointer(s), elem))
		if n == max {
			return n
		}
		n++
	}
	return n
}

func strlenS[T comparable](s []T) int {
	if s == nil {
		return 0
	}
	var zero T
	for i, v := range s {
		if v == zero {
			return i
		}
	}
	return len(s)
}

func ZeroTermLen[T comparable](s *T) int {
	return strlen(s)
}

func StrLen[T number](s *T) int {
	if s == nil {
		return 0
	}
	n := strlen(s)
	logMemReadString(unsafe.Pointer(s), uintptr(n+1))
	return n
}

func StrLenN[T number](s *T, max int) int {
	if s == nil {
		return 0
	}
	n := strlenN(s, max)
	logMemReadString(unsafe.Pointer(s), uintptr(max))
	return n
}

func StrLenS[T number](s []T) int {
	if len(s) == 0 {
		return 0
	}
	n := strlenS(s)
	logMemReadString(unsafe.Pointer(&s[0]), uintptr(len(s)))
	return n
}

func Strlen(ptr unsafe.Pointer) int {
	return StrLen((*byte)(ptr))
}
