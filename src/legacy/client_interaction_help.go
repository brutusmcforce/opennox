package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func interactionHelpRoot() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(interactionHelpRootWord)))
}
func interactionHelpText(id string) string {
	return GetClient().Cli().Strings().GetStringInFile(strman.ID(id), "chathelp.c")
}
func interactionHelpClose() int {
	w := interactionHelpRoot()
	if w == nil {
		return 0
	}
	Set_nox_server_sanctuaryHelp_54276(int((^uint32(w.ChildByID(4104).DrawData().Field0) >> 2) & 1))
	w.StackPop()
	w.Capture(false)
	w.Destroy()
	interactionHelpRootWord = 0
	GetClient().Cli().GUI.Focus(nil)
	if noxflags.HasGame(1) {
		return serverOptionsVisible(0)
	}
	return 0
}
func interactionHelpEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16391 {
		arg, _ := ev.EventArgsC()
		id := -2
		if child := (*gui.Window)(unsafe.Pointer(arg)); child != nil {
			id = int(child.ID())
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		if id == 4103 {
			interactionHelpClose()
		}
	}
	return gui.RawEventResp(1)
}
func interactionHelpOpen() uintptr {
	lang := nox_strman_get_lang_code()
	if nox_xxx_guiFontHeightMB_43F320(nil) > 10 {
		lang = 2
	}
	name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 164512+4*uintptr(uint32(lang)))))
	w := Nox_new_window_from_file(name, interactionHelpEvent)
	interactionHelpRootWord = uint32(uintptr(w.C()))
	if w == nil {
		return 0
	}
	w.SetParent(nil)
	w.ShowModal()
	w.StackPush()
	GetClient().Cli().GUI.Focus(w)
	w.SetPos(image.Pt((int(nox_win_width)-w.SizeVal.X)/2, (int(nox_win_height)-w.SizeVal.Y)/2))
	label := w.ChildByID(4102)
	id := "Sanchlp.wnd:ClientHelp"
	if noxflags.HasGame(1) {
		id = "Sanchlp.wnd:Help"
	}
	first := fmt.Sprintf(interactionHelpText(id), GoWString(sub_42E8E0(45, 1))) + " "
	second := fmt.Sprintf(interactionHelpText("cdecode.c:KeyToChat"), GoWString(sub_42E8E0(8, 1)))
	interactionWriteText(1304400, second)
	interactionWriteText(1304656, first+second)
	// The static text control can retain this mapped pointer.
	label.Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(memmap.PtrUint16(0x5D4594, 1304656))), 0))
	if noxflags.HasGame(1) {
		if serverOptionsRoot == 0 {
			serverOptionsConstruct()
		}
		serverOptionsVisible(1)
	}
	if noxflags.HasGame(4096) || questRuntimeWord(1556160) != 0 {
		return uintptr(uint32(interactionHelpClose()))
	}
	return 0
}

func nox_xxx_cliShowHelpGui_49C560() *uint32 {
	return (*uint32)(unsafe.Pointer(interactionHelpOpen()))
}

func nox_xxx_wnd_49C760(w, code int32, a *int32, b int32) int32 {
	return int32(gui.EventRespInt(interactionHelpEvent((*gui.Window)(unsafe.Pointer(uintptr(uint32(w)))), &gui.RawEvent{Event: int(code), Arg1: uintptr(unsafe.Pointer(a)), Arg2: uintptr(uint32(b))})))
}

func sub_49C7A0() int32 { return int32(interactionHelpClose()) }

func sub_49C810() int32 { return int32(bool2int(interactionHelpRootWord != 0)) }

func sub_48D4B0(v int32) int32 {
	*memmap.PtrUint32(0x5D4594, 1197304) = uint32(v)
	hidden := 1
	if v == 1 {
		hidden = 0
	}
	return int32(interactionVoteHide(hidden))
}
