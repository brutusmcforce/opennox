package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func visibilityUpdateSight(u *server.Object) {
	ud := u.UpdateData
	s := GetServer().S()
	if uint32(u.ObjFlags)&0x8000 != 0 && !monsterIsZombie(u) {
		return
	}
	radius := float32(250)
	if controlFlags(4096) {
		radius = 640
	}
	if !(float64(radius) >= float64(*visibilityFloat(ud, 1312))) {
		radius = *visibilityFloat(ud, 1312)
	}
	refresh := s.Frame()-*visibilityWord(ud, 1212) > 2*s.TickRate()
	if refresh {
		*visibilityWord(ud, 1212) = s.Frame()
	}
	changed := false
	slots := visibilitySlots(ud)
	for i := 0; i < int(*controlByte(ud, 1129)); i++ {
		v := slots[i]
		lost := uint32(v.ObjFlags)&0x8020 != 0 || !s.CanSee(u, v, 0)
		if !lost {
			dx, dy := float64(u.PosVec.X)-float64(v.PosVec.X), float64(u.PosVec.Y)-float64(v.PosVec.Y)
			limit := float32((float64(radius) + 30) * (float64(radius) + 30))
			lost = dy*dy+dx*dx > float64(limit)
			if !lost {
				dx, dy = float64(u.PosVec.X)-float64(u.PrevPos.X), float64(u.PosVec.Y)-float64(u.PrevPos.Y)
				lost = dy*dy+dx*dx > 1000
			}
			if !lost && refresh {
				lost = !s.CanInteract(u, v, 0)
			}
		}
		if lost {
			visibilityLost(u, i)
			i--
			changed = true
		}
	}
	target := *visibilityObject(ud, 1196)
	if target != nil && target.HasEnchant(server.EnchantID(28)) {
		changed = true
	}
	if (target == nil || s.Frame()-*visibilityWord(ud, 1204) > 2*s.TickRate()) && (*visibilityWord(ud, 1208) <= s.Frame() || s.Frame() == *memmap.PtrUint32(0x5D4594, 2487684)) {
		s.Map.EachObjInCircle(u.PosVec, radius, func(it *server.Object) bool { visibilityCandidate(it, u); return true })
		*visibilityWord(ud, 1204) = s.Frame()
		*visibilityWord(ud, 1212) = s.Frame()
		changed = true
	}
	if changed {
		prior := uint32(0)
		if v := *visibilityObject(ud, 1196); v != nil {
			prior = v.NetCode
		}
		visibilitySelectTarget(u)
		if v := *visibilityObject(ud, 1196); v != nil && prior != 0 && prior != v.NetCode {
			*visibilityWord(ud, 1200) = prior
		}
	}
	if *visibilityWord(ud, 1204) == s.Frame() {
		if *visibilityWord(ud, 1440)&0x400 != 0 || controlFlags(0x2000) || *visibilityObject(ud, 1196) != nil {
			*visibilityWord(ud, 1208) = s.Frame() + uint32(s.Rand.Logic.IntClamp(5, 10))
		} else {
			duration := int32(5 * s.TickRate())
			distance := float32(s.Sub5336D0(u))
			*visibilityFloat(ud, 524) = distance
			if distance < 0 {
				*visibilityWord(ud, 1208) = uint32(duration) + s.Frame()
			} else if float64(distance) > float64(radius) {
				*visibilityWord(ud, 1208) = uint32(int64((float64(distance)-float64(radius))*float64(duration)/(1000-float64(radius)))) + 10 + s.Frame()
			} else {
				*visibilityWord(ud, 1208) = uint32(s.Rand.Logic.IntClamp(5, 10)) + s.Frame()
			}
		}
	}
}
