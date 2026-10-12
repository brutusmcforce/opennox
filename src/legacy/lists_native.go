package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// The third word is a sentinel's own address or an item's signed sort key.
// Shared nodes use the existing address-stable storage and retain its layout.
type legacyListNode struct {
	next, prev *legacyListNode
	tag        uintptr
}

var _ = [1]struct{}{}[12-unsafe.Sizeof(legacyListNode{})]
var _ = [1]struct{}{}[8-unsafe.Offsetof(legacyListNode{}.tag)]

func listClear(p *legacyListNode) {
	p.next = p
	p.prev = p
	p.tag = uintptr(unsafe.Pointer(p))
}
func listInit(p *legacyListNode) *legacyListNode {
	p.next = p
	p.prev = p
	p.tag = 0
	return p
}
func listNext(p *legacyListNode) *legacyListNode {
	if p == nil {
		return nil
	}
	next := p.next
	if next != nil && next.tag == uintptr(unsafe.Pointer(next)) {
		return nil
	}
	return next
}
func listPrev(p *legacyListNode) *legacyListNode {
	prev := p.prev
	if prev.tag == uintptr(unsafe.Pointer(prev)) {
		return nil
	}
	return prev
}
func listAppend(before, p *legacyListNode) {
	if before == nil || p == nil {
		panic("nil list or item")
	}
	prev := before.prev
	p.next = before
	p.prev = prev
	before.prev = p
	// Preserve existing partially initialized lists instead of inventing a head.
	if prev != nil {
		prev.next = p
	}
}
func listPrepend(after, p *legacyListNode) *legacyListNode {
	p.prev = after
	p.next = after.next
	after.next = p
	p.next.prev = p
	return p
}
func listRemove(p *legacyListNode) {
	p.prev.next = p.next
	p.next.prev = p.prev
	p.next = p
	p.prev = p
}
func listAscending(head, p *legacyListNode) int {
	index := 0
	for it := listNext(head); it != nil; it = listNext(it) {
		if int32(p.tag) <= int32(it.tag) {
			listAppend(it, p)
			return index
		}
		index++
	}
	listAppend(head, p)
	return index
}
func listDescending(head, p *legacyListNode) {
	for it := listNext(head); it != nil; it = listNext(it) {
		if int32(it.tag) <= int32(p.tag) {
			listAppend(it, p)
			return
		}
	}
	listAppend(head, p)
}
func listAt(head *legacyListNode, index int) *legacyListNode {
	// This helper uses the literal head link, not the third-word sentinel marker.
	for p := head.next; p != head; p = p.next {
		if index == 0 {
			return p
		}
		index--
	}
	return nil
}

func ListFirst(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(listNext((*legacyListNode)(p)))
}
func ListNext(p unsafe.Pointer) unsafe.Pointer { return unsafe.Pointer(listNext((*legacyListNode)(p))) }
func ListClear(p unsafe.Pointer)               { listClear((*legacyListNode)(p)) }
func ListAppend(before, p unsafe.Pointer) {
	listAppend((*legacyListNode)(before), (*legacyListNode)(p))
}

type playerGroup struct {
	list    legacyListNode
	name    [10]uint16
	id      uint32
	team    unsafe.Pointer
	members legacyListNode
}
type playerGroupMember struct {
	list  legacyListNode
	index int32
}

var _ = [1]struct{}{}[52-unsafe.Sizeof(playerGroup{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(playerGroupMember{})]
var _ = [1]struct{}{}[12-unsafe.Offsetof(playerGroup{}.name)]
var _ = [1]struct{}{}[32-unsafe.Offsetof(playerGroup{}.id)]
var _ = [1]struct{}{}[40-unsafe.Offsetof(playerGroup{}.members)]
var _ = [1]struct{}{}[12-unsafe.Offsetof(playerGroupMember{}.index)]

func playerGroupsHead() *legacyListNode { return (*legacyListNode)(memmap.PtrOff(0x5D4594, 599460)) }
func playerGroupsInit() {
	guard := memmap.PtrUint32(0x5D4594, 599472)
	if *guard == 0 {
		listClear(playerGroupsHead())
		*guard = 1
	}
}
func playerGroupsFirst() *playerGroup {
	return (*playerGroup)(unsafe.Pointer(listNext(playerGroupsHead())))
}
func playerGroupsNext(p *playerGroup) *playerGroup {
	return (*playerGroup)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}
func playerGroupFind(id uint32) *playerGroup {
	for p := playerGroupsFirst(); p != nil; p = playerGroupsNext(p) {
		if p.id == id {
			return p
		}
	}
	return nil
}
func playerGroupAdd(id uint32, name *uint16) *playerGroup {
	if playerGroupFind(id) != nil {
		return nil
	}
	p := (*playerGroup)(legacyCalloc(1, 52))
	if p == nil {
		panic("player group allocation")
	}
	p.id = id
	alloc.StrCopyZero16P(p.name[:], name)
	listInit(&p.list)
	listClear(&p.members)
	listAppend(playerGroupsHead(), &p.list)
	return p
}
func playerGroupAddMember(p *playerGroup, index int32) {
	member := (*playerGroupMember)(legacyCalloc(1, 16))
	if member == nil {
		panic("player group member allocation")
	}
	member.index = index
	listInit(&member.list)
	listAppend(&p.members, &member.list)
}
func playerGroupRemoveMember(p *playerGroup, index int32) unsafe.Pointer {
	for node := listNext(&p.members); node != nil; node = listNext(node) {
		member := (*playerGroupMember)(unsafe.Pointer(node))
		if member.index != index {
			continue
		}
		listRemove(node)
		legacyFree(unsafe.Pointer(member))
		break
	}
	result := unsafe.Pointer(&p.members)
	if p.members.prev == &p.members {
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
	}
	return result
}
func playerGroupsFree() {
	for p := playerGroupsFirst(); p != nil; {
		next := playerGroupsNext(p)
		for member := listNext(&p.members); member != nil; {
			after := listNext(member)
			listRemove(member)
			legacyFree(unsafe.Pointer(member))
			member = after
		}
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
}
