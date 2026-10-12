package legacy

import "github.com/brutusmcforce/opennox/server"

func init() {
	server.RegisterObjectDeathGo("PlayerDie", deathKey(deathIdentityPlayer), func(u *server.Object) { playerDeath(u) }, 0)
	server.RegisterObjectDeathGo("PotionDie", deathKey(deathIdentityPotion), objectDeathPotion, 0)
	server.RegisterObjectDeathGo("ImpEggDie", deathKey(deathIdentityImpEgg), func(u *server.Object) { objectDeathImpEgg(u) }, 0)
	server.RegisterObjectDeathGo("GlyphDie", deathKey(deathIdentityGlyph), func(u *server.Object) { Nox_xxx_dieGlyph_54DF30(u) }, 0)
	server.RegisterObjectDeathGo("BarrelDie", deathKey(deathIdentityBarrel), objectDeathBarrel, 0)
	server.RegisterObjectDeathGo("CreateObjectDie", deathKey(deathIdentityCreateObject), func(u *server.Object) { objectDeathCreate(u, false) }, 132)
	server.RegisterObjectDeathGo("SpawnObjectDie", deathKey(deathIdentitySpawnObject), func(u *server.Object) { objectDeathCreate(u, true) }, 132)
	server.RegisterObjectDeathGo("PolypDie", deathKey(deathIdentityPolyp), diePolyp, 0)
	server.RegisterObjectDeathGo("MarkerDie", deathKey(deathIdentityMarker), objectDeathMarker, 0)
	server.RegisterObjectDeathGo("WeaponDie", deathKey(deathIdentityWeapon), objectDeathWeapon, 0)
	server.RegisterObjectDeathGo("ArmorDie", deathKey(deathIdentityArmor), objectDeathArmor, 0)
	server.RegisterObjectDeathGo("BoulderDie", deathKey(deathIdentityBoulder), objectDeathBoulder, 0)
	server.RegisterObjectDeathGo("GameBallDie", deathKey(deathIdentityGameBall), func(u *server.Object) { objectiveBallReset(u) }, 0)
	server.RegisterObjectDeathGo("MonsterGeneratorDie", deathKey(deathIdentityMonsterGenerator), generatorDeath, 0)

	server.RegisterObjectDeathParse("CreateObjectDie", resourceObjectParser("death", "spawn"))
	server.RegisterObjectDeathParse("SpawnObjectDie", resourceObjectParser("death", "spawn"))
}
