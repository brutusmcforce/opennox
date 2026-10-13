package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
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
	if !(float64(radius) >= float64(layout.As[server.MonsterUpdateData](ud).SightRange)) {
		radius = layout.As[server.MonsterUpdateData](ud).SightRange
	}
	refresh := s.Frame()-layout.As[server.MonsterUpdateData](ud).Field303 > 2*s.TickRate()
	if refresh {
		layout.As[server.MonsterUpdateData](ud).Field303 = s.Frame()
	}
	changed := false
	slots := visibilitySlots(ud)
	for i := 0; i < int(layout.As[server.MonsterUpdateData](ud).Field282_1); i++ {
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
	target := layout.As[server.MonsterUpdateData](ud).CurrentEnemy
	if target != nil && target.HasEnchant(server.EnchantID(28)) {
		changed = true
	}
	if (target == nil || s.Frame()-layout.As[server.MonsterUpdateData](ud).Field301 > 2*s.TickRate()) && (layout.As[server.MonsterUpdateData](ud).Field302 <= s.Frame() || s.Frame() == *memmap.PtrUint32(0x5D4594, 2487684)) {
		s.Map.EachObjInCircle(u.PosVec, radius, func(it *server.Object) bool { visibilityCandidate(it, u); return true })
		layout.As[server.MonsterUpdateData](ud).Field301 = s.Frame()
		layout.As[server.MonsterUpdateData](ud).Field303 = s.Frame()
		changed = true
	}
	if changed {
		prior := uint32(0)
		if v := layout.As[server.MonsterUpdateData](ud).CurrentEnemy; v != nil {
			prior = v.NetCode
		}
		visibilitySelectTarget(u)
		if v := layout.As[server.MonsterUpdateData](ud).CurrentEnemy; v != nil && prior != 0 && prior != v.NetCode {
			layout.As[server.MonsterUpdateData](ud).Field300 = prior
		}
	}
	if layout.As[server.MonsterUpdateData](ud).Field301 == s.Frame() {
		if *(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).StatusFlags))&0x400 != 0 || controlFlags(0x2000) || layout.As[server.MonsterUpdateData](ud).CurrentEnemy != nil {
			layout.As[server.MonsterUpdateData](ud).Field302 = s.Frame() + uint32(s.Rand.Logic.IntClamp(5, 10))
		} else {
			duration := int32(5 * s.TickRate())
			distance := float32(s.Sub5336D0(u))
			*(*float32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).Field131)) = distance
			if distance < 0 {
				layout.As[server.MonsterUpdateData](ud).Field302 = uint32(duration) + s.Frame()
			} else if float64(distance) > float64(radius) {
				layout.As[server.MonsterUpdateData](ud).Field302 = uint32(int64((float64(distance)-float64(radius))*float64(duration)/(1000-float64(radius)))) + 10 + s.Frame()
			} else {
				layout.As[server.MonsterUpdateData](ud).Field302 = uint32(s.Rand.Logic.IntClamp(5, 10)) + s.Frame()
			}
		}
	}
}
