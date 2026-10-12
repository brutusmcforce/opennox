package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

type durationCallbackID uint8

const (
	durationKeyNoxXxxSpellBlink2530310 durationCallbackID = iota
	durationKeyNoxXxxSpellBlink1530380
	durationKeySub52F460
	durationKeyNoxXxxCharmCreature15011F0
	durationKeyNoxXxxCharmCreatureFinish5013E0
	durationKeyNoxXxxCharmCreature2501690
	durationKeyNoxXxxSpellTurnUndeadCreate531310
	durationKeyNoxXxxSpellTurnUndeadUpdate531410
	durationKeyNoxXxxSpellTurnUndeadDelete531420
	durationKeyNoxXxxSpellDrainMana52E210
	durationKeyNoxXxxSpellEnergyBoltStop52E820
	durationKeyNoxXxxSpellEnergyBoltTick52E850
	durationKeyNullsub29
	durationKeyNoxXxxFirewalkTick52ED40
	durationKeySub52EF30
	durationKeySub52EFD0
	durationKeySub52F1D0
	durationKeySub52F220
	durationKeySub52F2E0
	durationKeyNoxXxxOnStartLightning52F820
	durationKeyNoxXxxOnFrameLightning52F8A0
	durationKeySub530100
	durationKeyNoxXxxCastShield152F5A0
	durationKeySub52F650
	durationKeySub52F670
	durationKeyNoxXxxSpellCreateMoonglow531A00
	durationKeySub531AF0
	durationKeyNoxXxxManaBomb530F90
	durationKeyNoxXxxManaBombBoom5310C0
	durationKeySub531290
	durationKeyNoxXxxPlasmaSmth531580
	durationKeyNoxXxxPlasmaShot531600
	durationKeySub5319E0
	durationKeySub531490
	durationKeySub5314F0
	durationKeySub531560
	durationKeyNoxXxxSummonStart500DA0
	durationKeyNoxXxxSummonFinish5010D0
	durationKeyNoxXxxSummonCancel5011C0
	durationKeySub530CA0
	durationKeySub530D30
	durationKeyNoxXxxSpellTagCreature530160
	durationKeySub530250
	durationKeySub530270
	durationKeySub5305D0
	durationKeySub530650
	durationKeyNoxXxxCastTele530820
	durationKeySub530880
	durationKeySub530A30SpellExecdur
	durationKeyNoxXxxCastTTT530B70
	durationKeyNoxXxxSpellWallCreate4FFA90
	durationKeyNoxXxxSpellWallUpdate500070
	durationKeyNoxXxxSpellWallDestroy500080
)

// Each callback identity has process lifetime and a distinct address. The
// one-byte slots are opaque keys; callback dispatch uses the server registry.
var durationCallbackSlots [53]byte

func durationCallbackKey(id durationCallbackID) unsafe.Pointer {
	return unsafe.Pointer(&durationCallbackSlots[id])
}

func init() {
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellBlink2530310), func(sp *server.DurSpell) int32 {
		return int32(sustainedBlinkStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellBlink1530380), func(sp *server.DurSpell) int32 {
		return int32(sustainedBlinkTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F460), func(sp *server.DurSpell) int32 {
		return int32(sustainedChannelLife(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCharmCreature15011F0), func(sp *server.DurSpell) int32 {
		return int32(spellEffectCharmStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCharmCreatureFinish5013E0), func(sp *server.DurSpell) int32 {
		return int32(spellEffectCharmFinish(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCharmCreature2501690), func(sp *server.DurSpell) int32 {
		return int32(spellEffectCharmCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadCreate531310), func(sp *server.DurSpell) int32 {
		return int32(sustainedTurnUndeadStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadUpdate531410), func(sp *server.DurSpell) int32 {
		return int32(sustainedTurnUndeadTick())
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadDelete531420), func(sp *server.DurSpell) int32 {
		return int32(sustainedTurnUndeadCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellDrainMana52E210), func(sp *server.DurSpell) int32 {
		return int32(sustainedDrainMana(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellEnergyBoltStop52E820), func(sp *server.DurSpell) int32 {
		return int32(sustainedEnergyStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellEnergyBoltTick52E850), func(sp *server.DurSpell) int32 {
		return int32(sustainedEnergyTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNullsub29), func(sp *server.DurSpell) int32 {
		return 0
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxFirewalkTick52ED40), func(sp *server.DurSpell) int32 {
		return int32(sustainedFirewalk(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52EF30), func(sp *server.DurSpell) int32 {
		return int32(sustainedForceStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52EFD0), func(sp *server.DurSpell) int32 {
		return int32(sustainedForceTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F1D0), func(sp *server.DurSpell) int32 {
		return int32(sustainedForceCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F220), func(sp *server.DurSpell) int32 {
		return int32(sustainedGreaterHealStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F2E0), func(sp *server.DurSpell) int32 {
		return int32(sustainedGreaterHealTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxOnStartLightning52F820), func(sp *server.DurSpell) int32 {
		return int32(sustainedLightningStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxOnFrameLightning52F8A0), func(sp *server.DurSpell) int32 {
		return int32(sustainedLightningTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530100), func(sp *server.DurSpell) int32 {
		return int32(int8(sustainedLightningCancel(sp.C())))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCastShield152F5A0), func(sp *server.DurSpell) int32 {
		return int32(sustainedShieldStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F650), func(sp *server.DurSpell) int32 {
		return int32(sustainedShieldTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub52F670), func(sp *server.DurSpell) int32 {
		return int32(sustainedShieldCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellCreateMoonglow531A00), func(sp *server.DurSpell) int32 {
		return int32(sustainedMoonglowStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub531AF0), func(sp *server.DurSpell) int32 {
		return int32(sustainedMoonglowCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxManaBomb530F90), func(sp *server.DurSpell) int32 {
		return int32(sustainedManaBombStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxManaBombBoom5310C0), func(sp *server.DurSpell) int32 {
		return int32(sustainedManaBombTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub531290), func(sp *server.DurSpell) int32 {
		return int32(sustainedManaBombCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxPlasmaSmth531580), func(sp *server.DurSpell) int32 {
		return int32(sustainedPlasmaStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxPlasmaShot531600), func(sp *server.DurSpell) int32 {
		return int32(sustainedPlasmaTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub5319E0), func(sp *server.DurSpell) int32 {
		return int32(sustainedPlasmaCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub531490), func(sp *server.DurSpell) int32 {
		return int32(sustainedOvalShieldStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub5314F0), func(sp *server.DurSpell) int32 {
		return int32(sustainedOvalShieldTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub531560), func(sp *server.DurSpell) int32 {
		return int32(sustainedOvalShieldCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSummonStart500DA0), func(sp *server.DurSpell) int32 {
		return int32(spellEffectSummonStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSummonFinish5010D0), func(sp *server.DurSpell) int32 {
		return int32(spellEffectSummonFinish(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSummonCancel5011C0), func(sp *server.DurSpell) int32 {
		spellEffectSummonCancel(sp.C())
		return 0
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530CA0), func(sp *server.DurSpell) int32 {
		return int32(sustainedSwapStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530D30), func(sp *server.DurSpell) int32 {
		return int32(sustainedSwapTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellTagCreature530160), func(sp *server.DurSpell) int32 {
		return int32(sustainedTagStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530250), func(sp *server.DurSpell) int32 {
		return int32(sustainedTagTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530270), func(sp *server.DurSpell) int32 {
		return int32(sustainedTagCancel(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub5305D0), func(sp *server.DurSpell) int32 {
		return int32(sustainedGlyphStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530650), func(sp *server.DurSpell) int32 {
		return int32(sustainedGlyphTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCastTele530820), func(sp *server.DurSpell) int32 {
		return int32(sustainedTeleportStart(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530880), func(sp *server.DurSpell) int32 {
		return int32(sustainedRandomGlyphTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeySub530A30SpellExecdur), func(sp *server.DurSpell) int32 {
		return int32(spellStartTeleport(sp))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxCastTTT530B70), func(sp *server.DurSpell) int32 {
		return int32(sustainedTeleportToPointTick(sp.C()))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellWallCreate4FFA90), func(sp *server.DurSpell) int32 {
		return int32(Nox_xxx_spellWallCreate_4FFA90(sp))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellWallUpdate500070), func(sp *server.DurSpell) int32 {
		return int32(Nox_xxx_spellWallUpdate_500070(sp))
	})
	server.RegisterDurSpellCallback(durationCallbackKey(durationKeyNoxXxxSpellWallDestroy500080), func(sp *server.DurSpell) int32 {
		Nox_xxx_spellWallDestroy_500080(sp)
		return 0
	})
}
