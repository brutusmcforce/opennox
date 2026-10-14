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
