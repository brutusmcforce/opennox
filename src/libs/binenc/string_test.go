package binenc

import "testing"

// A string filling its whole field (no terminator) decodes without junk.
func TestStringDecodeFull(t *testing.T) {
	var s String
	s.Decode([]byte("abcd"))
	if s.Value != "abcd" || s.Junk != nil {
		t.Fatalf("got %q %v", s.Value, s.Junk)
	}
	s.Decode([]byte("ab\x00x"))
	if s.Value != "ab" || string(s.Junk) != "x" {
		t.Fatalf("got %q %q", s.Value, s.Junk)
	}
}
