package legacy

import "github.com/brutusmcforce/opennox/libs/platform"

var (
	PlatformTicks          func() uint64
	Nox_ticks_reset_416D40 func()
)

func nox_platform_rand() int {
	// The remaining CRT-rand consumers scale by 0x7fff. Keep this ABI at the
	// original 15-bit range even when the Go platform supplies a wider int.
	return platform.RandInt() & 0x7fff
}
