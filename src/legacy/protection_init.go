package legacy

import (
	"time"

	"github.com/brutusmcforce/opennox/common/memmap"
)

func Sub_56F1C0() { initializeProtection() }

func initializeProtection() uint32 {
	protectionRandom.Seed(uint32(time.Now().Unix()))
	dword_5d4594_2516352 = nil
	dword_5d4594_2516348 = uint32(GetServer().S().Frame())
	dword_5d4594_2516344 = nil
	*memmap.PtrUint16(0x587000, 311204) = 0
	dword_5d4594_2516356 = 657757279
	dword_5d4594_2516348 ^= uint32(protectionRandom.Draw())
	dword_5d4594_2516328 = ^dword_5d4594_2516348
	*memmap.PtrUint32(0x5D4594, 2516340) = uint32(nox_xxx_protectionCreateInt_56F400(0))
	sub_56F250()
	result := uint32(nox_xxx_protectionCreateInt_56F400(1))
	*memmap.PtrUint32(0x5D4594, 2516332) = result
	return result
}
