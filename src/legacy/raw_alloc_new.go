package legacy

import (
	"reflect"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
)

// legacyNew allocates a zeroed T with legacyCalloc (free it with legacyFree).
func legacyNew[T any]() *T {
	p := legacyCalloc(1, unsafe.Sizeof(*new(T)))
	layout.TrackAs[T](p)
	return (*T)(p)
}

// mapRoomAlloc allocates a zeroed T with mapRoomCalloc (free it with mapRoomRelease).
func mapRoomAlloc[T any]() *T {
	p := mapRoomCalloc(1, unsafe.Sizeof(*new(T)))
	layout.TrackAs[T](p)
	return (*T)(p)
}

// spriteDataNew allocates zeroed sprite draw data of type T (see spriteDataAlloc).
func spriteDataNew[T any]() *T {
	p := spriteDataAlloc(1, unsafe.Sizeof(*new(T)))
	layout.TrackAs[T](p)
	return (*T)(p)
}

// legacyMake allocates n zeroed Ts with legacyCalloc (free it with legacyFree).
func legacyMake[T any](n int) *T {
	p := legacyCalloc(uintptr(n), unsafe.Sizeof(*new(T)))
	if n > 0 {
		layout.Track(p, uintptr(n)*unsafe.Sizeof(*new(T)), reflect.ArrayOf(n, reflect.TypeFor[T]()), "")
	}
	return (*T)(p)
}

// spriteDataMake allocates n zeroed Ts as sprite draw data (see spriteDataAlloc).
func spriteDataMake[T any](n int) *T {
	p := spriteDataAlloc(n, unsafe.Sizeof(*new(T)))
	if n > 0 {
		layout.Track(p, uintptr(n)*unsafe.Sizeof(*new(T)), reflect.ArrayOf(n, reflect.TypeFor[T]()), "")
	}
	return (*T)(p)
}

// mapRoomMake allocates n zeroed Ts with mapRoomCalloc (free it with mapRoomRelease).
func mapRoomMake[T any](n uint32) *T {
	p := mapRoomCalloc(n, unsafe.Sizeof(*new(T)))
	if n > 0 {
		layout.Track(p, uintptr(n)*unsafe.Sizeof(*new(T)), reflect.ArrayOf(int(n), reflect.TypeFor[T]()), "")
	}
	return (*T)(p)
}
