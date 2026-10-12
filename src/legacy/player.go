package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_playerDisconnByPlrID_4DEB00 func(id ntype.PlayerInd)
	Nox_xxx_playerCallDisconnect_4DEAB0 func(ind ntype.PlayerInd, v int8)
	Nox_xxx_playerCameraUnlock_4E6040   func(cplayer *server.Object)
	Nox_xxx_playerCameraFollow_4E6060   func(cplayer, cunit *server.Object)
	Nox_xxx_playerGetPossess_4DDF30     func(cplayer *server.Object) *server.Object
	Nox_xxx_playerGoObserver_4E6860     func(pl *server.Player, a2 int, a3 int) int
	Nox_xxx_playerObserveClear_4DDEF0   func(cplayer *server.Object)
	Nox_xxx_playerObserveMonster_4DDE80 func(cplayer, cunit *server.Object)
)

type nox_playerInfo = server.Player

func asPlayerS(p *nox_playerInfo) *server.Player {
	return (*server.Player)(unsafe.Pointer(p))
}

func AsPlayerP(p unsafe.Pointer) *server.Player {
	return (*server.Player)(p)
}

var _ = [1]struct{}{}[16-unsafe.Sizeof(server.ClassStats{})]

func nox_common_playerInfoGetByID_417040(id int) *nox_playerInfo {
	return (*nox_playerInfo)(GetServer().S().Players.ByID(id).C())
}

func nox_xxx_playerCameraUnlock_4E6040(cplayer *nox_object_t) {
	Nox_xxx_playerCameraUnlock_4E6040(asObjectS(cplayer))
}

func nox_xxx_playerCameraFollow_4E6060(cplayer, cunit *nox_object_t) {
	Nox_xxx_playerCameraFollow_4E6060(asObjectS(cplayer), asObjectS(cunit))
}

func Nox_xxx_scavengerTreasureMax_4D1600() uint32 {
	return uint32(sessionScavengerMaximum())
}

func Nox_xxx_netMsgFadeBeginPlayer(ind int, dir int, a3 int) {
	gameplayReportFade(int(ind), int(dir), int(a3))
}

func PrintToPlayers(text string) {
	cstr, free := CWString(text)
	defer free()
	textFormatAll(0, (*uint16)(unsafe.Pointer(cstr)))
}

func ClientPlayerNetCode() int {
	return int(nox_player_netCode_85319C)
}

func ClientSetPlayerNetCode(id int) {
	nox_player_netCode_85319C = uint32(id)
}

func Nox_xxx_playerForceDisconnect_4DE7C0(ind ntype.PlayerInd) {
	sessionPlayerDeparture(int32(ind))
}

func Get_nox_xxx_updatePlayerMonsterBot_4FAB20() unsafe.Pointer {
	return updateIdentityKey(updateIDPlayerMonsterBot)
}

func Nox_xxx_netNeedTimestampStatus_4174F0(pl *server.Player, v int) {
	playerStateAddStatus(pl, uint32(v))
}

func Sub_40A1F0(v int) {
	serverConfigTimerSet(int32(v))
}

func Nox_game_sendQuestStage_4D6960(v ntype.PlayerInd) {
	questRuntimeStageMessage(int(v), 14, 0)
}

func Nox_xxx_playerForceSendLessons_416E50(v int) {
	playerStateLessons(int32(v))
}

func Get_nox_xxx_updatePlayerObserver_4E62F0() unsafe.Pointer {
	return updateIdentityKey(updateIDPlayerObserver)
}

func Nox_xxx_playerRemoveSpawnedStuff_4E5AD0(u *server.Object) {
	stateRemoveSpawned(u)
}

func Nox_xxx_playerObserverFindGoodSlave0_4E6280(p *server.Player) *server.Object {
	return controlObserverSlave(p.C())
}

func Sub_4E6150(p *server.Player) *server.Object {
	return controlNextObserver(p.C())
}

func Get_nox_xxx_updatePlayer_4F8100() unsafe.Pointer {
	return updateIdentityKey(updateIDPlayer)
}

func Nox_xxx_playerUnsetStatus_417530(p *server.Player, a2 int) {
	playerStateRemoveStatus(p, uint32(a2))
}

func Nox_xxx_playerResetImportantCtr_4E4F40(v ntype.PlayerInd) {
	reliableResetRate(int(v))
}

func Get_dword_5d4594_1046492() int {
	return int(serverOptionsRoot)
}

func Nox_xxx_playerInitColors_461460(pl *server.Player) {
	clientPlayerColors(pl)
}

func Sub_425B30(a1 unsafe.Pointer, a2 ntype.PlayerInd) {
	playerGroupAddMember((*playerGroup)(a1), int32(a2))
}

func Sub_425A70(a1 int) unsafe.Pointer {
	return unsafe.Pointer(playerGroupFind(uint32(a1)))
}

func Sub_425AD0(a1 int, a2 *uint16) unsafe.Pointer {
	return unsafe.Pointer(playerGroupAdd(uint32(a1), a2))
}

func Sub_41D670(a1 string) {
	// The former service list has no population path.
}

func Sub_4DF3C0(p *server.Player) {
	matchRosterAssignTeam(p)
}

func Sub_40AA70(p *server.Player) int {
	return playerStateAdmission(p)
}

func Nox_xxx_netReportPlayerStatus_417630(p *server.Player) {
	playerStateReport(p)
}

func Sub_509C30(p *server.Player) {
	matchRosterRemember(p)
}

func Nox_xxx_playerLeaveObserver_0_4E6AA0(p *server.Player) {
	controlLeaveObserver(p.C())
}

func Nox_xxx_netGuiGameSettings_4DD9B0(a1 int, a2 *server.Settings2, a3 int) {
	matchRosterGUISettings(byte(a1), unsafe.Pointer(a2), a3)
}

func Sub_459AA0(a1 *server.Settings2) {
	serverOptionsRead(serverOptionsRecord(unsafe.Pointer(a1)))
}

func Nox_xxx_netNotifyRate_4D7F10(v ntype.PlayerInd) {
	gameplayReportRate(int(int32(v)))
}

func Nox_xxx_plrReadVals_4EEDC0(obj *server.Object, a2 int) {
	controlReadStats(obj, int32(a2))
}

func Nox_xxx_playerManaAdd_4EEB80(obj *server.Object, v int) {
	resourceAddMana(obj, int16(v))
}

func Nox_xxx_removePoison_4EE9D0(obj *server.Object) {
	resourceRemovePoison(obj)
}

func Sub_4FD0E0(obj *server.Object, sp spell.ID) int {
	return int(spellLifeCheckClass(obj, int32(sp)))
}

func Nox_xxx_checkPlrCantCastSpell_4FD150(obj *server.Object, sp spell.ID, a3 int) int {
	return int(spellLifeCantCast(obj, int32(sp), int32(a3)))
}

func Sub_4FCF90(obj *server.Object, sp spell.ID, a3 int) int {
	return int(spellLifeSpendMana(obj, int32(sp), int32(a3)))
}

func Sub_4D79A0(pli ntype.PlayerInd) {
	questRuntimeSlotMask(uint32(pli))
}

func Sub_4E80C0(pli ntype.PlayerInd) {
	matchRosterClearMask(byte(pli))
}

func Nox_xxx_player_4E3CE0() int {
	return int(questRuntimeCount())
}

func Sub_4E55F0(pli ntype.PlayerInd) {
	reliableRemoveRecipient(byte(pli))
}
