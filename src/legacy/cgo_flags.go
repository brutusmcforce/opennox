package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_moveUpdateSpecial_517970 func(obj *server.Object)
)

func nox_common_gameFlags_check_40A5C0(f uint32) bool {
	return noxflags.HasGame(noxflags.GameFlag(f))
}

func nox_xxx_CheckGameplayFlags_417DA0(v int) bool {
	return noxflags.HasGamePlay(noxflags.GameplayFlag(v))
}
