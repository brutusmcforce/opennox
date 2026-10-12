package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

func monsterCache(off uintptr, name string) uint32 {
	p := memmap.PtrUint32(0x5D4594, off)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID(name))
	}
	return *p
}
func monsterIsMimic(u *server.Object) bool {
	return uint32(u.TypeInd) == monsterCache(2488524, "Mimic")
}
func monsterIsPlant(u *server.Object) bool {
	return uint32(u.TypeInd) == monsterCache(2488528, "CarnivorousPlant")
}
func monsterIsZombie(u *server.Object) bool {
	z := memmap.PtrUint32(0x5D4594, 2488532)
	v := memmap.PtrUint32(0x5D4594, 2488536)
	if *z == 0 { // C refreshes both IDs only from the Zombie-cache miss.
		*z = uint32(GetServer().S().Types.IndByID("Zombie"))
		*v = uint32(GetServer().S().Types.IndByID("VileZombie"))
	}
	return uint32(u.TypeInd) == *z || uint32(u.TypeInd) == *v
}
func monsterCanCast(u *server.Object) bool { return u.UpdateDataMonster().StatusFlags&0x20 != 0 }
func monsterCanShoot(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	if u.SubClass()&0x10 != 0 {
		return ud.Field514&0x047f00fe != 0
	}
	return ud.MonsterDef.MissileName148[0] != 0
}
func monsterCanMelee(u *server.Object) bool {
	if monsterCanCast(u) {
		return false
	}
	if u.SubClass()&0x10 != 0 {
		return !monsterCanShoot(u)
	}
	return !(float64(u.UpdateDataMonster().MonsterDef.MeleeAttackRange112) <= 0)
}
func monsterHasShield(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	if u.SubClass()&0x10 != 0 {
		return ud.Field515&0x3000000 != 0
	}
	return ud.StatusFlags&4 != 0
}
func monsterMoving(u *server.Object) bool { return float64(u.SpeedBase) >= .0099999998 }
func monsterHeadSafe(u *server.Object) ai.ActionType {
	ud := u.UpdateDataMonster()
	return ud.AIStack[ud.AIStackInd].Type()
}
func monsterActionToAnimation(u *server.Object) int {
	ud := u.UpdateDataMonster()
	v := 8
	if i := ud.AIStackInd; i != -1 {
		switch ud.AIStack[i].Type() {
		case 7, 8, 10, 13, 29, 36, 37:
			v = int(((uint32(ud.StatusFlags) & 0x4000) | 0x30000) >> 14)
		case 9:
			v = 12
		case 16:
			v = 1
		case 17:
			v = 3
		case 18, 19, 20:
			v = 7
		case 21, 23:
			v = 5
		case 22:
			v = 6
		case 24:
			v = 13
		case 30:
			v = 9
		case 31:
			v = 10
		case 33, 35:
			v = 14
		case 34:
			v = 15
		}
	}
	if monsterIsMimic(u) && v == 8 {
		if ud.StatusFlags&0x40000 != 0 {
			return 0
		}
		return v
	}
	if monsterIsPlant(u) && v == 8 {
		if ud.CurrentEnemy == nil {
			return 14
		}
		return 8 // do not reach Zombie cache lookup on this C branch.
	}
	if monsterIsZombie(u) && v == 9 && ud.StatusFlags&0x80000 != 0 {
		return 15
	}
	return v
}
func monsterCalcDir(u *server.Object, p *float32) {
	if u == nil {
		return
	}
	q := [2]float32{*p - u.PosVec.X, *(*float32)(unsafe.Add(unsafe.Pointer(p), 4)) - u.PosVec.Y}
	u.Direction2 = server.Dir16(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&q[0]))))))
}
func monsterNPCAnim(u *server.Object) unsafe.Pointer {
	ud := u.UpdateDataMonster()
	if u.SubClass()&0x10 == 0 {
		if ud.Field119 == nil {
			return nil
		}
		return unsafe.Pointer(&ud.Field119[monsterActionToAnimation(u)])
	}
	a, b := 0, 0
	lookup := false
	switch ud.AIStack[ud.AIStackInd].Type() {
	case 16, 17:
		if ud.Field514&^3 != 0 {
			a = int(controlWeaponAnimation(uint32(ud.Field514 &^ 3)))
			lookup = true
		} else {
			a = int(byte(ud.Field517))
			lookup = true
		}
	case 18, 19, 20:
		a = 21
		lookup = true
	case 30:
		a = 1
		lookup = true
	case 31:
		a = 2
		lookup = true
	case 21:
		a = 40
		lookup = true
	case 23:
		if ud.Field514&0x7ff8000 != 0 {
			a = 30
			lookup = true
		} else {
			a = 47
			lookup = true
		}
	}
	if lookup {
		a, b = GetServer().S().PlayerAnimFrames(a)
	}
	*memmap.PtrUint32(0x5D4594, 2487964) = 0
	*memmap.PtrUint32(0x5D4594, 2487968) = 0
	*memmap.PtrUint32(0x5D4594, 2487972) = 0
	*memmap.PtrUint32(0x5D4594, 2487976) = 0
	*memmap.PtrUint8(0x5D4594, 2487973) = uint8(a)
	*memmap.PtrUint8(0x5D4594, 2487974) = uint8(b)
	return memmap.PtrOff(0x5D4594, 2487964)
}
func monsterMoveAudio(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.SubClass()&0x30 != 0 {
		a, b := GetServer().S().PlayerAnimFrames(4)
		n := uint32(u.NetCode) + GetServer().S().Frame()
		cur := n / uint32(b+1) % uint32(a)
		old := (n - 1) / uint32(b+1) % uint32(a)
		if cur != old && (cur == ud.MonsterDef.MoveSndFrameA100 || cur == ud.MonsterDef.MoveSndFrameB104) {
			combatSound(u, 18)
		}
		return
	}
	f := uint32(ud.Field120_1)
	if (f == ud.MonsterDef.MoveSndFrameA100 || f == ud.MonsterDef.MoveSndFrameB104) && ud.Field120_2 == 0 {
		combatSound(u, 18)
	}
}
func monsterAttackAtWill(u *server.Object) bool {
	return float64(u.UpdateDataMonster().Aggression) > .66000003
}
func monsterAggressionMid(u *server.Object) bool {
	v := float64(u.UpdateDataMonster().Aggression)
	return v < .66000003 && v > .33000001
}
func monsterAggressionLow(u *server.Object) bool {
	v := float64(u.UpdateDataMonster().Aggression)
	return v < .33000001 && v > .079999998
}
func monsterAggressionRetreat(u *server.Object) bool {
	return float64(u.UpdateDataMonster().Aggression) < .079999998
}
func monsterRunningStatus(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	return ud.StatusFlags&0x40 != 0 || ud.AIStack[ud.AIStackInd].Type() == 4
}
func monsterStartRunning(u *server.Object) uint32 {
	ud := u.UpdateDataMonster()
	if ud.StatusFlags&0x10000 == 0 {
		ud.StatusFlags |= 0x4000
	}
	return uint32(ud.StatusFlags)
}
func monsterStopRunning(u *server.Object) uint32 {
	ud := u.UpdateDataMonster()
	if ud.StatusFlags&0x8000 == 0 {
		ud.StatusFlags &^= 0x4000
	}
	return uint32(ud.StatusFlags)
}
func monsterHasFlag9(u *server.Object) bool { return u.UpdateDataMonster().StatusFlags&0x200 != 0 }
func monsterHasMissingHealth(u *server.Object) bool {
	h := u.HealthData
	return h != nil && h.Cur < h.Max && h.Max != 0
}
func monsterPoisoned(u *server.Object) bool { return u.Poison540 != 0 }
func monsterMoveAttempt(u *server.Object) bool {
	return uint32(GetServer().S().Frame()-u.UpdateDataMonster().Field127) < 3*GetServer().S().TickRate()
}
func monsterMimicMorph(u *server.Object) {
	ud := u.UpdateDataMonster()
	a := ud.AIStack[ud.AIStackInd]
	act := a.Type()
	dx := float64(math.Float32frombits(uint32(a.Args[0]))) - float64(u.PosVec.X)
	dy := float64(math.Float32frombits(uint32(a.Args[1]))) - float64(u.PosVec.Y)
	if act != 0 && (act != 4 || dy*dy+dx*dx > 64) {
		if act != 34 && ud.StatusFlags&0x40000 != 0 {
			u.MonsterPushAction(ai.ActionType(61))
			u.MonsterPushAction(ai.ActionType(34))
			GetServer().S().Audio.EventObj(sound.ID(460), u, 0, 0)
		}
	} else if ud.StatusFlags&0x40000 == 0 && uint32(GetServer().S().Frame()-ud.Field137) > GetServer().S().TickRate() {
		u.MonsterPushAction(ai.ActionType(61))
		u.MonsterPushAction(ai.ActionType(33))
		GetServer().S().Audio.EventObj(sound.ID(460), u, 0, 0)
	}
}

func monsterOrder(owner, unit *server.Object, order int) {
	if owner == nil {
		return
	}
	if unit != nil {
		monsterEnactOrder(owner, unit, order)
		return
	}
	if owner.Class().Has(object.ClassPlayer) && (order == 3 || order == 4 || order == 5) {
		unitLocalOrder(int(owner.UpdateDataPlayer().Player.PlayerInd), int32(order))
	}
	for v := owner.Field129; v != nil; v = v.Field128 {
		if v.Class().Has(object.ClassMonster) && v.UpdateDataMonster().StatusFlags&0x80 != 0 {
			monsterEnactOrder(owner, v, order)
		}
	}
}
func monsterEnactOrder(source, u *server.Object, order int) {
	if source == nil || !u.Class().Has(object.ClassMonster) {
		return
	}
	ud := u.UpdateDataMonster()
	if !(monsterIsZombie(u) && order == 0 || u.ObjFlags&0x8000 == 0) {
		return
	}
	if ud.MonsterDef == nil {
		ud.MonsterDef = monsterDefinitionByType(uint32(u.TypeInd))
		if ud.MonsterDef == nil {
			return
		}
	}
	switch order {
	case 0:
		if u.ObjOwner == source {
			spellEffectBanish(u)
		}
	case 1:
		if source.Class().Has(object.ClassPlayer) {
			Nox_xxx_playerObserveMonster_4DDE80(source, u)
		}
	case 2:
		if source.Class().Has(object.ClassPlayer) {
			visibilityMonsterCommand(u, source, "MonUtil.c:idle", 0)
		}
		combatSound(u, 17)
		ud.StatusFlags &^= 0x40
		ud.Aggression = math.Float32frombits(0x3f000000)
		u.ClearActionStack()
		u.MonsterPushAction(ai.ACTION_IDLE)
	case 3:
		if !monsterMoving(u) {
			return
		}
		if source.Class().Has(object.ClassPlayer) {
			visibilityMonsterCommand(u, source, "MonUtil.c:guarding", 0)
		}
		combatSound(u, 17)
		ud.Aggression = math.Float32frombits(0x3f000000)
		if monsterCanShoot(u) {
			ud.StatusFlags |= 0x40
		}
		ud.SightRange = math.Float32frombits(0x437a0000)
		u.ClearActionStack()
		u.MonsterPushAction(ai.ACTION_GUARD, u.PosVec, int(int16(u.Direction1)))
	case 4:
		if !monsterMoving(u) {
			return
		}
		if source.Class().Has(object.ClassPlayer) {
			visibilityMonsterCommand(u, source, "MonUtil.c:escorting", 0)
		}
		combatSound(u, 17)
		ud.StatusFlags &^= 0x40
		ud.Aggression = math.Float32frombits(0x3f547ae1)
		u.ClearActionStack()
		u.MonsterPushAction(ai.ACTION_ESCORT, source.PosVec, source)
	case 5:
		if !monsterMoving(u) {
			return
		}
		combatSound(u, 17)
		if source.Class().Has(object.ClassPlayer) {
			visibilityMonsterCommand(u, source, "MonUtil.c:Hunting", 0)
		}
		ud.StatusFlags &^= 0x40
		ud.Aggression = math.Float32frombits(0x3f547ae1)
		u.ClearActionStack()
		u.MonsterPushAction(ai.ACTION_HUNT)
	}
}

func nox_xxx_orderUnit_533900(a, b *nox_object_t, o int32) {
	monsterOrder(asObjectS(a), asObjectS(b), int(o))
}

func objectFromInt(p int32) *server.Object {
	return objectFromWord(uint32(p))
}

func objectFromWord(p uint32) *server.Object {
	return (*server.Object)(unsafe.Pointer(uintptr(p)))
}

func monsterCanBlockHead(u *server.Object) bool {
	a := monsterHeadSafe(u)
	return a == 0 || a == 1 || a == 4 || a == 25 || a == 26 || a == 27 || a == 23
}
