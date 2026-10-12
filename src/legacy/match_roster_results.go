package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func matchRosterWinner(minimum bool) int {
	s := GetServer().S()
	var team *server.Team
	var unit *server.Object
	best := int32(math.MinInt32)
	if minimum {
		best = math.MaxInt32
	}
	tied := false
	accepts := func(v int32) bool {
		if minimum {
			return v <= best
		}
		return v >= best
	}
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		score := int32(tm.Lessons)
		if accepts(score) {
			tied = score == best && team != nil
			best = score
			team = tm
		}
	}
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		pl := u.UpdateDataPlayer().Player
		if u.TeamVal.Has() || pl.Field3680&1 != 0 {
			continue
		}
		score := int32(pl.Lessons)
		if minimum {
			score = int32(pl.Field2140)
		}
		if accepts(score) {
			// The C selector tracks ties separately for the team and solo-player passes.
			tied = score == best && unit != nil
			best = score
			unit = u
		}
	}
	noxflags.SetGame(8)
	if tied {
		return gameplayReportDMTeamWinner(nil, 1)
	}
	if unit != nil {
		return gameplayReportDMWinner(unit, 1)
	}
	if team != nil {
		return gameplayReportDMTeamWinner(team, 1)
	}
	return 0
}
func matchRosterFlagWinner() int {
	s := GetServer().S()
	var team *server.Team
	best := int32(-1)
	tied := false
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		score := int32(tm.Lessons)
		if score >= best {
			tied = score == best && team != nil
			best = score
			team = tm
		}
	}
	noxflags.SetGame(8)
	if team == nil || tied {
		return gameplayReportFlagWinner(nil, 1)
	}
	if noxflags.HasGame(64) {
		return gameplayReportFlagballWinner(team)
	}
	return gameplayReportFlagWinner(team, 1)
}
func matchRosterCheckVictory() {
	s := GetServer().S()
	limit := uint32(uint16(serverConfigScore(int16(noxflags.GetGame()))))
	if noxflags.HasGame(1024) {
		if limit == 0 {
			return
		}
		var team *server.Team
		var unit *server.Object
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			pl := u.UpdateDataPlayer().Player
			if pl.Field3680&1 != 0 || pl.Field2140 >= limit {
				continue
			}
			if u.TeamVal.Has() {
				tm := s.Teams.ByID(u.TeamVal.ID)
				if team != nil {
					if team != tm {
						return
					}
				} else {
					team = tm
				}
			} else {
				if unit != nil || team != nil {
					return
				}
				unit = u
			}
		}
		if playerStateMultiple() == 0 {
			return
		}
		noxflags.SetGame(8)
		if team != nil {
			gameplayReportDMTeamWinner(team, 0)
		} else {
			gameplayReportDMWinner(unit, 0)
		}
		return
	}
	if noxflags.HasGame(512) || limit == 0 {
		return
	}
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		if int32(tm.Lessons) >= int32(limit) {
			noxflags.SetGame(8)
			gameplayReportDMTeamWinner(tm, 0)
			return
		}
	}
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		pl := u.UpdateDataPlayer().Player
		if pl.Field3680&1 == 0 && pl.Lessons >= int32(limit) {
			noxflags.SetGame(8)
			gameplayReportDMWinner(u, 0)
			return
		}
	}
}

// Preserve the low 16 bits of the C float -> signed long long -> ushort path.
// At exponent 39 and above, all finite float32 integers have zero low 16 bits;
// the masked out-of-range/NaN conversion also has zero low 16 bits.
func matchRosterTimerCoordinate(value float32) uint16 {
	bits := math.Float32bits(value)
	exponent := int((bits>>23)&255) - 127
	if exponent < 0 || exponent >= 39 {
		return 0
	}
	magnitude := uint32(0x800000) | (bits & 0x7fffff)
	if exponent >= 23 {
		magnitude <<= uint(exponent - 23)
	} else {
		magnitude >>= uint(23 - exponent)
	}
	out := uint16(magnitude)
	if bits>>31 != 0 {
		return -out
	}
	return out
}
func matchRosterCheckLimit() int {
	if Sub_40A1A0() == 0 {
		return 0
	}
	s := GetServer().S()
	if noxflags.HasGame(4096) {
		GetServer().SwitchMap(alloc.GoString(worldQuestPending()))
		gameplayTextPrivateAll(alloc.InternCString("chklimit.c:AutoExitToNextMap"))
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			if *equipmentWord(u.UpdateData, 312) == 0 && *equipmentWord(u.UpdateDataPlayer().Player.C(), int(layoutServerPlayer.Off(4792)) /* Field4792 */) == 1 {
				questRuntimeStageComplete(u)
			}
		}
	} else if !GetServer().GetFlag3592() {
		if noxflags.HasGame(96) {
			matchRosterFlagWinner()
		} else if noxflags.HasGame(272) {
			matchRosterWinner(false)
		} else if noxflags.HasGame(1024) {
			matchRosterWinner(true)
		}
	} else {
		noxflags.SetGame(0x4000000)
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			if u.UpdateData != nil {
				pl := u.UpdateDataPlayer().Player
				b := []byte{154, 0, 0, 0, 0}
				binary.LittleEndian.PutUint16(b[1:], matchRosterTimerCoordinate(*(*float32)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(3632)) /* Pos3632Vec.X */))))
				binary.LittleEndian.PutUint16(b[3:], matchRosterTimerCoordinate(*(*float32)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(3636)) /* Pos3632Vec.Y */))))
				s.NetList.AddToMsgListCli(ntype.PlayerInd(pl.PlayerInd), netlist.Kind1, b)
			}
			s.Audio.EventObj(sound.ID(582), u, 2, u.NetCode)
		}
	}
	return int(serverConfigTimerSet(int32(0)))
}
