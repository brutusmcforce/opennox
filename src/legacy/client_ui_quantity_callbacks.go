package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

const (
	uiAmountDialogExport = iota
	uiAmountDrop
	uiAmountShopBuy
	uiAmountShopSell
	uiAmountShopSellCancel
	uiAmountShopRepair
	uiAmountShopRepairCancel
	uiAmountShopActive
	uiAmountShopTooltip
	uiAmountTradeTooltip
)

var uiAmountIdentitySlots [10]byte

func uiAmountNativeKey(id int) unsafe.Pointer { return unsafe.Pointer(&uiAmountIdentitySlots[id]) }

func init() {
	gui.RegisterTooltipCallbackGo(uiAmountNativeKey(uiAmountShopTooltip), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiShopHover(uiInventoryPackedPoint(arg))
	})
	gui.RegisterTooltipCallbackGo(uiAmountNativeKey(uiAmountTradeTooltip), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiTradeHover(uiInventoryPackedPoint(arg))
	})
}

// Native callbacks retain the original five-word frame and stable identities.
var uiAmountCallbacks = make(map[unsafe.Pointer]func(uintptr, uintptr, uintptr, uintptr, uintptr))

func init() {
	uiAmountCallbacks[uiAmountNativeKey(uiAmountDrop)] = func(a1, a2, a3, a4, a5 uintptr) {
		p := (*[2]int32)(unsafe.Pointer(a1))
		uiInventoryDropQuantity(image.Pt(int(p[0]), int(p[1])), uint32(a2), uint32(a3), int(int32(a4)))
	}
	uiAmountCallbacks[uiAmountNativeKey(uiAmountShopBuy)] = func(a1, a2, a3, a4, a5 uintptr) {
		uiShopBuyAccept(uint16(a2), uint32(a3), uint32(a4))
	}
	uiAmountCallbacks[uiAmountNativeKey(uiAmountShopSell)] = func(a1, a2, a3, a4, a5 uintptr) {
		uiShopSellAccept(uint16(a2), uint16(a3), uint32(a4))
	}
	uiAmountCallbacks[uiAmountNativeKey(uiAmountShopSellCancel)] = func(a1, a2, a3, a4, a5 uintptr) {
		*uiShopWord(1098616) = 0
	}
	uiAmountCallbacks[uiAmountNativeKey(uiAmountShopRepair)] = func(a1, a2, a3, a4, a5 uintptr) {
		uiShopRepairAccept(uint16(a2))
	}
	uiAmountCallbacks[uiAmountNativeKey(uiAmountShopRepairCancel)] = func(a1, a2, a3, a4, a5 uintptr) {
		*uiShopWord(1098620) = 0
	}
}

func uiAmountCall(fn unsafe.Pointer, a1, a2, a3, a4, a5 uintptr) {
	if cb := uiAmountCallbacks[fn]; cb != nil {
		cb(a1, a2, a3, a4, a5)
		return
	}
	panic("unregistered quantity callback")
}
