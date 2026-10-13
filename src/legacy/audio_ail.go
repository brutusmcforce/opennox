package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"

	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"github.com/brutusmcforce/opennox/legacy/timer"
)

var (
	Sub_43EA20 func(a1 unsafe.Pointer) int
	Sub_43E9F0 func()
	Sub_43E940 func(a1 unsafe.Pointer) int
	Sub_43EFD0 func(a1 unsafe.Pointer) int
	Sub_43EC10 func() int
	Sub_43F130 func() ail.Driver
	Sub_43ED00 func(a1p unsafe.Pointer) int
	Sub_43D6A0 func()
	Sub_44D640 func()
	Sub_44D7E0 func(a1 int) int
	Sub_44D660 func(a1 string) bool
	Sub_43F060 func(a1p unsafe.Pointer) int
	Sub_43EC30 func(a1 unsafe.Pointer) int
	Sub_43ECB0 func(a1 unsafe.Pointer) int
)

type AudioSample struct {
	Dev    ail.Driver     // 0, 0
	Field1 unsafe.Pointer // 1, 4
	Smp    ail.Sample     // 2, 8
	Field3 uint32         // 3, 12
	Field4 unsafe.Pointer // 4, 16
	Data1  *byte          // 5, 20
	Data2  *byte          // 6, 24
	Flag7  uint32         // 7, 28
}

func Get_dword_587000_127004() unsafe.Pointer {
	return legacyGlobals.dword_587000_127004
}

func Sub_43F0E0(v unsafe.Pointer) int {
	return int(audioEventSampleFormat((*AudioStreamFormat)(v)))
}

func Sub_43EE00(v *AudioSample) {
	audioEventSampleRefill(v)
}

func Sub_43EDB0(v ail.Sample) {
	audioEventSampleEnded(v)
}

func Sub_413890() string {
	return ""
}

func Nox_xxx_parseSoundSetBin_424170(path string) int {
	return resourceSoundLoad(path)
}

func Sub_43DC00() {
	audioEventMusicDisable()
}

func Sub_44D960() {
	audioEventDialogDisable()
}

func Sub_453050() {
	*audioEventPlayback = 0
}

func Get_dword_5d4594_816376() ail.Driver {
	return dword_5d4594_816376
}

func Set_dword_5d4594_816376(drv ail.Driver) {
	dword_5d4594_816376 = drv
}

func Sub_486640(g *timer.TimerGroup, a2 int) int {
	return int(int32(audioStreamScaleVolume(g, uint32(a2))))
}

func Get_dword_5d4594_805984() unsafe.Pointer {
	return legacyGlobals.dword_5d4594_805984
}

func Set_dword_5d4594_805984(v unsafe.Pointer) {
	legacyGlobals.dword_5d4594_805984 = v
}

func Set_dword_587000_81128(v unsafe.Pointer) {
	legacyGlobals.dword_587000_81128 = v
}

func Sub_451850(a1 unsafe.Pointer, a2 unsafe.Pointer) {
	audioEventInit((*AudioStreamContext)(a1), (*audioStreamCatalog)(a2))
}

// RegisterMilesAudioDevice registers the Miles Sound System stream device.
func RegisterMilesAudioDevice() {
	audioStreamDeviceRegister(audioMilesDescriptor())
}

var audioMilesDesc *audioStreamDescriptor

// audioMilesDescriptor builds the Miles Sound System stream descriptor. In the original it is static data
// at 0x587000+93992, with its voice API at +93952.
func audioMilesDescriptor() *audioStreamDescriptor {
	if audioMilesDesc == nil {
		key := audioBridgeCallbackKey
		api := &audioStreamVoiceAPI{
			Init: key(0), Free: key(1), Start: key(2), Stop: key(3),
			Fields20: [3]unsafe.Pointer{key(4), key(5), key(6)},
			Update:   key(7), Restart: key(8),
		}
		audioMilesDesc = &audioStreamDescriptor{
			Name:    (*byte)(memmap.PtrOff(0x587000, 94036)), // "Miles Sound System"
			Field16: 0x6d445541,
			Init:    key(9), Free: key(10), ContextInit: key(11), ContextFree: key(12),
			VoiceAPI: api,
		}
	}
	return audioMilesDesc
}

func Sub_487D00(a1 unsafe.Pointer) {
	audioStreamByteRate((*AudioStreamFormat)(a1))
}

func Sub_487150(a1 int, a2 unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(audioStreamContextAcquire(int32(a1), (*AudioStreamFormat)(a2)))
}

func Sub_487790(a1 unsafe.Pointer, a2 int) int {
	return int(audioStreamVoiceCreateMany((*AudioStreamContext)(a1), int32(a2)))
}

// AudioVoiceFormat returns the format of a voice's buffer, or the format its context requested.
func AudioVoiceFormat(p unsafe.Pointer) unsafe.Pointer {
	v := (*audioStreamVoice)(p)
	if f := *(*unsafe.Pointer)(unsafe.Pointer(&v.Buffer.Format)); f != nil {
		return f
	}
	return unsafe.Pointer(&v.Context.RequestedFormat)
}

// AudioVoiceSample returns the sample of an audio stream voice.
func AudioVoiceSample(voice unsafe.Pointer) *AudioSample {
	return (*AudioSample)((*audioStreamVoice)(voice).Sample)
}

// AudioVoiceEffective returns the effective volume, rate and pan timers of an audio stream voice.
func AudioVoiceEffective(voice unsafe.Pointer) *timer.TimerGroup {
	return &(*audioStreamVoice)(voice).Effective
}

// AudioVoiceSetSample sets the sample of an audio stream voice.
func AudioVoiceSetSample(voice unsafe.Pointer, s *AudioSample) {
	(*audioStreamVoice)(voice).Sample = unsafe.Pointer(s)
}

// AudioVoiceEndCallback returns the end callback of an audio stream voice.
func AudioVoiceEndCallback(voice unsafe.Pointer) unsafe.Pointer {
	return (*audioStreamVoice)(voice).EndCallback
}
