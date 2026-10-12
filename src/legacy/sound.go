package legacy

import (
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_soundPlayerDamageSound_5328B0  func(obj, obj2 *server.Object) int
	Nox_xxx_soundDefaultDamageSound_532E20 func(obj, obj2 *server.Object) int
)

func Nox_xxx_clientPlaySoundSpecial_452D80(a1 sound.ID, a2 int) {
	audioEventPlay(int32(a1), int32(a2), 0, 0)
}

func Sub_4133D0(a1 *server.Object) int {
	return bool2int(runtimeMaterial(a1))
}
