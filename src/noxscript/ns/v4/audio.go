package ns

import (
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/audio"
)

// AudioEvent play an audio event at a location.
func AudioEvent(audio audio.Name, p Positioner) {
	if impl == nil {
		return
	}
	impl.AudioEvent(audio, p)
}

// PlayAudioFile plays a WAV file from the map folder (e.g. "voice/intro.wav") for all players, like a speech.
//
// The file is sent to each player the first time it plays (OpenNox players only), so keep it small.
func PlayAudioFile(name string) {
	if impl == nil {
		return
	}
	impl.PlayAudioFile(name)
}

// Music plays music.
func Music(music int, volume int) {
	if impl == nil {
		return
	}
	impl.Music(music, volume)
}
func MusicPushEvent() {
	if impl == nil {
		return
	}
	impl.MusicPushEvent()
}
func MusicPopEvent() {
	if impl == nil {
		return
	}
	impl.MusicPopEvent()
}
func MusicEvent() {
	if impl == nil {
		return
	}
	impl.MusicEvent()
}
