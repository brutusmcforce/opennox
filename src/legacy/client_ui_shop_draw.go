package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"fmt"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func uiShopInit() int {
	uiTradeViewportInit(&uiShopViewport)
	w := Nox_new_window_from_file("Shop.wnd", uiInventoryWindowEvent(uiShopPanel))
	*uiShopWord(1098576) = uiInventoryPointer(w.C())
	if w == nil {
		return 0
	}
	w.SetAllFuncs(uiInventoryWindowEvent(uiShopMouse), func(_ *gui.Window, _ *gui.WindowData) int { return uiShopDraw() }, nil)
	w.ChildByID(3806).SetTooltipFunc(uiAmountNativeKey(uiAmountShopTooltip))
	for i := range uiShopGrid() {
		uiShopGrid()[i].Drawable = nil
	}
	uiShopClear()
	w.SetPos(image.Pt(int(nox_win_width)-w.SizeVal.X, int(nox_win_height)-w.SizeVal.Y))
	w.Hide()
	uiWindowEnable(w, 0)
	for _, v := range []struct {
		off  uintptr
		name string
	}{
		{1098400, "ShopBase"},
		{1098404, "ShopTradeMode"},
		{1098408, "ShopIdentifyMode"},
		{1098412, "ShopRepairMode"},
		{1098416, "ShopRepairMode"},
		{1098420, "ShopExitMode"},
		{1098424, "ShopInventoryBar1"},
		{1098428, "ShopInventoryBar2"},
		{1098432, "ShopInventorySlider"},
		{1098436, "ShopInventorySliderSelected"},
		{1098448, "ShopInventoryUp"},
		{1098452, "ShopInventoryUpSelected"},
		{1098440, "ShopInventorydown"},
		{1098444, "ShopInventorydownSelected"},
		{1098456, "ShopTextBorder"},
		{1098460, "ShopkeeperPic"},
		{1098464, "ShopkeeperWarriorPic"},
		{1098468, "ShopkeeperConjurerPic"},
		{1098472, "ShopkeeperWizardPic"},
		{1098476, "ShopkeeperLandOfTheDeadPic"},
		{1098480, "ShopkeeperMagicShopPic"},
	} {
		*uiShopWord(v.off) = uiMeterLoadImage(v.name)
	}
	slider := w.ChildByID(3807)
	*uiShopWord(1098580) = uiInventoryPointer(slider.C())
	*uiShopWord(1098584) = uiInventoryPointer(w.ChildByID(3808).C())
	*uiShopWord(1098588) = uiInventoryPointer(w.ChildByID(3809).C())
	grid := w.ChildByID(3806)
	p := uiWindowPosition(grid)
	*uiShopWord(1098380) = uint32(p.X)
	*uiShopWord(1098384) = uint32(p.Y)
	*uiShopWord(1098388) = uint32(p.X + grid.SizeVal.X)
	*uiShopWord(1098392) = uint32(p.Y + grid.SizeVal.Y)
	thumb := slider.Field100()
	thumb.SizeVal = image.Pt(16, 12)
	thumb.DrawData().ImgPtVal = image.Pt(0, -15)
	handle := func(off uintptr) noxrender.ImageHandle {
		return noxrender.ImageHandle(memmap.WordPtr(*uiShopWord(off)))
	}
	gui.ButtonSetImage(slider, nil, nil, handle(1098432), handle(1098436), handle(1098436))
	*uiShopWord(1098592) = uint32((*gui.SliderData)(slider.WidgetData).Max)
	return 1
}
func uiShopPicture(typ uint32) uint32 {
	entries := []struct {
		off      uintptr
		typ, img string
	}{
		{1098396, "Shopkeeper", "ShopKeeperPic"},
		{1098560, "ShopkeeperWarriorsRealm", "ShopKeeperWarriorPic"},
		{1098556, "ShopkeeperConjurerRealm", "ShopKeeperConjurerPic"},
		{1098564, "ShopkeeperWizardRealm", "ShopKeeperWizardPic"},
		{1098572, "ShopkeeperLandOfTheDead", "ShopKeeperLandOfTheDeadPic"},
		{1098568, "ShopkeeperMagicShop", "ShopKeeperMagicShopPic"},
		{1098484, "ShopkeeperPurple", "ShopKeeperPurplePic"},
		{1097292, "ShopkeeperYellow", "ShopKeeperBrownPic"},
	}
	if *uiShopWord(1107040) == 0 {
		for _, e := range entries {
			*uiShopWord(e.off) = uint32(GetClient().Cli().Things.IndByID(e.typ))
		}
		*uiShopWord(1107040) = 1
	}
	name := "ShopKeeperPic"
	for _, e := range entries {
		if typ == *uiShopWord(e.off) {
			name = e.img
			break
		}
	}
	h := uiMeterLoadImage(name)
	uiShopWindow().ChildByID(3805).DrawData().BgImageHnd = noxrender.ImageHandle(memmap.WordPtr(h))
	return h
}
func uiShopStart(name *uint16, greeting string, typ uint32) int {
	child := uiShopWindow().ChildByID(3810)
	sessionQuitHide()
	*uiShopWord(1098624) = 1
	*uiShopWord(1098628) = 1
	w := uiShopWindow()
	w.Show()
	uiWindowEnable(w, 1)
	w.ShowModal()
	*uiShopWord(1098612) = uint32(Nox_client_getRenderGUI())
	Nox_client_setRenderGUI(0)
	uiInventoryOpenWindow()
	if name == nil {
		name = uiTradeTextAt(1107044)
	}
	uiTradeStoreText(1097300, 26, alloc.GoString16(name))
	uiTradeSetText(child, uiTradeTextAt(1097300))
	if greeting != "" {
		v, _ := GetServer().S().Strings().GetVariantInFile(strman.ID(greeting), "GUIShop.c")
		*uiShopWord(1098604) = uiInventoryPointer(unsafe.Pointer(alloc.InternCString16(v.Str)))
		*uiShopWord(1098608) = uiInventoryPointer(unsafe.Pointer(alloc.InternCString(v.Str2)))
	} else {
		*uiShopWord(1098604) = 0
		*uiShopWord(1098608) = 0
	}
	uiShopPicture(typ)
	if *uiShopWord(1098608) != 0 {
		Dialogs.PlayFile(alloc.GoString((*byte)(memmap.WordPtr(*uiShopWord(1098608)))), 100)
	}
	*uiShopWord(1107036) = 0
	return gui.EventRespInt(uiInventoryWindowValue(*uiShopWord(1098580)).Func94(gui.AsWindowEvent(16394, uintptr(*uiShopWord(1098592)), 0)))
}
func uiShopDraw() int {
	p := image.Pt(int(nox_win_width)-640, int(nox_win_height)-480)
	uiMeterImage(*uiShopWord(1098400), p)
	switch uiShopMode() {
	case 2:
		uiMeterImage(*uiShopWord(1098404), p)
		uiShopDrawStock()
	case 3:
		uiMeterImage(*uiShopWord(1098408), p)
		uiShopDrawText(p, 3)
	case 1:
		uiMeterImage(*uiShopWord(1098412), p)
		uiShopDrawText(p, 1)
	case 4:
		uiMeterImage(*uiShopWord(1098416), p)
		uiShopDrawText(p, 4)
	}
	return 1
}
func uiShopDrawText(p image.Point, mode int) int {
	w := uiShopWindow().ChildByID(3806)
	pos := uiWindowPosition(w)
	size := w.SizeVal
	if mode != 3 {
		nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
	}
	uiMeterImage(*uiShopWord(1098456), p)
	off := uintptr(1098604)
	id := ""
	if mode == 3 {
		off = 1098596
		id = "SellInstructions"
	} else if mode == 4 {
		off = 1098600
		id = "RepairInstructions"
	}
	if *uiShopWord(off) == 0 && id != "" {
		*uiShopWord(off) = uiInventoryPointer(unsafe.Pointer(alloc.InternCString16(uiShopString(id))))
	}
	ptr := *uiShopWord(off)
	if ptr == 0 {
		return 0
	}
	return nox_xxx_drawStringWrap_43FAF0(nil, (*wchar2_t)(memmap.WordPtr(ptr)), pos.X+8, pos.Y+8, size.X-16, size.Y-16)
}
func uiShopDrawStock() int {
	gold := uint32(sub_4674A0())
	objectRenderSaveClip()
	x, top := int(int32(*uiShopWord(1098380))), int(int32(*uiShopWord(1098384)))
	bottom := int(int32(*uiShopWord(1098392)))
	uiRenderCopyRect(x, top, int(int32(*uiShopWord(1098388)-*uiShopWord(1098380))), int(int32(*uiShopWord(1098392)-*uiShopWord(1098384))))
	y := int(int32(uint32(top) - *uiShopWord(1107036)))
	height := nox_xxx_guiFontHeightMB_43F320(nil)
	grid := uiShopGrid()
	drawText := func(v uint32, p image.Point) {
		s, free := alloc.CString16(fmt.Sprint(int32(v)))
		defer free()
		nox_xxx_drawStringWrap_43FAF0(nil, (*wchar2_t)(unsafe.Pointer(s)), p.X, p.Y, 320, 0)
	}
	for row := 0; row < 10; row++ {
		if y > top-50 {
			uiMeterImage(*uiShopWord(1098424 + uintptr(4*(row%2))), image.Pt(x, y))
			for col := 0; col < 6; col++ {
				c := &grid[col*10+row]
				if c.Count == 0 {
					continue
				}
				cx := x + col*50
				c.Drawable.PosVec = image.Pt(cx+25, y+25)
				c.Drawable.CallDraw(&uiShopViewport)
				if gold < c.Value {
					nox_client_drawRectFilledAlpha_49CF10(cx, y, 50, 50)
				}
				if c.Count > 1 {
					nox_xxx_drawSetTextColor_434390(int(nox_color_white_2523948))
					drawText(c.Count, image.Pt(cx+5, y+5))
				}
				nox_xxx_drawSetTextColor_434390(int(nox_color_yellow_2589772))
				drawText(c.Value, image.Pt(cx+5, y-height+45))
			}
		}
		y += 50
		if y >= bottom {
			break
		}
	}
	return objectRenderRestoreClip()
}
