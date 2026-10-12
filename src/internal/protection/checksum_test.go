package protection

import "testing"

func TestChecksum(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want uint32
	}{
		{"nil", nil, 0},
		{"empty", []byte{}, 0},
		{"one", []byte{255}, 0},
		{"two", []byte{255, 255}, 0},
		{"three", []byte{255, 255, 255}, 0},
		{"little_endian", []byte{1, 2, 3, 4}, 0x04030201},
		{"sign_bit", []byte{0, 0, 0, 128}, 0x80000000},
		{"all_bits", []byte{255, 255, 255, 255}, 0xffffffff},
		{"xor", []byte{1, 2, 3, 4, 255, 255, 255, 255}, 0xfbfcfdfe},
		{"cancel", []byte{1, 2, 3, 4, 1, 2, 3, 4}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Checksum(tc.data); got != tc.want {
				t.Fatalf("got %08x, want %08x", got, tc.want)
			}
		})
	}
}

func FuzzChecksum(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 128, 9, 8, 7})
	f.Fuzz(func(t *testing.T, data []byte) {
		whole := len(data) &^ 3
		got := Checksum(data)
		if got != Checksum(data[:whole]) {
			t.Fatal("trailing partial word affected checksum")
		}
		// XORing a word twice must cancel, including high-bit values.
		pair := append(append([]byte{}, data[:whole]...), 0x01, 0x80, 0xfe, 0xff, 0x01, 0x80, 0xfe, 0xff)
		if Checksum(pair) != got {
			t.Fatal("duplicate word did not cancel")
		}
	})
}
