package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"unsafe"
)

func mapPaintNode(raw uint32) *[5]uint32 { return (*[5]uint32)(mapRoomPointer(raw)) }
func mapPaintSubtileNew(tile, variation, border, edge int32) *[5]uint32 {
	head := mapPaintGlobal(paintFreeHead)
	if *head == 0 {
		p := mapRoomCalloc(1, 200)
		*head = mapRoomRaw(p)
		for i := 0; i < 9; i++ {
			(*[5]uint32)(unsafe.Add(p, i*20))[4] = mapRoomRaw(unsafe.Add(p, (i+1)*20))
		}
	}
	p := mapPaintNode(*head)
	*head = p[4]
	*p = [5]uint32{uint32(tile), uint32(variation), uint32(border), uint32(edge), 0}
	return p
}
func mapPaintSubtileFree(p *[5]uint32) uint32 {
	p[4] = *mapPaintGlobal(paintFreeHead)
	raw := mapRoomRaw(unsafe.Pointer(p))
	*mapPaintGlobal(paintFreeHead) = raw
	return raw
}
func mapPaintSubtileClear(p *[5]uint32) uint32 {
	for raw := p[4]; raw != 0; {
		n := mapPaintNode(raw)
		raw = n[4]
		mapPaintSubtileFree(n)
	}
	p[4] = 0
	return 0
}
func mapPaintSubtileAdd(base *[5]uint32, tile, variation, border, edge, merge int32) uint32 {
	if tile == 255 || border == 255 {
		prev := base
		p := mapPaintNode(base[4])
		if p != nil {
			for p[4] != 0 {
				prev = p
				p = mapPaintNode(p[4])
			}
			mapPaintSubtileFree(p)
			prev[4] = 0
		}
		return 1
	}
	if base[0] == 255 {
		return 1
	}
	tail := base
	found := false
	for p := mapPaintNode(base[4]); p != nil; p = mapPaintNode(p[4]) {
		tail = p
		if p[0] == uint32(tile) && p[1] == uint32(variation) && p[2] == uint32(border) {
			if merge == 0 {
				if p[3] == uint32(edge) {
					return 0
				}
			} else if mergeBorderEdge((*[4]uint32)(unsafe.Pointer(p)), edge) {
				found = true
			}
		}
	}
	if !found {
		if merge != 0 {
			edge = generateBorderEdge(border, edge)
		}
		tail[4] = mapRoomRaw(unsafe.Pointer(mapPaintSubtileNew(tile, variation, border, edge)))
		return 1
	}
	for {
		removed := false
		for p := mapPaintNode(base[4]); p != nil; p = mapPaintNode(p[4]) {
			if p[0] != uint32(tile) || p[1] != uint32(variation) || p[2] != uint32(border) {
				continue
			}
			prev := p
			for q := mapPaintNode(prev[4]); q != nil; q = mapPaintNode(prev[4]) {
				if q[0] == uint32(tile) && q[1] == uint32(variation) && q[2] == uint32(border) && mergeBorderEdge((*[4]uint32)(unsafe.Pointer(p)), normalizeBorderEdge(border, int32(q[3]))) {
					prev[4] = q[4]
					mapPaintSubtileFree(q)
					removed = true
				} else {
					prev = q
				}
			}
		}
		if !removed {
			return 1
		}
	}
}
func mapPaintFloorCell(x, y int32, p *[8]int32, side, merge int32) uint32 {
	tile, variation := int32(255), int32(0)
	if p != nil {
		tile = p[0]
		if p[2] != 0 && byte(p[5]) == 0 {
			p[3], p[4] = x, y
			p[5] = (p[5] &^ 255) | int32(byte(side))
		} else {
			dx, dy, adjust := x, y, int32(bool2int(side != 1))
			if p[2] != 0 {
				dx, dy = x-p[3], y-p[4]
				adjust = 0
				if int32(byte(p[5])) != side {
					if byte(p[5]) == 1 {
						adjust = 1
					} else {
						adjust = -1
					}
				}
			}
			def := &tileDefinitionsAll()[tile]
			w, h := int32(def.Field52), int32(def.Field53) // the variation grid size
			a, b := (dy+dx)%w, (adjust+dy-dx)%h
			if a < 0 {
				a += w
			}
			if b < 0 {
				b += h
			}
			variation = a + b*w
			if *mapPaintSelection(35916) == 1 {
				variation = int32(*mapPaintGlobal(paintVariation))
			}
		}
	}
	if side != 1 && side != 2 {
		return 0
	}
	cell := mapPaintCell(x, y)
	off := 1
	if side == 2 {
		off = 6
	}
	node := (*[5]uint32)(unsafe.Pointer(&cell[off]))
	if *mapPaintGlobal(paintSubtileMode) == 1 {
		border, edge := int32(255), int32(0)
		if p != nil {
			border, edge = p[6], p[7]
		}
		return mapPaintSubtileAdd(node, tile, variation, border, edge, merge)
	}
	if node[0] == uint32(tile) && node[1] == uint32(variation) {
		return 0
	}
	node[0], node[1] = uint32(tile), uint32(variation)
	if tile == 255 {
		mapPaintSubtileClear(node)
	}
	if p != nil {
		cell[0] |= uint32(side)
	} else {
		cell[0] &= ^uint32(side)
	}
	return 1
}
func mapPaintFloorRoute(x, y, sx, sy int32, p *[8]int32) uint32 {
	if x-1 <= 0 || x >= 127 || y-1 <= 0 || y >= 127 {
		return 0
	}
	if sx <= sy {
		if 46-sx <= sy {
			return mapPaintFloorCell(x, y, p, 2, 0)
		}
		return mapPaintFloorCell(x-1, y, p, 1, 0)
	}
	if 46-sx <= sy {
		return mapPaintFloorCell(x, y, p, 1, 0)
	}
	return mapPaintFloorCell(x, y-1, p, 2, 0)
}
func mapPaintWorldTile(pos *types.Pointf, border bool) uint32 {
	if border {
		*mapPaintGlobal(paintSubtileMode) = 1
		defer func() { *mapPaintGlobal(paintSubtileMode) = 0 }()
	}
	p := [8]int32{int32(*mapPaintSelection(35912)), int32(*mapPaintGlobal(paintVariation)), 0, -1, -1, 0, 255, 0}
	if border {
		p[6], p[7] = int32(*mapPaintGlobal(paintBorder)), int32(*mapPaintGlobal(paintBorderVariation))
	}
	x, y := float64(pos.X)+11.5, float64(pos.Y)+11.5
	gx, gy := int32(int64(x*0.021739131)), int32(int64(y*0.021739131))
	sx, sy := int32(int64(x))%46, int32(int64(float32(y)))%46
	if p[0] == 255 {
		return mapPaintFloorRoute(gx, gy, sx, sy, nil)
	}
	return mapPaintFloorRoute(gx, gy, sx, sy, &p)
}
func mapPaintWorldFloor(pos *types.Pointf) uint32  { return mapPaintWorldTile(pos, false) }
func mapPaintWorldBorder(pos *types.Pointf) uint32 { return mapPaintWorldTile(pos, true) }
func mapPaintFloorPoint(pos *types.Pointf) uint32 {
	if pos == nil {
		return 0
	}
	if *mapPaintSelection(35912) == 255 {
		return 1
	}
	var world types.Pointf
	mapPaintTransform(pos, &world)
	if mapPaintWorldFloor(&world) == 0 {
		return 0
	}
	if *mapPaintGlobal(paintAfterWalls) == 1 {
		return mapPaintWorldWalls(&world)
	}
	return 1
}
func mapPaintRect(_ unsafe.Pointer, pos *types.Pointf, w, h int32) uint32 {
	result := uint32(h)
	p := types.Pointf{Y: float32(float64(pos.Y) + 16.263456)}
	for j := int32(0); j < h; j++ {
		p.X = float32(float64(pos.X) + 16.263456)
		for i := int32(0); i < w; i++ {
			result = mapPaintFloorPoint(&p)
			p.X = float32(float64(p.X) + 32.526913)
		}
		p.Y = float32(float64(p.Y) + 32.526913)
	}
	return result
}
func mapPaintLineY(_ unsafe.Pointer, pos *types.Pointf, n int32) uint32 {
	result := mapRoomRaw(unsafe.Pointer(pos))
	p := types.Pointf{X: float32(float64(pos.X) + 16.263456), Y: float32(float64(pos.Y) + 16.263456)}
	for i := int32(0); i < n; i++ {
		result = mapPaintFloorPoint(&p)
		p.Y = float32(float64(p.Y) + 32.526913)
	}
	return result
}
