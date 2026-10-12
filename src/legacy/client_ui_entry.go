package legacy

import (
	"image"
	"unicode/utf16"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// These ASCII predicates match the exhaustive legacy C-locale capture.
func uiEntryDigit(v uint16) bool { return v >= '0' && v <= '9' }
func uiEntryAlnum(v uint16) bool {
	return uiEntryDigit(v) || v >= 'A' && v <= 'Z' || v >= 'a' && v <= 'z'
}

var uiEntryContext bool
var uiEntryActive *gui.Window
var _ [1056 - unsafe.Sizeof(gui.EntryFieldData{})]byte
var _ [unsafe.Sizeof(gui.EntryFieldData{}) - 1056]byte

func uiEntryData(w *gui.Window) *gui.EntryFieldData { return (*gui.EntryFieldData)(w.WidgetData) }
func uiEntryList(d *gui.EntryFieldData) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(d.Field_1048)))
}
func uiEntryIME() bool { lang := GetServer().S().Strings().Lang(); return lang == 6 || lang == 8 }
func uiEntryIMEAllowed(d *gui.EntryFieldData) bool {
	return uiEntryIME() && d.Field_1036 == 0 && d.Field_1032 == 0 && d.Field_1028 == 0
}
func uiEntryLength(v []uint16) int {
	for i, c := range v {
		if c == 0 {
			return i
		}
	}
	return len(v)
}
func uiEntryComposition(d *gui.EntryFieldData) {
	v := utf16.Encode([]rune(GetClient().GetTextEditBuf()))
	n := min(uiEntryLength(v), 255)
	copy(d.Text[256:256+n], v[:n])
	d.Text[256+n] = 0
	d.Field_1052 = d.Field_1052&0xffff | uint32(n)<<16
	uiEntryList(d).SetHidden(true)
}
func uiEntryAppend(d *gui.EntryFieldData, v uint16) {
	n := int(uint16(d.Field_1052))
	if n >= int(int16(d.Field_1040))-1 {
		return
	}
	d.Text[n] = v
	d.Text[n+1] = 0
	d.Field_1052 = d.Field_1052&0xffff0000 | uint32(uint16(n+1))
}
func uiEntryKey(w *gui.Window, key, state uintptr) gui.WindowEventResp {
	d := uiEntryData(w)
	switch key {
	case 1, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 87, 88, 199, 201, 207, 209:
		return nil
	case 14, 211:
		if key == 14 || !uiEntryIME() {
			if state != 2 {
				return gui.RawEventResp(1)
			}
			if d.Field_1052>>16 == 0 {
				n := uint16(d.Field_1052)
				if n > 0 {
					n--
					d.Text[n] = 0
					d.Field_1052 = d.Field_1052&0xffff0000 | uint32(n)
				}
				return gui.RawEventResp(1)
			}
		}
	case 15, 205, 208, 200, 203:
		return gui.RawEventResp(1)
	case 28, 156:
		if state == 2 && d.Field_1044 == 0 {
			uiRadioNotify(w, 16415, uintptr(w.C()), 0)
		}
		return gui.RawEventResp(1)
	}
	if uiEntryIMEAllowed(d) {
		uiEntryComposition(d)
		return gui.RawEventResp(1)
	}
	v := GetClient().Cli().Inp.KeyToWChar(keybind.Key(uint16(key)))
	if v == 0 || state != 2 {
		return gui.RawEventResp(1)
	}
	if d.Field_1028 != 0 {
		if !uiEntryDigit(v) {
			return gui.RawEventResp(1)
		}
	} else if d.Field_1032 != 0 {
		if !uiEntryAlnum(v) {
			return gui.RawEventResp(1)
		}
	}
	uiEntryAppend(d, uint16(v))
	return gui.RawEventResp(1)
}
func uiEntryInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	d := w.DrawData()
	switch ev.EventCode() {
	case 7:
		d.Field0 |= 2
		GetClient().Cli().GUI.Focus(w)
	case 8:
		if d.Style&0x100 != 0 {
			uiRadioNotify(w, 16384, uintptr(w.C()), 0)
		}
	case 17:
		if d.Style&0x100 != 0 {
			d.Field0 |= 2
			uiRadioNotify(w, 16389, uintptr(w.C()), 0)
			GetClient().Cli().GUI.Focus(w)
		}
	case 18:
		if d.Style&0x100 != 0 {
			d.Field0 &^= 2
			uiRadioNotify(w, 16390, uintptr(w.C()), 0)
		}
	case 21:
		return uiEntryKey(w, a, b)
	default:
		return nil
	}
	return gui.RawEventResp(1)
}
func uiEntryEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, _ := ev.EventArgsC()
	d := uiEntryData(w)
	switch ev.EventCode() {
	case 16413:
		return gui.RawEventResp(uintptr(w.WidgetData))
	case 16414:
		done := false
		for i := 0; i < 255; i++ {
			v := uint16(0)
			if !done {
				v = *(*uint16)(unsafe.Pointer(a + uintptr(i*2)))
				done = v == 0
			}
			d.Text[i] = v
		}
		d.Text[255] = 0
		d.Field_1052 = uint32(uiEntryLength(d.Text[:256]))
		d.Text[256] = 0
	case 2:
		if uiEntryActive == w {
			uiEntryActive = nil
			GetClient().SetTextInput(false)
		}
		uiEntryList(d).Destroy()
		alloc.FreePtr(w.WidgetData)
	case 23:
		if a != 0 {
			uiEntryActive = w
			GetClient().SetTextInput(true)
			w.DrawData().Field0 |= 6
		} else {
			uiEntryActive = nil
			GetClient().SetTextInput(false)
			w.DrawData().Field0 &^= 6
			uiEntryList(d).SetHidden(true)
			d.Text[256] = 0
			d.Field_1052 &= 0xffff
		}
		uiRadioNotify(w, 16387, a, uintptr(w.ID()))
		return gui.RawEventResp(1)
	}
	return nil
}
func uiEntryChar(v uint16) {
	if !uiEntryContext || uiEntryActive == nil {
		return
	}
	d := uiEntryData(uiEntryActive)
	if !uiEntryIMEAllowed(d) {
		return
	}
	d.Field_1044 = 1
	switch v {
	case 7, 8, 9, 11, 12:
	case 10, 13:
		d.Field_1044 = 0
	default:
		uiEntryAppend(d, v)
		uiEntryComposition(d)
	}
}
func uiEntryNew(parent *gui.Window, flags gui.StatusFlags, x, y, width, height int, draw *gui.WindowData, input *gui.EntryFieldData) *gui.Window {
	if draw.Style&gui.StyleEntryField == 0 {
		return nil
	}
	w := GetClient().Cli().GUI.NewWindowRaw(parent, flags, x, y, width, height, uiEntryEvent)
	if w == nil {
		return nil
	}
	w.SetAllFuncs(uiEntryInput, func(w *gui.Window, d *gui.WindowData) int { return uiEntryDraw(w, d, flags&128 != 0) }, nil)
	if draw.Window == nil {
		draw.Window = w
	}
	w.CopyDrawData(draw)
	clear(input.Text[:128])
	clear(input.Text[256:384])
	input.Field_1052 = 0
	input.Field_1044 = 0
	if int16(input.Field_1040) >= 256 {
		input.Field_1040 = 256
	}
	data, _ := alloc.New(gui.EntryFieldData{})
	*data = *input
	w.WidgetData = unsafe.Pointer(data)
	if !uiEntryIME() {
		data.Field_1048 = 0
		return w
	}
	ld := gui.WindowData{Style: 288, EnColorVal: draw.TextColorVal, HlColorVal: 0x80000000, SelColorVal: 0x80000000, TextColorVal: draw.TextColorVal}
	li := gui.ScrollListBoxData{Count: 128, Line_height: 10, Field_2: 1, Field_3: 1}
	list := Nox_gui_newScrollListBox_4A4310(nil, 17584, 0, height, 110, 119, &ld, &li)
	data.Field_1048 = uint32(uintptr(list.C()))
	if list == nil {
		return nil
	}
	list.Flags &^= 128
	uiListInit(list, uiListData(list))
	list.DrawData().BgColorVal = noxcolor.RGB5551Color(0, 0, 0).Color32()
	return w
}
func uiEntryDraw(w *gui.Window, draw *gui.WindowData, images bool) int {
	r := GetClient().R2()
	d := uiEntryData(w)
	d.Field_1044 = 0
	p := w.GlobalPos()
	x, width := p.X, w.SizeVal.X
	f := r.GetFonts().AsFont(draw.FontPtr)
	fh := r.FontHeight(f)
	y := p.Y + w.SizeVal.Y/2 - fh/2
	if w.Flags&0x2000 != 0 {
		r.SetTextSmooting(true)
	}
	if images {
		bg := draw.BgImageHnd
		if w.Flags&8 == 0 {
			bg = draw.DisImageHnd
		}
		if bg != nil {
			r.DrawImageAt(r.GetBag().AsImage(bg), p)
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(draw.TextColorVal))
	}
	if label := draw.Text(); label != "" {
		size := r.GetStringSizeWrapped(f, label, 0)
		if !images {
			r.Data().SetTextColor(noxcolor.RGBA5551(draw.TextColorVal))
		}
		r.DrawStringWrapped(f, label, image.Rect(x+2, y, x+2+width, y))
		x += size.X + 6
		width -= size.X + 6
	}
	if !images {
		if cap := int(d.Field_1042); cap > 0 && width > cap {
			width = cap
			x = p.X + w.SizeVal.X - cap
		}
		bg, border := draw.BgColorVal, draw.EnColorVal
		if w.Flags&8 == 0 {
			bg = draw.DisColorVal
		} else if draw.Field0&2 != 0 {
			border = draw.HlColorVal
		}
		if bg != 0x80000000 {
			effectColor(bg)
			r.DrawRectFilledOpaque(x+1, p.Y+1, width-2, w.SizeVal.Y-2, r.Data().Color2())
		}
		if border != 0x80000000 {
			effectColor(border)
			r.DrawBorder(x, p.Y, width, w.SizeVal.Y, r.Data().Color2())
		}
	}
	if draw.TextColorVal != 0x80000000 {
		var text [256]uint16
		n := uiEntryLength(d.Text[:256])
		copy(text[:], d.Text[:256])
		if d.Field_1024 != 0 {
			for i := 0; i < n; i++ {
				text[i] = '*'
			}
			text[n] = 0
		}
		toString := func(v []uint16) string { return string(utf16.Decode(v[:uiEntryLength(v)])) }
		main := text[:]
		s := toString(main)
		comp := toString(d.Text[256:])
		tw := r.GetStringSizeWrapped(f, s, 0).X
		cw := r.GetStringSizeWrapped(f, comp, 0).X
		if w.Flags&0x4000 != 0 && tw+cw > 0 && (images || width >= 10) && tw+cw+10 > width {
			for main[0] != 0 && tw+cw+10 > width {
				main = main[1:]
				s = toString(main)
				tw = r.GetStringSizeWrapped(f, s, 0).X
			}
		}
		if !images && d.Field_1048 != 0 {
			uiEntryList(d).SetPos(image.Pt(x+tw, fh+y))
		}
		wrap := 0
		if images {
			wrap = width
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(draw.TextColorVal))
		r.DrawStringWrapped(f, s, image.Rect(x+5, y, x+5+wrap, y))
		r.Data().SetTextColor(noxcolor.RGB5551Color(192, 0, 192))
		r.DrawStringWrapped(f, comp, image.Rect(x+tw+5, y, x+tw+5+wrap, y))
		if w == GetClient().Cli().GUI.Focused() {
			blink := memmap.PtrUint8(0x5D4594, 1193344)
			old := *blink
			*blink++
			if old&8 != 0 {
				effectColor(draw.TextColorVal)
				r.DrawRectFilledOpaque(x+tw+cw+5, y, 2, fh, r.Data().Color2())
			}
		}
	}
	r.SetTextSmooting(false)
	return 1
}
