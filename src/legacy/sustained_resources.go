package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func sustainedHasMana(u *server.Object) uint32 {
	if u.ObjClass&0x400000 != 0 && u.ObjSubClass&0x18 != 0 {
		if u.UpdateDataObelisk().Mana > 0 {
			return 1
		}
	} else if u.ObjClass&2 != 0 {
		if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](u.UpdateData).StatusFlags))&0x20 != 0 {
			return 1
		}
	} else if u.ObjClass&4 != 0 && resourceGetMana(u) != 0 {
		return 1
	}
	return 0
}
func sustainedTransferMana(u, t *server.Object, amount int32) uint32 {
	if u.ObjClass&4 != 0 && uint16(resourceGetMana(u)) >= u.UpdateDataPlayer().ManaMax {
		return 0
	}
	moved := amount
	if t.ObjClass&0x400000 != 0 && t.ObjSubClass&0x18 != 0 {
		if t.HasTeam() && !t.TeamPtr().SameAs(u.TeamPtr()) {
			return 0
		}
		pool := (*uint32)(unsafe.Pointer(&t.UpdateDataObelisk().Mana))
		if int32(*pool) <= amount {
			moved = int32(*pool)
			amount = moved
			if !sustainedGame(4096) {
				*pool = 0
			}
		} else if !sustainedGame(4096) {
			*pool -= uint32(amount)
		}
		if moved == 0 {
			return 0
		}
		if !sustainedGame(4096) {
			t.NeedSync()
		}
	} else if t.ObjClass&2 != 0 {
		if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](t.UpdateData).StatusFlags))&0x20 != 0 {
			moved = 1
			amount = 1
		}
	} else if t.ObjClass&4 != 0 {
		mana := int32(layout.As[server.PlayerUpdateData](t.UpdateData).ManaCur)
		if mana <= amount {
			moved = mana
			amount = mana
			resourceSubMana(t, mana)
		} else {
			resourceSubMana(t, amount)
		}
	}
	if moved <= 0 {
		return 0
	}
	if sustainedGame(4096) && u.ObjClass&4 != 0 {
		mult := GetServer().S().Players.Mult
		var scale float32
		switch controlPlayer(u).InfoBytes[66] {
		case 0:
			scale = mult.Warrior.Mana
		case 1:
			scale = mult.Wizard.Mana
		case 2:
			scale = mult.Conjurer.Mana
		default:
			resourceAddMana(u, int16(moved))
			return 1
		}
		moved = int32(int16(floatToInt32(float32(float64(amount) * float64(scale)))))
	}
	resourceAddMana(u, int16(moved))
	return 1
}
func sustainedFindMana(pos types.Pointf, u *server.Object) *server.Object {
	*memmap.PtrUint32(0x5d4594, 2487828) = 0
	*memmap.PtrUint32(0x5d4594, 2487876) = 1287568416
	*memmap.PtrFloat32(0x5d4594, 2487836) = pos.X
	*memmap.PtrFloat32(0x5d4594, 2487840) = pos.Y
	GetServer().S().Map.EachObjInCircle(pos, float32(spellEffectScalar("ManaDrainRange")), func(t *server.Object) bool { sustainedManaCandidate(t, u); return true })
	return sustainedPointer(*memmap.PtrUint32(0x5d4594, 2487828))
}
func sustainedManaCandidate(t, u *server.Object) {
	if t == u || sustainedHasMana(t) == 0 || t.ObjFlags&0x8020 != 0 {
		return
	}
	weight := 1.0
	if t.ObjClass&2 != 0 {
		if !((u == nil || sustainedEnemy(u, t)) && *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](t.UpdateData).StatusFlags))&0x20 != 0) {
			return
		}
	} else if t.ObjClass&0x400000 != 0 && t.ObjSubClass&0x18 != 0 {
	} else {
		if t.ObjClass&4 == 0 {
			return
		}
		if u != nil && !sustainedEnemy(u, t) {
			return
		}
		if Nox_xxx_unitsHaveSameTeam_4EC520(t, u) || *(*byte)(unsafe.Pointer(&controlPlayer(t).Field3680))&1 != 0 {
			return
		}
		weight = 0.5
	}
	pos := types.Ptf(*memmap.PtrFloat32(0x5d4594, 2487836), *memmap.PtrFloat32(0x5d4594, 2487840))
	dx := float64(pos.X) - float64(t.PosVec.X)
	dy := float64(pos.Y) - float64(t.PosVec.Y)
	distance := weight * (dy*dy + dx*dx)
	if distance < float64(*memmap.PtrFloat32(0x5d4594, 2487876)) && spellEffectTrace(pos, t.PosVec, 5) {
		*memmap.PtrUint32(0x5d4594, 2487828) = controlRaw(t)
		*memmap.PtrFloat32(0x5d4594, 2487876) = float32(distance)
	}
}
func sustainedDrainMana(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if u != nil {
		if spellLifeHasBuff(u, 8) {
			return 1
		}
	} else if d.Flag20 == 0 {
		return 1
	}
	if d.Flag20 != 0 {
		pos := d.Pos
		if u != nil {
			pos = u.PosVec
		}
		if t := sustainedFindMana(pos, u); t != nil {
			resourceSubMana(t, 50)
		}
		return 1
	}
	if u.ObjClass&4 != 0 && uint16(resourceGetMana(u)) >= uint16(resourceGetMaxMana(u)) {
		return 1
	}
	if u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	dx := float32(math.Abs(float64(d.Pos.X) - float64(u.PosVec.X)))
	dy := float32(math.Abs(float64(d.Pos.Y) - float64(u.PosVec.Y)))
	if sustainedHurtRecently(u) || dx >= 5 || dy >= 5 {
		return 1
	}
	d.Target48 = sustainedFindMana(u.PosVec, u)
	previous := durWordObject(layout.As[server.DurSpell](p).Field36)
	if d.Target48 == nil {
		if previous != nil {
			sustainedStopRay(d, previous)
		}
		return 1
	}
	value := float32(spellEffectTable("ManaDrainCoeff", int32(d.Level)-1) + float64(*(*float32)(durScratch(p, 72))))
	n := floatToInt32(value)
	*(*float32)(durScratch(p, 72)) = float32(float64(value) - float64(n))
	if d.Target48 != previous {
		if previous != nil {
			sustainedStopRay(d, previous)
		}
		spellLifeRayMessage(d)
	}
	if sustainedTransferMana(u, d.Target48, floatToInt32(value)) != 0 && sustainedFrame()%(sustainedFPS()>>1) == 0 {
		sustainedAudio(230, u)
		sustainedAudio(229, d.Target48)
	}
	d.Field36 = controlRaw(d.Target48)
	return 0
}
func sustainedGreaterHealStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	if d.Caster16 == nil && d.Flag20 == 0 {
		return 1
	}
	u := d.Caster16
	if d.Flag20 != 0 {
		if t := sustainedHealTarget(d, nil); t != nil {
			resourceAdjustHP(t, 20)
		}
		return 1
	}
	d.Target48 = sustainedHealTarget(d, u)
	if d.Target48 == nil {
		resourcePriority(u, "ExecDur.c:GreaterHealNoTarget")
		return 1
	}
	if sustainedEnemy(u, d.Target48) {
		return 1
	}
	spellLifeRayMessage(d)
	return 0
}
func sustainedGreaterHealTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u, t := d.Caster16, d.Target48
	if t == nil || t.ObjFlags&0x8020 != 0 {
		return 1
	}
	if u != nil && spellLifeHasBuff(u, 8) {
		return 1
	}
	if !sustainedInteract(u, t) || resourceGetMana(u) == 0 {
		return 1
	}
	if u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	if sustainedHurtRecently(u) || resourceGetMaxHP(t) == resourceGetHP(t) {
		return 1
	}
	value := float32(float64(*(*float32)(durScratch(p, 72))) + float64(*memmap.PtrFloat32(0x587000, 260360+uintptr(d.Level)*4)))
	if u != nil && u.ObjClass&4 != 0 {
		m := GetServer().S().Players.Mult
		var scale float32
		switch controlPlayer(u).InfoBytes[66] {
		case 0:
			scale = m.Warrior.Health
		case 1:
			scale = m.Wizard.Health
		case 2:
			scale = m.Conjurer.Health
		default:
			scale = 1
		}
		value = float32(float64(scale) * float64(value))
	}
	n := floatToInt32(value)
	*(*float32)(durScratch(p, 72)) = float32(float64(value) - float64(n))
	resourceAdjustHP(t, n)
	resourceSubMana(u, 1)
	return 0
}
func sustainedChannelLife(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u, t := d.Caster16, d.Target48
	if t == nil || t.ObjFlags&0x8020 != 0 {
		return 1
	}
	if d.Flag20 != 0 {
		resourceAddMana(t, 20)
		resourceDamage(t, 20)
		return 1
	}
	if u != nil && spellLifeHasBuff(u, 8) {
		return 1
	}
	if t.ObjClass&2 != 0 && spellLifeMoved(t, &d.Pos) != 0 {
		return 1
	}
	if resourceGetMaxMana(t) == resourceGetMana(t) || uint16(resourceGetHP(u)) <= 1 {
		return 1
	}
	if resourceGetHP(u) != 0 {
		value := float32(spellEffectTable("ChannelLifeCoeff", int32(d.Level)-1) + float64(*(*float32)(durScratch(p, 72))))
		n := floatToInt32(value)
		*(*float32)(durScratch(p, 72)) = float32(float64(value) - float64(n))
		resourceAddMana(t, int16(n))
		resourceDamage(u, 1)
	}
	return 0
}
func sustainedShieldStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	if u == nil || u.ObjFlags&0x8020 != 0 || u.ObjClass&2 != 0 && u.ObjSubClass&4 != 0 || u.ObjClass&6 == 0 {
		return 1
	}
	dur := sustainedTableInt("ShieldDuration", d.Level)
	spellLifeApplyBuff(u, 26, int16(dur), int8(d.Level))
	d.Frame68 = uint32(dur) + sustainedFrame()
	d.Field72 = sustainedTableInt("ShieldHealth", d.Level)
	return 0
}
func sustainedShieldTick(p unsafe.Pointer) uint32 {
	u := (*server.DurSpell)(p).Target48
	if u == nil {
		return 1
	}
	return uint32(bool2int(u.ObjFlags&0x8020 != 0))
}
func sustainedShieldCancel(p unsafe.Pointer) uint32 {
	u := (*server.DurSpell)(p).Target48
	if u == nil {
		return 0
	}
	return spellLifeBuffOff(u, 26)
}
func sustainedShieldAbsorb(u *server.Object, damage int32) {
	spells := &GetServer().S().Spells.Dur
	for d := spells.List; d != nil; d = d.Next {
		if d.Target48 != u || d.Spell != 51 {
			continue
		}
		if d.Target48 == nil || d.Target48.ObjFlags&0x8020 != 0 {
			spells.CancelSpell(d)
		} else if d.Field72-damage > 0 {
			d.Field72 -= damage
		} else {
			spells.CancelSpell(d)
		}
		return
	}
	spellLifeBuffOff(u, 26)
}
func sustainedShieldDamage(u *server.Object, damage *int32, kind int32, source *server.Object) {
	if u == nil || damage == nil {
		return
	}
	sustainedAudio(131, u)
	sustainedShieldFX(u, source)
	v := *damage
	if int32(uint16(resourceGetHP(u))) > v {
		*damage = v / 2
		if *damage == 0 {
			*damage = 1
		}
		sustainedShieldAbsorb(u, *damage)
		return
	}
	if u.ObjClass&4 != 0 && controlPlayer(u).InfoBytes[66] == 0 && u.UpdateDataPlayer().State == 16 && source != nil && stateFront(&u.PosVec, int32(int16(u.Direction1)), &source.PrevPos)&1 != 0 {
		*damage = 1
	} else {
		resourceSetHP(u, 2)
		*damage = 0
	}
	sustainedShieldAbsorb(u, 999999)
	*(*unsafe.Pointer)(unsafe.Pointer(&u.Obj130)) = source.CObj()
	u.Field131 = uint32(kind)
	u.Frame134 = sustainedFrame()
}
