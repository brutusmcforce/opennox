package opennox

import (
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/server"
)

func castSlow(spellID spell.ID, _, a3, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_SLOWED, lvl, args.Obj, spellBuffConf{
		DurOpt: "SlowEnchantDuration", Orig: a3, Offensive: true,
	})
}
