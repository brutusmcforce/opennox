package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func scoreboardEmptyEvent(*gui.Window, gui.WindowEvent) gui.WindowEventResp { return nil }
func scoreboardConstruct() *gui.Window {
	*scoreboardData.requested = 6
	scoreboardLoadClasses()
	c := GetClient().Cli()
	r := GetClient().R2()
	fh := r.FontHeight(nil)
	widths := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1084036), 5)
	measure := func(id string) int { return r.GetStringSizeWrapped(nil, scoreboardText(id), scoreboardScreenWidth()).X }
	widths[0] = 80
	widths[1] = uint32(max(18, measure("Flag")) + 14)
	widths[3] = uint32(max(measure("Score"), measure("HealthHeading")) + 7)
	widths[4] = uint32(max(measure("Ping"), r.GetStringSizeWrapped(nil, scoreboardLiteral(145972), scoreboardScreenWidth()).X) + 7)
	widths[2] = uint32(max(measure("Class"), measure("Warrior"), measure("Wizard"), measure("Conjurer"), measure("LivesHeading")) + 7)
	total := 0
	for _, w := range widths {
		total += int(w)
	}
	height := 439 - fh
	*scoreboardData.width = uint32(total)
	*scoreboardData.height = uint32(height)
	parent := c.GUI.NewWindowRaw(nil, 1560, 0, fh+40, 1, 1, nil)
	*scoreboardData.parent = parent
	parent.SetAllFuncs(scoreboardEmptyEvent, scoreboardDraw, nil)
	pd := parent.DrawData()
	pd.BgColorVal = 0x80000000
	pd.EnColorVal = 0x80000000
	pd.HlColorVal = 0x80000000
	pd.DisColorVal = 0x80000000
	pd.SelColorVal = 0x80000000
	draw := gui.WindowData{Style: 32, BgColorVal: 0x80000000, EnColorVal: 0x80000000, HlColorVal: 0x80000000, DisColorVal: 0x80000000, SelColorVal: 0x80000000, TextColorVal: memmap.Uint32(0x85B3FC, 940)}
	draw.SetText(alloc.GoString16(memmap.PtrUint16(0x5D4594, 1090136)))
	data := gui.ScrollListBoxData{Count: 64, Line_height: uint16(fh + 1), Field_2: 1}
	for side := 0; side < 2; side++ {
		group := uiListNew(parent, 1088, side*total, 3*fh+1, total, height-2*(fh+1), &draw, &data)
		ScoreboardGroups[side] = group
		x := 0
		for column, width := range widths {
			w := uiListNew(group, 1088, x, 2*fh, int(width), height-2*(fh+1), &draw, &data)
			scoreboardColumns[column][side] = w
			x += int(width)
		}
		group.SetFunc94(scoreboardEmptyEvent)
		for col := 0; col < 5; col++ {
			scoreboardColumn(side, col).SetParent(group)
		}
	}
	draw.Style = 2048
	makeText := func(id string, y int, color uint32) *gui.Window {
		draw.TextColorVal = color
		return c.GUI.NewStaticTextRaw(parent, 1088, 0, y, total, fh+1, &draw, &gui.StaticTextData{Text: alloc.InternCString16(scoreboardText(id))})
	}
	*scoreboardData.rank = makeText("yourrank", fh, scoreboardYellow())
	*scoreboardData.limit = makeText("WindowDir:Empty", 2*fh, scoreboardWhite())
	*scoreboardData.time = makeText("WindowDir:Empty", 3*fh, memmap.Uint32(0x85B3FC, 940))
	scoreboardTitle = makeText("TeamPlayerRank", 0, scoreboardTitleColor())
	*scoreboardData.dirty = 1
	return parent
}
