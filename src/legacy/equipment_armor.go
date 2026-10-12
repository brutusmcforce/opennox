package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func equipmentArmorMask(it *server.Object) int {
	return bool2int(it.ObjClass&0x2000000 == 0 || equipmentArmorBits(it)&0xc0d == 0)
}
func equipmentNPCDequipArmor(u, it *server.Object) int {
	if it.ObjClass&0x2000000 == 0 || it.ObjFlags&0x100 == 0 || !equipmentContains(u, it) {
		return 0
	}
	it.ObjFlags &^= 0x100
	if u.ObjSubClass&0x10 != 0 {
		equipmentNPCSync(u, it, 0)
	}
	it.ObjFlags &^= 0x10000000
	equipmentRecalculate(u)
	equipmentEffects(it, u, false)
	return 1
}
func equipmentDequipArmor(u, it *server.Object, report, broadcast int) int {
	if it.ObjClass&0x2000000 == 0 || it.ObjFlags&0x100 == 0 {
		return 0
	}
	flags := it.ObjFlags
	if u.ObjClass&2 != 0 {
		return equipmentNPCDequipArmor(u, it)
	}
	if u.ObjClass&4 == 0 || !equipmentContains(u, it) {
		return 0
	}
	ud := u.UpdateData
	it.ObjFlags = flags &^ 0x10000100
	*equipmentWord(equipmentPlayer(u), 0) &^= equipmentArmorBits(it)
	equipmentDequipReports(u, it, report, broadcast)
	equipmentRecalculate(u)
	equipmentEffects(it, u, false)
	// This is the original byte at object+48 (team linkage), not subclass+12.
	if *(*byte)(unsafe.Add(it.CObj(), 48))&2 != 0 {
		state := *(*byte)(unsafe.Add(ud, 88))
		if state >= 15 && state <= 17 {
			Nox_xxx_playerSetState_4FA020(u, 13)
		}
	}
	return 1
}
func equipmentNPCEquipArmor(u, it *server.Object) int {
	if it.ObjClass&0x2000000 == 0 || it.ObjFlags&0x100 != 0 || !equipmentContains(u, it) {
		return 0
	}
	if old := equipmentSameArmor(u, it); old != nil {
		equipmentNPCDequipArmor(u, old)
	}
	it.ObjFlags |= 0x100
	if u.ObjSubClass&0x10 != 0 {
		equipmentNPCSync(u, it, 1)
	}
	if equipmentArmorMask(it) == 0 {
		it.ObjFlags |= 0x10000000
	}
	equipmentRecalculate(u)
	equipmentEffects(it, u, true)
	if it.ObjSubClass&2 != 0 {
		equipmentRemoveWeapons(u)
	}
	return 1
}
func equipmentRemoveWeapons(u *server.Object) {
	if u == nil {
		return
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 != 0 && it.ObjClass&0x1001000 != 0 && equipmentWeaponBits(it)&0x7ffe40c != 0 {
			equipmentDequipWeapon(u, it, 1, 1)
		}
	}
}
func equipmentEquipArmor(u, it *server.Object, report, broadcast int) int {
	if it.ObjClass&0x2000000 == 0 || it.ObjFlags&0x100 != 0 {
		return 0
	}
	if u.ObjClass&2 != 0 {
		return equipmentNPCEquipArmor(u, it)
	}
	if u.ObjClass&4 == 0 {
		return 0
	}
	old := equipmentSameArmor(u, it)
	if !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		return equipmentAdmissionFail(u, "armor.c:ArmorEquipClassFail", broadcast)
	}
	if !equipmentCheckStrength(u, it) {
		return equipmentAdmissionFail(u, "armor.c:ArmorEquipStrengthFail", broadcast)
	}
	if old != nil {
		equipmentDequipArmor(u, old, 1, 1)
	}
	it.ObjFlags |= 0x100
	*equipmentWord(equipmentPlayer(u), 0) |= equipmentArmorBits(it)
	gameplayReportEquip(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), it, int(report))
	if equipmentArmorMask(it) == 0 {
		it.ObjFlags |= 0x10000000
	}
	equipmentRecalculate(u)
	equipmentEffects(it, u, true)
	if it.ObjSubClass&2 != 0 {
		equipmentRemoveWeapons(u)
	}
	return 1
}
func equipmentSameArmor(u, it *server.Object) *server.Object {
	for old := u.InvFirstItem; old != nil; old = old.InvNextItem {
		if old.ObjFlags&0x100 != 0 && old.ObjClass&0x2000000 != 0 && old.ObjSubClass == it.ObjSubClass {
			return old
		}
	}
	return nil
}
