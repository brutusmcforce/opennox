package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"strings"
	"unsafe"
)

var serverOptionsRoot, serverOptionsMapsWindow, serverOptionsMapControls, serverOptionsLimitsPanel, serverOptionsTeamsPanel uint32
var serverOptionsNameControl, serverOptionsScoreControl, serverOptionsTimeControl, serverOptionsMainPanel uint32
var serverOptionsAccessPanel, serverOptionsPlayersPanel, serverOptionsAdvancedControl, serverOptionsGeneralPanel uint32
var serverOptionsTabs2, serverOptionsTabs3 uint32
var serverOptionsFirstOpen uint32 = 1

func serverOptionsWord(off uintptr) *uint32 {
	switch off {
	case 1046492:
		return &serverOptionsRoot
	case 1046496:
		return &serverOptionsMapsWindow
	case 1046500:
		return &serverOptionsMapControls
	case 1046504:
		return &serverOptionsLimitsPanel
	case 1046508:
		return &serverOptionsTeamsPanel
	case 1046512:
		return &serverOptionsNameControl
	case 1046516:
		return &serverOptionsScoreControl
	case 1046520:
		return &serverOptionsTimeControl
	case 1046524:
		return &serverOptionsMainPanel
	case 1046528:
		return &serverOptionsAccessPanel
	case 1046532:
		return &serverOptionsPlayersPanel
	case 1046536:
		return &serverOptionsAdvancedControl
	case 1046540:
		return &serverOptionsGeneralPanel
	case 1046356:
		return &serverOptionsTabs2
	case 1046360:
		return &serverOptionsTabs3
	default:
		return memmap.PtrUint32(0x5D4594, off)
	}
}
func serverOptionsWindow(off uintptr) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*serverOptionsWord(off))))
}
func serverOptionsChild(id uint) *gui.Window { return serverOptionsWindow(1046492).ChildByID(id) }
func serverOptionsText(key string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(key), "guiserv.c")
}
func serverOptionsPtr(w *gui.Window) uintptr { return uintptr(w.C()) }
func serverOptionsSetText(w *gui.Window, code int, text string, last int) int {
	return teamUIEvent(w, code, uintptr(unsafe.Pointer(alloc.InternCString16(text))), uintptr(last))
}
func serverOptionsGetText(w *gui.Window, code int, row int) string {
	return alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(teamUIEvent(w, code, uintptr(row), 0)))))
}
func serverOptionsRecord(p unsafe.Pointer) []byte { return unsafe.Slice((*byte)(p), 58) }
func serverOptionsCurrent() []byte                { return serverOptionsRecord(unsafe.Pointer(teamUISettings())) }
func serverOptionsDirty(v int) int                { *serverOptionsWord(1046544) = uint32(v); return v }
func serverOptionsFormat(key string, args ...any) string {
	return fmt.Sprintf(strings.ReplaceAll(serverOptionsText(key), "%S", "%s"), args...)
}
func serverOptionsStoreText(off uintptr, n int, text string) uintptr {
	p := memmap.PtrUint16(0x5D4594, off)
	alloc.StrCopy16(unsafe.Slice(p, n), text)
	return uintptr(unsafe.Pointer(p))
}
func serverOptionsHidden(w *gui.Window, hidden int) int {
	if w == nil {
		return -2
	}
	w.SetHidden(hidden != 0)
	return 0
}
func serverOptionsVisible(hidden int) int {
	if serverOptionsRoot == 0 {
		return 0
	}
	return serverOptionsHidden(serverOptionsWindow(1046492), hidden)
}
func serverOptionsListHead() unsafe.Pointer { return memmap.PtrOff(0x5D4594, 1045956) }
