package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"unsafe"
)

func mapSectionRegionBounds(region *[8]uint32) [4]int32 {
	var b [4]int32
	geometryWallBounds(region, &b)
	for i := range b {
		b[i] /= 23
	}
	return b
}
func mapSectionFloorHalf(x, y int32, create bool) *[5]uint32 {
	half := uint32(2)
	gx, gy := 23*x/46, 23*y/46
	base := 6
	if y&1 != 0 {
		half = 1
		gx = (23*x - 23) / 46
		gy = (23*y + 23) / 46
		base = 1
	}
	if create && noxflags.HasGame(noxflags.GameFlag(0x400000)) {
		return mapPaintNode(*prefabWord(prefabTileNew(gx, gy, half), 0))
	}
	c := mapPaintCell(gx, gy)
	if create {
		c[0] |= half
	} else if c[0]&half == 0 {
		return nil
	}
	return (*[5]uint32)(unsafe.Pointer(&c[base]))
}
func mapSectionFloorClearFlags() {
	for x := int32(0); x < 128; x++ {
		for y := int32(0); y < 128; y++ {
			mapPaintCell(x, y)[0] &= 0xffffff00
		}
	}
}
func mapSectionFloorHistorical(region *[8]uint32) uint32 {
	x0, y0 := int32(mapSectionDword(0)), int32(mapSectionDword(0))
	width, height := int32(mapSectionDword(0)), int32(mapSectionDword(0))
	if region == nil {
		mapSectionFloorClearFlags()
	} else {
		b := mapSectionRegionBounds(region)
		x0, y0 = b[0], b[1]
	}
	readTile := func(p *[5]uint32) {
		mapSectionIO(unsafe.Pointer(p), 16)
		count := mapSectionByte(0)
		p[4] = 0
		tail := p
		for i := 0; i < int(count); i++ {
			node := mapPaintSubtileNew(0, 0, 0, 0)
			mapSectionIO(unsafe.Pointer(node), 16)
			tail[4] = mapRoomRaw(unsafe.Pointer(node))
			tail = node
		}
	}
	for y := y0; y < y0+height; y++ {
		for x := x0; x < x0+width; x++ {
			if region != nil {
				if mapSectionByte(0) != 0 {
					readTile(mapSectionFloorHalf(x, y, true))
				}
				continue
			}
			c := mapPaintCell(x, y)
			c[0] = c[0]&0xffffff00 | uint32(mapSectionByte(byte(c[0])))
			if c[0]&1 != 0 {
				readTile((*[5]uint32)(unsafe.Pointer(&c[1])))
			}
			if c[0]&2 != 0 {
				readTile((*[5]uint32)(unsafe.Pointer(&c[6])))
			}
		}
	}
	return 1
}
func mapSectionFloor(region *[8]uint32) uint32 {
	version := int16(mapSectionWord(4))
	if version > 4 || version < 3 {
		return 0
	}
	if version == 3 {
		return mapSectionFloorHistorical(region)
	}
	if cryptfile.Global().ReadOnly() {
		x0, y0 := int32(mapSectionDword(0)), int32(mapSectionDword(0))
		mapSectionDword(0)
		mapSectionDword(0)
		var dx, dy int32
		if region == nil {
			mapSectionFloorClearFlags()
		} else {
			b := mapSectionRegionBounds(region)
			dx, dy = b[0]-x0, b[1]-y0
		}
		for {
			coord := mapSectionWord(0)
			if coord == 0xffff {
				break
			}
			if region != nil {
				mapSectionTile(mapSectionFloorHalf(dx+int32(coord>>8), dy+int32(byte(coord)), true))
				continue
			}
			c := mapPaintCell(int32(coord>>8&127), int32(coord&127))
			if coord&0x8000 != 0 {
				c[0] |= 1
				mapSectionTile((*[5]uint32)(unsafe.Pointer(&c[1])))
			}
			if coord&0x80 != 0 {
				c[0] |= 2
				mapSectionTile((*[5]uint32)(unsafe.Pointer(&c[6])))
			}
		}
		return 1
	}
	var b [4]int32
	if region != nil {
		b = mapSectionRegionBounds(region)
	} else {
		b = [4]int32{128, 128, -1, -1}
		for x := int32(0); x < 128; x++ {
			for y := int32(0); y < 128; y++ {
				if byte(mapPaintCell(x, y)[0]) == 0 {
					continue
				}
				if x < b[0] {
					b[0] = x
				}
				if y < b[1] {
					b[1] = y
				}
				if x > b[2] {
					b[2] = x
				}
				if y > b[3] {
					b[3] = y
				}
			}
		}
		if b[2] < 0 {
			return 1
		}
	}
	mapSectionDword(uint32(b[0]))
	mapSectionDword(uint32(b[1]))
	mapSectionDword(uint32(b[2] - b[0] + 1))
	mapSectionDword(uint32(b[3] - b[1] + 1))
	for y := b[1]; y <= b[3]; y++ {
		for x := b[0]; x <= b[2]; x++ {
			if region != nil {
				if (x+y)&1 != 0 {
					continue
				}
				p := mapSectionFloorHalf(x, y, false)
				if p == nil {
					continue
				}
				point := [2]int32{23*x + 11, 23*y + 34}
				if geometryWallPoint(&point, (*[8]int32)(unsafe.Pointer(region))) == 0 {
					continue
				}
				mapSectionWord(uint16(x<<8 | y))
				mapSectionTile(p)
				continue
			}
			c := mapPaintCell(x, y)
			flags := c[0] & 3
			if flags == 0 {
				continue
			}
			coord := uint16(x<<8 | y)
			if flags&1 != 0 {
				coord |= 0x8000
			}
			if flags&2 != 0 {
				coord |= 0x80
			}
			mapSectionWord(coord)
			if flags&1 != 0 {
				mapSectionTile((*[5]uint32)(unsafe.Pointer(&c[1])))
			}
			if flags&2 != 0 {
				mapSectionTile((*[5]uint32)(unsafe.Pointer(&c[6])))
			}
		}
	}
	mapSectionWord(0xffff)
	return 1
}
