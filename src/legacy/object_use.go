package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_useConsume_53EE10      server.UseFunc
	Nox_xxx_useMushroom_53ECE0     server.UseFunc
	Nox_xxx_useCiderConfuse_53EF00 server.UseFunc
	Nox_xxx_useEnchant_53ED60      server.UseFunc
	Nox_xxx_useCast_53ED90         server.UseFunc
	Nox_xxx_usePotion_53EF70       server.UseFunc
)

func init() {
	server.RegisterObjectUseNative("ConsumeUse", itemIdentityKey(itemIDConsumeUse), func(obj, obj2 *server.Object) int32 { return int32(bool2int(Nox_xxx_useConsume_53EE10(obj, obj2))) }, server.DataOf[server.ConsumeUseData]())
	server.RegisterObjectUseNative("ConsumeConfuseUse", itemIdentityKey(itemIDConsumeConfuseUse), func(obj, obj2 *server.Object) int32 {
		return int32(bool2int(Nox_xxx_useCiderConfuse_53EF00(obj, obj2)))
	}, server.DataOf[server.ConsumeUseData]())
	server.RegisterObjectUseNative("CastUse", itemIdentityKey(itemIDCastUse), func(obj, obj2 *server.Object) int32 { return int32(bool2int(Nox_xxx_useCast_53ED90(obj, obj2))) }, server.DataOf[server.CastUseData]())
	server.RegisterObjectUseNative("EnchantUse", itemIdentityKey(itemIDEnchantUse), func(obj, obj2 *server.Object) int32 { return int32(bool2int(Nox_xxx_useEnchant_53ED60(obj, obj2))) }, server.DataOf[server.EnchantUseData]())
	server.RegisterObjectUseNative("MushroomUse", itemIdentityKey(itemIDMushroomUse), func(obj, obj2 *server.Object) int32 { return int32(bool2int(Nox_xxx_useMushroom_53ECE0(obj, obj2))) }, server.NoData)
	server.RegisterObjectUseNative("PotionUse", itemIdentityKey(itemIDPotionUse), func(obj, obj2 *server.Object) int32 { return int32(bool2int(Nox_xxx_usePotion_53EF70(obj, obj2))) }, server.DataOf[server.PotionUseData]())

	server.RegisterObjectUseNative("FireWandUse", itemIdentityKey(itemIDFireWandUse), func(u, it *server.Object) int32 { return int32(effectsFireWand(u, it)) }, server.NoData)
	server.RegisterObjectUseNative("ReadUse", itemIdentityKey(itemIDReadUse), func(u, it *server.Object) int32 { return int32(bool2int(unitRead(u, it, false))) }, server.DataOf[server.ReadUseData]())
	server.RegisterObjectUseNative("WarpReadUse", itemIdentityKey(itemIDWarpReadUse), func(u, it *server.Object) int32 { return int32(bool2int(unitRead(u, it, true))) }, server.DataOf[server.ReadUseData]())
	server.RegisterObjectUseNative("WandUse", itemIdentityKey(itemIDWandUse), func(u, it *server.Object) int32 { return int32(effectsLesserFireball(u, it)) }, server.DataOf[server.WandUseData]())
	server.RegisterObjectUseNative("WandCastUse", itemIdentityKey(itemIDWandCastUse), func(u, it *server.Object) int32 { return int32(effectsWandCast(u, it)) }, server.DataOf[server.WandUseData]())
	server.RegisterObjectUseNative("SpellRewardUse", itemIdentityKey(itemIDSpellRewardUse), func(u, it *server.Object) int32 { return int32(bookUseSpell(u, it)) }, server.DataOf[server.SpellRewardUseData]())
	server.RegisterObjectUseNative("AbilityRewardUse", itemIdentityKey(itemIDAbilityRewardUse), func(u, it *server.Object) int32 { return int32(bookUseAbility(u, it)) }, server.DataOf[server.AbilityRewardUseData]())
	server.RegisterObjectUseNative("FieldGuideUse", itemIdentityKey(itemIDFieldGuideUse), func(u, it *server.Object) int32 { return int32(bookUseGuide(u, it)) }, server.DataOf[server.FieldGuideUseData]())

	server.RegisterObjectUseParse("WandUse", resourceObjectParser("use", "wand"))
	server.RegisterObjectUseParse("WandCastUse", resourceObjectParser("use", "wandcast"))
}

func Get_nox_xxx_usePotion_53EF70() unsafe.Pointer {
	return itemIdentityKey(itemIDPotionUse)
}

func Get_nox_xxx_useSpellReward_53F9E0() unsafe.Pointer {
	return itemIdentityKey(itemIDSpellRewardUse)
}

func Get_nox_xxx_useAbilityReward_53FAE0() unsafe.Pointer {
	return itemIdentityKey(itemIDAbilityRewardUse)
}

func Get_nox_xxx_useEnchant_53ED60() unsafe.Pointer {
	return itemIdentityKey(itemIDEnchantUse)
}

func Get_nox_xxx_useCast_53ED90() unsafe.Pointer {
	return itemIdentityKey(itemIDCastUse)
}

func Get_sub_53F930() unsafe.Pointer {
	return itemIdentityKey(itemIDFieldGuideUse)
}
