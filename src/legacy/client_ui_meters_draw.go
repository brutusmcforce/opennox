package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"math"
	"strconv"
	"unsafe"
)

func uiMeterMainWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(dword_5d4594_1090276)))
}
func uiMeterImage(handle uint32, pos image.Point) {
	nox_client_drawImageAt_47D2C0(noxrender.ImageHandle(unsafe.Pointer(uintptr(handle))), pos.X, pos.Y)
}
func uiMeterSetIcon(w *gui.Window, handle uint32) {
	if w != nil {
		w.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(handle)))
	}
}
func uiMeterSetColor(v uint32) { GetClient().R2().Data().SetColor2(noxcolor.RGBA5551(v)) }
func uiMeterCross(x, y int) int {
	nox_client_drawPixel_49EFA0(x+1, y)
	nox_client_drawRectFilledOpaque_49CE30(x, y+1, 3, 1)
	nox_client_drawPixel_49EFA0(x+1, y+2)
	return 0
}
func uiMeterLabel(w *gui.Window) int {
	r := GetClient().R2()
	text := strconv.FormatInt(int64(int32(uiMeters()[uintptr(w.WidgetData)].Current)), 10)
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	face := r.GetFonts().AsFont(unsafe.Pointer(uintptr(dword_5d4594_1096288)))
	width := r.GetStringSizeWrapped(face, text, 0).X
	pos := uiWindowPosition(w)
	r.DrawString(face, text, image.Pt(pos.X-width/2+8, pos.Y+1))
	return 1
}

func uiMeterMiniBar(w *gui.Window) int {
	index := int(uintptr(w.WidgetData))
	m := &uiMeters()[index]
	if nox_xxx_clientIsObserver_4372E0() != 0 || nox_gameDisableMapDraw_5d4594_2650672 != 0 || noxflags.HasGame(9437184) {
		return 1
	}
	x := int(nox_win_width)/2 + 15
	if index != 0 {
		x += 6
	}
	y := int(nox_win_height)/2 - 48
	height := 0
	if m.Maximum != 0 {
		height = int(48 * int32(m.Current) / int32(m.Maximum))
	}
	uiMeterSetColor(uint32(nox_color_black_2650656))
	nox_client_drawRectFilledOpaque_49CE30(x, y, 2, 48)
	uiMeterSetColor(m.Color)
	nox_client_drawRectFilledOpaque_49CE30(x, y-height+48, 2, height)
	if index != 0 {
		uiMeterSetColor(memmap.Uint32(0x85B3FC, 944))
	} else if dword_5d4594_1096264 != 0 {
		uiMeterSetColor(memmap.Uint32(0x85B3FC, 984))
	} else {
		uiMeterSetColor(uint32(nox_color_violet_2598268))
	}
	nox_client_drawBorderLines_49CC70(x-1, y-1, 4, 50)
	return 1
}

type uiMeterBubble struct{ X, Y, Size, Speed, Active, Color int32 }

func uiMeterTube(w *gui.Window) int {
	index := int(uintptr(w.WidgetData))
	m := &uiMeters()[index]
	if index == 0 && dword_5d4594_1096264 != 0 {
		uiMeterImage(memmap.Uint32(0x5D4594, 1091900), uiWindowPosition(uiMeterMainWindow()))
	}
	pos := uiWindowPosition(w)
	pos.X += 5
	if get_dword_5d4594_3799468() != 0 {
		uiMeterSetColor(uint32(nox_color_black_2650656))
		nox_client_drawRectFilledOpaque_49CE30(pos.X, pos.Y, 15, 125)
	}
	if m.Maximum == 0 {
		nox_client_drawRectFilledAlpha_49CF10(pos.X, pos.Y, 15, 125)
		return 1
	}
	height := int(uint32(125) * m.Current / m.Maximum)
	nox_client_drawRectFilledAlpha_49CF10(pos.X, pos.Y, 15, 125-height)
	uiMeterSetColor(m.Color)
	nox_client_drawEnableAlpha_434560(1)
	nox_client_drawRectFilledOpaque_49CE30(pos.X, pos.Y-height+125, 15, height)
	nox_client_drawEnableAlpha_434560(0)
	nox_client_drawAddPoint_49F500(pos.X, pos.Y-height+125)
	nox_xxx_rasterPointRel_49F570(14, 0)
	nox_client_drawLineFromPoints_49E4B0()
	if index < 2 && nox_client_renderBubbles_80844 == 1 {
		bubbles := unsafe.Slice((*uiMeterBubble)(memmap.PtrOff(0x5D4594, 1093180+uintptr(index)*1536)), 64)
		for i := range bubbles {
			b := &bubbles[i]
			if b.Active == 0 {
				continue
			}
			y := int(b.Y >> 4)
			if y < 125-height {
				b.Active = 0
				continue
			}
			if dword_5d4594_1096264 != 0 {
				uiMeterSetColor(m.Color)
			} else {
				uiMeterSetColor(uint32(b.Color))
			}
			if b.Size <= 2 {
				nox_client_drawRectFilledOpaque_49CE30(pos.X+int(b.X), pos.Y+y, int(b.Size), int(b.Size))
			} else {
				uiMeterCross(pos.X+int(b.X), pos.Y+y)
			}
			b.Y -= b.Speed
		}
		if height > 1 {
			rng := GetServer().S().Rand.Other
			for i := range bubbles {
				b := &bubbles[i]
				if b.Active != 0 {
					continue
				}
				roll := rng.Int(1, 100)
				b.Size = 1
				if roll >= 80 {
					b.Size = 2
					if roll >= 95 {
						b.Size = 3
					}
				}
				b.X = int32(rng.Int(0, 14))
				if b.Size+b.X > 15 {
					b.X = 15 - b.Size
				}
				b.Y = 16 * (125 - b.Size)
				b.Speed = int32(rng.Int(4, 48))
				b.Active = 1
				shade := rng.Int(0, 64)
				if index != 0 {
					b.Color = int32(uiMeterColor(shade, shade, 255))
				} else {
					b.Color = int32(uiMeterColor(255, shade, shade))
				}
			}
		}
	}
	if index == 0 {
		frame := uint32(dword_5d4594_1096256)
		if int32(dword_5d4594_1096260) > 0 {
			frame += uint32(dword_5d4594_1096260)
			dword_5d4594_1096256 = uint32(frame)
			dword_5d4594_1096260--
			if frame>>3 >= 10 {
				frame = 0
				dword_5d4594_1096256 = 0
			}
		}
		uiMeterSetIcon(uiMeterMainWindow(), memmap.Uint32(0x5D4594, 1092996+4*uintptr(frame>>3)))
	}
	uiMeterAdvanceCharge()
	return 1
}

func uiMeterChargeRaster(w *gui.Window) int {
	m := &uiMeters()[uintptr(w.WidgetData)]
	pos := uiWindowPosition(w)
	row := func(index int) []byte {
		return unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 147904+uintptr(index)*8)), 8)
	}
	if int32(m.Maximum) < 1 {
		uiMeterSetColor(m.Color)
		nox_client_drawEnableAlpha_434560(1)
		for i := 0; i < 61; i++ {
			b := row(i)
			nox_client_drawRectFilledOpaque_49CE30(pos.X+int(b[0]), pos.Y+int(b[1]), int(b[2]), 1)
		}
		nox_client_drawEnableAlpha_434560(0)
		return 1
	}
	count := int(int32(m.Maximum))
	gap := 1
	if count > 30 {
		gap = 0
	}
	seq := uint32(GetClient().GetInputSeq())
	span := (gap+61)/count - gap
	next := 61 - span
	remainder := float32(float64(gap-count*((gap+61)/count)+61) / float64(count))
	fraction := float32(0.001)
	one := math.Float64frombits(uint64(qword_581450_9512))
	for i := 0; i < count; i++ {
		if uint32(i) >= m.Current {
			uiMeterSetColor(m.Alternate)
		} else {
			uiMeterSetColor(m.Color)
		}
		start := next
		next -= span + gap
		length := span
		if length <= 0 {
			length = 1
		}
		sum := float64(fraction) + float64(remainder)
		fraction = float32(sum)
		if sum >= one {
			next--
			start--
			length++
			fraction = float32(float64(fraction) - one)
		}
		nox_client_drawEnableAlpha_434560(1)
		if start < 0 {
			length += start
			start = 0
		}
		for j := 0; j < length && start+j < 61; j++ {
			b := row(start + j)
			stamp := (*uint32)(unsafe.Pointer(&b[4]))
			if *stamp != seq {
				nox_client_drawRectFilledOpaque_49CE30(pos.X+int(b[0]), pos.Y+int(b[1]), int(b[2]), 1)
				*stamp = seq
			}
		}
		nox_client_drawEnableAlpha_434560(0)
	}
	return 1
}

func uiMeterWeaponDraw(w *gui.Window) int {
	m := &uiMeters()[uintptr(w.WidgetData)]
	pos := uiWindowPosition(w)
	width, height := w.Size().X, w.Size().Y
	draw := m.Maximum != 0
	sector := 0
	if draw {
		sector = int((int32(m.Current) << 8) / int32(m.Maximum))
	}
	m.Alternate = m.Color
	if draw && sector >= 256 {
		item := uintptr(unsafe.Pointer(uiInventorySelectedWeapon()))
		sector = 1
		if item == 0 {
			draw = false
		} else {
			current := float64(*(*uint16)(unsafe.Pointer(item + 292)))
			maximum := float64(*(*uint16)(unsafe.Pointer(item + 294)))
			if current < maximum*memmap.Float64(0x581450, 9608) {
				m.Alternate = memmap.Uint32(0x85B3FC, 940)
			} else if current < maximum*math.Float64frombits(uint64(qword_581450_9544)) {
				m.Alternate = uint32(nox_color_yellow_2589772)
			} else {
				draw = false
			}
		}
	}
	if draw {
		nox_client_drawEnableAlpha_434560(1)
		GetClient().R2().Data().SetAlpha(64)
		sub_4AE6F0(pos.X+width/2, pos.Y+height/2, width/2, sector, int(m.Alternate))
		nox_client_drawEnableAlpha_434560(0)
	}
	uiInventoryCurrentWeaponDraw(w)
	return 1
}
