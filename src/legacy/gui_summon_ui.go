package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"image"
	"unsafe"
)

func summonWin(off uintptr) *gui.Window { return bookWindow(*summonWord(off)) }
func summonText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "guisumn.c")
}
func summonState() *byte { return (*byte)(unsafe.Pointer(summonWord(1321200))) }
func summonCreate() int {
	*summonWord(1321004) = 0
	*summonWord(1321000) = ^uint32(144)
	*summonWord(1320988) = uint32(nox_win_width) - 95
	*summonWord(1320992) = *summonWord(1321000)
	g := GetClient().Cli().GUI
	root := g.NewWindowRaw(nil, 8, int(nox_win_width)-95, -145, 87, 115, nil)
	*summonWord(1321032) = quickbarPointer(root.C())
	root.SetAllFuncs(bookEvent(func(*gui.Window, uint32, uint32) int { return 0 }), func(*gui.Window, *gui.WindowData) int { return 1 }, nil)
	box := g.NewWindowRaw(root, 136, 5, 38, 76, 76, nil)
	*summonWord(1321036) = quickbarPointer(box.C())
	box.SetAllFuncs(bookEvent(summonBoxEvent), func(w *gui.Window, _ *gui.WindowData) int { return summonDraw(w) }, bookTooltipKey(2))
	box.DrawData().SetTooltip(GetServer().S().Strings(), summonText("ToolTipSummon"))
	*summonWord(1320996) = uiMeterLoadImage("CreatureCageBottom")
	top := g.NewWindowRaw(box, 160, 0, 0, 1, 1, nil)
	top.DrawData().BgImageHnd = quickbarImage("CreatureCageTop")
	top.DrawData().ImgPtVal = image.Pt(-5, -38)
	small := g.NewWindowRaw(root, 8, 19, 0, 48, 39, nil)
	small.SetAllFuncs(bookEvent(func(*gui.Window, uint32, uint32) int { return 1 }), func(*gui.Window, *gui.WindowData) int { return 1 }, nil)
	for _, a := range []struct {
		off  uintptr
		name string
	}{{1321008, "HuntButtonLit"}, {1321012, "HuntButton"}, {1321016, "GuardButtonLit"}, {1321020, "GuardButton"}, {1321024, "EscortButtonLit"}, {1321028, "EscortButton"}} {
		*summonWord(a.off) = uiMeterLoadImage("CreatureCage" + a.name)
	}
	big := g.NewWindowRaw(nil, 168, int(*summonWord(1320988))+27, int(*summonWord(1320992))+12, 34, 34, nil)
	*summonWord(1321040) = quickbarPointer(big.C())
	big.DrawData().Style |= 256
	big.SetAllFuncs(bookEvent(summonBigEvent), nil, bookTooltipKey(3))
	big.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(*summonWord(1321028))))
	big.DrawData().SelImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(*summonWord(1321024))))
	big.DrawData().HlImageHnd = big.DrawData().SelImageHnd
	big.DrawData().ImgPtVal = image.Pt(-27, -12)
	*summonState() = 0
	root.Hide()
	big.Hide()
	for i := int32(0); i < 4; i++ {
		summonAt(i).Active = 0
	}
	summonClearGrid()
	*summonWord(1321044) = 0
	*summonWord(1321204) = 0
	*summonWord(1321196) = 0
	return 1
}
func summonSetCommand(command uint32) int {
	*memmap.PtrUint32(0x587000, 184448) = command
	var normal, lit uintptr
	switch command {
	case 3:
		normal, lit = 1321020, 1321016
	case 4:
		normal, lit = 1321028, 1321024
	case 5:
		normal, lit = 1321012, 1321008
	default:
		return int(command - 5)
	}
	w := summonWin(1321040)
	if w == nil {
		return -2
	}
	d := w.DrawData()
	d.BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(*summonWord(normal))))
	d.SelImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(*summonWord(lit))))
	d.HlImageHnd = d.SelImageHnd
	return 0
}
func summonClose() uint32 {
	if w := summonWin(1321044); w != nil {
		w.Destroy()
	}
	*summonWord(1321044) = 0
	*summonWord(1321204) = 0
	bookSound(920)
	summonWin(1321040).DrawData().Field0 &^= 2
	return *summonWord(1321040)
}
func summonBigEvent(_ *gui.Window, event, arg uint32) int {
	switch event {
	case 5, 6:
		return 1
	case 7:
		*summonWord(1321204) = 0
		summonMenu(bookPoint(arg))
		summonWin(1321040).DrawData().Field0 |= 2
		return 1
	case 17:
		summonWin(1321040).DrawData().Field0 |= 2
	case 18:
		summonWin(1321040).DrawData().Field0 &^= 2
	}
	return 0
}
func summonBoxEvent(w *gui.Window, event, arg uint32) int {
	switch event {
	case 5, 6:
		return 1
	case 7:
		p := bookPoint(arg)
		off := w.GlobalPos()
		xy := [2]int32{int32((p.X - off.X) / 38), int32((p.Y - off.Y) / 38)}
		*summonWord(1321204) = summonGet(&xy)
		if *summonWord(1321204) != 0 {
			summonMenu(p)
		}
		return 1
	}
	return 0
}
func summonOrder(record *summonRecord, command uint32) {
	var msg [4]byte
	msg[0] = 0x78
	if record != nil {
		binary.LittleEndian.PutUint16(msg[1:], uint16(record.Code))
	} else if command == 1 {
		return
	}
	msg[3] = byte(command)
	GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])
	summonClose()
	if command == 0 {
		bookSound(777)
	} else {
		bookSound(898)
	}
}
func summonCommandEvent(w *gui.Window, event, _ uint32) int {
	switch event {
	case 5, 6:
		return 1
	case 7:
		command := *summonCommandWord(w)
		if command != 2 && (*summonWord(1321204) != 0 || command != 1) {
			summonOrder((*summonRecord)(unsafe.Pointer(uintptr(*summonWord(1321204)))), command)
		}
		return 1
	}
	return 0
}
func summonMenu(p image.Point) {
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	*summonMenuWidth() = 0
	for i := uintptr(0); i < 6; i++ {
		if i == 2 {
			continue
		}
		text := summonText(GoStringP(*memmap.PtrPtr(0x587000, 184344+i*4)))
		sz := r.GetStringSizeWrapped(font, text, int(nox_win_width))
		if *summonMenuWidth() < uint32(sz.X) {
			*summonMenuWidth() = uint32(sz.X)
		}
	}
	*summonMenuWidth() += 8
	line := r.FontHeight(font) + 2
	width := int(*summonMenuWidth())
	height := 5*line + 12
	x, y := p.X-width/2, p.Y-height/2
	if x < 0 {
		x = 0
	} else if uint32(width+x) >= uint32(nox_win_width) {
		x = int(nox_win_width) - width - 1
	}
	if y < 0 {
		y = 0
	} else if height+y >= int(nox_win_height) {
		y = int(nox_win_height) - height - 1
	}
	g := GetClient().Cli().GUI
	menu := g.NewWindowRaw(nil, 40, x, y, width, height, nil)
	*summonWord(1321044) = quickbarPointer(menu.C())
	menu.SetAllFuncs(nil, func(w *gui.Window, _ *gui.WindowData) int { return Sub_4C26F0(w) }, nil)
	menu.ShowModal()
	yy := 0
	for i := uint32(0); i < 6; i++ {
		if i == 2 {
			continue
		}
		w := g.NewWindowRaw(menu, 8, 0, yy, width, line+1, nil)
		w.SetAllFuncs(bookEvent(summonCommandEvent), func(w *gui.Window, _ *gui.WindowData) int { return summonDrawMenu(w) }, nil)
		*summonCommandWord(w) = i
		yy += line + 2
	}
	bookSound(791)
}
func summonCommandTooltip() int {
	var name string
	switch *memmap.PtrUint32(0x587000, 184448) {
	case 3:
		name = "ccs:GUARD"
	case 4:
		name = "ccs:ESCORT"
	case 5:
		name = "ccs:HUNT"
	default:
		return 1
	}
	Nox_xxx_cursorSetTooltip_4776B0(summonText(name))
	return 1
}
func summonSlotTooltip(w *gui.Window, p image.Point) unsafe.Pointer {
	off := w.GlobalPos()
	xy := [2]int32{int32((p.X - off.X) / 38), int32((p.Y - off.Y) / 38)}
	value := summonGet(&xy)
	if value == 0 {
		return nil
	}
	r := (*summonRecord)(unsafe.Pointer(uintptr(value)))
	return unsafe.Pointer(nox_get_thing_pretty_name(int(r.Type)))
}
func summonAdd(code, typ uint32, quiet bool) byte {
	for r := summonFirst(); r != nil; r = summonNext(r) {
		if r.Code == code && r.Type == typ {
			return byte(summonAddress(r))
		}
	}
	r := summonAllocate()
	if r == nil {
		return 0
	}
	r.Code = code
	r.Type = typ
	r.Size = summonClass(int(typ))
	summonLayout()
	ret := *summonState()
	if ret == 0 || ret == 3 {
		if !quiet {
			bookSound(801)
		}
		*summonState() = 1
		summonWin(1321032).Show()
		summonWin(1321040).Show()
		ret = 0
	}
	*summonWord(1321196)++
	return ret
}
func summonRemove(code uint32, quiet bool) {
	r := summonFind(code)
	if r == nil {
		return
	}
	if summonAddress(r) == *summonWord(1321204) {
		summonClose()
	}
	if dr := GetClient().Cli().Objs.ByNetCodeDynamic(int(code)); dr != nil {
		dr.ObjFlags &^= 0x40000000
	}
	summonPaint((*[2]int32)(unsafe.Pointer(&r.X)), int32(r.Size), 0)
	summonDeactivate(r)
	summonLayout()
	*summonWord(1321196)--
	if *summonWord(1321196) == 0 {
		if !quiet {
			bookSound(802)
		}
		*summonState() = 3
	}
}

func summonCommandWord(w *gui.Window) *uint32 { return (*uint32)(unsafe.Add(w.C(), 32)) }
