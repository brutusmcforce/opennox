package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/server"
)

func Sub_460D40() bool {
	return dword_5d4594_1049508 != 0
}

func Sub_4D7140(a1 uint32) {
	dword_5d4594_1556144 = uint32(a1)
}

func Sub_43AF30() uint32 {
	return uint32(browserUI.hosting)
}

func Get_nox_wol_server_result_cnt_815088() int {
	return int(browserUI.resultCount)
}
func Get_dword_5d4594_815044() int {
	return int(browserUI.transition)
}
func Get_dword_587000_87412() int {
	return int(browserUI.region)
}
func Get_dword_5d4594_815000() *gui.Window {
	return browserUI.detailPanel
}
func Get_dword_587000_87408() int {
	return int(browserUI.listMode)
}
func Get_nox_game_createOrJoin_815048() int {
	return int(browserUI.creating)
}
func Get_dword_5d4594_814996() *gui.Window {
	return browserUI.label
}
func Get_dword_5d4594_814984() *gui.Window {
	return browserUI.mapWindow
}
func Get_dword_5d4594_2488604() int {
	return int(collisionActiveHead)
}
func Get_dword_5d4594_528252() int {
	return int(onlineRetryPending)
}
func Get_dword_5d4594_528260() int {
	return int(onlineRetryDeadline)
}
func Get_dword_587000_122852() unsafe.Pointer {
	return legacyGlobals.dword_587000_122852
}
func Get_dword_587000_93164() unsafe.Pointer {
	return legacyGlobals.dword_587000_93164
}
func Get_nox_client_translucentFrontWalls_805844() int {
	return int(nox_client_translucentFrontWalls_805844)
}
func Get_nox_client_highResFrontWalls_80820() int {
	return int(nox_client_highResFrontWalls_80820)
}
func Get_nox_client_highResFloors_154952() int {
	return int(nox_client_highResFloors_154952)
}
func Get_nox_client_fadeObjects_80836() int {
	return int(clientFadeObjects)
}
func Get_nox_client_renderBubbles_80844() int {
	return int(nox_client_renderBubbles_80844)
}
func Get_nox_server_sanctuaryHelp_54276() int {
	return int(nox_server_sanctuaryHelp_54276)
}
func Get_nox_server_sendMotd_108752() int {
	return int(nox_server_sendMotd_108752)
}
func Get_nox_server_connectionType_3596() int {
	return int(nox_server_connectionType_3596)
}
func Get_nox_server_resetQuestMinVotes_229988() int {
	return int(voteQuestResetSetting)
}
func Get_nox_server_kickQuestPlayerMinVotes_229992() int {
	return int(voteQuestKickSetting)
}
func Get_dword_5d4594_816340() int {
	return int(dword_5d4594_816340)
}
func Get_dword_587000_93156() int {
	return int(dword_587000_93156)
}
func Get_dword_587000_93160() int {
	return int(dword_587000_93160)
}
func Get_dword_5d4594_816348() int {
	return int(dword_5d4594_816348)
}
func Get_dword_5d4594_811904() int {
	return int(interactionDrawToggle)
}
func Get_dword_5d4594_1193156() int {
	return int(dword_5d4594_1193156)
}
func Get_dword_5d4594_1193188() int {
	return int(dword_5d4594_1193188)
}
func Get_dword_5d4594_3798840() int {
	return int(dword_5d4594_3798840)
}
func Get_dword_5d4594_3798824() int {
	return int(dword_5d4594_3798824)
}
func Get_dword_5d4594_3798836() int {
	return int(dword_5d4594_3798836)
}
func Get_dword_5d4594_3798820() int {
	return int(dword_5d4594_3798820)
}
func Get_nox_alloc_chat_1197364() unsafe.Pointer {
	return chatBubblePool
}
func Get_dword_5d4594_2660032() int {
	return int(dword_5d4594_2660032)
}
func Get_dword_5d4594_1090048() *gui.Window {
	return legacyGlobals.dword_5d4594_1090048
}
func Get_dword_5d4594_1090120() int {
	return int(dword_5d4594_1090120)
}
func Get_dword_5d4594_1090100() *gui.Window {
	return AsWindowP(memmap.WordPtr(legacyGlobals.dword_5d4594_1090100))
}
func Get_dword_5d4594_814624() unsafe.Pointer {
	return browserUI.selected
}
func Get_dword_5d4594_251560() unsafe.Pointer {
	return worldSecretHead
}
func Get_dword_5d4594_2489460() int {
	return int(dword_5d4594_2489460)
}
func Get_gameex_flags() uint32 {
	return uint32(gameex_flags)
}
func Get_gameex_flags_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&gameex_flags))
}
func Get_dword_5d4594_1064868() int {
	return int(interactionChatActive)
}
func Get_dword_5d4594_1316972() *gui.Window {
	return serverPanelsWindow(1316972)
}
func Get_dword_5d4594_1599588() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabPath))
}
func Get_dword_5d4594_1599592() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabAlternate))
}
func Get_dword_5d4594_1599540() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabObjects))
}
func Get_dword_5d4594_1599476() int {
	return int(*prefabGlobal(prefabPlaced))
}
func Get_dword_5d4594_1599548() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabWaypoints))
}
func Get_dword_5d4594_1599556() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabTiles))
}
func Get_dword_5d4594_1599532() unsafe.Pointer {
	return mapRoomPointer(*prefabGlobal(prefabWalls))
}
func Get_dword_5d4594_1200804() int {
	return int(clientGameMapFrame)
}
func Get_nox_gameDisableMapDraw_5d4594_2650672() int {
	return int(nox_gameDisableMapDraw_5d4594_2650672)
}
func Get_dword_5d4594_528256() int {
	return int(onlineRetryActive)
}
func Get_dword_5d4594_527988() int {
	return int(onlineSessionStatus)
}
func Get_nox_alloc_screenParticles_806044() unsafe.Pointer {
	return legacyGlobals.nox_alloc_screenParticles_806044
}
func Get_dword_5d4594_251744() int {
	return int(dword_5d4594_251744)
}
func Get_dword_5d4594_830236() *gui.Window {
	return AsWindowP(legacyGlobals.dword_5d4594_830236)
}
func Get_dword_5d4594_830232() *gui.Window {
	return AsWindowP(legacyGlobals.dword_5d4594_830232)
}
func Get_dword_5d4594_831260() int {
	return int(dword_5d4594_831260)
}
func Get_dword_5d4594_831220() int {
	return int(dword_5d4594_831220)
}
func Get_nox_color_black_2650656() uint32 {
	return uint32(nox_color_black_2650656)
}
func Get_nox_color_white_2523948() uint32 {
	return uint32(nox_color_white_2523948)
}
func Get_nox_color_violet_2598268() uint32 {
	return uint32(nox_color_violet_2598268)
}
func Get_nox_color_red_2589776() uint32 {
	return uint32(nox_color_red_2589776)
}
func Get_nox_color_green_2614268() uint32 {
	return uint32(nox_color_green_2614268)
}
func Get_nox_color_cyan_2649820() uint32 {
	return uint32(nox_color_cyan_2649820)
}
func Get_nox_color_blue_2650684() uint32 {
	return uint32(nox_color_blue_2650684)
}
func Get_nox_color_orange_2614256() uint32 {
	return uint32(nox_color_orange_2614256)
}
func Get_nox_color_yellow_2589772() uint32 {
	return uint32(nox_color_yellow_2589772)
}
func Get_nox_color_black_2650656_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_black_2650656))
}
func Get_nox_color_white_2523948_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_white_2523948))
}
func Get_nox_color_violet_2598268_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_violet_2598268))
}
func Get_nox_color_red_2589776_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_red_2589776))
}
func Get_nox_color_green_2614268_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_green_2614268))
}
func Get_nox_color_cyan_2649820_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_cyan_2649820))
}
func Get_nox_color_blue_2650684_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_blue_2650684))
}
func Get_nox_color_orange_2614256_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_orange_2614256))
}
func Get_nox_color_yellow_2589772_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_color_yellow_2589772))
}
func Get_dword_5d4594_823776() int {
	return int(dword_5d4594_823776)
}
func Get_dword_5d4594_831236() *gui.Window {
	return AsWindowP(memmap.WordPtr(legacyGlobals.dword_5d4594_831236))
}
func Get_dword_5d4594_1049844() int {
	return int(dword_5d4594_1049844)
}
func Get_dword_5d4594_3799452() int {
	return int(dword_5d4594_3799452)
}
func Get_dword_5d4594_1321236() *gui.Window {
	return legacyGlobals.dword_5d4594_1321236
}
func Get_dword_5d4594_1321240() *gui.Window {
	return legacyGlobals.dword_5d4594_1321240
}
func Get_dword_5d4594_1321244() *gui.Window {
	return legacyGlobals.dword_5d4594_1321244
}
func Get_dword_5d4594_1321248() *gui.Window {
	return legacyGlobals.dword_5d4594_1321248
}
func Get_dword_5d4594_1522616() *gui.Window {
	return legacyGlobals.dword_5d4594_1522616
}
func Get_dword_5d4594_1522620() *gui.Window {
	return legacyGlobals.dword_5d4594_1522620
}
func Get_dword_5d4594_1522624() *gui.Window {
	return legacyGlobals.dword_5d4594_1522624
}
func Get_dword_5d4594_1522628() *gui.Window {
	return legacyGlobals.dword_5d4594_1522628
}
func Get_dword_5d4594_1062452() *gui.Window {
	return legacyGlobals.dword_5d4594_1062452
}
func Get_dword_5d4594_2516356() int {
	return int(dword_5d4594_2516356)
}
func Get_dword_5d4594_1548524() int {
	return int(dword_5d4594_1548524)
}
func Get_dword_587000_155144() unsafe.Pointer {
	return unsafe.Pointer(audioStreamRoot())
}
func Get_nox_alloc_magicEnt_1569668() unsafe.Pointer {
	return legacyGlobals.nox_alloc_magicEnt_1569668
}
func Get_nox_client_gui_flag_1556112() int {
	return int(nox_client_gui_flag_1556112)
}

func Set_dword_5d4594_1599628(v int) {
	scriptBindingStringCount = uint32(v)
}

func Sub_41E2F0() int {
	return int(onlineSessionStatus)
}

func Get_nox_wol_wnd_world_814980() *gui.Window {
	return browserUI.world
}

func Get_dword_5d4594_815004() *gui.Window {
	return browserUI.detailList
}

func Nox_xxx_get_57AF20() int {
	return int(dword_5d4594_2523804)
}

func Set_dword_5d4594_1599644(v int) {
	*prefabGlobal(prefabScript) = uint32(v)
}
func Set_dword_5d4594_1599480(v uint32) {
	*prefabGlobal(prefabLoaded) = uint32(v)
}
func Set_dword_5d4594_1599476(v int) {
	*prefabGlobal(prefabPlaced) = uint32(v)
}
func Set_dword_5d4594_1599540(v unsafe.Pointer) {
	*prefabGlobal(prefabObjects) = mapRoomRaw(v)
}
func Set_dword_5d4594_1599532(v unsafe.Pointer) {
	*prefabGlobal(prefabWalls) = mapRoomRaw(v)
}
func Set_dword_5d4594_1599556(v unsafe.Pointer) {
	*prefabGlobal(prefabTiles) = mapRoomRaw(v)
}
func Set_dword_5d4594_1599548(v unsafe.Pointer) {
	*prefabGlobal(prefabWaypoints) = mapRoomRaw(v)
}
func Get_dword_5d4594_2386836() int {
	return int(scriptBindingMoverType)
}
func Set_dword_5d4594_2386836(v int) {
	scriptBindingMoverType = uint32(v)
}
func Set_dword_5d4594_2386848(v int) {
	scriptInventoryReserved = int32(v)
}
func Set_dword_5d4594_2386852(v int) {
	scriptInventoryNotice = uint32(v)
}
func Set_dword_5d4594_815104(v int) {
	browserUI.retry = int32(v)
}
func Set_nox_wol_server_result_cnt_815088(v int) {
	browserUI.resultCount = uint32(v)
}
func Set_qword_5d4594_815068(v uint64) {
	browserUI.refreshDeadline = uint64(v)
}
func Set_nox_client_connError_814552(v int) {
	browserUI.connectionError = uint32(v)
}
func Set_dword_5d4594_815056(v int) {
	browserUI.hasSelection = uint32(v)
}
func Set_nox_game_createOrJoin_815048(v int) {
	browserUI.creating = uint32(v)
}
func Set_dword_5d4594_527988(v int) {
	onlineSessionStatus = uint32(v)
}
func Set_dword_5d4594_10984(v int) {
	onlineConnectionFailed = uint32(v)
}
func Set_nox_alloc_screenParticles_806044(v unsafe.Pointer) {
	legacyGlobals.nox_alloc_screenParticles_806044 = v
}
func Set_nox_screenParticles_head(v unsafe.Pointer) {
	legacyGlobals.nox_screenParticles_head = (*Nox_screenParticle)(v)
}
func Set_dword_5d4594_806052(v unsafe.Pointer) {
	legacyGlobals.dword_5d4594_806052 = (*Nox_screenParticle)(v)
}
func Set_dword_5d4594_830236(v unsafe.Pointer) {
	legacyGlobals.dword_5d4594_830236 = v
}
func Set_nox_gameDisableMapDraw_5d4594_2650672(v int) {
	nox_gameDisableMapDraw_5d4594_2650672 = uint32(v)
}
func Set_dword_5d4594_831224(v int) {
	dword_5d4594_831224 = uint32(v)
}
func Set_dword_5d4594_1563096(v int) {
	orchestrationRestoreCleanup = uint32(v)
}
func Set_nox_win_width(v int) {
	nox_win_width = int32(v)
}
func Set_nox_win_height(v int) {
	nox_win_height = int32(v)
}
func Set_dword_5d4594_823776(v int) {
	dword_5d4594_823776 = uint32(v)
}
func Set_nox_draw_curDrawData_3799572(v *noxrender.RenderData) {
	legacyGlobals.nox_draw_curDrawData_3799572 = v
}
func Set_nox_color_black_2650656(v uint32) {
	nox_color_black_2650656 = uint32(v)
}
func Set_nox_color_white_2523948(v uint32) {
	nox_color_white_2523948 = uint32(v)
}
func Set_nox_color_violet_2598268(v uint32) {
	nox_color_violet_2598268 = uint32(v)
}
func Set_nox_color_red_2589776(v uint32) {
	nox_color_red_2589776 = uint32(v)
}
func Set_nox_color_green_2614268(v uint32) {
	nox_color_green_2614268 = uint32(v)
}
func Set_nox_color_cyan_2649820(v uint32) {
	nox_color_cyan_2649820 = uint32(v)
}
func Set_nox_color_blue_2650684(v uint32) {
	nox_color_blue_2650684 = uint32(v)
}
func Set_nox_color_orange_2614256(v uint32) {
	nox_color_orange_2614256 = uint32(v)
}
func Set_nox_color_yellow_2589772(v uint32) {
	nox_color_yellow_2589772 = uint32(v)
}
func Set_dword_5d4594_1305748(v int) {
	dword_5d4594_1305748 = uint32(v)
}
func Set_nox_pixbuffer_rows_3798784(v unsafe.Pointer) {
	legacyGlobals.nox_pixbuffer_rows_3798784 = (**byte)(v)
}
func Set_dword_5d4594_1193188(v int) {
	dword_5d4594_1193188 = uint32(v)
}
func Set_dword_5d4594_2649712(v uint32) {
	dword_5d4594_2649712 = uint32(v)
}
func Set_dword_5d4594_1548524(v int) {
	dword_5d4594_1548524 = uint32(v)
}
func Set_nox_client_gui_flag_1556112(v int) {
	nox_client_gui_flag_1556112 = uint32(v)
}
func Set_dword_5d4594_3798812(v int) {
	dword_5d4594_3798812 = uint32(v)
}
func Set_dword_5d4594_3798800(v int) {
	dword_5d4594_3798800 = uint32(v)
}
func Set_dword_5d4594_3798816(v int) {
	dword_5d4594_3798816 = uint32(v)
}
func Set_dword_5d4594_3798808(v int) {
	dword_5d4594_3798808 = uint32(v)
}
func Set_dword_5d4594_3798804(v int) {
	dword_5d4594_3798804 = uint32(v)
}
func Set_dword_5d4594_3798820(v int) {
	dword_5d4594_3798820 = uint32(v)
}
func Set_dword_5d4594_3798824(v int) {
	dword_5d4594_3798824 = uint32(v)
}
func Set_dword_5d4594_3798828(v int) {
	dword_5d4594_3798828 = uint32(v)
}
func Set_dword_5d4594_3798832(v int) {
	dword_5d4594_3798832 = uint32(v)
}
func Set_dword_5d4594_3798836(v int) {
	dword_5d4594_3798836 = uint32(v)
}
func Set_dword_5d4594_3798840(v int) {
	dword_5d4594_3798840 = uint32(v)
}
func Set_nox_video_tileBuf_ptr_3798796(v unsafe.Pointer) {
	legacyGlobals.nox_video_tileBuf_ptr_3798796 = v
}
func Set_nox_video_tileBuf_end_3798844(v unsafe.Pointer) {
	legacyGlobals.nox_video_tileBuf_end_3798844 = v
}
func Set_dword_5d4594_3799452(v int) {
	dword_5d4594_3799452 = uint32(v)
}
func Set_nox_alloc_chat_1197364(v unsafe.Pointer) {
	chatBubblePool = v
}
func Set_dword_587000_145664(v int) {
	dword_587000_145664 = uint32(v)
}
func Set_dword_587000_145668(v int) {
	dword_587000_145668 = uint32(v)
}
func Set_dword_5d4594_1090120(v int) {
	dword_5d4594_1090120 = uint32(v)
}
func Set_dword_5d4594_251744(v int) {
	dword_5d4594_251744 = uint32(v)
}
func Set_dword_5d4594_816340(v int) {
	dword_5d4594_816340 = uint32(v)
}
func Set_dword_5d4594_816348(v int) {
	dword_5d4594_816348 = uint32(v)
}
func Set_dword_5d4594_2489460(v int) {
	dword_5d4594_2489460 = uint32(v)
}
func Set_dword_5d4594_1548532(v unsafe.Pointer) {
	legacyGlobals.dword_5d4594_1548532 = v
}
func Set_dword_5d4594_831260(v int) {
	dword_5d4594_831260 = uint32(v)
}
func Set_nox_client_translucentFrontWalls_805844(v int) {
	nox_client_translucentFrontWalls_805844 = uint32(v)
}
func Set_nox_client_highResFrontWalls_80820(v int) {
	nox_client_highResFrontWalls_80820 = uint32(v)
}
func Set_nox_client_highResFloors_154952(v int) {
	nox_client_highResFloors_154952 = uint32(v)
}
func Set_nox_client_fadeObjects_80836(v int) {
	clientFadeObjects = uint32(v)
}
func Set_nox_client_renderBubbles_80844(v int) {
	nox_client_renderBubbles_80844 = uint32(v)
}
func Set_nox_server_sanctuaryHelp_54276(v int) {
	nox_server_sanctuaryHelp_54276 = uint32(v)
}
func Set_nox_server_sendMotd_108752(v int) {
	nox_server_sendMotd_108752 = uint32(v)
}
func Set_nox_server_connectionType_3596(v int) {
	nox_server_connectionType_3596 = uint32(v)
}
func Set_nox_server_resetQuestMinVotes_229988(v int) {
	voteQuestResetSetting = uint32(v)
}
func Set_nox_server_kickQuestPlayerMinVotes_229992(v int) {
	voteQuestKickSetting = uint32(v)
}
func Set_dword_5d4594_1308156(v unsafe.Pointer) {
	characterUI.modifiers[0] = v
}
func Set_dword_5d4594_1308160(v unsafe.Pointer) {
	characterUI.modifiers[1] = v
}
func Set_dword_5d4594_1308164(v unsafe.Pointer) {
	characterUI.modifiers[2] = v
}
func Set_dword_5d4594_2650676(v int) {
	dword_5d4594_2650676 = uint32(v)
}
func Set_dword_5d4594_2650680(v int) {
	dword_5d4594_2650680 = uint32(v)
}
func Set_qword_5d4594_814956(v uint64) {
	browserUI.connectionDeadline = uint64(v)
}
func Set_dword_5d4594_1200804(v uint32) {
	clientGameMapFrame = v
}
func Set_dword_5d4594_1200832(v int) {
	dword_5d4594_1200832 = uint32(v)
}
func Set_nox_alloc_magicEnt_1569668(v unsafe.Pointer) {
	legacyGlobals.nox_alloc_magicEnt_1569668 = v
}
func Set_dword_5d4594_1569672(v int) {
	dword_5d4594_1569672 = uint32(v)
}
func Inc_nox_xxx_waypointCounterMB_587000_154948() {
	nox_xxx_waypointCounterMB_587000_154948++
}
func Inc_dword_5d4594_2516356() {
	dword_5d4594_2516356++
}
func Inc_dword_5d4594_2386848() {
	scriptInventoryReserved++
}
func Inc_nox_wol_server_result_cnt_815088() {
	browserUI.resultCount++
}
func Set_dword_5d4594_1599588(v unsafe.Pointer) {
	*prefabGlobal(prefabPath) = mapRoomRaw(v)
}
func Set_dword_5d4594_1599592(v unsafe.Pointer) {
	*prefabGlobal(prefabAlternate) = mapRoomRaw(v)
}
func Get_dword_5d4594_1599592_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&dword_5d4594_2649712))
}
func Get_dword_5d4594_1123524() *gui.Window {
	return AsWindowP(interactionConversationRoot)
}
func Set_dword_5d4594_1123520(v int) {
	interactionConversationActive = uint32(v)
}
func Get_dword_5d4594_826036() *byte {
	return sessionMOTDFile
}
func Get_dword_8531A0_2576() *server.Player {
	return dword_8531A0_2576
}
func Set_dword_8531A0_2576(p *server.Player) {
	dword_8531A0_2576 = p
}
func Get_dword_8531A0_2572() uint32 {
	return uint32(dword_8531A0_2572)
}
func Get_dword_8531A0_2572_ptr() *uint32 {
	return (*uint32)(unsafe.Pointer(&dword_8531A0_2572))
}
func Set_dword_8531A0_2572(v uint32) {
	dword_8531A0_2572 = uint32(v)
}
func Get_dword_5d4594_2488620() uint32 {
	return uint32(dword_5d4594_2488620)
}
func Nox_xxx_playerGet_470A90() int {
	return int(dword_5d4594_1096252)
}

// nox_client_drawIsTick is set when the world is drawn as a part of a game tick.
// It's unset for additional frames drawn between the ticks (for high refresh rates),
// which must not advance any state (particle movement, timers, spawning effects, etc).
var nox_client_drawIsTick = true

func Set_nox_client_drawIsTick(v bool) {
	nox_client_drawIsTick = v
}

// nox_video_tileBuf_dirty marks the tile buffer cells (one per 46x46 floor tile) that were drawn into,
// so the GPU renderer only uploads those. It's owned by the caller of Set_nox_video_tileBuf_dirty.
var nox_video_tileBuf_dirty unsafe.Pointer

func Set_nox_video_tileBuf_dirty(v unsafe.Pointer) {
	nox_video_tileBuf_dirty = v
}
