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

func TestReaderAlign(t *testing.T) {
	data := make([]byte, 32)
	for _, c := range []struct {
		pos, off, want int
	}{
		{0, 0, 0},  // aligned: skip nothing
		{8, 0, 8},  // aligned
		{2, 0, 8},  // skip 6
		{2, 2, 6},  // stream offset 4: skip 4
		{6, 2, 6},  // stream offset 8: aligned
		{7, 2, 14}, // stream offset 9: skip 7
	} {
		r := NewReader(data)
		r.ReadNext(c.pos)
		if !r.Align(c.off) {
			t.Fatalf("pos %d off %d: align failed", c.pos, c.off)
		}
		if got := len(data) - r.Remaining(); got != c.want {
			t.Errorf("pos %d off %d: at %d, want %d", c.pos, c.off, got, c.want)
		}
	}
	// aligned at the end of the data: nothing to skip, not an error
	r := NewReader(data)
	r.ReadNext(32)
	if !r.Align(0) {
		t.Error("align at the end failed")
	}
}
