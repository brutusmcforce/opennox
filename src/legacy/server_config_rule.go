package legacy

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/server"
)

var serverConfigRuleWords [7]uint32

func serverConfigRuleWindow(index int) *gui.Window {
	return (*gui.Window)(memmap.WordPtr(serverConfigRuleWords[index]))
}
func serverConfigRuleOpen(parent *gui.Window, settings *byte) int {
	w := Nox_new_window_from_file("rulelist.wnd", serverConfigRuleProc)
	serverConfigRuleWords[0] = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	for i := 1; i < 7; i++ {
		serverConfigRuleWords[i] = memmap.PtrWord(w.ChildByID(uint(10169 + i)).C())
	}
	w.ChildByID(10176).SetDraw(serverConfigRuleDraw)
	w.SetParent(parent)
	list, slider, up, down := serverConfigRuleWindow(1), w.ChildByID(10179), w.ChildByID(10177), w.ChildByID(10178)
	en := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISlider")))
	lit := noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("UISliderLit")))
	slider.Field100Ptr.SizeVal = image.Pt(16, 10)
	gui.ButtonSetImage(slider, nil, nil, en, lit, lit)
	slider.DrawData().Window = list
	up.DrawData().Window = list
	down.DrawData().Window = list
	data := uiListData(list)
	data.Field_9 = slider.C()
	data.Field_7 = up.C()
	data.Field_8 = down.C()
	serverConfigRulePopulate(settings)
	return int(memmap.PtrWord(w.C())) // stored as a pointer word
}
func serverConfigRulePopulate(settings *byte) int {
	list := serverConfigRuleWindow(1)
	teamUIEvent(list, 16399, 0, 0)
	name := alloc.GoString(settings)
	paths, err := filepath.Glob(ifs.Normalize("maps\\" + name + "\\*.rul"))
	if err != nil || len(paths) == 0 {
		return -1
	}
	matched := false
	for _, path := range paths {
		// libc glob excludes dotfiles unless the pattern explicitly starts with a dot.
		filename := filepath.Base(path)
		if strings.HasPrefix(filename, ".") {
			continue
		}
		matched = true
		if st, err := os.Stat(path); err == nil && st.IsDir() {
			continue
		}
		title := strings.TrimSuffix(filename, filepath.Ext(filename))
		if serverConfigEqualFoldBytes(name, title) || serverConfigEqualFoldBytes("user", title) {
			continue
		}
		wide := ruleWideBytes([]byte(title))
		teamUIEvent(list, 16397, uintptr(unsafe.Pointer(&wide[0])), ^uintptr(0))
	}
	if !matched {
		return -1
	}
	return 1
}
func serverConfigRuleDraw(w *gui.Window, d *gui.WindowData) int {
	pos := uiWindowPosition(w)
	if w.Flags.Has(gui.StatusImage) {
		bookDrawImage(memmap.PtrWord(unsafe.Pointer(d.BgImageHnd)), pos)
	} else if d.BgColorVal != 0x80000000 {
		r := GetClient().R2()
		r.DrawRectFilledOpaque(pos.X, pos.Y, w.SizeVal.X, w.SizeVal.Y, r.Data().Color2())
	}
	selected := teamUIEvent(serverConfigRuleWindow(1), 16404, 0, 0) >= 0
	for _, i := range []int{4, 5, 6} {
		child := serverConfigRuleWindow(i)
		if !selected {
			if child.Flags&8 != 0 {
				uiWindowEnable(child, 0)
			}
		} else if child.Flags&8 == 0 && (i != 5 || !noxflags.HasGame(49152)) {
			uiWindowEnable(child, 1)
		}
	}
	entry := (*uint16)(unsafe.Pointer(uintptr(teamUIEvent(serverConfigRuleWindow(2), 16413, 0, 0))))
	focus := GetClient().Cli().GUI.Focused()
	if focus != nil && focus.ID() == 10171 {
		on := entry != nil && *entry != 0
		button := serverConfigRuleWindow(3)
		if on != (button.Flags&8 != 0) {
			uiWindowEnable(button, bool2int(on))
		}
	}
	return 1
}
func serverConfigRuleProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverConfigRuleEvent(w, e.EventCode(), a, int(b)))
}
func serverConfigRuleText(w *gui.Window, code, row int) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(teamUIEvent(w, code, uintptr(row), 0))))
}
func serverConfigRuleEvent(_ *gui.Window, event int, arg uintptr, value int) int {
	if uint32(event) > 0x4007 {
		return 1
	}
	if event != 16391 {
		if event == 16387 {
			child := serverConfigRuleWindow(0).ChildByID(uint(value))
			if child == nil || uint16(arg) == 1 {
				return 0
			}
			// Preserve the historical narrow view of UTF-16 text for this completion check.
			p := (*byte)(unsafe.Pointer(serverConfigRuleText(child, 16413, 0)))
			if p != nil && *p != 0 && value == 10171 {
				text := alloc.GoString(p)
				if serverConfigEqualFoldBytes(text, mapNameField(serverConfigSlotCurrent())) || serverConfigEqualFoldBytes(text, "user") {
					uiWindowEnable(serverConfigRuleWindow(3), 0)
				}
			}
		}
		return 1
	}
	child := (*gui.Window)(unsafe.Pointer(arg))
	id := child.ID()
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	list, entry := serverConfigRuleWindow(1), serverConfigRuleWindow(2)
	current := func() *server.Settings2 { return (*server.Settings2)(unsafe.Pointer(serverConfigSlotCurrent())) }
	selectedText := func() *uint16 {
		index := teamUIEvent(list, 16404, 0, 0)
		return serverConfigRuleText(list, 16406, index)
	}
	fileName := func(p *uint16, off uintptr) string {
		return serverConfigNarrow(p) + alloc.GoString(memmap.PtrUint8(0x587000, off))
	}
	switch id {
	case 10172:
		text := serverConfigRuleText(entry, 16413, 0)
		name := fileName(text, 191640)
		serverOptionsRead(serverOptionsRecord(unsafe.Pointer(current())))
		ruleWrite(name, current(), nil)
		found := false
		for i := 0; i < int(int16(uiListData(list).Field_11_0)); i++ {
			other := serverConfigRuleText(list, 16406, i)
			if serverConfigWideEqualFold(text, other) {
				found = true
				break
			}
		}
		off := uintptr(1523052)
		if !found {
			teamUIEvent(list, 16397, uintptr(unsafe.Pointer(text)), ^uintptr(0))
			off = 1523056
		}
		teamUIEvent(entry, 16414, uintptr(memmap.PtrOff(0x5D4594, off)), 0)
		if !found {
			uiWindowEnable(serverConfigRuleWindow(3), 0)
		}
	case 10173:
		name := fileName(selectedText(), 191592)
		serverOptionsRead(serverOptionsRecord(unsafe.Pointer(current())))
		ruleWrite(name, current(), nil)
		teamUIEvent(list, 16403, ^uintptr(0), 0)
	case 10174:
		st := current()
		name := fileName(selectedText(), 191608)
		ruleLoad(st, name, nil, 7, st.Field52)
		serverPanelsCopySettings(unsafe.Pointer(st))
		serverOptionsSettingsLabels(serverOptionsRecord(unsafe.Pointer(current())))
		teamUIEvent(list, 16403, ^uintptr(0), 0)
		serverOptionsDirty(1)
	case 10175:
		mapName := mapNameField(serverConfigSlotCurrent())
		index := teamUIEvent(list, 16404, 0, 0)
		name := fileName(serverConfigRuleText(list, 16406, index), 191624)
		_ = ifs.Remove("maps\\" + mapName + "\\" + name)
		teamUIEvent(list, 16398, uintptr(index), 0)
	}
	return 1
}
func serverConfigWideEqualFold(a, b *uint16) bool {
	for i := uintptr(0); ; i += 2 {
		x, y := *(*uint16)(unsafe.Add(unsafe.Pointer(a), i)), *(*uint16)(unsafe.Add(unsafe.Pointer(b), i))
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
		if x == 0 {
			return true
		}
	}
}
