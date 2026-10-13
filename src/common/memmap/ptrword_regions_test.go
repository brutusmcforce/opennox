package memmap

import (
	"testing"
	"unsafe"
)

func TestWordRegions(t *testing.T) {
	a, b := make([]byte, 100), make([]byte, 64)
	pa, pb := unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0])
	RegisterWordRegion(pa, 100)
	RegisterWordRegion(pb, 64)
	defer UnregisterWordRegion(pa)
	defer UnregisterWordRegion(pb)
	wa, wb := PtrWord(pa), PtrWord(pb)
	if wa&handleTag == handleTag || wa < regionStart {
		t.Fatalf("not a region word: 0x%X", wa)
	}
	// Word arithmetic stays inside the block, including its end.
	for _, off := range []uint32{0, 1, 76, 99, 100} {
		if got, want := WordPtr(wa+off), unsafe.Add(pa, off); got != want {
			t.Fatalf("WordPtr(base+%d) = %p, want %p", off, got, want)
		}
		if got := PtrWord(unsafe.Add(pa, off)); got != wa+off {
			t.Fatalf("PtrWord(base+%d) = 0x%X, want 0x%X", off, got, wa+off)
		}
	}
	if WordPtr(wb+10) != unsafe.Add(pb, 10) {
		t.Fatal("second region")
	}
	if wb >= wa && wb <= wa+100 || wa >= wb && wa <= wb+64 {
		t.Fatalf("pseudo ranges overlap: 0x%X 0x%X", wa, wb)
	}
	UnregisterWordRegion(pb)
	if w := PtrWord(pb); w&handleTag != handleTag {
		t.Fatalf("unregistered block still uses a region word: 0x%X", w)
	}
}
