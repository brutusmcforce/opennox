package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/client/keybind"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	InputSetKeyTimeoutLegacy                 func(key byte)
	InputKeyCheckTimeoutLegacy               func(key byte, dt uint32) bool
	Sub_416120                               func(key byte) bool
	Sub_416170                               func(key int) int
	Sub_416150                               func(key, ts int) int
	Nox_xxx_keybind_nameByTitle_42E960       func(title string) keybind.Key
	Nox_xxx_bindevent_bindNameByTitle_42EA40 func(title string) *keybind.BindEvent
	Sub_4C3B70                               func()
	Sub_4CBBF0                               func()
)

func nox_xxx_setKeybTimeout_4160D0(key int) int {
	InputSetKeyTimeoutLegacy(byte(key))
	return key
}

func nox_xxx_bookGet_430B40_get_mouse_prev_seq() int {
	return int(GetClient().GetInputSeq())
}

func nox_input_pollEvents_4453A0() int {
	// TODO
	//inpHandler.Tick()
	return 0
}

func NoxInputOnChar(c uint16) {
	uiEntryChar(c)
}
func Nox_xxx_clientIsObserver_4372E0() int {
	return int(nox_xxx_clientIsObserver_4372E0())
}
func Sub_4675B0() int {
	return uiInventoryMode()
}
func Sub_479590() int {
	return int(uiShopMode())
}
func Sub_478030() int {
	return int(uiShopActive())
}
func Sub_47A260() int {
	return int(sub_47A260())
}
func Nox_xxx_cursor_430B00() int {
	return int(nox_xxx_cursor_430B00())
}
func Sub_45D9B0() int {
	return int(*bookWord(1047520))
}
func Sub_45D870() {
	bookFinishAddition()
}
func Nox_xxx_sprite_4C3220(a1 *client.Drawable) int {
	return bool2int(summonFind(a1.NetCode32) != nil)
}
func Nox_xxx_wnd_46C2A0(a1 *gui.Window) int {
	return uiWindowHidden(a1)
}
func Nox_xxx_clientAskInfoMb_4BF050(a1 *client.Drawable) string {
	return GoWString((*wchar2_t)(unsafe.Pointer(uiItemTooltip(a1))))
}
func Nox_xxx_cursorSetTooltip_4776B0(a1 string) {
	if a1 == "" {
		uiCursorTooltip(nil)
		return
	}
	wstr, free := CWString(a1)
	defer free()
	uiCursorTooltip((*uint16)(unsafe.Pointer(wstr)))
}

func Sub_57B450(a1 *client.Drawable) int {
	return glyphItemAllowed(a1)
}
func Nox_xxx_clientPickup_46C140(a1 *client.Drawable) {
	nox_xxx_clientPickup_46C140((*nox_drawable)(a1.C()))
}
func Sub_46B630(a1 *gui.Window, a2 int, a3 int) *gui.Window {
	return uiWindowChildAt(a1, a2, a3)
}
