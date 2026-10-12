package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

var motionVelocityTypes [10]uint32

func motionProjectileStep(u *server.Object) {
	vx := (float64(u.ForceVec.X) + float64(u.VelVec.X)) * float64(u.Float28)
	vy := (float64(u.ForceVec.Y) + float64(u.VelVec.Y)) * float64(u.Float28)
	u.PosVec = u.NewPos
	u.VelVec = types.Pointf{X: float32(vx), Y: float32(vy)}
	u.NewPos.X = float32(vx + float64(u.NewPos.X))
	u.NewPos.Y = float32(float64(vy) + float64(u.NewPos.Y))
	u.Nox_xxx_objectUnkUpdateCoords_4E7290()
}
func motionFall(u *server.Object) {
	flags, z := u.ObjFlags, u.ZVal
	if flags&0x40000 != 0 {
		stateRaise(u, float32(float64(z)+float64(u.Field27)))
		u.Field27 = float32(float64(u.Field27) - 1)
		u.ForceVec, u.VelVec = types.Pointf{}, types.Pointf{}
		dx := float64(u.PosVec.X) - float64(u.Pos39.X)
		dyWide := float64(u.PosVec.Y) - float64(u.Pos39.Y)
		dy := dyWide
		length := math.Sqrt(dyWide*float64(dy) + dx*dx)
		if length > 0 {
			den := float64(float32(length))
			u.ForceVec = types.Pointf{X: float32(dx * -3 / den), Y: float32(float64(dy) * -3 / den)}
		}
		if u.ZVal < -50 {
			stateRaise(u, 90)
			u.ObjFlags &^= 0x40000
			Nox_xxx_unitMove_4E7010(u, types.Pointf{X: math.Float32frombits(u.Field41), Y: math.Float32frombits(u.Field42)})
		}
	} else if z != 0 || u.Field27 != 0 {
		if flags&0x800000 != 0 {
			stateRaise(u, float32(float64(u.ZVal)+float64(u.Field27)))
			if u.ZVal >= 0 {
				u.Field27 = float32(float64(u.Field27) - .5)
			} else {
				stateRaise(u, 0)
				bounce := -float64(u.Field27) * float64(math.Float32frombits(u.Field29)) * .1
				u.Field27 = float32(bounce)
				if bounce < 2 {
					stateRaise(u, 0)
					u.Field27 = 0
				}
			}
		} else if flags&0x100000 == 0 {
			if u.ZVal > 0 {
				if u.Field27 <= 0 {
					u.ObjFlags = flags | 0x20000
				}
				stateRaise(u, float32(float64(u.ZVal)+float64(u.Field27)))
				u.Field27 = float32(float64(u.Field27) - 1)
			}
			if u.ZVal <= 0 {
				down := u.Field27
				u.ObjFlags &^= 0x20000
				if down < 0 && u.ObjClass&1 == 0 {
					collisionActivate(u)
					if u.Field27 < -10 && u.ObjClass&4 != 0 {
						GetServer().S().Audio.EventObj(280, u, 0, 0)
					}
				}
				stateRaise(u, 0)
				u.Field27 = 0
			}
		}
	}
}
func motionVelocity(step float32) int32 {
	s := GetServer().S()
	if motionVelocityTypes[0] == 0 {
		for i, name := range [...]string{"SmallFlameCleanse", "SmallFlameCleanse", "MediumFlameCleanse", "FlameCleanse", "LargeFlameCleanse", "SmallBlueFlameCleanse", "SmallBlueFlameCleanse", "MediumBlueFlameCleanse", "BlueFlameCleanse", "LargeBlueFlameCleanse"} {
			motionVelocityTypes[i] = uint32(s.Types.IndByID(name))
		}
	}
	for p := collisionActiveHead; p != 0; {
		u := motionObject(p)
		collisionScan(u)
		p = collisionNextActive(u)
	}
	Nox_xxx_updateSprings_5113A0()
	for p := collisionActiveHead; p != 0; {
		u := motionObject(p)
		if u.ObjFlags&2 != 0 || u.ObjClass&2 != 0 && u.UpdateDataMonster().HasAction(ai.ActionType(67)) {
			u.Pos24, u.VelVec = types.Pointf{}, types.Pointf{}
		} else {
			fx, fy := float64(u.Pos24.X), float64(u.Pos24.Y)
			if u.Buffs&((1<<5)|(1<<25)|(1<<28)) == 0 {
				fx += float64(u.ForceVec.X)
				fy += float64(u.ForceVec.Y)
			}
			u.VelVec.X = float32(float64(u.VelVec.X) + (fx-float64(u.VelVec.X)*float64(u.Float28))*float64(step))
			u.VelVec.Y = float32(float64(u.VelVec.Y) + (float64(fy)-float64(u.VelVec.Y)*float64(u.Float28))*float64(step))
			next := types.Pointf{X: float32(float64(step)*float64(u.VelVec.X) + float64(u.NewPos.X)), Y: float32(float64(step)*float64(u.VelVec.Y) + float64(u.NewPos.Y))}
			flags := server.MapTraceFlags((uint32(u.ObjFlags)>>12)&4 | 1)
			for _, id := range motionVelocityTypes {
				if uint32(u.TypeInd) == id {
					flags |= 0x40
					break
				}
			}
			if s.MapTraceRayAt(u.NewPos, next, nil, nil, flags) {
				u.NewPos = next
			}
			if u.ObjFlags&0x4000 == 0 && u.HealthData != nil && tileAtPoint(u.NewPos) == 6 {
				normal := types.Pointf{}
				collisionAddHit(u, 6, &normal)
			}
			if math.Abs(float64(u.NewPos.X)-float64(u.PosVec.X)) > .0099999998 || float64(float32(math.Abs(float64(u.NewPos.Y)-float64(u.PosVec.Y)))) > .0099999998 {
				u.NeedSync()
				u.Nox_xxx_objectUnkUpdateCoords_4E7290()
				Nox_xxx_moveUpdateSpecial_517970(u)
			}
		}
		p = collisionNextActive(u)
	}
	return 0
}
