package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func collisionInitTypes() {
	types := &GetServer().S().Types
	collisionSmallFist = uint32(types.IndByID("SmallFist"))
	collisionMediumFist = uint32(types.IndByID("MediumFist"))
	collisionLargeFist = uint32(types.IndByID("LargeFist"))
	collisionMeteor = uint32(types.IndByID("Meteor"))
	collisionTypesReady = 1
}
func collisionHeightExcluded(u *server.Object) bool {
	v := uint32(u.TypeInd)
	return v == collisionSmallFist || v == collisionMediumFist || v == collisionLargeFist || v == collisionMeteor
}
func collisionElevator(e, u *server.Object, mode int32) {
	ud := e.UpdateData
	if collisionTypesReady == 0 {
		collisionInitTypes()
	}
	if mode == 0 || collisionHeightExcluded(u) {
		return
	}
	height := float64(int32(layout.As[server.ElevatorUpdateData](ud).Field_4))
	delta := float32(float64(u.ZVal) - height)
	if float64(worldAbsScratch(float32(delta))) > 10 {
		if height > float64(u.ZVal) {
			if u.Shape.Kind == server.ShapeKindCircle {
				collisionCircleBox(u, e, 0)
			} else if u.Shape.Kind == server.ShapeKindBox {
				geometryBoxBox(u, e)
			}
		}
	} else if collisionContains(&e.NewPos, (*[11]float32)(unsafe.Pointer(&e.Shape)), &u.NewPos) {
		u.ObjFlags = u.ObjFlags&^0x40000 | 0x100000
		stateRaise(u, float32(height+4))
		u.Field27 = 0
	}
}
func collisionShaft(e, u *server.Object) {
	if collisionTypesReady == 0 {
		collisionInitTypes()
	}
	linked := e.UpdateDataElevatorShaft().Partner // the elevator (see worldShaft)
	if linked == nil || collisionHeightExcluded(u) {
		return
	}
	if u.Shape.Kind == server.ShapeKindBox {
		if e.Shape.Box.W < u.Shape.Box.W || e.Shape.Box.H < u.Shape.Box.H {
			return
		}
	} else if u.Shape.Kind == server.ShapeKindCircle {
		diameter := float64(u.Shape.Circle.R) + float64(u.Shape.Circle.R)
		if diameter > float64(e.Shape.Box.W) || diameter > float64(e.Shape.Box.H) {
			return
		}
	}
	if !collisionContains(&e.NewPos, (*[11]float32)(unsafe.Pointer(&e.Shape)), &u.NewPos) {
		return
	}
	heightWide := float64(int32(layout.As[server.ElevatorUpdateData](linked.UpdateData).Field_4 - 64))
	height := float32(heightWide)
	// C spills height for later use, but subtracts the original wide value.
	delta := float32(float64(u.ZVal) - heightWide)
	if float64(worldAbsScratch(float32(delta))) > 10 {
		if height <= -10 {
			u.ObjFlags |= 0x40000
			u.Pos39 = e.NewPos
			u.Field41 = math.Float32bits(linked.NewPos.X)
			u.Field42 = math.Float32bits(linked.NewPos.Y)
		}
	} else {
		u.ObjFlags = u.ObjFlags&^0x40000 | 0x100000
		stateRaise(u, float32(float64(height)+4))
		u.Field27 = 0
	}
}
func collisionCircleWalls(u *server.Object) {
	if u.ObjClass&0x400000 != 0 {
		return
	}
	x1 := floatToInt32(float32(float64(u.CollideP1.X) * 0.043478262))
	y1 := floatToInt32(float32(float64(u.CollideP1.Y) * 0.043478262))
	x2 := floatToInt32(float32(float64(u.CollideP2.X) * 0.043478262))
	y2 := floatToInt32(float32(float64(u.CollideP2.Y) * 0.043478262))
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			grid := [2]int32{x, y}
			if geometryCircleWall(&grid, u) != 0 {
				collisionWallOpen(&grid, u)
			}
		}
	}
}
func collisionDrainAngles() {
	for p := collisionAngleHead; p != 0; {
		ud := (*server.DoorUpdateData)(memmap.WordPtr(p))
		torque := ud.Torque
		velocity := floatToInt32(float32(float64(torque) + float64(torque) + 0.5))
		if velocity < 0 {
			velocity = -velocity
		}
		if velocity > 4 {
			velocity = 4
		}
		angle := &ud.Angle
		base := int32(ud.BaseDir)
		current := int32(ud.Dir)
		if !(float64(torque) >= -0.0099999998) {
			stop := base - 12
			if stop < 0 {
				stop += 32
			}
			if current != stop {
				*angle -= int16(velocity)
				if *angle < 0 {
					*angle += 256
				}
			}
		} else if !(float64(torque) <= 0.0099999998) && uint32(current) != (uint32(base)+12)%32 {
			*angle += int16(velocity)
			if *angle >= 256 {
				*angle -= 256
			}
		}
		index := 32 * int32(*angle) / 256
		for index < 0 {
			index += 32
		}
		for index >= 32 {
			index -= 32
		}
		ud.Dir = uint32(index)
		ud.AngleQueued = 0
		ud.Torque = 0
		p = ud.AngleNext
	}
	collisionAngleHead = 0
}
