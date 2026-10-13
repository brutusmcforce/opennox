package platform

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestRandInt31(t *testing.T) {
	a, b := rand.New(rand.NewSource(1234)), rand.New(rand.NewSource(1234))
	for i := 0; i < 1000; i++ {
		got := RandInt31(a)
		if got < 0 || int64(got) >= 1<<31 {
			t.Fatalf("RandInt31 = %d, out of the 31-bit range", got)
		}
		if strconv.IntSize == 32 {
			// RandInt31 must keep the sequence that rand.Int gives on 32-bit platforms.
			if want := b.Int(); got != want {
				t.Fatalf("draw %d: RandInt31 = %d, rand.Int = %d", i, got, want)
			}
		} else if want := int(uint32(b.Int63()) & 0x7fffffff); got != want {
			t.Fatalf("draw %d: RandInt31 = %d, want %d", i, got, want)
		}
	}
}
