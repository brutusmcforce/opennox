package protection

import "math"

// protectionFloatInt preserves the 386 x87 signed-64-bit conversion followed by
// its low 32 bits. Nonfinite and out-of-range inputs yield the low bits of x87's
// integer-indefinite value, rather than relying on Go's target-specific casts.
func protectionFloatInt(x float64) uint32 {
	if math.IsNaN(x) || x < -0x1p63 || x >= 0x1p63 {
		return 0
	}
	return uint32(int64(x))
}

func FloatValue(v float32) uint32 { return protectionFloatInt(float64(v)) }

// AddFloatValue adds before truncation. The Go-hosted legacy C process uses x87
// double precision (53 significant bits), matching float64 round-to-nearest-even
// here. A standalone C process's default extended precision is different.
func AddFloatValue(old uint32, v float32) uint32 {
	return protectionFloatInt(float64(old) + float64(v))
}
