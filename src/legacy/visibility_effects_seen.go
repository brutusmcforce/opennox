package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
)

func visibilityWord(p unsafe.Pointer, off int) *uint32   { return rawAt[uint32](p, off) }
func visibilityFloat(p unsafe.Pointer, off int) *float32 { return rawAt[float32](p, off) }
func visibilityObject(p unsafe.Pointer, off int) **server.Object {
	return (**server.Object)(unsafe.Add(p, off))
}
func visibilitySlots(p unsafe.Pointer) []*server.Object {
	return unsafe.Slice((**server.Object)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field283)), 16)
}
func visibilityKillable(u *server.Object) int {
	h := u.HealthData
	return bool2int(h != nil && (h.Cur != 0 || h.Max == 0))
}
func visibilityFrameCopy(next bool) uint32 {
	frame := GetServer().S().Frame()
	if next {
		frame++
	}
	*memmap.PtrUint32(0x5D4594, 2487684) = frame
	return frame
}
func visibilityLost(u *server.Object, index int) int {
	ud := u.UpdateData
	slots := visibilitySlots(ud)
	v := slots[index]
	unitDebug(0, GetServer().S().Frame(), v.NetCode, GetServer().S().Types.ByInd(int(v.TypeInd)).ID())
	GetServer().NoxScriptC().ScriptCallback(&layout.As[server.MonsterUpdateData](ud).ScriptLostEnemy, slots[index], u, server.ScriptEventType(15))
	selected := &layout.As[server.MonsterUpdateData](ud).CurrentEnemy
	if slots[index] == *selected {
		layout.As[server.MonsterUpdateData](ud).Field300 = (*selected).NetCode
		*selected = nil
	}
	count := &layout.As[server.MonsterUpdateData](ud).Field282_1
	*count--
	result := int(*count)
	if index < int(*count) {
		result = int(uintptr(unsafe.Pointer(&slots[index]))) // the legacy return value is the address of the slot
		for i := index; i < int(*count); i++ {
			slots[i] = slots[i+1]
			result += 4
		}
	}
	return result
}
func visibilityRemove(u, v *server.Object) int {
	ud := u.UpdateData
	count := int(layout.As[server.MonsterUpdateData](ud).Field282_1)
	slots := visibilitySlots(ud)
	for i := 0; i < count; i++ {
		if slots[i] == v {
			return visibilityLost(u, i)
		}
	}
	return count
}
func visibilityContains(u, v *server.Object) bool {
	ud := u.UpdateData
	for i, it := range visibilitySlots(ud) {
		if i >= int(layout.As[server.MonsterUpdateData](ud).Field282_1) {
			break
		}
		if it == v {
			return true
		}
	}
	return false
}
func visibilitySelectTarget(u *server.Object) {
	ud := u.UpdateData
	selected := &layout.As[server.MonsterUpdateData](ud).CurrentEnemy
	*selected = nil
	best := float32(100000000)
	slots := visibilitySlots(ud)
	for i := 0; i < int(layout.As[server.MonsterUpdateData](ud).Field282_1); i++ {
		it := slots[i]
		if it == *(**server.Object)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).Field304)) {
			*selected = it
			return
		}
		if GetServer().S().IsEnemyTo(u, it) && visibilityKillable(it) != 0 {
			dx, dy := float64(it.PosVec.X)-float64(u.PosVec.X), float64(it.PosVec.Y)-float64(u.PosVec.Y)
			dist := dy*dy + dx*dx
			if dist < float64(best) {
				best = float32(dist)
				*selected = it
			}
		}
	}
}
func visibilitySee(u, v *server.Object) {
	ud := u.UpdateData
	count := &layout.As[server.MonsterUpdateData](ud).Field282_1
	slots := visibilitySlots(ud)
	if *count == 16 {
		farthest := float64(0)
		var far *server.Object
		for _, it := range slots {
			dx, dy := float64(it.PosVec.X)-float64(u.PosVec.X), float64(it.PosVec.Y)-float64(u.PosVec.Y)
			dist := dy*dy + dx*dx
			if dist > farthest {
				farthest = float64(float32(dist))
				far = it
			}
		}
		dx, dy := float64(v.PosVec.X)-float64(u.PosVec.X), float64(v.PosVec.Y)-float64(u.PosVec.Y)
		if farthest <= dy*dy+dx*dx {
			return
		}
		visibilityRemove(u, far)
	}
	slots[int(*count)] = v
	*count++
	s := GetServer().S()
	cooldown := &layout.As[server.MonsterUpdateData](ud).Field134
	if s.Frame() > *cooldown && s.IsEnemyTo(u, v) && (!monsterIsZombie(u) || uint32(u.ObjFlags)&0x8000 == 0) {
		if set := Nox_xxx_monsterGetSoundSet_424300(u); set != nil {
			s.Audio.EventObj(sound.ID(set.Sound(17)), u, 0, 0)
		}
		*cooldown = s.Frame() + uint32(s.Rand.Logic.IntClamp(int(int32(2*s.TickRate())), int(int32(4*s.TickRate()))))
	}
	GetServer().NoxScriptC().ScriptCallback(&layout.As[server.MonsterUpdateData](ud).ScriptEnemySighted, v, u, server.ScriptEventType(14))
}
func visibilityCandidate(candidate, u *server.Object) {
	ud := u.UpdateData
	if candidate == u || uint32(candidate.ObjClass)&6 == 0 || uint32(candidate.ObjFlags)&0x8020 != 0 {
		return
	}
	s := GetServer().S()
	status := (*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).StatusFlags))
	if *status&0x400 == 0 && !s.IsEnemyTo(u, candidate) {
		return
	}
	if visibilityContains(u, candidate) {
		return
	}
	if *status&0x100 == 0 {
		dx, dy := float64(candidate.PosVec.X)-float64(u.PosVec.X), float64(candidate.PosVec.Y)-float64(u.PosVec.Y)
		storedY := float32(dy)
		direction := memmap.PtrOff(0x587000, uintptr(194136+8*int(int16(u.Direction1))))
		length := float32(math.Sqrt(dy*float64(storedY)+dx*dx) + 0.001)
		dot := float64(storedY)/float64(length)*float64(*visibilityFloat(direction, 4)) + dx/float64(length)*float64(*visibilityFloat(direction, 0))
		if !(dot >= 0.5) {
			return
		}
	}
	if s.CanInteract(u, candidate, 0) {
		visibilitySee(u, candidate)
	}
}
func visibilityGlobalRemove(v *server.Object) int {
	for u := GetServer().S().Objs.UpdatableList; u != nil; u = u.UpdatableNext {
		if uint32(u.ObjClass)&2 != 0 && uint32(u.ObjFlags)&0x20 == 0 {
			visibilityRemove(u, v)
			visibilitySelectTarget(u)
		}
	}
	return 0
}
