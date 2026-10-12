package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func playerStateEligible(pl *server.Player) bool {
	return pl.Field3680&1 == 0 && (pl.PlayerInd != 31 || !noxflags.HasEngine(noxflags.EngineNoRendering))
}
func playerStateCompetitors() int {
	s := GetServer().S()
	n := 0
	if !noxflags.HasGamePlay(4) {
		for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
			if playerStateEligible(pl) {
				n++
			}
		}
		return n
	}
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			if teamRuntimeContains(u.TeamPtr(), tm.ID()) && playerStateEligible(u.UpdateDataPlayer().Player) {
				n++
				break
			}
		}
	}
	return n
}
func playerStateTeamPlayers(tm *server.Team) int {
	s := GetServer().S()
	n := 0
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		if teamRuntimeContains(u.TeamPtr(), tm.ID()) && playerStateEligible(u.UpdateDataPlayer().Player) {
			n++
		}
	}
	return n
}
func playerStateMultiple() int {
	s := GetServer().S()
	n := 0
	participates := func(pl *server.Player) bool { return pl.Field3680&1 == 0 || pl.Field3680&0x20 != 0 }
	if !noxflags.HasGamePlay(4) {
		for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
			if pl.PlayerUnit != nil && participates(pl) {
				n++
			}
		}
	} else {
		for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
			for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
				if teamRuntimeContains(u.TeamPtr(), tm.ID()) && participates(u.UpdateDataPlayer().Player) {
					n++
					break
				}
			}
		}
	}
	return bool2int(n > 1)
}
func playerStateReset() {
	s := GetServer().S()
	*memmap.PtrUint32(0x5D4594, 3520) = s.Frame()
	*memmap.PtrUint32(0x5D4594, 3536) = 0
	*memmap.PtrInt32(0x5D4594, 3476) = floatToInt32(float32(s.Balance.Float("SuddenDeathPlayerThreshold")))
	*memmap.PtrInt32(0x5D4594, 1392) = floatToInt32(float32(s.Balance.Float("SuddenDeathLifeTime")))
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if pl.Field3680&256 != 0 {
			playerStateRemoveStatus(pl, 256)
		}
	}
	noxflags.UnsetGame(0x4000000)
}
func playerStateElapsed() int {
	s := GetServer().S()
	return bool2int(s.Frame()-memmap.Uint32(0x5D4594, 3520) > 20*s.TickRate())
}
func playerStateThreshold() int32          { return memmap.Int32(0x5D4594, 3476) }
func playerStateReentry(value int32) int32 { *memmap.PtrInt32(0x5D4594, 3508) = value; return value }
func playerStateAdmission(pl *server.Player) int {
	s := GetServer().S()
	if pl != nil {
		if noxflags.HasGame(4096) {
			return bool2int(playerStateCompetitors() < 6)
		}
		if pl.Field3680&256 != 0 && memmap.Uint32(0x5D4594, 3508) == 0 {
			return 0
		}
		if serverConfigSpecialMode() != 0 || noxflags.HasGame(0x8000) {
			if pl.Field2068 == 0 {
				return 0
			}
			if s.Teams.ByXxx(int(pl.Field2068)) == nil {
				capacity := int(memmap.Uint8(0x5D4594, 371516+52))
				if (noxflags.HasGame(96) || noxflags.HasGame(16) && noxflags.HasGamePlay(4)) && capacity > 2 {
					capacity = 2
				}
				if int(byte(teamRuntimeGroupCount())) >= capacity {
					return 0
				}
				if noxflags.HasGame(96) && int32(byte(teamRuntimeGroupCount())) >= int32(teamRuntimeFlagCount) {
					return 0
				}
			}
		}
	}
	if noxflags.HasGame(128) || !noxflags.HasGame(1024) {
		return 1
	}
	for p := s.Players.First(); p != nil; p = s.Players.Next(p) {
		if int32(p.Field2140) > 0 {
			return 1 - playerStateElapsed()
		}
	}
	return 1
}
