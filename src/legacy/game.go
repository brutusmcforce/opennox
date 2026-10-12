package legacy

import (
	"context"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_GetEndgameDialog            func() string
	GameGetPlayState                    func() int
	ServerCheatGod                      func(enable bool)
	Nox_xxx_serverHost_43B4D0           func()
	Nox_xxx_netServerCmd_440950         func(id byte, cmd string)
	ExecConsoleCmd                      func(ctx context.Context, cmd string) bool
	Nox_xxx_gameIsNotMultiplayer_4DB250 func() bool
	Nox_xxx_gameSetSwitchSolo_4DB220    func(a1 int)
	Nox_xxx_gameIsSwitchToSolo_4DB240   func() bool
	Nox_xxx_gameSetWallsDamage_4E25A0   func(v int)
	GetDoDamageWalls                    func() bool
	Sub_41CC00                          func(s string)
	Nox_xxx_playerSendMOTD_4DD140       func(a1 ntype.PlayerInd)
	Nox_client_getChatMap_49FF40        func() string
	Nox_xxx_mapSwitchLevel_4D12E0       func(a1 bool)
)

func init() {
	gui.RegisterState(client.StateOptions, "Options", func() bool {
		return optionsMenu.construct() != 0
	})
	gui.RegisterState(client.StateClassSelect, "ClassSelect", func() bool {
		return characterShowClass() != 0
	})
	gui.RegisterState(client.StateColorSelect, "ColorSelect", func() bool {
		return characterShowColor() != 0
	})
	gui.RegisterState(client.StateServerList, "ServerList", func() bool {
		return nox_game_showGameSel_4379F0() != 0
	})
	gui.RegisterState(client.StateXxx, "StateXxx", func() bool {
		return nox_game_showGameSel_4379F0() != 0
	})
}

func nox_server_parseCmdText_443C80(cstr *wchar2_t, _ int) int {
	cmd := GoWString(cstr)
	if cmd == "" {
		return 0
	}
	res := ExecConsoleCmd(context.Background(), cmd)
	return bool2int(res)
}

func nox_xxx_gameIsSwitchToSolo_4DB240() int {
	return bool2int(Nox_xxx_gameIsSwitchToSolo_4DB240())
}

func sub_517590(x float32, y float32) int {
	return bool2int(GetServer().S().Map.ValidIndexPos(types.Ptf(x, y)))
}

func nox_xxx_gameSetWallsDamage_4E25A0(v int) {
	Nox_xxx_gameSetWallsDamage_4E25A0(v)
}

func mapDamageUnitsAround(pos types.Pointf, outer, inner float32, damage, kind int32, who, exclude *server.Object) {
	GetServer().Nox_xxx_mapDamageUnitsAround(pos, outer, inner, int(damage), object.DamageType(kind), who, objectAsInterface(exclude), GetDoDamageWalls())
}

func nox_game_decStateInd_43BDC0() {
	GetClient().GamePopState()
}

func sub_4537F0() {
	GetServer().S().Sub4537F0()
}

func nox_xxx_mapCheck_537110(a1, a2 *nox_object_t) int {
	return bool2int(GetServer().S().MapTraceVision(asObjectS(a1), asObjectS(a2)))
}

func Nox_xxx_sMakeScorch_537AF0(pos types.Pointf, a2 int) { motionScorch(&pos, int32(a2)) }

func Nox_xxx_getSomeMapName_4D0CF0() string {
	return GoStringP(unsafe.Pointer(mapCycleNext()))
}

func Nox_server_gameSettingsUpdated_40A670() {
	legacyGlobals.nox_server_gameSettingsUpdated = 1
}

func Sub_43AF50() {
	Set_dword_5d4594_2650652(0)
}

func Nox_server_testTwoPointsAndDirection_4E6E50(p1 types.Pointf, dir int16, p2 types.Pointf) int {
	cp1, free1 := alloc.New(types.Pointf{})
	defer free1()
	cp2, free2 := alloc.New(types.Pointf{})
	defer free2()
	*cp1, *cp2 = p1, p2
	return int(stateFront(cp1, int32(dir), cp2))
}

func Nox_xxx_mapLoadOrSaveMB_4DCC70(v int) {
	sessionLoadStateSet(int32(v))
}
func Sub_44E560() unsafe.Pointer {
	return unsafe.Pointer(briefingCreateWindow())
}
func Nox_xxx_serverOptionsGetServername_40A4C0() string {
	return GoStringP(unsafe.Pointer((*int8)(unsafe.Pointer(serverConfigNameGet()))))
}
func Nox_xxx_mapGetMapName_409B40() string {
	return GoStringP(unsafe.Pointer(sessionMapName()))
}
func Nox_xxx_servGetPlrLimit_409FA0() int {
	return int(int32(serverConfigLimitGet()))
}
func Nox_client_xxx_switchChatMap_43B510() {
	nox_client_xxx_switchChatMap_43B510()
}
func Nox_client_guiXxx_43A9D0() {
	nox_client_guiXxx_43A9D0()
}
func Sub_43B630() {
	sub_43B630()
}
func Sub_49FF20() {
	sub_49FF20()
}
func Sub_445450() {
	interactionMessagesClear()
}
func Sub_45DB90() {
	quickbarResetFlash()
}
func Nox_xxx_initTime_435570() {
	nox_xxx_initTime_435570()
}
func Nox_xxx_allocArrayHealthChanges_49A5F0() int {
	return bool2int(combatHealthInit())
}
func Nox_xxx_loadGuides_427070() int {
	return int(bookLoadGuides())
}
func Sub_494F00() int {
	return drawableEffectTypes()
}
func Nox_xxx_loadReflSheild_499360() int {
	return bool2int(presentationShieldInit())
}
func Nox_xxx_allocClassListFriends_495980() int {
	return bool2int(combatFriendInit())
}
func Sub_4958F0() {
	combatFeedInit()
}
func Sub_460380() {
	quickbarClearAbilities()
}
func Nox_xxx_cliPrepareGameplay1_460E60() int {
	return int(quickbarPrepare())
}
func Nox_xxx_cliPrepareGameplay2_4721D0() {
	nox_xxx_cliPrepareGameplay2_4721D0()
}
func Sub_4951C0() {
	combatAllyClear()
}
func Nox_xxx_netGameSettings_4DEF00() {
	matchRosterSettings()
}
func Nox_server_gameUnsetMapLoad_40A690() {
	serverConfigUpdatedClear()
}
func Sub_416650() int {
	return int(int32(serverConfigRecordState()))
}
func Sub_46DCC0() {
	scoreboardCollect()
}
func Sub_409B80() string {
	return GoStringP(unsafe.Pointer(sessionSelectedMap()))
}
func Sub_4EDD70() {
	orchestrationDropFlags()
}
func Sub_4573B0() {
	teamUIRequestsReset()
}

func Sub_455C30() int {
	return teamUICTFConstruct()
}
func Sub_456070() int {
	return teamUIBallConstruct()
}
func Nox_xxx_guiHealthManaInit_4714E0() int {
	return uiMeterInit()
}
func Nox_xxx_bookInit_45B9D0() int {
	return bookInit()
}
func Sub_476E20() unsafe.Pointer {
	return presentationPhonemeInit().C()
}
func Sub_4BFAD0() int {
	return int(sub_4BFAD0())
}
func Nox_xxx_wndCreateInventoryMB_465E00() uint32 {
	return uint32(uiInventoryCreateWindow())
}
func Nox_game_initOptionsInGame_4ADAD0() int {
	return optionsInGame.construct()
}
func Sub_48D000_initGuiKick() int {
	return voteGUIInit()
}
func Sub_4C3760() int {
	return bindingInGame.construct()
}
func Sub_4C09D0() int {
	return uiTradeInit()
}
func Sub_478110() int {
	return uiShopInit()
}
func Sub_49B3E0() int {
	return int(sub_49B3E0())
}
func Sub_4BFC90() int {
	return int(sub_4BFC90())
}
func Nox_gui_itemAmount_init_4BFEF0() int {
	return uiAmountInit()
}
func Sub_4799A0() int {
	return int(sub_4799A0())
}
func Sub_46A730() unsafe.Pointer {
	return unsafe.Pointer(sub_46A730())
}
func Sub_4C3500() int {
	return bool2int(bindingYesNo() != nil)
}
func Nox_xxx_guiDrawRank_46E870() uint32 {
	return memmap.PtrWord(unsafe.Pointer(scoreboardConstruct()))
}
func Nox_xxx_guiMotdLoad_4465C0() uint32 {
	return memmap.PtrWord(sessionMOTDOpen().C())
}
func Nox_xxx_guiSummonCreatureLoad_4C1D80() int {
	return summonCreate()
}
func Sub_4AB260() int {
	return int(sessionDisconnectOpen())
}
func Nox_xxx_guiChatIconLoad_445650() int {
	return int(nox_xxx_guiChatIconLoad_445650())
}
func Sub_4C3390() int {
	return int(sub_4C3390())
}
func Sub_48C980() int {
	return int(sub_48C980())
}
func Sub_4D22B0() {
	orchestrationTransitionPlayers()
}
func Sub_459870() unsafe.Pointer {
	return serverOptionsListHead()
}
func Nox_xxx_gamePlayIsAnyPlayers_40A8A0() int {
	return playerStateMultiple()
}
func Sub_40A250() {
	serverConfigTimerInit()
}
func Sub_4D2160() {
	orchestrationRoundFlags()
}
func Nox_xxx_mapLoadRequired_4DCC80() int {
	return int(orchestrationLoadState())
}
func Sub_40A970() {
	playerStateReset()
}
func Nox_gui_itemAmount_free_4C03E0() {
	uiAmountFree()
}
func Sub_4AE3B0() {
	optionsDestroy()
}
func Sub_48D450() {
	voteGUIClose()
}
func Sub_4C4220() {
	bindingDestroy()
}
func Nox_xxx_closeP2PTradeWnd_4C12A0() {
	uiTradeDestroy()
}
func Sub_4BFD10() {
	sub_4BFD10()
}
func Sub_49B490() {
	sub_49B490()
}
func Sub_478F80() {
	uiShopDestroy()
}
func Sub_479D10() {
	sub_479D10()
}
func Sub_4AB470() {
	sessionDisconnectClose()
}
func Sub_4C34A0() {
	sub_4C34A0()
}
func Sub_445770() {
	sub_445770()
}
func Sub_456240() {
	teamUIHUDDestroy(true)
}
func Sub_455EE0() {
	teamUIHUDDestroy(false)
}
func Sub_4505E0() {
	briefingDestroy()
}
func Sub_46A860() {
	sub_46A860()
}
func Sub_467980() {
	uiInventoryResetWindow()
}
func Sub_460D50() {
	quickbarDestroy()
}
func Nox_xxx_guiServerOptionsGetGametypeName_4573C0(a1 noxflags.GameFlag) string {
	return serverOptionsModeName(uint16(a1))
}
func Sub_40A180(a1 noxflags.GameFlag) int {
	return int(uint8(serverConfigMinutes(int16(a1))))
}
func Nox_xxx_servGamedataGet_40A020(a1 uint16) int {
	return int(int16(serverConfigScore(int16(a1))))
}
func Sub_41D1A0(a1 int) {
	onlineSessionBriefing(uint32(a1))
}
func Nox_xxx_netPlayerIncomingServ_4DDF60(a1 int) {
	sessionPlayerIncoming(int32(a1))
}
func Nox_xxx_plrLoad_41A480(a1 string) int {
	return playerFileClientLoad(a1)
}
func Sub_465DE0(a1 int) {
	uiInventorySetWindowLevel(a1)
}
func Sub_4E79B0(a1 int) {
	*memmap.PtrUint32(0x5d4594, 1567712) = uint32(a1)
}
func Nox_xxx_playerMakeDefItems_4EF7D0(a1 *server.Object, a2 int, a3 int) {
	controlDefaultItems(a1, int32(a2), int32(a3))
}
func Sub_4181F0(a1 int) {
	teamRuntimeBalance(a1 != 0)
}
func Sub_4AB4A0(a1 int) {
	sessionDisconnectIconShow(a1)
}
func Sub_4AB4D0(a1 int) {
	sessionDisconnectShow(a1)
}
func Sub_4721A0(a1 int) {
	sub_4721A0(a1)
}
func Sub_460EA0(a1 int) {
	quickbarVisible(a1 != 0)
}
func Nox_window_set_visible_unk5(a1 int) {
	nox_window_set_visible_unk5(int(int32(a1)))
}
func Sub_45D500(a1 int) {
	bookTemporaryShow(a1)
}
func Sub_455A00(a1 int) {
	teamUIHUDShow(false, a1)
}
func Sub_455F10(a1 int) {
	teamUIHUDShow(true, a1)
}
func Nox_xxx_mapFindPlayerStart_4F7AB0(a2 *server.Object) types.Pointf {
	var out types.Pointf
	controlFindStart(&out, a2)
	return out
}
func Sub_500510(a1 string) {
	questProgressNamespace(a1)
}
func Nox_xxx_mapSwitchLevel_4D12E0_tileFree() {
	sessionClearTiles()
}
func Sub_57A1E0(a1 *server.Settings2, a2 string, a3 unsafe.Pointer, a4 int, a5 noxflags.GameFlag) {
	ruleLoad(a1, ruleCString(a2), (*legacyListNode)(a3), byte(a4), uint16(a5))
}
func Sub_57AAA0(a1 string, a2 *server.Settings2, a3 unsafe.Pointer) {
	ruleWrite(ruleCString(a1), a2, (*legacyListNode)(a3))
}
func Sub_4EF660(a1 *server.Object) {
	orchestrationResetPlayer(a1)
}
func Sub_4DBA30(a1 bool) {
	orchestrationRestore(int32(bool2int(a1)))
}
func Nox_xxx_isUnit_4E5B50(a1 *server.Object) int {
	return bool2int(stateIsUnit(a1))
}
func Sub_4E5B80(a1 *server.Object) int {
	return bool2int(stateIsPixie(a1))
}
func Sub_4E81D0(a1 *server.Object) {
	stateResetPixie(a1)
}
func Sub_4D71E0(a1 int) {
	questRuntimeSetSoulFrame(uint32(a1))
}
func Nox_xxx_calcDistance_4E6C00(a1 *server.Object, a2 *server.Object) float32 {
	return float32(stateDistance(a1, a2))
}
func Get_nox_game_switchStates_43C0A0() unsafe.Pointer {
	return animationCallbackKey(animationKeySwitchStates)
}
func Get_nox_game_showOptions_4AA6B0() unsafe.Pointer {
	return animationCallbackKey(animationKeyShowOptions)
}
func Get_nox_game_showMainMenu_4A1C00() unsafe.Pointer {
	return animationCallbackKey(animationKeyShowMainMenu)
}
func Sub_41CAC0(a1 string, data []byte) {
	playerFileExtract(a1, unsafe.Pointer(&data[0]))
}
func Nox_xxx_spell_4FE680(a1 *server.Object, a2 float32) {
	spellLifeCounterBooks(a1, float32(a2))
}
