package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	WinMainMenuAnimOutStartFnc   func() int
	Sub_44E320                   func()
	WinMainMenuAnimOutDoneFnc    func() int
	Sub_4A24C0                   func(a1 int) int
	Nox_game_showMainMenu_4A1C00 func() int
	Sub_43BE40                   func(a1 int)
	Sub_43BE30                   func() int
	Sub_4A18E0                   func(a1 *gui.Window, a2, a3, a4 int) int
)

// MainMenuEvent preserves the C adapter's event decoding and mutable-hook order.
func MainMenuEvent(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	code := e.EventCode()
	a, b := e.EventArgsC()
	ret := Sub_4A18E0(w, code, int(a), int(b))
	if ret == 0 {
		return nil
	}
	return gui.RawEventResp(ret)
}
func nox_client_drawGeneralCallback_4A2200() int {
	if err := GetClient().DrawGeneral(false); err != nil {
		guiLog.Println(err)
		return 0
	}
	return 1
}

func Get_nox_client_drawGeneralCallback_4A2200() unsafe.Pointer {
	return animationCallbackKey(animationKeyDrawGeneral)
}

func Get_nox_game_showGameSel_4379F0() unsafe.Pointer {
	return animationCallbackKey(animationKeyBrowserShow)
}

func Sub_461440(v int) {
	*quickbarWord(1049688) = uint32(v)
}

func Sub_4D7440(v int) {
	questRuntimeSetWord(1556116, uint32(v))
}

func Nox_xxx_cliSetMinimapZoom_472520(v int) {
	minimapSetZoom(v)
}
func Nox_xxx_monsterListFree_5174F0() {
	monsterDefinitionFree()
}
func Nox_client_countSaveFiles_4DC550() int {
	return int(sessionSaveSlots())
}
func Sub_4A7A60(a1 int) {
	characterUI.defaults = uint32(a1)
}
func Nox_client_lockScreenBriefing_450160(a1 int, a2 int, a3 int) {
	briefingShow(a1, a2, byte(a3))
}
func Sub_43E8C0(a1 int) {
	*memmap.PtrUint32(0x5D4594, 816408) = uint32(a1)
}
func Sub_4D6F80(a1 int) {
	questRuntimeSetWord(1556164, uint32(a1))
}
func Sub_4A7A70(a1 int) {
	*memmap.PtrUint32(0x5D4594, 1308168) = uint32(a1)
}

func Get_nox_game_showSelChar_4A4DB0() unsafe.Pointer {
	return animationCallbackKey(animationKeyShowSelChar)
}
func Get_nox_game_showSelClass_4A4840() unsafe.Pointer {
	return animationCallbackKey(animationKeyShowSelClass)
}
func Get_winMainMenuAnimOutStartFnc() unsafe.Pointer {
	return animationCallbackKey(animationKeyMainMenuStart)
}
func Get_winMainMenuAnimOutDoneFnc() unsafe.Pointer {
	return animationCallbackKey(animationKeyMainMenuDone)
}
