package legacy

import (
	"encoding/binary"
	"unsafe"

	"fmt"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func teamRuntimeJoinText(t *server.Team, pl *server.Player) {
	sm := GetServer().S().Strings()
	var text string
	if noxflags.HasGame(4096) {
		text = fmt.Sprintf(sm.GetStringInFile("GeneralPrint:PlayerJoinQuest", "team.c"), pl.Name())
	} else {
		text = fmt.Sprintf(sm.GetStringInFile("NewMember", "team.c"), pl.Name(), t.Name())
	}
	Nox_xxx_printCentered_445490(text)
}
func teamRuntimeJoin(id server.TeamID, m *server.ObjectTeam, notify, code, relocate int) {
	s := GetServer().S()
	if teamRuntimeBallType == 0 {
		teamRuntimeBallType = uint32(s.Types.IndByID("GameBall"))
	}
	if m == nil {
		return
	}
	t := s.Teams.ByID(id)
	if t != nil {
		if teamRuntimeContains(m, id) {
			return
		}
	} else {
		t = s.Teams.Create(id)
	}
	m.ID = t.ID()
	m.Field0 = t.Field_44
	teamRuntimeSetFirst(t, m)
	if code == ClientPlayerNetCode() {
		teamUICTFSelect(byte(t.ID()))
	}
	if noxflags.HasGame(1) {
		if noxflags.HasGame(0x2000) {
			u := objectLookupByNetCode(uint32(code))
			pl := s.Players.ByID(code)
			if pl != nil {
				if noxflags.HasGame(0x8000) {
					statisticsParticipation((*server.Player)(pl.C()), 1)
				}
				if u != nil && u.ObjClass&4 != 0 {
					if relocate == 1 && !noxflags.HasGamePlay(2) && noxflags.HasGame(128) {
						flag := (*server.Object)(t.Field_72)
						var pos types.Pointf
						inventoryRandomPlacement(50, &flag.PosVec, &pos)
						Nox_xxx_unitMove_4E7010(u, pos)
					}
					if p := s.Players.ByID(code); p != nil {
						teamRuntimeJoinText(t, p)
					}
				}
			}
			if notify != 0 && u != nil && (pl != nil || uint32(u.TypeInd) == teamRuntimeBallType) {
				b := []byte{196, 1, 0, 0, 0, 0, 0, 0, 0, 0}
				binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
				binary.LittleEndian.PutUint16(b[6:], uint16(code))
				binary.LittleEndian.PutUint16(b[8:], u.TypeInd)
				teamUIPlayerTeam(code, int(t.ID()))
				gameplayReportSend(159, b, true, 1)
			}
		}
	} else {
		dr := GetClient().Cli().Objs.ByNetCodeDynamic(code)
		if dr != nil && dr.Class()&4 != 0 {
			if pl := s.Players.ByID(code); pl != nil {
				teamRuntimeJoinText(t, pl)
			}
		}
	}
	t.Field_48++
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		if u.NetCode == uint32(code) {
			teamRuntimeMembershipChanged(u)
			return
		}
	}
}
func teamRuntimeLeave(m *server.ObjectTeam, code int) {
	if m == nil {
		return
	}
	t := GetServer().S().Teams.ByID(m.ID)
	if t == nil || !teamRuntimeContains(m, m.ID) {
		return
	}
	if noxflags.HasGame(1) && noxflags.HasGame(0x2000) {
		teamUIPlayerTeam(code, 0)
		teamRuntimeSendID(159, 2, uint32(code))
	}
	teamRuntimeUnlink(t, m)
	t.Field_48--
	if (serverConfigSpecialMode() != 0 || noxflags.HasGame(0x8000)) && teamRuntimeCount(t) == 0 {
		if noxflags.HasGame(96) || (noxflags.HasGame(16) && noxflags.HasGamePlay(4)) {
			teamRuntimeSetGroup(t, 0)
			teamRuntimeSetName(t, (*uint16)(memmap.PtrOff(0x5D4594, 527664)), 0)
		} else {
			GetServer().TeamRemove(t, true)
		}
	}
}
func teamRuntimeSwitch(m *server.ObjectTeam, t *server.Team, code, relocate int) int {
	if m == nil || t == nil || !teamRuntimeContains(m, m.ID) {
		return 0
	}
	if noxflags.HasGame(1) && noxflags.HasGame(0x2000) {
		b := []byte{196, 3, 0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
		binary.LittleEndian.PutUint16(b[6:], uint16(code))
		gameplayReportSend(159, b, true, 1)
		teamUIPlayerTeam(code, int(t.ID()))
	}
	old := GetServer().S().Teams.ByID(m.ID)
	old.Field_48--
	teamRuntimeUnlink(old, m)
	teamRuntimeJoin(t.ID(), m, 0, code, relocate)
	if code == ClientPlayerNetCode() {
		teamUICTFSelect(byte(t.ID()))
	}
	return 1
}

// Keep pointer conversion local to the ABI boundary.
func teamRuntimeMember(p unsafe.Pointer) *server.ObjectTeam { return (*server.ObjectTeam)(p) }
