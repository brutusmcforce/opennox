package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func mapGrowthConfig() unsafe.Pointer { return memmap.PtrOff(0x5D4594, 1549796) }

func mapGrowthFill(r *mapRoom, depth, halls, width int32, origin *mapRoom) uint32 {
	var directions [4]int32
	any := false
	for i := range directions {
		if mapRoomHasRoomNeighbor(r, int32(i)) == 0 {
			directions[i] = int32(i) + 2
			any = true
		}
	}
	if !any {
		return 1
	}
	index := mapRoomRandomInt(0, 3)
	for attempts := 0; attempts < 8; {
		index = (index + 1) % 4
		kind := directions[index]
		if kind == 0 {
			continue
		}
		mean, spread := int32(*populationBlob(1549808)), int32(*populationBlob(1549812))
		w := mapRoomRandomInt(mean-spread, spread+mean)
		n := mapRoomPrepareHall(mapGrowthConfig(), kind, w)
		if n == nil {
			return 0
		}
		var p types.Pointf
		switch n.Kind {
		case 2:
			p.X = float32(mapRoomRandomX(r, n))
			p.Y = r.Pos.Y - n.Size.Y
		case 3:
			p.X = float32(mapRoomRandomX(r, n))
			p.Y = r.Size.Y + r.Pos.Y
		case 4:
			p.X = r.Size.X + r.Pos.X
			p.Y = float32(mapRoomRandomY(r, n))
		case 5:
			p.X = r.Pos.X - n.Size.X
			p.Y = float32(mapRoomRandomY(r, n))
		}
		mapRoomSetPos(n, &p)
		if mapRoomWithinBounds(mapGrowthConfig(), n) == 0 {
			mapRoomFree(n)
		} else {
			dir := mapRoomHallDirection(n.Kind)
			mapRoomConnect(n, r, mapRoomOpposite(dir))
			if other := mapRoomOverlap(n); other != nil {
				if mapRoomIsHall(other) != 0 || other.Flags&2 != 0 || other == origin || mapRoomRandomInt(1, 100) > int32(*populationBlob(1549844)) || mapRoomTrimHall(n, other) == 0 {
					mapRoomFree(n)
				} else {
					mapRoomAdd(n)
					mapRoomConnectBoth(n, other, mapRoomHallDirection(n.Kind))
					mapRoomConnect(r, n, mapRoomHallDirection(n.Kind))
				}
			} else {
				mapRoomAdd(n)
				if mapGrowthDispatch(n, depth, 1, w, r) != 0 {
					mapRoomConnect(r, n, mapRoomHallDirection(n.Kind))
				} else {
					mapRoomRemove(n)
					mapRoomFree(n)
				}
			}
		}
		attempts++
	}
	return 1
}

// Each accepted branch owns its candidate; rejected candidates are released only
// after their neighbor and geometry state has been established, as in the C path.
func mapGrowthAdmitHall(r, n *mapRoom, back, depth, halls, width int32, origin *mapRoom) bool {
	mapRoomConnect(n, r, back)
	if mapRoomWithinBounds(mapGrowthConfig(), n) == 0 {
		mapRoomFree(n)
		return false
	}
	other := mapRoomOverlap(n)
	if other != nil {
		if other.Kind != 1 || other.Flags&2 != 0 || other == origin || mapRoomRandomInt(1, 100) > int32(*populationBlob(1549844)) || mapRoomTrimHall(n, other) == 0 {
			mapRoomFree(n)
			return false
		}
	}
	mapRoomAdd(n)
	if other != nil {
		mapRoomConnectBoth(n, other, mapRoomHallDirection(n.Kind))
	} else if mapGrowthDispatch(n, depth, halls+1, width, origin) == 0 {
		mapRoomRemove(n)
		mapRoomFree(n)
		return false
	}
	mapRoomConnect(r, n, mapRoomOpposite(back))
	return true
}

func mapGrowthHall(r *mapRoom, depth, halls, width int32, origin *mapRoom) uint32 {
	choice := int32(1)
	if halls < int32(*populationBlob(1549816)) {
		choice = mapGrowthBranch(r)
	}
	if choice == 1 {
		for attempt := 0; attempt < 10; attempt++ {
			n := mapRoomPrepare(mapGrowthConfig())
			if n == nil {
				return 0
			}
			var p types.Pointf
			switch r.Kind {
			case 2:
				p.X = float32(mapRoomRandomReverseX(n, r))
				p.Y = r.Pos.Y - n.Size.Y
			case 3:
				p.X = float32(mapRoomRandomReverseX(n, r))
				p.Y = r.Size.Y + r.Pos.Y
			case 4:
				p.X = r.Size.X + r.Pos.X
				p.Y = float32(mapRoomRandomReverseY(n, r))
			case 5:
				p.X = r.Pos.X - n.Size.X
				p.Y = float32(mapRoomRandomReverseY(n, r))
			}
			mapRoomSetPos(n, &p)
			if mapRoomCanPlace(mapGrowthConfig(), n) == 0 {
				mapRoomFree(n)
				continue
			}
			mapRoomAdd(n)
			mapRoomConnectBoth(r, n, mapRoomHallDirection(r.Kind))
			mapGrowthDispatch(n, depth, 0, 0, n)
			return 1
		}
		return 0
	}
	left, right := false, false
	if choice == 2 || choice == 8 || choice == 32 || choice == 64 {
		kind := int32(*memmap.PtrUint32(0x587000, 197812+uintptr(uint32(r.Kind))*4))
		n := mapRoomPrepareHall(mapGrowthConfig(), kind, width)
		if n == nil {
			return 0
		}
		var p types.Pointf
		switch r.Kind {
		case 2:
			p = types.Ptf(r.Pos.X-n.Size.X, r.Pos.Y)
		case 3:
			p = types.Ptf(r.Size.X+r.Pos.X, float32(float64(r.Size.Y)+float64(r.Pos.Y)-float64(n.Size.Y)))
		case 4:
			p = types.Ptf(float32(float64(r.Size.X)+float64(r.Pos.X)-float64(n.Size.X)), r.Pos.Y-n.Size.Y)
		case 5:
			p = types.Ptf(r.Pos.X, r.Size.Y+r.Pos.Y)
		}
		mapRoomSetPos(n, &p)
		left = mapGrowthAdmitHall(r, n, mapRoomHallOppositeSide(r.Kind), depth, halls, width, origin)
		if !left && (choice == 2 || choice == 8) {
			return 0
		}
	}
	if choice == 4 || choice == 16 || choice == 32 || choice == 64 {
		kind := int32(*memmap.PtrUint32(0x587000, 197836+uintptr(uint32(r.Kind))*4))
		n := mapRoomPrepareHall(mapGrowthConfig(), kind, width)
		if n == nil {
			return 0
		}
		var p types.Pointf
		switch r.Kind {
		case 2:
			p = types.Ptf(r.Size.X+r.Pos.X, r.Pos.Y)
		case 3:
			p = types.Ptf(r.Pos.X-n.Size.X, float32(float64(r.Size.Y)+float64(r.Pos.Y)-float64(n.Size.Y)))
		case 4:
			p = types.Ptf(float32(float64(r.Size.X)+float64(r.Pos.X)-float64(n.Size.X)), r.Size.Y+r.Pos.Y)
		case 5:
			p = types.Ptf(r.Pos.X, r.Pos.Y-n.Size.Y)
		}
		mapRoomSetPos(n, &p)
		right = mapGrowthAdmitHall(r, n, mapRoomHallSide(r.Kind), depth, halls, width, origin)
		if !right && (choice == 4 || choice == 16) {
			return 0
		}
	}
	if !left && !right {
		return 0
	}
	if choice != 32 && choice != 8 && choice != 16 {
		return 1
	}
	n := mapRoomPrepareHall(mapGrowthConfig(), r.Kind, width)
	if n == nil {
		return 0
	}
	var p types.Pointf
	switch r.Kind {
	case 2:
		p = types.Ptf(r.Pos.X, r.Pos.Y-n.Size.Y)
	case 3:
		p = types.Ptf(r.Pos.X, r.Size.Y+r.Pos.Y)
	case 4:
		p = types.Ptf(r.Size.X+r.Pos.X, r.Pos.Y)
	case 5:
		p = types.Ptf(r.Pos.X-n.Size.X, r.Pos.Y)
	}
	mapRoomSetPos(n, &p)
	mapGrowthAdmitHall(r, n, mapRoomOpposite(mapRoomHallDirection(r.Kind)), depth, halls, width, origin)
	return 1
}
