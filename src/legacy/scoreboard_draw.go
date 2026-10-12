package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func scoreboardDraw(w *gui.Window, draw *gui.WindowData) int {
	if !noxflags.HasGame(8) && *scoreboardData.requested != 6 {
		*scoreboardData.mode = *scoreboardData.requested - 1
		if int32(*scoreboardData.mode) < 0 {
			*scoreboardData.mode = 5
		}
		Sub_4703F0()
		*scoreboardData.requested = 6
		if *scoreboardData.mode == 0 {
			return 1
		}
	}
	pos := w.GlobalPos()
	limit := scoreboardLessonLimit()
	if byte(w.Flags)&0x80 == 0 {
		if draw.BgColorVal != 0x80000000 {
			GetClient().R2().DrawRectFilledAlpha(pos.X, pos.Y, w.SizeVal.X, w.SizeVal.Y)
		}
	} else {
		GetClient().R2().DrawImageAt(GetClient().R2().GetBag().AsImage(draw.BgImageHnd), pos)
	}
	s := GetServer().S()
	if *scoreboardData.dirty != 0 || s.Frame() > memmap.Uint32(0x5D4594, 1090124)+s.TickRate() {
		GetClient().Cli().GUI.ValYYY = 1
		*memmap.PtrUint32(0x5D4594, 1090124) = s.Frame()
		*scoreboardData.localRow = 0xffffffff
		scoreboardClearRows()
		scoreboardCollect()
		*scoreboardData.dirty = 0
		var team *server.Team
		if ot := objectRenderTeam(int(*scoreboardData.localCode)); ot != nil {
			team = s.Teams.ByID(ot.ID)
		}
		showRank := false
		if p := Get_dword_8531A0_2576(); p != nil && !scoreboardObserver(p) {
			showRank = true
		}
		mode := *scoreboardData.mode
		if *scoreboardCount(1090116) != 0 && (mode == 2 || mode == 3) {
			scoreboardDrawTeams()
		}
		padding := int(int16(uiListData(scoreboardColumn(0, 0)).Field_11_1))
		if *scoreboardCount(1090117) != 0 && (mode == 2 || mode == 4 || mode == 5) {
			scoreboardHeadings(0, 0)
			count := int(*scoreboardCount(1090117))
			if mode == 4 {
				count = min(count, 3)
			}
			for i := 0; i < count; i++ {
				side := i >> 4
				if i == 16 {
					scoreboardHeadings(1, padding)
				}
				scoreboardDrawPlayer(i, side, limit)
			}
		} else if mode == 1 {
			scoreboardDrawQuest()
		}
		title := (*uint16)(unsafe.Pointer(draw))
		canRank := !noxflags.HasGame(noxflags.GameHost) || !noxflags.HasEngine(noxflags.EngineNoRendering)
		switch mode {
		case 1, 2, 4, 5:
			id := "TeamPlayerRank"
			switch mode {
			case 1:
				id = "Noxworld.c:Quest"
			case 4:
				id = "Top3"
			case 5:
				id = "WolRank"
			}
			title = alloc.InternCString16(scoreboardText(id))
			if canRank {
				if mode == 1 {
					scoreboardWrite(1086692, fmt.Sprintf("%s %d", scoreboardText("Noxworld.c:Stage"), briefingStage()))
				} else {
					scoreboardWrite(1086692, fmt.Sprintf("%s %d / %d", scoreboardText("yourrank"), scoreboardLocalRank(), *scoreboardCount(1090118)))
				}
			}
		case 3:
			rank := byte(0)
			if team != nil {
				showRank = true
				rank = scoreboardTeamRank(uint32(team.Lessons))
			}
			title = alloc.InternCString16(scoreboardText("Teams"))
			if canRank {
				scoreboardWrite(1086692, fmt.Sprintf("%s %d / %d", scoreboardText("yourteamrank"), rank, *scoreboardCount(1090116)))
			}
		}
		rw := *scoreboardData.rank
		if showRank {
			if scoreboardHidden(rw) {
				rw.Show()
			}
			scoreboardSetText(rw, memmap.PtrUint16(0x5D4594, 1086692))
		} else if !scoreboardHidden(rw) {
			rw.Hide()
		}
		scoreboardSetText(scoreboardTitle, title)
		scoreboardTimeHeading()
		scoreboardLessonHeading()
	}
	if int32(*scoreboardData.localRow) >= 0 {
		scoreboardHighlight()
	}
	return 1
}
func scoreboardDrawTeams() {
	for col, p := range []*uint16{alloc.InternCString16(scoreboardText("team")), memmap.PtrUint16(0x587000, 146512), memmap.PtrUint16(0x587000, 146516), alloc.InternCString16(scoreboardText("score")), memmap.PtrUint16(0x587000, 146564)} {
		scoreboardInsertSafe(scoreboardColumn(0, col), 9, p)
	}
	color := byte(0)
	for i := 0; i < int(*scoreboardCount(1090116)); i++ {
		r := &scoreboardTeams()[i]
		color = scoreboardTeamColor(byte(i))
		scoreboardCellName(0, 0, color, &r.Name[0])
		scoreboardCell(0, 1, color, scoreboardLiteral(146576))
		scoreboardCell(0, 2, color, scoreboardLiteral(146580))
		scoreboardCell(0, 3, color, fmt.Sprintf(scoreboardLiteral(146584), r.Score))
		scoreboardCell(0, 4, color, scoreboardLiteral(146592))
	}
	for col := 0; col < 5; col++ {
		scoreboardCell(0, col, color, scoreboardLiteral(146596+uintptr(4*col)))
	}
}
func scoreboardDrawPlayer(index, side, limit int) {
	r := &scoreboardPlayers()[index]
	eliminated := noxflags.HasGame(noxflags.GameModeElimination) && limit > 0 && int32(r.Score) >= int32(limit)
	color := byte(9)
	if r.Flags&1 == 0 || r.Flags&0x20 != 0 || eliminated {
		if r.Team == -1 {
			color = 3
			if r.Flags&0x20 != 0 || eliminated {
				color = 2
			}
		} else {
			color = scoreboardTeamColor(scoreboardTeamIndex(uint32(r.Team)))
			if r.Flags&0x20 != 0 || eliminated {
				color -= 2
			}
		}
	}
	if r.NetCode == *scoreboardData.localCode {
		*scoreboardData.localRow = uint32(int32(int16(uiListData(scoreboardColumn(side, 0)).Field_11_1)))
		*memmap.PtrUint32(0x5D4594, 1088996) = uint32(side)
	}
	scoreboardCellName(side, 0, color, &r.Name[0])
	scoreboardCellName(side, 2, color, scoreboardClassNames[r.Class])
	if *scoreboardData.mode != 5 || r.Score > 0 {
		scoreboardCell(side, 3, color, fmt.Sprintf(scoreboardLiteral(146632), int32(r.Score)))
	} else {
		scoreboardCell(side, 3, color, scoreboardLiteral(146640))
	}
	scoreboardCell(side, 4, color, fmt.Sprintf(scoreboardLiteral(146648), int32(r.Ping)))
	if *scoreboardData.mode == 5 {
		if p := GetServer().S().Players.ByID(int(r.NetCode)); p != nil {
			scoreboardCell(side, 1, color, "("+alloc.GoString((*byte)(unsafe.Pointer(&p.Field2096Buf[0])))+")")
		}
	} else {
		var statusColor byte
		p := scoreboardStatus(int(r.Objective), &statusColor)
		scoreboardInsertSafe(scoreboardColumn(side, 1), statusColor, p)
	}
}
