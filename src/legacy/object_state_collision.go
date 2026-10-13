package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
)

func stateCloseDoor(u *server.Object, p unsafe.Pointer) {
	if u.ObjClass&0x80 != 0 && u.UpdateDataDoor().GridX == *equipmentWord(p, 0) && u.UpdateDataDoor().GridY == *equipmentWord(p, 4) {
		u.UpdateDataDoor().KeyType = 0
		if bool(nox_common_gameFlags_check_40A5C0(4096)) {
			stateDoorNotify(u)
		}
	}
}
func stateDoorNotify(u *server.Object) int32 {
	u.UpdateDataDoor().Notified = 1
	return int32(questRuntimeCodeMessage(255, u))
}
func stateMonsterCollision(u, t *server.Object) unsafe.Pointer {
	return nox_xxx_scriptCallByEventBlock_502490(unsafe.Pointer(&layout.As[server.MonsterUpdateData](u.UpdateData).ScriptCollision), t.CObj(), u.CObj(), 22)
}
func stateMimicCollision(u, t *server.Object) unsafe.Pointer {
	if t != nil && t.ObjFlags&0x8000 == 0 && t.ObjClass&6 != 0 && nox_xxx_unitIsEnemyTo_5330C0(asObjectC(u), asObjectC(t)) != 0 && !u.MonsterActionIsScheduled(15) {
		if st := u.MonsterPushAction(43); st != nil {
			st.SetArgs(GetServer().S().Frame())
		}
		if st := u.MonsterPushAction(15); st != nil {
			st.SetArgs(t.PosVec, GetServer().S().Frame())
		}
	}
	return stateMonsterCollision(u, t)
}
func stateChargeImpact(u, t *server.Object) bool {
	if t == nil {
		return true
	}
	class := t.ObjClass
	if class&6 == 0 || t.HealthData.Cur == 0 && t.HealthData.Max != 0 {
		if class&0x400000 == 0 && int8(t.ObjFlags) >= 0 && float64(t.Mass) <= float64(u.Mass) {
			return false
		}
	}
	if class&0x80 != 0 {
		return t.UpdateDataDoor().KeyType != 0
	}
	return t.ObjFlags&9 == 0 && class&1 == 0
}
func stateChargeMoveBack(u *server.Object) {
	Nox_xxx_unitMove_4E7010(u, u.PrevPos)
}
func stateChargeStun(u *server.Object) {
	duration := floatToInt32(float32(nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerStunDuration"))))
	spellLifeApplyBuff(u, 5, int16(duration), 5)
}
func stateCharge(u, t *server.Object) {
	Nox_xxx_playerSetState_4FA020(u, 13)
	gameplayReportEarthquake(&u.PosVec, 10)
	Sub_4FC300(u, 1)
	if t != nil {
		damage := floatToInt32(float32(nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerDamage"))))
		if t.ObjClass&0x400000 == 0 {
			worldCollideMass(u, t)
		}
		projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 2)
		if t.ObjClass&0x20006 != 0 {
			stateChargeMoveBack(u)
			return
		}
		stateChargeStun(u)
	} else {
		wall := (*server.Wall)(layout.As[server.PlayerUpdateData](u.UpdateData).Field74)
		if wall != nil && nox_xxx_wallFlags(int(wall.Tile1))&5 == 0 {
			stateChargeMoveBack(u)
			return
		}
		GetServer().S().Audio.EventObj(sound.ID(171), u, 0, 0)
		stateChargeStun(u)
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, 100, 2)
	}
	pain := floatToInt32(float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerPainRatio"))) * float64(u.HealthData.Cur)))
	if pain < 1 {
		pain = 1
	}
	resourceDamage(u, pain)
	stateChargeMoveBack(u)
}
func statePlayerCollision(u, t *server.Object) {
	if GetServer().S().Abils.IsActive(u, 1) && stateChargeImpact(u, t) {
		stateCharge(u, t)
	}
	if t != nil && t.ObjClass&4 != 0 && t.ObjFlags&0x8000 == 0 && u.Buffs&(1<<16) != 0 && uint32(u.BuffsDur[16]) < 14*GetServer().S().TickRate() {
		spellLifeApplyBuff(t, 16, int16(uint16(15*GetServer().S().TickRate())), int8(u.BuffsPower[16]))
		spellLifeBuffOff(u, int32(16))
	}
}
