package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func questDeathPenalty(u *server.Object) {
	data := u.UpdateDataPlayer()
	Nox_xxx_playerSubGold_4FA5D0(u, int(data.Player.GoldVal>>1))
	questLoseGems(u)
	questLoseWeapon(u)
	questLoseArmor(u)
	if data.Player.PlayerClass() == 0 {
		questLoseArmor(u)
	}
	questLoseSpell(u)
	questLoseSpell(u)
	questLoseBeastScroll(u)
	questLoseBeastScroll(u)
	questLoseWarriorAbility(u)
}
func questLoseWeapon(u *server.Object) {
	data := u.UpdateDataPlayer()
	var selected *server.Object
	for t := u.FirstItem(); t != nil; t = t.NextItem() {
		if t.Flags()&0x100 != 0 && t.Class()&0x1001000 != 0 && t.SubClass()&2 == 0 {
			selected = t
			break
		}
	}
	if selected == nil || selected.SubClass()&0x10000 != 0 {
		return
	}
	if selected.SubClass()&0x104 != 0 {
		empty := true
		for _, p := range unsafe.Slice((*unsafe.Pointer)(selected.InitData), 4) {
			if p != nil {
				empty = false
			}
		}
		if empty {
			return
		}
	}
	found := false
	for t := u.FirstItem(); t != nil; t = t.NextItem() {
		if t.Class()&0x1001000 != 0 && t.Flags()&0x100 == 0 {
			// Keep visiting candidates after a match: the shared eligibility callback
			// is observable and can inspect current player state.
			if Nox_xxx_playerClassCanUseItem_57B3D0(t, data.Player.PlayerClass()) {
				found = true
			}
		}
	}
	if found {
		GetServer().DelayedDelete(selected)
	}
}
func questLoseArmor(u *server.Object) {
	core := GetServer().S()
	eligible := func(t *server.Object) bool {
		return t.Flags()&0x100 != 0 && t.Class()&0x2000000 != 0 && core.Armor.Sub_415D10(int(t.TypeInd))&0x405 == 0
	}
	count := 0
	for t := u.FirstItem(); t != nil; t = t.NextItem() {
		if eligible(t) {
			count++
		}
	}
	if count == 0 {
		return
	}
	selected := core.Rand.Logic.IntClamp(0, count-1)
	index := 0
	for t := u.FirstItem(); t != nil; t = t.NextItem() {
		if eligible(t) {
			if index == selected {
				GetServer().DelayedDelete(t)
				return
			}
			index++
		}
	}
}
func questKnowledgePacket(p *server.Player, opcode uint16, id int) int32 {
	word, free := alloc.New(uint32(0))
	defer free()
	*word = uint32(opcode) | uint32(id)<<16
	return int32(reliableEnqueue(int(p.PlayerInd), unsafe.Slice((*byte)(unsafe.Pointer(word)), 4), nil, 1, 0))
}
func questLoseSpell(u *server.Object) {
	p := u.UpdateDataPlayer().Player
	class := p.PlayerClass()
	if class != 1 && class != 2 {
		return
	}
	count := 0
	for i, v := range p.SpellLvl {
		if v != 0 && questEligibilityPenaltySpell(uint32(i)) {
			count++
		}
	}
	selected := GetServer().S().Rand.Logic.IntClamp(0, count-1)
	index := 0
	for i, v := range p.SpellLvl {
		if v != 0 && questEligibilityPenaltySpell(uint32(i)) {
			if index == selected {
				p.SpellLvl[i] = 0
				questKnowledgePacket(p, 0x11f0, i)
				return
			}
			index++
		}
	}
}
func questLoseBeastScroll(u *server.Object) {
	p := u.UpdateDataPlayer().Player
	if p.PlayerClass() != 2 {
		return
	}
	count := 0
	for i, v := range p.BeastScrollLvl {
		if v == 1 && questEligibilityPenaltyBeast(uint32(i)) {
			count++
		}
	}
	selected := GetServer().S().Rand.Logic.IntClamp(0, count-1)
	index := 0
	for i, v := range p.BeastScrollLvl {
		if v == 1 && questEligibilityPenaltyBeast(uint32(i)) {
			if index == selected {
				p.BeastScrollLvl[i] = 0
				questKnowledgePacket(p, 0x13f0, i)
				return
			}
			index++
		}
	}
}
func questLoseWarriorAbility(u *server.Object) int8 {
	p := u.UpdateDataPlayer().Player
	result := int32(p.PlayerClass())
	if result != 0 {
		return int8(result)
	}
	count := 0
	for i := 0; i < 6; i++ {
		if p.SpellLvl[i] != 0 && questEligibilityAbility(uint32(i)) {
			count++
		}
	}
	result = int32(GetServer().S().Rand.Logic.IntClamp(1, count))
	selected, index := int(result), 1
	for i := 1; i < 6; i++ {
		if p.SpellLvl[i] != 0 {
			result = int32(bool2int(questEligibilityAbility(uint32(i))))
			if result != 0 {
				if index == selected {
					p.SpellLvl[i] = 0
					return int8(questKnowledgePacket(p, 0x12f0, i))
				}
				index++
			}
		}
	}
	return int8(result)
}
func questLoseGems(u *server.Object) {
	if dword_5d4594_2491676 == 0 {
		core := GetServer().S()
		dword_5d4594_2491676 = uint32(core.Types.IndByID("Diamond"))
		*memmap.PtrUint32(0x5D4594, 2491680) = uint32(core.Types.IndByID("Emerald"))
		*memmap.PtrUint32(0x5D4594, 2491684) = uint32(core.Types.IndByID("Ruby"))
	}
	kind := func(t *server.Object) int {
		id := uint32(t.TypeInd)
		if id == uint32(dword_5d4594_2491676) {
			return 0
		}
		if id == *memmap.PtrUint32(0x5D4594, 2491680) {
			return 1
		}
		if id == *memmap.PtrUint32(0x5D4594, 2491684) {
			return 2
		}
		return -1
	}
	var count [3]int
	for t := u.FirstItem(); t != nil; t = t.NextItem() {
		if k := kind(t); k >= 0 {
			count[k]++
		}
	}
	var odd [3]bool
	for i, n := range count {
		odd[i] = n&1 != 0
		count[i] = n / 2
	}
	for t := u.FirstItem(); t != nil; {
		next := t.NextItem()
		if k := kind(t); k >= 0 {
			if odd[k] {
				cost := shopPrice(1, nil, t)
				GetServer().DelayedDelete(t)
				Nox_xxx_playerAddGold_4FA590(u, int(cost/2))
				odd[k] = false
			} else if count[k] > 0 {
				GetServer().DelayedDelete(t)
				count[k]--
			}
		}
		t = next
	}
}
