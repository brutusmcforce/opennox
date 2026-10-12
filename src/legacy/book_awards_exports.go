package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func nox_xxx_guide_427010(name *int8) int32 {
	return int32(bookGuideID(alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

func nox_xxx_guiCreatureGetName_427240(id int32) int32 {
	return int32(bookGuideCreatureName(int32(id)))
}
