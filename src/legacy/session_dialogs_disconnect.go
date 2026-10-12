package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var sessionDisconnectRoot, sessionDisconnectIcon *gui.Window

func sessionDisconnectOpen() int {
	img := Nox_xxx_gLoadImg("DisconnectIcon")
	width, height := int(nox_win_width), int(nox_win_height)
	icon := GetClient().Cli().GUI.NewWindowRaw(nil, 136, width-50, height/2+3, 50, 50, nil)
	sessionDisconnectIcon = icon
	icon.DrawData().BgImageHnd = img.C()
	icon.SetAllFuncs(nil, sessionDisconnectDraw, nil)
	w := Nox_new_window_from_file("discon.wnd", sessionDisconnectEvent)
	sessionDisconnectRoot = w
	w.SetFunc93(sessionDisconnectInput)
	w.SetParent(nil)
	// The legacy calculation reads EndPos (offsets24/28), including the resource
	// position. Preserve that placement for resources whose initial position is nonzero.
	w.SetPos(image.Pt(width/2-int(uint32(w.EndPos.X)/2), height/2-int(uint32(w.EndPos.Y)/2)))
	return 1
}
func sessionDisconnectInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() != 21 {
		return nil
	}
	a, _ := ev.EventArgsC()
	if a == 1 {
		return gui.RawEventResp(1)
	}
	if a == 57 {
		p := GetClient().GetMousePos()
		w.Func93(gui.AsWindowEvent(5, uintptr(uint32(p.X)|uint32(p.Y)<<16), 0))
	}
	return nil
}
func sessionDisconnectEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 23 {
		return gui.RawEventResp(1)
	}
	if ev.EventCode() != 16391 {
		return nil
	}
	a, _ := ev.EventArgsC()
	button := (*gui.Window)(unsafe.Pointer(a))
	switch button.ID() {
	case 576:
		Sub_43CF40()
	case 577:
		Sub_446380()
		if Get_dword_5d4594_2650652() != 0 && Sub_41E2F0() == 9 {
			Sub_41F4B0()
			Sub_41EC30()
			sessionMOTDFree(0)
			Nox_xxx____setargv_4_44B000()
		} else {
			browserNotice(true)
		}
		sessionDisconnectShow(0)
	}
	return nil
}
func sessionDisconnectDraw(w *gui.Window, d *gui.WindowData) int {
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(w.DrawData().BgImageHnd), uiWindowPosition(w).Add(w.DrawData().ImgPtVal))
	return 1
}
func sessionDisconnectClose() int {
	sessionDisconnectRoot.Destroy()
	sessionDisconnectIcon.Destroy()
	sessionDisconnectRoot = nil
	sessionDisconnectIcon = nil
	return 0
}
func sessionDisconnectIconShow(show int) int {
	if sessionDisconnectIcon == nil {
		return -2
	}
	sessionDisconnectIcon.SetHidden(show == 0)
	return 0
}
func sessionDisconnectShow(show int) int {
	w := sessionDisconnectRoot
	if show != 0 {
		GetClient().Nox_video_stopAllFades44E040()
		w.Show()
		w.ShowModal()
		w.StackPush()
		w.Focus()
		return uiWindowEnable(w, 1)
	}
	w.Hide()
	w.StackPop()
	GetClient().Cli().GUI.Focus(nil)
	return uiWindowEnable(w, 0)
}
