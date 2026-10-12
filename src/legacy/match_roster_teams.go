package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func matchRosterAssignTeam(pl *server.Player) {
	u := pl.PlayerUnit
	if u == nil {
		return
	}
	s := GetServer().S()
	if serverConfigSpecialMode() == 0 && !noxflags.HasGame(0x8000) {
		if byte(s.Teams.Count()) != 0 {
			tm := teamRuntimeLeast()
			if tm != nil && !u.TeamVal.Has() {
				Nox_xxx_createAtImpl_4191D0(tm.IDVal, u.TeamPtr(), 1, int(u.NetCode), 1)
			}
		}
		return
	}
	if pl.Field2068 == 0 {
		return
	}
	tm := s.Teams.ByXxx(int(pl.Field2068))
	if tm == nil {
		capacity := int(memmap.Uint8(0x5D4594, 371516+52))
		if (noxflags.HasGame(96) || noxflags.HasGame(16) && noxflags.HasGamePlay(4)) && capacity > 2 {
			capacity = 2
		}
		count := func() int { return int(byte(teamRuntimeGroupCount())) }
		if count() >= capacity {
			return
		}
		if noxflags.HasGame(96) && count() >= int(teamRuntimeFlagCount) {
			return
		}
		tm = teamRuntimeAvailable()
		if tm == nil {
			return
		}
		teamRuntimeSetName(tm, &pl.Field2072[0], 0)
		teamRuntimeSetGroup(tm, pl.Field2068)
		teamRuntimeAnnounce(tm)
	}
	if u.TeamVal.Has() {
		teamRuntimeSwitch(u.TeamPtr(), tm, int(u.NetCode), 0)
	} else {
		Nox_xxx_createAtImpl_4191D0(tm.IDVal, u.TeamPtr(), 1, int(u.NetCode), 0)
	}
}
