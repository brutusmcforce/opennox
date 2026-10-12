package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/protection"
)

func deleteProtectionRecord(id uint32) bool {
	head := protectionHead()
	r := protection.Find(head, uint32(dword_5d4594_2516348), id)
	if r == nil {
		return false
	}
	tail := dword_5d4594_2516352
	head, tail = protection.Unlink(head, tail, r)
	dword_5d4594_2516344 = head
	dword_5d4594_2516352 = tail
	dword_5d4594_2516328 ^= uint32(r.ID ^ r.Value)
	*memmap.PtrUint16(0x587000, 311204)--
	legacyFree(unsafe.Pointer(r))
	return true
}

func freeProtectionRecords() {
	for p := protectionHead(); p != nil; {
		next := p.Next
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	dword_5d4594_2516328 = 0
	*memmap.PtrUint16(0x587000, 311204) = 0
	dword_5d4594_2516348 = 0
	dword_5d4594_2516352 = nil
	dword_5d4594_2516344 = nil
}
