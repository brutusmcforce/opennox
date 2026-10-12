package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var uiMeterCallbackIdentities [11]byte

func uiMeterCallbackIdentity(op int) unsafe.Pointer {
	i := -1
	switch op {
	case 15:
		i = 0
	case 17:
		i = 1
	case 18:
		i = 2
	case 19:
		i = 3
	case 20:
		i = 4
	case 21:
		i = 5
	case 22:
		i = 6
	case 25:
		i = 7
	case 33:
		i = 8
	case 36:
		i = 9
	}
	if i < 0 {
		return nil
	}
	return unsafe.Pointer(&uiMeterCallbackIdentities[i])
}

func uiMeterInventoryWeaponDrawIdentity() unsafe.Pointer {
	return unsafe.Pointer(&uiMeterCallbackIdentities[10])
}

func init() {
	gui.RegisterTooltipCallbackGo(uiMeterCallbackIdentity(33), func(*gui.Window, *gui.WindowData, uintptr) { uiMeterWeaponTooltip() })
}

func uiMeterEventResult(result int) gui.WindowEventResp {
	if result == 0 {
		return nil
	}
	return gui.RawEventResp(uintptr(uint32(result)))
}

func uiMeterEventCode(e gui.WindowEvent) int {
	code := e.EventCode()
	a1, a2 := e.EventArgsC()
	_, _ = a1, a2 // The original C trampoline evaluates both ignored words.
	return int(int32(uint32(code)))
}

func uiMeterWeaponInputCallback(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	code := uiMeterEventCode(e)
	wptr := int(int32(memmap.PtrWord(w.C())))
	return uiMeterEventResult(sub_470E90(wptr, code))
}
func uiMeterPotionInputCallback(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	code := uiMeterEventCode(e)
	wptr := int(int32(memmap.PtrWord(w.C())))
	return uiMeterEventResult(nox_xxx_guiBottleSlotProc_471B90(wptr, code))
}
func uiMeterTubeInputCallback(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	code := uiMeterEventCode(e)
	wptr := int(int32(memmap.PtrWord(w.C())))
	return uiMeterEventResult(nox_xxx_guiHealthManaTubeProc_472100(wptr, code))
}

func uiMeterWeaponDrawCallback(w *gui.Window, _ *gui.WindowData) int { return uiMeterWeaponDraw(w) }
func uiMeterPotionDrawCallback(w *gui.Window, _ *gui.WindowData) int { return uiMeterPotionDraw(w) }
func uiMeterLabelDrawCallback(w *gui.Window, _ *gui.WindowData) int  { return uiMeterLabel(w) }
func uiMeterChargeRasterDrawCallback(w *gui.Window, _ *gui.WindowData) int {
	return uiMeterChargeRaster(w)
}
func uiMeterMiniBarDrawCallback(w *gui.Window, _ *gui.WindowData) int { return uiMeterMiniBar(w) }
func uiMeterTubeDrawCallback(w *gui.Window, _ *gui.WindowData) int    { return uiMeterTube(w) }
