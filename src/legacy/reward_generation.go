package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func rewardRoll(lo, hi int32) int32 {
	return int32(GetServer().S().Rand.Logic.IntClamp(int(lo), int(hi)))
}
func rewardWord(off uintptr) uint32        { return memmap.Uint32(0x587000, off) }
func rewardPtr(off uintptr) unsafe.Pointer { return memmap.LoadPtr(0x587000, off) }
func rewardByte(off uintptr) byte          { return memmap.Uint8(0x587000, off) }
func rewardTier(stage uint32) uint32 {
	if stage >= 10 {
		return 16
	}
	if stage == 0 {
		return 1
	}
	var weights [5]float32
	switch stage {
	case 1:
		weights[0], weights[1] = 87.5, 12.5
	case 9:
		weights[3], weights[4] = 12.5, 87.5
	default:
		i := stage >> 1
		if stage&1 != 0 {
			// The original overwrites its 75.0 assignment with 12.5 in this same slot.
			weights[i], weights[i+1] = 12.5, 12.5
		} else {
			weights[i-1], weights[i] = 50, 50
		}
	}
	roll := rewardRoll(0, 200)
	var total float64
	for i, w := range weights {
		total += float64(w)
		if float64(float32(float64(roll)*.5)) <= total {
			return 1 << uint(i)
		}
	}
	return 1
}

// rewardWeighted preserves table order and the single draw after a nonzero sum.
// All predicates read definition data; none mutate it or consume randomness.
func rewardWeighted(base, stride uintptr, eligible func(uintptr) bool) (uintptr, bool) {
	var total int32
	for off := base; rewardWord(off+4) != 0; off += stride {
		if eligible(off) {
			total += int32(rewardByte(off))
		}
	}
	if total == 0 {
		return 0, false
	}
	roll := rewardRoll(0, total-1)
	total = 0
	for off := base; rewardWord(off+4) != 0; off += stride {
		if eligible(off) {
			total += int32(rewardByte(off))
			if roll < total {
				return off, true
			}
		}
	}
	return 0, false
}
func rewardExplicit(values []byte) int32 {
	var n int32
	for _, v := range values {
		if v == 1 {
			n++
		}
	}
	if n == 0 {
		return 0
	}
	pick := rewardRoll(0, n-1)
	for i, v := range values {
		if v == 1 {
			if pick == 0 {
				return int32(i)
			}
			pick--
		}
	}
	return 0
}
func rewardBook(u *server.Object, stage uint32, kind int) *server.Object {
	var id int32
	d := u.InitDataRewardMarker()
	switch kind {
	case 1:
		if byte(d.Explicit)&2 != 0 {
			id = rewardExplicit(d.Abilities[:])
		} else {
			id = rewardRoll(1, 5)
		}
	default:
		flag, list, base := byte(1), d.Spells[:], uintptr(207104)
		if kind == 2 {
			flag, list, base = 4, d.Guides[:], 207792
		}
		if byte(d.Explicit)&flag != 0 {
			id = rewardExplicit(list)
		} else {
			tier := rewardTier(stage)
			off, ok := rewardWeighted(base, 12, func(off uintptr) bool { return rewardWord(off+8)&tier != 0 })
			if !ok {
				return nil
			}
			id = int32(rewardWord(off + 4))
		}
	}
	if id == 0 {
		return nil
	}
	name := "AbilityBook"
	if kind == 0 {
		wizard := playerSpellClassCheck(1, int32(id)) == 0
		conjurer := playerSpellClassCheck(2, int32(id)) == 0
		switch {
		case wizard && conjurer:
			name = "CommonSpellBook"
		case wizard:
			name = "WizardSpellBook"
		case conjurer:
			name = "ConjurerSpellBook"
		default:
			return nil
		}
	} else if kind == 2 {
		name = "FieldGuide"
	}
	it := GetServer().S().NewObjectByTypeID(name)
	if it == nil {
		return nil
	}
	if kind == 2 {
		name := alloc.GoString((*byte)(rewardPtr(uintptr(70500 + 4*id))))
		data := it.UseDataFieldGuide().CreatureBuf[:len(name)+1]
		copy(data, name)
		data[len(name)] = 0
	} else if kind == 0 {
		it.UseDataSpellReward().Spell = byte(id)
	} else {
		it.UseDataAbilityReward().Ability = byte(id)
	}
	return it
}
func rewardEquipmentRow(stage, category uint32) (uint32, uintptr, bool) {
	tier := rewardTier(stage)
	off, ok := rewardWeighted(208176, 20, func(off uintptr) bool {
		return uint32(rewardByte(off+12))&category != 0 && rewardWord(off+16)&tier != 0 && GetServer().S().Types.ByInd(int(rewardWord(off+8))).Allowed()
	})
	return tier, off, ok
}
func rewardPotion(stage uint32) *server.Object {
	_, off, ok := rewardEquipmentRow(stage, 4)
	if !ok || rewardWord(off+8) == 0 {
		return nil
	}
	return GetServer().S().NewObjectByTypeInd(int(rewardWord(off + 8)))
}
func rewardGem(stage uint32) *server.Object {
	tier := rewardTier(stage)
	if tier < 4 || rewardRoll(1, 100) <= 90 {
		name := "QuestGoldPile"
		if rewardRoll(1, 2) == 1 {
			name = "QuestGoldChest"
		}
		it := GetServer().S().NewObjectByTypeID(name)
		if it == nil {
			return nil
		}
		off := uintptr(211136)
		switch tier {
		case 2:
			off += 8
		case 4:
			off += 16
		case 8:
			off += 24
		case 16:
			off += 32
		}
		*equipmentWord(it.InitData, 0) = uint32(rewardRoll(int32(rewardWord(off)), int32(rewardWord(off+4))))
		return it
	}
	n := rewardRoll(1, 100)
	name := "RubyGem"
	if n >= 90 {
		name = "DiamondGem"
	} else if n >= 50 {
		name = "EmeraldGem"
	}
	return GetServer().S().NewObjectByTypeID(name)
}
func rewardMarker(u *server.Object, stage uint32) *server.Object {
	if uint32(u.TypeInd) == stateType(1568276, "RewardMarkerPlus") {
		stage += 2
	}
	chance := u.InitDataRewardMarker().Chance
	thresholds := [4]int32{75, 50, 25, 5}
	if chance >= 1 && chance <= 4 && rewardRoll(0, 100) > thresholds[chance-1] {
		return nil
	}
	mask := u.InitDataRewardMarker().Categories
	var total int32
	for i := uint(0); i < 8; i++ {
		if mask&(1<<i) != 0 {
			total += int32(rewardByte(207044 + uintptr(i*8)))
		}
	}
	if total == 0 {
		return nil
	}
	roll := uint32(rewardRoll(1, total))
	var sum, category uint32
	category = stage
	for i := uint(0); i < 8; i++ {
		if mask&(1<<i) != 0 {
			sum += uint32(rewardByte(207044 + uintptr(i*8)))
			if sum >= roll {
				category = 1 << i
				break
			}
		}
	}
	switch category {
	case 1:
		return rewardBook(u, stage, 0)
	case 2:
		return rewardBook(u, stage, 1)
	case 4:
		return rewardBook(u, stage, 2)
	case 8:
		return rewardEquipment(stage, false)
	case 16:
		return rewardEquipment(stage, true)
	case 64:
		return rewardPotion(stage)
	default:
		return rewardGem(stage)
	}
}
