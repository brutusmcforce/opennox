package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func playerStateLessons(send int32) {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		pl.Lessons = 0
		pl.Field2140 = 0
		if send != 0 && pl.PlayerUnit != nil {
			gameplayReportLesson(pl.PlayerUnit)
		}
	}
}

// The existing wide comparison runs in the process's C locale. Compare code
// units so non-ASCII UTF-16, including unmatched surrogates, stays unchanged.
func playerStateByName(name *uint16) *server.Player {
	if name == nil {
		return nil
	}
	fold := func(v uint16) uint16 {
		if v >= 'A' && v <= 'Z' {
			return v + ('a' - 'A')
		}
		return v
	}
	equal := func(a, b *uint16) bool {
		for i := uintptr(0); ; i += 2 {
			x := *(*uint16)(unsafe.Add(unsafe.Pointer(a), i))
			y := *(*uint16)(unsafe.Add(unsafe.Pointer(b), i))
			if fold(x) != fold(y) {
				return false
			}
			if x == 0 {
				return true
			}
		}
	}
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if equal(&pl.NameFinal[0], name) {
			return pl
		}
	}
	return nil
}
func playerStateTracks(index int, obj *server.Object) int {
	if index < 0 || index >= 32 || obj == nil {
		return 0
	}
	head := GetServer().S().Players.ByIndRaw(ntype.PlayerInd(index)).Field4580
	for m := head; m != nil; {
		if m.Field4 == obj {
			return 1
		}
		m = m.Field8
		if m == head {
			break
		}
	}
	return 0
}
func playerStateMark(obj *server.Object, flags uint32) {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		s.Players.Nox_xxx_netMarkMinimapObject_417190(pl.PlayerIndex(), obj, flags)
	}
}
func playerStateUnmark(obj *server.Object, flags uint32) {
	GetServer().S().Players.Nox_xxx_netUnmarkMinimapSpec_417470(obj, flags)
}
func playerStateStatusData(pl *server.Player) []byte {
	b := make([]byte, 7)
	b[0] = 106
	binary.LittleEndian.PutUint16(b[1:], uint16(pl.NetCodeVal))
	binary.LittleEndian.PutUint32(b[3:], pl.Field3680&0x423)
	return b
}
func playerStateReport(pl *server.Player) int32 {
	return int32(gameplayReportSend(255, playerStateStatusData(pl), true, 0))
}
func playerStateAllStatus(to *server.Player) {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		gameplayReportSend(int(to.PlayerInd), playerStateStatusData(pl), true, 0)
	}
}
func playerStateAddStatus(pl *server.Player, mask uint32) int32 {
	pl.Field3680 |= mask
	if !noxflags.HasGame(1) {
		return 0
	}
	if mask&0x423 != 0 {
		return playerStateReport(pl)
	}
	return 1
}
func playerStateRemoveStatus(pl *server.Player, mask uint32) int8 {
	pl.Field3680 &^= mask
	if !noxflags.HasGame(1) {
		return 0
	}
	result := int32(1)
	if mask&1 != 0 {
		result = int32(bool2int(noxflags.HasGame(128)))
		if result == 0 {
			result = int32(playerStateMultiple())
			if result != 0 {
				result = serverConfigTimerGet()
				if result == 0 {
					result = int32(serverConfigMinutes(int16(noxflags.GetGame())))
					if result != 0 {
						serverConfigTimerInit()
						result = serverConfigTimerSet(1)
					}
				}
			}
		}
	}
	if mask&0x423 != 0 {
		result = playerStateReport(pl)
	}
	return int8(result)
}
