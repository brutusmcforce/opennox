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

func teamUIPlayersConstruct(parent *gui.Window) int {
	lang := nox_strman_get_lang_code()
	r := GetClient().R2()
	if r.FontHeight(nil) > 10 {
		lang = 2
	}
	if teamUIWindow(1045684) != nil {
		return 0
	}
	name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 129048+4*uintptr(lang))))
	w := Nox_new_window_from_file(name, func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		return gui.RawEventResp(teamUIPlayersEvent(w, ev.EventCode(), (*gui.Window)(unsafe.Pointer(a)), int(b)))
	})
	*teamUIWord(1045684) = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	*teamUIWord(1045688) = memmap.PtrWord(w.ChildByID(10507).C())
	*teamUIWord(1045692) = memmap.PtrWord(w.ChildByID(10509).C())
	w.SetParent(parent)
	w.SetDraw(teamUIPlayersDraw)
	// The old wndRetNULL entry point is a no-op; it does not clear focus.
	listClear(teamUIHead(false))
	listClear(teamUIHead(true))
	en := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISlider")))
	lit := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISliderLit")))
	for i := 0; i < 2; i++ {
		list := w.ChildByID(uint(10501 + i))
		data := (*gui.ScrollListBoxData)(list.WidgetData)
		slider := w.ChildByID(uint(10517 + 3*i))
		up := w.ChildByID(uint(10515 + 3*i))
		down := w.ChildByID(uint(10516 + 3*i))
		slider.Field100Ptr.SizeVal = image.Pt(16, 10)
		gui.ButtonSetImage(slider, nil, nil, en, lit, lit)
		slider.DrawData().Window = list
		up.DrawData().Window = list
		down.DrawData().Window = list
		data.Field_9 = slider.C()
		data.Field_7 = up.C()
		data.Field_8 = down.C()
	}
	teamUIPlayersRefresh()
	if noxflags.HasGame(128) {
		teamUIEvent(w.ChildByID(10504), 16385, uintptr(unsafe.Pointer(alloc.InternCString16(teamUIText("Title1")))), ^uintptr(0))
	}
	return int(memmap.PtrWord(w.C())) // stored as a pointer word
}
func teamUISettingsLocked() bool {
	return int8(*(*byte)(unsafe.Add(unsafe.Pointer(teamUISettings()), 53))) < 0
}
func teamUIPlayersDraw(w *gui.Window, d *gui.WindowData) int {
	teamRuntimeObject(ClientPlayerNetCode())
	pos := uiWindowPosition(w)
	r := GetClient().R2()
	if !w.Flags.Has(128) {
		if d.BgColorVal != 0x80000000 {
			r.DrawRectFilledAlpha(pos.X, pos.Y, w.SizeVal.X, w.SizeVal.Y)
		}
	} else {
		bookDrawImage(memmap.PtrWord(unsafe.Pointer(d.BgImageHnd)), pos)
	}
	root := teamUIWindow(1045684)
	teams := root.ChildByID(10502)
	if noxflags.HasGame(1) && !noxflags.HasGame(0x8000) && !(noxflags.HasGame(128) && teamUISettingsLocked()) {
		enabled := 0
		if teamUIEvent(teams, 16404, 0, 0) >= 0 {
			selected := teamUIEvent(root.ChildByID(10501), 16404, 0, 0)
			if *(*int32)(unsafe.Pointer(uintptr(selected))) >= 0 { // multi-select: the selected indexes
				enabled = 1
			}
		}
		uiWindowEnable(root.ChildByID(10503), enabled)
	}
	if teamUIEvent(teams, 16404, 0, 0) < 0 {
		uiWindowEnable(teamUIWindow(1045688), 0)
		uiWindowEnable(teamUIWindow(1045692), 0)
	}
	return 1
}
