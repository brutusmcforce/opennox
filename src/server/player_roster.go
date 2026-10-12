package server

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

// EncodePlayerRoster fills the existing 129-byte roster message or a larger
// scratch buffer. The caller owns initialization of unused identifier padding.
func EncodePlayerRoster(buf []byte, pl *Player) {
	buf[0] = 45 // MSG_NEW_PLAYER
	binary.LittleEndian.PutUint16(buf[1:], uint16(pl.NetCode()))
	binary.LittleEndian.PutUint16(buf[100:], uint16(pl.Lessons))
	binary.LittleEndian.PutUint16(buf[102:], uint16(pl.Field2140))
	binary.LittleEndian.PutUint32(buf[104:], pl.ArmorEquip)
	binary.LittleEndian.PutUint32(buf[108:], pl.WeaponEquip)
	binary.LittleEndian.PutUint32(buf[112:], pl.Field3680&0x423)
	buf[116] = byte(pl.Field2152)
	buf[117] = byte(pl.Field2156)
	buf[118] = 0
	if pl.Field3676 == 3 {
		buf[118] = 1
	}
	alloc.StrCopy(buf[119:], pl.Field2096())
	*(*PlayerInfo)(unsafe.Pointer(&buf[3])) = *pl.Info()
}
