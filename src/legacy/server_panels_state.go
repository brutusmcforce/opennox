package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

// Private panel owners; addresses are only used to map the legacy layout.
var serverPanelsPrivate [26]uint32

func serverPanelsWord(off uintptr) *uint32 {
	switch off {
	case 1045460:
		return &serverPanelsPrivate[0]
	case 1045464:
		return &serverPanelsPrivate[1]
	case 1045468:
		return &serverPanelsPrivate[2]
	case 1045480:
		return &serverPanelsPrivate[3]
	case 1045484:
		return &serverPanelsPrivate[4]
	case 1045508:
		return &serverPanelsPrivate[5]
	case 1045516:
		return &serverPanelsPrivate[6]
	case 1045520:
		return &serverPanelsPrivate[7]
	case 1045528:
		return &serverPanelsPrivate[8]
	case 1045532:
		return &serverPanelsPrivate[9]
	case 1045536:
		return &serverPanelsPrivate[10]
	case 1045540:
		return &serverPanelsPrivate[11]
	case 1045544:
		return &serverPanelsPrivate[12]
	case 1045548:
		return &serverPanelsPrivate[13]
	case 1045552:
		return &serverPanelsPrivate[14]
	case 1045556:
		return &serverPanelsPrivate[15]
	case 1045576:
		return &serverPanelsPrivate[16]
	case 1045580:
		return &serverPanelsPrivate[17]
	case 1045584:
		return &serverPanelsPrivate[18]
	case 1045588:
		return &serverPanelsPrivate[19]
	case 1045596:
		return &serverPanelsPrivate[20]
	case 1309812:
		return &serverPanelsPrivate[21]
	case 1316704:
		return &serverPanelsPrivate[22]
	case 1316708:
		return &serverPanelsPrivate[23]
	case 1316712:
		return &serverPanelsPrivate[24]
	case 1316972:
		return &serverPanelsPrivate[25]
	default:
		return memmap.PtrUint32(0x5D4594, off)
	}
}
func serverPanelsWindow(off uintptr) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*serverPanelsWord(off))))
}
