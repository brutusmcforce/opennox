package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
)

var (
	Nox_xxx_spellCastByPlayer_4FEEF0 func()
)

func nox_xxx_spellCancelDurSpell_4FEB10(a1 int, a2 *nox_object_t) {
	GetServer().S().Spells.Dur.CancelFor(spell.ID(a1), asObjectS(a2))
}
