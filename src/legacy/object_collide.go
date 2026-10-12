package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_collideDeathBall_4E9E90 func(a1, a2 *server.Object, pos *types.Pointf)
	Nox_xxx_castCounterSpell_52BBB0 func(cspl spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int
	Nox_xxx_changeOwner_52BE40      func(a1, a2 *server.Object)
)

func init() {
	server.RegisterObjectCollideNative("DefaultCollide", collisionKey(collisionIdentityDefault), func(u *server.Object, a2, a3 uintptr) uint32 {
		return 0
	}, 0)
	server.RegisterObjectCollideNative("MonsterCollide", collisionKey(collisionIdentityMonster), func(u *server.Object, a2, a3 uintptr) uint32 {
		return uint32(uintptr(stateMonsterCollision(u, objectFromWord(uint32(a2)))))
	}, 0)
	server.RegisterObjectCollideNative("PlayerCollide", collisionKey(collisionIdentityPlayer), func(u *server.Object, a2, a3 uintptr) uint32 {
		statePlayerCollision(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("ProjectileCollide", collisionKey(collisionIdentityProjectile), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileGeneric(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("ProjectileSparkCollide", collisionKey(collisionIdentityProjectileSpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileGenericSpark(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("DoorCollide", collisionKey(collisionIdentityDoor), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideDoor(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("PickupCollide", collisionKey(collisionIdentityPickup), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollidePickup(u, objectFromWord(uint32(a2)))
	}, 0)
	server.RegisterObjectCollideNative("ExitCollide", collisionKey(collisionIdentityExit), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideExit(u, objectFromWord(uint32(a2)))
		return 0
	}, 88)
	server.RegisterObjectCollideNative("DamageCollide", collisionKey(collisionIdentityDamage), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDamageField(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("ManaDrainCollide", collisionKey(collisionIdentityManaDrain), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileManaDrain(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("BombCollide", collisionKey(collisionIdentityBomb), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileBomb(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("SparkExplosionCollide", collisionKey(collisionIdentitySparkExplosion), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileFireball(u, objectFromWord(uint32(a2)))
		return 0
	}, 1)
	server.RegisterObjectCollideNative("ChestCollide", collisionKey(collisionIdentityChest), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideChest(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("WallReflectCollide", collisionKey(collisionIdentityWallReflect), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSulphur(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("WallReflectSparkCollide", collisionKey(collisionIdentityWallReflectSpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileWallSpark(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("PixieCollide", collisionKey(collisionIdentityPixie), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectilePixie(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("OwnCollide", collisionKey(collisionIdentityOwn), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSparkOwner(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("SparkCollide", collisionKey(collisionIdentitySpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSpark(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("BarrelCollide", collisionKey(collisionIdentityBarrel), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideBarrel(u)
		return 0
	}, 0)
	server.RegisterObjectCollideNative("AudioEventCollide", collisionKey(collisionIdentityAudioEvent), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideAudio(u, objectFromWord(uint32(a2)))
		return 0
	}, 4)
	server.RegisterObjectCollideNative("TriggerCollide", collisionKey(collisionIdentityTrigger), func(u *server.Object, a2, a3 uintptr) uint32 {
		motionTrigger(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("TeleportCollide", collisionKey(collisionIdentityTeleport), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideTeleport(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("ElevatorCollide", collisionKey(collisionIdentityDefault), func(u *server.Object, a2, a3 uintptr) uint32 {
		return 0
	}, 8)
	server.RegisterObjectCollideNative("AwardSpellCollide", collisionKey(collisionIdentityAwardSpell), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollideSpellAward(u, objectFromWord(uint32(a2)))
	}, 4)
	server.RegisterObjectCollideNative("DieCollide", collisionKey(collisionIdentityDie), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDie(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("GlyphCollide", collisionKey(collisionIdentityGlyph), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideGlyph_4E9A00(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("SpellProjectileCollide", collisionKey(collisionIdentitySpellProjectile), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSpellProjectile(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("BoomCollide", collisionKey(collisionIdentityBoom), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileBoom(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("SignCollide", collisionKey(collisionIdentitySign), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSign(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("PentagramCollide", collisionKey(collisionIdentityPentagram), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollidePentagram(u)
	}, 0)
	server.RegisterObjectCollideNative("SpiderSpitCollide", collisionKey(collisionIdentitySpiderSpit), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileWeb(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("DeathBallCollide", collisionKey(collisionIdentityDeathBall), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideDeathBall_4E9E90(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("DeathBallFragmentCollide", collisionKey(collisionIdentityDeathBallFragment), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDeathFragment(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("TelekinesisCollide", collisionKey(collisionIdentityDefault), func(u *server.Object, a2, a3 uintptr) uint32 {
		return 0
	}, 0)
	server.RegisterObjectCollideNative("FistCollide", collisionKey(collisionIdentityFist), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileFist(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("TeleportWakeCollide", collisionKey(collisionIdentityTeleportWake), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTeleportWake(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("FlagCollide", collisionKey(collisionIdentityFlag), func(u *server.Object, a2, a3 uintptr) uint32 {
		objectiveFlagCollide(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("ChakramInMotionCollide", collisionKey(collisionIdentityChakramInMotion), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileChakram(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("ArrowCollide", collisionKey(collisionIdentityArrow), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileArrow(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("MonsterArrowCollide", collisionKey(collisionIdentityMonsterArrow), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileMonsterArrow(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("BearTrapCollide", collisionKey(collisionIdentityBearTrap), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTrap(u, objectFromWord(uint32(a2)), false)
		return 0
	}, 0)
	server.RegisterObjectCollideNative("PoisonGasTrapCollide", collisionKey(collisionIdentityPoisonGasTrap), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTrap(u, objectFromWord(uint32(a2)), true)
		return 0
	}, 0)
	server.RegisterObjectCollideNative("TrapDoorCollide", collisionKey(collisionIdentityTrapDoor), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideTrap(u, objectFromWord(uint32(a2)))
		return 0
	}, 28)
	server.RegisterObjectCollideNative("BallCollide", collisionKey(collisionIdentityBall), func(u *server.Object, a2, a3 uintptr) uint32 {
		objectiveBallCollide(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("HomeBaseCollide", collisionKey(collisionIdentityHomeBase), func(u *server.Object, a2, a3 uintptr) uint32 {
		return uint32(int32(objectiveHomeBase(u, objectFromWord(uint32(a2)))))
	}, 0)
	server.RegisterObjectCollideNative("CrownCollide", collisionKey(collisionIdentityCrown), func(u *server.Object, a2, a3 uintptr) uint32 {
		return objectiveCrownCollide(u, objectFromWord(uint32(a2)))
	}, 0)
	server.RegisterObjectCollideNative("UndeadKillerCollide", collisionKey(collisionIdentityUndeadKiller), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideUndead(u, objectFromWord(uint32(a2)), int32(a3) != 0)
		return 0
	}, 4)
	server.RegisterObjectCollideNative("YellowStarShotCollide", collisionKey(collisionIdentityYellowStarShot), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSulphurTrail(u, objectFromWord(uint32(a2)), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("MimicCollide", collisionKey(collisionIdentityMimic), func(u *server.Object, a2, a3 uintptr) uint32 {
		return uint32(uintptr(stateMimicCollision(u, objectFromWord(uint32(a2)))))
	}, 0)
	server.RegisterObjectCollideNative("HarpoonCollide", collisionKey(collisionIdentityHarpoon), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideHarpoon_4EB6A0(u, objectFromWord(uint32(a2)))
		return 0
	}, 8)
	server.RegisterObjectCollideNative("MonsterGeneratorCollide", collisionKey(collisionIdentityMonsterGenerator), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideGenerator(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)
	server.RegisterObjectCollideNative("SoulGateCollide", collisionKey(collisionIdentitySoulGate), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSoulGate(u, objectFromWord(uint32(a2)))
		return 0
	}, 4)
	server.RegisterObjectCollideNative("AnkhCollide", collisionKey(collisionIdentityAnkh), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideAnkh(u, objectFromWord(uint32(a2)))
		return 0
	}, 0)

	server.RegisterObjectCollideParse("ProjectileCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("ProjectileSparkCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("DamageCollide", resourceObjectParser("collide", "damage"))
	server.RegisterObjectCollideParse("ManaDrainCollide", resourceObjectParser("collide", "mana"))
	server.RegisterObjectCollideParse("SparkExplosionCollide", resourceObjectParser("collide", "spark"))
	server.RegisterObjectCollideParse("WallReflectCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("WallReflectSparkCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("PixieCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("AudioEventCollide", resourceObjectParser("collide", "audio"))
	server.RegisterObjectCollideParse("MonsterArrowCollide", resourceObjectParser("collide", "arrow"))
	server.RegisterObjectCollideParse("YellowStarShotCollide", resourceObjectParser("collide", "projectile"))
}

func nox_xxx_castCounterSpell_52BBB0(a1 int32, a2, a3, a4 *nox_object_t) {
	Nox_xxx_castCounterSpell_52BBB0(spell.ID(a1), asObjectS(a2), asObjectS(a3), asObjectS(a4), nil, 0)
}
