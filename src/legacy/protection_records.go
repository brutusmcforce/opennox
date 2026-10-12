package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/protection"
)

func protectionHead() *protection.Record {
	return dword_5d4594_2516344
}

func swapProtectionRecords(a, b *protection.Record) {
	if protection.Swap(a, b) {
		*memmap.PtrUint32(0x5D4594, 2516360)++
	}
}
