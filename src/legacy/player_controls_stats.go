package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/server"
)

func controlGiveAbilities(u *server.Object, level int8, notify int32) {
	if u == nil || level <= 0 {
		return
	}
	pl := controlPlayer(u)
	for i := 0; i < int(level); i++ {
		ability := *memmap.PtrUint32(0x587000, 206108+uintptr(4*i))
		if ability == 0 {
			continue
		}
		if controlFlags(4096) || questRuntimeWord(1556160) != 0 || questRuntimeWord(1556164) != 0 {
			pl.SpellLvl[i] = 0
		} else {
			bookAwardAbility(u, int32(ability), int32(notify))
		}
	}
}

// controlReadStats preserves the C double intermediates and explicit float stores.
// Integer fields and their protection records are updated in the original order.
func controlReadStats(u *server.Object, notify int32) int32 {
	d := u.UpdateData
	pl := controlPlayer(u)
	class := pl.InfoBytes[66]
	s := GetServer().S()
	base := s.Players.BaseStats()
	max := s.Players.ClassStats(player.Class(class))
	warrior := s.Players.ClassStats(player.Warrior)
	// put writes an unaligned word of the player info (strength at 2239, speed at 2235)
	put := func(off int, x float32) {
		*(*uint32)(unsafe.Pointer(&pl.InfoBytes[off-2185])) = uint32(floatToInt32(x))
	}
	hp := &u.HealthData
	if controlFlags(8192) {
		(*hp).Max = uint16(floatToInt32(max.Health))
		resourceSetHP(u, uint16(floatToInt32(float32(math.Abs(float64(max.Health))))))
		layout.As[server.PlayerUpdateData](d).ManaMax = uint16(floatToInt32(max.Mana))
		layout.As[server.PlayerUpdateData](d).ManaCur = layout.As[server.PlayerUpdateData](d).ManaMax
		put(2239, max.Strength)
		u.SpeedBase = float32(float64(max.Speed) * 0.000099999997)
		put(2235, max.Speed)
		if class == 0 && !controlFlags(4096) && sub_4D6F30() == 0 {
			controlGiveAbilities(u, 10, 0)
		}
	} else {
		level := int8(pl.Level)
		if level > 10 {
			level = 10
		}
		n := float64(int32(level) - 1)
		interp := func(a, b float32) float64 { return (float64(a)-float64(b))*n/9 + float64(b) }
		(*hp).Max = uint16(floatToInt32(float32(interp(max.Health, base.Health) + 0.5)))
		resourceSetHP(u, (*hp).Max)
		mana := interp(max.Mana, base.Mana)
		if mana > float64(max.Mana) {
			mana = float64(max.Mana)
		}
		layout.As[server.PlayerUpdateData](d).ManaMax = uint16(floatToInt32(float32(mana + 0.5)))
		layout.As[server.PlayerUpdateData](d).ManaCur = layout.As[server.PlayerUpdateData](d).ManaMax
		put(2239, float32(interp(max.Strength, base.Strength)+0.5))
		speed := interp(max.Speed, base.Speed)
		u.SpeedBase = float32(speed * 0.000099999997)
		put(2235, float32(speed+0.5))
		if class == 0 {
			controlGiveAbilities(u, level, notify)
		}
	}
	strength := float64(int32(*(*uint32)(unsafe.Pointer(&pl.InfoBytes[54])))) / float64(warrior.Strength)
	u.Mass = float32(strength*20 + 10)
	carry := int64((strength*1250 + 750) * *memmap.PtrFloat64(0x581450, 10216))
	*(*uint16)(unsafe.Pointer(&pl.Field3652)) = uint16(carry)
	u.CarryCapacity = uint16(carry)
	setProtectionRecord(int32(pl.ProtPlayerField2239), pl.Info().Field2239())
	setProtectionRecord(int32(pl.ProtPlayerField2235), pl.Info().Field2235())
	setProtectionRecord(int32(pl.ProtUnitManaMax), uint32(layout.As[server.PlayerUpdateData](d).ManaMax))
	setProtectionRecord(int32(pl.ProtUnitHPMax), uint32((*hp).Max))
	var weight uint32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		weight += uint32(it.Weight)
	}
	pl.Field3656 = 0
	if weight > uint32(u.CarryCapacity) {
		pl.Field3656 = 1
	}
	name := unsafe.Pointer(&pl.InfoBytes)
	n := 0
	for *controlHalf(name, n) != 0 {
		n += 2
	}
	result := protectionValidateString(name, uint32(n), int32(pl.ProtPlayerOrigName))
	pl.Field2184 = 1
	return int32(result)
}
func controlLevelFromXP(u *server.Object) int32 {
	pl := controlPlayer(u)
	level := 10
	if !controlFlags(8192) {
		i := 0
		for ; i <= 10; i++ {
			if float64(nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int32(i)))) > float64(u.Experience) {
				break
			}
		}
		level = i - 1
	}
	pl.Level = byte(level)
	setProtectionRecord(int32(pl.ProtPlayerLevel), uint32(byte(level)))
	return controlReadStats(u, 0)
}
func controlSetLevel(u *server.Object, level byte) {
	pl := controlPlayer(u)
	if int8(level) > 10 {
		level = 10
	}
	xp := nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int8(level)))
	u.Experience = float32(xp)
	xp = nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int8(level)))
	// XP protection stores a truncated numeric value, not its float bits.
	updateProtectionFloat(int32(pl.ProtUnitExperience), float32(xp), false)
	gameplayReportExperience(u)
	pl.Level = level
	setProtectionRecord(int32(pl.ProtPlayerLevel), uint32(level))
	controlReadStats(u, 0)
	if controlFlags(2048) {
		if pl.InfoBytes[66] == 0 {
			for i := 1; i < 6; i++ {
				if pl.SpellLvl[i] != 0 {
					quickbarBookSlot(3, uint32(i), i-1)
				}
			}
		}
		runtimePauseStart(u, 0)
	}
}
