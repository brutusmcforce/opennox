package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func (e optionsEditor) event(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	child := (*gui.Window)(unsafe.Pointer(a))
	value := int(int32(b))
	switch ev.EventCode() {
	case 16389:
		Nox_xxx_clientPlaySoundSpecial_452D80(920, 100)
		return gui.RawEventResp(1)
	case 16391:
		id := int(child.ID())
		if e && id >= 380 {
			return gui.RawEventResp(Nox_gui_menu_proc_ext(id))
		}
		switch {
		case id >= 311 && id <= 314:
			cut := [...]int{65, 75, 85, 100}[id-311]
			if e {
				Nox_video_setCutSize_4766A0(cut)
			} else {
				Nox_draw_setCutSize_476700(cut, 0)
			}
		case bool(e) && (id == 331 || id == 332):
			if id == 331 {
				*optionsWord(172880) = 8
				Nox_video_setFullScreen(0)
			} else {
				*optionsWord(172880) = 16
				Nox_video_setFullScreen(1)
			}
		case bool(e) && id == 333:
			Nox_video_setFullScreen(1)
		case bool(e) && id == 334:
			Nox_video_setFullScreen(0)
		case id == 341:
			if e {
				Sub_4AA9C0()
				anim := Get_nox_wnd_xxx_1309740()
				anim.FncDoneOutPtr = animationCallbackKey(animationKeyOptionsDone)
				anim.Func13Ptr = animationCallbackKey(animationKeyBindingsShow)
			} else {
				optionsClose(0)
				Sub_445C40()
				bindingShow()
			}
		case id >= 361 && id <= 363:
			optionsToggle(id - 361)
		case !bool(e) && id == 371:
			optionsClose(0)
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
		return gui.RawEventResp(1)
	case 16393, 16396:
		id := int(child.ID())
		if id >= 351 && id <= 353 {
			e.volume(child, ev.EventCode(), id-351, value)
		} else if ev.EventCode() == 16393 {
			switch id {
			case 316:
				Nox_video_setGammaSlider(value)
			}
		}
	}
	return gui.RawEventResp(0)
}
func optionsTabWidth() int {
	r := GetClient().R2()
	old := r.TabWidth()
	r.SetTabWidth(15)
	return old
}
func (e optionsEditor) construct() int {
	if e {
		GetClient().GameAddStateCode(300)
	}
	w := Nox_new_window_from_file("Options.wnd", e.event)
	optionsStore(e.rootOffset(), w)
	if w == nil {
		return 0
	}
	GetClient().NewGUIAdvOptsOn(w)
	if e {
		w.SetFunc93(func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			a, b := ev.EventArgsC()
			return gui.RawEventResp(Sub_4A18E0(w, ev.EventCode(), int(a), int(b)))
		})
		optionsTabWidth()
		anim := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -480, 0, 20, 0, -40)
		legacyGlobals.nox_wnd_xxx_1309740 = anim
		if anim == nil {
			return 0
		}
		anim.StateID = 300
		anim.Func12Ptr = animationCallbackKey(animationKeyOptionsStart)
		anim.FncDoneOutPtr = animationCallbackKey(animationKeyOptionsFinish)
	}
	for ch := 0; ch < 3; ch++ {
		slider := w.ChildByID(uint(351 + ch))
		slider.Field100Ptr.SizeVal = image.Pt(24, 20)
		hl := Nox_xxx_gLoadImg("OptionsVolumeSliderLit")
		sel := Nox_xxx_gLoadImg("OptionsVolumeSliderLit")
		en := Nox_xxx_gLoadImg("OptionsVolumeSlider")
		gui.ButtonSetImage(slider, nil, nil, noxrender.ImageHandle(en.C()), noxrender.ImageHandle(sel.C()), noxrender.ImageHandle(hl.C()))
		optionsSend(slider, 16395, 0, 0x4000)
		optionsSend(slider, 16394, uintptr(optionsTimer(ch).Current>>16), 0)
		check := w.ChildByID(uint(361 + ch))
		optionsStore(e.buttonOffset(ch), check)
		if optionsEnabled(ch) == 1 {
			check.DrawData().Field0 |= 4
		} else {
			check.DrawData().Field0 &^= 4
		}
	}
	if e {
		Sub_4A19F0("OptsBack.wnd:Back")
		Sub_4A1A40(0)
		optionsMenuRefresh()
		return 1
	}
	Nox_video_setMenuOptions(w)
	id := uint(334)
	if Nox_video_getFullScreen() != 0 {
		id = 333
	}
	optionsSend(w.ChildByID(id), 16392, 1, 0)
	for id := uint(320); id <= 332; id++ {
		c := w.ChildByID(id)
		if c != nil && c.DrawData().Field0&4 == 0 {
			uiWindowEnable(c, 0)
		}
	}
	g := GetClient().Cli().GUI
	overlay := g.NewWindowRaw(nil, 32, 0, 0, 1, 1, nil)
	optionsStore(1309824, overlay)
	for off := uintptr(174072); ; off += 16 {
		p := unsafe.Slice(memmap.PtrInt32(0x587000, off), 4)
		if p[0] == -1 {
			break
		}
		c := g.NewWindowRaw(overlay, 0, int(p[0]), int(p[1]), int(p[2]), int(p[3]), nil)
		c.SetDraw(optionsOverlayDraw)
	}
	// C converts the unsigned subtraction to signed before dividing.
	pos := image.Pt(int(int32(uint32(nox_win_width)-uint32(w.SizeVal.X)))/2, 0)
	w.SetPos(pos)
	overlay.SetPos(pos)
	w.ChildByID(371).SetHidden(false)
	w.SetHidden(true)
	overlay.SetHidden(true)
	return 1
}
func optionsMenuRefresh() int {
	w := optionsWindow(1309720)
	Nox_video_setMenuOptions(w)
	id := uint(331)
	if Nox_video_getFullScreen() != 0 {
		id = 332
	}
	optionsSend(w.ChildByID(id), 16392, 1, 0)
	cut := Nox_video_getCutSize_4766D0()
	*memmap.PtrUint32(0x587000, 172884) = uint32(cut)
	if c := w.ChildByID(optionsViewportID(cut)); c != nil {
		return optionsSend(c, 16392, 1, 0)
	}
	return 0
}
func optionsShow() int {
	Sub_413A00(1)
	w := optionsWindow(1309820)
	w.ShowModal()
	optionsWindow(1309824).ShowModal()
	GetClient().GUIAdvVideoOptsLoad()
	if c := w.ChildByID(optionsViewportID(Nox_video_getCutSize_4766D0())); c != nil {
		optionsSend(c, 16392, 0, 0)
	}
	return optionsTabWidth()
}
func optionsClose(cancel int) int {
	w := optionsWindow(1309820)
	if w == nil || w.Flags.IsHidden() {
		return 0
	}
	Sub_413A00(0)
	Dialogs.Sub_44D8F0()
	if cancel == 0 {
		WriteConfigLegacy("nox.cfg")
	}
	w.SetHidden(true)
	optionsWindow(1309824).SetHidden(true)
	if cancel == 0 {
		Sub_445C40()
	}
	return 1
}
func optionsVisible() int {
	w := optionsWindow(1309820)
	if w == nil || w.Flags.IsHidden() {
		return bindingVisible()
	}
	return 1
}
func optionsDestroy() int {
	w := optionsWindow(1309820)
	if w == nil {
		return -2
	}
	w.Destroy()
	*optionsWord(1309820) = 0
	return 0
}
func optionsMenuDone() int {
	a := Get_nox_wnd_xxx_1309740()
	fn := a.Func13Ptr
	a.Free()
	optionsWindow(1309720).Destroy()
	gui.CallAnimationCallback(fn)
	return 1
}
func optionsOverlayDraw(w *gui.Window, d *gui.WindowData) int {
	p := w.GlobalPos()
	GetClient().R2().DrawRectFilledAlpha(p.X, p.Y, w.SizeVal.X, w.SizeVal.Y)
	return 1
}
