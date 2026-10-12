package legacy

import (
	"fmt"
	"image"
	"unicode/utf16"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func (e bindingEditor) mainEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	fn := Sub_4C3A90
	if e {
		fn = Sub_4CBE70
	}
	return gui.RawEventResp(fn(int(uintptr(w.C())), ev.EventCode(), unsafe.Pointer(a), int(b)))
}
func (e bindingEditor) route(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	root, modal, base, selected, buffer := e.offsets()
	a, _ := ev.EventArgsC()
	switch ev.EventCode() {
	case 23:
		return gui.RawEventResp(1)
	case 16393:
		uiListEvent(w, ev)
		index := uiListIndex(bindingList(w))
		for i := 1; i < 4; i++ {
			bindingSend(bindingWindow(base+4*i), 16412, uintptr(index), 0)
		}
	case 16400:
		dst := (*gui.Window)(unsafe.Pointer(a))
		row := int(int32(bindingList(dst).Field_12))
		if row >= 0 {
			bindingStore(selected, dst)
			text := fmt.Sprintf("%s\n'%s'", GetClient().Strings().GetString("InputCfg.wnd:PressKey"), bindingGetText(bindingWindow(base+4), row))
			alloc.StrCopy16(unsafe.Slice(memmap.PtrUint16(0x5D4594, uintptr(buffer)), len(utf16.Encode([]rune(text)))+1), text)
			m := bindingWindow(modal)
			bindingSend(m.ChildByID(981), 16385, uintptr(memmap.PtrOff(0x5D4594, uintptr(buffer))), 0)
			m.ShowModal()
			m.Focus()
			m.StackPush()
		}
	case 16384, 16391:
		r := bindingWindow(root)
		if a == uintptr(r.ChildByID(921).C()) || a == uintptr(r.ChildByID(922).C()) {
			for i := 1; i < 4; i++ {
				bindingSend(bindingWindow(base+4*i), ev.EventCode(), a, 0)
			}
		}
	}
	return uiListEvent(w, ev)
}
func (e bindingEditor) construct() int {
	root, modal, base, _, buffer := e.offsets()
	if e {
		GetClient().GameAddStateCode(900)
	}
	w := Nox_new_window_from_file("InputCfg.wnd", e.mainEvent)
	bindingStore(root, w)
	if w == nil {
		return 0
	}
	if e {
		w.SetFunc93(func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			a, b := ev.EventArgsC()
			return gui.RawEventResp(Sub_4A18E0(w, ev.EventCode(), int(a), int(b)))
		})
		anim := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -480, 0, 20, 0, -40)
		legacyGlobals.nox_wnd_xxx_1522608 = anim
		if anim == nil {
			return 0
		}
		anim.StateID = 900
		anim.Func12Ptr = animationCallbackKey(animationKeyBindingsBack)
		anim.FncDoneOutPtr = animationCallbackKey(animationKeyBindingsDone)
	}
	for i := 0; i < 4; i++ {
		bindingStore(base+4*i, w.ChildByID(uint(910+i)))
	}
	spacer := bindingWindow(base)
	if spacer == nil {
		return 0
	}
	data := bindingList(spacer)
	(*gui.Window)(data.Field_7).SetID(921)
	(*gui.Window)(data.Field_8).SetID(922)
	(*gui.Window)(data.Field_9).SetID(920)
	spacer.SetFunc94(e.route)
	for i := 1; i < 4; i++ {
		bindingWindow(base + 4*i).SetParent(spacer)
	}
	bindingWindow(base + 8).SetFunc93(bindingFilter)
	bindingWindow(base + 12).SetFunc93(bindingFilter)
	for _, spec := range [][2]int{{921, 16408}, {922, 16409}} {
		button := w.ChildByID(uint(spec[0]))
		for i := 1; i < 4; i++ {
			bindingSend(bindingWindow(base+4*i), spec[1], uintptr(button.C()), 0)
		}
	}
	if e {
		Sub_4CBBF0()
	}
	for id, end := 971, 971+int(memmap.Uint8(0x5D4594, 1193128)); id < end; id++ {
		uiWindowEnable(w.ChildByID(uint(id)), 1)
	}
	bindingSend(w.ChildByID(uint(971+Nox_client_mousePriKey_430AF0())), 16392, 1, 0)
	// Preserve unsigned C centering, including endpoint normalization in SetPos.
	if !e {
		w.SetPos(image.Pt(int((uint32(nox_win_width)-uint32(w.SizeVal.X))/2), 0))
	}
	m := w.ChildByID(980)
	bindingStore(modal, m)
	m.SetParent(nil)
	m.SetFunc94(e.mainEvent)
	m.SetFunc93(e.modalEvent)
	m.SetHidden(true)
	if !e {
		m.SetPos(image.Pt(int((uint32(nox_win_width)-uint32(m.SizeVal.X))/2), m.Off.Y))
	}
	bindingSend(m.ChildByID(981), 16385, uintptr(memmap.PtrOff(0x5D4594, uintptr(buffer))), 0)
	if !e {
		uiWindowEnable(w.ChildByID(932), 1)
		w.SetHidden(true)
	}
	return 1
}
func bindingDestroy() int {
	bindingWindow(1321228).Destroy()
	bindingWindow(1321232).Destroy()
	for _, off := range []int{1321228, 1321232, 1321236, 1321240, 1321244, 1321248} {
		*bindingSlot(off) = nil
	}
	return 0
}
func bindingShow() { Sub_4C3B70(); bindingWindow(1321228).ShowModal(); Sub_413A00(1) }
func bindingVisible() int {
	w := bindingWindow(1321228)
	if w != nil && !w.Flags.IsHidden() {
		return 1
	}
	return 0
}
func bindingClose(cancel int) int {
	w := bindingWindow(1321228)
	if w == nil || w.Flags.IsHidden() {
		return 0
	}
	Sub_413A00(0)
	if cancel != 0 {
		w.SetHidden(true)
	} else {
		bindingInGame.apply()
		WriteConfigLegacy("nox.cfg")
		w.SetHidden(true)
		uiMeterBindings()
		optionsShow()
	}
	return 1
}
func bindingMenuBack() int {
	bindingMenu.apply()
	WriteConfigLegacy("nox.cfg")
	a := legacyGlobals.nox_wnd_xxx_1522608
	a.SetState(gui.AnimOut)
	Sub_43BE40(2)
	Nox_xxx_clientPlaySoundSpecial_452D80(923, 100)
	return 1
}
func bindingMenuDone() int {
	a := legacyGlobals.nox_wnd_xxx_1522608
	fn := a.Func13Ptr
	a.Free()
	bindingWindow(1522604).Destroy()
	bindingWindow(1522612).Destroy()
	gui.CallAnimationCallback(fn)
	return 1
}
func bindingYesNo() *gui.Window {
	w := Nox_new_window_from_file("yesno.wnd", nil)
	bindingStore(1321224, w)
	if w != nil {
		w.SetPos(image.Pt((int(nox_win_width)-320)/2, (int(nox_win_height)-240)/2))
	}
	w.SetHidden(true)
	return *bindingSlot(1321224)
}
