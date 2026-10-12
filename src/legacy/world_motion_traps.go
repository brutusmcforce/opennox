package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
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
	d := (*[3]uint32)(u.UpdateData)
	state := (*byte)(unsafe.Add(u.UpdateData, 8))
	if motionTrapScan(u) != 0 {
		if *state != 1 {
			*state = 1
			d[1] = 0
		}
	} else {
		*state = 0
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
		d := (*[2]int32)(arrow.CollideData)
		d[0], d[1] = damage, damage
	}
	if uint32(typ) == memmap.Uint32(0x5D4594, 2491768) {
		s.Audio.EventObj(889, u, 0, 0)
	}
}
func motionTrapUpdate(u *server.Object) int32 {
	result := uint32(u.ObjFlags)
	d := (*[13]uint32)(u.UpdateData)
	on := (*byte)(unsafe.Add(u.UpdateData, 48))
	if result&0x1000000 != 0 {
		if *on == 0 {
			d[0], d[1] = 0, 0
		}
		scan := d[0]
		*on = 1
		if scan == 0 {
			motionTrapState(u)
			d[0] = GetServer().S().TickRate()
		}
		if byte(d[2]) == 1 && d[1] == 0 {
			if memmap.Uint32(0x5D4594, 2491780) == 0 {
				*memmap.PtrUint32(0x5D4594, 2491780) = uint32(GetServer().S().Types.IndByID("ArrowTrap1"))
				*memmap.PtrUint32(0x5D4594, 2491784) = uint32(GetServer().S().Types.IndByID("ArrowTrap2"))
			}
			motionTrapProjectile(u, int32(d[3]))
			if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491780) {
				visibilityFXArrowTrap(u.PosVec, 1)
			} else if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2491784) {
				visibilityFXArrowTrap(u.PosVec, 2)
			}
			d[1] = 30
		}
		if d[0] != 0 {
			d[0]--
		}
		result = d[1]
		if result != 0 {
			result--
			d[1] = result
		}
	} else {
		*on = 0
	}
	return int32(result)
}
func motionTrigger(u, target *server.Object) {
	d := (*[14]uint32)(u.UpdateData)
	if u.ObjFlags&0x1000000 == 0 {
		d[1] = 0
		d[0] &^= 1
		return
	}
	if byte(d[2]) == 5 || target == nil || !(target.Mass > 0) {
		return
	}
	if d[11] != 0 && d[11]&uint32(target.ObjClass) == 0 || d[12] != 0 && d[12]&uint32(target.ObjClass) != 0 {
		return
	}
	allow := int8(byte(d[13]))
	deny := int8(byte(d[13] >> 8))
	// C promotes the signed allow byte against an unsigned object team byte.
	if allow != 0 && int32(target.TeamVal.ID) != int32(allow) || deny != 0 && int8(target.TeamVal.ID) == deny {
		return
	}
	if d[4] != 0xffffffff {
		result := GetServer().NoxScriptC().ScriptCallback((*server.ScriptCallback)(unsafe.Add(u.UpdateData, 12)), target, u, 1)
		if result == nil || *(*uint32)(result) == 0 {
			return
		}
	}
	d[1] = motionAddress(target)
	d[0] |= 1
}
