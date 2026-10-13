package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

// clientSequenceNode is a queued ordered message (see clientSequenceEnqueue); the message bytes follow it.
type clientSequenceNode struct {
	legacyListNode // tag: the sequence number
	Ticks          uint64
	Size           uint16
}

func clientSequenceData(p *legacyListNode) []byte {
	n := (*clientSequenceNode)(unsafe.Pointer(p))
	return unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(n), unsafe.Sizeof(clientSequenceNode{}))), n.Size)
}

func clientSequenceHead() *legacyListNode {
	return &clientSequenceList
}
func clientSequenceOld(current, sequence uint16) bool {
	if current < 10000 && sequence >= 65535-(10000-current) {
		return true
	}
	distance := int(current) - int(sequence)
	return distance > 0 && distance <= 10000
}
func clientSequenceContains(sequence uint16) bool {
	for p := listNext(clientSequenceHead()); p != nil; p = listNext(p) {
		if p.tag == uintptr(sequence) {
			return true
		}
	}
	return false
}
func clientSequenceEnqueue(data []byte) {
	seq := binary.LittleEndian.Uint16(data[1:])
	if clientSequenceOld(memmap.Uint16(0x5D4594, 1197360), seq) || clientSequenceContains(seq) {
		return
	}
	size := int(data[3])
	// Retain the foreign queue allocation/layout while shared globals remain ABI words.
	ptr := legacyCalloc(unsafe.Sizeof(clientSequenceNode{})+uintptr(size), 1)
	if ptr == nil {
		return
	}
	node := listInit((*legacyListNode)(ptr))
	node.tag = uintptr(seq)
	n := (*clientSequenceNode)(ptr)
	n.Size = uint16(size)
	n.Ticks = uint64(uint32(PlatformTicks()))
	copy(clientSequenceData(node), data[4:4+size])
	if memmap.Uint16(0x5D4594, 1197360) == seq {
		dword_5d4594_1197352 = memmap.PtrWord(ptr)
	}
	listAscending(clientSequenceHead(), node)
}
func clientSequencePoll() {
	ready := (*legacyListNode)(memmap.WordPtr(dword_5d4594_1197352))
	if ready == nil {
		pending := (*legacyListNode)(memmap.WordPtr(dword_5d4594_1197356))
		if pending != nil && uint64(uint32(PlatformTicks()))-(*clientSequenceNode)(unsafe.Pointer(pending)).Ticks > 30000 {
			*memmap.PtrUint16(0x5D4594, 1197360) = uint16(pending.tag)
			ready = pending
			dword_5d4594_1197352 = memmap.PtrWord(unsafe.Pointer(pending))
			dword_5d4594_1197356 = memmap.PtrWord(unsafe.Pointer(listNext(pending)))
		}
	}
	p := ready
	for p != nil && uint32(p.tag) == uint32(memmap.Uint16(0x5D4594, 1197360)) {
		next := listNext(p)
		if next == nil {
			next = listNext(clientSequenceHead())
			if next == p {
				next = nil
			}
		}
		data := clientSequenceData(p)
		Nox_xxx_netOnPacketRecvCli_48EA70(31, unsafe.SliceData(data), len(data))
		*memmap.PtrUint16(0x5D4594, 1197360)++
		listRemove(p)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	// Original polling drops the pending cursor when a gap is polled before its
	// timeout. Preserve this observable behavior separately from queue membership.
	dword_5d4594_1197356 = memmap.PtrWord(unsafe.Pointer(p))
	dword_5d4594_1197352 = 0
}
func clientSequenceInit() {
	listClear(clientSequenceHead())
	dword_5d4594_1197352 = 0
	dword_5d4594_1197356 = 0
	*memmap.PtrUint16(0x5D4594, 1197360) = 0
}
func clientSequenceFree() {
	for p := listNext(clientSequenceHead()); p != nil; {
		next := listNext(p)
		listRemove(p)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	listClear(clientSequenceHead())
	*memmap.PtrUint16(0x5D4594, 1197360) = 0
}
