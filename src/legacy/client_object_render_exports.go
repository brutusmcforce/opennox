package legacy

import "unsafe"

func sub_4C5020(packet int32) int32 {
	return int32(objectRenderBeamAppend(unsafe.Pointer(uintptr(packet))))
}
