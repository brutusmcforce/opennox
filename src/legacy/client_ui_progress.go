package legacy

import (
	"fmt"
	"image"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
)

func uiProgressNew(parent *gui.Window, status gui.StatusFlags, x, y, w, h int, draw *gui.WindowData) *gui.Window {
	if draw.Style&0x1000 == 0 {
		return nil
	}
	win := GetClient().Cli().GUI.NewWindowRaw(parent, status&^8, x, y, w, h, uiProgressEvent)
	if win == nil {
		return nil
	}
	if win.Flags&128 == 0 {
		win.SetAllFuncs(nil, uiProgressDraw, nil)
	} else {
		win.SetAllFuncs(nil, uiProgressImage, nil)
	}
	if draw.Window == nil {
		draw.Window = win
	}
	win.CopyDrawData(draw)
	return win
}
func uiProgressEvent(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16416 {
		a, _ := ev.EventArgsC()
		v := int32(a)
		if v >= 0 && v <= 100 {
			*(*uint32)(unsafe.Add(win.C(), 32)) = uint32(v)
		}
	}
	return nil
}
func uiProgressValue(win *gui.Window) int { return int(*(*int32)(unsafe.Add(win.C(), 32))) }
func uiProgressDraw(win *gui.Window, draw *gui.WindowData) int {
	r := GetClient().R2()
	p := win.GlobalPos()
	w, h := win.SizeVal.X, win.SizeVal.Y
	if draw.BgColorVal != 0x80000000 {
		effectColor(draw.BgColorVal)
		r.DrawRectFilledOpaque(p.X, p.Y, w, h, r.Data().Color2())
	}
	if draw.HlColorVal != 0x80000000 {
		effectColor(draw.HlColorVal)
		r.DrawRectFilledOpaque(p.X, p.Y, w*uiProgressValue(win)/100, h, r.Data().Color2())
	}
	if draw.TextColorVal != 0x80000000 {
		if win.Flags&0x2000 != 0 {
			r.SetTextSmooting(true)
		}
		text := fmt.Sprintf("%d%%", uiProgressValue(win))
		font := r.GetFonts().AsFont(draw.FontPtr)
		size := r.GetStringSizeWrapped(font, text, 64)
		x, y := p.X+w/2-size.X/2, p.Y+h/2-size.Y/2+1
		r.Data().SetTextColor(noxcolor.RGBA5551(draw.TextColorVal))
		r.DrawStringWrapped(font, text, image.Rect(x, y, x+w, y))
		r.SetTextSmooting(false)
	}
	if draw.EnColorVal != 0x80000000 {
		effectColor(draw.EnColorVal)
		r.DrawBorder(p.X, p.Y, w, h, r.Data().Color2())
	}
	return 1
}
func uiProgressImage(win *gui.Window, draw *gui.WindowData) int {
	p := win.GlobalPos()
	w := win.SizeVal.X * uiProgressValue(win) / 100
	objectRenderSaveClip()
	uiRenderCopyRect(p.X, p.Y, w, win.SizeVal.Y)
	if draw.BgImageHnd != nil && w > 0 {
		r := GetClient().R2()
		r.DrawImageAt(r.GetBag().AsImage(draw.BgImageHnd), p)
	}
	objectRenderRestoreClip()
	return 1
}
