package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var _ = [1]struct{}{}[16-unsafe.Sizeof(server.MonsterAnim{})]

func nox_xxx_animPlayerGetFrameRange_4F9F90(a1 int, a2, a3 *int32) {
	v1, v2 := GetServer().S().PlayerAnimFrames(a1)
	if a2 != nil {
		*a2 = int32(v1)
	}
	if a3 != nil {
		*a3 = int32(v2)
	}
}
