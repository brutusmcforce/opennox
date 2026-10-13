package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
)

// Dimensions are compared as unsigned words by the original branch selector.
func mapGrowthBranch(r *mapRoom) int32 {
	if r.Kind == 2 || r.Kind == 3 {
		if uint32(r.Width) >= uint32(r.Height) {
			return 1
		}
	} else if uint32(r.Height) >= uint32(r.Width) {
		return 1
	}
	if mapRoomRandomInt(1, 100) <= int32(*populationBlob(1549820)) {
		mapRoomRandomInt(1, 100)
		return 32
	}
	if mapRoomRandomInt(1, 100) <= int32(*populationBlob(1549824)) {
		return 1
	}
	if mapRoomRandomInt(1, 100) >= 50 {
		return 4
	}
	return 2
}
func mapGrowthDispatch(r *mapRoom, depth, halls, width int32, origin *mapRoom) uint32 {
	depth++
	if depth >= int32(*populationBlob(1549868)) {
		return 0
	}
	mapPopulationProgress(155)
	if r.Kind == 1 {
		return mapGrowthFill(r, depth, halls, width, origin)
	}
	return mapGrowthHall(r, depth, halls, width, origin)
}
func mapGrowthFrontiers() {
	for r := populationRoom(*mapGrowthFrontier()); r != nil; r = populationRoom(r.Reserved64[5]) {
		switch r.Kind {
		case 1:
			mapGrowthDispatch(r, 0, 0, 0, nil)
		case 2, 3:
			mapGrowthDispatch(r, 0, 0, r.Width, nil)
		case 4, 5:
			mapGrowthDispatch(r, 0, 0, r.Height, nil)
		}
	}
}
func mapGrowthInitial(cfg uint32) *mapRoom {
	*mapGrowthFrontier() = 0
	if *populationWord(cfg, 0) != 1 {
		r := mapRoomPrepare(mapRoomPointer(cfg))
		*populationGlobal(0) = mapRoomRaw(unsafe.Pointer(r))
		if r != nil {
			p := types.Ptf(0, float32(float64(int32(*populationWord(cfg, 68)))*32.526913-float64(r.Size.Y)+97.580734))
			mapRoomSetPos(r, &p)
			mapRoomAdd(r)
			r.Reserved64[5] = *mapGrowthFrontier()
			*mapGrowthFrontier() = mapRoomRaw(unsafe.Pointer(r))
		}
		return r
	}
	var rooms [128]*mapRoom
	var previous, anchor *mapRoom
	var pos types.Pointf
	count := 0
	for side := 0; side < 4; side++ {
		kind := int32(*memmap.PtrUint32(0x587000, 197924+uintptr(4*side)))
		for step := 0; step < 10; step++ {
			r := mapRoomAlloc[mapRoom]()
			rooms[count] = r
			if r == nil {
				return nil
			}
			r.Kind = kind
			switch kind {
			case 2:
				pos.Y = float32(float64(pos.Y) - 162.63457)
				r.Width = 4
				r.Height = 5
				mapRoomConnect(r, previous, 1)
				mapRoomConnect(previous, r, 0)
			case 3:
				if previous.Kind == 4 {
					pos.X = float32(float64(pos.X) + 32.526913)
					pos.Y = float32(float64(pos.Y) + 130.10765)
				} else {
					pos.Y = float32(float64(pos.Y) + 162.63457)
				}
				r.Width = 4
				r.Height = 5
				mapRoomConnect(r, previous, 0)
				mapRoomConnect(previous, r, 1)
			case 4:
				if previous != nil {
					pos.X = float32(float64(pos.X) + 162.63457)
					mapRoomConnect(r, previous, 3)
					mapRoomConnect(previous, r, 2)
				}
				r.Width = 5
				r.Height = 4
			case 5:
				if previous.Kind == 3 {
					pos.Y = float32(float64(pos.Y) + 32.526913)
				}
				pos.X = float32(float64(pos.X) - 162.63457)
				r.Width = 5
				r.Height = 4
				mapRoomConnect(r, previous, 2)
				mapRoomConnect(previous, r, 3)
				if step == 5 {
					anchor = r
				}
			}
			p := types.Ptf(float32(float64(pos.X)-878.22662), float32(float64(pos.Y)-878.22662))
			r.Size = types.Ptf(float32(float64(r.Width)*32.526913), float32(float64(r.Height)*32.526913))
			mapRoomSetPos(r, &p)
			mapRoomAdd(r)
			previous = r
			count++
		}
	}
	mapRoomConnect(rooms[count-1], rooms[0], 2)
	mapRoomConnect(rooms[0], rooms[count-1], 3)
	previous = anchor
	pos = types.Ptf(anchor.Pos.X, anchor.Size.Y+anchor.Pos.Y)
	for i := 0; i < 8; i++ {
		r := mapRoomAlloc[mapRoom]()
		rooms[count] = r
		if r == nil {
			return nil
		}
		r.Kind = 3
		r.Width = 4
		r.Height = 5
		mapRoomConnect(r, previous, 0)
		mapRoomConnect(previous, r, 1)
		r.Size = types.Ptf(float32(float64(r.Width)*32.526913), float32(float64(r.Height)*32.526913))
		mapRoomSetPos(r, &pos)
		mapRoomAdd(r)
		previous = r
		count++
		pos.Y = float32(float64(pos.Y) + 162.63457)
	}
	stride := int32(count / 5)
	if stride < 1 {
		stride = 1
	}
	for i := mapRoomRandomInt(0, stride); i < int32(count); i += mapRoomRandomInt(1, stride) {
		r := rooms[i]
		r.Reserved64[5] = *mapGrowthFrontier()
		*mapGrowthFrontier() = mapRoomRaw(unsafe.Pointer(r))
	}
	r := rooms[count-1]
	*populationGlobal(0) = mapRoomRaw(unsafe.Pointer(r))
	return r
}
