package legacy

import (
	"math"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
)

func mapPopulationMetadataIndex(name uint32) uint32 {
	for i := int32(0); i < int32(*populationGlobal(13)); i++ {
		if populationString(*populationGlobal(12)+uint32(i)*64) == populationString(name) {
			return uint32(i)
		}
	}
	return math.MaxUint32
}
func mapPopulationMetadataAt(index uint32) uint32 { return *populationGlobal(12) + index*64 }
func mapPopulationMetadataFree()                  { mapRoomRelease(mapRoomPointer(*populationGlobal(12))) }
func mapPopulationMetadataInit() uint32 {
	n := *prefabGlobal(prefabCount)
	*populationGlobal(13) = n
	p := mapRoomRaw(mapRoomCalloc(n, 64))
	*populationGlobal(12) = p
	if p == 0 {
		*populationGlobal(13) = 0
		return 0
	}
	ret := p
	for i := int32(0); i < int32(n); i++ {
		name := populationString(prefabMetadataAt(i))
		row := p + uint32(i)*64
		copy(unsafe.Slice((*byte)(mapRoomPointer(row)), len(name)+1), append([]byte(name), 0))
		fields := strings.FieldsFunc(name, func(r rune) bool { return r == '-' })
		if len(fields) >= 3 {
			for j := 0; ; j++ {
				c := *memmap.PtrUint8(0x587000, 255032+uintptr(j))
				if c == 0 {
					break
				}
				if strings.IndexByte(fields[2], c) >= 0 {
					*populationWord(row, 60) |= 1 << uint(j)
				}
			}
		}
		ret = uint32(i + 1)
	}
	return ret
}
func populationMarkerType(dir int) uint32 {
	if dir == 0 {
		return *populationGlobal(11)
	}
	return *populationBlob(uintptr(2487656 + dir*4))
}
func mapPopulationFillCachedObjects(cfg, mapping uint32) uint32 {
	for u := prefabObjectHead(); u != 0; u = prefabObjectNext(u) {
		for m := mapping; m != 0; m = *populationWord(m, 8) {
			if uint32(populationObject(u).TypeInd) == *populationWord(m, 0) {
				mapPopulationInventory(cfg, u, *populationWord(m, 4))
				break
			}
		}
	}
	return 0
}
func mapPopulationPlacePrefabInRoom(cfg, room, row, index uint32) uint32 {
	e := mapRoomAlloc[mapRoomExclusion]()
	if e == nil {
		return 0
	}
	width := float32(prefabDimension(int32(index), 64))
	height := float32(prefabDimension(int32(index), 68))
	w, h := int32(int64(float64(width)*0.030743772)), int32(int64(float64(height)*0.030743772))
	r := populationRoom(room)
	dx, dy := r.Width-w, r.Height-h
	if dx >= 0 && dy >= 0 {
		flags := *populationWord(mapPopulationMetadataAt(index), 60)
		pos := types.Pointf{X: float32(float64(mapRoomRandomInt(0, dx))*32.526913 + float64(r.Min.X)), Y: float32(float64(mapRoomRandomInt(0, dy))*32.526913 + float64(r.Min.Y))}
		if flags&1 != 0 {
			pos.Y = r.Min.Y
		} else if flags&2 != 0 {
			pos.Y = r.Max.Y - height
		}
		if flags&4 != 0 {
			pos.X = r.Max.X - width
		} else if flags&8 != 0 {
			pos.X = r.Min.X
		}
		if flags&16 != 0 {
			pos.X = float32(float64(dx/2)*32.526913 + float64(r.Min.X))
			pos.Y = float32(float64(dy/2)*32.526913 + float64(r.Min.Y))
		}
		e.Kind = 1
		e.Min = pos
		e.Max = types.Pointf{X: pos.X + width, Y: pos.Y + height}
		e.Prefab = index
		if mapRoomContainsRect(r, e) != 0 && mapRoomFindExclusionOverlap(r, e) == nil {
			populationLoadPrefab(int32(index))
			mapPopulationFillCachedObjects(cfg, *populationWord(row, 84))
			prefabInstantiate(&pos)
			e.Next = r.Exclusions
			r.Exclusions = e
			return 1
		}
	}
	mapRoomRelease(unsafe.Pointer(e))
	return 0
}
func mapPopulationPrefabInfo(cfg, prefab uint32) uint32 {
	index := prefabFindName(prefab)
	*populationWord(prefab, 68) = index
	if int32(index) == -1 {
		return 0
	}
	populationLoadPrefab(int32(index))
	*populationFloat(prefab, 60) = float32(prefabDimension(int32(index), 64))
	*populationFloat(prefab, 64) = float32(prefabDimension(int32(index), 68))
	for u := prefabObjectHead(); u != 0; u = prefabObjectNext(u) {
		var pos types.Pointf
		prefabRelativePosition(u, &pos)
		pos.X -= 1
		pos.Y -= 1
		var grid [2]int32
		mapRoomRound(&pos, &grid)
		typ := uint32(populationObject(u).TypeInd)
		for _, dir := range []int{0, 1, 3, 2} {
			if typ != populationMarkerType(dir) {
				continue
			}
			e := prefab + uint32(80+dir*16)
			if *populationWord(e, 12) == 0 {
				*populationWord(e, 0) = uint32(grid[0])
				*populationWord(e, 4) = uint32(grid[1])
				*populationWord(e, 8) = 1
				*populationWord(e, 12) = 1
			} else {
				axis := 0
				if dir >= 2 {
					axis = 1
				}
				if grid[axis] < int32(*populationWord(e, axis*4)) {
					*populationWord(e, axis*4) = uint32(grid[axis])
				}
				*populationWord(e, (1-axis)*4) = uint32(grid[1-axis])
				*populationWord(e, 8)++
			}
			break
		}
	}
	exits := 0
	for dir := 0; dir < 4; dir++ {
		e := prefab + uint32(80+dir*16)
		if *populationWord(e, 12) != 0 {
			if *populationWord(e, 8) > *populationWord(prefab, 144) {
				*populationWord(prefab, 144) = *populationWord(e, 8)
			}
			exits++
		}
	}
	if *populationWord(prefab, 144) > *populationWord(cfg, 12)+*populationWord(cfg, 16) || exits > 2 || ((*populationWord(prefab, 92) != 0 || *populationWord(prefab, 108) != 0) && (*populationWord(prefab, 124) != 0 || *populationWord(prefab, 140) != 0)) {
		return 0
	}
	return 1
}
func mapPopulationPrefabRoom(cfg, prefab uint32) uint32 {
	if mapPopulationPrefabInfo(cfg, prefab) == 0 {
		return 0
	}
	padding := *populationWord(prefab, 144)
	r := mapRoomNew(int32(4*padding+uint32(int64(float64(*populationFloat(prefab, 60))*0.030743772+0.5))), int32(4*padding+uint32(int64(float64(*populationFloat(prefab, 64))*0.030743772+0.5))))
	raw := mapRoomRaw(unsafe.Pointer(r))
	*populationWord(prefab, 148) = raw
	var pos types.Pointf
	populationNextPosition(r, &pos)
	mapRoomSetPos(r, &pos)
	valid := true
	if other := mapRoomOverlap(r); other != nil && mapRoomResolveOverlap(r, other) == 0 {
		valid = false
	}
	if mapRoomWithinBounds(mapRoomPointer(cfg), r) != 0 && valid {
		mapRoomAdd(r)
		r.Flags |= 2
		for dir := 0; dir < 4; dir++ {
			e := prefab + uint32(80+dir*16)
			if *populationWord(e, 12) != 0 {
				*populationWord(e, 0) += 2*padding + uint32(r.Grid[0])
				*populationWord(e, 4) += 2*padding + uint32(r.Grid[1])
			}
		}
		return 1
	}
	mapRoomFree(r)
	*populationWord(prefab, 148) = 0
	return 0
}
func mapPopulationSelectPrefabs(cfg uint32) uint32 {
	if *populationGlobal(11) == 0 {
		for dir, name := range []string{"ExitNorthMarker", "ExitSouthMarker", "ExitEastMarker", "ExitWestMarker"} {
			typ := uint32(GetServer().S().Types.IndByID(name))
			if dir == 0 {
				*populationGlobal(11) = typ
			} else {
				*populationBlob(uintptr(2487656 + 4*dir)) = typ
			}
		}
	}
	mapPopulationPrefabPositions(math.Float32frombits(cfg))
	count := 0
	for _, required := range []bool{true, false} {
		for p := *populationWord(cfg, 80); p != 0; p = *populationWord(p, 156) {
			if (*populationWord(p, 72) != 0) != required {
				continue
			}
			if count == 5 {
				if required {
					return 0
				}
				return 1
			}
			if mapPopulationPrefabRoom(cfg, p) == 0 {
				return 0
			}
			*populationWord(p, 76) = 1
			count++
		}
	}
	return 1
}
func mapPopulationCandidates(prefab, direction uint32) uint32 {
	e := prefab + 80 + direction*16
	px, py := int32(*populationWord(e, 0)), int32(*populationWord(e, 4))
	var rooms [6]*mapRoom
	var distances [6]int32
	count, worst := 0, 0
	for r := mapRoomHead(); r != nil; r = r.Next {
		if mapRoomIsHall(r) != 0 {
			continue
		}
		dx, dy := r.Grid[0]+r.Width/2-px, r.Grid[1]+r.Height/2-py
		if direction == 0 && dy >= 0 || direction == 1 && dy <= 0 || direction == 2 && dx <= 0 || direction == 3 && dx >= 0 {
			continue
		}
		d := dx*dx + dy*dy
		if count == 6 {
			if distances[worst] > d {
				rooms[worst] = r
				distances[worst] = d
			}
		} else {
			rooms[count] = r
			distances[count] = d
			count++
		}
		worst = 0
		for i := 1; i < count; i++ {
			if distances[i] > distances[worst] {
				worst = i
			}
		}
	}
	for i := 0; i < count-1; {
		if distances[i] > distances[i+1] {
			distances[i], distances[i+1] = distances[i+1], distances[i]
			rooms[i], rooms[i+1] = rooms[i+1], rooms[i]
			i = 0
		} else {
			i++
		}
	}
	for i := 0; i < count; i++ {
		p := mapRoomRaw(unsafe.Pointer(rooms[i]))
		*populationWord(p, 72) = 0
		if i+1 < count {
			*populationWord(p, 72) = mapRoomRaw(unsafe.Pointer(rooms[i+1]))
		}
	}
	if count == 0 {
		return 0
	}
	return mapRoomRaw(unsafe.Pointer(rooms[0]))
}
func mapPopulationConnectPrefabs(cfg uint32) uint32 {
	for p := *populationWord(cfg, 80); p != 0; p = *populationWord(p, 156) {
		if *populationWord(p, 76) == 0 {
			continue
		}
		r := populationRoom(*populationWord(p, 148))
		pos := r.Pos
		pos.X = float32(float64(int32(*populationWord(p, 144)))*65.053825 + float64(pos.X))
		pos.Y = float32(float64(int32(*populationWord(p, 144)))*65.053825 + float64(pos.Y))
		mapRoomRemove(r)
		mapRoomFree(r)
		r = mapRoomNew(int32(int64(float64(*populationFloat(p, 60))*0.030743772+0.5)), int32(int64(float64(*populationFloat(p, 64))*0.030743772+0.5)))
		*populationWord(p, 148) = mapRoomRaw(unsafe.Pointer(r))
		mapRoomSetPos(r, &pos)
		mapRoomAdd(r)
		r.Flags |= 2
		for dir := uint32(0); dir < 4; dir++ {
			if *populationWord(p, 92+int(16*dir)) == 0 {
				continue
			}
			candidate := mapPopulationCandidates(p, dir)
			for candidate != 0 && mapHallConnect(p, int32(dir), populationRoom(candidate)) == 0 {
				candidate = *populationWord(candidate, 72)
			}
			if candidate == 0 {
				return 0
			}
		}
	}
	return 1
}
func mapPopulationApplyPrefabs(cfg uint32) uint32 {
	for p := *populationWord(cfg, 80); p != 0; p = *populationWord(p, 156) {
		if *populationWord(p, 76) == 0 {
			continue
		}
		populationLoadPrefab(int32(*populationWord(p, 68)))
		for u := prefabObjectHead(); u != 0; {
			next := prefabObjectNext(u)
			typ := uint32(populationObject(u).TypeInd)
			for dir := 0; dir < 4; dir++ {
				if typ == populationMarkerType(dir) {
					prefabObjectRemove(u)
					break
				}
			}
			u = next
		}
		mapPopulationFillCachedObjects(cfg, *populationWord(p, 152))
		r := populationRoom(*populationWord(p, 148))
		prefabInstantiate(&r.Pos)
	}
	return 1
}
