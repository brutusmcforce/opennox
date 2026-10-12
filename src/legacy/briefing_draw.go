package legacy

import (
	"fmt"
	"image"
	"unicode/utf16"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"golang.org/x/image/font"
)

func briefingBlink(flag *uint32, mid, y int, measure, draw font.Face) int {
	frame := gameFrame()
	if frame%30 != 0 {
		if *flag != 1 {
			return int(frame / 30)
		}
	} else {
		if *flag == 1 {
			*flag = 0
			return 1
		}
		*flag = 1
	}
	r := GetClient().R2()
	s := briefingString("GeneralPrint:QuestSplash12")
	width := r.GetStringSizeWrapped(measure, s, 0).X
	nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
	return r.DrawString(draw, s, image.Pt(mid-width/2, y))
}
func briefingDrawTitle(data *gui.WindowData) int {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(data.FontPtr)
	x, y := int(nox_win_width)/2, int(nox_win_height)/2
	nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
	s := fmt.Sprintf("%s %d", briefingString("Noxworld.c:Stage"), int32(briefingStage()))
	sz := r.GetStringSizeWrapped(font, s, 0)
	r.DrawString(font, s, image.Pt(x-sz.X/2, y+3*(sz.Y-80)))
	if p := briefingCaption; p != nil {
		s = alloc.GoString16(p)
		sz = r.GetStringSizeWrapped(font, s, 0)
		r.DrawString(font, s, image.Pt(x-sz.X/2, y+3*(80-sz.Y)))
	}
	return briefingBlink((*uint32)(&nox_xxx_aSpellphoneme_3_587000_123008), (int(nox_win_width)-640)/2+320, (int(nox_win_height)-480)/2+462, font, font)
}
func briefingDrawStats(data *gui.WindowData) int {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(data.FontPtr)
	x, y := int(nox_win_width)/2, int(nox_win_height)/2
	ox, oy := (int(nox_win_width)-640)/2, (int(nox_win_height)-480)/2
	nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
	title := fmt.Sprintf("%s - %s XX1 %d", briefingString("GUIBrief.c:GauntletStatTitle"), briefingString("Noxworld.c:Stage"), int32(memmap.Uint32(0x5D4594, 831228)))
	sz := r.GetStringSizeWrapped(font, title, 0)
	height := sz.Y
	nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
	r.DrawString(font, title, image.Pt(x-sz.X/2, height+y-240))
	span := int(int32(memmap.Uint32(0x587000, 122968) - memmap.Uint32(0x587000, 122964)))
	step := int(int32(float32(float64(height) * 1.5)))
	rows := &briefingScores
	local, others, count := 0, 0, 0
	for i, row := range rows {
		if row.Player == nil {
			continue
		}
		count++
		if uint32(uintptr(unsafe.Pointer(row.Player))) == uint32(dword_8531A0_2576) {
			local = int(row.Found)
		} else {
			others += int(row.Found)
		}
		px := int(int32(memmap.Uint32(0x587000, 122960+uintptr(i*8)))) + x - 320
		py := int(int32(memmap.Uint32(0x587000, 122964+uintptr(i*8)))) + y - 240
		nox_xxx_drawSetTextColor_434390(int(nox_color_orange_2614256))
		s := fmt.Sprintf("%d) %s", i+1, row.Player.Name())
		units := utf16.Encode([]rune(s))
		limit := int(int32(memmap.Uint32(0x587000, 122968)-memmap.Uint32(0x587000, 122960))) + px - 16
		for {
			size := r.GetStringSizeWrapped(font, s, 0)
			if px+size.X < limit || len(units) <= 5 {
				break
			}
			units = units[:len(units)-1]
			s = string(utf16.Decode(units))
		}
		r.DrawStringWrapped(font, s, image.Rect(px, py, px+span-8, py+height))
		py += step + step/2
		for j, id := range []string{"GUIBrief.c:GeneratorsDestroyed", "GUIBrief.c:numSecretsFound", "GUIBrief.c:Kills", "GUIBrief.c:TotalScore"} {
			width := int(int32(dword_5d4594_832476))
			nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
			r.DrawStringWrapped(font, briefingString(id), image.Rect(px, py, px+width, py+height))
			value := int(row.Generators)
			if j == 1 {
				value = int(row.Secrets)
			} else if j == 2 {
				value = int(row.Kills)
			} else if j == 3 {
				value = int(int32(row.Total))
			}
			color := int(nox_color_green_2614268)
			if j == 3 {
				color = int(nox_color_blue_2650684)
			}
			nox_xxx_drawSetTextColor_434390(color)
			r.DrawStringWrapped(font, fmt.Sprintf(" %d", value), image.Rect(px+width, py, px+span-8, py+height))
			py += step
		}
	}
	s := briefingFormat("GeneralPrint:SecretsTotal", int32(memmap.Uint32(0x5D4594, 832356)))
	sz = r.GetStringSizeWrapped(font, s, 0)
	nox_xxx_drawSetTextColor_434390(int(nox_color_orange_2614256))
	r.DrawString(font, s, image.Pt(ox-sz.X/2+320, oy+3*(150-sz.Y)))
	if local != 0 {
		s = briefingFormat("GeneralPrint:SecretsFound", local)
	} else {
		s = briefingFormat("GeneralPrint:SecretsNoneFound")
	}
	if count > 1 {
		var other string
		if others != 0 {
			other = briefingFormat("GeneralPrint:SecretsFoundByFriends", others)
		} else {
			other = briefingString("GeneralPrint:SecretsNoneFoundByFriends")
		}
		s = fmt.Sprintf("%s - %s", s, other)
	}
	sz = r.GetStringSizeWrapped(font, s, 0)
	nox_xxx_drawSetTextColor_434390(int(nox_color_orange_2614256))
	r.DrawString(font, s, image.Pt(ox-sz.X/2+320, oy+2*(225-sz.Y)))
	return briefingBlink((*uint32)(&dword_587000_122956), ox+320, oy+450, font, font)
}

// briefingDrawInstructions preserves the original text/sprite order, including
// the SoulGate paragraph's height-based right edge and mixed prompt fonts.
func briefingDrawInstructions(data *gui.WindowData) int {
	vp := drawViewport()
	briefingInitSprites()
	Nox_client_resetScreenParticles_431510()
	Nox_xxx_bookHideMB_45ACA0(1)
	Sub_446780()
	r := GetClient().R2()
	titleFont := r.GetFonts().AsFont(data.FontPtr)
	bodyFont := r.GetFonts().AsFont(unsafe.Pointer(uintptr(dword_5d4594_832484)))
	ox, oy := (int(nox_win_width)-640)/2, (int(nox_win_height)-480)/2
	orange, white := int(nox_color_orange_2614256), int(nox_color_white_2523948)
	text := briefingString("GeneralPrint:QuestSplash1")
	x := ox - r.GetStringSizeWrapped(titleFont, text, 0).X/2 + 320
	nox_xxx_drawSetTextColor_434390(int(nox_color_black_2650656))
	for _, p := range []image.Point{image.Pt(x-1, oy+19), image.Pt(x+1, oy+19), image.Pt(x-1, oy+21), image.Pt(x+1, oy+21)} {
		r.DrawString(titleFont, text, p)
	}
	nox_xxx_drawSetTextColor_434390(orange)
	r.DrawString(titleFont, text, image.Pt(x, oy+20))
	sprites := briefingSprites()
	sprite := func(index, x, y int) {
		dr := (*client.Drawable)(unsafe.Pointer(uintptr(*sprites[index])))
		p := image.Pt(ox+x, oy+y)
		dr.PosVec = vp.ToWorldPos(p)
		dr.CallDraw(vp)
	}
	left := func(n string, x, y, right int) {
		s := briefingString("GeneralPrint:QuestSplash" + n + "a")
		nox_xxx_drawSetTextColor_434390(orange)
		r.DrawString(bodyFont, s, image.Pt(ox+x, oy+y))
		x = ox + x + r.GetStringSizeWrapped(bodyFont, s, 0).X + 4
		nox_xxx_drawSetTextColor_434390(white)
		s = briefingString("GeneralPrint:QuestSplash" + n + "b")
		r.DrawStringWrapped(bodyFont, s, image.Rect(x, oy+y, right, oy+y))
	}
	right := func(n string, y, edge, limit, fallback int) {
		a := briefingString("GeneralPrint:QuestSplash" + n + "a")
		aw := r.GetStringSizeWrapped(bodyFont, a, 0).X
		b := briefingString("GeneralPrint:QuestSplash" + n + "b")
		bw := r.GetStringSizeWrapped(bodyFont, b, 0).X
		x := ox + edge - bw - aw - 4
		if aw+bw > limit {
			x = ox + fallback
		}
		nox_xxx_drawSetTextColor_434390(orange)
		r.DrawString(bodyFont, a, image.Pt(x, oy+y))
		nox_xxx_drawSetTextColor_434390(white)
		if aw+bw <= limit {
			r.DrawString(bodyFont, b, image.Pt(ox+edge-bw, oy+y))
		} else {
			x += aw + 4
			r.DrawStringWrapped(bodyFont, b, image.Rect(x, oy+y, ox+edge, oy+y))
		}
	}
	sprite(1, 73, 123)
	left("2", 109, 76, ox+520)
	sprite(0, 565, 117)
	right("3", 115, 520, 390, 199)
	sprite(3, 133, 192)
	left("4", 157, 156, oy+630)
	sprite(2, 525, 222)
	right("7", 198, 500, 215, 250)
	sprite(9, 182, 262)
	sprite(11, 201, 251)
	sprite(10, 185, 234)
	left("5", 221, 240, ox+470)
	sprite(6, 484, 278)
	sprite(7, 503, 303)
	right("6", 286, 462, 350, 113)
	sprite(5, 186, 333)
	sprite(4, 219, 345)
	sprite(8, 220, 322)
	left("8", 241, 330, ox+550)
	for i, n := range []string{"9", "10", "11"} {
		a := briefingString("GeneralPrint:QuestSplash" + n + "a")
		b := briefingString("GeneralPrint:QuestSplash" + n + "b")
		aw := r.GetStringSizeWrapped(bodyFont, a, 0).X
		bw := r.GetStringSizeWrapped(bodyFont, b, 0).X
		x := ox - (aw+bw)/2 + 320
		y := oy + 370 + 25*i
		nox_xxx_drawSetTextColor_434390(orange)
		r.DrawString(bodyFont, a, image.Pt(x, y))
		nox_xxx_drawSetTextColor_434390(white)
		r.DrawString(bodyFont, b, image.Pt(x+aw+4, y))
	}
	return briefingBlink(memmap.PtrUint32(0x587000, 123012), ox+320, oy+450, bodyFont, titleFont)
}
