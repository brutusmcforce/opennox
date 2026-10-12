package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	Nox_gui_makeAnimation_43C5B0 func(win *gui.Window, x1, y1, x2, y2, in_dx, in_dy, out_dx, out_dy int) *gui.Anim
)

var _ = [1]struct{}{}[68-unsafe.Sizeof(gui.Anim{})]
