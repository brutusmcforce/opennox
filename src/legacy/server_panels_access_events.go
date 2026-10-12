package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func serverPanelsAccessProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsAccessEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsAccessEdit(child *gui.Window) {
	data := serverPanelsSettings()
	admission := data[100:]
	p := (*uint16)(unsafe.Pointer(uintptr(teamUIEvent(child, 16413, 0, 0))))
	text := alloc.GoString16(p)
	n := 0
	if text != "" {
		n = serverPanelsParseNumber(text)
		if n < 0 {
			n = 0
		}
	}
	switch child.ID() {
	case 10104:
		count := alloc.StrLen(p) + 1
		copy(unsafe.Slice((*uint16)(unsafe.Pointer(&data[78])), count), unsafe.Slice(p, count))
	case 10126, 10128:
		if n > 14 {
			n = 14
			serverOptionsSetText(child, 16414, "14", -1)
		}
		if child.ID() == 10126 {
			admission[1] = admission[1]&0xf0 | byte(n)
		} else {
			admission[1] = admission[1]&15 | byte(n<<4)
		}
	case 10130:
		binary.LittleEndian.PutUint16(admission[5:], uint16(n))
	case 10132:
		binary.LittleEndian.PutUint16(admission[7:], uint16(n))
	case 10133:
		if n < 1 {
			n = 1
			serverOptionsSetText(child, 16414, "1", -1)
		} else if n > 32 {
			n = 32
			serverOptionsSetText(child, 16414, "32", -1)
		}
		Nox_xxx_servSetPlrLimit_409F80(n)
		admission[4] = byte(n)
	case 10136:
		serverConfigPasswordSet((*uint16)(unsafe.Pointer(unsafe.Pointer(p))))
	}
}
func serverPanelsAccessShowList(allowed bool) {
	w := serverPanelsWindow(1045516)
	serverPanelsWindow(1045532).SetHidden(!allowed)
	serverPanelsHide(w, 10188, 10190, !allowed)
	serverPanelsWindow(1045528).SetHidden(allowed)
	serverPanelsHide(w, 10185, 10187, allowed)
}
func serverPanelsAccessEvent(_ *gui.Window, event int, arg uintptr, value int) int {
	w := serverPanelsWindow(1045516)
	child := (*gui.Window)(unsafe.Pointer(arg))
	data := serverPanelsSettings()
	admission := data[100:]
	switch event {
	case 16387:
		child = w.ChildByID(uint(value))
		if child != nil && uint16(arg) != 1 {
			serverPanelsAccessEdit(child)
		}
	case 16415:
		serverPanelsAccessEdit(child)
	case 16400:
		switch child.ID() {
		case 10123:
			admission[0] ^= byte(uint32(1) << uint(uint32(value)&31))
		case 10200:
			serverPanelsEnable(w, 10191, 10192, bool2int(serverPanelsAccessSelected()))
		}
	case 16391:
		id := child.ID()
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		switch id {
		case 10102:
			admission[0] ^= 0x10
			if admission[0]&0x10 != 0 {
				uiWindowEnable(w.ChildByID(10206), 1)
			} else {
				c := w.ChildByID(10206)
				uiWindowEnable(c, 0)
				c.DrawData().Field0 &^= 4
				w.ChildByID(10207).DrawData().Field0 |= 4
				serverPanelsAccessShowList(false)
			}
		case 10103:
			uiWindowEnable(serverPanelsWindow(1045556), bool2int(serverPanelsWindow(1045524).DrawData().Field0&4 == 0))
			admission[0] ^= 0x20
		case 10112:
			p := unsafe.Pointer(uintptr(teamUIEvent(serverPanelsWindow(1045540), 16413, 0, 0)))
			if serverPanelsWindow(1045528).Flags.IsHidden() {
				serverConfigAllowedAdd((*uint16)(unsafe.Pointer(p)))
			} else {
				serverConfigBlockedAdd(int32(0), (*uint16)(unsafe.Pointer(p)), (*byte)(unsafe.Pointer(nil)))
			}
			teamUIEvent(serverPanelsWindow(1045540), 16414, uintptr(memmap.PtrOff(0x5D4594, 1045600)), 0)
		case 10113:
			if serverPanelsWindow(1045520).DrawData().Field0&4 != 0 {
				list := serverPanelsWindow(1045532)
				index := teamUIEvent(list, 16404, 0, 0)
				teamUIEvent(list, 16398, uintptr(index), 0)
				serverConfigAllowedRemove(int32(index))
			} else {
				list := serverPanelsWindow(1045528)
				index := teamUIEvent(list, 16404, 0, 0)
				teamUIEvent(list, 16398, uintptr(index), 0)
				serverConfigBlockedRemove(int32(index))
			}
		case 10124:
			admission[2] ^= 0x80
		case 10125, 10127, 10129, 10131:
			index := (id - 10125) / 2
			check := serverPanelsWindow(1045560 + 4*uintptr(index))
			entry := serverPanelsWindow(1045576 + 4*uintptr(index))
			checked := check.DrawData().Field0&4 != 0
			uiWindowEnable(entry, bool2int(!checked))
			if checked {
				switch id {
				case 10125:
					admission[1] |= 15
				case 10127:
					admission[1] |= 0xf0
				case 10129:
					binary.LittleEndian.PutUint16(admission[5:], 0xffff)
				case 10131:
					binary.LittleEndian.PutUint16(admission[7:], 0xffff)
				}
				return 0
			}
			text := serverOptionsGetText(entry, 16413, 0)
			if text == "" {
				return 0
			}
			n := serverPanelsParseNumber(text)
			switch id {
			case 10125:
				admission[1] = admission[1]&0xf0 | byte(n)
			case 10127:
				admission[1] = admission[1]&15 | byte(n)
			case 10129:
				binary.LittleEndian.PutUint16(admission[5:], uint16(n))
			case 10131:
				binary.LittleEndian.PutUint16(admission[7:], uint16(n))
			}
		case 10191, 10192:
			list := serverPanelsWindow(1045536)
			selection := uiListSelection(uiListData(list))
			for i := 0; selection[i] >= 0; i++ {
				name := unsafe.Pointer(uintptr(teamUIEvent(list, 16406, uintptr(selection[i]), 0)))
				player := playerStateByName((*uint16)(name))
				if player == nil || player.PlayerInd == 31 {
					continue
				}
				if noxflags.HasGame(4096) {
					Sub_4DCFB0(player.PlayerUnit)
				} else if id == 10191 {
					Nox_xxx_playerCallDisconnect_4DEAB0(ntype.PlayerInd(player.PlayerInd), 4)
				} else {
					Nox_xxx_playerDisconnByPlrID_4DEB00(ntype.PlayerInd(player.PlayerInd))
				}
				if id == 10192 {
					serverConfigBlockedAdd(int32(0), (*uint16)(unsafe.Pointer(name)), (*byte)(unsafe.Pointer(unsafe.Add(unsafe.Pointer(player), 2112))))
				}
			}
		case 10206, 10207:
			allowed := id == 10206
			serverPanelsAccessShowList(allowed)
			if allowed {
				*serverPanelsWord(1045596) = *serverPanelsWord(1045532)
			} else {
				*serverPanelsWord(1045596) = *serverPanelsWord(1045528)
			}
		}
	}
	return 0
}
