package legacy

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func uiInventoryCellRef(v uint32) *uiInventoryCell {
	return (*uiInventoryCell)(unsafe.Pointer(uintptr(v)))
}
func uiInventoryPointer(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }
func uiInventorySetNext(dr, next *client.Drawable) {
	*(**client.Drawable)(unsafe.Add(dr.C(), 368)) = next
}
func uiInventoryPrev(dr *client.Drawable) *client.Drawable {
	return *(**client.Drawable)(unsafe.Add(dr.C(), 372))
}
func uiInventorySetPrev(dr, prev *client.Drawable) {
	*(**client.Drawable)(unsafe.Add(dr.C(), 372)) = prev
}
func uiInventoryError(id string) {
	s := GetServer().S().Strings().GetStringInFile(strman.ID(id), `C:\NoxPost\src\Client\Gui\guiinv.c`)
	Nox_xxx_printCentered_445490(s)
}
func uiInventoryCopyItem(dst, src *client.Drawable) {
	copy(unsafe.Slice((*byte)(unsafe.Add(dst.C(), 432)), 24), unsafe.Slice((*byte)(unsafe.Add(src.C(), 432)), 24))
	*(*uint32)(unsafe.Add(dst.C(), 292)) = *(*uint32)(unsafe.Add(src.C(), 292))
}
func uiInventoryClearAlternateFlags() uintptr {
	grid := uiInventoryGrid()
	for i := range grid {
		if grid[i].Count != 0 {
			grid[i].Alternate = 0
		}
	}
	// Preserve the legacy derived return address; it is not dereferenced.
	return uintptr(unsafe.Pointer(&grid[0])) + 15532
}
func uiInventorySetAlternate(cell *uiInventoryCell) int {
	old := uiInventoryCellRef(uint32(dword_5d4594_1062480))
	dword_5d4594_1062484 = 0
	if old != nil {
		dword_5d4594_1062484 = uint32(old.Codes[0])
	}
	dword_5d4594_1062480 = uint32(uiInventoryPointer(unsafe.Pointer(cell)))
	uiInventoryClearAlternateFlags()
	if cell == nil {
		return int(nox_xxx_clientReportSecondaryWeapon_4BF010(int32(0)))
	}
	cell.Drawable.NetCode32 = cell.Codes[0]
	cell.Alternate = 1
	return int(nox_xxx_clientReportSecondaryWeapon_4BF010(int32(uiInventoryPointer(cell.Drawable.C()))))
}

func uiInventoryAppend(code, typ uint32) *uiInventoryCell {
	if uint32(GetClient().Cli().Things.TypeByInd(int(typ)).ObjClass)&0x4000000 != 0 {
		return nil
	}
	grid := uiInventoryGrid()
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			if cell.Count != 0 && cell.Count < 32 && cell.Drawable.TypeIDVal == typ {
				cell.Codes[cell.Count] = code
				cell.Count++
				return cell
			}
		}
	}
	return nil
}
func uiInventoryNewStack(code, typ uint32, mods unsafe.Pointer, coords *[2]int32) int {
	grid := uiInventoryGrid()
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			if cell.Count != 0 {
				continue
			}
			dr := GetClient().Nox_new_drawable_for_thing(int(typ))
			cell.Drawable = dr
			if dr == nil {
				uiInventoryError("DrawablesExhausted")
				return 0
			}
			dr.ObjFlags |= 0x40000000
			cell.Codes[0] = code
			cell.Count = 1
			if uint32(dr.ObjClass)&0x13001000 != 0 {
				copy(unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 20), unsafe.Slice((*byte)(mods), 20))
			}
			if coords != nil {
				*coords = [2]int32{int32(col), int32(row)}
			}
			if sub_461930() != 0 && dword_5d4594_1062480 == 0 {
				class, sub := uint32(dr.ObjClass), uint32(dr.ObjSubClass)
				if class&0x1000000 != 0 && sub&2 == 0 || class&0x1000 != 0 {
					def := GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
					// Match the 386 shift-count masking before the original byte truncation.
					if def == nil || byte(uint32(1)<<(*(*byte)(unsafe.Add(uiMeterPlayer(), 2251))&31))&*(*byte)(unsafe.Add(unsafe.Pointer(def), 62)) != 0 {
						uiInventorySetAlternate(cell)
						cell.Alternate = 1
					}
				}
			}
			return 1
		}
	}
	return 0
}

func nox_xxx_spritePickup_461660(code, typ int32, mods unsafe.Pointer) int32 {
	t := uint32(typ)
	if t == uint32(dword_5d4594_1062560) || t == memmap.Uint32(0x5D4594, 1049728) || t == memmap.Uint32(0x5D4594, 1049724) || t == uint32(dword_5d4594_1062556) || t == uint32(dword_5d4594_1062564) {
		return 1
	}
	var coords [2]int32
	cell := uiInventoryAppend(uint32(code), t)
	if cell != nil {
		index := (uintptr(unsafe.Pointer(cell)) - uintptr(unsafe.Pointer(&uiInventoryGrid()[0]))) / unsafe.Sizeof(*cell)
		coords = [2]int32{int32(index / 21), int32(index % 21)}
		if uint32(cell.Drawable.ObjClass)&0x10 != 0 {
			uiMeterRefreshPotions()
		}
	} else {
		if uiInventoryNewStack(uint32(code), t, mods, &coords) == 0 {
			uiInventoryError("InventoryFull")
			return 0
		}
		cell = &uiInventoryGrid()[coords[1]+21*coords[0]]
		if uint32(cell.Drawable.ObjClass)&0x10 != 0 {
			uiMeterRefreshPotions()
		}
		if uint32(cell.Drawable.ObjClass)&0x3001000 != 0 {
			dword_5d4594_1062516 = 0
			if coords[1] >= 3 {
				dword_5d4594_1062516 = uint32(10 * (5*coords[1] - 10))
			}
		}
	}
	if tt := GetClient().Cli().Things.TypeByInd(int(t)); tt != nil && uint32(tt.ObjClass)&0x1001000 != 0 {
		sub_4673F0(int(coords[0]), int(coords[1]))
	}
	return 1
}
func uiInventoryRemoveStack(found *uiInventoryLookup) uintptr {
	cell := found.Cell
	for i := int(found.Index); i < int(cell.Count)-1; i++ {
		cell.Codes[i] = cell.Codes[i+1]
	}
	cell.Count--
	if cell.Count == 0 {
		GetClient().Nox_xxx_spriteDelete_45A4B0(cell.Drawable)
		cell.Drawable = nil
	}
	if cell.Alternate != 0 {
		uiInventorySetAlternate(nil)
		cell.Alternate = 0
		return uintptr(unsafe.Pointer(cell))
	}
	return 0
}
func uiInventoryCompact() uintptr {
	grid := uiInventoryGrid()
	base := uintptr(unsafe.Pointer(&grid[0]))
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			dst := &grid[row+21*col]
			if dst.Count != 0 {
				continue
			}
		retry:
			sr, sc := -1, -1
			for r := row; r < 20 && sr < 0; r++ {
				c0 := 0
				if r == row {
					c0 = col
				}
				for c := c0; c < 4; c++ {
					if grid[r+21*c].Count != 0 {
						sr, sc = r, c
						break
					}
				}
			}
			if sr < 0 {
				for r := row; r < 20; r++ {
					c0 := 0
					if r == row {
						c0 = col
					}
					for c := c0; c < 4; c++ {
						cell := &grid[r+21*c]
						cell.Equipped = 0
						cell.Alternate = 0
					}
				}
				return base + 15380
			}
			src := &grid[sr+21*sc]
			typ := src.Drawable.TypeIDVal
			if uint32(src.Drawable.ObjClass)&0x4000000 == 0 {
				for r := 0; r < 20; r++ {
					for c := 0; c < 4; c++ {
						other := &grid[r+21*c]
						if other.Count == 0 || other.Count == 32 || other.Drawable.TypeIDVal != typ || other == src {
							continue
						}
						for src.Count > 0 && other.Count < 32 {
							src.Count--
							other.Codes[other.Count] = src.Codes[src.Count]
							other.Count++
						}
						if src.Count == 0 {
							GetClient().Nox_xxx_spriteDelete_45A4B0(src.Drawable)
							src.Drawable = nil
							goto retry
						}
						break
					}
				}
			}
			*dst = *src
			if dst.Alternate != 0 {
				dword_5d4594_1062480 = uint32(uiInventoryPointer(unsafe.Pointer(dst)))
			}
			src.Count = 0
			src.Drawable = nil
			src.Equipped = 0
		}
	}
	return base + 3096
}

func sub_461A80(code int32) {
	if found := uiInventoryFindCode(uint32(code)); found != nil {
		potion := uint32(found.Cell.Drawable.ObjClass)&0x10 != 0
		uiInventoryRemoveStack(found)
		found.Cell.Equipped = 0
		uiInventoryCompact()
		if dr := uiInventoryUnlink(uint32(code)); dr != nil {
			GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		}
		if potion {
			uiMeterRefreshPotions()
		}
	} else if dr := uiInventoryDragged(); dr != nil && dr.NetCode32 == uint32(code) {
		if uint32(dr.ObjClass)&0x10 != 0 {
			uiMeterRefreshPotions()
		}
		GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		uiInventoryDraggedItem = nil
		dword_5d4594_1049856 = 0
		Nox_xxx_cursorResetDraggedItem_4776A0()
	} else {
		uiInventoryError("DroppedNotFound")
	}
}
