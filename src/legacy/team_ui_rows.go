package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"strings"
	"unsafe"
)

func teamUIPlayerFind(code int, remove bool) int {
	i := 0
	for row := teamUIFirst(false); row != nil; row = teamUINext(row) {
		if row.code == uint32(code) {
			if remove {
				listRemove(&row.list)
				alloc.Free(row)
			}
			return i
		}
		i++
	}
	return -1
}
func teamUITeamFind(name string, remove bool) int {
	i := 0
	key := teamUIASCIIName(name)
	for row := teamUIFirst(true); row != nil; row = teamUINext(row) {
		if teamUIASCIIName(alloc.GoString16S(row.name[:])) == key {
			if remove {
				listRemove(&row.list)
				alloc.Free(row)
			}
			return i
		}
		i++
	}
	return -1
}
func teamUIPlayerAdd(code int, name string) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	teamUINewRow(false, name, uint32(code))
	return teamUIEvent(w.ChildByID(10501), 16397, uintptr(unsafe.Pointer(alloc.InternCString16(name))), 3)
}
func teamUIPlayerRemove(code int) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	i := teamUIPlayerFind(code, true)
	if i < 0 {
		return i
	}
	return teamUIEvent(w.ChildByID(10501), 16398, uintptr(i), 0)
}
func teamUITeamColor(t *server.Team) byte {
	return *memmap.PtrUint8(0x587000, 128968+8*uintptr(byte(t.ColorInd)%10))
}
func teamUITeamLabel(t *server.Team, name string) string {
	if (noxflags.HasGame(96) || *(*byte)(unsafe.Add(unsafe.Pointer(teamUISettings()), 52))&96 != 0) && byte(t.ID()) < 3 {
		key := "BlueFlag"
		if t.ID() == 1 {
			key = "RedFlag"
		}
		return name + teamUIText(key)
	}
	return name
}
func teamUITeamAdd(name string) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	t := teamRuntimeFind(alloc.InternCString16(name))
	if t == nil {
		return 0
	}
	row := teamUINewRow(true, name, uint32(t.ID()))
	row.color = uint32(t.ColorInd)
	row.palette = teamUITeamColor(t)
	text := teamUITeamLabel(t, name)
	return teamUIEvent(w.ChildByID(10502), 16397, uintptr(unsafe.Pointer(alloc.InternCString16(text))), uintptr(row.palette))
}
func teamUIPlayerTeam(code, id int) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	i := teamUIPlayerFind(code, false)
	if i < 0 {
		return i
	}
	palette := byte(3)
	if id != 0 {
		for row := teamUIFirst(true); row != nil; row = teamUINext(row) {
			if row.code == uint32(id) {
				palette = row.palette
			}
		}
	}
	d := (*gui.ScrollListBoxData)(w.ChildByID(10501).WidgetData)
	items := unsafe.Slice(d.Items, int(d.Count))
	items[i].Field_129 = **(**uint32)(memmap.PtrOff(0x85B3FC, 132+4*uintptr(palette)))
	return int(uintptr(unsafe.Pointer(d.Items)))
}
func teamUIPlayersRefresh() int {
	w := teamUIWindow(1045684)
	teamUIFreeRows(false)
	teamUIFreeRows(true)
	players, teams := w.ChildByID(10501), w.ChildByID(10502)
	teamUIEvent(players, 16399, 0, 0)
	teamUIEvent(teams, 16399, 0, 0)
	uiWindowEnable(teamUIWindow(1045688), 0)
	uiWindowEnable(teamUIWindow(1045692), 0)
	s := GetServer().S()
	for t := s.Teams.First(); t != nil; t = s.Teams.Next(t) {
		teamUITeamAdd(t.Name())
	}
	for p := s.Players.First(); p != nil; p = s.Players.Next(p) {
		if p.PlayerInd == 31 && noxflags.HasEngine(noxflags.EngineNoRendering) {
			continue
		}
		code := int(p.NetCodeVal)
		teamUIPlayerAdd(code, p.Name())
		if member := teamRuntimeObject(code); member != nil && teamRuntimeContains(member, member.ID) {
			teamUIPlayerTeam(code, int(member.ID))
		}
	}
	enabled := 1
	if noxflags.HasGame(4096) {
		enabled = 0
	}
	uiWindowEnable(players, enabled)
	return uiWindowEnable(teams, enabled)
}
func teamUIPlayersRefreshOpen() int {
	if teamUIWindow(1045684) == nil {
		return 0
	}
	return teamUIPlayersRefresh()
}
func teamUITeamRemove(name string) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	list := w.ChildByID(10502)
	if i := teamUITeamFind(name, true); i >= 0 {
		teamUIEvent(list, 16404, 0, 0)
		teamUIEvent(list, 16398, uintptr(i), 0)
	}
	return teamUIPlayersRefresh()
}
func teamUITeamClear() int {
	if teamUIWindow(1045684) == nil {
		return 0
	}
	teamUIFreeRows(true)
	uiWindowEnable(teamUIWindow(1045688), 0)
	uiWindowEnable(teamUIWindow(1045692), 0)
	return teamUIPlayersRefresh()
}
func teamUITeamRename(t *server.Team, name string) int {
	w := teamUIWindow(1045684)
	if w == nil {
		return 0
	}
	list := w.ChildByID(10502)
	selected := teamUIEvent(list, 16404, 0, 0)
	index := 0
	row := teamUIFirst(true)
	for row != nil && row.code != uint32(t.ID()) {
		row = teamUINext(row)
		index++
	}
	if row == nil {
		return 0
	}
	alloc.StrCopy16(row.name[:], name)
	label := teamUITeamLabel(t, name)
	teamUIEvent(list, 16398, uintptr(index), 0)
	teamUIEvent(list, 16402, uintptr(index), 0)
	teamUIEvent(list, 16397, uintptr(unsafe.Pointer(alloc.InternCString16(label))), uintptr(teamUITeamColor(t)))
	return teamUIEvent(list, 16403, uintptr(selected), 0)
}
func teamUISelectedName(index int) string {
	value := teamUIEvent(teamUIWindow(1045684).ChildByID(10502), 16406, uintptr(index), 0)
	name := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(value))))
	name = strings.TrimLeft(name, "\t\n\r")
	if i := strings.IndexAny(name, "\t\n\r"); i >= 0 {
		name = name[:i]
	}
	return name
}
func teamUIPlayersDestroy(destroy bool) uintptr {
	w := teamUIWindow(1045684)
	if destroy && w != nil {
		w.Destroy()
	}
	teamUIFreeRows(false)
	old := uintptr(unsafe.Pointer(teamUIFirst(true)))
	teamUIFreeRows(true)
	*teamUIWord(1045684) = 0
	return old
}
func teamUIRequestsReset() { *teamUIWord(1045696) = 0 }
