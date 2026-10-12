package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func interactionGameOverRoot() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(interactionGameOverRootWord))
}
func interactionGameOverText(id string) string {
	return GetClient().Cli().Strings().GetStringInFile(strman.ID(id), "GUIGGOvr.c")
}
func interactionWriteText(off uintptr, text string) {
	ptr, free := alloc.CString16(text)
	defer free()
	interactionCopyText(memmap.PtrUint16(0x5D4594, off), ptr)
}
func interactionWindowText(w *gui.Window, text *uint16) int {
	if w != nil {
		w.Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(text)), 0))
	}
	return 0
}
func interactionGameOverOpen() int {
	w := Nox_new_window_from_file("GGOver.wnd", interactionGameOverEvent)
	interactionGameOverRootWord = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	w.Hide()
	uiWindowEnable(w, 0)
	return 1
}
func interactionGameOverHide() int {
	w := interactionGameOverRoot()
	nox_window_set_hidden((*nox_window)(w.C()), 1)
	uiWindowEnable(w, 0)
	GetClient().Cli().GUI.Focus(nil)
	return 0
}
func interactionGameOverEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16391 {
		arg, _ := ev.EventArgsC()
		id := -2
		if child := (*gui.Window)(unsafe.Pointer(arg)); child != nil {
			id = int(child.ID())
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		switch id {
		case 10701:
			Nox_client_quit_4460C0()
			interactionGameOverHide()
		case 10702:
			reliableClientSend(31, []byte{0xf0, 3}, nil, 1)
			interactionGameOverHide()
		}
	}
	return nil
}
func interactionGameOverDestroy() int {
	result := nox_xxx_windowDestroyMB_46C4E0((*nox_window)(interactionGameOverRoot().C()))
	interactionGameOverRootWord = 0
	return result
}
func interactionGameOverShow(data *uint16) uint32 {
	w := interactionGameOverRoot()
	w.Show()
	uiWindowEnable(w, 1)
	Nox_xxx_clientPlaySoundSpecial_452D80(1007, 100)
	w.SetPos(image.Pt(int(nox_win_width)/2-w.SizeVal.X/2, int(nox_win_height)/2-w.SizeVal.Y/2))
	values := unsafe.Slice(data, 4)
	for i, row := range []struct {
		off uintptr
		id  string
	}{{1302172, "GGOver.wnd:GeneratorsDestroyed"}, {1301916, "GGOver.wnd:NumSecretsFound"}, {1302428, "GGOver.wnd:Kills"}} {
		interactionWriteText(row.off, fmt.Sprintf(interactionGameOverText(row.id), int(values[i+1])))
	}
	interactionCopyText(memmap.PtrUint16(0x5D4594, 1303196), memmap.PtrUint16(0x5D4594, 1303460))
	for _, row := range []struct {
		id  uint
		off uintptr
	}{{10710, 1302172}, {10705, 1302940}, {10706, 1302684}, {10707, 1301916}, {10708, 1302428}, {10711, 1303196}} {
		interactionWindowText(w.ChildByID(row.id), memmap.PtrUint16(0x5D4594, row.off))
	}
	frame := gameFrame()
	*memmap.PtrUint32(0x5D4594, 1303456) = frame
	return frame
}
func interactionGameOverTick() int {
	w := interactionGameOverRoot()
	if w == nil {
		return 0
	}
	if hidden := uiWindowHidden(w); hidden != 0 {
		return hidden
	}
	remaining := int32(memmap.Uint32(0x5D4594, 1303456) + 30*gameFPS() - gameFrame())
	if remaining < 0 {
		remaining = 0
	}
	p := Get_dword_8531A0_2576()
	if p != nil && p.PlayerInd == 31 {
		interactionCopyText(memmap.PtrUint16(0x5D4594, 1301852), memmap.PtrUint16(0x5D4594, 1303464))
	} else {
		interactionWriteText(1301852, fmt.Sprintf("%s - %d", interactionGameOverText("Rules.c:Time"), uint32(remaining)/gameFPS()))
	}
	return interactionWindowText(w.ChildByID(10712), memmap.PtrUint16(0x5D4594, 1301852))
}

func sub_49B3E0() int32 { return int32(interactionGameOverOpen()) }

func sub_49B490() int32 { return int32(interactionGameOverDestroy()) }

func sub_49B6B0() int32 { return int32(interactionGameOverHide()) }

func sub_49B4B0(data *uint16) int32 {
	return int32(interactionGameOverShow((*uint16)(unsafe.Pointer(data))))
}

func sub_49B6E0() int32 { return int32(interactionGameOverTick()) }
