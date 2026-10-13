package opennox

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

func (s *Server) NetUpdateRemotePlrAudioEvents(obj *server.Object, v2 unsafe.Pointer, v18 int8) {
	s.Audio.ResetBitmap()
	var tm *server.Team
	if obj.HasTeam() {
		tm = s.Teams.ByID(obj.TeamVal.ID)
	}
	s.Audio.EachEvent(func(it *server.AudioEvent) {
		if it.Kind == 1 {
			tm2 := s.Teams.ByID(server.TeamID(it.Code))
			if tm == nil || tm2 == nil || tm != tm2 {
				return
			}
		} else if it.Kind == 2 {
			if obj.NetCode != it.Code {
				return
			}
		}
		v12 := int32(legacy.Sub_501C00(it.Pos, it.Obj))
		if v12 == int32(v18) || v12 == 0 {
			if int32(layout.As[server.PlayerUpdateData](v2).Field47_0) != 0 || it.Sound < sound.SoundSpellPhonemeUp || it.Sound > sound.SoundSpellPhonemeUpLeft || obj != it.Obj {
				fade := s.ai.soundFadePerc(it.Sound, it.Pos, layout.As[server.Player](unsafe.Pointer(layout.As[server.PlayerUpdateData](v2).Player)).Pos3632Vec) / 2
				if fade > 0 {
					if s.Audio.Field20(it.Sound) != 0 {
						s.Audio.AddAudio(it, fade)
					} else {
						s.netSendAudioEvent(obj, it, int16(fade))
					}
				}
			}
		}
	})
	s.netSendAudioEvents(obj)
}

func (s *Server) netSendAudioEvents(obj *server.Object) {
	s.Audio.EachEventBitmap(func(it *server.AudioEvent) {
		s.netSendAudioEvent(obj, it, int16(it.Perc))
	})
}
