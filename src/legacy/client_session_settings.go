package legacy

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func clientSessionSettings(op int, data []byte) {
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	switch op {
	case 175:
		dword_5d4594_1200768 = 0
		generation := dword(1)
		if uint32(serverConfig3512Get()) >= generation {
			return
		}
		serverConfig3512Set(int32(generation))
		ctf, quest := noxflags.HasGame(32), noxflags.HasGame(4096)
		if noxflags.HasGame(1) {
			return
		}
		noxflags.UnsetGame(524272)
		noxflags.SetGame(noxflags.GameFlag(dword(9)))
		serverConfigFlagsSet(int32(dword(13)))
		Nox_client_setVersion_409AE0(dword(5))
		serverConfigLimitSet(int32(data[17]))
		slot := unsafe.Slice(serverConfigSlot(0), 58)
		if slot[56] != data[19] || binary.LittleEndian.Uint16(slot[54:]) != uint16(data[18]) {
			dword_5d4594_1200768 = 1
		}
		binary.LittleEndian.PutUint16(slot[54:], uint16(data[18]))
		slot[56] = data[19]
		if int8(data[9]) >= 0 {
			binary.LittleEndian.PutUint16(slot[52:], uint16(dword(9)))
			alloc.StrCopy(slot, alloc.GoString(sessionMapName()))
		}
		if !ctf && noxflags.HasGame(32) {
			teamUICTFOpen(2)
		} else if !noxflags.HasGame(32) {
			teamUIHUDHide(false)
		}
		// C deliberately tests the previous CTF bit for the ball transition too.
		if !ctf && noxflags.HasGame(64) {
			teamUIBallOpen()
		} else if !noxflags.HasGame(64) {
			teamUIHUDHide(true)
		}
		if quest && !noxflags.HasGame(4096) {
			for p := GetServer().S().Players.First(); p != nil; p = GetServer().S().Players.Next(p) {
				playerStateRespawn(p, 255)
			}
		}
		if !noxflags.HasGame(0x20000) && GameGetPlayState() == 3 {
			Nox_game_exit_xxx2()
		}
	case 176:
		if noxflags.HasGame(1) {
			return
		}
		slot := unsafe.Slice(serverConfigSlot(0), 58)
		// Match strncpy's padding and lack of an extra terminator at capacity.
		name := data[1:16]
		zero := bytes.IndexByte(name, 0)
		copy(slot[9:24], name)
		if zero >= 0 {
			clear(slot[9+zero : 24])
		}
		if !bytes.Equal(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1200732), 20), data[17:37]) || memmap.Uint32(0x5D4594, 1200752) != dword(37) || memmap.Uint32(0x5D4594, 1200756) != dword(41) {
			dword_5d4594_1200768 = 1
		}
		copy(slot[24:44], data[17:37])
		serverPanelsSpellApply((*uint32)(unsafe.Pointer(&slot[24])))
		copy(slot[44:52], data[37:45])
		if dword(45) != 0 {
			serverConfigTimerSet(1)
			serverConfigTimerReset(int32(dword(45)))
		} else {
			serverConfigTimerSet(0)
		}
		copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1200708), 58), slot)
		if Nox_client_isConnected() && dword_5d4594_1200768 != 0 {
			Nox_xxx_printCentered_445490(clientGameProgressString("OptionsChanged"))
			audioEventPlay(310, 100, 0, 0)
		}
	case 177:
		slot := unsafe.Slice(serverConfigSlot(int32(data[1])), 58)
		old := slot[52]
		copy(slot, data[2:60])
		if (old>>5)&1 != (data[54]>>5)&1 {
			teamUIPlayersRefreshOpen()
		}
		if serverConfigAcquiredGet() == 0 {
			if serverOptionsRoot == 0 {
				if data[1] == 1 {
					serverConfigSlotCopy(1, 0)
				}
				if noxflags.HasGame(128) {
					if dword_5d4594_1200832 != 0 {
						p := Get_dword_8531A0_2576()
						Nox_xxx_printCentered_445490(clientGameProgressText("NameChange", alloc.GoString16(&p.NameFinal[0])))
						dword_5d4594_1200832 = 0
					}
					if nox_server_sanctuaryHelp_54276 != 0 {
						interactionHelpOpen()
					}
				}
			}
			serverConfigAcquiredSet(1)
		}
		serverOptionsRefresh()
	}
}
