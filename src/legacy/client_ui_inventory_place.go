package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func uiInventoryValidCell(col, row int) bool { return col >= 0 && col < 4 && row >= 0 && row < 21 }

func uiInventoryPlace(dr *client.Drawable, col, row int) int {
	if !uiInventoryValidCell(col, row) {
		return 0
	}
	cell := &uiInventoryGrid()[row+21*col]
	if cell.Count != 0 && (uint32(dr.ObjClass)&0x4000000 != 0 || cell.Drawable.TypeIDVal != dr.TypeIDVal) {
		return 0
	}
	if cell.Count >= 32 {
		return 0
	}
	if cell.Count == 0 {
		cell.Drawable = GetClient().Nox_new_drawable_for_thing(int(dr.TypeIDVal))
		if cell.Drawable == nil {
			uiInventoryError("DrawablesExhausted")
			return 0
		}
		cell.Drawable.ObjFlags |= 0x40000000
		uiInventoryCopyItem(cell.Drawable, dr)
	}
	cell.Codes[cell.Count] = dr.NetCode32
	cell.Count++
	cell.Equipped = 0
	for _, v := range uiInventoryEquipment() {
		for eq := uiInventoryDrawable(v); eq != nil; eq = uiInventoryNext(eq) {
			if eq.NetCode32 == dr.NetCode32 {
				cell.Equipped = 1
				if cell.Alternate != 0 {
					uiInventorySetAlternate(nil)
					cell.Alternate = 0
				}
				return 1
			}
		}
	}
	return 1
}

func uiInventorySetClick(col, row int) {
	dword_5d4594_1049796_inventory_click_column_index = uint32(col)
	dword_5d4594_1049800_inventory_click_row_index = uint32(row)
}
func uiInventoryDragCopy() {
	index := int(dword_5d4594_1049800_inventory_click_row_index) + 21*int(dword_5d4594_1049796_inventory_click_column_index)
	cell := &uiInventoryGrid()[index]
	if cell.Count == 0 {
		return
	}
	dr := GetClient().Nox_new_drawable_for_thing(int(cell.Drawable.TypeIDVal))
	uiInventorySetDragged(dr)
	if dr == nil {
		uiInventoryError("DrawablesExhausted")
		return
	}
	dr.ObjFlags |= 0x40000000
	dr.NetCode32 = cell.Codes[0]
	uiInventoryCopyItem(dr, cell.Drawable)
	uiInventoryRemoveStack(&uiInventoryLookup{Cell: cell})
}

func uiInventoryCapacity(typ, quantity int32) int32 {
	count := 0
	grid := uiInventoryGrid()
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			if cell.Count == 0 {
				count++
				continue
			}
			if cell.Drawable.TypeIDVal != uint32(typ) {
				continue
			}
			limit := int32(31)
			if uint32(cell.Drawable.ObjClass)&0x10 != 0 {
				limit = 3
				if noxflags.HasGame(6144) {
					limit = 9
				}
			}
			if uint32(cell.Drawable.ObjClass)&0x4000000 == 0 && int32(quantity)+int32(cell.Count) <= limit {
				count++
			}
		}
	}
	return int32(count)
}
func uiInventoryAlterWeapon() {
	playerDr := uiInventoryDrawable(memmap.Uint32(0x852978, 8))
	if playerDr == nil || GetClient().Cli().Cursor != 0 || quickbarAbilityAvailable(1) != 0 {
		return
	}
	player := uiMeterPlayer()
	if player == nil || *(*uint32)(unsafe.Add(playerDr.C(), 276)) == 34 {
		return
	}
	if int32(geometryRectInt((*[2]int32)(unsafe.Pointer(memmap.PtrOff(0x5D4594, 1062572))), (*[4]int32)(unsafe.Pointer(memmap.PtrOff(0x587000, 136336))))) == 1 {
		Nox_xxx_cursorSetDraggedItem_477690(nil)
	}
	alt := uiInventoryCellRef(uint32(dword_5d4594_1062480))
	dequip := func(dr *client.Drawable) {
		dword_5d4594_1062492 = uint32(uiInventoryPointer(dr.C()))
		uiInventoryDequipRequest(dr)
		audioEventPlay(895, 100, 0, 0)
	}
	if alt != nil && GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(int(alt.Drawable.TypeIDVal)) == 2 {
		typ := GetServer().S().Weapons.Sub_415840(2)
		if dr := uiInventoryEquippedType(uint32(typ)); dr != nil {
			dequip(dr)
		}
		return
	}
	mask := *(*uint32)(unsafe.Add(player, 4))
	for bit := uint32(1); bit < 27; bit++ {
		if bit == 1 || mask&(1<<bit) == 0 {
			continue
		}
		typ := GetServer().S().Weapons.Sub_415840(1 << bit)
		if dr := uiInventoryEquippedType(uint32(typ)); dr != nil {
			dequip(dr)
			return
		}
	}
	if alt != nil {
		alt.Drawable.NetCode32 = alt.Codes[0]
		uiInventoryEquipRequest(alt.Drawable)
		audioEventPlay(895, 100, 0, 0)
	}
}
