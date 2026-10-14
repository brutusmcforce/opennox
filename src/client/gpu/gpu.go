// Package gpu implements an OpenGL drawing backend for noxrender.
//
// Nox originally renders everything into a 16 bit CPU pixel buffer. This package implements the same
// drawing operations on the GPU: sprites are decoded once into texture atlases, and all draw calls are
// batched into textured quads. Each CPU pixel buffer is mapped to a GPU render target (framebuffer object),
// which is presented to the screen instead of uploading the pixel buffer.
//
// The OpenGL context is owned by the window (see libs/client/seat/sdl). All calls must be made from the
// thread where this context is current.
package gpu

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"strings"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/noximage"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

var _ noxrender.GPU = (*Renderer)(nil)

const (
	// A quad is drawn as an instance: the vertex shader makes its 6 vertices. That uploads 17 floats per quad
	// instead of 6 vertices of 9 floats (the shadow edges alone are thousands of quads per frame).
	instSize     = 17 // x0, y0, x1, y1, u0, v0, u1, v1, left r, g, b, a, right r, g, b, a, mode
	maxBatchQuad = 4096
)

const (
	modeNormal   = 0
	modeAdditive = 1
	modeColorize = 2
	modeRaw      = 3 // output color as-is, used for subtraction
)

type blendMode int

const (
	blendOver blendMode = iota
	blendSubtract
)

// target is a render target backing a CPU pixel buffer.
type target struct {
	size image.Point // physical size
	log  image.Point // logical size; drawing coordinates are scaled from it to the physical size
	fbo  uint32
	tex  uint32
	used uint64
}

type batchState struct {
	solid   bool // a batch of solid fills (see solidRenderer), not sprites
	page    *atlasPage
	clip    image.Rectangle
	blend   blendMode
	palette [16][3]uint8
}

// Renderer is an OpenGL drawing backend.
type Renderer struct {
	log *slog.Logger

	prog     uint32
	locSize  int32
	locTex   int32
	locMat   int32
	locPal   int32
	locSharp int32

	present        uint32
	locPresTex     int32
	locPresGamma   int32
	locPresSrc     int32
	locPresOpacity int32
	locPresTexSize int32
	locPresScale   int32

	vao, vbo uint32
	emptyVAO uint32 // used for drawing fullscreen quads
	libsVAO  int32

	verts []float32
	state batchState
	glc   glCache
	pal   [16][3]uint8
	palF  [16 * 3]float32
	palOK bool

	atlas   atlas
	floor   floorRenderer
	solid   solidRenderer
	white   texRegion
	sprites map[noxrender.Image16]texRegion
	broken  map[noxrender.Image16]struct{}
	glyphs  map[noxrender.GlyphKey]texRegion

	targets map[*noximage.Image16]*target
	phys    map[*noximage.Image16]image.Point // physical sizes of render targets, if they differ from the image size
	scratch *target                           // temporary target for Zoom
	cur     *target
	frame   uint64
	upload  []uint32 // reused buffer for Upload

	Gamma     float32
	Filtering bool
}

// New creates a new OpenGL renderer. The OpenGL context must be current.
func New(log *slog.Logger) (*Renderer, error) {
	r := &Renderer{
		log:     log,
		sprites: make(map[noxrender.Image16]texRegion),
		broken:  make(map[noxrender.Image16]struct{}),
		glyphs:  make(map[noxrender.GlyphKey]texRegion),
		targets: make(map[*noximage.Image16]*target),
		Gamma:   1,
	}
	// The window already initialized OpenGL functions, but it's cheap to make sure.
	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("cannot init OpenGL: %w", err)
	}
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &r.libsVAO)
	var err error
	r.prog, err = newProgram(spriteVert, spriteFrag)
	if err != nil {
		return nil, fmt.Errorf("sprite shader: %w", err)
	}
	r.locSize = gl.GetUniformLocation(r.prog, gl.Str("uSize\x00"))
	r.locTex = gl.GetUniformLocation(r.prog, gl.Str("uTex\x00"))
	r.locMat = gl.GetUniformLocation(r.prog, gl.Str("uMat\x00"))
	r.locPal = gl.GetUniformLocation(r.prog, gl.Str("uPal\x00"))
	r.locSharp = gl.GetUniformLocation(r.prog, gl.Str("uSharp\x00"))

	r.present, err = newProgram(presentVert, presentFrag)
	if err != nil {
		return nil, fmt.Errorf("present shader: %w", err)
	}
	r.locPresTex = gl.GetUniformLocation(r.present, gl.Str("uTex\x00"))
	r.locPresGamma = gl.GetUniformLocation(r.present, gl.Str("uGamma\x00"))
	r.locPresSrc = gl.GetUniformLocation(r.present, gl.Str("uSrc\x00"))
	r.locPresOpacity = gl.GetUniformLocation(r.present, gl.Str("uOpacity\x00"))
	r.locPresTexSize = gl.GetUniformLocation(r.present, gl.Str("uTexSize\x00"))
	r.locPresScale = gl.GetUniformLocation(r.present, gl.Str("uScale\x00"))
	if err = r.floor.init(); err != nil {
		return nil, fmt.Errorf("floor shader: %w", err)
	}
	if err = r.solid.init(); err != nil {
		return nil, fmt.Errorf("solid shader: %w", err)
	}

	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	const stride = instSize * 4
	for _, a := range []struct {
		loc, size int32
		off       uintptr
	}{
		{0, 4, 0},      // aDst
		{1, 4, 4 * 4},  // aUV
		{2, 4, 8 * 4},  // aCol (left)
		{4, 4, 12 * 4}, // aColR
		{3, 1, 16 * 4}, // aMode
	} {
		gl.EnableVertexAttribArray(uint32(a.loc))
		gl.VertexAttribPointerWithOffset(uint32(a.loc), a.size, gl.FLOAT, false, stride, a.off)
		gl.VertexAttribDivisor(uint32(a.loc), 1)
	}
	gl.GenVertexArrays(1, &r.emptyVAO)
	r.restore()

	// a white pixel for solid fills
	r.white = r.atlas.add(1, 1, []byte{0xff, 0xff, 0xff, 0xff}, []byte{0})
	r.verts = make([]float32, 0, maxBatchQuad*instSize)
	if errc := gl.GetError(); errc != gl.NO_ERROR {
		return nil, fmt.Errorf("OpenGL error: 0x%x", errc)
	}
	log.Info("GPU renderer initialized", "renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	return r, nil
}

// Close releases all GPU resources.
func (r *Renderer) Close() {
	for _, t := range r.targets {
		t.free()
	}
	r.targets = nil
	if r.scratch != nil {
		r.scratch.free()
		r.scratch = nil
	}
	r.atlas.free()
	r.floor.free()
	r.solid.free()
	gl.DeleteBuffers(1, &r.vbo)
	gl.DeleteVertexArrays(1, &r.vao)
	gl.DeleteVertexArrays(1, &r.emptyVAO)
	gl.DeleteProgram(r.prog)
	gl.DeleteProgram(r.present)
}

// restore OpenGL state expected by the window (libs/client/seat/opengl).
func (r *Renderer) restore() {
	r.glc.valid = false
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindVertexArray(uint32(r.libsVAO))
	gl.Disable(gl.BLEND)
	gl.Disable(gl.SCISSOR_TEST)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.ActiveTexture(gl.TEXTURE0)
}

func (t *target) free() {
	gl.DeleteFramebuffers(1, &t.fbo)
	gl.DeleteTextures(1, &t.tex)
}

func (r *Renderer) targetFor(img *noximage.Image16) *target {
	lsz := img.Rect.Size()
	sz := lsz
	if p, ok := r.phys[img]; ok {
		sz = p
	}
	t := r.targets[img]
	if t != nil && t.size == sz && t.log == lsz {
		return t
	}
	if t != nil {
		if r.cur == t {
			r.flush()
			r.cur = nil
		}
		t.free()
	}
	t = r.newTarget(sz)
	t.log = lsz
	r.targets[img] = t
	return t
}

// SetTargetSize sets the physical size of the render target associated with an image.
// Drawing coordinates are still given in the image coordinates, and are scaled to the physical size.
// Zero size resets the physical size to the image size.
func (r *Renderer) SetTargetSize(img *noximage.Image16, sz image.Point) {
	if sz == (image.Point{}) || sz == img.Rect.Size() {
		delete(r.phys, img)
	} else {
		if r.phys == nil {
			r.phys = make(map[*noximage.Image16]image.Point)
		}
		r.phys[img] = sz
	}
	if r.cur != nil && r.targets[img] == r.cur {
		r.cur = r.targetFor(img)
	}
}

// TargetScale implements noxrender.GPU.
func (r *Renderer) TargetScale() (sx, sy float64) {
	t := r.cur
	if t == nil || t.log.X <= 0 || t.log.Y <= 0 {
		return 1, 1
	}
	return float64(t.size.X) / float64(t.log.X), float64(t.size.Y) / float64(t.log.Y)
}

func (r *Renderer) newTarget(sz image.Point) *target {
	t := &target{size: sz, log: sz}
	gl.GenTextures(1, &t.tex)
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
	setTexParams(gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(sz.X), int32(sz.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.GenFramebuffers(1, &t.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, t.tex, 0)
	if st := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); st != gl.FRAMEBUFFER_COMPLETE {
		r.log.Error("cannot create render target", "size", sz, "status", st)
	}
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	r.glc.valid = false // see glCache
	r.log.Info("created render target", "size", sz)
	return t
}

// SetTarget implements noxrender.GPU.
func (r *Renderer) SetTarget(img *noximage.Image16) {
	t := r.targetFor(img)
	if r.cur == t {
		return
	}
	r.flush()
	r.cur = t
}

// pruneTargets removes render targets that were not used for a while.
func (r *Renderer) pruneTargets() {
	for img, t := range r.targets {
		if t != r.cur && r.frame-t.used > 600 {
			// the physical size is kept: the image may be drawn again, and the target is then recreated with it
			t.free()
			delete(r.targets, img)
		}
	}
}

// Forget releases the render target of an image that won't be used anymore.
func (r *Renderer) Forget(img *noximage.Image16) {
	if t := r.targets[img]; t != nil {
		if r.cur == t {
			r.flush()
			r.cur = nil
		}
		t.free()
		delete(r.targets, img)
	}
	delete(r.phys, img)
}

// Clear implements noxrender.GPU.
func (r *Renderer) Clear(c color.NRGBA) {
	r.clearTarget(float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255)
}

// ClearTransparent clears current target to a transparent color.
func (r *Renderer) ClearTransparent() {
	r.clearTarget(0, 0, 0, 0)
}

func (r *Renderer) clearTarget(cr, cg, cb, ca float32) {
	if r.cur == nil {
		return
	}
	r.flush()
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.cur.fbo)
	gl.Disable(gl.SCISSOR_TEST)
	gl.ClearColor(cr, cg, cb, ca)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	r.restore()
}

func (r *Renderer) setState(page *atlasPage, clip image.Rectangle, blend blendMode, pal *[16][3]uint8) {
	if pal == nil && !r.state.solid && page == r.state.page && clip == r.state.clip && blend == r.state.blend && len(r.verts)+instSize <= cap(r.verts) {
		return // the same batch; fills and glyphs keep the palette, so skip copying and comparing it
	}
	st := batchState{page: page, clip: clip, blend: blend}
	if pal != nil {
		st.palette = *pal
	} else {
		st.palette = r.state.palette
	}
	if st != r.state || len(r.verts)+instSize > cap(r.verts) {
		r.flush()
		r.state = st
	}
}

// setSolidState starts or continues a batch of solid fills.
func (r *Renderer) setSolidState(clip image.Rectangle, blend blendMode) {
	if r.state.solid && clip == r.state.clip && blend == r.state.blend && len(r.solid.data)+solidSize <= cap(r.solid.data) {
		return
	}
	r.flush()
	// the page and the palette stay, so a sprite batch after this can continue with them
	st := r.state
	st.solid, st.clip, st.blend = true, clip, blend
	r.state = st
}

func (r *Renderer) quad(dst image.Rectangle, reg texRegion, src image.Rectangle, c [4]float32, mode float32) {
	r.quadGrad(dst, reg, src, c, c, mode)
}

// quadGrad draws a quad with a horizontal color gradient from c to c2.
func (r *Renderer) quadGrad(dst image.Rectangle, reg texRegion, src image.Rectangle, c, c2 [4]float32, mode float32) {
	r.quadGradF(float32(dst.Min.X), float32(dst.Min.Y), float32(dst.Max.X), float32(dst.Max.Y), reg, src, c, c2, mode)
}

// quadGradF is similar to quadGrad, but accepts fractional coordinates of the destination.
func (r *Renderer) quadGradF(x0, y0, x1, y1 float32, reg texRegion, src image.Rectangle, c, c2 [4]float32, mode float32) {
	ps := float32(reg.page.size)
	u0 := float32(reg.rect.Min.X+src.Min.X) / ps
	v0 := float32(reg.rect.Min.Y+src.Min.Y) / ps
	u1 := float32(reg.rect.Min.X+src.Max.X) / ps
	v1 := float32(reg.rect.Min.Y+src.Max.Y) / ps
	r.verts = append(r.verts,
		x0, y0, x1, y1,
		u0, v0, u1, v1,
		c[0], c[1], c[2], c[3],
		c2[0], c2[1], c2[2], c2[3],
		mode,
	)
}

func (r *Renderer) flush() {
	n := len(r.verts)
	if r.state.solid {
		n = len(r.solid.data)
	}
	if n == 0 || r.cur == nil {
		r.verts = r.verts[:0]
		r.solid.data = r.solid.data[:0]
		return
	}
	t := r.cur
	t.used = r.frame
	st := &r.state
	// The OpenGL state stays set between batches (see glCache): a frame has many batches, and setting
	// all of the state for each of them took more time than drawing.
	c := &r.glc
	if !c.valid || c.target != t || c.size != t.size || c.log != t.log {
		gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
		gl.Viewport(0, 0, int32(t.size.X), int32(t.size.Y))
		gl.Enable(gl.BLEND)
		*c = glCache{valid: true, target: t, size: t.size, log: t.log}
	}
	if st.solid {
		if c.prog != r.solid.prog {
			gl.UseProgram(r.solid.prog)
			gl.BindVertexArray(r.solid.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, r.solid.vbo)
			gl.Uniform2f(r.solid.locSize, float32(t.log.X), float32(t.log.Y))
			c.prog = r.solid.prog
		}
		r.setBlendClip(t, st)
		r.solid.draw()
		return
	}
	if c.prog != r.prog {
		gl.UseProgram(r.prog)
		gl.BindVertexArray(r.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
		// coordinates are logical, the viewport is physical
		gl.Uniform2f(r.locSize, float32(t.log.X), float32(t.log.Y))
		gl.Uniform1i(r.locTex, 0)
		gl.Uniform1i(r.locMat, 1)
		// sprites are scaled when the target is larger than the image (e.g. scaled GUI)
		sharp := int32(0)
		if t.size != t.log {
			sharp = 1
		}
		gl.Uniform1i(r.locSharp, sharp)
		c.prog = r.prog
	}
	if !r.palOK || r.pal != st.palette {
		pal := &r.palF // a local would escape to the heap
		for i, c := range st.palette {
			pal[3*i+0] = float32(c[0]) / 255
			pal[3*i+1] = float32(c[1]) / 255
			pal[3*i+2] = float32(c[2]) / 255
		}
		gl.Uniform3fv(r.locPal, 16, &pal[0])
		r.pal, r.palOK = st.palette, true
	}
	if c.page != st.page {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, st.page.tex)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, st.page.mat)
		c.page = st.page
	}
	r.setBlendClip(t, st)
	gl.BufferData(gl.ARRAY_BUFFER, len(r.verts)*4, unsafe.Pointer(&r.verts[0]), gl.STREAM_DRAW) // gl.Ptr allocates
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(r.verts)/instSize))
	r.verts = r.verts[:0]
}

func (r *Renderer) setBlendClip(t *target, st *batchState) {
	c := &r.glc
	if !c.blendSet || c.blend != st.blend {
		switch st.blend {
		case blendSubtract:
			gl.BlendEquation(gl.FUNC_REVERSE_SUBTRACT)
			gl.BlendFunc(gl.ONE, gl.ONE)
		default:
			gl.BlendEquation(gl.FUNC_ADD)
			// colors are premultiplied by the shader
			gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
		}
		c.blend, c.blendSet = st.blend, true
	}
	if !c.clipSet || c.clip != st.clip {
		if !st.clip.Empty() {
			clip := st.clip
			if t.size != t.log {
				clip = image.Rect(
					clip.Min.X*t.size.X/t.log.X, clip.Min.Y*t.size.Y/t.log.Y,
					(clip.Max.X*t.size.X+t.log.X-1)/t.log.X, (clip.Max.Y*t.size.Y+t.log.Y-1)/t.log.Y,
				)
			}
			gl.Enable(gl.SCISSOR_TEST)
			gl.Scissor(int32(clip.Min.X), int32(t.size.Y-clip.Max.Y), int32(clip.Dx()), int32(clip.Dy()))
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
		c.clip, c.clipSet = st.clip, true
	}
}

// glCache is the OpenGL state that flush set. Anything else that changes the state (other drawing,
// texture uploads, creating targets) must call restore, which invalidates it.
type glCache struct {
	valid     bool
	prog      uint32 // the bound program: r.prog or r.solid.prog
	target    *target
	size, log image.Point
	page      *atlasPage
	blend     blendMode
	blendSet  bool
	clip      image.Rectangle
	clipSet   bool
}

func rgbaF(c color.NRGBA) [4]float32 {
	return [4]float32{float32(c.R) / 255, float32(c.G) / 255, float32(c.B) / 255, float32(c.A) / 255}
}

// Fill implements noxrender.GPU.
func (r *Renderer) Fill(rc image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	if rc.Empty() || r.cur == nil {
		return
	}
	white := r.whiteRegion()
	r.setState(white.page, clip, blendOver, nil)
	r.quad(rc, white, image.Rect(0, 0, 1, 1), rgbaF(c), modeNormal)
}

// FillRects implements noxrender.RectFiller.
func (r *Renderer) FillRects(rcs []image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	if r.cur == nil {
		return
	}
	col := rgbaF(c)
	// On scaled targets the sprite shader filters the white texel (sharp bilinear), so keep using it there.
	if t := r.cur; t.size == t.log {
		for _, rc := range rcs {
			if rc.Empty() {
				continue
			}
			if !solidFits(rc) {
				r.Fill(rc, c, clip)
				continue
			}
			r.setSolidState(clip, blendOver)
			r.solid.add(rc, &col)
		}
		return
	}
	for _, rc := range rcs {
		if rc.Empty() {
			continue
		}
		white := r.whiteRegion()
		r.setState(white.page, clip, blendOver, nil)
		r.quadGradF(float32(rc.Min.X), float32(rc.Min.Y), float32(rc.Max.X), float32(rc.Max.Y), white, image.Rect(0, 0, 1, 1), col, col, modeNormal)
	}
}

// Subtract implements noxrender.GPU.
func (r *Renderer) Subtract(rc image.Rectangle, v uint8, clip image.Rectangle) {
	if rc.Empty() || r.cur == nil {
		return
	}
	f := float32(v) / 255
	white := r.whiteRegion()
	r.setState(white.page, clip, blendSubtract, nil)
	r.quad(rc, white, image.Rect(0, 0, 1, 1), [4]float32{f, f, f, 0}, modeRaw)
}

// whiteRegion returns a white pixel on the page of the current batch if it has one,
// so fills between sprites keep the batch going.
func (r *Renderer) whiteRegion() texRegion {
	if p := r.state.page; p != nil && !p.white.Empty() {
		return texRegion{page: p, rect: p.white}
	}
	return r.white
}

func (r *Renderer) spriteRegion(img noxrender.Image16) (texRegion, bool) {
	if reg, ok := r.sprites[img]; ok {
		return reg, true
	}
	if _, ok := r.broken[img]; ok {
		return texRegion{}, false
	}
	d, err := noxrender.DecodeSprite(img.Type()&0x3F, img.Pixdata())
	if err != nil {
		if err != noxrender.ErrEmptySprite {
			r.log.Warn("cannot decode sprite", "err", err, "type", img.Type()&0x3F, "size", len(img.Pixdata()))
		}
		r.broken[img] = struct{}{}
		return texRegion{}, false
	}
	reg := r.atlas.add(d.W, d.H, d.RGBA, d.Mat)
	r.glc.valid = false // the upload changed the texture bindings
	r.sprites[img] = reg
	return reg, true
}

// Sprite implements noxrender.GPU.
func (r *Renderer) Sprite(img noxrender.Image16, pos image.Point, rows int, st *noxrender.SpriteState, clip image.Rectangle) {
	if r.cur == nil {
		return
	}
	reg, ok := r.spriteRegion(img)
	if !ok {
		return
	}
	w, h := reg.rect.Dx(), reg.rect.Dy()
	if rows < h {
		h = rows
	}
	if h <= 0 {
		return
	}
	mode := float32(modeNormal)
	switch st.Mode {
	case noxrender.SpriteAdditive:
		mode = modeAdditive
	case noxrender.SpriteColorize:
		mode = modeColorize
	}
	a := float32(st.Alpha) / 255
	c := [4]float32{float32(st.Mul[0]) / 255, float32(st.Mul[1]) / 255, float32(st.Mul[2]) / 255, a}
	r.setState(reg.page, clip, blendOver, &st.Materials)
	dst := image.Rect(pos.X, pos.Y, pos.X+w, pos.Y+h)
	if st.Gradient {
		l := [4]float32{st.GradL[0], st.GradL[1], st.GradL[2], a}
		rc := [4]float32{st.GradR[0], st.GradR[1], st.GradR[2], a}
		r.quadGrad(dst, reg, image.Rect(0, 0, w, h), l, rc, mode)
		return
	}
	r.quad(dst, reg, image.Rect(0, 0, w, h), c, mode)
}

// Glyph implements noxrender.GPU.
func (r *Renderer) Glyph(key noxrender.GlyphKey, dr image.Rectangle, mask image.Image, maskp image.Point, c color.NRGBA, clip image.Rectangle) {
	if r.cur == nil {
		return
	}
	reg := r.glyphRegion(key, dr, mask, maskp)
	r.setState(reg.page, clip, blendOver, nil)
	r.quad(dr, reg, image.Rect(0, 0, dr.Dx(), dr.Dy()), rgbaF(c), modeNormal)
}

// HasGlyph implements noxrender.GlyphCache.
func (r *Renderer) HasGlyph(key noxrender.GlyphKey) bool {
	_, ok := r.glyphs[key]
	return ok
}

// GlyphHiRes implements noxrender.GPU.
func (r *Renderer) GlyphHiRes(key noxrender.GlyphKey, dr image.Rectangle, mask image.Image, maskp image.Point, c color.NRGBA, clip image.Rectangle) {
	t := r.cur
	if t == nil {
		return
	}
	reg := r.glyphRegion(key, dr, mask, maskp)
	r.setState(reg.page, clip, blendOver, nil)
	sx := float32(t.log.X) / float32(t.size.X)
	sy := float32(t.log.Y) / float32(t.size.Y)
	col := rgbaF(c)
	r.quadGradF(float32(dr.Min.X)*sx, float32(dr.Min.Y)*sy, float32(dr.Max.X)*sx, float32(dr.Max.Y)*sy,
		reg, image.Rect(0, 0, dr.Dx(), dr.Dy()), col, col, modeNormal)
}

// glyphRegion returns an atlas region with a glyph mask, uploading it if necessary.
func (r *Renderer) glyphRegion(key noxrender.GlyphKey, dr image.Rectangle, mask image.Image, maskp image.Point) texRegion {
	reg, ok := r.glyphs[key]
	if !ok {
		w, h := dr.Dx(), dr.Dy()
		rgba := make([]byte, 4*w*h)
		mat := make([]byte, w*h)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				_, _, _, a := mask.At(maskp.X+x, maskp.Y+y).RGBA()
				p := rgba[4*(y*w+x):]
				p[0], p[1], p[2], p[3] = 0xff, 0xff, 0xff, byte(a>>8)
			}
		}
		reg = r.atlas.add(w, h, rgba, mat)
		r.glc.valid = false // the upload changed the texture bindings
		r.glyphs[key] = reg
	}
	return reg
}

// ReadPixels implements noxrender.GPU.
func (r *Renderer) ReadPixels(img *noximage.Image16) *image.NRGBA {
	r.flush()
	t := r.targetFor(img)
	out := image.NewNRGBA(image.Rectangle{Max: t.size})
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)
	gl.ReadPixels(0, 0, int32(t.size.X), int32(t.size.Y), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(out.Pix))
	r.restore()
	// OpenGL rows are bottom-up, and the alpha is not meaningful for the screen
	for y := 0; y < t.size.Y/2; y++ {
		a := out.Pix[y*out.Stride : (y+1)*out.Stride]
		b := out.Pix[(t.size.Y-1-y)*out.Stride : (t.size.Y-y)*out.Stride]
		for i := range a {
			a[i], b[i] = b[i], a[i]
		}
	}
	for i := 3; i < len(out.Pix); i += 4 {
		out.Pix[i] = 0xff
	}
	return out
}

// Composite draws the content of one render target over another one, scaling it to fit.
// The source must contain premultiplied colors (see ClearTransparent).
func (r *Renderer) Composite(dst, src *noximage.Image16) {
	r.flush()
	ts := r.targetFor(src)
	td := r.targetFor(dst)
	gl.BindFramebuffer(gl.FRAMEBUFFER, td.fbo)
	gl.Viewport(0, 0, int32(td.size.X), int32(td.size.Y))
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	r.drawTexture(texDraw{tex: ts.tex, size: ts.size, out: td.size, gamma: 1, opacity: 1})
	r.restore()
}

// Upload replaces the content of the render target associated with an image with the image pixels.
// It's used for images decoded on the CPU, e.g. movie frames.
func (r *Renderer) Upload(img *noximage.Image16) {
	r.flush()
	delete(r.phys, img)
	t := r.targetFor(img)
	sz := t.size
	lut := rgbaLUT()
	if n := sz.X * sz.Y; cap(r.upload) < n {
		r.upload = make([]uint32, n)
	}
	buf := r.upload[:sz.X*sz.Y]
	for y := 0; y < sz.Y; y++ {
		// render targets have the bottom-left origin
		row := img.Row(y)[:sz.X]
		dst := buf[sz.X*(sz.Y-1-y):][:sz.X]
		for x, v := range row {
			dst[x] = lut[v]
		}
	}
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 4)
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int32(sz.X), int32(sz.Y), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(buf))
	r.restore()
}

var rgbaTable struct {
	once sync.Once
	lut  *[1 << 16]uint32
}

// rgbaLUT returns a table that converts RGBA5551 pixels to opaque RGBA8 (in memory order R, G, B, A).
func rgbaLUT() *[1 << 16]uint32 {
	rgbaTable.once.Do(func() {
		lut := new([1 << 16]uint32)
		for v := range lut {
			c := noxcolor.RGBA5551(v).ColorNRGBA()
			lut[v] = binary.LittleEndian.Uint32([]byte{c.R, c.G, c.B, 0xff})
		}
		rgbaTable.lut = lut
	})
	return rgbaTable.lut
}

// Zoom replaces the content of the render target associated with an image with a part of it
// (in top-left origin coordinates), scaled to fill the whole target.
func (r *Renderer) Zoom(img *noximage.Image16, src image.Rectangle) {
	r.flush()
	t := r.targetFor(img)
	if r.scratch == nil || r.scratch.size != t.size {
		if r.scratch != nil {
			r.scratch.free()
		}
		r.scratch = r.newTarget(t.size)
	}
	w, h := int32(t.size.X), int32(t.size.Y)
	gl.Disable(gl.SCISSOR_TEST)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, t.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, r.scratch.fbo)
	gl.BlitFramebuffer(0, 0, w, h, 0, 0, w, h, gl.COLOR_BUFFER_BIT, gl.NEAREST)

	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.Viewport(0, 0, w, h)
	gl.Disable(gl.BLEND)
	// render targets have the bottom-left origin
	r.drawTexture(texDraw{
		tex:     r.scratch.tex,
		size:    t.size,
		src:     image.Rect(src.Min.X, t.size.Y-src.Max.Y, src.Max.X, t.size.Y-src.Min.Y),
		out:     t.size,
		gamma:   1,
		opacity: 1,
	})
	r.restore()
}

// Present draws the render target associated with an image to the screen.
// The view is a rectangle on the screen (top-left origin), and win is the window size.
func (r *Renderer) Present(img *noximage.Image16, view image.Rectangle, win image.Point) {
	r.flush()
	r.frame++
	if r.frame%300 == 0 {
		r.pruneTargets()
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Disable(gl.SCISSOR_TEST)
	gl.Viewport(0, 0, int32(win.X), int32(win.Y))
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	if img == nil {
		r.restore()
		return
	}
	t := r.targetFor(img)
	gl.Viewport(int32(view.Min.X), int32(win.Y-view.Max.Y), int32(view.Dx()), int32(view.Dy()))
	gl.Disable(gl.BLEND)
	r.drawTexture(texDraw{tex: t.tex, size: t.size, out: view.Size(), gamma: r.Gamma, opacity: 1})
	r.restore()
}

// texDraw describes a texture drawn over the whole viewport.
type texDraw struct {
	tex     uint32
	size    image.Point     // texture size
	src     image.Rectangle // part of the texture to draw (bottom-left origin); empty means the whole texture
	out     image.Point     // viewport size
	gamma   float32
	opacity float32
}

// drawTexture draws a texture over the whole viewport.
//
// With filtering enabled, the texture is scaled using sharp bilinear filtering: texels stay sharp,
// and only the edges between them are interpolated. Otherwise, the nearest neighbour filtering is used.
func (r *Renderer) drawTexture(d texDraw) {
	src := d.src
	if src.Empty() {
		src = image.Rectangle{Max: d.size}
	}
	gl.UseProgram(r.present)
	gl.BindVertexArray(r.emptyVAO)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, d.tex)
	f := int32(gl.NEAREST)
	if r.Filtering {
		f = gl.LINEAR
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, f)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, f)
	gl.Uniform1i(r.locPresTex, 0)
	gl.Uniform1f(r.locPresGamma, d.gamma)
	gl.Uniform1f(r.locPresOpacity, d.opacity)
	tw, th := float32(d.size.X), float32(d.size.Y)
	gl.Uniform4f(r.locPresSrc, float32(src.Min.X)/tw, float32(src.Min.Y)/th, float32(src.Dx())/tw, float32(src.Dy())/th)
	gl.Uniform2f(r.locPresTexSize, tw, th)
	if r.Filtering {
		gl.Uniform2f(r.locPresScale, float32(d.out.X)/float32(src.Dx()), float32(d.out.Y)/float32(src.Dy()))
	} else {
		gl.Uniform2f(r.locPresScale, 0, 0)
	}
	// the vertex shader generates a fullscreen triangle strip from vertex IDs
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	// restore nearest filtering for render targets
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func newProgram(vert, frag string) (uint32, error) {
	vs, err := compileShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fs, err := compileShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vs)
		return 0, err
	}
	p := gl.CreateProgram()
	gl.AttachShader(p, vs)
	gl.AttachShader(p, fs)
	gl.BindAttribLocation(p, 0, gl.Str("aPos\x00"))
	gl.BindAttribLocation(p, 0, gl.Str("aDst\x00"))
	gl.BindAttribLocation(p, 1, gl.Str("aUV\x00"))
	gl.BindAttribLocation(p, 2, gl.Str("aCol\x00"))
	gl.BindAttribLocation(p, 3, gl.Str("aMode\x00"))
	gl.BindAttribLocation(p, 4, gl.Str("aColR\x00"))
	gl.LinkProgram(p)
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var n int32
		gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &n)
		msg := strings.Repeat("\x00", int(n+1))
		gl.GetProgramInfoLog(p, n, nil, gl.Str(msg))
		gl.DeleteProgram(p)
		return 0, errors.New(strings.TrimRight(msg, "\x00"))
	}
	return p, nil
}

func compileShader(src string, typ uint32) (uint32, error) {
	s := gl.CreateShader(typ)
	csrc, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s, 1, csrc, nil)
	free()
	gl.CompileShader(s)
	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var n int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &n)
		msg := strings.Repeat("\x00", int(n+1))
		gl.GetShaderInfoLog(s, n, nil, gl.Str(msg))
		gl.DeleteShader(s)
		return 0, errors.New(strings.TrimRight(msg, "\x00"))
	}
	return s, nil
}

const spriteVert = `#version 330 core
in vec4 aDst;  // x0, y0, x1, y1
in vec4 aUV;   // u0, v0, u1, v1
in vec4 aCol;  // the color of the left side
in vec4 aColR; // the color of the right side
in float aMode;
uniform vec2 uSize;
out vec2 vUV;
out vec4 vCol;
flat out int vMode;
// the corners of the two triangles of a quad, in the order the quads were drawn before instancing
const vec2 corners[6] = vec2[6](vec2(0, 0), vec2(1, 0), vec2(0, 1), vec2(1, 0), vec2(1, 1), vec2(0, 1));
void main() {
	vec2 c = corners[gl_VertexID];
	// mix with 0 or 1 gives the exact end values
	vec2 pos = mix(aDst.xy, aDst.zw, c);
	vUV = mix(aUV.xy, aUV.zw, c);
	vCol = c.x == 0.0 ? aCol : aColR;
	vMode = int(aMode + 0.5);
	vec2 p = pos / uSize * 2.0 - 1.0;
	gl_Position = vec4(p.x, -p.y, 0.0, 1.0);
}
`

const spriteFrag = `#version 330 core
in vec2 vUV;
in vec4 vCol;
flat in int vMode;
uniform sampler2D uTex;
uniform sampler2D uMat;
uniform vec3 uPal[16];
uniform bool uSharp; // use sharp bilinear filtering for scaled sprites
out vec4 oColor;

// shade returns the premultiplied color of a texel.
vec4 shade(ivec2 tc) {
	vec4 t = texelFetch(uTex, tc, 0);
	int m = int(texelFetch(uMat, tc, 0).r * 255.0 + 0.5);
	vec3 c = t.rgb;
	if (m > 0) {
		c *= uPal[m - 1];
	}
	if (vMode == 2 && m == 0) {
		c = vCol.rgb * dot(t.rgb, vec3(0.297, 0.586, 0.109));
	} else {
		c *= vCol.rgb;
	}
	if (vMode == 1) {
		// additive
		return vec4(c * t.a, 0.0);
	}
	float a = t.a * vCol.a;
	return vec4(c * a, a);
}

void main() {
	if (vMode == 3) {
		oColor = vCol;
		return;
	}
	vec2 p = vUV * vec2(textureSize(uTex, 0)); // in texels
	if (!uSharp) {
		oColor = shade(ivec2(floor(p)));
		return;
	}
	// Sharp bilinear: with a non-integer scale, nearest filtering makes some texels larger than others.
	// Instead, texels keep a solid center, and only a one pixel wide band at texel edges is interpolated.
	vec2 s = 1.0 / max(fwidth(p), vec2(1e-6)); // output pixels per texel
	vec2 tl = floor(p);
	vec2 cd = p - tl - 0.5;
	vec2 rr = max(0.5 - 0.5 / s, vec2(0.0));
	vec2 q = tl + (cd - clamp(cd, -rr, rr)) * s; // bilinear sample position, relative to texel centers
	vec2 qi = floor(q);
	vec2 w = q - qi;
	ivec2 i = ivec2(qi);
	vec4 c00 = shade(i);
	vec4 c10 = shade(i + ivec2(1, 0));
	vec4 c01 = shade(i + ivec2(0, 1));
	vec4 c11 = shade(i + ivec2(1, 1));
	oColor = mix(mix(c00, c10, w.x), mix(c01, c11, w.x), w.y);
}
`

const presentVert = `#version 330 core
uniform vec4 uSrc; // source rectangle in texture coordinates: xy is the origin, zw is the size
out vec2 vUV;
void main() {
	vec2 p = vec2(float(gl_VertexID & 1), float((gl_VertexID >> 1) & 1));
	vUV = uSrc.xy + p * uSrc.zw;
	gl_Position = vec4(p * 2.0 - 1.0, 0.0, 1.0);
}
`

const presentFrag = `#version 330 core
in vec2 vUV;
uniform sampler2D uTex;
uniform float uGamma;
uniform float uOpacity;
uniform vec2 uTexSize; // texture size in texels
uniform vec2 uScale;   // output pixels per texel; zero disables sharp filtering
out vec4 oColor;
void main() {
	vec2 uv = vUV;
	if (uScale.x > 0.0) {
		// sharp bilinear: sample the texel center, except for a one pixel wide band at texel edges
		vec2 scale = max(uScale, vec2(1.0));
		vec2 t = uv * uTexSize;
		vec2 d = fract(t) - 0.5;
		vec2 r = 0.5 - 0.5 / scale;
		uv = (floor(t) + 0.5 + (d - clamp(d, -r, r)) * scale) / uTexSize;
	}
	vec4 c = texture(uTex, uv);
	c.rgb = pow(c.rgb, vec3(1.0 / uGamma));
	oColor = c * uOpacity;
}
`
