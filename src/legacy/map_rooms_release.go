package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

func mapRoomBeforeRelease(p unsafe.Pointer) { memmap.UnregisterWordRegion(p) }
