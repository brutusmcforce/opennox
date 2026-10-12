package legacy

import (
	"github.com/brutusmcforce/opennox/libs/client/keybind"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"strconv"
	"unsafe"
)

const meterStringSource = `C:\NoxPost\src\Client\Gui\guimeter.c`

func uiMeterString(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), meterStringSource)
}

type uiPotionSlot struct {
	Image   *client.Drawable
	Binding [4]uint16
	Type    uint32
	Count   uint16
	padding uint16
}

func uiMeterSlot(index int) *uiPotionSlot {
	return &uiPotionSlots[index].Slot
}
func uiMeterPotionDrawable(index int) *client.Drawable {
	return &uiPotionSlots[index].Drawable
}
func uiMeterLinkPotion(index int, typeID uint32) bool {
	slot := uiMeterSlot(index)
	typ := GetClient().Cli().Things.TypeByInd(int(typeID))
	if typ == nil {
		slot.Image = nil
		return false
	}
	dr := uiMeterPotionDrawable(index)
	GetClient().Cli().DrawableLinkThing(dr, typ.Index())
	slot.Image = dr
	dr.ObjFlags |= 0x40000000
	return true
}
func uiMeterPotionDraw(w *gui.Window) int {
	s := uiMeterSlot(int(uintptr(w.WidgetData)))
	pos := uiWindowPosition(w)
	r := GetClient().R2()
	font := r.GetFonts().AsFont(memmap.WordPtr(dword_5d4594_1096288))
	if s.Count != 0 {
		if s.Image != nil {
			s.Image.PosVec = pos.Add(image.Pt(14, 15))
			s.Image.CallDraw(&uiPotionViewport)
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
		r.DrawString(font, strconv.Itoa(int(s.Count)), image.Pt(pos.X-2, pos.Y-r.FontHeight(font)+10))
	}
	r.DrawString(font, alloc.GoString16(&s.Binding[0]), image.Pt(pos.X-2, pos.Y-r.FontHeight(font)+33))
	return 1
}

func uiMeterInputResult(event int) int {
	if event == 8 || event == 12 || event == 16 {
		return 0
	}
	return 1
}

func sub_470E90(window, event int) int {
	if event == 5 {
		uiInventoryAlterWeapon()
	}
	return uiMeterInputResult(event)
}

func nox_xxx_guiBottleSlotProc_471B90(window, event int) int {
	if event == 5 {
		w := (*gui.Window)(memmap.WordPtr(uint32(window)))
		if typ := uiMeterSlot(int(uintptr(w.WidgetData))).Type; typ != 0 {
			uiInventoryUsePotion(typ)
		}
	}
	return uiMeterInputResult(event)
}

func nox_xxx_guiHealthManaTubeProc_472100(window, event int) int {
	if event == 7 {
		old := uint32(dword_5d4594_1096252)
		dword_5d4594_1096252 = uint32(1 - old)
		uiMeterHide(uiMeters()[2].Window, old == 1)
		if memmap.Uint8(0x85B3FC, 12254) != 0 {
			uiMeterHide(uiMeters()[3].Window, dword_5d4594_1096252 == 0)
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(901, 100)
	}
	return uiMeterInputResult(event)
}
func uiMeterQuickPotion(index int) {
	if memmap.Uint32(0x5D4594, 1096672) == 0 {
		if typ := uiMeterSlot(index).Type; typ != 0 {
			uiInventoryUsePotion(typ)
		}
	}
}
func uiMeterBindings() unsafe.Pointer {
	player := uiMeterPlayer()
	if player == nil {
		return nil
	}
	fill := func(index int, event keybind.Event) {
		text := GetClient().GetCtrlEvent().Sub_42E8E0_go(event, 1)
		// The key label is deliberately limited to three UTF16 units.
		slot := uiMeterSlot(index)
		clear(slot.Binding[:])
		alloc.StrCopyZero16(slot.Binding[:], text)
	}
	fill(2, 38)
	fill(0, 36)
	if *(*byte)(unsafe.Add(player, 2251)) != 0 {
		fill(1, 37)
		return unsafe.Pointer(&uiMeterSlot(1).Binding[0])
	}
	return player
}

func uiMeterRefreshPotions() uintptr {
	count := func(typ uint32) uint16 { return uint16(uiInventoryTypeCount(typ)) }
	uiMeterSlot(2).Count = count(uint32(dword_5d4594_1096276))
	uiMeterSlot(1).Count = count(uint32(dword_5d4594_1096272))
	uiMeterSlot(2).Count = count(uint32(dword_5d4594_1096276))
	slot := uiMeterSlot(0)
	red := memmap.Uint32(0x5D4594, 1096268)
	slot.Count = count(red)
	if slot.Count != 0 {
		linked := uiMeterLinkPotion(0, red)
		slot.Type = red
		if linked {
			return uintptr(uint32(slot.Image.ObjFlags))
		}
		return 0
	}
	meat := uint32(dword_5d4594_1096284)
	slot.Count = count(meat)
	if slot.Count != 0 {
		uiMeterLinkPotion(0, meat)
		slot.Type = meat
		return uintptr(meat)
	}
	apple := uint32(dword_5d4594_1096280)
	slot.Count = count(apple)
	if slot.Count != 0 {
		linked := uiMeterLinkPotion(0, apple)
		slot.Type = apple
		if linked {
			return uintptr(uint32(slot.Image.ObjFlags))
		}
		return 0
	}
	slot.Image = nil
	slot.Type = 0
	return 0
}

func uiMeterWeaponTooltip() int {
	dr := uiInventoryCurrentWeapon()
	if dr == nil {
		uiCursorTooltip(alloc.InternCString16(uiMeterString("ToolTipCurWeapon")))
		return 1
	}
	// The next live mapped field is the poison color at1092992. Leave a
	// terminator within this512-unit buffer even for oversized localized names.
	dst := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1091968)), 512)
	alloc.StrCopyZero16P(dst, uiItemTooltip(dr))
	if uint32(dr.ObjSubClass)&12 != 0 {
		typ := memmap.PtrUint32(0x5D4594, 1096292)
		if *typ == 0 {
			*typ = uint32(GetServer().S().Types.IndByID("Quiver"))
		}
		q := uiInventoryEquippedType(*typ)
		if q != nil {
			appendText := func(p *uint16) { n := alloc.StrLenS(dst); alloc.StrCopyZero16P(dst[n:], p) }
			appendText(alloc.InternCString16("\n"))
			appendText(uiItemTooltip(q))
		}
	}
	uiCursorTooltip(&dst[0])
	return 1
}
