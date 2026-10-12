package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Sub_4E3B80 func(ind int) bool
)

func init() {
	server.DefaultDamage = damageIdentityKey(damageIDDefault)
	server.DefaultDamageSound = damageIdentityKey(damageIDDefaultSound)
	server.DefaultXfer = xferIdentityKey(xferIDDefault)

	server.RegisterObjectCreateGo("MonsterCreate", lifecycleCreateKey(createIDMonster), func(u *server.Object) { Nox_xxx_monsterCreateFn_54C480(u) })
	server.RegisterObjectCreateGo("ArmorCreate", lifecycleCreateKey(createIDArmor), func(u *server.Object) { createArmor(u) })
	server.RegisterObjectCreateGo("WeaponCreate", lifecycleCreateKey(createIDWeapon), func(u *server.Object) { createWeapon(u) })
	server.RegisterObjectCreateGo("ObeliskCreate", lifecycleCreateKey(createIDObelisk), func(u *server.Object) { createObelisk(u) })
	server.RegisterObjectCreateGo("AnimCreate", lifecycleCreateKey(createIDAnim), func(u *server.Object) { createAnim(u) })
	server.RegisterObjectCreateGo("TriggerCreate", lifecycleCreateKey(createIDTrigger), func(u *server.Object) { createTrigger(u) })
	server.RegisterObjectCreateGo("MonsterGeneratorCreate", lifecycleCreateKey(createIDMonsterGenerator), func(u *server.Object) { createMonsterGenerator(u) })
	server.RegisterObjectCreateGo("RewardMarkerCreate", lifecycleCreateKey(createIDRewardMarker), func(u *server.Object) { createRewardMarker(u) })

	server.RegisterObjectInitGo("MonsterInit", lifecycleInitKey(initIDMonster), func(u *server.Object) { Nox_xxx_unitMonsterInit_4F0040(u) }, 0)
	server.RegisterObjectInitGo("PlayerInit", lifecycleInitKey(initIDPlayer), func(u *server.Object) { controlInitPlayer(u) }, 0)
	server.RegisterObjectInitGo("SparkInit", lifecycleInitKey(initIDSpark), func(u *server.Object) { rewardInitSpark(u) }, 0)
	server.RegisterObjectInitGo("FrogInit", lifecycleInitKey(initIDFrog), func(u *server.Object) { rewardInitFrog(u) }, 0)
	server.RegisterObjectInitGo("ChestInit", lifecycleInitKey(initIDChest), func(u *server.Object) { rewardInitBreakable(u) }, 0)
	server.RegisterObjectInitGo("BoulderInit", lifecycleInitKey(initIDBoulder), func(u *server.Object) { rewardInitBoulder(u) }, 0)
	server.RegisterObjectInitGo("BreakInit", lifecycleInitKey(initIDBreak), func(u *server.Object) { rewardInitBreakable(u) }, 0)
	server.RegisterObjectInitGo("MonsterGeneratorInit", lifecycleInitKey(initIDMonsterGenerator), func(u *server.Object) { rewardInitGenerator(u) }, 0)
	server.RegisterObjectInitGo("ShopkeeperInit", lifecycleInitKey(initIDMonster), func(u *server.Object) { Nox_xxx_unitMonsterInit_4F0040(u) }, unsafe.Sizeof(server.ShopkeeperInitData{}))
	server.RegisterObjectInitGo("SkullInit", lifecycleInitKey(initIDSkull), func(u *server.Object) { rewardInitDirection(u, true) }, 8)
	server.RegisterObjectInitGo("DirectionInit", lifecycleInitKey(initIDDirection), func(u *server.Object) { rewardInitDirection(u, false) }, 8)
	server.RegisterObjectInitGo("GoldInit", lifecycleInitKey(initIDGold), func(u *server.Object) { rewardInitGold(u) }, unsafe.Sizeof(server.GoldInitData{}))
}

func nox_xxx_newObjectWithTypeInd_4E3450(ind int) *nox_object_t {
	s := GetServer().S()
	return asObjectC(s.NewObjectByTypeInd(ind))
}

func nox_xxx_getUnitName_4E39D0(cobj *nox_object_t) *int8 {
	return (*int8)(internCStr(GetServer().S().Types.ByInd(int(asObjectS(cobj).TypeInd)).ID()))
}

func sub_4E3B80(ind int) int { return bool2int(Sub_4E3B80(ind)) }

func nox_xxx_getUnitNameByThingType_4E3A80(ind int) *int8 {
	if ind == 0 {
		return nil
	}
	return (*int8)(internCStr(GetServer().S().Types.ByInd(ind).ID()))
}

func nox_xxx_newObjectByTypeID_4E3810(cstr *int8) *nox_object_t {
	obj := GetServer().S().NewObjectByTypeID(GoStringP(unsafe.Pointer(cstr)))
	if obj == nil {
		return nil
	}
	return asObjectC(obj)
}

func Get_nox_xxx_XFerInvLight_4F5AA0() unsafe.Pointer {
	return xferIdentityKey(xferIDInvisibleLight)
}
func Get_nox_xxx_unitInitGold_4F04B0() unsafe.Pointer {
	return lifecycleInitKey(initIDGold)
}
