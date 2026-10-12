package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

// The shop shares the trade cell layout, but Value is a unit price.
type uiShopCell = uiTradeCell

func uiShopWord(off uintptr) *uint32 {
	switch off {
	case 1098456:
		return (*uint32)(&dword_5d4594_1098456)
	case 1098576:
		return (*uint32)(&dword_5d4594_1098576)
	case 1098580:
		return (*uint32)(&dword_5d4594_1098580)
	case 1098592:
		return (*uint32)(&dword_5d4594_1098592)
	case 1098596:
		return (*uint32)(&dword_5d4594_1098596)
	case 1098600:
		return (*uint32)(&dword_5d4594_1098600)
	case 1098604:
		return (*uint32)(&dword_5d4594_1098604)
	case 1098616:
		return (*uint32)(&dword_5d4594_1098616)
	case 1098620:
		return (*uint32)(&dword_5d4594_1098620)
	case 1098624:
		return (*uint32)(&dword_5d4594_1098624)
	case 1098628:
		return (*uint32)(&dword_5d4594_1098628)
	case 1107036:
		return (*uint32)(&dword_5d4594_1107036)
	}
	return memmap.PtrUint32(0x5D4594, off)
}
func uiShopWindow() *gui.Window { return uiInventoryWindowValue(*uiShopWord(1098576)) }
func uiShopActive() uint32      { return *uiShopWord(1098624) }
func uiShopMode() uint32        { return *uiShopWord(1098628) }
func uiShopSetMode(mode uint32) {
	if uiShopMode() == 4 && mode != 4 {
		sub_467680()
	}
	*uiShopWord(1098628) = mode
}
func uiShopGrid() []uiShopCell {
	return unsafe.Slice((*uiShopCell)(memmap.PtrOff(0x5D4594, 1098636)), 60)
}
func uiShopFind(code uint32) *uiShopCell {
	grid := uiShopGrid()
	for row := 0; row < 10; row++ {
		for col := 0; col < 6; col++ {
			c := &grid[col*10+row]
			if c.Count != 0 {
				for _, v := range c.Codes {
					if v == code {
						return c
					}
				}
			}
		}
	}
	return nil
}
func uiShopDrawable(code uint32) *client.Drawable {
	if uiShopActive() != 0 {
		if c := uiShopFind(code); c != nil {
			return c.Drawable
		}
	}
	return nil
}
func uiShopEmpty() *uiShopCell {
	grid := uiShopGrid()
	for row := 0; row < 10; row++ {
		for col := 0; col < 6; col++ {
			c := &grid[col*10+row]
			if c.Count == 0 {
				return c
			}
		}
	}
	return nil
}
func uiShopFindType(typ uint32) *uiShopCell {
	grid := uiShopGrid()
	for row := 0; row < 10; row++ {
		for col := 0; col < 6; col++ {
			c := &grid[col*10+row]
			if c.Count != 0 && c.Count < 32 && c.Drawable.TypeIDVal == typ && uint32(c.Drawable.Class())&0x04000000 == 0 {
				return c
			}
		}
	}
	return uiShopEmpty()
}
func uiShopAdd(typ, code, value uint32, health uint16, mods unsafe.Pointer) uint32 {
	c := uiShopFindType(typ)
	if c == nil {
		return 0
	}
	if c.Drawable == nil {
		dr := GetClient().Nox_new_drawable_for_thing(int(typ))
		c.Drawable = dr
		if dr == nil {
			return 0
		}
		*(*uint32)(unsafe.Add(dr.C(), 120)) |= 0x40000000
		*(*uint16)(unsafe.Add(dr.C(), 292)) = health
		*(*uint16)(unsafe.Add(dr.C(), 294)) = health
		if uint32(dr.Class())&0x13001000 != 0 {
			for i, id := range unsafe.Slice((*byte)(mods), 4) {
				var p unsafe.Pointer
				if id != 255 {
					p = GetServer().S().Modif.Nox_xxx_modifGetDescById413330(int(id)).C()
				}
				*(*unsafe.Pointer)(unsafe.Add(dr.C(), 432+4*i)) = p
			}
		}
		c.Count = 0
	}
	if c.Count >= 32 {
		return 0
	}
	c.Codes[c.Count] = code
	c.Count++
	c.Value = value
	return c.Count
}
func uiShopRemove(code uint32) uint32 {
	c := uiShopFind(code)
	if c == nil {
		return 0
	}
	uiTradeRemoveCode(c, code)
	c.Count--
	if c.Count != 0 {
		return c.Count
	}
	v := GetClient().Nox_xxx_spriteDelete_45A4B0(c.Drawable)
	c.Drawable = nil
	c.Value = 0
	return uint32(v)
}
func uiShopClear() uint32 {
	grid := uiShopGrid()
	for row := 0; row < 10; row++ {
		for col := 0; col < 6; col++ {
			c := &grid[col*10+row]
			if c.Drawable != nil {
				GetClient().Nox_xxx_spriteDelete_45A4B0(c.Drawable)
			}
			*c = uiShopCell{}
		}
	}
	*uiShopWord(1107036) = 0
	return uiInventoryPointer(memmap.PtrOff(0x5D4594, 1100036))
}
func uiShopCancelQuantity() {
	cb := uiAmountAccept
	if cb == uiAmountNativeKey(uiAmountShopBuy) || cb == uiAmountNativeKey(uiAmountShopSell) || cb == uiAmountNativeKey(uiAmountShopRepair) {
		uiAmountCancel()
	}
}
func uiShopDestroy() int {
	uiShopCancelQuantity()
	uiShopClear()
	Dialogs.Sub_44D8F0()
	uiShopWindow().Destroy()
	for _, off := range []uintptr{1098576, 1098624, 1098596, 1098600, 1098604, 1098608, 1098616, 1098620} {
		*uiShopWord(off) = 0
	}
	*uiShopWord(1098628) = 1
	return 0
}
func uiShopCancelRequest() int {
	if uiShopActive() == 0 {
		return 0
	}
	uiTradeRequest(18)
	sub_467680()
	return 1
}
func uiShopClose() {
	if uiShopActive() == 0 {
		return
	}
	uiShopCancelQuantity()
	sub_467680()
	for _, off := range []uintptr{1098624, 1098628, 1098616, 1098620} {
		*uiShopWord(off) = 0
	}
	uiShopClear()
	Dialogs.Sub_44D8F0()
	w := uiShopWindow()
	w.Hide()
	uiWindowEnable(w, 0)
	uiInventoryCloseWindow()
	nox_client_setCursorType_477610(0)
	if Nox_client_getRenderGUI() == 0 && *uiShopWord(1098612) == 1 {
		Nox_client_setRenderGUI(1)
	}
}
func uiShopInside(p image.Point) bool {
	return int32(p.X) >= int32(*uiShopWord(1098380)) && int32(p.Y) >= int32(*uiShopWord(1098384)) && int32(p.X) <= int32(*uiShopWord(1098388)) && int32(p.Y) <= int32(*uiShopWord(1098392))
}

// Coordinate division is unsigned for stock clicks, as in the original C words.
func uiShopHitCell(p image.Point) *uiShopCell {
	col := (uint32(p.X) - *uiShopWord(1098380)) / 50
	row := (uint32(p.Y) - *uiShopWord(1098384) + *uiShopWord(1107036)) / 50
	if col >= 6 {
		col = 5
	}
	if row >= 10 {
		row = 9
	}
	return &uiShopGrid()[col*10+row]
}
func uiShopHit(p image.Point) *client.Drawable {
	if !uiShopInside(p) {
		return nil
	}
	c := uiShopHitCell(p)
	if c.Count == 0 {
		return nil
	}
	c.Drawable.NetCode32 = c.Codes[0]
	return c.Drawable
}
