package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

type roamAIAction struct{}

func (roamAIAction) Type() ai.ActionType     { return ai.ACTION_ROAM }
func (roamAIAction) End(*server.Object)      {}
func (roamAIAction) Update(u *server.Object) { roamUpdate(u) }

func (a roamAIAction) Start(u *server.Object) {
	ud := u.UpdateDataMonster()
	// The original start clears only this slot, despite its empty 16-step loop.
	ud.Field75 = 0
	ud.Field91 = 0
}
func (a roamAIAction) Cancel(u *server.Object) { u.UpdateDataMonster().AIStackHead().Args[0] = 0 }
func init() {
	server.RegisterAIAction(roamAIAction{})
}

// The shared monster layout stores these C-owned waypoint addresses as words.
func roamHistory(ud *server.MonsterUpdateData) *[16]uint32 {
	return (*[16]uint32)(unsafe.Pointer(&ud.Field75))
}
func roamWaypoint(raw uint32) *server.Waypoint {
	return (*server.Waypoint)(memmap.WordPtr(raw))
}
func roamWaypointWord(wp *server.Waypoint) uint32 { return memmap.PtrWord(unsafe.Pointer(wp)) }
func roamEligible(wp *server.Waypoint, mask byte) bool {
	return wp != nil && wp.IsEnabled() && wp.HasFlag2Mask(mask)
}
func roamInsert(ud *server.MonsterUpdateData, wp *server.Waypoint) {
	ud.Field91++
	if ud.Field91 >= 16 {
		ud.Field91 = 0
	}
	history := roamHistory(ud)
	word := roamWaypointWord(wp)
	history[ud.Field91] = word
	for i, v := range history {
		if v == word && uint32(i) != ud.Field91 {
			history[i] = 0
		}
	}
}
func roamPrevious(ud *server.MonsterUpdateData, mask byte) *server.Waypoint {
	history := roamHistory(ud)
	for k := 1; k < 16; k++ {
		i := int(ud.Field91) - k
		if i < 0 {
			i += 16
		}
		wp := roamWaypoint(history[i])
		if roamEligible(wp, mask) {
			return wp
		}
	}
	return nil
}
func roamSuccessor(ud *server.MonsterUpdateData, wp *server.Waypoint, mask byte) *server.Waypoint {
	history := roamHistory(ud)
	var candidates [32]*server.Waypoint
	n := 0
	for i := 0; i < int(wp.PointsCnt); i++ {
		next := wp.Points[i].Waypoint
		if !roamEligible(next, mask) {
			continue
		}
		seen := false
		for _, old := range history {
			if old == roamWaypointWord(next) {
				seen = true
				break
			}
		}
		if !seen {
			candidates[n] = next
			n++
		}
	}
	if n != 0 {
		return candidates[nox_common_randomInt_415FA0(0, n-1)]
	}
	for k := 1; k <= 16; k++ {
		i := int(ud.Field91) + k
		if i >= 16 {
			i -= 16
		}
		old := roamWaypoint(history[i])
		if old == nil {
			continue
		}
		for j := 0; j < int(wp.PointsCnt); j++ {
			if old == wp.Points[j].Waypoint {
				if roamEligible(old, mask) {
					return old
				}
				break
			}
		}
	}
	return nil
}
func roamDeadEnd(u *server.Object, wp *server.Waypoint) bool {
	ud := u.UpdateDataMonster()
	head := ud.AIStackHead()
	ud.Field2 = 0
	if wp.PointsCnt != 0 {
		next := roamSuccessor(ud, wp, byte(head.Args[2]))
		head.Args[0] = uintptr(unsafe.Pointer(next))
		if wp.PointsCnt != 0 && next != nil {
			roamInsert(ud, next)
			return true
		}
	}
	// Keep lookup before the logging gate, matching the original caller.
	name := GetServer().S().Types.ByInd(int(u.TypeInd)).ID()
	if noxflags.HasEngine(noxflags.EngineShowAI) {
		ai.Log.Printf("%d: %s(#%d) Reached dead end, giving up.\n", int32(GetServer().S().Frame()), name, int32(u.NetCode))
	}
	u.MonsterPopAction()
	return false
}
