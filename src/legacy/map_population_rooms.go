package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"unsafe"
)

func mapPopulationRoom(cfg, room uint32) {
	for group := *populationWord(mapRoomRaw(populationRoom(room).Decoration), 92); group != 0; group = *populationWord(group, 20) {
		mapPopulationGroup(cfg, room, group)
	}
}
func mapPopulationGroup(cfg, room, group uint32) {
	head := *populationWord(group, 8)
	if head == 0 {
		return
	}
	count := mapRoomRandomInt(int32(*populationWord(group, 0)), int32(*populationWord(group, 4)))
	total := int32(*populationWord(group, 12))
	if count > total {
		count = total
	}
	tail := head
	*populationWord(head, 92) = 0
	*populationWord(head, 96) = 0
	for node, i := *populationWord(head, 88), int32(1); i < total; i++ {
		if mapRoomRandomInt(1, 100) >= 50 {
			if mapRoomRandomInt(1, 100) >= 50 {
				*populationWord(tail, 92) = node
				*populationWord(node, 96) = tail
				*populationWord(node, 92) = 0
				tail = node
			} else {
				prev := *populationWord(tail, 96)
				*populationWord(node, 92) = tail
				*populationWord(node, 96) = prev
				if prev != 0 {
					*populationWord(prev, 92) = node
				} else {
					head = node
				}
				*populationWord(tail, 96) = node
			}
		} else if mapRoomRandomInt(1, 100) >= 50 {
			next := *populationWord(head, 92)
			if next != 0 {
				*populationWord(next, 96) = node
			} else {
				tail = node
			}
			*populationWord(node, 96) = head
			*populationWord(node, 92) = next
			*populationWord(head, 92) = node
		} else {
			*populationWord(node, 92) = head
			*populationWord(node, 96) = 0
			*populationWord(head, 96) = node
			head = node
		}
		node = *populationWord(node, 88)
	}
	for ; count != 0; count-- {
		mapPopulationSpawn(cfg, room, head)
		head = *populationWord(head, 92)
	}
}
func mapPopulationSpawn(cfg, room, row uint32) uint32 {
	r := populationRoom(room)
	var count int32
	if *populationWord(row, 64) != 0 {
		count = int32(int64(float64(r.Width*r.Height) * float64(*populationFloat(row, 76))))
		if lo := int32(*populationWord(row, 68)); count < lo {
			count = lo
		} else if hi := int32(*populationWord(row, 72)); count > hi {
			count = hi
		}
	} else {
		count = mapRoomRandomInt(int32(*populationWord(row, 68)), int32(*populationWord(row, 72)))
	}
	ret := *populationWord(row, 0)
	switch ret {
	case 0:
		for ; count > 0; count-- {
			u := mapPopulationMonster(room, row+4)
			if u != 0 {
				mapPopulationInventory(cfg, u, *populationWord(row, 80))
			}
			ret = uint32(count - 1)
		}
	case 1:
		index := mapPopulationMetadataIndex(row + 4)
		ret = index
		if int32(index) >= 0 {
			for ; count > 0; count-- {
				for attempt := 0; attempt < 3; attempt++ {
					if mapPopulationPlacePrefabInRoom(cfg, room, row, index) != 0 {
						break
					}
				}
				ret = uint32(count - 1)
			}
		}
	case 2:
		ret = mapRoomRaw(unsafe.Pointer(mapRoomRemoveGeneratedExclusions(r)))
	case 3, 4, 5:
		var pos types.Pointf
		ret = mapRoomRandomPoint(r, float32(0.89999998), &pos)
		if ret != 0 {
			switch *populationWord(row, 0) {
			case 3:
				ret = mapPopulationItem(row+4, *populationWord(cfg, 1100), *populationWord(cfg, 1104))
			case 4:
				ret = mapPopulationItem(row+4, *populationWord(cfg, 1108), *populationWord(cfg, 1112))
			case 5:
				ret = mapPopulationSpellbook(cfg, row+4)
			}
			if ret != 0 {
				ret = mapRoomRaw(unsafe.Pointer(mapPaintMoveObject(populationObject(ret), &pos)))
			}
		}
	}
	return ret
}
func mapPopulationMonster(room, name uint32) uint32 {
	r := populationRoom(room)
	center := types.Pointf{X: float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5), Y: float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)}
	var pos types.Pointf
	if mapRoomRandomPoint(r, float32(0.94999999), &pos) == 0 {
		return 0
	}
	mapPaintSelectObject((*byte)(mapRoomPointer(name)))
	u := mapPaintPlaceObject(&pos)
	// The original expression converts the class word interpreted as a float.
	if u != nil && byte(u.ObjClass)&2 != 0 {
		delta := types.Pointf{X: center.X - pos.X, Y: center.Y - pos.Y}
		angle := geometryVectorAngle(&delta)
		mapPaintOrientObject(u, geometryDirection4Index(angle))
	}
	return mapRoomRaw(unsafe.Pointer(u))
}
func mapPopulationWaypoint(point uint32) uint32 {
	return populationEnsureWaypoint(populationPoint(point))
}
func populationEnsureWaypoint(p *types.Pointf) uint32 {
	r := prefabFindWaypoint(p)
	if r != 0 {
		return r
	}
	return prefabCreateWaypoint(p)
}
func populationWaypointConnect(a, b *types.Pointf) {
	prefabConnectWaypoint(a, b)
}
func mapPopulationHallwayWaypoints(cfg uint32) uint32 {
	prefabSetWaypointKind(128)
	for r := mapRoomHead(); r != nil; r = r.Next {
		mapPopulationProgress(156)
		if r.Kind == 1 {
			for dir := 0; dir < 4; dir++ {
				for j := 0; j < int(r.Counts[dir]); j++ {
					n := r.Neighbors[dir][j]
					if mapRoomIsHall(n) != 0 && mapRoomHallDirection(n.Kind) == int32(dir) {
						var a, b types.Pointf
						mapRoomHallEntry(n, &a)
						mapRoomRememberPoint(r, &a)
						mapRoomHallCenter(n, &b)
						populationEnsureWaypoint(&a)
						populationEnsureWaypoint(&b)
						populationWaypointConnect(&a, &b)
						if *populationWord(cfg, 60) == 0 {
							populationWaypointConnect(&b, &a)
						}
					}
				}
			}
		} else {
			var a types.Pointf
			mapRoomHallCenter(r, &a)
			for dir := 0; dir < 4; dir++ {
				if mapRoomIsEntranceSide(r, int32(dir)) != 0 {
					continue
				}
				for j := 0; j < int(r.Counts[dir]); j++ {
					n := r.Neighbors[dir][j]
					var b types.Pointf
					if mapRoomIsHall(n) != 0 {
						mapRoomHallCenter(n, &b)
					} else {
						mapRoomHallExit(r, &b)
						mapRoomRememberPoint(n, &b)
					}
					populationEnsureWaypoint(&a)
					populationEnsureWaypoint(&b)
					populationWaypointConnect(&a, &b)
					if *populationWord(cfg, 60) == 0 {
						populationWaypointConnect(&b, &a)
					}
				}
			}
		}
	}
	return 0
}
func mapPopulationRoomWaypoints(cfg uint32) uint32 {
	if *populationWord(cfg, 60) != 0 {
		return cfg
	}
	for r := mapRoomHead(); r != nil; r = r.Next {
		if r.Kind != 1 {
			continue
		}
		mapPopulationProgress(156)
		for i := 0; i < int(r.PointCount); i++ {
			for j := 0; j < int(r.PointCount); j++ {
				if i != j {
					populationWaypointConnect(&r.Points[i], &r.Points[j])
				}
			}
		}
	}
	return 0
}
