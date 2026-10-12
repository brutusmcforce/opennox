package opennox

import "testing"

func TestMulDiv(t *testing.T) {
	cases := []struct{ a, b, c, want int }{
		{3, 1, 2, 2},   // 1.5 is rounded away from zero
		{-3, 1, 2, -2}, // -1.5 too
		{3, -1, 2, -2},
		{5, 1, -2, -3}, // -2.5
		{1, 1, 3, 0},
		{2, 1, 3, 1},
		{0x10000, 45875, 0x10000, 45875},
		{1 << 30, 4, 8, 1 << 29}, // the intermediate value doesn't fit into 32 bits
		{1, 1, 0, -1},
	}
	for _, c := range cases {
		if got := muldiv(c.a, c.b, c.c); got != c.want {
			t.Errorf("muldiv(%d, %d, %d) = %d, want %d", c.a, c.b, c.c, got, c.want)
		}
	}
}

func TestParticleRand(t *testing.T) {
	var pfx partFXes
	pfx.rand.u16 = 0x1352
	pfx.rand.u32 = 0x13527438
	seen := make(map[uint16]bool)
	for i := 0; i < 100; i++ {
		seen[pfx.randU16()] = true
		if v := pfx.randI32(-5, 5); v < -5 || v >= 5 {
			t.Fatalf("randI32 out of range: %d", v)
		}
	}
	if len(seen) < 50 {
		t.Errorf("randU16 repeats too often: %d unique values out of 100", len(seen))
	}
}
