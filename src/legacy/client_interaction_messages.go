package legacy

import (
	"image"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func interactionCopyText(dst, src *uint16) *uint16 {
	for i := uintptr(0); ; i += 2 {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(src), i))
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), i)) = v
		if v == 0 {
			return dst
		}
	}
}
func interactionMessagesClear() *uint16 {
	var result *uint16
	for i := uintptr(0); i < 3; i++ {
		off := i * 644
		result = interactionCopyText(memmap.PtrUint16(0x5D4594, 823804+off), memmap.PtrUint16(0x5D4594, 825740))
		*memmap.PtrUint32(0x5D4594, 824440+off) = 0
		*memmap.PtrUint8(0x5D4594, 824444+off) = 0
	}
	interactionMessageHead = 0
	return result
}
func interactionCentered(text *uint16) {
	if text == nil {
		return
	}
	slot := uint32(interactionMessageHead) + 1
	if slot == 3 {
		slot = 0
	}
	interactionMessageHead = uint32(slot)
	off := uintptr(slot * 644)
	dst := unsafe.Slice(memmap.PtrUint16(0x5D4594, 823804+off), 318)
	i := 0
	for i < 317 {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(text), 2*i))
		if v == 0 {
			break
		}
		dst[i] = v
		i++
	}
	dst[i] = 0
	*memmap.PtrUint32(0x5D4594, 824440+off) = gameFrame() + 5*gameFPS()
	*memmap.PtrUint8(0x5D4594, 824444+off) = 0
	format := alloc.InternCString16(GetServer().S().Strings().GetStringInFile("systemmsg", `C:\NoxPost\src\Client\Gui\guimsg.c`))
	textFormatConsole(byte(console.ColorRed), format, textFormatPointer(unsafe.Pointer(text)))
}
func interactionMessagesDraw() int32 {
	r := GetClient().R2()
	vp := drawViewport()
	y := 3*int32(vp.Size.Y)/4 + int32(vp.Screen.Min.Y) - 15
	slot := uint32(interactionMessageHead)
	face := r.GetFonts().AsFont(nil)
	for row := 0; row < 3; row++ {
		off := uintptr(slot * 644)
		frame := gameFrame()
		if memmap.Uint32(0x5D4594, 824440+off) < frame {
			return int32(frame)
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_black_2650656))
		text := alloc.GoString16(memmap.PtrUint16(0x5D4594, 823804+off))
		x := (int32(nox_win_width) - int32(r.GetStringSizeWrapped(face, text, 0).X)) / 2
		for i := uintptr(0); i < 4; i++ {
			dx := memmap.Int32(0x587000, 107848+8*i)
			dy := memmap.Int32(0x587000, 107852+8*i)
			r.DrawString(face, text, image.Pt(int(x+dx), int(y+dy)))
		}
		color := uint32(nox_color_white_2523948)
		if row != 0 {
			color = memmap.Uint32(0x5D4594, 2597996)
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(color))
		r.DrawString(face, text, image.Pt(int(x), int(y)))
		y -= 4 + int32(r.FontHeight(face))
		result := int32(slot)
		if slot == 0 {
			slot = 2
		} else {
			slot--
		}
		if row == 2 {
			return result
		}
	}
	panic("unreachable message row")
}
