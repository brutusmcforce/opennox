package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/protection"
)

func nox_xxx_protectData_56F5C0() int32 {
	frame := GetServer().S().Frame()
	oldKey := uint32(dword_5d4594_2516348)
	newKey := protectionRandom.Draw() ^ frame
	dword_5d4594_2516328 = uint32(^newKey)
	count := int(*memmap.PtrUint16(0x587000, 311204))
	head := protectionHead()
	for i := 0; i < count/4; i++ {
		a := GetServer().S().Rand.Logic.IntClamp(0, count/2)
		b := GetServer().S().Rand.Logic.IntClamp(count/2+1, count-1)
		if a != b {
			swapProtectionRecords(protection.At(head, int32(a)), protection.At(head, int32(b)))
		}
	}
	dword_5d4594_2516348 = 0
	dword_5d4594_2516328 = uint32(protection.Rekey(head, oldKey, newKey))
	*memmap.PtrUint32(0x5D4594, 2516364)++
	dword_5d4594_2516348 = uint32(newKey)
	return int32(newKey)
}
