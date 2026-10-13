package memmap

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Word regions.
//
// Some legacy code (the map generator) keeps heap blocks as pointer words and does address arithmetic on them,
// e.g. row := base + 76*i. A plain pointer handle does not allow that: base+76 is a different handle. Blocks
// registered with RegisterWordRegion get a contiguous range of pseudo addresses instead, so a word pointing
// inside the block is base+offset, as it was on 386. Pseudo addresses lie in [regionStart, handleTag), above
// every blob of the original binary.

const (
	regionStart = 0x80000000
	regionAlign = 16
)

type wordRegion struct {
	ptr  uintptr // real address
	addr uintptr // pseudo address
	size uintptr
}

type regionsT struct {
	sync.Mutex
	count  atomic.Int32
	byPtr  []*wordRegion // sorted by ptr
	byAddr []*wordRegion // sorted by addr
}

var regions regionsT

// RegisterWordRegion makes pointer words into the block [p, p+size] (including its end) support address
// arithmetic.
// The block must be unregistered with UnregisterWordRegion before it is freed.
func RegisterWordRegion(p unsafe.Pointer, size uintptr) {
	if p == nil || size == 0 {
		return
	}
	r := &regions
	r.Lock()
	defer r.Unlock()
	ptr := uintptr(p)
	// A block that was freed without unregistering may overlap: drop it.
	for i := 0; i < len(r.byPtr); {
		o := r.byPtr[i]
		if o.ptr < ptr+size && ptr < o.ptr+o.size { // real overlap
			r.remove(o)
			continue
		}
		i++
	}
	// First fit between the pseudo ranges in use.
	addr := uintptr(regionStart)
	for _, o := range r.byAddr {
		if addr+size < o.addr {
			break
		}
		addr = (o.addr + o.size + regionAlign) &^ (regionAlign - 1)
	}
	if addr+size > handleTag {
		panic("memmap: word regions exhausted")
	}
	reg := &wordRegion{ptr: ptr, addr: addr, size: size}
	i := sort.Search(len(r.byPtr), func(i int) bool { return r.byPtr[i].ptr > ptr })
	r.byPtr = append(r.byPtr, nil)
	copy(r.byPtr[i+1:], r.byPtr[i:])
	r.byPtr[i] = reg
	j := sort.Search(len(r.byAddr), func(i int) bool { return r.byAddr[i].addr > addr })
	r.byAddr = append(r.byAddr, nil)
	copy(r.byAddr[j+1:], r.byAddr[j:])
	r.byAddr[j] = reg
	r.count.Add(1)
}

// UnregisterWordRegion removes the block starting at p (see RegisterWordRegion).
func UnregisterWordRegion(p unsafe.Pointer) {
	if p == nil || regions.count.Load() == 0 {
		return
	}
	r := &regions
	r.Lock()
	defer r.Unlock()
	if o := r.findPtr(uintptr(p)); o != nil && o.ptr == uintptr(p) {
		r.remove(o)
	}
}

func (r *regionsT) remove(o *wordRegion) {
	for i, v := range r.byPtr {
		if v == o {
			r.byPtr = append(r.byPtr[:i], r.byPtr[i+1:]...)
			break
		}
	}
	for i, v := range r.byAddr {
		if v == o {
			r.byAddr = append(r.byAddr[:i], r.byAddr[i+1:]...)
			break
		}
	}
	r.count.Add(-1)
}

// findPtr returns the region that contains the real address ptr.
func (r *regionsT) findPtr(ptr uintptr) *wordRegion {
	i := sort.Search(len(r.byPtr), func(i int) bool { return r.byPtr[i].ptr > ptr }) - 1
	if i >= 0 && ptr <= r.byPtr[i].ptr+r.byPtr[i].size {
		return r.byPtr[i]
	}
	return nil
}

// regionWord returns the pseudo address of p if it is inside a word region.
func regionWord(p uintptr) (uint32, bool) {
	if regions.count.Load() == 0 {
		return 0, false
	}
	r := &regions
	r.Lock()
	defer r.Unlock()
	if o := r.findPtr(p); o != nil {
		return uint32(o.addr + (p - o.ptr)), true
	}
	return 0, false
}

// regionPtr decodes a pseudo address of a word region.
func regionPtr(w uint32) unsafe.Pointer {
	r := &regions
	r.Lock()
	defer r.Unlock()
	addr := uintptr(w)
	i := sort.Search(len(r.byAddr), func(i int) bool { return r.byAddr[i].addr > addr }) - 1
	if i >= 0 && addr <= r.byAddr[i].addr+r.byAddr[i].size {
		o := r.byAddr[i]
		return unsafe.Pointer(o.ptr + (addr - o.addr))
	}
	panic(fmt.Errorf("memmap: pointer word 0x%X is not in a word region", w))
}
