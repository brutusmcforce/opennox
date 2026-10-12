package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/console"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func Nox_cmd_ban(i int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	name := t[i]
	p := teamUIPlayerByName(name)
	if p == nil {
		serverConfigBlockedAdd(0, alloc.InternCString16(name), nil)
		consoleCommandCentered("banDisallow", name)
		return true
	}
	if p.PlayerInd == 31 {
		consoleCommandCentered("cantbanyourself")
		return true
	}
	if noxflags.HasGame(4096) {
		Sub_4DCFB0(p.PlayerUnit)
	} else {
		Nox_xxx_playerDisconnByPlrID_4DEB00(ntype.PlayerInd(p.PlayerInd))
	}
	serverConfigBlockedAdd(0, alloc.InternCString16(name), (*byte)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(2112)) /* SerialBuf */)))
	consoleCommandCentered("banned", p.Name())
	return true
}
func Nox_cmd_kick(i int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	p := teamUIPlayerByName(t[i])
	if p == nil {
		return true
	}
	if p.PlayerInd == 31 {
		consoleCommandCentered("cantkickyourself")
		return true
	}
	if noxflags.HasGame(4096) {
		Sub_4DCFB0(p.PlayerUnit)
	} else {
		Nox_xxx_playerCallDisconnect_4DEAB0(ntype.PlayerInd(p.PlayerInd), 4)
		consoleCommandCentered("kicked", p.Name())
	}
	return true
}
func Nox_cmd_list_users(_ int, _ []string) bool {
	consoleCommandPrint("userslist")
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		name := p.Name()
		if consoleCommandServer && p.Field3680&4 != 0 {
			name += ", " + consoleCommandText("SysMuted")
		}
		if p.Field3680&8 != 0 {
			name += ", " + consoleCommandText("ClientMuted")
		}
		GetConsole().Print(console.ColorRed, consoleCommandFormat(alloc.GoString16(memmap.PtrUint16(0x587000, 106604)), name))
	}
	return true
}
func consoleCommandMute(i int, t []string, mute bool) bool {
	if len(t) < 2 || len(t) > 3 {
		return false
	}
	mask := 8
	if consoleCommandServer && mapASCIIEqual(t[i], consoleCommandToken(7)) {
		if i+1 != len(t)-1 {
			return false
		}
		i++
		mask = 4
	}
	p := teamUIPlayerByName(t[i])
	found := p != nil && (mask == 4 || noxflags.HasGame(2)) && !(mute && mask == 8 && p.PlayerInd == 31)
	id := "UserNotFound"
	if found {
		if mute {
			Nox_xxx_netNeedTimestampStatus_4174F0(p, mask)
			id = "Muted"
		} else {
			Nox_xxx_playerUnsetStatus_417530(p, mask)
			id = "UnMuted"
		}
	}
	consoleCommandPrint(id, t[i])
	return true
}
func Nox_cmd_mute(i int, t []string) bool   { return consoleCommandMute(i, t, true) }
func Nox_cmd_unmute(i int, t []string) bool { return consoleCommandMute(i, t, false) }
func Nox_cmd_cheat_ability(_ int, _ []string) bool {
	if !noxflags.HasGame(8192) {
		players := &GetServer().S().Players
		for p := players.First(); p != nil; p = players.Next(p) {
			if p.PlayerUnit != nil {
				Nox_xxx_playerCancelAbils_4FC180(p.PlayerUnit)
			}
		}
	}
	return true
}
func Nox_cmd_cheat_level(_ int, t []string) bool {
	if noxflags.HasGame(8192) {
		return true
	}
	if len(t) < 3 {
		return false
	}
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		if p.PlayerUnit != nil {
			controlSetLevel(p.PlayerUnit, byte(consoleCommandNumber(t[2])))
		}
	}
	return true
}
