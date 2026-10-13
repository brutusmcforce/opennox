package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func sustainedLightningStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	d.Field72 = 0
	d.Field76 = 0
	u := d.Caster16
	if u != nil {
		GetServer().S().Spells.Dur.CancelFor(24, u)
	}
	if d.Flag20 == 0 && u.ObjClass&4 != 0 {
		if item := layout.As[server.PlayerUpdateData](u.UpdateData).Field26; item != nil && item.ObjSubClass&0x40000 != 0 && item.UseDataWand().Flags&4 != 0 {
			d.Field72 = int32(controlRaw(item))
		}
	}
	if d.Field72 != 0 {
		d.Flags88 |= 2
	}
	sustainedFX(129, d.Pos)
	return 0
}
func sustainedLightningArray(i uint32) *uint32 {
	return memmap.PtrUint32(0x5d4594, 2487844+uintptr(i)*4)
}
func sustainedLightningObject(i uint32) *server.Object {
	return sustainedPointer(*sustainedLightningArray(i))
}
func sustainedLightningCandidate(t, u *server.Object) {
	if t.ObjClass&0x20006 == 0 {
		return
	}
	owner := sustainedPointer(*sustainedGlobal(2))
	if owner != nil && !sustainedEnemy(owner, t) {
		return
	}
	if t.ObjClass&2 != 0 && t.ObjSubClass&0x8000 != 0 || t.ObjFlags&0x8020 != 0 || t == u || t == owner {
		return
	}
	for i := int32(0); i < int32(*sustainedGlobal(3)); i++ {
		if sustainedLightningObject(uint32(i)) == t {
			return
		}
	}
	if !sustainedInteract(u, t) {
		return
	}
	dx := float64(t.PosVec.X) - float64(u.PosVec.X)
	dy := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	dist := dy*dy + dx*dx
	if dist < float64(math.Float32frombits(*sustainedGlobal(5))) {
		*sustainedGlobal(5) = math.Float32bits(float32(dist))
		*sustainedGlobal(4) = controlRaw(t)
	}
}
func sustainedLightningTrapHit(t, u *server.Object) {
	if t == u || t.ObjClass&6 == 0 || t.ObjFlags&0x8020 != 0 || t.ObjClass&2 != 0 && t.ObjSubClass&0x8000 != 0 || u != nil && !sustainedEnemy(u, t) {
		return
	}
	pos := types.Ptf(*memmap.PtrFloat32(0x5d4594, 2487820), *memmap.PtrFloat32(0x5d4594, 2487824))
	if spellEffectTrace(pos, t.PosVec, 9) {
		projectileDamage(t, nil, nil, sustainedScalarInt("LightningGlyphDamage"), 17)
		sustainedFX(129, t.PosVec)
		spellEffectAudio(43, 0, t)
	}
}
func sustainedFreeSegments(d *server.DurSpell) {
	for d != nil {
		next := d.Next
		GetServer().S().Spells.Dur.FreeRecursive(d)
		d = next
	}
}
func sustainedSpendCharge(u, item *server.Object) bool {
	data := item.UseDataWand()
	charge := &data.Charges
	if *charge == 0 {
		return false
	}
	*charge--
	data.Percent = 100 * uint32(*charge) / uint32(data.MaxCharges)
	if u != nil && u.ObjClass&4 != 0 {
		sustainedState(u, 22)
		sustainedCharges(u, item)
	}
	return *charge != 0
}
func sustainedLightningTick(p unsafe.Pointer) uint32 {
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
	core := GetServer().S()
	if d.Flag20 != 0 {
		*memmap.PtrFloat32(0x5d4594, 2487820) = d.Pos.X
		*memmap.PtrFloat32(0x5d4594, 2487824) = d.Pos.Y
		core.Map.EachObjInCircle(d.Pos, radius, func(t *server.Object) bool { sustainedLightningTrapHit(t, u); return true })
		return 1
	}
	if u.ObjClass&4 != 0 && resourceGetMana(u) == 0 {
		return 1
	}
	if sustainedFrame()-d.Frame60 > 2 && sustainedHurtRecently(u) {
		return 1
	}
	if u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	sustainedFreeSegments(d.Sub104)
	d.Sub104 = d.Sub108
	d.Sub108 = nil
	*sustainedGlobal(4) = 0
	*sustainedGlobal(3) = 0
	*sustainedGlobal(2) = controlRaw(u)
	level := int32(*memmap.PtrUint32(0x587000, 260380+uintptr(d.Level)*4))
	for i := uint32(0); i < 5; i++ {
		*sustainedLightningArray(i) = 0
	}
	if u.ObjClass&4 != 0 {
		t := layout.As[server.PlayerUpdateData](u.UpdateData).CursorObj
		if t != nil && sustainedEnemy(u, t) && stateDistance(u, t) <= float64(radius) {
			*sustainedGlobal(4) = controlRaw(t)
		}
	}
	if *sustainedGlobal(4) == 0 {
		*sustainedGlobal(5) = math.Float32bits(radius * radius)
		core.Map.EachObjInCircle(d.Pos, radius, func(t *server.Object) bool { sustainedLightningCandidate(t, u); return true })
		if *sustainedGlobal(4) == 0 {
			for seg := d.Sub104; seg != nil; seg = seg.Next {
				if seg.Target48 != nil {
					sustainedStopRay(seg, seg.Target48)
				}
			}
			sustainedFreeSegments(d.Sub104)
			d.Sub104 = nil
			return 0
		}
	}
	appendTarget := func() {
		i := *sustainedGlobal(3)
		*sustainedLightningArray(i) = *sustainedGlobal(4)
		*sustainedGlobal(3) = i + 1
	}
	appendTarget()
	for i, scale := range []float64{0.94999999, 0.89999998, 0.85000002, 0.80000001} {
		if level <= int32(i+1) {
			continue
		}
		fromIndex := uint32(0)
		if i >= 2 {
			fromIndex = uint32(i - 1)
		}
		from := sustainedLightningObject(fromIndex)
		if from == nil {
			continue
		}
		*sustainedGlobal(4) = 0
		*sustainedGlobal(5) = math.Float32bits(radius * radius)
		core.Map.EachObjInCircle(from.PosVec, float32(float64(radius)*scale), func(t *server.Object) bool { sustainedLightningCandidate(t, from); return true })
		if *sustainedGlobal(4) != 0 {
			appendTarget()
		}
	}
	core.Spells.Dur.NewLightningSub(d, u, sustainedLightningObject(0))
	if level > 1 && sustainedLightningObject(1) != nil {
		core.Spells.Dur.NewLightningSub(d, sustainedLightningObject(0), sustainedLightningObject(1))
	}
	if level > 2 && sustainedLightningObject(2) != nil {
		core.Spells.Dur.NewLightningSub(d, sustainedLightningObject(0), sustainedLightningObject(2))
	}
	if level > 3 && sustainedLightningObject(3) != nil {
		from := sustainedLightningObject(1)
		if from == nil {
			from = sustainedLightningObject(2)
		}
		if from != nil {
			core.Spells.Dur.NewLightningSub(d, from, sustainedLightningObject(3))
		}
	}
	if level > 4 && sustainedLightningObject(4) != nil {
		from := sustainedLightningObject(2)
		if from == nil {
			from = sustainedLightningObject(1)
		}
		if from != nil {
			core.Spells.Dur.NewLightningSub(d, from, sustainedLightningObject(4))
		}
	}
	if sustainedLightningObject(0) == nil {
		return 0
	}
	value := float32(spellEffectScalar("LightningDamage") + float64(*(*float32)(durScratch(p, 76))))
	damage := floatToInt32(value)
	*(*float32)(durScratch(p, 76)) = float32(float64(value) - float64(damage))
	old := d.Sub104
	for seg := d.Sub108; seg != nil; seg = seg.Next {
		if old != nil {
			if seg.Target48 != old.Target48 || seg.Caster16 != old.Caster16 {
				if old.Target48 != nil {
					sustainedStopRay(old, old.Target48)
				}
				spellLifeRayMessage(seg)
			}
			old = old.Next
		} else {
			spellLifeRayMessage(seg)
		}
		if damage > 0 {
			projectileDamage(seg.Target48, d.Caster16, nil, damage, 17)
		}
		if seg.Target48.ObjFlags&0x8020 != 0 {
			sustainedFX(129, seg.Target48.PosVec)
		}
	}
	for ; old != nil; old = old.Next {
		if old.Target48 != nil {
			sustainedStopRay(old, old.Target48)
		}
	}
	if u.ObjClass&4 != 0 {
		if item := durWordObject(uint32(layout.As[server.DurSpell](p).Field72)); item != nil {
			if !sustainedSpendCharge(u, item) {
				return 1
			}
		} else {
			sustainedState(u, 10)
			resourceSubMana(u, 1)
			if resourceGetMana(u) == 0 {
				return 1
			}
		}
	}
	if sustainedFrame()%(sustainedFPS()/3) == 0 {
		sustainedAudio(78, u)
		sustainedAudio(78, sustainedLightningObject(0))
	}
	d.Frame68 = sustainedFrame() + uint32(sustainedScalarInt("LightningSearchTime"))
	return 0
}
func sustainedLightningCancel(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	sustainedFreeSegments(d.Sub108)
	d.Sub108 = nil
	sustainedFreeSegments(d.Sub104)
	d.Sub104 = nil
	var ret uint32
	if item := durWordObject(uint32(layout.As[server.DurSpell](p).Field72)); item != nil {
		word := &item.UseDataWand().Flags
		*word &^= 4
		ret = *word
	}
	return uint32(int32(int8(ret)))
}
func sustainedPlasmaStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	d.Field72 = 0
	d.Field76 = 0
	d.Target48 = nil
	sustainedFX(131, d.Pos)
	u := d.Caster16
	if u == nil || u.ObjClass&4 == 0 {
		return 0
	}
	item := layout.As[server.PlayerUpdateData](u.UpdateData).Field26
	if item == nil || item.ObjSubClass&0x4000000 == 0 {
		return 1
	}
	if item.UseDataWand().Flags&4 != 0 {
		d.Field72 = int32(controlRaw(item))
		d.Flags88 |= 2
	}
	return 0
}
func sustainedPlasmaCandidate(t, u *server.Object) {
	if t.ObjClass&0x20006 == 0 || t.ObjFlags&0x8020 != 0 || t == u || t.ObjClass&2 != 0 && t.ObjSubClass&0x8000 != 0 {
		return
	}
	if !sustainedEnemy(u, t) {
		return
	}
	_ = sustainedFront(u, t) // Original bitwise OR with 0xC makes this direction predicate always true.
	if !sustainedInteract(u, t) {
		return
	}
	dx := float64(t.PosVec.X) - float64(u.PosVec.X)
	dy := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	distance := dy*dy + dx*dx
	if distance < float64(*memmap.PtrFloat32(0x587000, 260404)) {
		*memmap.PtrFloat32(0x587000, 260404) = float32(distance)
		*sustainedGlobal(1) = controlRaw(t)
	}
}
func sustainedPlasmaTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	if *memmap.PtrUint32(0x5d4594, 2487936) == 0 {
		for i, name := range []string{"Hecubah", "HecubahWithOrb"} {
			*memmap.PtrUint32(0x5d4594, 2487936+uintptr(i)*4) = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	u := d.Caster16
	if u == nil || d.Flags88&0x20 != 0 || u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	t := d.Target48
	if t != nil {
		if t.ObjFlags&0x8020 != 0 {
			d.Target48 = nil
		} else {
			if sustainedFront(u, t)&2 != 0 {
				if d.Field76 == 0 {
					d.Field76 = uintptr(3 * sustainedFPS())
				}
			} else {
				d.Field76 = 0
			}
			if stateDistance(t, u) > 400 || !sustainedInteract(u, t) {
				d.Target48 = nil
			}
		}
	}
	if d.Target48 == nil {
		t = layout.As[server.PlayerUpdateData](u.UpdateData).CursorObj
		*sustainedGlobal(1) = 0
		if t != nil && sustainedEnemy(u, t) && stateDistance(u, t) <= 400 {
			*sustainedGlobal(1) = controlRaw(t)
		}
		if *sustainedGlobal(1) == 0 {
			*memmap.PtrUint32(0x587000, 260404) = 1209810944
			GetServer().S().Map.EachObjInCircle(u.PosVec, 400, func(t *server.Object) bool { sustainedPlasmaCandidate(t, u); return true })
		}
		d.Target48 = sustainedPointer(*sustainedGlobal(1))
		d.Field76 = 0
	}
	if d.Field76 != 0 {
		d.Field76--
		if d.Field76 == 0 {
			d.Target48 = nil
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
	if d.Target48 != old {
		if old != nil {
			sustainedStopRay(d, old)
		}
		spellLifeRayMessage(d)
	}
	key := "PlasmaDamage"
	id := uint32(d.Target48.TypeInd)
	if id == *memmap.PtrUint32(0x5d4594, 2487936) || id == *memmap.PtrUint32(0x5d4594, 2487940) {
		key = "PlasmaDamageHecubah"
	}
	projectileDamage(d.Target48, u, nil, sustainedScalarInt(key), 14)
	if d.Target48.ObjFlags&0x8020 != 0 {
		sustainedFX(131, d.Target48.PosVec)
	}
	d.Field36 = controlRaw(d.Target48)
	sustainedState(u, 22)
	if sustainedFrame()%(sustainedFPS()/3) == 0 {
		sustainedAudio(98, u)
		sustainedAudio(98, d.Target48)
	}
	if d.Field76 == 0 {
		d.Frame68 = sustainedFrame() + uint32(sustainedScalarInt("PlasmaSearchTime"))
	}
	if item := durWordObject(uint32(layout.As[server.DurSpell](p).Field72)); item != nil && !sustainedSpendCharge(u, item) {
		return 1
	}
	return 0
}
func sustainedPlasmaCancel(p unsafe.Pointer) uint32 {
	item := durWordObject(uint32(layout.As[server.DurSpell](p).Field72))
	if item == nil {
		return 0
	}
	item.UseDataWand().Flags &= ^uint32(4)
	return memmap.PtrWord(item.UseData.Ptr)
}
