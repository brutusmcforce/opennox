package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func controlRespawnItem(u *server.Object, name string, attrs *server.ModifierInitData, a, b int32) *server.Object {
	it := GetServer().S().NewObjectByTypeID(name)
	if it == nil {
		return nil
	}
	if fn := it.Init; fn != nil {
		it.CallInitWithArg(nil)
	}
	if attrs != nil {
		stateAttributes(it, attrs)
	}
	Nox_xxx_inventoryServPlace_4F36F0(u, it, int(a), int(b))
	it.ObjFlags &^= 0x80000
	if it.ObjClass&0x3001000 != 0 {
		it.UpdateDataWeaponArmor().Flags |= 1
	}
	return it
}

// The return is a legacy signed byte, including the low byte of a player/item
// address on some paths. Keep it intact for the remaining C callers.
func controlDefaultItems(u *server.Object, refresh, keep int32) int8 {
	d := u.UpdateData
	ud := layout.As[server.PlayerUpdateData](d)
	pl := controlPlayer(u)
	if refresh != 0 {
		resourceRemovePoison(u)
		resourceRestoreHP(u)
		resourceRefreshMana(u)
	}
	Nox_xxx_playerCancelAbils_4FC180(u)
	gameplayReportReset(u)
	ud.Field78, ud.Field79, ud.Field21 = nil, nil, 0
	u.Field541 = 0
	ud.Field19_1 = 0
	ud.Field19_0 = 0
	hp := u.HealthData.Cur
	for i := range ud.HealthPrev {
		ud.HealthPrev[i] = hp
	}
	u.ObjFlags &= 0xffeb3fe7
	nox_xxx_playerSetState_4FA020((*nox_object_t)(u.CObj()), 13)
	spellLifeClearBuffs(u)
	ud.Field47_0 = 0
	ud.SpellCastStart = 0
	ud.TrapSpells = [5]uint32{}
	ud.HarpoonBolt, ud.HarpoonTarg = nil, nil
	ud.Field67 = 0
	*(*byte)(unsafe.Pointer(&ud.TrapSpellsCnt)) = 0
	controlClearWaypoints(u)
	*(*unsafe.Pointer)(unsafe.Pointer(&u.Obj130)) = nil
	if controlFlags(8192) {
		matchRosterResyncMask(pl.PlayerInd)
	}
	result := int8(memmap.PtrWord(unsafe.Pointer(pl)))
	if pl == nil || pl.Field4700 != 0 {
		return result
	}
	gameplayReportTotalHealth(int(pl.PlayerInd), u)
	gameplayReportTotalMana(int(pl.PlayerInd), u)
	if keep != 0 {
		result = int8(controlRespawnNotify(u, 0))
	} else {
		for it := u.InvFirstItem; it != nil; {
			next := it.InvNextItem
			if equipmentArmorMask(it) != 0 || it.ObjFlags&0x100 == 0 || it.ObjClass&0x2000000 != 0 && nox_xxx_unitArmorInventoryEquipFlags_415C70((*nox_object_t)(it.CObj()))&0x808 != 0 {
				GetServer().DelayedDelete(it)
			}
			it = next
		}
		controlRespawnNotify(u, 1)
		desc := func(id int32) *server.ModifierEff {
			return (*server.ModifierEff)(nox_xxx_modifGetDescById_413330(int32(id)))
		}
		byName := func(name string) *server.ModifierEff {
			return desc(int32(nox_xxx_modifGetIdByName_413290((*int8)(internCStr(name)))))
		}
		base := byName("UserColor1")
		baseID := base.Ind4
		color := func(off int) *server.ModifierEff {
			return desc(int32(baseID + uint32(pl.InfoBytes[off])))
		}
		// The fifth modifier word is intentionally initialized along with descriptors.
		var attrs server.ModifierInitData
		mods := func(m ...*server.ModifierEff) server.ModifierInitData {
			var a server.ModifierInitData
			copy(a.Modifiers[:], m)
			return a
		}
		makeItem := func(name string) *server.Object { return controlRespawnItem(u, name, &attrs, 1, 0) }
		if (controlFlags(2560) || pl.InfoBytes[66] != 0) && pl.ArmorEquip&0x400 == 0 {
			attrs = mods(nil, color(84), color(85), nil, nil)
			makeItem("StreetShirt")
		}
		if *(*byte)(unsafe.Pointer(&pl.ArmorEquip))&4 == 0 {
			attrs = mods(nil, color(83), nil, nil, nil)
			makeItem("StreetPants")
		}
		if *(*byte)(unsafe.Pointer(&pl.ArmorEquip))&1 == 0 {
			attrs = mods(color(87), color(86), nil, nil, nil)
			makeItem("StreetSneakers")
		}
		class := pl.InfoBytes[66]
		if controlFlags(2048) {
			attrs = mods(byName("ArmorQuality1"), nil, nil, nil, nil)
			if class == 0 {
				attrs.Modifiers[1] = byName("Material1")
			}
			name := GoStringP(memmap.LoadPtr(0x587000, 206376+4*uintptr(class)))
			result = int8(controlRaw(makeItem(name)))
		} else if controlFlags(4096) && sessionMapState() >= 0 {
			attrs = server.ModifierInitData{}
			if class == 1 {
				attrs.Modifiers[2] = byName("Replenishment1")
			}
			name := GoStringP(memmap.LoadPtr(0x587000, 206388+4*uintptr(class)))
			result = int8(controlRaw(makeItem(name)))
		} else {
			result = int8(class)
			switch class {
			case 0:
				controlRespawnItem(u, "Longsword", nil, 1, 0)
				result = int8(controlRaw(controlRespawnItem(u, "WoodenShield", nil, 1, 0)))
			case 1:
				result = int8(controlRaw(controlRespawnItem(u, "WizardRobe", nil, 1, 0)))
			}
		}
	}
	pl = controlPlayer(u)
	if pl != nil {
		pl.Field4700 = 1
	}
	return result
}
func controlInitPlayer(u *server.Object) int8 {
	d := u.UpdateData
	resourceSubGold(u, resourceObjectGold(u))
	controlLevelFromXP(u)
	pl := controlPlayer(u)
	nox_xxx_spellAwardAll1_4EFD80((*nox_playerInfo)(pl))
	nox_xxx_spellAwardAll2_4EFC80((*nox_playerInfo)(pl))
	controlReadStats(u, 0)
	nox_xxx_spellAwardAll3_4EFE10((*nox_playerInfo)(pl))
	if controlFlags(4096) {
		layout.As[server.PlayerUpdateData](d).Field80 = uint32(floatToInt32(float32(nox_xxx_gamedataGetFloat_419D40(internCStr("QuestGameStartingExtraLives")))))
	}
	return controlDefaultItems(u, 1, 0)
}
func controlResetPlayer(u *server.Object) int32 {
	d := u.UpdateData
	pl := controlPlayer(u)
	nox_xxx_spellAwardAll1_4EFD80((*nox_playerInfo)(pl))
	nox_xxx_spellAwardAll2_4EFC80((*nox_playerInfo)(pl))
	pl.Level = 1
	Nox_xxx_playerCancelAbils_4FC180(u)
	controlReadStats(u, 0)
	nox_xxx_spellAwardAll3_4EFE10((*nox_playerInfo)(pl))
	mana := layout.As[server.PlayerUpdateData](d).ManaMax
	layout.As[server.PlayerUpdateData](d).ManaCur = mana
	layout.As[server.PlayerUpdateData](d).ManaPrev = mana
	setProtectionRecord(int32(pl.ProtUnitManaCur), uint32(mana))
	layout.As[server.PlayerUpdateData](d).TrapSpells = [5]uint32{}
	*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](d).TrapSpellsCnt)) = 0
	resourceRestoreHP(u)
	u.Field541 = 0
	u.ObjFlags &= 0xffeb3fe7
	nox_xxx_playerSetState_4FA020((*nox_object_t)(u.CObj()), 13)
	spellLifeClearBuffs(u)
	spellLifeCancelPlayer(u)
	resourceRemovePoison(u)
	controlClearWaypoints(u)
	gameplayReportTotalHealth(int(pl.PlayerInd), u)
	gameplayReportTotalMana(int(pl.PlayerInd), u)
	*(*unsafe.Pointer)(unsafe.Pointer(&u.Obj130)) = nil
	pl.Field3664 = 0xdeadface
	pl.Field3660 = 0xdeadface
	return -559023410
}
func controlTeamFlag(pl *server.Player) {
	u := pl.PlayerUnit
	team := unsafe.Pointer(nox_xxx_getTeamByID_418AB0(int(*(*byte)(unsafe.Pointer(&u.TeamVal.ID)))))
	flag := (*server.Object)(layout.As[server.Team](team).Field_76)
	if flag != nil && flag.InvHolder == nil {
		inventoryCrownPickup(u, flag, 1)
	}
}
func controlLeaveObserver(pl *server.Player) {
	if pl == nil {
		return
	}
	u := pl.PlayerUnit
	if u == nil || u.Update == updateIdentityKey(updateIDPlayerMonsterBot) {
		return
	}
	playerStateRemoveStatus((*server.Player)(pl), 289)
	spellLifeBuffOff(u, int32(0))
	u.Update = updateIdentityKey(updateIDPlayer)
	u.ObjFlags &^= 0x40
	nox_xxx_monsterMarkUpdate_4E8020((*nox_object_t)(pl.PlayerUnit.CObj()))
	if controlFlags(16) && bool(nox_xxx_CheckGameplayFlags_417DA0(4)) {
		controlTeamFlag(pl)
	}
	if controlFlags(49152) && !matchRosterHasIdentity((*server.Player)(pl)) {
		matchRosterRemember((*server.Player)(pl))
	}
	if controlFlags(4096) {
		for it := GetServer().S().Players.FirstUnit(); it != nil; it = GetServer().S().Players.NextUnit(it) {
			if controlPlayer(it).Field4792 == 1 {
				gameplayReportEnchant(int(int32(pl.PlayerInd)), it)
			}
		}
	}
}
func controlMakeCorpse(u *server.Object, settings unsafe.Pointer) {
	if dword_5d4594_2650652 == 0 || *equipmentWord(settings, 58) != 0 {
		playerCorpseCreate(u.PosVec, int32(int16(u.Direction1)))
	}
}
func controlRespawn(u *server.Object) int16 {
	settings := unsafe.Pointer((unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())))
	result := int16(memmap.PtrWord(settings))
	if u == nil {
		return result
	}
	d := u.UpdateData
	pl := controlPlayer(u)
	if pd := layout.As[server.PlayerUpdateData](d); controlFlags(4096) && pd.Field137 != 0 {
		return int16(pd.Field137)
	}
	if pl != nil {
		pl.Field4700 = 0
	}
	if controlFlags(4096) {
		controlDefaultItems(u, 1, 1)
		layout.As[server.PlayerUpdateData](d).QuestLevels[controlPlayer(u).PlayerInd] = 250
		nox_xxx_netPriMsgToPlayer_4DA2C0((*nox_object_t)(u.CObj()), (*int8)(internCStr("GeneralPrint:Respawn")), 0)
	} else {
		controlDefaultItems(u, 1, 0)
	}
	sound := 148
	if controlFlags(4096) {
		sound = 1006
	}
	nox_xxx_aud_501960(int32(int32(sound)), (*nox_object_t)(u.CObj()), 0, 0)
	controlMakeCorpse(u, settings)
	pos := u.PosVec
	if target := layout.As[server.PlayerUpdateData](d).Field77; controlFlags(4096) && target != nil {
		controlNearStart(target, &pos)
	} else {
		controlFindStart(&pos, u)
	}
	Nox_xxx_unitMove_4E7010(u, pos)
	if controlFlags(16) && bool(nox_xxx_CheckGameplayFlags_417DA0(4)) {
		controlTeamFlag((*server.Player)(unsafe.Pointer(pl)))
	}
	if controlFlags(8192) {
		spellLifeApplyBuff(u, 23, int16(5*uint16(GetServer().S().TickRate())), 5)
		return 1
	}
	return 0
}
func controlRespawnBot(u *server.Object) int32 {
	b := u.UpdateDataPlayer().BotData
	settings := unsafe.Pointer(serverConfigSettings())
	if u.HealthData.Cur == 0 {
		if GetServer().S().Frame()-b.Field137 < 2*GetServer().S().TickRate() {
			return 1
		}
		controlBotCreate(u)
		controlDefaultItems(u, 1, 0)
		controlMakeCorpse(u, settings)
		var pos types.Pointf
		controlFindStart(&pos, u)
		Nox_xxx_unitMove_4E7010(u, pos)
		nox_xxx_aud_501960(148, (*nox_object_t)(u.CObj()), 0, 0)
		if controlFlags(8192) {
			spellLifeApplyBuff(u, 23, int16(5*uint16(GetServer().S().TickRate())), 5)
		}
	}
	return 0
}
