package legacy

import "github.com/brutusmcforce/opennox/internal/protection"

func addProtectionRecord(id int32, delta uint32) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	old := r.Value
	value := (old ^ key) + delta
	r.Value = value ^ key
	dword_5d4594_2516328 ^= uint32(old ^ r.Value)
	return uint32(nox_xxx_protectData_56F5C0())
}
