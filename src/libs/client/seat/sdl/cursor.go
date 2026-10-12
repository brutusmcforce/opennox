package sdl

import (
	"errors"
	"fmt"
	"image"
	"image/draw"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/env"
)

var _ seat.CursorScreen = &Window{}

type cursor struct {
	c *sdl.Cursor
}

func (c *cursor) Free() {
	if c.c != nil {
		sdl.FreeCursor(c.c)
		c.c = nil
	}
}

// NewCursor implements seat.CursorScreen.
func (win *Window) NewCursor(img image.Image, hot image.Point) (seat.Cursor, error) {
	b := img.Bounds()
	if b.Empty() {
		return nil, errors.New("empty cursor image")
	}
	// ABGR8888 is R, G, B, A in memory on little endian, same as NRGBA
	surf, err := sdl.CreateRGBSurfaceWithFormat(0, int32(b.Dx()), int32(b.Dy()), 32, uint32(sdl.PIXELFORMAT_ABGR8888))
	if err != nil {
		return nil, fmt.Errorf("cannot create cursor surface: %w", err)
	}
	defer surf.Free()
	dst := &image.NRGBA{Pix: surf.Pixels(), Stride: int(surf.Pitch), Rect: image.Rect(0, 0, b.Dx(), b.Dy())}
	draw.Draw(dst, dst.Rect, img, b.Min, draw.Src)
	c := sdl.CreateColorCursor(surf, int32(hot.X), int32(hot.Y))
	if c == nil {
		return nil, fmt.Errorf("cannot create cursor: %w", sdl.GetError())
	}
	return &cursor{c: c}, nil
}

// SetCursor implements seat.CursorScreen.
func (win *Window) SetCursor(c seat.Cursor) {
	cur, _ := c.(*cursor)
	if cur != nil && cur.c == nil {
		cur = nil
	}
	win.cursor = cur
	if cur != nil {
		sdl.SetCursor(cur.c)
	}
	win.updateCursor()
}

// updateCursor shows or hides the system cursor.
func (win *Window) updateCursor() {
	show := win.cursor != nil
	if win.cursor == nil && win.mode == seat.Windowed && (env.IsDevMode() || env.IsE2E()) {
		// show the default cursor for debugging
		sdl.SetCursor(sdl.GetDefaultCursor())
		show = true
	}
	if show {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}
