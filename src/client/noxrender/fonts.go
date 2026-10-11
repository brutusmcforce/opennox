package noxrender

import (
	"embed"
	"image"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Bundled outline fonts that replace the original bitmap fonts. See fonts/OFL-*.txt for licenses.
//
//go:embed fonts/*.ttf
var bundledFonts embed.FS

const (
	bundledFontText  = "fonts/FiraSansCondensed-Regular.ttf"
	bundledFontLarge = "fonts/Philosopher-Regular.ttf"
)

// SmoothFonts enables outline fonts instead of the original bitmap fonts. It must be set before loading the fonts.
//
// Outline fonts can be drawn at a higher resolution than the GUI, so the text stays smooth when the GUI is scaled.
var SmoothFonts = true

// outlineFace is a font face backed by an outline (TrueType/OpenType) font.
//
// When it replaces an original bitmap font, it keeps the metrics of the original font so the text layout stays the same,
// and the original font is used for characters missing in the outline font.
type outlineFace struct {
	fnt  *opentype.Font
	size float64
	face font.Face // outline font at the base size
	orig font.Face // original bitmap font; optional
	// yOff moves the outline glyphs vertically, so that they are placed the same way as the original glyphs.
	yOff fixed.Int26_6
	// lineGap is an additional space between the lines of wrapped text.
	lineGap int
	// sharedOrig is set if the original font is owned by another face.
	sharedOrig bool

	mu     sync.Mutex
	buf    sfnt.Buffer
	scaled map[float64]font.Face
}

func newOutlineFace(fnt *opentype.Font, size float64, orig font.Face) (*outlineFace, error) {
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: size, DPI: 72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}
	return &outlineFace{fnt: fnt, size: size, face: face, orig: orig}, nil
}

// newReplacementFace creates an outline face that replaces an original bitmap font.
// The size is selected so that none of the sample texts is wider than with the original font.
func newReplacementFace(data []byte, orig font.Face, samples []string) (*outlineFace, error) {
	fnt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	const refSize = 100
	ref, err := opentype.NewFace(fnt, &opentype.FaceOptions{Size: refSize, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		return nil, err
	}
	defer ref.Close()
	size := 0.0
	for _, sample := range samples {
		sz := refSize * float64(font.MeasureString(orig, sample)) / float64(font.MeasureString(ref, sample))
		if size == 0 || sz < size {
			size = sz
		}
	}
	// letters must not be taller than the original ones either
	rr := refRune(orig)
	if top, bottom, ok := glyphRows(orig, rr); ok {
		if b, _, ok := ref.GlyphBounds(rr); ok && b.Max.Y > b.Min.Y {
			size = min(size, refSize*float64(fixed.I(bottom-top))/float64(b.Max.Y-b.Min.Y))
		}
	}
	f, err := newOutlineFace(fnt, size, orig)
	if err != nil {
		return nil, err
	}
	f.yOff = capCenterOffset(orig, f.face, rr)
	f.lineGap = 1
	return f, nil
}

// derive creates a copy of the face, scaled by a given factor.
func (f *outlineFace) derive(scale float64) (*outlineFace, error) {
	d, err := newOutlineFace(f.fnt, f.size*scale, f.orig)
	if err != nil {
		return nil, err
	}
	d.lineGap = f.lineGap
	d.sharedOrig = true
	if f.orig != nil {
		d.yOff = capCenterOffset(f.orig, d.face, refRune(f.orig))
	}
	return d, nil
}

// refRune returns a reference character for comparing the outline font with the original font:
// a capital letter, or a digit for fonts that have only digits.
func refRune(orig font.Face) rune {
	if _, _, ok := glyphRows(orig, 'H'); ok {
		return 'H'
	}
	return '0'
}

// capCenterOffset returns a vertical offset for the outline font that puts the center of the reference character
// at the same height as the center of the same character in the original font.
//
// Glyphs of original bitmap fonts are placed above the dot, including descenders,
// so their baseline is higher than the one of outline fonts.
func capCenterOffset(orig, outline font.Face, rr rune) fixed.Int26_6 {
	top, bottom, ok := glyphRows(orig, rr)
	if !ok {
		return 0
	}
	b, _, ok := outline.GlyphBounds(rr)
	if !ok {
		return 0
	}
	origCenter := fixed.I(top+bottom) / 2
	return origCenter - (b.Min.Y+b.Max.Y)/2
}

// glyphRows returns the vertical range of non-transparent pixels of a glyph, relative to the dot.
func glyphRows(f font.Face, r rune) (top, bottom int, ok bool) {
	dr, mask, maskp, _, ok := f.Glyph(fixed.Point26_6{}, r)
	if !ok || mask == nil {
		return 0, 0, false
	}
	top, bottom = dr.Max.Y, dr.Min.Y
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		for x := dr.Min.X; x < dr.Max.X; x++ {
			if _, _, _, a := mask.At(maskp.X+x-dr.Min.X, maskp.Y+y-dr.Min.Y).RGBA(); a != 0 {
				top, bottom = min(top, y), max(bottom, y+1)
				break
			}
		}
	}
	return top, bottom, top < bottom
}

// hasGlyph checks if the outline font has a glyph for a given rune.
func (f *outlineFace) hasGlyph(r rune) bool {
	if f.orig == nil {
		return true
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	x, err := f.fnt.GlyphIndex(&f.buf, r)
	return err == nil && x != 0
}

func (f *outlineFace) pick(r rune) font.Face {
	if f.hasGlyph(r) {
		return f.face
	}
	return f.orig
}

// scaledFace returns the outline font scaled by a given factor, for drawing at a higher resolution.
func (f *outlineFace) scaledFace(scale float64) font.Face {
	f.mu.Lock()
	defer f.mu.Unlock()
	if face, ok := f.scaled[scale]; ok {
		return face
	}
	face, err := opentype.NewFace(f.fnt, &opentype.FaceOptions{
		Size: f.size * scale, DPI: 72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		face = nil
	}
	if f.scaled == nil {
		f.scaled = make(map[float64]font.Face)
	}
	f.scaled[scale] = face
	return face
}

func (f *outlineFace) Close() error {
	f.mu.Lock()
	for _, face := range f.scaled {
		if face != nil {
			face.Close()
		}
	}
	f.scaled = nil
	f.mu.Unlock()
	if f.orig != nil && !f.sharedOrig {
		f.orig.Close()
	}
	return f.face.Close()
}

func (f *outlineFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if !f.hasGlyph(r) {
		return f.orig.Glyph(dot, r)
	}
	dot.Y += f.yOff
	return f.face.Glyph(dot, r)
}

func (f *outlineFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if !f.hasGlyph(r) {
		return f.orig.GlyphBounds(r)
	}
	bounds, advance, ok = f.face.GlyphBounds(r)
	bounds.Min.Y += f.yOff
	bounds.Max.Y += f.yOff
	return bounds, advance, ok
}

func (f *outlineFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return f.pick(r).GlyphAdvance(r)
}

func (f *outlineFace) Kern(r0, r1 rune) fixed.Int26_6 {
	if f.hasGlyph(r0) && f.hasGlyph(r1) {
		return f.face.Kern(r0, r1)
	}
	return 0
}

func (f *outlineFace) Metrics() font.Metrics {
	if f.orig != nil {
		return f.orig.Metrics()
	}
	return f.face.Metrics()
}
