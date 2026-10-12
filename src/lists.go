package opennox

import (
	"github.com/brutusmcforce/opennox/legacy"
	"unsafe"
)

func nox_common_list_getFirstSafe_425890(p unsafe.Pointer) unsafe.Pointer { return legacy.ListFirst(p) }
func nox_common_list_getNextSafe_4258A0(p unsafe.Pointer) unsafe.Pointer  { return legacy.ListNext(p) }
func nox_common_list_getNext_425940(p unsafe.Pointer) unsafe.Pointer      { return legacy.ListNext(p) }
func nox_common_list_clear_425760(p unsafe.Pointer)                       { legacy.ListClear(p) }
func nox_common_list_append_4258E0(p, q unsafe.Pointer)                   { legacy.ListAppend(p, q) }
