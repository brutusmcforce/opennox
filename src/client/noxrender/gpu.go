package noxrender

import (
	"encoding/binary"
	"image"
	"image/color"
	"math"

	"github.com/brutusmcforce/opennox/libs/noximage"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// SpriteMode selects how a sprite is combined with the render target.
type SpriteMode byte

const (
	// SpriteNormal draws sprite using its own alpha and a global alpha.
	SpriteNormal SpriteMode = iota
	// SpriteAdditive adds sprite color to the destination (used by particles).
	SpriteAdditive
	// SpriteColorize replaces sprite color with the color multiplier, scaled by the sprite luminance.
	SpriteColorize
)

// SpriteState is the render state for drawing a sprite on the GPU.
type SpriteState struct {
	Mode SpriteMode
	// Mul is a color multiplier (or a target color for colorize mode). White means no change.
	Mul [3]uint8
	// Alpha is a global alpha for the sprite.
	Alpha uint8
	// Materials are colors for indexed (material) pixels of the sprite.
	Materials [16][3]uint8
	// Gradient replaces Mul with a horizontal gradient from GradL to GradR (1 is no change).
	Gradient     bool
	GradL, GradR [3]float32
}

// GlyphKey identifies a rasterized glyph.
type GlyphKey struct {
	Face font.Face
	Rune rune
	Frac fixed.Point26_6 // fractional part of the dot
}

// TileCell is the size of the cells in FloorParams.Dirty: one floor tile.
const TileCell = 46

// FloorParams describes floor tiles for drawing them on the GPU.
//
// Floor tiles are pre-rendered by the game into a circular tile buffer. For a screen position (x, y)
// the index in the tile buffer is Pitch*(y+Off.Y) + (x+Off.X), modulo the buffer length.
type FloorParams struct {
	Tiles []uint16 // RGB555
	Pitch int
	Off   image.Point
	// Dirty has a flag for each TileCell x TileCell cell of Tiles (DirtyCols per row) that changed since the last call.
	// Only those cells are uploaded, and the flags are cleared. Nil means the whole buffer changed.
	Dirty     []byte
	DirtyCols int
	// Lit enables lighting. Light is a grid of RGB light multipliers, stored row by row.
	Lit      bool
	Light    []float32
	LightW   int
	LightH   int
	LightOff image.Point // light grid position (in pixels) is a screen position plus this offset
	Step     int         // light grid step, in pixels
}

// GPU is the drawing backend. All drawing goes through it.
// Pixel buffers are only used to identify render targets and their (logical) size; their pixels are not used.
//
// Colors are passed in non-premultiplied form. An empty clip rectangle means no clipping.
type GPU interface {
	// SetTarget selects a render target associated with a given pixel buffer.
	SetTarget(img *noximage.Image16)
	// Clear the current render target.
	Clear(c color.NRGBA)
	// Fill a rectangle with a color, blending it according to the alpha.
	Fill(rc image.Rectangle, c color.NRGBA, clip image.Rectangle)
	// Sprite draws an image at a given position. Only the first rows of the image are drawn.
	Sprite(img Image16, pos image.Point, rows int, st *SpriteState, clip image.Rectangle)
	// Glyph draws a glyph mask with a given color.
	Glyph(key GlyphKey, dr image.Rectangle, mask image.Image, maskp image.Point, c color.NRGBA, clip image.Rectangle)
	// TargetScale returns the ratio of the physical size of the current render target to its logical size.
	// All coordinates passed to the GPU are logical, but the target may be drawn at a higher resolution.
	TargetScale() (sx, sy float64)
	// GlyphHiRes is similar to Glyph, but dr is given in physical pixels of the render target (see TargetScale).
	GlyphHiRes(key GlyphKey, dr image.Rectangle, mask image.Image, maskp image.Point, c color.NRGBA, clip image.Rectangle)
	// Subtract a value from each color channel of the destination.
	Subtract(rc image.Rectangle, v uint8, clip image.Rectangle)
	// ReadPixels reads back the content of the render target associated with a given pixel buffer.
	ReadPixels(img *noximage.Image16) *image.NRGBA
}

// SetGPU sets the GPU drawing backend. Nothing is drawn without it.
func (r *NoxRender) SetGPU(g GPU) {
	r.gpu = g
	if g != nil && r.pix != nil {
		g.SetTarget(r.pix)
	}
}

// GPU returns current GPU backend, if any.
func (r *NoxRender) GPU() GPU {
	return r.gpu
}

func (r *NoxRender) gpuClip() image.Rectangle {
	if r.p.Clip() {
		return r.p.ClipRect()
	}
	return image.Rectangle{}
}

func rgb8(c Color16) [3]uint8 {
	return [3]uint8{uint8(c.R), uint8(c.G), uint8(c.B)}
}

func toNRGBA(cl color.Color, a uint8) color.NRGBA {
	c := color.NRGBAModel.Convert(cl).(color.NRGBA)
	c.A = a
	return c
}

func color16NRGBA(c Color16, a uint8) color.NRGBA {
	return color.NRGBA{R: uint8(c.R), G: uint8(c.G), B: uint8(c.B), A: a}
}

// gpuFill fills a rectangle. Alpha is an opacity of the color (0-256).
func (r *NoxRender) gpuFill(rc image.Rectangle, c Color16, alpha int) {
	if alpha <= 0 || rc.Empty() {
		return
	}
	if alpha > 255 {
		alpha = 255
	}
	r.gpu.Fill(rc, color16NRGBA(c, uint8(alpha)), image.Rectangle{})
}

func (r *NoxRender) gpuPixel(x, y int, c Color16, alpha int) {
	r.gpuFill(image.Rect(x, y, x+1, y+1), c, alpha)
}

// gpuDrawImage draws an image on the GPU. See DrawImage16 for the software implementation.
func (r *NoxRender) gpuDrawImage(img Image16, pos image.Point) {
	typ := img.Type() & 0x3F
	switch typ {
	case 2, 3, 4, 5, 6, 7, 8:
	default:
		return
	}
	data := img.Pixdata()
	if len(data) < 17 {
		return
	}
	width := int(binary.LittleEndian.Uint32(data[0:]))
	height := int(binary.LittleEndian.Uint32(data[4:]))
	pos.X += int(int32(binary.LittleEndian.Uint32(data[8:])))
	pos.Y += int(int32(binary.LittleEndian.Uint32(data[12:])))
	if r.dword_5d4594_3799484 != 0 {
		height -= int(r.dword_5d4594_3799484)
		if height <= 0 {
			return
		}
		r.dword_5d4594_3799476 = pos.Y + height
	}
	if typ != 2 && typ != 7 && r.HookImageDrawXxx != nil {
		r.HookImageDrawXxx(pos, image.Point{X: width, Y: height})
	}
	var st SpriteState
	st.Mul = [3]uint8{255, 255, 255}
	st.Alpha = 0xff
	switch typ {
	case 8:
		st.Mode = SpriteAdditive
	case 3, 4, 5, 6:
		d := r.p
		if d.IsAlphaEnabled() {
			st.Alpha = d.Alpha()
		}
		if d.Multiply14() {
			if !(d.IsAlphaEnabled() && st.Alpha == 0xff && d.Flag16()) {
				st.Mul = rgb8(d.ColorMultA())
			}
		} else if !d.IsAlphaEnabled() && d.Colorize17() {
			st.Mode = SpriteColorize
			st.Mul = rgb8(d.ColorMultA())
		}
		for i := range st.Materials {
			st.Materials[i] = rgb8(d.ColorMultOp(i))
		}
	}
	r.gpu.Sprite(img, pos, height, &st, r.gpuClip())
}

// DrawImageGradientGPU draws an image on the GPU, multiplying it by a horizontal color gradient.
// Cut is the number of rows to skip at the bottom of the image. It returns false if GPU rendering is disabled.
func (r *NoxRender) DrawImageGradientGPU(img Image16, pos image.Point, cut int, left, right [3]float32, clip image.Rectangle) bool {
	if r.gpu == nil {
		return false
	}
	if img == nil {
		return true
	}
	if p, ok := img.(*Image); ok && p == nil {
		return true
	}
	data := img.Pixdata()
	if len(data) < 17 {
		return true
	}
	height := int(binary.LittleEndian.Uint32(data[4:])) - cut
	if height <= 0 || clip.Empty() {
		return true
	}
	pos.X += int(int32(binary.LittleEndian.Uint32(data[8:])))
	pos.Y += int(int32(binary.LittleEndian.Uint32(data[12:])))
	st := SpriteState{
		Mul:      [3]uint8{255, 255, 255},
		Alpha:    0xff,
		Gradient: true,
		GradL:    left,
		GradR:    right,
	}
	r.gpu.Sprite(img, pos, height, &st, clip)
	return true
}

// gpuDrawString draws a string on the GPU. It mimics font.Drawer.DrawString.
//
// If the render target has a higher resolution than its logical size, outline fonts are drawn at the full resolution.
// The text layout is not affected: glyphs are positioned using the font at its original size.
func (r *NoxRender) gpuDrawString(s string, dot fixed.Point26_6, clip image.Rectangle) fixed.Int26_6 {
	face := r.text.face
	c := toNRGBA(r.p.TextColor(), 0xff)
	var (
		hi     font.Face
		sx, sy float64
	)
	of, _ := face.f.(*outlineFace)
	if of != nil {
		sx, sy = r.gpu.TargetScale()
		if sx > 1.01 || sy > 1.01 {
			// The scale changes with every window size, and each face fills the glyph atlas, which is never freed.
			// Glyphs are positioned with the exact scale, so a slightly different glyph size isn't noticeable.
			hi = of.scaledFace(math.Round(sx*8) / 8)
		}
	}
	prev := rune(-1)
	for _, ch := range s {
		if prev >= 0 {
			dot.X += face.Kern(prev, ch)
		}
		if hi != nil && of.hasGlyph(ch) {
			pdot := fixed.Point26_6{X: fixed.Int26_6(float64(dot.X) * sx), Y: fixed.Int26_6(float64(dot.Y+of.yOff) * sy)}
			dr, mask, maskp, _, ok := hi.Glyph(pdot, ch)
			if ok && !dr.Empty() {
				key := GlyphKey{Face: hi, Rune: ch, Frac: fixed.Point26_6{X: pdot.X & 0x3f, Y: pdot.Y & 0x3f}}
				r.gpu.GlyphHiRes(key, dr, mask, maskp, c, clip)
			}
			if adv, ok := face.GlyphAdvance(ch); ok {
				dot.X += adv
			}
			prev = ch
			continue
		}
		dr, mask, maskp, advance, ok := face.Glyph(dot, ch)
		if !ok {
			continue
		}
		key := GlyphKey{
			Face: face.f,
			Rune: ch,
			Frac: fixed.Point26_6{X: dot.X & 0x3f, Y: dot.Y & 0x3f},
		}
		if !dr.Empty() {
			r.gpu.Glyph(key, dr, mask, maskp, c, clip)
		}
		dot.X += advance
		prev = ch
	}
	return dot.X
}
