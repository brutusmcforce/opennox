package protection

// Rekey changes the encoding of each payload and returns the rebuilt checksum.
// It preserves the decoded values and all list links.
func Rekey(head *Record, oldKey, newKey uint32) uint32 {
	delta := oldKey ^ newKey
	sum := ^newKey
	for r := head; r != nil; r = r.Next {
		r.ID ^= delta
		r.Value ^= delta
		sum ^= r.ID ^ r.Value
	}
	return sum
}
