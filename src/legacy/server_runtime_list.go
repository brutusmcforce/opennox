package legacy

import (
	"unsafe"
)

func runtimeRejectedClear(head unsafe.Pointer) uintptr {
	first := listNext((*legacyListNode)(head))
	result := uintptr(unsafe.Pointer(first))
	for p := first; p != nil; {
		next := listNext(p)
		listRemove(p)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	return result
}
