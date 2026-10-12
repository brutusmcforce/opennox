package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func interactionKeyRoot() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(interactionKeyRootWord))
}
func interactionSetKeyState(v uint32) {
	if interactionKeyState == 0 && v == 1 {
		Nox_xxx_clientPlaySoundSpecial_452D80(1022, 100)
	}
	interactionKeyState = uint32(v)
}
func interactionKeyShow() int {
	if w := interactionKeyRoot(); w != nil {
		w.Show()
		uiWindowEnable(w, 1)
		w.SetPos(image.Pt(int(nox_win_width)/2-w.SizeVal.X/2, int(nox_win_height)/2-w.SizeVal.Y/2))
		GetClient().Cli().GUI.Focus(nil)
	}
	return 0
}
func interactionKeyHide() int {
	if w := interactionKeyRoot(); w != nil {
		w.Hide()
		GetClient().Cli().GUI.Focus(nil)
	}
	return 0
}
func interactionKeyUpdate(value uintptr) {
	state := uint32(interactionKeyState)
	if state != 0 {
		if state == 1 && value == 0 {
			interactionKeyHide()
			interactionSetKeyState(0)
		}
	} else if value == 1 {
		interactionKeyShow()
		interactionSetKeyState(1)
	}
}
func interactionKeyEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16391 {
		a, _ := ev.EventArgsC()
		id := -2
		if child := (*gui.Window)(unsafe.Pointer(a)); child != nil {
			id = int(child.ID())
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		if id == 10803 {
			interactionKeyHide()
		}
	}
	return nil
}
func interactionKeyOpen() int {
	w := Nox_new_window_from_file("SKey.wnd", interactionKeyEvent)
	interactionKeyRootWord = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	interactionSetKeyState(0)
	interactionKeyHide()
	return 1
}
func interactionKeyDestroy() {
	nox_xxx_windowDestroyMB_46C4E0((*nox_window)(interactionKeyRoot().C()))
	interactionKeyRootWord = 0
	interactionSetKeyState(0)
}
func interactionVoteOpen() int {
	img := Nox_xxx_gLoadImg("VoteInProgress")
	w := GetClient().Cli().GUI.NewWindowRaw(nil, 136, int(nox_win_width)-50, int(nox_win_height)/2-100, 50, 50, nil)
	interactionVoteIcon = memmap.PtrWord(w.C())
	w.DrawData().BgImageHnd = img.C()
	w.SetAllFuncs(nil, interactionIconDraw, nil)
	w.Hide()
	return 1
}
func interactionVoteHide(v int) int {
	return nox_window_set_hidden((*nox_window)(memmap.WordPtr(interactionVoteIcon)), v)
}
func interactionVoteDestroy() int {
	result := nox_xxx_windowDestroyMB_46C4E0((*nox_window)(memmap.WordPtr(interactionVoteIcon)))
	interactionVoteIcon = 0
	return result
}

func sub_4BFB70(v int32) { interactionSetKeyState(uint32(v)) }

func sub_4BFBB0(v uint32) { interactionKeyUpdate(uintptr(v)) }

func sub_4BFBF0() int32 { return int32(interactionKeyShow()) }

func sub_4BFC70() int32 { return int32(interactionKeyHide()) }

func sub_4BFC90() int32 { return int32(interactionKeyOpen()) }

func sub_4BFD10() { interactionKeyDestroy() }

func sub_4BFD30() int32 { return int32(interactionKeyState) }

func sub_4C3390() int32 { return int32(interactionVoteOpen()) }

func sub_4C3410(w *int32) int32 {
	return int32(interactionIconDraw((*gui.Window)(unsafe.Pointer(w)), nil))
}

func sub_4C3460(v int32) int32 { return int32(interactionVoteHide(int(v))) }

func sub_4C34A0() int32 { return int32(interactionVoteDestroy()) }
