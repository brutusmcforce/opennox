package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/client/keybind"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
)

type CtrlEventHandler interface {
	Reset()
	Sub_42E8E0_go(key keybind.Event, a2 int) string
}

var (
	Nox_client_parseConfigHotkeysLine_42CF50 func(s string) int
)

func sub_42E8E0(key, a2 int) *wchar2_t {
	s := GetClient().GetCtrlEvent().Sub_42E8E0_go(keybind.Event(key), a2)
	return internWStr(s)
}

func Nox_xxx_guiSpellTargetClickCheckSend_45DBB0() {
	quickbarSendPending()
}

func Nox_xxx_clientTrade_42E850(dr *client.Drawable) {
	nox_xxx_clientTrade_42E850((*nox_drawable)(dr.C()))
}

func Nox_xxx_clientTalk_42E7B0(dr *client.Drawable) {
	nox_xxx_clientTalk_42E7B0((*nox_drawable)(dr.C()))
}

func Nox_xxx_clientCollideOrUse_42E810(dr *client.Drawable) {
	nox_xxx_clientCollideOrUse_42E810((*nox_drawable)(dr.C()))
}

func Nox_xxx_packetGetMarshall_476F40() uint32 {
	return uint32(nox_xxx_packetGetMarshall_476F40())
}

func Nox_xxx_clientUpdateButtonRow_45E110(ind int) {
	quickbarSelectRow(ind)
}

func Nox_client_trapSetSelect_4604B0(ind int) {
	quickbarTrapSelect(ind)
}

func Nox_client_orderCreature(a1, a2 int) {
	summonOrder((*summonRecord)(unsafe.Pointer(uintptr(a1))), uint32(a2))
}

func Get_nox_gui_itemAmount_dialog_1319228() *gui.Window {
	return AsWindowP(legacyGlobals.nox_gui_itemAmount_dialog_1319228)
}

func Sub_4C01C0(a1, a2 int, a3 unsafe.Pointer, a4 int) {
	uiAmountPanel(uiInventoryWindowValue(uint32(a1)), a2, uintptr(a3), uintptr(a4))
}

func Nox_client_chatStart_46A430(v int) {
	nox_client_chatStart_46A430(int32(v))
}

func Nox_client_toggleSpellbook_45AC70() {
	bookToggle()
}

func Nox_client_toggleInventory_467C60() {
	uiInventoryToggleWindow()
}

func Nox_client_mapZoomIn_4724E0() {
	minimapZoomIn()
}

func Nox_client_mapZoomOut_472500() {
	minimapZoomOut()
}

func Nox_client_invAlterWeapon_4672C0() {
	uiInventoryAlterWeapon()
}

func Nox_client_quickHealthPotion_472220() {
	uiMeterQuickPotion(0)
}

func Nox_client_quickManaPotion_472240() {
	uiMeterQuickPotion(1)
}

func Nox_client_quickCurePoisonPotion_472260() {
	uiMeterQuickPotion(2)
}

func Nox_client_spellSetNext_4604F0() {
	quickbarMoveRow(1)
}

func Nox_client_spellSetPrev_460540() {
	quickbarMoveRow(-1)
}

func Nox_client_spellSetSelect_460590() {
	quickbarTimedRow()
}

func Nox_client_buildTrap_45E040() {
	quickbarBuildTrap()
}

func Nox_xxx_guiServerOptsLoad_457500() {
	serverOptionsConstruct()
}

func Sub_460630() {
	quickbarRememberMouseSequence()
}
func Nox_xxx_game_4DCCB0() int {
	return int(sessionSaveAllowed())
}
func Sub_554290() uint32 {
	return pingMinimum()
}
func Sub_554300() uint32 {
	return pingAverage()
}
func Nox_client_mousePriKey_430AF0() int {
	return int(nox_client_mousePriKey_430AF0())
}
func Sub_42EB90(a1 int) {
	*memmap.PtrUint32(0x5D4594, 754052) = uint32(a1)
}
func Nox_client_invokeSpellSlot_45DA50(a1 int) {
	quickbarInvokeSlot(a1)
}
