package opennox

import (
	"math/rand"
	"sync"
	"time"

	"github.com/brutusmcforce/opennox/libs/platform"
)

func newPlayformE2E() *platformE2E {
	p := &platformE2E{}
	p.RandSeed(1234)
	return p
}

type platformE2E struct {
	mu  sync.Mutex
	t   time.Duration
	rnd *rand.Rand
}

func (p *platformE2E) ticks() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.t
}

func (p *platformE2E) tick(dt time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t += dt
}

func (p *platformE2E) Ticks() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.t
}

func (p *platformE2E) Sleep(dt time.Duration) {
	if e2e.slow != 0 {
		// Real-time pace (multiplayer tests): waits for the network take real time. They don't advance the
		// fake clock, which counts frames: a 30ms network wait would otherwise skip the script ahead.
		time.Sleep(dt)
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t += dt
}

func (p *platformE2E) TimeSeed() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return int64(p.t)
}

func (p *platformE2E) RandInt() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return platform.RandInt31(p.rnd)
}

func (p *platformE2E) RandSeed(seed int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rnd = rand.New(rand.NewSource(seed))
}

func (p *platformE2E) RandSeedTime() {
	p.RandSeed(p.TimeSeed())
}
