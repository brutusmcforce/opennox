package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/protection"
)

func Nox_xxx_protectionCreateStructForInt_56F280(a1 int, a2 int) int {
	return int(createProtectionRecord(uint32(a1), uint32(a2)))
}
func Nox_xxx_protectionCreateStructForFloat_56F480(a1 int, a2 float32) int {
	return int(createProtectionRecord(uint32(a1), math.Float32bits(a2)))
}

func protectionStringChecksum(data unsafe.Pointer, size uint32) int32 {
	if data == nil {
		return 0
	}
	// Process complete words only. Chunking avoids converting an unsigned C
	// byte count greater than MaxInt into a negative Go slice length on 386.
	remaining := size &^ 3
	p := data
	var sum uint32
	for remaining != 0 {
		n := remaining
		if n > 1<<20 {
			n = 1 << 20
		}
		sum ^= protection.Checksum(unsafe.Slice((*byte)(p), int(n)))
		remaining -= n
		if remaining != 0 {
			p = unsafe.Add(p, n)
		}
	}
	return int32(sum)
}
