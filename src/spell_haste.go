package opennox

import (
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/server"
)

func castHaste(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_HASTED, lvl, args.Obj, spellBuffConf{
		DurOpt: "HasteEnchantDuration",
	})
}
