package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_unitIsUnitTT_4E7C80 func(a1 *server.Object, a2 int) int
	Nox_xxx_unitMove_4E7010     func(a1 *server.Object, pos types.Pointf)
)

func Nox_xxx_unitSetHP_4E4560(a1 *server.Object, a2 uint16) {
	resourceSetHP(a1, a2)
}
func Nox_xxx_mobInformOwnerHP_4EE4C0(a1 *server.Object) {
	resourceInformOwner(a1)
}
func Nox_xxx_protectMana_56F9E0(a1 int, a2 int16) {
	addProtectionRecord(int32(a1), uint32(int32(a2)))
}
func Nox_xxx_monsterWalkTo_514110(a1 *server.Object, a2 float32, a3 float32) {
	monsterControlWalk(a1, types.Pointf{X: a2, Y: a3})
}
func Nox_xxx_monsterLookAt_5125A0(a1 *server.Object, a2 int) {
	monsterControlLook(a1, int32(a2))
}
func Nox_xxx_unitFreeze_4E79C0(a1 *server.Object, a2 int) {
	stateFreeze(a1, int32(a2))
}
func Nox_xxx_unitUnFreeze_4E7A60(a1 *server.Object, a2 int) {
	stateUnfreeze(a1, int32(a2))
}
func Nox_xxx_scriptMonsterRoam_512930(a1 *server.Object) {
	scriptBindingRoam(a1)
}
func Nox_server_gotoHome(a1 *server.Object) {
	scriptBindingHome(a1)
}
func Nox_xxx_unitIdle_515820(a1 *server.Object) {
	monsterControlIdle(a1, false)
}
func Nox_xxx_unitSetFollow_5158C0(a1 *server.Object, a2 *server.Object) {
	monsterControlFollow(a1, a2)
}
func Nox_xxx_unitHunt_5157A0(a1 *server.Object) {
	monsterControlIdle(a1, true)
}
func Nox_xxx_playerSubGold_4FA5D0(a1 *server.Object, a2 int) {
	resourceSubGold(a1, uint32(a2))
}
func Nox_xxx_playerAddGold_4FA590(a1 *server.Object, a2 int) {
	resourceAddGold(a1, uint32(a2))
}
