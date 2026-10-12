package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Sub_4A1BE0                             func(a1 int) int
	Sub_41E300                             func(a1 int) int
	Sub_4A50A0                             func() int
	Sub_4A50D0                             func() int
	Nox_client_resetScreenParticles_431510 func()
	Sub_46D6F0                             func() int
	Sub_413A00                             func(a1 int)
	Sub_44A400                             func()
	Nox_game_showSelChar_4A4DB0            func() int
	Nox_savegame_sub_46D580                func()
	Sub_450580                             func()
	Sub_4DB170                             func(a1 bool, a2 unsafe.Pointer, a3 int)
)

type Nox_screenParticle struct {
	Draw_fnc unsafe.Pointer      // 0, 0, func(unsafe.Pointer, *Nox_screenParticle) int
	Field_4  uint32              // 1, 4
	Field_8  uint32              // 2, 8
	Field_12 uint32              // 3, 12, color
	Field_16 uint32              // 4, 16, vx
	Field_20 uint32              // 5, 20, vy
	Field_24 uint32              // 6, 24, x
	Field_28 uint32              // 7, 28, y
	Field_32 uint32              // 8, 32
	Field_36 uint32              // 9, 36
	Field_40 [4]byte             // 10, 40
	Field_44 *Nox_screenParticle // 11, 44, next
	Field_48 *Nox_screenParticle // 12, 48, prev
}

func sub_413A00(a1 int) { Sub_413A00(a1) }

func Sub_41E470() {
	onlineSessionRetry()
}
func Nox_xxx_reconStart_41E400() {
	onlineSessionStart()
}
func Sub_467C10() int {
	return uiInventoryCloseWindow()
}
func Nox_xxx_guiServerOptionsTryHide_4574D0() int {
	return serverOptionsTryClose()
}
func Sub_446780() {
	sessionMOTDClose()
}
func Sub_460660() {
	quickbarCancelCapture()
}
func Sub_467BB0() {
	uiInventoryOpenWindow()
}
func Sub_4602F0() {
	quickbarClearSlots()
}
func Nox_xxx_playerAnimCheck_4372B0() int {
	return int(nox_xxx_playerAnimCheck_4372B0())
}
func Sub_41DA10(a1 int) {
	// The legacy queue allocator is gone; every queue remains empty.
}
func Nox_xxx_bookHideMB_45ACA0(a1 int) int {
	return bookHide(a1)
}
func Sub_4AD9B0(a1 int) {
	optionsClose(a1)
}
func Nox_xxx_bookShowMB_45AD70(a1 int) {
	bookShow(a1)
}
func Sub_41A000(a1 string, a2 *server.SaveGameInfo) int {
	return int(sessionSaveMetadata(a1, a2))
}
func Get_sub_4A50A0() unsafe.Pointer {
	return animationCallbackKey(animationKeySelCharStart)
}
func Get_sub_4A50D0() unsafe.Pointer {
	return animationCallbackKey(animationKeySelCharDone)
}
