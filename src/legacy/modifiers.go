package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Sub_4A5E90_A                  func()
	Nox_xxx_fireRingEffect_4E05B0 func(a1 unsafe.Pointer, a2p, a3p, a4p *server.Object)
	Nox_xxx_blueFREffect_4E05F0   func(a1 unsafe.Pointer, a2p, a3p, a4p *server.Object)
)

var _ = [1]struct{}{}[88-unsafe.Sizeof(server.Modifier{})]

var _ = [1]struct{}{}[144-unsafe.Sizeof(server.ModifierEff{})]

func init() {
	registerNativeModifierCallbacks()
}

func nox_xxx_modifGetDescById_413330(a1 int32) unsafe.Pointer {
	return GetServer().S().Modif.Nox_xxx_modifGetDescById413330(int(a1)).C()
}

func nox_xxx_modifGetIdByName_413290(name *int8) int32 {
	return int32(GetServer().S().Modif.Nox_xxx_modifGetIdByName413290(GoStringP(unsafe.Pointer(name))))
}
