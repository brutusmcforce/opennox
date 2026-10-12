package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var (
	Sub_446380 func()
	Sub_445B40 func() int
)

func nox_xxx____setargv_4_44B000() {
	dword_5d4594_830272 = 1
}

func Sub_4D70B0() {
	questRuntimeSettings()
}

func Sub_509CB0() {
	matchRosterForget()
}

func Sub_41F4B0() {
	onlineSessionListCleanup()
}

func Sub_41EC30() {
	onlineSessionListCleanup()
}

func Sub_446490(v int) {
	sessionMOTDFree(v)
}

func Nox_xxx_guiServerOptionsHide_4597E0(v int) {
	serverOptionsClose(v)
}

func Sub_445C40() {
	sessionQuitToggle()
}

func Set_nox_wnd_quitMenu_825760(win *gui.Window) {
	sessionQuitRoot = win
}

func Nox_xxx____setargv_4_44B000() {
	nox_xxx____setargv_4_44B000()
}

func Get_dword_5d4594_2650652() int {
	return int(dword_5d4594_2650652)
}

func Set_dword_5d4594_2650652(v int) {
	dword_5d4594_2650652 = uint32(v)
}

func Sub_41CEE0(p unsafe.Pointer, a2 int) {
	playerFileClientWrite(p, a2)
}
