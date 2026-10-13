package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func worldElevatorSound(u *server.Object, up bool) {
	id := 0
	switch u.Material {
	case 8:
		id = 258
	case 16:
		if u.ObjSubClass&0x20 != 0 {
			id = 254
		} else if u.ObjSubClass&0x40 != 0 {
			id = 260
		} else {
			id = 252
		}
	case 32:
		if u.ObjSubClass&2 != 0 {
			id = 256
		} else {
			id = 250
		}
	}
	if id != 0 {
		if up {
			id--
		}
		inventorySound(id, u, 0, 0)
	}
}
func worldPlatformContains(platform, target *server.Object) bool {
	return collisionContains(&platform.PosVec, (*[11]float32)(unsafe.Pointer(&platform.Shape)), &target.PosVec)
}
func worldShaftCandidate(t, u *server.Object) {
	if !worldPlatformContains(u, t) {
		return
	}
	partner := u.UpdateDataElevatorShaft().Partner
	height := &partner.UpdateDataElevator().Field_4
	delta := float32(float64(t.ZVal) + 64 - float64(int32(*height)))
	if float64(worldAbsScratch(float32(delta))) < 10 {
		Nox_xxx_unitMove_4E7010(t, partner.PosVec)
		Nox_xxx_unitRaise_4E46F0(t, float32(int32(*height)))
	}
}
func worldElevatorCandidate(t, u *server.Object) {
	if !worldPlatformContains(u, t) {
		return
	}
	ud := u.UpdateDataElevator()
	partner := ud.Partner
	if partner == nil {
		return
	}
	delta := float32(float64(t.ZVal) - float64(int32(ud.Field_4)))
	if !(float64(worldAbsScratch(float32(delta))) < 10) {
		return
	}
	tooBig := false
	switch t.Shape.Kind {
	case server.ShapeKindBox:
		tooBig = partner.Shape.Box.W < t.Shape.Box.W || partner.Shape.Box.H < t.Shape.Box.H
	case server.ShapeKindCircle:
		diameter := float64(t.Shape.Circle.R) * 2
		tooBig = diameter > float64(partner.Shape.Box.W) || diameter > float64(partner.Shape.Box.H)
	}
	if tooBig {
		t.ObjFlags &^= 0x100000
		Nox_xxx_unitRaise_4E46F0(t, 0)
		return
	}
	// Reload the partner and height after retained calls, as the original does.
	Nox_xxx_unitMove_4E7010(t, ud.Partner.PosVec)
	Nox_xxx_unitRaise_4E46F0(t, float32(int32(ud.Field_4-64)))
}
func worldShaft(u *server.Object) byte {
	ud := u.UpdateDataElevatorShaft()
	partner := ud.Partner
	if partner == nil {
		return 0
	}
	collisionActivate(u)
	other := ud.Partner.UpdateDataElevator()
	state := &other.Field_3
	prior := &ud.State
	if *state == 1 {
		if other.Field_4 <= 32 {
			GetServer().S().Map.EachObjInCircle(u.PosVec, 64, func(t *server.Object) bool { worldShaftCandidate(t, u); return true })
		}
		if *prior != *state {
			worldElevatorSound(u, false)
		}
	} else if *state == 3 && *prior != *state {
		worldElevatorSound(u, true)
	}
	*prior = *state
	return *state
}
func worldElevator(u *server.Object) {
	core := GetServer().S()
	ud := u.UpdateDataElevator()
	state := &ud.Field_3
	height := &ud.Field_4
	sync := func() {
		u.NeedSync()
		if partner := ud.Partner; partner != nil {
			partner.NeedSync()
		}
	}
	switch *state {
	case 0:
		if u.ObjFlags&0x1000000 != 0 && core.Frame()-u.Field34 > uint32(core.TickRate()) {
			*state = 3
			worldElevatorSound(u, true)
		}
	case 1:
		if int32(*height) > 0 {
			*height -= 2
		} else {
			*state = 0
			u.Field34 = core.Frame()
		}
		sync()
		if int32(*height) <= 20 {
			u.ObjFlags |= 0x10
		}
	case 2:
		if u.ObjFlags&0x1000000 != 0 && core.Frame()-u.Field34 > uint32(core.TickRate()) {
			*state = 1
			worldElevatorSound(u, false)
		}
	case 3:
		*height += 2
		sync()
		if int32(*height) >= 20 {
			u.ObjFlags &^= 0x10
		}
		if int32(*height) >= 32 {
			core.Map.EachObjInCircle(u.PosVec, 64, func(t *server.Object) bool { worldElevatorCandidate(t, u); return true })
		}
		if *height >= 64 {
			*height = 64
			*state = 2
			u.Field34 = core.Frame()
		}
	}
}
func worldTeleportCandidate(t *server.Object, pos *types.Pointf, visible bool) {
	// This legacy filter converts the float interpretation of the class word.
	if uint32(effectsTruncWord(float64(*(*float32)(unsafe.Pointer(&t.ObjClass)))))&0x420000 != 0 {
		return
	}
	if visible {
		monsterPointFX(t, 137)
		inventorySound(147, t, 0, 0)
	}
	stateTeleport(t, pos)
	if visible {
		monsterPointFX(t, 137)
		inventorySound(147, t, 0, 0)
	}
}
func worldTeleportArea(u, destination *server.Object, visible bool) {
	r := float64(u.Shape.Circle.R)
	x, y := float64(u.PosVec.X), float64(u.PosVec.Y)
	rect := types.Rectf{Min: types.Ptf(float32(x-r), float32(y-r)), Max: types.Ptf(float32(x+r), float32(y+r))}
	GetServer().S().Map.EachObjInRect(rect, func(t *server.Object) bool { worldTeleportCandidate(t, &destination.PosVec, visible); return true })
}
func worldTeleport(u *server.Object) uint32 {
	ud := u.UpdateDataTransporter()
	state := &ud.State
	speed := &ud.Speed
	counter := &ud.Counter
	anim := &ud.Anim
	active := &ud.Active
	if *state != 0 {
		if *state <= 2 {
			if *counter == *speed {
				*anim++
				u.NeedSync()
				if *anim == 9 {
					*anim = 1
					*speed++
				}
				*counter = 0
			} else {
				*counter++
			}
		}
	} else {
		if *anim != 0 {
			u.NeedSync()
		}
		*anim = 0
	}
	result := uint32(*state)
	if *state != 0 {
		if result != 1 {
			result -= 2
			if result == 0 && *speed >= 4 {
				*anim = 0
				*state = 0
				*active = 0
				return result
			}
		} else {
			result = uint32(*speed)
			if *speed != 0 {
				if *speed == 4 {
					*state = 0
					*active = 0
					return result
				}
			} else if *anim == 8 {
				destination := ud.Destination
				result = memmap.PtrWord(destination.CObj())
				if destination != nil {
					worldTeleportArea(u, destination, true)
					*active = 0
					return result
				}
			}
		}
	} else if *active != 0 {
		destination := ud.Destination
		result = memmap.PtrWord(destination.CObj())
		if destination != nil && u.ObjFlags&0x1000000 != 0 {
			other := destination.UpdateDataTransporter()
			*state = 1
			*speed = 0
			*counter = 0
			u.Field34 = GetServer().S().Frame()
			other.State = 2
			other.Speed = 0
			other.Counter = 0
			result = GetServer().S().Frame()
			ud.Destination.Field34 = GetServer().S().Frame()
		}
	}
	*active = 0
	return result
}
func worldInvisibleTeleport(u *server.Object) uint32 {
	ud := u.UpdateDataTransporter()
	if ud.Active != 0 {
		destination := ud.Destination
		if destination != nil && u.ObjFlags&0x1000000 != 0 {
			worldTeleportArea(u, destination, false)
		}
	}
	ud.Active = 0
	return memmap.PtrWord(u.CObj())
}
func worldPush(u *server.Object) {
	ud := u.UpdateDataPush()
	spellEffectPushAround(u.PosVec, ud.Radius, 0, ud.Power, nil, nil, 0)
}
func worldIndexedDirection(u *server.Object) (int32, int32) {
	var out [2]int32
	geometryIndexedDirection(int32(int16(u.Direction1)), &out)
	return int32(out[0]), int32(out[1])
}
func worldBlowCandidate(t, u *server.Object) {
	if byte(effectsTruncWord(float64(*(*float32)(unsafe.Pointer(&t.ObjFlags)))))&0x20 != 0 || uint32(effectsTruncWord(float64(*(*float32)(unsafe.Pointer(&t.ObjClass)))))&0x400000 != 0 {
		return
	}
	x := float32(float64(t.PosVec.X) - float64(u.PosVec.X))
	yd := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	y := float32(yd)
	length := math.Sqrt(yd*float64(y)+float64(x)*float64(x)) + .1
	spilled := float32(length)
	if !(length < 400) {
		return
	}
	dx, dy := worldIndexedDirection(u)
	ax, ay := float64(x), float64(y)
	if ax < 0 {
		ax = -ax
	}
	if ay < 0 {
		ay = -ay
	}
	ratio := ay / ax
	switch dx + 3*dy + 4 {
	case 0:
		if !(x < 0 && y < 0 && ratio >= .57730001 && ratio <= .1732) {
			return
		}
	case 1:
		if !(y < 0 && ratio <= .3732) {
			return
		}
	case 2:
		if !(x > 0 && y < 0 && ratio >= .57730001 && ratio <= .1732) {
			return
		}
	case 3:
		if !(x < 0 && ratio <= .26789999) {
			return
		}
	case 5:
		if !(x > 0 && ratio <= .26789999) {
			return
		}
	case 6:
		if !(x < 0 && y > 0 && ratio >= .57730001 && ratio <= .1732) {
			return
		}
	case 7:
		if !(y > 0 && ratio <= .3732) {
			return
		}
	case 8:
		if !(x > 0 && y > 0 && ratio >= .57730001 && ratio <= .1732) {
			return
		}
	}
	if !GetServer().S().CanInteract(u, t, 0) {
		return
	}
	d := 400 - float64(spilled)
	strength := float32(d * d * d * .0000005)
	mass := float64(t.Mass)
	force := float64(strength) / mass
	vx, vy := movementDirectionVector(int32(int16(u.Direction1)))
	t.ForceVec.X = float32(force*float64(vx) + float64(t.ForceVec.X))
	t.ForceVec.Y = float32(force*float64(vy) + float64(t.ForceVec.Y))
}
func worldBlow(u *server.Object) {
	if u.ObjFlags&0x1000000 == 0 {
		return
	}
	dx, dy := worldIndexedDirection(u)
	x, y := float64(u.PosVec.X), float64(u.PosVec.Y)
	rect := types.Rectf{Min: u.PosVec, Max: u.PosVec}
	if dx >= 0 {
		rect.Max.X = float32(x + 400)
	} else {
		rect.Min.X = float32(x - 400)
	}
	if dy >= 0 {
		rect.Max.Y = float32(y + 400)
	} else {
		rect.Min.Y = float32(y - 400)
	}
	GetServer().S().Map.EachObjAndMissileInRect(rect, func(t *server.Object) bool { worldBlowCandidate(t, u); return true })
}
