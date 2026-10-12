package alloc

import "unsafe"

// RawCalloc allocates zeroed unmanaged memory without registering it in allocs.
func RawCalloc(num, size uintptr) unsafe.Pointer { return rawCalloc(num, size) }

// RawFree releases unmanaged memory without consulting or changing allocs.
func RawFree(ptr unsafe.Pointer) { rawFree(ptr) }
