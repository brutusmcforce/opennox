package opennox

import (
	"context"
	"image"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/noximage"
	"github.com/opennox/libs/platform"
)

// World zoom.
//
// The world is always drawn at full resolution into the world buffer (see video_uiscale.go).
// When zoomed in, the center of the world buffer is scaled up to fill the whole buffer,
// right after the world is drawn and before the GUI is composed on top of it.
// This way the zoom doesn't affect the game state, the server or the GUI.
//
// Zoom is controlled by ZoomIn and ZoomOut key bindings (mouse wheel by default, see input_bindings.go).
//
// Zooming out further than 100% would require drawing a larger part of the world,
// which the server may not send to the client, so it's not supported.

const (
	zoomMin    = 1.0
	zoomMax    = 3.0
	zoomFactor = 1.25 // zoom change for one step
	// zoomSpeed controls how fast the zoom approaches the target value (per second).
	zoomSpeed = 15.0
)

func init() {
	noxConsole.Register(&console.Command{
		Token: "zoom",
		Help:  "zoom [percent]: get or set the world zoom (100-300)",
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
	if z.cur < zoomMin {
		return zoomMin
	}
	return z.cur
}

func (c *Client) setZoom(v float64) {
	if math.IsNaN(v) {
		return
	}
	c.zoom.target = max(zoomMin, min(zoomMax, v))
}

// zoomStep zooms the world in or out by one step.
func (c *Client) zoomStep(in bool) {
	if in {
		c.setZoom(max(c.zoom.target, zoomMin) * zoomFactor)
	} else {
		c.setZoom(max(c.zoom.target, zoomMin) / zoomFactor)
	}
}

// updateZoom advances the zoom animation and makes sure the world buffer exists when it's needed.
// It must be called before drawing the world.
func (c *Client) updateZoom() {
	z := &c.zoom
	if z.target < zoomMin {
		z.target = zoomMin
	}
	now := platform.Ticks()
	dt := (now - z.last).Seconds()
	z.last = now
	if z.cur < zoomMin || dt <= 0 || dt > 1 {
		dt = 1
	}
	z.cur += (z.target - z.cur) * min(1, dt*zoomSpeed)
	if math.Abs(z.cur-z.target) < 0.002 {
		z.cur = z.target
	}
	videoNeedWorldBuffer(z.cur > zoomMin)
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
	if l == zoomMin {
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
	if l == zoomMin {
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
	if img == nil || c.gpu == nil || c.zoom.level() == zoomMin {
		return
	}
	src := c.zoomRect(img.Size())
	if src.Empty() || src == img.Rect {
		return
	}
	c.gpu.Zoom(img, src)
}
