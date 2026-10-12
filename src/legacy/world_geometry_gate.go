package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func geometryGateBox(gate, u *server.Object, mode int32) {
	ud := gate.UpdateData
	originX := (float64(gate.NewPos.Y) + float64(gate.NewPos.X)) * 0.70710677
	originY := float32((float64(gate.NewPos.Y) - float64(gate.NewPos.X)) * 0.70710677)
	center := geometryRotated(u.NewPos)
	rect := geometryRotatedBounds(u, center)
	direction := int32(*(*int16)(unsafe.Add(ud, 40)))
	index := direction + 128
	if index >= 256 {
		index = direction - 128
	}
	line := [4]float32{float32(originX), originY,
		float32(float64(memmap.Float32(0x587000, uintptr(194136+8*index)))*32 + originX),
		float32(float64(memmap.Float32(0x587000, uintptr(194140+8*index)))*32 + float64(originY))}
	var bounds [4]float32
	if originX >= float64(line[2]) {
		bounds[2] = float32(originX)
		bounds[0] = line[2]
	} else {
		bounds[0] = float32(originX)
		bounds[2] = line[2]
	}
	if originY >= line[3] {
		bounds[1] = line[3]
		bounds[3] = originY
	} else {
		bounds[3] = line[3]
		bounds[1] = originY
	}
	if !(bounds[0] <= rect[2] && bounds[1] <= rect[3] && bounds[2] >= rect[0] && bounds[3] >= rect[1]) {
		return
	}
	var midpoint types.Pointf
	if geometryClipCenter(&line, &bounds, &rect, &midpoint) == 0 {
		return
	}
	dx, dy := float64(center.X)-float64(midpoint.X), float64(center.Y)-float64(midpoint.Y)

	distance64 := math.Sqrt(dy*dy + dx*dx)
	distance := float32(distance64)
	if distance64 == 0 {
		return
	}
	ray := [4]float32{center.X, center.Y, midpoint.X, midpoint.Y}
	normal := types.Pointf{float32(dx / float64(distance)), float32(dy / float64(distance))}
	var crossing types.Pointf
	if geometryRectCrossings(&ray, &rect, &crossing, 1, 1) != 1 {
		return
	}
	cy, cx := float64(crossing.Y)-float64(center.Y), float64(crossing.X)-float64(center.X)
	radius := math.Sqrt(cy*cy + cx*cx)
	if radius == 0 {
		return
	}
	nx := (float64(normal.X) - float64(normal.Y)) * 0.70710677
	normal.Y = float32((float64(normal.Y) + float64(normal.X)) * 0.70710677)
	normal.X = float32(nx)
	depth64 := radius - float64(distance)
	depth := float32(depth64)
	if !(depth64 > 0) {
		return
	}
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
			torqueIndex := index + 32
			if torqueIndex >= 256 {
				torqueIndex -= 256
			}
			torque := float64(depth) * float64(u.Mass)
			angular := (*float32)(unsafe.Add(ud, 32))
			cross := float64(memmap.Int32(0x587000, uintptr(192088+8*torqueIndex)))*(float64(u.NewPos.Y)-float64(gate.NewPos.Y)) - float64(memmap.Int32(0x587000, uintptr(192092+8*torqueIndex)))*(float64(u.NewPos.X)-float64(gate.NewPos.X))
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
