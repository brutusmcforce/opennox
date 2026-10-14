package noxrender

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func TestLoadFontClosesFile(t *testing.T) {
	for _, ext := range []string{".ttf", ".otf"} {
		t.Run(ext, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "default")
			require.NoError(t, os.WriteFile(path+ext, goregular.TTF, 0600))
			face, err := loadFont(path, 10)
			require.NoError(t, err)
			defer face.Close()

			// Windows refuses removal while the source file is still open.
			require.NoError(t, os.Remove(path+ext))
			advance, ok := face.GlyphAdvance('A')
			require.True(t, ok, "the font must remain usable after removing its source")
			require.Positive(t, advance)
		})
	}
}

// The ASCII caches of outline faces must give the same results as the font.
func TestOutlineFaceCaches(t *testing.T) {
	data, err := bundledFonts.ReadFile(bundledFontText)
	require.NoError(t, err)
	fnt, err := opentype.Parse(data)
	require.NoError(t, err)
	f, err := newOutlineFace(fnt, 12, basicfont.Face7x13)
	require.NoError(t, err)
	defer func() { _ = f.Close() }()
	ref, err := opentype.NewFace(fnt, &opentype.FaceOptions{Size: 12, DPI: 72, Hinting: font.HintingNone})
	require.NoError(t, err)
	defer func() { _ = ref.Close() }()
	var buf sfnt.Buffer
	has := func(r rune) bool {
		x, err := fnt.GlyphIndex(&buf, r)
		return err == nil && x != 0
	}
	runes := []rune{-1, 0, 'A', 'V', 'a', 'w', '.', 0x7f, 0x80, 'é', 'Ж', 0x2603}
	for r := rune(32); r < 127; r++ {
		runes = append(runes, r)
	}
	for pass := 0; pass < 2; pass++ { // filling the caches, then from the caches
		kerned := 0
		for _, r0 := range runes {
			require.Equal(t, has(r0), f.hasGlyph(r0), "hasGlyph(%q)", r0)
			for _, r1 := range runes {
				want := fixed.Int26_6(0)
				if has(r0) && has(r1) {
					want = ref.Kern(r0, r1)
				}
				if want != 0 {
					kerned++
				}
				require.Equal(t, want, f.Kern(r0, r1), "Kern(%q, %q)", r0, r1)
			}
		}
		require.Positive(t, kerned)
	}
}
