package legacy

import (
	"math"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

type movementAIAction struct {
	typ         ai.ActionType
	update      func(*server.Object)
	popOnCancel bool
}

func (a movementAIAction) Type() ai.ActionType     { return a.typ }
func (a movementAIAction) Start(*server.Object)    {}
func (a movementAIAction) End(*server.Object)      {}
func (a movementAIAction) Update(u *server.Object) { a.update(u) }
func (a movementAIAction) Cancel(u *server.Object) {
	if a.popOnCancel {
		u.MonsterPopAction()
	}
}
func init() {
	for _, a := range []movementAIAction{
		{ai.ACTION_FACE_LOCATION, facingLocation, true},
		{ai.ACTION_FACE_OBJECT, facingObject, true},
		{ai.ACTION_FACE_ANGLE, facingAngle, true},
		{ai.ACTION_SET_ANGLE, setFacingAngle, true},
		{ai.ACTION_RANDOM_WALK, randomWalk, false},
		{ai.ACTION_CONFUSED, confusedMovement, false},
	} {
		server.RegisterAIAction(a)
	}
}

func movementDirectionVector(dir int32) (float32, float32) {
	offset := uintptr(dir * 8)
	v := (*[2]float32)(memmap.PtrOff(0x587000, 194136+offset))
	return v[0], v[1]
}
func facingDot(u *server.Object, v types.Pointf) bool {
	cos, sin := movementDirectionVector(int32(int16(u.Direction1)))
	// C spills the first product across the second table lookup.
	first := float32(float64(sin) * float64(v.Y))
	return float64(cos)*float64(v.X)+float64(first) > 0.89999998
}
func facingTurn(u *server.Object, cross float64, v types.Pointf) {
	if cross >= 0 {
		u.Direction2 = server.Dir16(uint8(u.Direction2 + 8))
	} else {
		u.Direction2 = server.Dir16(uint8(u.Direction2 - 8))
	}
	if facingDot(u, v) {
		u.MonsterPopAction()
	}
}
func facingLocationTo(u *server.Object, target types.Pointf) {
	dx, dy := float64(target.X)-float64(u.PosVec.X), float64(target.Y)-float64(u.PosVec.Y)
	length := float32(math.Sqrt(dx*dx+dy*dy) + 0.001)
	nx, ny := float32(dx/float64(length)), dy/float64(length)
	cos, sin := movementDirectionVector(int32(int16(u.Direction1)))
	// The turn retains normalized Y as double; the dot helper receives float32.
	facingTurn(u, ny*float64(cos)-float64(nx)*float64(sin), types.Pointf{X: nx, Y: float32(ny)})
}
func facingLocation(u *server.Object) {
	facingLocationTo(u, u.UpdateDataMonster().AIStackHead().ArgPos(0))
}
func facingObject(u *server.Object) {
	if target := u.UpdateDataMonster().AIStackHead().ArgObj(0); target != nil {
		facingLocationTo(u, target.PosVec)
	} else {
		u.MonsterPopAction()
	}
}
func facingAngle(u *server.Object) {
	x, y := movementDirectionVector(int32(u.UpdateDataMonster().AIStackHead().ArgU32(0)))
	cos, sin := movementDirectionVector(int32(int16(u.Direction1)))
	facingTurn(u, float64(y)*float64(cos)-float64(sin)*float64(x), types.Pointf{X: x, Y: y})
}
func setFacingAngle(u *server.Object) {
	dir := server.Dir16(uint8(u.UpdateDataMonster().AIStackHead().ArgU32(0)))
	u.Direction1, u.Direction2 = dir, dir
	u.MonsterPopAction()
}
func randomWalkDirection(u *server.Object) server.Dir16 {
	dir := uint8(int32(int16(u.Direction1)) + int32(nox_common_randomInt_415FA0(-20, 20)))
	if uint32(u.ObjSubClass)&0x400 == 0 {
		x, y := movementDirectionVector(int32(dir))
		probe := types.Pointf{X: float32(float64(x)*30 + float64(u.PosVec.X)), Y: float32(float64(y)*30 + float64(u.PosVec.Y))}
		if tileAtPoint(probe) == 6 {
			dir += 64
		}
	}
	return server.Dir16(dir)
}
func randomWalk(u *server.Object) {
	ud := u.UpdateDataMonster()
	dir := randomWalkDirection(u)
	u.Direction1, u.Direction2 = dir, dir
	speed := float64(u.SpeedCur)
	if uint32(ud.StatusFlags)&0x4000 != 0 {
		speed *= float64(ud.MonsterDef.RunMultiplier96)
	}
	// C spills the selected running speed to float32 before both products.
	speed = float64(float32(speed))
	x, y := movementDirectionVector(int32(dir))
	u.ForceVec = types.Pointf{X: float32(speed * float64(x)), Y: float32(speed * float64(y))}
	monsterMoveAudio(u)
}
func confusedMovement(u *server.Object) {
	if nox_common_randomInt_415FA0(0, 100) >= 15 {
		randomWalk(u)
		return
	}
	if monsterCanMelee(u) {
		if !monsterCanShoot(u) || nox_common_randomInt_415FA0(0, 100) < 50 {
			u.MonsterPushActionImpl(ai.ActionType(16), "go", 0)
			return
		}
	} else if !monsterCanShoot(u) {
		return
	}
	if st := u.MonsterPushActionImpl(ai.ActionType(17), "go", 0); st != nil {
		x, y := movementDirectionVector(int32(int16(u.Direction1)))
		st.SetArgs(types.Pointf{X: float32(float64(x)*10 + float64(u.PosVec.X)), Y: float32(float64(y)*10 + float64(u.PosVec.Y))}, uint32(0))
	}
}
