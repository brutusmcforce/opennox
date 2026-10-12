package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"github.com/brutusmcforce/opennox/legacy/music"
	"github.com/brutusmcforce/opennox/legacy/timer"
)

var (
	MusicModule *music.Module
)

func initMusic() {
	var (
		dword_5d4594_816368      *uint32            = (*uint32)(&dword_5d4594_816368)
		dword_5d4594_816372      *uint32            = (*uint32)(&dword_5d4594_816372)
		dword_5d4594_816376      *ail.Driver        = (*ail.Driver)(unsafe.Pointer(&dword_5d4594_816376))
		dword_587000_93156       *uint32            = (*uint32)(&dword_587000_93156)
		dword_587000_93160       *uint32            = (*uint32)(&dword_587000_93160)
		dword_5d4594_816340      *uint32            = (*uint32)(&dword_5d4594_816340)
		counter_5d4594_816244    *timer.TimerGroup  = memmap.PtrT[timer.TimerGroup](0x5d4594, 816244)
		ptr_counter_587000_81128 **timer.TimerGroup = (**timer.TimerGroup)(unsafe.Pointer(&legacyGlobals.dword_587000_81128))
		dword_5d4594_816348      *uint32            = (*uint32)(&dword_5d4594_816348)
	)
	MusicModule = music.NewModule(
		"music",
		dword_5d4594_816368,
		dword_5d4594_816372,
		dword_5d4594_816376,
		dword_587000_93156,
		dword_587000_93160,
		dword_5d4594_816340,
		counter_5d4594_816244,
		ptr_counter_587000_81128,
		dword_5d4594_816348,
		Sub_43F130,
		checkDialogs,
		Sub_413890,
		PlatformTicks,
	)
}

// Stop playing music
func Sub_43D990() {
	MusicModule.SetNextMusic(music.MusicState{})
}

func checkDialogs() bool {
	return Dialogs.IsFallbackMode() != 0 && Dialogs.Sub_44D930()
}

func Sub_43D9B0(a1, a2 uint32) {
	MusicModule.SetNextMusic(music.MusicState{MusicIdx: a1, Volume: a2, Position: 0, D: 0})
}
