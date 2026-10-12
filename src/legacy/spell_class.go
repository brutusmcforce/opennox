package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/things"
)

func playerSpellClassCheck(class, ind int32) int32 {
	flags := GetServer().S().Spells.Flags(spell.ID(ind))
	// Keep the raw int: narrowing to the byte-sized class enum accepts invalid inputs.
	var allowed things.SpellFlags
	switch class {
	case 1:
		allowed = things.SpellClassWizard
	case 2:
		allowed = things.SpellClassConjurer
	default:
		return 9
	}
	if flags&(things.SpellClassAny|allowed) != 0 {
		return 0
	}
	return 9
}
