package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

var teamRuntimeFlagCount uint32
var teamRuntimeBallType uint32

// These fields retain the shared 386 Team layout; list links are C-allocated addresses.
func teamRuntimeWord(t *server.Team, off uintptr) *uint32 { return (*uint32)(unsafe.Add(t.C(), off)) }
func teamRuntimeFirst(t *server.Team) *server.ObjectTeam {
	if t == nil {
		return nil
	}
	return (*server.ObjectTeam)(memmap.WordPtr(t.Field_44))
}
func teamRuntimeNext(m *server.ObjectTeam) *server.ObjectTeam {
	if m == nil {
		return nil
	}
	return (*server.ObjectTeam)(memmap.WordPtr(m.Field0))
}
func teamRuntimeSetFirst(t *server.Team, m *server.ObjectTeam) {
	t.Field_44 = memmap.PtrWord(m.C()) // a list of words: members link through ObjectTeam.Field0
}
func teamRuntimeSetName(t *server.Team, name *uint16, flag uint32) {
	if t == nil {
		return
	}
	buf := unsafe.Slice((*uint16)(t.C()), 22)
	var text [20]uint16
	for i := range text {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(name), 2*i))
		if c == 0 {
			break
		}
		text[i] = c
	}
	copy(buf[:20], text[:])
	buf[20] = 0
	t.Field_68 = flag
}
func teamRuntimeSetGroup(t *server.Team, group uint32) {
	if t != nil {
		t.Field_60 = group
	}
}
func teamRuntimeObject(code int) *server.ObjectTeam {
	if noxflags.HasGame(1) {
		if u := objectLookupByNetCode(uint32(code)); u != nil {
			return u.TeamPtr()
		}
	} else {
		if dr := GetClient().Cli().Objs.ByNetCodeDynamic(code); dr != nil {
			return dr.TeamPtr()
		}
	}
	return nil
}
func teamRuntimeContains(m *server.ObjectTeam, id server.TeamID) bool {
	if m == nil || m.ID != id {
		return false
	}
	t := GetServer().S().Teams.ByID(id)
	for p := teamRuntimeFirst(t); p != nil; p = teamRuntimeNext(p) {
		if p == m {
			return true
		}
	}
	return false
}
func teamRuntimeCount(t *server.Team) int {
	if t == nil {
		return 0
	}
	s := GetServer().S()
	n := 0
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if teamRuntimeContains(teamRuntimeObject(int(pl.NetCodeVal)), t.ID()) {
			n++
		}
	}
	return n
}
func teamRuntimeLeast() *server.Team {
	s := GetServer().S()
	var best *server.Team
	count := byte(32)
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		n := byte(teamRuntimeCount(t))
		if n < count {
			best, count = t, n
		}
	}
	return best
}
func teamRuntimeAvailable() *server.Team {
	s := GetServer().S()
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		if t.Field_60 == 0 {
			return t
		}
	}
	return s.Teams.Create(0)
}
func teamRuntimeFind(name *uint16) *server.Team {
	s := GetServer().S()
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		if textCompareWide((*uint16)(t.C()), name) == 0 {
			return t
		}
	}
	return nil
}
func teamRuntimeGroupCount() int8 {
	s := GetServer().S()
	var n byte
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		if t.Field_60 != 0 {
			n++
		}
	}
	return int8(n)
}
func teamRuntimeMembershipChanged(u *server.Object) {
	pl := u.UpdateDataPlayer().Player
	gameplayReportFriendReset(int(pl.PlayerInd))
	statePlayerVisibility(int32(pl.PlayerInd))
	u.Nox_xxx_monsterMarkUpdate_4E8020()
	for it := u.FirstOwned516(); it != nil; it = it.NextOwned512() {
		if it.ObjClass&6 != 0 {
			it.Nox_xxx_monsterMarkUpdate_4E8020()
		}
	}
}

// Unlink owns list and visibility updates. Its caller owns the member counter.
func teamRuntimeUnlink(t *server.Team, m *server.ObjectTeam) {
	if t == nil || m == nil {
		return
	}
	if teamRuntimeFirst(t) == m {
		teamRuntimeSetFirst(t, teamRuntimeNext(m))
	} else {
		p := teamRuntimeFirst(t)
		for p != nil && teamRuntimeNext(p) != m {
			p = teamRuntimeNext(p)
		}
		if p == nil {
			return
		}
		p.Field0 = m.Field0
	}
	m.Field0 = 0
	m.ID = 0
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		if u.TeamPtr() == m {
			teamRuntimeMembershipChanged(u)
			return
		}
	}
}
func teamRuntimeClearPlayers(t *server.Team) {
	if t == nil || teamRuntimeCount(t) <= 0 {
		return
	}
	if noxflags.HasGame(1) {
		teamRuntimeSendID(159, 5, uint32(t.ID()))
	}
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		m := teamRuntimeObject(int(pl.NetCodeVal))
		if m != nil && m.ID == t.ID() {
			teamUIPlayerTeam(int(pl.NetCodeVal), 0)
			teamRuntimeUnlink(t, m)
			if m.ID == 0 {
				t.Field_48--
			}
		}
	}
}
