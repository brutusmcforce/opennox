package legacy

import (
	"unsafe"
)

// Matching allocator boundary for spriteDataAlloc's raw libc calloc allocations.
func spriteDataFree(p unsafe.Pointer) {
	if p != nil {
		legacyFree(p)
	}
}

func spriteVectorFree(p unsafe.Pointer) {
	for direction := 0; direction < 8; direction++ {
		slot := direction
		if slot >= 4 {
			slot++
		} // middle word is metadata, not an owned frame list
		spriteDataFree(*(*unsafe.Pointer)(unsafe.Add(p, slot*4)))
	}
}
func spriteDataFreeKind(p unsafe.Pointer, kind int) {
	switch kind {
	case 2, 3:
		spriteDataFree(*(*unsafe.Pointer)(unsafe.Add(p, 4)))
	case 4:
		for i := 0; i < 5; i++ {
			spriteDataFree(*(*unsafe.Pointer)(unsafe.Add(p, 4+4*i)))
		}
	case 5:
		spriteVectorFree(unsafe.Add(p, 4))
	case 6:
		for state := 0; state < 55; state++ {
			for slot := 0; slot < 54; slot++ {
				group := *(*unsafe.Pointer)(unsafe.Add(p, 52+264*state+4*slot))
				if group != nil {
					spriteVectorFree(unsafe.Add(group, 4))
					spriteDataFree(group)
				}
			}
		}
	case 7, 8:
		groups := 16
		if kind == 8 {
			groups = 3
		}
		for i := 0; i < groups; i++ {
			spriteVectorFree(unsafe.Add(p, 8+48*i))
		}
	}
	spriteDataFree(p)
}
