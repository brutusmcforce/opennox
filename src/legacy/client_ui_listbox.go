package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"math"
	"unicode/utf16"
	"unsafe"
)

func uiListData(w *gui.Window) *gui.ScrollListBoxData { return (*gui.ScrollListBoxData)(w.WidgetData) }
func uiListRows(d *gui.ScrollListBoxData) []gui.ScrollListBoxItem {
	return unsafe.Slice(d.Items, int(d.Count))
}
func uiListSelection(d *gui.ScrollListBoxData) []int32 {
	return unsafe.Slice((*int32)(memmap.WordPtr(d.Field_12)), int(d.Count)+1)
}
func uiListRemoveSelection(d *gui.ScrollListBoxData, index int) {
	s := uiListSelection(d)
	copy(s[index:int(d.Count)], s[index+1:int(d.Count)+1])
	s[int(d.Count)-1] = -1
}
func uiListIndex(d *gui.ScrollListBoxData) int {
	if d.Count == 0 {
		return 0
	}
	rows := uiListRows(d)
	top := int32(int16(d.Field_13_1))
	if rows[0].Field_0 > uint32(top) {
		return 0
	}
	for i := 0; i < int(int16(d.Field_11_0)); i++ {
		if i+1 >= int(d.Count) {
			break
		}
		if int32(rows[i+1].Field_0) > top {
			return i + 1
		}
	}
	return 0
}
func uiListScroll(w *gui.Window, delta int, snap bool) int {
	d := uiListData(w)
	index := uiListIndex(d) + delta
	if index < 0 {
		index = 0
	} else if index >= int(int16(d.Field_11_0)) {
		index = int(int16(d.Field_11_0)) - 1
	}
	if snap {
		d.Field_13_1 = 0
		if index > 0 {
			d.Field_13_1 = uint16(uiListRows(d)[index-1].Field_0) + 1
		}
	}
	slider := (*gui.Window)(d.Field_9)
	if slider == nil {
		return 0
	}
	sd := (*gui.SliderData)(slider.WidgetData)
	maximum := int32(d.Field_10) - int32(int16(d.Field_13_0)) + 3
	if maximum < 0 {
		maximum = 0
	}
	sd.Max = uint32(maximum)
	sd.Field2 = math.Float32bits(float32(float64(slider.Size().Y-slider.Field100Ptr.Size().Y) / float64(maximum)))
	if snap {
		return gui.EventRespInt(slider.Func94(&gui.RawEvent{Event: 16394, Arg1: uintptr(sd.Max - uint32(int32(int16(d.Field_13_1))))}))
	}
	return int(uintptr(slider.WidgetData))
}
func uiListRecalculate(w *gui.Window) int {
	d := uiListData(w)
	total := uint32(0)
	if d.Field_11_0 > 0 {
		// The C loop always visits its first row before the signed count comparison.
		for i := 0; ; i++ {
			row := &uiListRows(d)[i]
			total += uint32(uint8(row.Field_130)) + 1
			row.Field_0 = total
			if i+1 >= int(int16(d.Field_11_0)) {
				break
			}
		}
	}
	d.Field_10 = total
	return uiListScroll(w, 0, true)
}
func uiListHit(w *gui.Window, packed uint32, signed bool) (int, bool) {
	d := uiListData(w)
	y := w.GlobalPos().Y
	if w.DrawData().Text() != "" {
		y += GetClient().R2().FontHeight(w.DrawData().Font()) + 1
	}
	limit := int32(int16(d.Field_13_0)) + int32(int16(d.Field_13_1))
	mouse := uint32(packed>>16) + uint32(int32(int16(d.Field_13_1))) - uint32(y)
	rows := uiListRows(d)
	for i := 0; ; i++ {
		beyond := false
		if i > 0 {
			if signed {
				beyond = int32(rows[i-1].Field_0) > limit
			} else {
				beyond = rows[i-1].Field_0 > uint32(limit)
			}
		}
		if beyond || i == int(int16(d.Field_11_0)) {
			return i, false
		}
		if signed {
			if int32(rows[i].Field_0) > int32(mouse) {
				return i, true
			}
		} else if rows[i].Field_0 > mouse {
			return i, true
		}
	}
}
func uiListSingleInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	code := ev.EventCode()
	d := uiListData(w)
	notify := func(code int, value uint32) { uiRadioNotify(w, code, uintptr(w.C()), uintptr(value)) }
	if code == 21 {
		switch a {
		case 15, 205, 203:
			return gui.RawEventResp(1)
		case 28, 57:
			if b == 1 {
				notify(16400, d.Field_12)
			}
			return gui.RawEventResp(1)
		case 200:
			if b != 2 {
				return gui.RawEventResp(1)
			}
			code = 19
		case 208:
			if b != 2 {
				return gui.RawEventResp(1)
			}
			code = 20
		default:
			return nil
		}
	}
	switch code {
	case 5, 17, 18:
	case 6, 7:
		old := d.Field_12
		index, hit := uiListHit(w, uint32(a), true)
		d.Field_12 = 0xffffffff
		if hit && !(uint32(index) == old && d.Field_5 == 0) {
			d.Field_12 = uint32(index)
		}
		if int32(d.Field_12) < 0 && d.Field_5 != 0 {
			d.Field_12 = old
		}
		notify(16400, d.Field_12)
	case 8:
		if w.DrawData().Style&0x100 != 0 {
			notify(16384, 0)
		}
	case 10, 11:
		index, hit := uiListHit(w, uint32(a), false)
		if !hit {
			index = -1
		}
		notify(16401, uint32(index))
	case 19, 20:
		sel := int32(d.Field_12)
		if sel == -1 {
			d.Field_12 = 0
			uiListScroll(w, 0, true)
		} else if code == 19 && sel > 0 {
			d.Field_12--
			if uiListRows(d)[d.Field_12].Field_0 < uint32(int32(int16(d.Field_13_1))) {
				uiListScroll(w, -1, true)
			}
		} else if code == 20 && sel < int32(int16(d.Field_11_0))-1 {
			d.Field_12++
			if uiListRows(d)[d.Field_12].Field_0 > uint32(int32(int16(d.Field_13_0))+int32(int16(d.Field_13_1))) {
				uiListScroll(w, 1, true)
			}
		}
		notify(16400, d.Field_12)
	default:
		return nil
	}
	return gui.RawEventResp(1)
}
func uiListMultiInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, _ := ev.EventArgsC()
	d := uiListData(w)
	switch ev.EventCode() {
	case 5, 17, 18:
	case 6, 7:
		index, hit := uiListHit(w, uint32(a), false)
		if !hit {
			index = -1
		}
		s := uiListSelection(d)
		i := 0
		for s[i] >= 0 {
			if s[i] == int32(index) {
				uiListRemoveSelection(d, i)
				uiRadioNotify(w, 16400, uintptr(w.C()), uintptr(index))
				return gui.RawEventResp(1)
			}
			i++
		}
		s[i] = int32(index)
		if i < int(d.Count) {
			s[i+1] = -1
		}
		uiRadioNotify(w, 16400, uintptr(w.C()), uintptr(index))
	case 8:
		if w.DrawData().Style&0x100 != 0 {
			uiRadioNotify(w, 16384, uintptr(w.C()), 0)
		}
	case 10, 11:
		index, hit := uiListHit(w, uint32(a), false)
		if !hit {
			index = -1
		}
		uiRadioNotify(w, 16401, uintptr(w.C()), uintptr(index))
	case 19:
		if d.Field_7 != nil && int16(d.Field_13_1) > 0 {
			uiListScroll(w, -1, true)
		}
	case 20:
		if d.Field_8 != nil && int32(int16(d.Field_13_0))+int32(int16(d.Field_13_1)) <= int32(d.Field_10) {
			uiListScroll(w, 1, true)
		}
	case 21:
		if a != 15 {
			return nil
		}
	default:
		return nil
	}
	return gui.RawEventResp(1)
}

// uiListCopy preserves the C bounded-copy padding and untouched final unit
// for short source strings. The final unit is set only when the bound is met.
func uiListCopy(dst []uint16, src uintptr, bound int) {
	done := false
	for i := 0; i < bound; i++ {
		v := uint16(0)
		if !done {
			v = *(*uint16)(unsafe.Pointer(src + uintptr(2*i)))
			done = v == 0
		}
		dst[i] = v
	}
	if !done {
		dst[bound] = 0
	}
}
func uiListAppend(w *gui.Window, src uintptr, color int32) int {
	d := uiListData(w)
	row := &uiListRows(d)[int(int16(d.Field_11_1))]
	row.Field_129 = w.DrawData().TextColorVal
	if color >= 0 && color < 17 {
		row.Field_129 = *ColorTable[color]
	}
	if src != 0 {
		uiListCopy(row.Text[:], src, 255)
		row.Text[255] = 0
		n := uiEntryLength(row.Text[:])
		if n > 0 && row.Text[n-1] == 10 {
			row.Text[n-1] = 0
		}
	} else {
		row.Text[0] = ' '
		row.Text[1] = 0
	}
	r := GetClient().R2()
	height := r.FontHeight(w.DrawData().Font())
	if w.Flags&0x4000 == 0 {
		width := w.SizeVal.X - 7
		if d.Field_3 != 0 {
			width -= 10
		}
		height = r.GetStringSizeWrapped(w.DrawData().Font(), string(utf16.Decode(row.Text[:uiEntryLength(row.Text[:])])), width).Y
	}
	row.Field_130 = row.Field_130&0xffffff00 | uint32(uint8(height))
	d.Field_11_0++
	d.Field_11_1++
	return uiListRecalculate(w)
}
func uiListEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	code := ev.EventCode()
	d := uiListData(w)
	ret := func(v int) gui.WindowEventResp { return gui.RawEventResp(v) }
	switch code {
	case 2:
		if d != nil {
			alloc.FreePtr(unsafe.Pointer(d.Items))
			if d.Field_4 != 0 {
				alloc.FreePtr(memmap.WordPtr(d.Field_12))
			}
			alloc.FreePtr(unsafe.Pointer(d))
		}
		w.WidgetData = nil
	case 23:
		if a != 0 {
			w.DrawData().Field0 |= 2
		} else {
			w.DrawData().Field0 &^= 2
		}
		uiRadioNotify(w, 16387, a, uintptr(w.ID()))
		return ret(1)
	case 16385:
		uiListCopy(unsafe.Slice((*uint16)(unsafe.Add(w.C(), int(layoutGuiWindow.Off(108)) /* drawData.text */)), 64), a, 63)
	case 16388:
		up, down, slider := (*gui.Window)(d.Field_7), (*gui.Window)(d.Field_8), (*gui.Window)(d.Field_9)
		if up != nil {
			up.SetPos(image.Pt(int(a)-up.SizeVal.X, 0))
		}
		if down != nil {
			down.SetPos(image.Pt(int(a)-down.SizeVal.X, int(b)-down.SizeVal.Y))
		}
		if slider != nil {
			slider.SetPos(image.Pt(int(a)-slider.SizeVal.X, up.SizeVal.Y))
			uiWindowResize(slider, slider.SizeVal.X, int(b)-2*slider.Field100Ptr.SizeVal.Y)
		}
		d.Field_13_0 = uint16(b)
		if w.DrawData().Text() != "" {
			d.Field_13_0 -= uint16(GetClient().R2().FontHeight(w.DrawData().Font()))
		}
	case 16384, 16391:
		if a == uintptr(d.Field_7) {
			if d.Field_13_1 > 0 {
				uiListScroll(w, -1, true)
			}
		} else if a == uintptr(d.Field_8) && uint32(d.Field_13_0)+uint32(d.Field_13_1) <= d.Field_10 {
			uiListScroll(w, 1, true)
		}
	case 16393:
		slider := (*gui.Window)(d.Field_9)
		sd := (*gui.SliderData)(slider.WidgetData)
		v := int16(uint16(sd.Max) - uint16(b))
		maximum := int32(d.Field_10) - int32(int16(d.Field_13_0)) + 1
		d.Field_13_1 = uint16(v)
		if int32(v) > maximum {
			d.Field_13_1 = uint16(maximum)
		}
		uiListScroll(w, 0, false)
	case 16397:
		cursor, count := int(int16(d.Field_11_1)), int(int16(d.Field_11_0))
		if cursor != count {
			if count >= int(d.Count) {
				if d.Field_2 == 0 {
					return nil
				}
				w.Func94(&gui.RawEvent{Event: 16411, Arg1: 1})
			}
			rows := uiListRows(d)
			for i := int(int16(d.Field_11_0)) - 1; i >= int(int16(d.Field_11_1)); i-- {
				rows[i+1] = rows[i]
			}
		} else if cursor >= int(d.Count) {
			if d.Field_2 == 0 {
				return nil
			}
			w.Func94(&gui.RawEvent{Event: 16411, Arg1: 1})
		}
		uiListAppend(w, a, int32(b))
		if d.Field_1 != 0 {
			for uiListRows(d)[int(d.Field_11_1)-1].Field_0 >= uint32(d.Field_13_1)+uint32(d.Field_13_0) {
				old := d.Field_13_1
				uiListScroll(w, 1, true)
				if d.Field_13_1 == old {
					break
				}
			}
		}
		if d.Field_4 != 0 {
			s := uiListSelection(d)
			for i := 0; s[i] >= 0; i++ {
				if int32(d.Field_11_1) < s[i] {
					s[i]++
				}
			}
		} else if int32(d.Field_11_1) < int32(d.Field_12) {
			d.Field_12++
		}
		return ret(1)
	case 16398:
		index := int(int32(a))
		count := int(d.Field_11_0)
		if index < 0 || index >= count {
			return nil
		}
		rows := uiListRows(d)
		copy(rows[index:count-1], rows[index+1:count])
		d.Field_11_0--
		d.Field_11_1 = d.Field_11_0
		rows[d.Field_11_0] = gui.ScrollListBoxItem{}
		if d.Field_4 != 0 {
			s := uiListSelection(d)
			for i := 0; s[i] >= 0; i++ {
				if int32(index) < s[i] {
					s[i]--
				} else if int32(index) == s[i] {
					uiListRemoveSelection(d, i)
					i--
				}
			}
		} else if int32(index) < int32(d.Field_12) {
			d.Field_12--
		} else if int32(index) == int32(d.Field_12) {
			d.Field_12 = 0xffffffff
		}
		uiListRecalculate(w)
	case 16399:
		clear(uiListRows(d))
		if a != 1 {
			d.Field_13_1 = 0
		}
		if d.Field_4 != 0 {
			s := uiListSelection(d)
			for i := 0; i < int(d.Count); i++ {
				s[i] = -1
			}
		} else {
			d.Field_12 = 0xffffffff
		}
		d.Field_11_1 = 0
		d.Field_11_0 = 0
		d.Field_10 = 0
		uiListScroll(w, 0, true)
	case 16402:
		index := int32(a)
		if index < 0 {
			d.Field_11_1 = 0
		} else if index <= int32(int16(d.Field_11_0)) {
			d.Field_11_1 = uint16(a)
		} else {
			d.Field_11_1 = d.Field_11_0
		}
	case 16403, 16405:
		index := int(int32(a))
		if index < 0 || index >= int(d.Count) {
			if d.Field_4 != 0 {
				s := uiListSelection(d)
				for i := 0; i < int(d.Count); i++ {
					s[i] = -1
				}
			} else if code == 16403 {
				d.Field_12 = 0xffffffff
			}
			return nil
		}
		row := &uiListRows(d)[index]
		if row.Text[0] == 0 {
			return nil
		}
		if d.Field_4 != 0 {
			s := uiListSelection(d)
			if code == 16403 {
				s[0] = int32(index)
				s[1] = -1
				return nil
			}
			i := 0
			for s[i] >= 0 {
				if s[i] == int32(index) {
					uiListRemoveSelection(d, i)
					return nil
				}
				i++
			}
			s[i] = int32(index)
			s[i+1] = -1
			return nil
		}
		if code == 16405 {
			return nil
		}
		d.Field_12 = uint32(index)
		if row.Field_0 < uint32(d.Field_13_1) {
			w.Func94(&gui.RawEvent{Event: 16412, Arg1: a})
		} else if row.Field_0 > uint32(int32(d.Field_13_1)+int32(int16(d.Field_13_0))) {
			d.Field_13_1 = 0
			if index > 0 {
				d.Field_13_1 = uint16(row.Field_0) - d.Field_13_0
			}
			uiListScroll(w, 0, true)
		}
	case 16404:
		return ret(int(d.Field_12))
	case 16406:
		index := int(int32(a))
		if index >= 0 && index < int(d.Count) {
			return ret(int(uintptr(unsafe.Pointer(&uiListRows(d)[index].Text[0]))))
		}
	case 16407:
		index := int(int32(b))
		if index >= 0 && index < int(d.Count) {
			uiListCopy(uiListRows(d)[index].Text[:], a, 255)
		}
	case 16408:
		d.Field_7 = unsafe.Pointer(a)
	case 16409:
		d.Field_8 = unsafe.Pointer(a)
	case 16410:
		d.Field_9 = unsafe.Pointer(a)
	case 16411:
		index := int(int32(a))
		count := int(d.Field_11_0)
		if index < 0 || index > count {
			return nil
		}
		rows := uiListRows(d)
		copy(rows[:count-index], rows[index:count])
		d.Field_11_0 -= uint16(index)
		d.Field_11_1 = d.Field_11_0
		if d.Field_4 != 0 {
			s := uiListSelection(d)
			for i := 0; s[i] >= 0; i++ {
				if int32(index) < s[i] {
					uiListRemoveSelection(d, i)
					i--
				} else {
					s[i] -= int32(index)
				}
			}
		} else if int32(d.Field_12) > 0 {
			d.Field_12 -= uint32(index)
		}
		if d.Field_13_1 > 0 {
			uiListScroll(w, -index, true)
		}
		uiListRecalculate(w)
		return ret(1)
	case 16412:
		index := int32(a) - 1
		d.Field_13_1 = 0
		if index >= 0 && index < int32(d.Count) {
			d.Field_13_1 = uint16(uiListRows(d)[index].Field_0 + 1)
		}
		view := int32(int16(d.Field_13_0))
		total := int32(d.Field_10)
		if int32(d.Field_13_1)+view >= total {
			d.Field_13_1 = uint16(total - view)
		}
		uiListScroll(w, 0, true)
	}
	return nil
}

func uiListInit(w *gui.Window, d *gui.ScrollListBoxData) {
	if w == nil {
		return
	}
	input := uiListSingleInput
	if d.Field_4 != 0 {
		input = uiListMultiInput
	}
	images := w.Flags&128 != 0
	w.SetAllFuncs(input, func(w *gui.Window, draw *gui.WindowData) int { return uiListDraw(w, draw, images) }, nil)
}
func uiListDraw(w *gui.Window, draw *gui.WindowData, images bool) int {
	r := GetClient().R2()
	d := uiListData(w)
	p := w.GlobalPos()
	x, y := p.X, p.Y
	width, height := w.SizeVal.X, w.SizeVal.Y
	font := draw.Font()
	fh := r.FontHeight(font)
	if w.Flags&0x2000 != 0 {
		r.SetTextSmooting(true)
	}
	objectRenderSaveClip()
	if d.Field_3 != 0 {
		width -= 10
	}
	if images {
		bg := draw.BgImageHnd
		if w.Flags&8 == 0 {
			bg = draw.DisImageHnd
		}
		if bg != nil {
			r.DrawImageAt(r.GetBag().AsImage(bg), p)
		}
	}
	if draw.Text() != "" {
		r.Data().SetTextColor(noxcolor.RGBA5551(draw.TextColorVal))
		limit := fh
		if images {
			limit = 0
		}
		r.DrawStringWrapped(font, draw.Text(), image.Rect(x+1, y, x+1+width, y+limit))
		y += fh + 1
		height -= fh + 1
	}
	border := draw.EnColorVal
	if !images {
		bg := draw.BgColorVal
		if w.Flags&8 == 0 {
			bg = draw.DisColorVal
		} else if draw.Field0&2 != 0 {
			border = draw.HlColorVal
		}
		if bg != 0x80000000 {
			effectColor(bg)
			r.DrawRectFilledOpaque(x, y, width, height, r.Data().Color2())
		}
	}
	if border != 0x80000000 {
		effectColor(border)
		r.DrawBorder(x, y, width, height, r.Data().Color2())
	}
	uiRenderCopyRect(x, y, width, height)
	top := int32(int16(d.Field_13_1))
	limit := top + int32(int16(d.Field_13_0))
	ry := y - int(top)
	if draw.TextColorVal != 0x80000000 {
		rows := uiListRows(d)
		for i := 0; ; i++ {
			if (i > 0 && rows[i-1].Field_0 > uint32(limit)) || i == int(int16(d.Field_11_0)) {
				break
			}
			row := &rows[i]
			rh := int(uint8(row.Field_130)) + 1
			if row.Field_0 >= uint32(top) {
				r.Data().SetTextColor(noxcolor.RGBA5551(row.Field_129))
				selected := uint32(i) == d.Field_12
				if d.Field_4 != 0 {
					selected = false
					s := uiListSelection(d)
					for j := 0; s[j] >= 0; j++ {
						if s[j] == int32(i) {
							selected = true
							break
						}
					}
				}
				if selected && draw.SelColorVal != 0x80000000 {
					effectColor(draw.SelColorVal)
					r.DrawRectFilledOpaque(x, ry, width, rh, r.Data().Color2())
				}
				r.Data().SetTextColor(noxcolor.RGBA5551(row.Field_129))
				units := row.Text[:uiEntryLength(row.Text[:])]
				text := string(utf16.Decode(units))
				if w.Flags&0x4000 != 0 {
					for r.GetStringSizeWrapped(font, text, 0).X > width-7 && len(units) > 0 {
						units = units[:len(units)-1]
						text = string(utf16.Decode(units))
					}
				}
				r.DrawStringWrapped(font, text, image.Rect(x+5, ry+2, x+5+width-7, ry+2+rh))
			}
			ry += rh
		}
	}
	objectRenderRestoreClip()
	r.SetTextSmooting(false)
	return 1
}
func uiListNew(parent *gui.Window, flags gui.StatusFlags, x, y, width, height int, draw *gui.WindowData, opts *gui.ScrollListBoxData) *gui.Window {
	fh := GetClient().R2().FontHeight(GetClient().R2().GetFonts().AsFont(draw.FontPtr))
	if int(opts.Line_height) < fh {
		opts.Line_height = uint16(fh)
	}
	if draw.Style&gui.StyleScrollListBox == 0 {
		return nil
	}
	w := GetClient().Cli().GUI.NewWindowRaw(parent, flags, x, y, width, height, uiListEvent)
	uiListInit(w, opts)
	if w == nil {
		return nil
	}
	if draw.Window == nil {
		draw.Window = w
	}
	w.CopyDrawData(draw)
	rows, _ := alloc.Make([]gui.ScrollListBoxItem{}, int(opts.Count))
	opts.Items = unsafe.SliceData(rows)
	opts.Field_13_0 = uint16(height)
	label := draw.Text() != ""
	if label {
		opts.Field_13_0 -= uint16(fh)
	}
	opts.Field_13_1 = 0
	opts.Field_12 = 0xffffffff
	opts.Field_11_1 = 0
	opts.Field_11_0 = 0
	opts.Field_10 = 0
	if opts.Field_4 != 0 {
		selection, _ := alloc.Make([]int32{}, int(opts.Count)+1)
		for i := range selection {
			selection[i] = -1
		}
		opts.Field_12 = memmap.PtrWord(unsafe.Pointer(unsafe.SliceData(selection)))
	}
	if opts.Field_3 != 0 {
		childFlags := flags&0xffffefef | 9
		offset, available := 0, height
		if label {
			offset = fh + 1
			available -= fh + 1
		}
		images := flags&128 != 0
		buttonHeight := 10
		if images {
			buttonHeight = 13
		}
		button := func(up bool) *gui.Window {
			d := gui.WindowData{Style: 1, Window: w}
			name := "Down"
			cy := offset + available - buttonHeight
			if up {
				name = "Up"
				cy = offset
			}
			if images {
				d.BgImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultLB" + name + "Button").C())
				d.HlImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultLB" + name + "ButtonLit").C())
				d.DisImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultLB" + name + "ButtonDis").C())
				d.SelImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultLB" + name + "ButtonLit").C())
			} else {
				d.BgColorVal = Get_nox_color_black_2650656()
				d.DisColorVal = d.BgColorVal
				d.EnColorVal = Get_nox_color_orange_2614256()
				d.TextColorVal = d.EnColorVal
				d.HlColorVal = Get_nox_color_white_2523948()
				d.SelColorVal = Get_nox_color_yellow_2589772()
				d.SetText(GetServer().S().Strings().GetStringInFile(strman.ID("WindowDir:"+name), "C:\\NoxPost\\src\\Client\\Gui\\Gadgets\\listbox.c"))
			}
			return NewButtonOrCheckbox(w, childFlags, width-10, cy, 10, buttonHeight, &d)
		}
		opts.Field_7 = button(true).C()
		opts.Field_8 = button(false).C()
		d := gui.WindowData{Style: 8, Window: w}
		sliderWidth := 10
		if images {
			sliderWidth = 9
			d.BgImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultSliderThumb").C())
			d.HlImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultSliderThumbLit").C())
			d.DisImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultSliderThumbDis").C())
			d.SelImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg("DefaultSliderThumbLit").C())
		} else {
			d.BgColorVal = Get_nox_color_black_2650656()
			d.DisColorVal = d.BgColorVal
			d.HlColorVal = d.BgColorVal
			d.EnColorVal = Get_nox_color_orange_2614256()
			d.SelColorVal = d.EnColorVal
		}
		opts.Field_9 = uiSliderNew(w, childFlags, width-sliderWidth, offset+buttonHeight, sliderWidth, available-2*buttonHeight, &d, &gui.SliderData{}).C()
	}
	data, _ := alloc.New(gui.ScrollListBoxData{})
	*data = *opts
	w.WidgetData = unsafe.Pointer(data)
	return w
}
