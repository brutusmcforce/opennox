package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func teamUIJoin(t *server.Team) {
	if p := Get_dword_8531A0_2576(); p != nil && (p.WeaponEquip&1 != 0 || p.Field3680&1 != 0) {
		return
	}
	host := noxflags.HasGame(1)
	code := ClientPlayerNetCode()
	member := teamRuntimeObject(code)
	if member == nil {
		return
	}
	if teamRuntimeContains(member, member.ID) {
		if !host {
			teamRuntimeRequest(t, member, int16(code), 11)
		} else {
			teamRuntimeSwitch(member, t, code, 1)
			if !noxflags.HasGame(128) {
				if pl := GetServer().S().Players.ByID(code); pl != nil {
					var pos types.Pointf
					controlFindStart(&pos, pl.PlayerUnit)
					Nox_xxx_unitMove_4E7010(pl.PlayerUnit, pos)
				}
			}
		}
	} else if host {
		teamRuntimeJoin(t.ID(), member, 1, code, 1)
	} else {
		teamRuntimeRequest(t, member, int16(code), 10)
	}
	*teamUIWord(1045696)++
}
func teamUIPlayerByName(name string) *server.Player {
	s := GetServer().S()
	name = teamUIASCIIName(name)
	for p := s.Players.First(); p != nil; p = s.Players.Next(p) {
		if teamUIASCIIName(p.Name()) == name {
			return p
		}
	}
	return nil
}
func teamUIPlayersEvent(w *gui.Window, code int, child *gui.Window, value int) int {
	if code != 16391 && code != 16400 {
		return 0
	}
	id := child.ID()
	if code == 16400 && id == 10502 {
		if teamUIEvent(child, 16404, 0, 0) < 0 || noxflags.HasGame(0x8000) || teamUISettingsLocked() {
			uiWindowEnable(teamUIWindow(1045688), 0)
			uiWindowEnable(teamUIWindow(1045692), 0)
		} else {
			if noxflags.HasGame(1) {
				uiWindowEnable(teamUIWindow(1045692), 1)
			}
			enabled := 0
			if noxflags.HasGame(128) || *teamUIWord(1045696) == 0 {
				enabled = 1
			}
			uiWindowEnable(teamUIWindow(1045688), enabled)
		}
		if noxflags.HasGame(1) && noxflags.HasEngine(noxflags.EngineNoRendering) {
			uiWindowEnable(teamUIWindow(1045688), 0)
		}
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	root := teamUIWindow(1045684)
	selectedTeam := func() *server.Team {
		i := teamUIEvent(root.ChildByID(10502), 16404, 0, 0)
		return teamRuntimeFind(alloc.InternCString16(teamUISelectedName(i)))
	}
	switch id {
	case 10509:
		Nox_xxx_dialogMsgBoxCreate_449A10(root, teamUIText("Rename"), teamUIText("NewName"), 163, nil, nil)
	case 10507:
		if t := selectedTeam(); t != nil {
			teamUIJoin(t)
		}
	case 4001:
		if t := selectedTeam(); t != nil {
			name := (*uint16)(unsafe.Pointer(uintptr(Sub_449E60(-88))))
			if teamRuntimeFind(name) == nil {
				teamRuntimeRename(t, name)
			}
		}
	case 10503:
		list := root.ChildByID(10502)
		index := teamUIEvent(list, 16404, 0, 0)
		if index >= 0 {
			if t := teamRuntimeFind(alloc.InternCString16(teamUISelectedName(index))); t != nil {
				players := root.ChildByID(10501)
				selected := (*int32)(unsafe.Pointer(uintptr(teamUIEvent(players, 16404, 0, 0))))
				for cursor := selected; *cursor >= 0; cursor = (*int32)(unsafe.Add(unsafe.Pointer(cursor), 4)) {
					text := teamUIEvent(players, 16406, uintptr(*cursor), 0)
					if p := teamUIPlayerByName(alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(text))))); p != nil && p.Field3680&1 == 0 && p.WeaponEquip&1 == 0 {
						if member := teamRuntimeObject(int(p.NetCodeVal)); member != nil {
							if teamRuntimeContains(member, member.ID) {
								teamRuntimeSwitch(member, t, int(p.NetCodeVal), 1)
							} else {
								teamRuntimeJoin(t.ID(), member, 1, int(p.NetCodeVal), 1)
							}
						}
					}
				}
			}
		}
		teamUIEvent(list, 16403, ^uintptr(0), 0)
		teamUIEvent(root.ChildByID(10501), 16403, ^uintptr(0), 0)
	}
	return 0
}
