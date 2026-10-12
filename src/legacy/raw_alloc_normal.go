//go:build !safe

package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func legacyCalloc(num, size uintptr) unsafe.Pointer { return alloc.RawCalloc(num, size) }
func legacyRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return alloc.RawRealloc(ptr, size)
}
func legacyFree(ptr unsafe.Pointer) { alloc.RawFree(ptr) }

func legacyMalloc(size uintptr) unsafe.Pointer { return alloc.RawMalloc(size) }
