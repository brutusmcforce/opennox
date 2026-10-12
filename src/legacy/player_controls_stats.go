package legacy

import (
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
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
			*equipmentWord(pl, 3696+4*i) = 0
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
	class := *controlByte(pl, 2251)
	s := GetServer().S()
	base := s.Players.BaseStats()
	max := s.Players.ClassStats(player.Class(class))
	warrior := s.Players.ClassStats(player.Warrior)
	put := func(off int, x float32) { *equipmentWord(pl, off) = uint32(floatToInt32(x)) }
	hp := controlPtr(u.CObj(), 556)
	if controlFlags(8192) {
		*controlHalf(*hp, 4) = uint16(floatToInt32(max.Health))
		resourceSetHP(u, uint16(floatToInt32(float32(math.Abs(float64(max.Health))))))
		*controlHalf(d, 8) = uint16(floatToInt32(max.Mana))
		*controlHalf(d, 4) = *controlHalf(d, 8)
		put(2239, max.Strength)
		*(*float32)(unsafe.Add(u.CObj(), 548)) = float32(float64(max.Speed) * 0.000099999997)
		put(2235, max.Speed)
		if class == 0 && !controlFlags(4096) && sub_4D6F30() == 0 {
			controlGiveAbilities(u, 10, 0)
		}
	} else {
		level := int8(*controlByte(pl, 3684))
		if level > 10 {
			level = 10
		}
		n := float64(int32(level) - 1)
		interp := func(a, b float32) float64 { return (float64(a)-float64(b))*n/9 + float64(b) }
		*controlHalf(*hp, 4) = uint16(floatToInt32(float32(interp(max.Health, base.Health) + 0.5)))
		resourceSetHP(u, *controlHalf(*hp, 4))
		mana := interp(max.Mana, base.Mana)
		if mana > float64(max.Mana) {
			mana = float64(max.Mana)
		}
		*controlHalf(d, 8) = uint16(floatToInt32(float32(mana + 0.5)))
		*controlHalf(d, 4) = *controlHalf(d, 8)
		put(2239, float32(interp(max.Strength, base.Strength)+0.5))
		speed := interp(max.Speed, base.Speed)
		*(*float32)(unsafe.Add(u.CObj(), 548)) = float32(speed * 0.000099999997)
		put(2235, float32(speed+0.5))
		if class == 0 {
			controlGiveAbilities(u, level, notify)
		}
	}
	strength := float64(int32(*equipmentWord(pl, 2239))) / float64(warrior.Strength)
	*(*float32)(unsafe.Add(u.CObj(), 120)) = float32(strength*20 + 10)
	carry := int64((strength*1250 + 750) * *memmap.PtrFloat64(0x581450, 10216))
	*controlHalf(pl, 3652) = uint16(carry)
	*controlHalf(u.CObj(), 490) = uint16(carry)
	for _, v := range [][2]int{{4624, 2239}, {4620, 2235}} {
		setProtectionRecord(int32(*equipmentWord(pl, v[0])), *equipmentWord(pl, v[1]))
	}
	setProtectionRecord(int32(*equipmentWord(pl, 4600)), uint32(*controlHalf(d, 8)))
	setProtectionRecord(int32(*equipmentWord(pl, 4592)), uint32(*controlHalf(*hp, 4)))
	var weight uint32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		weight += uint32(*controlByte(it.CObj(), 488))
	}
	*equipmentWord(pl, 3656) = 0
	if weight > uint32(*controlHalf(u.CObj(), 490)) {
		*equipmentWord(pl, 3656) = 1
	}
	name := unsafe.Add(pl, 2185)
	n := 0
	for *controlHalf(name, n) != 0 {
		n += 2
	}
	result := protectionValidateString(name, uint32(n), int32(*equipmentWord(pl, 4628)))
	*controlByte(pl, 2184) = 1
	return int32(result)
}
func controlLevelFromXP(u *server.Object) int32 {
	pl := controlPlayer(u)
	level := 10
	if !controlFlags(8192) {
		i := 0
		for ; i <= 10; i++ {
			if float64(nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int32(i)))) > float64(*(*float32)(unsafe.Add(u.CObj(), 28))) {
				break
			}
		}
		level = i - 1
	}
	*controlByte(pl, 3684) = byte(level)
	setProtectionRecord(int32(*equipmentWord(pl, 4644)), uint32(byte(level)))
	return controlReadStats(u, 0)
}
func controlSetLevel(u *server.Object, level byte) {
	pl := controlPlayer(u)
	if int8(level) > 10 {
		level = 10
	}
	xp := nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int8(level)))
	*(*float32)(unsafe.Add(u.CObj(), 28)) = float32(xp)
	xp = nox_xxx_gamedataGetFloatTable_419D70(internCStr("XPTable"), int(int8(level)))
	// XP protection stores a truncated numeric value, not its float bits.
	updateProtectionFloat(int32(*equipmentWord(pl, 4604)), float32(xp), false)
	gameplayReportExperience(u)
	*controlByte(pl, 3684) = level
	setProtectionRecord(int32(*equipmentWord(pl, 4644)), uint32(level))
	controlReadStats(u, 0)
	if controlFlags(2048) {
		if *controlByte(pl, 2251) == 0 {
			for i := 1; i < 6; i++ {
				if *equipmentWord(pl, 3696+4*i) != 0 {
					quickbarBookSlot(3, uint32(i), i-1)
				}
			}
		}
		runtimePauseStart(u, 0)
	}
}
