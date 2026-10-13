package legacy

import "github.com/brutusmcforce/opennox/server"

func init() {
	server.RegisterObjectDeathGo("PlayerDie", deathKey(deathIdentityPlayer), func(u *server.Object) { playerDeath(u) }, server.NoData)
	server.RegisterObjectDeathGo("PotionDie", deathKey(deathIdentityPotion), objectDeathPotion, server.NoData)
	server.RegisterObjectDeathGo("ImpEggDie", deathKey(deathIdentityImpEgg), func(u *server.Object) { objectDeathImpEgg(u) }, server.NoData)
	server.RegisterObjectDeathGo("GlyphDie", deathKey(deathIdentityGlyph), func(u *server.Object) { Nox_xxx_dieGlyph_54DF30(u) }, server.NoData)
	server.RegisterObjectDeathGo("BarrelDie", deathKey(deathIdentityBarrel), objectDeathBarrel, server.NoData)
	server.RegisterObjectDeathGo("CreateObjectDie", deathKey(deathIdentityCreateObject), func(u *server.Object) { objectDeathCreate(u, false) }, server.DataOf[server.CreateObjectDeathData]())
	server.RegisterObjectDeathGo("SpawnObjectDie", deathKey(deathIdentitySpawnObject), func(u *server.Object) { objectDeathCreate(u, true) }, server.DataOf[server.CreateObjectDeathData]())
	server.RegisterObjectDeathGo("PolypDie", deathKey(deathIdentityPolyp), diePolyp, server.NoData)
	server.RegisterObjectDeathGo("MarkerDie", deathKey(deathIdentityMarker), objectDeathMarker, server.NoData)
	server.RegisterObjectDeathGo("WeaponDie", deathKey(deathIdentityWeapon), objectDeathWeapon, server.NoData)
	server.RegisterObjectDeathGo("ArmorDie", deathKey(deathIdentityArmor), objectDeathArmor, server.NoData)
	server.RegisterObjectDeathGo("BoulderDie", deathKey(deathIdentityBoulder), objectDeathBoulder, server.NoData)
	server.RegisterObjectDeathGo("GameBallDie", deathKey(deathIdentityGameBall), func(u *server.Object) { objectiveBallReset(u) }, server.NoData)
	server.RegisterObjectDeathGo("MonsterGeneratorDie", deathKey(deathIdentityMonsterGenerator), generatorDeath, server.NoData)

	server.RegisterObjectDeathParse("CreateObjectDie", resourceObjectParser("death", "spawn"))
	server.RegisterObjectDeathParse("SpawnObjectDie", resourceObjectParser("death", "spawn"))
}
