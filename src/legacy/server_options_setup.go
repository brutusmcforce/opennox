package legacy

import (
	"encoding/binary"
	"strconv"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func serverOptionsHideRange(w *gui.Window, first, last uint, hidden bool) {
	for id := first; id <= last; id++ {
		serverOptionsHidden(w.ChildByID(id), bool2int(hidden))
	}
}
func serverOptionsEnableRange(w *gui.Window, first, last uint, enabled bool) {
	for id := first; id <= last; id++ {
		uiWindowEnable(w.ChildByID(id), bool2int(enabled))
	}
}
func serverOptionsSetup(data []byte) int {
	serverOptionsSetText(serverOptionsWindow(1046512), 16414, GoStringP(unsafe.Pointer(serverConfigNameGet())), 0)
	serverOptionsTeamCount()
	id := uint(10119)
	if noxflags.HasGame(128) {
		id = 10122
	}
	serverOptionsChild(id).SetHidden(true)
	root, teams := serverOptionsWindow(1046492), serverOptionsWindow(1046508)
	if noxflags.HasGame(1) {
		uiWindowEnable(teams, 1)
		if noxflags.HasGame(128) {
			initial := serverOptionsRecord(unsafe.Pointer(serverConfigSlot(int32(1))))
			mode := binary.LittleEndian.Uint16(initial[52:]) & 0x17f0
			serverOptionsMapList(int(mode), alloc.GoString(&initial[0]), true)
			data = serverOptionsCurrent()
			binary.LittleEndian.PutUint16(data[52:], mode|(binary.LittleEndian.Uint16(data[52:])&0xe80f))
			if teamUISettingsLocked() {
				uiWindowEnable(teams, 0)
			} else if noxflags.HasGamePlay(4) {
				teams.ChildByID(10330).DrawData().Field0 |= 4
			} else {
				serverOptionsEnableRange(teams, 10331, 10333, false)
			}
		} else {
			last := uint(10333)
			if noxflags.HasGamePlay(4) {
				last = 10331
			}
			serverOptionsEnableRange(teams, 10330, last, false)
			serverOptionsMapList(int(binary.LittleEndian.Uint16(data[52:])), alloc.GoString(&data[0]), false)
		}
		mode := binary.LittleEndian.Uint16(data[52:])
		serverOptionsSetText(serverOptionsWindow(1046516), 16414, strconv.Itoa(int(uint16(Nox_xxx_servGamedataGet_40A020(mode)))), 0)
		serverOptionsSetText(serverOptionsWindow(1046520), 16414, strconv.Itoa(int(byte(Sub_40A180(noxflags.GameFlag(mode))))), 0)
		serverOptionsLabels(data)
	} else {
		uiWindowEnable(root.ChildByID(10161), 0)
		uiWindowEnable(serverOptionsWindow(1046536), 0)
		serverOptionsWindow(1046536).Flags |= 8
		uiWindowEnable(serverOptionsWindow(1046504), 1)
		serverOptionsEnableRange(serverOptionsWindow(1046504), 10134, 10135, false)
		serverOptionsChild(10141).SetHidden(true)
		serverOptionsMapList(int(binary.LittleEndian.Uint16(data[52:])), alloc.GoString(&data[0]), false)
		serverOptionsLimits(data)
	}
	if noxflags.HasGamePlay(2) {
		teams.ChildByID(10331).DrawData().Field0 |= 4
	}
	if noxflags.HasGamePlay(1) {
		teams.ChildByID(10333).DrawData().Field0 |= 4
	}
	serverOptionsSetText(root.ChildByID(10119), 16385, serverOptionsModeName(binary.LittleEndian.Uint16(data[52:])), 0)
	serverOptionsSettingsLabels(data)
	return serverOptionsQuest(int((binary.LittleEndian.Uint16(data[52:]) >> 12) & 1))
}
