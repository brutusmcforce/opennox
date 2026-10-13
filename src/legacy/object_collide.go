package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"

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
	}, server.NoData)
	server.RegisterObjectCollideNative("MonsterCollide", collisionKey(collisionIdentityMonster), func(u *server.Object, a2, a3 uintptr) uint32 {
		return memmap.PtrWord(stateMonsterCollision(u, objectArg(a2)))
	}, server.NoData)
	server.RegisterObjectCollideNative("PlayerCollide", collisionKey(collisionIdentityPlayer), func(u *server.Object, a2, a3 uintptr) uint32 {
		statePlayerCollision(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("ProjectileCollide", collisionKey(collisionIdentityProjectile), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileGeneric(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("ProjectileSparkCollide", collisionKey(collisionIdentityProjectileSpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileGenericSpark(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("DoorCollide", collisionKey(collisionIdentityDoor), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideDoor(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("PickupCollide", collisionKey(collisionIdentityPickup), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollidePickup(u, objectArg(a2))
	}, server.NoData)
	server.RegisterObjectCollideNative("ExitCollide", collisionKey(collisionIdentityExit), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideExit(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ExitCollideData]())
	server.RegisterObjectCollideNative("DamageCollide", collisionKey(collisionIdentityDamage), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDamageField(u, objectArg(a2))
		return 0
	}, server.DataOf[server.DamageCollideData]())
	server.RegisterObjectCollideNative("ManaDrainCollide", collisionKey(collisionIdentityManaDrain), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileManaDrain(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ManaDrainCollideData]())
	server.RegisterObjectCollideNative("BombCollide", collisionKey(collisionIdentityBomb), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileBomb(u, objectArg(a2))
		return 0
	}, server.DataOf[server.UnknownCollideData]())
	server.RegisterObjectCollideNative("SparkExplosionCollide", collisionKey(collisionIdentitySparkExplosion), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileFireball(u, objectArg(a2))
		return 0
	}, server.DataOf[server.SparkExplosionCollideData]())
	server.RegisterObjectCollideNative("ChestCollide", collisionKey(collisionIdentityChest), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideChest(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("WallReflectCollide", collisionKey(collisionIdentityWallReflect), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSulphur(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("WallReflectSparkCollide", collisionKey(collisionIdentityWallReflectSpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileWallSpark(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("PixieCollide", collisionKey(collisionIdentityPixie), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectilePixie(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("OwnCollide", collisionKey(collisionIdentityOwn), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSparkOwner(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("SparkCollide", collisionKey(collisionIdentitySpark), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSpark(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("BarrelCollide", collisionKey(collisionIdentityBarrel), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideBarrel(u)
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("AudioEventCollide", collisionKey(collisionIdentityAudioEvent), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideAudio(u, objectArg(a2))
		return 0
	}, server.DataOf[server.AudioEventCollideData]())
	server.RegisterObjectCollideNative("TriggerCollide", collisionKey(collisionIdentityTrigger), func(u *server.Object, a2, a3 uintptr) uint32 {
		motionTrigger(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("TeleportCollide", collisionKey(collisionIdentityTeleport), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideTeleport(u, objectArg(a2))
		return 0
	}, server.DataOf[server.TeleportCollideData]())
	server.RegisterObjectCollideNative("ElevatorCollide", collisionKey(collisionIdentityDefault), func(u *server.Object, a2, a3 uintptr) uint32 {
		return 0
	}, server.DataOf[server.UnknownCollideData]())
	server.RegisterObjectCollideNative("AwardSpellCollide", collisionKey(collisionIdentityAwardSpell), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollideSpellAward(u, objectArg(a2))
	}, server.DataOf[server.AwardSpellCollideData]())
	server.RegisterObjectCollideNative("DieCollide", collisionKey(collisionIdentityDie), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDie(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("GlyphCollide", collisionKey(collisionIdentityGlyph), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideGlyph_4E9A00(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("SpellProjectileCollide", collisionKey(collisionIdentitySpellProjectile), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSpellProjectile(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("BoomCollide", collisionKey(collisionIdentityBoom), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileBoom(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("SignCollide", collisionKey(collisionIdentitySign), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSign(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("PentagramCollide", collisionKey(collisionIdentityPentagram), func(u *server.Object, a2, a3 uintptr) uint32 {
		return worldCollidePentagram(u)
	}, server.NoData)
	server.RegisterObjectCollideNative("SpiderSpitCollide", collisionKey(collisionIdentitySpiderSpit), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileWeb(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("DeathBallCollide", collisionKey(collisionIdentityDeathBall), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideDeathBall_4E9E90(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("DeathBallFragmentCollide", collisionKey(collisionIdentityDeathBallFragment), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileDeathFragment(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("TelekinesisCollide", collisionKey(collisionIdentityDefault), func(u *server.Object, a2, a3 uintptr) uint32 {
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("FistCollide", collisionKey(collisionIdentityFist), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileFist(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("TeleportWakeCollide", collisionKey(collisionIdentityTeleportWake), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTeleportWake(u, objectArg(a2))
		return 0
	}, server.DataOf[server.TeleportWakeCollideData]())
	server.RegisterObjectCollideNative("FlagCollide", collisionKey(collisionIdentityFlag), func(u *server.Object, a2, a3 uintptr) uint32 {
		objectiveFlagCollide(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("ChakramInMotionCollide", collisionKey(collisionIdentityChakramInMotion), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileChakram(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("ArrowCollide", collisionKey(collisionIdentityArrow), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileArrow(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("MonsterArrowCollide", collisionKey(collisionIdentityMonsterArrow), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileMonsterArrow(u, objectArg(a2))
		return 0
	}, server.DataOf[server.MonsterArrowCollideData]())
	server.RegisterObjectCollideNative("BearTrapCollide", collisionKey(collisionIdentityBearTrap), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTrap(u, objectArg(a2), false)
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("PoisonGasTrapCollide", collisionKey(collisionIdentityPoisonGasTrap), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileTrap(u, objectArg(a2), true)
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("TrapDoorCollide", collisionKey(collisionIdentityTrapDoor), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideTrap(u, objectArg(a2))
		return 0
	}, server.DataOf[server.HoleCollideData]())
	server.RegisterObjectCollideNative("BallCollide", collisionKey(collisionIdentityBall), func(u *server.Object, a2, a3 uintptr) uint32 {
		objectiveBallCollide(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("HomeBaseCollide", collisionKey(collisionIdentityHomeBase), func(u *server.Object, a2, a3 uintptr) uint32 {
		return uint32(int32(objectiveHomeBase(u, objectArg(a2))))
	}, server.NoData)
	server.RegisterObjectCollideNative("CrownCollide", collisionKey(collisionIdentityCrown), func(u *server.Object, a2, a3 uintptr) uint32 {
		return objectiveCrownCollide(u, objectArg(a2))
	}, server.NoData)
	server.RegisterObjectCollideNative("UndeadKillerCollide", collisionKey(collisionIdentityUndeadKiller), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideUndead(u, objectArg(a2), int32(a3) != 0)
		return 0
	}, server.DataOf[server.UndeadKillerCollideData]())
	server.RegisterObjectCollideNative("YellowStarShotCollide", collisionKey(collisionIdentityYellowStarShot), func(u *server.Object, a2, a3 uintptr) uint32 {
		projectileSulphurTrail(u, objectArg(a2), (*types.Pointf)(unsafe.Pointer(a3)))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("MimicCollide", collisionKey(collisionIdentityMimic), func(u *server.Object, a2, a3 uintptr) uint32 {
		return memmap.PtrWord(stateMimicCollision(u, objectArg(a2)))
	}, server.NoData)
	server.RegisterObjectCollideNative("HarpoonCollide", collisionKey(collisionIdentityHarpoon), func(u *server.Object, a2, a3 uintptr) uint32 {
		Nox_xxx_collideHarpoon_4EB6A0(u, objectArg(a2))
		return 0
	}, server.DataOf[server.ProjectileCollideData]())
	server.RegisterObjectCollideNative("MonsterGeneratorCollide", collisionKey(collisionIdentityMonsterGenerator), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideGenerator(u, objectArg(a2))
		return 0
	}, server.NoData)
	server.RegisterObjectCollideNative("SoulGateCollide", collisionKey(collisionIdentitySoulGate), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideSoulGate(u, objectArg(a2))
		return 0
	}, server.DataOf[server.SoulGateCollideData]())
	server.RegisterObjectCollideNative("AnkhCollide", collisionKey(collisionIdentityAnkh), func(u *server.Object, a2, a3 uintptr) uint32 {
		worldCollideAnkh(u, objectArg(a2))
		return 0
	}, server.NoData)

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
