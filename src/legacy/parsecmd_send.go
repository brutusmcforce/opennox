package legacy

import "unsafe"
import "github.com/brutusmcforce/opennox/server"

func Nox_xxx_serverHandleClientConsole_443E90(a1 *server.Player, a2 byte, cmd string) {
	consoleCommandRemote(a1, a2, cmd)
}
func Nox_xxx_cmdSayDo_46A4B0(text string, a2 int) {
	cstr, free := CWString(text)
	defer free()
	interactionSay((*uint16)(unsafe.Pointer(cstr)), int(int32(a2)))
}
