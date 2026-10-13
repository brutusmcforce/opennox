package noxfont

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/brutusmcforce/opennox/libs/noxtest"
)

func TestFonts(t *testing.T) {
	path := noxtest.DataPath(t)
	hashes := []string{
		"34f4189ac58fd184b9fb71fe6300dfaf",
		"8e2d5ee796c8c48043ddd7e4e02194d3",
		"34f4189ac58fd184b9fb71fe6300dfaf",
		"d9dd152bdb61bd06e11c4089609fe730",
	}
	for i, name := range []string{
		DefaultFile, LargeFile, SmallFile, NumbersFile,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(path, name+Ext)
			data, err := os.ReadFile(path)
			must.NoError(t, err)
			fnt, err := Decode(bytes.NewReader(data))
			must.NoError(t, err)
			t.Run("encode", func(t *testing.T) {
				got, err := fnt.Encode()
				must.NoError(t, err)
				must.Eq(t, data, got)
			})
			t.Run("draw", func(t *testing.T) {
				text := "The quick brown fox jumps over the lazy dog."
				if name == NumbersFile {
					text = "0123456789"
				}
				b, _ := font.BoundString(fnt, text)
				dy := (b.Max.Y - b.Min.Y).Ceil()
				img := image.NewGray(image.Rect(0, 0, b.Max.X.Ceil(), dy))
				dr := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(color.White),
					Face: fnt,
					Dot:  fixed.P(0, dy),
				}
				dr.DrawString(text)
				noxtest.WritePNG(t, name+".png", img, hashes[i])
			})
		})
	}
}
