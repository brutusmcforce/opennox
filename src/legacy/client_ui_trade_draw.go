package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

func uiTradeInit() int {
	w := Nox_new_window_from_file("Trade.wnd", uiInventoryWindowEvent(uiTradePanel))
	dword_5d4594_1320940 = uint32(uiInventoryPointer(w.C()))
	if w == nil {
		return 0
	}
	w.SetAllFuncs(uiInventoryWindowEvent(uiTradeMouse), func(_ *gui.Window, _ *gui.WindowData) int { return uiTradeDraw() }, nil)
	w.DrawData().SetTooltip(GetServer().S().Strings(), uiTradeString("TradeMain"))
	for _, v := range []struct {
		id  uint
		key string
	}{{3702, "TradePlayerName"}, {3703, "TradeVendorName"}, {3708, "TradePlayerAccept"}, {3709, "TradeVendorAccept"}, {3710, "TradeCancel"}} {
		w.ChildByID(v.id).DrawData().SetTooltip(GetServer().S().Strings(), uiTradeString(v.key))
	}
	for _, id := range []uint{3704, 3705} {
		w.ChildByID(id).SetTooltipFunc(uiAmountNativeKey(uiAmountTradeTooltip))
	}
	w.Hide()
	uiWindowEnable(w, 0)
	for side := 0; side < 2; side++ {
		cells := uiTradeGrid(side)
		for _, i := range uiTradeScanOrder {
			cells[i].Drawable = nil
			cells[i].Count = 0
		}
	}
	uiTradeStoreText(1319972, 64, uiTradeString("TotalValueLabel"))
	uiTradeViewportInit(1320188)
	for i, name := range []string{"TradeBase", "TradeLeftAcceptPushed", "TradeLeftAcceptLit", "TradeRightAcceptLit", "TradeCancelLit", "TradeGold"} {
		*memmap.PtrUint32(0x5D4594, 1320164+uintptr(4*i)) = uiMeterLoadImage(name)
	}
	return 1
}
func uiTradeDraw() int {
	w := uiTradeWindow()
	pos := uiWindowPosition(w)
	draw := func(off uintptr, p image.Point) { uiMeterImage(memmap.Uint32(0x5D4594, off), p) }
	draw(1320164, pos)
	for _, v := range [][2]uintptr{{3711, 1320240}, {3712, 1320868}, {3713, 1320100}} {
		uiTradeSetText(w.ChildByID(uint(v[0])), uiTradeTextAt(v[1]))
	}
	for _, off := range []uintptr{183696, 183704} {
		draw(1320184, pos.Add(image.Pt(int(memmap.Int32(0x587000, off))-64, int(memmap.Int32(0x587000, off+4))-64)))
	}
	if dword_5d4594_1320944 != 0 {
		draw(1320172, pos)
	} else if w.ChildByID(3708).DrawData().Field0&4 != 0 {
		draw(1320168, pos)
	}
	if dword_5d4594_1320948 != 0 {
		draw(1320176, pos)
	}
	if memmap.Uint32(0x5D4594, 1320960) != 0 || w.ChildByID(3710).DrawData().Field0&4 != 0 {
		draw(1320180, pos)
	}
	height := nox_xxx_guiFontHeightMB_43F320(nil)
	for side := 0; side < 2; side++ {
		origin := uiWindowPosition(w.ChildByID(uint(3704 + side)))
		cells := uiTradeGrid(side)
		for _, i := range uiTradeScanOrder {
			c := &cells[i]
			if c.Count == 0 {
				continue
			}
			p := origin.Add(image.Pt((i/2)*50, (i%2)*50))
			c.Drawable.PosVec = p.Add(image.Pt(25, 25))
			c.Drawable.CallDraw((*noxrender.Viewport)(memmap.PtrOff(0x5D4594, 1320188)))
			nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
			GetClient().R2().DrawString(GetClient().R2().GetFonts().AsFont(nil), fmt.Sprint(int32(c.Count)), p.Add(image.Pt(5, 5)))
			nox_xxx_drawSetTextColor_434390(int(nox_color_yellow_2589772))
			GetClient().R2().DrawString(GetClient().R2().GetFonts().AsFont(nil), fmt.Sprint(int32(c.Value)), p.Add(image.Pt(5, 50-height-5)))
		}
	}
	return 1
}
