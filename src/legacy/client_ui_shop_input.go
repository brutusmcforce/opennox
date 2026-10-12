package legacy

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func uiShopString(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "GUIShop.c")
}
func uiShopPanel(w *gui.Window, event int, a, b uintptr) int {
	if event == 16393 {
		*uiShopWord(1107036) = *uiShopWord(1098592) - uint32(b)
		return 0
	}
	if event != 16391 || *bookWord(1047520) != 0 {
		return 0
	}
	id := uiInventoryWindowValue(uint32(a)).ID()
	uiTradeSound(766)
	switch id {
	case 3801:
		if uiShopActive() != 0 {
			uiShopCancelRequest()
		}
	case 3802:
		if uiShopActive() != 0 {
			if uiShopMode() == 4 {
				sub_467680()
			}
			nox_client_setCursorType_477610(12)
			*uiShopWord(1098628) = 3
			if !uiInventoryWindowOpenState() {
				uiInventoryOpenWindow()
			}
		}
	case 3803:
		if uiShopActive() != 0 {
			uiInventoryRepairMode()
			*uiShopWord(1098628) = 4
		}
	case 3804:
		if uiShopActive() != 0 {
			if uiShopMode() == 4 {
				sub_467680()
			}
			nox_client_setCursorType_477610(11)
			*uiShopWord(1098628) = 2
		}
	case 3808, 3809:
		scroll, max := *uiShopWord(1107036), *uiShopWord(1098592)
		if id == 3808 {
			if int32(scroll)-50 >= 0 {
				scroll = scroll - 50 - (scroll-50)%50
			} else {
				scroll = 0
			}
		} else {
			if int32(scroll)+50 <= int32(max) {
				scroll = scroll + 50 - (scroll+50)%50
			} else {
				scroll = max
			}
		}
		*uiShopWord(1107036) = scroll
		uiInventoryWindowValue(*uiShopWord(1098580)).Func94(gui.AsWindowEvent(16394, uintptr(max-scroll), 0))
	}
	return 0
}
func uiShopMouse(w *gui.Window, event int, a, b uintptr) int {
	if *bookWord(1047520) != 0 {
		return 1
	}
	p := uiInventoryPackedPoint(a)
	switch event {
	case 5:
		if uiShopInside(p) && uiShopMode() == 2 {
			uiShopBuyShow(p)
		}
	case 19, 20:
		if uiShopMode() == 2 {
			off := uintptr(1098584)
			if event == 20 {
				off = 1098588
			}
			uiShopWindow().Func94(gui.AsWindowEvent(16391, uintptr(*uiShopWord(off)), 0))
		}
	default:
		return 0
	}
	return 1
}
func uiShopHover(p image.Point) int {
	if uiShopMode() == 2 {
		col := (uint32(p.X) - *uiShopWord(1098380)) / 50
		// The original tooltip uses signed Y division, unlike the stock click.
		row := int32(uint32(p.Y)-*uiShopWord(1098384)+*uiShopWord(1107036)) / 50
		if col >= 6 {
			col = 5
		}
		if row >= 10 {
			row = 9
		}
		c := &uiShopGrid()[int(col)*10+int(row)]
		if c.Count != 0 {
			c.Drawable.NetCode32 = c.Codes[0]
			uiCursorTooltip(uiItemTooltip(c.Drawable))
		}
	}
	return 1
}
func uiShopMods(dr *client.Drawable) unsafe.Pointer {
	if uint32(dr.Class())&0x13001000 != 0 {
		return unsafe.Add(dr.C(), 432)
	}
	return nil
}
func uiShopBuyShow(p image.Point) {
	c := uiShopHitCell(p)
	if c.Count == 0 {
		return
	}
	gold := uint32(sub_4674A0())
	count := c.Count
	if c.Value != 0 && gold/c.Value < count {
		count = gold / c.Value
	}
	if count == 0 {
		uiShopGoldWarning(c.Value - gold)
		return
	}
	uiAmountPrice(1, c.Value)
	uiAmountShow(alloc.InternCString16(uiShopString("BuyLabel")), p.X, p.Y, c.Codes[c.Count-1], c.Drawable.TypeIDVal, uiShopMods(c.Drawable), count, 0, uiAmountNativeKey(uiAmountShopBuy), nil)
}
func uiShopCarryWarning() {
	uiTradeSound(925)
	Nox_xxx_printCentered_445490(uiShopString("pickup.c:CarryingTooMuch"))
}
func uiShopBuySingle(typ uint32, code uint16) {
	if uiInventoryCapacity(int32(typ), 1) != 0 {
		uiInventoryTrade(22, code)
	} else {
		uiShopCarryWarning()
	}
}
func uiShopBuyMultiple(typ, count uint32) {
	if uiInventoryCapacity(int32(typ), int32(count)) != 0 {
		uiShopMultiple(23, uint16(typ), byte(count))
	} else {
		uiShopCarryWarning()
	}
}
func uiShopBuyAccept(code uint16, typ, count uint32) {
	if count == 1 {
		uiShopBuySingle(typ, code)
	} else if count != 0 {
		uiShopBuyMultiple(typ, count)
	}
}
func uiShopMultiple(op byte, typ uint16, count byte) int {
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, []byte{201, op, byte(typ), byte(typ >> 8), count}))
}
func uiShopSellAccept(code, typ uint16, count uint32) int {
	*uiShopWord(1098616) = 0
	if count == 0 {
		return 0
	}
	if count == 1 {
		return uiInventoryTrade(24, code)
	}
	return uiShopMultiple(25, typ, byte(count))
}
func uiShopRepairAccept(code uint16) int {
	ret := uiInventoryTrade(26, code)
	*uiShopWord(1098620) = 0
	return ret
}
func uiShopGoldWarning(amount uint32) {
	uiTradeStoreText(1097352, 514, fmt.Sprintf(uiShopString("NotEnoughGold"), int32(amount)))
	Nox_xxx_dialogMsgBoxCreate_449A10(uiShopWindow(), uiShopString("ShopInformationTitle"), alloc.GoString16(uiTradeTextAt(1097352)), 33, nil, nil)
	uiTradeSound(925)
}
func uiShopSellShow(code, value uint32) int {
	p := GetClient().GetMousePos()
	ret := int(*uiShopWord(1098616))
	if ret == 1 {
		return ret
	}
	dr := uiInventoryItem(code)
	if dr == nil {
		return 0
	}
	uiAmountPrice(1, value)
	old := uiAmountItem()
	ret = uiAmountShow(alloc.InternCString16(uiShopString("SellLabel")), p.X, p.Y, code, dr.TypeIDVal, uiShopMods(dr), uint32(sub_467700(int(code))), 0, uiAmountNativeKey(uiAmountShopSell), uiAmountNativeKey(uiAmountShopSellCancel))
	if uiAmountItem() != old {
		*uiShopWord(1098616) = 1
	}
	return ret
}
func uiShopRepairShow(code, value uint32) {
	p := GetClient().GetMousePos()
	gold := uint32(sub_4674A0())
	if *uiShopWord(1098620) == 1 {
		return
	}
	dr := uiInventoryItem(code)
	if dr == nil {
		return
	}
	uiAmountPrice(1, value)
	if value > gold {
		uiShopGoldWarning(value - gold)
		sub_467680()
		return
	}
	old := uiAmountItem()
	uiAmountShow(alloc.InternCString16(uiShopString("RepairLabel")), p.X, p.Y, code, dr.TypeIDVal, uiShopMods(dr), 1, 0, uiAmountNativeKey(uiAmountShopRepair), uiAmountNativeKey(uiAmountShopRepairCancel))
	if uiAmountItem() != old {
		*uiShopWord(1098620) = 1
	}
}
