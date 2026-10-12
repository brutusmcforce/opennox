package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func temporaryRefWord(p unsafe.Pointer, off int) **server.Object {
	return (**server.Object)(unsafe.Add(p, off))
}
func temporaryFloat(p unsafe.Pointer, off int) *float32 { return (*float32)(unsafe.Add(p, off)) }
func temporarySpark(pos, vel types.Pointf, stage, life int32, z float32, owner *server.Object) *server.Object {
	core := GetServer().S()
	u := core.NewObjectByTypeID("Spark")
	if u == nil {
		return nil
	}
	cd, ud := u.CollideData, u.UpdateData
	var creator server.Obj
	if owner != nil {
		creator = owner
	}
	GetServer().CreateObjectAt(u, creator, pos)
	u.Field34 = core.Frame()
	*equipmentWord(ud, 0) = uint32(life)
	*equipmentWord(ud, 4) = uint32(life)
	*equipmentWord(ud, 12) = uint32(stage)
	class := uint32(u.ObjClass)&^0x2000 | 0x80000
	// The original expression converts the flag word's float interpretation.
	flags := uint32(effectsTruncWord(float64(*temporaryFloat(u.CObj(), 16)))) & 0xff7fffbf
	u.ObjClass = object.Class(class)
	u.ObjFlags = object.Flags(flags)
	collision := uint32(0)
	heightRate := float32(0)
	switch stage {
	case 0:
		u.ObjFlags = object.Flags(flags | 0x40)
		u.ObjClass = object.Class(class &^ 0x80000)
	case 1:
		u.ObjClass = object.Class(class | 0x2000)
		u.ObjFlags = object.Flags(flags | 0x800000)
		collision = 3
		heightRate = 7
	case 2:
		u.ObjFlags = object.Flags(flags | 0x800040)
		heightRate = 7
	case 4:
		u.ObjClass = object.Class(class &^ 0x80000)
	}
	*equipmentWord(cd, 0) = collision
	Nox_xxx_unitRaise_4E46F0(u, 28)
	u.Field27 = z
	u.Field29 = math.Float32bits(heightRate)
	u.VelVec = vel
	return u
}
func temporarySparkUpdate(u *server.Object) {
	life := equipmentWord(u.UpdateData, 4)
	if int32(*life) <= 0 {
		GetServer().DelayedDelete(u)
		return
	}
	*life--
	u.Float28 = math.Float32frombits(1064514355)
	if *equipmentWord(u.UpdateData, 12) == 4 {
		u.Float28 = 1
	}
}
func temporaryTrail(u *server.Object) *server.Object {
	dx, dy := movementDirectionVector(int32(int16(u.Direction1)))
	u.ForceVec.X = float32(float64(dx)*float64(u.SpeedCur)*.25 + float64(u.ForceVec.X))
	u.ForceVec.Y = float32(float64(dy)*float64(u.SpeedCur)*.25 + float64(u.ForceVec.Y))
	sx := float32((float64(u.PosVec.X) - float64(u.PrevPos.X)) * .25)
	sy := float32((float64(u.PosVec.Y) - float64(u.PrevPos.Y)) * .25)
	rng := GetServer().S().Rand.Logic
	var out *server.Object
	for i := 0; i < 4; i++ {
		px := float32(floatToInt32(float32(float64(i)*float64(sx) + float64(u.PrevPos.X))))
		py := float32(floatToInt32(float32(float64(i)*float64(sy) + float64(u.PrevPos.Y))))
		for j := 0; j < 2; j++ {
			vy := float32(rng.FloatClamp(-2, 2))
			vx := float32(rng.FloatClamp(-2, 2))
			y := float32(rng.FloatClamp(-4, 4) + float64(py))
			x := float32(rng.FloatClamp(-4, 4) + float64(px))
			out = temporarySpark(types.Ptf(x, y), types.Ptf(vx, vy), 1, 6, 0, u)
		}
	}
	return out
}
func temporaryLifetime(u *server.Object) {
	if GetServer().S().Frame()-u.Field32 > *equipmentWord(u.UpdateData, 0) {
		u.ObjFlags |= 0x8000
		if u.Death != nil {
			u.CallDeath()
		} else {
			GetServer().DelayedDelete(u)
		}
	}
}
func temporarySpellFly(u *server.Object) {
	core := GetServer().S()
	ud := u.UpdateData
	target := temporaryRefWord(ud, 4)
	key := "UnTargetedSpellLifetime"
	if *target != nil {
		key = "TargetedSpellLifetime"
	}
	life := floatToInt32(float32(core.Balance.Float(key)))
	if core.Frame()-u.Field32 > uint32(life) {
		monsterPointFX(u, 150)
		Sub_4E71F0(u)
		return
	}
	for _, off := range []int{0, 8} {
		p := temporaryRefWord(ud, off)
		if *p != nil && (*p).ObjFlags&0x20 != 0 {
			*p = nil
		}
	}
	if *target != nil && (*target).ObjFlags&0x8020 != 0 {
		*target = nil
	}
	if *target == nil && (core.Frame()-u.Field34 > uint32(int32(core.TickRate())>>2) || u.Field34 == u.Field32) {
		*target = core.Nox_xxx_spellFlySearchTarget(nil, u, core.Spells.Flags(spell.ID(*equipmentWord(ud, 12))), 600, 0, *temporaryRefWord(ud, 0))
		u.Field34 = core.Frame()
	}
	epsilon := math.Float64frombits(uint64(qword_581450_10176))
	if t := *target; t != nil {
		x := float32(float64(t.PosVec.X) - float64(u.PosVec.X))
		yd := float64(t.PosVec.Y) - float64(u.PosVec.Y)
		y := float32(yd)
		length := float32(math.Sqrt(yd*float64(y)+float64(x)*float64(x)) + epsilon)
		u.Float28 = math.Float32frombits(1063675494)
		u.ForceVec.X = float32(float64(x) * float64(u.SpeedCur) / float64(length))
		u.ForceVec.Y = float32(float64(y) * float64(u.SpeedCur) / float64(length))
	} else {
		x, y := float64(u.VelVec.X), float64(u.VelVec.Y)
		length := float32(math.Sqrt(x*x+y*y) + epsilon)
		u.Float28 = math.Float32frombits(1063675494)
		u.ForceVec.X = float32(float64(u.SpeedCur) * x / float64(length))
		u.ForceVec.Y = float32(float64(u.SpeedCur) * y / float64(length))
	}
}
func temporaryAntiCandidate(t, u *server.Object) {
	if t.ObjFlags&0x20 != 0 || t.ObjClass&1 != 0 && t.ObjSubClass&2 == 0 || t == u || u.ObjOwner != nil && t.HasOwner(u.ObjOwner) {
		return
	}
	if !GetServer().S().CanInteract(u, t, 0) {
		return
	}
	x, y := float64(u.PosVec.X)-float64(t.PosVec.X), float64(u.PosVec.Y)-float64(t.PosVec.Y)
	dist := y*y + x*x
	if dist < float64(*memmap.PtrFloat32(0x5d4594, 2488672)) {
		*memmap.PtrFloat32(0x5d4594, 2488672) = float32(dist)
		*memmap.PtrPtr(0x5d4594, 2488668) = t.CObj()
	}
}
func temporaryAntiSpell(u *server.Object) {
	core := GetServer().S()
	frame, fps := core.Frame(), uint32(core.TickRate())
	if frame-u.Field32 > 5*fps {
		GetServer().DelayedDelete(u)
		return
	}
	target := temporaryRefWord(u.UpdateData, 4)
	if *target != nil && (*target).ObjFlags&0x20 != 0 {
		*target = nil
		frame = core.Frame()
		fps = uint32(core.TickRate())
	}
	if *target == nil && frame-u.Field34 > fps>>2 {
		*memmap.PtrUint32(0x5d4594, 2488668) = 0
		*memmap.PtrUint32(0x5d4594, 2488672) = 1287568416
		core.Map.EachMissileInCircle(u.PosVec, 600, func(t *server.Object) bool { temporaryAntiCandidate(t, u); return true })
		*target = (*server.Object)(*memmap.PtrPtr(0x5d4594, 2488668))
		u.Field34 = core.Frame()
	}
	if t := *target; t != nil {
		x, y := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		length := math.Sqrt(y*y+x*x) + .1
		u.ForceVec.X = float32(x * float64(u.SpeedCur) / length)
		u.ForceVec.Y = float32(y * float64(u.SpeedCur) / length)
		if length < 10 {
			inventorySound(20, u, 0, 0)
			GetServer().DelayedDelete(u)
			GetServer().DelayedDelete(*target)
		}
	} else {
		x, y := float64(u.VelVec.X), float64(u.VelVec.Y)
		length := math.Sqrt(x*x+y*y) + .1
		u.ForceVec.X = float32(float64(u.SpeedCur) * x / length)
		u.ForceVec.Y = float32(float64(u.SpeedCur) * y / length)
	}
	pos := u.PosVec
	rng := core.Rand.Logic
	vy := float32(rng.FloatClamp(-2, 2))
	vx := float32(rng.FloatClamp(-2, 2))
	life := rng.IntClamp(15, 30)
	y := float32(rng.FloatClamp(-4, 4) + float64(pos.Y))
	x := float32(rng.FloatClamp(-4, 4) + float64(pos.X))
	temporarySpark(types.Ptf(x, y), types.Ptf(vx, vy), 3, int32(life), 0, u)
}
func temporaryMagicMissile(u *server.Object) uint32 {
	core := GetServer().S()
	ud := u.UpdateData
	caster := *temporaryRefWord(ud, 0)
	if caster.ObjFlags&0x20 != 0 || core.Frame()-u.Field32 > 3*uint32(core.TickRate()) {
		return u.CallCollideResult(0, 0)
	}
	target := temporaryRefWord(ud, 4)
	if *target != nil && (*target).ObjFlags&0x8020 != 0 {
		*target = nil
	}
	if *target == nil && byte(core.Frame())&7 == 0 {
		*target = core.Nox_xxx_spellFlySearchTarget(nil, u, 32, 600, 0, caster)
		if *target == u.ObjOwner {
			*target = nil
		}
	}
	if t := *target; t != nil {
		dir := int16(u.Direction2)
		dx, dy := movementDirectionVector(int32(dir))
		first := float32((float64(t.PosVec.Y) - float64(u.PosVec.Y)) * float64(dx))
		dot := float64(first) - (float64(t.PosVec.X)-float64(u.PosVec.X))*float64(dy)
		if dot >= 0 {
			u.Direction2 = server.Dir16(uint16(dir) + 42)
			if int16(u.Direction2) >= 256 {
				for u.Direction2 >= 256 {
					u.Direction2 -= 256
				}
			}
		} else {
			n := int16(uint16(dir) - 42)
			for n < 0 {
				n += 256
			}
			u.Direction2 = server.Dir16(n)
		}
	}
	dir := int32(int16(u.Direction2))
	dx, dy := movementDirectionVector(dir)
	u.ForceVec.X = float32(float64(dx) * float64(u.SpeedCur))
	u.ForceVec.Y = float32(float64(dy) * float64(u.SpeedCur))
	u.Float28 = math.Float32frombits(1061997773)
	return uint32(8 * dir)
}
func temporaryChakram(u *server.Object) {
	core := GetServer().S()
	ud := u.UpdateData
	it := u.InvFirstItem
	if it == nil || it.ObjFlags&0x20 != 0 {
		GetServer().DelayedDelete(u)
		return
	}
	prior := temporaryRefWord(ud, 12)
	if *prior != nil && (*prior).ObjFlags&0x20 != 0 {
		*prior = nil
	}
	target := temporaryRefWord(ud, 8)
	returning := (*byte)(unsafe.Add(ud, 24))
	owner := u.ObjOwner
	if owner == nil || owner.ObjFlags&0x20 != 0 {
		*returning = 1
		*target = nil
	} else {
		*temporaryFloat(ud, 16) = owner.PosVec.X
		*temporaryFloat(ud, 20) = owner.PosVec.Y
		if !core.MapTraceVision(u, owner) {
			*target = nil
		} else if *returning == 0 {
			*target = owner
		} else {
			goto expiry
		}
		if *returning == 0 {
			if *target != nil && (*target).ObjFlags&0x8020 != 0 {
				*target = nil
				*returning = 1
			} else {
				x := float64(*temporaryFloat(ud, 16)) - float64(u.PosVec.X)
				yd := float64(*temporaryFloat(ud, 20)) - float64(u.PosVec.Y)
				y := float32(yd)
				length := float32(math.Sqrt(yd*float64(y)+x*x) + math.Float64frombits(uint64(qword_581450_10176)))
				u.VelVec.X = float32(x * float64(u.SpeedCur) / float64(length))
				u.VelVec.Y = float32(float64(y) * float64(u.SpeedCur) / float64(length))
			}
		}
	}
expiry:
	if core.Frame()-u.Field32 > 5*uint32(core.TickRate()) {
		*returning = 1
		*target = nil
	}
}
