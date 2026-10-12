package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

func mapPaintWallGrid(pos *types.Pointf) (image.Point, bool) {
	var p types.Pointf
	if mapPaintTransform(pos, &p) == 0 {
		return image.Point{}, false
	}
	x, y := int32(int64(float64(p.X)*0.043478262)), int32(int64(float64(p.Y)*0.043478262))
	return image.Pt(int(x), int(y)), (x+y)&1 == 0 && x >= 0 && x < 256 && y >= 0 && y < 256
}
func mapPaintFindWall(pos *types.Pointf, out *uint32) uint32 {
	pt, ok := mapPaintWallGrid(pos)
	if !ok || out == nil {
		return 0
	}
	w := GetServer().S().Walls.GetWallAtGrid(pt)
	if w == nil {
		return 0
	}
	*out = mapRoomRaw(w.C())
	return 1
}
func mapPaintMakeWall(pos *types.Pointf) uint32 {
	if *mapPaintSelection(35948) == 255 {
		return 1
	}
	pt, ok := mapPaintWallGrid(pos)
	if !ok {
		return 0
	}
	walls := &GetServer().S().Walls
	w := walls.GetWallAtGrid(pt)
	existing := w != nil
	if !existing {
		w = walls.CreateAtGrid(pt)
		if w == nil {
			return 0
		}
	}
	w.Tile1 = byte(*mapPaintSelection(35948))
	dir := byte(*mapPaintSelection(35952))
	if existing && *mapPaintGlobal(paintWallMerge) == 1 {
		dir = mapPaintWallCompose(dir, w.Dir0)
	}
	w.Dir0 = dir
	w.Field2 = byte(*mapPaintSelection(35956))
	if *mapPaintGlobal(paintWallCycle) != 0 {
		x := pt.X
		if existing {
			x = int(w.X5)
		}
		w.Field2 = byte(x % int(walls.DefByInd(int(w.Tile1)).Field749))
	}
	if w.Field2 >= walls.DefByInd(int(w.Tile1)).Variations(int(w.Dir0), 0) {
		w.Field2 = 0
	}
	w.Flags4 = 128
	return 1
}
func mapPaintUnlinkSecret(w *server.Wall) {
	if w.Data != nil {
		worldSecretRemove(w.Data)
		w.Data = nil
	}
}
func mapPaintEraseWall(pos *types.Pointf) uint32 {
	pt, ok := mapPaintWallGrid(pos)
	if !ok {
		return 0
	}
	walls := &GetServer().S().Walls
	w := walls.GetWallAtGrid(pt)
	if w == nil {
		return 0
	}
	mapPaintUnlinkSecret(w)
	walls.DeleteAtGrid(pt)
	return 1
}
func mapPaintWallLineX(pos *types.Pointf, n int32) uint32 {
	p := *pos
	mapPaintWallDirection(0)
	r := uint32(n)
	for i := int32(0); i < n; i++ {
		r = mapPaintMakeWall(&p)
		p.X = float32(float64(p.X) + 32.526913)
	}
	return r
}
func mapPaintWallLineY(pos *types.Pointf, n int32) uint32 {
	p := *pos
	r := mapPaintWallDirection(1)
	for i := int32(0); i < n; i++ {
		r = mapPaintMakeWall(&p)
		p.Y = float32(float64(p.Y) + 32.526913)
	}
	return r
}
func mapPaintEraseLine(pos *types.Pointf, n int32, vertical bool) uint32 {
	p := *pos
	r := mapRoomRaw(unsafe.Pointer(pos))
	for i := int32(0); i < n; i++ {
		r = mapPaintEraseWall(&p)
		if vertical {
			p.Y = float32(float64(p.Y) + 32.526913)
		} else {
			p.X = float32(float64(p.X) + 32.526913)
		}
	}
	return r
}
func mapPaintEraseLineX(pos *types.Pointf, n int32) uint32 { return mapPaintEraseLine(pos, n, false) }
func mapPaintEraseLineY(pos *types.Pointf, n int32) uint32 { return mapPaintEraseLine(pos, n, true) }
func mapPaintWallCorner(pos *types.Pointf) uint32 {
	points := [4]types.Pointf{{X: pos.X, Y: float32(float64(pos.Y) - 32.526913)}, {X: pos.X, Y: float32(float64(pos.Y) + 32.526913)}, {X: float32(float64(pos.X) - 32.526913), Y: pos.Y}, {X: float32(float64(pos.X) + 32.526913), Y: pos.Y}}
	mask := 0
	var raw uint32
	for i := range points {
		if mapPaintFindWall(&points[i], &raw) != 0 {
			mask |= 1 << i
		}
	}
	dirs := [16]int32{-1, -1, -1, 1, -1, 10, 9, 6, -1, 7, 8, 4, 0, 3, 5, 2}
	if dirs[mask] < 0 {
		return uint32(mask - 3)
	}
	mapPaintWallDirection(dirs[mask])
	return mapPaintMakeWall(pos)
}
func mapPaintWorldWalls(pos *types.Pointf) uint32 {
	p := [2]int32{int32(int64(float64(pos.X) * 0.043478262)), int32(int64(float64(pos.Y) * 0.043478262))}
	valid := func() bool { return p[0] > 0 && p[0] < 255 && p[1] > 0 && p[1] < 255 }
	if valid() {
		mapPaintGridWalls(&p)
	}
	p[0]--
	if valid() {
		mapPaintGridWalls(&p)
	}
	p[0]++
	p[1]--
	if valid() {
		mapPaintGridWalls(&p)
	}
	return 1
}
func mapPaintGridWalls(point *[2]int32) uint32 {
	x0, y0 := int32(0), int32(0)
	if v := point[0] - 1; v > 0 {
		x0 = v & 0xfffe
	}
	if v := point[1] - 1; v > 0 {
		y0 = v & 0xfffe
	}
	walls := &GetServer().S().Walls
	for y := y0; y <= y0+4; y++ {
		for x := x0; x <= x0+4; x++ {
			if w := walls.GetWallAtGrid(image.Pt(int(x), int(y))); w != nil && w.Flags4&0x8c == 0 {
				mapPaintUnlinkSecret(w)
				walls.DeleteAtGrid(w.GridPos())
			}
		}
	}
	emit := func(x, y int32, mask uint32) {
		dir := memmap.Uint32(0x587000, 255052+4*uintptr(mask))
		if dir == 255 {
			return
		}
		pt := image.Pt(int(x), int(y))
		w := walls.GetWallAtGrid(pt)
		if w != nil {
			if w.Flags4&0x8c != 0 || w.Dir0 == byte(dir) || uint32(w.Tile1) == *mapPaintSelection(35948) {
				return
			}
		} else {
			w = walls.CreateAtGrid(pt)
			if w == nil {
				return
			}
			w.Tile1 = byte(*mapPaintSelection(35948))
		}
		w.Dir0 = byte(dir)
		w.Field2 = 0
		if dir == 0 || dir == 1 {
			w.Field2 = byte(int(w.X5) % int(walls.DefByInd(int(w.Tile1)).Field749))
		}
	}
	for y := y0; y <= y0+4; y += 2 {
		for x := x0; x <= x0+4; x += 2 {
			gx, gy := x/2, y/2
			if gx >= 128 || gy >= 128 {
				continue
			}
			c := mapPaintCell(gx, gy)
			mask := uint32(0)
			if c[6] != 255 {
				mask |= 8
			}
			if c[1] != 255 {
				mask |= 2
			}
			if gy > 0 && mapPaintCell(gx, gy-1)[6] != 255 {
				mask |= 4
			}
			if gx > 0 && mapPaintCell(gx-1, gy)[1] != 255 {
				mask |= 1
			}
			emit(x, y, mask)
			mask = 0
			if gy > 0 && mapPaintCell(gx, gy-1)[6] != 255 {
				mask |= 2
			}
			if gx > 0 {
				if mapPaintCell(gx-1, gy)[1] != 255 {
					mask |= 8
				}
				if gy > 0 {
					c = mapPaintCell(gx-1, gy-1)
					if c[1] != 255 {
						mask |= 4
					}
					if c[6] != 255 {
						mask |= 1
					}
				}
			}
			emit(x-1, y-1, mask)
		}
	}
	return uint32(x0 + 6)
}
