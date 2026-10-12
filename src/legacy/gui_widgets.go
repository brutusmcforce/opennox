package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	NewButtonOrCheckbox func(parent *gui.Window, status gui.StatusFlags, px, py, w, h int, draw *gui.WindowData) *gui.Window
)

func Nox_gui_newScrollListBox_4A4310(par *gui.Window, status gui.StatusFlags, px, py, w, h int, draw *gui.WindowData, tdata *gui.ScrollListBoxData) *gui.Window {
	return uiListNew(par, status, px, py, w, h, draw, tdata)
}

func Nox_gui_newEntryField_488500(par *gui.Window, status gui.StatusFlags, px, py, w, h int, draw *gui.WindowData, tdata *gui.EntryFieldData) *gui.Window {
	return uiEntryNew(par, status, px, py, w, h, draw, tdata)
}

func Nox_gui_newSlider_4B4EE0(par *gui.Window, status gui.StatusFlags, px, py, w, h int, draw *gui.WindowData, tdata *gui.SliderData) *gui.Window {
	return uiSliderNew(par, status, px, py, w, h, draw, tdata)
}

func Nox_gui_newProgressBar_4CAF10(par *gui.Window, status gui.StatusFlags, px, py, w, h int, draw *gui.WindowData) *gui.Window {
	return uiProgressNew(par, status, px, py, w, h, draw)
}

func Nox_xxx_wndRadioButtonSetAllFn_4A87E0(win *gui.Window) {
	uiRadioInit(win)
}

var (
	Nox_xxx_wndRadioButtonProcPre_4A93C0 = uiRadioEvent
)
