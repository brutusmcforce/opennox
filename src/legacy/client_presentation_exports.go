package legacy

import "unsafe"

func nox_xxx_clientAddRayEffect_49C160(event int32) {
	presentationRayAdd((*[7]byte)(unsafe.Pointer(uintptr(uint32(event)))))
}

func nox_xxx_clientRemoveRayEffect_49C450(event int32) {
	presentationRayRemove((*[7]byte)(unsafe.Pointer(uintptr(uint32(event)))))
}

func nox_xxx_fxDrawTurnUndead_499880(pos *int16) {
	presentationTurnUndead((*[2]int16)(unsafe.Pointer(pos)))
}
