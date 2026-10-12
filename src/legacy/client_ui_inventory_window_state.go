package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func uiInventoryCloseIdentify() int {
	w := uiInventoryIdentifyWindow()
	if w == nil || w.Flags.IsHidden() {
		return 0
	}
	w.Hide()
	dword_5d4594_1063116, dword_5d4594_1063120 = 0, 0
	alloc.StrCopyZero16(unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1063124)), 256), uiInventoryText("thing.db:IdentifyDescription"))
	w.ChildByID(9156).Func94(&gui.RawEvent{Event: 16399})
	uiInventoryMainWindow().Capture(false)
	dword_5d4594_1049864 = 0
	nox_client_setCursorType_477610(0)
	return 1
}

func uiInventoryOpenIdentify() int {
	uiInventoryIdentifyWindow().Show()
	dword_5d4594_1049864 = 5
	nox_client_setCursorType_477610(6)
	if !uiInventoryMainWindow().Capture(true) {
		return -4
	}
	return 0
}
func uiInventoryWindowOpenState() bool {
	state := memmap.Uint8(0x5D4594, 1049868)
	return state == 1 || state == 2
}

func uiInventoryOpenWindow() int {
	if v := int(sessionQuitShown()); v != 0 {
		return v
	}
	if v := optionsVisible(); v != 0 {
		return v
	}
	if v := int(nox_xxx_guiCursor_477600()); v != 0 {
		return v
	}
	if v := Nox_xxx_playerAnimCheck_4372B0(); v != 0 {
		return v
	}
	if v := Nox_xxx_get_57AF20(); v != 0 {
		return v
	}
	state := memmap.PtrUint8(0x5D4594, 1049868)
	if *state == 0 || *state == 3 {
		*state = 1
		audioEventPlay(789, 100, 0, 0)
	}
	dword_5d4594_1062512 = dword_5d4594_1062516
	return int(int32(dword_5d4594_1062516))
}

func uiInventoryCloseWindow() int {
	if uiInventoryMode() == 6 {
		return 1
	}
	if !uiInventoryWindowOpenState() {
		return 0
	}
	*memmap.PtrUint8(0x5D4594, 1049868) = 3
	audioEventPlay(790, 100, 0, 0)
	if uiInventoryMode() == 5 {
		uiInventoryCloseIdentify()
	}
	uiInventoryCancelDrag()
	return 1
}

func uiInventoryToggleWindow() int {
	if uiInventoryWindowOpenState() {
		return uiInventoryCloseWindow()
	}
	return uiInventoryOpenWindow()
}
func uiInventoryRepairMode() int {
	uiInventoryCloseIdentify()
	dword_5d4594_1049864 = 6
	nox_client_setCursorType_477610(8)
	if uiInventoryWindowOpenState() {
		return 1
	}
	return uiInventoryOpenWindow()
}

func uiInventoryResetClosedScroll() int {
	if uiInventoryWindowOpenState() {
		return 1
	}
	dword_5d4594_1062516 = 0
	w := (*gui.Window)(unsafe.Pointer(uintptr(dword_5d4594_1062508)))
	if w == nil {
		return 0
	}
	max := (*gui.SliderData)(w.WidgetData).Max
	return gui.EventRespInt(w.Func94(&gui.RawEvent{Event: 16394, Arg1: uintptr(uint32(max))}))
}
func uiInventorySetWindowLevel(level int) int {
	dword_5d4594_1049844 = uint32(level)
	return int(nox_xxx_inventoryNameSignInit_4671E0())
}

func uiInventoryCancelDrag() int {
	ret := 0
	if dr := uiInventoryDragged(); dr != nil {
		if dword_5d4594_1049856 == 0 && uiInventoryPlace(dr, int(int32(dword_5d4594_1049796_inventory_click_column_index)), int(int32(dword_5d4594_1049800_inventory_click_row_index))) == 0 {
			nox_xxx_spritePickup_461660(int32(dr.NetCode32), int32(dr.TypeIDVal), unsafe.Add(dr.C(), 432))
			if found := uiInventoryFindCode(dr.NetCode32); found != nil {
				restored := found.Cell
				restored.Equipped = 0
				for _, v := range uiInventoryEquipment() {
					for eq := uiInventoryDrawable(v); eq != nil; eq = uiInventoryNext(eq) {
						if eq.NetCode32 != dr.NetCode32 {
							continue
						}
						restored.Equipped = 1
						if restored.Alternate != 0 {
							uiInventorySetAlternate(nil)
							restored.Alternate = 0
						}
						break
					}
				}
			}
		}
		// Inventory restoration copies the item into another drawable. Only an
		// equipment drag borrows its drawable; inventory drags own this temporary.
		if dword_5d4594_1049856 == 0 {
			GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		}
		uiInventorySetDragged(nil)
		dword_5d4594_1049856 = 0
		Nox_xxx_cursorSetDraggedItem_477690(nil)
		ret = 1
	}
	parent := (*gui.Window)(unsafe.Pointer(legacyGlobals.dword_5d4594_1062452))
	captured := GetClient().Cli().GUI.Captured()
	if uiWindowIsChild(parent, captured) {
		captured.Capture(false)
	}
	return ret
}
