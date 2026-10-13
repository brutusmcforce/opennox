package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
)

func mapHallCount() *uint32       { return populationBlob(2491608) }
func mapHallAt(i uint32) *mapRoom { return populationRoom(*populationBlob(2491612 + uintptr(i)*4)) }
func mapHallStore(r *mapRoom) {
	*populationBlob(2491612 + uintptr(*mapHallCount())*4) = mapRoomRaw(unsafe.Pointer(r))
}
func mapHallPos(r *mapRoom, x, y int32) {
	p := types.Ptf(float32(float64(x)*32.526913), float32(float64(y)*32.526913))
	mapRoomSetPos(r, &p)
}

func mapHallAdmit(a, b *mapRoom, ad, bd int32) uint32 {
	for i := int32(0); i < int32(*mapHallCount()); i++ {
		if mapRoomOverlap(mapHallAt(uint32(i))) != nil {
			for j := int32(0); j < int32(*mapHallCount()); j++ {
				mapRoomFree(mapHallAt(uint32(j)))
			}
			return 0
		}
	}
	for i := int32(0); i < int32(*mapHallCount()); i++ {
		mapRoomAdd(mapHallAt(uint32(i)))
	}
	mapRoomConnectBoth(a, mapHallAt(0), ad)
	mapRoomConnectBoth(b, mapHallAt(*mapHallCount()-1), bd)
	return 1
}

func mapHallVertical(a, b *mapRoom, p, q *[2]int32, w int32) uint32 {
	dx := q[0] - p[0]
	distance := dx
	if distance < 0 {
		distance = -distance
	}
	*mapHallCount() = 0
	if distance != 0 {
		if distance < 3 {
			return 0
		}
		mid := p[1] - (p[1]-q[1])/2 - 1
		r := mapRoomNewHall(2, w, p[1]-mid)
		mapHallStore(r)
		mapHallPos(r, p[0], mid)
		*mapHallCount()++
		if dx <= 0 {
			r = mapRoomNewHall(5, w, p[0]-q[0])
			mapHallStore(r)
			mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 3)
			mapHallPos(r, q[0], mid)
		} else {
			r = mapRoomNewHall(4, w, q[0]-p[0])
			mapHallStore(r)
			mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 2)
			mapHallPos(r, w+p[0], mid)
		}
		*mapHallCount()++
		r = mapRoomNewHall(2, w, mid-q[1]-1)
		mapHallStore(r)
		mapHallPos(r, q[0], q[1]+1)
		mapRoomConnectBoth(mapHallAt(1), mapHallAt(2), 0)
	} else {
		r := mapRoomNewHall(2, w, p[1]-q[1]-1)
		mapHallStore(r)
		mapHallPos(r, q[0], q[1]+1)
	}
	*mapHallCount()++
	return mapHallAdmit(a, b, 0, 1)
}

func mapHallHorizontal(a, b *mapRoom, p, q *[2]int32, w int32) uint32 {
	dy := q[1] - p[1]
	distance := dy
	if distance < 0 {
		distance = -distance
	}
	*mapHallCount() = 0
	if distance != 0 {
		if distance < 3 {
			return 0
		}
		mid := p[0] + (q[0]-p[0])/2
		r := mapRoomNewHall(4, w, mid-p[0]+w-1)
		mapHallStore(r)
		mapHallPos(r, p[0]+1, p[1])
		*mapHallCount()++
		if dy <= 0 {
			r = mapRoomNewHall(2, w, p[1]-q[1])
			mapHallStore(r)
			mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 0)
			mapHallPos(r, mid, q[1])
		} else {
			r = mapRoomNewHall(3, w, q[1]-p[1])
			mapHallStore(r)
			mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 1)
			mapHallPos(r, mid, p[1]+w)
		}
		*mapHallCount()++
		r = mapRoomNewHall(4, w, q[0]-mid-w)
		mapHallStore(r)
		mapHallPos(r, w+mid, q[1])
		mapRoomConnectBoth(mapHallAt(1), mapHallAt(2), 2)
	} else {
		r := mapRoomNewHall(4, w, q[0]-p[0]-1)
		mapHallStore(r)
		mapHallPos(r, p[0]+1, p[1])
	}
	*mapHallCount()++
	return mapHallAdmit(a, b, 2, 3)
}

func mapHallBendNorth(a, b *mapRoom, p, q *[2]int32, w int32) uint32 {
	if q[1] > p[1]-w {
		return 0
	}
	*mapHallCount() = 0
	dx := q[0] - p[0]
	r := mapRoomNewHall(2, w, p[1]-q[1])
	mapHallStore(r)
	mapHallPos(r, p[0], q[1])
	*mapHallCount()++
	var bd int32
	if dx <= 0 {
		r = mapRoomNewHall(5, w, p[0]-q[0]-1)
		mapHallStore(r)
		mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 3)
		mapHallPos(r, q[0]+1, q[1])
		bd = 2
	} else {
		r = mapRoomNewHall(4, w, q[0]-p[0]-w)
		mapHallStore(r)
		mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 2)
		mapHallPos(r, w+p[0], q[1])
		bd = 3
	}
	*mapHallCount()++
	return mapHallAdmit(a, b, 0, bd)
}

func mapHallBendSouth(a, b *mapRoom, p, q *[2]int32, w int32) uint32 {
	if q[1] <= p[1] {
		return 0
	}
	dx := q[0] - p[0]
	*mapHallCount() = 0
	r := mapRoomNewHall(3, w, q[1]-p[1]+w-1)
	mapHallStore(r)
	mapHallPos(r, p[0], p[1]+1)
	*mapHallCount()++
	var bd int32
	if dx <= 0 {
		r = mapRoomNewHall(5, w, p[0]-q[0]-1)
		mapHallStore(r)
		mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 3)
		mapHallPos(r, q[0]+1, q[1])
		bd = 2
	} else {
		r = mapRoomNewHall(4, w, q[0]-p[0]-w)
		mapHallStore(r)
		mapRoomConnectBoth(mapHallAt(0), mapHallAt(1), 2)
		mapHallPos(r, w+p[0], q[1])
		bd = 3
	}
	*mapHallCount()++
	return mapHallAdmit(a, b, 1, bd)
}

func mapHallConnect(prefab uint32, dir int32, candidate *mapRoom) uint32 {
	marker := (*[2]int32)(unsafe.Pointer(populationWord(prefab, 80+int(dir)*16)))
	w := int32(*populationWord(prefab, 88+int(dir)*16))
	room := populationRoom(*populationWord(prefab, 148))
	var q [2]int32
	switch dir {
	case 0, 1:
		q[0] = candidate.Grid[0] + candidate.Width/2
		q[1] = candidate.Grid[1]
		if dir == 0 {
			q[1] += candidate.Height - 1
		}
		limit := candidate.Width - w
		for i := int32(0); i <= limit; i++ {
			var ok uint32
			if dir == 0 {
				ok = mapHallVertical(room, candidate, marker, &q, w)
			} else {
				ok = mapHallVertical(candidate, room, &q, marker, w)
			}
			if ok != 0 {
				return 1
			}
			q[0]++
			if q[0] > candidate.Grid[0]+limit {
				q[0] = candidate.Grid[0]
			}
		}
		q[1] = candidate.Grid[1] + candidate.Height/2
		limit = candidate.Height - w
		for i := int32(0); i <= limit; i++ {
			q[0] = candidate.Grid[0]
			var ok uint32
			if dir == 0 {
				ok = mapHallBendNorth(room, candidate, marker, &q, w)
			} else {
				ok = mapHallBendSouth(room, candidate, marker, &q, w)
			}
			if ok != 0 {
				return 1
			}
			q[0] += candidate.Width - 1
			if dir == 0 {
				ok = mapHallBendNorth(room, candidate, marker, &q, w)
			} else {
				ok = mapHallBendSouth(room, candidate, marker, &q, w)
			}
			if ok != 0 {
				return 1
			}
			q[1]++
			if q[1] > candidate.Grid[1]+limit {
				q[1] = candidate.Grid[1]
			}
		}
	case 2, 3:
		q[0] = candidate.Grid[0]
		if dir == 3 {
			q[0] += candidate.Width - 1
		}
		q[1] = candidate.Grid[1] + candidate.Height/2
		// Original horizontal search uses Width for the vertical retry bound.
		limit := candidate.Width - w
		for i := int32(0); i <= limit; i++ {
			var ok uint32
			if dir == 2 {
				ok = mapHallHorizontal(room, candidate, marker, &q, w)
			} else {
				ok = mapHallHorizontal(candidate, room, &q, marker, w)
			}
			if ok != 0 {
				return 1
			}
			q[1]++
			if q[1] > candidate.Grid[1]+limit {
				q[1] = candidate.Grid[1]
			}
		}
		q[0] = candidate.Grid[0] + candidate.Width/2
		for i := int32(0); i <= limit; i++ {
			q[1] = candidate.Grid[1]
			if mapHallBendNorth(candidate, room, &q, marker, w) != 0 {
				return 1
			}
			q[1] += candidate.Height - 1
			if mapHallBendSouth(candidate, room, &q, marker, w) != 0 {
				return 1
			}
			q[0]++
			if q[0] > candidate.Grid[0]+limit {
				q[0] = candidate.Grid[0]
			}
		}
	}
	return 0
}
