package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func browserString(key string) string  { return serverPanelsText("noxworld.c", key) }
func browserHidden(w *gui.Window) bool { return w == nil || w.GetFlags().IsHidden() }
func browserTick() int {
	switch uint32(browserUI.connectionState) {
	case 0:
		if (!browserHidden(browserWindow(uint32(uintptr(browserUI.mapWindow)))) || !browserHidden(browserWindow(uint32(browserUI.overview))) || !browserHidden(browserWindow(uint32(browserUI.gameList)))) && browserUI.creating == 0 && browserUI.transition == 0 && browserUI.hosting == 0 && browserHidden(browserWindow(uint32(uintptr(browserUI.detailPanel)))) {
			if uint64(uint32(PlatformTicks())) > uint64(browserUI.refreshDeadline) {
				Nox_client_refreshServerList_4378B0()
			} else {
				Sub_438770_waitList()
			}
		}
	case 2:
		browserConnectionError()
		browserUI.connectionState = 1
	case 3:
		if uint64(uint32(PlatformTicks())) >= uint64(browserUI.connectionDeadline) {
			browserUI.connectionError = 8
			browserUI.connectionState = 2
		}
	case 4:
		browserUI.connectionState = 3
		Sub_449E30(browserString("TestCon"))
		// The original unsigned-int addition wraps before widening to the deadline.
		browserUI.connectionDeadline = uint64(uint32(PlatformTicks()) + 20000)
	case 5:
		Sub_449E00(browserString("Password"))
		Sub_449E30(browserString("PasswordRequired"))
		Sub_449EA0(7)
		Sub_44A360(0)
		browserUI.connectionState = 6
		Sub_4A24C0(1)
	case 7:
		Sub_44A360(1)
		Sub_449E30(browserString("Connected"))
		Sub_449EA0(0)
		GetClient().SetDrawFunc(func() bool { return Nox_xxx_cliDrawConnectedLoop_43B360() != 0 })
		browserUI.connectionState = 1
	case 8:
		browserUI.connectionState = 9
		// This path widens the clock before adding, unlike the test-connection path.
		*memmap.PtrUint64(0x5D4594, 814972) = uint64(uint32(PlatformTicks())) + 1000
	case 9:
		if uint64(uint32(PlatformTicks())) > memmap.Uint64(0x5D4594, 814972) {
			Nox_client_joinGame_438A90()
		}
	case 10:
		// The legacy child lookup has a nil parent and always returns nil.
		Sub_449E60(4)
	}
	return 1
}
func browserConnectionError() int {
	if Sub_44A4A0() == 0 {
		Nox_xxx_dialogMsgBoxCreate_449A10(browserUI.world, "", "", 0, nil, nil)
	}
	code := uint32(browserUI.connectionError)
	if code != 8 && code != 9 && code != 10 {
		Sub_449E00(browserString("ConnError"))
		code = uint32(browserUI.connectionError)
	}
	key := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 87416+uintptr(code)*4)))
	Sub_449E30(browserString(key))
	browserUI.transition = 0
	Sub_449EA0(1)
	Sub_44A360(1)
	return Sub_4A24C0(1)
}
func browserAttemptConnect() {
	browserWindow(uint32(uintptr(browserUI.detailPanel))).StackPop()
	Nox_xxx_dialogMsgBoxCreate_449A10(browserUI.world, "", browserString("AttemptingConn"), 34, nil, nil)
}
func browserNotice(timeout bool) {
	pending := &browserUI.pendingKicked
	key := "Kicked"
	if timeout {
		pending = &browserUI.pendingTimeout
		key = "Timeout"
	}
	if browserUI.world == nil {
		*pending = 1
		return
	}
	Nox_xxx_dialogMsgBoxCreate_449A10(nil, browserString("Notification"), browserString(key), 33, nil, nil)
	Sub_44A360(1)
	*pending = 0
}

func sub_438770() int32 { return int32(browserTick()) }

func sub_43B630() *uint32 { browserAttemptConnect(); return nil }
