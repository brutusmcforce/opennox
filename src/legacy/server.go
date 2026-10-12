package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/libs/types"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_mapToGameFlags func(v int) noxflags.GameFlag
	Sub_40A1A0         func() int
)

type Server interface {
	S() *server.Server
	SetUpdateFunc2(fnc func() bool)
	ServerPort() int
	ServStartCountdown(dt int, text strman.ID)
	SwitchMap(name string)
	GetFlag3592() bool
	Sub40A040settings(a1 int, a2 int)
	CreateObjectAt(obj, owner server.Obj, pt types.Pointf)
	NetUpdateRemotePlrAudioEvents(obj *server.Object, v2 unsafe.Pointer, v18 int8)
	Nox_xxx_mapDamageUnitsAround(pos types.Pointf, r1, r2 float32, dmg int, dtyp object.DamageType, who *server.Object, a7 server.Obj, damageWalls bool)
	Nox_xxx_mapReset5028E0()
	Nox_xxx_free503F40()
	ObjectDeleteLast(obj *server.Object)
	ObjectsAddPending()
	FinalizeDeletingObjects()
	TeamsResetYyy() int
	TeamsRemoveActive(hooks bool) int
	TeamRemove(t *server.Team, netUpd bool)
	DelayedDelete(obj *server.Object)
	Sub504720(a1, a2 uint32) int32
	ApplyForce(obj *server.Object, vec types.Pointf, force float64)
	PlayerSpell(u *server.Object)
	Nox_script_event_playerLeave(pl *server.Player)
	NoxScriptC() NoxScript
	Nox_xxx_spellAccept4FD400(spellID spell.ID, a2, obj3, obj4 *server.Object, sa *server.SpellAcceptArg, lvl int) bool
	Nox_xxx_generateRetreatPath_50CA00(path []types.Pointf, u *server.Object, a4 *types.Pointf) int
	Nox_xxx_creatureSetDetailedPath_50D220(obj *server.Object, a2 *types.Pointf)
	Sub_50CB20(a1 *server.Object, a2 *types.Pointf) *server.Waypoint
	Sub_50B810(obj *server.Object, pos *types.Pointf) bool
	Nox_xxx_mapDamageToWalls_534FC0(rect image.Rectangle, pos types.Pointf, rad float32, dmg int, dtyp object.DamageType, who *server.Object) bool
	Nox_xxx_damageToMap_534BC0(gx, gy int, dmg int, dtyp object.DamageType, who *server.Object) int
	Nox_xxx_wall_4DF1E0(a1 int)
}

var (
	GetServer func() Server
)

func sub_40A040_settings(a1 int16, a2 uint8) {
	GetServer().Sub40A040settings(int(a1), int(a2))
}

func sub_40A300() int {
	return bool2int(GetServer().GetFlag3592())
}

func nox_xxx_mapLoad_4D2450(a1 *int8) {
	GetServer().SwitchMap(GoStringP(unsafe.Pointer(a1)))
}

func gameFPS() uint32 {
	return GetServer().S().TickRate()
}

func gameFrame() uint32 {
	return GetServer().S().Frame()
}
func Sub_409A70(a1 int) int {
	return int(sessionModeIndex(int16(a1)))
}
func Nox_xxx_netInformTextMsg2_4DA180(a1 int, a2 unsafe.Pointer) {
	gameplayTextInformationAll(int(int32(a1)), a2)
}
func Nox_xxx_netReportUnitHeight_4D9020(a1 ntype.PlayerInd, a2 *server.Object) {
	gameplayReportHeight(int(int32(a1)), a2)
}
func Sub_511100(a1 int) {
	motionSentryReport(int32(a1))
}
func Nox_xxx_netUpdateRemotePlr_501CA0(a1 *server.Object) {
	objectReportRemoteAudio(a1)
}
func Nox_xxx_netSendObjects2Plr_519410(a1 *server.Object, a2 *server.Object) {
	objectReportRecipient(a1, a2)
}
func Sub_4D6770(a1 ntype.PlayerInd) {
	questRuntimeScoreboard(int(a1))
}
func Sub_4D6880(a1 int, a2 int) {
	questRuntimeStageMessage(a1, 13, uint32(a2))
}
func Sub_4D60B0() {
	questRuntimeResetAll()
}
func Sub_4CFDF0(a1 int) {
	sessionMapStateSet(int32(a1))
}
func Nox_xxx_playerMapTracksObj_4173D0(a1 int, a2 *server.Object) int {
	return playerStateTracks(a1, a2)
}
func Sub_519710(a1 unsafe.Pointer) int {
	return objectReportSchedule((*server.PlayerUpdateData)(a1))
}
func Nox_xxx_updateUnits_51B100_D() {
	orchestrationWalls()
}
func Nox_xxx_decay_511750() {
	motionDecayTick()
}
func Nox_server_checkVictory_509A60() {
	matchRosterCheckVictory()
}
func Nox_xxx_allocHitArray_5486D0() {
	collisionResetHits()
}
func Nox_xxx_updateObjectsVelocity_5118A0(a1 float32) {
	motionVelocity(a1)
}
func Sub_548B60() {
	collisionDrainAngles()
}
func Sub_537700() *server.Object {
	return collisionPopActive()
}
func Nox_xxx_collide_548740() {
	collisionDispatch()
}
func Nox_xxx_updatePoison_4EE8F0(a1 *server.Object, a2 int) {
	resourceReducePoison(a1, int32(a2))
}
func Nox_xxx_updateUnitBuffs_4FF620(a1 *server.Object) {
	spellLifeUpdateBuffs(a1)
}
func Nox_xxx_reconAttempt_41E390() {
	onlineSessionAttempt()
}
func Sub_5096F0() {
	matchRosterCheckLimit()
}
func Sub_416720() {
	serverConfigExpire()
}
func Sub_40AA00() int {
	return playerStateElapsed()
}
func Sub_40AA40() int {
	return int(playerStateThreshold())
}
func Nox_xxx_countNonEliminatedPlayersInTeam_40A830(a1 *server.Team) int {
	return playerStateTeamPlayers(a1)
}
func Sub_40A770() int {
	return playerStateCompetitors()
}
func Sub_40A6B0() int {
	return int(int32(serverConfigRateDirtyGet()))
}
func Sub_40A6A0(a1 int) {
	serverConfigRateDirtySet(int32(a1))
}
func Nox_xxx_netReportAllLatency_4D3050() {
	sessionReportLatency()
}
func Sub_4183C0() {
	teamRuntimeNearest()
}
func Sub_4D7150() {
	questRuntimeObserverDeadline()
}
func Sub_4D71F0() {
	questRuntimeSoulTimeout()
}
func Nox_server_checkWarpGate_4D7600() {
	questRuntimeWarpTick()
}
func Sub_4D7A80() {
	questRuntimeDepartureTick()
}
func Sub_4EC720() {
	itemRespawnTick()
}
func Sub_50D890() {
	spawnPolicyTick()
}
func Sub_4E4170() {
	orchestrationDifficulty()
}
func Nox_xxx_voteUptate_506F30() {
	voteTick()
}

func Nox_xxx_protectData_56F5C0() {
	nox_xxx_protectData_56F5C0()
}
func Sub_57B140() bool {
	return orchestrationTimeout()
}
func Sub_57B0A0() {
	runtimePauseStop()
}
func Sub_4DF120(a1 unsafe.Pointer) {
	matchRosterWall((*server.Wall)(a1), 59)
}
func Sub_4DF180(a1 unsafe.Pointer) {
	matchRosterWall((*server.Wall)(a1), 60)
}
func Nox_xxx_netUpdateObjectSpecial_527E50(a1 *server.Object, a2 *server.Object) {
	visibilitySpecialUpdate(a1, a2)
}
func Sub_4D15C0() {
	sessionScavengerReset()
}
func Sub_4D7B40() {
	questRuntimeDepartureReset()
}
func Sub_41E4B0(a1 int) {
	onlineSessionMarker(uint32(a1))
}
func Sub_4E4ED0() {
	reliableResetSequences()
}
func Nox_motd_4463E0(a1 int) {
	sessionMOTDRead(a1)
}
func Sub_4259C0() {
	playerGroupsInit()
}
func Sub_518770() int {
	return objectReportInit()
}
func Nox_xxx_allocItemRespawnArray_4ECA60() int {
	return itemRespawnInit()
}
func Nox_xxx_registerShopClasses_50E2A0() int {
	return shopInit()
}
func Nox_xxx_allocMonsterRelatedArrays_50D780() int {
	return spawnPolicyInit()
}
func Nox_xxx_allocVoteArray_5066D0() int {
	return voteInit()
}
func Nox_xxx_monsterList_517520() int {
	return monsterDefinitionBind()
}
func Sub_416920() {
	serverConfigAdmissionInit()
}
func Nox_xxx_allocPendingOwnsArray_516EE0() int {
	return monsterPendingInit()
}
func Sub_421B10() {
	mapPolygonReset()
}
func Sub_516F10() {
	monsterPendingFree()
}
func Sub_4259F0() {
	playerGroupsFree()
}
func Sub_4ECA90() {
	itemRespawnFree()
}
func Sub_506720() {
	voteClose()
}
func Sub_50D820() {
	spawnPolicyFree()
}
func Nox_xxx_deleteShopInventories_50E300() {
	shopFree()
}
func Sub_416950() {
	serverConfigAdmissionClose()
}
func Sub_56F3B0() {
	freeProtectionRecords()
}
func Sub_51D0E0() {
	sub_51D0E0()
}
func Sub_502DF0() {
	sub_502DF0()
}
func Sub_481410() {
	tileCompositionReset()
}
func Sub_4D0550(a1 string) {
	commandRulesPath(&a1)
}
func Sub_4165F0(a1 int, a2 int) {
	serverConfigSlotCopy(int32(a1), int32(a2))
}
func Sub_4D7430() int {
	return int(questRuntimeWord(1556116))
}
func Sub_4D76F0() int {
	return int(questRuntimeWord(1556124))
}
func Sub_4D76E0(a1 int) {
	questRuntimeSetWord(1556124, uint32(a1))
}
func Nox_xxx_ruleSetNoRespawn_40A5E0(a1 int) {
	serverConfigRespawnSet(int32(a1))
}
func Sub_455C10() {
	teamUIHUDHide(false)
}
func Sub_456050() {
	teamUIHUDHide(true)
}
func Nox_xxx_mapFindCrown_4CFC30() {
	sessionClearCrowns()
}
func Nox_xxx_mapInfoSetCapflag_417EA0() int {
	return teamUIMapCTF()
}
func Nox_xxx_mapInfoSetFlagball_417F30() int {
	return int(teamUIMapBall())
}
func Nox_xxx_mapInfoSetKotr_4180D0() int {
	return teamRuntimeCrown()
}
func Nox_xxx_teamAssignFlags_418640() {
	teamRuntimeAssignFlags()
}
func Nox_xxx_toggleAllTeamFlags_418690(a1 int) {
	teamRuntimeToggle(a1 != 0)
}
func Sub_410730() {
	worldSecretClear()
}
func Nox_xxx_j_allocHitArray_511840() {
	collisionResetHits()
}
func Nox_xxx_decayDestroy_5117B0() {
	motionDecayClear()
}
func Sub_510E50() {
	motionSentryHead = 0
}
func Sub_4D1610() {
	sessionScavengerMaximumReset()
}
func Sub_4EC5B0() {
	itemRespawnReset()
}
func Sub_50E360() {
	shopReset()
}
func Sub_50D7E0() {
	spawnPolicyReset()
}
func Sub_4E4F80() {
	reliableTrim()
}
func Sub_4ECFE0() {
	netCodeCacheFlush()
}
func Sub_511E20() {
	monsterCacheClear()
}
func Sub_516F30() {
	monsterPendingClear()
}
func Nox_xxx_playerBotCreate_4FA700(u *server.Object) {
	controlBotCreate(u)
}
func Sub_4F7950(u *server.Object) {
	controlClearWaypoints(u)
}
func Sub_4E5BF0(a1 bool) {
	stateCleanup(int32(bool2int(a1)))
}
func Nox_xxx_Fn_4FCAC0(a1 bool, a2 int) {
	spellLifeReset(int32(bool2int(a1)), int32(a2))
}
func Nox_xxx_shopCancelSession_510DC0(a1 *server.TradeSession) {
	shopCancel(a1)
}
func Sub_4FCEB0(a1 bool) {
	spellLifeCancelDurations(int32(bool2int(a1)))
}

var _ = [1]struct{}{}[60-unsafe.Sizeof(server.TileDef{})]

func Get_nox_tile_defs_arr() []server.TileDef {
	ptr := (*server.TileDef)(unsafe.Pointer(&worldTileDefinitions[0]))
	sz := int(worldTileDefinitionCount)
	return unsafe.Slice(ptr, len(worldTileDefinitions))[:sz]
}
