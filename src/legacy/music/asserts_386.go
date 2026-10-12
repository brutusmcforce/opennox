//go:build 386

package music

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"unsafe"
)

var _ = [1]struct{}{}[16-unsafe.Sizeof(MusicState{})]
var _ = [1]struct{}{}[20-unsafe.Sizeof(MusicBlock{})]
