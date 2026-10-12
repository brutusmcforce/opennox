package server

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestDecodePlayerInputShort(t *testing.T) {
	// truncated input from a client must not panic
	for _, data := range [][]byte{
		{0xFF},
		{0x01, 0x00, 0x00},
		{0x01, 0, 0, 0}, // Orientation with its argument cut off
	} {
		_ = netDecodePlayerInput(data, nil)
	}
	must.Len(t, 1, netDecodePlayerInput([]byte{0x05, 0, 0, 0, 0x01}, nil)) // MoveRight with its 1-byte argument
}
