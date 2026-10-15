package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_inventoryServPlace_4F36F0 func(obj, it *server.Object, a3, a4 int) bool
)

func Nox_xxx_getObjectByScrName_4DA4F0(name string) *server.Object {
	return objectLookupByName(name)
}
func Nox_server_scriptMoveTo_5123C0(a1 *server.Object, a2 *server.Waypoint) {
	scriptBindingMove(a1, a2)
}
func Nox_xxx_playerCanCarryItem_513B00(a1 *server.Object, a2 *server.Object) {
	scriptInventoryCarry(a1, a2)
}

func Sub_516D00(a1 *server.Object) {
	monsterControlRevive(a1)
}
func Nox_xxx_netSendChat_528AC0(a1 *server.Object, a2 string, a3 uint16) {
	gameplayTextChat(a1, gameplayTextUnits(alloc.InternCString16(a2)), a3)
}

// NetSendChatTeam sends a team chat message (as players send with the team chat key) from an object
// to the players of its team. It does nothing if the object has no team.
func NetSendChatTeam(u *server.Object, text string, dur uint16) {
	if !u.TeamVal.Has() {
		return
	}
	team := u.TeamVal.ID
	b := gameplayTextMessage(gameplayTextUnits(alloc.InternCString16(text)), 1, gameplayReportCode(u), uint16(int64(u.PosVec.X)), uint16(int64(u.PosVec.Y)), dur)
	players := &GetServer().S().Players
	for p := players.FirstUnit(); p != nil; p = players.NextUnit(p) {
		if Nox_xxx_teamCompare2_419180(p.TeamPtr(), team) != 0 {
			gameplayReportDirect(gameplayReportRecipient(p), b)
		}
	}
}
