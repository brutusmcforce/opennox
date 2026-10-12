package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/things"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_spellAwardAll1_4EFD80     func(p *server.Player)
	Nox_xxx_spellAwardAll2_4EFC80     func(p *server.Player)
	Nox_xxx_spellAwardAll3_4EFE10     func(p *server.Player)
	Nox_xxx_spellTitle_424930         func(ind int) (string, bool)
	Nox_xxx_spellDescription_424A30   func(ind int) (string, bool)
	Nox_xxx_spellIcon_424A90          func(ind int) unsafe.Pointer
	Nox_xxx_spellIconHighlight_424AB0 func(ind int) unsafe.Pointer
	Nox_xxx_castSpellByUser_4FDD20    func(a1 int, a2 *server.Object, a3 unsafe.Pointer) int
	Nox_xxx_spellWallCreate_4FFA90    func(sp *server.DurSpell) int
	Nox_xxx_spellWallUpdate_500070    func(sp *server.DurSpell) int
	Nox_xxx_spellWallDestroy_500080   func(sp *server.DurSpell)
)

func nox_xxx_getEnchantSpell_424920(enc int) int {
	return int(server.EnchantID(enc).Spell())
}

func nox_xxx_spellAwardAll1_4EFD80(p *nox_playerInfo) { Nox_xxx_spellAwardAll1_4EFD80(asPlayerS(p)) }

func nox_xxx_spellAwardAll2_4EFC80(p *nox_playerInfo) { Nox_xxx_spellAwardAll2_4EFC80(asPlayerS(p)) }

func nox_xxx_spellAwardAll3_4EFE10(p *nox_playerInfo) { Nox_xxx_spellAwardAll3_4EFE10(asPlayerS(p)) }

func nox_xxx_spellGetAud44_424800(ind, a2 int) int {
	return int(GetServer().S().Spells.DefByInd(spell.ID(ind)).GetAudio(a2))
}

func nox_xxx_spellTitle_424930(ind int) *wchar2_t {
	s, ok := Nox_xxx_spellTitle_424930(ind)
	if !ok {
		return nil
	}
	return internWStr(s)
}

func nox_xxx_spellDescription_424A30(ind int) *wchar2_t {
	s, ok := Nox_xxx_spellDescription_424A30(ind)
	if !ok {
		return nil
	}
	return internWStr(s)
}

func nox_xxx_spellManaCost_4249A0(ind, a2 int) int {
	return GetServer().S().Spells.ManaCost(spell.ID(ind), a2)
}

func nox_xxx_spellHasFlags_424A50(ind, flags int) bool {
	return GetServer().S().Spells.HasFlags(spell.ID(ind), things.SpellFlags(flags))
}

func nox_xxx_spellFlags_424A70(ind int) uint32 {
	return uint32(GetServer().S().Spells.Flags(spell.ID(ind)))
}

func nox_xxx_spellIcon_424A90(ind int) unsafe.Pointer { return Nox_xxx_spellIcon_424A90(ind) }

func nox_xxx_spellIconHighlight_424AB0(ind int) unsafe.Pointer {
	return Nox_xxx_spellIconHighlight_424AB0(ind)
}

func nox_xxx_spellFirstValid_424AD0() int { return int(GetServer().S().Spells.FirstValid()) }

func nox_xxx_spellNextValid_424AF0(ind int) int {
	return int(GetServer().S().Spells.NextValid(spell.ID(ind)))
}

func nox_xxx_spellIsValid_424B50(ind int) bool {
	return GetServer().S().Spells.DefByInd(spell.ID(ind)).IsValid()
}

func nox_xxx_spellIsEnabled_424B70(ind int) bool {
	return GetServer().S().Spells.DefByInd(spell.ID(ind)).IsEnabled()
}

func nox_xxx_castSpellByUser_4FDD20(a1 int, a2 *nox_object_t, a3 unsafe.Pointer) int {
	return Nox_xxx_castSpellByUser_4FDD20(a1, asObjectS(a2), a3)
}

func Nox_xxx_spellCastByBook_4FCB80() {
	spellLifeCastBooks()
}
func Nox_xxx_playerResetProtectionCRC_56F7D0(a1 uint32, a2 int) {
	nox_xxx_playerResetProtectionCRC_56F7D0(int32(a1), int32(a2))
}
func Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(a1 uint32, a2 int, a3 int) {
	nox_xxx_playerAwardSpellProtectionCRC_56FCE0(int32(a1), int32(a2), int32(a3))
}
func Nox_xxx_playerApplyProtectionCRC_56FD50(a1 uint32, a2 unsafe.Pointer, a3 int) {
	nox_xxx_playerApplyProtectionCRC_56FD50(int32(a1), a2, int32(a3))
}
func Nox_xxx_spellGrantToPlayer_4FB550(a1 *server.Object, a2 spell.ID, a3 int, a4 int, a5 int) int {
	return int(bookAwardSpell(a1, int32(a2), int32(a3), int32(a4), int32(a5)))
}
func Nox_xxx_gameCaptureMagic_4FDC10(a1 spell.ID, a2 *server.Object) int {
	return int(spellLifeCaptureAllowed(int32(a1), a2))
}
func Nox_xxx_createSpellFly_4FDDA0(a1 *server.Object, a2 *server.Object, a3 spell.ID) {
	spellLifeCreateFly(a1, a2, int32(a3))
}
func Nox_xxx_spellGetPower_4FE7B0(a1 spell.ID, a2 *server.Object) int {
	return int(spellLifePower(int32(a1), a2))
}

func Nox_xxx_spellArachna_52DC80(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectArachna(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castBurn_52C3E0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectBurn(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_spellCastCleansingFlame_52D5C0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectCleansingFlame(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castConfuse_52C1E0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectConfuse(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castCurePoison_52CDB0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectCurePoison(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castEquake_52DE40(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectQuake(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castFireball_52C790(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectFireball(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castFist_52D3C0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectFist(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castFumble_52C060(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectFumble(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52BEB0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectInversion(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52DD50(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectLesserHeal(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castLock_52CE90(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectLock(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52CA80(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectPortal(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52CBD0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectNamedPortal(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castMeteor_52D9D0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectMeteor(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castMeteorShower_52D8A0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectMeteorShower(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castPixies_540440(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return spellStartPixies(spellID, a3, a4, lvl)
}

func Nox_xxx_castPoison_52C720(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectPoison(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castPull_52BFA0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectPull(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castPush_52C000(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectPush(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castSpellWinkORrestoreHealth_52BF20(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectRestoreHealth(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52BF50(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectRestoreMana(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_useShock_52C5A0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectShock(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castStun_52C2C0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectStun(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castTelekinesis_52D330(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectTelekinesis(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Nox_xxx_castToxicCloud_52DB60(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectToxicCloud(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Sub_52CCD0(spellID spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int {
	return int(spellEffectDetonateGlyph(int32(spellID), a2, a3, a4, unsafe.Pointer(sa), int32(lvl)))
}

func Get_nox_xxx_spellBlink2_530310() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellBlink2530310)
}

func Get_nox_xxx_spellBlink1_530380() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellBlink1530380)
}

func Get_sub_52F460() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F460)
}

func Get_nox_xxx_charmCreature1_5011F0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCharmCreature15011F0)
}

func Get_nox_xxx_charmCreatureFinish_5013E0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCharmCreatureFinish5013E0)
}

func Get_nox_xxx_charmCreature2_501690() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCharmCreature2501690)
}

func Get_nox_xxx_spellTurnUndeadCreate_531310() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadCreate531310)
}

func Get_nox_xxx_spellTurnUndeadUpdate_531410() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadUpdate531410)
}

func Get_nox_xxx_spellTurnUndeadDelete_531420() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellTurnUndeadDelete531420)
}

func Get_nox_xxx_spellDrainMana_52E210() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellDrainMana52E210)
}

func Get_nox_xxx_spellEnergyBoltStop_52E820() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellEnergyBoltStop52E820)
}

func Get_nox_xxx_spellEnergyBoltTick_52E850() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellEnergyBoltTick52E850)
}

func Get_nullsub_29() unsafe.Pointer {
	return durationCallbackKey(durationKeyNullsub29)
}

func Get_nox_xxx_firewalkTick_52ED40() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxFirewalkTick52ED40)
}

func Get_sub_52EF30() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52EF30)
}

func Get_sub_52EFD0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52EFD0)
}

func Get_sub_52F1D0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F1D0)
}

func Get_sub_52F220() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F220)
}

func Get_sub_52F2E0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F2E0)
}

func Get_nox_xxx_onStartLightning_52F820() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxOnStartLightning52F820)
}

func Get_nox_xxx_onFrameLightning_52F8A0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxOnFrameLightning52F8A0)
}

func Get_sub_530100() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530100)
}

func Get_nox_xxx_castShield1_52F5A0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCastShield152F5A0)
}

func Get_sub_52F650() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F650)
}

func Get_sub_52F670() unsafe.Pointer {
	return durationCallbackKey(durationKeySub52F670)
}

func Get_nox_xxx_spellCreateMoonglow_531A00() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellCreateMoonglow531A00)
}

func Get_sub_531AF0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub531AF0)
}

func Get_nox_xxx_manaBomb_530F90() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxManaBomb530F90)
}

func Get_nox_xxx_manaBombBoom_5310C0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxManaBombBoom5310C0)
}

func Get_sub_531290() unsafe.Pointer {
	return durationCallbackKey(durationKeySub531290)
}

func Get_nox_xxx_plasmaSmth_531580() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxPlasmaSmth531580)
}

func Get_nox_xxx_plasmaShot_531600() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxPlasmaShot531600)
}

func Get_sub_5319E0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub5319E0)
}

func Get_sub_531490() unsafe.Pointer {
	return durationCallbackKey(durationKeySub531490)
}

func Get_sub_5314F0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub5314F0)
}

func Get_sub_531560() unsafe.Pointer {
	return durationCallbackKey(durationKeySub531560)
}

func Get_nox_xxx_summonStart_500DA0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSummonStart500DA0)
}

func Get_nox_xxx_summonFinish_5010D0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSummonFinish5010D0)
}

func Get_nox_xxx_summonCancel_5011C0() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSummonCancel5011C0)
}

func Get_sub_530CA0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530CA0)
}

func Get_sub_530D30() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530D30)
}

func Get_nox_xxx_spellTagCreature_530160() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellTagCreature530160)
}

func Get_sub_530250() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530250)
}

func Get_sub_530270() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530270)
}

func Get_sub_5305D0() unsafe.Pointer {
	return durationCallbackKey(durationKeySub5305D0)
}

func Get_sub_530650() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530650)
}

func Get_nox_xxx_castTele_530820() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCastTele530820)
}

func Get_sub_530880() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530880)
}

func Get_sub_530A30_spell_execdur() unsafe.Pointer {
	return durationCallbackKey(durationKeySub530A30SpellExecdur)
}

func Get_nox_xxx_castTTT_530B70() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxCastTTT530B70)
}

func Get_nox_xxx_spellWallCreate_4FFA90() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellWallCreate4FFA90)
}

func Get_nox_xxx_spellWallUpdate_500070() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellWallUpdate500070)
}

func Get_nox_xxx_spellWallDestroy_500080() unsafe.Pointer {
	return durationCallbackKey(durationKeyNoxXxxSpellWallDestroy500080)
}
