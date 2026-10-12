package legacy

import "github.com/brutusmcforce/opennox/server"

func init() {
	server.RegisterObjectDamageValueGo("DefaultDamage", damageIdentityKey(damageIDDefault), damageDefault)
	server.RegisterObjectDamageValueGo("SkeletonDamage", damageIdentityKey(damageIDSkeleton), damageSkeleton)
	server.RegisterObjectDamageValueGo("PlayerDamage", damageIdentityKey(damageIDPlayer), damagePlayer)
	server.RegisterObjectDamageValueGo("StoneDamage", damageIdentityKey(damageIDStone), damageDefault)
	server.RegisterObjectDamageValueGo("MechGolemDamage", damageIdentityKey(damageIDMechGolem), damageMechGolem)
	server.RegisterObjectDamageValueGo("FlammableDamage", damageIdentityKey(damageIDFlammable), damageFlammable)
	server.RegisterObjectDamageValueGo("BlackPowderDamage", damageIdentityKey(damageIDBlackPowder), damageBlackPowder)
	server.RegisterObjectDamageValueGo("ArmorDamage", damageIdentityKey(damageIDArmor), damageArmor)
	server.RegisterObjectDamageValueGo("WeaponDamage", damageIdentityKey(damageIDWeapon), damageWeapon)
	server.RegisterObjectDamageValueGo("BallDamage", damageIdentityKey(damageIDBall), func(_, _, _ *server.Object, _, _ int32) int32 { return 0 })
	server.RegisterObjectDamageValueGo("MonsterGeneratorDamage", damageIdentityKey(damageIDMonsterGenerator), damageGenerator)

	server.RegisterObjectDamageSoundGo("DefaultDamageSound", damageIdentityKey(damageIDDefaultSound), func(u, other *server.Object) {
		Nox_xxx_soundDefaultDamageSound_532E20(u, other)
	})
	server.RegisterObjectDamageSoundGo("PlayerDamageSound", damageIdentityKey(damageIDPlayerSound), func(u, other *server.Object) {
		Nox_xxx_soundPlayerDamageSound_5328B0(u, other)
	})
}
