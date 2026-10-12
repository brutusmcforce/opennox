package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

const (
	inventoryCallbackButton   = 2
	inventoryCallbackAlt      = 3
	inventoryCallbackStatus   = 4
	inventoryCallbackHover    = 19
	inventoryCallbackDrawAlt  = 34
	inventoryCallbackDrawCur  = 36
	inventoryCallbackMode     = 37
	inventoryCallbackIdentify = 38
)

var inventoryCallbackSlots [39]byte

func inventoryCallbackKey(index int) unsafe.Pointer {
	return unsafe.Pointer(&inventoryCallbackSlots[index])
}

func init() {
	gui.RegisterTooltipCallbackGo(inventoryCallbackKey(inventoryCallbackButton), func(_ *gui.Window, _ *gui.WindowData, _ uintptr) {
		uiInventoryButtonTooltip()
	})
	gui.RegisterTooltipCallbackGo(inventoryCallbackKey(inventoryCallbackAlt), func(_ *gui.Window, _ *gui.WindowData, _ uintptr) {
		uiInventoryAlternateTooltip()
	})
	gui.RegisterTooltipCallbackGo(inventoryCallbackKey(inventoryCallbackStatus), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiInventoryStatusTooltip(uiInventoryPackedPoint(uintptr(uint32(arg))))
	})
	gui.RegisterTooltipCallbackGo(inventoryCallbackKey(inventoryCallbackHover), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiInventoryHover(uiInventoryPackedPoint(uintptr(uint32(arg))))
	})
	gui.RegisterTooltipCallbackGo(inventoryCallbackKey(inventoryCallbackMode), func(w *gui.Window, _ *gui.WindowData, _ uintptr) {
		sub_466E20((*uint32)(w.C()))
	})
}
