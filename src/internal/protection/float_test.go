package protection

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

// x87Oracle models the Go-hosted legacy x87 precision (53 significant bits),
// round-to-nearest-even addition, and signed-64-bit integer conversion.
func x87Oracle(old uint32, v float32) uint32 {
	x := float64(v)
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0
	}
	a := new(big.Float).SetPrec(53).SetMode(big.ToNearestEven).SetUint64(uint64(old))
	b := new(big.Float).SetPrec(53).SetFloat64(x)
	a.Add(a, b)
	n, _ := a.Int(nil)
	if !n.IsInt64() {
		return 0
	}
	return uint32(n.Int64())
}

func TestFloatValues(t *testing.T) {
	olds := []uint32{0, 1, 2, 3, 255, 256, 257, 0x7fffffff, 0x80000000, 0xffffffff}
	vals := []uint32{0, 0x80000000, 1, 0x80000001, 0x007fffff, 0x807fffff, 0x3f000000, 0xbf000000, 0x3fc00000, 0xbfc00000, 0x4f7fffff, 0x4f800000, 0xcf800000, 0x5effffff, 0x5f000000, 0xdeffffff, 0xdf000000, 0xdf000001, 0x7f800000, 0xff800000, 0x7fc00000, 0x7f800001, 0xffc00001}
	for e := -70; e <= -15; e++ {
		v := float32(-math.Ldexp(1, e))
		vals = append(vals, math.Float32bits(v), math.Float32bits(math.Nextafter32(v, float32(math.Inf(-1)))), math.Float32bits(math.Nextafter32(v, 0)))
	}
	check := func(old, bits uint32) {
		t.Helper()
		v := math.Float32frombits(bits)
		if got, want := FloatValue(v), x87Oracle(0, v); got != want {
			t.Fatalf("set bits=%08x got=%08x want=%08x", bits, got, want)
		}
		if got, want := AddFloatValue(old, v), x87Oracle(old, v); got != want {
			t.Fatalf("add old=%08x bits=%08x got=%08x want=%08x", old, bits, got, want)
		}
	}
	for _, old := range olds {
		for _, v := range vals {
			check(old, v)
		}
	}
	gen := rand.New(rand.NewSource(0x56fa40))
	for i := 0; i < 100000; i++ {
		check(gen.Uint32(), gen.Uint32())
	}
}
