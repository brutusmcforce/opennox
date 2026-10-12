package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

type guardAIAction struct{}

func (guardAIAction) Type() ai.ActionType     { return ai.ACTION_GUARD }
func (guardAIAction) Start(*server.Object)    {}
func (guardAIAction) End(*server.Object)      {}
func (guardAIAction) Cancel(*server.Object)   {}
func (guardAIAction) Update(u *server.Object) { guardUpdate(u) }

type escortAIAction struct{}

func (escortAIAction) Type() ai.ActionType     { return ai.ACTION_ESCORT }
func (escortAIAction) Start(*server.Object)    {}
func (escortAIAction) End(u *server.Object)    { escortStopRunning(u) }
func (escortAIAction) Cancel(u *server.Object) { escortStopRunning(u) }
func (escortAIAction) Update(u *server.Object) { escortUpdate(u) }
func init()                                    { server.RegisterAIAction(guardAIAction{}); server.RegisterAIAction(escortAIAction{}) }

func aiAttackAtWill(u *server.Object) bool {
	return float64(u.UpdateDataMonster().Aggression) > .66000003
}
func aiRespondsToThreat(u *server.Object) bool {
	a := float64(u.UpdateDataMonster().Aggression)
	return a < .66000003 && a > .33000001
}

// These caches are still shared with real C predicate callers.
func aiCachedType(u *server.Object, offset uintptr, name string) bool {
	index := memmap.PtrUint32(0x5D4594, offset)
	if *index == 0 {
		*index = uint32(GetServer().S().Types.IndByID(name))
	}
	return uint32(u.TypeInd) == *index
}
func aiIsMimic(u *server.Object) bool { return aiCachedType(u, 2488524, "Mimic") }
func aiIsPlant(u *server.Object) bool { return aiCachedType(u, 2488528, "CarnivorousPlant") }
func aiPushFight(u *server.Object) {
	if st := u.MonsterPushAction(ai.ACTION_FIGHT); st != nil {
		st.SetArgs(u.UpdateDataMonster().CurrentEnemy.PosVec, GetServer().S().Frame())
	}
}
func escortStopRunning(u *server.Object) { monsterStopRunning(u) }
func escortResolve(u *server.Object) *server.Object {
	ud := u.UpdateDataMonster()
	name := (*byte)(unsafe.Pointer(&ud.Field341))
	value := alloc.GoString(name)
	defer func() { *name = 0 }()
	switch value {
	case "**OWNER**":
		return u.ObjOwner
	case "**PLAYER**":
		players := &GetServer().S().Players
		count := 0
		for it := players.FirstUnit(); it != nil; it = players.NextUnit(it) {
			count++
		}
		index := nox_common_randomInt_415FA0(0, count-1)
		for it := players.FirstUnit(); it != nil; it = players.NextUnit(it) {
			if index == 0 {
				return it
			}
			index--
		}
		return nil
	default:
		return Nox_xxx_getObjectByScrName_4DA4F0(value)
	}
}
func escortUpdate(u *server.Object) {
	ud := u.UpdateDataMonster()
	head := ud.AIStackHead()
	if head.ArgObj(2) == nil {
		target := escortResolve(u)
		head.Args[2] = uintptr(unsafe.Pointer(target))
		if target == nil {
			u.MonsterPopAction()
			return
		}
		head.Args[0], head.Args[1] = uintptr(math.Float32bits(target.PosVec.X)), uintptr(math.Float32bits(target.PosVec.Y))
	}
	if aiAttackAtWill(u) {
		if ud.CurrentEnemy != nil {
			aiPushFight(u)
			return
		}
		if u.MonsterLookAtDamager() {
			return
		}
	} else if aiRespondsToThreat(u) {
		if u.Sub_545E60() != 0 || u.MonsterLookAtDamager() {
			return
		}
	}
	target := head.ArgObj(2)
	dx, dy := float64(u.PosVec.X)-float64(target.PosVec.X), float64(u.PosVec.Y)-float64(target.PosVec.Y)
	radius := float64(ud.Field329) + 30
	if radius*radius >= dy*dy+dx*dx {
		if !aiAttackAtWill(u) || investigateHeardSound(u) == 0 {
			if !u.HasEnchant(29) {
				Nox_xxx_mobHealSomeone_5411A0(u)
			}
		}
	} else {
		if aiRespondsToThreat(u) || aiAttackAtWill(u) {
			u.MonsterPushAction(ai.DEPENDENCY_NOT_UNDER_ATTACK)
		}
		if aiAttackAtWill(u) {
			u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
		}
		if st := u.MonsterPushAction(ai.DEPENDENCY_OBJECT_FARTHER_THAN); st != nil {
			st.SetArgs(ud.Field329, uint32(0), head.ArgObj(2))
		}
		if st := u.MonsterPushAction(ai.ACTION_MOVE_TO); st != nil {
			st.SetArgs(head.ArgObj(2).PosVec, head.ArgObj(2))
		}
	}
}
func guardUpdate(u *server.Object) {
	ud := u.UpdateDataMonster()
	head := ud.AIStackHead()
	if !(float64(ud.Aggression) < .079999998) && u.Sub_545E60() != 0 {
		return
	}
	if ud.CurrentEnemy != nil {
		if aiAttackAtWill(u) {
			aiPushFight(u)
			return
		}
		if aiRespondsToThreat(u) {
			home := head.ArgPos(0)
			enemy := ud.CurrentEnemy.PosVec
			dx, dy := float64(home.X)-float64(enemy.X), float64(home.Y)-float64(enemy.Y)
			if float64(ud.SightRange)*float64(ud.SightRange) > dy*dy+dx*dx {
				if aiIsPlant(u) {
					if st := u.MonsterPushAction(ai.DEPENDENCY_ENEMY_CLOSER_THAN); st != nil {
						st.SetArgs(float32(float64(ud.SightRange) * 1.05))
					}
				} else {
					if st := u.MonsterPushAction(ai.DEPENDENCY_UNDER_ATTACK); st != nil {
						st.SetArgs(uint32(0))
					}
					if st := u.MonsterPushAction(ai.DEPENDENCY_LOCATION_CLOSER_THAN); st != nil {
						st.SetArgs(float32(float64(ud.SightRange)*1.5), uint32(0), head.ArgPos(0))
					}
					u.MonsterPushAction(ai.DEPENDENCY_OR)
				}
				aiPushFight(u)
			}
		}
	}
	if !aiIsMimic(u) && u.MonsterLookAtDamager() {
		return
	}
	if aiAttackAtWill(u) {
		if investigateHeardSound(u) != 0 {
			return
		}
	} else if !aiIsMimic(u) && heardSoundAction(u) != 0 {
		return
	}
	if (GetServer().S().Frame()+u.NetCode)&15 == 0 {
		home := head.ArgPos(0)
		dx, dy := float64(home.X)-float64(u.PosVec.X), float64(home.Y)-float64(u.PosVec.Y)
		if dy*dy+dx*dx > 64 {
			if aiAttackAtWill(u) {
				u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
				u.MonsterPushAction(ai.DEPENDENCY_NO_INTERESTING_SOUND)
			}
			u.MonsterPushAction(ai.DEPENDENCY_NOT_UNDER_ATTACK)
			if st := u.MonsterPushAction(ai.ACTION_MOVE_TO); st != nil {
				st.SetArgs(head.ArgPos(0), uint32(0))
			}
			return
		}
		if !aiIsMimic(u) && ud.Field282_1 > 0 {
			target := ud.CurrentEnemy
			if target == nil {
				// The seen list holds object pointers (see visibilitySlots).
				seen := unsafe.Slice((**server.Object)(unsafe.Pointer(&ud.Field283)), int(ud.Field282_1))
				target = seen[0]
				for _, candidate := range seen {
					if candidate.ObjClass.Has(object.ClassPlayer) {
						target = candidate
						break
					}
				}
			}
			dx, dy := float64(target.PosVec.X)-float64(u.PosVec.X), float64(target.PosVec.Y)-float64(u.PosVec.Y)
			length := math.Sqrt(dy*dy+dx*dx) + .001
			if !facingDot(u, types.Pointf{X: float32(dx / length), Y: float32(dy / length)}) {
				if st := u.MonsterPushAction(ai.ACTION_FACE_OBJECT); st != nil {
					st.SetArgs(target)
				}
			}
		} else {
			direction := head.ArgU32(2)
			x, y := movementDirectionVector(int32(direction))
			if !facingDot(u, types.Pointf{X: x, Y: y}) {
				if !aiIsMimic(u) {
					u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
				}
				if st := u.MonsterPushAction(ai.ACTION_FACE_ANGLE); st != nil {
					st.SetArgs(direction)
				}
			}
		}
	}
	if aiIsMimic(u) || GetServer().S().Frame()-ud.Field137 <= GetServer().S().TickRate()>>1 || u.PosVec.X == u.PrevPos.X && u.PosVec.Y == u.PrevPos.Y {
		if !u.HasEnchant(29) {
			Nox_xxx_mobHealSomeone_5411A0(u)
		}
	} else {
		if st := u.MonsterPushAction(ai.ACTION_FACE_LOCATION); st != nil {
			st.SetArgs(u.PrevPos)
		}
	}
}
