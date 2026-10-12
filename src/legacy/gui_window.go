package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

type nox_window = gui.Window

func AsWindowP(win unsafe.Pointer) *gui.Window {
	w := (*gui.Window)(win)
	return w
}

func asWindow(win *nox_window) *gui.Window {
	return AsWindowP(unsafe.Pointer(win))
}

func get_dword_5d4594_3799468() int {
	return GetClient().Cli().GUI.ValXXX
}

func nox_xxx_windowDestroyMB_46C4E0(a1 *nox_window) int {
	win := asWindow(a1)
	if win == nil {
		return -2
	}
	win.Destroy()
	return 0
}

func nox_window_set_hidden(p *nox_window, hidden int) int {
	if p == nil {
		return -2
	}
	win := asWindow(p)
	if hidden != 0 {
		win.Hide()
	} else {
		win.Show()
	}
	return 0
}

func Nox_xxx_wnd_46ABB0(p *gui.Window, v int) int {
	return uiWindowEnable(p, v)
}
