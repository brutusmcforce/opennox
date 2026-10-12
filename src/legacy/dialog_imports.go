package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"

	"github.com/brutusmcforce/opennox/legacy/dialog"
	"github.com/brutusmcforce/opennox/legacy/timer"
)

func Set_dword_5d4594_831080(v uint32) {
	*memmap.PtrUint32(0x5D4594, 831080) = v
}

var (
	Dialogs *dialog.Dialog
)

func initDialog() {
	Dialogs = dialog.NewDialog(
		"dialog",
		(*uint32)(&dword_587000_122848),
		(*uint32)(&dword_5d4594_830864), (*uint32)(&dword_5d4594_831084), (*uint32)(&dword_587000_122856), (*uint32)(&dword_5d4594_831076), (*ail.Driver)(unsafe.Pointer(&dword_5d4594_831092)),
		memmap.PtrUint32(0x5D4594, 830860),
		func() *strman.StringManager {
			return GetServer().S().Strings()
		},
		memmap.PtrT[timer.TimerGroup](0x5D4594, 830876),
		memmap.PtrT[timer.TimerGroup](0x5D4594, 830980),
		(*timer.TimerGroup)(legacyGlobals.dword_587000_81128),
		(*timer.TimerGroup)(legacyGlobals.dword_587000_122852),
		Sub_413890,

		func() ail.Driver {
			return Sub_43F130()
		},
		func() { MusicModule.Sub_43DBD0() },
		func() { MusicModule.Sub_43DBE0() },
		func() bool { return MusicModule.Sub_43DC40() },
		Set_dword_5d4594_831080,
		func() uint32 { return uint32(Get_dword_587000_93160()) },
	)
}
