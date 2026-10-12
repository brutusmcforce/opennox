package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"math"
)

func effectCurveColor(color int) int { *effectMapped(1316980) = uint32(color); return color }
func effectCurveGlow(enabled, color, size int, intensity int8) int8 {
	*effectMapped(1316988) = uint32(enabled)
	*effectMapped(1316984) = uint32(color)
	*effectMapped(1316992) = uint32(size)
	*memmap.PtrUint8(0x5D4594, 1316996) = byte(intensity)
	return intensity
}
func effectCurveCoefficients(points [4]image.Point) (x, y [4]int32) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			coef := *memmap.PtrInt32(0x581450, 9872+uintptr(16*i+4*j))
			x[i] += int32(points[j].X) * coef
			y[i] += int32(points[j].Y) * coef
		}
	}
	return
}

// The rasterizer uses forward differences. Two named coefficient words are
// separate from the mapped matrix, just as they were in C; do not substitute
// an idealized Hermite rasterizer without an explicit rendering change.
func effectCurveRaster(points [4]image.Point, steps, thick int) {
	x, y := effectCurveCoefficients(points)
	h := 1 / float64(steps)
	*memmap.PtrFloat32(0x587000, 180484) = float32(h)
	h2, h3 := float32(h*h), float32(h*h*h)
	dword_587000_180480 = uint32(math.Float32bits(h2))
	dword_587000_180476 = uint32(math.Float32bits(h3))
	*memmap.PtrFloat32(0x587000, 180496) = h2 + h2
	*memmap.PtrFloat32(0x587000, 180492) = float32(float64(h3) * 6)
	*memmap.PtrFloat32(0x587000, 180508) = float32(float64(h3) * 6)
	var dx, dy [4]float32
	for i := 0; i < 4; i++ {
		var m [4]float32
		for j := range m {
			m[j] = *memmap.PtrFloat32(0x587000, 180460+uintptr(16*i+4*j))
		}
		dx[i] = float32(float64(x[0])*float64(m[0]) + float64(x[1])*float64(m[1]) + float64(x[2])*float64(m[2]) + float64(x[3])*float64(m[3]))
		// C stores only the first two Y coefficients in float variables; their
		// products and sum round before the remaining double terms are added.
		first := float32(float32(y[0])*m[0]) + float32(float32(y[1])*m[1])
		dy[i] = float32(float64(first) + float64(y[2])*float64(m[2]) + float64(y[3])*float64(m[3]))
	}
	previous := points[0]
	for i := 0; i < steps; i++ {
		dx[0] += dx[1]
		dx[1] = dx[2] + dx[1]
		dx[2] = dx[3] + dx[2]
		dy[0] += dy[1]
		dy[1] = dy[2] + dy[1]
		dy[2] = dy[3] + dy[2]
		next := image.Pt(effectFloatInt(dx[0]), effectFloatInt(dy[0]))
		r := GetClient().R2()
		r.AddPoint(previous)
		r.AddPoint(next)
		effectColor(*effectMapped(1316980))
		r.DrawLineFromPoints(r.Data().Color2())
		if thick != 0 {
			// Both branches of the old orientation comparison compared the same
			// absolute expression, so these offsets are always horizontal.
			for _, offset := range []int{-1, 1} {
				effectLine(next.Add(image.Pt(offset, 0)), previous.Add(image.Pt(offset, 0)))
			}
		}
		if *effectMapped(1316988) != 0 {
			r.AddPoint(previous)
			r.AddPoint(next)
			nox_draw_set54RGB32_434040(int(*effectMapped(1316984)))
			r.Data().SetField262(int(*effectMapped(1316992)))
			r.DrawParticles49ED80(int(*memmap.PtrUint8(0x5D4594, 1316996)))
		}
		previous = next
	}
}
func effectCurveSegments(points [4]image.Point, steps int, shift float32, segment func(image.Point, image.Point)) {
	x, y := effectCurveCoefficients(points)
	if steps <= 0 {
		return
	}
	h := 1 / float64(steps)
	offset := h * float64(shift)
	increment := float32(offset + h)
	var xf, yf [4]float32
	for i := range xf {
		xf[i] = float32(x[i])
		yf[i] = float32(y[i])
	}
	var t float32
	previous := points[0]
	for i := 0; i < steps; i++ {
		// The C compiler keeps t, t² and t³ in wider registers for X,
		// then spills them to float32 across the first converter call.
		// Y uses those spilled values, but wider products and accumulation.
		nextT := float64(increment) + float64(t)
		if nextT > 1 {
			nextT = 1
		}
		squared := nextT * nextT
		cubed := squared * nextT
		px := float32(float64(xf[1])*squared + float64(xf[0])*cubed + float64(xf[2])*nextT + float64(xf[3]))
		py := float32(float64(yf[1])*float64(float32(squared)) + float64(yf[0])*float64(float32(cubed)) + float64(yf[2])*float64(float32(nextT)) + float64(yf[3]))
		next := image.Pt(effectFloatInt(px), effectFloatInt(py))
		t = float32(nextT) - float32(offset)
		segment(previous, next)
		previous = next
	}
}
