// Package noxtest implements different test helper functions.
package noxtest

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/datapath"
)

func DataPath(t testing.TB, sub ...string) string {
	path := datapath.Data(sub...)
	if path == "" || path == "." || !datapath.Found() {
		t.Skip("cannot detect Nox path and NOX_DATA is not set")
	}
	return path
}

func WritePNG(t testing.TB, path string, img image.Image, exp string) {
	f, err := os.Create(path)
	must.NoError(t, err)
	defer f.Close()
	err = png.Encode(f, img)
	must.NoError(t, err)
	got := ImageHash(img)
	if exp != "" {
		must.EqOp(t, exp, got)
	} else {
		t.Logf("%s: %s", path, got)
	}
}

// ImageHash returns a hash of the image's size and pixels (as non-premultiplied RGBA), which unlike its encoded
// form does not change with the encoder.
func ImageHash(img image.Image) string {
	b := img.Bounds()
	rgba := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Rect, img, b.Min, draw.Src)
	h := md5.New()
	_ = binary.Write(h, binary.LittleEndian, [2]int32{int32(b.Dx()), int32(b.Dy())})
	h.Write(rgba.Pix)
	return hex.EncodeToString(h.Sum(nil))
}
