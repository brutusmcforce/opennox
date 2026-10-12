package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

func waypointFromRaw(a1 int32) *server.Waypoint {
	return (*server.Waypoint)(unsafe.Pointer(uintptr(uint32(a1))))
}

func waypointRaw(wp *server.Waypoint) int32 {
	return int32(uint32(uintptr(unsafe.Pointer(wp))))
}

func waypointEnabledMask(wp *server.Waypoint, mask byte) bool {
	return wp.IsEnabled() && wp.HasFlag2Mask(mask)
}
