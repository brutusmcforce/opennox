package legacy

import "github.com/brutusmcforce/opennox/client/gui"

func nox_client_setCursorType_477610(v int) int {
	GetClient().Nox_client_setCursorType(gui.Cursor(v))
	return v
}
