package opennox

import (
	"encoding/binary"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

// testCursorSprite creates a 2x2 white RGB555 sprite with a given draw offset.
func testCursorSprite(off image.Point) *noxrender.Image {
	data := make([]byte, 17+2*4)
	binary.LittleEndian.PutUint32(data[0:], 2)
	binary.LittleEndian.PutUint32(data[4:], 2)
	binary.LittleEndian.PutUint32(data[8:], uint32(int32(off.X)))
	binary.LittleEndian.PutUint32(data[12:], uint32(int32(off.Y)))
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint16(data[17+2*i:], 0x7FFF)
	}
	return noxrender.NewRawImage(2, data)
}

func TestCursorImage(t *testing.T) {
	img := testCursorSprite(image.Pt(cursorSize-1, cursorSize))
	white := [3]uint8{255, 255, 255}

	c, hot, err := cursorImage(sysCursorKey{img: img, mul: white, scale: 100})
	require.NoError(t, err)
	require.Equal(t, image.Pt(2, 2), c.Bounds().Size())
	require.Equal(t, image.Pt(1, 0), hot)

	// red tint, as used for hostile targets
	c, _, err = cursorImage(sysCursorKey{img: img, mul: [3]uint8{255, 0, 0}, scale: 100})
	require.NoError(t, err)
	require.Equal(t, color.NRGBA{R: 255, A: 255}, color.NRGBAModel.Convert(c.At(0, 0)))

	// scaled for stretched menus
	c, hot, err = cursorImage(sysCursorKey{img: img, mul: white, scale: 300})
	require.NoError(t, err)
	require.Equal(t, image.Pt(6, 6), c.Bounds().Size())
	require.Equal(t, image.Pt(3, 0), hot)

	c, _, err = cursorImage(sysCursorKey{img: img, mul: white, scale: 250})
	require.NoError(t, err)
	require.Equal(t, image.Pt(5, 5), c.Bounds().Size())
}
