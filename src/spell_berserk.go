package opennox

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

func nox_xxx_warriorBerserker_53FEB0(u *server.Object) {
	if u == nil {
		return
	}
	s := noxServer
	legacy.Nox_xxx_cancelAllSpells_4FEE90(u)
	if u.HasEnchant(server.ENCHANT_CONFUSED) {
		u.Direction2 = server.Dir16(byte(int(u.Direction1) + 4*s.Rand.Logic.IntClamp(-8, 8))) // wraps, as the original
	}
	if u.Class().Has(object.ClassPlayer) {
		if ud := u.UpdateDataPlayer(); ud != nil {
			ud.Field59_0 = 0
		}
	}
	nox_xxx_playerSetState_4FA020(u, server.PlayerState1)
	asObjectS(u).DisableEnchant(server.ENCHANT_INVISIBLE)
	asObjectS(u).DisableEnchant(server.ENCHANT_INVULNERABLE)
	s.Spells.Dur.CancelFor(spell.SPELL_OVAL_SHIELD, u)
	s.abilities.netAbilReportActive(u, server.AbilityBerserk, true)
}
