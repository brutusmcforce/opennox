package legacy

import (
	"github.com/brutusmcforce/opennox/server"
)

func waypointEnabledMask(wp *server.Waypoint, mask byte) bool {
	return wp.IsEnabled() && wp.HasFlag2Mask(mask)
}
