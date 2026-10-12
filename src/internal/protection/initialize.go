package protection

// Initialize fills a newly allocated record and incorporates its encrypted
// payload into the manager checksum. Nil models a failed allocation and must
// leave the checksum untouched. Insertion into the list remains the caller's job.
func Initialize(r *Record, id, value, key uint32, sum *uint32) bool {
	if r == nil {
		return false
	}
	*r = Record{ID: id ^ key, Value: value ^ key}
	*sum ^= r.ID ^ r.Value
	return true
}
