package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// The legacy sync ABI returns the one-past pointer for the 32 player words.
func stateSyncEnd(u *server.Object) unsafe.Pointer { return unsafe.Pointer(&u.Init) }
func stateSync(u *server.Object, flag, bit uint32) unsafe.Pointer {
	if u.ObjClass&0x20400004 != 0 {
		for i := range u.Field140 {
			u.Field140[i] = u.Field140[i]&0xfffff000 | flag
		}
	} else {
		u.Sub_4E4500(flag, bit, u.Sub_4E4C90(uint(bit)))
	}
	return stateSyncEnd(u)
}
func stateOnOff(u *server.Object, on bool) unsafe.Pointer {
	u.NeedSync()
	if on {
		u.ObjFlags |= 0x1000000
	} else {
		u.ObjFlags &^= 0x1000000
	}
	return stateSync(u, 0x40000, 4)
}
func stateRaise(u *server.Object, z float32) {
	if u.ZVal != z {
		u.NeedSync()
		u.ZVal = z
		stateSync(u, 0x400000, 64)
	}
}
func stateAnimation(u *server.Object, frame uint32) unsafe.Pointer {
	u.NeedSync()
	u.Field33 = frame
	return stateSync(u, 0x10000, 1)
}
func stateBuffs(u *server.Object, flags uint32) unsafe.Pointer {
	u.NeedSync()
	u.Buffs = flags
	if u.ObjClass&4 != 0 {
		pl := *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 276))
		setProtectionRecord(int32(*equipmentWord(pl, 4612)), flags)
	}
	return stateSync(u, 0x800000, 128)
}
func stateType(offset uintptr, name string) uint32 {
	p := memmap.PtrUint32(0x5d4594, offset)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID(name))
	}
	return *p
}
func stateAttributes(u *server.Object, attrs unsafe.Pointer) unsafe.Pointer {
	if u.ObjClass&0x1000 == 0 || u.ObjSubClass&0x47f0000 == 0 {
		any := false
		for _, v := range unsafe.Slice((*uint32)(attrs), 4) {
			any = any || v != 0
		}
		if !any {
			return unsafe.Add(attrs, 16)
		}
	}
	typ := stateType(1564960, "TeamBase")
	if u.ObjClass&0x13001000 == 0 && uint32(u.TypeInd) != typ {
		return unsafe.Pointer(uintptr(typ))
	}
	u.NeedSync()
	copy(unsafe.Slice((*byte)(u.InitData), 20), unsafe.Slice((*byte)(attrs), 20))
	return stateSync(u, 0x2000000, 512)
}
func stateOn(u *server.Object) int8 {
	if u.ObjFlags&0x1000000 == 0 && u.ObjClass&0x4000 != 0 {
		nox_xxx_aud_501960(235, asObjectC(u), 0, 0)
	}
	stateOnOff(u, true)
	if u.ObjClass&0x10042000 != 0 {
		u.ObjFlags &^= 0x40
	}
	if u.ObjClass&1 == 0 {
		return int8(collisionActivate(u))
	}
	return int8(u.ObjClass)
}
func stateOff(u *server.Object) int32 {
	if u.ObjFlags&0x1000000 != 0 && u.ObjClass&0x4000 != 0 {
		nox_xxx_aud_501960(236, asObjectC(u), 0, 0)
	}
	stateOnOff(u, false)
	if u.ObjClass&0x10042000 != 0 {
		u.ObjFlags |= 0x40
		return int32(u.ObjFlags)
	}
	return int32(u.ObjClass)
}
func stateChecksum(u *server.Object) uint32 {
	var sum uint32
	for _, off := range [...]int{340, 248, 120, 128, 132, 136, 148, 152, 108, 104, 100, 96, 92, 88, 84, 80, 76, 72, 68, 64, 60, 16, 20, 36, 40, 44, 56} {
		sum ^= *equipmentWord(u.CObj(), off)
	}
	sum ^= uint32(int32(int16(u.Direction1))) ^ uint32(int32(int16(u.Direction2))) ^ uint32(u.TypeInd) ^ uint32(*(*byte)(unsafe.Add(u.CObj(), 52)))
	if h := u.HealthData; h != nil {
		for _, v := range unsafe.Slice((*uint16)(unsafe.Pointer(h)), 3) {
			sum ^= uint32(v)
		}
	}
	return sum
}
