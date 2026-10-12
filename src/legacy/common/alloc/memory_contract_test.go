package alloc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"
)

// Every contract invokes the public helper on bounded, valid storage. Frozen
// comparison hashes normalize only the nonportable nonzero return magnitude.
func TestMemoryFillCopyContracts(t *testing.T) {
	for off := 0; off < 16; off++ {
		for n := 0; n <= 128; n++ {
			for _, v := range []byte{0, 1, 0x7f, 0x80, 0xff} {
				dst := bytes.Repeat([]byte{0xa5}, 160)
				want := bytes.Clone(dst)
				for i := off; i < off+n; i++ {
					want[i] = v
				}
				ptr := unsafe.Pointer(&dst[off])
				if got := Memset(ptr, v, uintptr(n)); got != ptr || !bytes.Equal(dst, want) {
					t.Fatalf("fill off=%d n=%d byte=%d", off, n, v)
				}
			}
			for _, srcOff := range []int{0, 1, 7, 15} {
				src := make([]byte, 160)
				for i := range src {
					src[i] = byte(i*37 + 19)
				}
				before := bytes.Clone(src)
				dst := bytes.Repeat([]byte{0xa5}, 160)
				want := bytes.Clone(dst)
				copy(want[off:off+n], src[srcOff:srcOff+n])
				ptr := unsafe.Pointer(&dst[off])
				if got := Memcpy(ptr, unsafe.Pointer(&src[srcOff]), uintptr(n)); got != ptr || !bytes.Equal(dst, want) || !bytes.Equal(src, before) {
					t.Fatalf("copy dst=%d src=%d n=%d", off, srcOff, n)
				}
			}
		}
	}

	for _, n := range []int{255, 256, 1023, 4096, 65535} {
		src := make([]byte, n+32)
		dst := bytes.Repeat([]byte{0xa5}, n+32)
		for i := range src {
			src[i] = byte(i*53 + 31)
		}
		want := bytes.Clone(dst)
		copy(want[15:15+n], src[1:1+n])
		before := bytes.Clone(src)
		ptr := unsafe.Pointer(&dst[15])
		if Memcpy(ptr, unsafe.Pointer(&src[1]), uintptr(n)) != ptr || !bytes.Equal(dst, want) || !bytes.Equal(src, before) {
			t.Fatalf("large copy n=%d", n)
		}
		for i := 15; i < 15+n; i++ {
			want[i] = 0xff
		}
		if Memset(ptr, 0xff, uintptr(n)) != ptr || !bytes.Equal(dst, want) {
			t.Fatalf("large fill n=%d", n)
		}
	}
}

func TestStringCopyAppendContracts(t *testing.T) {
	for _, dstOff := range []int{0, 1, 3, 7, 15} {
		for _, srcOff := range []int{0, 1, 3, 7, 15} {
			for _, n := range []int{0, 1, 2, 7, 15, 16, 31, 32, 63, 64} {
				src := bytes.Repeat([]byte{0x5a}, 160)
				for i := 0; i < n; i++ {
					src[srcOff+i] = byte((i*37+17)%255 + 1)
				}
				src[srcOff+n] = 0
				before := bytes.Clone(src)
				dst := bytes.Repeat([]byte{0xa5}, 192)
				want := bytes.Clone(dst)
				copy(want[dstOff:], src[srcOff:srcOff+n+1])
				ptr := unsafe.Pointer(&dst[dstOff])
				if Strcpy(ptr, unsafe.Pointer(&src[srcOff])) != ptr || !bytes.Equal(dst, want) || !bytes.Equal(src, before) {
					t.Fatalf("strcpy dst=%d src=%d n=%d", dstOff, srcOff, n)
				}
				for _, prefix := range []int{0, 1, 7, 31, 64} {
					dst = bytes.Repeat([]byte{0xa5}, 192)
					for i := 0; i < prefix; i++ {
						dst[dstOff+i] = 0x80
					}
					dst[dstOff+prefix] = 0
					want = bytes.Clone(dst)
					copy(want[dstOff+prefix:], src[srcOff:srcOff+n+1])
					ptr = unsafe.Pointer(&dst[dstOff])
					if Strcat(ptr, unsafe.Pointer(&src[srcOff])) != ptr || !bytes.Equal(dst, want) || !bytes.Equal(src, before) {
						t.Fatalf("strcat dst=%d src=%d n=%d prefix=%d", dstOff, srcOff, n, prefix)
					}
				}
			}
		}
	}
}

func TestMemoryCompareContracts(t *testing.T) {
	const wantHash = "3973015d60c667cee95867d663027cd09a040df0a6cef253faf482792440cb29" // Freeze from original libc after sign normalization.
	h := sha256.New()
	var encoded [4]byte
	count := 0
	check := func(a, b []byte, offA, offB, n int) {
		t.Helper()
		beforeA, beforeB := bytes.Clone(a), bytes.Clone(b)
		got := Memcmp(unsafe.Pointer(&a[offA]), unsafe.Pointer(&b[offB]), uintptr(n))
		order := bytes.Compare(a[offA:offA+n], b[offB:offB+n])
		if (got < 0) != (order < 0) || (got > 0) != (order > 0) || !bytes.Equal(a, beforeA) || !bytes.Equal(b, beforeB) {
			t.Fatalf("memcmp off=%d/%d n=%d got=%d order=%d", offA, offB, n, got, order)
		}
		canonical := 0
		if got < 0 {
			canonical = -1
		} else if got > 0 {
			canonical = 1
		}
		for _, v := range []int{len(a), len(b), offA, offB, n, canonical} {
			binary.LittleEndian.PutUint32(encoded[:], uint32(int32(v)))
			h.Write(encoded[:])
		}
		h.Write(a)
		h.Write(b)
		count++
	}
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			check([]byte{byte(x)}, []byte{byte(y)}, 0, 0, 1)
		}
	}
	for _, oa := range []int{0, 1, 3, 7, 15} {
		for _, ob := range []int{0, 1, 3, 7, 15} {
			for n := 0; n <= 128; n++ {
				a, b := make([]byte, 160), make([]byte, 160)
				for i := 0; i < n; i++ {
					a[oa+i] = byte(i*37 + 91)
					b[ob+i] = a[oa+i]
				}
				check(a, b, oa, ob, n)
				for _, pos := range []int{0, n / 2, n - 1} {
					if pos < 0 || pos >= n {
						continue
					}
					old := b[ob+pos]
					b[ob+pos] = old ^ 0xff
					check(a, b, oa, ob, n)
					check(b, a, ob, oa, n)
					b[ob+pos] = old
				}
			}
		}
	}
	gotHash := fmt.Sprintf("%x", h.Sum(nil))
	t.Logf("memcmp cases=%d ordering-sha256=%s", count, gotHash)
	if wantHash != "" && gotHash != wantHash {
		t.Fatalf("memcmp contract hash %s want %s", gotHash, wantHash)
	}
}

func TestStringCompareContracts(t *testing.T) {
	const wantHash = "2fa41584258301b12796e740845f54350ffa51746dd346a2e1d52310c1a6b968" // Freeze from original libc after sign normalization.
	h := sha256.New()
	var encoded [4]byte
	count := 0
	check := func(a, b []byte, oa, ob int) {
		t.Helper()
		beforeA, beforeB := bytes.Clone(a), bytes.Clone(b)
		got := Strcmp(unsafe.Pointer(&a[oa]), unsafe.Pointer(&b[ob]))
		na, nb := bytes.IndexByte(a[oa:], 0), bytes.IndexByte(b[ob:], 0)
		if na < 0 || nb < 0 {
			t.Fatal("invalid contract string")
		}
		order := bytes.Compare(a[oa:oa+na], b[ob:ob+nb])
		if (got < 0) != (order < 0) || (got > 0) != (order > 0) || !bytes.Equal(a, beforeA) || !bytes.Equal(b, beforeB) {
			t.Fatalf("strcmp offsets=%d/%d got=%d order=%d", oa, ob, got, order)
		}
		canonical := 0
		if got < 0 {
			canonical = -1
		} else if got > 0 {
			canonical = 1
		}
		for _, v := range []int{len(a), len(b), oa, ob, canonical} {
			binary.LittleEndian.PutUint32(encoded[:], uint32(int32(v)))
			h.Write(encoded[:])
		}
		h.Write(a)
		h.Write(b)
		count++
	}
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			check([]byte{byte(x), 0}, []byte{byte(y), 0}, 0, 0)
		}
	}
	for _, oa := range []int{0, 1, 3, 7, 15} {
		for _, ob := range []int{0, 1, 3, 7, 15} {
			for n := 0; n <= 128; n++ {
				a, b := bytes.Repeat([]byte{0xa5}, 160), bytes.Repeat([]byte{0x5a}, 160)
				for i := 0; i < n; i++ {
					a[oa+i] = byte(i%255 + 1)
					b[ob+i] = a[oa+i]
				}
				a[oa+n] = 0
				b[ob+n] = 0
				check(a, b, oa, ob)
				for _, pos := range []int{0, n / 2, n - 1} {
					if pos < 0 || pos >= n {
						continue
					}
					old := b[ob+pos]
					b[ob+pos] = old ^ 0x80
					check(a, b, oa, ob)
					check(b, a, ob, oa)
					b[ob+pos] = 0
					check(a, b, oa, ob)
					check(b, a, ob, oa)
					b[ob+pos] = old
				}
			}
		}
	}
	gotHash := fmt.Sprintf("%x", h.Sum(nil))
	t.Logf("strcmp cases=%d ordering-sha256=%s", count, gotHash)
	if wantHash != "" && gotHash != wantHash {
		t.Fatalf("strcmp contract hash %s want %s", gotHash, wantHash)
	}
}
