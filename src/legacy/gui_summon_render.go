package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func summonIcon(id int) uint32 {
	name := nox_get_thing_name(id)
	if name == nil {
		return 0
	}
	return uint32(bookGuideCage(bookGuideID(GoStringP(unsafe.Pointer(name)))))
}
func summonDraw(w *gui.Window) int {
	mouse := GetClient().GetMousePos()
	r := GetClient().R2()
	switch *summonState() {
	case 1:
		*summonWord(1320992) += 20
		if *summonWord(1320992) >= *summonWord(1321004) {
			*summonWord(1320992) = *summonWord(1321004)
			*summonState() = 2
		}
	case 3:
		*summonWord(1320992) -= 20
		if *summonWord(1320992) <= *summonWord(1321000) {
			*summonWord(1320992) = *summonWord(1321000)
			*summonState() = 0
			summonWin(1321032).Hide()
			summonWin(1321040).Hide()
		}
	}
	origin := image.Pt(int(*summonWord(1320988)), int(*summonWord(1320992)))
	summonWin(1321032).SetPos(origin)
	summonWin(1321040).SetPos(origin.Add(image.Pt(27, 12)))
	pos := w.GlobalPos()
	r.FontHeight(r.GetFonts().AsFont(nil))
	if img := *summonWord(1320996); img != 0 {
		bookDrawImage(img, origin)
	}
	ww, hh := 0, 0
	for i := int32(0); i < 4; i++ {
		s := summonAt(i)
		if s.Active == 0 {
			continue
		}
		switch s.Size {
		case 1:
			ww, hh = 1, 1
		case 2:
			ww, hh = 1, 2
		case 4:
			ww, hh = 2, 2
		}
		x, y := pos.X+38*int(s.X)+2, pos.Y+38*int(s.Y)+2
		width, height := 38*ww-4, 38*hh-4
		if s.Flash != 0 {
			r.Data().SetColor2(noxcolor.RGBA5551(nox_color_yellow_2589772))
			s.Flash = 0
			r.DrawRectFilledOpaque(x, y, width, height, r.Data().Color2())
		} else if img := summonIcon(int(s.Type)); img != 0 {
			bookDrawImage(img, image.Pt(x, y))
		} else {
			cx, cy := x+width/2, y+height/2
			r.Data().SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 956)))
			r.DrawPointRad(image.Pt(cx, cy), 9, r.Data().Color2())
			r.DrawCircle(cx, cy, 9, noxcolor.RGBA5551(memmap.Uint32(0x852978, 4)))
		}
		if cur, max, alt, ok := Sub_495180(int(s.Code)); ok {
			filled := 0
			if max != 0 {
				filled = height * cur / max
			}
			bg, fg := uint32(nox_color_violet_2598268), memmap.Uint32(0x85B3FC, 940)
			if alt {
				bg, fg = memmap.Uint32(0x85B3FC, 984), uint32(dword_8531A0_2572)
			}
			r.Data().SetColor2(noxcolor.RGBA5551(bg))
			r.DrawRectFilledOpaque(width+x-2, y, 2, height, r.Data().Color2())
			r.Data().SetColor2(noxcolor.RGBA5551(fg))
			r.DrawRectFilledOpaque(width+x-2, height+y-filled, 2, filled, r.Data().Color2())
		}
	}
	big := summonPointIn(summonWin(1321040), mouse)
	if summonPointIn(w, mouse) || big || *summonWord(1321212) == 1 {
		xy := [2]int32{int32((mouse.X - pos.X) / 38), int32((mouse.Y - pos.Y) / 38)}
		selected := summonGet(&xy)
		*summonWord(1321212) = 0
		for s := summonFirst(); s != nil; s = summonNext(s) {
			if dr := GetClient().Cli().Objs.ByNetCodeDynamic(int(s.Code)); dr != nil {
				if summonAddress(s) == selected || big {
					dr.ObjFlags |= 0x40000000
					*summonWord(1321212) = 1
				} else {
					dr.ObjFlags &^= 0x40000000
				}
			}
		}
	}
	if menu := summonWin(1321044); menu != nil && !summonPointIn(menu, mouse) {
		summonClose()
	}
	return 1
}
func summonOutline(pos image.Point, fg, bg uint32, text string) int {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	r.Data().SetTextColor(noxcolor.RGBA5551(bg))
	for off := uintptr(184520); off < 184552; off += 8 {
		r.DrawString(font, text, pos.Add(image.Pt(int(memmap.Int32(0x587000, off)), int(memmap.Int32(0x587000, off+4)))))
	}
	r.Data().SetTextColor(noxcolor.RGBA5551(fg))
	return r.DrawString(font, text, pos)
}
func summonDrawMenu(w *gui.Window) int {
	if *summonWord(1321208) == 0 {
		*summonWord(1321208) = uint32(GetServer().S().Types.IndByID("CarnivorousPlant"))
	}
	command := *summonCommandWord(w)
	selected := *summonWord(1321204)
	if selected == 0 && command == 1 {
		return 1
	}
	text := summonText(GoStringP(memmap.LoadPtr(0x587000, 184344+uintptr(command)*4)))
	pos := w.GlobalPos()
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	size := r.GetStringSizeWrapped(font, text, 0)
	mouse := GetClient().GetMousePos()
	r.FontHeight(font)
	x := int((*summonMenuWidth()-uint32(size.X))/2) + 1
	pos = pos.Add(image.Pt(x, 3))
	if summonPointIn(w, mouse) {
		summonOutline(pos, uint32(nox_color_yellow_2589772), uint32(nox_color_black_2650656), text)
		last := memmap.PtrUint32(0x587000, 184552)
		if command != *last {
			*last = command
			bookSound(920)
		}
		return 1
	}
	var fg uint32
	if selected != 0 {
		if summonMobile((*summonRecord)(memmap.WordPtr(selected))) != 0 || command != 4 && command != 5 {
			fg = uint32(nox_color_white_2523948)
		} else {
			fg = memmap.Uint32(0x85B3FC, 956)
		}
	} else {
		if command != 4 && command != 5 || summonAnyMobile() != 0 {
			fg = uint32(nox_color_blue_2650684)
		} else {
			fg = memmap.Uint32(0x85B3FC, 956)
		}
	}
	summonOutline(pos, fg, uint32(nox_color_black_2650656), text)
	return 1
}

// The legacy hit test uses global coordinates; Window.PointIn uses local ones.
func summonPointIn(w *gui.Window, p image.Point) bool {
	return uiWindowPointIn(w, int32(p.X), int32(p.Y))
}
