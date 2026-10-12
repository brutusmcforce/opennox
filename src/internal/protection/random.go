package protection

const (
	randomScale     float64 = 0x1p-32
	randomC0        float64 = 5115
	randomC1        float64 = 1776
	randomC2        float64 = 1492
	randomC3        float64 = 2111111111
	randomMaxUint32         = ^uint32(0)
)

// Random is the legacy protection floating RNG. Production only seeds it via
// Seed, which keeps all state finite; deliberately injected nonfinite states
// are outside its compatibility contract.
type Random struct {
	State          [5]float64
	Span, Min, Max uint32
}

// Seed initializes the five-word xorshift state and advances it 19 times.
func (r *Random) Seed(seed uint32) {
	if seed == 0 {
		seed = randomMaxUint32
	}
	for i := range r.State {
		v := (((seed << 13) ^ seed) >> 17) ^ (seed << 13) ^ seed
		seed = 32*v ^ v
		r.State[i] = float64(seed) * randomScale
	}
	for i := 0; i < 19; i++ {
		r.Next()
	}
	r.Min, r.Max, r.Span = 0, 99, 100
}

// Next advances the floating state. The original C evaluates floor(v), then
// discards that result and stores v-v, so floor intentionally has no Go call.
func (r *Random) Next() float64 {
	r.State[3] = r.State[2]
	r.State[2] = r.State[1]
	r.State[1] = r.State[0]

	// Keep every product and addition at float64 precision and in C's
	// left-associative order. This is PC53-compatible and prevents FMA changes.
	p0 := float64(r.State[0] * randomC0)
	p1 := float64(r.State[2] * randomC1)
	p2 := float64(r.State[3] * randomC2)
	p3 := float64(r.State[3] * randomC3)
	v := float64(p0 + p1)
	v = float64(v + p2)
	v = float64(v + p3)
	v = float64(v + r.State[4])

	r.State[0] = float64(v - v)
	r.State[4] = float64(v * randomScale)
	return r.State[0]
}

// Range advances the state and returns the legacy inclusive range draw.
func (r *Random) Range(min, max uint32) uint32 {
	r.Max = max
	r.Min = min
	r.Span = max - min + 1
	v := protectionFloatInt(float64(r.Next()) * float64(r.Span))
	if v < r.Span {
		return r.Min + v
	}
	return r.Span + r.Min
}

// Draw is the protection key draw, equivalent to Range(1, UINT_MAX).
func (r *Random) Draw() uint32 {
	return r.Range(1, randomMaxUint32)
}
