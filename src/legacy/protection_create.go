package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/protection"
)

func createProtectionRecord(id, value uint32) int32 {
	r := (*protection.Record)(legacyCalloc(1, uintptr(unsafe.Sizeof(protection.Record{}))))
	if !protection.Initialize(r, id, value, uint32(dword_5d4594_2516348), (*uint32)(unsafe.Pointer(&dword_5d4594_2516328))) {
		return 0
	}
	return insertProtectionRecord(r)
}
