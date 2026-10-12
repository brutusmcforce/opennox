package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

func Sub_46AF00(a1 *gui.Window) string {
	return GoWString((*wchar2_t)(uiWindowText(a1)))
}
func Sub_46AF40(a1 *gui.Window) unsafe.Pointer {
	return uiWindowFont(a1)
}
