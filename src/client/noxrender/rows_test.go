package noxrender

import (
	"image"
	"math/rand/v2"
	"testing"
)

// Merged spans must cover exactly the pixels of the spans, once each.
func TestMergeRows(t *testing.T) {
	rnd := rand.New(rand.NewPCG(3, 4))
	for iter := 0; iter < 500; iter++ {
		var spans []image.Rectangle
		// a few column ranges that change now and then, like the visible floor
		ranges := [][2]int{{0, 100}, {150, 300}}
		for y := 0; y < 200; y++ {
			if rnd.IntN(10) == 0 { // the spans of a row don't overlap
				ranges[0][0] = rnd.IntN(100)
				ranges[1][0] = 150 + rnd.IntN(150)
			}
			for _, r := range ranges {
				if r[1] > r[0] && rnd.IntN(20) != 0 {
					spans = append(spans, image.Rect(r[0], y, r[1], y+1))
				}
			}
		}
		want := make(map[image.Point]int)
		for _, s := range spans {
			for x := s.Min.X; x < s.Max.X; x++ {
				want[image.Pt(x, s.Min.Y)]++
			}
		}
		rects := MergeRows(nil, spans)
		got := make(map[image.Point]int)
		for _, r := range rects {
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					got[image.Pt(x, y)]++
				}
			}
		}
		if len(got) != len(want) {
			t.Fatalf("iteration %d: %d pixels, want %d", iter, len(got), len(want))
		}
		for p, n := range got {
			if n != 1 || want[p] != 1 {
				t.Fatalf("iteration %d: pixel %v covered %d times (want %d)", iter, p, n, want[p])
			}
		}
		if iter == 0 {
			t.Logf("%d spans into %d rectangles", len(spans), len(rects))
		}
	}
}
