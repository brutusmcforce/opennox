package gpu

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	atlasSize = 2048
	atlasPad  = 1
)

// atlasPage is a pair of textures: RGBA colors and material indexes.
type atlasPage struct {
	size    int
	tex     uint32
	mat     uint32
	shelves []shelf
	nextY   int
}

type shelf struct {
	y, h, x int
}

// texRegion is a region of an atlas page.
type texRegion struct {
	page *atlasPage
	rect image.Rectangle
}

func newAtlasPage(size int) *atlasPage {
	p := &atlasPage{size: size}
	gl.GenTextures(1, &p.tex)
	gl.BindTexture(gl.TEXTURE_2D, p.tex)
	setTexParams(gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(size), int32(size), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.GenTextures(1, &p.mat)
	gl.BindTexture(gl.TEXTURE_2D, p.mat)
	setTexParams(gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R8, int32(size), int32(size), 0, gl.RED, gl.UNSIGNED_BYTE, nil)
	return p
}

func (p *atlasPage) free() {
	gl.DeleteTextures(1, &p.tex)
	gl.DeleteTextures(1, &p.mat)
}

func setTexParams(filter int32) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}

// alloc reserves a region in the page using a simple shelf packing.
func (p *atlasPage) alloc(w, h int) (image.Rectangle, bool) {
	pw, ph := w+2*atlasPad, h+2*atlasPad
	if pw > p.size || ph > p.size {
		return image.Rectangle{}, false
	}
	best := -1
	for i, s := range p.shelves {
		if s.h >= ph && s.x+pw <= p.size && (best < 0 || s.h < p.shelves[best].h) {
			best = i
		}
	}
	if best < 0 || p.shelves[best].h > ph*2 {
		// open a new shelf, if possible
		if p.nextY+ph <= p.size {
			p.shelves = append(p.shelves, shelf{y: p.nextY, h: ph})
			p.nextY += ph
			best = len(p.shelves) - 1
		} else if best < 0 {
			return image.Rectangle{}, false
		}
	}
	s := &p.shelves[best]
	r := image.Rect(s.x+atlasPad, s.y+atlasPad, s.x+atlasPad+w, s.y+atlasPad+h)
	s.x += pw
	return r, true
}

func (p *atlasPage) upload(r image.Rectangle, rgba, mat []byte) {
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.BindTexture(gl.TEXTURE_2D, p.tex)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba))
	gl.BindTexture(gl.TEXTURE_2D, p.mat)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()), gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(mat))
}

// atlas manages multiple atlas pages.
type atlas struct {
	pages []*atlasPage
	big   []*atlasPage
}

// add uploads an image to the atlas and returns its region.
func (a *atlas) add(w, h int, rgba, mat []byte) texRegion {
	if w+2*atlasPad > atlasSize || h+2*atlasPad > atlasSize {
		// dedicated texture
		p := newAtlasPage(max(w, h) + 2*atlasPad)
		a.big = append(a.big, p)
		r, _ := p.alloc(w, h)
		p.upload(r, rgba, mat)
		return texRegion{page: p, rect: r}
	}
	for _, p := range a.pages {
		if r, ok := p.alloc(w, h); ok {
			p.upload(r, rgba, mat)
			return texRegion{page: p, rect: r}
		}
	}
	p := newAtlasPage(atlasSize)
	a.pages = append(a.pages, p)
	r, _ := p.alloc(w, h)
	p.upload(r, rgba, mat)
	return texRegion{page: p, rect: r}
}

func (a *atlas) free() {
	for _, p := range a.pages {
		p.free()
	}
	for _, p := range a.big {
		p.free()
	}
	a.pages, a.big = nil, nil
}
