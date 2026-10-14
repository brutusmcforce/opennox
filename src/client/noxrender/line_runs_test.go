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
