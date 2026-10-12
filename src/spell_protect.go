package opennox

import (
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/server"
)

func castProtectElectricity(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_PROTECT_FROM_ELECTRICITY, lvl, args.Obj, spellBuffConf{
		DurOpt: "ProtectElectricityEnchantDuration",
	})
}

func castProtectPoison(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_PROTECT_FROM_POISON, lvl, args.Obj, spellBuffConf{
		DurOpt: "ProtectPoisonEnchantDuration",
	})
}

func castProtectFire(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_PROTECT_FROM_FIRE, lvl, args.Obj, spellBuffConf{
		DurOpt: "ProtectFireEnchantDuration",
	})
}
