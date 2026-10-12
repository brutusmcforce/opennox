package protection

// Unlink detaches a known member without freeing it or changing its payload.
// The caller owns checksum/count updates and allocation lifetime.
func Unlink(head, tail, r *Record) (*Record, *Record) {
	if r.Prev != nil {
		r.Prev.Next = r.Next
	} else {
		head = r.Next
	}
	if r.Next != nil {
		r.Next.Prev = r.Prev
	} else {
		tail = r.Prev
	}
	return head, tail
}
