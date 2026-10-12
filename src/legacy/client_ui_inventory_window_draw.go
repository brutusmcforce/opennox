package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"strconv"
	"unsafe"
)

func uiInventoryDrawWindow(w *gui.Window) int {
	// Position uses the previous animation offset, before advancing this frame.
	w.Parent().SetPos(image.Pt(0, int(int32(dword_587000_136184))))
	pos := uiWindowPosition(w)
	nox_xxx_guiFontHeightMB_43F320(nil)
	x, y := pos.X+10, pos.Y+234
	mask := memmap.Uint32(0x5D4594, 1062540)
	for i := 0; i < 30; i++ {
		if mask&(uint32(1)<<i) != 0 {
			spell := int32(nox_xxx_getEnchantSpell_424920(int(int32(i))))
			uiMeterImage(memmap.PtrWord(nox_xxx_spellIcon_424A90(int(spell))), image.Pt(x, y))
			x += 35
		}
	}
	extra := memmap.Uint8(0x5D4594, 1062536)
	for i := 0; i < 6; i++ {
		if extra&(1<<i) != 0 {
			nox_client_drawImageAt_47D2C0(runtimeModifierIcon(byte(1<<i)), x, y)
			x += 35
		}
	}
	if noxflags.HasGame(4096) && dword_5d4594_1050008 != 0 {
		x += 6
		y += 5
		ref := AsImageRefP(memmap.WordPtr(dword_5d4594_1050008))
		anim := (*ImageRefAnim)(ref.Field_24)
		frames := anim.Images()
		uiMeterImage(memmap.PtrWord(unsafe.Pointer(frames[GetServer().S().Frame()%uint32(len(frames))])), image.Pt(x-58, y-53))
		r := GetClient().R2()
		r.Data().SetTextColor(noxcolor.RGBA5551(memmap.Uint32(0x852978, 0)))
		r.DrawString(uiInventorySmallFont(), "X "+strconv.FormatInt(int64(memmap.Int32(0x5D4594, 1050012)), 10), image.Pt(x+20, y+9))
		copy(unsafe.Slice((*int32)(memmap.PtrOff(0x5D4594, 1049812)), 4), []int32{int32(x - 30), int32(y - 20), int32(x + 30), int32(y + 20)})
	}
	if noxflags.HasGame(4096) && sub_4BFD30() != 0 {
		xx, yy := x+66, y+5
		uiMeterImage(memmap.Uint32(0x5D4594, 1050004), image.Pt(xx-64, yy-58))
		copy(unsafe.Slice((*int32)(memmap.PtrOff(0x5D4594, 1049828)), 4), []int32{int32(xx - 30), int32(yy - 20), int32(xx + 30), int32(yy + 20)})
	}
	state := memmap.PtrUint8(0x5D4594, 1049868)
	if *state != 0 {
		if pos.Y+163 > 0 {
			objectRenderSaveClip()
			uiRenderCopyRect(pos.X+254, pos.Y+13, 260, 150)
			switch memmap.Uint8(0x5D4594, 1049869) {
			case 0:
				uiInventoryDrawTray(int32(pos.X+254), int32(pos.Y+13))
			case 1:
				journalDraw(pos.X+254, pos.Y+13, int(dword_5d4594_1062512))
			}
			objectRenderRestoreClip()
		}
		if uiInventoryMode() == 5 {
			uiInventoryIdentify(pos)
			uiMeterImage(memmap.Uint32(0x5D4594, 1049912), pos)
		} else {
			switch memmap.Uint8(0x5D4594, 1049870) {
			case 0:
				point := [2]int32{int32(pos.X), int32(pos.Y)}
				sub_4BF7E0((*uint32)(unsafe.Pointer(&point[0])))
				uiMeterImage(memmap.Uint32(0x5D4594, 1049908), pos)
			case 1:
				uiInventoryStats(pos)
				uiMeterImage(memmap.Uint32(0x5D4594, 1049912), pos)
			}
			r := GetClient().R2()
			r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
			r.DrawStringWrapped(r.GetFonts().AsFont(nil), uiInventoryWideAt(1062588), image.Rect(pos.X+13, pos.Y+17, pos.X+209, pos.Y+17))
		}
	}
	switch *state {
	case 1:
		next := int32(dword_587000_136184) + 64
		dword_587000_136184 = uint32(next)
		if next > 0 {
			dword_587000_136184 = 0
			*state = 2
		}
	case 3:
		next := int32(dword_587000_136184) - 32
		dword_587000_136184 = uint32(next)
		if next <= -225 {
			dword_587000_136184 = uint32(^uint32(224))
			*state = 0
			switch memmap.Uint8(0x5D4594, 1049869) {
			case 0:
				dword_5d4594_1062516 = dword_5d4594_1062512
			case 1:
				dword_5d4594_1062520 = dword_5d4594_1062512
			}
			uiInventoryResetPanelControls()
		}
	}
	if uiInventoryWindowOpenState() && Nox_xxx_playerAnimCheck_4372B0() != 0 {
		uiInventoryCloseWindow()
	}
	return 1
}
