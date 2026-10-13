package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

func shopVendor(s *shopSession) *server.Object {
	if s.Units[0].ObjClass&4 != 0 {
		return s.Units[1]
	}
	return s.Units[0]
}

// The original price routine stores its running values in float32, while
// intermediate products/divisions use the hosted x87 precision-53 arithmetic.
// Keep those stores explicit: in particular a quest modifier product remains
// double until it is added, and the quest sell multiplier is not a float local.
func shopPrice(mode int, s *shopSession, u *server.Object) int32 {
	core := GetServer().S()
	value := float32(u.Worth)
	diamond := memmap.PtrUint32(0x5D4594, 2386504)
	if *diamond == 0 {
		*diamond = uint32(core.Types.IndByID("Diamond"))
		*memmap.PtrUint32(0x5D4594, 2386508) = uint32(core.Types.IndByID("Ruby"))
		*memmap.PtrUint32(0x5D4594, 2386512) = uint32(core.Types.IndByID("Emerald"))
	}
	var vendor unsafe.Pointer
	if s != nil && s.Kind != 0 {
		vendor = shopVendor(s).InitData
	}
	class, sub := uint32(u.ObjClass), uint32(u.ObjSubClass)
	quest := noxflags.HasGame(noxflags.GameModeQuest)
	if class&0x100 != 0 {
		if sub&1 != 0 {
			value = float32(uint16(core.Spells.Price(spell.ID(u.UseDataSpellReward().Spell))))
		} else if sub&2 != 0 {
			if typ := core.Types.ByID(u.UseDataFieldGuide().Creature()); typ != nil && typ.Worth >= 0 {
				value = float32(typ.Worth)
			}
			if quest {
				value = float32(core.Balance.Float("QuestGuideWorthMultiplier") * float64(value))
			}
		}
	}
	if class&0x13001000 != 0 {
		for _, mod := range itemModifiers(u) {
			if mod == nil {
				continue
			}
			price := float64(float32(mod.Price20))
			if quest {
				price *= core.Balance.Float("QuestModifierWorthMultiplier")
			}
			value = float32(price + float64(value))
		}
	}
	base := value
	if class&0x1000000 != 0 && sub&0x82 != 0 {
		use := u.UseDataAmmo()
		if use.Infinite == 0 && use.MaxCharges != 0 {
			key := "DefaultAmmoAmount"
			if quest {
				key = "DefaultAmmoAmountQuest"
			}
			def := floatToInt32(float32(core.Balance.Float(key)))
			if int32(use.Charges) > def {
				per := float32(float64(value) / float64(def))
				value = float32(float64(int32(use.Charges)-def)*float64(per) + float64(value))
			} else {
				value = float32(float64(use.Charges) / float64(def) * float64(value))
			}
		}
	} else if class&0x1000 != 0 && sub&0x47f0000 != 0 {
		use := u.UseDataWand()
		if use.MaxCharges != 0 {
			value = float32(float64(use.Charges) / float64(use.MaxCharges) * float64(value))
		}
	}
	if vendor != nil && s != nil && s.Kind != 0 && mode != 0 {
		mult := float64((*server.ShopkeeperInitData)(vendor).Field1716)
		value, base = float32(float64(value)*mult), float32(float64(base)*mult)
	}
	if hp := u.HealthData; hp != nil && hp.Max != 0 {
		value = float32(float64(hp.Cur) / float64(hp.Max) * float64(value))
	}
	if vendor != nil && mode == 0 {
		mult := float64((*server.ShopkeeperInitData)(vendor).Field1720)
		if quest {
			mult = core.Balance.Float("QuestSellMultiplier")
		}
		typ := uint32(u.TypeInd)
		if typ != *diamond && typ != *memmap.PtrUint32(0x5D4594, 2386508) && typ != *memmap.PtrUint32(0x5D4594, 2386512) {
			value, base = float32(mult*float64(value)), float32(mult*float64(base))
		}
	}
	if quest && mode == 0 && class&0x3001000 != 0 && byte(u.UpdateDataWeaponArmor().Flags)&1 != 0 {
		base, value = 0, 1
	} else if !(value >= 1) { // includes NaN, matching the original forward jump
		value = 1
	}
	if base < 1 {
		base = 1
	}
	if mode == 2 {
		repair := core.Balance.Float("RepairCoefficient") * (float64(base) - float64(value))
		value = float32(repair)
		if repair < 1 {
			value = 1
		}
	}
	return floatToInt32(value)
}

func shopStockKey(n *shopItem) uint32 {
	mask, priority := uint32(16781312), uint32(255)
	off := uintptr(10308)
	for {
		sub := *memmap.PtrUint32(0x581450, off)
		if uint32(n.Object.ObjClass)&mask != 0 && (sub == 0 || sub&uint32(n.Object.ObjSubClass) != 0) {
			break
		}
		mask = *memmap.PtrUint32(0x581450, off+4)
		off += 8
		priority--
		if mask == 0 {
			break
		}
	}
	return n.Value | priority<<24
}

func shopStockMatches(u *server.Object, e *server.ShopStockEntry) int32 {
	if u == nil || e == nil || uint32(u.TypeInd) != e.Type {
		return 0
	}
	match := int32(1)
	if uint32(u.ObjClass)&0x13001000 != 0 {
		mods := itemModifiers(u)
		for i, mod := range mods {
			if mod != e.Modifiers[i] {
				match = 0
			}
		}
	}
	if u.ObjClass&0x100 != 0 {
		if u.ObjSubClass&5 != 0 {
			if uint32(bookRewardIndex(u)) != e.Reward {
				return 0
			}
		} else if u.UseDataFieldGuide().Creature() != GetServer().S().Types.ByInd(int(e.Reward)).ID() {
			match = 0
		}
	}
	return match
}
func shopStockIndex(s *shopSession, u *server.Object) int32 {
	data := layout.As[server.ShopkeeperInitData](shopVendor(s).InitData)
	for i := 0; i < int(data.StockCount); i++ {
		if shopStockMatches(u, &data.Stock[i]) != 0 {
			return int32(i)
		}
	}
	return -1
}
