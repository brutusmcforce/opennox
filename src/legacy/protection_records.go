package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/protection"
)

// The C manager stores payloads at words 0/1 and links at words 2/3.
var (
	_ [16 - unsafe.Sizeof(protection.Record{})]byte
	_ [unsafe.Sizeof(protection.Record{}) - 16]byte
	_ [4 - unsafe.Alignof(protection.Record{})]byte
	_ [unsafe.Alignof(protection.Record{}) - 4]byte
	_ [0 - unsafe.Offsetof(protection.Record{}.ID)]byte
	_ [4 - unsafe.Offsetof(protection.Record{}.Value)]byte
	_ [unsafe.Offsetof(protection.Record{}.Value) - 4]byte
	_ [8 - unsafe.Offsetof(protection.Record{}.Next)]byte
	_ [unsafe.Offsetof(protection.Record{}.Next) - 8]byte
	_ [12 - unsafe.Offsetof(protection.Record{}.Prev)]byte
	_ [unsafe.Offsetof(protection.Record{}.Prev) - 12]byte
)

func protectionHead() *protection.Record {
	return *(**protection.Record)(unsafe.Pointer(&dword_5d4594_2516344))
}

func swapProtectionRecords(a, b *protection.Record) {
	if protection.Swap(a, b) {
		*memmap.PtrUint32(0x5D4594, 2516360)++
	}
}
