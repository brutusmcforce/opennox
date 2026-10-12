package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/protection"
)

func insertProtectionRecord(r *protection.Record) int32 {
	count := memmap.PtrUint16(0x587000, 311204)
	if *count == 0 {
		(*count)++
		dword_5d4594_2516352 = uint32(uintptr(unsafe.Pointer(r)))
		dword_5d4594_2516344 = uint32(uintptr(unsafe.Pointer(r)))
		return 1
	}

	target := GetServer().S().Rand.Logic.IntClamp(0, int(*count)-1)
	before := protectionHead()
	for i := 0; i < target; i++ {
		before = before.Next
	}
	head := protection.InsertBefore(protectionHead(), before, r)
	dword_5d4594_2516344 = uint32(uintptr(unsafe.Pointer(head)))
	(*count)++
	return 1
}
