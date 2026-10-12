package legacy

/*
#include <stdlib.h>
#include "common__crypt.h"
*/
import "C"

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/cryptfile"
)

// The size of the write buffer for the C code (see nox_crypt_wbuf in common__crypt.h), and the number of fields in it.
const (
	cryptWBufSize   = 64 << 10
	cryptWBufFields = 16 << 10
)

func init() {
	cryptfile.Sync = cryptWBufDrain
	cryptfile.OnGlobal = func(cf *cryptfile.CryptFile) {
		b := &C.nox_crypt_wbuf
		if cf == nil || cf.ReadOnly() {
			b.on = 0
			return
		}
		if b.data == nil {
			b.data = (*C.uchar)(C.malloc(cryptWBufSize))
			b.cap = cryptWBufSize
			b.sizes = (*C.uint)(C.malloc(4 * cryptWBufFields))
			b.capsizes = cryptWBufFields
		}
		b.on = 1
	}
}

// cryptWBufDrain writes the data buffered by the C code to the global crypt file.
func cryptWBufDrain() {
	b := &C.nox_crypt_wbuf
	n, nf := int(b.len), int(b.nsizes)
	if nf == 0 {
		return
	}
	b.len, b.nsizes = 0, 0 // first, since writing syncs again
	if f := cryptfile.Global(); f != nil {
		data := unsafe.Slice((*byte)(unsafe.Pointer(b.data)), n)
		sizes := unsafe.Slice((*uint32)(unsafe.Pointer(b.sizes)), nf)
		_ = f.WriteFields(data, sizes)
	}
}
