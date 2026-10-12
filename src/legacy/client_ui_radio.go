package legacy

import (
	"image"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
)

var _ [4 - unsafe.Sizeof(gui.RadioButtonData{})]byte
var _ [unsafe.Sizeof(gui.RadioButtonData{}) - 4]byte

func uiRadioInit(win *gui.Window) {
	draw := uiRadioDraw
	if win.Flags&128 != 0 {
		draw = uiRadioImage
	}
	win.SetAllFuncs(uiRadioInput, draw, nil)
}
func uiRadioNotify(win *gui.Window, code int, a, b uintptr) {
	win.DrawData().Window.Func94(&gui.RawEvent{Event: code, Arg1: a, Arg2: b})
}
func uiRadioDeselectSiblings(win, parent *gui.Window) {
	if parent == nil {
		return
	}
	for w := parent.Field100Ptr; w != nil; w = w.Prev() {
		if w != win && w.DrawData().Group() == win.DrawData().Group() {
			w.DrawData().Field0 &^= 4
		}
	}
}
func uiRadioInput(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	d := win.DrawData()
	switch ev.EventCode() {
	case 5:
	case 6, 7:
		parent := win.Parent()
		if d.Field0&4 != 0 {
			if d.Field0&2 != 0 {
				return gui.RawEventResp(1)
			}
			return nil
		}
		uiRadioNotify(win, 16391, uintptr(win.C()), a)
		uiRadioDeselectSiblings(win, parent)
		d.Field0 |= 4
	case 8:
		uiRadioNotify(win, 16384, uintptr(win.C()), a)
	case 17:
		if d.Style&0x100 != 0 {
			d.Field0 |= 2
			uiRadioNotify(win, 16389, uintptr(win.C()), a)
			GetClient().Cli().GUI.Focus(win)
		}
	case 18:
		if d.Style&0x100 != 0 {
			d.Field0 &^= 2
			uiRadioNotify(win, 16390, uintptr(win.C()), a)
		}
	case 21:
		switch a {
		case 15, 205, 208, 200, 203: // original tab-navigation calls are no-ops
		case 28, 57:
			if b != 2 {
				return gui.RawEventResp(1)
			}
			parent := win.Parent()
			if d.Field0&4 != 0 {
				return gui.RawEventResp(1)
			}
			uiRadioNotify(win, 16391, uintptr(win.C()), 0)
			uiRadioDeselectSiblings(win, parent)
			d.Field0 ^= 4
		default:
			return nil
		}
	default:
		return nil
	}
	return gui.RawEventResp(1)
}
func uiRadioEvent(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, _ := ev.EventArgsC()
	switch ev.EventCode() {
	case 23:
		if a == 0 {
			win.DrawData().Field0 &^= 2
		}
		uiRadioNotify(win, 16387, a, uintptr(win.ID()))
		return gui.RawEventResp(1)
	case 16385:
		// Preserve the exact 63 UTF-16 code-unit copy, including unpaired surrogates.
		dst := unsafe.Slice((*uint16)(unsafe.Add(win.C(), 108)), 64)
		done := false
		for i := 0; i < 63; i++ {
			v := uint16(0)
			if !done {
				v = *(*uint16)(unsafe.Pointer(a + uintptr(2*i)))
				done = v == 0
			}
			dst[i] = v
		}
		dst[63] = 0
	case 16392:
		parent := win.Parent()
		if win.DrawData().Field0&4 == 0 {
			if a == 1 {
				uiRadioNotify(win, 16391, uintptr(win.C()), 0)
			}
			uiRadioDeselectSiblings(win, parent)
			win.DrawData().Field0 |= 4
		}
	}
	return nil
}
func uiRadioDraw(win *gui.Window, d *gui.WindowData) int {
	r := GetClient().R2()
	p := win.GlobalPos()
	w, h := win.SizeVal.X, win.SizeVal.Y
	border, bg := d.EnColorVal, d.BgColorVal
	if win.Flags&8 != 0 {
		if d.Field0&2 != 0 {
			border = d.HlColorVal
		}
	} else {
		bg = d.DisColorVal
	}
	x, y := p.X+4, p.Y+int(uint32(h)/2)-5
	font := r.GetFonts().AsFont(d.FontPtr)
	fh := r.FontHeight(font)
	if bg != 0x80000000 {
		effectColor(bg)
		r.DrawRectFilledOpaque(x, y, 10, 10, r.Data().Color2())
	}
	if border != 0x80000000 {
		effectColor(border)
		r.DrawBorder(x, y, 10, 10, r.Data().Color2())
	}
	if d.Field0&4 != 0 && d.SelColorVal != 0x80000000 {
		effectColor(d.SelColorVal)
		r.DrawRectFilledOpaque(x+1, y+1, 8, 8, r.Data().Color2())
	}
	if win.Flags&0x2000 != 0 {
		r.SetTextSmooting(true)
	}
	text := d.Text()
	if *(*uint32)(win.WidgetData) != 0 {
		size := r.GetStringSizeWrapped(font, text, 0)
		if d.TextColorVal != 0x80000000 {
			r.Data().SetTextColor(noxcolor.RGBA5551(d.TextColorVal))
			tx, ty := p.X+int(uint32(w)/2)-size.X/2, p.Y+int(uint32(h)/2)-r.FontHeight(font)/2
			r.DrawStringWrapped(font, text, image.Rect(tx, ty, tx+w, ty))
		}
	} else if d.TextColorVal != 0x80000000 {
		r.Data().SetTextColor(noxcolor.RGBA5551(d.TextColorVal))
		tx, ty := x+14, y-fh/2+5
		r.DrawStringWrapped(font, text, image.Rect(tx, ty, tx+w, ty))
	}
	r.SetTextSmooting(false)
	return 1
}
func uiRadioImage(win *gui.Window, d *gui.WindowData) int {
	r := GetClient().R2()
	p := win.GlobalPos()
	w, h := win.SizeVal.X, win.SizeVal.Y
	fg, bg := d.EnImageHnd, d.BgImageHnd
	if win.Flags&8 != 0 {
		if d.Field0&2 != 0 {
			fg = d.HlImageHnd
		}
	} else {
		bg = d.DisImageHnd
	}
	at := p.Add(d.ImagePoint())
	if bg != nil {
		r.DrawImageAt(r.GetBag().AsImage(bg), at)
	}
	if d.Field0&4 != 0 {
		fg = d.SelImageHnd
	}
	if fg != nil {
		r.DrawImageAt(r.GetBag().AsImage(fg), at)
	}
	if win.Flags&0x2000 != 0 {
		r.SetTextSmooting(true)
	}
	font := r.GetFonts().AsFont(d.FontPtr)
	text := d.Text()
	if *(*uint32)(win.WidgetData) != 0 {
		size := r.GetStringSizeWrapped(font, text, 0)
		if d.TextColorVal != 0x80000000 {
			r.Data().SetTextColor(noxcolor.RGBA5551(d.TextColorVal))
			x, y := p.X+int(uint32(w)/2)-size.X/2, p.Y+int(uint32(h)/2)-r.FontHeight(font)/2
			r.DrawStringWrapped(font, text, image.Rect(x, y, x+w, y))
		}
	} else {
		y := p.Y + (h-r.FontHeight(font))/2
		if d.TextColorVal != 0x80000000 {
			r.Data().SetTextColor(noxcolor.RGBA5551(d.TextColorVal))
			x := p.X + 28
			r.DrawStringWrapped(font, text, image.Rect(x, y, x+w, y))
		}
	}
	r.SetTextSmooting(false)
	return 1
}
