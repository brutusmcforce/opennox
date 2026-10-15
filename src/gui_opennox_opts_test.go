//go:build !server

package opennox

import "testing"

// Values are taken from window layouts of a fresh start at 200% (1280x720) and 250% (1024x576) UI scale.
func TestRelayoutShift(t *testing.T) {
	cases := []struct {
		name           string
		off, sz        int
		old, cur, want int
	}{
		{"left edge", 0, 166, 1280, 1024, 0},
		{"near left", 121, 32, 1280, 1024, 0},
		{"right edge", 1230, 50, 1280, 1024, -256},
		{"centered", 320, 640, 1280, 1024, -128},
		{"near center", 480, 2, 1280, 1024, -128},
		{"top", 50, 300, 720, 576, 0},
		{"spell bar stack", 428, 250, 720, 576, -144},
		{"bottom edge", 646, 74, 720, 576, -144},
		{"grow", 646, 74, 576, 720, 144},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := relayoutShift(c.off, c.sz, c.old, c.cur); got != c.want {
				t.Errorf("relayoutShift(%d, %d, %d, %d) = %d, want %d", c.off, c.sz, c.old, c.cur, got, c.want)
			}
		})
	}
}

func TestVideoSliderIndex(t *testing.T) {
	s := &videoSlider{Choices: []int{100, 125, 150}}
	for _, c := range []struct{ v, want int }{
		{100, 0},
		{125, 1},
		{150, 2},
		{110, 0}, // not in the list: the closest one
		{140, 2},
		{50, 0},
		{300, 2},
	} {
		if got := s.index(c.v); got != c.want {
			t.Errorf("index(%d) = %d, want %d", c.v, got, c.want)
		}
	}
	// slider positions run over twice the steps; a position selects the closest choice
	for _, c := range []struct{ pos, want int }{{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 2}} {
		if got := sliderIndex(c.pos); got != c.want {
			t.Errorf("sliderIndex(%d) = %d, want %d", c.pos, got, c.want)
		}
	}
	for i := range s.Choices {
		if got := sliderIndex(sliderPos(i)); got != i {
			t.Errorf("sliderIndex(sliderPos(%d)) = %d", i, got)
		}
	}
}
