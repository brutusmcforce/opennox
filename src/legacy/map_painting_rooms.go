package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"unsafe"
)

func mapPaintRoomWalls(config unsafe.Pointer, r *mapRoom) {
	if r.Flags&2 != 0 {
		return
	}
	layout := mapPaintLayout(r.Decoration)
	mapPaintRoomFloor(config, r, layout)
	Sub_526CA0(GoStringP(layout))
	p := r.Min
	mapPaintWallDirection(8)
	mapPaintMakeWall(&p)
	p.X = float32(float64(p.X) + 32.526913)
	mapPaintWallLineX(&p, r.Width-1)
	p = types.Pointf{X: r.Max.X, Y: r.Min.Y}
	mapPaintWallDirection(9)
	mapPaintMakeWall(&p)
	p = types.Pointf{X: r.Min.X, Y: r.Max.Y}
	mapPaintWallDirection(7)
	mapPaintMakeWall(&p)
	p.X = float32(float64(p.X) + 32.526913)
	mapPaintWallLineX(&p, r.Width-1)
	p = r.Max
	mapPaintWallDirection(10)
	mapPaintMakeWall(&p)
	p = types.Pointf{X: r.Min.X, Y: float32(float64(r.Min.Y) + 32.526913)}
	mapPaintWallLineY(&p, r.Height-1)
	p = types.Pointf{X: r.Max.X, Y: float32(float64(r.Min.Y) + 32.526913)}
	mapPaintWallLineY(&p, r.Height-1)
	mapPaintMergeWalls(0)
	for dir := 0; dir < 4; dir++ {
		for i := 0; i < int(r.Counts[dir]); i++ {
			mapPaintJoinWalls(r, r.Neighbors[dir][i], int32(dir))
		}
	}
	mapPaintMergeWalls(1)
	mapPaintRoomDoors(config, r)
}
func mapPaintJoinWalls(r, t *mapRoom, dir int32) uint32 {
	var p, q types.Pointf
	switch dir {
	case 0, 1:
		p.X = max(r.Min.X, t.Min.X)
		p.Y = r.Min.Y
		if dir == 1 {
			p.Y = r.Max.Y
		}
		q = p
		if dir == 1 {
			q.Y = float32(float64(q.Y) - 32.526913)
		}
		n, w := r.Width, r.Size.X
		if r.Width >= t.Width {
			n, w = t.Width, t.Size.X
		}
		mapRoomAddExclusion(r, &q, w, float32(32.526913))
		q = types.Pointf{X: float32(float64(p.X) + 32.526913), Y: p.Y}
		mapPaintEraseLineX(&q, n-1)
		mapPaintWallCorner(&p)
		p.X = min(r.Max.X, t.Max.X)
		return mapPaintWallCorner(&p)
	case 2, 3:
		p.X = r.Max.X
		if dir == 3 {
			p.X = r.Min.X
		}
		p.Y = max(r.Min.Y, t.Min.Y)
		q = p
		if dir == 2 {
			q.X = float32(float64(q.X) - 32.526913)
		}
		n, h := r.Height, r.Size.Y
		if r.Height >= t.Height {
			n, h = t.Height, t.Size.Y
		}
		mapRoomAddExclusion(r, &q, float32(32.526913), h)
		q = types.Pointf{X: p.X, Y: float32(float64(p.Y) + 32.526913)}
		mapPaintEraseLineY(&q, n-1)
		mapPaintWallCorner(&p)
		p.Y = min(r.Max.Y, t.Max.Y)
		return mapPaintWallCorner(&p)
	default:
		return uint32(dir)
	}
}
