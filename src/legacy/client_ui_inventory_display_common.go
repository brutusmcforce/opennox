package legacy

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"golang.org/x/image/font"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

func uiInventoryText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), `C:\NoxPost\src\Client\Gui\guiinv.c`)
}
func uiInventoryFormatInts(id string, a, b int) string {
	var dst [256]uint16
	textFormatBuffer(dst[:], (*uint16)(unsafe.Pointer(internWStr(uiInventoryText(id)))), textFormatWord(uint32(a)), textFormatWord(uint32(b)))
	return alloc.GoString16(&dst[0])
}
func uiInventoryFormatFloat(id string, a float64) string {
	var dst [256]uint16
	textFormatBuffer(dst[:], (*uint16)(unsafe.Pointer(internWStr(uiInventoryText(id)))), textFormatReal(a))
	return alloc.GoString16(&dst[0])
}
func uiInventoryFormatLiteral(text string) string {
	var dst [256]uint16
	textFormatBuffer(dst[:], (*uint16)(unsafe.Pointer(internWStr(text))))
	return alloc.GoString16(&dst[0])
}
func uiInventorySmallFont() font.Face {
	return GetClient().R2().GetFonts().AsFont(memmap.WordPtr(dword_5d4594_1063636))
}
func uiInventoryMainWindow() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(dword_5d4594_1062456))
}
func uiInventoryIdentifyWindow() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(dword_5d4594_1062476))
}
func uiInventoryAlternate() *uiInventoryCell {
	return (*uiInventoryCell)(memmap.WordPtr(dword_5d4594_1062480))
}
func uiInventorySelectedItem() *client.Drawable {
	return uiInventoryDrawable(uint32(dword_5d4594_1063116))
}

// drawableModifier returns modifier slot i (0-3) of an item drawable.
func drawableModifier(dr *client.Drawable, i int) **server.ModifierEff {
	return &layout.As[client.DrawableUnionItem](unsafe.Pointer(&dr.Union)).Modifiers[i]
}

// drawableSetModifiers sets the modifiers of an item drawable; Field16 goes to Field_112.
func drawableSetModifiers(dr *client.Drawable, m *server.ModifierInitData) {
	d := layout.As[client.DrawableUnionItem](unsafe.Pointer(&dr.Union))
	d.Modifiers = m.Modifiers
	*(*uint32)(unsafe.Pointer(&d.Field_112_0)) = m.Field16
}

// drawableModifiers returns the modifiers of an item drawable (see drawableSetModifiers).
func drawableModifiers(dr *client.Drawable) *server.ModifierInitData {
	m := new(server.ModifierInitData)
	d := layout.As[client.DrawableUnionItem](unsafe.Pointer(&dr.Union))
	m.Modifiers = d.Modifiers
	m.Field16 = *(*uint32)(unsafe.Pointer(&d.Field_112_0))
	return m
}

func uiInventoryItemModifier(dr *client.Drawable, index int) *server.ModifierEff {
	return *drawableModifier(dr, index)
}
func uiInventoryElementValue(dr *client.Drawable, fire bool) float64 {
	if dr == nil || uint32(dr.Class())&0x13001000 == 0 {
		return 0
	}
	var fn unsafe.Pointer = modifierKey(modifierIDLightningEffect)
	if fire {
		fn = modifierKey(modifierIDFireEffect)
	}
	for i := 2; i < 4; i++ {
		if m := uiInventoryItemModifier(dr, i); m != nil && m.AttackPreHit52.Fnc == fn {
			return float64(m.AttackPreHit52.Valf)
		}
	}
	return 0
}
func uiInventoryScaledDurability(dr *client.Drawable, current, maximum *float32) uint32 {
	*current = float32(dr.Field_73_1)
	*maximum = float32(dr.Field_73_2)
	ret := uiInventoryPointer(dr.C())
	if uint32(dr.Class())&0x13001000 != 0 {
		if m := uiInventoryItemModifier(dr, 1); m != nil && m.Defend76.Fnc == modifierKey(modifierIDDurabilityMultiplierEffect) {
			*current = float32(int32(*current * m.Defend76.Valf))
			*maximum = float32(int32(*maximum * m.Defend76.Valf))
			ret = uint32(int32(*maximum))
		}
	}
	return ret
}

func sub_463420(v int32) int32 { *memmap.PtrUint32(0x5D4594, 1050012) = uint32(v); return v }

func uiInventoryHitRect(p image.Point, offset uintptr) bool {
	r := unsafe.Slice((*int32)(memmap.PtrOff(0x587000, offset)), 4)
	return p.X >= int(r[0]) && p.X <= int(r[2]) && p.Y >= int(r[1]) && p.Y <= int(r[3])
}
func uiInventoryHoverText(p image.Point) *uint16 {
	if uiInventoryHitRect(p, 136336) {
		if memmap.Uint8(0x5D4594, 1049870) == 1 {
			return nil
		}
		slot := uiInventoryEquipmentAt(p)
		if slot == -1 {
			return alloc.InternCString16(uiInventoryText("DollRegionError"))
		}
		if dr := uiInventoryDrawable(uiInventoryEquipment()[slot]); dr != nil {
			return uiItemTooltip(dr)
		}
		return alloc.InternCString16(uiInventoryText("ToolTipDrag"))
	}
	if memmap.Uint8(0x5D4594, 1049869) != 0 {
		return nil
	}
	var col, row int
	if uiInventoryHitRect(p, 136368) {
		col, row = (p.Y-13)/50, 20
	} else if uiInventoryHitRect(p, 136352) {
		col = (p.X - 314) / 50
		row = (p.Y + int(dword_5d4594_1062512) - 13) / 50
	} else {
		return nil
	}
	dword_5d4594_1049796_inventory_click_column_index = uint32(col)
	dword_5d4594_1049800_inventory_click_row_index = uint32(row)
	if col >= 0 && col < 4 && row >= 0 && row < 21 {
		cell := &uiInventoryGrid()[col*21+row]
		if cell.Count != 0 {
			cell.Drawable.NetCode32 = cell.Codes[0]
			return uiItemTooltip(cell.Drawable)
		}
	}
	return nil
}

func sub_466E20(w *uint32) int32 {
	var key string
	switch uint32(*w) {
	case 9105:
		key = "JournalModeTT"
	case 9106:
		key = "InventoryModeTT"
	case 9107:
		key = "StatsModeTT"
	case 9108:
		key = "PaperDollModeTT"
	case 9111:
		key = "CloseInventoryTT"
	default:
		return 0
	}
	uiCursorTooltip(alloc.InternCString16(uiInventoryText(key)))
	return 1
}

func nox_xxx_inventoryNameSignInit_4671E0() int {
	dst := (*uint16)(memmap.PtrOff(0x5D4594, 1062588))
	textCopy((*uint16)(unsafe.Pointer(dst)), (*uint16)(memmap.PtrOff(0x5D4594, 1063676)))
	p := uiMeterPlayer()
	level := 0
	if noxflags.HasGame(4096) || questRuntimeWord(1556160) != 0 || questRuntimeWord(1556164) != 0 {
		level = int(min(uint32(dword_5d4594_1049844), 10))
	} else if p != nil {
		level = int(*(*int8)(unsafe.Pointer(&layout.As[server.Player](unsafe.Pointer(p)).Level)))
	}
	if p == nil {
		return level
	}
	class := layout.As[server.Player](unsafe.Pointer(p)).InfoBytes[66]
	className := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 29456+uintptr(class)*4)))
	title := uiInventoryText(fmt.Sprintf("experience:%s%d", className, level))
	return textFormatBuffer(unsafe.Slice(dst, 256), (*uint16)(unsafe.Pointer(internWStr(uiInventoryText("ElaborateNameFormat")))), textFormatPointer(unsafe.Pointer(&layout.As[server.Player](unsafe.Pointer(p)).NameFinal)), textFormatPointer(unsafe.Pointer(internWStr(title))))
}

func sub_467750(code int32, status int8) int32 {
	if code != 0 {
		if found := uiInventoryFindCode(uint32(code)); found != nil {
			if old := uiInventoryAlternate(); old != nil {
				old.Alternate = 0
			}
			dword_5d4594_1062480 = uint32(uiInventoryPointer(unsafe.Pointer(found.Cell)))
			found.Cell.Alternate = 1
			return 1
		}
	} else if old := uiInventoryAlternate(); old != nil {
		old.Alternate = 0
		dword_5d4594_1062480 = 0
	}
	if status != 0 {
		if status != 1 {
			return 0
		}
		Nox_xxx_printCentered_445490(uiInventoryText("Weapon2CantUse"))
		if dword_5d4594_1062484 == 0 {
			return 0
		}
		if found := uiInventoryFindCode(uint32(dword_5d4594_1062484)); found != nil {
			uiInventorySetAlternate(found.Cell)
			return 0
		}
	}
	dword_5d4594_1062484 = 0
	return 0
}
