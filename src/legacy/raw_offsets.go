package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
)

// rawAt returns the T at the 386 offset off from p. It is the body of the raw offset helpers (equipmentWord,
// spellLifeWord, ...), which are only valid for memory that keeps the 386 layout: blobs and RawData blocks.
// With NOX_LAYOUTCHECK, an access into a Go struct whose field at off moved on this architecture is reported
// with the helper's call site.
func rawAt[T any](p unsafe.Pointer, off int) *T {
	if layout.Checking {
		var zero T
		layout.CheckRaw(p, off, unsafe.Sizeof(zero), 2) // skip rawAt and the helper
	}
	return (*T)(unsafe.Add(p, off))
}
