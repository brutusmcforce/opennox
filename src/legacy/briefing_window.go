package legacy

import (
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func briefingParent() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(legacyGlobals.dword_5d4594_831236)))
}
func briefingWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(nox_wnd_briefing_831232)))
}
func briefingPlayVoice() int {
	result := int(dword_5d4594_831224)
	*memmap.PtrUint32(0x5D4594, 831248) = 1
	if result != 0 {
		result = Dialogs.PlayFile(alloc.GoString((*byte)(unsafe.Pointer(uintptr(dword_5d4594_831240)))), 100)
		dword_5d4594_831244 = 1
	}
	return result
}

// Preserve the old C callback adapter's nil response for a zero result.
func briefingEvent(fn func(*gui.Window, int, uintptr, uintptr) int) gui.WindowFunc {
	return func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		if v := fn(w, ev.EventCode(), a, b); v != 0 {
			return gui.RawEventResp(uintptr(v))
		}
		return nil
	}
}
func briefingBackgroundEvent(_ *gui.Window, event int, _, _ uintptr) int {
	return bool2int(event == 23)
}
func briefingCreateWindow() *gui.Window {
	ox, oy := (int(nox_win_width)-640)/2, (int(nox_win_height)-480)/2
	*memmap.PtrUint32(0x5D4594, 831284) = uint32(ox)
	*memmap.PtrUint32(0x5D4594, 831288) = uint32(oy)
	parent := GetClient().Cli().GUI.NewWindowRaw(nil, 56, 0, 0, int(nox_win_width), int(nox_win_height), briefingEvent(briefingBackgroundEvent))
	legacyGlobals.dword_5d4594_831236 = uint32(uintptr(parent.C()))
	if parent == nil {
		return nil
	}
	parent.SetFunc93(briefingEvent(briefingInput))
	parent.DrawData().BgColorVal = uint32(nox_color_black_2650656)
	w := Nox_new_window_from_file("Briefing.wnd", nil)
	nox_wnd_briefing_831232 = uint32(uintptr(w.C()))
	if w == nil {
		return nil
	}
	w.SetParent(parent)
	w.SetPos(image.Pt(ox, oy))
	w.ChildByID(1010).SetDraw(briefingDrawWindow)
	briefingLoadChapters()
	return parent
}
func briefingFadeBack() {
	nox_gameDisableMapDraw_5d4594_2650672 = 0
	GetClient().R2().FadeInScreen(Nox_client_getIntroScreenDuration_44E3B0(), true, Sub_44E320)
	nox_gameDisableMapDraw_5d4594_2650672 = 1
}
func briefingInput(_ *gui.Window, event int, _, _ uintptr) int {
	if memmap.Uint32(0x5D4594, 831248) != 0 && event != 18 && event != 17 && memmap.Uint32(0x5D4594, 527720) != 1 {
		Sub_450580()
		if dword_5d4594_831220 != 0 {
			if dword_5d4594_831220 == 255 {
				briefingFadeBack()
			}
		} else {
			briefingParent().Capture(false)
			GetClient().Cli().GUI.Focus(nil)
			dword_5d4594_831256 = 1
			Nox_savegame_sub_46D580()
		}
		briefingParent().SetFunc93(nil)
	}
	return 1
}
func briefingScrollSpeed() float64 {
	if dword_5d4594_831220 == 255 {
		return 1
	}
	return 0
}
func briefingCompletedDraw(_ *gui.Window, _ *gui.WindowData) int { return 1 }
func briefingDrawWindow(w *gui.Window, data *gui.WindowData) int {
	scroll := float32(float64(math.Float32frombits(uint32(dword_5d4594_831276))) - briefingScrollSpeed())
	dword_5d4594_831276 = uint32(math.Float32bits(scroll))
	w.SetPos(image.Pt(0, int(int32(scroll))))
	uiRenderCopyRect(int(int32(memmap.Uint32(0x5D4594, 831284))), int(int32(memmap.Uint32(0x5D4594, 831288))), 640, 480)
	mode := memmap.Uint8(0x5D4594, 832472)
	if mode&1 != 0 {
		briefingDrawStats(data)
	} else if mode&2 != 0 {
		briefingDrawTitle(data)
	} else if mode&4 != 0 {
		briefingDrawInstructions(data)
	} else {
		gui.StaticTextDrawNoImage(w, data)
	}
	rect := GetClient().Cli().Render().PixBufferRect()
	uiRenderCopyRect(0, 0, rect.Dx(), rect.Dy())
	// Read transition state after invoking the selected drawing owner.
	scroll = math.Float32frombits(uint32(dword_5d4594_831276))
	mode = memmap.Uint8(0x5D4594, 832472)
	if int64(scroll) <= int64(int32(memmap.Uint32(0x5D4594, 831280))) && dword_5d4594_831244 == 1 && !Dialogs.Sub_44D930() && PlatformTicks()-*memmap.PtrUint64(0x5D4594, 831292) > uint64(Nox_client_getBriefDuration()) && (memmap.Uint32(0x5D4594, 832488) == 1 || mode&5 == 0) {
		if mode&2 != 0 && dword_5d4594_832480 != 0 {
			Nox_xxx_clientPlaySoundSpecial_452D80(582, 100)
		}
		Sub_450580()
		if dword_5d4594_831220 != 0 {
			briefingFadeBack()
		} else if dword_5d4594_831256 == 0 {
			briefingParent().Capture(false)
			dword_5d4594_831256 = 1
			Nox_savegame_sub_46D580()
			briefingParent().SetFunc93(nil)
		}
		w.SetDraw(briefingCompletedDraw)
	}
	return 1
}
func briefingShow(chapter, begin int, mode byte) int {
	dword_5d4594_831260 = 1
	dword_5d4594_831244 = 1
	p := (*server.Player)(unsafe.Pointer(uintptr(dword_8531A0_2576)))
	*memmap.PtrUint32(0x5D4594, 832488) = 0
	*memmap.PtrUint32(0x5D4594, 832472) = 0
	dword_5d4594_831256 = 0
	*memmap.PtrUint8(0x5D4594, 831252) = byte(chapter + 1)
	class := begin
	if p != nil {
		class = int(p.PlayerClass())
	}
	parent := briefingParent()
	parent.ShowModal()
	parent.Capture(true)
	GetClient().Cli().GUI.Focus(parent)
	parent.SetFunc93(briefingEvent(briefingInput))
	w := briefingWindow()
	child := w.ChildByID(1010)
	child.SetDraw(briefingDrawWindow)
	dword_5d4594_831220 = uint32(begin)
	textData := (*gui.StaticTextData)(child.WidgetData)
	setImage := func(v uint32) { w.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(v))) }
	setText := func(p unsafe.Pointer) { child.Func94(gui.AsWindowEvent(16385, uintptr(p), 0)) }
	music := uint32(0)
	if chapter == 255 {
		setImage(0)
		setText(*memmap.PtrPtr(0x5D4594, 831268))
		child.Flags &^= 0x2000
		var voice *byte
		if s := Nox_xxx_GetEndgameDialog(); s != "" {
			voice = alloc.InternCString(s)
		}
		dword_5d4594_831240 = uint32(uintptr(unsafe.Pointer(voice)))
		music = 24
		dword_5d4594_831220 = 255
	} else if chapter == 254 {
		// The old unused map-filename getter has no side effects.
		setImage(briefingLoadImage("GauntletStartMines"))
		setText(memmap.PtrOff(0x5D4594, 832540))
		if mode&1 != 0 {
			*memmap.PtrUint32(0x5D4594, 832472) = 1
			setImage(briefingLoadImage("MenuSystemBG"))
		} else if mode&4 != 0 {
			*memmap.PtrUint32(0x5D4594, 832472) = 4
			setImage(briefingLoadImage("GauntletInstructionBackground"))
		} else {
			*memmap.PtrUint32(0x5D4594, 832472) = 2
			setImage(memmap.Uint32(0x5D4594, 832460))
			if caption := *memmap.PtrPtr(0x5D4594, 832464); caption != nil {
				setText(caption)
			}
		}
		Nox_client_resetScreenParticles_431510()
		Nox_xxx_bookHideMB_45ACA0(1)
		Sub_446780()
		dword_5d4594_831240 = 0
		dword_5d4594_831244 = 1
		dword_5d4594_831224 = 0
	} else {
		row := &(*[33]briefingChapter)(memmap.PtrOff(0x5D4594, 831300))[class*11+chapter]
		slide := &row.Loss
		if begin != 0 {
			slide = &row.Begin
		}
		setImage(slide.Image)
		setText(unsafe.Pointer(slide.Text))
		music = slide.Duration
		dword_5d4594_831240 = uint32(uintptr(unsafe.Pointer(slide.Voice)))
	}
	r := GetClient().R2()
	font := r.GetFonts().AsFont(child.DrawData().FontPtr)
	height := r.GetStringSizeWrapped(font, alloc.GoString16(textData.Text), 480).Y
	*memmap.PtrUint32(0x5D4594, 831280) = uint32(height)
	uiWindowResize(child, 640, height)
	if chapter == 255 {
		dword_5d4594_831276 = 1140457472
		*memmap.PtrUint32(0x5D4594, 831280) = uint32(-20 - height)
	} else {
		height = (480 - r.FontHeight(font) - height) / 2
		*memmap.PtrUint32(0x5D4594, 831280) = uint32(height)
		dword_5d4594_831276 = uint32(math.Float32bits(float32(height)))
	}
	Sub_431290()
	dword_5d4594_831224 = 0
	if begin != 0 || chapter == 255 {
		if chapter != 254 {
			dword_5d4594_831224 = 1
		}
	} else if chapter != 254 {
		seen := (*uint32)(unsafe.Add(p.C(), 4408+4*chapter))
		if *seen != 0 {
			parent.Capture(false)
			dword_5d4594_831256 = 1
			Nox_savegame_sub_46D580()
			parent.SetFunc93(nil)
		} else {
			dword_5d4594_831224 = 1
			*seen = 1
		}
	}
	MusicModule.Sub_43DD70(music, 50)
	*memmap.PtrUint64(0x5D4594, 831292) = PlatformTicks()
	*memmap.PtrUint32(0x5D4594, 831248) = 0
	r.FadeOutScreen(Nox_client_getIntroScreenDuration_44E3B0(), true, func() { briefingPlayVoice() })
	result := int(dword_5d4594_831224)
	if result != 0 {
		dword_5d4594_831244 = 0
	}
	return result
}
func briefingDestroy() int {
	if w := briefingParent(); w != nil {
		w.Destroy()
		legacyGlobals.dword_5d4594_831236 = 0
		nox_wnd_briefing_831232 = 0
	}
	dword_5d4594_831260 = 0
	dword_5d4594_832484 = 0
	sprites := briefingSprites()
	result := 0
	for _, i := range []int{3, 1, 2, 0, 4, 5, 6, 7, 8, 9, 10, 11} {
		p := sprites[i]
		if *p != 0 {
			v := GetClient().Nox_xxx_spriteDelete_45A4B0((*client.Drawable)(unsafe.Pointer(uintptr(*p))))
			if i == 11 {
				result = v
			}
		}
		*p = 0
	}
	return result
}
