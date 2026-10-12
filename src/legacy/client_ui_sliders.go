package legacy

import (
	"image"
	"math"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)




func uiSliderNew(parent *gui.Window, flags gui.StatusFlags, x, y, w, h int, draw *gui.WindowData, input *gui.SliderData) *gui.Window {
	horizontal := draw.Style&16 != 0
	if !horizontal && draw.Style&8 == 0 {
		return nil
	}
	win := GetClient().Cli().GUI.NewWindowRaw(parent, flags|0x100, x, y, w, h, func(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		return uiSliderEvent(win, ev, horizontal)
	})
	if win == nil {
		return nil
	}
	imageMode := win.Flags&128 != 0
	win.SetAllFuncs(func(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		return uiSliderInput(win, ev, horizontal)
	}, func(win *gui.Window, d *gui.WindowData) int { return uiSliderDraw(win, d, horizontal, imageMode) }, nil)
	if draw.Window == nil {
		draw.Window = win
	}
	win.CopyDrawData(draw)
	thumb := gui.WindowData{Style: 1, Window: win, TextColorVal: draw.TextColorVal}
	childFlags := (flags|0x100)&^0x10 | 0x0c
	if childFlags&128 == 0 {
		thumb.SelColorVal = draw.SelColorVal
		thumb.EnColorVal = draw.BgColorVal
		thumb.BgColorVal = draw.EnColorVal
	} else {
		thumb.EnImageHnd = draw.EnImageHnd
		thumb.DisImageHnd = draw.DisImageHnd
		thumb.SelImageHnd = draw.SelImageHnd
		thumb.BgImageHnd = draw.BgImageHnd
		thumb.HlImageHnd = draw.HlImageHnd
	}
	cw, ch := w, 10
	span := h
	if horizontal {
		cw, ch = 10, h
		span = w
	}
	NewButtonOrCheckbox(win, childFlags, 0, 0, cw, ch, &thumb)
	if input.Max == input.Min {
		input.Max = input.Min + 1
	}
	input.Field2 = math.Float32bits(float32(float64(int32(span-10)) / float64(int32(input.Max-input.Min))))
	data, _ := alloc.New(gui.SliderData{})
	*data = *input
	win.WidgetData = unsafe.Pointer(data)
	return win
}

// The original converts through int64 before narrowing to 32 bits. This also
// preserves the low zero word of x87's indefinite result for NaN/infinity.
func uiSliderInt(v float64) int               { return int(int32(int64(v))) }
func uiSliderScale(d *gui.SliderData) float64 { return float64(math.Float32frombits(d.Field2)) }
func uiSliderPos(win *gui.Window, horizontal bool, v int) {
	p := image.Pt(0, v)
	if horizontal {
		p = image.Pt(v, 0)
	}
	win.Field100Ptr.SetPos(p)
}
func uiSliderNotify(win *gui.Window, code int, a, b uintptr) {
	win.DrawData().Window.Func94(gui.AsWindowEvent(code, a, b))
}
func uiSliderValueNotify(win *gui.Window, code int, d *gui.SliderData) {
	uiSliderNotify(win, code, uintptr(win.C()), uintptr(d.Field3))
}
func uiSliderAxis(win *gui.Window, horizontal bool, packed uintptr) (coord, origin, span, thumb int) {
	p := win.GlobalPos()
	if horizontal {
		return int(uint16(packed)), p.X, win.SizeVal.X, win.Field100Ptr.SizeVal.X
	}
	return int(uint32(packed) >> 16), p.Y, win.SizeVal.Y, win.Field100Ptr.SizeVal.Y
}
func uiSliderEvent(win *gui.Window, ev gui.WindowEvent, horizontal bool) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	d := (*gui.SliderData)(win.WidgetData)
	switch ev.EventCode() {
	case 2:
		if d != nil {
			alloc.Free(d)
		}
	case 23:
		draw := win.DrawData()
		if a != 0 {
			draw.Field0 |= 2
		} else {
			draw.Field0 &^= 2
		}
		uiSliderNotify(win, 16387, a, uintptr(win.ID()))
		return gui.RawEventResp(1)
	case 16391:
		if !horizontal {
			uiSliderValueNotify(win, 16396, d)
		}
	case 16394:
		v := uint32(a)
		if int32(v) >= int32(d.Min) && int32(v) <= int32(d.Max) {
			d.Field3 = v
			delta := v
			if !horizontal {
				delta = d.Max - v
			}
			uiSliderPos(win, horizontal, uiSliderInt(float64(int32(delta))*uiSliderScale(d)))
		}
	case 16395:
		d.Min, d.Max, d.Field3 = uint32(a), uint32(b), uint32(a)
		span := win.SizeVal.Y
		if horizontal {
			span = win.SizeVal.X
		}
		delta := float64(int32(d.Max - d.Min))
		d.Field2 = math.Float32bits(float32(float64(int32(span-10)) / delta))
		pos := 0
		if !horizontal {
			pos = uiSliderInt(delta * uiSliderScale(d))
		}
		uiSliderPos(win, horizontal, pos)
	case 16384:
		coord, origin, span, thumb := uiSliderAxis(win, horizontal, b)
		if coord < origin {
			uiSliderPos(win, horizontal, 0)
			if horizontal {
				d.Field3 = d.Min
			} else {
				d.Field3 = d.Max
			}
		} else if coord >= origin+span {
			uiSliderPos(win, horizontal, span-thumb)
			if horizontal {
				d.Field3 = d.Max
			} else {
				d.Field3 = d.Min
			}
		} else {
			value := int32(uiSliderInt(float64(coord-origin) / uiSliderScale(d)))
			if value > int32(d.Max) {
				value = int32(d.Max)
			}
			d.Field3 = uint32(value)
			if !horizontal {
				d.Field3 = d.Max - d.Field3
			}
		}
		uiSliderValueNotify(win, 16393, d)
	}
	return nil
}
func uiSliderInput(win *gui.Window, ev gui.WindowEvent, horizontal bool) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	d := (*gui.SliderData)(win.WidgetData)
	switch ev.EventCode() {
	case 5:
	case 6, 7:
		coord, origin, span, _ := uiSliderAxis(win, horizontal, a)
		if coord > span+origin-10 {
			coord = span + origin - 10
		}
		uiSliderPos(win, horizontal, coord-origin-5)
		win.Func94(&gui.RawEvent{Event: 16384, Arg2: a})
	case 8, 17, 18:
		draw := win.DrawData()
		if draw.Style&0x100 != 0 {
			switch ev.EventCode() {
			case 8:
				uiSliderNotify(win, 16384, uintptr(win.C()), 0)
			case 17:
				draw.SetBackgroundColor(noxcolor.RGBA5551(draw.HlColorVal))
				uiSliderNotify(win, 16389, uintptr(win.C()), 0)
				GetClient().Cli().GUI.Focus(win)
			case 18:
				draw.SetBackgroundColor(noxcolor.RGBA5551(draw.EnColorVal))
				uiSliderNotify(win, 16390, uintptr(win.C()), 0)
			}
		}
	case 21:
		key := uint32(a)
		// These tab-navigation helpers are no-ops in the original engine.
		if key == 15 || (horizontal && (key == 200 || key == 208)) || (!horizontal && (key == 203 || key == 205)) {
			return gui.RawEventResp(1)
		}
		increase, decrease := uint32(200), uint32(208)
		if horizontal {
			increase, decrease = 203, 205
		}
		if key != increase && key != decrease {
			return nil
		}
		if b != 2 {
			return gui.RawEventResp(1)
		}
		if key == increase {
			if int32(d.Field3) >= int32(d.Max)-1 {
				return gui.RawEventResp(1)
			}
			d.Field3 += 2
		} else {
			if int32(d.Field3) <= int32(d.Min)+1 {
				return gui.RawEventResp(1)
			}
			d.Field3 -= 2
		}
		uiSliderValueNotify(win, 16393, d)
		delta := d.Max - d.Field3
		if horizontal {
			delta = d.Field3 - d.Min
		}
		uiSliderPos(win, horizontal, uiSliderInt(float64(int32(delta))*uiSliderScale(d)))
	default:
		return nil
	}
	return gui.RawEventResp(1)
}
func uiSliderDraw(win *gui.Window, draw *gui.WindowData, horizontal, imageMode bool) int {
	p := win.GlobalPos()
	r := GetClient().R2()
	if imageMode {
		// The original vertical image callback only obtains the position.
		if horizontal {
			img := draw.BgImageHnd
			if win.Flags&8 == 0 {
				img = draw.DisImageHnd
			}
			if img != nil {
				r.DrawImageAt(r.GetBag().AsImage(img), p)
			}
		}
		return 1
	}
	w, h := win.SizeVal.X, win.SizeVal.Y
	track, bg := draw.EnColorVal, draw.BgColorVal
	if win.Flags&8 != 0 {
		if draw.Field0&2 != 0 {
			if !horizontal {
				track = draw.SelColorVal
			}
			if draw.HlColorVal != 0x80000000 {
				effectColor(draw.HlColorVal)
				r.DrawBorder(p.X, p.Y, w, h, r.Data().Color2())
			}
		}
	} else {
		bg = draw.DisColorVal
	}
	if bg != 0x80000000 {
		effectColor(bg)
		r.DrawRectFilledOpaque(p.X+1, p.Y+1, w-2, h-2, r.Data().Color2())
	}
	if track != 0x80000000 {
		effectColor(track)
		if horizontal {
			r.DrawRectFilledOpaque(p.X, p.Y+int(uint32(h)/2)-1, w, 3, r.Data().Color2())
		} else {
			r.DrawRectFilledOpaque(p.X+int(uint32(w)/2)-1, p.Y+4, 3, h-8, r.Data().Color2())
		}
	}
	return 1
}
