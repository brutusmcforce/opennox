package memmap

import "testing"

// BlobByAddr's table of blob starts must find the same blobs as searching them.
func TestBlobByAddr(t *testing.T) {
	saved, savedBases := blobs, blobBases
	defer func() { blobs, blobBases = saved, savedBases }()
	blobs, blobBases = nil, [256]*Blob{}
	RegisterBlob(0x587000, "a", 0x100)
	RegisterBlob(0x587100, "b", 0x10) // right after a
	RegisterBlob(0x5D4594, "c", 0x1000)
	RegisterBlob(0x600094, "d", 8) // the same low byte as c
	search := func(addr uintptr) *Blob {
		for _, b := range blobs {
			if b.Contains(addr) {
				return b
			}
		}
		return nil
	}
	for _, base := range []uintptr{0x587000, 0x587100, 0x5D4594, 0x600094} {
		for _, d := range []uintptr{0, 1, 0xff, 0x100, 0x101, 0x1000} {
			for _, addr := range []uintptr{base + d, base - d} {
				if got, want := BlobByAddr(addr), search(addr); got != want {
					t.Fatalf("0x%x: %v, want %v", addr, got, want)
				}
			}
		}
	}
}
