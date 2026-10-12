// Package protection implements the legacy player/item protection checksum.
package protection

import "encoding/binary"

// Checksum XORs complete little-endian words. Trailing bytes are deliberately
// ignored to preserve the legacy contract; this is not a polynomial CRC.
func Checksum(data []byte) uint32 {
	var sum uint32
	for len(data) >= 4 {
		sum ^= binary.LittleEndian.Uint32(data)
		data = data[4:]
	}
	return sum
}
