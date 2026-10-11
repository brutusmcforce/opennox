package gpu

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func spriteHeader(w, h int) []byte {
	b := make([]byte, spriteHeaderSize)
	binary.LittleEndian.PutUint32(b[0:], uint32(w))
	binary.LittleEndian.PutUint32(b[4:], uint32(h))
	return b
}

func TestDecodeSpriteRLE(t *testing.T) {
	data := spriteHeader(4, 1)
	// 1 transparent pixel, 1 RGB555 pixel (pure red), 2 grayscale pixels with material 3
	data = append(data, 0x01, 1)
	data = append(data, 0x02, 1, 0x00, 0x7C)
	data = append(data, 0x34, 2, 0x80, 0xFF)
	d, err := decodeSprite(3, data)
	require.NoError(t, err)
	require.Equal(t, 4, d.W)
	require.Equal(t, 1, d.H)
	require.Equal(t, []byte{
		0, 0, 0, 0,
		0xFF, 0, 0, 0xFF, // same conversion as SplitColor16
		0x80, 0x80, 0x80, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
	}, d.RGBA)
	require.Equal(t, []byte{0, 0, 4, 4}, d.Mat)
}

func TestDecodeSpriteEmpty(t *testing.T) {
	_, err := decodeSprite(3, spriteHeader(0, 0))
	require.ErrorIs(t, err, errEmptySprite)
}

func TestDecodeSpriteTruncated(t *testing.T) {
	data := append(spriteHeader(4, 1), 0x02, 4, 0x00)
	_, err := decodeSprite(3, data)
	require.Error(t, err)
}
