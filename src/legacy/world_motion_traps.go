package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func motionTrapCandidate(target, trap *server.Object) {
	s := GetServer().S()
	if target.ObjClass&6 != 0 && target.ObjFlags&0x8020 == 0 && s.IsEnemyTo(trap, target) && s.MapTraceVision(target, trap) {
		if stateFront(&trap.PosVec, int32(int16(trap.Direction1)), &target.PosVec)&1 != 0 {
			*memmap.PtrUint32(0x5D4594, 2491764) = 1
		}
	}
}
func motionTrapScan(u *server.Object) int32 {
	p := u.PosVec
	r := types.Rectf{Min: types.Pointf{X: float32(float64(p.X) - 350), Y: float32(float64(p.Y) - 350)}, Max: types.Pointf{X: float32(float64(p.X) + 350), Y: float32(float64(p.Y) + 350)}}
	*memmap.PtrUint32(0x5D4594, 2491764) = 0
	GetServer().S().Map.EachObjInRect(r, func(t *server.Object) bool { motionTrapCandidate(t, u); return true })
	return int32(memmap.Uint32(0x5D4594, 2491764))
}
func motionTrapState(u *server.Object) {
	d := u.UpdateDataSkull()
	if motionTrapScan(u) != 0 {
		if d.State != 1 {
			d.State = 1
			d.Cooldown = 0
		}
	} else {
		d.State = 0
	}
}
func motionTrapProjectile(u *server.Object, typ int32) {
	s := GetServer().S()
	radius := float64(u.Shape.Circle.R) + 4
	dir := int32(int16(u.Direction1))
	pos := types.Pointf{X: float32(radius*float64(memmap.Float32(0x587000, uintptr(194136+8*dir))) + float64(u.PosVec.X)), Y: float32(radius*float64(memmap.Float32(0x587000, uintptr(194140+8*dir))) + float64(u.PosVec.Y))}
	arrow := s.NewObjectByTypeInd(int(typ))
	if arrow == nil {
		return
	}
	GetServer().CreateObjectAt(arrow, u, pos)
	arrow.Direction1, arrow.Direction2 = u.Direction1, u.Direction1
	dir = int32(int16(u.Direction1))
	arrow.VelVec.X = float32(float64(memmap.Float32(0x587000, uintptr(194136+8*dir))) * float64(arrow.SpeedCur))
	arrow.VelVec.Y = float32(float64(memmap.Float32(0x587000, uintptr(194140+8*dir))) * float64(arrow.SpeedCur))
	if memmap.Uint32(0x5D4594, 2491768) == 0 {
		*memmap.PtrUint32(0x5D4594, 2491768) = uint32(s.Types.IndByID("MercArcherArrow"))
		*memmap.PtrUint32(0x5D4594, 2491772) = uint32(s.Types.IndByID("ArrowTrap1"))
		*memmap.PtrUint32(0x5D4594, 2491776) = uint32(s.Types.IndByID("ArrowTrap2"))
	}
	if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491772) || uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491776) {
		damage := floatToInt32(float32(s.Balance.Float("ArrowTrapDamage")))
		d := arrow.CollideDataMonsterArrow()
		d.DamageCoop, d.Damage = uint32(damage), uint32(damage)
	}
	if uint32(typ) == memmap.Uint32(0x5D4594, 2491768) {
		s.Audio.EventObj(889, u, 0, 0)
	}
}
func motionTrapUpdate(u *server.Object) int32 {
	result := uint32(u.ObjFlags)
	d := u.UpdateDataSkull()
	if result&0x1000000 != 0 {
		if d.Active == 0 {
			d.ScanDelay, d.Cooldown = 0, 0
		}
		scan := d.ScanDelay
		d.Active = 1
		if scan == 0 {
			motionTrapState(u)
			d.ScanDelay = GetServer().S().TickRate()
		}
		if d.State == 1 && d.Cooldown == 0 {
			if memmap.Uint32(0x5D4594, 2491780) == 0 {
				*memmap.PtrUint32(0x5D4594, 2491780) = uint32(GetServer().S().Types.IndByID("ArrowTrap1"))
				*memmap.PtrUint32(0x5D4594, 2491784) = uint32(GetServer().S().Types.IndByID("ArrowTrap2"))
			}
			motionTrapProjectile(u, int32(d.MissileType))
			if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491780) {
				visibilityFXArrowTrap(u.PosVec, 1)
			} else if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491784) {
				visibilityFXArrowTrap(u.PosVec, 2)
			}
			d.Cooldown = 30
		}
		if d.ScanDelay != 0 {
			d.ScanDelay--
		}
		result = d.Cooldown
		if result != 0 {
			result--
			d.Cooldown = result
		}
	} else {
		d.Active = 0
	}
	return int32(result)
}
func motionTrigger(u, target *server.Object) {
	d := u.UpdateDataTrigger()
	if u.ObjFlags&0x1000000 == 0 {
		d.Activator = 0
		d.Flags &^= 1
		return
	}
	if d.State == 5 || target == nil || !(target.Mass > 0) {
		return
	}
	if d.AllowClass != 0 && d.AllowClass&uint32(target.ObjClass) == 0 || d.IgnoreClass != 0 && d.IgnoreClass&uint32(target.ObjClass) != 0 {
		return
	}
	allow := d.AllowTeam
	deny := d.IgnoreTeam
	// C promotes the signed allow byte against an unsigned object team byte.
	if allow != 0 && int32(target.TeamVal.ID) != int32(allow) || deny != 0 && int8(target.TeamVal.ID) == deny {
		return
	}
	if d.ScriptCollide.Func != -1 {
		result := GetServer().NoxScriptC().ScriptCallback(&d.ScriptCollide, target, u, 1)
		if result == nil || *(*uint32)(result) == 0 {
			return
		}
	}
	d.Activator = motionAddress(target)
	d.Flags |= 1
}
