package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/protection"
)

func nox_xxx_playerAwardSpellProtectionCRC_56FCE0(id, index, enabled int32) int32 {
	if id < 657757279 {
		return id
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return id
	}
	old := r.Value
	r.Value = ((old ^ key) | protection.Bit(int32(index), int32(enabled))) ^ key
	dword_5d4594_2516328 ^= uint32(old ^ r.Value)
	return int32(r.Value)
}

func nox_xxx_playerApplyProtectionCRC_56FD50(id int32, data unsafe.Pointer, count int32) int32 {
	if id < 657757279 {
		return 0
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	var flags uint32
	if count > 1 {
		flags = protection.Flags(unsafe.Slice((*int32)(data), int(count)))
	}
	if flags^key == r.Value {
		return 1
	}
	return 0
}
