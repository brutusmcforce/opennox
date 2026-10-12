package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"unsafe"
)

// Go owns the browser state. Fixed-width ABI types preserve the existing legacy
// callers and fixture word views while their C global symbols retire.
var browserUI = struct {
	listMode           uint32
	region             int32
	polygonsReady      uint32
	popup              uint32
	popupCount         uint32
	connectionState    uint32
	selected           unsafe.Pointer
	mapWindow          unsafe.Pointer
	overview           uint32
	filter             uint32
	label              unsafe.Pointer
	detailPanel        unsafe.Pointer
	detailList         *gui.Window
	playersColumn      uint32
	modeColumn         uint32
	mapColumn          uint32
	pingColumn         uint32
	statusColumn       uint32
	transition         uint32
	hosting            uint32
	hasSelection       uint32
	pendingKicked      uint32
	pendingTimeout     uint32
	retry              int32
	connectionError    uint32
	creating           uint32
	animation          *gui.Anim
	resultCount        uint32
	sort               uint32
	gameList           uint32
	world              *gui.Window
	connectionDeadline uint64
	refreshDeadline    uint64
	servers            *legacyListNode
	joinButton         *gui.Window // was 0x5D4594+815008
	popupWindow        *gui.Window // was 0x5D4594+815036
}{listMode: 1, region: -1, sort: 6}
