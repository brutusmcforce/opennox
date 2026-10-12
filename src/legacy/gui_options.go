package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	Nox_video_setMenuOptions func(win *gui.Window)
	Nox_gui_menu_proc_ext    func(id int) int
	Sub_4A19F0               func(name strman.ID)
	Sub_4AAA10               func() int
	Sub_4C3A90               func(a1, a2 int, a3 unsafe.Pointer, a4 int) int
	Sub_4CBE70               func(a1, a2 int, a3 unsafe.Pointer, a4 int) int
	Sub_4A1A40               func(a1 int)
)

func Sub_4CBD30() {
	bindingMenu.apply()
}

func Sub_430AA0(v int) {
	interactionMouseMode(int32(v))
}

func Sub_4C35B0(v int) {
	bindingClose(v)
}

func Get_nox_wnd_xxx_1309740() *gui.Anim {
	return legacyGlobals.nox_wnd_xxx_1309740
}

func Get_dword_5d4594_1309720() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(legacyGlobals.dword_5d4594_1309720)))
}
