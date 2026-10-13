package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
)

func mapRoomRandomX(r, t *mapRoom) float64 {
	return float64(mapRoomRandomInt(0, r.Width-t.Width))*32.526913 + float64(r.Pos.X)
}
func mapRoomRandomY(r, t *mapRoom) float64 {
	return float64(mapRoomRandomInt(0, r.Height-t.Height))*32.526913 + float64(r.Pos.Y)
}
func mapRoomRandomReverseX(r, t *mapRoom) float64 {
	return float64(t.Pos.X) - float64(mapRoomRandomInt(0, r.Width-t.Width))*32.526913
}
func mapRoomRandomReverseY(r, t *mapRoom) float64 {
	return float64(t.Pos.Y) - float64(mapRoomRandomInt(0, r.Height-t.Height))*32.526913
}
func mapRoomAddExclusion(r *mapRoom, pos *types.Pointf, w, h float32) *mapRoomExclusion {
	e := mapRoomAlloc[mapRoomExclusion]()
	if e != nil {
		e.Min = *pos
		e.Max = types.Ptf(w+pos.X, h+pos.Y)
		e.Next = r.Exclusions
		r.Exclusions = e
	}
	return e
}
func mapRoomRemoveGeneratedExclusions(r *mapRoom) *mapRoomExclusion {
	var prev *mapRoomExclusion
	for e := r.Exclusions; e != nil; {
		next := e.Next
		if e.Kind == 1 {
			if prev != nil {
				prev.Next = next
			} else {
				r.Exclusions = next
			}
			mapRoomRelease(unsafe.Pointer(e))
		} else {
			prev = e
		}
		e = next
	}
	return nil
}
func mapRoomContainsRect(r *mapRoom, e *mapRoomExclusion) uint32 {
	return uint32(bool2int(float64(e.Min.X)+0.5 >= float64(r.Min.X) && float64(e.Max.X)-0.5 <= float64(r.Max.X) && float64(e.Min.Y)+0.5 >= float64(r.Min.Y) && float64(e.Max.Y)-0.5 <= float64(r.Max.Y)))
}

// Like the original x87 comparisons, retain the half-unit insets in double
// precision. A float32 store here can turn touching rectangles into overlaps.
func mapRoomFindExclusionOverlap(r *mapRoom, t *mapRoomExclusion) *mapRoomExclusion {
	for e := r.Exclusions; e != nil; e = e.Next {
		if e != t && float64(t.Min.X)+0.5 < float64(e.Max.X)-0.5 &&
			float64(t.Max.X)-0.5 > float64(e.Min.X)+0.5 &&
			float64(t.Min.Y)+0.5 < float64(e.Max.Y)-0.5 &&
			float64(t.Max.Y)-0.5 > float64(e.Min.Y)+0.5 {
			return e
		}
	}
	return nil
}

func mapRoomRandomPoint(r *mapRoom, scale float32, out *types.Pointf) uint32 {
	// The x87 build keeps these centers in double stack slots across RNG calls.
	cx := (float64(r.Min.X) + float64(r.Max.X)) * 0.5
	cy := (float64(r.Min.Y) + float64(r.Max.Y)) * 0.5
	hx := float32((float64(r.Max.X) - float64(r.Min.X)) * 0.5)
	hy := float32((float64(r.Max.Y) - float64(r.Min.Y)) * 0.5)
	for i := 0; i < 10; i++ {
		out.X = float32(mapRoomRandomFloat(-hx, hx)*float64(scale) + float64(cx))
		out.Y = float32(mapRoomRandomFloat(-hy, hy)*float64(scale) + float64(cy))
		if mapRoomPointExcluded(r, out) == 0 {
			return 1
		}
	}
	return 0
}
func mapRoomPointExcluded(r *mapRoom, p *types.Pointf) uint32 {
	for e := r.Exclusions; e != nil; e = e.Next {
		if !(p.X < e.Min.X || p.X > e.Max.X || p.Y < e.Min.Y || p.Y > e.Max.Y) {
			return 1
		}
	}
	return 0
}
func mapRoomRememberPoint(r *mapRoom, p *types.Pointf) int8 {
	if r.PointCount < 16 {
		for i := 0; i < int(r.PointCount); i++ {
			if mapRoomNear(r.Points[i].X, p.X) != 0 && mapRoomNear(r.Points[i].Y, p.Y) != 0 {
				return 1
			}
		}
		r.Points[r.PointCount] = *p
		r.PointCount++
	}
	return int8(r.PointCount)
}
func mapRoomNear(a, b float32) uint32 {
	return uint32(bool2int(math.Abs(float64(a)-float64(b)) < *memmap.PtrFloat64(0x581450, 10432)))
}
func mapRoomIsHall(r *mapRoom) uint32 { return uint32(bool2int(r.Kind >= 2 && r.Kind <= 5)) }
func mapRoomHallDirection(kind int32) int32 {
	switch kind {
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 3
	}
	return 0
}
func mapRoomOpposite(dir int32) int32 {
	return int32(*memmap.PtrUint32(0x587000, 254952+uintptr(dir)*4))
}
func mapRoomHallSide(kind int32) int32 {
	switch kind {
	case 2:
		return 3
	case 3:
		return 2
	case 5:
		return 1
	}
	return 0
}
func mapRoomHallOppositeSide(kind int32) int32 { return mapRoomOpposite(mapRoomHallSide(kind)) }

// Keep the double intermediates across trimming iterations; Size is stored as
// float32 on each iteration, while the moving endpoint is rounded only at exit.
func mapRoomTrimHall(r, t *mapRoom) uint32 {
	pos := r.Pos
	switch r.Kind {
	case 2, 3:
		if r.Min.X < t.Min.X || r.Max.X > t.Max.X {
			return 0
		}
		if r.Kind == 2 {
			y := float64(pos.Y)
			for r.Height > 1 {
				y += 32.526913
				r.Size.Y = float32(float64(r.Size.Y) - 32.526913)
				r.Height--
				if !(y+0.1 <= float64(t.Max.Y)) {
					break
				}
			}
			pos.Y = float32(y)
		} else {
			for r.Height > 1 {
				h := float64(r.Size.Y) - 32.526913
				r.Height--
				r.Size.Y = float32(h)
				if h+float64(pos.Y)-0.1 < float64(t.Min.Y) {
					break
				}
			}
		}
		mapRoomSetPos(r, &pos)
	case 4, 5:
		if r.Min.Y < t.Min.Y || r.Max.Y > t.Max.Y {
			return 0
		}
		if r.Kind == 5 {
			x := float64(pos.X)
			for r.Width > 1 {
				x += 32.526913
				r.Size.X = float32(float64(r.Size.X) - 32.526913)
				r.Width--
				if !(x+0.1 <= float64(t.Max.X)) {
					break
				}
			}
			pos.X = float32(x)
		} else {
			for r.Width > 1 {
				w := float64(r.Size.X) - 32.526913
				r.Width--
				r.Size.X = float32(w)
				if !(w+float64(pos.X)-0.1 >= float64(t.Min.X)) {
					break
				}
			}
		}
		mapRoomSetPos(r, &pos)
	}
	return uint32(bool2int(mapRoomOverlap(r) == nil))
}
func mapRoomHallExit(r *mapRoom, out *types.Pointf) uint32 {
	ret := mapRoomRaw(unsafe.Pointer(r))
	switch r.Kind {
	case 2, 3:
		out.X = float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5)
		if r.Kind == 2 {
			out.Y = r.Min.Y
			ret = math.Float32bits(r.Min.Y)
		} else {
			out.Y = r.Max.Y
		}
	case 4, 5:
		if r.Kind == 4 {
			out.X = r.Max.X
		} else {
			out.X = r.Min.X
		}
		out.Y = float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)
	}
	return ret
}
func mapRoomHallEntry(r *mapRoom, out *types.Pointf) uint32 {
	ret := mapRoomRaw(unsafe.Pointer(r))
	switch r.Kind {
	case 2, 3:
		out.X = float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5)
		if r.Kind == 2 {
			out.Y = r.Max.Y
			ret = math.Float32bits(r.Max.Y)
		} else {
			out.Y = r.Min.Y
		}
	case 4, 5:
		if r.Kind == 4 {
			out.X = r.Min.X
		} else {
			out.X = r.Max.X
		}
		out.Y = float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)
	}
	return ret
}
func mapRoomHallCenter(r *mapRoom, out *types.Pointf) *mapRoom {
	side := 2
	if r.Kind == 2 || r.Kind == 3 {
		side = 0
	}
	a, b := r.Neighbors[side][0], r.Neighbors[side+1][0]
	if a != nil && b != nil && a.Kind == 1 && b.Kind == 1 {
		out.X = float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5)
		out.Y = float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)
		return r
	}
	switch r.Kind {
	case 2, 3:
		out.X = float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5)
		if r.Kind == 2 {
			out.Y = float32(float64(r.Size.X)*0.5 + float64(r.Min.Y))
		} else {
			out.Y = float32(float64(r.Max.Y) - float64(r.Size.X)*0.5)
		}
	case 4, 5:
		if r.Kind == 4 {
			out.X = float32(float64(r.Max.X) - float64(r.Size.Y)*0.5)
		} else {
			out.X = float32(float64(r.Size.Y)*0.5 + float64(r.Min.X))
		}
		out.Y = float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)
	}
	return r
}
func mapRoomNewHall(kind, width, length int32) *mapRoom {
	r := mapRoomAlloc[mapRoom]()
	if r != nil {
		r.Kind = kind
		switch kind {
		case 2, 3:
			r.Width = width
			if r.Width < 1 {
				r.Width = 1
			}
			r.Height = length
		case 4, 5:
			r.Width = length
			r.Height = width
			if r.Height < 1 {
				r.Height = 1
			}
		}
		r.Size = types.Ptf(float32(float64(r.Width)*32.526913), float32(float64(r.Height)*32.526913))
	}
	return r
}
func mapRoomPrepareHall(config unsafe.Pointer, kind, width int32) *mapRoom {
	return mapRoomNewHall(kind, width, mapRoomRandomCentered(*mapRoomWord(config, 4), *mapRoomWord(config, 8)))
}
