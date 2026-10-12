package opennox

import (
	"github.com/brutusmcforce/opennox/libs/object"

	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

func nox_xxx_warriorTreadLightly_5400B0(u *server.Object, dur int) {
	if u == nil {
		return
	}
	if !u.Class().Has(object.ClassPlayer) {
		return
	}
	s := noxServer
	if ud := u.UpdateDataPlayer(); ud != nil {
		if ud.State == server.PlayerState5 {
			nox_xxx_playerSetState_4FA020(u, server.PlayerState13)
		}
		asObjectS(u).ApplyEnchant(server.ENCHANT_SNEAK, dur, int(ud.Player.SpellLvl[server.AbilityTreadLightly]))
		s.abilities.netAbilReportActive(u, server.AbilityTreadLightly, true)
		legacy.Nox_xxx_frameCounterSetCopyToNextFrame_5281D0()
	}
}
