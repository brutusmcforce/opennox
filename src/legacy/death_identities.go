package legacy

import "unsafe"

// These nonzero-sized linker globals provide unique, stable, non-executable
// identities for death callbacks stored in object/type records.
var deathIdentitySlots [14]struct{ slot byte }

const (
	deathIdentityPlayer           = 0
	deathIdentityPotion           = 1
	deathIdentityImpEgg           = 2
	deathIdentityGlyph            = 3
	deathIdentityBarrel           = 4
	deathIdentityCreateObject     = 5
	deathIdentitySpawnObject      = 6
	deathIdentityPolyp            = 7
	deathIdentityMarker           = 8
	deathIdentityWeapon           = 9
	deathIdentityArmor            = 10
	deathIdentityBoulder          = 11
	deathIdentityGameBall         = 12
	deathIdentityMonsterGenerator = 13
)

func deathKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&deathIdentitySlots[id])
}
