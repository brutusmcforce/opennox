package legacy

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func teamUIHidden(w *gui.Window, hidden bool) int {
	if w == nil {
		return -2
	}
	w.SetHidden(hidden)
	return 0
}
func teamUIHUDShow(ball bool, show int) int {
	root, visible := uintptr(1045604), uintptr(1045608)
	if ball {
		root, visible = 1045636, 1045640
	}
	w := teamUIWindow(root)
	if show != 0 && *teamUIWord(visible) != 0 && w != nil && w.Flags.IsHidden() {
		return w.ShowModal()
	}
	return teamUIHidden(w, true)
}
func teamUIHUDHide(ball bool) int {
	off := uintptr(1045608)
	if ball {
		off = 1045640
	}
	*teamUIWord(off) = 0
	return teamUIHUDShow(ball, 0)
}
func teamUIHUDDestroy(ball bool) int {
	root, visible := uintptr(1045604), uintptr(1045608)
	if ball {
		root, visible = 1045636, 1045640
	}
	if w := teamUIWindow(root); w != nil {
		w.Destroy()
	}
	*teamUIWord(root), *teamUIWord(visible) = 0, 0
	return 0
}
func teamUICTFConstruct() int {
	if teamUIWindow(1045604) != nil {
		return 1
	}
	w := Nox_new_window_from_file("GUI_CTF.wnd", nil)
	*teamUIWord(1045604) = uint32(uintptr(w.C()))
	if w == nil {
		return 0
	}
	sm := GetServer().S().Strings()
	for id := uint(8811); id <= 8826; id++ {
		child := w.ChildByID(id)
		child.SetAllFuncs(nil, teamUICTFDraw, nil)
		child.DrawData().SetTooltip(sm, sm.GetStringInFile("FlagHomeTT", "GUI_CTF.c"))
	}
	teamUIHUDShow(false, 0)
	*teamUIWord(1045632) = uiMeterLoadImage("FlagTeamBorder")
	return 1
}
func teamUIBallConstruct() int {
	if teamUIWindow(1045636) != nil {
		return 1
	}
	w := Nox_new_window_from_file("gui_fb.wnd", nil)
	*teamUIWord(1045636) = uint32(uintptr(w.C()))
	w.SetAllFuncs(nil, teamUIBallDraw, nil)
	if w == nil {
		return 0
	}
	teamUIHUDShow(true, 0)
	*teamUIWord(1045648) = uiMeterLoadImage("FlagTeamBorder")
	return 1
}
func teamUIScreen() image.Point {
	var w, h, d int
	Nox_xxx_gameGetScreenBoundaries_43BEB0_get_video_mode(&w, &h, &d)
	max := VideoGetMaxSize()
	if w > max.X {
		w = max.X
	}
	if h > max.Y {
		h = max.Y
	}
	return image.Pt(w, h)
}
func teamUICTFOpen(count byte) int8 {
	teamRuntimeObject(ClientPlayerNetCode())
	if teamUICTFConstruct() == 0 {
		return 0
	}
	w := teamUIWindow(1045604)
	screen := teamUIScreen()
	clear(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1045612), 16))
	*teamUIWord(1045608) = 1
	for id := uint(8811); id <= 8826; id++ {
		teamUIHidden(w.ChildByID(id), true)
	}
	*memmap.PtrUint8(0x5D4594, 1045628) = count
	var last *gui.Window
	for i := 0; i < int(count); i++ {
		last = w.ChildByID(uint(8811 + i))
		teamUIHidden(last, false)
	}
	size := w.SizeVal
	x := screen.X - int(uint32(size.X)) - 91
	if count <= 4 {
		x = screen.X - int(uint32(size.X)/2) - 91
	}
	w.Off.X = x
	w.EndPos.X = x + size.X
	pos := image.Pt(x, w.Off.Y)
	if count <= 4 && last == nil {
		*teamUIWord(1045608) = 0
		return int8(count)
	}
	if count <= 4 {
		pos.Y = screen.Y - 40*int(count)
	} else if count <= 8 {
		pos.Y = screen.Y - size.Y/2
	} else {
		pos.Y = screen.Y - size.Y
	}
	w.Off.Y = pos.Y
	w.EndPos.Y = pos.Y + size.Y
	return int8(w.ShowModal())
}
func teamUIBallOpen() int {
	teamRuntimeObject(ClientPlayerNetCode())
	if teamUIBallConstruct() == 0 {
		return 0
	}
	w := teamUIWindow(1045636)
	screen := teamUIScreen()
	*teamUIWord(1045640) = 1
	w.Off = image.Pt(screen.X-int(uint32(w.SizeVal.X)/3)-91, screen.Y-120)
	w.EndPos = w.Off.Add(w.SizeVal)
	*memmap.PtrUint8(0x5D4594, 1045644) = 0
	return w.ShowModal()
}
func teamUICTFSelect(id byte) int {
	root := teamUIWindow(1045604)
	for i := 0; i < int(*memmap.PtrUint8(0x5D4594, 1045628)); i++ {
		if w := root.ChildByID(uint(8811 + i)); w != nil {
			w.Flags &^= 32
		}
	}
	w := root.ChildByID(uint(id) + 8810)
	if w == nil {
		return -2
	}
	old := w.Flags
	w.Flags |= 32
	return int(old)
}
func teamUICTFTooltip(id, state byte) {
	*memmap.PtrUint8(0x5D4594, 1045611+uintptr(id)) = state
	w := teamUIWindow(1045604).ChildByID(uint(id) + 8810)
	if w == nil {
		return
	}
	var key string
	switch state {
	case 0:
		key = "FlagHomeTT"
	case 1:
		key = "TheirFlagCarriedTT"
		if w.Flags.Has(32) {
			key = "YourFlagCarriedTT"
		}
	case 2:
		key = "FlagAwayTT"
	default:
		return
	}
	sm := GetServer().S().Strings()
	w.DrawData().SetTooltip(sm, sm.GetStringInFile(strman.ID(key), "GUI_CTF.c"))
}
func teamUICTFDraw(w *gui.Window, d *gui.WindowData) int {
	pos := uiWindowPosition(w)
	h := d.BgImageHnd
	switch *memmap.PtrUint8(0x5D4594, 1045612+uintptr(byte(w.ID())-107)) {
	case 1:
		h = d.EnImageHnd
	case 2:
		h = d.DisImageHnd
	}
	if h != nil {
		bookDrawImage(uint32(uintptr(h)), pos)
	}
	if border := *teamUIWord(1045632); w.Flags.Has(32) && border != 0 {
		bookDrawImage(border, pos.Sub(image.Pt(4, 4)))
	}
	return 1
}
func teamUIBallDraw(w *gui.Window, d *gui.WindowData) int {
	pos := uiWindowPosition(w)
	if h := d.BgImageHnd; h != nil {
		bookDrawImage(uint32(uintptr(h)), pos)
	}
	if border := *teamUIWord(1045648); w.Flags.Has(32) && border != 0 {
		bookDrawImage(border, pos.Sub(image.Pt(4, 4)))
	}
	return 1
}
func teamUIMapCTF() int {
	if !teamRuntimeScan() {
		return 0
	}
	teamUICTFOpen(2)
	return 1
}
func teamUIMapBall() int8 {
	if !teamRuntimeScan() {
		return 0
	}
	teamUIBallOpen()
	return int8(objectiveBallReset(nil))
}
