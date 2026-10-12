package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/binfile"
)

func asMemfileP(p unsafe.Pointer) *binfile.MemFile {
	return (*binfile.MemFile)(p)
}
