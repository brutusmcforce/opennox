package opennox

import (
	"github.com/brutusmcforce/opennox/legacy/timer"
)

func init() {
	timer.PlatformTicks = platformTicks
}
