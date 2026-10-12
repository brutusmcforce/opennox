package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

type nox_team_t = server.Team

func asTeamP(p unsafe.Pointer) *server.Team {
	if p == nil {
		return nil
	}
	return (*server.Team)(p)
}

func nox_xxx_getTeamByID_418AB0(a1 int) *nox_team_t {
	return (*nox_team_t)(GetServer().S().Teams.ByID(server.TeamID(a1)).C())
}

func Sub_459CD0() {
	serverOptionsTeamCount()
}
func Sub_456FA0() {
	teamUITeamClear()
}
func Sub_418E40(t *server.Team, p *server.ObjectTeam) {
	teamRuntimeUnlink(t, p)
}
func Sub_456EA0(name string) {
	teamUITeamRemove(name)
}
