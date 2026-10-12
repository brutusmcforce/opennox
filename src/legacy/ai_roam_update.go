package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

func roamDebug(u *server.Object, message string) {
	name := GetServer().S().Types.ByInd(int(u.TypeInd)).ID()
	if noxflags.HasEngine(noxflags.EngineShowAI) {
		ai.Log.Printf("%d: %s(#%d) %s\n", int32(GetServer().S().Frame()), name, int32(u.NetCode), message)
	}
}
func roamUpdate(u *server.Object) {
	ud := u.UpdateDataMonster()
	// These C predicates use double thresholds, not float32-rounded constants.
	attackAtWill := func() bool { return float64(ud.Aggression) > .66000003 }
	reacts := float64(ud.Aggression) < .66000003 && float64(ud.Aggression) > .33000001
	if (reacts || attackAtWill()) && u.Sub_545E60() != 0 {
		return
	}
	if uint32(ud.StatusFlags)&0x20000 != 0 && ud.CurrentEnemy == nil && u.HasEnchant(0) && GetServer().S().Frame()&31 == 0 && nox_common_randomInt_415FA0(0, 100) < 10 {
		if st := u.MonsterPushActionImpl(ai.DEPENDENCY_ENEMY_FARTHER_THAN, "go", 0); st != nil {
			st.SetArgs(float32(150))
		}
		if st := u.MonsterPushActionImpl(ai.DEPENDENCY_IS_ENCHANTED, "go", 0); st != nil {
			st.SetArgs(uint32(0))
		}
		if st := u.MonsterPushActionImpl(ai.ACTION_WAIT, "go", 0); st != nil {
			st.SetArgs(GetServer().S().Frame() + GetServer().S().TickRate()*uint32(nox_common_randomInt_415FA0(3, 10)))
		}
		return
	}
	if attackAtWill() {
		if ud.CurrentEnemy != nil {
			if st := u.MonsterPushActionImpl(ai.ACTION_FIGHT, "go", 0); st != nil {
				st.SetArgs(ud.CurrentEnemy.PosVec, GetServer().S().Frame())
			}
			return
		}
		if investigateHeardSound(u) != 0 {
			return
		}
	}
	head := ud.AIStackHead()
	mask := byte(head.Args[2])
	wp := roamWaypoint(uint32(head.Args[0]))
	if !roamEligible(wp, mask) {
		head.Args[0] = 0
		wp = nil
	}
	if wp == nil {
		wp = GetServer().S().Sub_518460(u.PosVec, mask, true)
		head.Args[0] = uintptr(unsafe.Pointer(wp))
		if wp == nil && byte(head.Args[2]) == 128 {
			wp = GetServer().Sub_50CB20(u, &u.PosVec)
			head.Args[0] = uintptr(unsafe.Pointer(wp))
		}
		if wp == nil {
			roamDebug(u, "Cannot find any waypoints")
			u.MonsterPopAction()
			if st := u.MonsterPushActionImpl(ai.ACTION_WAIT, "go", 0); st != nil {
				st.SetArgs(GetServer().S().Frame() + uint32(nox_common_randomInt_415FA0(int(int32(GetServer().S().TickRate())), int(int32(2*GetServer().S().TickRate())))))
			}
			return
		}
		roamInsert(ud, wp)
	}
	dx, dy := float64(wp.PosVec.X)-float64(u.PosVec.X), float64(wp.PosVec.Y)-float64(u.PosVec.Y)
	if dy*dy+dx*dx <= 64 {
		if !roamDeadEnd(u, wp) {
			return
		}
		ud.Field70 = 0
	}
	if ud.Field2 == 0 {
		wp = roamWaypoint(uint32(head.Args[0]))
		GetServer().Nox_xxx_creatureSetDetailedPath_50D220(u, &wp.PosVec)
	}
	if byte(ud.Field71) == 2 {
		roamDebug(u, "Cannot compute path to waypoint, choosing other")
		previous := roamPrevious(ud, mask)
		if previous == nil {
			roamDebug(u, "No previous waypoint, giving up.")
			u.MonsterPopAction()
			return
		}
		if !roamDeadEnd(u, previous) {
			return
		}
	}
	if pathActuallyMove(u) {
		ud.Field2 = 0
	}
	monsterMoveAudio(u)
}
