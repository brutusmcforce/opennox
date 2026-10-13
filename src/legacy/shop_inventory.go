package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func shopRepairQuote(u *server.Object, s *shopSession, code uint32) uint32 {
	var item *server.Object
	for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode == code {
			break
		}
	}
	if item == nil {
		return 0
	}
	charge := false
	if item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x47f0000 != 0 {
		data := item.UseDataWand()
		charge = data.Charges < data.MaxCharges && data.MaxCharges != 0
	}
	hp := item.HealthData
	if (hp != nil && hp.Cur != hp.Max && hp.Max != 0) || charge {
		b := [8]byte{0xc9, 31}
		binary.LittleEndian.PutUint16(b[2:], uint16(code))
		binary.LittleEndian.PutUint32(b[4:], uint32(shopPrice(2, s, item)))
		return shopSend(u, b[:], 0)
	}
	GetServer().S().Audio.EventObj(925, u, 0, 0)
	return memmap.PtrWord(item.CObj())
}
func shopRepair(u *server.Object, s *shopSession, code uint32) uint32 {
	gold := shopGetGold(u)
	if u.InvFirstItem == nil {
		return gold
	}
	var item *server.Object
	for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode == code {
			break
		}
	}
	if item == nil {
		return code
	}
	shopSubGold(u, uint32(shopPrice(2, s, item)))
	resourceSetHP(item, item.HealthData.Max)
	ind := int32(uint8(u.UpdateDataPlayer().Player.PlayerInd))
	gameplayReportItemHealth(int(ind), item)
	if item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x47f0000 != 0 {
		data := item.UseDataWand()
		if effectsRecharge(item, 100) != 0 {
			gameplayReportCharges(int(ind), item, data.Charges, data.MaxCharges)
		}
	}
	gameplayReportPlayerStat(int(ind), u)
	GetServer().S().Audio.EventObj(803, u, 2, u.NetCode)
	return code
}
func shopSell(u *server.Object, s *shopSession, typ int32, count uint32) {
	shopGetGold(u)
	for i := uint32(0); i < count; i++ {
		var item *server.Object
		for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
			if int32(item.TypeInd) == typ {
				break
			}
		}
		if item == nil {
			return
		}
		inventoryRemove(u, item)
		GetServer().DelayedDelete(item)
		shopAddGold(u, uint32(shopPrice(0, s, item)))
		gameplayReportPlayerStat(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), u)
		if i+1 == count {
			GetServer().S().Audio.EventObj(307, u, 2, u.NetCode)
			return
		}
	}
}
func shopLoad(s *shopSession) {
	var vendor *server.Object
	for _, u := range s.Units {
		if u != nil && u.ObjClass&2 != 0 && u.ObjSubClass&8 != 0 {
			vendor = u
			break
		}
	}
	if vendor == nil {
		return
	}
	core := GetServer().S()
	if noxflags.HasGame(noxflags.GameModeQuest) {
		cutoff := uint32(floatToInt32(float32(core.Balance.Float("ShopAnkhCutoffStage"))))
		if uint32(Nox_game_getQuestStage_4E3CC0()) < cutoff {
			if u := core.NewObjectByTypeID("AnkhTradable"); u != nil {
				shopAdd(s, u)
			}
		}
		for off := uintptr(234816); ; off += 4 {
			name := memmap.LoadPtr(0x587000, off)
			if name == nil {
				break
			}
			u := asObjectS(nox_xxx_newObjectByTypeID_4E3810((*int8)(name)))
			if u != nil {
				shopAdd(s, u)
			}
		}
		marker := core.NewObjectByTypeID("RewardMarker")
		if marker == nil {
			return
		}
		stage := uint32(Nox_game_getQuestStage_4E3CC0()) + 2
		for _, category := range [...]uint32{8, 8, 8, 8, 16, 16, 16, 16, 1, 1, 1, 1, 1, 1, 4, 4, 4} {
			marker.InitDataRewardMarker().Categories = category
			if u := rewardMarker(marker, stage); u != nil {
				shopAdd(s, u)
			}
		}
		if core.Rand.Logic.Int(0, 100) > 90 {
			marker.InitDataRewardMarker().Categories = 2
			if u := rewardMarker(marker, stage); u != nil {
				shopAdd(s, u)
			}
		}
		GetServer().DelayedDelete(marker)
		return
	}
	data := vendor.InitData
	if data == nil {
		return
	}
	stock := layout.As[server.ShopkeeperInitData](data)
	for i := 0; i < int(stock.StockCount); i++ {
		e := &stock.Stock[i]
		for j := 0; j < int(e.Count); j++ {
			u := asObjectS(nox_xxx_newObjectWithTypeInd_4E3450(int(int32(e.Type))))
			if u == nil {
				continue
			}
			if uint32(u.ObjClass)&0x13001000 != 0 {
				// C copied an uninitialized fifth word. Define only that former
				// indeterminate tail as zero; keep the four actual modifiers.
				mods := server.ModifierInitData{Modifiers: e.Modifiers}
				stateAttributes(u, &mods)
			}
			if u.Xfer == xferIdentityKey(xferIDSpellReward) {
				u.UseDataSpellReward().Spell = byte(e.Reward)
			}
			if u.Xfer == xferIdentityKey(xferIDAbilityReward) {
				u.UseDataAbilityReward().Ability = byte(e.Reward)
			}
			if u.Xfer == xferIdentityKey(xferIDFieldGuide) {
				alloc.Strcpy(unsafe.Pointer(&u.UseDataFieldGuide().CreatureBuf[0]), unsafe.Pointer(nox_xxx_getUnitNameByThingType_4E3A80(int(int32(e.Reward)))))
			}
			shopAdd(s, u)
		}
	}
}
