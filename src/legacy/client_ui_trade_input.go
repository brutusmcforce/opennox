package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func uiTradeInside(w *gui.Window, p image.Point) bool {
	origin, size := image.Point{}, image.Point{}
	if w != nil {
		origin = uiWindowPosition(w)
		size = w.SizeVal
	}
	return p.X >= origin.X && p.Y >= origin.Y && p.X <= origin.X+size.X && p.Y <= origin.Y+size.Y
}
func uiTradeHit(side int, p image.Point) *uiTradeCell {
	origin := uiWindowPosition(uiTradeWindow().ChildByID(uint(3704 + side)))
	cells := uiTradeGrid(side)
	for _, i := range uiTradeScanOrder {
		x, y := origin.X+(i/2)*50, origin.Y+(i%2)*50
		if p.X >= x && p.Y >= y && p.X <= x+50 && p.Y <= y+50 {
			return &cells[i]
		}
	}
	return nil
}
func uiTradeHover(p image.Point) int {
	side := -1
	w := uiTradeWindow()
	if uiTradeInside(w.ChildByID(3704), p) {
		side = 0
	} else if uiTradeInside(w.ChildByID(3705), p) {
		side = 1
	}
	if side >= 0 {
		if c := uiTradeHit(side, p); c != nil && c.Drawable != nil {
			c.Drawable.NetCode32 = c.Codes[0]
			uiCursorTooltip(uiItemTooltip(c.Drawable))
		}
	}
	return 1
}
func uiTradeMouse(w *gui.Window, event int, a, b uintptr) int {
	p := uiInventoryPackedPoint(a)
	main := uiTradeWindow()
	if event == 5 {
		if uiTradeInside(main.ChildByID(3704), p) {
			if c := uiTradeHit(0, p); c != nil && c.Count != 0 {
				main.Capture(true)
				InputSetKeyTimeoutLegacy(2)
				*(*[2]int32)(memmap.PtrOff(0x5D4594, 1319276)) = [2]int32{int32(p.X), int32(p.Y)}
				dr := c.Drawable
				dword_5d4594_1320968 = uint32(uiInventoryPointer(dr.C()))
				dr.NetCode32 = c.Codes[c.Count-1]
				c.Codes[c.Count-1] = 0
				Nox_xxx_cursorSetDraggedItem_477690(dr)
				c.Count--
				if c.Count == 0 {
					c.Drawable = nil
				}
				dword_5d4594_1320972 = uint32(uiInventoryPointer(unsafe.Pointer(c)))
				*memmap.PtrUint32(0x5D4594, 1320304) = 0
				uiTradeSound(791)
			}
		}
		return 1
	}
	if event <= 5 || event > 7 {
		return 0
	}
	if GetClient().Cli().GUI.Captured() == main {
		main.Capture(false)
	}
	dr := uiTradeDragged()
	if dr == nil {
		return 1
	}
	origin := memmap.Uint32(0x5D4594, 1320304)
	if origin < 2 {
		side := int(origin)
		key := byte(2 + side)
		dx := memmap.Int32(0x5D4594, 1319276) - int32(p.X)
		dy := memmap.Int32(0x5D4594, 1319280) - int32(p.Y)
		if !uiTradeInside(main.ChildByID(uint(3704+side)), p) || (!InputKeyCheckTimeoutLegacy(key, uint32(GetServer().S().TickRate())/3) && dx*dx+dy*dy < 100) {
			uiTradeRemoveRequest(dr)
		}
	}
	c := uiTradeSource()
	c.Codes[c.Count] = dr.NetCode32
	c.Count++
	if c.Count == 1 {
		c.Drawable = dr
	}
	Nox_xxx_cursorResetDraggedItem_4776A0()
	dword_5d4594_1320968 = 0
	dword_5d4594_1320972 = 0
	return 1
}
func uiTradePanel(w *gui.Window, event int, a, b uintptr) int {
	if event == 16391 {
		id := uiInventoryWindowValue(uint32(a)).ID()
		uiTradeSound(766)
		switch id {
		case 3708:
			uiTradeRequest(17)
		case 3710:
			uiTradeRequest(14)
		}
	}
	return 0
}
