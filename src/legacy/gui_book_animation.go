package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func bookFloat(off uintptr) *float32 { return (*float32)(unsafe.Pointer(bookWord(off))) }
func bookNormalizeStep() {
	geometryNormalize((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&bookVector))))
	speed := float32(10)
	if nox_win_width < 750 {
		speed = 6
	} else if nox_win_width < 1000 {
		speed = 8
	}
	bookVector[0] *= speed
	bookVector[1] *= speed
}
func bookStoreInBar() {
	quickbarBookSlot(uintptr(*bookWord(1046676)), *bookWord(1047524), int(*bookWord(1046852)))
}
func bookAdd(kind, id int) {
	if *bookWord(1047520) == 1 {
		return
	}
	quickbarCloseExpanded()
	*bookWord(1046612) = uint32(quickbarMain().Selected)
	if kind == 2 && bool(nox_xxx_spellHasFlags_424A50(id, 0x15000)) {
		return
	}
	if quickbarContains(uint32(id)) == 1 {
		return
	}
	var slot int
	switch kind {
	case 2, 4:
		slot = int(quickbarFirstEmpty(false))
	case 3:
		slot = int(quickbarFirstEmpty(true))
	default:
		return
	}
	*bookWord(1046852) = uint32(slot)
	if slot == -1 {
		return
	}
	*bookWord(1047520) = 1
	bookPageComplete(kind == 4)
	from := (*image.Point)(unsafe.Pointer(bookWord(1046844)))
	bookIconPosition(from)
	*bookFloat(1046636) = float32(from.X)
	*bookFloat(1046640) = float32(from.Y)
	moving := bookWindow(*bookWord(1046956))
	moving.SetPos(*from)
	quickbarSlotPosition(slot, (*[2]int32)(unsafe.Pointer(bookWord(1046668))))
	*bookWord(1047524) = uint32(id)
	*bookWord(1046652) = 0
	if kind == 3 {
		*bookWord(1046652) = 1
	}
	*bookWord(1046676) = uint32(kind)
	nox_client_setRenderGUI(1)
	*bookWord(1046680) = 0
	target := *(*image.Point)(unsafe.Pointer(bookWord(1046668)))
	effectCurveSegments([4]image.Point{*from, target, {X: 400, Y: -500}, {X: 350, Y: 400}}, 19, 0, func(a, b image.Point) { bookPath(a, b) })
	*bookWord(1046628) = 0
	bookVector[0] = *bookFloat(1046692) - *bookFloat(1046684)
	bookVector[1] = *bookFloat(1046696) - *bookFloat(1046688)
	bookNormalizeStep()
	bookHideWindow(moving, false)
	moving.ShowModal()
	if noxflags.HasGame(noxflags.GameModeCoop) {
		runtimePauseStart(nil, int32(kind))
	}
	*bookWord(1046648) = uint32(nox_xxx_bookGet_430B40_get_mouse_prev_seq())
	if !noxflags.HasGame(noxflags.GameModeCoop) || sessionQuitShown() == 1 || nox_xxx_gameGet_4DB1B0() == 1 {
		bookFinishAddition()
	}
}
func bookDrawAddition(w *gui.Window) int {
	kind := 0
	if *bookWord(1046652) != 0 {
		kind = 3
	}
	if *bookWord(1047520) == 0 {
		return 1
	}
	start := *bookWord(1046648)
	if start != 0 && uint32(nox_xxx_bookGet_430B40_get_mouse_prev_seq())-start < uint32(2*GetServer().S().TickRate()) {
		return 1
	}
	pos := w.GlobalPos()
	if start != 0 {
		for i := 0; i < 50; i++ {
			timer := byte(effectRand(3, 6))
			size := byte(effectRand(2, 5))
			vy := effectRand(-10, -1)
			vx := effectRand(-10, 10)
			y := pos.Y + effectRand(0, 30)
			x := pos.X + effectRand(0, 30)
			screenParticleCreate(kind, x, y, vx, vy, 1, size, timer, 2, 1)
		}
		bookSound(795)
		*bookWord(1046648) = 0
	}
	for i := 0; i < 2; i++ {
		timer := byte(effectRand(2, 4))
		size := byte(effectRand(1, 2))
		y := pos.Y + effectRand(0, 30)
		x := pos.X + effectRand(0, 30)
		screenParticleCreate(kind, x, y, 0, 0, 0, size, timer, 1, 1)
	}
	var img uint32
	if *bookWord(1046652) == 1 {
		img = bookAbilityImage(int(*bookWord(1047524)))
	} else {
		img = bookSpellImage(int(*bookWord(1047524)))
	}
	bookDrawImage(img, pos)
	*bookFloat(1046636) += bookVector[0]
	*bookFloat(1046640) += bookVector[1]
	if float64(int32(*bookWord(1046668))) <= float64(*bookFloat(1046636)) && float64(int32(*bookWord(1046672))) <= float64(*bookFloat(1046640)) {
		bookStoreInBar()
		bookStopAddition()
	} else if *bookFloat(1046636) > memmap.Float32(0x5D4594, 1046692+8*uintptr(*bookWord(1046628))) {
		*bookWord(1046628)++
		index := int32(*bookWord(1046628))
		count := int32(*bookWord(1046680))
		if index < count {
			if index <= count-1 {
				off := 8 * uintptr(uint32(index))
				bookVector[0] = *bookFloat(1046692 + off) - *bookFloat(1046684 + off)
				bookVector[1] = *bookFloat(1046696 + off) - *bookFloat(1046688 + off)
				bookNormalizeStep()
			}
		} else {
			bookStoreInBar()
			bookStopAddition()
		}
	}
	w.SetPos(image.Pt(effectFloatInt(*bookFloat(1046636)), effectFloatInt(*bookFloat(1046640))))
	return 1
}
func bookFinishAddition() {
	if *bookWord(1047520) == 0 {
		return
	}
	kind := 0
	if *bookWord(1046652) == 1 {
		kind = 3
	}
	x, y := effectFloatInt(*bookFloat(1046636)), effectFloatInt(*bookFloat(1046640))
	dx, dy := (int(int32(*bookWord(1046668)))-x)/50, (int(int32(*bookWord(1046672)))-y)/50
	for i := 0; i < 50; i++ {
		size := byte(effectRand(3, 4))
		py := y - effectRand(0, 30)
		px := x + effectRand(0, 30)
		screenParticleCreate(kind, px, py, 0, 0, 1, size, 0, 0, 1)
		x += dx
		y += dy
	}
	bookStoreInBar()
	bookStopAddition()
}
