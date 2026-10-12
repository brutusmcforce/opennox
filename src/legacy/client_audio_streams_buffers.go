package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

// Counts are positive at every production owner; pool blocks retain their 4-byte free-list headers.
func audioStreamPoolNew(count, size int32) *unsafe.Pointer {
	stride := uint32(size) + 4
	p, _ := alloc.Calloc(1, uintptr(uint32(count)*stride+4))
	if p == nil {
		return nil
	}
	head := (*unsafe.Pointer)(p)
	*head = unsafe.Add(p, 4)
	for i := int32(0); i < count; i++ {
		node := (*unsafe.Pointer)(unsafe.Add(p, 4+uintptr(uint32(i)*stride)))
		if i+1 == count {
			*node = nil
		} else {
			*node = unsafe.Add(unsafe.Pointer(node), uintptr(stride))
		}
	}
	return head
}
func audioStreamPoolFree(p unsafe.Pointer) {
	if p != nil {
		alloc.FreePtr(p)
	}
}
func audioStreamPoolPop(p *unsafe.Pointer) unsafe.Pointer {
	item := *p
	if item == nil {
		return nil
	}
	*p = *(*unsafe.Pointer)(item)
	return unsafe.Add(item, 4)
}
func audioStreamPoolPush(p *unsafe.Pointer, item unsafe.Pointer) unsafe.Pointer {
	node := unsafe.Add(item, -4)
	*(*unsafe.Pointer)(node) = *p
	*p = node
	return node
}
func audioStreamBufferInit(p *audioStreamBuffer) {
	p.Data = 0
	p.Length = 0
	p.Format = nil
	p.Field24 = 0
	listClear(&p.Chunks)
}
func audioStreamBufferAppend(p *audioStreamBuffer, c *audioStreamChunk) uint32 {
	listAppend(&p.Chunks, &c.Node)
	p.Length += c.Length
	c.Owner = p
	return p.Length
}
func audioStreamBufferFirst(p *audioStreamBuffer) *audioStreamChunk {
	return (*audioStreamChunk)(unsafe.Pointer(listNext(&p.Chunks)))
}
func audioStreamChunkInit(p *audioStreamChunk, data uint32, length uint32) *audioStreamChunk {
	p.Length = length
	p.Data = data
	listInit(&p.Node)
	p.Owner = nil
	return p
}
func audioStreamChunkClear(p *audioStreamChunk) *audioStreamChunk { p.Owner = nil; return p }
func audioStreamByteRate(p *audioStreamFormat) uint32 {
	result := int32(p.Rate * p.Channels * p.Width)
	if p.Encoding == 1 {
		result >>= 2
	}
	p.ByteRate = uint32(result)
	return uint32(result)
}
func audioStreamCatalogEntry(p *audioStreamCatalog, index int32) *audioStreamEntry {
	return (*audioStreamEntry)(unsafe.Add(unsafe.Pointer(p.Entries), uintptr(uint32(index)*36)))
}
func audioStreamReadFormat(p *audioStreamCatalog, index int32, out *audioStreamFormat) uint32 {
	e := audioStreamCatalogEntry(p, index)
	out.Kind = 4
	out.Rate = e.Rate
	out.Channels = 1 + (e.Flags & 1)
	out.Extra = e.Extra
	if e.Flags&8 != 0 {
		out.Encoding = 2
		out.Width = 2
	} else {
		out.Encoding = 0
		out.Width = 1 + ((e.Flags >> 2) & 1)
	}
	return out.Width
}
func audioStreamScaleVolume(p unsafe.Pointer, percent uint32) uint32 {
	return percent * (*(*uint32)(unsafe.Add(p, 36)) >> 16) / 100
}
func audioStreamCacheNew(cat *audioStreamCatalog, budget, entries, payload int32) *audioStreamCache {
	p, _ := alloc.New(audioStreamCache{})
	p.Catalog = cat
	p.Payload = payload
	p.Blocks = audioStreamPoolNew(budget/(payload+24), payload+24)
	p.Entries = audioStreamPoolNew(entries, 84)
	listClear(&p.List)
	if p.Blocks != nil && p.Entries != nil {
		return p
	}
	audioStreamCacheFree(p)
	return nil
}
func audioStreamCacheFree(p *audioStreamCache) {
	for it := listNext(&p.List); it != nil; it = listNext(&p.List) {
		audioStreamCacheDrop((*audioStreamCacheEntry)(unsafe.Pointer(it)))
	}
	audioStreamPoolFree(unsafe.Pointer(p.Blocks))
	audioStreamPoolFree(unsafe.Pointer(p.Entries))
	alloc.Free(p)
}
func audioStreamCacheFind(p *audioStreamCache, index int32) *audioStreamCacheEntry {
	for n := p.List.next; n != &p.List; n = n.next {
		e := (*audioStreamCacheEntry)(unsafe.Pointer(n))
		if e.Index == index && e.Valid != 0 {
			return e
		}
	}
	return nil
}
func audioStreamCacheEvict(p *audioStreamCache) int32 {
	for n := listPrev(&p.List); n != nil; n = listPrev(n) {
		e := (*audioStreamCacheEntry)(unsafe.Pointer(n))
		if e.References == 0 {
			audioStreamCacheDrop(e)
			return 1
		}
	}
	return 0
}
func audioStreamCacheRef(p *audioStreamCacheEntry) *audioStreamCacheEntry { p.References++; return p }
func audioStreamCacheUnref(p *audioStreamCacheEntry) int32 {
	result := int32(p.References - 1)
	p.References = uint32(result)
	if result < 0 {
		p.References = 0
	}
	return result
}
func audioStreamCacheReferences(p *audioStreamCacheEntry) uint32 { return p.References }
func audioStreamCacheDrop(p *audioStreamCacheEntry) unsafe.Pointer {
	if p.Node.prev != &p.Node {
		listRemove(&p.Node)
	}
	for c := audioStreamBufferFirst(&p.Buffer); c != nil; c = audioStreamBufferFirst(&p.Buffer) {
		listRemove(&c.Node)
		audioStreamChunkClear(c)
		audioStreamPoolPush(p.Cache.Blocks, unsafe.Pointer(c))
	}
	return audioStreamPoolPush(p.Cache.Entries, unsafe.Pointer(p))
}
func audioStreamCacheBuffer(p *audioStreamCacheEntry) *audioStreamBuffer { return &p.Buffer }
func audioStreamCacheLoad(p *audioStreamCache, index int32) *audioStreamCacheEntry {
	if e := audioStreamCacheFind(p, index); e != nil {
		listRemove(&e.Node)
		listPrepend(&p.List, &e.Node)
		return e
	}
	if audioStreamOpen(p.Catalog, index) == 0 {
		return nil
	}
	e := (*audioStreamCacheEntry)(audioStreamPoolPop(p.Entries))
	if e == nil {
		audioStreamCacheEvict(p)
		e = (*audioStreamCacheEntry)(audioStreamPoolPop(p.Entries))
		if e == nil {
			audioStreamClose(p.Catalog)
			return nil
		}
	}
	e.Index = index
	e.Cache = p
	listInit(&e.Node)
	e.References = 0
	audioStreamBufferInit(&e.Buffer)
	e.Buffer.Format = &e.Format
	remaining := p.Catalog.Remaining
	for remaining != 0 {
		n := p.Payload
		if n > int32(remaining) {
			n = int32(remaining)
		}
		c := (*audioStreamChunk)(audioStreamPoolPop(p.Blocks))
		if c == nil {
			for audioStreamCacheEvict(p) != 0 {
				c = (*audioStreamChunk)(audioStreamPoolPop(p.Blocks))
				if c != nil {
					break
				}
			}
		}
		if c == nil {
			audioStreamCacheDrop(e)
			return nil
		}
		audioStreamChunkInit(c, memmap.PtrWord(unsafe.Add(unsafe.Pointer(c), 24)), uint32(n))
		audioStreamBufferAppend(&e.Buffer, c)
		got := audioStreamRead(p.Catalog, memmap.WordPtr(c.Data), n)
		if got != n {
			audioStreamCacheDrop(e)
			return nil
		}
		remaining -= uint32(got)
	}
	audioStreamReadFormat(p.Catalog, e.Index, &e.Format)
	listPrepend(&p.List, &e.Node)
	e.Valid = 1
	audioStreamClose(p.Catalog)
	return e
}
