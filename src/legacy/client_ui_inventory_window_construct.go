package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

func uiInventoryImage(v uint32) noxrender.ImageHandle {
	return noxrender.ImageHandle(memmap.WordPtr(v))
}
func uiInventoryNewScrollControls(parent *gui.Window) int {
	for i, cfg := range [][4]int{{522, 2, 20, 25}, {522, 148, 20, 25}} {
		off := uintptr(1049940 + i*8)
		en, lit := uiInventoryImage(memmap.Uint32(0x5D4594, off)), uiInventoryImage(memmap.Uint32(0x5D4594, off+4))
		d := gui.WindowData{Style: 1, Window: parent, EnImageHnd: en, HlImageHnd: lit, SelImageHnd: lit, TextColorVal: 0x80000000}
		w := NewButtonOrCheckbox(parent, 1161, cfg[0], cfg[1], cfg[2], cfg[3], &d)
		*memmap.PtrUint32(0x5D4594, 1062500+uintptr(i*4)) = uiInventoryPointer(w.C())
		if w == nil {
			return 0
		}
		w.SetID(uint(9102 + i))
	}
	d := gui.WindowData{Style: 8, Window: parent, BgColorVal: 0x80000000, EnColorVal: 0x80000000, HlColorVal: 0x80000000, DisColorVal: 0x80000000, SelColorVal: 0x80000000}
	data := gui.SliderData{Max: 850}
	w := uiSliderNew(parent, 1033, 524, 42, 16, 91, &d, &data)
	dword_5d4594_1062508 = uint32(uiInventoryPointer(w.C()))
	if w == nil {
		return 0
	}
	w.SetFunc93(uiInventoryWindowEvent(uiInventoryTrackEvents))
	thumb := w.Field100()
	thumb.SetFunc93(uiInventoryWindowEvent(uiInventoryThumbEvents))
	// The C constructor sets dimensions without recomputing EndPos.
	thumb.SizeVal = image.Pt(16, 16)
	thumb.DrawData().ImgPtVal = image.Pt(0, -15)
	en := uiInventoryImage(memmap.Uint32(0x5D4594, 1049956))
	lit := uiInventoryImage(memmap.Uint32(0x5D4594, 1049960))
	gui.ButtonSetImage(w, nil, nil, en, lit, lit)
	return 1
}
func uiInventoryNewModeControls(parent *gui.Window) int {
	cfg := []struct {
		x, y, w, h int
		id         uint
		bg, sel    uint32
	}{{243, 170, 34, 34, 9105, 0, uint32(dword_5d4594_1049976)}, {5, 186, 34, 34, 9107, uint32(dword_5d4594_1049992), uint32(dword_5d4594_1049996)}, {547, 2, 16, 16, 9111, 0, memmap.Uint32(0x5D4594, 1049972)}}
	for i, c := range cfg {
		d := gui.WindowData{Style: 1, Window: parent, BgImageHnd: uiInventoryImage(c.bg), SelImageHnd: uiInventoryImage(c.sel), ImgPtVal: image.Pt(-c.x, -c.y)}
		w := NewButtonOrCheckbox(parent, 1161, c.x, c.y, c.w, c.h, &d)
		switch i {
		case 0:
			dword_5d4594_1062528 = uint32(uiInventoryPointer(w.C()))
		case 1:
			dword_5d4594_1062524 = uint32(uiInventoryPointer(w.C()))
		case 2:
			*memmap.PtrUint32(0x5D4594, 1062532) = uiInventoryPointer(w.C())
		}
		if w == nil {
			return 0
		}
		w.SetID(c.id)
		w.SetTooltipFunc(inventoryCallbackKey(inventoryCallbackMode))
	}
	return 1
}
func uiInventoryNewIdentifyWindow(parent *gui.Window) int {
	w := Nox_new_window_from_file("identify.wnd", nil)
	dword_5d4594_1062476 = uint32(uiInventoryPointer(w.C()))
	if w == nil {
		return 0
	}
	uiWindowResize(w, 200, 200)
	w.SetParent(parent)
	w.SetPos(image.Pt(51, 15))
	w.ChildByID(9155).SetDraw(func(w *gui.Window, d *gui.WindowData) int {
		return sub_466F50(w, d)
	})
	return 1
}
func uiInventoryCreateWindow() int {
	uiInventoryLoadImages()
	nox_xxx_inventoryNameSignInit_4671E0()
	dword_5d4594_1063636 = memmap.PtrWord(GetClient().R2().GetFonts().FontPtrByName("small"))
	vp := &uiInventoryViewport
	vp.Screen = image.Rect(0, 0, int(nox_win_width), int(nox_win_height))
	vp.Size = image.Pt(int(nox_win_width), int(nox_win_height))
	vp.World.Min = image.Point{}
	g := GetClient().Cli().GUI
	drawOne := func(*gui.Window, *gui.WindowData) int { return 1 }
	root := g.NewWindowRaw(nil, 552, 0, 0, 563, 264, nil)
	legacyGlobals.dword_5d4594_1062452 = root
	root.SetAllFuncs(nil, drawOne, nil)
	status := g.NewWindowRaw(root, 8, 0, 224, int(nox_win_width), 40, nil)
	status.SetAllFuncs(func(*gui.Window, gui.WindowEvent) gui.WindowEventResp { return gui.RawEventResp(0) }, drawOne, inventoryCallbackKey(inventoryCallbackStatus))
	main := g.NewWindowRaw(root, 40, 0, 0, 563, 224, uiInventoryWindowEvent(uiInventoryPanelEvents))
	dword_5d4594_1062456 = uint32(uiInventoryPointer(main.C()))
	main.SetAllFuncs(uiInventoryWindowEvent(uiInventoryMainEvents), func(w *gui.Window, _ *gui.WindowData) int { return uiInventoryDrawWindow(w) }, inventoryCallbackKey(inventoryCallbackHover))
	main.DrawData().Style |= 0x100
	catcher := g.NewWindowRaw(root, 40, 0, 0, 1, 1, nil)
	*memmap.PtrUint32(0x5D4594, 1062472) = uiInventoryPointer(catcher.C())
	catcher.SetAllFuncs(uiInventoryWindowEvent(uiInventoryMainEvents), drawOne, nil)
	alt := g.NewWindowRaw(main, 40, 173, 174, 50, 50, nil)
	dword_5d4594_1062468 = uint32(uiInventoryPointer(alt.C()))
	alt.SetAllFuncs(uiInventoryWindowEvent(uiInventoryAlternateEvents), func(w *gui.Window, _ *gui.WindowData) int { return sub_4625D0(w) }, inventoryCallbackKey(inventoryCallbackAlt))
	alt.DrawData().Style |= 0x100
	if uiInventoryNewScrollControls(main) == 0 || uiInventoryNewModeControls(main) == 0 || uiInventoryNewIdentifyWindow(main) == 0 {
		return 0
	}
	current := g.NewWindowRaw(nil, gui.StatusFlags(0x408|gui.StatusBelow), -1, int(nox_win_height)-127, 111, 127, nil)
	legacyGlobals.nox_win_unk5 = current
	if current == nil {
		return 0
	}
	current.SetAllFuncs(uiInventoryWindowEvent(uiInventoryWindowAdmission), func(w *gui.Window, _ *gui.WindowData) int {
		return int(nox_xxx_inventoryDrawProc_466580(w))
	}, nil)
	current.DrawData().BgImageHnd = uiInventoryImage(uiMeterLoadImage("CurrentWeapon"))
	current.DrawData().HlImageHnd = uiInventoryImage(uiMeterLoadImage("CurrentWeaponLit"))
	current.DrawData().ImgPtVal = image.Pt(-1, 0)
	nox_win_init_cur_weapon(current, 24, 51, 53, 53)
	sub_471160(current, 79, 40, 20, 127)
	sub_470D70()
	button := g.NewWindowRaw(current, 8, 5, 11, 28, 29, nil)
	button.SetAllFuncs(uiInventoryWindowEvent(uiInventoryToggleButton), drawOne, inventoryCallbackKey(inventoryCallbackButton))
	grid := uiInventoryGrid()
	clear(grid[:])
	if dword_5d4594_1062560 == 0 {
		dword_5d4594_1062560 = uint32(GetClient().Cli().Things.IndByID("Gold"))
		*memmap.PtrUint32(0x5D4594, 1049728) = uint32(GetClient().Cli().Things.IndByID("QuestGoldPile"))
		*memmap.PtrUint32(0x5D4594, 1049724) = uint32(GetClient().Cli().Things.IndByID("QuestGoldChest"))
	}
	create := func(index int, typ uint32) {
		cell := &grid[index]
		cell.Drawable = GetClient().Nox_new_drawable_for_thing(int(typ))
		if cell.Drawable != nil {
			cell.Count = 1
		}
	}
	create(20, uint32(dword_5d4594_1062560))
	if dword_5d4594_1062564 == 0 {
		dword_5d4594_1062564 = uint32(GetClient().Cli().Things.IndByID("Identify"))
	}
	create(41, uint32(dword_5d4594_1062564))
	if dword_5d4594_1062556 == 0 {
		dword_5d4594_1062556 = uint32(GetClient().Cli().Things.IndByID("AutoMap"))
	}
	create(62, uint32(dword_5d4594_1062556))
	return int(uiInventoryPointer(main.C()))
}

func uiInventoryLoadImages() uintptr {
	*memmap.PtrUint32(0x5D4594, 1049908) = uiMeterLoadImage("InventoryBase")
	*memmap.PtrUint32(0x5D4594, 1049912) = uiMeterLoadImage("InventoryIdentifyBase")
	*memmap.PtrUint32(0x5D4594, 1049916) = uiMeterLoadImage("InventoryTray1")
	*memmap.PtrUint32(0x5D4594, 1049920) = uiMeterLoadImage("InventoryTray2")
	*memmap.PtrUint32(0x5D4594, 1049924) = uiMeterLoadImage("InventoryTray3")
	*memmap.PtrUint32(0x5D4594, 1049928) = uiMeterLoadImage("InventoryTraySpecial")
	*memmap.PtrUint32(0x5D4594, 1049932) = uiMeterLoadImage("InventoryTrayIdentifyLit")
	*memmap.PtrUint32(0x5D4594, 1049936) = uiMeterLoadImage("InventoryTrayMapLit")
	*memmap.PtrUint32(0x5D4594, 1049940) = uiMeterLoadImage("InventoryUpButton")
	*memmap.PtrUint32(0x5D4594, 1049944) = uiMeterLoadImage("InventoryUpButtonLit")
	*memmap.PtrUint32(0x5D4594, 1049948) = uiMeterLoadImage("InventoryDownButton")
	*memmap.PtrUint32(0x5D4594, 1049952) = uiMeterLoadImage("InventoryDownButtonLit")
	*memmap.PtrUint32(0x5D4594, 1049956) = uiMeterLoadImage("InventorySliderButton")
	*memmap.PtrUint32(0x5D4594, 1049960) = uiMeterLoadImage("InventorySliderButtonLit")
	*memmap.PtrUint32(0x5D4594, 1049964) = uiMeterLoadImage("InventoryEquipRing")
	*memmap.PtrUint32(0x5D4594, 1049968) = uiMeterLoadImage("InventoryQuickItemRing")
	*memmap.PtrUint32(0x5D4594, 1049972) = uiMeterLoadImage("InventoryCloseButtonLit")
	dword_5d4594_1049976 = uint32(uiMeterLoadImage("InventoryJournalButtonLit"))
	*memmap.PtrUint32(0x5D4594, 1049980) = uiMeterLoadImage("InventoryInventoryButton")
	*memmap.PtrUint32(0x5D4594, 1049984) = uiMeterLoadImage("InventoryInventoryButtonLit")
	*memmap.PtrUint32(0x5D4594, 1049988) = uiMeterLoadImage("InventoryDollButtonLit")
	dword_5d4594_1049992 = uint32(uiMeterLoadImage("InventoryStatsButton"))
	dword_5d4594_1049996 = uint32(uiMeterLoadImage("InventoryStatsButtonLit"))
	*memmap.PtrUint32(0x5D4594, 1050000) = uiMeterLoadImage("GUIFist")
	*memmap.PtrUint32(0x5D4594, 1050004) = uiMeterLoadImage("SharedKeyMode")
	ref := Nox_xxx_gLoadAnim("ExtraLives")
	dword_5d4594_1050008 = uint32(uiInventoryPointer(ref.C()))
	return uintptr(ref.C())
}
