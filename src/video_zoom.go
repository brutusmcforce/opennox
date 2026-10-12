package opennox

import (
	"context"
	"image"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/noximage"
	"github.com/brutusmcforce/opennox/libs/platform"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// World zoom.
//
// The world is always drawn at full resolution into the world buffer (see video_uiscale.go).
// When zoomed in, the center of the world buffer is scaled up to fill the whole buffer,
// right after the world is drawn and before the GUI is composed on top of it.
// This way the zoom doesn't affect the game state, the server or the GUI.
//
// When zoomed out, the world buffer is made larger than the screen (in logical pixels), so that a larger part
// of the world is drawn, and the GPU draws it scaled down into a render target of the screen size.
// It also requires a larger view on the server, so zooming out is only possible when hosting the game.
//
// Zoom is controlled by ZoomIn and ZoomOut key bindings (mouse wheel by default, see input_bindings.go).

const (
	zoomOutMin = 0.5 // the lowest zoom, when zooming out is possible
	zoomMax    = 3.0
	zoomFactor = 1.25 // zoom change for one step
	// zoomSpeed controls how fast the zoom approaches the target value (per second).
	zoomSpeed = 15.0
)

// worldViewMax is the largest part of the world that can be drawn.
var worldViewMax = image.Pt(worldViewMaxWidth, worldViewMaxHeight)

func init() {
	noxConsole.Register(&console.Command{
		Token: "zoom",
		Help:  "zoom [percent]: get or set the world zoom (50-300)",
		Flags: console.ClientServer,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) == 0 {
				c.Printf(console.ColorLightYellow, "zoom = %.0f%%", noxClient.zoom.target*100)
				return true
			} else if len(tokens) > 1 {
				return false
			}
			v, err := strconv.ParseFloat(strings.TrimSuffix(tokens[0], "%"), 64)
			if err != nil {
				return false
			}
			noxClient.setZoom(v / 100)
			c.Printf(console.ColorLightYellow, "zoom = %.0f%%", noxClient.zoom.target*100)
			return true
		},
	})
}

type worldZoom struct {
	target float64       // requested zoom
	cur    float64       // current zoom, approaching the target
	last   time.Duration // time of the last update
}

// level returns the current zoom factor.
func (z *worldZoom) level() float64 {
	if z.cur <= 0 {
		return 1
	}
	return z.cur
}

// zoomMin returns the lowest zoom that is currently possible.
func (c *Client) zoomMin() float64 {
	if !noxflags.HasGame(noxflags.GameHost) {
		// the server only sends the part of the world that fits the screen
		return 1
	}
	if noxflags.HasGame(noxflags.GameOnline) && getServerSettings().Flags102.Has(server.ServerLimitMaxRes) {
		// the server limits the view of all players
		return 1
	}
	return zoomOutMin
}

func (c *Client) setZoom(v float64) {
	if math.IsNaN(v) {
		return
	}
	c.zoom.target = max(c.zoomMin(), min(zoomMax, v))
}

// zoomStep zooms the world in or out by one step.
func (c *Client) zoomStep(in bool) {
	cur := c.zoom.target
	if cur <= 0 {
		cur = 1
	}
	if in {
		c.setZoom(cur * zoomFactor)
	} else {
		c.setZoom(cur / zoomFactor)
	}
}

// updateZoom advances the zoom animation and makes sure the world buffer exists when it's needed.
// It must be called before drawing the world.
func (c *Client) updateZoom() {
	z := &c.zoom
	if lo := c.zoomMin(); z.target < lo {
		z.target = lo
	}
	now := platform.Ticks()
	dt := (now - z.last).Seconds()
	z.last = now
	if z.cur <= 0 || dt <= 0 || dt > 1 {
		dt = 1
	}
	if z.target < 1 || z.cur < 1 {
		// zooming out changes the size of the world view, so it's not animated
		z.cur = z.target
	}
	z.cur += (z.target - z.cur) * min(1, dt*zoomSpeed)
	if math.Abs(z.cur-z.target) < 0.002 {
		z.cur = z.target
	}
	c.updateWorldView()
}

// worldPhysSize returns the size of the screen area used for the world, in physical pixels of the render target.
func worldPhysSize() image.Point {
	if sz := uiScaleState.worldSize; sz != (image.Point{}) {
		return sz
	}
	if noxPixBuffer.img != nil {
		return noxPixBuffer.img.Size()
	}
	return image.Point{}
}

// updateWorldView creates, resizes or frees the world buffer according to the current zoom.
func (c *Client) updateWorldView() {
	st := &uiScaleState
	if st.inWorld || noxPixBuffer.img == nil {
		return
	}
	phys := worldPhysSize()
	l := c.zoom.level()
	view := phys
	if l < 1 {
		view = image.Pt(int(float64(phys.X)/l+0.5), int(float64(phys.Y)/l+0.5))
		view.X = min(view.X, max(worldViewMax.X, phys.X))
		view.Y = min(view.Y, max(worldViewMax.Y, phys.Y))
	}
	if st.worldSize == (image.Point{}) && view == phys && l == 1 {
		// no UI scaling and no zoom: the world is drawn directly
		videoFreeWorldBuffer()
	} else if st.world == nil || st.world.Size() != view {
		videoAllocWorldBuffer(view)
		if g := c.gpu; g != nil {
			g.SetTargetSize(st.world, phys)
		}
	}
	sz := videoWorldSize()
	if sz != st.viewSize {
		c.resizeWorldView(sz)
	}
	c.syncHostView(sz)
}

// syncHostView makes sure that the server sends all objects in the world view of the local player.
// The server resets the view when the player joins (e.g. on map change), so it's checked on every frame.
func (c *Client) syncHostView(sz image.Point) {
	if !noxflags.HasGame(noxflags.GameHost) {
		return
	}
	pl := noxServer.Players.ByInd(server.HostPlayerIndex)
	if pl == nil {
		return
	}
	hx, hy := uint16(sz.X/2), uint16(sz.Y/2)
	if pl.Field10 < hx || pl.Field12 < hy {
		pl.Field10, pl.Field12 = hx, hy
	}
}

// resizeWorldView updates the client state that depends on the size of the world view.
func (c *Client) resizeWorldView(sz image.Point) {
	st := &uiScaleState
	videoLog.Printf("world view: %v -> %v", st.viewSize, sz)
	st.viewSize = sz
	c.resizeTileBuf(sz)
	c.Sight.Free()
	c.Sight.Init(sz.X, sz.Y)
	vp := c.Viewport()
	vp.Screen = image.Rect(0, 0, sz.X-1, sz.Y-1)
	vp.Size = sz
	c.nox_draw_setCutSize(nox_video_getCutSize(), 0)
	if vp.World.Max != (image.Point{}) {
		// Keep the camera centered. The server only sends the camera position on ticks,
		// and the line of sight is computed from the view center, so an off-center view
		// would briefly reveal objects behind walls. The view center is stored in World.Max,
		// see nox_xxx_cliUpdateCameraPos_435600.
		prev := vp.World.Min
		vp.World.Min = vp.World.Max.Sub(vp.Size.Div(2))
		d := vp.World.Min.Sub(prev)
		c.interp.shiftCamera(d)
	}
	if noxflags.HasGame(noxflags.GameHost) {
		// the view may also be smaller now
		if pl := noxServer.Players.ByInd(server.HostPlayerIndex); pl != nil {
			pl.Field10, pl.Field12 = uint16(sz.X/2), uint16(sz.Y/2)
		}
	}
}

// zoomRect returns a part of the world buffer of a given size that is visible with the current zoom.
func (c *Client) zoomRect(sz image.Point) image.Rectangle {
	l := c.zoom.level()
	w, h := int(float64(sz.X)/l+0.5), int(float64(sz.Y)/l+0.5)
	x, y := (sz.X-w)/2, (sz.Y-h)/2
	return image.Rect(x, y, x+w, y+h)
}

// worldScreenToWindow converts a position in the world buffer to a position on the (full resolution) screen.
func (c *Client) worldScreenToWindow(p image.Point) image.Point {
	l := c.zoom.level()
	if l <= 1 {
		// the world buffer covers the whole screen
		return p
	}
	ctr := videoWorldSize().Div(2)
	return image.Point{
		X: ctr.X + int(math.Round(float64(p.X-ctr.X)*l)),
		Y: ctr.Y + int(math.Round(float64(p.Y-ctr.Y)*l)),
	}
}

// windowToWorldScreen converts a position on the (full resolution) screen to a position in the world buffer.
func (c *Client) windowToWorldScreen(p image.Point) image.Point {
	l := c.zoom.level()
	if l <= 1 {
		return p
	}
	ctr := videoWorldSize().Div(2)
	return image.Point{
		X: ctr.X + int(math.Floor(float64(p.X-ctr.X)/l)),
		Y: ctr.Y + int(math.Floor(float64(p.Y-ctr.Y)/l)),
	}
}

// applyZoom scales the visible part of the world buffer to fill the whole buffer.
func (c *Client) applyZoom(img *noximage.Image16) {
	if img == nil || c.gpu == nil || c.zoom.level() <= 1 {
		return
	}
	src := c.zoomRect(img.Size())
	if src.Empty() || src == img.Rect {
		return
	}
	c.gpu.Zoom(img, src)
}
