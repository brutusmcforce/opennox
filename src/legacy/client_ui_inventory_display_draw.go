package legacy

import (
	"image"
	"math"
	"strconv"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/client/keybind"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/server"
)

func uiInventoryDrawItem(dr *client.Drawable, pos image.Point) {
	dr.PosVec = pos
	dr.CallDraw(&uiInventoryViewport)
}

func sub_4625D0(w *gui.Window) int {
	if uiInventoryMode() == 5 {
		return 1
	}
	pos := uiWindowPosition(w)
	size := w.SizeVal
	if pos.Y+size.Y > 0 {
		r := GetClient().R2()
		r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
		if cell := uiInventoryAlternate(); cell != nil && cell.Drawable != nil {
			uiInventoryDrawItem(cell.Drawable, pos.Add(image.Pt(size.X/2, size.Y/2)))
		}
		text := GetClient().GetCtrlEvent().Sub_42E8E0_go(keybind.Event(35), 1)
		r.DrawString(uiInventorySmallFont(), text, pos.Add(image.Pt(22, 41)))
	}
	return 1
}

func uiInventoryCurrentWeaponDraw(w *gui.Window) int {
	pos := uiWindowPosition(w.Parent())
	if dr := uiInventoryCurrentWeapon(); dr != nil {
		uiInventoryDrawItem(dr, pos.Add(image.Pt(51, 81)))
	} else if dword_5d4594_1062496 == 0 && dword_5d4594_1062492 == 0 {
		uiMeterImage(memmap.Uint32(0x5D4594, 1050000), pos.Add(image.Pt(21, 50)))
	}
	return 1
}

func nox_xxx_inventoryDrawProc_466580(w *gui.Window) int {
	pos := uiWindowPosition(w)
	img := w.DrawData().BgImageHnd
	if memmap.Uint8(0x5D4594, 1049868) != 0 {
		img = w.DrawData().HlImageHnd
	}
	uiMeterImage(memmap.PtrWord(unsafe.Pointer(img)), pos)
	r := GetClient().R2()
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	text := GetClient().GetCtrlEvent().Sub_42E8E0_go(keybind.Event(35), 1)
	r.DrawString(uiInventorySmallFont(), text, pos.Add(image.Pt(19, 102)))
	return 1
}

func sub_466F50(w *gui.Window, data *gui.WindowData) int {
	dr := uiInventorySelectedItem()
	if dr == nil {
		return 1
	}
	class := uint32(dr.Class())
	if class&0x13001000 != 0 {
		var def *server.Modifier
		if class&0x11001000 != 0 {
			def = GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
		} else {
			def = GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
		}
		if def != nil {
			for i := 1; i < 7; i++ {
				c := def.Colors12[i]
				nox_draw_setMaterial_4340A0(i, int(c.R), int(c.G), int(c.B))
			}
			indexes := def.ColorIndexes()
			for i := 0; i < 4; i++ {
				if m := uiInventoryItemModifier(dr, i); m != nil {
					c := m.Color24
					nox_draw_setMaterial_4340A0(int(indexes[i]), int(c.R), int(c.G), int(c.B))
				}
			}
		}
	}
	pos := uiWindowPosition(w)
	uiMeterImage(memmap.PtrWord(unsafe.Pointer(data.BgImageHnd)), pos.Add(data.ImgPtVal))
	return 1
}

func uiInventoryDrawTray(ax, ay int32) int32 {
	x, top := int(ax), int(ay)
	r := GetClient().R2()
	small := uiInventorySmallFont()
	uiMeterImage(memmap.Uint32(0x5D4594, 1049928), image.Pt(x, top))
	text := strconv.FormatInt(int64(int32(dword_5d4594_1062552)), 10)
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_yellow_2589772))
	width := r.GetStringSizeWrapped(small, text, 0).X
	r.DrawString(small, text, image.Pt(x-width+43, top+36))
	if uiInventoryMode() == 5 {
		uiMeterImage(memmap.Uint32(0x5D4594, 1049932), image.Pt(x, top+50))
	}
	if sub_473670() != 0 {
		uiMeterImage(memmap.Uint32(0x5D4594, 1049936), image.Pt(x, top+100))
	}
	y := top - int(dword_5d4594_1062512)
	for row := 0; row < 20; row++ {
		if y > top-50 {
			uiMeterImage(memmap.Uint32(0x5D4594, 1049916+uintptr(row%3)*4), image.Pt(x+60, y))
			for col := 0; col < 4; col++ {
				cell := &uiInventoryGrid()[col*21+row]
				dr := cell.Drawable
				if cell.Count == 0 || dr == nil {
					continue
				}
				left := x + 60 + col*50
				r.Data().SetAlphaEnabled(true)
				r.Data().SetAlpha(64)
				current := dr.Field_73_1
				maximum := dr.Field_73_2
				var color uint32
				shade := false
				if float64(current) < float64(maximum)*memmap.Float64(0x581450, 9608) {
					color = memmap.Uint32(0x85B3FC, 940)
					shade = true
				} else if float64(current) < float64(maximum)*math.Float64frombits(uint64(qword_581450_9544)) {
					color = uint32(nox_color_yellow_2589772)
					shade = true
				}
				if shade && color != 0x80000000 {
					uiMeterSetColor(color)
					nox_client_drawRectFilledOpaque_49CE30(left, y, 50, 50)
				}
				r.Data().SetAlpha(128)
				var decoration uint32
				if cell.Equipped != 0 {
					decoration = 1049964
				} else if cell.Alternate != 0 {
					decoration = 1049968
				} else if cell.Codes[0] == uint32(dword_5d4594_1062488) {
					if alt := uiInventoryAlternate(); alt != nil && alt.Drawable != nil && uint32(alt.Drawable.Class())&0x1000000 != 0 && uint32(alt.Drawable.SubClass())&12 != 0 {
						decoration = 1049968
					}
				}
				if decoration != 0 {
					uiMeterImage(memmap.Uint32(0x5D4594, uintptr(decoration)), image.Pt(left, y))
				}
				r.Data().SetAlphaEnabled(false)
				uiInventoryDrawItem(dr, image.Pt(left+25, y+25))
				item := layout.As[client.DrawableUnionItem](unsafe.Pointer(&dr.Union))
				charges, maxCharges := uint16(item.Field_112_0), uint16(item.Field_112_2)
				if uiInventoryMode() == 6 {
					depleted := uint32(dr.Class())&0x1000 != 0 && charges < maxCharges && maxCharges != 0
					if (current == maximum || maximum == 0) && !depleted {
						nox_client_drawRectFilledAlpha_49CF10(left, y, 50, 50)
					}
				}
				if cell.Count > 1 {
					r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
					r.DrawString(small, strconv.Itoa(int(cell.Count)), image.Pt(left+6, y+6))
				}
				if uint32(dr.Class())&0x13001000 != 0 && int16(charges) >= 0 {
					text = strconv.Itoa(int(int16(charges)))
					r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_blue_2650684))
					width = r.GetStringSizeWrapped(small, text, 0).X
					r.DrawString(small, text, image.Pt(left-width+44, y+6))
				}
			}
		}
		y += 50
		if y > top+150 {
			break
		}
	}
	return int32(top + 150)
}
