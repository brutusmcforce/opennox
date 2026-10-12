package legacy

import (
	"math"
	"unsafe"
)

func mapPopulationDistance(room, previous uint32, distance float32) {
	mapPopulationProgress(156)
	extra := float64(0)
	if previous != 0 {
		p := populationRoom(previous)
		switch p.Kind {
		case 1:
			extra = math.Sqrt(float64(p.Size.X)*float64(p.Size.X) + float64(p.Size.Y)*float64(p.Size.Y))
		case 2, 3:
			extra = float64(p.Size.Y)
		case 4, 5:
			extra = float64(p.Size.X)
		default:
			extra = float64(distance)
		}
	}
	next := float32(extra + float64(distance))
	r := populationRoom(room)
	marker := (*byte)(unsafe.Pointer(&r.Reserved220))
	if *marker != 1 || next < *populationFloat(room, 356) {
		*populationFloat(room, 356) = next
		*marker = 1
		for dir := 0; dir < 4; dir++ {
			for j := 0; j < int(r.Counts[dir]); j++ {
				if n := r.Neighbors[dir][j]; n != nil {
					mapPopulationDistance(mapRoomRaw(unsafe.Pointer(n)), room, next)
				}
			}
		}
	}
}
func mapPopulationFindFarthest(room uint32) {
	r := populationRoom(room)
	marker := (*byte)(unsafe.Pointer(&r.Reserved220))
	if *marker == 2 {
		return
	}
	if mapPopulationDistanceMax() < float64(*populationFloat(room, 356)) && r.Kind == 1 {
		*populationGlobal(3) = room
		*populationGlobal(4) = *populationWord(room, 356)
	}
	*marker = 2
	for dir := 0; dir < 4; dir++ {
		for j := 0; j < int(r.Counts[dir]); j++ {
			if n := r.Neighbors[dir][j]; n != nil {
				mapPopulationFindFarthest(mapRoomRaw(unsafe.Pointer(n)))
			}
		}
	}
}
func mapPopulationSort() uint32 {
	*populationGlobal(5) = 0
	for r := mapRoomHead(); r != nil; r = r.Next {
		raw := mapRoomRaw(unsafe.Pointer(r))
		cur := *populationGlobal(5)
		var prev uint32
		for cur != 0 && !(*populationFloat(raw, 356) < *populationFloat(cur, 356)) {
			prev = cur
			cur = *populationWord(cur, 64)
		}
		*populationWord(raw, 64) = cur
		*populationWord(raw, 68) = prev
		if prev == 0 {
			*populationGlobal(5) = raw
		} else {
			*populationWord(prev, 64) = raw
		}
		if cur != 0 {
			*populationWord(cur, 68) = raw
		}
	}
	return 0
}
func mapPopulationThemes(start uint32) uint32 {
	*populationGlobal(4) = 0
	*populationGlobal(3) = start
	mapPopulationFindFarthest(start)
	populationRoom(*populationGlobal(3)).Flags |= 1
	var count int32
	for r := mapRoomHead(); r != nil; r = r.Next {
		count++
		raw := mapRoomRaw(unsafe.Pointer(r))
		*populationFloat(raw, 360) = float32(float64(*populationFloat(raw, 356)) / mapPopulationDistanceMax())
	}
	*populationGlobal(5) = 0
	mapPopulationSort()
	for p, rank := start, int32(0); p != 0; p, rank = *populationWord(p, 64), rank+100 {
		pct := int32(100)
		if count > 1 {
			pct = rank / (count - 1)
		}
		flag := uint32(1)
		if p != start {
			switch {
			case pct < 30:
				flag = 2
			case pct < 60:
				flag = 4
			case pct < 90:
				flag = 8
			default:
				flag = 16
			}
		}
		*populationWord(p, 364) = flag
	}
	*populationWord(*populationGlobal(3), 364) = 32
	return *populationGlobal(3)
}
