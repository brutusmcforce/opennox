package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// bookRewardIndex returns the spell (SpellRewardUse) or the ability (AbilityRewardUse) taught by a book.
func bookRewardIndex(u *server.Object) byte {
	if u.ObjSubClass&1 != 0 {
		return u.UseDataSpellReward().Spell
	}
	return u.UseDataAbilityReward().Ability
}

func bookUseGuide(u, item *server.Object) int {
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	p := u.UpdateDataPlayer().Player
	id := bookGuideID(item.UseDataFieldGuide().Creature())
	if noxflags.HasGame(4096) && p.PlayerClass() != 2 {
		bookAwardPrivate(u, "pickup.c:ObjectEquipClassFail")
		return 0
	}
	if p.BeastScrollLvl[id] != 0 {
		bookAwardPrivate(u, "objcoll.c:AlreadyHaveGuide")
		return 0
	}
	bookAwardGuide(u, id, 1)
	GetServer().DelayedDelete(item)
	return 1
}
func bookUseSpell(u, item *server.Object) int {
	id := int32(item.UseDataSpellReward().Spell)
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	p := u.UpdateDataPlayer().Player
	class := int32(p.PlayerClass())
	if class != 1 && class != 2 || playerSpellClassCheck(class, id) != 0 {
		bookAwardPrivate(u, "use.c:SpellRewardClassFail")
		GetServer().S().Audio.EventObj(925, u, 2, u.NetCode)
		return 0
	}
	auto := int32(0)
	if noxflags.HasGame(6144) && p.SpellLvl[id] == 0 {
		auto = 1
	}
	if bookAwardSpell(u, id, 1, auto, 0) != 0 {
		GetServer().DelayedDelete(item)
	} else {
		GetServer().S().Audio.EventObj(925, u, 2, u.NetCode)
	}
	return 1
}
func bookUseAbility(u, item *server.Object) int {
	id := int32(item.UseDataAbilityReward().Ability)
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	p := u.UpdateDataPlayer().Player
	if p.PlayerClass() != 0 {
		bookAwardPrivate(u, "pickup.c:ObjectEquipClassFail")
		GetServer().S().Audio.EventObj(925, u, 2, u.NetCode)
		return 0
	}
	auto := int32(0)
	if noxflags.HasGame(6144) && p.SpellLvl[id] == 0 {
		auto = 1
	}
	if bookAwardAbility(u, id, auto) != 0 {
		GetServer().DelayedDelete(item)
	} else {
		GetServer().S().Audio.EventObj(925, u, 2, u.NetCode)
	}
	return 1
}
