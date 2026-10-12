package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func collisionGateCircle(gate, u *server.Object, mode int32) {
	ud := gate.UpdateData
	direction := int32(*(*int16)(unsafe.Add(ud, 40)))
	index := direction + 160
	if index >= 256 {
		index = direction - 96
	}
	ix, iy := memmap.Int32(0x587000, uintptr(192088+8*index)), memmap.Int32(0x587000, uintptr(192092+8*index))
	line := [4]float32{gate.NewPos.X, gate.NewPos.Y, float32(float64(2*ix) + float64(gate.NewPos.X)), float32(float64(2*iy) + float64(gate.NewPos.Y))}
	var closest types.Pointf
	projectLineClamped(&line, &u.NewPos, &closest, 32)
	dx, dy := float64(u.NewPos.X)-float64(closest.X), float64(u.NewPos.Y)-float64(closest.Y)
	wide := math.Sqrt(dy*dy + dx*dx)
	distance := float32(wide)
	if wide == 0 {
		distance = 0.1
	}
	if !(distance < u.Shape.Circle.R) {
		return
	}
	normal := types.Pointf{float32(dx / float64(distance)), float32(dy / float64(distance))}
	depth := float32(float64(u.Shape.Circle.R) - float64(distance))
	geometryHit(u, gate, &normal)
	s := GetServer().S()
	*equipmentWord(ud, 44) = s.Frame()
	if mode == 1 {
		velocity := float32(-float64(normal.Y)*float64(u.VelVec.Y) - float64(normal.X)*float64(u.VelVec.X))
		damp := math.Sqrt(float64(u.Mass) * float64(geometryWallForce()) * 4)
		force := damp*float64(velocity)*0.25 + float64(depth)*float64(geometryWallForce())
		u.Sub548600(types.Pointf{float32(force * float64(normal.X)), float32(force * float64(normal.Y))})
	}
	if u.ObjFlags&0x8000000 != 0 {
		if u.ObjFlags&8 == 0 {
			geometryActivate(u)
		}
		u.ObjFlags &^= 0x8000000
	}
	geometryActivate(gate)
	if !gate.TeamPtr().Has() || *equipmentWord(ud, 12) != *equipmentWord(ud, 4) || gate.TeamPtr().SameAs(u.TeamPtr()) {
		if mode == 0 && *(*byte)(unsafe.Add(ud, 1)) == 0 && (gate.ObjOwner == nil || gate.ObjOwner == u) {
			torque := (float64(u.Shape.Circle.R) - float64(distance)) * float64(u.Mass)
			angular := (*float32)(unsafe.Add(ud, 32))
			cross := float64(memmap.Int32(0x587000, uintptr(192088+8*index)))*(float64(u.NewPos.Y)-float64(gate.NewPos.Y)) - float64(memmap.Int32(0x587000, uintptr(192092+8*index)))*(float64(u.NewPos.X)-float64(gate.NewPos.X))
			if cross <= 0 {
				*angular = float32(torque + float64(*angular))
			} else {
				*angular = float32(float64(*angular) - torque)
			}
			worldAngleQueue(ud)
			s.Objs.AddToUpdatable(gate)
		}
	} else if s.Frame() > gate.Field34 {
		gate.Field34 = s.Frame() + s.TickRate()
		_ = s.Teams.ByID(gate.TeamVal.ID)
		gameplayTextPrivate(u, (*byte)(unsafe.Pointer(internCStr("objcoll.c:GateLockedMechanism"))), 0)
	}
}
