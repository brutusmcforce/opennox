package legacy

import (
	"os"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/log"

	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	guiLog   = log.New("gui")
	guiDebug = os.Getenv("NOX_DEBUG_GUI") == "true"
)

var (
	Nox_client_gui_set_flag_815132 func(v int)
	Nox_client_onClientStatusA     func(v int)
	Nox_client_setRenderGUI        func(v int)
	Nox_client_getRenderGUI        func() int
)

var _ = [1]struct{}{}[332-unsafe.Sizeof(gui.WindowData{})]

func nox_client_onClientStatusA(v int) { Nox_client_onClientStatusA(v) }

func nox_client_setRenderGUI(v int) { Nox_client_setRenderGUI(v) }

func nox_color_rgb_4344A0(r, g, b int) uint32 {
	return noxcolor.RGB5551Color(byte(r), byte(g), byte(b)).Color32()
}

func nox_set_color_rgb_434430(r, g, b int) {
	GetClient().R2().Data().SetColor2(noxcolor.RGB5551Color(byte(r), byte(g), byte(b)))
}

func Sub_46A4A0() int {
	return int(sub_46A4A0())
}

func Nox_xxx_wndEditProc_487D70(a1 *gui.Window, ev gui.WindowEvent) gui.RawEventResp {
	return gui.RawEventResp(gui.EventRespInt(uiEntryInput(a1, ev)))
}

func Nox_gui_xxx_check_446360() int {
	return int(sessionQuitShown())
}
