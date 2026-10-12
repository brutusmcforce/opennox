package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func damageBall(source, u *server.Object, amount int32) {
	if u.ObjClass&4 == 0 || amount < 30 {
		return
	}
	typ := memmap.PtrUint32(0x5d4594, 1563316)
	if *typ == 0 {
		*typ = uint32(GetServer().S().Types.IndByID("GameBall"))
	}
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) != *typ {
			continue
		}
		it.ObjFlags &^= 0x40
		nox_xxx_objectApplyForce_52DF80((*float32)(unsafe.Pointer(&u.PosVec)), asObjectC(it), 30)
		nox_xxx_unitClearOwner_4EC300(asObjectC(it))
		objectiveRememberOwner(it, u)
		ind := *(*byte)(unsafe.Add(source.CObj(), 52))
		if it.TeamVal.Has() {
			t := nox_xxx_getTeamByID_418AB0(int(ind))
			if t != nil {
				teamRuntimeSwitch(it.TeamPtr(), asTeamP(unsafe.Pointer(t)), int(it.NetCode), 0)
			}
		} else {
			teamRuntimeJoin(server.TeamID(ind), it.TeamPtr(), 1, int(it.NetCode), 0)
		}
		inventorySound(926, u, 0, 0)
		return
	}
}
func damageSetSource(u, source *server.Object, kind int32) {
	*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 520)) = source.CObj()
	*equipmentWord(u.CObj(), 524) = uint32(kind)
	*equipmentWord(u.CObj(), 536) = GetServer().S().Frame()
}
func damageDefault(u, source, weapon *server.Object, amount, kind int32) int32 {
	if u == nil {
		return 1
	}
	frame := GetServer().S().Frame()
	if u.Buffs&(1<<23) != 0 {
		if frame&3 == 0 {
			inventorySound(71, u, 0, 0)
		}
		return 1
	}
	if u.ObjClass&2 != 0 {
		*equipmentWord(u.UpdateData, 2188) = 0
		if kind == 1 || kind == 12 || kind == 7 || kind == 14 || kind == 6 {
			*equipmentWord(u.UpdateData, 1440) |= 0x80000
		}
	}
	actual := weapon
	if actual == nil {
		actual = source
	}
	if u.ObjFlags&0x8000 != 0 {
		if monsterIsZombie(u) {
			damageSetSource(u, actual, kind)
		}
		return 1
	}
	if !bool(nox_xxx_CheckGameplayFlags_417DA0(1)) {
		if owner := source.FindOwnerChainPlayer(); owner != nil && owner.ObjClass&6 != 0 && !GetServer().S().IsEnemyTo(u, owner) && (u != owner || bool(nox_common_gameFlags_check_40A5C0(4096))) {
			return 1
		}
	}
	if source != nil && weapon != nil && !GetServer().S().IsEnemyTo(u, source) && u.ObjClass&6 != 0 && damageMelee(source, weapon) && !damageFriendlyWeapon(weapon) || u.ObjFlags&2 != 0 {
		return 1
	}
	if source != nil && u.Buffs&(1<<22) != 0 && source.ObjClass&6 != 0 && weapon != nil && damageMelee(source, weapon) {
		inventorySound(135, source, 0, 0)
		spellLifeBuffOff(u, int32(22))
		n := floatToInt32(float32(nox_xxx_gamedataGetFloatTable_419D70(internCStr("ShockDamage"), 4)))
		projectileDamage(source, u, nil, n, 9)
		if source.ObjClass&4 != 0 {
			nox_xxx_playerSetState_4FA020(asObjectC(source), 23)
		}
	}
	if u.ObjClass&2 != 0 {
		sub := uint32(u.ObjSubClass)
		if sub&0x200 != 0 && kind == 5 {
			return 1
		}
		if sub&0x400 != 0 {
			if kind == 1 || kind == 12 {
				return 1
			}
			if kind == 7 {
				amount /= 2
			}
		} else if sub&0x800 != 0 && (kind == 9 || kind == 17) {
			return 1
		}
	}
	if kind == 1 || kind == 12 || kind == 7 {
		protect := effectsProtection(u, modifierKey(modifierIDFireProtectEngage), 17, "FireSpellProtection", .5, .60000002)
		if protect != 0 && GetServer().S().Frame()&3 == 0 {
			inventorySound(104, u, 0, 0)
		}
		amount = floatToInt32(float32((1 - float64(float32(protect))) * float64(amount)))
		if amount == 0 {
			amount = 1
		}
	}
	if kind == 9 || kind == 17 {
		protect := effectsProtection(u, modifierKey(modifierIDLightningProtectEngage), 20, "ElectricitySpellProtection", .5, .60000002)
		if protect != 0 && GetServer().S().Frame()&3 == 0 {
			inventorySound(108, u, 0, 0)
		}
		amount = floatToInt32(float32((1 - float64(float32(protect))) * float64(amount)))
		if amount == 0 {
			amount = 1
		}
		if u.ObjClass&4 != 0 {
			*(*uint16)(unsafe.Add(u.UpdateData, 160)) = 2
		} else if u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0 {
			*(*byte)(unsafe.Add(u.UpdateData, 2094)) = 2
		}
	}
	if source == nil {
		*equipmentWord(u.CObj(), 528) = 0
		*equipmentWord(u.CObj(), 532) = 0
		if kind == 12 {
			spellLifeBuffOff(u, int32(0))
		}
	} else {
		from := source
		if weapon != nil && weapon.ObjClass&0x1001000 == 0 {
			from = weapon
		}
		*equipmentWord(u.CObj(), 528) = *equipmentWord(from.CObj(), 72)
		*equipmentWord(u.CObj(), 532) = *equipmentWord(from.CObj(), 76)
		if u.ObjClass&2 != 0 && (weapon != nil && weapon != source || weapon == nil && (kind == 10 || kind == 2)) {
			*equipmentWord(u.UpdateData, 2188) = 1
			*equipmentWord(u.UpdateData, 2184) = uint32(actual.TypeInd)
		}
		spellLifeBuffOff(u, int32(0))
	}
	value, free := alloc.New(int32(0))
	*value = amount
	defer free()
	if u.ObjClass&4 != 0 || u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0 {
		damageDefend(u, source, weapon, value, kind)
	}
	damageSetSource(u, actual, kind)
	if u.ObjClass&2 != 0 {
		*equipmentWord(u.UpdateData, 1440) |= 0x200
		if *equipmentWord(u.UpdateData, 2188) == 0 {
			*equipmentWord(u.UpdateData, 2188) = 2
			*equipmentWord(u.UpdateData, 2184) = uint32(kind)
		}
	}
	if weapon != nil && weapon.ObjClass&0x1001000 != 0 {
		damagePre(u, source, weapon, value)
	}
	if u != weapon || u.ObjClass&0x1001000 == 0 {
		play := true
		if source != nil && source.ObjClass&2 != 0 && source.UpdateData != nil {
			set := resourceMonsterSound(source)
			if set != nil {
				sound := *(*int32)(unsafe.Add(set, 32))
				if sound != 0 && nox_xxx_getSevenDwords3_501940(int32(sound)) > 0 {
					play = false
				}
			}
		}
		if play {
			fn := *(*unsafe.Pointer)(unsafe.Add(u.CObj(), 720))
			if fn != nil {
				u.CallDamageSound(actual)
			} else {
				Nox_xxx_soundDefaultDamageSound_532E20(u, actual)
			}
		}
	}
	if source != nil && u.ObjClass&6 != 0 && source.Buffs&(1<<13) != 0 {
		inventorySound(163, weapon, 0, 0)
		coefficient := float64(nox_xxx_gamedataGetFloatTable_419D70(internCStr("VampirismCoeff"), int(int32(uint32(source.BuffsPower[13])-1))))
		heal := uint16(floatToInt32(float32(coefficient * float64(*value))))
		if heal < 1 {
			heal = 1
		}
		resourceAdjustHP(source, int32(heal))
		pos := [4]int32{floatToInt32(source.PosVec.X), floatToInt32(source.PosVec.Y), floatToInt32(u.PosVec.X), floatToInt32(u.PosVec.Y)}
		visibilityFXVampire(162, pos, heal)
	}
	damageBall(source, u, *value)
	if u.ObjClass&4 != 0 && *value >= 20 {
		state := *(*byte)(unsafe.Add(u.UpdateData, 88))
		if state != 1 && state != 15 {
			nox_xxx_playerSetState_4FA020(asObjectC(u), 30)
		}
	}
	if bool(nox_common_gameFlags_check_40A5C0(6144)) {
		controlGuideDamage(source, u, (*int32)(unsafe.Pointer(value)))
	}
	if source != nil {
		mob := source
		if mob.ObjClass&2 == 0 {
			mob = mob.ObjOwner
			if mob != nil && mob.ObjClass&2 == 0 {
				mob = nil
			}
		}
		if mob != nil && GetServer().S().IsEnemyTo(u, mob) {
			unitDamageTimer(mob)
		}
	}
	if u.Buffs&(1<<26) != 0 && kind != 5 {
		if kind != 15 || source != u {
			sustainedShieldDamage(u, value, kind, actual)
		}
		if *value == 0 {
			return 0
		}
	}
	resourceDamage(u, *value)
	return 1
}
func damageMechGolem(u, source, weapon *server.Object, amount, kind int32) int32 {
	if kind == 9 || kind == 17 {
		amount *= 2
	}
	return damageDefault(u, source, weapon, amount, kind)
}

func damageFlammable(u, source, weapon *server.Object, amount, kind int32) int32 {
	if kind == 1 || kind == 12 || kind == 7 {
		amount = 9999999
	}
	return damageDefault(u, source, weapon, amount, kind)
}

func damageBlackPowderAllowed(kind int32) bool {
	return kind == 0 || kind == 1 || kind == 2 || kind == 12
}

func damageBlackPowder(u, source, weapon *server.Object, amount, kind int32) int32 {
	if !damageBlackPowderAllowed(kind) {
		return 0
	}
	amount = 999999
	return damageDefault(u, source, weapon, amount, kind)
}

func damageWeapon(u, source, weapon *server.Object, amount, kind int32) int32 {
	if u != nil && u.ObjClass&0x1001000 != 0 && (u.InvHolder != nil || kind == 12) {
		return damageDefault(u, source, weapon, amount, kind)
	}
	return 0
}
func damageArmor(u, source, weapon *server.Object, amount, kind int32) int32 {
	if u.ObjClass&0x2000000 == 0 || u.InvHolder == nil && kind != 12 {
		return 0
	}
	if kind == 2 && *(*byte)(unsafe.Add(u.CObj(), 24))&0x10 != 0 {
		amount *= 2
	}
	if amount == 0 {
		return 0
	}
	return damageDefault(u, source, weapon, amount, kind)
}
func damageSkeleton(u, source, weapon *server.Object, amount, kind int32) int32 {
	if u.Buffs&(1<<23) != 0 {
		if GetServer().S().Frame()&3 == 0 {
			inventorySound(71, u, 0, 0)
		}
		return 1
	}
	if source != nil {
		actual := weapon
		if actual == nil {
			actual = source
		}
		if stateFront(&u.PosVec, int32(int16(u.Direction1)), &actual.PrevPos)&1 != 0 && monsterControlHead(u) == 21 && uint32(*(*byte)(unsafe.Add(u.UpdateData, 481))) > uint32(*(*byte)(unsafe.Add(u.UpdateData, 480))>>1) {
			inventorySound(878, u, 0, 0)
			return 1
		}
	}
	return damageDefault(u, source, weapon, amount, kind)
}
func damageGenerator(u, source, weapon *server.Object, amount, kind int32) int32 {
	ud := u.UpdateData
	if u.ObjFlags&0x8020 != 0 {
		return 0
	}
	frame := GetServer().S().Frame()
	if frame-*equipmentWord(u.CObj(), 536) > 20 || frame%30 == 0 {
		pos := u.PosVec
		if weapon != nil {
			pos.X = float32(float64(weapon.PosVec.X) - float64(u.PosVec.X) + 0.0099999998)
			pos.Y = float32(float64(weapon.PosVec.Y) - float64(u.PosVec.Y) + 0.0099999998)
			geometryNormalize((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&pos))))
			pos.X = float32(float64(pos.X)*22 + float64(u.PosVec.X))
			pos.Y = float32(float64(pos.Y)*22 + float64(u.PosVec.Y))
		}
		visibilityFXPointExtra(240, 26, pos)
		inventorySound(1001, u, 0, 0)
	}
	old := u.HealthData.Cur
	result := damageDefault(u, source, weapon, amount, kind)
	if u.HealthData.Cur < old {
		nox_xxx_scriptCallByEventBlock_502490(unsafe.Add(ud, 48), source.CObj(), u.CObj(), 23)
	}
	if u.ObjFlags&0x8020 == 0 {
		hp := int32(u.HealthData.Cur)
		if hp <= floatToInt32(float32(float64(u.HealthData.Max)*.333)) {
			if *equipmentWord(u.CObj(), 20)&0x100 != 0 {
				nox_xxx_unitUnsetXStatus_4E4780(asObjectC(u), 256)
			}
			if *equipmentWord(u.CObj(), 20)&0x200 == 0 {
				nox_xxx_unitSetXStatus_4E4800(asObjectC(u), 512)
			}
		} else if hp <= floatToInt32(float32(float64(u.HealthData.Max)*.66600001)) && *equipmentWord(u.CObj(), 20)&0x100 == 0 {
			nox_xxx_unitSetXStatus_4E4800(asObjectC(u), 256)
		}
	}
	return result
}
