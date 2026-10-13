package legacy

import (
	"encoding/binary"
	"strconv"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

var sessionFilterRoot, sessionFilterControls *gui.Window

func sessionFilterValues(slot int) []uint32 {
	return unsafe.Slice(memmap.PtrUint32(0x5D4594, 1193388+uintptr(slot)*44), 11)
}
func sessionFilterMode(slot int) *uint32 { return memmap.PtrUint32(0x5D4594, 1193372+uintptr(slot)*4) }
func sessionFilterConfig(slot, mode int, p unsafe.Pointer) int {
	*sessionFilterMode(slot) = uint32(mode)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&sessionFilterValues(slot)[0])), 44), unsafe.Slice((*byte)(p), 44))
	return 11 * slot
}
func sessionFilterAccept(p unsafe.Pointer) int {
	mode := *sessionFilterMode(0)
	if mode != 1 && mode != 2 {
		return 1
	}
	rec := (*Nox_gui_server_ent_t)(p)
	f := sessionFilterValues(0)
	switch mode {
	case 1:
		return bool2int(rec.StatusVal&0x30 == 0 && rec.VersionVal == 0x000f039a)
	case 2:
		if f[0] != 0 {
			ping := uint32(rec.PingVal)
			if ping > f[4] && ping != 9999 {
				return 0
			}
		}
		if f[1] != 0 && rec.StatusVal&0x10 != 0 {
			return 0
		}
		if f[2] != 0 && rec.StatusVal&0x20 != 0 {
			return 0
		}
		if int8(rec.Field_25_2) < 0 && f[3] > uint32(rec.Field_25_2&127) {
			return 0
		}
		if f[5] != 0 {
			var st server.Settings2
			copy(st.Field0[:], unsafe.Slice(&rec.MapNameBuf[0], 24)) // the map and server names
			ruleLoad(&st, "user.rul", nil, 5, binary.LittleEndian.Uint16(rec.FlagsVal[:]))
			for i, v := range st.Field24.Vals {
				if v != binary.LittleEndian.Uint32(rec.Field_33_3[4*i:]) {
					return 0
				}
			}
			for i, v := range st.Field44 {
				if v != rec.Field_38_3[i] {
					return 0
				}
			}
			if st.Field48 != binary.LittleEndian.Uint32(rec.Field_39_3[:]) {
				return 0
			}
		}
		if f[10] != 0 && rec.VersionVal != 0x000f039a {
			return 0
		}
	}
	return 1
}
func sessionFilterSave() int {
	if *sessionFilterMode(0) == 2 {
		f := sessionFilterValues(0)
		bit := func(id uint) uint32 { return sessionFilterRoot.ChildByID(id).DrawData().Field0 >> 2 & 1 }
		f[0] = bit(10028)
		f[4] = uint32(serverPanelsParseNumber(serverOptionsGetText(sessionFilterRoot.ChildByID(10031), 16413, 0)))
		f[1], f[2] = bit(10029), bit(10030)
		f[3] = 0
		if bit(10015) != 0 {
			f[3] = 0x80
			if bit(10016) != 0 {
				f[3] |= 1
			} else {
				f[3] |= 2
			}
		}
		f[5], f[10] = bit(10014), bit(10018)
	}
	if sessionFilterRoot == nil {
		return -2
	}
	sessionFilterRoot.Hide()
	return 0
}
func sessionFilterUpdateControls() {
	sessionFilterControls.Show()
	uiWindowEnable(sessionFilterControls.ChildByID(10031), int(sessionFilterControls.ChildByID(10028).DrawData().Field0>>2&1))
	serverPanelsEnable(sessionFilterControls, 10016, 10017, int(sessionFilterControls.ChildByID(10015).DrawData().Field0>>2&1))
}
func sessionFilterOpen(parent *gui.Window) *gui.Window {
	w := Nox_new_window_from_file("filter.wnd", sessionFilterEvent)
	sessionFilterRoot = w
	if w == nil {
		return nil
	}
	sessionFilterControls = w.ChildByID(10012)
	w.SetParent(parent)
	sessionFilterControls.SetParent(w)
	sessionFilterControls.SetFunc94(sessionFilterEvent)
	f := sessionFilterValues(0)
	for _, pair := range [][2]uint{{0, 10028}, {1, 10029}, {2, 10030}, {3, 10015}, {5, 10014}, {10, 10018}} {
		if f[pair[0]] != 0 {
			w.ChildByID(pair[1]).DrawData().Field0 |= 4
		}
	}
	serverOptionsSetText(w.ChildByID(10031), 16414, strconv.FormatInt(int64(int32(f[4])), 10), -1)
	resolution := uint(10016)
	if f[3]&2 != 0 {
		resolution = 10017
	}
	w.ChildByID(resolution).DrawData().Field0 |= 4
	mode := *sessionFilterMode(0)
	switch mode {
	case 0, 1:
		sessionFilterControls.Hide()
	case 2:
		sessionFilterUpdateControls()
	}
	if mode <= 2 {
		w.ChildByID(10024 + uint(mode)).DrawData().Field0 |= 4
	}
	return w
}
func sessionFilterEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 23 {
		return gui.RawEventResp(1)
	}
	if ev.EventCode() != 16391 {
		return nil
	}
	a, _ := ev.EventArgsC()
	button := (*gui.Window)(unsafe.Pointer(a))
	id := button.ID()
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	switch id {
	case 10015:
		serverPanelsEnable(sessionFilterRoot, 10016, 10017, int(^button.DrawData().Field0>>2&1))
	case 10024, 10025:
		*sessionFilterMode(0) = uint32(id - 10024)
		sessionFilterControls.Hide()
	case 10026:
		*sessionFilterMode(0) = 2
		sessionFilterUpdateControls()
	case 10028:
		uiWindowEnable(sessionFilterRoot.ChildByID(10031), int(^button.DrawData().Field0>>2&1))
	}
	return nil
}
func sessionFilterClose() int {
	if sessionFilterRoot == nil {
		return 0
	}
	sessionFilterSave()
	sessionFilterRoot.Capture(false)
	sessionFilterRoot.Destroy()
	sessionFilterRoot = nil
	return 0
}
