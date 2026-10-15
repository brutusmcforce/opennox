package legacy

import (
	"context"
	"strings"
	"unicode/utf16"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// AFKObserver lets idle players go to observer mode automatically (the clients ask for it after 90 s without input).
var AFKObserver = true

func consoleCommandRemote(pl *server.Player, action byte, text string) int {
	if pl == nil || pl.PlayerUnit == nil {
		return 0
	}
	text = consoleCommandString(text)
	consoleCommandSender = pl
	defer func() { consoleCommandSender = nil }()
	if action != 0 && action != 4 && action != 5 && noxflags.HasGame(49152) {
		return 1
	}
	switch action {
	case 0:
		if noxflags.HasGame(8|4096) || pl.Field3680&1 != 0 {
			return 1
		}
		force := text == "\uf00d\uf0ad"
		if force && !AFKObserver {
			return 1 // the automatic request of an idle client (opennox#691)
		}
		if Nox_xxx_playerGoObserver_4E6860(pl, bool2int(force), 0) != 1 {
			return 1
		}
		if Nox_xxx_gamePlayIsAnyPlayers_40A8A0() != 0 {
			if !force {
				Nox_xxx_netNeedTimestampStatus_4174F0(pl, 256)
			}
			if noxflags.HasGame(1024) && Sub_40A770() == 1 {
				matchRosterWinner(true)
			}
		}
		if pl.PlayerUnit != nil {
			teamRuntimeLeave(pl.PlayerUnit.TeamPtr(), pl.NetCode())
		}
		consoleCommandPrint("observermode", consoleCommandText("set"))
	case 1, 2:
		// The original check compares the player slot with the local network code.
		if Sub_4D12A0(int(pl.PlayerInd)) == 0 && int(pl.PlayerInd) != ClientPlayerNetCode() && !noxflags.HasGame(2048) {
			return 1
		}
		if action == 2 {
			consoleCommandPrint("RemoteSysop", pl.Name(), text)
			if text != "" {
				ExecConsoleCmd(context.Background(), text)
			}
			return 1
		}
		fields := strings.FieldsFunc(text, func(r rune) bool { return r == ' ' })
		name := "(null)"
		if len(fields) > 1 {
			name = fields[1]
		}
		index := GetServer().S().NoxScriptVM.ScriptIndexByName(consoleCommandNarrow(name))
		if index != -1 && consoleCommandSender != nil {
			consoleCommandPrint("ExecutingFunction", name)
			p := consoleCommandSender
			if err := GetServer().S().NoxScriptVM.CallByIndex(index, p.PlayerUnit, p.PlayerUnit); err != nil {
				scriptLog.Println(err)
			}
		} else {
			consoleCommandPrint("InvalidFunction", name)
		}
	case 3, 5:
		flag := byte(0)
		if action == 5 {
			flag = 16
		}
		gameplayTextAll(flag, utf16.Encode([]rune(consoleCommandFormat(text))))
		if action == 5 {
			players := &GetServer().S().Players
			for p := players.First(); p != nil; p = players.Next(p) {
				if p.PlayerUnit != nil {
					GetServer().S().Audio.EventObj(902, p.PlayerUnit, 0, 0)
				}
			}
		}
	case 4:
		if pl.Field3680&1 == 0 && !noxflags.HasEngine(noxflags.EngineReplayRead) {
			if noxflags.HasGame(1) {
				consoleCommandPrint("notinobserver")
			}
			return 1
		}
		if text == "" {
			Nox_xxx_playerCameraUnlock_4E6040(pl.PlayerUnit)
			return 1
		}
		players := &GetServer().S().Players
		for p := players.First(); p != nil; p = players.Next(p) {
			if mapASCIIEqual(text, p.Name()) {
				Nox_xxx_playerCameraFollow_4E6060(pl.PlayerUnit, p.PlayerUnit)
			}
		}
	default:
		consoleCommandPrint("invalidattempt", pl.Name(), text)
	}
	return 1
}
