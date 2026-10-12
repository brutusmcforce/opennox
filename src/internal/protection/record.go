package protection

// Record matches the four-word legacy protection record on 386. Legacy owns
// its C-heap allocation and list links until the manager itself is converted.
type Record struct {
	ID, Value  uint32
	Next, Prev *Record
}

// Find returns the first record with the requested decoded ID.
func Find(head *Record, key, id uint32) *Record {
	for p := head; p != nil; p = p.Next {
		if p.ID == id^key {
			return p
		}
	}
	return nil
}

// At returns a record by zero-based list position, or nil for a missing index.
func At(head *Record, index int32) *Record {
	if index < 0 {
		return nil
	}
	for p := head; p != nil; p = p.Next {
		if index == 0 {
			return p
		}
		index--
	}
	return nil
}

// Swap exchanges payloads, never links. A self-swap is still successful, as
// the legacy caller counts it. Missing operands leave both records unchanged.
func Swap(a, b *Record) bool {
	if a == nil || b == nil {
		return false
	}
	a.ID, b.ID = b.ID, a.ID
	a.Value, b.Value = b.Value, a.Value
	return true
}
