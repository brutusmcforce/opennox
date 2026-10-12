package alloc

import "unsafe"

//go:noescape
func rawClear(ptr unsafe.Pointer, n uintptr)
