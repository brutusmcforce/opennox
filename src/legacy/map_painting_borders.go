package legacy

import "github.com/brutusmcforce/opennox/libs/types"

func mapPaintBorderNeighbor(x, y, side, key int32, p *[8]int32, edge int32) {
	if x <= 0 || x >= 127 || y <= 0 || y >= 127 || side&1 != 0 && y == 1 || side&2 != 0 && x == 1 {
		return
	}
	cell := mapPaintCell(x, y)
	if side&2 != 0 && cell[6] == uint32(key) {
		return
	}
	if side&1 != 0 && cell[1] == uint32(key) {
		return
	}
	if p != nil {
		p[7] = edge
		mapPaintFloorCell(x, y, p, side, 1)
	}
}
func mapPaintBorderPoint(pos *types.Pointf) uint32 {
	if *mapPaintGlobal(paintBorder) == 255 {
		return 1
	}
	var world types.Pointf
	if mapPaintTransform(pos, &world) == 0 {
		return 0
	}
	p := [8]int32{int32(*mapPaintSelection(35912)), int32(*mapPaintGlobal(paintVariation)), 0, 0, 0, 0, int32(*mapPaintGlobal(paintBorder)), int32(*mapPaintGlobal(paintBorderVariation))}
	*mapPaintGlobal(paintSubtileMode) = 1
	xy := [2]int32{int32(int64(world.X)), int32(int64(world.Y))}
	r := mapPaintFloodBorder(&xy, &p, 46)
	*mapPaintGlobal(paintSubtileMode) = 0
	return r
}
func mapPaintFloodBorder(point *[2]int32, p *[8]int32, size int32) uint32 {
	fx, fy := float64(point[0])+11.5, float64(point[1])+11.5
	x, y := int32(int64(fx*0.021739131)), int32(int64(fy*0.021739131))
	sx, sy := int32(int64(fx))%46, int32(int64(float32(fy)))%46
	*mapPaintSelection(22200) = 0
	*mapPaintGlobal(paintWorkCount) = 0
	key := size
	if x > 0 && x < 127 && y > 0 && y < 127 {
		side := int32(1)
		if sx <= sy {
			if size-sx <= sy {
				side = 2
			} else {
				x--
			}
		} else if size-sx > sy {
			y--
			side = 2
		}
		off := 1
		if side == 2 {
			off = 6
		}
		key = int32(mapPaintCell(x, y)[off])
		pushTileFill(x, y, side, key)
	}
	stack := tileFillStack()
	for i := int32(0); i < int32(*mapPaintGlobal(paintWorkCount)); i++ {
		q := stack[i]
		x, y, side := int32(q.X), int32(q.Y), int32(q.Flags)
		if side&2 != 0 {
			pushTileFill(x, y, 1, key)
			pushTileFill(x, y+1, 1, key)
			pushTileFill(x-1, y, 1, key)
			pushTileFill(x-1, y+1, 1, key)
		} else if side&1 != 0 {
			pushTileFill(x+1, y, 2, key)
			pushTileFill(x+1, y-1, 2, key)
			pushTileFill(x, y, 2, key)
			pushTileFill(x, y-1, 2, key)
		}
	}
	var px, py, side uint32
	for popTileFill(&px, &py, &side) {
		x, y := int32(px), int32(py)
		emit := func(dx, dy, s, e int32) { mapPaintBorderNeighbor(x+dx, y+dy, s, key, p, e) }
		if side&2 != 0 {
			emit(0, 0, 1, 1)
			emit(0, 1, 1, 4)
			emit(-1, 0, 1, 3)
			emit(-1, 1, 1, 6)
			emit(0, -1, 2, 0)
			emit(1, 0, 2, 2)
			emit(-1, 0, 2, 5)
			emit(0, 1, 2, 7)
		} else if side&1 != 0 {
			emit(1, 0, 2, 4)
			emit(1, -1, 2, 1)
			emit(0, 0, 2, 6)
			emit(0, -1, 2, 3)
			emit(0, -1, 1, 0)
			emit(1, 0, 1, 2)
			emit(-1, 0, 1, 5)
			emit(0, 1, 1, 7)
		}
	}
	return uint32(bool2int(*mapPaintSelection(22200) == 0))
}
