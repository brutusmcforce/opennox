package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

var mapSectionMaxX, mapSectionMaxY uint32

func mapSectionLoadWall(x, y int32, dir byte, merge bool, generation bool) *server.Wall {
	var w *server.Wall
	if generation {
		w = (*server.Wall)(mapRoomPointer(*prefabWord(prefabWallNew(byte(x), byte(y)), 0)))
	} else {
		w = GetServer().S().Walls.GetWallAtGrid(image.Pt(int(x), int(y)))
		if w != nil {
			if merge {
				dir = mapPaintWallCompose(w.Dir0, dir)
			}
		} else {
			w = GetServer().S().Walls.CreateAtGrid(image.Pt(int(x), int(y)))
			if w == nil {
				return nil
			}
		}
	}
	w.Dir0 = dir
	return w
}
func mapSectionWallSprite(w *server.Wall) {
	w.Field2 = 0
	if w.Dir0 == 0 || w.Dir0 == 1 {
		w.Field2 = w.X5 % 3
	}
	if GetServer().S().Walls.DefByInd(int(w.Tile1)).Sprite(int(w.Dir0), int(w.Field2), int(w.Flags4>>2&2)) == nil {
		w.Field2 = 0
	}
}
func mapSectionWallFields(w *server.Wall, version int16, merge bool) {
	defs := &GetServer().S().Walls
	if version < 2 {
		w.Tile1 = 0
		w.Health7 = defs.DefByInd(0).Health41
		w.Field8 = w.Field8&0xff00 | 1
		return
	}
	w.Tile1 = mapSectionByte(w.Tile1)
	if defs.DefByInd(int(w.Tile1)) == nil {
		// a broken map or save: an unknown tile, use the first one as for old map versions
		w.Tile1 = 0
	}
	if version >= 3 {
		w.Field2 = mapSectionByte(w.Field2)
	} else {
		mapSectionWallSprite(w)
	}
	d := defs.DefByInd(int(w.Tile1))
	if merge && w.Field2 >= d.Variations(int(w.Dir0), 0) {
		w.Field2 = 0
	}
	w.Health7 = d.Health41
	if version < 4 {
		w.Field8 = w.Field8&0xff00 | 1
	} else if version == 6 {
		mapSectionByte(0)
		w.Field8 = w.Field8&0xff00 | 1
	} else {
		w.Field8 = w.Field8&0xff00 | uint16(mapSectionByte(byte(w.Field8)))
	}
	w.Field12 = 0
	if version >= 5 && version != 6 {
		w.Field12 = uint32(mapSectionByte(0))
	}
}
func mapSectionWalls(region *[8]uint32) uint32 {
	load := memmap.Uint32(0x5D4594, 739992)
	merge := load&1 != 0
	magic := memmap.PtrUint8(0x5D4594, 741372)
	if *magic == 0 {
		*magic = byte(GetServer().S().Walls.DefIndByName("MagicWallSystemUseOnly"))
	}
	*memmap.PtrUint32(0x5D4594, 741352) = 0
	*memmap.PtrUint32(0x5D4594, 741344) = 0
	version := int16(mapSectionWord(7))
	if version > 7 {
		return 0
	}
	minX, minY := memmap.PtrInt32(0x5D4594, 741360), memmap.PtrInt32(0x5D4594, 741368)
	var width, height int32
	if !cryptfile.Global().ReadOnly() {
		if region != nil {
			b := mapSectionRegionBounds(region)
			*minX, *minY = b[0], b[1]
			mapSectionMaxX, mapSectionMaxY = uint32(b[2]), uint32(b[3])
		} else {
			*minX, *minY = 256, 256
			mapSectionMaxX, mapSectionMaxY = 0, 0
			GetServer().S().Walls.EachWallXxx(func(w *server.Wall) bool {
				x, y := int32(w.X5), int32(w.Y6)
				if x < *minX {
					*minX = x
				}
				if y < *minY {
					*minY = y
				}
				if x > int32(mapSectionMaxX) {
					mapSectionMaxX = uint32(x)
				}
				if y > int32(mapSectionMaxY) {
					mapSectionMaxY = uint32(y)
				}
				return true
			})
		}
		width, height = int32(mapSectionMaxX)-*minX+1, int32(mapSectionMaxY)-*minY+1
	}
	mapSectionIO(unsafe.Pointer(minX), 4)
	mapSectionIO(unsafe.Pointer(minY), 4)
	width = int32(mapSectionDword(uint32(width)))
	height = int32(mapSectionDword(uint32(height)))
	if cryptfile.Global().ReadOnly() {
		var dx, dy int32
		if region != nil {
			b := mapSectionRegionBounds(region)
			dx, dy = b[0]-*minX, b[1]-*minY
			*minX, *minY = b[0], b[1]
		}
		gen := noxflags.HasGame(noxflags.GameFlag(0x400000))
		loadOne := func(x, y int32, encoded byte) bool {
			w := mapSectionLoadWall(x, y, encoded&127, merge, gen)
			if w == nil {
				return false
			}
			if encoded&128 != 0 {
				w.Flags4 |= 128
			}
			mapSectionWallFields(w, version, merge)
			return true
		}
		if version < 6 {
			for y := int32(0); y <= height; y++ {
				for x := int32(0); x <= width; x++ {
					dir := mapSectionByte(255)
					if dir == 255 {
						continue
					}
					wx, wy := *minX+x, *minY+y
					if gen {
						wx, wy = int32(byte(*minX))+x, int32(byte(*minY))+y
					}
					if !loadOne(wx, wy, dir) {
						return 0
					}
				}
			}
		} else {
			for {
				x := mapSectionByte(255)
				if x == 255 {
					break
				}
				y := mapSectionByte(0)
				dir := mapSectionByte(0)
				if !loadOne(int32(byte(int32(int8(x))+dx)), int32(byte(int32(y)+dy)), dir) {
					return 0
				}
			}
		}
		return 1
	}
	for y := *minY; y <= *minY+height; y++ {
		for x := *minX; x <= *minX+width; x++ {
			w := GetServer().S().Walls.GetWallAtGrid(image.Pt(int(x), int(y)))
			if w == nil || w.Tile1 == *magic {
				continue
			}
			if region != nil {
				pt := [2]int32{23 * int32(w.X5), 23 * int32(w.Y6)}
				if geometryWallPoint(&pt, (*[8]int32)(unsafe.Pointer(region))) == 0 {
					continue
				}
			}
			mapSectionByte(w.X5)
			mapSectionByte(w.Y6)
			// The loader restores wall flag 0x80 from the high bit of the direction, but the writer never stored it:
			// the decompiled check was "(int)flags >= 0" on an unsigned byte. Keep the files the same as before.
			mapSectionByte(w.Dir0)
			mapSectionByte(w.Tile1)
			mapSectionByte(w.Field2)
			pt := [2]int32{23*int32(w.X5) + 11, 23*int32(w.Y6) + 11}
			p := mapPolygonFind(&pt, 0, false)
			if p == nil {
				p = mapPolygonFindEdge(&pt, 0, 10)
			}
			light := byte(100)
			if p != nil {
				light = p.Level
			}
			mapSectionByte(light)
			tag := byte(w.Field12)
			if noxflags.HasGame(noxflags.GameFlag(0x200000)) {
				tag = 0
			}
			mapSectionByte(tag)
		}
	}
	mapSectionByte(255)
	return 1
}
