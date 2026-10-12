package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func interactionConversationWindow() *gui.Window { return (*gui.Window)(interactionConversationRoot) }
func interactionConversationChoice(choice byte) {
	*memmap.PtrUint8(0x5D4594, 1123516) = choice
	GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, []byte{0xd0, 2, choice})
}
func interactionConversationCancel() int {
	if uiWindowHidden(interactionConversationWindow()) == 1 {
		return 0
	}
	interactionConversationChoice(0)
	return 1
}
func interactionConversationOpen() int {
	*memmap.PtrUint32(0x5D4594, 1107052) = uint32(noxcolor.RGB5551Color(240, 128, 64))
	w := Nox_new_window_from_file("Dialog.wnd", interactionConversationEvent)
	interactionConversationRoot = w.C()
	if w == nil {
		return 0
	}
	w.SetFunc93(interactionConversationInput)
	w.SetDraw(interactionConversationDraw)
	w.SetTooltipFunc(finalTooltipKey(tooltipConversation))
	slider, up, down, list := w.ChildByID(3904), w.ChildByID(3903), w.ChildByID(3902), w.ChildByID(3901)
	hilite := Nox_xxx_gLoadImg("UISliderLit")
	selected := Nox_xxx_gLoadImg("UISliderLit")
	normal := Nox_xxx_gLoadImg("UISlider")
	gui.ButtonSetImage(slider, nil, nil, normal.C(), selected.C(), hilite.C())
	slider.DrawData().Window = list
	up.DrawData().Window = list
	down.DrawData().Window = list
	data := uiListData(list)
	data.Field_9 = slider.C()
	data.Field_7 = up.C()
	data.Field_8 = down.C()
	w.ChildByID(3906).SetDraw(interactionConversationBlink)
	uiWindowEnable(w, 0)
	w.Hide()
	interactionConversationActive = 0
	return 1
}
func interactionConversationEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() != 16391 {
		return nil
	}
	arg, _ := ev.EventArgsC()
	id := -2
	if child := (*gui.Window)(unsafe.Pointer(arg)); child != nil {
		id = int(child.ID())
	}
	if *bookWord(1047520) != 0 {
		return nil
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	switch id {
	case 3906:
		interactionConversationCancel()
	case 3907:
		Dialogs.PlayFile(alloc.GoString(DialogVoice), 100)
	case 3908:
		interactionConversationChoice(1)
	case 3909:
		interactionConversationChoice(2)
	}
	return nil
}
func interactionConversationInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	// The original mouse cases only computed and discarded a containment result.
	return gui.RawEventResp(1)
}
func interactionConversationBlink(w *gui.Window, d *gui.WindowData) int {
	seq := byte(GetClient().GetInputSeq())
	if !Dialogs.Sub_44D930() && seq&127 < 30 && seq&8 != 0 {
		p := uiWindowPosition(w)
		uiRenderBorder(p.X, p.Y, w.SizeVal.X, w.SizeVal.Y-2, int(memmap.Uint32(0x5D4594, 1107052)), 4)
	}
	return gui.ButtonDrawNoImg(w, d)
}
func interactionConversationDraw(w *gui.Window, d *gui.WindowData) int {
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(d.BgImageHnd), image.Pt(int(nox_win_width)-640, int(nox_win_height)-480))
	return 1
}
func interactionConversationDestroy() int {
	nox_xxx_windowDestroyMB_46C4E0((*nox_window)(interactionConversationWindow().C()))
	interactionConversationRoot = nil
	interactionConversationActive = 0
	return 0
}

func sub_479950() int32 { return int32(interactionConversationCancel()) }

func sub_4799A0() int32 { return int32(interactionConversationOpen()) }

func sub_479BE0(w *uint32, code int32, a uint32, b int32) int32 { return 1 }

func sub_479D10() int32 { return int32(interactionConversationDestroy()) }

func sub_47A260() int32 { return int32(interactionConversationActive) }
