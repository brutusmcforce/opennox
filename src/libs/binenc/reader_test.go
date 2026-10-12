package binenc

import "testing"

func TestReaderHugeLength(t *testing.T) {
	// a 32-bit length with the high bit set is negative as an int on 32-bit systems
	data := []byte{0x00, 0x00, 0x00, 0x80, 1, 2, 3, 4}
	if _, ok := NewReader(data).ReadBytes32(); ok {
		t.Fatal("expected an error")
	}
	if _, ok := NewReader(data).ReadString32(); ok {
		t.Fatal("expected an error")
	}
	r := NewReader(data)
	if _, ok := r.ReadWString(-1); ok {
		t.Fatal("expected an error")
	}
	if _, ok := r.ReadWString(1 << 30); ok {
		t.Fatal("expected an error")
	}
	if _, ok := r.ReadNext(-1); ok {
		t.Fatal("expected an error")
	}
}
