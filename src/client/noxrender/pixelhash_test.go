package noxrender

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"image"
	"testing"

	"github.com/brutusmcforce/opennox/libs/noximage"
	"github.com/stretchr/testify/require"
)

// pixelHash16 hashes dimensions and visible framebuffer words in little-endian
// order. It is independent of PNG encoding, color-library expansion and padding.
func pixelHash16(p *noximage.Image16) string {
	h := sha256.New()
	var size [8]byte
	binary.LittleEndian.PutUint32(size[:4], uint32(p.Rect.Dx()))
	binary.LittleEndian.PutUint32(size[4:], uint32(p.Rect.Dy()))
	h.Write(size[:])
	var word [2]byte
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		start := p.PixOffset(p.Rect.Min.X, y)
		for _, v := range p.Pix[start : start+p.Rect.Dx()] {
			binary.LittleEndian.PutUint16(word[:], v)
			h.Write(word[:])
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func TestPixelHash16(t *testing.T) {
	p := noximage.NewImage16(image.Rect(0, 0, 2, 2))
	copy(p.Pix, []uint16{0, 1, 0x8000, 0xffff})
	want := pixelHash16(p)
	p.Pix[2] ^= 1
	require.NotEqual(t, want, pixelHash16(p), "one framebuffer bit must change the hash")
	p.Pix[2] ^= 1
	// A subimage has a nonzero origin and padded stride, but identical visible data.
	padded := noximage.NewImage16(image.Rect(0, 0, 4, 4))
	sub := padded.SubImage(image.Rect(1, 1, 3, 3))
	copy(sub.Pix[:2], p.Pix[:2])
	copy(sub.Pix[4:6], p.Pix[2:])
	require.Equal(t, want, pixelHash16(sub))
	row := noximage.NewImage16(image.Rect(0, 0, 4, 1))
	copy(row.Pix, p.Pix)
	require.NotEqual(t, want, pixelHash16(row), "dimensions are part of the reference")
}
