package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func inventoryCrownDrop(u, it *server.Object, pos *types.Pointf) int {
	core := GetServer().S()
	if noxflags.HasGame(16) && noxflags.GetGamePlay()&4 != 0 && it.TeamVal.ID != 0 {
		time := core.Frame()
		chosen := it
		for other := core.Players.FirstUnit(); other != nil; other = core.Players.NextUnit(other) {
			since := other.UpdateDataPlayer().Field66
			if Nox_xxx_teamCompare2_419180(&other.TeamVal, it.TeamVal.ID) != 0 && since < time {
				time = since
				chosen = other
			}
		}
		if chosen != nil && chosen != u {
			it.UpdateDataCrown().SetPendingObj(chosen) // the crown goes to this player next
		}
	}
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	core.ObjClearOwner(it)
	Nox_xxx_spellBuffOff_4FF5B0(u, 30)
	inventoryMessage(11, u, uint32(u.TeamVal.ID))
	playerStateMark(it, 1)
	return 1
}
func inventoryTreasureDrop(u, it *server.Object, pos *types.Pointf) int {
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	if u.ObjClass&4 != 0 && noxflags.HasGame(64) {
		pl := u.UpdateDataPlayer().Player
		pl.Field2152--
		pl.Field2156 = uint32(sessionScavengerMaximum())
		gameplayReportScavenger(u)
		inventorySound(308, u, 0, 0)
	}
	return 1
}
