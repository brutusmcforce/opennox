package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func inventoryPickup(u, it *server.Object, arg int) int {
	return bool2int(Nox_xxx_pickupDefault_4F31E0(u, it, arg, 0))
}
func inventoryFoodPickup(u, it *server.Object, arg int) int {
	if u == nil || it == nil {
		return 0
	}
	if sub_419E60(asObjectC(u)) == 0 && it.ObjSubClass&0x84 == 0 {
		it.Use.CallDiscard(u, it)
	}
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		inventoryFoodSound(215640, u, it)
	}
	return rv
}
func inventoryUsePickup(u, it *server.Object, arg int) int {
	effectsUse(u, it)
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	return inventoryPickup(u, it, arg)
}
func inventoryTrapPickup(u, it *server.Object, arg int) int {
	if nox_xxx_unitHasThatParent_4EC4F0(asObjectC(it), asObjectC(u)) != 0 {
		rv := inventoryPickup(u, it, arg)
		if rv != 0 {
			inventorySound(824, u, 0, 0)
		}
		return rv
	}
	if u.ObjClass&4 != 0 {
		inventorySound(925, u, 2, int(u.NetCode))
	}
	return 0
}
func inventoryBookPickup(u, it *server.Object, arg int, ability bool) int {
	if noxflags.HasGame(6144) {
		effectsUse(u, it)
	}
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		sound := 826
		if !ability && it.ObjSubClass&1 == 0 {
			sound = 828
		}
		inventorySound(sound, u, 0, 0)
	}
	return rv
}
func inventoryAnkhPickup(u, it *server.Object) int {
	if u.ObjClass&4 == 0 {
		return 0
	}
	*(*uint32)(unsafe.Add(u.UpdateData, 320))++
	GetServer().DelayedDelete(it)
	inventorySound(1004, u, 0, 0)
	return 1
}
func inventoryCrownPickup(u, it *server.Object, arg int) int {
	data := it.UpdateData
	if u.ObjClass&4 == 0 {
		return 0
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		*(*uint32)(unsafe.Add(u.UpdateData, 264)) = GetServer().S().Frame()
		GetServer().S().ObjSetOwner(u, it)
		spellLifeApplyBuff(u, 30, 0, 5)
		inventorySound(313, u, 0, 0)
		inventoryMessage(10, u, uint32(u.TeamVal.ID))
		playerStateUnmark(it, 1)
	}
	*(*uint32)(unsafe.Add(data, 4)) = 0
	return rv
}
func inventoryTreasurePickup(u, it *server.Object, arg int) int {
	if inventoryPickup(u, it, arg) == 0 {
		return 0
	}
	if u.ObjClass&4 == 0 || !noxflags.HasGame(64) {
		return 1
	}
	pl := u.UpdateDataPlayer().Player
	inventorySound(307, u, 0, 0)
	pl.Field2152++
	pl.Field2156 = uint32(sessionScavengerMaximum())
	gameplayReportScavenger(u)
	core := GetServer().S()
	if !u.TeamVal.Has() {
		if pl.Field2152 == uint32(sessionScavengerMaximum()) {
			noxflags.SetGame(8)
			gameplayReportChangeScore(u, 1)
			gameplayReportLesson(u)
			for other := core.Players.FirstUnit(); other != nil; other = core.Players.NextUnit(other) {
				if other != u {
					gameplayReportEliminationDeath(other)
					gameplayReportLesson(other)
				}
			}
		}
		return 1
	}
	if team := core.Teams.ByID(u.TeamVal.ID); team != nil {
		var total uint32
		for other := core.Players.FirstUnit(); other != nil; other = core.Players.NextUnit(other) {
			if Nox_xxx_teamCompare2_419180(&other.TeamVal, team.ID()) != 0 {
				total += other.UpdateDataPlayer().Player.Field2152
			}
		}
		if total == uint32(sessionScavengerMaximum()) {
			noxflags.SetGame(8)
		}
	}
	return 1
}
