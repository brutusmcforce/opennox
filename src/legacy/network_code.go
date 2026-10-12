package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
)

func drawableUnitCode(dr *client.Drawable) uint32 {
	if dr == nil {
		return 0
	}
	code := dr.NetCode32
	if code >= 0x8000 {
		return 0
	}
	if uint32(dr.ObjClass)&0x20400000 != 0 {
		code |= 0x8000
	}
	return code
}

func nox_xxx_netGetUnitCodeCli_578B00(a1 int32) uint32 {
	return uint32(drawableUnitCode((*client.Drawable)(unsafe.Pointer(uintptr(uint32(a1))))))
}

func nox_xxx_netClearHighBit_578B30(a1 int16) int32 {
	return int32(uint16(a1) & 0x7fff)
}

func nox_xxx_netTestHighBit_578B70(a1 uint32) uint32 {
	return (a1 >> 15) & 1
}
