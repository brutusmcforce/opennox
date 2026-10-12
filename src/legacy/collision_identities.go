package legacy

import "unsafe"

// These nonzero-sized linker globals give every callback a distinct, stable
// identity. They contain no Go pointers and remain pinned for process lifetime.
// They are data addresses, never executable C callbacks.
var collisionIdentitySlots [51]struct{ slot byte }

const (
	collisionIdentityDefault           = 0
	collisionIdentityMonster           = 1
	collisionIdentityPlayer            = 2
	collisionIdentityProjectile        = 3
	collisionIdentityProjectileSpark   = 4
	collisionIdentityDoor              = 5
	collisionIdentityPickup            = 6
	collisionIdentityExit              = 7
	collisionIdentityDamage            = 8
	collisionIdentityManaDrain         = 9
	collisionIdentityBomb              = 10
	collisionIdentitySparkExplosion    = 11
	collisionIdentityChest             = 12
	collisionIdentityWallReflect       = 13
	collisionIdentityWallReflectSpark  = 14
	collisionIdentityPixie             = 15
	collisionIdentityOwn               = 16
	collisionIdentitySpark             = 17
	collisionIdentityBarrel            = 18
	collisionIdentityAudioEvent        = 19
	collisionIdentityTrigger           = 20
	collisionIdentityTeleport          = 21
	collisionIdentityAwardSpell        = 22
	collisionIdentityDie               = 23
	collisionIdentityGlyph             = 24
	collisionIdentitySpellProjectile   = 25
	collisionIdentityBoom              = 26
	collisionIdentitySign              = 27
	collisionIdentityPentagram         = 28
	collisionIdentitySpiderSpit        = 29
	collisionIdentityDeathBall         = 30
	collisionIdentityDeathBallFragment = 31
	collisionIdentityFist              = 32
	collisionIdentityTeleportWake      = 33
	collisionIdentityFlag              = 34
	collisionIdentityChakramInMotion   = 35
	collisionIdentityArrow             = 36
	collisionIdentityMonsterArrow      = 37
	collisionIdentityBearTrap          = 38
	collisionIdentityPoisonGasTrap     = 39
	collisionIdentityTrapDoor          = 40
	collisionIdentityBall              = 41
	collisionIdentityHomeBase          = 42
	collisionIdentityCrown             = 43
	collisionIdentityUndeadKiller      = 44
	collisionIdentityYellowStarShot    = 45
	collisionIdentityMimic             = 46
	collisionIdentityHarpoon           = 47
	collisionIdentityMonsterGenerator  = 48
	collisionIdentitySoulGate          = 49
	collisionIdentityAnkh              = 50
)

func collisionKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&collisionIdentitySlots[id])
}
