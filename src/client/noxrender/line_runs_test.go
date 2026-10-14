package noxrender

import (
	"image"
	"image/color"
	"log/slog"
	"math/rand/v2"
	"testing"
)

// fillRecorder records which pixels Fill covers, and how often.
type fillRecorder struct {
	GPU
	cover map[image.Point]int
	fills int
}

func (g *fillRecorder) Fill(rc image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	g.fills++
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			g.cover[image.Pt(x, y)]++
		}
	}
}

// refLineAlpha is the per-pixel line DrawLineAlpha draws (Bresenham, the last point excluded).
func refLineAlpha(p1, p2 image.Point) []image.Point {
	var out []image.Point
	width, dx := p2.X-p1.X, 1
	if p2.X < p1.X {
		dx, width = -1, -width
	}
	height, dy := p2.Y-p1.Y, 1
	if p2.Y < p1.Y {
		dy, height = -1, -height
	}
	p := p1
	if width < height {
		v := 2*width - height
		for i := 0; i < height; i++ {
			out = append(out, p)
			p.Y += dy
			if v >= 0 {
				p.X += dx
				v += 2 * (width - height)
			} else {
				v += 2 * width
			}
		}
	} else {
		v := 2*height - width
		for i := 0; i < width; i++ {
			out = append(out, p)
			p.X += dx
			if v >= 0 {
				v += 2 * (height - width)
				p.Y += dy
			} else {
				v += 2 * height
			}
		}
	}
	return out
}

// Lines are drawn as runs of pixels; they must cover exactly the pixels of the per-pixel line, once each.
func TestLineRuns(t *testing.T) {
	r := NewRender(slog.New(slog.DiscardHandler), nil)
	data, free := NewRenderData()
	defer free()
	r.SetData(data)
	data.SetAlphaEnabled(true)
	data.SetAlpha(128)
	g := &fillRecorder{}
	r.gpu = g
	rnd := rand.New(rand.NewPCG(1, 2))
	pixels, fills := 0, 0
	for i := 0; i < 2000; i++ {
		p1 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		p2 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		if i%4 == 0 { // shallow lines, like the sight edges
			p2.Y = p1.Y + rnd.IntN(20) - 10
		}
		g.cover, g.fills = make(map[image.Point]int), 0
		r.DrawLineAlpha(p1, p2, color.Black)
		want := refLineAlpha(p1, p2)
		if len(g.cover) != len(want) {
			t.Fatalf("%v-%v: %d pixels, want %d", p1, p2, len(g.cover), len(want))
		}
		for _, p := range want {
			if g.cover[p] != 1 {
				t.Fatalf("%v-%v: pixel %v drawn %d times", p1, p2, p, g.cover[p])
			}
		}
		pixels += len(want)
		fills += g.fills
	}
	t.Logf("%d pixels in %d fills", pixels, fills)
}

// batchRecorder also implements RectFiller, the path the GPU renderer takes.
type batchRecorder struct {
	fillRecorder
}

func (g *batchRecorder) FillRects(rcs []image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	for _, rc := range rcs {
		g.Fill(rc, c, clip)
	}
}

func TestLineRunsBatched(t *testing.T) {
	r := NewRender(slog.New(slog.DiscardHandler), nil)
	data, free := NewRenderData()
	defer free()
	r.SetData(data)
	data.SetAlphaEnabled(true)
	data.SetAlpha(128)
	g := &batchRecorder{}
	r.gpu = g
	rnd := rand.New(rand.NewPCG(5, 6))
	for i := 0; i < 500; i++ {
		p1 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		p2 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		g.cover = make(map[image.Point]int)
		r.DrawLineAlpha(p1, p2, color.Black)
		want := refLineAlpha(p1, p2)
		if len(g.cover) != len(want) {
			t.Fatalf("%v-%v: %d pixels, want %d", p1, p2, len(g.cover), len(want))
		}
		for _, p := range want {
			if g.cover[p] != 1 {
				t.Fatalf("%v-%v: pixel %v drawn %d times", p1, p2, p, g.cover[p])
			}
		}
	}
}

// colorRecorder records the colors filled into each pixel, in order.
type colorRecorder struct {
	GPU
	px map[image.Point][]color.NRGBA
}

func (g *colorRecorder) Fill(rc image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			p := image.Pt(x, y)
			g.px[p] = append(g.px[p], c)
		}
	}
}

func (g *colorRecorder) FillRects(rcs []image.Rectangle, c color.NRGBA, clip image.Rectangle) {
	for _, rc := range rcs {
		g.Fill(rc, c, clip)
	}
}

// DrawLineFade must draw what the loop it replaces (sub_498B50) draws, also when the copies are clipped.
func TestLineFade(t *testing.T) {
	r := NewRender(slog.New(slog.DiscardHandler), nil)
	data, free := NewRenderData()
	defer free()
	r.SetData(data)
	data.SetAlphaEnabled(true)
	data.SetClip(true)
	data.SetClipRect2(image.Rect(20, 20, 380, 380))
	rnd := rand.New(rand.NewPCG(11, 12))
	cl := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	fast := 0
	for i := 0; i < 1000; i++ {
		p1 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		p2 := image.Pt(rnd.IntN(400), rnd.IntN(400))
		if i%2 == 0 { // mostly inside, like the sight edges
			p1 = image.Pt(40+rnd.IntN(320), 40+rnd.IntN(320))
			p2 = image.Pt(40+rnd.IntN(320), 40+rnd.IntN(320))
		}
		step := []image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}[rnd.IntN(4)]
		want := &colorRecorder{px: make(map[image.Point][]color.NRGBA)}
		r.gpu = want
		alpha := byte(208)
		for k := 0; k < 10; k++ {
			data.SetAlpha(alpha)
			alpha -= 20
			r.DrawVector(p1.Add(step.Mul(k)), p2.Sub(p1), cl)
		}
		wantAlpha := data.Alpha()
		got := &colorRecorder{px: make(map[image.Point][]color.NRGBA)}
		r.gpu = got
		r.DrawLineFade(p1, p2, step, 10, 208, 20, cl)
		if data.Alpha() != wantAlpha {
			t.Fatalf("alpha after: %d, want %d", data.Alpha(), wantAlpha)
		}
		if len(got.px) != len(want.px) {
			t.Fatalf("%v-%v: %d pixels, want %d", p1, p2, len(got.px), len(want.px))
		}
		for p, w := range want.px {
			g := got.px[p]
			if len(g) != len(w) {
				t.Fatalf("%v-%v: pixel %v: %v, want %v", p1, p2, p, g, w)
			}
			for j := range w {
				if g[j] != w[j] {
					t.Fatalf("%v-%v: pixel %v: %v, want %v", p1, p2, p, g, w)
				}
			}
		}
		r1 := data.ClipRect2()
		if clipFlags(p1, r1)|clipFlags(p2, r1)|clipFlags(p1.Add(step.Mul(9)), r1)|clipFlags(p2.Add(step.Mul(9)), r1) == 0 {
			fast++
		}
	}
	if fast < 400 {
		t.Fatalf("only %d lines took the fast path", fast)
	}
}
