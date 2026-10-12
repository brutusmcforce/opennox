package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"slices"
	"strconv"
	"unsafe"
)

func serverPanelsAccessOpen(parent *gui.Window) int {
	if *serverPanelsWord(1045516) != 0 {
		return 0
	}
	w := Nox_new_window_from_file(serverPanelsResource(127824), serverPanelsAccessProc)
	*serverPanelsWord(1045516) = uint32(serverOptionsPtr(w))
	if w == nil {
		return 0
	}
	GetClient().R2().SetTabWidth(100)
	w.SetParent(parent)
	for _, v := range []struct {
		off uintptr
		id  uint
	}{{1045520, 10102}, {1045524, 10103}, {1045532, 10109}, {1045528, 10105}, {1045536, 10200}, {1045540, 10111}, {1045544, 10112}, {1045548, 10113}, {1045556, 10104}, {1045560, 10125}, {1045564, 10127}, {1045568, 10129}, {1045572, 10131}, {1045576, 10126}, {1045580, 10128}, {1045584, 10130}, {1045588, 10132}, {1045552, 10123}, {1045592, 10133}} {
		*serverPanelsWord(v.off) = uint32(serverOptionsPtr(w.ChildByID(v.id)))
	}
	en := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISlider")))
	lit := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISliderLit")))
	for _, v := range []struct{ list, slider, up, down uint }{{10109, 10190, 10188, 10189}, {10105, 10187, 10185, 10186}, {10200, 10203, 10201, 10202}} {
		list, slider, up, down := w.ChildByID(v.list), w.ChildByID(v.slider), w.ChildByID(v.up), w.ChildByID(v.down)
		slider.Field100Ptr.SizeVal = image.Pt(16, 10)
		gui.ButtonSetImage(slider, nil, nil, en, lit, lit)
		slider.DrawData().Window = list
		up.DrawData().Window = list
		down.DrawData().Window = list
		data := uiListData(list)
		data.Field_9 = slider.C()
		data.Field_7 = up.C()
		data.Field_8 = down.C()
	}
	serverPanelsAccessPopulate()
	serverPanelsAccessMeasure()
	if noxflags.HasGame(1) {
		w.SetDraw(serverPanelsAccessDraw)
	}
	return int(serverOptionsPtr(w))
}
func serverPanelsAccessMeasure() int {
	w := serverPanelsWindow(1045516).ChildByID(10123)
	r := GetClient().R2()
	height := r.FontHeight(w.DrawData().Font()) + 1
	w.EndPos.Y = w.Off.Y + 4*height + 2
	w.SizeVal.Y = 4*height + 2
	width := 0
	for _, key := range []string{"WARRIOR", "WIZARD", "CONJURER"} {
		if x := r.GetStringSizeWrapped(w.DrawData().Font(), serverPanelsText("access.c", key), 0).X; x > width {
			width = x
		}
	}
	w.SizeVal.X = width + 7
	w.EndPos.X = w.Off.X + w.SizeVal.X
	return w.EndPos.X
}
func serverPanelsAccessPopulate() uintptr {
	w := serverPanelsWindow(1045516)
	data := serverPanelsSettings()
	serverOptionsSetText(w.ChildByID(10136), 16414, Nox_xxx_sysopGetPass_40A630(), 0)
	for _, v := range []struct {
		off          int
		entry, check uint
	}{{105, 10130, 10129}, {107, 10132, 10131}} {
		n := binary.LittleEndian.Uint16(data[v.off:])
		if n != 0xffff {
			uiWindowEnable(w.ChildByID(v.entry), 1)
			w.ChildByID(v.check).DrawData().Field0 |= 4
			serverOptionsSetText(w.ChildByID(v.entry), 16414, strconv.Itoa(int(n)), 0)
		}
	}
	if int8(data[102]) < 0 {
		w.ChildByID(10124).DrawData().Field0 |= 4
	}
	if data[100]&0x20 != 0 {
		uiWindowEnable(w.ChildByID(10104), 1)
		w.ChildByID(10103).DrawData().Field0 |= 4
	}
	teamUIEvent(w.ChildByID(10104), 16414, uintptr(unsafe.Pointer(&data[78])), 0)
	closed := w.ChildByID(10102)
	if Sub_4D6F30() != 0 {
		uiWindowEnable(closed, 0)
	} else {
		uiWindowEnable(closed, 1)
		if data[100]&0x10 != 0 {
			closed.DrawData().Field0 = 4
		} else {
			uiWindowEnable(w.ChildByID(10206), 0)
		}
	}
	w.ChildByID(10207).DrawData().Field0 |= 4
	*serverPanelsWord(1045596) = *serverPanelsWord(1045528)
	classes := w.ChildByID(10123)
	for _, key := range []string{"WARRIOR", "WIZARD", "CONJURER"} {
		serverOptionsSetText(classes, 16397, serverPanelsText("access.c", key), -1)
	}
	if data[100]&0x10 != 0 {
		w.ChildByID(10109).SetHidden(false)
		w.ChildByID(10105).SetHidden(true)
	}
	if data[100]&7 != 0 {
		sel := uiListSelection(uiListData(classes))
		n := 0
		for i := 0; i < 3; i++ {
			if data[100]&(1<<uint(i)) != 0 {
				sel[n] = int32(i)
				n++
				sel[n] = -1
			}
		}
	}
	serverOptionsSetText(w.ChildByID(10133), 16414, strconv.Itoa(int(data[104])), 0)
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		if p.PlayerInd != 31 || !noxflags.HasEngine(noxflags.EngineNoRendering) {
			serverPanelsPlayerAdd(&p.NameFinal[0])
		}
	}
	return serverPanelsAccessRefresh()
}
func serverPanelsAccessDraw(w *gui.Window, d *gui.WindowData) int {
	serverPanelsBackground(w, d)
	entry, add := serverPanelsWindow(1045540), serverPanelsWindow(1045544)
	if entry.Flags.IsEnabled() {
		enabled := serverOptionsGetText(entry, 16413, 0) != ""
		if add.Flags.IsEnabled() != enabled {
			uiWindowEnable(add, bool2int(enabled))
		}
	}
	selected := teamUIEvent(serverPanelsWindow(1045596), 16404, 0, 0)
	remove := serverPanelsWindow(1045548)
	if remove.Flags.IsEnabled() != (selected >= 0) {
		uiWindowEnable(remove, bool2int(selected >= 0))
	}
	return 1
}
func serverPanelsAccessClose(destroy bool) int {
	if destroy {
		serverPanelsWindow(1045516).Destroy()
	}
	*serverPanelsWord(1045516) = 0
	return 0
}
func serverPanelsAccessRefresh() uintptr {
	w := serverPanelsWindow(1045516)
	if w == nil {
		return 0
	}
	serverOptionsSetText(serverPanelsWindow(1045592), 16414, strconv.Itoa(int(serverConfigLimitGet())), 0)
	if noxflags.HasGame(1) {
		allowed, blocked := serverPanelsWindow(1045532), serverPanelsWindow(1045528)
		teamUIEvent(blocked, 16399, 0, 0)
		teamUIEvent(allowed, 16399, 0, 0)
		for node := listNext(serverConfigAllowedHead()); node != nil; node = listNext(node) {
			teamUIEvent(allowed, 16397, uintptr(unsafe.Add(unsafe.Pointer(node), 12)), ^uintptr(0))
		}
		for node := listNext(serverConfigBlockedHead()); node != nil; node = listNext(node) {
			p := unsafe.Add(unsafe.Pointer(node), 12)
			if *(*byte)(unsafe.Add(unsafe.Pointer(node), 72)) != 0 {
				teamUIEvent(blocked, 16397, uintptr(p), ^uintptr(0))
			} else {
				serverOptionsSetText(blocked, 16397, "*"+alloc.GoString16((*uint16)(p)), -1)
			}
		}
	}
	return 0
}
func serverPanelsPlayerAdd(name *uint16) int {
	if *serverPanelsWord(1045516) == 0 {
		return 0
	}
	return teamUIEvent(serverPanelsWindow(1045516).ChildByID(10200), 16397, uintptr(unsafe.Pointer(name)), 3)
}
func serverPanelsPlayerLookup(name *uint16) int {
	list := serverPanelsWindow(1045536)
	n := int(int16(uiListData(list).Field_11_0))
	target := unsafe.Slice(name, alloc.StrLen(name))
	for i := 0; i < n; i++ {
		p := (*uint16)(unsafe.Pointer(uintptr(teamUIEvent(list, 16406, uintptr(i), 0))))
		if slices.Equal(unsafe.Slice(p, alloc.StrLen(p)), target) {
			return i
		}
	}
	return -1
}
func serverPanelsPlayerRemove(name *uint16) {
	if *serverPanelsWord(1045516) == 0 {
		return
	}
	index := serverPanelsPlayerLookup(name)
	if index < 0 {
		return
	}
	teamUIEvent(serverPanelsWindow(1045536), 16398, uintptr(index), 0)
	if !serverPanelsAccessSelected() {
		serverPanelsEnable(serverPanelsWindow(1045516), 10191, 10192, 0)
	}
}
func serverPanelsSelectedPlayer(list *gui.Window, index int32) *server.Player {
	p := unsafe.Pointer(uintptr(teamUIEvent(list, 16406, uintptr(index), 0)))
	return playerStateByName((*uint16)(p))
}
func serverPanelsAccessSelected() bool {
	list := serverPanelsWindow(1045536)
	sel := uiListSelection(uiListData(list))
	for _, index := range sel {
		if index < 0 {
			break
		}
		p := serverPanelsSelectedPlayer(list, index)
		if p != nil && p.PlayerInd != 31 {
			return true
		}
	}
	return false
}
