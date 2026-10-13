package legacy

import (
	"encoding/binary"
	"strings"

	"github.com/brutusmcforce/opennox/server"
)

func visibilitySpecialUpdate(viewer, u *server.Object) int {
	ind := int(viewer.UpdateDataPlayer().Player.PlayerInd)
	if u == nil || ind >= 32 {
		return 1
	}
	flags := &u.Field140[ind]
	if *flags&0x0FFF0000 == 0 {
		return 0
	}
	if *flags&0x10000 != 0 {
		gameplayReportAnimation(ind, u)
		*flags &^= 0x10000
	}
	if *flags&0x20000 != 0 {
		if !GetServer().S().IsEnemyTo(viewer, u) {
			gameplayReportCurrentHealth(ind, u)
		}
		*flags &^= 0x20000
	}
	if *flags&0x40000 != 0 {
		gameplayReportHidden(ind, u)
		*flags &^= 0x40000
	}
	if *flags&0x80000 != 0 {
		gameplayReportXStatus(ind, u)
		*flags &^= 0x80000
	}
	if *flags&0x400000 != 0 {
		gameplayReportHeight(ind, u)
		*flags &^= 0x400000
	}
	if *flags&0x800000 != 0 {
		gameplayReportEnchant(ind, u)
		*flags &^= 0x800000
	}
	if *flags&0x2000000 != 0 {
		gameplayReportTeamBase(ind, u)
		*flags &^= 0x2000000
	}
	if *flags&0x4000000 != 0 {
		gameplayReportNPC(ind, u)
		*flags &^= 0x4000000
	}
	return 1
}
func visibilityDestroyReport(u *server.Object) {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		ind := int(pl.PlayerInd)
		if uint32(1)<<uint(ind&31)&u.Field37 != 0 {
			code := gameplayReportCode(u)
			gameplayReportSend(ind, []byte{byte(u.Field5)>>6 | 49, byte(code), byte(code >> 8)}, false, 1)
		}
		if uint32(u.ObjClass)&6 != 0 {
			gameplayReportFriend(ind, u, 0)
		}
	}
}
func visibilityOutOfSight(to int, u *server.Object) int {
	return gameplayReportSend(to, gameplayReportID(50, u, 3), false, 1)
}
func visibilityInShadows(to int, u *server.Object) int {
	return gameplayReportSend(to, gameplayReportID(51, u, 3), false, 1)
}
func visibilityMonsterCommand(u, source *server.Object, command string, id uint16) int {
	if i := strings.IndexByte(command, 0); i >= 0 {
		command = command[:i]
	}
	b := []byte{168}
	b = binary.LittleEndian.AppendUint16(b, gameplayReportCode(u))
	b = append(b, 8)
	b = binary.LittleEndian.AppendUint16(b, uint16(int64(u.PosVec.X)))
	b = binary.LittleEndian.AppendUint16(b, uint16(int64(u.PosVec.Y)))
	n := byte(len(command) + 1)
	b = append(b, n)
	b = binary.LittleEndian.AppendUint16(b, id)
	b = append(b, command...)
	b = append(b, 0)
	b = b[:11+int(n)]
	if source == nil {
		return visibilityFXBroadcast(b)
	}
	if uint32(source.ObjClass)&4 != 0 {
		return gameplayReportDirect(int(source.UpdateDataPlayer().Player.PlayerInd), b)
	}
	return int(uintptr(source.CObj()))
}
func visibilityClearChats() int { return visibilityFXBroadcast([]byte{202, 173, 222}) }
