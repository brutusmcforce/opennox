package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func characterClassNextColor() {
	characterUI.classAnim.Func13Ptr = animationCallbackKey(animationKeyShowSelColor)
}
func characterMenuEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	return gui.RawEventResp(Sub_4A18E0(w, ev.EventCode(), int(a), int(b)))
}
func characterShowClass() int {
	questProgressReset("*:*")
	Sub_4A1BE0(1)
	characterUI.classHost = unsafe.Pointer(Nox_xxx_getHostInfoPtr_431770())
	GetClient().GameAddStateCode(600)
	w := Nox_new_window_from_file("SelClass.wnd", characterClassEvent)
	characterUI.classRoot = w
	if w == nil {
		return 0
	}
	w.SetFunc93(characterMenuEvent)
	a := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -460, 0, 20, 0, -40)
	characterUI.classAnim = a
	if a == nil {
		return 0
	}
	a.StateID = 600
	a.Func12Ptr = animationCallbackKey(animationKeyClassStart)
	a.FncDoneOutPtr = animationCallbackKey(animationKeyClassDone)
	for _, id := range []uint{601, 603, 602} {
		w.ChildByID(id).SetDraw(func(w *gui.Window, _ *gui.WindowData) int { return characterClassDraw(w) })
	}
	characterAcceptButton = w.ChildByID(610)
	*memmap.PtrUint32(0x5D4594, 1307740) = 0
	Sub_4A19F0("OptsBack.wnd:Back")
	quickbarClearSlots()
	return 1
}
func characterShowColor() int {
	GetClient().GameAddStateCode(700)
	characterUI.colorHost = unsafe.Pointer(Nox_xxx_getHostInfoPtr_431770())
	*(*byte)(unsafe.Add(characterUI.colorHost, 67)) = 0
	w := Nox_new_window_from_file("SelColor.wnd", characterColorEvent)
	characterUI.colorRoot = w
	if w == nil {
		return 0
	}
	w.SetFunc93(characterMenuEvent)
	a := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -440, 0, 20, 0, -40)
	characterUI.colorAnim = a
	if a == nil {
		return 0
	}
	a.StateID = 700
	a.Func12Ptr = animationCallbackKey(animationKeyColorStart)
	a.FncDoneOutPtr = animationCallbackKey(animationKeyColorDone)
	characterSetup()
	draw := func(w *gui.Window, _ *gui.WindowData) int { return characterPaletteDraw(w) }
	for id := uint(720); id <= 729; id++ {
		w.ChildByID(id).SetDraw(draw)
	}
	for id := uint(761); id <= 792; id++ {
		w.ChildByID(id).SetDraw(draw)
	}
	if characterUI.defaults != 0 {
		optionsSend(characterUI.controls[14], 16414, uintptr(unsafe.Pointer(characterDefaultName())), 0)
	}
	characterUI.palette = w.ChildByID(760)
	characterUI.palette.SetFunc94(characterColorEvent)
	characterUI.palette.SetFunc93(characterPaletteOutside)
	characterUI.palette.SetParent(nil)
	Sub_4BFAD0()
	w.ChildByID(740).SetDraw(func(w *gui.Window, _ *gui.WindowData) int { return characterPreview(w) })
	return 1
}
