package legacy

import "unsafe"

type damageIdentityID uint8

const (
	damageIDDefault damageIdentityID = iota
	damageIDSkeleton
	damageIDPlayer
	damageIDStone
	damageIDMechGolem
	damageIDFlammable
	damageIDBlackPowder
	damageIDArmor
	damageIDWeapon
	damageIDBall
	damageIDMonsterGenerator
	damageIDDefaultSound
	damageIDPlayerSound
	damageIdentityCount
)

// Each identity is a distinct address in Go static storage. The integer slot is
// intentionally not interpreted as a callback; server registries dispatch it
// through their existing typed Go maps.
var damageIdentitySlots [damageIdentityCount]uintptr

func damageIdentityKey(id damageIdentityID) unsafe.Pointer {
	return unsafe.Pointer(&damageIdentitySlots[id])
}
