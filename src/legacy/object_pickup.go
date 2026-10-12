package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/player"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_pickupDefault_4F31E0         server.PickupFunc
	Nox_objectPickupAudEvent_4F3D50      server.PickupFunc
	Nox_xxx_pickupPotion_4F37D0          server.PickupFunc
	Nox_xxx_playerClassCanUseItem_57B3D0 func(item *server.Object, cl player.Class) bool
	Sub_57B370                           func(cl object.Class, sub object.SubClass, typ int) byte
)

func init() {
	server.RegisterObjectPickup("DefaultPickup", itemIdentityKey(itemIDDefaultPickup), func(who, it *server.Object, a3, a4 int) bool {
		return Nox_xxx_pickupDefault_4F31E0(who, it, a3, a4)
	})
	server.RegisterObjectPickup("FoodPickup", itemIdentityKey(itemIDFoodPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryFoodPickup(who, it, a3) != 0 })
	server.RegisterObjectPickup("UsePickup", itemIdentityKey(itemIDUsePickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryUsePickup(who, it, a3) != 0 })
	server.RegisterObjectPickup("ArmorPickup", itemIdentityKey(itemIDArmorPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryArmorPickup(who, it, a3, a4) != 0 })
	server.RegisterObjectPickup("WeaponPickup", itemIdentityKey(itemIDWeaponPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryWeaponPickup(who, it, a3, a4) != 0 })
	server.RegisterObjectPickup("OblivionPickup", itemIdentityKey(itemIDOblivionPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryOblivionPickup(who, it, a3, a4) != 0 })
	server.RegisterObjectPickup("TreasurePickup", itemIdentityKey(itemIDTreasurePickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryTreasurePickup(who, it, a3) != 0 })
	server.RegisterObjectPickup("TrapPickup", itemIdentityKey(itemIDTrapPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryTrapPickup(who, it, a3) != 0 })
	server.RegisterObjectPickup("PotionPickup", itemIdentityKey(itemIDPotionPickup), func(who, it *server.Object, a3, a4 int) bool {
		return Nox_xxx_pickupPotion_4F37D0(who, it, a3, a4)
	})
	server.RegisterObjectPickup("GoldPickup", itemIdentityKey(itemIDGoldPickup), func(who, it *server.Object, a3, a4 int) bool { return resourceGoldPickup(who, it, a3) })
	server.RegisterObjectPickup("AmmoPickup", itemIdentityKey(itemIDAmmoPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryAmmoPickup(who, it, a3, a4) != 0 })
	server.RegisterObjectPickup("SpellBookPickup", itemIdentityKey(itemIDSpellBookPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryBookPickup(who, it, a3, false) != 0 })
	server.RegisterObjectPickup("AbilityBookPickup", itemIdentityKey(itemIDAbilityBookPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryBookPickup(who, it, a3, true) != 0 })
	server.RegisterObjectPickup("CrownPickup", itemIdentityKey(itemIDCrownPickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryCrownPickup(who, it, a3) != 0 })
	server.RegisterObjectPickup("AudEventPickup", itemIdentityKey(itemIDAudEventPickup), func(who, it *server.Object, a3, a4 int) bool {
		return Nox_objectPickupAudEvent_4F3D50(who, it, a3, a4)
	})
	server.RegisterObjectPickup("AnkhTradablePickup", itemIdentityKey(itemIDAnkhTradablePickup), func(who, it *server.Object, a3, a4 int) bool { return inventoryAnkhPickup(who, it) != 0 })
}

func sub_419E60(u *nox_object_t) int {
	return bool2int(GetServer().S().Players.CheckXxx(asObjectS(u)))
}
