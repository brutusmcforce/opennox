package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unicode/utf16"
	"unsafe"
)

func clientSessionTeam(data []byte) int {
	if len(data) < 2 {
		return -1
	}
	size := 0
	switch data[1] {
	case 0:
		if len(data) < 18 {
			return -1
		}
		size = 18 + 2*int(data[15])
	case 1, 3, 8:
		size = 10
	case 2, 5, 6:
		size = 6
	case 4:
		size = 46
	case 7, 9:
		size = 2
	case 12:
		size = 5
	default:
		return -1
	}
	if len(data) < size {
		return -1
	}
	if !Nox_client_isConnected() {
		return size
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	s := GetServer().S()
	switch data[1] {
	case 0:
		// The C temporary held only 20 units; keep the full wire field safely
		// terminated and let the team owner apply its established storage limit.
		name := make([]uint16, int(data[15])+1)
		for i := 0; i < int(data[15]); i++ {
			name[i] = word(18 + 2*i)
		}
		if data[17] != 0 {
			name = append(utf16.Encode([]rune(s.Teams.TeamTitle(server.TeamColor(data[16])))), 0)
		}
		t := s.Teams.ByID(server.TeamID(dword(2)))
		if t == nil {
			t = s.Teams.Create(server.TeamID(data[2]))
		}
		if t == nil {
			return size
		}
		teamRuntimeSetName(t, &name[0], 0)
		teamRuntimeSetGroup(t, dword(6))
		teamRuntimeLessons(t, int(dword(10)))
		t.ColorInd = server.TeamColor(data[16])
		teamUITeamAdd(alloc.GoString16(&name[0]))
		if data[14]&1 != 0 {
			if m := teamRuntimeObject(ClientPlayerNetCode()); m != nil {
				if noxflags.HasGame(1) {
					teamRuntimeJoin(t.ID(), m, 1, ClientPlayerNetCode(), 1)
				} else {
					teamRuntimeRequest(t, m, int16(ClientPlayerNetCode()), 10)
				}
			}
		}
	case 1:
		code := word(6)
		dr := clientSessionDrawable(code)
		if dr == nil {
			dr = GetClient().Nox_xxx_spriteCreate_48E970(int(word(8)), code&0x7fff, 0, 0)
		}
		if dr != nil {
			if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
				teamRuntimeJoin(t.ID(), dr.TeamPtr(), 0, int(code), 0)
				teamUIPlayerTeam(int(code), int(t.ID()))
			}
		}
	case 2:
		code := int(dword(2))
		if m := teamRuntimeObject(code); m != nil {
			teamRuntimeLeave(m, code)
			teamUIPlayerTeam(code, 0)
		}
	case 3:
		code := int(word(6))
		if m := teamRuntimeObject(code); m != nil {
			if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
				if teamRuntimeSwitch(m, t, code, 0) != 0 {
					teamUIPlayerTeam(code, int(t.ID()))
				}
			}
		}
	case 4:
		if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
			var name [21]uint16
			for i := 0; i < 20; i++ {
				name[i] = word(6 + 2*i)
			}
			teamRuntimeRename(t, &name[0])
		}
	case 5:
		if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
			teamRuntimeClearPlayers(t)
		}
	case 6:
		if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
			name := t.Name()
			GetServer().TeamRemove(t, false)
			teamUITeamRemove(name)
		}
	case 7:
		GetServer().TeamsRemoveActive(false)
		teamUITeamClear()
	case 8:
		if t := s.Teams.ByID(server.TeamID(dword(2))); t != nil {
			teamRuntimeLessons(t, int(dword(6)))
		}
	case 9:
		GetServer().TeamsResetYyy()
	case 12:
		if pl := s.Players.ByID(int(word(2))); pl != nil {
			*(*byte)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(2282)) /* Field2282 */)) = data[4]
		}
	}
	return size
}
