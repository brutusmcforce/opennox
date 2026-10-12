package opennox

import (
	"image/color"
	"time"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/platform"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

// selFadeDur is the time it takes for the selection circle to fade from one color to the next.
const selFadeDur = 250 * time.Millisecond

// selCircle tracks the color of the selection circle.
//
// Originally, the circle picked a random color from a small palette each time it was drawn,
// which looks like a fast flicker. Instead, we smoothly fade between random colors of the same palette.
type selCircle struct {
	dr       *client.Drawable
	hostile  bool
	from, to color.Color
	start    time.Duration
}

var (
	selColors        = []color.Color{nox_color_white_2523948, nox_color_blue_2650684, nox_color_cyan_2649820}
	selColorsHostile = []color.Color{nox_color_violet_2598268, nox_color_red}
)

func (c *Client) selRandomColor(hostile bool, prev color.Color) color.Color {
	pal := selColors
	if hostile {
		pal = selColorsHostile
	}
	i := c.srv.Rand.Other.Int(0, len(pal)-1)
	if pal[i] == prev {
		// always change the color, so that the circle keeps sparkling
		i = (i + 1) % len(pal)
	}
	return pal[i]
}

func (c *Client) drawSelection(vp *noxrender.Viewport) {
	dr := legacy.Get_dword_5d4594_1096640()
	if dr == nil {
		c.sel.dr = nil
		return
	}
	s := &c.sel
	now := platform.Ticks()
	hostile := dr.ObjClass&6 != 0 && legacy.Sub_495A80(dr.NetCode32) == 0
	if dr != s.dr || hostile != s.hostile || now-s.start > 2*selFadeDur {
		s.dr, s.hostile = dr, hostile
		s.from = c.selRandomColor(hostile, nil)
		s.to = c.selRandomColor(hostile, s.from)
		s.start = now
	}
	for now-s.start >= selFadeDur {
		s.from, s.to = s.to, c.selRandomColor(hostile, s.to)
		s.start += selFadeDur
	}
	cl := lerpColor(s.from, s.to, float64(now-s.start)/float64(selFadeDur))
	if (dr.ObjClass & 0x200) == 0 {
		c.drawSelectionCursor(vp, cl, dr)
	}
}

func lerpColor(a, b color.Color, t float64) color.Color {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	lerp := func(x, y uint32) byte {
		return byte((float64(x>>8)*(1-t) + float64(y>>8)*t) + 0.5)
	}
	return noxcolor.RGB5551Color(lerp(ar, br), lerp(ag, bg), lerp(ab, bb))
}

func (c *Client) drawSelectionCursor(vp *noxrender.Viewport, cl color.Color, dr *client.Drawable) {
	if dr == nil {
		return
	}
	rad := 6
	if dr.Shape.Kind == server.ShapeKindCircle {
		rad += int(dr.Shape.Circle.R)
	} else if dr.Shape.Kind == server.ShapeKindBox {
		rad += int(dr.Shape.Box.H+dr.Shape.Box.W) / 2
	}
	cx := dr.PosVec.X + vp.Screen.Min.X - vp.World.Min.X
	cy := dr.PosVec.Y + vp.Screen.Min.Y - vp.World.Min.Y
	c.r.Data().SetColor2(cl)
	c.r.Data().SetAlphaEnabled(true)
	c.r.Data().SetAlphaEnabled(true)
	c.r.Data().SetAlpha(128)
	c.r.DrawCircle(cx, cy, rad-1, cl)
	c.r.DrawCircle(cx, cy, rad+1, cl)
	c.r.Data().SetAlphaEnabled(false)
	c.r.DrawCircle(cx, cy, rad, cl)
}
