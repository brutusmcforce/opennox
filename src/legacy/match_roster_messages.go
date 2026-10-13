package legacy

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// legacyRosterProtocolVersion preserves the unconditional NOX_HIGH_RES C flag
// used by the legacy roster in every qualified build profile.
const legacyRosterProtocolVersion uint32 = 0x000F039A

var matchRosterFlagType uint32

func matchRosterGUISettings(mode byte, data unsafe.Pointer, to int) int {
	b := make([]byte, 60)
	b[0] = 177
	b[1] = mode
	copy(b[2:], unsafe.Slice((*byte)(data), 58))
	return gameplayReportSend(to, b, true, 0)
}
func matchRosterInventory(to int, pl *server.Player) {
	u := pl.PlayerUnit
	if u == nil {
		return
	}
	if matchRosterFlagType == 0 {
		matchRosterFlagType = uint32(GetServer().S().Types.IndByID("Flag"))
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if uint32(it.ObjFlags)&0x100 != 0 || uint32(it.TypeInd) == matchRosterFlagType {
			gameplayReportEquipment(to, it)
		}
	}
}
func matchRosterPlayerIDs(pl *server.Player) int {
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		if u.NetCode == pl.NetCodeVal || u.UpdateDataPlayer().IsCamping == 0 {
			continue
		}
		code := uint16(s.GetUnitNetCode(u))
		b := []byte{210, byte(code), byte(code >> 8), byte(u.TypeInd), byte(u.TypeInd >> 8), 1, 2}
		gameplayReportSend(int(pl.PlayerInd), b, false, 1)
	}
	return 0
}
func matchRosterSendPlayers(to int) {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if int(pl.PlayerInd) == to || pl.PlayerInd == 31 && noxflags.HasEngine(noxflags.EngineNoRendering) {
			continue
		}
		var b [132]byte
		server.EncodePlayerRoster(b[:], pl)
		gameplayReportSend(to, b[:129], true, 0)
		matchRosterInventory(to, pl)
	}
}
func matchRosterMinimap(to int) int {
	s := GetServer().S()
	ball, crown := memmap.PtrUint32(0x5D4594, 1563284), memmap.PtrUint32(0x5D4594, 1563272)
	if *ball == 0 {
		*ball = uint32(s.Types.IndByID("GameBall"))
	}
	if *crown == 0 {
		*crown = uint32(s.Types.IndByID("Crown"))
	}
	for u := s.Objs.First(); u != nil; u = u.Next() {
		if uint32(u.ObjFlags)&4 != 0 && (uint32(u.TypeInd) == *ball || uint32(u.TypeInd) == *crown) {
			s.Players.Nox_xxx_netMarkMinimapObject_417190(ntype.PlayerInd(to), u, 1)
		}
	}
	return 0
}
func matchRosterSettings() int {
	settings := memmap.PtrOff(0x5D4594, 371380)
	mode := *controlHalf(settings, 52)
	var a [20]byte
	var b [49]byte
	a[0] = 175
	binary.LittleEndian.PutUint32(a[1:], GetServer().S().Frame())
	// The legacy C flag is unconditional today, including default/server targets.
	binary.LittleEndian.PutUint32(a[5:], legacyRosterProtocolVersion)
	binary.LittleEndian.PutUint32(a[9:], uint32(noxflags.GetGame())&0x7fff0)
	binary.LittleEndian.PutUint32(a[13:], uint32(dword_5d4594_3484))
	a[17] = byte(memmap.Uint32(0x5D4594, 3464))
	a[18] = byte(int16(serverConfigScore(int16(mode))))
	a[19] = byte(serverConfigMinutes(int16(mode)))
	b[0] = 176
	alloc.StrCopy(b[1:17], alloc.GoString(memmap.PtrUint8(0x5D4594, 1324)))
	copy(b[17:45], unsafe.Slice((*byte)(unsafe.Add(settings, 24)), 28))
	if int32(serverConfigTimerGet()) != 0 && (sub_40A300() != 0 || a[19] != 0) {
		binary.LittleEndian.PutUint32(b[45:], memmap.Uint32(0x5D4594, 3468)-uint32(PlatformTicks()))
	}
	gameplayReportSend(159, a[:], true, 0)
	return gameplayReportSend(159, b[:], true, 0)
}
func matchRosterWall(w *server.Wall, opcode byte) int {
	s := GetServer().S()
	b := []byte{opcode, byte(w.Field10), byte(w.Field10 >> 8)}
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		gameplayReportSend(int(u.UpdateDataPlayer().Player.PlayerInd), b, false, 1)
	}
	return 0
}
func matchRosterTeamRoster(to int) {
	if to == 31 {
		return
	}
	s := GetServer().S()
	for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
		teamRuntimeDescribe(tm, to)
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			if teamRuntimeContains(u.TeamPtr(), tm.IDVal) {
				teamRuntimeMemberReport(u.TeamPtr(), to, int(u.NetCode))
			}
		}
	}
}
func matchRosterSimpleObject(to int, u *server.Object) int {
	b := []byte{47, 0, 0, byte(u.TypeInd), byte(u.TypeInd >> 8), 0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[1:], uint16(GetServer().S().GetUnitNetCode(u)))
	binary.LittleEndian.PutUint16(b[5:], uint16(floatToInt32(u.PosVec.X)))
	binary.LittleEndian.PutUint16(b[7:], uint16(floatToInt32(u.PosVec.Y)))
	return gameplayReportSend(to, b, true, 1)
}
