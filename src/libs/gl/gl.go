// Package gl is a minimal OpenGL 3.3 core binding without cgo.
//
// It has the functions and constants OpenNox uses, with the same names and types as github.com/go-gl/gl/v3.3-core/gl.
// The functions are loaded from the current context by Init (with SDL_GL_GetProcAddress) and called with purego.
// funcs.go and consts.go were generated from the go-gl signatures and constants; add new functions the same way.
package gl

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/brutusmcforce/opennox/libs/sdl"
)

// glClearColor takes float arguments, which purego.SyscallN can't pass, so it's bound with purego.RegisterFunc.
var glClearColor func(red, green, blue, alpha float32)

// Init loads the OpenGL functions of the current context. It must be called after the context is created.
func Init() error {
	var missing []string
	for name, p := range procNames() {
		if *p = uintptr(sdl.GLGetProcAddress(name)); *p == 0 {
			missing = append(missing, name)
		}
	}
	p := uintptr(sdl.GLGetProcAddress("glClearColor"))
	if p == 0 {
		missing = append(missing, "glClearColor")
	} else {
		purego.RegisterFunc(&glClearColor, p)
	}
	if len(missing) != 0 {
		sort.Strings(missing)
		return fmt.Errorf("OpenGL functions not found: %s", strings.Join(missing, ", "))
	}
	return nil
}

func call(f uintptr, args ...uintptr) uintptr {
	r, _, _ := purego.SyscallN(f, args...)
	return r
}

func cbool(v bool) uintptr {
	if v {
		return 1
	}
	return 0
}

func ClearColor(red, green, blue, alpha float32) {
	glClearColor(red, green, blue, alpha)
}

// Ptr returns a pointer to the data of a pointer or a slice, like go-gl's Ptr.
func Ptr(data any) unsafe.Pointer {
	if data == nil {
		return nil
	}
	switch v := data.(type) {
	case unsafe.Pointer:
		return v
	case uintptr:
		return unsafe.Pointer(v) //nolint:govet // an offset into a buffer object, like go-gl
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Pointer:
		return v.UnsafePointer()
	case reflect.Slice:
		if v.Len() == 0 {
			return nil
		}
		return v.Index(0).Addr().UnsafePointer()
	}
	panic(fmt.Errorf("gl.Ptr: unsupported type %T", data))
}

// PtrOffset returns an offset into a buffer object as a pointer, like go-gl's PtrOffset.
func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset)) //nolint:govet // an offset, not a Go pointer
}

// Str returns a pointer to a NUL-terminated string, like go-gl's Str.
func Str(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		panic("gl.Str: string must be NUL-terminated: " + str)
	}
	return unsafe.StringData(str)
}

// Strs returns an array of pointers to NUL-terminated strings, like go-gl's Strs.
// The strings must stay alive while the array is used; free does nothing (it's for compatibility).
func Strs(strs ...string) (cstrs **uint8, free func()) {
	if len(strs) == 0 {
		panic("gl.Strs: expected at least 1 string")
	}
	ptrs := make([]*uint8, len(strs))
	for i, s := range strs {
		ptrs[i] = Str(s)
	}
	return &ptrs[0], func() {}
}

// GoStr copies a NUL-terminated string, like go-gl's GoStr.
func GoStr(cstr *uint8) string {
	if cstr == nil {
		return ""
	}
	n := 0
	for *(*uint8)(unsafe.Add(unsafe.Pointer(cstr), n)) != 0 {
		n++
	}
	return string(unsafe.Slice(cstr, n))
}
