package legacy

import "github.com/brutusmcforce/opennox/internal/protection"

func updateProtectionFloat(id int32, value float32, add bool) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	old := r.Value
	var next uint32
	if add {
		next = protection.AddFloatValue(old^key, value)
	} else {
		next = protection.FloatValue(value)
	}
	r.Value = next ^ key
	dword_5d4594_2516328 ^= uint32(old ^ r.Value)
	return uint32(nox_xxx_protectData_56F5C0())
}
