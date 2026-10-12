package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func findBorderName(name *byte) int32 {
	count := int32(dword_5d4594_251572)
	if name == nil || count <= 0 {
		return -1
	}
	want := GoStringP(unsafe.Pointer(name))
	// The loader caps active count at the physical 64 rows. Unlike tile
	// lookup, this scan is exact, case-sensitive and stops at the first match.
	for i := int32(0); i < count; i++ {
		p := (*byte)(memmap.PtrOff(0x85B3FC, 28644+60*uintptr(i)))
		if GoStringP(unsafe.Pointer(p)) == want {
			return i
		}
	}
	return -1
}

func selectBorderName(name *byte) bool {
	dword_5d4594_2489436 = 0
	// Only the sentinel is case-insensitive; preserve the shared C locale
	// comparator rather than applying Unicode folding to border names.
	if textCompareNarrow((*byte)(unsafe.Pointer(alloc.InternCString("NONE"))), (*byte)(unsafe.Pointer(name))) == 0 {
		dword_5d4594_3835356 = 255
		return true
	}
	return selectBorderPrimary(findBorderName(name))
}

func selectBorderPrimary(index int32) bool {
	if index < 0 || index >= int32(dword_5d4594_251572) {
		return false
	}
	dword_5d4594_3835356 = uint32(index)
	dword_5d4594_2489436 = 1
	return true
}

func selectBorderVariation(variation int32) bool {
	if dword_5d4594_2489436 == 0 {
		return true
	}
	count := int32(dword_5d4594_251572)
	selected := uint32(dword_5d4594_3835356)
	// Approved fix: validate the selected border's limit, not the unrelated
	// row indexed by variation. Reject invalid selected rows before access.
	if count <= 0 || selected >= uint32(count) || selected >= 64 || variation < 0 ||
		variation >= int32(memmap.Uint16(0x85B3FC, 28688+60*uintptr(selected))) {
		return false
	}
	dword_5d4594_3835360 = uint32(variation)
	return true
}
