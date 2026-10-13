package legacy

import (
	"encoding/binary"
	"fmt"
	"image"
	"strconv"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

var serverPanelsRefreshCallbacks = [4]func() int{nil, serverPanelsSpellRefresh, serverPanelsObjectRefresh, serverPanelsObjectRefresh}

func serverPanelsSettings() []byte { return unsafe.Slice(memmap.PtrUint8(0x5D4594, 371516), 184) }
func serverPanelsParseNumber(text string) int {
	return int(textDecimal((*uint16)(unsafe.Pointer(alloc.InternCString16(text)))))
}
func serverPanelsAdvancedOpen(settings unsafe.Pointer) int {
	w := Nox_new_window_from_file("advanced.wnd", serverPanelsAdvancedProc)
	*serverPanelsWord(1316708) = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	w.SetParent(nil)
	w.StackPush()
	w.ShowModal()
	w.SetFunc93(func(*gui.Window, gui.WindowEvent) gui.WindowEventResp { return gui.RawEventResp(1) })
	w.Focus()
	return serverPanelsAdvancedSetup(settings)
}
func serverPanelsAdvancedSetup(settings unsafe.Pointer) int {
	w := serverPanelsWindow(1316708)
	if noxflags.HasGame(1) {
		w.ChildByID(10167).DrawData().Field0 |= 4
		*serverPanelsWord(1316704) = 0
	} else {
		w.ChildByID(10164).DrawData().Field0 |= 4
		uiWindowEnable(w.ChildByID(10167), 0)
		*serverPanelsWord(1316704) = 1
	}
	serverPanelsCopySettings(settings)
	return serverPanelsAdvancedTab()
}
func serverPanelsCopySettings(settings unsafe.Pointer) {
	serverPanelsSpellStore((*uint32)(unsafe.Add(settings, 24)))
	serverPanelsWeaponStore((*uint32)(unsafe.Add(settings, 44)))
	serverPanelsArmorStore(*(*uint32)(unsafe.Add(settings, 48)))
}
func serverPanelsAdvancedTab() int {
	parent := serverPanelsWindow(1316708)
	switch *serverPanelsWord(1316704) {
	case 0:
		*serverPanelsWord(1316712) = uint32(serverConfigRuleOpen(parent, (*byte)(unsafe.Pointer(unsafe.Pointer(&serverOptionsCurrent()[0])))))
	case 1:
		*serverPanelsWord(1316712) = uint32(serverPanelsSpellOpen(parent))
	case 2:
		*serverPanelsWord(1316712) = uint32(serverPanelsObjectOpen(parent, 0x1000000))
	case 3:
		*serverPanelsWord(1316712) = uint32(serverPanelsObjectOpen(parent, 0x2000000))
	}
	GetClient().Cli().GUI.Focus(serverPanelsWindow(1316712))
	return 0
}
func serverPanelsAdvancedProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsAdvancedEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsAdvancedEvent(_ *gui.Window, event int, arg uintptr, _ int) int {
	if event != 16391 {
		return 1
	}
	child := (*gui.Window)(unsafe.Pointer(arg))
	id := child.ID()
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	if id == 10148 {
		data := serverOptionsCurrent()
		copy(data[24:44], unsafe.Slice((*byte)(unsafe.Pointer(serverPanelsSpellPointer())), 20))
		binary.LittleEndian.PutUint32(data[44:], *serverPanelsWeaponPointer())
		binary.LittleEndian.PutUint32(data[48:], serverPanelsArmorLoad())
		serverPanelsAdvancedClose()
		return 1
	}
	if id < 10164 || id > 10167 {
		return 1
	}
	if w := serverPanelsWindow(1316712); w != nil {
		w.Destroy()
		*serverPanelsWord(1316712) = 0
	}
	tab := uint32(id - 10163)
	if id == 10167 {
		tab = 0
	}
	*serverPanelsWord(1316704) = tab
	serverPanelsAdvancedTab()
	return 1
}
func serverPanelsAdvancedClose() int {
	w := serverPanelsWindow(1316708)
	if w == nil {
		return 0
	}
	w.StackPop()
	w.Destroy()
	*serverPanelsWord(1316708) = 0
	*serverPanelsWord(1316712) = 0
	GetClient().Cli().GUI.Focus(nil)
	return 0
}
func serverPanelsAdvancedUpdate(settings unsafe.Pointer) int {
	if *serverPanelsWord(1316708) == 0 {
		return 0
	}
	return serverPanelsAdvancedRefresh(settings)
}
func serverPanelsAdvancedRefresh(settings unsafe.Pointer) int {
	serverPanelsCopySettings(settings)
	if f := serverPanelsRefreshCallbacks[*serverPanelsWord(1316704)]; f != nil {
		return f()
	}
	return 0
}
func serverPanelsAdvancedServerOpen() int {
	w := Nox_new_window_from_file(serverPanelsResource(180048), serverPanelsAdvancedServerProc)
	*serverPanelsWord(1316972) = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	w.SetParent(nil)
	w.StackPush()
	w.ShowModal()
	w.Focus()
	w.SetFunc93(func(*gui.Window, gui.WindowEvent) gui.WindowEventResp { return gui.RawEventResp(1) })
	// The legacy nil-root ChildByID query returns nil; its offset helper writes 0,0.
	w.SetPos(image.Pt(15, 80))
	list := w.ChildByID(2104)
	list.DrawData().Window = w
	list.SetFunc94(serverPanelsAdvancedServerProc)
	return serverPanelsAdvancedServerPopulate(serverPanelsSettings())
}
func serverPanelsAdvancedServerPopulate(data []byte) int {
	w := serverPanelsWindow(1316972)
	for _, v := range []struct {
		id  uint
		off int
	}{{2102, 58}, {2103, 62}} {
		c := w.ChildByID(v.id)
		if binary.LittleEndian.Uint32(data[v.off:]) != 0 {
			c.DrawData().Field0 |= 4
		} else {
			c.DrawData().Field0 &^= 4
		}
	}
	entry := w.ChildByID(2110)
	serverOptionsSetText(entry, 16414, strconv.Itoa(int(int32(binary.LittleEndian.Uint32(data[70:])))), -1)
	mode := binary.LittleEndian.Uint32(data[66:])
	if mode < 4 {
		uiWindowEnable(entry, bool2int(mode == 3))
		w.ChildByID(2106 + uint(mode)).DrawData().Field0 |= 4
	}
	return 0
}
func serverPanelsAudioSetting(value int) int {
	binary.LittleEndian.PutUint32(serverPanelsSettings()[74:], uint32(value))
	text := fmt.Sprintf(serverPanelsText("advserv.c", "AudCullDesc"), value)
	p := serverOptionsStoreText(1316716, 128, text)
	return teamUIEvent(serverPanelsWindow(1316972).ChildByID(2120), 16385, p, ^uintptr(0))
}
func serverPanelsAdvancedServerProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsAdvancedServerEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsAdvancedServerEvent(_ *gui.Window, event int, arg uintptr, value int) int {
	w := serverPanelsWindow(1316972)
	data := serverPanelsSettings()
	child := (*gui.Window)(unsafe.Pointer(arg))
	if event > 16391 {
		if event == 16393 {
			serverPanelsAudioSetting(value)
			Nox_xxx_gameSetAudioFadeoutMb_501AC0(value)
		} else if event == 16415 {
			text := serverOptionsGetText(child, 16413, 0)
			if text != "" {
				n := serverPanelsParseNumber(text)
				if n < 0 {
					n = 0
				}
				if child.ID() == 2110 {
					binary.LittleEndian.PutUint32(data[70:], uint32(n))
				}
			}
		}
		return 1
	}
	if event != 16391 {
		if event == 16387 {
			child = w.ChildByID(uint(value))
			if child == nil || uint16(arg) == 1 {
				return 0
			}
			text := serverOptionsGetText(child, 16413, 0)
			if text != "" {
				n := serverPanelsParseNumber(text)
				if n < 0 {
					n = 0
				}
				if value == 2110 {
					binary.LittleEndian.PutUint32(data[70:], uint32(n))
				}
			}
		}
		return 1
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	switch id := child.ID(); id {
	case 2102, 2103:
		off := 58 + int(id-2102)*4
		binary.LittleEndian.PutUint32(data[off:], binary.LittleEndian.Uint32(data[off:])^1)
	case 2106, 2107, 2108, 2109:
		uiWindowEnable(w.ChildByID(2110), bool2int(id == 2109))
		binary.LittleEndian.PutUint32(data[66:], uint32(id-2106))
		serverConfigRateDirtySet(int32(1))
	case 2130:
		serverPanelsAdvancedServerClose()
	}
	return 1
}
func serverPanelsAdvancedServerClose() int {
	w := serverPanelsWindow(1316972)
	if w == nil {
		return 0
	}
	w.StackPop()
	w.Destroy()
	*serverPanelsWord(1316972) = 0
	GetClient().Cli().GUI.Focus(nil)
	return 0
}
