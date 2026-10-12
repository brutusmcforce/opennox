//go:build safe

package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func legacyCalloc(num, size uintptr) unsafe.Pointer {
	ptr, _ := alloc.Calloc(int(num), size)
	return ptr
}
func legacyRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer { return alloc.Realloc(ptr, size) }
func legacyFree(ptr unsafe.Pointer)                                 { alloc.FreePtr(ptr) }

func legacyMalloc(size uintptr) unsafe.Pointer { p, _ := alloc.Malloc(size); return p }
