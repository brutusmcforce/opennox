package opennox

import (
	"encoding/binary"
	"image"
	"math"

	"golang.org/x/image/draw"

	"github.com/brutusmcforce/opennox/libs/client/seat"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

// System cursor.
//
// The original game draws the cursor as a sprite on every frame, so the cursor lags behind the mouse
// by at least a frame. Instead, the cursor sprites are given to the system as cursor images,
// and the system draws the cursor at the current mouse position.
//
// The game still decides which cursor to show: nox_video_cursorDrawImpl_477A30 runs as usual,
// but drawCursorImage only records the cursor sprite. The system cursor is updated when the frame is presented.
// If the seat doesn't support system cursors, or creating one fails, the cursor is drawn by the game as before.

// sysCursorKey identifies a cursor image.
type sysCursorKey struct {
	img   *noxrender.Image
	mul   [3]uint8 // color multiplier
	scale int      // in percents
}

type sysCursorState struct {
	scr     seat.CursorScreen // nil if system cursors are not supported or failed
	view    image.Rectangle   // the screen area of the presented image
	cache   map[sysCursorKey]seat.Cursor
	want    sysCursorKey // cursor recorded for the current frame
	hasWant bool
	cur     seat.Cursor // current cursor; nil if hidden
	missed  int         // frames presented without a cursor
}

// initSysCursor enables system cursors, if the seat supports them.
func (c *Client) initSysCursor(s seat.Seat) {
	st := &c.sysCursor
	st.scr, _ = s.(seat.CursorScreen)
	st.cache = make(map[sysCursorKey]seat.Cursor)
	if st.scr != nil {
		c.Win.OnViewResize(func(view image.Rectangle) {
			st.view = view
		})
	}
}

// freeSysCursor frees all system cursors.
func (c *Client) freeSysCursor() {
	st := &c.sysCursor
	if st.scr != nil {
		st.scr.SetCursor(nil)
	}
	for _, cur := range st.cache {
		cur.Free()
	}
	*st = sysCursorState{}
}

// disableSysCursor switches back to the cursor drawn by the game.
func (c *Client) disableSysCursor() {
	c.freeSysCursor()
	videoLog.Printf("system cursor disabled, the game will draw the cursor")
}

// hasSysCursor checks if the system cursor is used instead of drawing the cursor.
func (c *Client) hasSysCursor() bool {
	return c.sysCursor.scr != nil
}

// drawCursorImage draws a cursor sprite at a given position, or records it for the system cursor.
func (c *Client) drawCursorImage(img *noxrender.Image, pos image.Point) {
	if !dword_5d4594_1193672 || img == nil {
		return // cursor is hidden
	}
	if !c.hasSysCursor() {
		c.r.DrawImageAt(img, pos)
		return
	}
	mul := [3]uint8{255, 255, 255}
	if d := c.r.Data(); d.Multiply14() && !d.Flag16() {
		cm := d.ColorMultA()
		mul = [3]uint8{uint8(cm.R), uint8(cm.G), uint8(cm.B)}
	}
	st := &c.sysCursor
	st.want = sysCursorKey{img: img, mul: mul, scale: int(math.Round(c.sysCursorScale() * 100))}
	st.hasWant = true
}

// sysCursorScale returns the scale of the cursor image, so that it has the same size as the cursor drawn by the game.
func (c *Client) sysCursorScale() float64 {
	if uiScaleState.world != nil {
		// the game draws the cursor at the screen resolution
		return 1
	}
	view := c.sysCursor.view
	if noxPixBuffer.img == nil || view.Empty() {
		return 1
	}
	sz := noxPixBuffer.img.Size()
	if sz.X <= 0 || sz.Y <= 0 {
		return 1
	}
	return max(1, min(float64(view.Dx())/float64(sz.X), float64(view.Dy())/float64(sz.Y)))
}

// updateSysCursor shows the cursor recorded for the frame. It must be called when the frame is presented.
func (c *Client) updateSysCursor() {
	st := &c.sysCursor
	if st.scr == nil {
		return
	}
	var cur seat.Cursor
	if st.hasWant {
		st.missed = 0
		cur = c.sysCursorFor(st.want)
		if st.scr == nil {
			return // failed and disabled
		}
	} else {
		// some frames are presented twice, so hide the cursor only if the game stopped drawing it
		st.missed++
		if st.missed < 3 {
			return
		}
	}
	st.hasWant = false
	if cur != st.cur {
		st.cur = cur
		st.scr.SetCursor(cur)
	}
}

// sysCursorFor returns a system cursor for a given cursor image, creating it if necessary.
func (c *Client) sysCursorFor(key sysCursorKey) seat.Cursor {
	st := &c.sysCursor
	if cur, ok := st.cache[key]; ok {
		return cur
	}
	img, hot, err := cursorImage(key)
	if err != nil {
		videoLog.Printf("cannot decode cursor image: %v", err)
		st.cache[key] = nil
		return nil
	}
	cur, err := st.scr.NewCursor(img, hot)
	if err != nil {
		videoLog.Printf("cannot create system cursor: %v", err)
		c.disableSysCursor()
		return nil
	}
	st.cache[key] = cur
	return cur
}

// cursorImage converts a cursor sprite to an image. It returns the image and the position of the pointer in it.
func cursorImage(key sysCursorKey) (image.Image, image.Point, error) {
	data := key.img.Pixdata()
	d, err := noxrender.DecodeSprite(key.img.Type()&0x3F, data)
	if err != nil {
		return nil, image.Point{}, err
	}
	// the game draws cursors at the mouse position minus cursorSize, with the sprite offset added
	off := image.Point{
		X: int(int32(binary.LittleEndian.Uint32(data[8:]))),
		Y: int(int32(binary.LittleEndian.Uint32(data[12:]))),
	}
	hot := image.Pt(cursorSize, cursorSize).Sub(off)
	src := &image.NRGBA{Pix: d.RGBA, Stride: 4 * d.W, Rect: image.Rect(0, 0, d.W, d.H)}
	if key.mul != [3]uint8{255, 255, 255} {
		for i := 0; i < len(src.Pix); i += 4 {
			for j := 0; j < 3; j++ {
				src.Pix[i+j] = uint8(uint(src.Pix[i+j]) * uint(key.mul[j]) / 255)
			}
		}
	}
	var img image.Image = src
	if key.scale != 100 {
		s := float64(key.scale) / 100
		img = scaleSharp(src, s)
		hot = image.Pt(int(float64(hot.X)*s), int(float64(hot.Y)*s))
	}
	b := img.Bounds()
	hot.X = max(0, min(hot.X, b.Dx()-1))
	hot.Y = max(0, min(hot.Y, b.Dy()-1))
	return img, hot, nil
}

// scaleSharp scales an image with sharp edges: it's scaled up by the integer part of the scale first,
// and the rest is done with a bilinear filter.
func scaleSharp(src *image.NRGBA, s float64) image.Image {
	sz := src.Rect.Size()
	var img image.Image = src
	if k := int(s); k > 1 {
		up := image.NewNRGBA(image.Rect(0, 0, sz.X*k, sz.Y*k))
		draw.NearestNeighbor.Scale(up, up.Rect, src, src.Rect, draw.Src, nil)
		img = up
	}
	dsz := image.Pt(int(math.Round(float64(sz.X)*s)), int(math.Round(float64(sz.Y)*s)))
	if dsz == img.Bounds().Size() {
		return img
	}
	dst := image.NewNRGBA(image.Rectangle{Max: dsz})
	draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Src, nil)
	return dst
}
