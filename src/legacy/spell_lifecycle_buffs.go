package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func spellLifeHasBuff(u *server.Object, buff int32) bool {
	return u.HasEnchant(server.EnchantID(buff))
}
func spellLifeBuffTimer(u *server.Object, buff int32) int32 {
	return int32(u.EnchantDur(server.EnchantID(buff)))
}
func spellLifeBuffPower(u *server.Object, buff int32) int8 {
	return int8(u.EnchantPower(server.EnchantID(buff)))
}
func spellLifeBuffAudio(u *server.Object, buff, phase int32) {
	id := server.EnchantID(buff).Spell()
	aud := GetServer().S().Spells.DefByInd(id).GetAudio(int(phase))
	GetServer().S().Audio.EventObj(aud, u, 0, 0)
}
func spellLifeBuffOff(u *server.Object, buff int32) uint32 {
	mask := uint32(1) << uint32(buff)
	if u.Buffs&mask == 0 {
		return mask
	}
	stateBuffs(u, u.Buffs&^mask)
	u.BuffsDur[buff] = 0
	u.BuffsPower[buff] = 0
	if buff != 16 && buff != 30 {
		spellLifeBuffAudio(u, buff, 2)
	}
	return 0
}
func spellLifeClearBuffs(u *server.Object) {
	stateBuffs(u, 0)
	clear(u.BuffsDur[:])
	clear(u.BuffsPower[:])
}
func spellLifeApplyBuff(u *server.Object, buff int32, dur int16, power int8) {
	hec := *memmap.PtrUint32(0x5d4594, 1569740)
	if hec == 0 {
		hec = stateType(1569740, "Hecubah")
		*memmap.PtrUint32(0x5d4594, 1569744) = uint32(GetServer().S().Types.IndByID("Necromancer"))
	}
	nec := *memmap.PtrUint32(0x5d4594, 1569744)
	if u == nil {
		return
	}
	if uint32(u.TypeInd) == hec && buff == 29 {
		return
	}
	if controlFlags(4096) && uint32(u.TypeInd) == hec && buff == 3 {
		GetServer().S().Audio.EventObj(582, u, 0, 0)
		return
	}
	if controlFlags(4096) && uint32(u.TypeInd) == nec && buff == 3 {
		GetServer().S().Audio.EventObj(595, u, 0, 0)
		return
	}
	if u.ObjClass&2 != 0 && u.ObjSubClass&0x1000 != 0 && buff == 11 && !controlFlags(2048) {
		if uint32(u.TypeInd) == hec {
			GetServer().S().Audio.EventObj(582, u, 0, 0)
		} else if uint32(u.TypeInd) == nec {
			GetServer().S().Audio.EventObj(595, u, 0, 0)
		}
		return
	}
	if u.ObjFlags&0x8022 != 0 {
		return
	}
	if spellLifeHasBuff(u, buff) && spellLifeBuffTimer(u, buff) == 0 {
		return
	}
	if buff != 0 {
		spellLifeBuffOff(u, 0)
	}
	u.BuffsDur[buff] = uint16(dur)
	u.BuffsPower[buff] = byte(power)
	stateBuffs(u, u.Buffs|(uint32(1)<<uint32(buff)))
	spellLifeBuffAudio(u, buff, 1)
}
func spellLifeUpdateBuffs(u *server.Object) {
	if u.Buffs == 0 {
		return
	}
	fps := GetServer().S().TickRate()
	for i := int32(0); i < 32; i++ {
		if !spellLifeHasBuff(u, i) {
			continue
		}
		if i == 16 && uint32(u.BuffsDur[16])%uint32(fps) == uint32(fps)-1 {
			GetServer().S().Audio.EventObj(26, u, 0, 0)
		}
		if u.BuffsDur[i] == 0 {
			continue
		}
		u.BuffsDur[i]--
		if u.BuffsDur[i] != 0 {
			continue
		}
		if i == 7 {
			u.ObjFlags &^= 0x40
		} else if i == 16 {
			u.Obj130 = nil
			u.Field131 = 13
			resourceDamage(u, 9999999)
			GetServer().S().Audio.EventObj(779, u, 0, 0)
			if u.ObjClass&4 != 0 {
				gameplayReportEliminationDeath(u)
				gameplayReportLesson(u)
			}
		}
		spellLifeBuffOff(u, i)
		u.BuffsPower[i] = 0
	}
	if spellLifeHasBuff(u, 9) {
		u.SpeedCur = float32(float64(u.SpeedCur) * 1.25)
	}
}
