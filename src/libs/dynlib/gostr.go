package dynlib

import "unsafe"

// GoString copies a NUL-terminated native string.
func GoString(p uintptr) string {
	if p == 0 {
		return ""
	}
	n := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(p), n)) != 0 { //nolint:govet // native memory, not managed by Go
		n++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(p)), n)) //nolint:govet // native memory
}
