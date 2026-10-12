package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
)

func generateBorderEdge(index, edge int32) int32 {
	off := 60 * uintptr(index)
	w := int32(memmap.Uint8(0x85B3FC, 28696+off))
	h := int32(memmap.Uint8(0x85B3FC, 28697+off))
	if w == 3 && h == 3 {
		return edge
	}
	// Use the existing Logic RNG and its IntClamp contract, including the
	// zero-draw reversed bound at dimension2. Do not eagerly draw a value.
	switch edge {
	case 0:
		return 0
	case 1:
		return int32(nox_common_randomInt_415FA0(1, int(w)-2))
	case 2:
		return w - 1
	case 3:
		return w + 2*int32(nox_common_randomInt_415FA0(0, int(h)-3))
	case 4:
		return w + 2*int32(nox_common_randomInt_415FA0(0, int(h)-3)) + 1
	case 5:
		return w + 2*h - 4
	case 6:
		return int32(nox_common_randomInt_415FA0(1, int(w)-2)) + w + 2*h - 4
	case 7:
		return 2*(w+h) - 5
	default:
		return edge + 2*(w+h) - 12
	}
}

func mergeBorderEdge(record *[4]uint32, category int32) bool {
	// Existing callers supply a valid edge row and normalized category.
	class := normalizeBorderEdge(int32(record[2]), int32(record[3]))
	mapped := memmap.Uint32(0x587000, 282736+4*uintptr(category+12*class))
	if mapped == 255 {
		return false
	}
	edge := uint32(generateBorderEdge(int32(record[2]), int32(mapped)))
	if edge != record[3] {
		record[3] = edge
	}
	return true
}
