package protection

import "testing"

func TestBitset(t *testing.T) {
	for _, tc := range []struct {
		index, enabled int32
		want           uint32
	}{
		{0, 1, 1}, {1, -1, 2}, {31, 1, 0x80000000}, {32, 1, 1}, {33, 2, 2}, {2147483647, 1, 0x80000000}, {31, 0, 0},
	} {
		if got := Bit(tc.index, tc.enabled); got != tc.want {
			t.Fatalf("%+v: %08x", tc, got)
		}
	}
	if Flags(nil) != 0 || Flags([]int32{1}) != 0 || Flags([]int32{1, 0, -1}) != 4 {
		t.Fatal("ignored zero or truthy values")
	}
	values := make([]int32, 66)
	values[1], values[33], values[65] = 1, -1, 2
	if Flags(values) != 2 {
		t.Fatal("folded collisions")
	}
	values[32], values[63] = 1, 1
	if Flags(values) != 0x80000003 {
		t.Fatal("word boundaries")
	}
}
