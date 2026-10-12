package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func interactionIconDraw(w *gui.Window, _ *gui.WindowData) int {
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(w.DrawData().BgImageHnd), uiWindowPosition(w).Add(w.DrawData().ImgPtVal))
	return 1
}
func interactionChatIconOpen() int {
	img := Nox_xxx_gLoadImg("ChatIcon")
	*memmap.PtrUint32(0x5D4594, 825748) = uint32(uintptr(img.C()))
	w := GetClient().Cli().GUI.NewWindowRaw(nil, 136, int(nox_win_width)-50, int(nox_win_height)/2-50, 50, 50, nil)
	interactionChatIcon = uint32(uintptr(w.C()))
	w.DrawData().BgImageHnd = img.C()
	w.SetAllFuncs(nil, interactionIconDraw, nil)
	sm := GetClient().Cli().Strings()
	w.DrawData().SetTooltip(sm, sm.GetStringInFile("chatmode", "chaticon.c"))
	return 1
}
func interactionChatIconHidden(hidden int) int {
	return nox_window_set_hidden((*nox_window)(unsafe.Pointer(uintptr(interactionChatIcon))), hidden)
}
func interactionChatIconDestroy() int {
	result := nox_xxx_windowDestroyMB_46C4E0((*nox_window)(unsafe.Pointer(uintptr(interactionChatIcon))))
	interactionChatIcon = 0
	return result
}
func interactionObserverIconDraw(w *gui.Window, d *gui.WindowData) int {
	if Get_dword_8531A0_2576() != nil {
		root := (*gui.Window)(unsafe.Pointer(uintptr(interactionObserverIcon)))
		sm := GetClient().Cli().Strings()
		root.DrawData().SetTooltip(sm, sm.GetStringInFile("observermode", "guiobs.c"))
		interactionIconDraw(w, d)
	}
	return 1
}
func interactionObserverIconOpen() int {
	img := Nox_xxx_gLoadImg("ObserverIcon")
	*memmap.PtrUint32(0x5D4594, 1193716) = uint32(uintptr(img.C()))
	w := GetClient().Cli().GUI.NewWindowRaw(nil, 136, int(nox_win_width)-50, int(nox_win_height)/2-100, 50, 50, nil)
	interactionObserverIcon = uint32(uintptr(w.C()))
	w.DrawData().BgImageHnd = img.C()
	w.SetAllFuncs(nil, interactionObserverIconDraw, nil)
	return 1
}
func interactionObserverIconHidden(hidden int) int {
	return nox_window_set_hidden((*nox_window)(unsafe.Pointer(uintptr(interactionObserverIcon))), hidden)
}

func nox_xxx_guiChatMode_4456E0(w *int32) int32 {
	return int32(interactionIconDraw((*gui.Window)(unsafe.Pointer(w)), nil))
}

func nox_xxx_guiChatShowHide_445730(v int32) int32 { return int32(interactionChatIconHidden(int(v))) }

func sub_445770() int32 { return int32(interactionChatIconDestroy()) }

func nox_xxx_guiChatIconLoad_445650() int32 { return int32(interactionChatIconOpen()) }

func sub_48C9F0(w *int32) int32 {
	return int32(interactionObserverIconDraw((*gui.Window)(unsafe.Pointer(w)), nil))
}

func sub_48C980() int32 { return int32(interactionObserverIconOpen()) }

func nox_xxx_showObserverWindow_48CA70(v int32) int32 {
	return int32(interactionObserverIconHidden(int(v)))
}
