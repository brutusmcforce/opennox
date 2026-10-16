package opennox

import (
	"image"
	"testing"
)

func TestConfigVideoSize(t *testing.T) {
	for _, c := range []struct{ w, h, ew, eh int }{
		{1920, 1080, 1920, 1080},
		{640, 480, 640, 480},
		{0, 0, 1024, 768},
		{-5, 100000, 1024, 768},
		{639, 480, 1024, 768},
		{7680, 4320, noxMaxWidth, noxMaxHeight},
		{5120, 1440, noxMaxWidth, 1440},
	} {
		if got := configVideoSize(c.w, c.h); got != (image.Point{X: c.ew, Y: c.eh}) {
			t.Errorf("%dx%d: got %v, expected %dx%d", c.w, c.h, got, c.ew, c.eh)
		}
	}
}
