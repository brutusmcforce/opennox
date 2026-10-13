package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

func stateRemoveSpawned(u *server.Object) {
	if u.ObjClass&4 != 0 {
		pl := layout.As[server.PlayerUpdateData](u.UpdateData).Player
		switch pl.InfoBytes[66] {
		case 1:
			controlRemoveGlyphs(u)
		case 2:
			controlRemoveCreatures(u)
		}
	}
	for it := u.Field129; it != nil; {
		next := it.Field128
		if it.ObjClass&1 != 0 || sub_4E3B80(int(int32(it.TypeInd))) == 0 {
			GetServer().DelayedDelete(it)
		}
		it = next
	}
}
func stateIsUnit(u *server.Object) bool {
	return !bool(nox_common_gameFlags_check_40A5C0(0x2000)) && u.ObjClass&2 != 0 && u.ObjSubClass&0x100 != 0
}
func stateIsPixie(u *server.Object) bool {
	typ := stateType(1565592, "Pixie")
	if u == nil || u.ObjClass&1 == 0 || !bool(nox_common_gameFlags_check_40A5C0(2048)) || uint32(u.TypeInd) != typ {
		return false
	}
	owner := u.FindOwnerChainPlayer()
	return owner != nil && owner.ObjClass&4 != 0
}
func stateCleanup(mode int32) {
	typ := stateType(1565596, "Moonglow")
	for u := GetServer().S().Objs.List; u != nil; {
		next := u.ObjNext
		if mode == 0 || u.ObjClass&4 == 0 && (u.InvHolder == nil || u.InvHolder.ObjClass&4 == 0) &&
			(uint32(u.TypeInd) != typ || u.ObjOwner == nil || u.ObjOwner.ObjClass&4 == 0) && !stateIsUnit(u) {
			GetServer().DelayedDelete(u)
		}
		u = next
	}
	for u := GetServer().S().Objs.MissileList; u != nil; {
		next := u.ObjNext
		if mode != 1 || !stateIsPixie(u) {
			GetServer().DelayedDelete(u)
		}
		u = next
	}
}
func stateRememberAttacker(u, t *server.Object) {
	if u == nil || t == nil || u == t || u.ObjClass&4 == 0 || t.ObjClass&4 == 0 {
		return
	}
	from := layout.As[server.PlayerUpdateData](u.UpdateData).Player
	to := layout.As[server.PlayerUpdateData](t.UpdateData).Player
	to.Field3604 = uint32(from.PlayerInd)
	to.Field3608 = GetServer().S().Frame()
	to.Field3600 = 1
}
func stateFreeze(u *server.Object, force int32) int8 {
	result := int8(u.ObjFlags)
	if u.ObjFlags&2 != 0 {
		return result
	}
	u.ObjFlags |= 2
	if u.ObjClass&4 != 0 {
		p := memmap.PtrUint32(0x5d4594, 1567712)
		if *p == 0 {
			*p = uint32(force)
		}
		gameplayReportPlayerStatus(u)
		nox_xxx_playerSetState_4FA020(asObjectC(u), 13)
		stateRaise(u, 0)
		sub_50B510()
		for it := u.Field129; it != nil; it = it.Field128 {
			if it.ObjClass&2 != 0 && *(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](it.UpdateData).StatusFlags))&0x80 != 0 {
				stateFreeze(it, force)
			}
		}
	}
	result = int8(u.ObjClass)
	if u.ObjClass&2 != 0 {
		result = int8(u.ObjFlags)
		if u.ObjFlags&0x8000 == 0 {
			result = int8(memmap.PtrWord(unsafe.Pointer(u.MonsterPushAction(0))))
		}
	}
	return result
}
func stateUnfreeze(u *server.Object, force int32) int8 {
	result := int8(u.ObjFlags)
	if u.ObjFlags&2 == 0 {
		return result
	}
	if u.ObjClass&4 != 0 {
		p := memmap.PtrUint32(0x5d4594, 1567712)
		if *p != 0 && force == 0 {
			return int8(*p)
		}
		*p = 0
		u.ObjFlags &^= 2
		result = int8(gameplayReportPlayerStatus(u))
		for it := u.Field129; it != nil; it = it.Field128 {
			if it.ObjClass&2 != 0 {
				result = int8(memmap.PtrWord(it.UpdateData))
				if *(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](it.UpdateData).StatusFlags))&0x80 != 0 {
					result = stateUnfreeze(it, force)
				}
			}
		}
	} else {
		u.ObjFlags &^= 2
		result = int8(u.ObjFlags)
	}
	if u.ObjClass&2 != 0 {
		result = int8(u.ObjFlags)
		if u.ObjFlags&0x8000 == 0 {
			result = int8(u.MonsterPopAction())
		}
	}
	return result
}
func statePet(u, t *server.Object) {
	if u == nil || t == nil {
		return
	}
	t.ObjSubClass |= 0x80
	pl := u.UpdateDataPlayer().Player
	ind := int32(pl.PlayerInd)
	gameplayReportMonitor(int(ind), t)
	nox_xxx_netMarkMinimapObject_417190(int(ind), asObjectC(t), 1)
	nox_xxx_unitSetOwner_4EC290(asObjectC(u), asObjectC(t))
}
func stateRemoveMonitors(u, t *server.Object) {
	ud := u.UpdateData
	if u == nil || t == nil {
		return
	}
	t.ObjSubClass &^= 0x80
	pl := layout.As[server.PlayerUpdateData](ud).Player
	ind := int32(pl.PlayerInd)
	gameplayReportUnmonitor(int(ind), t)
	nox_xxx_netUnmarkMinimapObj_417300(int(ind), asObjectC(t), 1)
	nox_xxx_unitClearOwner_4EC300(asObjectC(t))
}
func stateOwns(u *server.Object, off uintptr, name string) bool {
	typ := stateType(off, name)
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) == typ {
			return true
		}
	}
	return false
}
func stateCount(u *server.Object, class, sub uint32) int32 {
	if u == nil || class == 0 || sub == 0 {
		return 0
	}
	var n int32
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.ObjClass)&class != 0 && uint32(it.ObjSubClass)&sub != 0 {
			n++
		}
	}
	return n
}
func stateEqual(u, t *server.Object) bool {
	if u == nil || t == nil || u.TypeInd != t.TypeInd {
		return false
	}
	if u.ObjClass&0x13001000 != 0 {
		if [4]*server.ModifierEff(itemModifiers(u)) != [4]*server.ModifierEff(itemModifiers(t)) {
			return false
		}
	}
	if u.ObjClass&0x100 == 0 {
		return true
	}
	if u.ObjSubClass&1 != 0 || u.ObjSubClass&2 == 0 {
		return bookRewardIndex(u) == bookRewardIndex(t)
	}
	return u.UseDataFieldGuide().Creature() == t.UseDataFieldGuide().Creature()
}
func statePostCreate(u *server.Object) {
	u.Field35 = 0
	u.Field36 = 0
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if pl.PlayerUnit != nil && nox_xxx_unitIsHostileMimic_4E7F90(asObjectC(pl.PlayerUnit), asObjectC(u)) == 1 {
			mask := uint32(1) << uint(pl.Index())
			u.Field35 |= mask
			u.Field36 |= mask
		}
	}
}
func statePlayerVisibility(ind int32) {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(ind))
	if pl == nil {
		return
	}
	mask := uint32(1) << uint32(ind&31)
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		u.Field35 &^= mask
		u.Field36 &^= mask
		if u.ObjClass&6 != 0 && pl.PlayerUnit != nil {
			enemy := nox_xxx_unitIsHostileMimic_4E7F90(asObjectC(pl.PlayerUnit), asObjectC(u)) == 1
			if enemy {
				if u.Field36&mask == 0 {
					u.Field36 |= mask
					u.Field35 |= mask
				}
			} else if u.Field36&mask != 0 {
				u.Field36 &^= mask
				u.Field35 |= mask
			}
		}
	}
}
func stateResetPixie(u *server.Object) int32 {
	typ := stateType(1567728, "Pixie")
	if u != nil && uint32(u.TypeInd) == typ {
		u.UpdateDataPixie().Target = nil
		return int32(memmap.PtrWord(u.UpdateData))
	}
	return int32(typ)
}
