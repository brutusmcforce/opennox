package legacy

import "math"

// floatToInt32 preserves the masked x87 truncation result using integer bits.
// All remaining callers are Go; no C converter is retained.
func floatToInt32(value float32) int32 {
	bits := math.Float32bits(value)
	exponent := int((bits>>23)&255) - 127
	if exponent < 0 {
		return 0
	}
	if exponent >= 31 {
		return -2147483648
	}
	magnitude := uint32(0x800000) | (bits & 0x7fffff)
	if exponent >= 23 {
		magnitude <<= uint(exponent - 23)
	} else {
		magnitude >>= uint(23 - exponent)
	}
	if bits>>31 != 0 {
		return -int32(magnitude)
	}
	return int32(magnitude)
}
