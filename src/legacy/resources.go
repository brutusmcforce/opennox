package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

// Return values preserve the original ABI, including checksum and pointer bits.
func resourceSetHP(u *server.Object, amount uint16) uint32 {
	if u.HealthData == nil {
		return 0
	}
	u.NeedSync()
	if h := u.InvHolder; h != nil && h.ObjClass&4 != 0 {
		toggleProtectionObject(int32(h.UpdateDataPlayer().Player.Prot4632), u)
	}
	u.HealthData.Cur = amount
	if u.ObjClass&4 != 0 {
		setProtectionRecord(int32(u.UpdateDataPlayer().Player.ProtUnitHPCur), uint32(amount))
	}
	if u.ObjClass&2 != 0 && u.ObjSubClass&0x80 != 0 {
		if u.ObjClass&0x20400004 != 0 {
			for i, v := range u.Field140 {
				u.Field140[i] = v&0xfffff000 | 0x20000
			}
		} else {
			u.Sub_4E4500(0x20000, 2, u.Sub_4E4C90(2))
		}
	}
	if h := u.InvHolder; h != nil {
		if h.ObjClass&4 != 0 {
			return toggleProtectionObject(int32(h.UpdateDataPlayer().Player.Prot4632), u)
		}
		return memmap.PtrWord(h.CObj())
	}
	return 0
}
func resourceInformOwner(u *server.Object) {
	if u != nil && u.ObjOwner != nil && u.ObjOwner.ObjClass&4 != 0 {
		gameplayReportCurrentHealth(int(int32(uint8(u.ObjOwner.UpdateDataPlayer().Player.PlayerInd))), u)
	}
}
func resourceAdjustHP(u *server.Object, dv int32) {
	if noxflags.HasGame(0x4000000) || u.HealthData == nil {
		return
	}
	h := u.HealthData
	if h.Cur < h.Max {
		v := int32(h.Cur) + dv
		val := uint16(v)
		if v > int32(h.Max) {
			val = h.Max
		}
		resourceSetHP(u, val)
		if u.ObjClass&2 != 0 {
			resourceInformOwner(u)
		}
	}
}
func resourceDamage(u *server.Object, amount int32) {
	if u == nil || u.HealthData == nil || u.HealthData.Max == 0 || (u.ObjClass&4 != 0 && noxflags.HasEngine(noxflags.EngineGodMode)) {
		return
	}
	if u.ObjClass&4 != 0 && u.UpdateData != nil {
		p := u.UpdateDataPlayer()
		if p.Player.PlayerClass() == 0 && p.HarpoonTarg != nil {
			Nox_xxx_harpoonBreakForPlr_537520(u)
		}
	}
	if int32(u.HealthData.Cur) > amount {
		resourceSetHP(u, uint16(int32(u.HealthData.Cur)-amount))
	} else {
		resourceSetHP(u, 0)
		if u.ObjFlags&0x8000 == 0 {
			u.ObjFlags |= 0x8000
			Nox_xxx_spellBuffOff_4FF5B0(u, 16)
			if !monsterIsZombie(u) {
				unitMonsterReward(u)
			}
			if u.ObjClass&2 != 0 {
				monsterControlDeath(u)
			} else if u.Death != nil {
				u.CallDeath()
			} else {
				GetServer().DelayedDelete(u)
			}
		}
	}
	if u.ObjClass&2 != 0 {
		resourceInformOwner(u)
	}
}
func resourceRestoreHP(u *server.Object) {
	if u == nil || u.HealthData == nil {
		return
	}
	resourceSetHP(u, u.HealthData.Max)
	u.HealthData.Field2 = u.HealthData.Cur
	if u.ObjClass&2 != 0 {
		resourceInformOwner(u)
	}
}
func resourceHPHistory(u *server.Object) {
	if u == nil || u.ObjClass&4 == 0 || u.HealthData == nil {
		return
	}
	ud := layout.As[server.PlayerUpdateData](u.UpdateData)
	for i := range ud.HealthPrev {
		ud.HealthPrev[i] = u.HealthData.Cur
	}
	u.UpdateDataPlayer().Field19_0 = u.HealthData.Cur
}
func resourceGetHP(u *server.Object) int16 {
	if u == nil || u.HealthData == nil {
		return 0
	}
	return int16(u.HealthData.Cur)
}
func resourceGetMaxHP(u *server.Object) int16 {
	if u == nil || u.HealthData == nil {
		return 0
	}
	return int16(u.HealthData.Max)
}
func resourceSetMaxHP(u *server.Object, v uint16) uint32 {
	if u == nil {
		return 0
	}
	if u.HealthData != nil {
		u.HealthData.Max = v
	}
	return memmap.PtrWord(unsafe.Pointer(u.HealthData))
}

func resourcePoison(u *server.Object, amount, max int32) bool {
	if u == nil {
		return false
	}
	old := int32(u.Poison540)
	if u.ObjFlags&2 != 0 || u.HasEnchant(23) || (u.ObjClass&4 != 0 && u.UpdateDataPlayer().Player.Field3680&1 != 0) || (u.ObjClass&2 != 0 && u.ObjSubClass&0x200 != 0) {
		return false
	}
	// Preserve the C double multiplication followed by its explicit float spill.
	chance := floatToInt32(float32(effectsProtection(u, modifierKey(modifierIDPoisonProtectEngage), 18, "PoisonSpellProtection", .69999999, .89999998) * 100))
	if int32(nox_common_randomInt_415FA0(0, 100)) < chance {
		resourcePriority(u, "Health.c:ResistPoison")
		return false
	}
	next := old + amount
	if next > max {
		if old <= max {
			next = max
		} else {
			next = old
		}
	}
	if next != old {
		resourceSetPoison(u, next)
		GetServer().S().Audio.EventObj(100, u, 0, 0)
	}
	if old == 0 && next > 0 && u.HealthData != nil {
		u.HealthData.Field16 = GetServer().S().Frame()
	}
	return true
}
func resourcePriority(u *server.Object, key string) {
	gameplayTextPrivate(u, (*byte)(unsafe.Pointer(internCStr(key))), 0)
}
func resourcePoisonReport(u *server.Object, active bool) {
	if u.ObjClass&2 == 0 {
		return
	}
	var owner *server.Object
	if noxflags.HasGame(2048) && u.ObjSubClass&0x10 != 0 {
		if p := GetServer().S().Players.ByInd(31); p != nil {
			owner = p.PlayerUnit
		}
	} else if u.ObjSubClass&0x80 != 0 {
		owner = u.ObjOwner
	}
	if owner != nil {
		gameplayReportPoison(owner, u, byte(bool2int(active)))
	}
}
func resourceClearPoison(u *server.Object, fade bool) {
	u.Poison540 = 0
	if u.HealthData != nil {
		u.HealthData.Field16 = 0
	}
	if u.ObjClass&4 != 0 {
		Nox_xxx_playerUnsetStatus_417530(u.UpdateDataPlayer().Player, 1024)
		if fade {
			resourcePriority(u, "Health.c:PoisonFade")
		}
	} else {
		resourcePoisonReport(u, false)
	}
}
func resourceReducePoison(u *server.Object, v int32) {
	if u == nil {
		return
	}
	if int32(u.Poison540) > v {
		u.Poison540 -= byte(v)
	} else {
		resourceClearPoison(u, true)
	}
}
func resourceRemovePoison(u *server.Object) {
	if u != nil && u.Poison540 != 0 {
		resourceClearPoison(u, false)
	}
}
func resourceSetPoison(u *server.Object, v int32) {
	if u == nil {
		return
	}
	if u.Poison540 == 0 && v > 0 && u.HealthData != nil {
		u.HealthData.Field16 = GetServer().S().Frame()
	}
	u.Poison540 = byte(v)
	if u.ObjClass&4 != 0 {
		p := u.UpdateDataPlayer().Player
		if v != 0 {
			Nox_xxx_netNeedTimestampStatus_4174F0(p, 1024)
		} else {
			Nox_xxx_playerUnsetStatus_417530(p, 1024)
		}
	} else {
		resourcePoisonReport(u, v != 0)
	}
	u.Field542 = 0
	if v != 0 {
		u.Field542 = 1000
	}
}
func resourceAddMana(u *server.Object, v int16) uint16 {
	if u == nil || u.ObjClass&4 == 0 {
		return uint16(memmap.PtrWord(unsafe.Pointer(u)))
	}
	p := u.UpdateDataPlayer()
	p.ManaPrev = p.ManaCur
	p.ManaCur += uint16(v)
	if p.ManaCur > p.ManaMax {
		p.ManaCur = p.ManaMax
	}
	addProtectionRecord(int32(p.Player.ProtUnitManaCur), uint32(int32(v)))
	out := p.ManaMax
	if p.ManaCur > out {
		out = uint16(setProtectionRecord(int32(p.Player.ProtUnitManaCur), uint32(out)))
	}
	return out
}
func resourceSubMana(u *server.Object, v int32) uint32 {
	if u == nil || u.ObjClass&4 == 0 {
		return memmap.PtrWord(unsafe.Pointer(u))
	}
	p := u.UpdateDataPlayer()
	if noxflags.HasEngine(noxflags.EngineGodMode) {
		return memmap.PtrWord(u.UpdateData)
	}
	p.ManaPrev = p.ManaCur
	if int32(p.ManaCur) > v {
		p.ManaCur = uint16(int32(p.ManaCur) - v)
	} else {
		p.ManaCur = 0
	}
	// The original compares the UPDATED mana, and truncates the delta to short.
	delta := int16(-int32(p.ManaCur))
	if int32(p.ManaCur) > v {
		delta = -int16(v)
	}
	return addProtectionRecord(int32(p.Player.ProtUnitManaCur), uint32(int32(delta)))
}
func resourceGetMana(u *server.Object) int16 {
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 0 {
		return int16(u.UpdateDataPlayer().ManaCur)
	}
	if u.ObjClass&2 != 0 {
		return 1000
	}
	return 0
}
func resourceGetMaxMana(u *server.Object) int16 {
	if u == nil || u.ObjClass&4 == 0 {
		return 0
	}
	return int16(u.UpdateDataPlayer().ManaMax)
}
func resourceSetMaxMana(u *server.Object, v uint16) uint32 {
	if u == nil || u.ObjClass&4 == 0 {
		return memmap.PtrWord(unsafe.Pointer(u))
	}
	u.UpdateDataPlayer().ManaMax = v
	return memmap.PtrWord(u.UpdateData)
}
func resourceRefreshMana(u *server.Object) uint32 {
	if u == nil || u.ObjClass&4 == 0 {
		return memmap.PtrWord(unsafe.Pointer(u))
	}
	p := u.UpdateDataPlayer()
	p.ManaPrev = p.ManaCur
	p.ManaCur = p.ManaMax
	return addProtectionRecord(int32(p.Player.ProtUnitManaCur), uint32(int32(int16(p.ManaMax))))
}
func resourceAddGold(u *server.Object, v uint32) uint32 {
	p := u.UpdateDataPlayer().Player
	p.GoldVal += v
	return addProtectionRecord(int32(p.ProtPlayerGold), v)
}
func resourceSubGold(u *server.Object, v uint32) uint32 {
	p := u.UpdateDataPlayer().Player
	if p.GoldVal >= v {
		p.GoldVal -= v
	} else {
		p.GoldVal = 0
	}
	return addProtectionRecord(int32(p.ProtPlayerGold), -v)
}
func resourceSetGold(u *server.Object, v int32) {
	if u == nil || u.ObjClass&4 == 0 {
		return
	}
	p := u.UpdateDataPlayer().Player
	if v >= 0 || p.GoldVal >= uint32(-v) {
		p.GoldVal += uint32(v)
		addProtectionRecord(int32(p.ProtPlayerGold), uint32(v))
	} else {
		p.GoldVal = 0
		setProtectionRecord(int32(p.ProtPlayerGold), 0)
	}
}
func resourceGetGold(u *server.Object) uint32 { return u.UpdateDataPlayer().Player.GoldVal }
func resourceObjectGold(u *server.Object) uint32 {
	if u == nil || u.ObjClass&4 == 0 {
		return 0
	}
	return resourceGetGold(u)
}
func resourceGoldPickup(u, item *server.Object, flags int) bool {
	if u.ObjClass&4 != 0 {
		gold := (*uint32)(item.InitData)
		resourceAddGold(u, *gold)
		GetServer().DelayedDelete(item)
		text := GetServer().S().Strings().GetStringInFile(strman.ID("GoldPickup"), `C:\NoxPost\src\Server\Object\pickdrop\pickup.c`)
		textFormatLine(u, (*uint16)(unsafe.Pointer(internWStr(text))), textFormatWord(uint32(*gold)))
	} else if !Nox_xxx_pickupDefault_4F31E0(u, item, flags, 0) {
		return false
	}
	GetServer().S().Audio.EventObj(307, u, 0, 0)
	return true
}
