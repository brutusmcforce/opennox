package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_server_questAllowDefault              func() bool
	Nox_server_questNextStageThreshold_4D74F0 func(lvl int) int
	Sub_4D6F30                                func() int
)

func sub_4D6F30() int { return Sub_4D6F30() }
func Nox_game_getQuestStage_4E3CC0() int {
	return int(questRuntimeStage())
}
func Nox_xxx_getQuestStage_51A930() int {
	return int(questProgressStage())
}
func Sub_4E3D50() {
	questRuntimeDifficulty()
}
func Sub_4E3DD0() {
	questRuntimeScaleHealth()
}
func Sub_4F1F20() {
	orchestrationRewards()
}
func Nox_xxx_getQuestMapFile_4D0F60() string {
	return GoStringP(unsafe.Pointer(mapQuestChoose()))
}
func Sub_4D0F30() {
	mapQuestReset()
}
func Nox_game_setQuestStage_4E3CD0(a1 int) {
	questRuntimeSetStage(uint32(a1))
}
func Sub_51A920(a1 int) {
	questProgressSetStage(uint32(a1))
}
func Sub_51A1F0(a1 int) {
	questProgressPrepare(a1)
}
func Sub_4D10F0(a1 string) {
	mapQuestPlayed((*byte)(unsafe.Pointer(internCStr(a1))))
}
func Sub_4D7520(a1 int) {
	questRuntimeGateSet(uint32(a1))
}
func Sub_4D9CF0(a1 int) {
	gameplayReportQuestStart(int(int32(a1)))
}
func Sub_4D6000(a1 *server.Object) {
	questRuntimeReset(a1)
}
func Sub_4D9D20(a1 int, a2 *server.Object) {
	gameplayReportQuestObject(int(int32(a1)), a2)
}
func Nox_xxx_unitInitPlayer_4EFE80(a1 *server.Object) {
	controlInitPlayer(a1)
}
func Nox_xxx_playerRespawnItem_4EF750(a1 *server.Object, a2 string, a3 unsafe.Pointer, a4 int, a5 int) *server.Object {
	return controlRespawnItem(a1, GoStringP(unsafe.Pointer(internCStr(a2))), a3, int32(a4), int32(a5))
}
func Nox_xxx_modifSetItemAttrs_4E4990(a1 *server.Object, a2 unsafe.Pointer) {
	stateAttributes(a1, a2)
}
func Sub_4184D0(a1 unsafe.Pointer) {
	teamRuntimeAnnounce(asTeamP(a1))
}
func Nox_xxx_createAtImpl_4191D0(a1 server.TeamID, a2 *server.ObjectTeam, a3 int, a4 int, a5 int) {
	teamRuntimeJoin(a1, a2, a3, a4, a5)
}
func Nox_xxx_cliPlrInfoLoadFromFile_41A2E0(a1 string, a2 ntype.PlayerInd) int {
	return playerFileServerLoad(a1, int(a2))
}
