package opennox

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/server"
)

func castInfravision(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_INFRAVISION, lvl, args.Obj, spellBuffConf{
		DurOpt: "InfravisionEnchantDuration",
	})
}

func nox_xxx_warriorInfravis_540110(u *server.Object, dur int) {
	if u == nil {
		return
	}
	if !u.Class().Has(object.ClassPlayer) {
		return
	}
	s := noxServer
	if ud := u.UpdateDataPlayer(); ud != nil {
		asObjectS(u).ApplyEnchant(server.ENCHANT_INFRAVISION, dur, int(ud.Player.SpellLvl[server.AbilityInfravis]))
		s.abilities.netAbilReportActive(u, server.AbilityInfravis, true)
	}
}
