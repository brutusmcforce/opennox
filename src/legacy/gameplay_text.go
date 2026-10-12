package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// Keep UTF16 code units intact; protocol strings are not necessarily valid Unicode.
func gameplayTextUnits(p *uint16) []uint16 {
	n := 0
	for *(*uint16)(unsafe.Add(unsafe.Pointer(p), 2*n)) != 0 {
		n++
	}
	return unsafe.Slice(p, n)
}
func gameplayTextByteEncoding(text []uint16) bool {
	for _, v := range text {
		if v > 255 {
			return false
		}
	}
	return true
}
func gameplayTextMessage(text []uint16, flags byte, code, x, y, extra uint16) []byte {
	if gameplayTextByteEncoding(text) {
		flags |= 2
	} else {
		flags |= 4
	}
	width := 1
	if flags&4 != 0 {
		width = 2
	}
	count := byte(len(text) + 1)
	b := make([]byte, 11+width*int(count))
	b[0] = 168
	b[3] = flags
	b[8] = count
	binary.LittleEndian.PutUint16(b[1:], code)
	binary.LittleEndian.PutUint16(b[4:], x)
	binary.LittleEndian.PutUint16(b[6:], y)
	binary.LittleEndian.PutUint16(b[9:], extra)
	for i := 0; i < int(count) && i < len(text); i++ {
		if width == 2 {
			binary.LittleEndian.PutUint16(b[11+2*i:], text[i])
		} else {
			b[11+i] = byte(text[i])
		}
	}
	return b
}
func gameplayTextFanout(b []byte) int {
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		gameplayReportDirect(gameplayReportRecipient(u), b)
	}
	return 0
}
func gameplayTextLine(u *server.Object, text []uint16) int {
	if u == nil || uint32(u.ObjClass)&4 == 0 {
		return gameplayReportPtr(u.CObj())
	}
	return gameplayReportDirect(gameplayReportRecipient(u), gameplayTextMessage(text, 0, 0, 0, 0, 0))
}
func gameplayTextAll(flags byte, text []uint16) int {
	return gameplayTextFanout(gameplayTextMessage(text, flags, 0, 0, 0, 0))
}
func gameplayTextInformation(to, kind int, data unsafe.Pointer) int {
	switch kind {
	case 0, 1, 2, 12, 13, 16, 20, 21:
		b := []byte{169, byte(kind), 0, 0, 0, 0}
		copy(b[2:], unsafe.Slice((*byte)(data), 4))
		return gameplayReportDirect(to, b)
	case 17:
		return gameplayReportDirect(to, []byte{169, 17})
	default:
		return kind
	}
}
func gameplayTextInformationAll(kind int, data unsafe.Pointer) int {
	switch kind {
	case 3, 4, 8, 18, 19, 21:
		b := []byte{169, byte(kind), 0, 0, 0, 0}
		copy(b[2:], unsafe.Slice((*byte)(data), 4))
		return gameplayTextFanout(b)
	case 5, 6, 7, 9, 10, 11, 14:
		size := 10
		if kind == 14 {
			size = 11
		}
		b := unsafe.Slice((*byte)(data), size)
		b[0], b[1] = 169, byte(kind)
		return gameplayTextFanout(b)
	default:
		return kind
	}
}
func gameplayTextPrivate(u *server.Object, text *byte, flag byte) {
	if u == nil || uint32(u.ObjClass)&4 == 0 || text == nil || GetServer().S().Players.CheckXxx(u) {
		return
	}
	s := alloc.GoString(text)
	if len(s) == 0 || len(s) > 48 {
		return
	}
	b := make([]byte, len(s)+4)
	b[0], b[1], b[2] = 169, 15, flag
	copy(b[3:], s)
	gameplayReportDirect(gameplayReportRecipient(u), b)
}
func gameplayTextPrivateAll(text *byte) int {
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		gameplayTextPrivate(u, text, 0)
	}
	return 0
}
func gameplayTextChat(u *server.Object, text []uint16, extra uint16) int {
	return gameplayTextFanout(gameplayTextMessage(text, 0, gameplayReportCode(u), uint16(int64(u.PosVec.X)), uint16(int64(u.PosVec.Y)), extra))
}
