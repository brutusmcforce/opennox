package legacy

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func uiAmountWindow() *gui.Window {
	return (*gui.Window)(legacyGlobals.nox_gui_itemAmount_dialog_1319228)
}
func uiAmountItem() *client.Drawable {
	return (*client.Drawable)(legacyGlobals.nox_gui_itemAmount_item_1319256)
}
func uiTradeTextAt(off uintptr) *uint16 { return (*uint16)(memmap.PtrOff(0x5D4594, off)) }
func uiTradeStoreText(off uintptr, capacity int, text string) int {
	return alloc.StrCopy16(unsafe.Slice(uiTradeTextAt(off), capacity), text)
}
func uiTradeSetText(w *gui.Window, text *uint16) int {
	if w != nil {
		w.Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(text)), 0))
	}
	return 0
}
func uiTradeViewportInit(off uintptr) {
	for _, v := range [][2]uint32{{0, 0}, {4, 0}, {8, uint32(nox_win_width)}, {12, uint32(nox_win_height)}, {32, uint32(nox_win_width)}, {36, uint32(nox_win_height)}, {16, 0}, {20, 0}} {
		*memmap.PtrUint32(0x5D4594, off+uintptr(v[0])) = v[1]
	}
}
func uiAmountToggle() {
	w := uiAmountWindow()
	if dword_5d4594_1319268 == 1 {
		w.Hide()
		uiWindowEnable(w, 0)
		w.StackPop()
		if dr := uiAmountItem(); dr != nil {
			GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		}
		legacyGlobals.nox_gui_itemAmount_item_1319256 = nil
		dword_5d4594_1319268 = 0
	} else {
		uiWindowEnable(w, 1)
		w.StackPush()
		w.ShowModal()
		dword_5d4594_1319268 = 1
	}
}
func uiAmountMouse(w *gui.Window, event int, a, b uintptr) int {
	switch event {
	case 5, 9, 13:
		p := uiInventoryPackedPoint(a)
		if !uiWindowPointIn(w, int32(p.X), int32(p.Y)) {
			uiAmountCancel()
		}
		return 1
	case 6, 7, 10, 11, 14, 15:
		return 1
	}
	return 0
}
func uiAmountCount() uint32 {
	w := uiInventoryWindowValue(uint32(dword_5d4594_1319232))
	return uint32(textDecimal((*uint16)(unsafe.Pointer(uiWindowText(w)))))
}
func uiAmountCallback(off uintptr) {
	count := uiAmountCount()
	if count > uint32(dword_5d4594_1319248) {
		count = uint32(dword_5d4594_1319248)
	}
	if fn := *memmap.PtrPtr(0x5D4594, off); fn != nil {
		// Keep the temporary point in the shared allocator through callback completion.
		p, free := alloc.New([2]int32{})
		defer free()
		pos := uiWindowPosition(uiAmountWindow())
		*p = [2]int32{int32(pos.X) + int32(dword_587000_183456), int32(pos.Y) + int32(dword_587000_183460)}
		uiAmountCall(fn, uintptr(unsafe.Pointer(p)), uintptr(memmap.Uint32(0x5D4594, 1319244)), uintptr(memmap.Uint32(0x5D4594, 1319240)), uintptr(count), uintptr(memmap.Uint32(0x5D4594, 1319252)))
	}
	uiAmountToggle()
}
func uiAmountCancel() int {
	if dword_5d4594_1319268 != 1 {
		return 0
	}
	uiAmountCallback(1319100)
	return 1
}
func uiAmountInit() int {
	dword_5d4594_1319264 = 0
	w := Nox_new_window_from_file("MultMove.wnd", uiInventoryWindowEvent(uiAmountPanel))
	legacyGlobals.nox_gui_itemAmount_dialog_1319228 = w.C()
	if w == nil {
		return 0
	}
	w.SetAllFuncs(uiInventoryWindowEvent(uiAmountMouse), func(w *gui.Window, _ *gui.WindowData) int { return uiAmountDraw(w) }, nil)
	dword_5d4594_1319232 = uint32(uiInventoryPointer(w.ChildByID(3601).C()))
	dword_5d4594_1319236 = uint32(uiInventoryPointer(w.ChildByID(3607).C()))
	w.Hide()
	uiWindowEnable(w, 0)
	uiTradeViewportInit(1319108)
	for i, name := range []string{"MultiMoveBase", "MultiMoveUpLit", "MultiMoveDownLit", "MultiMoveYesPressed", "MultiMoveNoPressed", "MultiMoveBaseNoTag", "MultiMoveYesPressedNoTag", "MultiMoveNoPressedNoTag"} {
		*memmap.PtrUint32(0x5D4594, 1319196+uintptr(4*i)) = uiMeterLoadImage(name)
	}
	return 1
}
func uiAmountDraw(w *gui.Window) int {
	nox_client_drawRectFilledAlpha_49CF10(0, 0, int(nox_win_width), int(nox_win_height))
	pos := uiWindowPosition(w)
	priced := dword_5d4594_1319264 != 0
	off := uintptr(1319196)
	if !priced {
		off = 1319216
	}
	uiMeterImage(memmap.Uint32(0x5D4594, off), pos)
	dr := uiAmountItem()
	dr.PosVec = pos.Add(image.Pt(int(int32(dword_587000_183456)), int(int32(dword_587000_183460))))
	dr.CallDraw((*noxrender.Viewport)(memmap.PtrOff(0x5D4594, 1319108)))
	for _, v := range [][3]uintptr{{3603, 1319204, 1319204}, {3602, 1319200, 1319200}, {3604, 1319208, 1319220}, {3605, 1319212, 1319224}} {
		if uiAmountWindow().ChildByID(uint(v[0])).DrawData().Field0&4 != 0 {
			off := v[1]
			if !priced {
				off = v[2]
			}
			uiMeterImage(memmap.Uint32(0x5D4594, off), pos)
		}
	}
	return 1
}
func uiAmountPanel(w *gui.Window, event int, a, b uintptr) int {
	if event != 16391 || dword_5d4594_1319268 != 1 {
		return 0
	}
	id := uiInventoryWindowValue(uint32(a)).ID()
	uiTradeSound(766)
	switch id {
	case 3602:
		count := uiAmountCount() + 1
		if count > uint32(dword_5d4594_1319248) {
			return 0
		}
		uiAmountUpdateText(count)
	case 3603:
		count := int32(uiAmountCount())
		if count > 1 {
			uiAmountUpdateText(uint32(count - 1))
		}
	case 3604, 3606:
		uiAmountCallback(1319160)
	case 3605:
		uiAmountCancel()
	}
	return 0
}
func uiAmountUpdateText(count uint32) {
	uiTradeStoreText(1319164, 16, fmt.Sprint(int32(count)))
	uiTradeSetText(uiInventoryWindowValue(uint32(dword_5d4594_1319232)), uiTradeTextAt(1319164))
	if dword_5d4594_1319264 != 0 {
		uiTradeStoreText(1319068, 16, fmt.Sprint(int32(uint32(dword_5d4594_1319260)*count)))
		uiTradeSetText(uiInventoryWindowValue(uint32(dword_5d4594_1319236)), uiTradeTextAt(1319068))
	}
}
func uiAmountFree() {
	if dr := uiAmountItem(); dr != nil {
		GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
	}
	legacyGlobals.nox_gui_itemAmount_item_1319256 = nil
	uiAmountWindow().Destroy()
	legacyGlobals.nox_gui_itemAmount_dialog_1319228 = nil
	dword_5d4594_1319232 = 0
	dword_5d4594_1319236 = 0
	dword_5d4594_1319264 = 0
	dword_5d4594_1319268 = 0
}
func uiAmountShow(title *uint16, x, y int, code, typ uint32, mods unsafe.Pointer, maximum, extra uint32, accept, cancel unsafe.Pointer) int {
	dr := GetClient().Nox_new_drawable_for_thing(int(typ))
	if dr == nil {
		return 0
	}
	if dword_5d4594_1319268 == 1 {
		uiAmountToggle()
	}
	legacyGlobals.nox_gui_itemAmount_item_1319256 = dr.C()
	uiTradeSetText(uiAmountWindow().ChildByID(3606), title)
	dr.ObjFlags |= 0x40000000
	if mods != nil {
		copy(unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 20), unsafe.Slice((*byte)(mods), 20))
	}
	*memmap.PtrPtr(0x5D4594, 1319160) = accept
	*memmap.PtrPtr(0x5D4594, 1319100) = cancel
	*memmap.PtrUint32(0x5D4594, 1319240) = typ
	*memmap.PtrUint32(0x5D4594, 1319244) = code
	dword_5d4594_1319248 = uint32(maximum)
	*memmap.PtrUint32(0x5D4594, 1319252) = extra
	uiAmountToggle()
	uiAmountPosition(x, y)
	uiTradeStoreText(1319164, 16, "1")
	uiTradeSetText(uiInventoryWindowValue(uint32(dword_5d4594_1319232)), uiTradeTextAt(1319164))
	text := alloc.GoString16(uiTradeTextAt(1319272))
	if dword_5d4594_1319264 != 0 {
		text = fmt.Sprint(int32(dword_5d4594_1319260))
	}
	uiTradeStoreText(1319068, 16, text)
	return uiTradeSetText(uiInventoryWindowValue(uint32(dword_5d4594_1319236)), uiTradeTextAt(1319068))
}
func uiAmountPosition(x, y int) int {
	w := uiAmountWindow()
	x -= int(int32(dword_587000_183456))
	y -= int(int32(dword_587000_183460))
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+w.SizeVal.X >= int(nox_win_width) {
		x = int(nox_win_width) - w.SizeVal.X
	}
	if y+w.SizeVal.Y >= int(nox_win_height) {
		y = int(nox_win_height) - w.SizeVal.Y
	}
	w.SetPos(image.Pt(x, y))
	return 0
}
func uiAmountPrice(enabled, unit uint32) uint32 {
	dword_5d4594_1319264 = uint32(enabled)
	dword_5d4594_1319260 = uint32(unit)
	return enabled
}
