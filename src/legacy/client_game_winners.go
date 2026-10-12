package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func clientGameScoreLimit() uint32 {
	if noxflags.HasGame(noxflags.GameHost) {
		return uint32(uint16(serverConfigScore(int16(noxflags.GetGame()))))
	}
	return uint32(memmap.Uint16(0x5D4594, 371434))
}
func clientGamePlayerScore(winner *server.Player) {
	s := GetServer().S()
	limit := clientGameScoreLimit()
	elimination := noxflags.HasGame(1024)
	for p := s.Players.First(); p != nil; p = s.Players.Next(p) {
		if p.Field3680&1 != 0 {
			continue
		}
		if elimination {
			if p == winner {
				if p.Field2140 >= limit {
					p.Field2140 = limit - 1
				}
			} else if p.Field2140 < limit {
				p.Field2140 = limit
			}
		} else {
			if p == winner {
				p.Lessons = int32(limit)
			} else if uint32(p.Lessons) >= limit {
				p.Lessons = int32(limit - 1)
			}
		}
	}
}
func clientGameTeamScore(winner *server.Team) {
	s := GetServer().S()
	limit := clientGameScoreLimit()
	elimination := noxflags.HasGame(1024)
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		if !elimination {
			if tm == winner {
				tm.Lessons = int(limit)
			} else if uint32(tm.Lessons) >= limit {
				tm.Lessons = int(limit - 1)
			}
		}
	}
	update := func(member *server.ObjectTeam, code uint32) {
		if teamRuntimeContains(member, winner.ID()) {
			return
		}
		p := s.Players.ByID(int(code))
		if p == nil || p.Field3680&1 != 0 {
			return
		}
		if elimination {
			if p.Field2140 < limit {
				p.Field2140 = limit
			}
		} else if uint32(p.Lessons) >= limit {
			p.Lessons = int32(limit - 1)
		}
	}
	if noxflags.HasGame(noxflags.GameHost) {
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			update(u.TeamPtr(), u.NetCode)
		}
	} else {
		for dr := GetClient().Cli().Objs.FirstPlayerList(); dr != nil; dr = dr.Field_104 {
			update(dr.TeamPtr(), dr.NetCode32)
		}
	}
}

func clientGameWinner(kind byte, data []byte) int {
	s := GetServer().S()
	code := binary.LittleEndian.Uint16(data[1:])
	var team *server.Team
	if kind != 88 {
		team = s.Teams.ByID(server.TeamID(code))
	}
	if !Nox_client_isConnected() || binary.LittleEndian.Uint32(data[4:]) <= clientGameMapFrame {
		return 8
	}
	audioEventPlay(309, 100, 0, 0)
	if !noxflags.HasGame(noxflags.GameHost) {
		noxflags.SetGame(8)
	}
	if kind == 88 {
		if !noxflags.HasEngine(noxflags.EngineNoRendering) {
			Sub_470510()
		}
		interactionInitTime()
	} else {
		interactionInitTime()
		if !noxflags.HasEngine(noxflags.EngineNoRendering) {
			Sub_470510()
		}
	}
	text := clientGameProgressText
	message := ""
	mode := uint32(0)
	if kind != 86 && data[3] == 1 {
		message = text("TimeLimitReached")
	}
	elimination := noxflags.HasGame(1024)
	localCode := ClientPlayerNetCode()
	if kind == 88 {
		id := int(code & 0x7fff)
		p := s.Players.ByID(id)
		if !elimination {
			if code == 0 {
				message += text("DM_Tie")
			} else if p == nil {
				mode = 1
			} else {
				if data[3] == 0 {
					clientGamePlayerScore(p)
				}
				if id != localCode {
					message += text("DM_Loss", p.Name())
					mode = 1
				} else if *(*byte)(unsafe.Add(p.C(), 2252)) == 0 {
					message += text("DM_MaleVictory")
				} else {
					message += text("DM_FemaleVictory")
				}
			}
		} else if code == 0 {
			message += text("HL_Tie")
		} else {
			message += text("HL_Header")
			if id == localCode {
				message += text("HL_Victory", clientGameProgressString("HL_You"))
			} else {
				if p != nil {
					message += text("HL_Victory", p.Name())
					if data[3] == 0 {
						clientGamePlayerScore(p)
					}
				}
				mode = 1
			}
		}
	} else {
		localTeam := teamRuntimeObject(localCode)
		name := "(null)"
		if team != nil {
			name = team.Name()
		}
		switch kind {
		case 86:
			if teamRuntimeContains(localTeam, server.TeamID(byte(code))) {
				message = text("TeamWon")
			} else {
				message = text("FB_Victory", text("teamformat", name))
				mode = 1
			}
		case 87:
			if team == nil {
				message += text("CTF_Tie")
			} else {
				message += text("CTF_Victory", name)
				if !teamRuntimeContains(localTeam, team.ID()) {
					mode = 1
				}
			}
		case 89:
			if team == nil {
				if elimination {
					message += text("HL_Tie")
				} else {
					message += text("DM_Tie")
				}
			} else {
				own := teamRuntimeContains(localTeam, team.ID())
				if elimination {
					message += text("HL_Header")
					if own {
						message += text("HL_Victory", clientGameProgressString("HL_YourTeam"))
					} else {
						// The original performs the second formatting pass after embedding the
						// localized team template in the victory format.
						message += consoleCommandFormat(text("HL_Victory", clientGameProgressString("Team")), name)
						mode = 1
					}
				} else if own {
					message += text("DM_TeamVictory")
				} else {
					message += text("DM_Loss", text("Team", name))
					mode = 1
				}
				if data[3] == 0 {
					clientGameTeamScore(team)
				}
			}
		}
	}
	interactionTextState(alloc.InternCString16(message), mode)
	serverOptionsClose(0)
	return 8
}
