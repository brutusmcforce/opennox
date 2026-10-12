package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Sub_4FC670                         func(a1 int)
	Nox_xxx_playerExecuteAbil_4FBB70   func(cu *server.Object, a2 int)
	Sub_4FC0B0                         func(a1 *server.Object, a2 int)
	Nox_xxx_playerCancelAbils_4FC180   func(cu *server.Object)
	Sub_4FC300                         func(cu *server.Object, a2 int)
	Nox_xxx_abilityGetName_0_425260    func(ca int) string
	Nox_xxx_abilityCooldown_4252D0     func(ca int) int
	Sub_4252F0                         func(ca int) string
	Nox_xxx_spellGetAbilityIcon_425310 func(abil, icon int) noxrender.ImageHandle
	Nox_xxx_bookFirstKnownAbil_425330  func() int
	Nox_xxx_bookNextKnownAbil_425350   func(a1 int) int
	Sub_425450                         func(a1 int) int
	Nox_xxx_netAbilRepotState_4D8100   func(a1 *server.Object, a2 server.Ability, a3 byte)
)

func sub_4FC300(cu *nox_object_t, a2 int) { Sub_4FC300(asObjectS(cu), a2) }

func sub_4FC440(a1 *nox_object_t, a2 int) {
	GetServer().S().Abils.Sub4FC440(asObjectS(a1), server.Ability(a2))
}

func nox_xxx_abilityGetName_0_425260(ca int) *wchar2_t {
	return internWStr(Nox_xxx_abilityGetName_0_425260(ca))
}

func nox_common_playerIsAbilityActive_4FC250(a1 *nox_object_t, a2 int) int {
	return bool2int(GetServer().S().Abils.IsActive(asObjectS(a1), server.Ability(a2)))
}

func nox_xxx_probablyWarcryCheck_4FC3E0(a1 *nox_object_t, a2 int) int {
	return bool2int(GetServer().S().Abils.IsActiveVal(asObjectS(a1), server.Ability(a2)))
}

func nox_xxx_abilityCooldown_4252D0(ca int) int { return Nox_xxx_abilityCooldown_4252D0(ca) }

func sub_4252F0(ca int) *wchar2_t { return internWStr(Sub_4252F0(ca)) }

func nox_xxx_spellGetAbilityIcon_425310(abil, icon int) noxrender.ImageHandle {
	return Nox_xxx_spellGetAbilityIcon_425310(abil, icon)
}

func nox_xxx_bookFirstKnownAbil_425330() int { return Nox_xxx_bookFirstKnownAbil_425330() }

func nox_xxx_bookNextKnownAbil_425350(a1 int) int { return Nox_xxx_bookNextKnownAbil_425350(a1) }

func sub_425450(a1 int) int { return Sub_425450(a1) }
