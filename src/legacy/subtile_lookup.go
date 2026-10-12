package legacy

import "unsafe"

func subtileContains(point *[2]int32, category int32) bool {
	x, y := point[0], point[1]
	diagonal := int32(46) - y
	switch category {
	case 0:
		return x <= y && x >= diagonal
	case 1:
		return x <= y
	case 2:
		return x <= y && x <= diagonal
	case 3:
		return x >= diagonal
	case 4:
		return x <= diagonal
	case 5:
		return x >= y && x >= diagonal
	case 6:
		return x >= y
	case 7:
		return x >= y && x <= diagonal
	case 8:
		return x <= y || x >= diagonal
	case 9:
		return x <= y || x <= diagonal
	case 10:
		return x >= y || x <= diagonal
	case 11:
		return x >= y || x >= diagonal
	default:
		return false
	}
}

func findSubtileAt(head *[5]uint32, point *[2]int32, fallback int32) int32 {
	result := fallback
	// The list belongs to C. Keep its 32-bit links and physical word layout.
	for p := head; p != nil; p = (*[5]uint32)(unsafe.Pointer(uintptr(p[4]))) {
		category := normalizeBorderEdge(int32(p[2]), int32(p[3]))
		if subtileContains(point, category) {
			result = int32(p[0])
		}
	}
	return result
}
