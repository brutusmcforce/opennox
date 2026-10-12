package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func interactionChatRoot() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(interactionChatRootWord))
}
func interactionChatEdit() *gui.Window {
	return (*gui.Window)(memmap.WordPtr(interactionChatEditWord))
}
func interactionChatData() unsafe.Pointer { return memmap.WordPtr(interactionChatDataWord) }
func interactionChatStart(team uint32) {
	if noxflags.HasGame(2048) || interactionChatActive != 0 {
		return
	}
	data := interactionChatData()
	*(*uint16)(data) = 0
	*(*uint16)(unsafe.Add(data, 1052)) = 0
	interactionChatRoot().ShowModal()
	interactionChatRoot().StackPush()
	GetClient().Cli().GUI.Focus(interactionChatEdit())
	interactionChatActive = 1
	*memmap.PtrUint32(0x5D4594, 1064872) = team
}
func interactionSay(text *uint16, team int) uint32 {
	code := ClientPlayerNetCode()
	dr := GetClient().Cli().Objs.ByNetCodeDynamic(code)
	units := gameplayTextUnits(text)
	original := len(units)
	for len(units) != 0 && units[0] == ' ' {
		units = units[1:]
	}
	if len(units) == 0 {
		return uint32(original)
	}
	x, y := uint16(0xffff), uint16(0xffff)
	if dr != nil {
		x, y = uint16(dr.PosVec.X), uint16(dr.PosVec.Y)
	}
	var flags byte
	if team != 0 {
		flags = 1
	}
	// The shared formatter already preserves raw UTF-16 and byte-count wrapping.
	msg := gameplayTextMessage(units, flags, uint16(code), x, y, 0)
	return uint32(bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg)))
}
func interactionChatDraw(w *gui.Window, d *gui.WindowData) int {
	root, edit := interactionChatRoot(), interactionChatEdit()
	root.ShowModal()
	GetClient().Cli().GUI.Focus(edit)
	r := GetClient().R2()
	face := r.GetFonts().AsFont(nil)
	data := interactionChatData()
	width := int32(r.GetStringSizeWrapped(face, alloc.GoString16((*uint16)(data)), 0).X) + int32(r.GetStringSizeWrapped(face, alloc.GoString16((*uint16)(unsafe.Add(data, 512))), 0).X) + 10
	if width < 100 {
		width = 100
	} else if width > 320 {
		width = 320
	}
	root.SetPos(image.Pt(int((int32(nox_win_width)-width)/2), root.Off.Y))
	uiWindowResize(w, int(width), 20)
	return uiEntryDraw(w, d, false)
}
func interactionChatClose() int {
	root := interactionChatRoot()
	if uiWindowHidden(root) != 0 {
		return 0
	}
	edit := interactionChatEdit()
	g := GetClient().Cli().GUI
	if g.Focused() == edit {
		g.Focus(nil)
	}
	root.StackPop()
	root.Hide()
	root.Flags &^= 8
	edit.Flags &^= 8
	g.ValYYY = 1
	interactionChatActive = 0
	return 1
}
func interactionChatOpen() *gui.Window {
	x, y := int32(nox_win_width)/2, 2*int32(nox_win_height)/3
	*memmap.PtrUint32(0x5D4594, 1064876) = uint32(x)
	*memmap.PtrUint32(0x5D4594, 1064880) = uint32(y)
	root := Nox_new_window_from_file("GuiChat.wnd", interactionChatEvent)
	interactionChatRootWord = memmap.PtrWord(root.C())
	if root == nil {
		return nil
	}
	root.SetPos(image.Pt(int(x), int(y)))
	edit := root.ChildByID(9201)
	interactionChatEditWord = memmap.PtrWord(edit.C())
	if edit == nil {
		return nil
	}
	edit.SetDraw(interactionChatDraw)
	edit.SetFunc93(interactionChatKey)
	interactionChatDataWord = memmap.PtrWord(edit.WidgetData)
	return root
}
func interactionChatKey(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	if ev.EventCode() != 21 || a != 1 {
		return uiEntryInput(w, ev)
	}
	if b == 2 {
		nox_xxx_consoleEsc_49B7A0()
	}
	return gui.RawEventResp(1)
}
func interactionChatEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16415 {
		data := interactionChatData()
		if *(*uint16)(unsafe.Add(data, 1052)) != 0 {
			interactionSay((*uint16)(data), int(memmap.Uint32(0x5D4594, 1064872)))
		}
		interactionChatClose()
	}
	return nil
}
func interactionChatDestroy() int {
	result := 0
	if root := interactionChatRoot(); root != nil {
		root.Destroy()
		interactionChatRootWord = 0
	}
	interactionChatEditWord = 0
	interactionChatDataWord = 0
	interactionChatActive = 0
	*memmap.PtrUint32(0x5D4594, 1064872) = 0
	return result
}

func sub_469FA0() int32 { return int32(memmap.Uint32(0x5D4594, 1064848)) }

func nox_client_chatStart_46A430(v int32) { interactionChatStart(uint32(v)) }

func sub_46A4A0() int32 { return int32(interactionChatActive) }

func nox_xxx_cmdSayDo_46A4B0(text *uint16, v int32) uint32 {
	return uint32(interactionSay((*uint16)(unsafe.Pointer(text)), int(v)))
}

func sub_46A6A0() int32 { return int32(interactionChatClose()) }

func sub_46A730() *uint32 { return (*uint32)(interactionChatOpen().C()) }

func sub_46A860() int32 { return int32(interactionChatDestroy()) }
