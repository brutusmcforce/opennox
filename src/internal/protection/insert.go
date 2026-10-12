package protection

// InsertBefore links a new record before a known non-nil member of the list.
// Existing payloads and the tail are unchanged; the caller updates the count.
func InsertBefore(head, before, r *Record) *Record {
	r.Prev, r.Next = before.Prev, before
	before.Prev = r
	if r.Prev != nil {
		r.Prev.Next = r
	} else {
		head = r
	}
	return head
}
