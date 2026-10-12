package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	Sub_44A4A0                        func() int
	Nox_xxx_dialogMsgBoxCreate_449A10 func(win *gui.Window, title, text string, a4 gui.DialogFlags, a5, a6 func())
	Sub_449E00                        func(a1 string) int
	Sub_449E30                        func(a1 string) int
	Sub_449E60                        func(a1 int8) int
	Sub_449EA0                        func(a1 gui.DialogFlags)
	Sub_44A4E0                        func() int
	Sub_44A4B0                        func()
	Sub_44A360                        func(a1 int)
)

func Sub_41DA70(a1, a2 int) {
	onlineSessionEvent()
}
func Sub_445C20() {
	sessionQuitHide()
}
