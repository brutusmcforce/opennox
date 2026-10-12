package legacy

import (
	"image"
	"strconv"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/player"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"golang.org/x/image/font"
)

func uiInventoryStats(pos image.Point) {
	r := GetClient().R2()
	normal := r.GetFonts().AsFont(nil)
	small := uiInventorySmallFont()
	height := r.FontHeight(normal)
	smallHeight := r.FontHeight(small)
	baseline := int(float32(float64(height-smallHeight)*0.5 + 0.5))
	p := uiMeterPlayer()
	if p == nil {
		return
	}
	class := *(*byte)(unsafe.Add(p, 2251))
	stats := *GetServer().S().Players.ClassStats(player.Class(class))
	warrior := *GetServer().S().Players.ClassStats(0)
	attr := func(off uintptr) uint32 { return *(*uint32)(unsafe.Add(p, off)) }
	textColor := func(c uint32) { r.Data().SetTextColor(noxcolor.RGBA5551(c)) }
	text := func(face font.Face, s string, x, y, width int) {
		r.DrawStringWrapped(face, s, image.Rect(x, y, x+width, y))
	}
	fill := func(color uint32, x, y, width, h int) {
		uiMeterSetColor(color)
		nox_client_drawRectFilledOpaque_49CE30(x, y, width, h)
	}
	white := uint32(nox_color_white_2523948)
	x, y := pos.X+13, pos.Y+15
	textColor(white)
	fill(uint32(nox_color_black_2650656), pos.X+11, y, 200, 200)
	y += 2*height + 3
	level := int(*(*int8)(unsafe.Add(p, 3684)))
	text(normal, uiInventoryFormatInts("StatsLevel", level, 0), x, y, 200)
	y += height + 1
	if noxflags.HasGame(2048) {
		xp := int(int32(int64(GetServer().S().Balance.FloatInd("XPTable", level+1))))
		text(normal, uiInventoryFormatInts("StatsEXP", int(memmap.Int32(0x5D4594, 1062544)), xp), x, y, 200)
	}
	y += 2*height + 2
	text(normal, uiInventoryText("StatsHealth"), x, y, 200)
	fill(uint32(nox_color_violet_2598268), x+60, y, 90, height)
	width := int(float32(float64(int32(90*attr(2247))) / float64(stats.Health)))
	fill(uint32(nox_color_red_2589776), x+60, y, width, height)
	currentHealth := int32(uiMeters()[0].Current)
	width = int(float32(float64(90*currentHealth) / float64(stats.Health)))
	fill(memmap.Uint32(0x85B3FC, 940), x+60, y, width, height)
	s := uiInventoryFormatInts("MinMaxFormat", int(int32(attr(2247))), int(stats.Health))
	width = r.GetStringSizeWrapped(small, s, 0).X
	text(small, s, x-width+193, y+baseline, 200)
	text(small, strconv.Itoa(int(currentHealth)), x+45, y+baseline, 200)
	y += height + 1
	if class != 0 {
		fill(memmap.Uint32(0x85B3FC, 944), x+60, y, 90, height)
		width = int(float32(float64(int32(90*attr(2243))) / float64(stats.Mana)))
		text(normal, uiInventoryText("StatsMana"), x, y, 200)
		fill(uint32(nox_color_blue_2650684), x+60, y, width, height)
		mana := int32(uiMeters()[1].Current)
		width = int(float32(float64(90*mana) / float64(stats.Mana)))
		fill(uint32(nox_color_cyan_2649820), x+60, y, width, height)
		s = uiInventoryFormatInts("MinMaxFormat", int(int32(attr(2243))), int(stats.Mana))
		width = r.GetStringSizeWrapped(small, s, 0).X
		text(small, s, x-width+193, y+baseline, 200)
		text(small, strconv.Itoa(int(mana)), x+45, y+baseline, 200)
		y += height + 1
	}
	fill(memmap.Uint32(0x85B3FC, 956), x+60, y, 90, height)
	width = int(float32(float64(int32(90*attr(2239))) / float64(stats.Strength)))
	text(normal, uiInventoryText("StatsStrength"), x, y, 200)
	fill(memmap.Uint32(0x5D4594, 2597996), x+60, y, width, height)
	s = uiInventoryFormatInts("MinMaxFormat", int(int32(attr(2239))), int(stats.Strength))
	width = r.GetStringSizeWrapped(small, s, 0).X
	text(small, s, x-width+193, y+baseline, 200)
	text(small, strconv.Itoa(int(int32(attr(2239)))), x+45, y+baseline, 200)
	y += height + 1
	fill(uint32(nox_color_orange_2614256), x+60, y, 90, height)
	width = int(float32(float64(int32(90*attr(2235)))/float64(stats.Speed) + 0.5))
	text(normal, uiInventoryText("StatsSpeed"), x, y, 200)
	fill(uint32(nox_color_yellow_2589772), x+60, y, width, height)
	textColor(white)
	speed := float32(float64(memmap.Float32(0x5D4594, 1063100)) / (float64(warrior.Speed) * 0.000001))
	if memmap.Uint8(0x5D4594, 1062541)&2 != 0 {
		speed = float32((float64(width)+float64(speed))*1.25 - float64(width))
	}
	if memmap.Uint8(0x5D4594, 1062540)&16 != 0 {
		speed = float32((float64(width)+float64(speed))*0.5 - float64(width))
	}
	if speed < 0 {
		fill(memmap.Uint32(0x85B3FC, 944), x+60+width+int(speed), y, int(-speed), height)
	} else if speed > 0 {
		extra := int(speed)
		if width+extra > 90 {
			extra = 90 - width
		}
		fill(uint32(nox_color_yellow_2589772), x+60+width, y, extra, height)
		textColor(uint32(nox_color_blue_2650684))
	}
	adjustment := float32(float64(speed) * 100.0 * 0.011111111)
	maxPercent := int(float32(float64(stats.Speed) * 100.0 / float64(warrior.Speed)))
	currentPercent := int(float32(float64(int32(attr(2235)))*100.0/float64(warrior.Speed) + float64(adjustment) + 0.5))
	s = uiInventoryFormatInts("MinMaxFormat", currentPercent, maxPercent)
	width = r.GetStringSizeWrapped(small, s, 0).X
	text(small, s, x-width+193, y+baseline, 200)
	textColor(white)
	text(small, strconv.Itoa(currentPercent), x+45, y+baseline, 200)
	textColor(white)
	y += 2*height + 2
	lang := GetServer().S().Strings().Lang()
	if lang == 6 || lang == 8 {
		x += 39
	}
	textColor(white)
	s = uiInventoryText("StatsArmorLabel")
	armorWidth := r.GetStringSizeWrapped(normal, s, 0).X
	text(normal, s, x, y, 0)
	armor := int(float32(float64(memmap.Float32(0x5D4594, 1062548))*1000.0 + 0.5))
	text(normal, uiInventoryFormatInts("MinMaxFormat", armor, 1000), x+armorWidth+5, y, 0)
	y += height + 1
	weight := 0
	for row := 0; row < 21; row++ {
		for col := 0; col < 4; col++ {
			cell := &uiInventoryGrid()[col*21+row]
			if cell.Count != 0 {
				weight += int(cell.Count) * int(*(*byte)(unsafe.Add(cell.Drawable.C(), int(layoutClientDrawable.Off(298)) /* Weight */)))
			}
		}
	}
	textColor(white)
	s = uiInventoryText("DollWeight")
	width = r.GetStringSizeWrapped(normal, s, 0).X
	text(normal, uiInventoryText("DollWeight"), x+armorWidth-width, y, 0)
	capacity := int(*(*uint16)(unsafe.Add(p, 3652)))
	if weight > capacity {
		white = memmap.Uint32(0x85B3FC, 940)
	}
	textColor(white)
	text(normal, uiInventoryFormatInts("MinMaxFormat", weight, capacity), x+armorWidth+5, y, 0)
}
