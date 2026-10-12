package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func serverPanelsResource(off uintptr) string {
	lang := nox_strman_get_lang_code()
	if GetClient().R2().FontHeight(nil) > 10 {
		lang = 2
	}
	return alloc.GoString((*byte)(memmap.LoadPtr(0x587000, off+4*uintptr(lang))))
}
func serverPanelsGeneralOpen(parent *gui.Window) int {
	if *serverPanelsWord(1309812) != 0 {
		return 0
	}
	w := Nox_new_window_from_file(serverPanelsResource(173556), serverPanelsGeneralProc)
	*serverPanelsWord(1309812) = uint32(serverOptionsPtr(w))
	if w == nil {
		return 0
	}
	w.SetParent(parent)
	w.SetDraw(func(*gui.Window, *gui.WindowData) int { return 1 })
	if noxflags.HasGame(1056) {
		uiWindowEnable(w.ChildByID(10306), 0)
	}
	serverPanelsGeneralRefresh()
	w.ChildByID(10319).SetHidden(true)
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		uiWindowEnable(w.ChildByID(10304), 0)
	}
	return int(serverOptionsPtr(w))
}
func serverPanelsGeneralRefresh() int {
	w := serverPanelsWindow(1309812)
	if w == nil {
		return 0
	}
	if Nox_server_doPlayersAutoRespawn_40A5F0() != 0 {
		child := w.ChildByID(10301)
		child.DrawData().Field0 |= 4
		if noxflags.HasGame(1024) {
			uiWindowEnable(child, 0)
		}
	}
	if Get_nox_server_sendMotd_108752() != 0 {
		w.ChildByID(10302).DrawData().Field0 |= 4
	}
	if mapCycleEnabled() != 0 {
		w.ChildByID(10304).DrawData().Field0 |= 4
	}
	if int32(serverConfigFlagsQuery(int32(2))) != 0 {
		w.ChildByID(10305).DrawData().Field0 |= 4
	}
	if int32(serverConfigFlagsQuery(int32(0x2000))) != 0 {
		w.ChildByID(10306).DrawData().Field0 |= 4
	}
	return int(serverOptionsPtr(w))
}
func serverPanelsGeneralProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsGeneralEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsGeneralEvent(_ *gui.Window, event int, arg uintptr, value int) int {
	w := serverPanelsWindow(1309812)
	child := (*gui.Window)(unsafe.Pointer(arg))
	switch event {
	case 16391:
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		switch child.ID() {
		case 10301:
			Nox_xxx_ruleSetNoRespawn_40A5E0(bool2int(Nox_server_doPlayersAutoRespawn_40A5F0() == 0))
		case 10302:
			Set_nox_server_sendMotd_108752(Get_nox_server_sendMotd_108752() ^ 1)
		case 10304:
			mapCycleSetEnabled(uint32(bool2int(mapCycleEnabled() == 0)))
		case 10305:
			serverConfigFlagsToggle(int32(2))
		case 10306:
			serverConfigFlagsToggle(int32(0x2000))
			if int32(serverConfigFlagsQuery(int32(0x2000))) == 0 {
				gameplayReportResetAll()
			}
		case 10316:
			list := w.ChildByID(10317)
			list.SetHidden(false)
			list.ShowModal()
			list.Capture(true)
		case 10319:
			serverPanelsAdvancedServerOpen()
		}
	case 16393:
		Nox_xxx_rateUpdate_40A6D0(4 - value)
	case 16400:
		if child.ID() != 10317 {
			return 0
		}
		selection := teamUIEvent(child, 16404, 0, 0)
		if selection < 0 || selection >= int(int16(uiListData(child).Field_11_0)) {
			return 0
		}
		text := teamUIEvent(child, 16406, uintptr(value), 0)
		teamUIEvent(w.ChildByID(10316), 16385, uintptr(text), ^uintptr(0))
		Set_nox_server_connectionType_3596(selection + 1)
		Nox_xxx_rateUpdate_40A6D0(int(int32(serverConfigConnectionRate(int32(selection + 1)))))
		teamUIEvent(w.ChildByID(10312), 16394, uintptr(4-int(int32(serverConfigRateGet()))), 0)
		child.SetHidden(true)
		child.Capture(false)
	}
	return 0
}
func serverPanelsGeneralClose() int {
	serverPanelsWindow(1309812).Destroy()
	*serverPanelsWord(1309812) = 0
	return serverPanelsAdvancedServerClose()
}
