package legacy

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

// clientGameNotice preserves the inform.c string namespace and message lengths.
func clientGameNotice(data []byte) int {
	text := func(key string) string { return GetServer().S().Strings().GetStringInFile(strman.ID(key), "inform.c") }
	print := func(key string, args ...any) { Nox_xxx_printCentered_445490(consoleCommandFormat(text(key), args...)) }
	word := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	player := func() *server.Player { return GetServer().S().Players.ByID(int(word(2))) }
	title := func(off int) string { return GetServer().S().Teams.TeamTitle(server.TeamColor(word(off))) }
	switch data[1] {
	case 0, 2:
		bookAwardClientMessage(int32(word(2)), data[1] == 2)
		return 6
	case 1:
		name, ok := Nox_xxx_spellTitle_424930(int(word(2)))
		if !ok {
			name = "(null)"
		}
		print("plyrspel.c:SpellCastSuccess", name)
		return 6
	case 3, 4, 18, 19, 20:
		if p := player(); p != nil {
			key := map[byte]string{3: "netserv.c:PlayerTimeout", 4: "objcoll.c:FlagRetrieveNotice", 18: "objcoll.c:PlayerExited", 19: "objcoll.c:PlayerExitedWarp", 20: "GeneralPrint:SecretFoundOther"}[data[1]]
			print(key, p.Name())
			if data[1] == 4 {
				audioEventPlay(305, 100, 0, 0)
			}
		}
		return 6
	case 5, 6, 7:
		if p := player(); p != nil {
			key := map[byte]string{5: "objcoll.c:FlagCaptureNotice", 6: "objcoll.c:FlagPickupNotice", 7: "drop.c:FlagDropNotice"}[data[1]]
			print(key, p.Name(), title(6))
			sound := map[byte]int32{5: 306, 6: 303, 7: 304}[data[1]]
			audioEventPlay(sound, 100, 0, 0)
		}
		return 10
	case 8:
		print("update.c:FlagRespawnNotice", title(2))
		audioEventPlay(305, 100, 0, 0)
		return 6
	case 9:
		team := GetServer().S().Teams.ByID(server.TeamID(word(6)))
		p := player()
		if team != nil {
			if p == nil {
				print("objcoll.c:FlagBallUnknownNotice", team.Name())
			} else {
				print("objcoll.c:FlagBallNotice", p.Name(), team.Name())
			}
		}
		return 10
	case 10, 11:
		p := player()
		team := GetServer().S().Teams.ByID(server.TeamID(word(6)))
		if p != nil {
			if data[1] == 10 {
				if team == nil {
					print("pickup.c:PickUpCrown", p.Name())
				} else {
					print("pickup.c:PickUpTeamCrown", p.Name(), team.Name())
				}
			} else {
				if team == nil {
					print("drop.c:DropCrown", p.Name())
				} else {
					print("drop.c:DropTeamCrown", p.Name(), team.Name())
				}
			}
		}
		return 10
	case 12:
		Nox_xxx_printCentered_445490(text("Netserv.c:InObservationMode"))
		if word(2) != 0 {
			Nox_xxx_printCentered_445490(text("Netserv.c:PressJump"))
		}
		return 6
	case 13:
		keys := [...]string{"atckexec.c:PlayerStunned", "atckexec.c:PlayerConfused", "atckexec.c:PlayerPoisoned", "player.c:TooHeavy"}
		if id := word(2); id < uint32(len(keys)) {
			Nox_xxx_printCentered_445490(text(keys[id]))
		}
		return 6
	case 14:
		combatFeedAdd(unsafe.Pointer(&data[0]))
		return 11
	case 15:
		n := bytes.IndexByte(data[3:], 0)
		if n < 0 {
			return 0
		}
		if GameGetPlayState() == 3 {
			value := text(string(data[3 : 3+n]))
			if data[2] != 0 {
				print("use.c:SignSays", value)
			} else {
				Nox_xxx_printCentered_445490(value)
			}
		}
		return n + 4
	case 16:
		print("pickup.c:WrongTeam", title(2))
		return 6
	case 17:
		Nox_xxx_dialogMsgBoxCreate_449A10(nil, text("guiserv.c:Notice"), text("Noxworld.c:ErrChangedClass"), 33, nil, nil)
		return 2
	case 21:
		stage := consoleCommandFormat(text("GeneralPrint:AdvanceToStage1"), int32(word(2)))
		print("use.c:SignSays", stage)
		return 6
	default:
		return 0
	}
}
