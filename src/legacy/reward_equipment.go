package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// Modifier rows contain an ID byte, descriptor/name pointers, a tier mask and
// separate armor/weapon exclusion masks. The descriptor supplies allowed masks.
func rewardModifierMatches(off uintptr, tier, itemMask uint32, armor bool, pos uint32) bool {
	if tier&rewardWord(off+12) == 0 {
		return false
	}
	mod := (*server.ModifierEff)(rewardPtr(off + 4))
	allow, deny := mod.AllowWeapons28, rewardWord(off+20)
	if armor {
		allow, deny = mod.AllowArmor32, rewardWord(off+16)
	}
	return itemMask&allow != 0 && itemMask&deny == 0 && (pos == 0 || mod.AllowPos36&pos != 0)
}
func rewardModifierCount(base uintptr, tier, itemMask uint32, armor bool, pos uint32) int32 {
	var n int32
	for off := base; rewardWord(off+8) != 0; off += 24 {
		if rewardModifierMatches(off, tier, itemMask, armor, pos) {
			n++
		}
	}
	return n
}
func rewardModifierPick(base uintptr, tier, itemMask uint32, armor bool, pos uint32) (unsafe.Pointer, uintptr, bool) {
	n := rewardModifierCount(base, tier, itemMask, armor, pos)
	if n == 0 {
		return nil, 0, false
	}
	pick := rewardRoll(0, n-1)
	for off := base; rewardWord(off+8) != 0; off += 24 {
		if rewardModifierMatches(off, tier, itemMask, armor, pos) {
			if pick == 0 {
				return rewardPtr(off + 4), off, true
			}
			pick--
		}
	}
	return nil, 0, false
}
func rewardModifierFlags(tier uint32) uint32 {
	var n int32
	switch tier {
	case 2:
		n = rewardRoll(0, 1)
	case 4:
		n = rewardRoll(0, 2)
	case 8:
		n = rewardRoll(1, 3)
	case 16:
		n = rewardRoll(2, 4)
	default:
		return 0
	}
	switch n {
	case 0:
		return 0
	case 1:
		v := rewardRoll(1, 100)
		if v <= 20 {
			return 4
		}
		if v > 50 {
			return 2
		}
		return 1
	case 2:
		v := rewardRoll(1, 100)
		if v <= 12 {
			return 5
		}
		if v > 25 {
			return 3
		}
		return 6
	case 3:
		return 7
	default:
		return 15
	}
}
func rewardEnchantTier(stage uint32) uint32 {
	v := int32(stage >> 1)
	if v < 1 {
		v = 1
	}
	if v >= 5 {
		v = 4
	}
	lo, hi := v-1, v+1
	if lo < 1 {
		lo = 1
	}
	if hi >= 5 {
		hi = 4
	}
	return uint32(rewardRoll(lo, hi))
}
func rewardEquipment(stage uint32, armor bool) *server.Object {
	category := uint32(1)
	if armor {
		category = 2
	}
	tier, off, ok := rewardEquipmentRow(stage, category)
	if !ok {
		return nil
	}
	typ := rewardWord(off + 8)
	if typ == 0 {
		return nil
	}
	core := GetServer().S()
	var itemMask uint32
	if armor {
		itemMask = core.Armor.Sub_415D10(int(typ))
	} else {
		itemMask = core.Weapons.Nox_xxx_ammoCheck_415880(int(typ))
	}
	it := core.NewObjectByTypeInd(int(typ))
	if it == nil {
		return nil
	}
	var attrs server.ModifierInitData
	if !armor && it.ObjClass&0x1000 != 0 && it.ObjSubClass&0x47f0000 != 0 {
		if *(*byte)(rewardPtr(off + 4)) == '#' {
			id := core.Modif.Nox_xxx_modifGetIdByName413290("Replenishment1")
			attrs.Modifiers[2] = core.Modif.Nox_xxx_modifGetDescById413330(id)
			stateAttributes(it, &attrs)
		}
		return it
	}
	flags := rewardModifierFlags(tier)
	if flags == 0 {
		return it
	}
	base := uintptr(210704)
	if armor {
		base = 210848
	}
	if flags&1 != 0 && rewardModifierCount(base, tier, itemMask, armor, 0) == 0 {
		if flags&2 == 0 {
			flags |= 2
		} else if flags&4 == 0 {
			flags |= 4
		} else if flags&8 == 0 {
			flags |= 8
		}
		flags &^= 1
	}
	if flags&2 != 0 && rewardModifierCount(210992, tier, itemMask, armor, 0) == 0 {
		if flags&4 == 0 {
			flags |= 4
		} else if flags&8 == 0 {
			flags |= 8
		}
		flags &^= 2
	}
	if flags&4 != 0 && rewardModifierCount(209336, tier, itemMask, armor, 1) == 0 {
		flags &^= 12
	}
	if flags == 0 {
		return it
	}
	if flags&1 != 0 {
		ptr, _, _ := rewardModifierPick(base, tier, itemMask, armor, 0)
		attrs.Modifiers[0] = (*server.ModifierEff)(ptr)
	}
	if flags&2 != 0 {
		ptr, _, _ := rewardModifierPick(210992, tier, itemMask, armor, 0)
		attrs.Modifiers[1] = (*server.ModifierEff)(ptr)
	}
	first := uintptr(209336)
	if flags&4 != 0 {
		modTier := tier
		if !armor {
			modTier = rewardEnchantTier(stage)
		}
		ptr, off, ok := rewardModifierPick(209336, modTier, itemMask, armor, 1)
		if ok {
			attrs.Modifiers[2] = (*server.ModifierEff)(ptr)
			first = off
		}
	}
	if flags&8 != 0 {
		modTier := rewardEnchantTier(stage)
		ptr, off, ok := rewardModifierPick(209336, modTier, itemMask, armor, 2)
		// ID equality, not pointer equality, suppresses the second enchantment.
		if ok && rewardByte(first) != rewardByte(off) {
			attrs.Modifiers[3] = (*server.ModifierEff)(ptr)
		}
	}
	stateAttributes(it, &attrs)
	return it
}
