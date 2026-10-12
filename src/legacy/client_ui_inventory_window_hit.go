package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func uiInventoryPackedPoint(p uintptr) image.Point {
	return image.Pt(int(uint16(p)), int(uint16(p>>16)))
}
func uiInventoryEquipmentAt(pos image.Point) int {
	pos = pos.Sub(image.Pt(11, 15))
	for i := 0; i < 9; i++ {
		if !uiInventoryHitRect(pos, 136192+uintptr(i*16)) {
			continue
		}
		if i == 6 {
			for dr := uiInventoryDrawable(uiInventoryEquipment()[8]); dr != nil; dr = uiInventoryNext(dr) {
				if uint32(dr.Class())&0x2000000 != 0 && uint32(dr.SubClass())&2 != 0 {
					return 8
				}
			}
			return 5
		}
		if i != 0 || uiInventoryEquipment()[0] != 0 {
			return i
		}
	}
	return -1
}
func uiInventoryButtonTooltip() int {
	key := "OpenInventoryTT"
	if memmap.Uint8(0x5D4594, 1049868) == 2 {
		key = "CloseInventoryTT"
	}
	uiInventoryTooltipKey(key)
	return 1
}

func uiInventoryAlternateTooltip() int {
	if cell := uiInventoryAlternate(); cell != nil {
		uiCursorTooltip(uiItemTooltip(cell.Drawable))
	} else {
		uiInventoryTooltipKey("ToolTipWeapon2Area")
	}
	return 1
}

func uiInventoryHover(pos image.Point) int { uiCursorTooltip(uiInventoryHoverText(pos)); return 1 }

func uiInventoryStatusTooltip(pos image.Point) int {
	index := 0
	for x := 40; x <= pos.X; x += 35 {
		index++
	}
	mask := memmap.Uint32(0x5D4594, 1062540)
	for i := 0; i < 32; i++ {
		if mask&(uint32(1)<<i) != 0 && i != 31 {
			index--
		}
		if index < 0 {
			spell := nox_xxx_getEnchantSpell_424920(i)
			uiCursorTooltip((*uint16)(unsafe.Pointer(nox_xxx_spellTitle_424930(int(spell)))))
			return 1
		}
	}
	extra := memmap.Uint8(0x5D4594, 1062536)
	for i := 0; i < 6; i++ {
		if extra&(1<<i) != 0 {
			index--
		}
		if index < 0 {
			uiCursorTooltip(runtimeModifierLabel(byte(1 << i)))
			return 1
		}
	}
	if noxflags.HasGame(4096) {
		hit := func(off uintptr) bool {
			r := unsafe.Slice((*int32)(memmap.PtrOff(0x5D4594, off)), 4)
			return pos.X >= int(r[0]) && pos.X <= int(r[2]) && pos.Y >= int(r[1]) && pos.Y <= int(r[3])
		}
		if hit(1049812) {
			uiInventoryTooltipKey("thing.db:AnkhGUI")
			return 1
		}
		if hit(1049828) && sub_4BFD30() == 1 {
			uiInventoryTooltipKey("GeneralPrint:TooltipKeyIcon")
			return 1
		}
	}
	uiCursorTooltip(nil)
	return 1
}

func uiInventoryTooltipKey(key string) {
	uiCursorTooltip((*uint16)(unsafe.Pointer(internWStr(uiInventoryText(key)))))
}
