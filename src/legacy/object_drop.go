package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_objectDropAudEvent_4EE2F0 server.DropFunc
)

func init() {
	registerInventoryDrop("DefaultDrop", itemIdentityKey(itemIDDefaultDrop), inventoryDefaultDrop)
	registerInventoryDrop("ArmorDrop", itemIdentityKey(itemIDArmorDrop), func(u, it *server.Object, p *types.Pointf) int { return inventoryEquipmentDrop(u, it, p, true) })
	registerInventoryDrop("WeaponDrop", itemIdentityKey(itemIDWeaponDrop), func(u, it *server.Object, p *types.Pointf) int { return inventoryEquipmentDrop(u, it, p, false) })
	registerInventoryDrop("TreasureDrop", itemIdentityKey(itemIDTreasureDrop), inventoryTreasureDrop)
	registerInventoryDrop("GlyphDrop", itemIdentityKey(itemIDGlyphDrop), inventoryGlyphDrop)
	registerInventoryDrop("PotionDrop", itemIdentityKey(itemIDPotionDrop), inventoryPotionDrop)
	registerInventoryDrop("TrapDrop", itemIdentityKey(itemIDTrapDrop), inventoryTrapDrop)
	registerInventoryDrop("FoodDrop", itemIdentityKey(itemIDFoodDrop), inventoryFoodDrop)
	registerInventoryDrop("CrownDrop", itemIdentityKey(itemIDCrownDrop), inventoryCrownDrop)
	server.RegisterObjectDrop("AudEventDrop", itemIdentityKey(itemIDAudEventDrop), func(obj, obj2 *server.Object, pos types.Pointf) bool {
		return Nox_objectDropAudEvent_4EE2F0(obj, obj2, pos)
	})
	inventoryNativeDrops[itemIdentityKey(itemIDAudEventDrop)] = func(obj, obj2 *server.Object, pos *types.Pointf) int {
		return bool2int(Nox_objectDropAudEvent_4EE2F0(obj, obj2, *pos))
	}
	registerInventoryDrop("AnkhTradableDrop", itemIdentityKey(itemIDAnkhTradableDrop), inventoryDefaultDrop)
}

func Nox_xxx_dropDefault_4ED290(obj1 *server.Object, obj2 *server.Object, a3 *types.Pointf) int {
	return inventoryDefaultDrop(obj1, obj2, a3)
}

var inventoryNativeDrops = make(map[unsafe.Pointer]func(*server.Object, *server.Object, *types.Pointf) int)

func registerInventoryDrop(name string, ptr unsafe.Pointer, fn func(*server.Object, *server.Object, *types.Pointf) int) {
	inventoryNativeDrops[ptr] = fn
	server.RegisterObjectDrop(name, ptr, func(u, it *server.Object, pos types.Pointf) bool { return fn(u, it, &pos) != 0 })
}
