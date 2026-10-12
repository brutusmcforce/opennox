package legacy

import "unsafe"

type modifierIdentityID uint8

const (
	modifierIDDamageMultiplierEffect       modifierIdentityID = iota // nox_xxx_effectDamageMultiplier_4E04C0
	modifierIDStunEffect                   modifierIdentityID = iota // nox_xxx_stunEffect_4E04D0
	modifierIDFireEffect                   modifierIdentityID = iota // nox_xxx_fireEffect_4E0550
	modifierIDFireRingEffect               modifierIdentityID = iota // nox_xxx_fireRingEffect_4E05B0
	modifierIDBlueFireRingEffect           modifierIdentityID = iota // nox_xxx_blueFREffect_4E05F0
	modifierIDFrostEffect                  modifierIdentityID = iota // nullsub_38
	modifierIDRecoilEffect                 modifierIdentityID = iota // nox_xxx_recoilEffect_4E0640
	modifierIDConfuseEffect                modifierIdentityID = iota // nox_xxx_confuseEffect_4E0670
	modifierIDLightningEffect              modifierIdentityID = iota // nox_xxx_lightngEffect_4E06F0
	modifierIDDrainManaEffect              modifierIdentityID = iota // nox_xxx_drainMEffect_4E0740
	modifierIDVampirismEffect              modifierIdentityID = iota // nox_xxx_vampirismEffect_4E07C0
	modifierIDPoisonEffect                 modifierIdentityID = iota // nox_xxx_poisonEffect_4E0850
	modifierIDPanicEffect                  modifierIdentityID = iota // nullsub_39
	modifierIDSympathyEffect               modifierIdentityID = iota // nox_xxx_sympathyEffect_4E08E0
	modifierIDReadinessEffect              modifierIdentityID = iota // nullsub_22
	modifierIDProjectileSpeedEffect        modifierIdentityID = iota // nox_xxx_effectProjectileSpeed_4E09B0
	modifierIDReplenishmentEffect          modifierIdentityID = iota // nullsub_36
	modifierIDArmorMultiplierEffect        modifierIdentityID = iota // sub_4E0370
	modifierIDDurabilityMultiplierEffect   modifierIdentityID = iota // sub_4E0380
	modifierIDResilienceEffect             modifierIdentityID = iota // nullsub_40
	modifierIDInversionEffect              modifierIdentityID = iota // nox_xxx_inversionEffect_4E03D0
	modifierIDGripEffect                   modifierIdentityID = iota // nox_xxx_gripEffect_4E0480
	modifierIDBreakingEffect               modifierIdentityID = iota // nullsub_41
	modifierIDPunctureProneEffect          modifierIdentityID = iota // nullsub_42
	modifierIDRegenerationUpdate           modifierIdentityID = iota // nox_xxx_effectRegeneration_4E01D0
	modifierIDParasiteUpdate               modifierIdentityID = iota // nullsub_43
	modifierIDAttractionUpdate             modifierIdentityID = iota // nullsub_44
	modifierIDContinualReplenishmentUpdate modifierIdentityID = iota // nox_xxx_attribContinualReplen_4E02C0
	modifierIDBrillianceEngage             modifierIdentityID = iota // sub_4DFB50
	modifierIDBrillianceDisengage          modifierIdentityID = iota // sub_4DFB80
	modifierIDSpeedEngage                  modifierIdentityID = iota // nox_xxx_effectSpeedEngage_4DFC30
	modifierIDSpeedDisengage               modifierIdentityID = iota // nox_xxx_effectSpeedDisengage_4DFCA0
	modifierIDFireProtectEngage            modifierIdentityID = iota // sub_4DFD10
	modifierIDFireProtectDisengage         modifierIdentityID = iota // nox_xxx_modifFireProtection_4DFD40
	modifierIDLightningProtectEngage       modifierIdentityID = iota // nox_xxx_buff_4DFD80
	modifierIDLightningProtectDisengage    modifierIdentityID = iota // sub_4DFDB0
	modifierIDPoisonProtectEngage          modifierIdentityID = iota // nox_xxx_checkPoisonProtectEnch_4DFDE0
	modifierIDPoisonProtectDisengage       modifierIdentityID = iota // sub_4DFE10
	modifierIDRegenerationEngage           modifierIdentityID = iota // sub_4E0140
	modifierIDRegenerationDisengage        modifierIdentityID = iota // sub_4E0170
)

var modifierIdentitySlots [40]byte

func modifierKey(id modifierIdentityID) unsafe.Pointer {
	return unsafe.Pointer(&modifierIdentitySlots[id])
}
