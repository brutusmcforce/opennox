package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/layout"
)

func uiInventoryEquipmentSlot(dr *client.Drawable) int {
	class, sub := uint32(dr.ObjClass), uint32(dr.ObjSubClass)
	if class&0x1000000 != 0 && sub&2 != 0 {
		return 0
	}
	if class&0x2000000 != 0 {
		switch {
		case sub&1 != 0:
			return 1
		case sub&0x144 != 0:
			return 2
		case sub&0x90 != 0:
			return 3
		case sub&0x20 != 0:
			return 4
		case sub&2 != 0:
			return 8
		case sub&8 != 0:
			return 5
		}
	}
	if class&0x1000000 != 0 {
		if sub&4 != 0 {
			return 8
		}
		return 7
	}
	if class&0x1000 != 0 {
		return 7
	}
	return 9
}
func uiInventoryInsertEquipment(dr *client.Drawable, slot int) *client.Drawable {
	equip := uiInventoryEquipment()
	head := uiInventoryDrawable(equip[slot])
	var after *client.Drawable
	if uint32(dr.ObjClass)&0x2000000 != 0 {
		sub := uint32(dr.ObjSubClass)
		if sub&0x140 != 0 || sub&0x10 != 0 {
			after = head
			if after != nil {
				for uiInventoryNext(after) != nil {
					after = uiInventoryNext(after)
				}
				if sub&0x140 != 0 && sub&0x40 != 0 && uint32(after.ObjClass)&0x2000000 != 0 && uint32(after.ObjSubClass)&0x100 != 0 {
					after = uiInventoryPrev(after)
				}
			}
		}
	}
	if after != nil {
		next := uiInventoryNext(after)
		if next != nil {
			uiInventorySetPrev(next, dr)
		}
		uiInventorySetNext(dr, next)
		uiInventorySetNext(after, dr)
		uiInventorySetPrev(dr, after)
		return after
	}
	uiInventorySetPrev(dr, nil)
	uiInventorySetNext(dr, head)
	if head != nil {
		uiInventorySetPrev(head, dr)
	}
	equip[slot] = uiInventoryPointer(dr.C())
	return head
}
func uiInventoryUnlink(code uint32) *client.Drawable {
	equip := uiInventoryEquipment()
	for slot, v := range equip {
		for dr := uiInventoryDrawable(v); dr != nil; dr = uiInventoryNext(dr) {
			if dr.NetCode32 != code {
				continue
			}
			prev, next := uiInventoryPrev(dr), uiInventoryNext(dr)
			if prev != nil {
				uiInventorySetNext(prev, next)
			} else {
				equip[slot] = uiInventoryPointer(unsafe.Pointer(next))
			}
			if next != nil {
				uiInventorySetPrev(next, prev)
			}
			ammo := GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(int(dr.TypeIDVal))
			if uint32(dr.ObjClass)&0x1000 != 0 || ammo == 2 || ammo == 128 {
				sub_470D70()
			}
			return dr
		}
	}
	return nil
}

func sub_462040(code int32) {
	found := uiInventoryFindCode(uint32(code))
	var src *client.Drawable
	if found != nil {
		src = found.Cell.Drawable
	} else {
		src = uiInventoryDragged()
		if src == nil || src.NetCode32 != uint32(code) {
			uiInventoryError("EquippedNotFound")
			return
		}
	}
	slot := uiInventoryEquipmentSlot(src)
	if slot == 9 {
		uiInventoryError("TooManyEquipped")
		return
	}
	dr := GetClient().Nox_new_drawable_for_thing(int(src.TypeIDVal))
	if dr == nil {
		uiInventoryError("DrawablesExhausted")
		return
	}
	dr.NetCode32 = uint32(code)
	dr.ObjFlags |= 0x40000000
	uiInventoryCopyItem(dr, src)
	uiInventoryInsertEquipment(dr, slot)
	if found != nil {
		found.Cell.Equipped = 1
		if found.Cell.Alternate != 0 {
			uiInventorySetAlternate(nil)
			found.Cell.Alternate = 0
		}
	}
	if uint32(dr.ObjClass)&0x1000000 != 0 && uint32(dr.ObjSubClass)&0xC != 0 {
		var quiver *uiInventoryCell
		if dword_5d4594_1062488 != 0 {
			if q := uiInventoryFindCode(uint32(dword_5d4594_1062488)); q != nil {
				quiver = q.Cell
			}
		}
		if quiver == nil {
			grid := uiInventoryGrid()
			for row := 0; row < 20 && quiver == nil; row++ {
				for col := 0; col < 4; col++ {
					cell := &grid[row+21*col]
					if cell.Count != 0 && uint32(cell.Drawable.ObjClass)&0x1000000 != 0 && uint32(cell.Drawable.ObjSubClass) == 2 {
						quiver = cell
						break
					}
				}
			}
		}
		if quiver != nil {
			quiver.Drawable.NetCode32 = quiver.Codes[0]
			uiInventoryEquipRequest(quiver.Drawable)
		}
	}
	if slot == 0 {
		dword_5d4594_1062488 = uint32(dr.NetCode32)
	}
	item := layout.As[client.DrawableUnionItem](unsafe.Pointer(&dr.Union))
	if charge := item.Field_112_0; charge >= 0 {
		sub_470D90(int(charge), int(item.Field_112_2))
	}
	if dword_5d4594_1062496 != 0 {
		if next := uiInventoryFindCode(uint32(dword_5d4594_1062496)); next != nil {
			next.Cell.Alternate = 1
			uiInventorySetAlternate(next.Cell)
			dword_5d4594_1062496 = 0
		}
	}
}

func sub_4624D0(code int32) int32 {
	dr := uiInventoryUnlink(uint32(code))
	if dr == nil {
		return 0
	}
	found := uiInventoryFindCode(uint32(code))
	if found == nil {
		return int32(GetClient().Nox_xxx_spriteDelete_45A4B0(dr))
	}
	found.Cell.Equipped = 0
	alt := uiInventoryCellRef(uint32(dword_5d4594_1062480))
	if uint32(dword_5d4594_1062492) != uiInventoryPointer(dr.C()) {
		if GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(int(dr.TypeIDVal))&0xC != 0 && alt != nil && GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(int(alt.Drawable.TypeIDVal)) == 2 {
			alt.Alternate = 0
			uiInventorySetAlternate(nil)
		}
	} else {
		dword_5d4594_1062492 = 0
		if alt != nil {
			dword_5d4594_1062496 = uint32(dr.NetCode32)
			alt.Drawable.NetCode32 = alt.Codes[0]
			uiInventoryEquipRequest(alt.Drawable)
		} else {
			uiInventorySetAlternate(found.Cell)
			found.Cell.Alternate = 1
		}
	}
	return int32(GetClient().Nox_xxx_spriteDelete_45A4B0(dr))
}
