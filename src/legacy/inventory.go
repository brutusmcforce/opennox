package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func inventoryInt(u *server.Object) int32 { return int32(memmap.PtrWord(u.CObj())) }
func inventoryWeight(u *server.Object) {
	var weight uint32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		weight += uint32(it.Weight)
	}
	u.UpdateDataPlayer().Player.Field3656 = uint32(bool2int(weight > uint32(u.CarryCapacity)))
}
func inventoryRemove(u, it *server.Object) {
	if u == nil || it == nil {
		return
	}
	report := int32(1)
	if u.ObjClass&4 != 0 {
		pl := u.UpdateDataPlayer().Player
		if !noxflags.HasGame(4096) && u.ObjFlags&0x8000 != 0 && it.ObjClass&0x13001000 != 0 {
			report = 0
		}
		if it.ObjClass&0x10000000 != 0 && noxflags.HasGame(32) {
			pl.WeaponEquip &^= 1
			if report == 1 {
				gameplayReportDequipFlags(255, it)
			}
		}
		equipmentDequipArmor(u, it, 0, int(report))
		equipmentDequipWeapon(u, it, 0, int(report))
		gameplayReportDrop(int(uint8(pl.PlayerInd)), it)
		toggleProtectionObject(int32(pl.Prot4632), it)
	} else if u.ObjClass&2 != 0 {
		if u.ObjSubClass&0x10 != 0 && it.ObjClass&0x10000000 != 0 && noxflags.HasGame(32) {
			equipmentNPCSync(u, it, 0)
		}
		equipmentDequipArmor(u, it, 1, 1)
		equipmentDequipWeapon(u, it, 1, 1)
	}
	if prev := it.Field125; prev != nil {
		prev.InvNextItem = it.InvNextItem
	} else {
		u.InvFirstItem = it.InvNextItem
	}
	if next := it.InvNextItem; next != nil {
		next.Field125 = it.Field125
	}
	it.InvHolder = nil
	GetServer().S().ObjClearOwner(it)
	if u.ObjClass&4 != 0 {
		inventoryWeight(u)
	}
}
func inventoryInsert(u, it *server.Object, report int) {
	if u == nil || it == nil || u.ObjFlags&0x20 != 0 || it.ObjFlags&0x20 != 0 {
		return
	}
	it.Field125 = nil
	it.InvNextItem = u.InvFirstItem
	if next := u.InvFirstItem; next != nil {
		next.Field125 = it
	}
	u.InvFirstItem = it
	it.InvHolder = u
	GetServer().S().ObjSetOwner(u, it)
	if u.ObjClass&4 != 0 {
		pl := u.UpdateDataPlayer().Player
		if report != 0 {
			gameplayReportPickup(int(uint8(pl.PlayerInd)), it)
		}
		toggleProtectionObject(int32(pl.Prot4632), it)
		inventoryWeight(u)
	}
	if it.ObjClass&64 != 0 {
		nox_xxx_aud_501960(820, asObjectC(u), 0, 0)
	}
}
func inventoryDropEligible(u, it *server.Object) bool {
	return it.ObjFlags&0x20 != 0 || u.ObjClass&6 == 0 || it.ObjFlags&0x10000000 == 0
}
func inventoryDroppable(it *server.Object) bool {
	if it == nil {
		return false
	}
	if dword_5d4594_2488728 == 0 {
		equipmentInitDropTable()
	}
	for off := uintptr(279432); *memmap.PtrUint32(0x587000, off) != 0; off += 12 {
		if *memmap.PtrUint32(0x587000, off+4) == uint32(it.TypeInd) {
			return true
		}
	}
	return false
}
