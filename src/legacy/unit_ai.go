package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_gameSetAudioFadeoutMb_501AC0 func(v int)
	Nox_xxx_unitUpdateMonster_50A5C0     func(a1 *server.Object)
)

func nox_xxx_monsterPopAction_50A160(a1 *nox_object_t) int {
	return asObjectS(a1).MonsterPopAction()
}

func sub_50B510() {
	GetServer().S().AI.Paths.Sub_50B510()
}

func Nox_xxx_mobSearchEdible_544A00(a1 *server.Object, a2 float32) int {
	if u := lifecycleFoodSearch(a1, a2, false); u != nil {
		return int(uintptr(u.CObj()))
	}
	return 0
}
func Nox_xxx_weaponGetStaminaByType_4F7E80(a1 int) int {
	return int(controlWeaponStamina(uint32(a1)))
}
func Nox_xxx_mobGetMoveAttemptTime_534810(a1 *server.Object) int {
	return bool2int(monsterMoveAttempt(a1))
}
func Nox_xxx_unitIsDangerous_547120(a1 *server.Object, a2 *server.Object) {
	monsterDangerous(a1, a2)
}
func Nox_xxx_checkIsKillable_528190(a1 *server.Object) int {
	return int(visibilityKillable(a1))
}

// PolygonLevelAt returns the level of the map polygon at pos, if there is one.
func PolygonLevelAt(pos image.Point) (byte, bool) {
	p := mapPolygonFind(&[2]int32{int32(pos.X), int32(pos.Y)}, 0, false)
	if p == nil {
		return 0, false
	}
	return p.Level, true
}
func Nox_xxx_mobAction_50A910(a1 *server.Object) {
	monsterControlRefresh(a1)
}
func Nox_xxx_monsterGetSoundSet_424300(a1 *server.Object) *server.MonsterSoundSet {
	return resourceMonsterSound(a1)
}
func Nox_xxx_monsterPlayHurtSound_532800(a1 *server.Object) {
	unitHurtSound(a1)
}
func Nox_xxx_mobAction_5469B0(a1 *server.Object) {
	monsterIdleAudio(a1)
}
func Nox_xxx_unitUpdateSightMB_5281F0(a1 *server.Object) {
	visibilityUpdateSight(a1)
}
func Nox_xxx_monsterMainAIFn_547210(a1 *server.Object) {
	monsterMainAI(a1)
}
func Nox_xxx_updateNPCAnimData_50A850(a1 *server.Object) {
	monsterControlAnimation(a1)
}
func Nox_xxx_monsterPolygonEnter_421FF0(a1 *server.Object) {
	mapPolygonMonster(a1)
}
func Nox_xxx_monsterMimicCheckMorph_534950(a1 *server.Object) {
	monsterMimicMorph(a1)
}
func Sub_5466F0(a1 *server.Object) int {
	return investigateHeardSound(a1)
}
func Nox_xxx_mobHealSomeone_5411A0(a1 *server.Object) {
	monsterHealSomeone(a1)
}
func Nox_xxx_mobActionCast_5413B0(a1 *server.Object, a2 int) {
	monsterActionCast(a1, a2)
}
