package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func monsterControlHead(u *server.Object) uint32 {
	ud := u.UpdateDataMonster()
	// The legacy empty-head read is a raw word before the first stack item.
	return *(*uint32)(unsafe.Add(u.UpdateData, layoutServerMonsterUpdateData.Off(uintptr(552+24*int(ud.AIStackInd)))))
}
func monsterControlEnsure(u *server.Object, action uint32) *server.AIStackItem {
	if u.ObjClass&2 != 0 && monsterControlHead(u) != action {
		return u.MonsterPushAction(ai.ActionType(action))
	}
	return nil
}
func monsterControlAlive(u *server.Object) bool {
	return u != nil && u.ObjClass&2 != 0 && u.ObjFlags&0x8000 == 0
}
func monsterControlLook(u *server.Object, index int32) uint32 {
	angle := geometryDirection4Angle(index)
	if !monsterControlAlive(u) {
		return uint32(angle)
	}
	x := float32(float64(memmap.Float32(0x587000, uintptr(194136+8*angle)))*10 + float64(u.PosVec.X))
	y := float32(float64(memmap.Float32(0x587000, uintptr(194140+8*angle)))*10 + float64(u.PosVec.Y))
	p := u.MonsterPushAction(ai.ACTION_FACE_LOCATION, x, y)
	return memmap.PtrWord(p.C())
}
func monsterControlWalk(u *server.Object, p types.Pointf) {
	if !monsterControlAlive(u) {
		return
	}
	u.ClearActionStack()
	u.MonsterPushAction(ai.ACTION_REPORT, 8)
	u.MonsterPushAction(ai.ACTION_FAR_MOVE_TO, p, 0)
}
func monsterControlPatrol(u *server.Object, p, q types.Pointf, dist float32) {
	if !monsterControlAlive(u) {
		return
	}
	ud := u.UpdateDataMonster()
	delta := types.Pointf{X: float32(q.X - p.X), Y: float32(q.Y - p.Y)}
	u.ClearActionStack()
	if a := u.MonsterPushAction(ai.ACTION_GUARD); a != nil {
		a.SetArgs(p, geometryVectorAngle(&delta))
	}
	ud.SightRange = dist
}
func monsterControlIdle(u *server.Object, hunt bool) {
	if !monsterControlAlive(u) {
		return
	}
	u.ClearActionStack()
	act := ai.ACTION_IDLE
	if hunt {
		act = ai.ACTION_HUNT
	}
	u.MonsterPushAction(act)
}
func monsterControlFollow(u, t *server.Object) {
	if !monsterControlAlive(u) || t == nil || u == t {
		return
	}
	u.ClearActionStack()
	u.MonsterPushAction(ai.ACTION_ESCORT, t.PosVec, t)
}
func monsterControlMelee(u *server.Object, p *types.Pointf) {
	if !monsterControlAlive(u) || !monsterCanMelee(u) {
		return
	}
	u.ClearActionStack()
	u.MonsterPushAction(ai.ACTION_REPORT, 16)
	u.MonsterPushAction(ai.ACTION_MELEE_ATTACK)
	r := float32(float64(u.UpdateDataMonster().MonsterDef.MeleeAttackRange112) + float64(u.Shape.Circle.R))
	u.MonsterPushAction(ai.DEPENDENCY_LOCATION_FARTHER_THAN, r, 0, *p)
	u.MonsterPushAction(ai.ACTION_MOVE_TO, *p, 0)
}
func monsterControlMissile(u *server.Object, p *types.Pointf) {
	if !monsterControlAlive(u) || p == nil || !monsterCanShoot(u) {
		return
	}
	u.ClearActionStack()
	u.MonsterPushAction(ai.ACTION_REPORT, 17)
	u.MonsterPushAction(ai.ACTION_MISSILE_ATTACK, *p, 0)
}
func monsterControlByte(u *server.Object, b *byte) uint32 {
	if u == nil {
		return 0
	}
	if u.ObjClass&2 == 0 {
		return motionAddress(u)
	}
	ud := u.UpdateDataMonster()
	ud.Field333 = ud.Field333&^255 | uint32(*b)
	return memmap.PtrWord(u.UpdateData)
}
func monsterControlFight(u, t *server.Object) {
	if !monsterControlAlive(u) || t == nil || u == t {
		return
	}
	u.ClearActionStack()
	u.UpdateDataMonster().Field304 = t
	visibilityFrameCopy(true)
	u.MonsterPushAction(ai.ACTION_REPORT, 15)
	u.MonsterPushAction(ai.ACTION_FIGHT, t.PosVec, GetServer().S().Frame())
}
func monsterControlFlee(u *server.Object, from types.Pointf, duration uint32) {
	if !monsterControlAlive(u) || u.UpdateDataMonster().HasAction(ai.ACTION_FLEE) {
		return
	}
	u.MonsterPushAction(ai.ACTION_REPORT, 24)
	u.MonsterPushAction(ai.DEPENDENCY_TIME, GetServer().S().Frame()+duration)
	u.MonsterPushAction(ai.ACTION_FLEE, from, 0)
}
func monsterControlWait(u *server.Object, duration uint32) {
	if !monsterControlAlive(u) {
		return
	}
	u.MonsterPushAction(ai.ACTION_REPORT, 1)
	u.MonsterPushAction(ai.ACTION_WAIT, GetServer().S().Frame()+duration)
}
