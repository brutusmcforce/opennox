package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/log"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// QuestAllowAnyItems keeps items in the saved quest inventories that can't be found in quests (e.g. placed on
// custom quest maps); by default they are dropped when the character is loaded.
var QuestAllowAnyItems = false

var questItemLog = log.New("quest")

// Eligibility shares the production reward tables, but uses their enabled
// fields rather than reward tier/chance filtering.
func questEligibilityListed(base uintptr, id uint32) bool {
	for off := base; rewardWord(off) != 0; off += 12 {
		if rewardWord(off) == id && rewardWord(off+4) != 0 {
			return true
		}
	}
	return false
}
func questEligibilityPenaltySpell(id uint32) bool {
	return questEligibilityListed(207108, id) && id != 0 && id != 34 && id != 27 && id != 9 && id != 41
}
func questEligibilityPenaltyBeast(id uint32) bool {
	return questEligibilityListed(207796, id) && id != 0
}
func questEligibilityAbility(id uint32) bool { return id > 0 && id < 6 }
func questEligibilitySpell(id uint32) bool {
	return questEligibilityListed(207108, id) || (id >= 46 && id <= 49) || (id >= 122 && id <= 125) || (id >= 75 && id <= 114)
}
func questEligibilityBeast(id uint32) bool {
	allowed := questEligibilityListed(207796, id)
	for off := uintptr(207032); rewardPtr(off) != nil; off += 4 {
		p := rewardPtr(off)
		if *(*uint32)(p) == 0 {
			continue
		}
		// The first word is a marker. Unlike the ordinary table, a terminating
		// zero in the following IDs participates in comparison.
		for pos := 4; ; pos += 4 {
			value := *(*uint32)(unsafe.Add(p, pos))
			if value == id {
				allowed = true
				break
			}
			if value == 0 {
				break
			}
		}
	}
	return allowed
}
func questEligibilityBook(u *server.Object) bool {
	sub := uint32(u.ObjSubClass)
	if sub&1 != 0 {
		return questEligibilityListed(207108, uint32(u.UseDataSpellReward().Spell))
	}
	if sub&2 != 0 {
		return questEligibilityListed(207796, uint32(bookGuideID(u.UseDataFieldGuide().Creature())))
	}
	if sub&4 != 0 {
		return questEligibilityAbility(uint32(u.UseDataAbilityReward().Ability))
	}
	return false
}
func questEligibilityItem(u *server.Object) bool {
	types := unsafe.Slice(memmap.PtrUint32(0x5d4594, 1568328), 7)
	if types[0] == 0 {
		for i, name := range []string{"Diamond", "Emerald", "Ruby", "SulphorousFlareWand", "StreetSneakers", "StreetPants", "StreetShirt"} {
			types[i] = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	class := uint32(u.ObjClass)
	if class&0x40 != 0 {
		return false
	}
	if class&0x10 != 0 {
		return uint32(u.ObjSubClass)&0x1ff78 != 0
	}
	if class&0x100 != 0 {
		return questEligibilityBook(u)
	}
	typ := uint32(u.TypeInd)
	if typ == types[0] || typ == types[1] || typ == types[2] {
		return true
	}
	listed := false
	for off := uintptr(208180); rewardWord(off) != 0; off += 20 {
		if rewardWord(off+4) == typ {
			listed = true
			break
		}
	}
	if typ == types[3] || typ == types[4] || typ == types[5] || typ == types[6] {
		return questEligibilitySpecial(u)
	}
	if !listed {
		return false
	}
	if class&0x3000000 != 0 {
		return questEligibilityModifiers(u)
	}
	return true
}
func questEligibilityMods(u *server.Object) []*server.ModifierEff {
	return itemModifiers(u)
}
func questEligibilityMask(u *server.Object) (uint32, bool) {
	core := GetServer().S()
	if uint32(u.ObjClass)&0x1000000 != 0 {
		return core.Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(u), true
	}
	return core.Armor.Nox_xxx_unitArmorInventoryEquipFlags_415C70(u), false
}
func questEligibilityModRow(base uintptr, m *server.ModifierEff) uintptr {
	for off := base; rewardWord(off+8) != 0; off += 24 {
		if rewardPtr(off+4) == m.C() {
			return off
		}
	}
	return 0
}
func questEligibilityModMask(off uintptr, m *server.ModifierEff, mask uint32, weapon bool) bool {
	if weapon {
		return mask&m.AllowWeapons28 != 0 && mask&rewardWord(off+20) == 0
	}
	return mask&m.AllowArmor32 != 0 && mask&rewardWord(off+16) == 0
}
func questEligibilityMaterial(u *server.Object) bool {
	m := questEligibilityMods(u)[0]
	if m == nil {
		return true
	}
	mask, weapon := questEligibilityMask(u)
	base := uintptr(210848)
	if weapon {
		base = 210704
	}
	off := questEligibilityModRow(base, m)
	return off != 0 && questEligibilityModMask(off, m, mask, weapon)
}
func questEligibilityQuality(u *server.Object) bool {
	m := questEligibilityMods(u)[1]
	if m == nil {
		return true
	}
	mask, weapon := questEligibilityMask(u)
	off := questEligibilityModRow(210992, m)
	return off != 0 && questEligibilityModMask(off, m, mask, weapon)
}
func questEligibilityEffects(u *server.Object) bool {
	mods := questEligibilityMods(u)
	if questEligibilityReplenish[0] == nil {
		core := GetServer().S()
		for i, name := range []string{"Replenishment1", "Replenishment2", "Replenishment3", "Replenishment4"} {
			m := core.Modif.Nox_xxx_modifGetDescById413330(core.Modif.Nox_xxx_modifGetIdByName413290(name))
			questEligibilityReplenish[i] = m.C()
		}
	}
	special := questEligibilityReplenish
	mask, weapon := questEligibilityMask(u)
	for slot := 2; slot < 4; slot++ {
		m := mods[slot]
		if m == nil {
			continue
		}
		if m.AllowPos36&uint32(1<<(slot-2)) == 0 {
			return false
		}
		off := questEligibilityModRow(209336, m)
		replenishment := false
		for _, p := range special {
			if m.C() == p {
				replenishment = true
				break
			}
		}
		if replenishment {
			allowed := false
			for row := uintptr(208180); rewardPtr(row) != nil; row += 20 {
				if *(*byte)(rewardPtr(row)) == '#' && rewardWord(row+4) == uint32(u.TypeInd) {
					allowed = true
				}
			}
			if !allowed {
				return false
			}
		} else if off == 0 || !questEligibilityModMask(off, m, mask, weapon) {
			return false
		}
	}
	return true
}
func questEligibilityModifiers(u *server.Object) bool {
	return questEligibilityMaterial(u) && questEligibilityQuality(u) && questEligibilityEffects(u)
}
func questEligibilitySpecial(u *server.Object) bool {
	cached := &questEligibilitySpecialMod
	if *cached == nil {
		modif := &GetServer().S().Modif
		*cached = modif.Nox_xxx_modifGetDescById413330(modif.Nox_xxx_modifGetIdByName413290("Replenishment1")).C()
	}
	class := uint32(u.ObjClass)
	core := GetServer().S()
	if class&0x1000000 != 0 && core.Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(u)&0x10000 != 0 {
		mods := questEligibilityMods(u)
		if mods[0] != nil || mods[1] != nil || mods[2].C() != *cached || mods[3] != nil {
			return false
		}
	}
	if class&0x2000000 != 0 && core.Armor.Nox_xxx_unitArmorInventoryEquipFlags_415C70(u)&0x405 != 0 {
		for _, m := range questEligibilityMods(u) {
			if m == nil {
				continue
			}
			name := m.Name()
			if len(name) < 8 {
				return false
			}
			for i, want := range []byte("usercolo") {
				c := name[i]
				if c >= 'A' && c <= 'Z' {
					c += 'a' - 'A'
				}
				if c != want {
					return false
				}
			}
		}
	}
	return true
}
func questEligibilityInventoryLimit(u *server.Object) bool {
	types := unsafe.Slice(memmap.PtrUint32(0x5d4594, 1568356), 13)
	if types[0] == 0 {
		for i, name := range []string{"RedPotion", "BluePotion", "CurePoisonPotion", "HastePotion", "InvisibilityPotion", "ShieldPotion", "VampirismPotion", "FireProtectPotion", "ShockProtectPotion", "PoisonProtectPotion", "InvulnerabilityPotion", "InfravisionPotion", "InfinitePainWand"} {
			types[i] = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	if u == nil || uint32(u.ObjClass)&4 == 0 {
		return true
	}
	for _, typ := range types[:12] {
		if equipmentCount(u, int(typ)) > 9 {
			return false
		}
	}
	limit := floatToInt32(float32(GetServer().S().Balance.Float("ForceOfNatureStaffLimit")))
	return equipmentCount(u, int(types[12])) <= int(limit)
}
