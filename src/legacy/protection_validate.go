package legacy

import (
	"github.com/brutusmcforce/opennox/internal/protection"
	"unsafe"
)

func protectionValidateString(data unsafe.Pointer, size uint32, id int32) int32 {
	if id < 657757279 {
		return 0
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	if uint32(key)^uint32(protectionStringChecksum(data, size)) != uint32(r.Value) {
		return 0
	}
	return 1
}
