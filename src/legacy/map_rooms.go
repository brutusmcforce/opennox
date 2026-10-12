package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/platform"
	"github.com/brutusmcforce/opennox/libs/types"
	"unsafe"
)

// These records retain the 32-bit map-generation layout and address links.
type mapRoom struct {
	Kind                int32
	Grid                [2]int32
	Width, Height       int32
	Pos, Size, Min, Max types.Pointf
	Flags               uint32
	Next, Prev          *mapRoom
	Reserved64          [6]uint32
	Neighbors           [4][8]*mapRoom
	Counts              [4]byte
	Reserved220         uint32
	Points              [16]types.Pointf
	PointCount          byte
	Reserved353         [11]byte
	ThemeFlags          uint32
	Exclusions          *mapRoomExclusion
	Decoration          unsafe.Pointer
}
type mapRoomExclusion struct {
	Kind     int32
	Min, Max types.Pointf
	Prefab   uint32
	Next     *mapRoomExclusion
}
type mapRoomCellData struct {
	Flags    uint32
	X, Y     int32
	Reserved uint32
	Room     *mapRoom
}

func mapRoomGlobalWord(i int) *uint32 {
	switch i {
	case 0:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487532))
	case 1:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487536))
	case 2:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487540))
	case 3:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487556))
	case 4:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487560))
	}
	panic("map room global")
}
func mapRoomPointer(v uint32) unsafe.Pointer              { return memmap.WordPtr(v) }
func mapRoomRaw(p unsafe.Pointer) uint32                  { return memmap.PtrWord(p) }
func mapRoomWord(p unsafe.Pointer, off int) *int32        { return (*int32)(unsafe.Add(p, off)) }
func mapRoomFloatWord(p unsafe.Pointer, off int) *float32 { return (*float32)(unsafe.Add(p, off)) }
func mapRoomRef(p unsafe.Pointer, off int) *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Add(p, off))
}
func mapRoomByte(p unsafe.Pointer, off int) *byte { return (*byte)(unsafe.Add(p, off)) }

// Keep allocation and release in the same legacy profile-specific domain.
func mapRoomCalloc(count uint32, size uintptr) unsafe.Pointer {
	return legacyCalloc(uintptr(count), uintptr(size))
}
func mapRoomRelease(p unsafe.Pointer) { mapRoomBeforeRelease(p); legacyFree(p) }
func mapRoomScratchAlloc() uint32 {
	p := mapRoomCalloc(1, 8192)
	*mapRoomGlobalWord(3) = mapRoomRaw(p)
	return uint32(bool2int(p != nil))
}
func mapRoomScratchFree()   { mapRoomRelease(mapRoomPointer(*mapRoomGlobalWord(3))) }
func mapRoomHead() *mapRoom { return (*mapRoom)(mapRoomPointer(*mapRoomGlobalWord(4))) }
func mapRoomNext(r *mapRoom) *mapRoom {
	if r == nil {
		return nil
	}
	return r.Next
}
func mapRoomAdd(r *mapRoom) int32 {
	r.Prev = nil
	r.Next = mapRoomHead()
	if r.Next != nil {
		r.Next.Prev = r
	}
	*mapRoomGlobalWord(4) = mapRoomRaw(unsafe.Pointer(r))
	return mapRoomOccupy(r)
}
func mapRoomRemove(r *mapRoom) int32 {
	if r.Prev != nil {
		r.Prev.Next = r.Next
	} else {
		*mapRoomGlobalWord(4) = mapRoomRaw(unsafe.Pointer(r.Next))
	}
	if r.Next != nil {
		r.Next.Prev = r.Prev
	}
	return mapRoomVacate(r)
}
func mapRoomNew(w, h int32) *mapRoom {
	r := (*mapRoom)(mapRoomCalloc(1, 376))
	if r != nil {
		r.Kind = 1
		r.Width = w
		r.Height = h
		r.Size = types.Ptf(float32(float64(w)*32.526913), float32(float64(h)*32.526913))
	}
	return r
}
func mapRoomPrepare(config unsafe.Pointer) *mapRoom {
	mean, spread := *mapRoomWord(config, 32), *mapRoomWord(config, 36)
	minimum := mean - spread
	w := mapRoomRandomCentered(mean, spread)
	if w < 5 {
		w = 5
	}
	h := mapRoomRandomCentered(w, int32(int64(float64(w)*0.30000001)))
	if h < 5 {
		h = 5
	}
	if h < minimum {
		h = minimum
	}
	if *mapRoomWord(config, 56) != 0 {
		mapRoomRandomInt(0, 2)
	}
	return mapRoomNew(w, h)
}
func mapRoomFree(r *mapRoom) {
	for node := r.Exclusions; node != nil; {
		next := node.Next
		mapRoomRelease(unsafe.Pointer(node))
		node = next
	}
	mapRoomRelease(unsafe.Pointer(r))
}
func mapRoomFreeAll() *mapRoom {
	for r := mapRoomHead(); r != nil; {
		next := r.Next
		mapRoomFree(r)
		r = next
	}
	*mapRoomGlobalWord(4) = 0
	return nil
}
func mapRoomHasRoomNeighbor(r *mapRoom, dir int32) uint32 {
	for i := 0; i < int(r.Counts[dir]); i++ {
		if n := r.Neighbors[dir][i]; n != nil && n.Kind == 1 {
			return 1
		}
	}
	return 0
}
func mapRoomConnect(r, t *mapRoom, dir int32) uint32 {
	n := r.Counts[dir]
	if n >= 8 {
		return 0
	}
	r.Counts[dir]++
	r.Neighbors[dir][n] = t
	return 1
}
func mapRoomConnectBoth(r, t *mapRoom, dir int32) uint32 {
	mapRoomConnect(r, t, dir)
	return mapRoomConnect(t, r, mapRoomOpposite(dir))
}
func mapRoomIsEntranceSide(r *mapRoom, dir int32) uint32 {
	var want int32
	switch r.Kind {
	case 1:
		return 0
	case 2:
		want = 1
	case 3:
		want = 0
	case 4:
		want = 3
	case 5:
		want = 2
	default:
		want = int32(mapRoomRaw(unsafe.Pointer(r)))
	}
	return uint32(bool2int(dir == want))
}
func mapRoomSeed(seed uint32) { platform.RandSeed(int64(seed)) }
func mapRoomRandomInt(lo, hi int32) int32 {
	v := int32(uint32(lo) + uint32(hi-lo+1)*uint32(nox_platform_rand())/0x7fff)
	if v > hi {
		v = hi
	}
	return v
}
func mapRoomRandomCentered(mean, spread int32) int32 {
	if spread == 0 {
		return mean
	}
	scale := int32(10000) / spread
	for {
		v := mapRoomRandomInt(-spread, spread)
		y := mapRoomRandomInt(0, 10000/spread)
		if y >= scale*(spread-v)/spread {
			continue
		}
		if y >= scale*(v+spread)/spread {
			continue
		}
		return v + mean
	}
}
func mapRoomRandomFloat(lo, hi float32) float64 {
	return float64(uint32(nox_platform_rand()))*(float64(hi)-float64(lo))*0.000030518509 + float64(lo)
}
