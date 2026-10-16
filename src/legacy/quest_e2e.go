package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// Exports for the quest e2e test (e2e/quest.yaml).

// E2EQuestItemAllowed reports whether a quest character may keep the item when it's loaded (see questEligibilityItem).
func E2EQuestItemAllowed(u *server.Object) bool {
	return questEligibilityItem(u)
}

// E2EItemDurability returns the durability a new weapon or armor of the item's type has in the current game mode
// (scaled by QuestDurabilityMultiplier in quests, see creationScaleDurability), and its base durability (the one
// a map editor saves). ok is false for other items.
func E2EItemDurability(u *server.Object) (dur, base uint16, ok bool) {
	core := GetServer().S()
	var mod *server.Modifier
	switch {
	case u.ObjClass&0x2000000 != 0: // armor
		mod = core.Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd))
	case u.ObjClass&0x1000000 != 0: // weapon
		mod = core.Modif.Nox_xxx_getProjectileClassById413250(int(u.TypeInd))
	}
	if mod == nil || u.HealthData == nil {
		return 0, 0, false
	}
	base = uint16(mod.Durability52)
	dur = base
	if noxflags.HasGame(4096) {
		scale := float64(float32(core.Balance.Float("QuestDurabilityMultiplier")))
		dur = uint16(int32(float32(float64(dur) * scale)))
	}
	return dur, base, true
}
