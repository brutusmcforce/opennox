package legacy

import (
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func Nox_cmd_set_sysop(i int, t []string) bool {
	if len(t) != 3 {
		return false
	}
	Nox_xxx_sysopSetPass_40A610(t[i])
	consoleCommandPrint("sysoppasswordset")
	return true
}
func Nox_cmd_set_cycle(i int, t []string) bool {
	if len(t) != 3 {
		return false
	}
	switch {
	case mapASCIIEqual(t[i], "on"):
		mapCycleSetEnabled(1)
		consoleCommandPrint("MapCycleOn")
	case mapASCIIEqual(t[i], "off"):
		mapCycleSetEnabled(0)
		consoleCommandPrint("MapCycleOff")
	default:
		return false
	}
	serverPanelsGeneralRefresh()
	return true
}
func consoleCommandFlag(text, id string, mask int32) bool {
	token := ""
	switch {
	case mapASCIIEqual(text, "on"):
		serverConfigFlagsAdd(mask)
		token = "cmd_token:on"
	case mapASCIIEqual(text, "off"):
		serverConfigFlagsRemove(mask)
		token = "cmd_token:off"
	default:
		return false
	}
	serverConfigUpdatedSet()
	consoleCommandPrint(id, consoleCommandText(token))
	return true
}
func Nox_cmd_set_weapons(i int, t []string) bool {
	return len(t) == 3 && consoleCommandFlag(t[i], "weapons", 1)
}
func Nox_cmd_set_staffs(i int, t []string) bool {
	return len(t) == 3 && consoleCommandFlag(t[i], "staffs", 16)
}
func Nox_cmd_set_mnstrs(i int, t []string) bool {
	if len(t) != 3 && len(t) != 4 {
		return false
	}
	if mapASCIIEqual(t[2], consoleCommandToken(6)) {
		return len(t) == 4 && consoleCommandFlag(t[3], "monsterrespawn", 8)
	}
	return consoleCommandFlag(t[i], "monsters", 4)
}
func Nox_cmd_set_name(i int, t []string) bool {
	if len(t) < 3 {
		return false
	}
	parts := make([]string, 0, len(t)-i)
	for _, s := range t[i:] {
		parts = append(parts, consoleCommandNarrow(s))
	}
	sep := alloc.GoString(memmap.PtrUint8(0x587000, 104484))
	name := strings.Join(parts, sep)
	if name != "" {
		Nox_xxx_gameSetServername_40A440(name)
		consoleCommandPrint("setgamename", consoleCommandWiden(name))
	}
	return true
}
func Nox_cmd_set_spell(_ int, t []string) bool {
	if len(t) != 4 {
		return false
	}
	if noxflags.HasGame(128) {
		consoleCommandPrint("NotInChat", t[2])
		return true
	}
	s := GetServer().S()
	id := s.Spells.ByTitle(t[2])
	if id == 0 {
		id = spell.ParseID(consoleCommandNarrow(t[2]))
		if id < 0 {
			id = 0
		}
	}
	if id == 0 {
		consoleCommandPrint("invalidspell", t[2])
		return false
	}
	switch {
	case mapASCIIEqual(t[3], "on"):
		mode := *(*uint16)(unsafe.Add(unsafe.Pointer(serverConfigSlotCurrent()), 52))
		if (noxflags.HasGame(64) || mode&64 != 0) && id == 132 {
			return true
		}
		if s.Spells.DefByInd(id).IsEnabled() {
			return true
		}
		s.Spells.Enable(id, true)
		serverConfigUpdatedSet()
		consoleCommandPrint("spellenabled", t[2])
		return true
	case mapASCIIEqual(t[3], "off"):
		if !s.Spells.DefByInd(id).IsEnabled() {
			return true
		}
		s.Spells.Enable(id, false)
		serverConfigUpdatedSet()
		consoleCommandPrint("spelldisabled", t[2])
		return true
	}
	return false
}
func Nox_cmd_set_players(_ int, t []string) bool {
	if len(t) != 3 {
		return false
	}
	n := consoleCommandNumber(t[2])
	if n < 0 {
		n = 0
	}
	if n > 999 {
		n = 999
	}
	if serverConfigLimitGet() == n {
		return true
	}
	serverConfigLimitSet(n)
	serverPanelsAccessRefresh()
	consoleCommandCentered("playersset", n)
	return true
}
func Nox_cmd_set_time(_ int, t []string) bool {
	if len(t) != 3 {
		return false
	}
	mode := *(*int16)(unsafe.Add(unsafe.Pointer(serverConfigSlotCurrent()), 52))
	GetServer().Sub40A040settings(int(mode), int(byte(consoleCommandNumber(t[2]))))
	return true
}
func Nox_cmd_set_lessons(_ int, t []string) bool {
	if len(t) != 3 {
		return false
	}
	mode := *(*uint16)(unsafe.Add(unsafe.Pointer(serverConfigSlotCurrent()), 52))
	Sub_409FB0_settings(mode, uint16(consoleCommandNumber(t[2])))
	return true
}
func serverConfigScoreSet(mode int16, n uint16) {
	p := memmap.PtrUint16(0x5D4594, 3488+2*uintptr(Sub_409A70(int(mode))))
	if *p == n {
		return
	}
	if n > 999 {
		n = 999
	}
	*p = n
	serverConfigUpdatedSet()
	if Nox_client_isConnected() {
		consoleCommandCentered("parsecmd.c:lessonsset", n)
	}
}
func consoleCommandQuality(kind, online int) bool {
	Set_dword_5d4594_2650652(online)
	serverConfigRateSet(serverConfigConnectionRate(int32(kind)))
	Set_nox_server_connectionType_3596(kind)
	serverPanelsGeneralRefresh()
	return true
}
func Nox_cmd_set_qual_modem(_ int, _ []string) bool { return consoleCommandQuality(4, 1) }
func Nox_cmd_set_qual_isdn(_ int, _ []string) bool  { return consoleCommandQuality(3, 1) }
func Nox_cmd_set_qual_cable(_ int, _ []string) bool { return consoleCommandQuality(2, 1) }
func Nox_cmd_set_qual_t1(_ int, _ []string) bool    { return consoleCommandQuality(1, 1) }
func Nox_cmd_set_qual_lan(_ int, _ []string) bool   { return consoleCommandQuality(1, 0) }
