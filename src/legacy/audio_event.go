package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
)

func nox_xxx_getSevenDwords3_501940(i int32) int32 {
	return GetServer().S().Audio.Field12(sound.ID(i))
}

func nox_xxx_aud_501960(a1 int32, a2p *nox_object_t, a3 int, a4 int32) {
	GetServer().S().Audio.EventObj(sound.ID(a1), asObjectS(a2p), a3, uint32(a4))
}

func Sub_501C00(p types.Pointf, obj *server.Object) int { return objectReportSoundLevel(p, obj) }
