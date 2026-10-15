package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

type lifecycleAIAction struct{ typ ai.ActionType }

func (a lifecycleAIAction) Type() ai.ActionType { return a.typ }
func (a lifecycleAIAction) Start(u *server.Object) {
	switch a.typ {
	case ai.ACTION_DYING:
		lifecycleDyingStart(u)
	case ai.ACTION_DEAD:
		lifecycleDeadStart(u)
	}
}
func (a lifecycleAIAction) Update(u *server.Object) {
	switch a.typ {
	case ai.ACTION_HUNT:
		lifecycleHunt(u)
	case ai.ACTION_PICKUP_OBJECT:
		lifecyclePickup(u)
	case ai.ACTION_DYING, ai.ACTION_GET_UP:
		if u.UpdateDataMonster().Field120_3 != 0 {
			u.MonsterPopAction()
		}
	case ai.ACTION_DEAD:
		lifecycleDeadUpdate(u)
	}
}
func (a lifecycleAIAction) End(u *server.Object) {
	if a.typ == ai.ACTION_DYING {
		lifecycleBurnDeleteCheck(u)
	}
}
func (a lifecycleAIAction) Cancel(*server.Object) {}

func init() {
	for _, typ := range []ai.ActionType{ai.ACTION_HUNT, ai.ACTION_PICKUP_OBJECT, ai.ACTION_DYING, ai.ACTION_DEAD, ai.ACTION_GET_UP} {
		server.RegisterAIAction(lifecycleAIAction{typ})
	}
}

func lifecycleHunt(u *server.Object) {
	u.UpdateDataMonster().Aggression = math.Float32frombits(0x3f547ae1)
	if st := u.MonsterPushAction(ai.ACTION_ROAM); st != nil {
		st.Args[0] = 0
		st.Args[2] = st.Args[2]&^0xff | 0x80
	}
}
func lifecyclePickup(u *server.Object) {
	h := u.UpdateDataMonster().AIStackHead()
	if t := h.ArgObj(0); t != nil {
		dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		if dx*dx+dy*dy < 5625 && GetServer().S().CanInteract(u, t, 0) {
			Nox_xxx_inventoryServPlace_4F36F0(u, t, 1, 1)
			// C reloads the action argument after inventory placement.
			if t = h.ArgObj(0); t != nil && t.SubClass()&0x10 != 0 {
				effectsUse(u, t)
			}
		}
	}
	u.MonsterPopAction()
}
func lifecycleDyingStart(u *server.Object) {
	ud := u.UpdateDataMonster()
	combatSound(u, 15)
	// The game passes no caller. Go scripts get the killer instead, e.g. to keep score (opennox#678);
	// the map's NoxScript keeps the original behavior.
	var killer *server.Object
	if GetServer().S().NoxScriptVM.IsGoFunc(ud.ScriptDeath.Func) {
		killer = lifecycleKiller(u)
	}
	GetServer().NoxScriptC().ScriptCallback(&ud.ScriptDeath, killer, u, server.ScriptEventType(7))
	if p := ud.MonsterDef.DieFunc228; p != nil {
		monsterCallbackResult(p, u)
	}
}
// lifecycleKiller returns the unit that dealt the last damage to u: the topmost unit owning the source (the caster
// of a missile, the master of a summoned creature), or the source itself if no unit owns it (e.g. a trap).
func lifecycleKiller(u *server.Object) *server.Object {
	src := u.Obj130
	killer := src
	for it := src; it != nil; it = it.ObjOwner {
		if it.Class().HasAny(object.MaskUnits) {
			killer = it
		}
	}
	if killer == u {
		return nil
	}
	return killer
}
func lifecycleIsZombie(u *server.Object) bool {
	return monsterIsZombie(u)
}
func lifecycleBurnDeleteCheck(u *server.Object) {
	if lifecycleIsZombie(u) && u.UpdateDataMonster().StatusFlags&0x80000 != 0 {
		lifecycleBurnDelete(u)
	}
}
func lifecycleBurnDelete(u *server.Object) {
	if owner := u.ObjOwner; owner != nil && owner.Class().Has(object.ClassPlayer) {
		pl := owner.UpdateDataPlayer().Player
		u.ObjSubClass &^= 0x80
		gameplayReportShield(int(pl.PlayerInd), u)
		GetServer().S().Players.Nox_xxx_netUnmarkMinimapObj_417300(ntype.PlayerInd(pl.PlayerInd), u, 1)
	}
	unitMonsterReward(u)
	Nox_xxx_sMakeScorch_537AF0(u.PosVec, 1)
	GetServer().DelayedDelete(u)
}
func lifecycleDeadStart(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.Field131 == 14 && u.SubClass()&0x10000 != 0 {
		lifecycleReleasedSoul(u)
	}
	u.VelVec = types.Pointf{}
	u.ForceVec = types.Pointf{}
	u.Pos24 = types.Pointf{}
	if p := ud.MonsterDef.DeadFunc232; p != nil {
		monsterCallbackDiscard(p, u)
	}
	if lifecycleIsZombie(u) {
		lo := floatToInt32(float32(GetServer().S().Balance.FloatInd("ZombieDeadDuration", 0)))
		hi := floatToInt32(float32(GetServer().S().Balance.FloatInd("ZombieDeadDuration", 1)))
		ud.Field123 = uint32(nox_common_randomInt_415FA0(int(lo), int(hi)))
		u.ObjFlags |= 0x10
	} else {
		u.ObjFlags |= 0x18
	}
}
func lifecycleReleasedSoul(u *server.Object) {
	id := memmap.PtrUint32(0x5D4594, 2489456)
	if *id == 0 {
		*id = uint32(GetServer().S().Types.IndByID("ReleasedSoul"))
	}
	if soul := GetServer().S().NewObjectByTypeInd(int(*id)); soul != nil {
		GetServer().CreateObjectAt(soul, nil, u.PosVec)
		soul.Direction1, soul.Direction2 = u.Direction1, u.Direction1
	}
}
func lifecycleDeadUpdate(u *server.Object) {
	ud := u.UpdateDataMonster()
	if lifecycleIsZombie(u) {
		if ud.StatusFlags&0x80000 != 0 {
			visibilityFXSpark(u.PosVec, 100)
			lifecycleBurnDelete(u)
		} else if GetServer().S().Frame()-ud.Field137 > ud.Field123 && ud.StatusFlags&0x100000 == 0 && ud.CurrentEnemy != nil {
			lifecycleRaiseZombie(u)
		}
		return
	}
	GetServer().S().Objs.RemoveFromUpdatable(u)
	u.Update = nil // C clears Object+744; Object+748 UpdateData remains live.
	lifecycleReset(u)
	if ud.MonsterDef.StatusFlags92&1 != 0 {
		GetServer().DelayedDelete(u)
	}
}
func lifecycleReset(u *server.Object) {
	ud := u.UpdateDataMonster()
	ud.Field523_2 = 0
	if n := int32(ud.Field74); n > 0 {
		clear(unsafe.Slice(&ud.Field75, int(n)))
	}
	ud.Field74, ud.Field91 = 0, 0
	if n := int(ud.Field282_1); n > 0 {
		clear(unsafe.Slice(&ud.Field283, n))
	}
	ud.Field282_1 = 0
	ud.CurrentEnemy = nil
}
func lifecycleRaiseZombie(u *server.Object) uint32 {
	if !lifecycleIsZombie(u) {
		return 0
	}
	if got := u.UpdateDataMonster().AIStackHead().Type(); got != ai.ACTION_DEAD {
		return uint32(got)
	}
	u.MonsterPopAction()
	u.MonsterPushAction(ai.DEPENDENCY_UNINTERRUPTABLE)
	u.MonsterPushAction(ai.ACTION_GET_UP)
	GetServer().S().Audio.EventObj(469, u, 0, 0)
	resourceRestoreHP(u)
	u.ObjFlags &= 0xffff7fa7
	return uint32(u.ObjFlags)
}

// lifecycleFood holds the best food candidates found by lifecycleFoodSearch (usable and other).
// Was 0x5D4594+2489440 and +2489452.
var lifecycleFood [2]*server.Object

func lifecycleCandidate(slot int) *server.Object         { return lifecycleFood[slot] }
func lifecycleSetCandidate(slot int, u *server.Object) { lifecycleFood[slot] = u }
func lifecycleFoodSearch(u *server.Object, radius float32, usable bool) *server.Object {
	var candidateOff int
	var minOff uintptr
	if usable {
		candidateOff, minOff = 0, 2489448
	} else {
		candidateOff, minOff = 1, 2489444
	}
	lifecycleFood[candidateOff] = nil
	*memmap.PtrFloat32(0x5D4594, minOff) = 10000000
	GetServer().S().Map.EachObjInCircle(u.PosVec, radius, func(t *server.Object) bool {
		if usable {
			if !t.Class().Has(0x1000000) || !Nox_xxx_playerClassCanUseItem_57B3D0(t, player.Class(0)) {
				return true
			}
		} else {
			sub := uint32(t.SubClass())
			if !t.Class().Has(0x10) || sub&4 != 0 || (sub&0x80 != 0 && u.Poison540 == 0) ||
				(sub&8 != 0 && (sub&0x10 == 0 || !noxflags.HasGame(0x2000) || u.SubClass()&0x10 == 0)) {
				return true
			}
		}
		if !GetServer().S().CanInteract(u, t, 0) {
			return true
		}
		dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		if dist := dx*dx + dy*dy; dist < float64(*memmap.PtrFloat32(0x5D4594, minOff)) {
			*memmap.PtrFloat32(0x5D4594, minOff) = float32(dist)
			lifecycleSetCandidate(candidateOff, t)
		}
		return true
	})
	return lifecycleCandidate(candidateOff)
}
