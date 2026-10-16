package sdl

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Audio playback through SDL3 audio streams, all bound to one logical device that is opened once: SDL mixes them and
// converts the data to the device's format. A stream is paused by unbinding it, which keeps its queued data.
// (A logical device per stream would be simpler, but opening and closing one takes several to tens of ms.)

// AudioFormat is SDL_AudioFormat.
type AudioFormat uint32

// AUDIO_S16 is signed 16-bit samples in the native byte order (little endian on all supported platforms).
const AUDIO_S16 AudioFormat = 0x8010

// AudioSpec is SDL_AudioSpec.
type AudioSpec struct {
	Format   AudioFormat
	Channels int32
	Freq     int32
}

const audioDeviceDefaultPlayback = 0xFFFFFFFF

var audioFn struct {
	once    sync.Once
	err     error
	openDev uintptr
	create  uintptr
	bind    uintptr
	unbind  uintptr
	devOf   uintptr
	put     uintptr
	queued  uintptr
	clear   uintptr
	format  uintptr
	destroy uintptr
	dev     uint32                                  // the shared logical playback device
	gain    func(stream uintptr, gain float32) bool // float argument: purego.SyscallN can't pass it
}

func loadAudio() error {
	audioFn.once.Do(func() {
		if audioFn.err = Load(); audioFn.err != nil {
			return
		}
		lib := loadedLib
		audioFn.err = lib.Syms(map[string]*uintptr{
			"SDL_OpenAudioDevice":      &audioFn.openDev,
			"SDL_CreateAudioStream":    &audioFn.create,
			"SDL_BindAudioStream":      &audioFn.bind,
			"SDL_UnbindAudioStream":    &audioFn.unbind,
			"SDL_GetAudioStreamDevice": &audioFn.devOf,
			"SDL_PutAudioStreamData":   &audioFn.put,
			"SDL_GetAudioStreamQueued": &audioFn.queued,
			"SDL_ClearAudioStream":     &audioFn.clear,
			"SDL_SetAudioStreamFormat": &audioFn.format,
			"SDL_DestroyAudioStream":   &audioFn.destroy,
		})
		if audioFn.err != nil {
			return
		}
		// default format: SDL picks the device's own
		audioFn.dev = uint32(call(audioFn.openDev, audioDeviceDefaultPlayback, 0))
		if audioFn.dev == 0 {
			audioFn.err = GetError()
			return
		}
		if p := lib.Sym("SDL_SetAudioStreamGain"); p != 0 {
			purego.RegisterFunc(&audioFn.gain, p)
		}
	})
	return audioFn.err
}

// AudioStream is an SDL_AudioStream; it plays while bound to the shared device (see Resume).
type AudioStream struct {
	p uintptr
}

// OpenAudioStream creates a playback stream for data in the given format, on the default device (opened with the
// first stream; SDL_INIT_AUDIO must be initialized). It starts paused: call Resume to play.
func OpenAudioStream(spec AudioSpec) (*AudioStream, error) {
	if err := loadAudio(); err != nil {
		return nil, err
	}
	// the output format is set by the device when the stream is bound
	p := call(audioFn.create, uintptr(unsafe.Pointer(&spec)), uintptr(unsafe.Pointer(&spec)))
	if p == 0 {
		return nil, GetError()
	}
	return &AudioStream{p: p}, nil
}

// Put queues audio data in the stream's input format.
func (s *AudioStream) Put(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return errorFrom(call(audioFn.put, s.p, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data))))
}

// PutInt16 queues 16-bit samples (AUDIO_S16).
func (s *AudioStream) PutInt16(data []int16) error {
	if len(data) == 0 {
		return nil
	}
	return s.Put(unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*2))
}

// Queued returns the number of input bytes not played yet.
func (s *AudioStream) Queued() int {
	n := int32(call(audioFn.queued, s.p))
	return max(int(n), 0)
}

// Clear drops the queued data.
func (s *AudioStream) Clear() {
	call(audioFn.clear, s.p)
}

// SetFormat changes the format of the data put from now on.
func (s *AudioStream) SetFormat(spec AudioSpec) error {
	return errorFrom(call(audioFn.format, s.p, uintptr(unsafe.Pointer(&spec)), 0))
}

// SetGain sets the volume (1 is unchanged).
func (s *AudioStream) SetGain(gain float32) {
	if audioFn.gain != nil {
		audioFn.gain(s.p, gain)
	}
}

// Pause stops playing the stream; the queued data stays.
func (s *AudioStream) Pause() {
	if !s.Paused() {
		call(audioFn.unbind, s.p)
	}
}

// Resume plays the stream.
func (s *AudioStream) Resume() {
	if s.Paused() {
		call(audioFn.bind, uintptr(audioFn.dev), s.p)
	}
}

// Paused reports whether the stream is paused (not bound to the device).
func (s *AudioStream) Paused() bool {
	return uint32(call(audioFn.devOf, s.p)) == 0
}

// Destroy stops and frees the stream.
func (s *AudioStream) Destroy() {
	if s.p != 0 {
		call(audioFn.destroy, s.p)
		s.p = 0
	}
}
