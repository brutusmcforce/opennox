package legacy

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func serverOptionsApply() int8 {
	data := serverOptionsCurrent()
	if data[52]&0x60 != 0 {
		s := GetServer().S()
		for i := s.Teams.Count(); i < 2; i++ {
			t := s.Teams.Create(0)
			color := 1
			if i != 0 {
				color = 2
			}
			if noxflags.HasGame(128) && teamUISettingsLocked() {
				teamRuntimeSetName(t, memmap.PtrUint16(0x5D4594, 1046564), 0)
			} else {
				teamRuntimeSetName(t, alloc.InternCString16(s.Teams.TeamTitle(server.TeamColor(color))), 1)
			}
			t.ColorInd = server.TeamColor(color)
			teamRuntimeAnnounce(t)
		}
		noxflags.SetGamePlay(4)
	}
	selected := alloc.GoString(&data[0])
	cycle := data[57] != 0
	result := int8(data[57])
	if cycle {
		p := sessionSelectedMap()
		selected = alloc.GoString(p)
		result = int8(memmap.PtrWord(unsafe.Pointer(p)))
	}
	if selected == "" {
		return result
	}
	different := !mapASCIIEqual(selected, GoStringP(unsafe.Pointer(sessionMapName())))
	serverConfigSlotCopy(int32(1), int32(0))
	current := serverOptionsRecord(unsafe.Pointer((*int8)(unsafe.Pointer(serverConfigSlotSelect(int32(0))))))
	Nox_xxx_gameSetServername_40A440(alloc.GoString(&current[9]))
	mode := binary.LittleEndian.Uint16(current[52:])
	if mode&0x1000 == 0 {
		Sub_409FB0_settings(mode, binary.LittleEndian.Uint16(current[54:]))
		sub_40A040_settings(int16(mode), uint8(current[56]))
	}
	if noxflags.HasGame(128) {
		noxflags.UnsetGame(49152)
		noxflags.SetGame(noxflags.GameFlag(mode & 0xc000))
	}
	current[57] = 0
	// Selection may alias record zero, overwritten by the settings copy above.
	// Keep the comparison from before that copy, but read the live requested name.
	if !cycle {
		selected = alloc.GoString(&data[0])
	}
	if different {
		suffix := alloc.GoString(memmap.PtrUint8(0x587000, 131668))
		nox_xxx_mapLoad_4D2450((*int8)(internCStr(selected + suffix)))

		serverConfigSlotSelect(int32(1))
	} else {
		GetServer().S().Spells.EnableAll()
		sub_4537F0()
		head := (*legacyListNode)(serverOptionsListHead())
		ruleWrite("user.rul", (*server.Settings2)(unsafe.Pointer((*int8)(unsafe.Pointer(serverConfigSlot(int32(1)))))), head)
		commandRulesMap(alloc.GoString(sessionMapFilename()))
		ruleLoad((*server.Settings2)(unsafe.Pointer((*int8)(unsafe.Pointer(serverConfigSlot(int32(0)))))), "user.rul", head, 3, uint16(noxflags.GetGame()))
	}
	return int8(serverOptionsClose(0))
}
