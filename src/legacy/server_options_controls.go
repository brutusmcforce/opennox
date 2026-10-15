package legacy

import (
	"encoding/binary"
	"strconv"
	"strings"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func serverOptionsLimits(data []byte) int {
	serverOptionsSetText(serverOptionsWindow(1046516), 16414, strconv.Itoa(int(binary.LittleEndian.Uint16(data[54:]))), 0)
	return serverOptionsSetText(serverOptionsWindow(1046520), 16414, strconv.Itoa(int(data[56])), 0)
}
func serverOptionsQuest(value int) int {
	enabled := bool2int(value != 1)
	for _, id := range []uint{10152, 10141} {
		uiWindowEnable(serverOptionsWindow(1046504).ChildByID(id), enabled)
	}
	uiWindowEnable(serverOptionsWindow(1046516), enabled)
	uiWindowEnable(serverOptionsWindow(1046520), enabled)
	if value == 1 {
		serverOptionsChild(10122).SetHidden(true)
	}
	return serverOptionsHidden(serverOptionsChild(10183), value)
}
func serverOptionsLabels(data []byte) int {
	serverOptionsStoreText(1045968, 64, serverOptionsFormat("SettingsMsg", mapNameField(&data[0])))
	modeText := serverOptionsText("GameType")
	if !noxflags.HasGame(128) {
		modeText = serverOptionsFormat("GameTypeIs", serverOptionsModeName(uint16(noxflags.GetGame())))
	}
	serverOptionsStoreText(1046096, 64, modeText)
	key := "OptsMessage"
	if noxflags.HasGame(1) {
		key = "GoMessage"
	}
	serverOptionsStoreText(1046224, 64, serverOptionsFormat(key))
	teamUIEvent(serverOptionsChild(10121), 16385, uintptr(memmap.PtrOff(0x5D4594, 1045968)), ^uintptr(0))
	teamUIEvent(serverOptionsChild(10118), 16385, uintptr(memmap.PtrOff(0x5D4594, 1046096)), ^uintptr(0))
	return teamUIEvent(serverOptionsChild(10117), 16385, uintptr(memmap.PtrOff(0x5D4594, 1046224)), ^uintptr(0))
}
func serverOptionsSettingsLabels(data []byte) int {
	mode := binary.LittleEndian.Uint16(data[52:])
	label := "Servopts.wnd:KillLimit"
	special := false
	if mode&0x20 != 0 {
		label = "Servopts.wnd:CaptureLimit"
		special = true
	} else if mode&0x400 != 0 {
		label = "Servopts.wnd:DeathLimit"
		special = true
	}
	if data[57] == 0 && noxflags.HasGame(1) && (special || !noxflags.HasGame(49152)) {
		uiWindowEnable(serverOptionsWindow(1046516), 1)
		uiWindowEnable(serverOptionsWindow(1046520), 1)
	}
	serverOptionsWindow(1046516).DrawData().SetText(serverOptionsText(label))
	serverOptionsLabels(data)
	serverOptionsName(alloc.GoString(&data[9]))
	if noxflags.HasGame(1) && !noxflags.HasGame(49152) {
		for _, off := range []uintptr{1046500, 1046504} {
			uiWindowEnable(serverOptionsWindow(off), bool2int(data[57] == 0))
		}
	}
	check := serverOptionsChild(10122).DrawData()
	if data[57] != 0 {
		check.Field0 |= 4
	} else {
		check.Field0 &^= 4
	}
	serverPanelsSpellStore((*uint32)(unsafe.Pointer(&data[24])))
	serverPanelsWeaponStore((*uint32)(unsafe.Pointer(&data[44])))
	return int(serverPanelsArmorStore(binary.LittleEndian.Uint32(data[48:])))
}
func serverOptionsRead(data []byte) uintptr {
	name := serverOptionsGetText(serverOptionsWindow(1046512), 16413, 0)
	clear(data[9:24])
	copy(data[9:24], name)
	binary.LittleEndian.PutUint16(data[52:], uint16(serverOptionsSelectedMode()))
	copy(data[24:44], unsafe.Slice((*byte)(unsafe.Pointer(serverPanelsSpellPointer())), 20))
	copy(data[44:48], unsafe.Slice((*byte)(unsafe.Pointer(serverPanelsWeaponPointer())), 4))
	binary.LittleEndian.PutUint32(data[48:], uint32(serverPanelsArmorLoad()))
	for _, v := range []struct {
		off uintptr
		dst int
	}{{1046516, 54}, {1046520, 56}} {
		s := serverOptionsGetText(serverOptionsWindow(v.off), 16413, 0)
		if s == "" {
			continue
		}
		n := int(textDecimal((*uint16)(unsafe.Pointer(alloc.InternCString16(s)))))
		if v.dst == 54 {
			binary.LittleEndian.PutUint16(data[54:], uint16(n))
		} else {
			data[56] = byte(n)
		}
	}
	data[57] = byte((serverOptionsChild(10122).DrawData().Field0 >> 2) & 1)
	selected := teamUIEvent(serverOptionsWindow(1046496), 16404, 0, 0)
	if selected < 0 {
		data[0] = 0
		return uintptr(selected)
	}
	name = serverOptionsMapToken(serverOptionsGetText(serverOptionsWindow(1046496), 16406, selected))
	mapNameSetField(&data[0], name)
	return 0
}
func serverOptionsMapToken(text string) string {
	v := strings.TrimLeft(text, "\t")
	if i := strings.IndexByte(v, '\t'); i >= 0 {
		v = v[:i]
	}
	return v
}
func serverOptionsRefresh() int {
	data := serverOptionsRecord(unsafe.Pointer(serverConfigSlot(int32(1))))
	serverPanelsSpellStore((*uint32)(unsafe.Pointer(&data[24])))
	serverPanelsWeaponStore((*uint32)(unsafe.Pointer(&data[44])))
	serverPanelsArmorStore(binary.LittleEndian.Uint32(data[48:]))
	if serverOptionsRoot == 0 {
		return 0
	}
	serverOptionsSettingsLabels(data)
	serverPanelsAdvancedUpdate(unsafe.Pointer(&data[0]))
	serverOptionsMapList(int(binary.LittleEndian.Uint16(data[52:])), mapNameField(&data[0]), false)
	serverOptionsLimits(data)
	return serverOptionsSetText(serverOptionsChild(10119), 16385, serverOptionsModeName(binary.LittleEndian.Uint16(data[52:])), 0)
}
