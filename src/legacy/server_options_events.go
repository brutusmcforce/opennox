package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func serverOptionsEvent(_ *gui.Window, event int, arg uintptr, value int) int {
	if event > 16391 {
		if event != 16400 {
			return 1
		}
		child := (*gui.Window)(unsafe.Pointer(arg))
		switch child.ID() {
		case 10120:
			data := serverOptionsCurrent()
			selected := teamUIEvent(child, 16404, 0, 0)
			if selected < 0 || selected >= int(int16((*gui.ScrollListBoxData)(child.WidgetData).Field_11_0)) {
				return 1
			}
			text := serverOptionsGetText(child, 16406, value)
			serverOptionsSetText(serverOptionsChild(10119), 16385, text, -1)
			before := binary.LittleEndian.Uint16(data[52:]) & 0x17f0
			mode := (binary.LittleEndian.Uint16(data[52:]) & 0xe80f) | uint16(serverOptionsModeFromName(text))
			binary.LittleEndian.PutUint16(data[52:], mode)
			if before != mode&0x17f0 {
				serverOptionsMapList(int(mode), alloc.GoString(memmap.PtrUint8(0x5D4594, 1046556)), true)
			}
			child.SetHidden(true)
			child.Capture(false)
			binary.LittleEndian.PutUint16(data[54:], uint16(Nox_xxx_servGamedataGet_40A020(mode)))
			data[56] = byte(Sub_40A180(noxflags.GameFlag(mode)))
			serverOptionsLimits(data)
			serverOptionsDirty(1)
			if data[53]&0x10 != 0 {
				GetServer().S().Spells.EnableAll()
				sub_4537F0()
				serverOptionsResetMap()
			}
			serverOptionsQuest(int((binary.LittleEndian.Uint16(data[52:]) >> 12) & 1))
		case 10114:
			if teamUIEvent(child, 16404, 0, 0) < 0 {
				return 1
			}
			name := serverOptionsMapToken(serverOptionsGetText(child, 16406, value))
			data := serverOptionsCurrent()
			mapNameSetField(&data[0], name)
			ruleLoad((*server.Settings2)(unsafe.Pointer(&data[0])), "user.rul", nil, 7, binary.LittleEndian.Uint16(data[52:]))
			serverOptionsSettingsLabels(data)
			serverOptionsDirty(1)
		}
		return 1
	}
	if event == 16391 {
		child := (*gui.Window)(unsafe.Pointer(arg))
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		root := serverOptionsWindow(1046492)
		switch child.ID() {
		case 10119:
			list := serverOptionsChild(10120)
			list.SetHidden(false)
			list.Capture(true)
			list.ShowModal()
			serverOptionsPopulate()
		case 10122:
			if !noxflags.HasGame(49152) {
				on := bool2int(child.DrawData().Field0&4 != 0)
				uiWindowEnable(serverOptionsWindow(1046500), on)
				uiWindowEnable(serverOptionsWindow(1046504), on)
			}
			serverOptionsDirty(1)
		case 10141:
			serverOptionsResetMap()
		case 10145:
			data := serverOptionsRecord(unsafe.Pointer(serverConfigSlot(int32(1))))
			serverOptionsRead(data)
			serverConfigSlotCopy(int32(1), int32(0))
			count := GetServer().S().Teams.Count()
			mode := binary.LittleEndian.Uint16(data[52:])
			tooMany := false
			if noxflags.HasGame(128) && mode&0x60 != 0 {
				if int16(mode) >= 0 && count < 2 {
					Nox_xxx_dialogMsgBoxCreate_449A10(root, serverOptionsBufferText(1046560), serverOptionsText("NeedTeams"), gui.DialogFlags(56), func() { serverOptionsApply() }, nil)
					Sub_44A360(1)
					return 1
				}
				tooMany = count > 2
			} else if noxflags.HasGame(128) && noxflags.HasGame(16) && noxflags.HasGamePlay(4) && count > 2 {
				tooMany = true
			}
			if tooMany {
				Nox_xxx_dialogMsgBoxCreate_449A10(root, serverOptionsText("Notice"), serverOptionsText("TooManyTeams"), gui.DialogFlags(33), nil, nil)
				Sub_44A360(1)
				return 1
			}
			serverOptionsApply()
		case 10146:
			data := serverOptionsCurrent()
			if noxflags.HasGame(128) {
				if int8(data[53]) < 0 {
					GetServer().TeamsRemoveActive(true)
				}
				binary.LittleEndian.PutUint16(data[52:], binary.LittleEndian.Uint16(data[52:])&0x3fff)
			}
			serverOptionsClose(0)
		case 10149:
			serverOptionsClose(0)
		case 10152:
			serverPanelsAdvancedOpen(unsafe.Pointer(serverConfigSlot(int32(1))))
		case 10159:
			parent := child.Parent()
			child.SetParent(nil)
			child.SetParent(parent)
			if noxflags.HasGame(1) {
				serverOptionsChild(10196).SetHidden(false)
				serverOptionsHideRange(root, 10161, 10163, false)
				serverOptionsEnableRange(root, 10161, 10163, true)
				serverOptionsTab(0)
				serverOptionsChild(10163).DrawData().Field0 |= 4
				serverOptionsChild(10161).DrawData().Field0 &^= 4
				serverOptionsChild(10162).DrawData().Field0 &^= 4
			} else {
				serverOptionsChild(10196).SetHidden(true)
				serverOptionsPlayersPanel = uint32(teamUIPlayersConstruct(root))
			}
			serverOptionsChild(10141).SetHidden(true)
			serverOptionsWindow(1046524).SetHidden(true)
		case 10160:
			parent := child.Parent()
			child.SetParent(nil)
			child.SetParent(parent)
			if serverOptionsPlayersPanel != 0 {
				teamUIPlayersDestroy(true)
				serverOptionsPlayersPanel = 0
			}
			if serverOptionsAccessPanel != 0 {
				serverPanelsAccessClose(true)
				serverOptionsAccessPanel = 0
			}
			if serverOptionsAdvancedControl != 0 {
				serverPanelsGeneralClose()
				serverOptionsGeneralPanel = 0
			}
			serverOptionsHideRange(root, 10161, 10163, true)
			if noxflags.HasGame(1) {
				serverOptionsChild(10141).SetHidden(false)
			}
			serverOptionsChild(10196).SetHidden(true)
			serverOptionsWindow(1046524).SetHidden(false)
		case 10161:
			serverOptionsTab(1)
		case 10162:
			serverOptionsTab(2)
		case 10163:
			serverOptionsTab(0)
		case 10330:
			if noxflags.HasGamePlay(4) {
				teamRuntimeToggle(false)
				GetServer().TeamsRemoveActive(true)
				serverOptionsEnableRange(serverOptionsWindow(1046508), 10331, 10333, false)
			} else {
				teamRuntimeCreateMap()
				if noxflags.HasGamePlay(2) {
					teamRuntimeBalance(false)
				} else {
					teamRuntimeToggle(true)
				}
				serverOptionsEnableRange(serverOptionsWindow(1046508), 10331, 10333, true)
			}
		case 10331:
			if noxflags.HasGamePlay(2) {
				noxflags.UnsetGamePlay(2)
				teamRuntimeToggle(true)
			} else {
				teamRuntimeToggle(false)
				teamRuntimeEnable()
			}
		case 10332:
			teamRuntimeBalance(true)
		case 10333:
			if noxflags.HasGamePlay(1) {
				noxflags.UnsetGamePlay(1)
			} else {
				noxflags.SetGamePlay(1)
			}
		}
		return 1
	}
	if event == 23 {
		return 0
	}
	if event != 16387 {
		return 1
	}
	child := serverOptionsChild(uint(value))
	data := serverOptionsCurrent()
	if child == nil {
		return 0
	}
	if uint16(arg) == 1 {
		if value == 10101 {
			child.DrawData().TextColorVal = uint32(nox_color_white_2523948)
		}
		return 1
	}
	text := serverOptionsGetText(child, 16413, 0)
	if text == "" {
		return 1
	}
	n := int(resourceAtoi(text))
	if n < 0 {
		n = 0
	}
	switch value {
	case 10101:
		child.DrawData().TextColorVal = uiMeterColor(230, 165, 65)
		alloc.StrCopyZero(data[9:24], text)
		Nox_xxx_gameSetServername_40A440(text)
	case 10134:
		binary.LittleEndian.PutUint16(data[54:], uint16(n))
		serverOptionsDirty(1)
	case 10135:
		if n > 255 {
			n = 255
			serverOptionsSetText(child, 16414, "255", -1)
		}
		data[56] = byte(n)
		serverOptionsDirty(1)
	}
	return 1
}
