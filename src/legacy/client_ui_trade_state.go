package legacy

import (
	"encoding/binary"
	"fmt"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

type uiTradeCell struct {
	Drawable *client.Drawable
	Count    uint32
	Codes    [32]uint32
	Value    uint32
}

var _ [140 - unsafe.Sizeof(uiTradeCell{})]byte
var _ [unsafe.Sizeof(uiTradeCell{}) - 140]byte
var _ [136 - unsafe.Offsetof(uiTradeCell{}.Value)]byte
var _ [unsafe.Offsetof(uiTradeCell{}.Value) - 136]byte
var uiTradeScanOrder = [4]int{0, 2, 1, 3}

func uiTradeGrid(side int) []uiTradeCell {
	off := uintptr(1319284)
	if side != 0 {
		off = 1320308
	}
	return unsafe.Slice((*uiTradeCell)(memmap.PtrOff(0x5D4594, off)), 4)
}
func uiTradeWindow() *gui.Window       { return uiInventoryWindowValue(uint32(dword_5d4594_1320940)) }
func uiTradeActive() uint32            { return uint32(dword_5d4594_1320964) }
func uiTradeDragged() *client.Drawable { return uiInventoryDrawable(uint32(dword_5d4594_1320968)) }
func uiTradeSource() *uiTradeCell {
	return (*uiTradeCell)(unsafe.Pointer(uintptr(dword_5d4594_1320972)))
}
func uiTradeSound(id int) { audioEventPlay(int32(id), 100, 0, 0) }
func uiTradeString(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "guitrade.c")
}
func uiTradeRequest(op byte) int {
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, []byte{201, op}))
}
func uiTradeRemoveRequest(dr *client.Drawable) int {
	return uiInventoryTrade(16, uint16(drawableUnitCode(dr)))
}
func uiTradeDestroy() int {
	uiTradeWindow().Destroy()
	dword_5d4594_1320940 = 0
	dword_5d4594_1320964 = 0
	return 0
}
func uiTradeShow() int {
	if uiTradeActive() != 0 {
		return uiTradeWindow().ShowModal()
	}
	return 0
}
func uiTradeStart(data unsafe.Pointer) int {
	if dword_8531A0_2576 == 0 || uiTradeActive() == 1 {
		return 0
	}
	dword_5d4594_1320964 = 1
	uiTradeReset()
	uiTradeStoreText(1319844, 64, alloc.GoString16((*uint16)(unsafe.Add(data, 2))))
	w := uiTradeWindow()
	uiWindowEnable(w, 1)
	uiTradeShow()
	w.SetPos(image.Pt(198, 193))
	uiTradeSetText(w.ChildByID(3702), (*uint16)(unsafe.Pointer(uintptr(dword_8531A0_2576)+4704)))
	uiTradeSetText(w.ChildByID(3703), uiTradeTextAt(1319844))
	uiInventoryOpenWindow()
	return 1
}
func uiTradeReset() int {
	if dr := uiTradeDragged(); dr != nil {
		owned := false
		for side := 0; side < 2; side++ {
			for _, c := range uiTradeGrid(side) {
				if c.Drawable == dr {
					owned = true
				}
			}
		}
		Nox_xxx_cursorResetDraggedItem_4776A0()
		if !owned {
			GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		}
	}
	w := uiTradeWindow()
	if w != nil && GetClient().Cli().GUI.Captured() == w {
		w.Capture(false)
	}
	for _, off := range []uintptr{1320240, 1320868, 1320100} {
		*uiTradeTextAt(off) = 0
	}
	uiTradeSetText(w.ChildByID(3702), uiTradeTextAt(1320976))
	uiTradeSetText(w.ChildByID(3703), uiTradeTextAt(1320980))
	for side := 0; side < 2; side++ {
		cells := uiTradeGrid(side)
		for _, i := range uiTradeScanOrder {
			c := &cells[i]
			if c.Drawable != nil {
				GetClient().Nox_xxx_spriteDelete_45A4B0(c.Drawable)
			}
			c.Drawable = nil
			c.Count = 0
		}
	}
	// Static-text widgets borrow their C string; keep reset labels alive.
	text := alloc.InternCString16(fmt.Sprintf(alloc.GoString16(uiTradeTextAt(1319972)), 0))
	for _, id := range []uint{3711, 3712, 3713} {
		uiTradeSetText(w.ChildByID(id), text)
	}
	dword_5d4594_1320944 = 0
	dword_5d4594_1320948 = 0
	dword_5d4594_1320968 = 0
	dword_5d4594_1320972 = 0
	dword_5d4594_1320932 = 0
	dword_5d4594_1320936 = 0
	return 0
}
func uiTradeFinish() int {
	if uiTradeActive() == 0 {
		return 0
	}
	uiTradeReset()
	w := uiTradeWindow()
	w.Hide()
	uiWindowEnable(w, 0)
	dword_5d4594_1320964 = 0
	return uiInventoryCloseWindow()
}
func uiTradePrepare() int {
	ret := int(int32(uiTradeActive()))
	if ret != 0 && dword_8531A0_2576 != 0 {
		uiTradeReset()
		w := uiTradeWindow()
		uiTradeSetText(w.ChildByID(3702), (*uint16)(unsafe.Pointer(uintptr(dword_8531A0_2576)+4704)))
		ret = uiTradeSetText(w.ChildByID(3703), uiTradeTextAt(1319844))
	}
	return ret
}
func uiTradeFindCode(c *uiTradeCell, code uint32) int {
	if c.Count == 0 {
		return 0
	}
	for _, v := range c.Codes {
		if v == code {
			return 1
		}
	}
	return 0
}
func uiTradeRemoveCode(c *uiTradeCell, code uint32) uint32 {
	for i, v := range c.Codes {
		if v == code {
			ret := uint32(i)
			if i < 31 {
				ret = c.Codes[31]
				copy(c.Codes[i:], c.Codes[i+1:])
			}
			c.Codes[31] = 0
			return ret
		}
	}
	return 32
}
func uiTradeCellFits(typ uint32, c *uiTradeCell) int {
	if c.Count == 0 {
		return 1
	}
	if c.Count >= 32 {
		return 0
	}
	return bool2int(c.Drawable.TypeIDVal == typ && uint32(c.Drawable.Class())&0x13001000 == 0)
}
func uiTradeFindCell(side int, typ uint32) *uiTradeCell {
	cells := uiTradeGrid(side)
	for _, i := range uiTradeScanOrder {
		c := &cells[i]
		if c.Drawable != nil && c.Count < 32 && c.Drawable.TypeIDVal == typ && uint32(c.Drawable.Class())&0x13001000 == 0 {
			return c
		}
	}
	for _, i := range uiTradeScanOrder {
		c := &cells[i]
		if c.Drawable == nil {
			return c
		}
	}
	return nil
}
func uiTradeAdd(data unsafe.Pointer) uint32 {
	if uiTradeActive() == 0 {
		return 0
	}
	b := unsafe.Slice((*byte)(data), 15)
	dword_5d4594_1320944 = 0
	dword_5d4594_1320948 = 0
	side := 1
	selected := uint32(dword_5d4594_1320936)
	if b[2] == 1 {
		side = 0
		selected = uint32(dword_5d4594_1320932)
	}
	typ := uint32(binary.LittleEndian.Uint16(b[3:]))
	c := (*uiTradeCell)(unsafe.Pointer(uintptr(selected)))
	if c == nil || uiTradeCellFits(typ, c) == 0 {
		c = uiTradeFindCell(side, typ)
	}
	if c == nil || c.Count >= 32 {
		return 0
	}
	if c.Drawable == nil {
		dr := GetClient().Nox_new_drawable_for_thing(int(typ))
		if dr == nil {
			return 0
		}
		c.Drawable = dr
		if uint32(dr.Class())&0x13001000 != 0 {
			for i, id := range b[11:15] {
				*(*unsafe.Pointer)(unsafe.Add(dr.C(), 432+4*i)) = GetServer().S().Modif.Nox_xxx_modifGetDescById413330(int(id)).C()
			}
		}
		c.Count = 0
		c.Value = 0
	}
	c.Codes[c.Count] = uint32(binary.LittleEndian.Uint16(b[5:]))
	c.Count++
	c.Value += binary.LittleEndian.Uint32(b[7:])
	dword_5d4594_1320932 = 0
	dword_5d4594_1320936 = 0
	return c.Value
}
func uiTradeRemove(data unsafe.Pointer) uint32 {
	if uiTradeActive() == 0 {
		return 0
	}
	code := uint32(binary.LittleEndian.Uint16(unsafe.Slice((*byte)(unsafe.Add(data, 2)), 2)))
	for side := 0; side < 2; side++ {
		cells := uiTradeGrid(side)
		for _, i := range uiTradeScanOrder {
			c := &cells[i]
			if uiTradeFindCode(c, code) == 0 {
				continue
			}
			c.Value -= c.Value / c.Count
			uiTradeRemoveCode(c, code)
			c.Count--
			ret := c.Count
			if ret == 0 {
				ret = uint32(GetClient().Nox_xxx_spriteDelete_45A4B0(c.Drawable))
				c.Drawable = nil
			}
			return ret
		}
	}
	Nox_xxx_printCentered_445490(uiTradeString("TradeGUIItemNotFound"))
	return 0
}
func uiTradeMoney(data unsafe.Pointer) int {
	if uiTradeActive() == 0 {
		return 0
	}
	b := unsafe.Slice((*byte)(data), 14)
	value := func(i int) int32 { return int32(binary.LittleEndian.Uint32(b[i:])) }
	uiTradeStoreText(1320240, 32, fmt.Sprint(value(2)))
	text := alloc.GoString16(uiTradeTextAt(1320984))
	if value(6) != 0 {
		text = fmt.Sprintf("(%d)", value(6))
	}
	uiTradeStoreText(1320868, 32, text)
	return uiTradeStoreText(1320100, 32, fmt.Sprint(value(10)))
}
func uiTradeAcceptance(data unsafe.Pointer) uint32 {
	if uiTradeActive() == 0 {
		return 0
	}
	bits := *(*byte)(unsafe.Add(data, 2))
	dword_5d4594_1320944 = uint32(bits & 1)
	dword_5d4594_1320948 = uint32((bits >> 1) & 1)
	return uiInventoryPointer(data)
}
