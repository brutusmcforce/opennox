package protection

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

func randomReferenceSeed(seed uint32) [5]float64 {
	x := uint64(seed)
	if x == 0 {
		x = 0xffffffff
	}
	var out [5]float64
	for i := range out {
		x = (x ^ (x << 13)) & 0xffffffff
		x = (x ^ (x >> 17)) & 0xffffffff
		x = (x ^ (x << 5)) & 0xffffffff
		out[i] = math.Ldexp(float64(x), -32)
	}
	for i := 0; i < 19; i++ {
		randomReferenceNext(&out)
	}
	return out
}

// Reference operations round to 53 bits using arbitrary-precision arithmetic.
// Every state store then rounds to binary64, including eventual underflow.
func randomReferenceNext(s *[5]float64) {
	s[3], s[2], s[1] = s[2], s[1], s[0]
	f := func(v float64) *big.Float { return new(big.Float).SetPrec(53).SetMode(big.ToNearestEven).SetFloat64(v) }
	acc := f(s[0])
	acc.Mul(acc, f(5115))
	for _, term := range [][2]float64{{s[2], 1776}, {s[3], 1492}, {s[3], 2111111111}} {
		p := f(term[0])
		p.Mul(p, f(term[1]))
		acc.Add(acc, p)
	}
	acc.Add(acc, f(s[4]))
	// Reachable seeded state is finite and nonnegative, so v-v is positive zero.
	s[0] = 0
	acc.Mul(acc, f(0x1p-32))
	s[4], _ = acc.Float64()
}

func TestProtectionRandom(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56ff00))
	seeds := []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff}
	for i := 0; i < 500; i++ {
		seeds = append(seeds, gen.Uint32())
	}
	pairs := [][2]uint32{{0, 0}, {0, 99}, {1, 0xffffffff}, {0xffffffff, 0}, {1, 0}, {0, 0xffffffff}, {0x80000000, 0x7fffffff}, {0x7fffffff, 0x80000000}}
	for _, seed := range seeds {
		var got Random
		got.Seed(seed)
		state := randomReferenceSeed(seed)
		check := func(step int) {
			t.Helper()
			for i := range state {
				if math.Float64bits(got.State[i]) != math.Float64bits(state[i]) {
					t.Fatalf("seed=%08x step=%d index=%d got=%016x want=%016x", seed, step, i, math.Float64bits(got.State[i]), math.Float64bits(state[i]))
				}
			}
		}
		check(-1)
		if got.Span != 100 || got.Min != 0 || got.Max != 99 {
			t.Fatal("seed range", got)
		}
		for step := 0; step < 80; step++ {
			randomReferenceNext(&state)
			switch step % 3 {
			case 0:
				span, min, max := got.Span, got.Min, got.Max
				if math.Float64bits(got.Next()) != 0 {
					t.Fatal("nonzero legacy draw")
				}
				if got.Span != span || got.Min != min || got.Max != max {
					t.Fatal("Next changed range")
				}
			case 1:
				p := pairs[step%len(pairs)]
				if got.Range(p[0], p[1]) != p[0] || got.Min != p[0] || got.Max != p[1] || got.Span != uint32(uint64(p[1])+1-uint64(p[0])) {
					t.Fatal("range", p, got)
				}
			case 2:
				if got.Draw() != 1 || got.Span != 0xffffffff || got.Min != 1 || got.Max != 0xffffffff {
					t.Fatal("draw", got)
				}
			}
			check(step)
		}
		if got.State != ([5]float64{}) {
			t.Fatal("expected eventual underflow to zero", got.State)
		}
	}
}
