package legacy

import "unsafe"

func nox_xxx_plrLoad_41A480(path *int8) int32 {
	return int32(playerFileClientLoad(GoStringP(unsafe.Pointer(path))))
}
