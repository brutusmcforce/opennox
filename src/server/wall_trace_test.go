package server

import (
	"math"
	"testing"
	"time"

	"github.com/brutusmcforce/opennox/libs/types"
)

// Rays to points outside the map must end: a ray going left to a negative X looped forever,
// so did the straight cases with a NaN end point.
func TestMapTraceRayOutside(t *testing.T) {
	s := &Server{}
	s.Walls = serverWalls{byPos: make([]*Wall, wallsPerBucket*WallGridSize), indexY: make([]*Wall, WallGridSize)}
	nan := float32(math.NaN())
	cases := []struct{ p1, p2 types.Pointf }{
		{types.Ptf(60, 60), types.Ptf(-4, 95)},
		{types.Ptf(60, 60), types.Ptf(17, -33)},
		{types.Ptf(60, 60), types.Ptf(-18, -31)},
		{types.Ptf(60, 60), types.Ptf(-1e9, 70)},
		{types.Ptf(60, 60), types.Ptf(70, 1e9)},
		{types.Ptf(5800, 5800), types.Ptf(1e12, 5810)},
		{types.Ptf(100, 100), types.Ptf(3000, 2000)},
	}
	for _, c := range cases {
		done := make(chan bool, 1)
		go func() { done <- s.MapTraceRayAt(c.p1, c.p2, nil, nil, 0) }()
		select {
		case ok := <-done:
			if !ok {
				t.Errorf("%v -> %v: blocked without walls", c.p1, c.p2)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("%v -> %v: does not end", c.p1, c.p2)
		}
	}
	for _, p2 := range []types.Pointf{types.Ptf(65, nan), types.Ptf(nan, 65), types.Ptf(nan, nan)} {
		if s.MapTraceRayAt(types.Ptf(60, 60), p2, nil, nil, 0) {
			t.Errorf("%v: not blocked", p2)
		}
	}
}

// A diagonal ray to a far away point must not walk millions of cells.
func TestMapTraceRayFar(t *testing.T) {
	s := &Server{}
	s.Walls = serverWalls{byPos: make([]*Wall, wallsPerBucket*WallGridSize), indexY: make([]*Wall, WallGridSize)}
	start := time.Now()
	for i := 0; i < 100; i++ {
		s.MapTraceRayAt(types.Ptf(60, 60), types.Ptf(70, 1e9), nil, nil, 0)
		s.MapTraceRayAt(types.Ptf(60, 60), types.Ptf(1e9, 80), nil, nil, 0)
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("too slow: %v", d)
	}
}
