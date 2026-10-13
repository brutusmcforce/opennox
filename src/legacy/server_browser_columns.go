package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func browserColumns() []*gui.Window {
	return []*gui.Window{browserUI.playersColumn, browserUI.modeColumn, browserUI.mapColumn, browserUI.pingColumn, browserUI.statusColumn}
}
func browserColumnsSend(code int, a, b uintptr) int {
	ret := 0
	for _, w := range browserColumns() {
		ret = optionsSend(w, code, a, b)
	}
	return ret
}
func browserColumnsInit() int {
	main := browserUI.gameList
	main.SetFunc94(browserListEvent)
	main.SetFunc93(browserListInput)
	for _, w := range browserColumns() {
		w.SetParent(main)
		w.SetFunc93(browserListInput)
	}
	d := uiListData(browserUI.detailList)
	(*gui.Window)(d.Field_7).SetID(10035)
	(*gui.Window)(d.Field_8).SetID(10036)
	(*gui.Window)(d.Field_9).SetID(10032)
	d = uiListData(main)
	browserColumnsSend(16408, uintptr(d.Field_7), 0)
	browserColumnsSend(16409, uintptr(d.Field_8), 0)
	return browserColumnsSend(16410, uintptr(d.Field_9), 0)
}
func browserListInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	code := ev.EventCode()
	if code == 19 || code == 20 {
		p := uiListData(w).Field_7
		if code == 20 {
			p = uiListData(w).Field_8
		}
		optionsSend(browserUI.gameList, 16391, uintptr(p), 0)
		browserColumnsSend(16391, uintptr(p), 0)
		return nil
	}
	return uiListSingleInput(w, ev)
}
func browserListEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	code := ev.EventCode()
	a, b := ev.EventArgsC()
	switch code {
	case 23:
		return gui.RawEventResp(1)
	case 16400:
		id := windowArg(a).ID()
		if id >= 10038 && id <= 10042 {
			browserColumnsSend(16403, b, 0)
			if uint32(b) < uint32(browserUI.resultCount) {
				p := browserListAt(int32(b))
				browserUI.selected = p
				pos := GetClient().GetMousePos()
				point := [2]uint32{uint32(pos.X), uint32(pos.Y)}
				browserInfoPopup(&point, p)
			}
		}
	case 16403, 16412:
		browserColumnsSend(code, a, 0)
	case 16398, 16399:
		browserColumnsSend(code, a, b)
	case 16384:
		main := browserUI.gameList
		if windowArg(a) == main.ChildByID(10043) || windowArg(a) == main.ChildByID(10044) {
			browserColumnsSend(code, a, 0)
		}
	case 16393:
		d := uiListData(w)
		uiListEvent(w, ev)
		browserColumnsSend(16412, uintptr(uiListIndex(d)), 0)
	}
	return uiListEvent(w, ev)
}
func browserInfoPosition(x, y int32, out *[2]uint32) {
	px, py := x-65, y-20
	if px+130 > 600 {
		px = 470
	}
	if py+120 > 451 {
		py = 331
	}
	low := uint32(27)
	if browserUI.listMode == 1 {
		low = 55
	}
	if uint32(py) < low {
		py = int32(low)
	}
	if px < 216 {
		px = 216
	}
	out[0], out[1] = uint32(px), uint32(py)
}
func browserInfoPopup(point *[2]uint32, record unsafe.Pointer) {
	browserInfoPosition(int32(point[0]), int32(point[1]), point)
	w := browserUI.detailPanel
	w.SetParent(nil)
	w.ShowModal()
	w.StackPush()
	w.SetPos(image.Pt(int(int32(point[0])), int(int32(point[1]))))
	browserDetails(record)
	browserUI.hasSelection = 1
	*memmap.PtrUint16(0x5D4594, 814604) = *(*uint16)(unsafe.Pointer(&(*Nox_gui_server_ent_t)(record).PortVal))
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		uiWindowEnable(browserUI.joinButton, 0)
	}
}
