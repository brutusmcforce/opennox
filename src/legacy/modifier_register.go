package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func registerNativeModifierCallbacks() {
	server.RegisterModifDamageEffect("DamageMultiplierEffect", modifierKey(modifierIDDamageMultiplierEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("StunEffect", modifierKey(modifierIDStunEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("FireEffect", modifierKey(modifierIDFireEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("FireRingEffect", modifierKey(modifierIDFireRingEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("BlueFireRingEffect", modifierKey(modifierIDBlueFireRingEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("FrostEffect", modifierKey(modifierIDFrostEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("RecoilEffect", modifierKey(modifierIDRecoilEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("ConfuseEffect", modifierKey(modifierIDConfuseEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("LightningEffect", modifierKey(modifierIDLightningEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("DrainManaEffect", modifierKey(modifierIDDrainManaEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("VampirismEffect", modifierKey(modifierIDVampirismEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("PoisonEffect", modifierKey(modifierIDPoisonEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("PanicEffect", modifierKey(modifierIDPanicEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("SympathyEffect", modifierKey(modifierIDSympathyEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("ReadinessEffect", modifierKey(modifierIDReadinessEffect), server.ModEffectParseInt)
	server.RegisterModifDamageEffect("ProjectileSpeedEffect", modifierKey(modifierIDProjectileSpeedEffect), server.ModEffectParseFloat)
	server.RegisterModifDamageEffect("ReplenishmentEffect", modifierKey(modifierIDReplenishmentEffect), server.ModEffectParseInt)
	server.RegisterModifDefendEffect("ArmorMultiplierEffect", modifierKey(modifierIDArmorMultiplierEffect), server.ModEffectParseFloat)
	server.RegisterModifDefendEffect("DurabilityMultiplierEffect", modifierKey(modifierIDDurabilityMultiplierEffect), server.ModEffectParseFloat)
	server.RegisterModifDefendEffect("ResilienceEffect", modifierKey(modifierIDResilienceEffect), server.ModEffectParseFloat)
	server.RegisterModifDefendEffect("InversionEffect", modifierKey(modifierIDInversionEffect), server.ModEffectParseInt)
	server.RegisterModifDefendEffect("GripEffect", modifierKey(modifierIDGripEffect), server.ModEffectParseInt)
	server.RegisterModifDefendEffect("BreakingEffect", modifierKey(modifierIDBreakingEffect), server.ModEffectParseFloat)
	server.RegisterModifDefendEffect("PunctureProneEffect", modifierKey(modifierIDPunctureProneEffect), server.ModEffectParseFloat)
	server.RegisterModifUpdateEffect("RegenerationUpdate", modifierKey(modifierIDRegenerationUpdate), server.ModEffectParseInt)
	server.RegisterModifUpdateEffect("ParasiteUpdate", modifierKey(modifierIDParasiteUpdate), server.ModEffectParseInt)
	server.RegisterModifUpdateEffect("AttractionUpdate", modifierKey(modifierIDAttractionUpdate), server.ModEffectParseInt)
	server.RegisterModifUpdateEffect("ContinualReplenishmentUpdate", modifierKey(modifierIDContinualReplenishmentUpdate), server.ModEffectParseInt)
	server.RegisterModifEngageEffect("BrillianceEngage", modifierKey(modifierIDBrillianceEngage), server.ModEffectParseInt)
	server.RegisterModifEngageEffect("BrillianceDisengage", modifierKey(modifierIDBrillianceDisengage), server.ModEffectParseInt)
	server.RegisterModifEngageEffect("SpeedEngage", modifierKey(modifierIDSpeedEngage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("SpeedDisengage", modifierKey(modifierIDSpeedDisengage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("FireProtectEngage", modifierKey(modifierIDFireProtectEngage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("FireProtectDisengage", modifierKey(modifierIDFireProtectDisengage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("LightningProtectEngage", modifierKey(modifierIDLightningProtectEngage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("LightningProtectDisengage", modifierKey(modifierIDLightningProtectDisengage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("PoisonProtectEngage", modifierKey(modifierIDPoisonProtectEngage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("PoisonProtectDisengage", modifierKey(modifierIDPoisonProtectDisengage), server.ModEffectParseFloat)
	server.RegisterModifEngageEffect("RegenerationEngage", modifierKey(modifierIDRegenerationEngage), nil)
	server.RegisterModifEngageEffect("RegenerationDisengage", modifierKey(modifierIDRegenerationDisengage), nil)
	server.RegisterModifierEffect5(modifierKey(modifierIDDamageMultiplierEffect), func(m *server.ModifierEff, a, b, c *server.Object, p unsafe.Pointer) {
		v := (*float32)(p)
		*v = float32(float64(m.Attack40.Valf) * float64(*v))
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDStunEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) { effectsStatus(m, b, c, true) })
	server.RegisterModifierEffect5(modifierKey(modifierIDFireEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) {
		GetServer().S().Nox_xxx_fireEffect_4E0550(unsafe.Pointer(m), a, b, c)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDFireRingEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) {
		Nox_xxx_fireRingEffect_4E05B0(unsafe.Pointer(m), a, b, c)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDBlueFireRingEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) {
		Nox_xxx_blueFREffect_4E05F0(unsafe.Pointer(m), a, b, c)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDFrostEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {})
	server.RegisterModifierEffect5(modifierKey(modifierIDRecoilEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) { effectsRecoil(m, a, c) })
	server.RegisterModifierEffect5(modifierKey(modifierIDConfuseEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) { effectsStatus(m, b, c, false) })
	server.RegisterModifierEffect5(modifierKey(modifierIDLightningEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) { effectsLightning(m, a, b, c) })
	server.RegisterModifierEffect5(modifierKey(modifierIDDrainManaEffect), func(m *server.ModifierEff, a, b, c *server.Object, p unsafe.Pointer) {
		modifierDrainNative(m, b, c, p)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDVampirismEffect), func(m *server.ModifierEff, a, b, c *server.Object, p unsafe.Pointer) {
		modifierVampirismNative(m, b, c, p)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDPoisonEffect), func(m *server.ModifierEff, a, b, c *server.Object, _ unsafe.Pointer) { effectsPoison(m, b, c) })
	server.RegisterModifierEffect5(modifierKey(modifierIDPanicEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {})
	server.RegisterModifierEffect5(modifierKey(modifierIDSympathyEffect), func(m *server.ModifierEff, a, b, c *server.Object, p unsafe.Pointer) {
		modifierSympathyNative(m, b, c, p)
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDReadinessEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {})
	server.RegisterModifierEffect5(modifierKey(modifierIDProjectileSpeedEffect), func(m *server.ModifierEff, a, b, c *server.Object, p unsafe.Pointer) {
		u := (*server.Object)(p)
		u.SpeedCur = float32(float64(m.Attack40.Valf) * float64(u.SpeedCur))
	})
	server.RegisterModifierEffect5(modifierKey(modifierIDReplenishmentEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {})
	server.RegisterModifierEffect6(modifierKey(modifierIDArmorMultiplierEffect), func(m *server.ModifierEff, a, b, c, d *server.Object, p unsafe.Pointer) {
		v := (*float32)(p)
		*v = float32(float64(m.Defend76.Valf) * float64(*v))
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDDurabilityMultiplierEffect), func(m *server.ModifierEff, a, b, c, d *server.Object, p unsafe.Pointer) {
		v := (*float32)(p)
		*v = float32((1 - float64(m.Defend76.Valf) + 1) * float64(*v))
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDResilienceEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDInversionEffect), func(m *server.ModifierEff, a, b, c, d *server.Object, p unsafe.Pointer) {
		effectsGrip(m, (*int32)(p), true)
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDGripEffect), func(m *server.ModifierEff, a, b, c, d *server.Object, p unsafe.Pointer) {
		effectsGrip(m, (*int32)(p), false)
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDBreakingEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {
	})
	server.RegisterModifierEffect6(modifierKey(modifierIDPunctureProneEffect), func(*server.ModifierEff, *server.Object, *server.Object, *server.Object, *server.Object, unsafe.Pointer) {
	})
	server.RegisterModifierEffect3(modifierKey(modifierIDRegenerationUpdate), func(m *server.ModifierEff, a, b *server.Object) int32 { effectsRegeneration(m, a); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDParasiteUpdate), func(*server.ModifierEff, *server.Object, *server.Object) int32 { return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDAttractionUpdate), func(*server.ModifierEff, *server.Object, *server.Object) int32 { return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDContinualReplenishmentUpdate), func(m *server.ModifierEff, a, b *server.Object) int32 { effectsReplenish(m, a); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDBrillianceEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsEngageFlag(u, 8, 75); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDBrillianceDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsDisengageFlag(u, 8, 76); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDSpeedEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsSpeed(m, u, true); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDSpeedDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsSpeed(m, u, false); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDFireProtectEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsEngageFlag(u, 1, 102); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDFireProtectDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 {
		if u != nil && it != nil {
			effectsDisengageFlag(u, 1, 103)
		}
		return 0
	})
	server.RegisterModifierEffect3(modifierKey(modifierIDLightningProtectEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsEngageFlag(u, 4, 106); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDLightningProtectDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsDisengageFlag(u, 4, 107); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDPoisonProtectEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsEngageFlag(u, 2, 110); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDPoisonProtectDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsDisengageFlag(u, 2, 111); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDRegenerationEngage), func(m *server.ModifierEff, u, it *server.Object) int32 { effectsEngageFlag(u, 32, 123); return 0 })
	server.RegisterModifierEffect3(modifierKey(modifierIDRegenerationDisengage), func(m *server.ModifierEff, u, it *server.Object) int32 {
		if u != nil && u.ObjClass&4 != 0 {
			effectsDisengageFlag(u, 32, 124)
		}
		return 0
	})
}

func modifierDrainNative(m *server.ModifierEff, u, target *server.Object, p unsafe.Pointer) {
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	effectsDrainMana(m, u, target, *(*int32)(p))
}
func modifierVampirismNative(m *server.ModifierEff, u, target *server.Object, p unsafe.Pointer) {
	if u == nil || target == nil || target.ObjClass&6 == 0 || target.ObjClass&2 != 0 && target.ObjSubClass&0x40 != 0 {
		return
	}
	effectsVampirism(m, u, target, *(*int32)(p))
}
func modifierSympathyNative(m *server.ModifierEff, u, target *server.Object, p unsafe.Pointer) {
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	effectsSympathy(m, u, target, *(*int32)(p))
}
