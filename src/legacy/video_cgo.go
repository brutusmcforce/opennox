package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/client"
)

var (
	Nox_video_getCutSize_4766D0             func() int
	Nox_video_setCutSize_4766A0             func(v int)
	Nox_video_setGammaSlider                func(v int)
	Sub_43BE50_get_video_mode_id            func() int
	Get_video_mode_string                   func(id int) string
	Nox_getBackbufWidth                     func() int
	Nox_getBackbufHeight                    func() int
	Nox_video_getFullScreen                 func() int
	Nox_video_setFullScreen                 func(v int)
	Sub_430C30_set_video_max                func(w, h int)
	VideoGetMaxSize                         func() image.Point
	Nox_video_callCopyBackBuffer_4AD170     func()
	Nox_getBackbufferPitch                  func() int
	Nox_client_clearScreen_440900           func()
	Nox_draw_setCutSize_476700              func(cutPerc int, a2 int)
	Nox_xxx_bookSaveSpellForDragDrop_477640 func(a1, a2 int)
	Nox_xxx_bookSpellDnDclear_477660        func()
	Nox_xxx_bookGetSpellDnDType_477670      func() int
	Nox_xxx_cursorSetDraggedItem_477690     func(a1 *client.Drawable)
	Nox_xxx_cursorResetDraggedItem_4776A0   func()
	Sub_478000                              func() int
)

func sub_478000() int { return Sub_478000() }
func Sub_4AEE30() {
	runtimeMeterWave()
}
func Nox_xxx_guiSpell_460650() int {
	return int(*quickbarWord(1047928))
}
func Sub_4611A0() int {
	return int(*quickbarWord(1047932))
}
func Sub_467CD0() {
	uiInventoryCancelDrag()
}
func Sub_49F6D0(a1 int) int {
	return int(uiRenderFlag(uint32(a1)))
}
func Sub_430B50(a1 int, a2 int, a3 int, a4 int) {
	uiRenderBounds(a1, a2, a3, a4)
}
func Sub_495A80(a1 uint32) int {
	return bool2int(combatFriendHas(a1))
}
