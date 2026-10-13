package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func teamRuntimeNameUnits(t *server.Team) []uint16 {
	raw := unsafe.Slice((*uint16)(t.C()), 22)
	n := 0
	for n < len(raw) && raw[n] != 0 {
		n++
	}
	return raw[:n]
}
func teamRuntimeSendID(to int, op byte, id uint32) {
	b := []byte{196, op, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[2:], id)
	gameplayReportSend(to, b, true, 1)
}
func teamRuntimeDescribe(t *server.Team, to int) {
	if t == nil {
		return
	}
	text := teamRuntimeNameUnits(t)
	b := make([]byte, 18+2*len(text))
	b[0] = 196
	binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
	binary.LittleEndian.PutUint32(b[6:], t.Field_60)
	binary.LittleEndian.PutUint32(b[10:], uint32(t.Lessons))
	if noxflags.HasGame(512) {
		b[14] = 1
	}
	b[15] = byte(len(text))
	b[16] = byte(t.ColorInd)
	b[17] = byte(t.Field_68)
	for i, c := range text {
		binary.LittleEndian.PutUint16(b[18+2*i:], c)
	}
	gameplayReportSend(to, b, true, 1)
}
func teamRuntimeAnnounce(t *server.Team) {
	if t != nil {
		teamUITeamAdd(t.Name())
		teamRuntimeDescribe(t, 159)
	}
}
func teamRuntimeRename(t *server.Team, name *uint16) {
	if t == nil {
		return
	}
	teamUITeamRename(t, alloc.GoString16(name))
	t.Field_68 = 0
	if noxflags.HasGame(1) {
		var b [46]byte
		b[0], b[1] = 196, 4
		binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
		alloc.StrCopy16P(unsafe.Slice((*uint16)(unsafe.Pointer(&b[6])), 20), name)
		gameplayReportSend(159, b[:], true, 1)
	}
	alloc.StrCopy16P(unsafe.Slice((*uint16)(t.C()), 22), name)
}
func teamRuntimeMemberReport(m *server.ObjectTeam, to, code int) {
	if m == nil {
		return
	}
	u := objectLookupByNetCode(uint32(code))
	if u == nil {
		return
	}
	b := []byte{196, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[2:], uint32(m.ID))
	binary.LittleEndian.PutUint16(b[6:], uint16(code))
	binary.LittleEndian.PutUint16(b[8:], u.TypeInd)
	gameplayReportSend(to, b, true, 1)
}
func teamRuntimeRequest(t *server.Team, m *server.ObjectTeam, code int16, op byte) int8 {
	result := int8(memmap.PtrWord(t.C()))
	if t == nil || m == nil {
		return result
	}
	result = int8(t.ID())
	if m.ID == t.ID() {
		return result
	}
	b := []byte{196, op, 0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
	binary.LittleEndian.PutUint16(b[6:], uint16(code))
	return int8(reliableClientSend(31, b, nil, 1))
}

// The legacy entry point sends through the reliable queue directly. The public
// Server.TeamChangeLessons API uses a replaceable send hook instead.
func teamRuntimeLessons(t *server.Team, value int) {
	if t == nil {
		return
	}
	t.Lessons = value
	if noxflags.HasGame(1) {
		b := []byte{196, 8, 0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b[2:], uint32(t.ID()))
		binary.LittleEndian.PutUint32(b[6:], uint32(value))
		gameplayReportSend(159, b, true, 1)
	}
}
