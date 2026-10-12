package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func equipmentNPCDequipWeapon(u, it *server.Object) int {
	ud := u.UpdateData
	if it.ObjClass&0x1001000 == 0 || it.ObjFlags&0x100 == 0 || !equipmentContains(u, it) {
		return 0
	}
	*(*byte)(unsafe.Add(ud, 2068)) = 0
	it.ObjFlags &^= 0x100
	if u.ObjSubClass&0x10 != 0 {
		equipmentNPCSync(u, it, 0)
	}
	if it.ObjSubClass&0xc != 0 {
		equipmentDequipAmmo(u, 1, 1)
	}
	if it.ObjSubClass&2 == 0 {
		*equipmentWord(ud, 2064) = 0
	}
	equipmentEffects(it, u, false)
	spellLifeCancelWand(u, it)
	return 1
}
func equipmentDequipAmmo(u *server.Object, report, broadcast int) {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 != 0 && equipmentWeaponBits(it) == 2 {
			equipmentDequipWeapon(u, it, report, broadcast)
			return
		}
	}
}
func equipmentDequipReports(u, it *server.Object, report, broadcast int) {
	if report != 0 {
		gameplayReportDequipItem(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), it)
	}
	if broadcast != 0 {
		gameplayReportDequipFlags(255, it)
	}
}
func equipmentDequipWeapon(u, it *server.Object, report, broadcast int) int {
	bits := equipmentWeaponBits(it)
	if u.ObjClass&2 != 0 {
		return equipmentNPCDequipWeapon(u, it)
	}
	if u.ObjClass&4 == 0 || it.ObjClass&0x1001000 == 0 {
		return 0
	}
	ud := u.UpdateData
	if *(*byte)(unsafe.Add(ud, 88)) == 1 {
		Nox_xxx_playerSetState_4FA020(u, 13)
	}
	active := *(**server.Object)(unsafe.Add(ud, 104))
	if active == nil || active != it && bits != 2 {
		return 0
	}
	if bits&0xc != 0 {
		equipmentDequipAmmo(u, report, broadcast)
	}
	spellLifeCancelWand(u, it)
	if bits == 2 {
		it.ObjFlags &^= 0x100
		*equipmentWord(equipmentPlayer(u), 4) &^= 2
		equipmentDequipReports(u, it, report, broadcast)
	} else if active = *(**server.Object)(unsafe.Add(ud, 104)); active != nil {
		active.ObjFlags &^= 0x100
		*equipmentWord(equipmentPlayer(u), 4) &^= equipmentWeaponBits(active)
		if report != 0 {
			gameplayReportDequipItem(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), active)
		}
		if broadcast != 0 {
			gameplayReportDequipFlags(255, it)
		}
		*equipmentWord(ud, 104) = 0
	}
	equipmentEffects(it, u, false)
	// Keep the community shield-selection behavior controlled by gameex bit 2.
	if gameex_flags&2 != 0 {
		equipmentSaveShield(u)
		secondary := *(**server.Object)(unsafe.Add(u.UpdateData, 108))
		if secondary == nil || equipmentWeaponBits(secondary)&0x7ffe40c == 0 {
			saved := *(**server.Object)(unsafe.Add(equipmentPlayer(u), 2500))
			if saved != nil && byte(saved.ObjFlags) == 16 {
				equipmentTryEquip(u, saved)
			} else if found := equipmentFindShield(u); found != nil {
				equipmentTryEquip(u, found)
			}
		}
	}
	return 1
}
func equipmentNPCEquipWeapon(u, it *server.Object) int {
	ud := u.UpdateData
	if it.ObjClass&0x1001000 == 0 || it.ObjFlags&0x100 != 0 || !equipmentContains(u, it) {
		return 0
	}
	*(*byte)(unsafe.Add(ud, 2068)) = 0
	if it.ObjSubClass&0xc == 0 {
		equipmentDequipAmmo(u, 1, 1)
	}
	sub := it.ObjSubClass
	if sub&2 == 0 {
		for old := u.InvFirstItem; old != nil; old = old.InvNextItem {
			if old.ObjFlags&0x100 != 0 && old.ObjClass&0x1001000 != 0 && (sub&0xc == 0 || old.ObjSubClass&2 == 0) {
				equipmentNPCDequipWeapon(u, old)
				break
			}
		}
	}
	it.ObjFlags |= 0x100
	if u.ObjSubClass&0x10 != 0 {
		equipmentNPCSync(u, it, 1)
	}
	if it.ObjSubClass&2 == 0 {
		*equipmentWord(ud, 2064) = uint32(uintptr(it.CObj()))
	}
	equipmentEffects(it, u, true)
	if equipmentWeaponBits(it)&0x7ffe40c != 0 {
		equipmentRemoveShields(u)
	}
	return 1
}
func equipmentRemoveShields(u *server.Object) {
	if u == nil {
		return
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x2000000 != 0 && it.ObjFlags&0x100 != 0 && equipmentArmorBits(it)&0x3000000 != 0 {
			equipmentDequipArmor(u, it, 1, 1)
		}
	}
}
func equipmentAdmissionFail(u *server.Object, key string, broadcast int) int {
	inventoryPriMessage(u, key)
	if broadcast != 0 {
		inventorySound(925, u, 2, int(u.NetCode))
	}
	return 0
}
func equipmentEquipWeapon(u, it *server.Object, report, broadcast int) int {
	bits := equipmentWeaponBits(it)
	if it.ObjClass&0x1001000 == 0 || it.ObjFlags&0x100 != 0 {
		return 0
	}
	if u.ObjClass&2 != 0 {
		return equipmentNPCEquipWeapon(u, it)
	}
	if u.ObjClass&4 == 0 {
		return 0
	}
	ud := u.UpdateData
	if nox_xxx_probablyWarcryCheck_4FC3E0(asObjectC(u), 2) != 0 || nox_xxx_probablyWarcryCheck_4FC3E0(asObjectC(u), 1) != 0 {
		return 0
	}
	if !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		return equipmentAdmissionFail(u, "weapon.c:WeaponEquipClassFail", broadcast)
	}
	if !equipmentCheckStrength(u, it) {
		return equipmentAdmissionFail(u, "weapon.c:WeaponEquipStrengthFail", broadcast)
	}
	if !equipmentContains(u, it) {
		return 0
	}
	if *(*byte)(unsafe.Add(ud, 88)) == 1 {
		Nox_xxx_playerSetState_4FA020(u, 13)
	}
	if bits == 2 {
		if *equipmentWord(equipmentPlayer(u), 4)&0xc == 0 && equipmentEquipBow(u) == 0 {
			return equipmentAdmissionFail(u, "weapon.c:BowNotFound", broadcast)
		}
		equipmentDequipAmmo(u, 1, 1)
	}
	active := *(**server.Object)(unsafe.Add(ud, 104))
	if active != nil && bits != 2 && equipmentDequipWeapon(u, active, 1, 1) == 0 {
		return 0
	}
	it.ObjFlags |= 0x100
	*equipmentWord(equipmentPlayer(u), 4) |= bits
	gameplayReportEquip(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), it, int(report))
	if bits != 2 {
		*equipmentWord(ud, 104) = uint32(uintptr(it.CObj()))
	}
	if it.ObjClass&0x1000 != 0 && it.ObjSubClass&0x47f0000 != 0 {
		equipmentReportCharges(u, it, 108, 109)
	} else if it.ObjClass&0x1000000 != 0 {
		if bits&0x82 != 0 {
			equipmentReportCharges(u, it, 1, 0)
		} else if bits&0xc != 0 {
			*(*byte)(it.UseData.Ptr) = 0
		}
	}
	equipmentEffects(it, u, true)
	if bits&0x7ffe40c != 0 {
		equipmentRemoveShields(u)
	}
	return 1
}
func equipmentReportCharges(u, it *server.Object, a, b int) {
	gameplayReportCharges(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), it, byte(*(*byte)(unsafe.Add(it.UseData.Ptr, a))), byte(*(*byte)(unsafe.Add(it.UseData.Ptr, b))))
}
func equipmentEquipBow(u *server.Object) int {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if equipmentWeaponBits(it)&0xc != 0 {
			return equipmentEquipWeapon(u, it, 1, 1)
		}
	}
	return 0
}
