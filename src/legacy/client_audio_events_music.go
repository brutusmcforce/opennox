package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/music"
	"unsafe"
)

func audioEventMusicSlot(index uint32) *music.MusicState {
	return (*music.MusicState)(memmap.PtrOff(0x5D4594, 815772+uintptr(16*(index+6*(*audioEventMusicLevel)))))
}
func audioEventMusicSave() int32 {
	if int32(*audioEventMusicCount) < 6 {
		*audioEventMusicSlot(*audioEventMusicCount) = MusicModule.GetCurrentBlock()
		*audioEventMusicCount++
		return 1
	}
	*audioEventMusicCount = 6
	return 0
}
func audioEventMusicRestore() {
	if *audioEventMusicCount > 0 {
		*audioEventMusicCount--
		MusicModule.SetNextMusic(*audioEventMusicSlot(*audioEventMusicCount))
	}
	*audioEventMusicCount = 0
}
func audioEventMusicEnter() int32 {
	if int32(*audioEventMusicLevel) < 3 {
		audioEventMusicSave()
		v := *audioEventMusicLevel
		count := *audioEventMusicCount
		*audioEventMusicCount = 0
		*memmap.PtrUint32(0x5D4594, 816076+uintptr(4*v)) = count
		*audioEventMusicLevel = v + 1
		return int32(v + 1)
	}
	*audioEventMusicLevel = 3
	return 3
}
func audioEventMusicLeave() {
	if *audioEventMusicLevel > 0 {
		*audioEventMusicLevel--
		*audioEventMusicCount = *memmap.PtrUint32(0x5D4594, 816076+uintptr(4*(*audioEventMusicLevel)))
		audioEventMusicRestore()
	} else {
		*audioEventMusicLevel = 0
	}
}
func audioEventMusicSlotPointer(index int32) unsafe.Pointer {
	return unsafe.Pointer(audioEventMusicSlot(uint32(index)))
}
func audioEventByte() uint8         { return *memmap.PtrUint8(0x5D4594, 831252) }
func audioEventSetByte(v int8) int8 { *memmap.PtrUint8(0x5D4594, 831252) = uint8(v); return v }

func audioEventMusicSetCount(v int32) int32 { *audioEventMusicCount = uint32(v); return v }
