package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

const (
	paintVariation = iota
	paintSubtileMode
	paintBorder
	paintBorderVariation
	paintAfterWalls
	paintWallMerge
	paintWallCycle
	paintObjectType
	paintObjectCounter
	paintFreeHead
	paintWorkCount
)

func mapPaintGlobal(i int) *uint32 {
	switch i {
	case paintVariation:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835348))
	case paintSubtileMode:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835352))
	case paintBorder:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835356))
	case paintBorderVariation:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835360))
	case paintAfterWalls:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835364))
	case paintWallMerge:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835368))
	case paintWallCycle:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835372))
	case paintObjectType:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835388))
	case paintObjectCounter:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_3835392))
	case paintFreeHead:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_588084))
	case paintWorkCount:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487248))
	}
	panic("map painting global")
}
func mapPaintCell(x, y int32) *[11]uint32 {
	rows := (*[128]*[128][11]uint32)(unsafe.Pointer(worldTileGrid))
	return &rows[x][y]
}
func mapPaintSelection(off uintptr) *uint32 { return memmap.PtrUint32(0x973F18, off) }
func mapPaintTransform(in, out *types.Pointf) uint32 {
	if in == nil || out == nil {
		return 0
	}
	// Reload after storing X: legacy callers may alias both points.
	out.X = float32((float64(in.Y)+float64(in.X))*0.70710677 + 2957)
	out.Y = float32((float64(in.Y)-float64(in.X))*0.70710677 + 2956)
	if out.X <= 80.5 {
		out.X = 82.5
	}
	if out.Y <= 80.5 {
		out.Y = 81.5
	}
	if out.X >= 5853.5 {
		out.X = 5851.5
	}
	if out.Y >= 5853.5 {
		out.Y = 5852.5
	}
	return 1
}
func mapPaintDirection(v int32) int32 {
	switch v {
	case 0:
		return 3
	case 1:
		return 0
	case 2:
		return 1
	case 3:
		return 6
	case 5:
		return 2
	case 6:
		return 7
	case 7:
		return 8
	case 8:
		return 5
	}
	return -1
}
func mapPaintWallCompose(a, b byte) byte {
	return memmap.Uint8(0x587000, 71276+13*uintptr(a)+uintptr(b))
}
func mapPaintAfterWalls(v int32) uint32 {
	if v != 0 && v != 1 {
		return 0
	}
	*mapPaintGlobal(paintAfterWalls) = uint32(v)
	return 1
}
func mapPaintMergeWalls(v int32) uint32 {
	if v != 0 && v != 1 {
		return 0
	}
	*mapPaintGlobal(paintWallMerge) = uint32(v)
	return 1
}
func mapPaintWallDirection(v int32) uint32 {
	if v < 0 || v >= 15 {
		*mapPaintSelection(35952) = 0
		return 0
	}
	*mapPaintSelection(35952) = uint32(v)
	return 1
}
func mapPaintLayout(d unsafe.Pointer) unsafe.Pointer {
	n := mapRoomRandomInt(0, *mapRoomWord(d, 88)-1)
	p := *mapRoomRef(d, 84)
	for ; n != 0; n-- {
		p = *mapRoomRef(p, 124)
	}
	return p
}
