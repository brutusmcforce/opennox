package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func equipmentObject(u unsafe.Pointer) *server.Object { return (*server.Object)(u) }

func equipmentWord(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
func equipmentPlayer(u *server.Object) unsafe.Pointer {
	return unsafe.Pointer(u.UpdateDataPlayer().Player)
}
func equipmentWeaponBits(it *server.Object) uint32 {
	return uint32(nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)))
}
func equipmentArmorBits(it *server.Object) uint32 {
	return uint32(nox_xxx_unitArmorInventoryEquipFlags_415C70(asObjectC(it)))
}
func equipmentContains(u, it *server.Object) bool {
	for p := u.InvFirstItem; p != nil; p = p.InvNextItem {
		if p == it {
			return true
		}
	}
	return false
}
func equipmentStrength(u *server.Object) int32 {
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 0 {
		return int32(*equipmentWord(equipmentPlayer(u), 2239))
	}
	if u.ObjClass&2 == 0 {
		return 0
	}
	if u.ObjSubClass&0x10 != 0 {
		return int32(*(*byte)(unsafe.Add(u.UpdateData, 1324)))
	}
	return 30
}
func equipmentCheckStrength(u, it *server.Object) bool {
	if nox_cheat_allowall != 0 {
		return true
	}
	if u.ObjClass&4 == 0 {
		return false
	}
	strength := equipmentStrength(u)
	var m *server.Modifier
	if it.ObjClass&0x2000000 != 0 {
		m = GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(it.TypeInd))
	} else {
		m = GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(it.TypeInd))
	}
	return m != nil && strength >= int32(m.ReqStrength60)
}

// Only the last two modifier slots participate. The last slot also determines
// the return value even when it has no callback (or is nil).
func equipmentEffects(it, u *server.Object, engage bool) int {
	var result int
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		result = int(uintptr(unsafe.Pointer(m)))
		if m != nil {
			fn := m.Disengage116
			if engage {
				fn = m.Engage112
			}
			if fn != nil {
				result = int(server.CallModifierEffect3Result(fn, m, u, it))
			}
		}
	}
	return result
}
func equipmentDefend(it *server.Object) float64 {
	if it.ObjClass&0x2000000 == 0 {
		return 0
	}
	m := GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(it.TypeInd))
	if m == nil {
		return 0
	}
	value := m.DamageCoeffOrArmor64
	if effect := *(**server.ModifierEff)(it.InitData); effect != nil && effect.Defend76.Fnc != nil {
		// The callback receives a writable float, never a pointer into a Go object.
		p, free := alloc.New(float32(0))
		defer free()
		*p = value
		server.CallModifierEffect6(effect.Defend76.Fnc, effect, it, nil, it, nil, unsafe.Pointer(p))
		value = *p
	}
	return float64(value)
}
func equipmentRecalculate(u *server.Object) int {
	var value float32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x2000000 != 0 && it.ObjFlags&0x100 != 0 {
			value = float32(equipmentDefend(it) + float64(value))
		}
	}
	if float64(value) > math.Float64frombits(uint64(qword_581450_9512)) {
		value = 1
	}
	if u.ObjClass&4 != 0 {
		*(*float32)(unsafe.Add(u.UpdateData, 228)) = value
		return int(uintptr(u.UpdateData))
	}
	if u.ObjClass&2 != 0 {
		*(*float32)(unsafe.Add(u.UpdateData, 2072)) = value
		return int(math.Float32bits(value))
	}
	return int(u.ObjClass)
}
func equipmentNPCSync(u, it *server.Object, value int) unsafe.Pointer {
	u.NeedSync()
	off := 2060
	var bits uint32
	if it.ObjClass&0x1001000 != 0 {
		off, bits = 2056, equipmentWeaponBits(it)
	} else {
		bits = equipmentArmorBits(it)
	}
	p := equipmentWord(u.UpdateData, off)
	if value == 1 {
		*p |= bits
	} else {
		*p &^= bits
	}
	if u.ObjClass&0x20400004 == 0 {
		u.Sub_4E4500(0x4000000, 1024, true)
		return stateSyncEnd(u)
	}
	for i := 0; i < 32; i++ {
		v := equipmentWord(u.CObj(), 560+i*4)
		*v = (*v & 0xfffff000) | 0x4000000
	}
	return unsafe.Add(u.CObj(), int(layoutServerObject.Off(688)) /* Init */)
}
func equipmentCount(u *server.Object, typ int) int {
	n := 0
	if u != nil {
		for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
			if typ == 0 || int(it.TypeInd) == typ && it.ObjFlags&0x20 == 0 {
				n++
			}
		}
	}
	return n
}
func equipmentDuplicate(u, it *server.Object) int {
	if u == nil || it == nil {
		return 0
	}
	for p := u.InvFirstItem; p != nil; p = p.InvNextItem {
		if bool2int(stateEqual(p, it)) != 0 {
			return 1
		}
	}
	return 0
}
func equipmentTryEquip(u, it *server.Object) int {
	if equipmentEquipWeapon(u, it, 1, 1) != 0 || equipmentEquipArmor(u, it, 1, 1) != 0 {
		return 1
	}
	return 0
}
func equipmentTryDequip(u, it *server.Object) int {
	if equipmentDequipWeapon(u, it, 1, 1) != 0 || equipmentDequipArmor(u, it, 1, 1) != 0 {
		return 1
	}
	return 0
}
func equipmentSaveShield(u *server.Object) {
	if u == nil {
		return
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x2000000 != 0 && it.ObjFlags&0x100 != 0 && equipmentArmorBits(it)&0x3000000 != 0 {
			*(*unsafe.Pointer)(unsafe.Pointer(equipmentWord(equipmentPlayer(u), 2500))) = it.CObj()
		}
	}
}
func equipmentFindShield(u *server.Object) *server.Object {
	if u != nil {
		for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
			if it.ObjClass&0x2000000 != 0 && it.ObjFlags == 16 {
				return it
			}
		}
	}
	return nil
}
func equipmentPickupSound(u, it *server.Object) {
	if u == nil || it == nil {
		return
	}
	sound := 0
	if it.ObjClass&0x1000 != 0 {
		sound = 830
	} else if it.Material&16 != 0 {
		sound = 842
	} else if it.Material&8 != 0 {
		sound = 844
	}
	if sound != 0 {
		inventorySound(sound, u, 0, 0)
	}
}
func equipmentDropSound(it *server.Object) {
	if it == nil {
		return
	}
	sound := 0
	if it.ObjClass&0x1000 != 0 {
		sound = 831
	} else if it.Material&16 != 0 {
		sound = 843
	} else if it.Material&8 != 0 {
		sound = 845
	}
	if sound != 0 {
		inventorySound(sound, it, 0, 0)
	}
}
func equipmentArmorDropSound(it *server.Object) {
	if it == nil {
		return
	}
	sound := 0
	switch {
	case it.Material&16 != 0:
		sound = 805
	case it.Material&8 != 0:
		sound = 811
	case it.Material&4 != 0:
		sound = 808
	case it.Material&2 != 0:
		sound = 814
		if it.ObjSubClass&32 != 0 {
			sound = 817
		}
	}
	if sound != 0 {
		inventorySound(sound, it, 0, 0)
	}
}
func equipmentSecondary(u, it *server.Object) {
	if u == nil {
		return
	}
	if it != nil && (!Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) || !equipmentCheckStrength(u, it)) {
		gameplayReportSecondary(int(int32(uint8(u.UpdateDataPlayer().Player.PlayerInd))), nil, byte(1))
	}
	*(*unsafe.Pointer)(unsafe.Pointer(equipmentWord(u.UpdateData, 108))) = it.CObj()
}
func equipmentInitDropTable() {
	for off := uintptr(279432); ; off += 12 {
		name := memmap.LoadPtr(0x587000, off)
		if name == nil {
			break
		}
		*memmap.PtrUint32(0x587000, off+4) = uint32(GetServer().S().Types.IndByID(alloc.GoString((*byte)(name))))
	}
	dword_5d4594_2488728 = 1
}
func equipmentDropPolicy(it *server.Object, mask int) int {
	if it == nil {
		return 0
	}
	if dword_5d4594_2488728 == 0 {
		equipmentInitDropTable()
	}
	for off := uintptr(279432); *memmap.PtrUint32(0x587000, off) != 0; off += 12 {
		if *memmap.PtrUint32(0x587000, off+4) == uint32(it.TypeInd) {
			return bool2int(*memmap.PtrUint32(0x587000, off+8)&uint32(mask) != 0)
		}
	}
	return 0
}
