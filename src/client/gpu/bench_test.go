package gpu

import (
	"image"
	"image/color"
	"log/slog"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/brutusmcforce/opennox/libs/gl"

	"github.com/brutusmcforce/opennox/client/noxrender"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/noximage"
	"github.com/brutusmcforce/opennox/libs/sdl"
)

// The benchmarks need an OpenGL 3.3 context, so they create a hidden SDL window.
// They are skipped when there is no display (e.g. in CI).

var benchGL struct {
	once sync.Once
	win  *sdl.Window
	ctx  sdl.GLContext
	r    *Renderer
	err  error
}

// withGL makes the shared benchmark context current on the calling goroutine's thread.
func withGL(b testing.TB) *Renderer {
	b.Helper()
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		b.Skip("no display")
	}
	runtime.LockOSThread()
	b.Cleanup(func() {
		_ = benchGL.win.GLMakeCurrent(nil)
		runtime.UnlockOSThread()
	})
	g := &benchGL
	g.once.Do(func() {
		if g.err = sdl.Init(sdl.INIT_VIDEO); g.err != nil {
			return
		}
		_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
		_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
		_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
		g.win, g.err = sdl.CreateWindow("bench", 0, 0, 64, 64, sdl.WINDOW_OPENGL|sdl.WINDOW_HIDDEN)
		if g.err != nil {
			return
		}
		g.ctx, g.err = g.win.GLCreateContext()
		if g.err != nil {
			return
		}
		g.r, g.err = New(slog.New(slog.DiscardHandler))
	})
	if g.err != nil {
		b.Skip("no OpenGL:", g.err)
	}
	if err := g.win.GLMakeCurrent(g.ctx); err != nil {
		b.Fatal(err)
	}
	return g.r
}

// floorBench describes the tile buffer the game uses for a given world view size (see nox_xxx_tileInitBuf_430DB0).
type floorBench struct {
	name string
	view image.Point
}

var floorSizes = []floorBench{
	{"1080p", image.Pt(1920, 1080)},
	{"4K", image.Pt(3840, 2160)},
	{"max-zoom", image.Pt(5198, 4094)},
}

func newFloorParams(view image.Point) *noxrender.FloorParams {
	const cell = noxrender.TileCell
	cols, rows := view.X/cell+4, view.Y/cell+3
	w, h := cols*cell, rows*cell
	tiles := make([]uint16, w*h)
	for i := range tiles {
		tiles[i] = uint16(i)
	}
	return &noxrender.FloorParams{Tiles: tiles, Pitch: w, Dirty: make([]byte, cols*rows), DirtyCols: cols}
}

// floorChange marks the cells that the game redraws in a frame.
type floorChange struct {
	name string
	mark func(p *noxrender.FloorParams)
}

var floorChanges = []floorChange{
	{"full", func(p *noxrender.FloorParams) {
		for i := range p.Dirty {
			p.Dirty[i] = 1
		}
	}},
	{"static", func(p *noxrender.FloorParams) {}},
	{"scroll-x", func(p *noxrender.FloorParams) {
		for i := 0; i < len(p.Dirty); i += p.DirtyCols {
			p.Dirty[i+3] = 1 // a new column of tiles
		}
	}},
	{"scroll-y", func(p *noxrender.FloorParams) {
		clear(p.Dirty[:p.DirtyCols])
		for i := range p.DirtyCols {
			p.Dirty[2*p.DirtyCols+i] = 1 // a new row of tiles
		}
	}},
	{"scroll-xy", func(p *noxrender.FloorParams) {
		for i := 0; i < len(p.Dirty); i += p.DirtyCols {
			p.Dirty[i+3] = 1
		}
		for i := range p.DirtyCols {
			p.Dirty[2*p.DirtyCols+i] = 1
		}
	}},
}

func floorSpans(view image.Point) []image.Rectangle {
	spans := make([]image.Rectangle, view.Y)
	for y := range spans {
		spans[y] = image.Rect(0, y, view.X, y+1)
	}
	return spans
}

// BenchmarkFloor draws a full screen of floor tiles, like the game does every frame,
// with different parts of the tile buffer changed since the previous frame.
func BenchmarkFloor(b *testing.B) {
	for _, sz := range floorSizes {
		for _, ch := range floorChanges {
			b.Run(sz.name+"/"+ch.name, func(b *testing.B) {
				r := withGL(b)
				img := noximage.NewImage16(image.Rectangle{Max: sz.view})
				p := newFloorParams(sz.view)
				spans := floorSpans(sz.view)
				r.SetTarget(img)
				r.Floor(p, spans) // allocate the textures
				gl.Finish()
				b.ResetTimer()
				for range b.N {
					ch.mark(p)
					r.Floor(p, spans)
					gl.Finish()
				}
				b.StopTimer()
				r.Forget(img)
			})
		}
	}
}

// TestFloorUploadDirty checks that uploading only the dirty cells keeps the texture equal to the tile buffer.
func TestFloorUploadDirty(t *testing.T) {
	r := withGL(t)
	view := image.Pt(700, 500)
	p := newFloorParams(view)
	img := noximage.NewImage16(image.Rectangle{Max: view})
	spans := floorSpans(view)
	r.SetTarget(img)
	defer r.Forget(img)
	r.Floor(p, spans)
	rnd := rand.New(rand.NewPCG(1, 2))
	got := make([]uint16, len(p.Tiles))
	for iter := range 50 {
		// change some cells (and mark them), plus whole rows and columns like scrolling does
		for range rnd.IntN(20) {
			cx, cy := rnd.IntN(p.DirtyCols), rnd.IntN(len(p.Dirty)/p.DirtyCols)
			fillCell(p, cx, cy, uint16(iter*7+cx))
		}
		if iter%3 == 0 {
			cx := rnd.IntN(p.DirtyCols)
			for cy := range len(p.Dirty) / p.DirtyCols {
				fillCell(p, cx, cy, uint16(iter))
			}
		}
		r.Floor(p, spans)
		gl.BindTexture(gl.TEXTURE_2D, r.floor.tiles)
		gl.PixelStorei(gl.PACK_ALIGNMENT, 1)
		gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RED_INTEGER, gl.UNSIGNED_SHORT, gl.Ptr(got))
		for i := range got {
			if got[i] != p.Tiles[i] {
				t.Fatalf("iteration %d: pixel %d,%d: got %d, want %d", iter, i%p.Pitch, i/p.Pitch, got[i], p.Tiles[i])
			}
		}
		for _, d := range p.Dirty {
			if d != 0 {
				t.Fatalf("iteration %d: dirty flags are not cleared", iter)
			}
		}
	}
}

func fillCell(p *noxrender.FloorParams, cx, cy int, v uint16) {
	const cell = noxrender.TileCell
	for y := cy * cell; y < (cy+1)*cell; y++ {
		for x := cx * cell; x < (cx+1)*cell; x++ {
			p.Tiles[y*p.Pitch+x] = v + uint16(x^y)
		}
	}
	p.Dirty[cy*p.DirtyCols+cx] = 1
}

// BenchmarkUpload uploads a CPU-decoded image (a movie frame) to its render target.
func BenchmarkUpload(b *testing.B) {
	for _, sz := range []floorBench{{"640x480", image.Pt(640, 480)}, {"1080p", image.Pt(1920, 1080)}} {
		b.Run(sz.name, func(b *testing.B) {
			r := withGL(b)
			img := noximage.NewImage16(image.Rectangle{Max: sz.view})
			for i := range img.Pix {
				img.Pix[i] = uint16(i)
			}
			r.Upload(img)
			gl.Finish()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				r.Upload(img)
				gl.Finish()
			}
			b.StopTimer()
			r.Forget(img)
		})
	}
}

// TestUploadColors checks Upload against the direct color conversion, including the bottom-left origin.
func TestUploadColors(t *testing.T) {
	r := withGL(t)
	sz := image.Pt(256, 256) // all 65536 colors
	img := noximage.NewImage16(image.Rectangle{Max: sz})
	for i := range img.Pix {
		img.Pix[i] = uint16(i)
	}
	r.Upload(img)
	defer r.Forget(img)
	got := make([]byte, 4*sz.X*sz.Y)
	gl.BindTexture(gl.TEXTURE_2D, r.targets[img].tex)
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(got))
	for y := range sz.Y {
		for x := range sz.X {
			c := noxcolor.RGBA5551(img.Pix[y*sz.X+x]).ColorNRGBA()
			p := got[4*((sz.Y-1-y)*sz.X+x):]
			if p[0] != c.R || p[1] != c.G || p[2] != c.B || p[3] != 0xff {
				t.Fatalf("pixel %d,%d: got %v, want %v", x, y, p[:4], c)
			}
		}
	}
}

// Sprite fills (scaled targets) between sprites on a later atlas page must not switch pages (each switch flushes the batch),
// and still draw the right color.
func TestFillKeepsPage(t *testing.T) {
	r := withGL(t)
	img := noximage.NewImage16(image.Rect(0, 0, 64, 64))
	r.SetTarget(img)
	defer r.Forget(img)
	// a region that needs a page of its own
	big := atlasSize - 8
	reg := r.atlas.add(big, big, make([]byte, big*big*4), make([]byte, big*big))
	if reg.page == r.white.page {
		reg = r.atlas.add(big, big, make([]byte, big*big*4), make([]byte, big*big))
	}
	if reg.page == r.white.page || reg.page.white.Empty() {
		t.Fatal("expected a new atlas page with a white pixel")
	}
	full := image.Rect(0, 0, 64, 64)
	r.setState(reg.page, full, blendOver, nil)
	r.fillSprite(image.Rect(4, 4, 20, 20), color.NRGBA{R: 255, A: 255}, full)
	if r.state.page != reg.page {
		t.Fatal("the fill switched the atlas page")
	}
	out := r.ReadPixels(img)
	if c := out.NRGBAAt(10, 10); c.R < 240 || c.G > 15 || c.B > 15 {
		t.Fatalf("fill color: %v", c)
	}
}

// Solid fills (Fill and FillRects on an unscaled target) must draw exactly the pixels of the sprite shader.
func TestSolidFillsMatchFill(t *testing.T) {
	r := withGL(t)
	full := image.Rect(0, 0, 128, 96)
	rnd := rand.New(rand.NewPCG(7, 8))
	type fill struct {
		rc image.Rectangle
		c  color.NRGBA
	}
	var fills []fill
	for i := 0; i < 300; i++ {
		x, y := rnd.IntN(120), rnd.IntN(90)
		fills = append(fills, fill{
			rc: image.Rect(x, y, x+1+rnd.IntN(20), y+1+rnd.IntN(3)),
			c:  color.NRGBA{R: uint8(rnd.IntN(256)), G: uint8(rnd.IntN(256)), B: uint8(rnd.IntN(256)), A: uint8(rnd.IntN(256))},
		})
	}
	draw := func(solid bool) []byte {
		img := noximage.NewImage16(full)
		r.SetTarget(img)
		defer r.Forget(img)
		r.Clear(color.NRGBA{R: 40, G: 90, B: 160, A: 255})
		for i, f := range fills {
			switch {
			case !solid:
				r.fillSprite(f.rc, f.c, full)
			case i%2 == 0:
				r.FillRects([]image.Rectangle{f.rc}, f.c, full)
			default:
				r.Fill(f.rc, f.c, full)
			}
		}
		if solid != r.state.solid {
			t.Fatal("the fills didn't use the expected batch")
		}
		return r.ReadPixels(img).Pix
	}
	want, got := draw(false), draw(true)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pixel %d channel %d: got %d, want %d", i/4, i%4, got[i], want[i])
		}
	}
}

// FillRects with more rectangles than a batch takes, and with ones that don't fit a solid instance
// in between, must draw the same as filling them one by one.
func TestFillRectsOverflow(t *testing.T) {
	r := withGL(t)
	full := image.Rect(0, 0, 128, 96)
	rnd := rand.New(rand.NewPCG(9, 10))
	var rcs []image.Rectangle
	for i := 0; i < maxSolidQuad*2+100; i++ {
		x, y := rnd.IntN(126), rnd.IntN(94)
		rcs = append(rcs, image.Rect(x, y, x+1+rnd.IntN(2), y+1+rnd.IntN(2)))
		if i%5000 == 0 {
			rcs = append(rcs, image.Rect(-40000, y, 40000, y+1)) // not a solid instance
		}
	}
	c := color.NRGBA{R: 200, G: 120, B: 30, A: 40}
	draw := func(batched bool) []byte {
		img := noximage.NewImage16(full)
		r.SetTarget(img)
		defer r.Forget(img)
		r.Clear(color.NRGBA{R: 40, G: 90, B: 160, A: 255})
		if batched {
			r.FillRects(rcs, c, full)
		} else {
			for _, rc := range rcs {
				r.fillSprite(rc, c, full)
			}
		}
		return r.ReadPixels(img).Pix
	}
	want, got := draw(false), draw(true)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pixel %d channel %d: got %d, want %d", i/4, i%4, got[i], want[i])
		}
	}
}

// Recorded batches must draw the same as drawing each of them at once: in order, on their own targets,
// also across clears and many batches.
func TestDeferredBatchesOrder(t *testing.T) {
	r := withGL(t)
	full := image.Rect(0, 0, 96, 64)
	imgs := []*noximage.Image16{noximage.NewImage16(full), noximage.NewImage16(full)}
	defer func() {
		for _, img := range imgs {
			r.Forget(img)
		}
	}()
	draw := func(eager bool) [][]byte {
		rnd := rand.New(rand.NewPCG(13, 14))
		for _, img := range imgs {
			r.SetTarget(img)
			r.Clear(color.NRGBA{R: 30, G: 60, B: 90, A: 255})
		}
		for i := 0; i < 3000; i++ {
			x, y := rnd.IntN(90), rnd.IntN(60)
			rc := image.Rect(x, y, x+1+rnd.IntN(8), y+1+rnd.IntN(4))
			c := color.NRGBA{R: uint8(rnd.IntN(256)), G: uint8(rnd.IntN(256)), B: uint8(rnd.IntN(256)), A: uint8(rnd.IntN(256))}
			clip := full
			if rnd.IntN(4) == 0 {
				clip = image.Rect(10, 10, 60, 50)
			}
			switch rnd.IntN(12) {
			case 0:
				r.SetTarget(imgs[rnd.IntN(2)])
			case 1:
				r.Subtract(rc, uint8(rnd.IntN(256)), clip)
			case 2:
				if rnd.IntN(20) == 0 {
					r.Clear(c)
				}
			case 3, 4, 5:
				r.fillSprite(rc, c, clip)
			default:
				r.Fill(rc, c, clip)
			}
			if eager {
				r.submit()
			}
		}
		var out [][]byte
		for _, img := range imgs {
			out = append(out, r.ReadPixels(img).Pix)
		}
		return out
	}
	want, got := draw(true), draw(false)
	for k := range want {
		for i := range want[k] {
			if got[k][i] != want[k][i] {
				t.Fatalf("target %d pixel %d channel %d: got %d, want %d", k, i/4, i%4, got[k][i], want[k][i])
			}
		}
	}
}
