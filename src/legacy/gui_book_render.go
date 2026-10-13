package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"strconv"
	"unsafe"
)

func bookText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "guibook.c")
}
func bookRankText(rank uint32) string {
	return bookFormatInt(bookText("PowerLevel"), int(rank))
}
func bookFormatInt(format string, number int) string {
	var dst [256]uint16
	textFormatBuffer(dst[:], (*uint16)(unsafe.Pointer(internWStr(format))), textFormatWord(uint32(number)))
	return alloc.GoString16(&dst[0])
}
func bookCreatureName(id int) string {
	return GoWStringP(memmap.WordPtr(uint32(bookGuideCreatureName(int32(id)))))
}
func bookDrawList(w *gui.Window) int {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	pos := w.Off
	backPos := pos.Sub(image.Pt(24, 76))
	height := r.FontHeight(font)
	*bookWord(1046656) = uint32(height + 2)
	r.Data().SetTextColor(noxcolor.RGBA5551(*bookWord(1046880)))
	bookDrawImage(*bookWord(1046856), backPos)
	tab := uintptr(1046644)
	if *bookWord(1046872) != 0 {
		tab = 1046660
	}
	bookDrawImage(*bookWord(tab), backPos)
	defer func() {
		if *bookContents() == 0 && int32(*bookWord(1046932)) >= int32(*bookWord(1047508))-int32(*bookWord(1047512))-1 {
			bookHideWindow(bookWindow(*bookWord(1046948)), true)
		}
	}()
	mode := *bookWord(1046868)
	if mode == 2 || mode == 3 {
		off := uintptr(1046924)
		if mode == 3 {
			off = 1046928
		}
		GetClient().Nox_video_drawAnimatedImageOrCursorAt((*ImageRef)(memmap.WordPtr(*bookWord(off))), backPos)
		return 1
	}
	p := bookPlayer
	class := bookClass(p)
	if *bookContents() != 0 {
		perColumn := int(int32(141 / *bookWord(1046656) - 1))
		if *bookWord(1046936) == 0 {
			bookHideWindow(bookWindow(*bookWord(1046944)), true)
		}
		x, y := pos.X+78, pos.Y+19
		for row := 0; row < 2*perColumn; row++ {
			index := uint32(row) + 2*uint32(perColumn)**bookWord(1046936)
			if index >= *bookWord(1047508)-*bookWord(1047512) {
				break
			}
			if row == perColumn {
				x, y = pos.X+199, pos.Y+19
			}
			r.Data().SetTextColor(noxcolor.RGBA5551(*bookWord(1046880)))
			id := int(*bookWord(1046960 + 4*uintptr(index)))
			var text string
			if mode == 1 {
				if class == 2 && !bool(nox_xxx_spellIsEnabled_424B70(id+74)) {
					r.Data().SetTextColor(noxcolor.RGBA5551(*bookWord(1046884)))
				}
				text = bookCreatureName(id)
			} else if class != 0 {
				if !bool(nox_xxx_spellIsEnabled_424B70(id)) {
					r.Data().SetTextColor(noxcolor.RGBA5551(*bookWord(1046884)))
				}
				text = GoWString(nox_xxx_spellTitle_424930(id))
			} else {
				text = GoWString(nox_xxx_abilityGetName_0_425260(id))
			}
			width := r.GetStringSizeWrapped(font, text, 128).X
			r.DrawString(font, text, image.Pt(x-width/2, y))
			y += int(*bookWord(1046656))
		}
		return 1
	}
	id := int(*bookWord(1046960 + 4*uintptr(*bookWord(1046932))))
	wrap := func(text string, x, y, width int) { r.DrawStringWrapped(font, text, image.Rect(x, y, x+width, y)) }
	centered := func(text string, y int) {
		width := r.GetStringSizeWrapped(font, text, 0).X
		wrap(text, (108-width)/2+pos.X+24, y, 128)
	}
	title := func(text string, y int) image.Point {
		size := r.GetStringSizeWrappedStyle(font, text, 108)
		x := pos.X + 24
		if size.Y <= height {
			x += (108 - size.X) / 2
		}
		r.DrawStringWrappedStyle(font, text, image.Rect(x, y, x+128, y))
		return size
	}
	var desc string
	hasDesc := false
	if mode == 1 {
		size := int(bookGuideSize(int32(id)))
		label := bookText("Size") + " "
		imageWidth, dy := 0, 0
		switch size {
		case 1:
			label += bookText("Small")
			imageWidth, dy = 38, 19
		case 2:
			label += bookText("Medium")
			imageWidth = 38
		case 4:
			label += bookText("Large")
			imageWidth = 76
		}
		if class == 2 && (*bookPlayerWord(p, 4232, 0) != 0 || (noxflags.HasGame(noxflags.GameFlag(0x2000)) && !noxflags.HasGame(noxflags.GameModeQuest))) {
			centered(label, pos.Y+51)
		}
		y := pos.Y + 51 + height + 2
		textSize := title(bookCreatureName(id), y)
		bookDrawImage(uint32(bookGuideImage(int32(id))), image.Pt((108-imageWidth)/2+pos.X+24, y+dy+textSize.Y+2))
		raw := bookGuideDescription(int32(id))
		hasDesc = raw != 0
		desc = GoWStringP(memmap.WordPtr(uint32(raw)))
	} else if class == 0 {
		Sub_425450(id)
		title(GoWString(nox_xxx_abilityGetName_0_425260(id)), pos.Y+53)
		raw := sub_4252F0(id)
		hasDesc = raw != nil
		desc = GoWString(raw)
	} else {
		flags := uint16(nox_xxx_spellFlags_424A70(id))
		textSize := title(GoWString(nox_xxx_spellTitle_424930(id)), pos.Y+53)
		mana := nox_xxx_spellManaCost_4249A0(id, 1)
		label := bookText("ManaCost") + " "
		if mana != 0 {
			label += strconv.Itoa(mana)
		} else if bool(nox_xxx_spellHasFlags_424A50(id, 0x800000)) {
			label += "0"
		} else {
			label += "*"
		}
		y := pos.Y + 53 + textSize.Y + 4
		centered(label, y)
		y += height + 2
		for _, v := range []struct {
			flag uint16
			key  string
		}{{0x100, "SpellInstant"}, {4, "SpellTargeted"}, {8, "SpellAtLocation"}, {0x20, "SpellHostile"}} {
			if flags&v.flag != 0 {
				centered(bookText(v.key), y)
				y += height
			}
		}
		centered(bookRankText(*bookPlayerWord(p, 3696, id)), y+height)
		raw := nox_xxx_spellDescription_424A30(id)
		hasDesc = raw != nil
		desc = GoWString(raw)
	}
	if hasDesc {
		h := r.GetStringSizeWrapped(font, desc, 92).Y
		y := (141-h)/2 + pos.Y + 17
		if y > pos.Y+52 {
			y = pos.Y + 52
		}
		wrap(desc, pos.X+153, y, 92)
	}
	return 1
}
