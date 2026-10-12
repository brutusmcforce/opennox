package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_monsterCreateFn_54C480 func(u *server.Object)
)

func Nox_xxx_monsterDefByTT_517560(typ int) *server.MonsterDef {
	return monsterDefinitionByType(uint32(typ))
}

func Nox_xxx_monsterAutoSpells_54C0C0(u *server.Object) {
	monsterAutoSpells(u)
}

func Nox_xxx_getDefaultSoundSet_424350(name string) unsafe.Pointer {
	return resourceSoundByName(name)
}
