package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func matchRosterClearMask(slot byte) int {
	mask := uint32(1) << (slot & 31)
	for u := GetServer().S().Objs.First(); u != nil; u = u.Next() {
		u.Field35 &^= mask
		u.Field36 &^= mask
	}
	return 0
}
func matchRosterResyncMask(slot byte) int {
	mask := uint32(1) << (slot & 31)
	for u := GetServer().S().Objs.First(); u != nil; u = u.Next() {
		u.Field38 |= mask
		if uint32(u.ObjFlags)&0x20 == 0 && uint32(u.ObjClass)&0x20400006 == 0 {
			u.Field37 &^= mask
		}
	}
	return 0
}
func matchRosterReportMask(slot int) {
	mask := uint32(1) << uint(slot&31)
	for u := GetServer().S().Objs.First(); u != nil; u = u.Next() {
		u.Field38 |= mask
		if uint32(u.ObjClass)&0x20400000 == 0 {
			u.Field37 &^= mask
		}
		u.Field140[slot] &= 0xfff
		if uint32(u.ObjClass)&0x20400000 != 0 {
			for bit := uint(1); bit < 0x10000; bit <<= 1 {
				if u.Sub_4E4C90(bit) {
					u.Field140[slot] |= uint32(bit) << 16
				}
			}
		}
	}
}
func matchRosterFlagBase() unsafe.Pointer { return memmap.PtrOff(0x5D4594, 1567736) }
func matchRosterFlagRecord(index byte) unsafe.Pointer {
	return memmap.PtrOff(0x5D4594, 1567740+6*uintptr(index))
}
func matchRosterFlagState(index, a, b byte, value uint16) int {
	p := matchRosterFlagRecord(index)
	*controlByte(p, 0) = index
	*controlByte(p, 1) = b
	*controlByte(p, 2) = a
	*controlHalf(p, 4) = value
	return gameplayReportFlag(255, index, a, b, value)
}

type matchRosterIdentity struct {
	Name  string
	Group uint32
	Class byte
}

var matchRosterRemembered []matchRosterIdentity
var matchRosterRememberedInit uint32

func matchRosterRemember(pl *server.Player) {
	matchRosterRememberedInit = 1
	id := matchRosterIdentity{pl.Field2096(), pl.Field2068, byte(pl.PlayerClass())}
	matchRosterRemembered = append(matchRosterRemembered, matchRosterIdentity{})
	copy(matchRosterRemembered[1:], matchRosterRemembered[:len(matchRosterRemembered)-1])
	matchRosterRemembered[0] = id
}
func matchRosterForget() bool {
	had := len(matchRosterRemembered) != 0
	matchRosterRemembered = nil
	return had
}
func matchRosterIdentityAllowed(name string, class byte, group uint32) bool {
	for _, id := range matchRosterRemembered {
		if id.Name == name && (int(int8(class)) != int(id.Class) || group != id.Group) {
			return false
		}
	}
	return true
}
func matchRosterHasIdentity(pl *server.Player) bool {
	name := pl.Field2096()
	for _, id := range matchRosterRemembered {
		if id.Name == name {
			return true
		}
	}
	return false
}
