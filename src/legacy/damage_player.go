package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func damagePlayer(u, source, weapon *server.Object, amount, kind int32) int32 {
	if u.ObjFlags&0x8002 != 0 {
		return 0
	}
	if u.Buffs&(1<<23) != 0 {
		if GetServer().S().Frame()&3 == 0 {
			inventorySound(71, u, 0, 0)
		}
		return 1
	}
	if bool(nox_common_gameFlags_check_40A5C0(2048)) && source.FindOwnerChainPlayer() == u && u.ObjClass&4 != 0 && kind != 15 {
		return 0
	}
	actual := source
	if weapon != nil {
		actual = weapon
	}
	player := u.ObjClass&4 != 0
	ud := u.UpdateData
	var weapons, armor uint32
	var coeff float32
	if player {
		data := unsafe.Pointer(u.UpdateDataPlayer().Player)
		weapons = layout.As[server.Player](data).WeaponEquip
		armor = layout.As[server.Player](data).ArmorEquip
		coeff = *(*float32)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).Field57))
		if *(*byte)(unsafe.Pointer(&layout.As[server.Player](data).Field3680))&1 != 0 {
			return 0
		}
	} else if u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0 {
		weapons = layout.As[server.MonsterUpdateData](ud).Field514
		armor = layout.As[server.MonsterUpdateData](ud).Field515
		coeff = *(*float32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).Field518))
	} else {
		return 0
	}
	// the last damage record: whether it is set and the damage source kind
	var flagAt, kindAt *uint32
	if player {
		pd := layout.As[server.PlayerUpdateData](ud)
		flagAt, kindAt = &pd.Field76, &pd.Field75
	} else {
		md := layout.As[server.MonsterUpdateData](ud)
		flagAt, kindAt = &md.Field547, &md.Field546
	}
	*flagAt = 0
	if player && Nox_xxx_playerGetPossess_4DDF30(u) != nil {
		Nox_xxx_playerObserveClear_4DDEF0(u)
	}
	if u.Buffs&(1<<27) != 0 && actual != nil {
		if actual.ObjClass&1 != 0 && projectileFront(u, actual) {
			damageReflect(actual, u)
			if actual.ObjSubClass&0x40 == 0 {
				nox_xxx_unitClearOwner_4EC300(asObjectC(actual))
				nox_xxx_unitSetOwner_4EC290(asObjectC(u), asObjectC(actual))
			}
			if actual.ObjClass&1 != 0 && actual.ObjSubClass&2 != 0 {
				Nox_xxx_changeOwner_52BE40(actual, u)
			}
			if kind == 16 {
				projectileFX(132, u)
			}
			inventorySound(122, u, 0, 0)
			return 0
		} else if actual.ObjClass&1 == 0 && (kind == 16 || kind == 17) && projectileFront(u, actual) {
			if kind == 16 {
				projectileFX(132, u)
			}
			inventorySound(122, u, 0, 0)
			return 0
		}
	}
	if source != nil {
		if dword_5d4594_1563320 == 0 {
			dword_5d4594_1563320 = uint32(GetServer().S().Types.IndByID("SmallFist"))
			for i, name := range []string{"MediumFist", "LargeFist", "Meteor", "ToxicCloud", "SmallToxicCloud"} {
				*memmap.PtrUint32(0x5d4594, 1563324+uintptr(i*4)) = uint32(GetServer().S().Types.IndByID(name))
			}
		}
		eligible := uint32(actual.TypeInd) != uint32(dword_5d4594_1563320)
		end := 3
		if weapon != nil {
			end = 5
		}
		for i := 0; i < end; i++ {
			if uint32(actual.TypeInd) == memmap.Uint32(0x5d4594, 1563324+uintptr(i*4)) {
				eligible = false
			}
		}
		if weapon != nil && weapon != source || weapon == nil && (kind == 10 || kind == 2) {
			*flagAt = 1
			*kindAt = uint32(actual.TypeInd)
		}
		front := stateFront(&u.PosVec, int32(int16(u.Direction1)), &actual.PrevPos)&1 != 0
		if kind != 15 && eligible && front {
			var state byte
			if player {
				state = *(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).State))
			} else {
				// the original reads the same offset (88) of the monster data
				state = *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](ud).Path[9].Y))
			}
			shield := (player && state == 16 || !player && monsterControlHead(u) == 21) && armor&0x3000000 != 0
			if !shield && weapons&0x400 == 0 && state == 1 && controlActionState(u) == 45 && armor&0x3000000 != 0 && gameex_flags&0x10 != 0 {
				shield = true
			}
			if shield {
				if kind != 9 && kind != 17 {
					inventorySound(878, u, 0, 0)
					if actual.ObjClass&1 != 0 && actual.ObjSubClass&0x70 == 0 {
						damageReflect(actual, u)
						if actual.ObjClass&1 != 0 && actual.ObjSubClass&2 == 0 {
							nox_xxx_unitClearOwner_4EC300(asObjectC(actual))
							nox_xxx_unitSetOwner_4EC290(asObjectC(u), asObjectC(actual))
						}
					}
					value := float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("ItemDamageFromBlockPercentage"))) * float64(amount))
					damageBlockingItem(u, source, weapon, 2, value, kind, false)
					return 0
				}
			} else {
				canBlock := state == 13 || state == 18 || state == 19 || state == 20
				if !player {
					canBlock = monsterCanBlockHead(u)
				}
				if weapons&0x400 != 0 && (actual.ObjClass&1 != 0 || kind == 0 || kind == 11) && canBlock {
					if actual.ObjClass&1 != 0 {
						damageReflect(actual, u)
						if actual.ObjClass&1 != 0 && actual.ObjSubClass&2 == 0 {
							nox_xxx_unitClearOwner_4EC300(asObjectC(actual))
							nox_xxx_unitSetOwner_4EC290(asObjectC(u), asObjectC(actual))
						}
					}
					inventorySound(890, u, 0, 0)
					if player {
						nox_xxx_playerSetState_4FA020(asObjectC(u), int(int32(nox_common_randomInt_415FA0(18, 20))))
					} else {
						monsterControlEnsure(u, 23)
					}
					value := float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("ItemDamageFromBlockPercentage"))) * float64(amount))
					damageBlockingItem(u, source, weapon, 1024, value, kind, true)
					return 0
				}
				skipGreat := weapons&0x400 != 0 && (actual.ObjClass&1 != 0 || kind == 0 || kind == 11) && !canBlock
				if player && state == 0 && gameex_flags&4 != 0 {
					canBlock = true
				}
				if !skipGreat && weapons&0x7ff8000 != 0 && (kind == 0 || kind == 11) && actual.ObjClass&1 == 0 && canBlock {
					inventorySound(894, u, 0, 0)
					if player {
						nox_xxx_playerSetState_4FA020(asObjectC(u), 21)
					} else {
						monsterControlEnsure(u, 23)
					}
					value := float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("ItemDamageFromBlockPercentage"))) * float64(amount))
					damageBlockingItem(u, source, weapon, 134184960, value, kind, true)
					return 0
				}
			}
		}
	}
	n := amount
	recordKind := false
	switch kind {
	case 0, 3, 7, 8, 10, 11:
		damageFraction(u, &n, float32((1-float64(coeff))*float64(amount)))
		damageInventory(u, source, weapon, amount-n, kind)
		recordKind = true
	case 1, 12:
		damageInventory(u, source, weapon, amount, kind)
		recordKind = true
	case 2:
		damageFraction(u, &n, float32((1-float64(coeff)*.5)*float64(amount)))
		damageInventory(u, source, weapon, amount-n, 2)
		if *flagAt == 0 {
			*flagAt = 2
			*kindAt = 2
		}
	case 4, 5, 6, 14, 15, 16:
		recordKind = true
	case 9, 17:
		damageFraction(u, &n, float32(damageConductivity(u)*float64(amount)))
		damageInventory(u, source, weapon, amount, kind)
		recordKind = true
	}
	if recordKind && *flagAt == 0 {
		*flagAt = 2
		*kindAt = uint32(kind) // the original passes the type as the bits of a float argument
	}
	if amount > 0 && n == 0 {
		n = 1
	}
	if bool(noxflags.HasEngine(noxflags.EngineGodMode)) && u.ObjClass&4 != 0 {
		return 1
	}
	if bool(nox_common_gameFlags_check_40A5C0(4096)) {
		old := n
		n = floatToInt32(float32(float64(questRuntimeFloat(202032)) * float64(n)))
		if old > 0 && n < 1 {
			n = 1
		}
	}
	return damageDefault(u, source, weapon, n, kind)
}
