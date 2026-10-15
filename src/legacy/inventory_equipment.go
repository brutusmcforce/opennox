package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func inventoryPriMessage(u *server.Object, key string) {
	gameplayTextPrivate(u, (*byte)(unsafe.Pointer(internCStr(key))), 0)
}

// AutoEquip reports if weapons and armor that the player unit u picks up are equipped automatically,
// as the game always did (opennox#741).
var AutoEquip = func(u *server.Object) bool { return true }

func inventoryWeaponPickup(u, it *server.Object, arg, equip int) int {
	if u.ObjClass&4 != 0 && noxflags.HasGame(4096) && it.ObjSubClass&0x200000 != 0 {
		limit := float32(nox_xxx_gamedataGetFloat_419D40(internCStr("ForceOfNatureStaffLimit")))
		if int32(equipmentCount(u, int(it.ObjFlags))) >= floatToInt32(limit) {
			inventoryPriMessage(u, "pickup.c:MaxSameItem")
			inventorySound(925, u, 0, 0)
			return 0
		}
	}
	if !noxflags.HasGame(2048|4096) && serverConfigFlagsQuery(int32(2)) != 0 {
		duplicate := it.ObjSubClass&0x82 == 0 && equipmentDuplicate(u, it) != 0
		if it.ObjSubClass&0x40 != 0 {
			for owned := u.Field129; owned != nil; owned = owned.Field128 {
				if owned.ObjClass&0x1000000 != 0 && owned.ObjSubClass&0x40 != 0 {
					duplicate = true
					break
				}
			}
		}
		if duplicate {
			if u.ObjClass&4 != 0 {
				inventoryPriMessage(u, "weapon.c:CannotPickupDuplicateWeapon")
				inventorySound(925, u, 2, int(u.NetCode))
			}
			return 0
		}
	}
	if !noxflags.HasGame(2048|4096) && u.ObjClass&4 != 0 && !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		inventoryPriMessage(u, "weapon.c:WeaponEquipClassFail")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryPickup(u, it, arg) != 1 {
		return 0
	}
	if u.ObjClass&4 != 0 {
		ud := u.UpdateData
		equipped := 0
		auto := AutoEquip(u)
		if auto && layout.As[server.PlayerUpdateData](ud).Field26 == nil && sub_419E60(asObjectC(u)) == 0 && nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) != 2 {
			equipped = equipmentEquipWeapon(u, it, equip, 0)
		}
		if auto && sub_419E60(asObjectC(u)) == 0 && nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) == 2 {
			flags := u.UpdateDataPlayer().Player.WeaponEquip
			if flags&0xC != 0 && flags&2 == 0 {
				equipped = equipmentEquipWeapon(u, it, equip, 0)
			}
		}
		if equipped == 0 {
			if it.ObjClass&0x1000 != 0 && it.ObjSubClass&0x47F0000 != 0 {
				equipmentReportWandCharges(u, it)
			} else if it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x82 != 0 {
				equipmentReportAmmoCharges(u, it)
			}
		}
	}
	equipmentPickupSound(u, it)
	motionDecayRemove(it)
	return 1
}
func inventoryAmmoPickup(u, it *server.Object, arg, equip int) int {
	bits := uint32(nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)))
	if u.ObjClass&4 != 0 && bits&0x82 != 0 {
		data := it.UseDataAmmo()
		mods := [4]*server.ModifierEff(itemModifiers(it))
		for old := u.InvFirstItem; old != nil; old = old.InvNextItem {
			if old.TypeInd != it.TypeInd || old.ObjClass&0x1000000 == 0 || uint32(nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(old)))&bits == 0 {
				continue
			}
			prev := old.UseDataAmmo()
			same := [4]*server.ModifierEff(itemModifiers(old)) == mods
			if !same || prev.Infinite != 0 || int(data.MaxCharges)+int(prev.MaxCharges) > 250 {
				continue
			}
			prev.Charges += data.Charges
			prev.MaxCharges += data.MaxCharges
			gameplayReportCharges(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), old, prev.Charges, prev.MaxCharges)
			GetServer().DelayedDelete(it)
			equipmentPickupSound(u, it)
			return 1
		}
	}
	return inventoryWeaponPickup(u, it, arg, equip)
}
func inventoryOblivionPickup(u, it *server.Object, arg, equip int) int {
	rv := inventoryWeaponPickup(u, it, arg, equip)
	if rv == 1 && u.ObjClass&4 != 0 && sub_419E60(asObjectC(u)) == 0 {
		for i, key := range []string{"weapon.c:PickupHalberdOblivion", "weapon.c:PickupHeartOblivion", "weapon.c:PickupWierdlingOblivion", "weapon.c:PickupOrbOblivion"} {
			if uint32(it.ObjSubClass)&(0x800000<<i) != 0 {
				inventoryPriMessage(u, key)
				inventorySound(914+i, u, 0, 0)
				break
			}
		}
		runtimePauseStart(u, 1)
		equipmentTryEquip(u, it)
	}
	return rv
}
func inventoryPlain(it *server.Object) bool {
	for _, v := range itemModifiers(it) {
		if v != nil {
			return false
		}
	}
	return true
}
func inventoryArmorPickup(u, it *server.Object, arg, equip int) int {
	if *memmap.PtrUint32(0x5D4594, 2488712) == 0 {
		*memmap.PtrUint32(0x5D4594, 2488712) = uint32(GetServer().S().Types.IndByID("StreetSneakers"))
		*memmap.PtrUint32(0x5D4594, 2488716) = uint32(GetServer().S().Types.IndByID("WizardRobe"))
		dword_5d4594_2488720 = uint32(GetServer().S().Types.IndByID("WoodenShield"))
		dword_5d4594_2488724 = uint32(GetServer().S().Types.IndByID("SteelShield"))
	}
	if !noxflags.HasGame(2048|4096) && serverConfigFlagsQuery(int32(2)) != 0 && equipmentDuplicate(u, it) != 0 {
		inventoryPriMessage(u, "armor.c:CannotPickupDuplicateArmor")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if !noxflags.HasGame(2048|4096) && u.ObjClass&4 != 0 && !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		inventoryPriMessage(u, "armor.c:ArmorEquipClassFail")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryPickup(u, it, arg) != 1 {
		return 0
	}
	if u.ObjClass&4 != 0 {
		old := equipmentSameArmor(u, it)
		if sub_419E60(asObjectC(u)) == 0 && AutoEquip(u) {
			wood, steel := uint32(dword_5d4594_2488720), uint32(dword_5d4594_2488724)
			sneakers, robe := *memmap.PtrUint32(0x5D4594, 2488712), *memmap.PtrUint32(0x5D4594, 2488716)
			replace := false
			switch {
			case old != nil && uint32(old.TypeInd) != wood && uint32(old.TypeInd) != steel && uint32(old.TypeInd) != sneakers && uint32(old.TypeInd) != robe:
			case old != nil && uint32(old.TypeInd) == robe:
				replace = inventoryPlain(old)
			case it.ObjSubClass&2 == 0:
				replace = true
			case old == nil:
				weapon := *(**server.Object)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](u.UpdateData).Field26))
				replace = weapon == nil || weapon.ObjSubClass&0x7FFE40C == 0
			case uint32(old.TypeInd) == wood:
				replace = inventoryPlain(old)
			case uint32(old.TypeInd) == steel:
				if uint32(it.TypeInd) == wood {
					replace = inventoryPlain(old) && !inventoryPlain(it)
				} else if uint32(it.TypeInd) == steel {
					replace = inventoryPlain(old)
				}
			}
			if replace {
				equipmentEquipArmor(u, it, equip, 0)
			}
		}
	}
	switch {
	case it.Material&0x10 != 0:
		inventorySound(804, u, 0, 0)
	case it.Material&8 != 0:
		inventorySound(810, u, 0, 0)
	case it.Material&4 != 0:
		inventorySound(807, u, 0, 0)
	case it.Material&2 != 0:
		sound := 813
		if it.ObjSubClass&0x20 != 0 {
			sound = 816
		}
		inventorySound(sound, u, 0, 0)
	}
	motionDecayRemove(it)
	return 1
}
