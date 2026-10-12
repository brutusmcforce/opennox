package legacy

import (
	"github.com/brutusmcforce/opennox/server"
)

func appendWaypointLink(source, target *server.Waypoint, kind int8) bool {
	count := source.PointsCnt
	// C deliberately leaves the last of the 32 physical slots unused.
	if count >= 31 || source == target {
		return false
	}
	for i := byte(0); i < count; i++ {
		p := &source.Points[i]
		// C promotes stored uint8 and incoming signed char to int. Kinds with
		// the high bit set therefore do not match an existing stored byte.
		if p.Waypoint == target && int(p.Ind) == int(kind) {
			return false
		}
	}
	source.Points[count].Waypoint = target
	source.Points[count].Ind = byte(kind)
	source.PointsCnt = count + 1
	return true
}
