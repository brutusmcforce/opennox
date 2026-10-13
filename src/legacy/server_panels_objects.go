package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/strman"
)

func serverPanelsText(file, key string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(key), file)
}
func serverPanelsEnable(w *gui.Window, first, last uint, on int) {
	for id := first; id <= last; id++ {
		uiWindowEnable(w.ChildByID(id), on)
	}
}
func serverPanelsHide(w *gui.Window, first, last uint, hidden bool) {
	for id := first; id <= last; id++ {
		w.ChildByID(id).SetHidden(hidden)
	}
}
func serverPanelsOwner(w, parent *gui.Window) {
	w.SetParent(parent)
	if parent == nil {
		parent = w
	}
	w.DrawData().Window = parent
}
func serverPanelsBackground(w *gui.Window, d *gui.WindowData) int {
	pos := uiWindowPosition(w)
	if w.Flags.Has(gui.StatusImage) {
		bookDrawImage(memmap.PtrWord(unsafe.Pointer(d.BgImageHnd)), pos)
	} else if d.BgColorVal != 0x80000000 {
		GetClient().R2().DrawRectFilledAlpha(pos.X, pos.Y, w.SizeVal.X, w.SizeVal.Y)
	}
	return 1
}
func serverPanelsObjectOpen(parent *gui.Window, kind int) int {
	w := Nox_new_window_from_file("objlst.wnd", serverPanelsObjectProc)
	*serverPanelsWord(1045468) = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	w.SetDraw(serverPanelsBackground)
	serverPanelsOwner(w, parent)
	list := w.ChildByID(1510)
	*serverPanelsWord(1045464) = memmap.PtrWord(list.C())
	serverPanelsObjectLayout()
	teamUIEvent(list, 16399, 0, 0)
	title, count := "", 0
	switch kind {
	case 0x1000000:
		*serverPanelsWord(1045460) = 0
		title = serverPanelsText("objlst.c", "Weapons")
		for i := uint(2); i < 27; i++ {
			p := uintptr(unsafe.Pointer(runtimeEquipmentLabel(false, uint32(1)<<i)))
			if p != 0 {
				teamUIEvent(list, 16397, p, ^uintptr(0))
				count++
			}
		}
	case 0x2000000:
		*serverPanelsWord(1045460) = 1
		title = serverPanelsText("objlst.c", "servopts.wnd:Armor")
		for i := uint(0); i < 26; i++ {
			p := uintptr(unsafe.Pointer(runtimeEquipmentLabel(true, uint32(1)<<i)))
			if p != 0 {
				teamUIEvent(list, 16397, p, ^uintptr(0))
				count++
			}
		}
	}
	serverOptionsSetText(list, 16385, title, 0)
	teamUIEvent(list, 16408, serverOptionsPtr(w.ChildByID(1513)), 0)
	teamUIEvent(list, 16409, serverOptionsPtr(w.ChildByID(1514)), 0)
	*serverPanelsWord(1045472 + 4*uintptr(*serverPanelsWord(1045460))) = uint32(count)
	serverPanelsObjectRefresh()
	if !noxflags.HasGame(1) || noxflags.HasGame(49152) {
		serverPanelsEnable(w, 1515, 1533, 0)
	}
	return int(memmap.PtrWord(w.C())) // stored as a pointer word
}
func serverPanelsObjectLayout() *gui.Window {
	list := serverPanelsWindow(1045464)
	height := GetClient().R2().FontHeight(list.DrawData().Font()) + 1
	uiWindowResize(list, list.SizeVal.X, 15*height+2)
	y := list.Off.Y + height + 2
	var last *gui.Window
	for id := uint(1520); id <= 1533; id++ {
		last = serverPanelsWindow(1045468).ChildByID(id)
		last.Off.Y = y
		y += height
	}
	return last
}
func serverPanelsObjectRefresh() int {
	index := uiListIndex(uiListData(serverPanelsWindow(1045464)))
	result := uint32(0)
	for id := uint(1520); id <= 1533; id++ {
		child := serverPanelsWindow(1045468).ChildByID(id)
		for (uint32(1)<<uint(uint32(index)&31))&0x33 != 0 {
			index++
		}
		if int32(index) >= int32(*serverPanelsWord(1045472 + 4*uintptr(*serverPanelsWord(1045460)))) {
			result = uint32(serverOptionsHidden(child, 1))
		} else {
			child.SetHidden(false)
			result = child.DrawData().Field0
			if serverPanelsClass(byte(index)) {
				result |= 4
			} else {
				result &^= 4
			}
			child.DrawData().Field0 = result
		}
		index++
	}
	return int(int8(byte(result)))
}
func serverPanelsObjectProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsObjectEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsObjectEvent(_ *gui.Window, event int, arg uintptr, _ int) int {
	w := serverPanelsWindow(1045468)
	list := serverPanelsWindow(1045464)
	child := (*gui.Window)(unsafe.Pointer(arg))
	if event == 16384 {
		if child == w.ChildByID(1513) || child == w.ChildByID(1514) {
			teamUIEvent(list, 16384, arg, 0)
			serverPanelsObjectRefresh()
		}
		return 0
	}
	if event != 16391 {
		return 0
	}
	id := child.ID()
	switch {
	case id == 1513 || id == 1514:
		teamUIEvent(list, 16384, arg, 0)
		serverPanelsObjectRefresh()
		return 0
	case id == 1515 || id == 1516:
		value := uint32(0)
		if id == 1515 {
			value = 0xffffffff
		}
		if *serverPanelsWord(1045460) != 0 {
			serverPanelsArmorStore(value)
		} else {
			*serverPanelsWeaponPointer() = value
		}
		serverPanelsObjectRefresh()
		serverOptionsDirty(1)
	case id >= 1520 && id <= 1533:
		index := uiListIndex(uiListData(list)) + int(id) - 1520
		p := unsafe.Pointer(uintptr(teamUIEvent(list, 16406, uintptr(index), 0)))
		on := bool2int(child.DrawData().Field0&4 == 0)
		if *serverPanelsWord(1045460) != 0 {
			mask := runtimeEquipmentMask(true, (*uint16)(p))
			serverPanelsWordMask(serverPanelsWord(1045456), mask, on)
		} else {
			value := int32(runtimeEquipmentMask(false, (*uint16)(p)))
			off := uintptr(0)
			if value > 0 {
				for {
					next := value >> 8
					if next > 0 {
						value = next
					}
					off++
					if next <= 0 {
						break
					}
				}
			}
			serverPanelsByte((*byte)(unsafe.Add(unsafe.Pointer(serverPanelsWord(1045448)), 3+off)), byte(value), on)
		}
		serverOptionsDirty(1)
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	return 0
}
