package platform

import (
	"math/rand"
	"time"
)

// New returns a platform based on the real hardware.
func New() Platform {
	p := &hardware{
		start: time.Now(),
	}
	p.RandSeedTime()
	return p
}

type hardware struct {
	start time.Time
	rand  *rand.Rand
}

func (p *hardware) Ticks() time.Duration {
	return time.Since(p.start)
}

func (*hardware) Sleep(dt time.Duration) {
	time.Sleep(dt)
}

func (*hardware) TimeSeed() int64 {
	return time.Now().UnixNano()
}

func (p *hardware) RandInt() int {
	return RandInt31(p.rand)
}

// RandInt31 returns what r.Int() returns on 32-bit platforms (the low 31 bits of r.Int63()),
// so the game draws the same random numbers on every architecture.
func RandInt31(r *rand.Rand) int {
	return int(uint32(r.Int63()) &^ (1 << 31))
}

func (p *hardware) RandSeed(seed int64) {
	p.rand = rand.New(rand.NewSource(seed))
}

func (p *hardware) RandSeedTime() {
	p.RandSeed(p.TimeSeed())
}
