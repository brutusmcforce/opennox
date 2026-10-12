package protection

// Bit returns the legacy folded bit for an enabled spell/ability. IDs are
// nonnegative; indices separated by 32 intentionally share a bit.
func Bit(index int32, enabled int32) uint32 {
	if enabled == 0 {
		return 0
	}
	return uint32(1) << (uint32(index) & 31)
}

// Flags folds nonzero entries into 32 bits. Entry zero is deliberately ignored.
func Flags(values []int32) uint32 {
	var out uint32
	for i := 1; i < len(values); i++ {
		out |= Bit(int32(i), values[i])
	}
	return out
}
