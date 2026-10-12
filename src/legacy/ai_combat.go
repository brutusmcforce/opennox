package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

type combatAIAction struct{ typ ai.ActionType }

func (a combatAIAction) Type() ai.ActionType { return a.typ }
func (a combatAIAction) Start(u *server.Object) {
	switch a.typ {
	case ai.ACTION_FIGHT:
		combatFightStart(u)
	case ai.ACTION_MELEE_ATTACK:
		combatMeleeStart(u)
	case ai.ACTION_MISSILE_ATTACK:
		combatMissileStart(u)
	}
}
func (a combatAIAction) End(u *server.Object) {
	if a.typ == ai.ACTION_FIGHT {
		u.UpdateDataMonster().StatusFlags &^= 0x100
		escortStopRunning(u)
	}
}
func (a combatAIAction) Cancel(u *server.Object) {
	if a.typ != ai.ACTION_FIGHT {
		u.MonsterPopAction()
	}
}
func (a combatAIAction) Update(u *server.Object) {
	switch a.typ {
	case ai.ACTION_FIGHT:
		combatFight(u)
	case ai.ACTION_MELEE_ATTACK:
		combatMelee(u)
	case ai.ACTION_MISSILE_ATTACK:
		combatMissile(u)
	case ai.ACTION_BLOCK_ATTACK:
		combatBlock(u)
	case ai.ACTION_BLOCK_FINISH, ai.ACTION_WEAPON_BLOCK:
		if u.UpdateDataMonster().Field120_3 != 0 {
			u.MonsterPopAction()
		}
	}
}
func init() {
	for _, a := range []ai.ActionType{ai.ACTION_FIGHT, ai.ACTION_MELEE_ATTACK, ai.ACTION_MISSILE_ATTACK, ai.ACTION_BLOCK_ATTACK, ai.ACTION_BLOCK_FINISH, ai.ACTION_WEAPON_BLOCK} {
		server.RegisterAIAction(combatAIAction{a})
	}
}
func combatPtr(u *server.Object) int32 { return int32(uintptr(u.CObj())) }
func combatSound(u *server.Object, index int) {
	if p := Nox_xxx_monsterGetSoundSet_424300(u); p != nil {
		id := *(*uint32)(unsafe.Add(p, 4*index))
		GetServer().S().Audio.EventObj(sound.ID(id), u, 0, 0)
	}
}
func combatChase(u, t *server.Object) {
	u.MonsterPushAction(ai.DEPENDENCY_NO_NEW_ENEMY, t)
	u.MonsterPushAction(ai.DEPENDENCY_ALIVE, t)
	u.MonsterPushAction(ai.ACTION_MOVE_TO, t.PosVec, t)
}
func combatMeleeChain(u, t *server.Object) {
	d := u.UpdateDataMonster().MonsterDef
	u.MonsterPushAction(ai.DEPENDENCY_NO_NEW_ENEMY, t)
	u.MonsterPushAction(ai.DEPENDENCY_ALIVE, t)
	if monsterCanShoot(u) {
		u.MonsterPushAction(ai.DEPENDENCY_OBJECT_CLOSER_THAN, float32(float64(d.MissileAttackRange212)*.60000002), 0, t)
	}
	u.MonsterPushAction(ai.DEPENDENCY_CAN_SEE, t)
	u.MonsterPushAction(ai.ACTION_MELEE_ATTACK)
	u.MonsterPushAction(ai.ACTION_FACE_OBJECT, t)
	u.MonsterPushAction(ai.DEPENDENCY_OBJECT_FARTHER_THAN, d.MeleeAttackRange112, 0, t)
	if u.SubClass()&0x10 != 0 {
		u.MonsterPushAction(ai.DEPENDENCY_WAIT_FOR_STAMINA)
		u.MonsterPushAction(ai.DEPENDENCY_OR)
	}
	u.MonsterPushAction(ai.ACTION_MOVE_TO, t.PosVec, t)
}
func combatMissileChain(u, t *server.Object) {
	d := u.UpdateDataMonster().MonsterDef
	u.MonsterPushAction(ai.DEPENDENCY_NO_NEW_ENEMY, t)
	u.MonsterPushAction(ai.DEPENDENCY_CAN_SEE, t)
	u.MonsterPushAction(ai.ACTION_MISSILE_ATTACK, t.PosVec, t)
	u.MonsterPushAction(ai.ACTION_FACE_OBJECT, t)
	if !monsterRunningStatus(u) {
		u.MonsterPushAction(ai.DEPENDENCY_BLOCKED_LINE_OF_FIRE, t)
		u.MonsterPushAction(ai.DEPENDENCY_OBJECT_FARTHER_THAN, d.MissileAttackRange212, 0, t)
		u.MonsterPushAction(ai.DEPENDENCY_OR)
		u.MonsterPushAction(ai.ACTION_MOVE_TO, t.PosVec, t)
	}
}
func combatChoose(u, t *server.Object) {
	if u.Buffs&(1<<29) == 0 && monsterCastRelated2(u, t) {
		return
	}
	if monsterCanShoot(u) {
		if monsterCanMelee(u) && float64(stateDistance(u, t)) < float64(u.UpdateDataMonster().MonsterDef.MissileAttackRange212)*.5 {
			combatMeleeChain(u, t)
		} else {
			combatMissileChain(u, t)
		}
	} else if monsterCanMelee(u) {
		combatMeleeChain(u, t)
	} else if !monsterCanCast(u) {
		combatChase(u, t)
	}
}
func combatFightStart(u *server.Object) {
	ud := u.UpdateDataMonster()
	combatSound(u, 5)
	GetServer().NoxScriptC().ScriptCallback(&ud.ScriptChangeFocus, ud.CurrentEnemy, u, server.ScriptEventType(13))
	ud.StatusFlags |= 0x100
	visibilityFrameCopy(false)
	Nox_xxx_unitUpdateSightMB_5281F0(u)
	navigationStartRunning(u)
}
func combatFight(u *server.Object) {
	s := GetServer().S()
	ud := u.UpdateDataMonster()
	h := ud.AIStackHead()
	if s.Frame()-uint32(h.Args[2]) > 10*s.TickRate() {
		u.MonsterPopAction()
		return
	}
	if t := ud.CurrentEnemy; t != nil {
		h.Args[2] = uintptr(s.Frame())
		if visibilityKillable(t) == 0 {
			u.MonsterPopAction()
			return
		}
		if u.Buffs&(1<<29) != 0 || (!monsterBuffSelf(u) && !monsterCastOffensive(u, t)) {
			combatChoose(u, t)
		}
		return
	}
	if monsterRunningStatus(u) {
		u.MonsterPopAction()
		return
	}
	p := types.Pointf{X: math.Float32frombits(uint32(h.Args[0])), Y: math.Float32frombits(uint32(h.Args[1]))}
	found := memmap.PtrUint32(0x5D4594, 2487944)
	*found = 0
	s.Map.EachObjInCircle(p, 30, func(t *server.Object) bool {
		if t.NetCode == ud.Field300 && t.ObjFlags&0x8000 != 0 {
			*found = 1
		}
		return true
	})
	if *found != 0 {
		u.MonsterPopAction()
		if ud.Field98 == ud.Field300 {
			ud.Field97 = 0
		}
		return
	}
	dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
	if dy*dy+dx*dx < 64 {
		u.MonsterPopAction()
		return
	}
	u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
	u.MonsterPushAction(ai.ACTION_MOVE_TO, p, 0)
}
func combatBlock(u *server.Object) {
	s := GetServer().S()
	h := u.UpdateDataMonster().AIStackHead()
	if monsterShieldThreat(u) != nil {
		h.Args[0] = uintptr(s.Frame() + s.TickRate()/2)
	}
	if s.Frame() > uint32(h.Args[0]) {
		u.MonsterPopAction()
		if u.SubClass()&0x10 == 0 {
			u.MonsterPushAction(ai.ACTION_BLOCK_FINISH)
		}
	}
}
func combatClearBuffs(u *server.Object) {
	Nox_xxx_spellBuffOff_4FF5B0(u, 0)
	Nox_xxx_spellBuffOff_4FF5B0(u, 23)
}
func combatWait(u *server.Object, r float32) {
	ud := u.UpdateDataMonster()
	u.MonsterPushAction(ai.DEPENDENCY_OBJECT_CLOSER_THAN, float32(float64(r)*1.2), 0, ud.CurrentEnemy)
	u.MonsterPushAction(ai.ACTION_WAIT, ud.Field128)
}
func combatMeleeStart(u *server.Object) {
	s := GetServer().S()
	ud := u.UpdateDataMonster()
	d := ud.MonsterDef
	if s.Frame() < ud.Field128 {
		combatWait(u, d.MeleeAttackRange112)
		return
	}
	if u.SubClass()&0x10 != 0 {
		cost := int(controlWeaponStamina(uint32(int32(ud.Field514))))
		if cost > int(ud.Field282_0) {
			ud.Field282_0 -= byte(cost)
		} else {
			ud.Field282_0 = 0
		}
	}
	radius := math.Float32frombits(uint32(dword_587000_261388))
	dword_5d4594_2487948 = 0
	*memmap.PtrFloat32(0x5D4594, 2487952) = float32(float64(radius) + 1)
	p := u.PosVec
	s.Map.EachObjInRect(types.Rectf{Min: types.Pointf{X: p.X - radius, Y: p.Y - radius}, Max: types.Pointf{X: p.X + radius, Y: p.Y + radius}}, func(t *server.Object) bool { combatScan(t, u); return true })
	t := (*server.Object)(unsafe.Pointer(uintptr(dword_5d4594_2487948)))
	if t != nil && !s.IsEnemyTo(u, t) && u.SubClass()&0x10 != 0 && ud.Field516 != 0 && (*server.Object)(unsafe.Pointer(uintptr(ud.Field516))).SubClass()&0x4000 != 0 {
		combatDebug(u, "Tried to MELEE_ATTACK but friend in the way")
		u.MonsterPopAction()
		u.MonsterPushAction(ai.ACTION_FACE_ANGLE, int32(int16(u.Direction1)))
		if a := u.MonsterPushAction(ai.DEPENDENCY_TIME); a != nil {
			a.Args[0] = uintptr(s.Frame() + uint32(s.Rand.Logic.IntClamp(int(s.TickRate()/4), int(s.TickRate()/2))))
		}
		u.MonsterPushAction(ai.ACTION_FLEE, t.PosVec, 0)
		return
	}
	combatClearBuffs(u)
	u.Field34 = s.Frame()
	ud.Field128 = s.Frame() + uint32(s.Rand.Logic.IntClamp(int(d.MeleeAttackDelayMin128), int(d.MeleeAttackDelayMax132)))
	combatSound(u, 6)
}
func combatScan(t, u *server.Object) {
	if !t.Class().HasAny(object.ClassMonster|object.ClassPlayer) || t.ObjFlags&0x8000 != 0 {
		return
	}
	dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
	dy32 := float32(dy)
	// The C square precedes the float32 Y spill used by the facing test.
	length := float32(math.Sqrt(dy*dy+dx*dx) + .000099999997)
	min := memmap.PtrFloat32(0x5D4594, 2487952)
	if float64(length) < float64(*min) {
		dir := unsafe.Slice(memmap.PtrFloat32(0x587000, uintptr(194136+8*int(int16(u.Direction1)))), 2)
		if float64(dy32)/float64(length)*float64(dir[1])+dx/float64(length)*float64(dir[0]) > .5 {
			dword_5d4594_2487948 = uint32(uintptr(t.CObj()))
			*min = length
		}
	}
}
func combatMelee(u *server.Object) {
	ud := u.UpdateDataMonster()
	d := ud.MonsterDef
	if u.SubClass()&0x10 != 0 {
		if ud.StatusFlags&0x20000 != 0 {
			controlMorphToPlayer((*server.Object)(unsafe.Pointer(u.CObj())))
		}
		r := attackPlayer(u)
		if ud.StatusFlags&0x20000 != 0 {
			controlMorphFromPlayer((*server.Object)(unsafe.Pointer(u.CObj())))
		}
		if r == 0 {
			u.MonsterPopAction()
		}
		return
	}
	if d.MeleeStrikeFunc236 == nil {
		combatDebug(u, "Tried to MELEE_ATTACK but cannot")
		u.MonsterPopAction()
		return
	}
	if uint32(ud.Field120_1) == d.MeleeAttackFrame108 && ud.Field120_2 == 0 {
		hit := monsterCallbackResult(d.MeleeStrikeFunc236, u)
		if hit != 0 {
			combatSound(u, 8)
		} else {
			combatSound(u, 9)
		}
	}
	if ud.Field120_3 != 0 {
		u.MonsterPopAction()
	}
}
func combatMissileStart(u *server.Object) {
	s := GetServer().S()
	ud := u.UpdateDataMonster()
	d := ud.MonsterDef
	if s.Frame() < ud.Field128 {
		combatWait(u, d.MissileAttackRange212)
		return
	}
	combatClearBuffs(u)
	u.Field34 = s.Frame()
	ud.Field128 = s.Frame() + uint32(s.Rand.Logic.IntClamp(int(d.MissileAttackDelayMin220), int(d.MissileAttackDelayMax224)))
	combatSound(u, 10)
}
func combatMissile(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.SubClass()&0x10 != 0 {
		if attackPlayer(u) == 0 {
			u.MonsterPopAction()
		}
		return
	}
	d := ud.MonsterDef
	h := ud.AIStackHead()
	if uint32(ud.Field120_1) == d.MissileAttackFrame216 && ud.Field120_2 == 0 {
		name := d.MissileName148[:]
		for i, c := range name {
			if c == 0 {
				name = name[:i]
				break
			}
		}
		if proj := GetServer().S().NewObjectByTypeID(string(name)); proj != nil {
			p := types.Pointf{X: math.Float32frombits(uint32(h.Args[0])), Y: math.Float32frombits(uint32(h.Args[1]))}
			if h.Args[2] != 0 {
				spatialPredict(u, motionObject(uint32(h.Args[2])), proj.SpeedCur, &p)
			}
			dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
			// The compiled C retains both deltas in x87 through length and velocity.
			length := float32(math.Sqrt(dy*dy+dx*dx) + .1)
			radius := float32(float64(u.Shape.Circle.R) + 4)
			vel := types.Pointf{X: float32(dx * float64(proj.SpeedCur) / float64(length)), Y: float32(dy * float64(proj.SpeedCur) / float64(length))}
			dir := unsafe.Slice(memmap.PtrFloat32(0x587000, uintptr(194136+8*int(int16(u.Direction1)))), 2)
			spawnY := float64(radius)*float64(dir[1]) + float64(u.PosVec.Y)
			spawn := types.Pointf{X: float32(float64(radius)*float64(dir[0]) + float64(u.PosVec.X)), Y: float32(spawnY)}
			// C rounds spawn X before the ray addition, but keeps spawn Y wide.
			end := types.Pointf{X: spawn.X + vel.X, Y: float32(spawnY + float64(vel.Y))}
			if GetServer().S().MapTraceRayAt(u.PosVec, end, nil, nil, 5) {
				GetServer().CreateObjectAt(proj, u, spawn)
				proj.VelVec = vel
				proj.Direction1 = u.Direction1
				proj.Direction2 = u.Direction1
			} else {
				GetServer().DelayedDelete(proj)
			}
		}
		combatSound(u, 11)
	}
	if ud.Field120_3 != 0 {
		u.MonsterPopAction()
	}
}

func combatDebug(u *server.Object, message string) {
	if noxflags.HasEngine(noxflags.EngineShowAI) {
		ai.Log.Printf("%d: %s(#%d) %s\n", GetServer().S().Frame(), GetServer().S().Types.ByInd(int(u.TypeInd)).ID(), u.NetCode, message)
	}
}
