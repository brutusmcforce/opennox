package server

import (
	"github.com/brutusmcforce/opennox/noxscript/ns/v4"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/audio"

	"github.com/brutusmcforce/opennox/common/sound"
)

func (s NoxScriptNS) AudioEvent(audio audio.Name, p ns.Positioner) {
	if p == nil {
		return
	}
	snd := sound.ByName(string(audio))
	s.s.Audio.EventPos(snd, p.Pos(), 0, 0)
}
