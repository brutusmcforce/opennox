package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// Population records still cross the 32-bit C map-generator boundary. Keep
// stored addresses as words, and use the existing room/object owners for them.
func populationWord(p uint32, off int) *uint32 { return rawAt[uint32](mapRoomPointer(p), off) }
func populationFloat(p uint32, off int) *float32 {
	return rawAt[float32](mapRoomPointer(p), off)
}
func populationRoom(p uint32) *mapRoom         { return (*mapRoom)(mapRoomPointer(p)) }
func populationObject(p uint32) *server.Object { return (*server.Object)(mapRoomPointer(p)) }
func populationPoint(p uint32) *types.Pointf   { return (*types.Pointf)(mapRoomPointer(p)) }
func populationString(p uint32) string         { return alloc.GoString((*byte)(mapRoomPointer(p))) }
func populationBlob(off uintptr) *uint32       { return memmap.PtrUint32(0x5D4594, off) }
func populationGlobal(i int) *uint32 {
	switch i {
	case 0:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1550916))
	case 1:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487564))
	case 2:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487568))
	case 3:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487576))
	case 4:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487580))
	case 5:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487584))
	case 6:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487620))
	case 7:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487624))
	case 8:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487628))
	case 9:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487632))
	case 10:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487652))
	case 11:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487656))
	case 12:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487672))
	case 13:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487676))
	}
	panic("population global")
}
func mapPopulationStart() uint32    { return *populationGlobal(0) }
func mapPopulationFarthest() uint32 { return *populationGlobal(3) }
func mapPopulationDistanceMax() float64 {
	return float64(*(*float32)(unsafe.Pointer(populationGlobal(4))))
}
