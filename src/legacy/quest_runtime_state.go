package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func questRuntimeDeadlinePtr() *uint32  { return (*uint32)(unsafe.Pointer(&dword_5d4594_1556144)) }
func questRuntimeSoulFramePtr() *uint32 { return (*uint32)(unsafe.Pointer(&dword_5d4594_1556136)) }
func questRuntimePreviousStage(stage uint32) uint32 {
	old := uint32(dword_5d4594_1556128)
	*memmap.PtrUint32(0x5D4594, 1556132) = old
	dword_5d4594_1556128 = uint32(stage)
	return old
}
func questRuntimeWord(off uintptr) uint32 { return memmap.Uint32(0x5D4594, off) }
func questRuntimeSetWord(off uintptr, value uint32) uint32 {
	*memmap.PtrUint32(0x5D4594, off) = value
	return value
}
func questRuntimeSetSoulFrame(value uint32) uint32 { *questRuntimeSoulFramePtr() = value; return value }
func questRuntimeStage() uint32                    { return memmap.Uint32(0x587000, 202028) }
func questRuntimeSetStage(value uint32)            { *memmap.PtrUint32(0x587000, 202028) = value }
func questRuntimeSlotMask(index uint32) uint32 {
	mask := ^(uint32(1) << (index & 31))
	*memmap.PtrUint32(0x5D4594, 1556300) &= mask
	return mask
}
func questRuntimeBookAllowed(id, first int) bool {
	return noxflags.HasGame(noxflags.GameModeQuest) || id < first || id > first+3
}
func questRuntimeSettings() unsafe.Pointer {
	if p := questSettings; p != nil {
		return p
	}
	p := memmap.PtrOff(0x5D4594, 371516)
	*controlByte(p, 100) &^= 16
	return p
}
func questRuntimeMapName(off uintptr) *byte { return (*byte)(memmap.PtrOff(0x973F18, off)) }
func questRuntimeStageMessage(to int, code byte, extra uint32) int {
	data := make([]byte, 69)
	data[0] = 240
	data[1] = code
	if code == 13 && extra != 0 {
		data[4] |= 1
	}
	if memmap.Uint32(0x5D4594, 2388656) != 0 {
		data[4] |= 2
	}
	binary.LittleEndian.PutUint16(data[2:], uint16(questRuntimeStage()))
	copy(data[5:37], GoStringP(unsafe.Pointer(questRuntimeMapName(3838))))
	copy(data[37:69], GoStringP(unsafe.Pointer(questRuntimeMapName(3806))))
	return gameplayReportSend(to, data, true, 1)
}
func questRuntimeSharedMessage(to int, value byte) int {
	return gameplayReportSend(to, []byte{240, 24, value}, true, 1)
}
func questRuntimeHighestMessage(to int, stage uint16) int {
	return gameplayReportSend(to, []byte{240, 29, byte(stage), byte(stage >> 8)}, false, 1)
}
func questRuntimeCodeMessage(to int, u *server.Object) int {
	code := *controlHalf(unsafe.Pointer(u), int(layoutServerObject.Off(40)) /* Extent */)
	return gameplayReportSend(to, []byte{240, 15, byte(code), byte(code >> 8)}, false, 1)
}
