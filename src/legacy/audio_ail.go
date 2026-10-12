package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
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

var _ = [1]struct{}{}[32-unsafe.Sizeof(AudioSample{})]

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
	return int(audioEventSampleFormat((*audioStreamFormat)(v)))
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
	return ail.Driver(dword_5d4594_816376)
}

func Set_dword_5d4594_816376(drv ail.Driver) {
	dword_5d4594_816376 = uint32(drv)
}

func Sub_486640(a1 unsafe.Pointer, a2 int) int {
	return int(int32(audioStreamScaleVolume(a1, uint32(a2))))
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
	audioEventInit((*audioStreamContext)(a1), (*audioStreamCatalog)(a2))
}

func Sub_486FA0(a1 int) {
	audioStreamDeviceRegister((*audioStreamDescriptor)(unsafe.Pointer(uintptr(a1))))
}

func Sub_487D00(a1 unsafe.Pointer) {
	audioStreamByteRate((*audioStreamFormat)(a1))
}

func Sub_487150(a1 int, a2 unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(audioStreamContextAcquire(int32(a1), (*audioStreamFormat)(a2)))
}

func Sub_487790(a1 unsafe.Pointer, a2 int) int {
	return int(audioStreamVoiceCreateMany((*audioStreamContext)(a1), int32(a2)))
}
