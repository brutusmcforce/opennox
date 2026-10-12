package legacy

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"strings"
	"unsafe"
)

var teamUICTFRoot, teamUIBallVisible, teamUIPlayersRoot, teamUIJoinControl, teamUIRenameControl uint32

func teamUIWord(off uintptr) *uint32 {
	switch off {
	case 1045604:
		return &teamUICTFRoot
	case 1045636:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1045636))
	case 1045640:
		return &teamUIBallVisible
	case 1045684:
		return &teamUIPlayersRoot
	case 1045688:
		return &teamUIJoinControl
	case 1045692:
		return &teamUIRenameControl
	default:
		return memmap.PtrUint32(0x5D4594, off)
	}
}
func teamUIWindow(off uintptr) *gui.Window {
	return (*gui.Window)(memmap.WordPtr(*teamUIWord(off)))
}
func teamUISettings() *byte {
	return memmap.PtrUint8(0x5D4594, 371380+58*uintptr(*memmap.PtrUint32(0x5D4594, 371688)))
}
func teamUIEvent(w *gui.Window, code int, a, b uintptr) int {
	return gui.EventRespInt(w.Func94(gui.AsWindowEvent(code, a, b)))
}
func teamUIText(key string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(key), "playrlst.c")
}
func teamUIASCIIName(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r - 'A' + 'a'
		}
		return r
	}, s)
}

// Keep the existing row layout while legacy callbacks and fixture snapshots share
// these address-stable intrusive lists. All row algorithms and ownership are Go.
type teamUIRow struct {
	list    legacyListNode
	name    [24]uint16
	code    uint32
	palette byte
	_       [3]byte
	color   uint32
}






func teamUIHead(teams bool) *legacyListNode {
	if teams {
		return &teamUITeamList
	}
	return &teamUIPlayerList
}
func teamUIFirst(teams bool) *teamUIRow {
	return (*teamUIRow)(unsafe.Pointer(listNext(teamUIHead(teams))))
}
func teamUINext(row *teamUIRow) *teamUIRow { return (*teamUIRow)(unsafe.Pointer(listNext(&row.list))) }
func teamUIFreeRows(teams bool) {
	for row := teamUIFirst(teams); row != nil; {
		next := teamUINext(row)
		listRemove(&row.list)
		alloc.Free(row)
		row = next
	}
}
func teamUINewRow(teams bool, name string, code uint32) *teamUIRow {
	row, _ := alloc.New(teamUIRow{})
	alloc.StrCopy16(row.name[:], name)
	row.code = code
	listInit(&row.list)
	listAppend(teamUIHead(teams), &row.list)
	return row
}
