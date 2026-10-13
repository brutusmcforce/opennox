package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func scoreboardDrawQuest() byte {
	scoreboardClearRows()
	scoreboardHeadings(0, 0)
	row := 1
	for i := 0; i < int(*scoreboardCount(1090117)); i++ {
		r := &scoreboardPlayers()[i]
		p := GetServer().S().Players.ByID(int(r.NetCode))
		if p == nil || p.Field4792 == 0 {
			continue
		}
		field := func(off uintptr) byte { return *(*byte)(unsafe.Add(unsafe.Pointer(p), layoutServerPlayer.Off(off))) } // off is a 386 offset
		color := byte(4)
		if r.Flags&1 != 0 {
			color = 9
		} else if r.Team != -1 {
			color = scoreboardTeamColor(scoreboardTeamIndex(uint32(r.Team)))
		}
		if r.NetCode == *scoreboardData.localCode {
			*scoreboardData.localRow = uint32(int32(int16(uiListData(scoreboardColumn(0, 0)).Field_11_1)))
			*memmap.PtrUint32(0x5D4594, 1088996) = 0
		}
		scoreboardCellName(0, 0, color, &r.Name[0])
		scoreboardCell(0, 2, color, fmt.Sprintf(scoreboardLiteral(147836), field(4816)))
		health := int(field(2282))
		if r.NetCode == *scoreboardData.localCode {
			maximum := float32(sub_470CD0())
			health = int(int32(float32(float64(sub_470CC0()) / float64(maximum) * 100)))
		}
		healthColor := byte(4)
		if health <= 25 {
			healthColor = 6
		} else if health <= 50 {
			healthColor = 15
		}
		scoreboardCell(0, 3, healthColor, fmt.Sprintf(scoreboardLiteral(147844), health))
		scoreboardCellName(0, 4, color, scoreboardClassNames[r.Class])
		var statusColor byte
		status := scoreboardStatus(int(r.Objective), &statusColor)
		scoreboardInsertSafe(scoreboardColumn(0, 1), statusColor, status)
		if field(4824) != 0 {
			scoreboardDrawKey(row, 0, scoreboardWhite())
		}
		if field(4825) != 0 {
			dx := 0
			if field(4824) == 1 {
				dx = 15
			}
			scoreboardDrawKey(row, dx, scoreboardYellow())
		}
		row++
	}
	*scoreboardData.dirty = 1
	return *scoreboardCount(1090117)
}
func scoreboardDrawKey(row, dx int, color uint32) {
	w := scoreboardColumn(0, 1)
	pos := w.GlobalPos()
	height := int(int16(uiListData(w).Line_height))
	x, y := pos.X+5+dx, pos.Y+height/2+height*row-1
	nox_client_drawSetColor_434460(int(color))
	nox_video_drawCircleColored_4C3270(x, y, 2, int(color))
	nox_client_drawAddPoint_49F500(x+2, y)
	nox_client_drawAddPoint_49F500(x+9, y)
	nox_client_drawLineFromPoints_49E4B0()
	nox_client_drawAddPoint_49F500(x+9, y)
	nox_client_drawAddPoint_49F500(x+9, y+3)
	nox_client_drawLineFromPoints_49E4B0()
	nox_client_drawAddPoint_49F500(x+7, y)
	nox_client_drawAddPoint_49F500(x+7, y+2)
	nox_client_drawLineFromPoints_49E4B0()
}
