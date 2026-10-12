package memmap

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// Pointer words.
//
// The blobs keep the 32-bit layout of the original binary, so a slot in them has room for 4 bytes only.
// Pointers stored in blobs are therefore encoded as 32-bit pointer words, which works on every architecture:
//
//   - 0 is nil;
//   - a pointer into a blob is stored as its address in the original binary (below handleTag);
//   - any other pointer is stored as handleTag | index into a registry of such pointers.
//
// Code that reads a pointer word as a raw pointer gets an address that is not mapped (on Linux),
// so mistakes crash instead of silently reading wrong memory.
//
// Pointer words are also used for pointers that legacy code keeps in 32-bit integers (variables and
// struct fields). The registry keeps plain addresses: like the integers it replaces, a word does not keep
// its object alive, and the same address always gets the same word.

const handleTag = 0xC0000000

// encodeWords selects encoded pointer words. On 386 a pointer fits in 32 bits, so words are the plain
// pointer values unless NOX_PTRWORDS is set; tests set it to find code that mixes words and raw pointers,
// which only works on 386 and breaks on 64-bit platforms.
var encodeWords = runtime.GOARCH != "386" || os.Getenv("NOX_PTRWORDS") != ""

// EncodedWords reports if pointer words are encoded (see PtrWord).
func EncodedWords() bool { return encodeWords }

var handles struct {
	sync.Mutex
	list  []uintptr
	index map[uintptr]uint32
}

// PtrWord encodes a pointer as a 32-bit pointer word.
func PtrWord(p unsafe.Pointer) uint32 {
	if p == nil {
		return 0
	}
	if !encodeWords {
		return uint32(uintptr(p))
	}
	if b, off := BlobByPtr(p); b != nil {
		return uint32(b.Addr + off)
	}
	h := &handles
	h.Lock()
	defer h.Unlock()
	if w, ok := h.index[uintptr(p)]; ok {
		return w
	}
	if h.index == nil {
		h.index = make(map[uintptr]uint32)
	}
	if len(h.list) >= handleTag>>1 {
		panic("memmap: too many pointer handles")
	}
	w := handleTag | uint32(len(h.list))
	h.list = append(h.list, uintptr(p))
	h.index[uintptr(p)] = w
	return w
}

// WordPtr decodes a 32-bit pointer word.
func WordPtr(w uint32) unsafe.Pointer {
	if w == 0 {
		return nil
	}
	if !encodeWords {
		return unsafe.Pointer(uintptr(w))
	}
	if w&handleTag == handleTag {
		h := &handles
		h.Lock()
		defer h.Unlock()
		i := w &^ handleTag
		if int(i) >= len(h.list) {
			panic(fmt.Errorf("memmap: invalid pointer handle 0x%X", w))
		}
		return unsafe.Pointer(h.list[i])
	}
	return Ptr(uintptr(w))
}

// LoadPtr reads a pointer word from a blob.
func LoadPtr(base, off uintptr) unsafe.Pointer {
	return WordPtr(Uint32(base, off))
}

// StorePtr writes a pointer word to a blob.
func StorePtr(base, off uintptr, p unsafe.Pointer) {
	*PtrUint32(base, off) = PtrWord(p)
}
