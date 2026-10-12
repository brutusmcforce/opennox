//go:build !386

package alloc

import "unsafe"

func rawClear(ptr unsafe.Pointer, n uintptr) { clear(unsafe.Slice((*byte)(ptr), n)) }
