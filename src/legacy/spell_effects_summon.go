package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func spellEffectSummonPosition(record, output unsafe.Pointer) int32 {
	if record == nil {
		return 0
	}
	u := layout.As[server.DurSpell](record).Caster16
	if u == nil || output == nil {
		return 0
	}
	pos := layout.As[server.DurSpell](record).Pos2
	if u.ObjClass&4 != 0 {
		dx := float64(pos.X) - float64(u.PosVec.X)
		dy := float64(pos.Y) - float64(u.PosVec.Y)
		fy := float32(dy)
		distance := math.Sqrt(dy*float64(fy) + dx*dx)
		if distance > 50 {
			length := float32(distance)
			pos.X = float32(dx*50/float64(length) + float64(u.PosVec.X))
			pos.Y = float32(float64(fy)*50/float64(length) + float64(u.PosVec.Y))
		}
		if spellEffectTrace(u.PosVec, pos, 9) && !worldTileWater(pos) {
			*(*types.Pointf)(output) = pos
			return 1
		}
		spellEffectInform(u)
		return 0
	}
	if !spellEffectTrace(u.PosVec, pos, 9) {
		pos = u.PosVec
	}
	*(*types.Pointf)(output) = pos
	return 1
}
func spellEffectSummonStart(record unsafe.Pointer) int32 {
	u := layout.As[server.DurSpell](record).Caster16
	guide := int32(layout.As[server.DurSpell](record).Spell) - 74
	if u == nil || u.ObjFlags&0x8020 != 0 || layout.As[server.DurSpell](record).Flag20 != 0 {
		return 1
	}
	if u.ObjClass&4 != 0 {
		if controlFlags(4608) && controlPlayer(u).BeastScrollLvl[int(guide)] == 0 {
			resourcePriority(u, "Summon.c:NeedGuideToSummon")
			return 1
		}
		if !Nox_xxx_checkSummonedCreaturesLimit_500D70(u, int(guide)) {
			resourcePriority(u, "Summon.c:CreatureControlFailed")
			return 1
		}
	}
	var pos types.Pointf
	if spellEffectSummonPosition(record, unsafe.Pointer(&pos)) == 0 {
		return 1
	}
	typ := GetServer().S().Types.IndByID(alloc.GoString(bookGuideName(int32(guide))))
	seq := memmap.PtrUint16(0x5d4594, 1570276)
	n := *seq
	*seq++
	if *seq >= 65000 {
		*seq = 0
	}
	*(*uint16)(unsafe.Pointer(summonByte(record, 72))) = uint16(typ)
	*(*types.Pointf)(unsafe.Pointer(summonByte(record, 74))) = pos
	*summonByte(record, 82) = byte(u.Direction1)
	*(*uint16)(unsafe.Pointer(summonByte(record, 83))) = n
	*summonByte(record, 85) = 0
	duration := int32(memmap.PtrWord(record))
	switch bookGuideSize(int32(guide)) {
	case 1:
		duration = floatToInt32(float32(spellEffectTable("SummonDuration", 0)))
	case 2:
		duration = floatToInt32(float32(spellEffectTable("SummonDuration", 1)))
	case 4:
		duration = floatToInt32(float32(spellEffectTable("SummonDuration", 2)))
	}
	layout.As[server.DurSpell](record).Frame68 = uint32(duration) + GetServer().S().Frame()
	visibilityFXSummonStart(uint16(n), pos, byte(u.Direction1), uint16(typ), uint16(duration))
	return 0
}
func spellEffectSummonFinish(record unsafe.Pointer) int32 {
	owner := layout.As[server.DurSpell](record).Caster16
	if owner == nil || owner.ObjFlags&0x8020 != 0 {
		return 1
	}
	if layout.As[server.DurSpell](record).Frame68-1 != GetServer().S().Frame() {
		return 0
	}
	guide := int32(layout.As[server.DurSpell](record).Spell) - 74
	if owner.ObjClass&4 != 0 {
		if controlFlags(512) && controlPlayer(owner).BeastScrollLvl[int(guide)] == 0 {
			resourcePriority(owner, "Summon.c:NeedGuideToSummon")
			return 1
		}
		if !Nox_xxx_checkSummonedCreaturesLimit_500D70(owner, int(guide)) {
			resourcePriority(owner, "Summon.c:CreatureControlFailed")
			return 1
		}
	}
	u := Nox_xxx_unitDoSummonAt_5016C0(int(*(*uint16)(unsafe.Pointer(summonByte(record, 72)))), *(*types.Pointf)(unsafe.Pointer(summonByte(record, 74))), owner, server.Dir16(*summonByte(record, 82)))
	if u != nil {
		GetServer().S().Audio.EventObj(899, u, 0, 0)
	}
	*summonByte(record, 85) = 1
	return 1
}
func spellEffectSummonCancel(record unsafe.Pointer) {
	if *summonByte(record, 85) == 0 {
		visibilityFXSummonCancel(*(*uint16)(unsafe.Pointer(summonByte(record, 83))))
		GetServer().S().Audio.EventPos(900, *(*types.Pointf)(durScratch(record, 74)), 0, 0)
	}
}
func spellEffectCharmStart(record unsafe.Pointer) int32 {
	s := GetServer().S()
	source := layout.As[server.DurSpell](record).Caster16
	target := layout.As[server.DurSpell](record).Target48
	level := int32(layout.As[server.DurSpell](record).Level)
	if layout.As[server.DurSpell](record).Flag20 != 0 {
		spellLifeApplyBuff(target, 3, int16(floatToInt32(float32(spellEffectScalar("ConfuseEnchantDuration")))), int8(level))
		spellEffectAlert(source, target)
		return 1
	}
	pos := layout.As[server.DurSpell](record).Pos2
	target = s.Nox_xxx_spellFlySearchTarget(&pos, source, s.Spells.Flags(spell.ID(layout.As[server.DurSpell](record).Spell)), 300, 0, source)
	*(*unsafe.Pointer)(unsafe.Pointer(&layout.As[server.DurSpell](record).Target48)) = target.CObj()
	if target == nil {
		resourcePriority(source, "Summon.c:CharmNoCreatureCloseEnough")
		s.Audio.EventObj(16, source, 0, 0)
		return 1
	}
	if target.ObjClass&2 != 0 && !server.Nox_xxx_creatureIsMonitored_500CC0(source, target) {
		guide := int32(bookGuideCharmable(uint32(target.TypeInd)))
		if source.ObjClass&4 != 0 && !spellCharmAll {
			if guide == 0 {
				resourcePriority(source, "Summon.c:CreatureNotCharmable")
				layout.As[server.DurSpell](record).Target48 = nil
				s.Audio.EventObj(16, source, 0, 0)
				return 1
			}
			if controlPlayer(source).BeastScrollLvl[int(guide)] == 0 {
				layout.As[server.DurSpell](record).Target48 = nil
				resourcePriority(source, "Summon.c:NeedGuideToCharm")
				s.Audio.EventObj(16, source, 0, 0)
				return 1
			}
		}
		duration := int32(memmap.PtrWord(record))
		size := int32(bookGuideSize(int32(guide)))
		switch {
		case size <= 1:
			duration = floatToInt32(float32(spellEffectTable("CharmSmallDuration", level-1)))
		case size == 2:
			duration = floatToInt32(float32(spellEffectTable("CharmMediumDuration", level-1)))
		case size == 4:
			duration = floatToInt32(float32(spellEffectTable("CharmLargeDuration", level-1)))
		}
		layout.As[server.DurSpell](record).Frame68 = uint32(duration) + s.Frame()
		spellLifeApplyBuff(target, 28, int16(uint16(duration)+1), 5)
		spellLifeRayMessage((*server.DurSpell)(record))
		return 0
	}
	layout.As[server.DurSpell](record).Target48 = nil
	s.Audio.EventObj(16, source, 0, 0)
	return 1
}
func spellEffectCharmFinish(record unsafe.Pointer) int32 {
	s := GetServer().S()
	source := layout.As[server.DurSpell](record).Caster16
	u := layout.As[server.DurSpell](record).Target48
	fail := func() int32 { s.Audio.EventObj(16, source, 0, 0); return 1 }
	if u == nil || u.ObjFlags&0x8020 != 0 {
		return fail()
	}
	distance := stateDistance(source, u)
	if distance > 300 {
		if source.ObjClass&4 != 0 {
			resourcePriority(source, "Summon.c:CharmBrokenDistance")
		}
		return fail()
	}
	if layout.As[server.DurSpell](record).Frame68-1 != s.Frame() {
		return 0
	}
	if !spellCharmAll {
		if u.ObjSubClass&0x2000 != 0 {
			resourcePriority(source, "Summon.c:CreatureControlImpossible")
			return fail()
		}
		if source.ObjClass&4 != 0 {
			guide := int(bookGuideCharmable(uint32(u.TypeInd)))
			if !Nox_xxx_checkSummonedCreaturesLimit_500D70(source, guide) {
				resourcePriority(source, "Summon.c:CreatureControlFailed")
				return fail()
			}
		}
	}
	spellLifeBuffOff(u, 28)
	if u.FindOwnerChainPlayer() == source {
		if source.ObjClass&4 != 0 {
			resourcePriority(source, "Summon.c:CreatureAlreadyOwned")
		}
		return fail()
	}
	u.SetOwner(nil)
	u.SetOwner(source)
	if spellLifeHasTeam(u) {
		teamRuntimeLeave(u.TeamPtr(), int(u.NetCode))
	}
	u.UpdateDataMonster().StatusFlags |= 0x80
	if controlFlags(4096) {
		hp := resourceGetHP(u)
		var max uint16
		if def := u.UpdateDataMonster().MonsterDef; def != nil {
			max = *(*uint16)(unsafe.Pointer(&def.HealthQuest72))
		} else {
			max = s.Types.ByInd(int(u.TypeInd)).Health().Max
		}
		u.HealthData.Max = max
		if uint16(hp) > max {
			resourceSetHP(u, max)
		}
	}
	if source.ObjClass&4 != 0 {
		monsterOrder(source, u, int(controlPlayer(source).SummonOrderAll))
		u.ObjSubClass |= 0x180
		ind := int(controlPlayer(source).PlayerInd)
		Nox_xxx_netReportAcquireCreature_4D91A0(ind, u)
		s.Players.Nox_xxx_netMarkMinimapObject_417190(source.ControllingPlayer().PlayerIndex(), u, 1)
		Nox_xxx_netSendSimpleObject2_4DF360(ind, u)
		if controlFlags(4096) {
			spawnPolicyRelease(u)
		}
	} else {
		monsterOrder(source, u, 4)
	}
	spellEffectAudio(int32(layout.As[server.DurSpell](record).Spell), 1, source)
	return 1
}
func spellEffectCharmCancel(record unsafe.Pointer) int32 {
	u := layout.As[server.DurSpell](record).Target48
	if u != nil && u.ObjFlags&0x8020 == 0 {
		spellLifeBuffOff(u, 5)
		return int32(spellLifeBuffOff(u, 28))
	}
	return int32(memmap.PtrWord(u.CObj()))
}
func spellEffectBanish(u *server.Object) {
	typ := stateType(1570280, "Glyph")
	if u == nil {
		return
	}
	for it := u.InvFirstItem; it != nil; {
		next := it.InvNextItem
		if uint32(it.TypeInd) == typ {
			GetServer().DelayedDelete(it)
		}
		it = next
	}
	visibilityFXPoint(129, u.PosVec)
	nox_xxx_scriptCallByEventBlock_502490(unsafe.Pointer(&u.UpdateDataMonster().ScriptDeath), nil, u.CObj(), 7)
	GetServer().DelayedDelete(u)
}

// summonByte returns byte off (72-87, as in the original) of the summon record that a summon spell keeps in its
// duration spell (see durScratch).
func summonByte(record unsafe.Pointer, off int) *byte { return (*byte)(durScratch(record, off)) }
