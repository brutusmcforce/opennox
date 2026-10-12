package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
)

func objectReportDirection(u *server.Object) byte {
	var dir [2]int32
	geometryIndexedDirection(int32(int16(u.Direction1)), &dir)
	index := byte(dir[0]) + 3*byte(dir[1]) + 4
	if index > 3 {
		index--
	}
	return index << 4
}
func objectReportBase(u *server.Object, op byte, size int) []byte {
	b := make([]byte, size)
	b[0] = op
	binary.LittleEndian.PutUint16(b[1:], uint16(GetServer().S().GetUnitNetCode(u)))
	binary.LittleEndian.PutUint16(b[3:], u.TypeInd)
	binary.LittleEndian.PutUint16(b[5:], uint16(floatToInt32(u.PosVec.X)))
	binary.LittleEndian.PutUint16(b[7:], uint16(floatToInt32(u.PosVec.Y)))
	return b
}
func objectReportHealth(to int, u *server.Object, offset int) {
	if u.HealthData == nil || GetServer().S().Frame()-u.Frame134 <= 2 {
		return
	}
	previous := (*uint16)(unsafe.Add(u.UpdateData, offset+2*to))
	if current := u.HealthData.Cur; current != *previous {
		gameplayReportHealthDelta(to, uint16(u.NetCode), int16(current-*previous))
		*previous = u.HealthData.Cur
	}
}
func objectReportSimple(to int, u *server.Object) int {
	if uint32(u.ObjClass)&0x20000 != 0 {
		objectReportHealth(to, u, 96)
	}
	return bool2int(Nox_netlist_addToMsgListSrv(ntype.PlayerInd(to), objectReportBase(u, 47, 9)))
}
func objectReportPhantom(to int, u *server.Object) int {
	b := objectReportBase(u, 48, 11)
	b[9] = objectReportDirection(u)
	b[10] = 255
	if uint32(u.ObjClass)&1 != 0 && uint32(u.ObjSubClass)&0x30 != 0 {
		b[10] = byte(u.Direction1 >> 3)
	}
	return bool2int(Nox_netlist_addToMsgListSrv(ntype.PlayerInd(to), b))
}
func objectReportMonster(to int, u *server.Object, health int) int {
	if health != 0 {
		objectReportHealth(to, u, 412)
	}
	b := objectReportBase(u, 48, 11)
	ud := u.UpdateDataMonster()
	b[9] = byte(monsterActionToAnimation(u))
	b[10] = ud.Field120_1
	s := GetServer().S()
	if uint32(u.ObjSubClass)&0x10 != 0 && ud.Field523_2 != 0 && s.Rand.Logic.IntClamp(0, 10) >= 8 {
		b[9] = 14
		frames, _ := s.PlayerAnimFrames(50)
		b[10] = byte(s.Rand.Logic.IntClamp(0, frames))
	}
	b[9] |= objectReportDirection(u)
	return bool2int(Nox_netlist_addToMsgListSrv(ntype.PlayerInd(to), b))
}
func objectReportPlayer(viewer, u *server.Object, updates, health int) int {
	s := GetServer().S()
	viewerData := viewer.UpdateDataPlayer()
	ud := u.UpdateDataPlayer()
	to := int(viewerData.Player.PlayerInd)
	if health != 0 {
		objectReportHealth(to, u, 12)
	}
	if viewer == u {
		gameplayReportAnything(viewer)
	}
	*(*uint32)(unsafe.Add(unsafe.Pointer(viewerData.Player), 4452+4*int(ud.Player.PlayerInd))) = s.Frame()
	b := objectReportBase(u, 195, 12)
	b[11] = byte(controlActionState(u))
	b[10] = 255
	switch ud.State {
	case 1, 10, 2, 15, 16, 17, 14, 20, 18, 19, 21, 22, 24, 25, 27, 28, 29, 26, 30, 32:
		b[10] = *(*byte)(unsafe.Add(u.UpdateData, 236))
	}
	if *(*uint16)(unsafe.Add(u.UpdateData, 160)) != 0 && s.Rand.Logic.IntClamp(0, 10) >= 8 {
		b[11], b[10] = 50, 255
	}
	if (ud.State == 3 || ud.State == 4) && u.Field131 == 16 {
		b[11], b[10] = 51, 255
	}
	reaction := (*uint32)(unsafe.Add(u.UpdateData, 164))
	if ud.State != 30 && *reaction != 0 {
		frames, delay := s.PlayerAnimFrames(52)
		elapsed := s.Frame() - *reaction
		frame := int32(elapsed / (uint32(delay) + 1))
		if frame >= int32(frames) || elapsed >= 4 {
			*reaction = 0
		} else {
			b[11], b[10] = 52, byte(frame)
		}
	}
	b[9] = objectReportDirection(u)
	to = int(viewerData.Player.PlayerInd)
	if updates != 0 {
		return bool2int(Nox_netlist_addToMsgListSrv(ntype.PlayerInd(to), b))
	}
	return bool2int(s.NetList.AddToMsgListCli(ntype.PlayerInd(to), netlist.Kind1, b))
}
