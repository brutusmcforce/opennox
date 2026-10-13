package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
)

func mapPaintRoomFloor(config unsafe.Pointer, r *mapRoom, layout unsafe.Pointer) uint32 {
	selectTileName((*byte)(unsafe.Add(layout, 60)))
	mapPaintRect(config, &r.Pos, r.Width, r.Height)
	size := [2]int32{}
	var p types.Pointf
	switch r.Kind {
	case 1:
		size = [2]int32{r.Width - 2, r.Height - 2}
		p = types.Pointf{X: float32(float64(r.Pos.X) + 32.526913), Y: float32(float64(r.Pos.Y) + 32.526913)}
	case 2:
		size[0] = r.Width - 2
		p.X = float32(float64(r.Pos.X) + 32.526913)
		if n := r.Neighbors[0][0]; n != nil && n.Kind != 1 {
			p.Y = r.Pos.Y
			size[1] = r.Height
		} else {
			size[1] = r.Height - 1
			p.Y = float32(float64(r.Pos.Y) + 32.526913)
		}
		if n := r.Neighbors[1][0]; n != nil {
			if n.Kind != 1 {
				size[1]++
			} else {
				size[1]--
			}
		}
	case 3:
		size[0] = r.Width - 2
		p.X = float32(float64(r.Pos.X) + 32.526913)
		p.Y = r.Pos.Y
		size[1] = r.Height - 1
		if n := r.Neighbors[1][0]; n != nil && n.Kind != 1 {
			size[1] = r.Height
		}
		if n := r.Neighbors[0][0]; n != nil {
			if n.Kind == 1 {
				size[1]--
				p.Y = float32(float64(p.Y) + 32.526913)
			} else {
				size[1]++
				p.Y = float32(float64(p.Y) - 32.526913)
			}
		}
	case 4:
		p.X = r.Pos.X
		p.Y = float32(float64(r.Pos.Y) + 32.526913)
		size[1] = r.Height - 2
		size[0] = r.Width - 1
		if n := r.Neighbors[2][0]; n != nil && n.Kind != 1 {
			size[0] = r.Width
		}
		if n := r.Neighbors[3][0]; n != nil {
			if n.Kind != 1 {
				p.X = float32(float64(p.X) - 32.526913)
				size[0]++
			} else {
				size[0]--
				p.X = float32(float64(p.X) + 32.526913)
			}
		}
	case 5:
		p.Y = float32(float64(r.Pos.Y) + 32.526913)
		size[1] = r.Height - 2
		if n := r.Neighbors[3][0]; n != nil && n.Kind != 1 {
			size[0] = r.Width
			p.X = r.Pos.X
		} else {
			size[0] = r.Width - 1
			p.X = float32(float64(r.Pos.X) + 32.526913)
		}
		if n := r.Neighbors[2][0]; n != nil {
			if n.Kind == 1 {
				size[0]--
			} else {
				size[0]++
			}
		}
	}
	for pat := mapRoomRef(layout, 120); pat != nil; pat = mapRoomRef(pat, 124) {
		if size[0] <= 0 || size[1] <= 0 {
			break
		}
		switch *mapRoomWord(pat, 0) {
		case 1:
			mapPaintPatternDiamond(config, pat, &p, &size)
		case 2:
			mapPaintPatternRandom(config, pat, &p, &size)
		default:
			mapPaintPatternFill(config, pat, &p, &size)
		}
		size[0] -= 2
		size[1] -= 2
		p.X = float32(float64(p.X) + 32.526913)
		p.Y = float32(float64(p.Y) + 32.526913)
	}
	return uint32(size[0])
}
func mapPaintPatternCenter(pat unsafe.Pointer, pos *types.Pointf, size *[2]int32) uint32 {
	selectBorderName((*byte)(unsafe.Add(pat, 4)))
	p := types.Pointf{X: float32(float64(size[0])*16.263456 + float64(pos.X)), Y: float32(float64(size[1])*16.263456 + float64(pos.Y))}
	return mapPaintBorderPoint(&p)
}
func mapPaintPatternFill(config, pat unsafe.Pointer, pos *types.Pointf, size *[2]int32) uint32 {
	selectTileName((*byte)(unsafe.Add(pat, 64)))
	mapPaintRect(config, pos, size[0], size[1])
	return mapPaintPatternCenter(pat, pos, size)
}
func mapPaintPatternDiamond(config, pat unsafe.Pointer, pos *types.Pointf, size *[2]int32) uint32 {
	selectTileName((*byte)(unsafe.Add(pat, 64)))
	third := size[0] / 3
	offset := float64(third) * 32.526913
	span := size[0] - 2*third
	p := types.Pointf{X: float32(offset + float64(pos.X)), Y: pos.Y}
	for i := int32(0); i <= size[1]/2; i++ {
		mapPaintRect(config, &p, span, 1)
		span += 2
		p.X = float32(float64(p.X) - 32.526913)
		p.Y = float32(float64(p.Y) + 32.526913)
		if span >= size[0] {
			span = size[0]
			p.X = pos.X
		}
	}
	span = size[0] - 2*third
	p.X = float32(offset + float64(pos.X))
	p.Y = float32(float64(size[1]-1)*32.526913 + float64(pos.Y))
	for i := int32(0); i < size[1]/2; i++ {
		mapPaintRect(config, &p, span, 1)
		span += 2
		p.X = float32(float64(p.X) - 32.526913)
		p.Y = float32(float64(p.Y) - 32.526913)
		if span >= size[0] {
			span = size[0]
			p.X = pos.X
		}
	}
	return mapPaintPatternCenter(pat, pos, size)
}
func mapPaintPatternRandom(config, pat unsafe.Pointer, pos *types.Pointf, size *[2]int32) {
	if size[0] < 3 || size[1] < 3 {
		return
	}
	selectTileName((*byte)(unsafe.Add(pat, 64)))
	p := types.Pointf{X: float32(float64(pos.X) + 32.526913), Y: float32(float64(pos.Y) + 32.526913)}
	mapPaintRect(config, &p, max(size[0]-2, 1), max(size[1]-2, 1))
	if size[0] < 4 {
		p.X = float32(float64(pos.X) + 32.526913)
		p.Y = pos.Y
		mapPaintRect(config, &p, 1, 1)
		p.Y = float32(float64(p.Y) + 65.053825)
		mapPaintRect(config, &p, 1, 1)
	} else {
		n := mapRoomRandomInt(1, size[0]-3)
		start := mapRoomRandomInt(1, size[0]-n-2)
		p.Y = pos.Y
		p.X = float32(float64(start)*32.526913 + float64(pos.X))
		mapPaintRect(config, &p, n, 1)
		n = mapRoomRandomInt(1, size[0]-3)
		start = mapRoomRandomInt(1, size[0]-n-2)
		p.X = float32(float64(start)*32.526913 + float64(pos.X))
		p.Y = float32(float64(size[1]-1)*32.526913 + float64(pos.Y))
		mapPaintRect(config, &p, n, 1)
	}
	if size[1] < 4 {
		p.Y = float32(float64(pos.Y) + 32.526913)
		p.X = pos.X
		mapPaintRect(config, &p, 1, 1)
		p.X = float32(float64(p.X) + 65.053825)
		mapPaintRect(config, &p, 1, 1)
	} else {
		n := mapRoomRandomInt(1, size[1]-3)
		start := mapRoomRandomInt(1, size[1]-n-2)
		p.X = pos.X
		p.Y = float32(float64(start)*32.526913 + float64(pos.Y))
		mapPaintLineY(config, &p, n)
		n = mapRoomRandomInt(1, size[1]-3)
		start = mapRoomRandomInt(1, size[1]-n-2)
		p.X = float32(float64(size[0]-1)*32.526913 + float64(pos.X))
		p.Y = float32(float64(start)*32.526913 + float64(pos.Y))
		mapPaintLineY(config, &p, n)
	}
	mapPaintPatternCenter(pat, pos, size)
}
