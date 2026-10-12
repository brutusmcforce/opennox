package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
)

func questRuntimeObserverDeadline() uint32 {
	deadline := *questRuntimeDeadlinePtr()
	if deadline == 0 || GetServer().S().Frame() <= deadline {
		return deadline
	}
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		d := u.UpdateData
		p := *controlPtr(d, 276)
		if *equipmentWord(p, 3680)&17 == 17 && *equipmentWord(p, 4792) == 1 && *controlPtr(d, 312) == nil && *controlPtr(d, 316) == nil {
			Sub_4DF3C0((*server.Player)(p))
			controlLeaveObserver(*controlPtr(d, 276))
			Nox_xxx_playerCameraUnlock_4E6040(u)
		}
	}
	return 0
}
func questRuntimeSoulTimeout() uint32 {
	start := *questRuntimeSoulFramePtr()
	core := GetServer().S()
	if start == 0 || core.Frame()-start < 9000 {
		return start
	}
	any := false
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if *controlPtr(u.UpdateData, 308) != nil {
			any = true
		}
	}
	if !any {
		return 0
	}
	count := questRuntimeCount()
	if count <= 1 {
		return uint32(count)
	}
	*questRuntimeSoulFramePtr() = 0
	if core.Doors.Sub_4D72C0() {
		return 1
	}
	core.Doors.Sub_4D72B0(true)
	return uint32(questRuntimeSharedMessage(255, 1))
}
func questRuntimeGateReturn(u *server.Object) {
	if u == nil || u.ObjClass&4 == 0 {
		return
	}
	d := u.UpdateData
	gate := controlObject(d, 316)
	if gate == nil {
		return
	}
	destination := gate.CollideData
	controlLeaveObserver(*controlPtr(d, 276))
	Nox_xxx_playerCameraUnlock_4E6040(u)
	*controlPtr(d, 316) = nil
	Nox_xxx_unitMove_4E7010(u, *(*types.Pointf)(unsafe.Add(destination, 80)))
	GetServer().S().Audio.EventObj(sound.ID(312), u, 2, u.NetCode)
	visibilityFXPoint(129, *(*types.Pointf)(unsafe.Add(destination, 80)))
}
func questRuntimeGateSet(next uint32) int8 {
	p := memmap.PtrUint32(0x5D4594, 1556120)
	previous := *p
	if previous != 1 || next != 0 {
		*p = next
		return int8(previous)
	}
	core := GetServer().S()
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if *equipmentWord(controlPlayer(u), 4792) != 0 && *controlPtr(u.UpdateData, 316) != nil {
			questRuntimeGateReturn(u)
		}
	}
	result := int8(0)
	for u := core.Objs.First(); u != nil; {
		nextObject := u.Next()
		result = int8(u.ObjClass)
		if u.ObjClass&0x20 != 0 && u.ObjSubClass&2 != 0 {
			result = int8(stateOff(u))
		}
		u = nextObject
	}
	*p = 0
	return result
}
func questRuntimeWarpTick() {
	count := questRuntimeCount()
	core := GetServer().S()
	if count == 0 || core.Frame()-memmap.Uint32(0x5D4594, 1556108) < 30 {
		return
	}
	inGate := 0
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if *equipmentWord(controlPlayer(u), 4792) != 0 && *controlPtr(u.UpdateData, 316) != nil {
			inGate++
		}
	}
	if count != inGate || worldQuestMaybeWarp() {
		return
	}
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if *equipmentWord(controlPlayer(u), 4792) != 0 && *controlPtr(u.UpdateData, 316) != nil {
			questRuntimeGateReturn(u)
			text := "Gauntlet.c:WarpRestrictedMulti"
			if count <= 1 {
				text = "Gauntlet.c:WarpRestrictedSolo"
			}
			worldCollideMessage(u, text)
		}
	}
}
