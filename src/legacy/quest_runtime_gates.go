package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

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
		d := u.UpdateDataPlayer()
		p := d.Player
		if p.Field3680&17 == 17 && p.Field4792 == 1 && d.Field78 == nil && d.Field79 == nil {
			Sub_4DF3C0(p)
			controlLeaveObserver(d.Player)
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
		if u.UpdateDataPlayer().Field77 != nil {
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
	gate := (*server.Object)(layout.As[server.PlayerUpdateData](d).Field79)
	if gate == nil {
		return
	}
	destination := gate.CollideDataExit()
	controlLeaveObserver((*server.Player)(layout.As[server.PlayerUpdateData](d).Player))
	Nox_xxx_playerCameraUnlock_4E6040(u)
	*(*unsafe.Pointer)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](d).Field79)) = nil
	Nox_xxx_unitMove_4E7010(u, destination.Dest)
	GetServer().S().Audio.EventObj(sound.ID(312), u, 2, u.NetCode)
	visibilityFXPoint(129, destination.Dest)
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
		if controlPlayer(u).Field4792 != 0 && layout.As[server.PlayerUpdateData](u.UpdateData).Field79 != nil {
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
		if controlPlayer(u).Field4792 != 0 && layout.As[server.PlayerUpdateData](u.UpdateData).Field79 != nil {
			inGate++
		}
	}
	if count != inGate || worldQuestMaybeWarp() {
		return
	}
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if controlPlayer(u).Field4792 != 0 && layout.As[server.PlayerUpdateData](u.UpdateData).Field79 != nil {
			questRuntimeGateReturn(u)
			text := "Gauntlet.c:WarpRestrictedMulti"
			if count <= 1 {
				text = "Gauntlet.c:WarpRestrictedSolo"
			}
			worldCollideMessage(u, text)
		}
	}
}
