package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

var geometryDirectionThreshold int32 = 6

func geometrySegments(a, b *[4]int32) int32 {
	x, y, endX, endY := a[0], a[1], a[2], a[3]
	bx, by, bEndX, bEndY := b[0], b[1], b[2], b[3]
	dx, bdx, dy, bdy := endX-x, bx-bEndX, endY-y, by-bEndY
	loX, hiX := x, endX
	if dx < 0 {
		loX, hiX = endX, x
	}
	if bdx <= 0 {
		if hiX < bx || bEndX < loX {
			return 0
		}
	} else if hiX < bEndX || bx < loX {
		return 0
	}
	loY, hiY := y, endY
	if dy < 0 {
		loY, hiY = endY, y
	}
	if bdy <= 0 {
		if hiY < by || bEndY < loY {
			return 0
		}
	} else if hiY < bEndY || by < loY {
		return 0
	}
	px, py := x-bx, y-by
	cross, denom := bdy*px-bdx*py, bdx*dy-dx*bdy
	if denom <= 0 {
		if cross > 0 || cross < denom {
			return 0
		}
	} else if cross < 0 || cross > denom {
		return 0
	}
	cross = dx*py - dy*px
	if denom <= 0 {
		return int32(bool2int(cross <= 0 && cross >= denom))
	}
	return int32(bool2int(cross >= 0 && cross <= denom))
}
func geometryProjectEdge(p *[2]int32, s *[4]int32, distance float32) int32 {
	x, y := float64(s[0]), float64(s[1])
	dx := float32(float64(s[2]) - x)
	dy64 := float64(s[3]) - y
	dy := float32(dy64)
	length := float32(math.Sqrt(dy64*float64(dy) + float64(dx)*float64(dx)))
	nx64, ny64 := float64(dx)/float64(length), float64(dy)/float64(length)
	px, py := float32(float64(p[0])-x), float32(float64(p[1])-y)
	if float32(math.Abs(float64(px)*-ny64+float64(py)*nx64)) > distance {
		return 0
	}
	nx, ny := float32(nx64), float32(ny64)
	along := float64(py)*float64(ny) + float64(px)*float64(nx)
	if along > float64(length) || along < 0 {
		return 0
	}
	ox := float32(along*float64(nx) + float64(float32(x)))
	oy := float32(along*float64(ny) + float64(float32(y)))
	p[0] = floatToInt32(ox)
	p[1] = floatToInt32(oy)
	return 1
}

func geometryWallPoint(p *[2]int32, s *[8]int32) int32 {
	if p[0] < 58 || p[1] < 58 || p[0] > 5888 || p[1] > 5888 {
		return 0
	}
	x, y := p[0]/23, p[1]/23
	top, left, leftY, right, rightY := s[1]/23, s[2]/23, s[3]/23, s[4]/23, s[5]/23
	if y < top || y > s[7]/23 || x < left || x > right {
		return 0
	}
	topX := s[0] / 23
	if y > leftY {
		if x < y+left-leftY {
			return 0
		}
	} else if x < topX+top-y {
		return 0
	}
	if y > rightY {
		if x > right+rightY-y {
			return 0
		}
	} else if x > y+topX-top {
		return 0
	}
	return 1
}
func geometryWallBounds(a *[8]uint32, out *[4]int32) int32 {
	// Preserve the original unsigned comparisons and sequential aliased writes.
	y := a[1]
	if y >= a[7] {
		out[3] = int32(y)
		out[1] = int32(a[7])
	} else {
		out[1] = int32(y)
		out[3] = int32(a[7])
	}
	x := a[2]
	if x >= a[4] {
		out[2] = int32(x)
		out[0] = int32(a[4])
	} else {
		out[0] = int32(x)
		out[2] = int32(a[4])
	}
	if out[0] < 0 {
		out[0] = 0
	}
	if out[1] < 0 {
		out[1] = 0
	}
	if out[2] >= 5888 {
		out[2] = 5887
	}
	if out[3] >= 5888 {
		out[3] = 5887
	}
	return 0
}
func geometryRectInt(p *[2]int32, r *[4]int32) int32 {
	return int32(bool2int(p[0] >= r[0] && p[0] <= r[2] && p[1] >= r[1] && p[1] <= r[3]))
}
func geometryRectFloat(p *types.Pointf, r *[4]float32) int32 {
	return int32(bool2int(p.X >= r[0] && p.X <= r[2] && p.Y >= r[1] && p.Y <= r[3]))
}
func geometryMapCoordinates(p, out *types.Pointf) int32 {
	if p == nil || out == nil {
		return 0
	}
	if p.X <= 80.5 {
		p.X = 82.5
	}
	if p.Y <= 80.5 {
		p.Y = 81.5
	}
	if p.X >= 5853.5 {
		p.X = 5851.5
	}
	if p.Y >= 5853.5 {
		p.Y = 5852.5
	}
	out.X = float32((float64(p.X) - 1 - float64(p.Y)) * 0.70710677)
	out.Y = float32((float64(p.Y) + float64(p.X) - 5912) * 0.70710677)
	return 1
}
func geometryDirectionAngle(p *[2]uint32) int32 {
	return int32(memmap.Uint32(0x587000, uintptr(uint32(230072)+4*(p[0]+3*p[1]))))
}
func geometryIndexedDirection(index int32, out *[2]int32) int32 {
	x := memmap.Int32(0x587000, uintptr(192088+8*index))
	if x <= geometryDirectionThreshold {
		out[0] = int32(bool2int(x >= -geometryDirectionThreshold)) - 1
	} else {
		out[0] = 1
	}
	y := memmap.Int32(0x587000, uintptr(192092+8*index))
	result := geometryDirectionThreshold
	if y <= geometryDirectionThreshold {
		result = -geometryDirectionThreshold
		if y >= -geometryDirectionThreshold {
			out[1] = 0
		} else {
			out[1] = -1
		}
	} else {
		out[1] = 1
	}
	return result
}
func geometryDirection4Angle(index int32) int32 {
	return int32(memmap.Uint32(0x587000, uintptr(230056+4*(index%9))))
}
func geometryDirection4Index(index int32) int32 {
	var out [2]int32
	geometryIndexedDirection(index, &out)
	return out[1] + out[0] + 2*out[1] + 4
}
func geometryVectorAngle(p *types.Pointf) int32 {
	v := float32((math.Atan2(float64(p.Y), float64(p.X))+6.2831855)*40.743664 + 0.5)
	result := floatToInt32(v)
	if result < 0 {
		result += int32(uint32(255-result) >> 8 << 8)
	}
	if result >= 256 {
		result += -256 * int32(uint32(result)>>8)
	}
	return result
}
func geometryNormalize(p *types.Pointf) {
	length := float64(float32(math.Sqrt(float64(p.X)*float64(p.X) + float64(p.Y)*float64(p.Y))))
	p.X = float32(float64(p.X) / length)
	p.Y = float32(float64(p.Y) / length)
}
func geometryProjectPositive(origin *types.Pointf, lo, hi float32, clampLo, clampHi int32, p, out *types.Pointf) int32 {
	v := (float64(p.X) - float64(origin.X) + float64(p.Y) - float64(origin.Y)) * 0.70709997 * 0.70709997
	if v < float64(lo) {
		if clampLo == 0 {
			return 0
		}
		v = float64(lo)
	}
	if v > float64(hi) {
		if clampHi == 0 {
			return 0
		}
		v = float64(hi)
	}
	out.X = float32(v + float64(origin.X))
	out.Y = float32(v + float64(origin.Y))
	return 1
}
func geometryProjectNegative(origin *types.Pointf, lo, hi float32, clampLo, clampHi int32, p, out *types.Pointf) int32 {
	v := 23 - ((float64(p.Y)-float64(origin.Y))*0.70709997-(float64(p.X)-(float64(origin.X)+23))*0.70709997)*0.70709997
	if v < float64(lo) {
		if clampLo == 0 {
			return 0
		}
		v = float64(lo)
	}
	if v > float64(hi) {
		if clampHi == 0 {
			return 0
		}
		v = float64(hi)
	}
	out.X = float32(v + float64(origin.X))
	out.Y = float32(23 - v + float64(origin.Y))
	return 1
}
func geometryQuadrant(p, origin *types.Pointf) int32 {
	x, y := float64(p.X)-float64(origin.X), float64(p.Y)-float64(origin.Y)
	if float64(x) <= 16.263456 {
		if y <= 0 {
			return 1
		}
		return 8
	}
	if y <= 0 {
		return 2
	}
	return 4
}

// The C calculation keeps its products wide until storing each corner. The
// existing ShapeBox.Calc rounds each product earlier and is not interchangeable.
func geometryShapeBox(s *server.Shape) {
	const mul = float32(0.35354999)
	px, py := float64(s.Box.W)*float64(mul), float64(s.Box.H)*float64(mul)
	s.Box.LeftTop = float32(-px + py)
	s.Box.LeftTop2 = s.Box.LeftTop
	s.Box.LeftBottom = float32(-px - py)
	s.Box.LeftBottom2 = s.Box.LeftBottom
	s.Box.RightTop = float32(px + py)
	s.Box.RightTop2 = s.Box.RightTop
	s.Box.RightBottom = float32(px - py)
	s.Box.RightBottom2 = s.Box.RightBottom
}
