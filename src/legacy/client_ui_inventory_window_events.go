package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func uiInventoryWindowValue(v uint32) *gui.Window { return (*gui.Window)(unsafe.Pointer(uintptr(v))) }
func uiInventoryWindowEvent(fn func(*gui.Window, int, uintptr, uintptr) int) gui.WindowFunc {
	return func(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
		a, b := e.EventArgsC()
		return gui.RawEventResp(uintptr(fn(w, e.EventCode(), a, b)))
	}
}
func uiInventorySliderValue(w *gui.Window, code int, a, b uint32) int {
	return gui.EventRespInt(w.Func94(&gui.RawEvent{Event: code, Arg1: uintptr(a), Arg2: uintptr(b)}))
}
func uiInventoryPanelEvents(w *gui.Window, event int, a, b uintptr) int {
	slider := uiInventoryWindowValue(uint32(dword_5d4594_1062508))
	if event == 16393 {
		dword_5d4594_1062512 = uint32((*gui.SliderData)(slider.WidgetData).Max - uint32(b))
		return 0
	}
	if event != 16391 {
		return 0
	}
	control := uiInventoryWindowValue(uint32(a))
	switch control.ID() {
	case 9102:
		v := int32(dword_5d4594_1062512) - 25
		if v < 0 {
			v = 0
		} else {
			v -= v % 50
		}
		dword_5d4594_1062512 = uint32(v)
		uiInventorySliderValue(slider, 16394, (*gui.SliderData)(slider.WidgetData).Max-uint32(v), 0)
		audioEventPlay(766, 100, 0, 0)
	case 9103:
		v := int32(dword_5d4594_1062512) + 50
		max := int32((*gui.SliderData)(slider.WidgetData).Max)
		if v > max {
			v = max
		} else {
			v -= v % 50
		}
		dword_5d4594_1062512 = uint32(v)
		uiInventorySliderValue(slider, 16394, uint32(max-v), 0)
		audioEventPlay(766, 100, 0, 0)
	case 9105:
		height := int32(sub_469FA0()) - 150
		if uiInventoryMode() == 5 {
			return 0
		}
		if height < 0 {
			height = 0
		}
		*memmap.PtrUint8(0x5D4594, 1049869) = 1
		dword_5d4594_1062516 = dword_5d4594_1062512
		dword_5d4594_1062512 = dword_5d4594_1062520
		uiInventorySliderValue(slider, 16395, 0, uint32(height))
		uiInventorySliderValue(slider, 16394, (*gui.SliderData)(slider.WidgetData).Max-uint32(dword_5d4594_1062512), 0)
		uiWindowSetBackgroundImage(uint32(dword_5d4594_1062528), uint32(memmap.Uint32(0x5D4594, 1049980)))
		uiWindowSetSelectedImage(uint32(dword_5d4594_1062528), uint32(memmap.Uint32(0x5D4594, 1049984)))
		uiInventoryWindowValue(uint32(dword_5d4594_1062528)).SetID(9106)
	case 9106:
		*memmap.PtrUint8(0x5D4594, 1049869) = 0
		dword_5d4594_1062520 = dword_5d4594_1062512
		dword_5d4594_1062512 = dword_5d4594_1062516
		uiInventorySliderValue(slider, 16395, 0, 850)
		uiInventorySliderValue(slider, 16394, (*gui.SliderData)(slider.WidgetData).Max-uint32(dword_5d4594_1062512), 0)
		uiWindowSetBackgroundImage(uint32(dword_5d4594_1062528), 0)
		uiWindowSetSelectedImage(uint32(dword_5d4594_1062528), uint32(dword_5d4594_1049976))
		uiInventoryWindowValue(uint32(dword_5d4594_1062528)).SetID(9105)
	case 9107:
		if uiInventoryMode() == 5 {
			return 0
		}
		*memmap.PtrUint8(0x5D4594, 1049870) = 1
		uiWindowSetBackgroundImage(uint32(dword_5d4594_1062524), 0)
		uiWindowSetSelectedImage(uint32(dword_5d4594_1062524), uint32(memmap.Uint32(0x5D4594, 1049988)))
		uiInventoryWindowValue(uint32(dword_5d4594_1062524)).SetID(9108)
		uiInventoryWindowValue(uint32(dword_5d4594_1062468)).Hide()
	case 9108:
		if uiInventoryMode() == 5 {
			return 0
		}
		*memmap.PtrUint8(0x5D4594, 1049870) = 0
		uiWindowSetBackgroundImage(uint32(dword_5d4594_1062524), uint32(dword_5d4594_1049992))
		uiWindowSetSelectedImage(uint32(dword_5d4594_1062524), uint32(dword_5d4594_1049996))
		uiInventoryWindowValue(uint32(dword_5d4594_1062524)).SetID(9107)
		uiInventoryWindowValue(uint32(dword_5d4594_1062468)).Show()
	case 9111:
		uiInventoryCloseWindow()
	}
	return 0
}
func uiInventoryToggleButton(w *gui.Window, event int, a, b uintptr) int {
	if event == 5 || event == 6 {
		return 1
	}
	if event == 7 {
		uiInventoryToggleWindow()
		return 1
	}
	return 0
}
func uiInventoryThumbEvents(w *gui.Window, event int, a, b uintptr) int {
	if uiInventoryDragged() != nil {
		return uiInventoryMainEvents(w, event, a, b)
	}
	return int(int32(gui.EventRespInt(gui.ButtonProc2(w, gui.AsWindowEvent(int(int32(event)), uintptr(uint32(int32(a))), uintptr(uint32(int32(b))))))))
}
func uiInventoryTrackEvents(w *gui.Window, event int, a, b uintptr) int {
	if uiInventoryDragged() != nil {
		return uiInventoryMainEvents(w, event, a, b)
	}
	return int(int32(gui.EventRespInt(uiSliderInput(w, &gui.RawEvent{Event: int(int32(event)), Arg1: uintptr(uint32(a)), Arg2: uintptr(uint32(int32(b)))}, false))))
}
func uiInventoryWindowAdmission(w *gui.Window, event int, a, b uintptr) int {
	return bool2int(event != 8 && event != 12 && event != 16)
}
