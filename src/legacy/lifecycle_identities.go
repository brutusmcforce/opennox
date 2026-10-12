package legacy

import "unsafe"

const (
	createIDMonster = iota
	createIDArmor
	createIDWeapon
	createIDObelisk
	createIDAnim
	createIDTrigger
	createIDMonsterGenerator
	createIDRewardMarker
)

const (
	initIDMonster = iota
	initIDPlayer
	initIDSpark
	initIDFrog
	initIDChest
	initIDBoulder
	initIDBreak
	initIDMonsterGenerator
	initIDSkull
	initIDDirection
	initIDGold
)

var (
	objectCreateIdentitySlots [8]struct{ slot byte }
	objectInitIdentitySlots   [11]struct{ slot byte }
)

func lifecycleCreateKey(id int) unsafe.Pointer { return unsafe.Pointer(&objectCreateIdentitySlots[id]) }
func lifecycleInitKey(id int) unsafe.Pointer   { return unsafe.Pointer(&objectInitIdentitySlots[id]) }
