package legacy

import (
	"image"
	"math"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func temporaryAreaDamage(u *server.Object, outer, inner float32, damage, kind int) {
	mapDamageUnitsAround(u.PosVec, float32(outer), float32(inner), int32(damage), int32(kind), u, nil)
}
func temporaryPowderBarrel(u *server.Object) {
	core := GetServer().S()
	stamp, start := u.Field34, u.Field32
	if stamp == start {
		u.Field34 = core.Frame() + uint32(core.Rand.Logic.IntClamp(1, 5))
		return
	}
	if stamp != core.Frame() {
		if core.Frame()-start >= uint32(core.TickRate()) {
			GetServer().DelayedDelete(u)
		}
		return
	}
	if *memmap.PtrUint32(0x5d4594, 2488684) == 0 {
		*memmap.PtrUint32(0x5d4594, 2488684) = uint32(core.Types.IndByID("SmallFlame"))
		*memmap.PtrUint32(0x5d4594, 2488688) = uint32(core.Types.IndByID("MediumFlame"))
	}
	temporaryAreaDamage(u, 100, 30, 30, 7)
	spellEffectPushAround(u.PosVec, 100, 30, 60, u, nil, 0)
	pos := u.PosVec
	rng := core.Rand.Logic
	for i := 0; i < 4; i++ {
		variant := rng.IntClamp(0, 1)
		radius := float32(rng.FloatClamp(0, 15) + 10)
		dir := rng.IntClamp(0, 255)
		dx, dy := movementDirectionVector(int32(dir))
		end := types.Ptf(float32(float64(radius)*float64(dx)+float64(u.PosVec.X)), float32(float64(radius)*float64(dy)+float64(u.PosVec.Y)))
		if !core.MapTraceRayAt(pos, end, nil, nil, 1) {
			continue
		}
		typ := *memmap.PtrUint32(0x5d4594, 2488684+uintptr(4*variant))
		flame := core.Objs.NewObject(core.Types.ByInd(int(typ)))
		if flame != nil {
			GetServer().CreateObjectAt(flame, nil, end)
			Nox_xxx_unitSetDecayTime_511660(flame, int(uint32(core.TickRate())*uint32(rng.IntClamp(5, 20))))
		}
	}
}
func temporaryOneSecond(u *server.Object) {
	if GetServer().S().Frame()-u.Field32 >= uint32(GetServer().S().TickRate()) {
		GetServer().DelayedDelete(u)
	}
}
func temporaryWaterCandidate(t *server.Object, pos types.Pointf) {
	if t.ObjClass&0x2000 == 0 {
		return
	}
	x, y := float64(pos.X)-float64(t.PosVec.X), float64(pos.Y)-float64(t.PosVec.Y)
	if math.Sqrt(y*y+x*x)-float64(t.Shape.Circle.R) <= 40 {
		GetServer().DelayedDelete(t)
		*memmap.PtrUint32(0x5d4594, 2488664) = 1
	}
}
func temporaryWaterBarrel(u *server.Object) {
	age := GetServer().S().Frame() - u.Field32
	if age == 8 {
		*memmap.PtrUint32(0x5d4594, 2488664) = 0
		rect := types.Rectf{Min: types.Ptf(float32(float64(u.PosVec.X)-40), float32(float64(u.PosVec.Y)-40)), Max: types.Ptf(float32(float64(u.PosVec.X)+40), float32(float64(u.PosVec.Y)+40))}
		GetServer().S().Map.EachObjInRect(rect, func(t *server.Object) bool { temporaryWaterCandidate(t, u.PosVec); return true })
		if *memmap.PtrUint32(0x5d4594, 2488664) != 0 {
			inventorySound(283, u, 0, 0)
		}
	} else if age >= 30 {
		GetServer().DelayedDelete(u)
	}
}
func temporarySelfDestruct(u *server.Object) {
	if GetServer().S().Frame()-u.Field32 > 2 {
		GetServer().DelayedDelete(u)
	}
}
func temporaryPowderBurn(u *server.Object) {
	core := GetServer().S()
	stamp, start := u.Field34, u.Field32
	if stamp == start {
		u.Field34 = core.Frame() + 3
	} else if stamp == core.Frame() {
		temporaryAreaDamage(u, 15, 15, 1, 1)
	} else if core.Frame()-start >= 2*uint32(core.TickRate()) {
		GetServer().DelayedDelete(u)
	}
}
func temporaryDeathFragment(u *server.Object) {
	core := GetServer().S()
	age := core.Frame() - u.Field32
	if int32(age) > 10 {
		temporaryAreaDamage(u, 25, 0, 20, 2)
	}
	// gameFPS returns unsigned, so the second C comparison is unsigned too.
	if age > 2*uint32(core.TickRate()) {
		GetServer().DelayedDelete(u)
	}
}
func temporaryCursor(owner *server.Object) types.Pointf {
	p := equipmentPlayer(owner)
	return types.Ptf(float32(int32(p.CursorVec.X)), float32(int32(p.CursorVec.Y)))
}
func temporaryMoonglow(u *server.Object) {
	owner := u.ObjOwner
	if owner == nil {
		GetServer().DelayedDelete(u)
		return
	}
	if owner.ObjClass&4 == 0 {
		return
	}
	core := GetServer().S()
	if owner.ObjFlags&0x20 != 0 || core.Frame()-u.Field32 > 300*uint32(core.TickRate()) {
		GetServer().DelayedDelete(u)
		Nox_xxx_spellBuffOff_4FF5B0(u.ObjOwner, 1)
	} else {
		pos := temporaryCursor(owner)
		if sub_517590(float32(float32(pos.X)), float32(float32(pos.Y))) != 0 {
			Nox_xxx_unitMove_4E7010(u, pos)
		}
	}
}
func temporaryTelekinesis(u *server.Object) {
	core := GetServer().S()
	owner := u.ObjOwner
	if owner == nil {
		// without a caster (e.g. spawned by a script) there is nothing to move it
		GetServer().DelayedDelete(u)
		return
	}
	if owner.ObjFlags&0x8020 != 0 || core.Frame()-u.Field32 > 20*uint32(core.TickRate()) {
		sound := core.Spells.DefByInd(127).GetAudio(2)
		inventorySound(int(sound), u.ObjOwner, 0, 0)
		GetServer().DelayedDelete(u)
		Nox_xxx_spellBuffOff_4FF5B0(u.ObjOwner, 24)
	} else if owner.ObjClass&4 != 0 { // only players have a cursor to follow
		end := temporaryCursor(owner)
		if core.MapTraceRayAt(owner.PosVec, end, nil, nil, 5) {
			Nox_xxx_unitMove_4E7010(u, end)
		}
	}
}
func temporaryScorch(u *server.Object) {
	motionScorch(&u.PosVec, 2)
}
func temporaryFist(u *server.Object) {
	if u.ZVal <= 0 && int32(u.ObjFlags) >= 0 {
		inventorySound(48, u, 0, 0)
		temporaryScorch(u)
		radius := float64(u.Shape.Circle.R)
		x := float32(radius + float64(u.PosVec.X))
		u.ObjFlags |= 0x80000000
		pos := types.Ptf(x, u.PosVec.Y)
		visibilityFXPoint(138, pos)
		pos = types.Ptf(float32(float64(u.PosVec.X)-float64(u.Shape.Circle.R)), u.PosVec.Y)
		visibilityFXPoint(138, pos)
		pos = types.Ptf(u.PosVec.X, float32(float64(u.Shape.Circle.R)+float64(u.PosVec.Y)))
		visibilityFXPoint(138, pos)
		gameplayReportEarthquake(&u.PosVec, 30)
	}
	if u.ZVal >= 200 && int32(u.ObjFlags) < 0 {
		GetServer().DelayedDelete(u)
	}
	if GetServer().S().Frame()-u.Field32 > 3*uint32(GetServer().S().TickRate()) {
		GetServer().DelayedDelete(u)
	}
}
func temporaryFlameCleanse(u *server.Object) {
	core := GetServer().S()
	if core.Frame() >= u.Field34 || !core.MapTraceRayAt(u.Pos39, u.PosVec, nil, nil, 65) || core.Frame()-u.Field32 > 3 && u.PrevPos == u.PosVec {
		GetServer().DelayedDelete(u)
	}
}
func temporaryMeteorShower(u *server.Object) {
	core := GetServer().S()
	frame := core.Frame()
	if frame-u.Field32 >= 5*uint32(core.TickRate()) {
		GetServer().DelayedDelete(u)
		return
	}
	if u.Field34 > frame {
		return
	}
	data := u.UpdateDataDamage()
	rng := core.Rand.Logic
	dir := rng.IntClamp(0, 255)
	r := rng.FloatClamp(4, 12)
	r *= r
	dx, dy := movementDirectionVector(int32(dir))
	end := types.Ptf(float32(r*float64(dx)+float64(u.PosVec.X)), float32(r*float64(dy)+float64(u.PosVec.Y)))
	ray := [4]float32{u.PosVec.X, u.PosVec.Y, end.X, end.Y}
	if spatialRay(&ray) {
		m := core.NewObjectByTypeID("Meteor")
		if m != nil {
			m.UpdateDataDamage().Damage = data.Damage
			GetServer().CreateObjectAt(m, u, end)
			m.Field5 |= 0x20
			Nox_xxx_unitRaise_4E46F0(m, 255)
			m.Field27 = math.Float32frombits(0xc1000000)
		}
	}
	u.Field34 = core.Frame() + uint32(rng.IntClamp(4, 8))
}
func temporaryMeteorExplode(u *server.Object) {
	if !(u.ZVal <= 0) {
		return
	}
	damage := &u.UpdateDataDamage().Damage
	inventorySound(87, u, 0, 0)
	temporaryScorch(u)
	gameplayReportEarthquake(&u.PosVec, 10)
	if fx := GetServer().S().NewObjectByTypeID("MeteorExplode"); fx != nil {
		GetServer().CreateObjectAt(fx, nil, u.PosVec)
	}
	owner := u.FindOwnerChainPlayer()
	mapDamageUnitsAround(u.PosVec, 80, 30, int32(*damage), 7, owner, nil)
	x1 := float32(float64(u.PosVec.X) - 80)
	y1 := float32(float64(u.PosVec.Y) - 80)
	x2 := float32(float64(u.PosVec.X) + 80)
	y2 := float32(float64(u.PosVec.Y) + 80)
	rect := image.Rect(int(floatToInt32(x1)/23), int(floatToInt32(y1)/23), int(floatToInt32(x2)/23), int(floatToInt32(y2)/23))
	GetServer().Nox_xxx_mapDamageToWalls_534FC0(rect, u.PosVec, 80, int(int32(*damage)), 7, u)
	GetServer().DelayedDelete(u)
}
func temporaryCloudCandidate(t, u *server.Object, small bool) {
	if t.ObjClass&6 == 0 || !GetServer().S().MapTraceRayAt(u.PosVec, t.PosVec, nil, nil, 5) {
		return
	}
	damage := 0
	if !small {
		damage = GetServer().S().Rand.Logic.IntClamp(3, 10)
	}
	owner := u.FindOwnerChainPlayer()
	_ = t.CallDamageValue(owner, u, int32(damage), 5)
	owner = u.FindOwnerChainPlayer()
	if GetServer().S().IsEnemyTo(owner, t) {
		resourcePoison(t, 1, 1)
	}
}
func temporaryCloud(u *server.Object, small bool) {
	core := GetServer().S()
	life := &u.UpdateDataToxicCloud().Life
	if u.Field34 < core.Frame() {
		radius := float32(75)
		if small {
			radius = 35
		}
		core.Map.EachObjInCircle(u.PosVec, radius, func(t *server.Object) bool { temporaryCloudCandidate(t, u, small); return true })
		u.Field34 = core.Frame() + uint32(core.Rand.Logic.IntClamp(5, 10))
	}
	if *life != 0 {
		*life--
	}
	if *life == 0 {
		GetServer().DelayedDelete(u)
	}
}
func temporaryArachnaphobia(u *server.Object) {
	core := GetServer().S()
	frame := core.Frame()
	if u.Field34 < core.Frame() {
		id := memmap.PtrUint32(0x5d4594, 2488708)
		if *id == 0 {
			*id = uint32(core.Types.IndByID("SmallSpider"))
		}
		if spider := core.Objs.NewObject(core.Types.ByInd(int(*id))); spider != nil {
			var owner server.Obj
			if u.ObjOwner != nil {
				owner = u.ObjOwner
			}
			GetServer().CreateObjectAt(spider, owner, u.PosVec)
		}
		u.Field34 = core.Frame() + uint32(core.Rand.Logic.IntClamp(1, 5))
		frame = core.Frame()
	}
	if frame-u.Field32 > 3*uint32(core.TickRate()) {
		GetServer().DelayedDelete(u)
	}
}
func temporaryExpire(u *server.Object) {
	if u.Field32 > GetServer().S().Frame() || u.Field34 < GetServer().S().Frame() {
		GetServer().DelayedDelete(u)
	}
}
func temporaryBreak(u *server.Object, open bool) uint32 {
	result := u.Field5
	core := GetServer().S()
	if result&2 != 0 {
		result = uint32(u.ObjFlags)
		if result&0x8000 != 0 {
			u.UnsetXStatus(2)
			u.SetXStatus(4)
			if open {
				result = core.Frame()
				u.Field34 = result + 2*uint32(core.TickRate())
			} else {
				result = uint32(u.ObjFlags) | 0x40
				u.Field34 = core.Frame() + 2*uint32(core.TickRate())
				u.ObjFlags |= 0x40
			}
		}
	} else if result&4 != 0 {
		result = core.Frame()
		if core.Frame() > u.Field34 {
			u.UnsetXStatus(4)
			u.SetXStatus(8)
		}
	} else if result&8 != 0 {
		core.Objs.RemoveFromUpdatable(u)
	}
	return result
}
func temporaryBreakRemove(u *server.Object) {
	core := GetServer().S()
	switch u.Field5 {
	case 2:
		if u.ObjFlags&0x8000 != 0 {
			u.UnsetXStatus(2)
			u.SetXStatus(4)
			u.Field34 = core.Frame() + 2*uint32(core.TickRate())
			u.ObjFlags |= 0x40
		}
	case 4:
		if core.Frame() > u.Field34 {
			u.UnsetXStatus(4)
			u.SetXStatus(8)
		}
	case 8:
		core.Objs.RemoveFromUpdatable(u)
		GetServer().DelayedDelete(u)
	}
}
