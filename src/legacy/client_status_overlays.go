package legacy

import (
	"fmt"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/strman"
	flags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"strings"
	"unsafe"
)

func clientOverlayString(id strman.ID) string {
	return GetServer().S().Strings().GetStringInFile(id, "client.c")
}
func clientLoadingOverlay() {
	if interactionFrameGate() == 0 {
		flags.UnsetEngine(flags.EngineFlag9)
	} else {
		flags.SetEngine(flags.EngineFlag9)
	}
	cache := memmap.PtrUint32(0x5D4594, 814540)
	if *cache == 0 {
		*cache = uiMeterLoadImage("MenuSystemBG")
	}
	if !flags.HasEngine(flags.EngineFlag9) {
		return
	}
	r := GetClient().R2()
	font := r.GetFonts().AsFont(r.GetFonts().FontPtrByName("large"))
	uiMeterImage(*cache, image.Point{})
	text := clientOverlayString("InProgress")
	size := r.GetStringSizeWrapped(font, text, 0)
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	r.DrawString(font, text, image.Pt((int(nox_win_width)-size.X)/2, int(nox_win_height)/2))
}
func clientWinnerOverlay() {
	text := alloc.GoString16(memmap.PtrUint16(0x5D4594, 811376))
	x, y := (int(nox_win_width)-310)/2, (int(nox_win_height)-200)/2
	mode := memmap.Uint32(0x5D4594, 811060)
	if mode > 1 {
		return
	}
	cache := memmap.PtrUint32(0x5D4594, 811888+4*uintptr(mode))
	if *cache == 0 {
		name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 85712+4*uintptr(mode))))
		*cache = uiMeterLoadImage(name)
	}
	uiMeterImage(*cache, image.Pt(x, y))
	r := GetClient().R2()
	size := r.GetStringSizeWrapped(nil, text, 220)
	lineY := y + (49-size.Y)/2 + 143
	white := noxcolor.RGBA5551(nox_color_white_2523948)
	r.Data().SetColor2(white)
	for _, line := range strings.FieldsFunc(text, func(c rune) bool { return c == '\n' || c == '\r' }) {
		r.Data().SetTextColor(white)
		size = r.GetStringSizeWrapped(nil, line, 0)
		r.DrawStringWrapped(nil, line, image.Rect(x+45+(220-size.X)/2, lineY, x+45+(220-size.X)/2+220, lineY))
		lineY += r.FontHeight(nil)
	}
}
func clientDebugOverlay() {
	r := GetClient().R2()
	height := r.FontHeight(nil)
	origin := drawViewport().Screen.Min.Add(image.Pt(10, 90))
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	store := func(text string) string {
		dst := unsafe.Slice(memmap.PtrUint16(0x5D4594, 811120), 80)
		alloc.StrCopyZero16(dst, text)
		return alloc.GoString16(&dst[0])
	}
	r.DrawString(nil, store(browserNarrowText(alloc.GoString(sessionMapFilename()))), origin)
	dr := LocalPlayerDrawable
	p := Get_dword_8531A0_2576()
	if dr == nil || p == nil {
		return
	}
	origin.Y += height
	r.DrawString(nil, store(fmt.Sprintf("X:%d\tY:%d", dr.PosVec.X, dr.PosVec.Y)), origin)
	class := p.Info().PlayerClass()
	name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 29456+4*uintptr(class))))
	label := clientOverlayString(strman.ID(name))
	text := fmt.Sprintf(strings.ReplaceAll(clientOverlayString("PlayerInfo"), "%S", "%s"), int8(p.Level), label)
	origin.Y += height
	r.DrawString(nil, store(text), origin)
}
