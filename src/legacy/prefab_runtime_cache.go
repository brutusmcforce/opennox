package legacy

import (
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func prefabCacheNode(slot, size, nextOff, prevOff int) uint32 {
	p := mapRoomRaw(mapRoomCalloc(1, uintptr(size)))
	if p == 0 {
		return 0
	}
	head := (*prefabGlobal(slot))
	*prefabWord(p, nextOff) = head
	if head != 0 {
		*prefabWord(head, prevOff) = p
	}
	(*prefabGlobal(slot)) = p
	return p
}
func prefabTileNew(x, y int32, orientation uint32) uint32 {
	node := mapRoomRaw(mapRoomCalloc(1, 24))
	if node == 0 {
		return 0
	}
	data := mapRoomRaw(mapRoomCalloc(1, 20))
	if data == 0 {
		mapRoomRelease(mapRoomPointer(node))
		return 0
	}
	*prefabWord(node, 0) = data
	*prefabWord(node, 16) = (*prefabGlobal(prefabTiles))
	if head := (*prefabGlobal(prefabTiles)); head != 0 {
		*prefabWord(head, 20) = node
	}
	(*prefabGlobal(prefabTiles)) = node
	*memmap.PtrUint32(0x5D4594, 1599560)++
	px, py := float64(x)*46, float32(float64(y)*46)
	*prefabWord(node, 12) = uint32(byte(orientation))
	*prefabFloat(node, 4) = float32(px)
	*prefabFloat(node, 8) = py
	if byte(orientation) == 1 {
		*prefabFloat(node, 4) = float32(px + 23)
	} else {
		*prefabFloat(node, 8) = py + 23
	}
	return node
}
func prefabWallNew(x, y byte) uint32 {
	node := prefabCacheNode(prefabWalls, 12, 4, 8)
	if node == 0 {
		return 0
	}
	data := mapRoomRaw(mapRoomCalloc(1, unsafe.Sizeof(server.Wall{}))) // used as *server.Wall
	*prefabWord(node, 0) = data
	if data == 0 {
		head := *prefabWord(node, 4)
		*prefabGlobal(prefabWalls) = head
		if head != 0 {
			*prefabWord(head, 8) = 0
		}
		mapRoomRelease(mapRoomPointer(node))
		return 0
	}
	*(*byte)(unsafe.Add(mapRoomPointer(data), 5)) = x
	*(*byte)(unsafe.Add(mapRoomPointer(data), 6)) = y
	return node
}
func prefabWallFind(x, y int32) uint32 {
	for n := (*prefabGlobal(prefabWalls)); n != 0; n = *prefabWord(n, 4) {
		data := mapRoomPointer(*prefabWord(n, 0))
		if int32(*(*byte)(unsafe.Add(data, 5))) == x && int32(*(*byte)(unsafe.Add(data, 6))) == y {
			return n
		}
	}
	return 0
}
func prefabWaypointNew(index uint32, x, y float32) uint32 {
	node := prefabCacheNode(prefabWaypoints, 12, 4, 8)
	if node == 0 {
		return 0
	}
	wp, _ := alloc.New(server.Waypoint{})
	wp.Index = index
	wp.PosVec = types.Ptf(x, y)
	wp.Flags = 0x1000000
	*prefabWord(node, 0) = mapRoomRaw(unsafe.Pointer(wp))
	if head := *prefabWord(node, 4); head != 0 {
		wp.WpNext = prefabWaypoint(*prefabWord(head, 0))
		wp.WpNext.WpPrev = wp
	}
	return node
}
func prefabObjectNew(raw uint32) uint32 {
	node := prefabCacheNode(prefabObjects, 12, 4, 8)
	if node == 0 {
		return 0
	}
	*prefabWord(node, 0) = raw
	obj := populationObject(raw)
	obj.ObjPrev = nil
	obj.ObjNext = nil
	if head := *prefabWord(node, 4); head != 0 {
		obj.ObjNext = populationObject(*prefabWord(head, 0))
		obj.ObjNext.ObjPrev = obj
	}
	return node
}
func prefabObjectHead() uint32 {
	if (*prefabGlobal(prefabLoaded)) != (*prefabGlobal(prefabSelected)) || (*prefabGlobal(prefabLoaded)) == math.MaxUint32 || (*prefabGlobal(prefabPlaced)) == 1 {
		if prefabLoad(int32((*prefabGlobal(prefabSelected)))) == 0 {
			return 0
		}
	}
	if head := (*prefabGlobal(prefabObjects)); head != 0 {
		return *prefabWord(head, 0)
	}
	return 0
}
func prefabObjectNext(raw uint32) uint32 {
	if raw == 0 {
		return 0
	}
	if next := populationObject(raw).ObjNext; next != nil {
		return mapRoomRaw(unsafe.Pointer(next))
	}
	return 0
}
func prefabNodeNext(raw uint32) uint32 {
	if raw == 0 {
		return 0
	}
	return *prefabWord(raw, 4)
}
func prefabObjectRemove(raw uint32) uint32 {
	if raw == 0 {
		return 0
	}
	for n := (*prefabGlobal(prefabObjects)); n != 0; n = *prefabWord(n, 4) {
		if *prefabWord(n, 0) != raw {
			continue
		}
		next, prev := *prefabWord(n, 4), *prefabWord(n, 8)
		if next != 0 {
			*prefabWord(next, 8) = prev
		}
		if prev != 0 {
			*prefabWord(prev, 4) = next
		}
		if (*prefabGlobal(prefabObjects)) == n {
			(*prefabGlobal(prefabObjects)) = next
		}
		u := populationObject(raw)
		if u.ObjNext != nil {
			u.ObjNext.ObjPrev = u.ObjPrev
		}
		if u.ObjPrev != nil {
			u.ObjPrev.ObjNext = u.ObjNext
		}
		GetServer().S().Objs.FreeObject(u)
		mapRoomRelease(mapRoomPointer(n))
		return 1
	}
	return 0
}
func prefabPlaceWaypoints(dx, dy int32) uint32 {
	x, y := float32(dx), float32(dy)
	for n := (*prefabGlobal(prefabWaypoints)); n != 0; n = *prefabWord(n, 4) {
		wp := prefabWaypoint(*prefabWord(n, 0))
		wp.PosVec.X += x
		wp.PosVec.Y += y
		GetServer().S().Sub_579E90(wp)
	}
	return 1
}
func prefabPlaceObjects(dx, dy int32) uint32 {
	x, y := float32(dx), float32(dy)
	for n := (*prefabGlobal(prefabObjects)); n != 0; n = *prefabWord(n, 4) {
		u := populationObject(*prefabWord(n, 0))
		u.PosVec.X += x
		u.PosVec.Y += y
		GetServer().CreateObjectAt(u, nil, u.PosVec)
		u.ObjFlags |= object.Flags(0x80000000)
	}
	return 1
}
func prefabPlaceTiles(dx, dy int32) uint32 {
	transparent := memmap.PtrUint32(0x587000, 229704)
	if *transparent == math.MaxUint32 {
		for i, p := range tileDefinitionsAll() {
			if i >= prefabTileDefinitionCount() {
				break
			}
			if alloc.GoString(&p.NameBuf[0]) == "TransparentFloor" {
				*transparent = uint32(i)
				break
			}
		}
	}
	selection := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 35912)), 72)
	saved := append([]byte(nil), selection...)
	defer copy(selection, saved)
	x, y := float32(dx), float32(dy)
	for n := (*prefabGlobal(prefabTiles)); n != 0; n = *prefabWord(n, 16) {
		pos := types.Ptf(x+*prefabFloat(n, 4), y+*prefabFloat(n, 8))
		data := *prefabWord(n, 0)
		tile := *prefabWord(data, 0)
		selectTileImage(int32(tile))
		selectTileVariation(int32(*prefabWord(data, 4)))
		setTileFlag(1)
		if tile != *transparent {
			mapPaintWorldFloor(&pos)
		}
		for sub := *prefabWord(data, 16); sub != 0; sub = *prefabWord(sub, 16) {
			selectBorderPrimary(int32(*prefabWord(sub, 8)))
			selectBorderVariation(int32(*prefabWord(sub, 12)))
			selectTileImage(int32(*prefabWord(sub, 0)))
			selectTileVariation(int32(*prefabWord(sub, 4)))
			mapPaintWorldBorder(&pos)
		}
	}
	return 1
}
func prefabPlaceWalls(dx, dy int32) uint32 {
	s := GetServer().S()
	for n := (*prefabGlobal(prefabWalls)); n != 0; n = *prefabWord(n, 4) {
		src := (*server.Wall)(mapRoomPointer(*prefabWord(n, 0)))
		pos := image.Pt(int((dx+23*int32(src.X5))/23), int((dy+23*int32(src.Y6))/23))
		dst := s.Walls.GetWallAtGrid(pos)
		if dst != nil {
			if *mapPaintGlobal(paintWallMerge) != 0 {
				dst.Dir0 = mapPaintWallCompose(dst.Dir0, src.Dir0)
			} else {
				dst.Dir0 = src.Dir0
			}
		} else {
			dst = s.Walls.CreateAtGrid(pos)
			if dst == nil {
				return 0
			}
			dst.Dir0 = src.Dir0
		}
		dst.Tile1 = src.Tile1
		dst.Field2 = src.Field2
		dst.Health7 = src.Health7
		dst.Flags4 |= src.Flags4 & 0x80
		if dst.Field2 >= s.Walls.DefByInd(int(dst.Tile1)).Variations(int(dst.Dir0), 0) {
			dst.Field2 = 0
		}
		prefabTransferWallData(src, dst)
		dst.Flags4 |= src.Flags4 & 0x40
	}
	return 1
}
