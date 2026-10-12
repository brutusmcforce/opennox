package opennox

import (
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	"github.com/brutusmcforce/opennox/libs/spell"

	"github.com/brutusmcforce/opennox/server"
)

func castInvisibility(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_INVISIBLE, lvl, args.Obj, spellBuffConf{
		DurOpt: "InvisibilityEnchantDuration", Defensive: true, PointFX: netmsg.MSG_FX_SMOKE_BLAST,
	})
}
