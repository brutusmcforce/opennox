package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func controlRespawnItem(u *server.Object, name string, attrs unsafe.Pointer, a, b int32) *server.Object {
	it := GetServer().S().NewObjectByTypeID(name)
	if it == nil {
		return nil
	}
	if fn := *controlPtr(it.CObj(), 688); fn != nil {
		it.CallInitWithArg(nil)
	}
	if attrs != nil {
		stateAttributes(it, attrs)
	}
	Nox_xxx_inventoryServPlace_4F36F0(u, it, int(a), int(b))
	it.ObjFlags &^= 0x80000
	if it.ObjClass&0x3001000 != 0 {
		*equipmentWord(it.UpdateData, 4) |= 1
	}
	return it
}

// The return is a legacy signed byte, including the low byte of a player/item
// address on some paths. Keep it intact for the remaining C callers.
func controlDefaultItems(u *server.Object, refresh, keep int32) int8 {
	d := u.UpdateData
	pl := controlPlayer(u)
	if refresh != 0 {
		resourceRemovePoison(u)
		resourceRestoreHP(u)
		resourceRefreshMana(u)
	}
	Nox_xxx_playerCancelAbils_4FC180(u)
	gameplayReportReset(u)
	for _, off := range []int{312, 316, 84} {
		*equipmentWord(d, off) = 0
	}
	*controlByte(u.CObj(), 541) = 0
	*controlHalf(d, 78) = 0
	*controlHalf(d, 76) = 0
	hp := *controlHalf(*controlPtr(u.CObj(), 556), 0)
	for i := 0; i < 32; i++ {
		*controlHalf(d, 12+2*i) = hp
	}
	u.ObjFlags &= 0xffeb3fe7
	nox_xxx_playerSetState_4FA020((*nox_object_t)(u.CObj()), 13)
	spellLifeClearBuffs(u)
	*controlByte(d, 188) = 0
	for _, off := range []int{216, 192, 196, 200, 204, 208, 136, 132, 268} {
		*equipmentWord(d, off) = 0
	}
	*controlByte(d, 212) = 0
	controlClearWaypoints(u)
	*controlPtr(u.CObj(), 520) = nil
	if controlFlags(8192) {
		matchRosterResyncMask(*controlByte(pl, 2064))
	}
	result := int8(uintptr(pl))
	if pl == nil || *equipmentWord(pl, 4700) != 0 {
		return result
	}
	gameplayReportTotalHealth(int(*controlByte(pl, 2064)), u)
	gameplayReportTotalMana(int(*controlByte(pl, 2064)), u)
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
		desc := func(id int32) uint32 { return uint32(uintptr(nox_xxx_modifGetDescById_413330(int32(id)))) }
		byName := func(name string) uint32 {
			return desc(int32(nox_xxx_modifGetIdByName_413290((*int8)(internCStr(name)))))
		}
		base := byName("UserColor1")
		baseID := *equipmentWord(unsafe.Pointer(uintptr(base)), 4)
		color := func(off int) uint32 { return desc(int32(baseID + uint32(*controlByte(pl, 2185+off)))) }
		// The fifth modifier word is intentionally initialized along with descriptors.
		attrs := [5]uint32{}
		makeItem := func(name string) *server.Object { return controlRespawnItem(u, name, unsafe.Pointer(&attrs[0]), 1, 0) }
		if (controlFlags(2560) || *controlByte(pl, 2251) != 0) && *equipmentWord(pl, 0)&0x400 == 0 {
			attrs = [5]uint32{0, color(84), color(85), 0, 0}
			makeItem("StreetShirt")
		}
		if *controlByte(pl, 0)&4 == 0 {
			attrs = [5]uint32{0, color(83), 0, 0, 0}
			makeItem("StreetPants")
		}
		if *controlByte(pl, 0)&1 == 0 {
			attrs = [5]uint32{color(87), color(86), 0, 0, 0}
			makeItem("StreetSneakers")
		}
		class := *controlByte(pl, 2251)
		if controlFlags(2048) {
			attrs = [5]uint32{byName("ArmorQuality1"), 0, 0, 0, 0}
			if class == 0 {
				attrs[1] = byName("Material1")
			}
			name := GoStringP(unsafe.Pointer((*int8)(*controlPtr(memmap.PtrOff(0x587000, 206376), 4*int(class)))))
			result = int8(controlRaw(makeItem(name)))
		} else if controlFlags(4096) && sessionMapState() >= 0 {
			attrs = [5]uint32{}
			if class == 1 {
				attrs[2] = byName("Replenishment1")
			}
			name := GoStringP(unsafe.Pointer((*int8)(*controlPtr(memmap.PtrOff(0x587000, 206388), 4*int(class)))))
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
		*equipmentWord(pl, 4700) = 1
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
		*equipmentWord(d, 320) = uint32(floatToInt32(float32(nox_xxx_gamedataGetFloat_419D40(internCStr("QuestGameStartingExtraLives")))))
	}
	return controlDefaultItems(u, 1, 0)
}
func controlResetPlayer(u *server.Object) int32 {
	d := u.UpdateData
	pl := controlPlayer(u)
	nox_xxx_spellAwardAll1_4EFD80((*nox_playerInfo)(pl))
	nox_xxx_spellAwardAll2_4EFC80((*nox_playerInfo)(pl))
	*controlByte(pl, 3684) = 1
	Nox_xxx_playerCancelAbils_4FC180(u)
	controlReadStats(u, 0)
	nox_xxx_spellAwardAll3_4EFE10((*nox_playerInfo)(pl))
	mana := *controlHalf(d, 8)
	*controlHalf(d, 4) = mana
	*controlHalf(d, 6) = mana
	setProtectionRecord(int32(*equipmentWord(pl, 4596)), uint32(mana))
	for i := 192; i <= 208; i += 4 {
		*equipmentWord(d, i) = 0
	}
	*controlByte(d, 212) = 0
	resourceRestoreHP(u)
	*controlByte(u.CObj(), 541) = 0
	u.ObjFlags &= 0xffeb3fe7
	nox_xxx_playerSetState_4FA020((*nox_object_t)(u.CObj()), 13)
	spellLifeClearBuffs(u)
	spellLifeCancelPlayer(u)
	resourceRemovePoison(u)
	controlClearWaypoints(u)
	gameplayReportTotalHealth(int(*controlByte(pl, 2064)), u)
	gameplayReportTotalMana(int(*controlByte(pl, 2064)), u)
	*controlPtr(u.CObj(), 520) = nil
	*equipmentWord(pl, 3664) = 0xdeadface
	*equipmentWord(pl, 3660) = 0xdeadface
	return -559023410
}
func controlTeamFlag(pl unsafe.Pointer) {
	u := controlObject(pl, 2056)
	team := unsafe.Pointer(nox_xxx_getTeamByID_418AB0(int(*controlByte(u.CObj(), 52))))
	flag := controlObject(team, 76)
	if flag != nil && flag.InvHolder == nil {
		inventoryCrownPickup(u, flag, 1)
	}
}
func controlLeaveObserver(pl unsafe.Pointer) {
	if pl == nil {
		return
	}
	u := controlObject(pl, 2056)
	if u == nil || *controlPtr(u.CObj(), 744) == updateIdentityKey(updateIDPlayerMonsterBot) {
		return
	}
	playerStateRemoveStatus((*server.Player)(pl), 289)
	spellLifeBuffOff(u, int32(0))
	*controlPtr(u.CObj(), 744) = updateIdentityKey(updateIDPlayer)
	u.ObjFlags &^= 0x40
	nox_xxx_monsterMarkUpdate_4E8020((*nox_object_t)(controlObject(pl, 2056).CObj()))
	if controlFlags(16) && bool(nox_xxx_CheckGameplayFlags_417DA0(4)) {
		controlTeamFlag(pl)
	}
	if controlFlags(49152) && !matchRosterHasIdentity((*server.Player)(pl)) {
		matchRosterRemember((*server.Player)(pl))
	}
	if controlFlags(4096) {
		for it := GetServer().S().Players.FirstUnit(); it != nil; it = GetServer().S().Players.NextUnit(it) {
			if *equipmentWord(controlPlayer(it), 4792) == 1 {
				gameplayReportEnchant(int(int32(*controlByte(pl, 2064))), it)
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
	result := int16(uintptr(settings))
	if u == nil {
		return result
	}
	d := u.UpdateData
	pl := controlPlayer(u)
	if controlFlags(4096) && *equipmentWord(d, 548) != 0 {
		return int16(*equipmentWord(d, 548))
	}
	if pl != nil {
		*equipmentWord(pl, 4700) = 0
	}
	if controlFlags(4096) {
		controlDefaultItems(u, 1, 1)
		*controlByte(d, 452+int(*controlByte(controlPlayer(u), 2064))) = 250
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
	if target := controlObject(d, 308); controlFlags(4096) && target != nil {
		controlNearStart(target, &pos)
	} else {
		controlFindStart(&pos, u)
	}
	Nox_xxx_unitMove_4E7010(u, pos)
	if controlFlags(16) && bool(nox_xxx_CheckGameplayFlags_417DA0(4)) {
		controlTeamFlag(pl)
	}
	if controlFlags(8192) {
		spellLifeApplyBuff(u, 23, int16(5*uint16(GetServer().S().TickRate())), 5)
		return 1
	}
	return 0
}
func controlRespawnBot(u *server.Object) int32 {
	b := *controlPtr(u.UpdateData, 292)
	settings := unsafe.Pointer((unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())))
	if *controlHalf(*controlPtr(u.CObj(), 556), 0) == 0 {
		if GetServer().S().Frame()-*equipmentWord(b, 548) < 2*GetServer().S().TickRate() {
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
