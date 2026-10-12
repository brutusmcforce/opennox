package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func teamRuntimeScan() bool {
	s := GetServer().S()
	teamRuntimeFlagCount = 0
	for u := s.Objs.List; u != nil; u = u.ObjNext {
		if u.ObjClass&0x10000000 != 0 {
			teamRuntimeFlagCount++
		}
	}
	if teamRuntimeFlagCount > 0 && !noxflags.HasGame(0x8000) {
		teamRuntimeBalance(false)
	}
	return teamRuntimeFlagCount > 0
}
func teamRuntimeBalance(reset bool) {
	s := GetServer().S()
	if reset {
		for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
			teamRuntimeClearPlayers(t)
		}
	}
	var units []*server.Object
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		u := pl.PlayerUnit
		if u != nil && (int(pl.NetCodeVal) != ClientPlayerNetCode() || !noxflags.HasEngine(noxflags.EngineNoRendering)) {
			flags := pl.Field3680
			if (flags&1 == 0 || flags&0x20 != 0) && !u.TeamVal.Has() {
				units = append(units, u)
			}
		}
	}
	if len(units) > 1 {
		for i := 0; i < 50; i++ {
			a := s.Rand.Logic.IntClamp(0, len(units)-1)
			b := s.Rand.Logic.IntClamp(0, len(units)-1)
			units[a], units[b] = units[b], units[a]
		}
	}
	for _, u := range units {
		if serverConfigSpecialMode() != 0 {
			group := u.UpdateDataPlayer().Player.Field2068
			if group != 0 {
				if t := s.Teams.ByXxx(int(group)); t != nil {
					teamRuntimeJoin(t.ID(), u.TeamPtr(), 1, int(u.NetCode), 0)
				}
			}
		} else {
			t := teamRuntimeLeast()
			teamRuntimeJoin(t.ID(), u.TeamPtr(), 1, int(u.NetCode), 1)
		}
	}
}
func teamRuntimeEnable() int {
	if byte(GetServer().S().Teams.Count()) == 0 {
		return 0
	}
	noxflags.SetGamePlay(2)
	teamRuntimeBalance(false)
	return 1
}
func teamRuntimeNearest() int {
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		if u.UpdateData != nil {
			pl := u.UpdateDataPlayer().Player
			flags := pl.Field3680
			if (int(pl.NetCodeVal) == ClientPlayerNetCode() && noxflags.HasEngine(noxflags.EngineNoRendering)) || (flags&1 != 0 && flags&0x20 == 0) {
				continue
			}
		}
		old := s.Teams.ByID(u.TeamVal.ID)
		var best *server.Team
		distance := float32(1e9)
		for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
			if flag := (*server.Object)(t.Field_72); flag != nil {
				// Match the original float subtraction, double square/sum and float best value.
				dx := float64(float32(flag.PosVec.X - u.PosVec.X))
				dy := float64(float32(flag.PosVec.Y - u.PosVec.Y))
				d := dx*dx + dy*dy
				if d < float64(distance) {
					distance = float32(d)
					best = t
				}
			}
		}
		if best != nil {
			if old != nil {
				if best != old {
					teamRuntimeSwitch(u.TeamPtr(), best, int(u.NetCode), 0)
				}
			} else {
				teamRuntimeJoin(best.ID(), u.TeamPtr(), 1, int(u.NetCode), 0)
			}
		}
	}
	return 0
}
func teamRuntimeCreateMap() int {
	s := GetServer().S()
	noxflags.SetGamePlay(4)
	s.Teams.Create(0)
	s.Teams.Create(0)
	for u := s.Objs.List; u != nil; u = u.ObjNext {
		if u.ObjClass&0x10000000 != 0 {
			if t := s.Teams.ByID(u.TeamVal.ID); t != nil {
				color := objectiveFlagID(u)
				teamRuntimeSetName(t, (*uint16)(unsafe.Pointer(internWStr(s.Teams.TeamTitle(server.TeamColor(color))))), 1)
				t.ColorInd = server.TeamColor(color)
				teamRuntimeAnnounce(t)
				t.Field_72 = u.CObj()
			}
		}
	}
	return 0
}
func teamRuntimeAssignFlags() int {
	s := GetServer().S()
	for u := s.Objs.List; u != nil; u = u.ObjNext {
		if u.ObjClass&0x10000000 != 0 {
			color := objectiveFlagID(u)
			if t := s.Teams.ByID(u.TeamVal.ID); t != nil {
				t.ColorInd = server.TeamColor(color)
				t.Field_72 = u.CObj()
			}
		}
	}
	return 0
}
func teamRuntimeToggle(on bool) {
	s := GetServer().S()
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		if u := (*server.Object)(t.Field_72); u != nil {
			if on {
				stateOn(u)
			} else {
				stateOff(u)
			}
		}
	}
}
func teamRuntimeCrown() int {
	s := GetServer().S()
	typ := memmap.Uint32(0x5D4594, 527652)
	if typ == 0 {
		typ = uint32(s.Types.IndByID("Crown"))
		*memmap.PtrUint32(0x5D4594, 527652) = typ
	}
	for u := s.Objs.List; u != nil; {
		next := u.ObjNext
		if u.ObjClass&0x10000000 != 0 {
			GetServer().DelayedDelete(u)
		}
		u = next
	}
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		*teamRuntimeWord(t, 76) = 0
	}
	n := 0
	for u := s.Objs.List; u != nil; {
		next := u.ObjNext
		if uint32(u.TypeInd) == typ {
			n++
			*(*uint32)(unsafe.Add(u.UpdateData, 4)) = 0
			bound := u.TeamVal.ID != 0
			if bound != noxflags.HasGamePlay(4) {
				GetServer().DelayedDelete(u)
				itemRespawnRemove(u)
			} else if !bound {
				playerStateMark(u, 1)
			} else if t := s.Teams.ByID(u.TeamVal.ID); t != nil {
				*teamRuntimeWord(t, 76) = uint32(uintptr(u.CObj()))
				playerStateMark(u, 1)
			}
		}
		u = next
	}
	if n > 0 {
		return 1
	}
	return 0
}
