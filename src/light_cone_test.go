package opennox

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/libs/common"
)

// A cone of light lights the light grid beyond the original 57x45 points (high resolutions).
func TestLightConeHighRes(t *testing.T) {
	legacy.InitBlobData() // the light constants
	c := &Client{}
	// a light far right and down in the grid, shining to the right
	center := image.Pt(150*common.GridStep, 100*common.GridStep)
	g := func(x, y int) image.Point { return image.Pt((x<<16)/common.GridStep, (y<<16)/common.GridStep) }
	a1 := g(center.X, center.Y)
	a2 := g(center.X+150, center.Y-60)
	a3 := g(center.X+150, center.Y+60)
	// like the client, in both orders: only one fills the triangle
	white := noxrender.RGB{R: 255, G: 255, B: 255}
	c.sub_4696B0(a1, a2, a3, center, 40<<16, white)
	c.sub_4696B0(a1, a3, a2, center, 40<<16, white)
	lit := 0
	for x := 150; x < lightGridW; x++ {
		for y := 95; y < 106; y++ {
			if c.tiles.nox_arr2_853BC0[x][y].R > 0 {
				lit++
			}
		}
	}
	require.NotZero(t, lit)
}

// Filling the light grid only resets the cells that lights changed: the result must be the whole grid.
func TestFillLightGrid(t *testing.T) {
	c := &Client{}
	check := func(cl noxrender.RGB) {
		t.Helper()
		for x := range c.tiles.nox_arr2_853BC0 {
			for y, v := range c.tiles.nox_arr2_853BC0[x] {
				if v != cl {
					t.Fatalf("cell %d,%d: %v, want %v", x, y, v, cl)
				}
			}
		}
	}
	amb := noxrender.RGB{R: 50 << 16, G: 60 << 16, B: 70 << 16}
	light := noxrender.RGB{R: 255, G: 128, B: 1}
	for i, cl := range []noxrender.RGB{amb, amb, {R: 255 << 16, G: 255 << 16, B: 255 << 16}, amb} {
		c.fillLightGrid(cl)
		check(cl)
		// lights in a few places, including the corners
		for _, p := range []image.Point{{0, 0}, {lightGridW - 1, lightGridH - 1}, {10 + i, 20}, {100, 3 * i}} {
			c.sub_4695E0(p.X, p.Y, light, 100<<8, i == 2) // darker on white
		}
		require.NotEqual(t, cl, c.tiles.nox_arr2_853BC0[10+i][20])
	}
	c.fillLightGrid(amb)
	check(amb)
}
