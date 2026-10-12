package alloc

import (
	"math/bits"
	"os"
	"sync"
	"unsafe"
)

// Retain up to 1 MiB per class, or two spans for the largest class (18 MiB total).
const rawSpanSize = 64 * 1024
const rawClasses = 17 // Powers of two, 16 bytes through 1 MiB.
const rawDirectClass = rawClasses
const rawUnused = ^uintptr(0)

var rawPageSize = uintptr(os.Getpagesize())

type rawSpan struct {
	data              []byte
	pointer           unsafe.Pointer // Exact returned address for a direct mapping.
	stride            uintptr
	class             int
	affinity          *rawAffinity
	free              []uint32
	sizes             []uintptr // Requested sizes; rawUnused marks a released slot.
	live, initialized int
	prev, next        *rawSpan
	linked, mapped    bool
}

// Pool tokens are reuse hints only. GC or scheduling may discard affinity;
// allocation ownership and correctness never depend on token identity.
type rawAffinity struct{ marker byte }

var rawAffinities = sync.Pool{New: func() any { return new(rawAffinity) }}

func rawAffinityHint() *rawAffinity {
	hint := rawAffinities.Get().(*rawAffinity)
	rawAffinities.Put(hint)
	return hint
}

type rawClass struct {
	sync.Mutex
	partial *rawSpan
	empty   [16]*rawSpan
	evict   uint32
	_       [64]byte // Keep neighboring class locks off the same cache line.
}

// Payloads are OS memory, never managed Go storage. Registries root span metadata
// and locate its size-class lock; allocation state is stored in integer arrays.
// Registry entries change only when spans are mapped/unmapped, not for every block.
var rawHeap struct {
	classes [rawClasses + 1]rawClass
	pages   sync.Map // Page address -> multi-block span.
	exact   sync.Map // Returned address -> single-block span or direct mapping.
}

func rawLink(c *rawClass, s *rawSpan) {
	if s.linked {
		panic("raw span already linked")
	}
	s.prev = nil
	s.next = c.partial
	if c.partial != nil {
		c.partial.prev = s
	}
	c.partial = s
	s.linked = true
}
func rawUnlink(c *rawClass, s *rawSpan) {
	if !s.linked {
		panic("raw span not linked")
	}
	if s.prev != nil {
		s.prev.next = s.next
	} else {
		c.partial = s.next
	}
	if s.next != nil {
		s.next.prev = s.prev
	}
	s.prev = nil
	s.next = nil
	s.linked = false
}
func rawMapping(n uintptr) []byte {
	const max = ^uintptr(0) >> 1
	if n == 0 || n > max-(rawPageSize-1) {
		return nil
	}
	n = (n + rawPageSize - 1) &^ (rawPageSize - 1)
	return rawMap(int(n))
}
func rawRegister(s *rawSpan) {
	if s.stride == 0 || len(s.sizes) == 1 {
		p := s.pointer
		if p == nil {
			p = unsafe.Pointer(&s.data[0])
		}
		rawHeap.exact.Store(uintptr(p), s)
	} else {
		base := uintptr(unsafe.Pointer(&s.data[0]))
		for off := uintptr(0); off < uintptr(len(s.data)); off += rawPageSize {
			rawHeap.pages.Store(base+off, s)
		}
	}
	s.mapped = true
}
func rawUnregister(s *rawSpan) {
	if s.stride == 0 || len(s.sizes) == 1 {
		p := s.pointer
		if p == nil {
			p = unsafe.Pointer(&s.data[0])
		}
		rawHeap.exact.Delete(uintptr(p))
	} else {
		base := uintptr(unsafe.Pointer(&s.data[0]))
		for off := uintptr(0); off < uintptr(len(s.data)); off += rawPageSize {
			rawHeap.pages.Delete(base + off)
		}
	}
	s.mapped = false
	if err := rawUnmap(s.data); err != nil {
		panic(err)
	}
}
func rawFind(ptr unsafe.Pointer) *rawSpan {
	address := uintptr(ptr)
	if s, ok := rawHeap.pages.Load(address &^ (rawPageSize - 1)); ok {
		return s.(*rawSpan)
	}
	if s, ok := rawHeap.exact.Load(address); ok {
		return s.(*rawSpan)
	}
	panic("invalid raw allocation pointer")
}

// Call while holding this span's class lock. A valid caller-owned block keeps
// its span mapped between registry lookup and lock acquisition.
func rawSlot(s *rawSpan, ptr unsafe.Pointer) int {
	if !s.mapped {
		panic("released raw span")
	}
	if s.stride == 0 {
		if ptr != s.pointer || s.live != 1 {
			panic("invalid raw allocation pointer")
		}
		return 0
	}
	off := uintptr(ptr) - uintptr(unsafe.Pointer(&s.data[0]))
	if off >= uintptr(len(s.data)) || off&(s.stride-1) != 0 {
		panic("invalid raw allocation pointer")
	}
	slot := int(off >> uint(s.class+4))
	if s.sizes[slot] == rawUnused {
		panic("released raw allocation")
	}
	return slot
}
func rawAllocateState(size, alignment uintptr) (unsafe.Pointer, bool) {
	const max = ^uintptr(0) >> 1
	if size > max || alignment == 0 || alignment&(alignment-1) != 0 {
		return nil, false
	}
	need := size
	if need < alignment {
		need = alignment
	}
	if need < 16 {
		need = 16
	}
	class := rawDirectClass
	if need <= 1<<20 && alignment <= rawPageSize {
		class = bits.Len(uint(need-1)) - 4
	}
	var affinity *rawAffinity
	if class >= 12 && class < rawDirectClass {
		affinity = rawAffinityHint()
	}
	c := &rawHeap.classes[class]
	c.Lock()
	defer c.Unlock()
	if class != rawDirectClass {
		stride := uintptr(1) << uint(class+4)
		s := c.partial
		if affinity != nil {
			// Prefer recent worker-local reuse of large buffers.
			s = nil
			for candidate := c.partial; candidate != nil; candidate = candidate.next {
				if candidate.affinity == affinity {
					s = candidate
					break
				}
			}
		}
		if s == nil {
			spanBytes := uintptr(rawSpanSize)
			if stride > spanBytes {
				spanBytes = stride
			}
			data := rawMapping(spanBytes)
			if data == nil {
				return nil, false
			}
			count := int(spanBytes / stride)
			s = &rawSpan{data: data, stride: stride, class: class, free: make([]uint32, count), sizes: make([]uintptr, count)}
			for i := range s.free {
				s.free[i] = uint32(count - 1 - i)
				s.sizes[i] = rawUnused
			}
			rawRegister(s)
			rawLink(c, s)
		}
		if s.live == 0 {
			for i, idle := range c.empty {
				if idle == s {
					c.empty[i] = nil
					break
				}
			}
		}
		s.affinity = affinity
		last := len(s.free) - 1
		slot := s.free[last]
		s.free = s.free[:last]
		s.sizes[slot] = size
		s.live++
		// Previously unused slots in an anonymous mapping are already zero. Every
		// allocation, including malloc/realloc, consumes that one-time guarantee.
		fresh := int(slot) == s.initialized
		if fresh {
			s.initialized++
		}
		if len(s.free) == 0 {
			rawUnlink(c, s)
		}
		return unsafe.Add(unsafe.Pointer(&s.data[0]), uintptr(slot)*stride), fresh
	}
	extra := uintptr(0)
	if alignment > rawPageSize {
		extra = alignment - 1
	}
	if size > max-extra {
		return nil, false
	}
	n := size + extra
	if n == 0 {
		n = 1
	}
	data := rawMapping(n)
	if data == nil {
		return nil, false
	}
	base := uintptr(unsafe.Pointer(&data[0]))
	offset := (-base) & (alignment - 1)
	p := unsafe.Add(unsafe.Pointer(&data[0]), offset)
	s := &rawSpan{data: data, pointer: p, class: class, sizes: []uintptr{size}, live: 1}
	rawRegister(s)
	return p, true
}
func rawAllocate(size, alignment uintptr) unsafe.Pointer {
	p, _ := rawAllocateState(size, alignment)
	return p
}
func rawCalloc(num, size uintptr) unsafe.Pointer {
	if num != 0 && size > (^uintptr(0)>>1)/num {
		return nil
	}
	n := num * size
	p, fresh := rawAllocateState(n, 16)
	if p != nil && !fresh {
		rawClear(p, n)
	}
	return p
}
func rawFree(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	s := rawFind(ptr)
	c := &rawHeap.classes[s.class]
	c.Lock()
	defer c.Unlock()
	slot := rawSlot(s, ptr)
	if s.stride == 0 {
		s.live = 0
		s.sizes[0] = rawUnused
		rawUnregister(s)
		return
	}
	wasFull := len(s.free) == 0
	s.sizes[slot] = rawUnused
	s.free = append(s.free, uint32(slot))
	s.live--
	if wasFull {
		rawLink(c, s)
	}
	if s.live == 0 {
		limit := (1 << 20) / len(s.data)
		if limit < 2 {
			limit = 2
		}
		for i, idle := range c.empty[:limit] {
			if idle == nil {
				c.empty[i] = s
				return
			}
		}
		// Rotate eviction so expired affinity hints cannot monopolize the cache.
		index := c.evict % uint32(limit)
		c.evict++
		victim := c.empty[index]
		c.empty[index] = s
		rawUnlink(c, victim)
		rawUnregister(victim)
	}

}

// rawResize keeps a block in place where appropriate and snapshots its logical
// size for a move. Never hold two class locks, or a lock during the byte copy.
func rawResize(ptr unsafe.Pointer, size uintptr) (oldSize, capacity uintptr, retained bool) {
	s := rawFind(ptr)
	c := &rawHeap.classes[s.class]
	c.Lock()
	defer c.Unlock()
	slot := rawSlot(s, ptr)
	capacity = s.stride
	if capacity == 0 {
		capacity = uintptr(len(s.data)) - (uintptr(ptr) - uintptr(unsafe.Pointer(&s.data[0])))
	}
	oldSize = s.sizes[slot]
	if size <= capacity && (capacity <= 32768 || size > capacity/4) {
		s.sizes[slot] = size
		return oldSize, capacity, true
	}
	return oldSize, capacity, false
}
func RawRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	if ptr == nil {
		return rawAllocate(size, 16)
	}
	if size == 0 {
		rawFree(ptr)
		return nil
	}
	oldSize, capacity, retained := rawResize(ptr, size)
	if retained {
		return ptr
	}
	next := rawAllocate(size, 16)
	if next == nil {
		// A shrinking move may fail without preventing the existing storage from
		// satisfying the request. Preserve the original block on every failure.
		if size <= capacity {
			s := rawFind(ptr)
			c := &rawHeap.classes[s.class]
			c.Lock()
			slot := rawSlot(s, ptr)
			s.sizes[slot] = size
			c.Unlock()
			return ptr
		}
		return nil
	}
	n := oldSize
	if n > size {
		n = size
	}
	copy(unsafe.Slice((*byte)(next), n), unsafe.Slice((*byte)(ptr), n))
	rawFree(ptr)
	return next
}

// RawMalloc retains the original non-recoverable allocation failure disposition.
func RawMalloc(size uintptr) unsafe.Pointer {
	p := rawAllocate(size, 16)
	if p == nil {
		os.Stderr.WriteString("runtime: C malloc failed\n")
		os.Exit(2)
	}
	return p
}
