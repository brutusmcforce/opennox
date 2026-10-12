//go:build 386

package timer

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"unsafe"
)

var _ = [1]struct{}{}[32-unsafe.Sizeof(Timer{})]
var _ = [1]struct{}{}[96-unsafe.Sizeof(TimerGroup{})]
