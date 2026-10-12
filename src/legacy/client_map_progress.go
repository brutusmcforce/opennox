package legacy

import (
	"encoding/binary"
	"fmt"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func clientMapCompassInit() int {
	*effectMapped(1309664) = 0
	for i := 0; i < 4; i++ {
		*effectMapped(1309644 + 4*uintptr(i)) = memmap.PtrWord(unsafe.Pointer(Nox_xxx_gLoadImg(fmt.Sprintf("Compass%d", i+1)).C()))
	}
	for i := 0; i < 32; i++ {
		*effectMapped(1309516 + 4*uintptr(i)) = memmap.PtrWord(unsafe.Pointer(Nox_xxx_gLoadImg(fmt.Sprintf("CompassMainArrow%d", i+1)).C()))
	}
	return 1
}

func clientMapProgress(data []byte) int16 {
	value := binary.LittleEndian.Uint16(data[1:])
	if uint32(value) == *effectMapped(1309668) {
		return int16(value)
	}
	*effectMapped(1309668) = uint32(value)
	audioEventPlay(897, 50, 0, 0)
	width, height := int(nox_win_width), int(nox_win_height)
	uiRenderBounds(0, 0, width-1, height-1)
	Nox_client_clearScreen_440900()
	r := GetClient().R2()
	frame := *effectMapped(1309672)
	pos := image.Pt(width/2-160, height/2-120)
	for _, offset := range []uintptr{1309644 + 4*uintptr(frame%4), 1309516 + 4*uintptr(frame)} {
		img := r.GetBag().AsImage(noxrender.ImageHandle(memmap.WordPtr(*effectMapped(offset))))
		r.DrawImageAt(img, pos)
	}
	var key string
	switch data[0] {
	case 155:
		key = "Generating"
	case 156:
		key = "Assembling"
	case 157:
		key = "Populating"
	}
	if key != "" {
		text := GetServer().S().Strings().GetStringInFile(strman.ID(key), "guigen.c")
		*effectMapped(1309660) = memmap.PtrWord(unsafe.Pointer(alloc.InternCString16(text)))
	}
	text := alloc.GoString16((*uint16)(memmap.WordPtr(*effectMapped(1309660))))
	face := r.GetFonts().AsFont(nil)
	size := r.GetStringSizeWrapped(face, text, 0)
	x, y := (width-size.X)/2, height/2-(2*r.FontHeight(face)+70)
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	r.Data().SetColor(noxcolor.RGBA5551(memmap.Uint32(0x852978, 4)))
	r.DrawStringHL(face, text, image.Pt(x, y))
	*effectMapped(1309672)++
	if *effectMapped(1309672) >= 32 {
		*effectMapped(1309672) = 0
	}
	Nox_video_callCopyBackBuffer_4AD170()
	return int16(value)
}
