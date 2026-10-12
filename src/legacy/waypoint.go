package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	_ = [1]struct{}{}[516-unsafe.Sizeof(server.Waypoint{})]
)
