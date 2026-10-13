package legacy

import (
	"math"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

type navigationAIAction struct {
	typ         ai.ActionType
	update      func(*server.Object)
	run, cancel bool
}

func (a navigationAIAction) Type() ai.ActionType     { return a.typ }
func (a navigationAIAction) Update(u *server.Object) { a.update(u) }
func (a navigationAIAction) Start(u *server.Object) {
	if a.run {
		navigationStartRunning(u)
	}
}
func (a navigationAIAction) End(u *server.Object) {
	if a.run {
		escortStopRunning(u)
	}
}
func (a navigationAIAction) Cancel(u *server.Object) {
	if a.cancel {
		escortStopRunning(u)
	}
}
func init() {
	for _, a := range []navigationAIAction{
		{ai.ACTION_MOVE_TO, navigationMove, false, false},
		{ai.ACTION_FAR_MOVE_TO, navigationFarMove, false, false},
		{ai.ACTION_DODGE, navigationDodge, false, false},
		{ai.ACTION_FLEE, navigationFlee, true, false},
		{ai.ACTION_MOVE_TO_HOME, navigationMove, true, true},
		{ai.ACTION_RETREAT, navigationRetreat, false, false},
		{ai.ACTION_RETREAT_TO_MASTER, navigationMaster, true, false},
	} {
		server.RegisterAIAction(a)
	}
}
func navigationStartRunning(u *server.Object) { monsterStartRunning(u) }

// The movement gate uses base speed; dodge modifies current speed separately.
func navigationMoving(u *server.Object) bool { return monsterMoving(u) }
func navigationPrevious(u *server.Object) ai.ActionType {
	ud := u.UpdateDataMonster()
	for i := int(ud.AIStackInd) - 1; i >= 0; i-- {
		if a := ud.AIStack[i].Type(); !a.IsCondition() {
			return a
		}
	}
	return ai.ActionType(38)
}
func navigationAudio(u *server.Object) { monsterMoveAudio(u) }
func navigationMove(u *server.Object) {
	ud := u.UpdateDataMonster()
	if !navigationMoving(u) {
		u.MonsterPopAction()
		return
	}
	head := ud.AIStackHead()
	if navigationPrevious(u) == ai.ACTION_ESCORT {
		radius := float64(ud.Field329) * 3
		inner := float64(float32(radius))
		outer := radius + 30
		p := head.ArgPos(0)
		dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
		dist := dy*dy + dx*dx
		if dist >= inner*inner {
			if dist > outer*outer {
				navigationStartRunning(u)
			}
		} else {
			escortStopRunning(u)
		}
	}
	if pathSetMove(u) {
		status := byte(ud.Field71)
		core := GetServer().S()
		retry := status == 2 || status == 1 && uint32(core.Frame()-ud.Field135) < 5*core.TickRate()
		if status == 1 {
			ud.Field135 = core.Frame()
		}
		if status == 0 && pathTakeStatus() == 0 && head.Args[2] == 0 {
			monsterCalcDir(u, head.ArgPos(0))
			u.MonsterPopAction()
		}
		if retry {
			if st := u.MonsterPushAction(ai.ActionType(41)); st != nil {
				st.Args[0] = uintptr(core.Frame() + uint32(nox_common_randomInt_415FA0(int(2*core.TickRate()), int(4*core.TickRate()))))
			}
			u.MonsterPushAction(ai.ACTION_RANDOM_WALK)
			ud.StatusFlags |= 0x200000
		}
		if byte(ud.Field71) != 0 {
			if st := u.MonsterPushAction(ai.ACTION_WAIT); st != nil {
				st.Args[0] = uintptr(core.Frame() + uint32(nox_common_randomInt_415FA0(int(core.TickRate()>>1), int(core.TickRate()))))
			}
		}
	}
	navigationAudio(u)
}
func navigationFarMove(u *server.Object) {
	if (aiRespondsToThreat(u) || aiAttackAtWill(u)) && u.Sub_545E60() != 0 {
		return
	}
	if aiAttackAtWill(u) && u.UpdateDataMonster().CurrentEnemy != nil {
		aiPushFight(u)
	}
	navigationMove(u)
}
func navigationDodge(u *server.Object) {
	if !navigationMoving(u) {
		u.MonsterPopAction()
		return
	}
	if u.HasEnchant(3) || u.HasEnchant(5) || u.HasEnchant(28) {
		return
	}
	ud := u.UpdateDataMonster()
	p := ud.AIStackHead().ArgPos(0)
	dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
	// x87 retains both deltas and the speed product despite float32 stores.
	// Only the length denominator is reloaded from float32.
	length := math.Sqrt(dy*dy+dx*dx) + .000099999997
	div := float32(length)
	if length >= 8 {
		speed := float64(ud.MonsterDef.RunMultiplier96) * float64(u.SpeedCur)
		u.SpeedCur = float32(speed)
		u.ForceVec.X = float32(speed * dx / float64(div))
		u.ForceVec.Y = float32(dy * speed / float64(div))
	} else {
		u.MonsterPopAction()
	}
}
func navigationFlee(u *server.Object) {
	if !navigationMoving(u) {
		u.MonsterPopAction()
		return
	}
	ud := u.UpdateDataMonster()
	head := ud.AIStackHead()
	core := GetServer().S()
	if enemy := ud.CurrentEnemy; enemy != nil {
		head.SetArgs(enemy.PosVec)
		dx, dy := float64(enemy.PosVec.X)-float64(u.PosVec.X), float64(enemy.PosVec.Y)-float64(u.PosVec.Y)
		radius := float64(ud.FleeRange)
		if radius*radius > dy*dy+dx*dx && uint32(core.Frame()-ud.Field70) > core.TickRate()>>1 {
			ud.Field2 = 0
		}
		if monsterCanCast(u) && !u.HasEnchant(29) && !monsterAggressionLow(u) && !u.Sub_534440() {
			var cast bool
			if ud.HasAction(ai.ACTION_RETREAT) {
				cast = monsterCastRelated(u)
			} else {
				cast = monsterBuffSelf(u)
			}
			if !cast && core.CanInteract(u, ud.CurrentEnemy, 0) {
				monsterCastOffensive(u, ud.CurrentEnemy)
			}
		}
	}
	frame := core.Frame()
	if ud.Field2 != 0 && uint32(core.Frame()-ud.Field70) > 2*core.TickRate() {
		ud.Field2 = 0
		frame = core.Frame()
	}
	move := ud.Field2 != 0 || uint32(frame-ud.Field70) <= 10
	if !move {
		target := head.ArgPos(0)
		count := GetServer().Nox_xxx_generateRetreatPath_50CA00(ud.Path[:], u, &target)
		ud.Field2 = uint32(count)
		ud.Field70 = core.Frame()
		ud.Field67 = 0
		move = count > 1
	}
	if move {
		if pathActuallyMove(u) {
			ud.Field2 = 0
		}
		navigationAudio(u)
	} else {
		head.SetArgs(u.PosVec)
	}
}
func navigationCanResume(u *server.Object) bool {
	ratio := 1.0
	if u.HealthData.Max != 0 {
		ratio = float64(u.HealthData.Cur) / float64(u.HealthData.Max)
	}
	return ratio >= float64(u.UpdateDataMonster().ResumeLevel)
}
func navigationCanCast2(u *server.Object) bool {
	return u.UpdateDataMonster().StatusFlags&0x20 != 0 && u.HasEnchant(29)
}
func navigationShouldRetreat(u *server.Object) bool {
	return !navigationCanResume(u) || navigationCanCast2(u)
}
func navigationRetreat(u *server.Object) {
	ud := u.UpdateDataMonster()
	if !navigationShouldRetreat(u) {
		u.MonsterPopAction()
		return
	}
	if ud.CurrentEnemy != nil {
		if u.HasEnchant(29) || !monsterCastRelated(u) {
			core := GetServer().S()
			if st := u.MonsterPushAction(ai.ActionType(41)); st != nil {
				st.Args[0] = uintptr(core.Frame() + uint32(nox_common_randomInt_415FA0(int(4*core.TickRate()), int(6*core.TickRate()))))
			}
			if st := u.MonsterPushAction(ai.ACTION_FLEE); st != nil {
				st.SetArgs(ud.CurrentEnemy.PosVec, uint32(0))
			}
		}
	} else if !navigationCanResume(u) {
		navigationEdibles(u)
	}
}
func navigationEdibles(u *server.Object) {
	radius := math.Float32frombits(1132068864)
	if noxflags.HasGame(noxflags.GameFlag(4096)) {
		radius = math.Float32frombits(1142947840)
	}
	food := (*server.Object)(unsafe.Pointer(uintptr(Nox_xxx_mobSearchEdible_544A00(u, radius))))
	u.MonsterPushAction(ai.ActionType(64))
	u.MonsterPushAction(ai.ActionType(56))
	if food != nil {
		if st := u.MonsterPushAction(ai.ActionType(48)); st != nil {
			st.SetArgs(food.PosVec, food)
		}
		if st := u.MonsterPushAction(ai.ACTION_PICKUP_OBJECT); st != nil {
			st.SetArgs(food)
		}
		if st := u.MonsterPushAction(ai.ACTION_MOVE_TO); st != nil {
			st.SetArgs(food.PosVec, food)
		}
	} else {
		u.MonsterPushAction(ai.ActionType(58))
		if st := u.MonsterPushAction(ai.ACTION_ROAM); st != nil {
			st.Args[0] = 0
			st.Args[2] = (st.Args[2] &^ 0xff) | 0x80
		}
	}
}
func navigationMaster(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.ObjOwner == nil || !navigationShouldRetreat(u) {
		u.MonsterPopAction()
		return
	}
	owner := u.ObjOwner
	dx, dy := float64(u.PosVec.X)-float64(owner.PosVec.X), float64(u.PosVec.Y)-float64(owner.PosVec.Y)
	radius := float64(ud.Field329) + 30
	if radius*radius < dy*dy+dx*dx {
		if st := u.MonsterPushAction(ai.ActionType(49)); st != nil {
			st.Args[0] = uintptr(math.Float32bits(ud.Field329))
			st.Args[2] = uintptr(unsafe.Pointer(u.ObjOwner))
		}
		if st := u.MonsterPushAction(ai.ACTION_MOVE_TO); st != nil {
			st.SetArgs(u.ObjOwner.PosVec, u.ObjOwner)
		}
	}
}
