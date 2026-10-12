package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func browserMouseDraw(_ *gui.Window, _ *gui.WindowData) int {
	p := GetClient().GetMousePos()
	detail := browserWindow(uint32(uintptr(browserUI.detailPanel)))
	if !browserHidden(detail) {
		enlarged := *detail
		enlarged.Off.X -= 32
		enlarged.Off.Y -= 32
		enlarged.SizeVal.X += 64
		enlarged.SizeVal.Y += 64
		if browserUI.transition == 0 && !serverOptionsPointIn(&enlarged, p) {
			detail.SetHidden(true)
			optionsSend(browserWindow(uint32(browserUI.gameList)), 16403, 0xffffffff, 0)
			browserUI.hasSelection = 0
			detail.StackPop()
			GetClient().Cli().GUI.Focus(browserUI.world)
		}
	}
	if browserUI.popup != 0 && !serverOptionsPointIn(browserWindow(memmap.Uint32(0x5D4594, 815036)), p) {
		browserPopupClose()
		GetClient().Cli().GUI.Focus(browserUI.world)
	}
	if browserUI.creating != 0 && sub_438DD0(uint32(p.X), uint32(p.Y)) != 0 {
		nox_client_setCursorType_477610(9)
	} else if Sub_44A4A0() == 0 {
		nox_client_setCursorType_477610(0)
	}
	return 1
}
func browserMapInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	switch ev.EventCode() {
	case 5:
		if w.ID() == 10020 && browserUI.creating == 1 {
			browserCreateAt(uint32(a)&0xffff, uint32(a)>>16)
			return gui.RawEventResp(1)
		}
	case 21:
		if a != 1 {
			if a == 57 {
				p := GetClient().GetMousePos()
				w.Func93(&gui.RawEvent{Event: 5, Arg1: uintptr(uint32(p.X) | uint32(p.Y)<<16)})
			}
			return nil
		}
		if b == 2 {
			Sub_4373A0()
		}
		return gui.RawEventResp(1)
	}
	return nil
}
func browserCreateAt(x, y uint32) int {
	hit := int(sub_438DD0(uint32(x), uint32(y)))
	if hit == 0 {
		return 0
	}
	off := uintptr(uint32(browserUI.region)) * 8
	*memmap.PtrUint16(0x5D4594, 814916) = uint16(x + uint32(memmap.Uint16(0x587000, 87528+off)) - 216)
	*memmap.PtrUint16(0x5D4594, 814918) = uint16(y + uint32(memmap.Uint16(0x587000, 87530+off)) - 27)
	browserChooseCharacter()
	if memmap.Uint32(0x5D4594, 815092)&2 != 0 {
		questRuntimeSetWord(1556160, 1)
	}
	count := 0
	if questRuntimeWord(1556160) != 0 {
		count = Nox_client_countPlayerFiles04_4DC7D0()
	} else {
		count = int(sessionCharacterCount())
	}
	if count != 0 {
		Sub_4A7A70(1)
		Nox_game_showSelChar_4A4DB0()
		return nox_client_setCursorType_477610(0)
	}
	Sub_4A7A70(0)
	characterShowClass()
	nox_client_setCursorType_477610(0)
	return hit
}
func browserMapDraw(w *gui.Window, _ *gui.WindowData) int {
	pos := uiWindowPosition(w)
	d := w.DrawData()
	img := d.BgImageHnd
	if d.Field0&6 != 0 {
		img = d.HlImageHnd
	}
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(img), pos.Add(d.ImgPtVal))
	for c := w.Field100Ptr; c != nil; c = c.Prev() {
		if !c.Flags.IsHidden() && c.DrawData().Style == 2048 {
			r.Data().SetTextColor(noxcolor.RGBA5551(Get_nox_color_white_2523948()))
			text := alloc.GoString16(*(**uint16)(c.WidgetData))
			r.DrawStringStyle(c.DrawData().Font(), text, pos.Add(c.Off))
		}
	}
	return 1
}

func sub_439D00(w *int32, code int32, a uint32, b int32) int32 {
	return int32(gui.EventRespInt(browserMapInput((*gui.Window)(unsafe.Pointer(w)), &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))})))
}
