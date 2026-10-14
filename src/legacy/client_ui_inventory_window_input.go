package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func uiInventorySourceCell() *uiInventoryCell {
	return &uiInventoryGrid()[int(dword_5d4594_1049800_inventory_click_row_index)+21*int(dword_5d4594_1049796_inventory_click_column_index)]
}
func uiInventoryRestoreDrag() int {
	return uiInventoryPlace(uiInventoryDragged(), int(int32(dword_5d4594_1049796_inventory_click_column_index)), int(int32(dword_5d4594_1049800_inventory_click_row_index)))
}
func uiInventoryFinishDrag() {
	Nox_xxx_cursorSetDraggedItem_477690(nil)
	if dword_5d4594_1049856 == 0 {
		GetClient().Nox_xxx_spriteDelete_45A4B0(uiInventoryDragged())
	}
	uiInventorySetDragged(nil)
	dword_5d4594_1049856 = 0
}
func uiInventoryTradeClick(pos image.Point) int {
	if !uiInventoryHitRect(pos, 136352) {
		return 0
	}
	col, row := (pos.X-314)/50, (pos.Y+int(int32(dword_5d4594_1062512))-13)/50
	cell := &uiInventoryGrid()[row+21*col]
	if cell.Count == 0 {
		return 0
	}
	return uiInventoryTrade(28, uint16(cell.Codes[cell.Count-1]))
}
func uiInventoryStartDrag(w *gui.Window, pos image.Point) {
	if memmap.Uint8(0x5D4594, 1049868) != 2 {
		return
	}
	if uiInventoryHitRect(pos, 136336) {
		// no slot (-1) is no item: C read the word before the array here
		var code uint32
		if slot := uiInventoryEquipmentAt(pos); slot >= 0 {
			code = uiInventoryEquipment()[slot]
		}
		uiInventorySetDragged(uiInventoryDrawable(code))
		dword_5d4594_1049856 = 1
		return
	}
	dword_5d4594_1049856 = 0
	if uiInventoryHitRect(pos, 136368) {
		if (pos.Y-13)/50 == 2 {
			Nox_client_toggleMap_473610()
		}
	} else if uiInventoryHitRect(pos, 136352) {
		col, row := (pos.X-314)/50, (pos.Y+int(int32(dword_5d4594_1062512))-13)/50
		uiInventorySetClick(col, row)
		if uiInventoryValidCell(col, row) {
			uiInventoryDragCopy()
		}
	}
}
func uiInventoryDropAt(pos image.Point) int {
	p := [2]int32{int32(pos.X), int32(pos.Y)}
	return int(nox_xxx_clientDrop_465BE0(&p))
}
func uiInventoryDropQuantity(pos image.Point, code, typ uint32, count int) {
	if count == 0 {
		return
	}
	world := Sub_473970(pos)
	if found := uiInventoryFindCode(code); found != nil {
		for i := 0; i < count; i++ {
			dr := found.Cell.Drawable
			uiInventorySetDragged(dr)
			dr.NetCode32 = found.Cell.Codes[i]
			if uiTradeActive() == 0 {
				uiInventoryDropAt(world)
			}
			uiInventorySetDragged(nil)
		}
	}
}

func uiInventoryMainEvents(w *gui.Window, event int, a, b uintptr) int {
	screen := uiInventoryPackedPoint(a)
	pos := screen.Sub(uiWindowPosition(uiInventoryMainWindow()))
	if *bookWord(1047520) != 0 || memmap.Uint8(0x5D4594, 1049868) != 2 {
		return 1
	}
	hit := func(off uintptr) bool { return uiInventoryHitRect(pos, off) }
	switch event {
	case 5:
		if Nox_xxx_playerAnimCheck_4372B0() != 0 {
			return 1
		}
		switch uiInventoryMode() {
		case 5:
			if hit(136352) {
				col, row := (pos.X-314)/50, (pos.Y+int(int32(dword_5d4594_1062512))-13)/50
				if !uiInventoryValidCell(col, row) {
					return 1
				}
				cell := &uiInventoryGrid()[row+21*col]
				if cell.Count != 0 {
					dword_5d4594_1063116 = uint32(uiInventoryPointer(cell.Drawable.C()))
					cell.Drawable.NetCode32 = cell.Codes[0]
				} else {
					dword_5d4594_1063116 = 0
				}
				return 1
			}
			if uiShopActive() != 0 && uiShopMode() == 2 {
				if uiShopInside(pos) {
					dword_5d4594_1063116 = uint32(uiInventoryPointer(uiShopHit(pos).C()))
					return 1
				}
			}
		case 6:
			if hit(136352) {
				col, row := (pos.X-314)/50, (pos.Y+int(int32(dword_5d4594_1062512))-13)/50
				if uiInventoryValidCell(col, row) {
					cell := &uiInventoryGrid()[row+21*col]
					if cell.Count != 0 {
						cell.Drawable.NetCode32 = cell.Codes[0]
						if cell.Drawable != nil {
							uiInventoryTrade(30, uint16(cell.Codes[0]))
							return 1
						}
					}
				}
			}
		default:
			if memmap.Uint8(0x5D4594, 1049870) != 1 || !hit(136336) {
				if memmap.Uint8(0x5D4594, 1049869) != 0 || hit(136384) || hit(136400) {
					if hit(136384) || hit(136400) {
						return 0
					}
				} else {
					uiInventoryMainWindow().Capture(true)
					if Sub_479590() == 3 {
						uiInventoryTradeClick(pos)
					} else {
						uiInventoryStartDrag(uiInventoryMainWindow(), pos)
					}
					if dr := uiInventoryDragged(); dr != nil {
						Nox_xxx_cursorSetDraggedItem_477690(dr)
						InputSetKeyTimeoutLegacy(0)
						*(*[2]int32)(memmap.PtrOff(0x5D4594, 1062572)) = [2]int32{int32(pos.X), int32(pos.Y)}
						audioEventPlay(791, 100, 0, 0)
						return 1
					}
				}
			}
		}
		return 1
	case 7:
		if Nox_xxx_playerAnimCheck_4372B0() != 0 || uiInventoryMode() == 6 {
			return 1
		}
		if memmap.Uint8(0x5D4594, 1049869) == 0 && hit(136368) && (pos.Y-13)/50 == 1 {
			if uiInventoryMode() != 5 {
				uiInventoryOpenIdentify()
			} else {
				uiInventoryCloseIdentify()
			}
			return 1
		}
		fallthrough
	case 6:
		if Nox_xxx_playerAnimCheck_4372B0() != 0 || uiInventoryMode() == 6 {
			return 1
		}
		if uiInventoryMode() == 5 {
			if GetClient().Cli().CursorPrev == 7 {
				uiInventoryCloseIdentify()
				return 1
			}
		} else {
			uiInventoryMainWindow().Capture(false)
		}
		if uiInventoryMode() == 4 {
			world := Sub_473970(screen)
			dr := GetClient().Nox_drawable_find(world, 20)
			if dr != nil {
				center := Sub_473970(image.Pt(int(nox_win_width)/2, int(nox_win_height)/2))
				d := center.Sub(dr.Pos())
				if d.X*d.X+d.Y*d.Y <= 5625 {
					dword_5d4594_1049864 = 0
					return 1
				}
				uiInventoryError("ObjectTooFar")
			} else {
				uiInventoryError("NoObject")
			}
			dword_5d4594_1049864 = 0
			return 1
		}
		dr := uiInventoryDragged()
		if dr == nil {
			return 1
		}
		defer uiInventoryFinishDrag()
		if !uiWindowPointIn(uiInventoryMainWindow(), int32(screen.X), int32(screen.Y)) || hit(136384) || hit(136400) {
			world := Sub_473970(screen)
			if dword_5d4594_1049856 == 1 {
				if uiTradeActive() == 0 {
					uiInventoryDropAt(world)
				}
			} else {
				count := int(uiInventorySourceCell().Count)
				if count != 0 {
					uiInventoryMainWindow().Capture(false)
					var mods unsafe.Pointer
					if uint32(dr.Class())&0x13001000 != 0 {
						mods = unsafe.Pointer(&dr.Union)
					}
					uiAmountPrice(0, 0)
					uiAmountShow((*uint16)(unsafe.Pointer(internWStr(uiInventoryText("DropLabel")))), screen.X, screen.Y, dr.NetCode32, dr.TypeIDVal, mods, uint32(count+1), 0, uiAmountNativeKey(uiAmountDrop), nil)
				} else if uiTradeActive() == 0 {
					uiInventoryDropAt(world)
				}
			}
			if dword_5d4594_1049856 == 0 {
				uiInventoryRestoreDrag()
			}
			return 1
		}
		delta := image.Pt(int(memmap.Int32(0x5D4594, 1062572)), int(memmap.Int32(0x5D4594, 1062576))).Sub(pos)
		if !InputKeyCheckTimeoutLegacy(0, uint32(GetServer().S().TickRate())/3) && delta.X*delta.X+delta.Y*delta.Y < 100 {
			if !hit(136352) {
				return 1
			}
			if uiTradeActive() == 0 {
				if uint32(dr.Class())&0x3001000 != 0 {
					cell := uiInventorySourceCell()
					if cell.Alternate != 0 {
						uiInventorySetAlternate(nil)
						cell.Alternate = 0
					} else if cell.Equipped != 0 {
						uiInventoryDequipRequest(dr)
					} else {
						nox_xxx_clientKeyEquip_465C30(int32(dword_5d4594_1049796_inventory_click_column_index), int32(dword_5d4594_1049800_inventory_click_row_index))
					}
				} else {
					uiInventoryUse(dr)
				}
			}
			uiInventoryRestoreDrag()
			return 1
		}
		if hit(136336) && memmap.Uint8(0x5D4594, 1049870) == 0 {
			if dword_5d4594_1049856 == 0 {
				uiInventoryEquipRequest(dr)
				uiInventoryRestoreDrag()
			}
			return 1
		}
		if !hit(136352) {
			uiInventoryRestoreDrag()
			return 1
		}
		typ := dr.TypeIDVal
		if typ == uint32(dword_5d4594_1062560) || typ == memmap.Uint32(0x5D4594, 1049728) || typ == memmap.Uint32(0x5D4594, 1049724) || typ == uint32(dword_5d4594_1062556) || typ == uint32(dword_5d4594_1062564) {
			uiInventoryRestoreDrag()
			return 1
		}
		col, row := (pos.X-314)/50, (pos.Y+int(int32(dword_5d4594_1062512))-13)/50
		dword_5d4594_1049804, dword_5d4594_1049808 = uint32(col), uint32(row)
		if !uiInventoryValidCell(col, row) {
			return 1
		}
		target := &uiInventoryGrid()[row+21*col]
		if dword_5d4594_1049856 != 0 {
			other := target.Drawable
			compatible := false
			if target.Count != 0 && other != nil {
				a, b := uint32(other.Class()), uint32(dr.Class())
				compatible = a&0x2000000 != 0 && b&0x2000000 != 0 && other.SubClass() == dr.SubClass() || a&0x1001000 != 0 && b&0x1001000 != 0
			}
			*memmap.PtrUint32(0x5D4594, 1049860) = 1
			if compatible {
				other.NetCode32 = target.Codes[0]
				uiInventoryEquipRequest(other)
			} else {
				uiInventoryDequipRequest(dr)
			}
			return 1
		}
		if uiInventorySourceCell().Count != 0 {
			uiInventoryRestoreDrag()
			return 1
		}
		if uiInventoryPlace(dr, col, row) == 0 {
			uiInventoryRestoreDrag()
			return 1
		}
		audioEventPlay(792, 100, 0, 0)
		source := uiInventorySourceCell()
		if source.Alternate != 0 {
			target.Alternate = source.Alternate
			source.Alternate = 0
			dword_5d4594_1062480 = uint32(uiInventoryPointer(unsafe.Pointer(target)))
		}
		uiInventoryCompact()
		return 1
	case 8:
		return 1
	case 9:
		if uiInventoryMode() == 5 {
			uiInventoryCloseIdentify()
			return 1
		}
		return 0
	case 19, 20:
		if Nox_xxx_playerAnimCheck_4372B0() != 0 {
			return 1
		}
		if hit(136384) || hit(136400) {
			return bool2int(uiInventoryMode() == 5)
		}
		off := uintptr(1062500)
		if event == 20 {
			off += 4
		}
		uiInventoryMainWindow().Func94(&gui.RawEvent{Event: 16391, Arg1: uintptr(memmap.LoadPtr(0x5D4594, off))}) // the up/down button window
		return 1
	default:
		return bool2int(uiInventoryMode() == 5)
	}
}

func uiInventoryAlternateEvents(w *gui.Window, event int, a, b uintptr) int {
	if uiInventoryMode() == 6 {
		return 1
	}
	switch event {
	case 5, 8:
		return 1
	case 6:
		dr := uiInventoryDragged()
		if dr != nil {
			pos := uiInventoryPackedPoint(a)
			inside := uiWindowPointIn(uiInventoryMainWindow(), int32(pos.X), int32(pos.Y))
			if inside {
				if dword_5d4594_1049856 != 0 {
					if uint32(dr.Class())&0x1001000 != 0 {
						if uiInventoryAlternate() != nil {
							uiInventoryAlterWeapon()
						} else {
							dword_5d4594_1062492 = uint32(uiInventoryPointer(dr.C()))
							uiInventoryDequipRequest(dr)
						}
					}
				} else if uint32(dr.Class())&0x1001000 == 0 || uiInventorySourceCell().Equipped != 0 {
					uiInventoryRestoreDrag()
				} else {
					if nox_xxx_ammoCheck_415880(int(dr.TypeIDVal)) == 2 {
						bow, crossbow := sub_461600(sub_415840(4)), sub_461600(sub_415840(8))
						if bow == 0 && crossbow == 0 {
							uiInventoryRestoreDrag()
							uiInventoryFinishDrag()
							return 1
						}
					}
					if alt := uiInventoryAlternate(); alt != nil {
						alt.Alternate = 0
					}
					uiInventoryRestoreDrag()
					uiInventorySetAlternate(uiInventorySourceCell())
					uiInventoryAlternate().Alternate = 1
				}
			}
			uiInventoryFinishDrag()
		}
		uiInventoryMainWindow().Capture(false)
		return 1
	case 7:
		if uiInventoryAlternate() != nil {
			uiInventoryAlterWeapon()
		}
		return 1
	default:
		return 0
	}
}
