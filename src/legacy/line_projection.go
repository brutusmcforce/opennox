package legacy

import (
	"math"

	"github.com/brutusmcforce/opennox/libs/types"
)

// projectionValue preserves x87 NaN payloads even with GO386=softfloat.
// Finite operations use the original PC53 precision; round32 marks C spills.
type projectionValue float64

func projectionLoad(v float32) projectionValue {
	b := math.Float32bits(v)
	if b&0x7fffffff > 0x7f800000 {
		return projectionValue(math.Float64frombits(uint64(b&0x80000000)<<32 | 0x7ff8000000000000 | uint64(b&0x7fffff)<<29))
	}
	return projectionValue(v)
}
func (v projectionValue) store() float32 {
	b := math.Float64bits(float64(v))
	if math.IsNaN(float64(v)) {
		return math.Float32frombits(uint32(b>>32)&0x80000000 | 0x7fc00000 | uint32(b>>29)&0x7fffff)
	}
	return float32(v)
}
func (v projectionValue) round32() projectionValue { return projectionLoad(v.store()) }

func projectionResult(a, b, result projectionValue) projectionValue {
	an, bn := math.IsNaN(float64(a)), math.IsNaN(float64(b))
	if an && bn {
		// x87 chooses the larger significand, then the positive sign on a tie.
		ab, bb := math.Float64bits(float64(a)), math.Float64bits(float64(b))
		ap, bp := ab&0x000fffffffffffff, bb&0x000fffffffffffff
		if bp > ap || bp == ap && bb < ab {
			return b
		}
		return a
	}
	if an {
		return a
	}
	if bn {
		return b
	}
	if math.IsNaN(float64(result)) {
		return projectionValue(math.Float64frombits(0xfff8000000000000))
	}
	return result
}
func (a projectionValue) add(b projectionValue) projectionValue { return projectionResult(a, b, a+b) }
func (a projectionValue) sub(b projectionValue) projectionValue { return projectionResult(a, b, a-b) }
func (a projectionValue) mul(b projectionValue) projectionValue { return projectionResult(a, b, a*b) }
func (a projectionValue) div(b projectionValue) projectionValue { return projectionResult(a, b, a/b) }

func projectionBounds(line *[4]float32) (minX, maxX, minY, maxY projectionValue) {
	x1, x2 := projectionLoad(line[0]), projectionLoad(line[2])
	y1, y2 := projectionLoad(line[1]), projectionLoad(line[3])
	// These ordered branches deliberately preserve the C handling of NaNs.
	if x1 >= x2 {
		minX, maxX = x2, x1
	} else {
		minX, maxX = x1, x2
	}
	if y1 >= y2 {
		minY, maxY = y2, y1
	} else {
		minY, maxY = y1, y2
	}
	return
}

func projectLineClamped(line *[4]float32, point, out *types.Pointf, length float32) {
	x1, y1 := projectionLoad(line[0]), projectionLoad(line[1])
	dx, dy := projectionLoad(line[2]).sub(x1), projectionLoad(line[3]).sub(y1)
	dot := projectionLoad(point.Y).sub(y1).mul(dy).add(projectionLoad(point.X).sub(x1).mul(dx))
	denom := projectionLoad(length).mul(projectionLoad(length))
	// The X quotient spills to float32, but its sum with old X stays wider.
	x := dx.mul(dot).div(denom).round32().add(x1)
	out.X = x.store()
	// Reload Y after storing X: output may alias the line's coordinates.
	y := dot.mul(dy).div(denom).add(projectionLoad(line[1])).round32()
	out.Y = y.store()
	minX, maxX, minY, maxY := projectionBounds(line)
	if x >= minX {
		if x > maxX {
			out.X = maxX.store()
		}
	} else {
		out.X = minX.store()
	}
	if minY <= y {
		if y > maxY {
			out.Y = maxY.store()
		}
	} else {
		out.Y = minY.store()
	}
}

func projectLine(line *[4]float32, point, out *types.Pointf) bool {
	x1, y1 := projectionLoad(line[0]), projectionLoad(line[1])
	dx, dy := projectionLoad(line[2]).sub(x1), projectionLoad(line[3]).sub(y1)
	denom := dy.mul(dy).add(dx.mul(dx)).round32()
	dot := projectionLoad(point.Y).sub(y1).mul(dy).add(projectionLoad(point.X).sub(x1).mul(dx))
	x := dx.mul(dot).div(denom).add(x1).round32()
	out.X = x.store()
	// Only Y uses the spilled dot product. Its final sum stays wider for bounds.
	roundedDot := dot.round32()
	y := dy.mul(roundedDot).div(denom).add(projectionLoad(line[1]))
	out.Y = y.store()
	minX, maxX, minY, maxY := projectionBounds(line)
	return minX <= x && x <= maxX && y >= minY && y <= maxY
}
