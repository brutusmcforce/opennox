package legacy

import (
	"image"
	"strconv"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func quickbarImage(name string) noxrender.ImageHandle {
	return noxrender.ImageHandle(unsafe.Pointer(uintptr(uiMeterLoadImage(name))))
}

func quickbarInitWindow(b *quickbarRecord, x, y, count, padding int, proc gui.WindowFunc, draw gui.WindowDrawFunc) int {
	g := GetClient().Cli().GUI
	b.Window = g.NewWindowRaw(nil, 1224, x-10, y-15, 199, 59, nil)
	b.Window.DrawData().BgImageHnd = quickbarImage("QuickBarBase")
	b.Window.DrawData().ImgPtVal = image.Pt(-61, -18)
	b.Window.SetFunc93(bookEvent(quickbarGenericEvent))
	xoff := padding + 10
	for i := 0; i < count; i++ {
		w := g.NewWindowRaw(b.Window, 1032, xoff, 15, 34, 34, nil)
		b.Slots[i] = w
		w.SetAllFuncs(proc, draw, nil)
		*quickbarUserData(w) = quickbarPointer(unsafe.Pointer(b))
		xoff += int(*memmap.PtrInt32(0x587000, 133488+uintptr(4*i)))
	}
	b.Current = quickbarRowPointer(b, b.Selected)
	return int(b.Selected)
}

func quickbarNewNamed(off uintptr, parent *gui.Window, flags gui.StatusFlags, x, y, width, height int) *gui.Window {
	w := GetClient().Cli().GUI.NewWindowRaw(parent, flags, x, y, width, height, nil)
	*quickbarWord(off) = quickbarPointer(w.C())
	return w
}

func quickbarSetTrapImages(class int) {
	name, tooltip := "QuickBarTrap", "ToolTipLayTrap"
	if class == 2 {
		name, tooltip = "QuickBarBomber", "ToolTipSummonBomber"
	}
	d := quickbarWindow(1049504).DrawData()
	d.BgImageHnd = quickbarImage(name)
	d.HlImageHnd = quickbarImage(name + "Hit")
	button, toggle := quickbarWindow(1049520), quickbarWindow(1049500)
	button.DrawData().SetTooltip(GetClient().Cli().Strings(), quickbarText(tooltip))
	button.Flags |= 8
	toggle.DrawData().SetTooltip(GetClient().Cli().Strings(), quickbarText("ToolTipTrapConstruct"))
	toggle.Flags |= 8
}

func quickbarRevealTrap() int {
	p := quickbarPlayer()
	if p == 0 {
		return 0
	}
	class := bookClass(p)
	if class != 1 && class != 2 {
		return class - 2
	}
	quickbarSetTrapImages(class)
	w := quickbarWindow(1049500)
	if w == nil {
		return -2
	}
	w.SetDraw(nil)
	return 0
}

func quickbarAddTrap(animate int) uintptr {
	if animate == 0 {
		return uintptr(quickbarRevealTrap())
	}
	w := quickbarWindow(1049500)
	ret := uintptr(w.C())
	if w.Flags&8 == 0 {
		*quickbarWord(1049536) = uint32(nox_win_height + 1)
		for i := 0; i < 50; i++ {
			rng := GetServer().S().Rand.Other
			timer := byte(rng.Int(4, 6))
			size := byte(rng.Int(3, 6))
			vy := rng.Int(-20, -5)
			vx := rng.Int(-5, 5)
			y := rng.Int(0, 20) + quickbarWindow(1049504).Off.Y + 10
			x := rng.Int(0, 20) + quickbarWindow(1049504).Off.X + 10
			ret = uintptr(unsafe.Pointer(screenParticleCreate(0, x, y, vx, vy, 1, size, timer, 2, 1)))
		}
	}
	return ret
}

func quickbarBookTooltip() int {
	name := "OpenSpellbookTT"
	if bookShown() != 0 {
		name = "CloseSpellbookTT"
	}
	Nox_xxx_cursorSetTooltip_4776B0(quickbarText(name))
	return 1
}

func quickbarDirectionTooltip(w *gui.Window) int {
	data := *quickbarUserData(w)
	name := "ToolTipCastAtOther"
	if quickbarDirectionLit(quickbarAt(1048196+uintptr(data>>16)*256), int(uint16(data))) != 0 {
		name = "ToolTipCastOnMe"
	}
	Nox_xxx_cursorSetTooltip_4776B0(quickbarText(name))
	return 1
}

func quickbarCreate() int {
	*quickbarWord(1047916) = 0
	*quickbarByte(1047920) = 0
	Sub_416170(5)
	*quickbarWord(1047924), *quickbarWord(1047928), *quickbarWord(1049480) = 0, 0, 0
	*quickbarByte(1049488) = 0
	r := GetClient().R2()
	height := r.FontHeight(r.GetFonts().AsFont(nil))
	x, y := (int(nox_win_width)-320)/2, int(nox_win_height)-74
	*quickbarWord(1047548), *quickbarWord(1047552) = uint32(x), uint32(y)
	*quickbarWord(1049684) = quickbarPointer(r.GetFonts().FontPtrByName("small"))
	for i, value := range []int{x, y - 17, x + 320, int(nox_win_height)} {
		*memmap.PtrUint32(0x587000, 133656+uintptr(i*4)) = uint32(value)
	}
	p := quickbarPlayer()
	if p == 0 {
		return 0
	}
	class := bookClass(p)
	draw := gui.WindowDrawFunc(quickbarDrawAbility)
	if class != 0 {
		draw = quickbarDrawSpell
	}
	quickbarInitWindow(quickbarMain(), x+69, y+32, 5, 0, bookEvent(quickbarSlotEvent), draw)
	for row, yy := 3, y+32-60; row >= 0; row, yy = row-1, yy-60 {
		b := quickbarAt(1048196 + uintptr(row)*256)
		quickbarInitWindow(b, x+69, yy, 5, 0, bookEvent(quickbarSlotEvent), quickbarDrawSpell)
		bookHideWindow(b.Window, true)
		b.Selected, b.Current = 0, &b.Rows[0]
	}
	flags := gui.StatusFlags(1672)
	if class != 0 {
		flags = 1160
	}
	right := quickbarNewNamed(1049504, nil, flags, x+260, y, 45, 66)
	right.DrawData().ImgPtVal = image.Pt(-263, 0)
	*quickbarWord(1049536) = uint32(nox_win_height - 74)
	right.SetFunc93(bookEvent(quickbarGenericEvent))
	if class != 0 {
		button := quickbarNewNamed(1049520, right, 1032, 9, 33, 32, 32)
		button.SetAllFuncs(bookEvent(quickbarTrapButtonEvent), quickbarDrawSlide, nil)
		toggle := quickbarNewNamed(1049500, right, 1160, 0, 19, 12, 12)
		toggle.SetFunc93(bookEvent(quickbarTrapEvent))
		toggle.DrawData().ImgPtVal = image.Pt(-265, -23)
		toggle.DrawData().BgImageHnd = quickbarImage("QuickBarTrapButton")
		toggle.DrawData().HlImageHnd = quickbarImage("QuickBarTrapButtonLit")
	}
	if class == 0 {
		right.DrawData().BgImageHnd = quickbarImage("QuickBarWarriorRight")
		right.DrawData().HlImageHnd = quickbarImage("QuickBarWarriorRight")
	} else if class == 1 || class == 2 {
		unavailable := *bookPlayerWord(p, 3832, 0) == 0 &&
			(!noxflags.HasGame(noxflags.GameFlag(0x2000)) || noxflags.HasGame(noxflags.GameFlag(4096)) || questRuntimeWord(1556160) != 0 || questRuntimeWord(1556164) != 0)
		if unavailable {
			right.DrawData().BgImageHnd = quickbarImage("QuickBarWarriorRight")
			right.DrawData().HlImageHnd = quickbarImage("QuickBarWarriorRight")
			quickbarWindow(1049520).Flags &^= 8
			quickbarWindow(1049500).Flags &^= 8
			quickbarWindow(1049500).SetDraw(quickbarDrawOne)
		} else {
			quickbarSetTrapImages(class)
		}
	}
	g := GetClient().Cli().GUI
	control := func(parent *gui.Window, x, y, width, height, dx, dy int, img, text string, action uint32) {
		w := g.NewWindowRaw(parent, 1032, x, y, width, height, nil)
		w.DrawData().BgImageHnd = nil
		w.DrawData().HlImageHnd = quickbarImage(img)
		w.SetAllFuncs(bookEvent(quickbarRowEvent), quickbarDrawControl, nil)
		w.DrawData().ImgPtVal = image.Pt(dx, dy)
		w.DrawData().SetTooltip(GetClient().Cli().Strings(), quickbarText(text))
		*quickbarUserData(w) = action
	}
	if class != 0 {
		trap := quickbarAt(1047940)
		quickbarInitWindow(trap, x+122, y-17, 3, 21, bookEvent(quickbarSlotEvent), quickbarDrawSpell)
		trap.Window.DrawData().BgImageHnd = quickbarImage("QuickBarTrapTray")
		trap.Window.DrawData().ImgPtVal = image.Pt(-40, -20)
		*quickbarWord(1048192) |= 1
		bookHideWindow(trap.Window, true)
		*quickbarWord(1049484) = 0
		label := g.NewWindowRaw(trap.Window, 1032, 20, -7, 110, height, nil)
		label.SetAllFuncs(nil, func(w *gui.Window, _ *gui.WindowData) int { return quickbarRowLabel(w, true) }, nil)
		control(trap.Window, 15, 12, 10, 14, -55, -32, "QuickBarTrapTrayUpLit", "ToolTipPrevTrap", 3)
		control(trap.Window, 15, 32, 10, 14, -55, -52, "QuickBarTrapTrayDownLit", "ToolTipNextTrap", 4)
	}
	left := quickbarNewNamed(1049508, nil, 1032, x-1, y+26, 61, 48)
	left.DrawData().ImgPtVal = image.Pt(1, -26)
	name := "QuickBarWarriorLeft"
	if class != 0 {
		name = "QuickBarSpellSetBase"
	}
	left.DrawData().BgImageHnd = quickbarImage(name)
	left.DrawData().HlImageHnd = quickbarImage(name)
	if class != 0 {
		*quickbarUserData(left) = 5
	}
	left.SetAllFuncs(bookEvent(quickbarGenericEvent), quickbarDrawControl, nil)
	book := quickbarNewNamed(1049524, left, 1160, 0, 9, 29, 30)
	book.DrawData().BgImageHnd = quickbarImage("SpellbookButton")
	book.DrawData().HlImageHnd = quickbarImage("SpellbookButtonLit")
	button := quickbarNewNamed(1049528, book, 1064, 1, 2, 28, 28)
	button.SetAllFuncs(bookEvent(quickbarBookEvent), quickbarDrawOne, bookTooltipKey(0))
	button.DrawData().SetTooltip(GetClient().Cli().Strings(), quickbarText("OpenSpellBookTT"))
	if class != 0 {
		control(left, 30, 0, 15, 19, -29, -26, "QuickBarSpellSetUpLit", "ToolTipPrevSpellSet", 0)
		control(left, 30, 29, 15, 19, -29, -55, "QuickBarSpellSetDownLit", "ToolTipNextSpellSet", 1)
		control(left, 48, 16, 13, 17, -47, -42, "QuickBarSpellSetMaxLit", "ToolTipAllSpellSets", 2)
		modifier := quickbarNewNamed(1049516, nil, 1032, 0, 0, 1, 1)
		modifier.SetAllFuncs(bookEvent(func(*gui.Window, uint32, uint32) int { return 0 }), func(*gui.Window, *gui.WindowData) int { quickbarModifier(); return 1 }, nil)
		title := quickbarNewNamed(1049512, nil, 1152, x, y, 2, 2)
		title.DrawData().BgImageHnd = quickbarImage("QuickBarTitle")
		label := g.NewWindowRaw(title, 8, 115, 6, 101, 14, nil)
		label.SetAllFuncs(nil, func(w *gui.Window, _ *gui.WindowData) int { return quickbarRowLabel(w, false) }, nil)
	}
	for row := 0; row < 5; row++ {
		b := quickbarAt(1048196 + uintptr(row)*256)
		pos := b.Window.Off.Add(image.Pt(10, 5))
		advance := 0
		for i := 0; i < 5; i++ {
			// Only the main row increments the legacy nugget counter.
			nugget := 0
			if row == 4 {
				nugget = i
			}
			w := g.NewWindowRaw(nil, 1160, pos.X, pos.Y, 30, 10, nil)
			name := "QuickBarNugget" + strconv.Itoa(nugget+1)
			if class == 0 {
				name = "QuickBarWarriorNugget" + strconv.Itoa(nugget+1)
			}
			w.DrawData().BgImageHnd = quickbarImage(name)
			if class != 0 {
				name += alloc.GoString((*byte)(memmap.PtrOff(0x587000, 134996)))
			}
			w.DrawData().HlImageHnd = quickbarImage(name)
			w.DrawData().ImgPtVal = image.Pt(-70-advance, -23)
			w.SetAllFuncs(bookEvent(quickbarGenericEvent), nil, nil)
			b.Directions[i] = w
			quickbarDirection(b, nugget)
			if class != 0 {
				direction := g.NewWindowRaw(w, 1032, 12, 0, 10, 10, nil)
				direction.SetAllFuncs(bookEvent(quickbarDirectionEvent), quickbarDrawOne, bookTooltipKey(1))
				*quickbarUserData(direction) = uint32(nugget | row<<16)
			}
			if row == 4 {
				label := g.NewWindowRaw(w, 1032, 13, 39, 10, 10, nil)
				label.SetAllFuncs(nil, quickbarKeyLabel, nil)
				*quickbarUserData(label) = uint32(nugget)
			} else {
				bookHideWindow(w, true)
			}
			delta := int(*memmap.PtrInt32(0x587000, 133488+uintptr(4*i)))
			pos.X += delta
			advance += delta
		}
	}
	if class == 0 {
		for i := 0; i < 5; i++ {
			for field := 0; field < 5; field++ {
				*quickbarWord(1047788 + uintptr(i*24+field*4)) = *memmap.PtrUint32(0x587000, 133536+uintptr(i*24+field*4))
			}
			*quickbarWord(1047788 + uintptr(i*24+20)) = 0
		}
	}
	quickbarSelectRow(0)
	return 1
}
