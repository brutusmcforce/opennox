package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unicode/utf16"
	"unsafe"
)

func scoreboardLiteral(off uintptr) string { return alloc.GoString16(memmap.PtrUint16(0x587000, off)) }
func scoreboardWrite(off uintptr, text string) *uint16 {
	p := memmap.PtrUint16(0x5D4594, off)
	units := utf16.Encode([]rune(text))
	dst := unsafe.Slice(p, len(units)+1)
	copy(dst, units)
	dst[len(units)] = 0
	return p
}
func scoreboardInsert(w *gui.Window, color byte, p *uint16) int {
	w.Func94(gui.AsWindowEvent(16397, uintptr(unsafe.Pointer(p)), uintptr(color)))
	return 1
}
func scoreboardInsertSafe(w *gui.Window, color byte, p *uint16) int {
	if p == nil {
		p = alloc.InternCString16(scoreboardText("InternalError"))
	}
	if p == nil {
		return 0
	}
	return scoreboardInsert(w, color, p)
}
func scoreboardClearRows() int {
	var out gui.WindowEventResp
	for side := 0; side < 2; side++ {
		for column := 0; column < 5; column++ {
			out = scoreboardColumn(side, column).Func94(gui.AsWindowEvent(16399, 1, 0))
		}
	}
	return gui.EventRespInt(out)
}
func scoreboardCell(side, column int, color byte, text string) {
	scoreboardInsert(scoreboardColumn(side, column), color, scoreboardWrite(1089000, text))
}
func scoreboardCellName(side, column int, color byte, p *uint16) {
	dst := memmap.PtrUint16(0x5D4594, 1089000)
	scoreboardCopyName(dst, p)
	scoreboardInsert(scoreboardColumn(side, column), color, dst)
}
func scoreboardHeadings(side, padding int) int {
	for row := 0; row < padding; row++ {
		for col := 0; col < 5; col++ {
			scoreboardInsertSafe(scoreboardColumn(side, col), 9, memmap.PtrUint16(0x587000, 147124+uintptr(4*col)))
		}
	}
	values := [5]string{scoreboardText("player"), scoreboardLiteral(147188), scoreboardText("class"), scoreboardText("score"), scoreboardText("ping")}
	switch *scoreboardData.mode {
	case 1:
		values[2] = scoreboardText("LivesHeading")
		values[3] = scoreboardText("HealthHeading")
		values[4] = scoreboardText("class")
	case 5:
		values[3] = scoreboardText("rank")
	}
	out := 0
	for col, text := range values {
		out = scoreboardInsertSafe(scoreboardColumn(side, col), 9, alloc.InternCString16(text))
	}
	return out
}
func scoreboardStatus(kind int, color *byte) *uint16 {
	id := ""
	*color = 4
	switch kind {
	case 1:
		id = "King"
	case 2:
		id = "Flag"
		*color = 7
	case 3:
		id = "Flag"
		*color = 13
	case 4:
		id = "Ball"
	default:
		return memmap.PtrUint16(0x587000, 147724)
	}
	return scoreboardWrite(1090024, "<"+scoreboardText(id))
}
func scoreboardSetText(w *gui.Window, p *uint16) int {
	return gui.EventRespInt(w.Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(p)), 0)))
}
func scoreboardHide(w *gui.Window) int {
	if scoreboardHidden(w) {
		return 1
	}
	if w == nil {
		return -2
	}
	w.Hide()
	return 0
}
func scoreboardTimeHeading() byte {
	if Sub_40A220() != 0 && (!noxflags.HasGame(noxflags.GameHost) || sub_40A300() != 0 || Sub_40A180(noxflags.GetGame()) != 0) {
		w := *scoreboardData.time
		if scoreboardHidden(w) {
			w.Show()
		}
		ms := int32(Sub_40A230())
		return byte(scoreboardSetText(w, scoreboardWrite(1084068, fmt.Sprintf(scoreboardText("TimeRemaining"), ms/60000, ms%60000/1000))))
	}
	return byte(scoreboardHide(*scoreboardData.time))
}
func scoreboardLessonLimit() int {
	if noxflags.HasGame(noxflags.GameHost) {
		return int(uint16(Nox_xxx_servGamedataGet_40A020(uint16(noxflags.GetGame()))))
	}
	return int(memmap.Uint16(0x5D4594, 371380+54))
}
func scoreboardLessonHeading() int {
	w := *scoreboardData.limit
	if noxflags.HasGame(4224) {
		return scoreboardHide(w)
	}
	if scoreboardHidden(w) {
		w.Show()
	}
	return scoreboardSetText(w, scoreboardWrite(1083972, fmt.Sprintf(scoreboardText("LessonLimit"), scoreboardLessonLimit())))
}
func scoreboardHighlight() {
	w := scoreboardColumn(int(memmap.Uint32(0x5D4594, 1088996)), 0)
	p := w.GlobalPos()
	h := int(int16(uiListData(w).Line_height))
	p.Y += int(int32(*scoreboardData.localRow))*h + h/2
	nox_client_drawSetColor_434460(int(scoreboardYellow()))
	nox_xxx_drawPointMB_499B70(p.X+1, p.Y, 3)
}
