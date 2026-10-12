package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

// ClientMainMenuDraw draws the backdrop and advances the original 48-byte records.
func ClientMainMenuDraw(win *gui.Window, draw *gui.WindowData) int {
	pos := win.GlobalPos()
	if !win.Flags.Has(gui.StatusImage) {
		if draw.BgColorVal != 0x80000000 {
			uiMeterSetColor(draw.BgColorVal)
			nox_client_drawRectFilledOpaque_49CE30(pos.X, pos.Y, win.SizeVal.X, win.SizeVal.Y)
		}
	} else {
		pos = pos.Add(draw.ImagePoint())
		handle := draw.BgImageHnd
		if draw.Field0&2 != 0 {
			handle = draw.HlImageHnd
		}
		uiMeterImage(uint32(uintptr(handle)), pos)
	}
	if memmap.Uint32(0x587000, 168836) != 0 {
		random := GetServer().S().Rand.Other
		for off := uintptr(168832); ; off += 48 {
			active := memmap.PtrUint32(0x587000, off+16)
			timer := memmap.PtrUint32(0x587000, off+36)
			flash := memmap.PtrUint32(0x587000, off+40)
			quiet := memmap.PtrUint32(0x587000, off+44)
			if *quiet != 0 {
				*quiet--
			}
			if *flash != 0 {
				*flash--
				*quiet = uint32(random.Int(60, 120))
			}
			oldTimer := *timer
			*timer--
			if oldTimer == 1 {
				if *active != 0 {
					*active = 0
					*timer = uint32(random.Int(int(int32(memmap.Uint32(0x587000, off+20))), int(int32(memmap.Uint32(0x587000, off+24)))))
					*quiet = uint32(random.Int(60, 90))
				} else {
					*active = 1
					*timer = uint32(random.Int(int(int32(memmap.Uint32(0x587000, off+28))), int(int32(memmap.Uint32(0x587000, off+32)))))
				}
			} else if *active == 0 && *quiet == 0 && *flash == 0 && random.Int(0, 100) > 75 {
				*flash = uint32(random.Int(4, 8))
			}
			if memmap.Uint32(0x587000, off+48) == 0 {
				break
			}
		}
	}
	if memmap.Uint32(0x587000, 168832) != 0 {
		for off := uintptr(168832); ; off += 48 {
			if memmap.Uint32(0x587000, off+16) == 0 && memmap.Uint32(0x587000, off+40) == 0 {
				uiMeterImage(memmap.Uint32(0x587000, off+4), image.Pt(int(int32(memmap.Uint32(0x587000, off+8))), int(int32(memmap.Uint32(0x587000, off+12)))))
			}
			if memmap.Uint32(0x587000, off+48) == 0 {
				break
			}
		}
	}
	return 1
}
