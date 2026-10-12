package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func nox_xxx_netInformTextMsg_4DA0F0(to, kind int32, data *int32) int32 {
	return int32(gameplayTextInformation(int(to), int(kind), unsafe.Pointer(data)))
}

func nox_xxx_netPriMsgToPlayer_4DA2C0(u *nox_object_t, text *int8, flag int8) {
	gameplayTextPrivate((*server.Object)(unsafe.Pointer(u)), (*byte)(unsafe.Pointer(text)), byte(flag))
}
