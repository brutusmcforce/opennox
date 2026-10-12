package legacy

import "unsafe"

type itemIdentityID uint8

const (
	itemIDConsumeUse         itemIdentityID = 0
	itemIDConsumeConfuseUse  itemIdentityID = 1
	itemIDCastUse            itemIdentityID = 2
	itemIDEnchantUse         itemIdentityID = 3
	itemIDMushroomUse        itemIdentityID = 4
	itemIDPotionUse          itemIdentityID = 5
	itemIDFireWandUse        itemIdentityID = 6
	itemIDReadUse            itemIdentityID = 7
	itemIDWarpReadUse        itemIdentityID = 8
	itemIDWandUse            itemIdentityID = 9
	itemIDWandCastUse        itemIdentityID = 10
	itemIDSpellRewardUse     itemIdentityID = 11
	itemIDAbilityRewardUse   itemIdentityID = 12
	itemIDFieldGuideUse      itemIdentityID = 13
	itemIDDefaultDrop        itemIdentityID = 14
	itemIDArmorDrop          itemIdentityID = 15
	itemIDWeaponDrop         itemIdentityID = 16
	itemIDTreasureDrop       itemIdentityID = 17
	itemIDGlyphDrop          itemIdentityID = 18
	itemIDPotionDrop         itemIdentityID = 19
	itemIDTrapDrop           itemIdentityID = 20
	itemIDFoodDrop           itemIdentityID = 21
	itemIDCrownDrop          itemIdentityID = 22
	itemIDAudEventDrop       itemIdentityID = 23
	itemIDAnkhTradableDrop   itemIdentityID = 24
	itemIDDefaultPickup      itemIdentityID = 25
	itemIDFoodPickup         itemIdentityID = 26
	itemIDUsePickup          itemIdentityID = 27
	itemIDArmorPickup        itemIdentityID = 28
	itemIDWeaponPickup       itemIdentityID = 29
	itemIDOblivionPickup     itemIdentityID = 30
	itemIDTreasurePickup     itemIdentityID = 31
	itemIDTrapPickup         itemIdentityID = 32
	itemIDPotionPickup       itemIdentityID = 33
	itemIDGoldPickup         itemIdentityID = 34
	itemIDAmmoPickup         itemIdentityID = 35
	itemIDSpellBookPickup    itemIdentityID = 36
	itemIDAbilityBookPickup  itemIdentityID = 37
	itemIDCrownPickup        itemIdentityID = 38
	itemIDAudEventPickup     itemIdentityID = 39
	itemIDAnkhTradablePickup itemIdentityID = 40
)

type itemIdentitySlot struct {
	marker byte
}

var itemIdentitySlots [41]itemIdentitySlot

func itemIdentityKey(id itemIdentityID) unsafe.Pointer {
	return unsafe.Pointer(&itemIdentitySlots[id])
}
