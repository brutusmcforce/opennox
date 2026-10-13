package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func sustainedHealTarget(d *server.DurSpell, u *server.Object) *server.Object {
	s := GetServer().S()
	return s.Nox_xxx_spellFlySearchTarget(&d.Pos2, u, s.Spells.Flags(spell.ID(d.Spell)), 400, 1, u)
}
func sustainedEnergyStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	if d.Caster16 != nil {
		GetServer().S().Spells.Dur.CancelFor(43, d.Caster16)
	}
	sustainedFX(130, d.Pos)
	return 0
}
func sustainedEnergySearch(pos types.Pointf, radius float32, u *server.Object, trap bool) *server.Object {
	*memmap.PtrFloat32(0x5d4594, 2487868) = pos.X
	*memmap.PtrFloat32(0x5d4594, 2487872) = pos.Y
	*sustainedGlobal(6) = 0
	*sustainedGlobal(0) = math.Float32bits(radius * radius)
	*memmap.PtrUint32(0x5d4594, 2487832) = uint32(bool2int(trap))
	GetServer().S().Map.EachObjInCircle(pos, radius, func(t *server.Object) bool { sustainedEnergyCandidate(t, u); return true })
	return sustainedPointer(*sustainedGlobal(6))
}
func sustainedEnergyCandidate(t, u *server.Object) {
	if t.ObjClass&0x20006 == 0 || t.ObjFlags&0x8020 != 0 || t == u || t.ObjClass&2 != 0 && t.ObjSubClass&0x8000 != 0 {
		return
	}
	if u != nil && (!sustainedEnemy(u, t) || *memmap.PtrUint32(0x5d4594, 2487832) == 0 && (sustainedFront(u, t)&1 == 0 || !sustainedInteract(u, t))) {
		return
	}
	dx := float64(t.PosVec.X) - float64(*memmap.PtrFloat32(0x5d4594, 2487868))
	dy := float64(t.PosVec.Y) - float64(*memmap.PtrFloat32(0x5d4594, 2487872))
	dist := dy*dy + dx*dx
	if dist < float64(math.Float32frombits(*sustainedGlobal(0))) {
		*sustainedGlobal(0) = math.Float32bits(float32(dist))
		*sustainedGlobal(6) = controlRaw(t)
	}
}
func sustainedEnergyTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if u != nil {
		if spellLifeHasBuff(u, 8) {
			return 1
		}
	} else if d.Flag20 == 0 {
		return 1
	}
	radius := float32(spellEffectScalar("LightningRange"))
	if d.Flag20 != 0 {
		t := sustainedEnergySearch(d.Pos, radius, u, true)
		if t != nil {
			projectileDamage(t, d.Obj12, nil, sustainedScalarInt("EnergyBoltGlyphDamage"), 17)
			spellEffectAudio(24, 0, sustainedPointer(*sustainedGlobal(6)))
			sustainedFX(130, sustainedPointer(*sustainedGlobal(6)).PosVec)
		}
		return 1
	}
	if u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	if sustainedFrame()-d.Frame60 > 2 && sustainedHurtRecently(u) {
		return 1
	}
	t := d.Target48
	valid := t != nil && t.ObjFlags&0x8020 == 0 && sustainedFront(u, t)&1 != 0 && stateDistance(t, u) <= float64(radius) && sustainedInteract(u, t)
	if !valid {
		d.Target48 = nil
		if u.ObjClass&4 != 0 {
			t = layout.As[server.PlayerUpdateData](u.UpdateData).CursorObj
			if t != nil && sustainedEnemy(u, t) && stateDistance(u, t) <= float64(radius) {
				d.Target48 = t
			}
		}
		if d.Target48 == nil {
			d.Target48 = sustainedEnergySearch(u.PosVec, radius, u, false)
		}
	}
	old := durWordObject(layout.As[server.DurSpell](p).Field36)
	if d.Target48 == nil {
		if old != nil {
			sustainedStopRay(d, old)
			d.Field36 = 0
		}
		return 0
	}
	value := float32(spellEffectTable("EnergyBoltDamage", int32(d.Level)-1) + float64(*(*float32)(durScratch(p, 72))))
	n := floatToInt32(value)
	*(*float32)(durScratch(p, 72)) = float32(float64(value) - float64(n))
	if d.Target48 != old {
		if old != nil {
			sustainedStopRay(d, old)
		}
		spellLifeRayMessage(d)
	}
	projectileDamage(d.Target48, d.Caster16, nil, floatToInt32(value), 17)
	if d.Target48.ObjFlags&0x8020 != 0 {
		sustainedFX(130, d.Target48.PosVec)
	}
	d.Field36 = controlRaw(d.Target48)
	if u.ObjClass&4 != 0 {
		sustainedState(u, 10)
	}
	if sustainedFrame()%(sustainedFPS()/3) == 0 {
		sustainedAudio(32, u)
		sustainedAudio(32, d.Target48)
	}
	d.Frame68 = sustainedFrame() + uint32(sustainedScalarInt("LightningSearchTime"))
	if u.ObjClass&4 != 0 {
		sustainedState(u, 10)
		resourceSubMana(u, 1)
		if resourceGetMana(u) == 0 {
			return 1
		}
	}
	if d.Target48.ObjFlags&0x8000 != 0 {
		d.Frame68 = sustainedFrame() + 1
	}
	return 0
}
func sustainedFirewalk(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	if u == nil || u.ObjFlags&0x8020 != 0 {
		return 1
	}
	if *memmap.PtrUint32(0x5d4594, 2487888) == 0 {
		for i, name := range []string{"SmallFlame", "MediumFlame", "Flame"} {
			*memmap.PtrUint32(0x5d4594, 2487888+uintptr(i)*4) = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	if d.Frame60 == d.Frame64 {
		*(*types.Pointf)(durScratch(p, 72)) = u.PosVec
		*(*types.Pointf)(durScratch(p, 80)) = u.PosVec
		d.Frame64++
		return 0
	}
	last := *(*types.Pointf)(durScratch(p, 72))
	prev := *(*types.Pointf)(durScratch(p, 80))
	dx := float64(u.PosVec.X) - float64(last.X)
	dy := float64(u.PosVec.Y) - float64(last.Y)
	limit := 0
	if d.Level >= 2 {
		limit = 1
		if d.Level >= 4 {
			limit = 2
		}
	}
	if math.Sqrt(dy*dy+dx*dx)-float64(u.Shape.Circle.R) > 15 {
		pos := last
		for i := 0; i < 2; i++ {
			id := *memmap.PtrUint32(0x5d4594, 2487888+uintptr(GetServer().S().Rand.Logic.IntClamp(0, limit))*4)
			if flame := spellEffectNew(id); flame != nil {
				spellEffectCreate(flame, nil, pos)
				GetServer().S().Audio.EventPos(sound.ID(46), pos, 0, 0)
				Nox_xxx_unitSetDecayTime_511660(flame, int(25*sustainedFPS()))
			}
			x := pos.X - prev.X
			y := pos.Y - prev.Y
			if x != 0 && y != 0 {
				pos.X = float32(float64(pos.X) - float64(x)*0.5)
				pos.Y = float32(float64(pos.Y) - float64(y)*0.5)
			}
		}
		*(*types.Pointf)(durScratch(p, 80)) = last
		*(*types.Pointf)(durScratch(p, 72)) = u.PosVec
	}
	return 0
}
func sustainedForceStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if d.Flag20 != 0 {
		u = d.Obj24
	} else if u.ObjClass&4 != 0 {
		if item := layout.As[server.PlayerUpdateData](u.UpdateData).Field26; item != nil && item.ObjSubClass&0x200000 != 0 && item.UseDataWand().Flags&4 != 0 {
			d.Flags88 |= 2
		}
	}
	*(*uint16)(durScratch(p, 72)) = uint16(u.Direction1)
	if charge := sustainedNew("ForceOfNatureCharge"); charge != nil {
		spellEffectCreate(charge, nil, u.PosVec)
		layout.As[server.DurSpell](p).Field76 = uintptr(charge.CObj())
	}
	return 0
}
func sustainedDirection(dir uint16) types.Pointf {
	off := uintptr(int32(int16(dir)) * 8)
	return types.Ptf(*memmap.PtrFloat32(0x587000, 194136+off), *memmap.PtrFloat32(0x587000, 194140+off))
}
func sustainedForceTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if u != nil && spellLifeHasBuff(u, 8) {
		return 1
	}
	if d.Frame68-7 == sustainedFrame() {
		if charge := (*server.Object)(unsafe.Pointer(layout.As[server.DurSpell](p).Field76)); charge != nil {
			sustainedDelete(charge)
			d.Field76 = 0
		}
	}
	if d.Frame68-1 != sustainedFrame() {
		if d.Flag20 == 0 && u != nil && u.ObjClass&4 != 0 {
			sustainedState(u, 10)
		}
		return 0
	}
	ball := sustainedNew("DeathBall")
	if ball == nil {
		return 1
	}
	pos := d.Pos
	dir := *(*uint16)(durScratch(p, 72))
	distance := 0.0
	if d.Flag20 == 0 {
		pos = u.PosVec
		dir = uint16(u.Direction1)
		switch *(*uint32)(unsafe.Pointer(&u.Shape.Kind)) {
		case 1:
			distance = 4
		case 2:
			distance = float64(u.Shape.Circle.R) + 4
		case 3:
			distance = math.Max(float64(u.Shape.Box.W), float64(u.Shape.Box.H)) + 4
		default:
			distance = 24
		}
	}
	vec := sustainedDirection(dir)
	dest := types.Ptf(float32(float64(vec.X)*distance+float64(pos.X)), float32(float64(vec.Y)*distance+float64(pos.Y)))
	if !spellEffectTrace(pos, dest, 5) {
		dest = pos
	}
	spellEffectCreate(ball, u, dest)
	if d.Flag20 == 1 {
		ball.VelVec = types.Pointf{}
	} else {
		speed := ball.SpeedCur
		ball.VelVec = types.Ptf(vec.X*speed, vec.Y*speed)
	}
	ball.Direction1 = server.Dir16(dir)
	ball.Direction2 = server.Dir16(dir)
	sustainedAudio(38, u)
	return 1
}
func sustainedForceCancel(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	if charge := (*server.Object)(unsafe.Pointer(layout.As[server.DurSpell](p).Field76)); charge != nil {
		sustainedDelete(charge)
	}
	u := d.Caster16
	ret := controlRaw(u)
	if u != nil && u.ObjClass&4 != 0 {
		item := layout.As[server.PlayerUpdateData](u.UpdateData).Field26
		ret = controlRaw(item)
		if item != nil && item.ObjSubClass&0x200000 != 0 {
			ret = memmap.PtrWord(item.UseData.Ptr)
			item.UseDataWand().Flags &= ^uint32(4)
		}
	}
	return ret
}
func sustainedManaBombStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	id := stateType(2487920, "ManaBombCharge")
	d.Field72 = sustainedTableInt("ManaBombInitPower", d.Level)
	d.Field76 = 0
	d.Field84 = 0
	u := d.Caster16
	if u != nil && u.ObjClass&4 == 0 || d.Flag20 != 0 {
		d.Frame68 = sustainedFrame() + uint32(sustainedScalarInt("ManaBombGlyphDuration"))
	} else {
		for _, buff := range []int32{5, 14, 29} {
			spellLifeApplyBuff(u, buff, int16(10*sustainedFPS()), 5)
		}
		d.Field80 = *(*uint32)(unsafe.Pointer(&u.Mass))
		*(*uint32)(unsafe.Pointer(&u.Mass)) = 1203982323
		u.VelVec, u.ForceVec, u.Pos24 = types.Pointf{}, types.Pointf{}, types.Pointf{}
	}
	if charge := spellEffectNew(id); charge != nil {
		spellEffectCreate(charge, nil, d.Pos)
		layout.As[server.DurSpell](p).Field76 = uintptr(charge.CObj())
	}
	return 0
}
func sustainedManaBombTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if u == nil && d.Flag20 == 0 {
		return 1
	}
	explode := d.Frame68-1 == sustainedFrame() || u != nil && u.ObjClass&4 != 0 && resourceGetMana(u) == 0
	if charge := (*server.Object)(unsafe.Pointer(layout.As[server.DurSpell](p).Field76)); charge != nil {
		remove := d.Frame68-sustainedFrame() < 10
		if d.Flag20 == 0 && u != nil && u.ObjClass&4 != 0 {
			remove = uint16(resourceGetMana(u)) < 15
		}
		if remove {
			sustainedDelete(charge)
			d.Field76 = 0
		}
	}
	if explode {
		pos := d.Pos
		if d.Flag20 == 0 {
			if u == nil {
				return 1
			}
			pos = u.PosVec
		}
		nox_xxx_gameSetWallsDamage_4E25A0(1)
		inner := float32(spellEffectScalar("ManaBombInRadius"))
		outer := float32(spellEffectScalar("ManaBombOutRadius"))
		mapDamageUnitsAround(pos, outer, inner, int32(d.Field72), 15, u, nil)
		gameplayReportEarthquake(&pos, int(sustainedScalarInt("ManaBombShakeMag")))
		sustainedFX(129, pos)
		sustainedFX(154, pos)
		sustainedAudio(81, u)
		d.Field84 = 1
		return 1
	}
	d.Field72 += sustainedTableInt("ManaBombDeltaPower", d.Level)
	if d.Flag20 == 0 && u != nil && u.ObjClass&4 != 0 {
		resourceSubMana(u, int32(d.Level))
	}
	return 0
}
func sustainedManaBombCancel(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	if charge := (*server.Object)(unsafe.Pointer(layout.As[server.DurSpell](p).Field76)); charge != nil {
		sustainedDelete(charge)
		d.Field76 = 0
	}
	u := d.Caster16
	if d.Flag20 == 0 && u != nil && u.ObjClass&4 != 0 {
		for _, buff := range []int32{5, 14, 29} {
			spellLifeBuffOff(u, buff)
		}
		*(*uint32)(unsafe.Pointer(&u.Mass)) = d.Field80
	}
	if d.Field84 == 0 {
		return sustainedFX(163, d.Pos)
	}
	return d.Field84
}
func sustainedTurnUndeadStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	d.Field72 = sustainedTableInt("TurnUndeadKillPoints", d.Level)
	id := stateType(2487924, "UndeadKiller")
	u := d.Caster16
	if d.Flag20 != 0 {
		u = d.Obj24
	}
	pos := u.PosVec
	for dir := uint16(0); dir < 256; dir += 6 {
		if killer := spellEffectNew(id); killer != nil {
			killer.CollideDataUndeadKiller().Spell = memmap.PtrWord(p) // the duration spell, as a pointer word
			spellEffectCreate(killer, d.Caster16, pos)
			killer.Direction1 = server.Dir16(dir)
			killer.Direction2 = server.Dir16(dir)
			vec := sustainedDirection(dir)
			killer.VelVec = types.Ptf(vec.X*4, vec.Y*4)
			*(*uint32)(unsafe.Pointer(&killer.Float28)) = 0
		}
	}
	sustainedFX(160, pos)
	return 0
}
func sustainedTurnUndeadTick() uint32 { return 0 }
func sustainedTurnUndeadCancel(p unsafe.Pointer) uint32 {
	id := stateType(2487928, "UndeadKiller")
	for u := asObjectS(nox_server_getFirstObject_4DA790()); u != nil; u = asObjectS(nox_server_getNextObject_4DA7A0(asObjectC(u))) {
		if uint32(u.TypeInd) == id && u.CollideDataUndeadKiller().SpellPtr() == p {
			sustainedDelete(u)
		}
	}
	return 0
}
