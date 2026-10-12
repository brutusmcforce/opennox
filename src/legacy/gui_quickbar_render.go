package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func quickbarWhite() {
	GetClient().R2().Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
}

func quickbarKeyLabel(w *gui.Window, _ *gui.WindowData) int {
	text := GetClient().GetCtrlEvent().Sub_42E8E0_go(keybind.Event(*quickbarUserData(w)+28), 1)
	// The live C adapter returns an interned pointer even for an empty label.
	quickbarWhite()
	r := GetClient().R2()
	r.DrawString(r.GetFonts().AsFont(memmap.WordPtr(*quickbarWord(1049684))), text, w.GlobalPos())
	return 1
}

func quickbarRowLabel(w *gui.Window, trap bool) int {
	if !trap && *quickbarWord(1049476) != 0 {
		return 1
	}
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	pos := w.GlobalPos()
	name, row := "SpellSet", int(quickbarMain().Selected)+1
	if trap {
		name, row = "TrapSet", int(quickbarRecords[0].Selected)+1
	} else {
		pos.Y += -r.FontHeight(font) + r.FontHeight(r.GetFonts().AsFont(memmap.WordPtr(*quickbarWord(1049684)))) + 1
	}
	text := bookFormatInt(quickbarText(name), row)
	size := r.GetStringSizeWrapped(font, text, 0)
	pos.X += (w.SizeVal.X - size.X) / 2
	quickbarWhite()
	// the shadow of DrawStringHL (nox_xxx_drawSetColor_4343E0 in the original)
	r.Data().SetColor(noxcolor.RGBA5551(*memmap.PtrUint32(0x852978, 4)))
	r.DrawStringHL(font, text, pos)
	return 1
}

func quickbarHighlight(slot int) bool {
	frame := GetServer().S().Frame()
	return frame > 10 && uint32(slot) == *memmap.PtrUint32(0x587000, 133484) && frame-*quickbarWord(1049540) < 10
}

func quickbarNoIcon(pos image.Point) {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	quickbarWhite()
	r.DrawString(font, quickbarText("NoIcon"), pos.Add(image.Pt(6, r.FontHeight(font)+2)))
}

func quickbarDrawAbility(w *gui.Window, _ *gui.WindowData) int {
	b := (*quickbarRecord)(memmap.WordPtr(*quickbarUserData(w)))
	i := 0
	for i < 5 && b.Slots[i] != w {
		i++
	}
	if i == 5 {
		return 0
	}
	pos := w.GlobalPos()
	id := quickbarMain().Current[i].ID
	if id == 0 {
		return 1
	}
	highlight := quickbarHighlight(i)
	lit := 0
	if highlight {
		lit = 1
	}
	img := quickbarPointer(unsafe.Pointer(nox_xxx_spellGetAbilityIcon_425310(int(id), lit)))
	w.DrawData().SetTooltip(GetClient().Cli().Strings(), GoWString(nox_xxx_abilityGetName_0_425260(int(id))))
	if img != 0 {
		bookDrawImage(img, pos)
	} else {
		quickbarNoIcon(pos)
	}
	off := uintptr(1047764 + 24*id)
	if *quickbarWord(off + 12) != 0 || *quickbarWord(off + 8) == 0 || Nox_xxx_playerAnimCheck_4372B0() != 0 {
		delay := uint32(nox_xxx_abilityCooldown_4252D0(int(id))) / uint32(GetServer().S().TickRate())
		if delay != 0 && !highlight {
			height := uint32(34) - (GetServer().S().Frame()-*quickbarWord(off + 20))/delay
			GetClient().R2().DrawRectFilledAlpha(pos.X, pos.Y, 34, int(int32(height)))
		}
	}
	return 1
}

func quickbarDrawSpell(w *gui.Window, _ *gui.WindowData) int {
	b := (*quickbarRecord)(memmap.WordPtr(*quickbarUserData(w)))
	i := 0
	for i < 5 && b.Slots[i] != w {
		i++
	}
	if i == 5 {
		return 0
	}
	pos := w.GlobalPos()
	id := b.Current[i].ID
	if id == 0 {
		w.DrawData().SetTooltip(GetClient().Cli().Strings(), GoWStringP(memmap.PtrOff(0x5D4594, 1049716)))
		return 1
	}
	img := bookSpellImage(int(id))
	if quickbarHighlight(i) {
		img = quickbarPointer(nox_xxx_spellIconHighlight_424AB0(int(id)))
	}
	flash := quickbarByte(1049544 + uintptr(id))
	if int8(*flash) > 0 {
		*flash--
	}
	quickbarSlotTooltip(b, i)
	if img != 0 {
		bookDrawImage(img, pos)
	} else {
		quickbarNoIcon(pos)
	}
	mana := int32(nox_xxx_cliGetMana_470DD0())
	cost := int32(nox_xxx_spellManaCost_4249A0(int(id), 1))
	if b == quickbarAt(1047940) {
		for j := 0; j < i; j++ {
			if previous := b.Current[j].ID; previous != 0 {
				mana -= int32(nox_xxx_spellManaCost_4249A0(int(previous), 2))
			}
		}
	}
	p := LocalPlayerDrawable
	buffed := p != nil && *(*uint32)(unsafe.Add(p.C(), int(layoutClientDrawable.Off(124)) /* Buffs */))&(1<<29) != 0
	if !bool(nox_xxx_spellIsEnabled_424B70(int(id))) || Nox_xxx_playerAnimCheck_4372B0() != 0 || buffed {
		GetClient().R2().DrawRectFilledAlpha(pos.X, pos.Y, 30, 30)
	} else if mana < cost && cost != 0 {
		GetClient().R2().DrawRectFilledAlpha(pos.X, pos.Y, 30, int(30*(cost-mana)/cost))
	}
	return 1
}

func quickbarDrawControl(w *gui.Window, d *gui.WindowData) int {
	pos := w.GlobalPos().Add(d.ImgPtVal)
	if int32(*quickbarWord(1049700)) > 0 && *quickbarWord(1049704) == *quickbarUserData(w) {
		bookDrawImage(quickbarPointer(unsafe.Pointer(d.HlImageHnd)), pos)
		*quickbarWord(1049700)--
		if *quickbarWord(1049700) == 0 && *quickbarWord(1049704) <= 2 {
			left := quickbarWindow(1049508).DrawData()
			left.BgImageHnd = left.HlImageHnd
		}
	} else if d.BgImageHnd != nil {
		bookDrawImage(quickbarPointer(unsafe.Pointer(d.BgImageHnd)), pos)
	}
	quickbarLit(quickbarWindow(1049524), bookShown() != 0)
	return 1
}

func quickbarDrawSlide(_ *gui.Window, _ *gui.WindowData) int {
	w := quickbarWindow(1049504)
	pos := w.Off
	diff := int32(pos.Y) - int32(*quickbarWord(1049536))
	if diff < 0 {
		pos.Y++
		w.SetPos(pos)
	} else if diff > 0 {
		pos.Y--
		w.SetPos(pos)
	} else if *quickbarWord(1049536) > uint32(nox_win_height) {
		quickbarRevealTrap()
		*quickbarWord(1049536) = uint32(nox_win_height - 74)
	}
	return 1
}
