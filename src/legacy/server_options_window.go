package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func serverOptionsDispatch(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverOptionsEvent(w, e.EventCode(), a, int(b)))
}
func serverOptionsConstruct() int {
	if Nox_gui_xxx_check_446360() != 0 {
		return 1
	}
	if serverOptionsRoot != 0 {
		Nox_xxx_clientPlaySoundSpecial_452D80(231, 100)
		serverOptionsTryClose()
		return 1
	}
	if noxflags.HasGame(1) {
		serverOptionsDirty(1)
		if serverOptionsFirstOpen != 0 {
			listClear((*legacyListNode)(serverOptionsListHead()))
		}
	}
	lang := nox_strman_get_lang_code()
	r := GetClient().R2()
	if r.FontHeight(nil) > 10 {
		lang = 2
	}
	name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 129760+4*uintptr(lang))))
	w := Nox_new_window_from_file(name, serverOptionsDispatch)
	serverOptionsRoot = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	r.SetTabWidth(100)
	w.SetPos(image.Pt(int(nox_win_width)-w.SizeVal.X-10, 0))
	w.SetFunc93(func(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
		a, b := e.EventArgsC()
		return gui.RawEventResp(serverOptionsKey(w, e.EventCode(), int(a), int(b)))
	})
	w.SetDraw(serverOptionsDraw)
	for _, v := range []struct {
		off uintptr
		id  uint
	}{{1046512, 10101}, {1046496, 10114}, {1046500, 10183}, {1046504, 10197}, {1046508, 10199}, {1046524, 10150}, {1046516, 10134}, {1046520, 10135}, {1046536, 10153}} {
		*serverOptionsWord(v.off) = memmap.PtrWord(w.ChildByID(v.id).C())
	}
	w.ChildByID(10331).SetTooltipFunc(finalTooltipKey(tooltipTeamAssign))
	w.ChildByID(10333).SetTooltipFunc(finalTooltipKey(tooltipTeamDamage))
	for _, off := range []uintptr{1046524, 1046532, 1046536, 1046500, 1046504, 1046508} {
		if child := serverOptionsWindow(off); child != nil {
			child.DrawData().Window = w
			child.SetFunc94(serverOptionsDispatch)
		}
	}
	*serverOptionsWord(1046352) = uiMeterLoadImage("UITabs1")
	serverOptionsTabs2 = uiMeterLoadImage("UITabs2")
	serverOptionsTabs3 = uiMeterLoadImage("UITabs3")
	list := serverOptionsWindow(1046496)
	data := (*gui.ScrollListBoxData)(list.WidgetData)
	en := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISlider")))
	lit := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISliderLit")))
	slider, up, down := w.ChildByID(10182), w.ChildByID(10180), w.ChildByID(10181)
	slider.Field100Ptr.SizeVal = image.Pt(16, 10)
	gui.ButtonSetImage(slider, nil, nil, en, lit, lit)
	slider.DrawData().Window = list
	up.DrawData().Window = list
	down.DrawData().Window = list
	data.Field_9 = slider.C()
	data.Field_7 = up.C()
	data.Field_8 = down.C()
	w.ChildByID(10160).DrawData().Field0 |= 4
	initial := serverOptionsRecord(unsafe.Pointer((*int8)(unsafe.Pointer(serverConfigSlotSelect(int32(0))))))
	if noxflags.HasGame(1) {
		Sub_4161E0()
	}
	serverOptionsLoadModes()
	serverOptionsSetup(initial)
	serverOptionsMeasure()
	if noxflags.HasGame(1) {
		serverConfigSlotCopy(int32(0), int32(1))
		serverConfigSlotSelect(int32(1))
	} else {
		tab := w.ChildByID(10159)
		serverOptionsSetText(tab, 16385, serverOptionsText("servopts.wnd:teams"), -1)
		tab.DrawData().SetTooltip(GetServer().S().Strings(), serverOptionsText("servopts.wnd:TeamTT"))
		w.ChildByID(10149).SetHidden(false)
		serverOptionsWindow(1046508).SetHidden(true)
		serverOptionsHideRange(w, 10145, 10146, true)
	}
	if serverOptionsFirstOpen != 0 {
		if Sub_4D6F30() != 0 || questRuntimeWord(1556160) != 0 {
			nox_server_parseCmdText_443C80(internWStr("execrul OTQuest.rul"), 1)
		} else if noxflags.HasEngine(noxflags.EngineNoRendering) {
			nox_server_parseCmdText_443C80(internWStr("execrul server.rul"), 1)
		}
	}
	serverOptionsFirstOpen = 0
	return 1
}
func serverOptionsClose(clearRules int) uintptr {
	if serverOptionsRoot != 0 {
		w := serverOptionsWindow(1046492)
		g := w.GUI()
		if w.IsChild(g.Focused()) {
			g.Focus(nil)
		}
		w.Destroy()
		serverOptionsRoot = 0
		teamUIPlayersDestroy(false)
		serverPanelsAdvancedServerClose()
		serverOptionsPlayersPanel = 0
		serverPanelsAccessClose(false)
		serverOptionsAccessPanel = 0
		serverPanelsGeneralClose()
		serverOptionsGeneralPanel = 0
		serverOptionsAdvancedControl = 0
	}
	if clearRules != 0 {
		p := runtimeRejectedClear(serverOptionsListHead())
		serverOptionsFirstOpen = 1
		return p
	}
	return 0
}
func serverOptionsTryClose() int {
	if serverOptionsRoot == 0 {
		return 0
	}
	serverOptionsClose(0)
	serverOptionsRoot = 0
	return 1
}
func serverOptionsKey(_ *gui.Window, event, key, state int) int {
	if event == 21 {
		if key != 1 {
			return 0
		}
		if state == 2 {
			Nox_xxx_clientPlaySoundSpecial_452D80(231, 100)
			serverOptionsClose(0)
		}
	}
	return 1
}
func serverOptionsDraw(w *gui.Window, _ *gui.WindowData) int {
	p := uiWindowPosition(w)
	GetClient().R2().DrawRectFilledAlpha(p.X, p.Y+25, w.SizeVal.X, w.SizeVal.Y-25)
	list := serverOptionsChild(10120)
	if !list.Flags.IsHidden() && !serverOptionsPointIn(list, GetClient().GetMousePos()) {
		list.Capture(false)
		list.SetHidden(true)
	}
	return 1
}

// The legacy hit test uses absolute position and live size, not cached EndPos.
func serverOptionsPointIn(w *gui.Window, p image.Point) bool {
	off := uiWindowPosition(w)
	return p.X >= off.X && p.X <= off.X+w.SizeVal.X && p.Y >= off.Y && p.Y <= off.Y+w.SizeVal.Y
}
