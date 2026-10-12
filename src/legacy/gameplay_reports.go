package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func gameplayReportWord(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
func gameplayReportPtr(p unsafe.Pointer) int               { return int(uintptr(p)) }
func gameplayReportPlayer(u *server.Object) *server.Player { return (*server.Player)(controlPlayer(u)) }
func gameplayReportRecipient(u *server.Object) int         { return int(gameplayReportPlayer(u).PlayerInd) }
func gameplayReportCode(u *server.Object) uint16           { return uint16(GetServer().S().GetUnitNetCode(u)) }
func gameplayReportID(op byte, u *server.Object, size int) []byte {
	b := make([]byte, size)
	b[0] = op
	binary.LittleEndian.PutUint16(b[1:], gameplayReportCode(u))
	return b
}
func gameplayReportSend(to int, b []byte, ordered bool, priority int) int {
	var order byte
	if ordered {
		order = 1
	}
	return reliableEnqueue(to, b, nil, priority, order)
}
func gameplayReportCoalesce(to int, b []byte) int { return reliableCoalesce(to, b, nil, 1) }
func gameplayReportDirect(to int, b []byte) int {
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(to), netlist.Kind1, b))
}

func gameplayReportCreature(to int, order byte) int {
	return gameplayReportSend(to, []byte{237, order}, true, 1)
}
func gameplayReportRate(to int) int {
	return gameplayReportSend(to, []byte{236, byte(*memmap.PtrUint32(0x587000, 4728))}, true, 1)
}
func gameplayReportPoison(player, unit *server.Object, value byte) int {
	b := gameplayReportID(218, unit, 4)
	b[3] = value
	return gameplayReportSend(gameplayReportRecipient(player), b, false, 1)
}
func gameplayReportExperience(unit *server.Object) int {
	if uint32(unit.ObjClass)&4 != 0 {
		b := []byte{110, 0, 0, 0, 0}
		value := *(*float32)(unsafe.Add(unit.CObj(), 28))
		binary.LittleEndian.PutUint32(b[1:], uint32(int64(value)))
		gameplayReportSend(gameplayReportRecipient(unit), b, false, 1)
	}
	return 0
}
func gameplayReportAnimation(to int, unit *server.Object) int {
	b := gameplayReportID(107, unit, 7)
	binary.LittleEndian.PutUint32(b[3:], *gameplayReportWord(unit.CObj(), 132))
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportXStatus(to int, unit *server.Object) int {
	b := gameplayReportID(101, unit, 7)
	binary.LittleEndian.PutUint32(b[3:], *gameplayReportWord(unit.CObj(), 20))
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportPlayerStatus(unit *server.Object) int {
	b := []byte{102, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[1:], uint32(unit.ObjFlags))
	return gameplayReportSend(gameplayReportRecipient(unit), b, false, 1)
}
func gameplayReportCharges(to int, unit *server.Object, current, max byte) int {
	b := gameplayReportID(100, unit, 5)
	b[3], b[4] = current, max
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportDequipItem(to int, unit *server.Object) int {
	return gameplayReportSend(to, gameplayReportID(97, unit, 3), true, 0)
}
func gameplayReportTotalHealth(to int, unit *server.Object) int {
	h := unit.HealthData
	if h == nil {
		return 0
	}
	b := gameplayReportID(221, unit, 7)
	binary.LittleEndian.PutUint16(b[3:], h.Cur)
	binary.LittleEndian.PutUint16(b[5:], h.Max)
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportCurrentHealth(to int, unit *server.Object) int {
	h := unit.HealthData
	if h == nil {
		return 0
	}
	b := gameplayReportID(65, unit, 4)
	b[3] = byte(h.Cur >> 1)
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportTeam(to int, unit *server.Object) int {
	h := unit.HealthData
	if h == nil {
		return 0
	}
	b := []byte{196, 12, 0, 0, byte(uint32(h.Cur) * 100 / uint32(h.Max))}
	binary.LittleEndian.PutUint16(b[2:], gameplayReportCode(unit))
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportHealthDelta(to int, id uint16, delta int16) int {
	if delta >= 0 {
		return int(delta)
	}
	b := []byte{66, 0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[1:], id)
	binary.LittleEndian.PutUint16(b[3:], uint16(delta))
	return int(int16(gameplayReportSend(to, b, false, 1)))
}
func gameplayReportItemHealth(to int, unit *server.Object) int {
	h := unit.HealthData
	if h == nil || h.Max == 0 {
		return gameplayReportPtr(unsafe.Pointer(h))
	}
	b := gameplayReportID(68, unit, 7)
	binary.LittleEndian.PutUint16(b[3:], h.Cur)
	binary.LittleEndian.PutUint16(b[5:], h.Max)
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportStamina(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	return gameplayReportSend(to, []byte{71, *controlByte(unit.UpdateData, 91)}, true, 1)
}
func gameplayReportObjectByte(to int, unit *server.Object) int {
	return gameplayReportSend(to, []byte{91, *controlByte(unit.CObj(), 440)}, false, 1)
}
func gameplayReportPlayerStat(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	b := []byte{74, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[1:], *gameplayReportWord(controlPlayer(unit), 2164))
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportTotalMana(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 || *controlByte(controlPlayer(unit), 2251) == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	b := gameplayReportID(222, unit, 7)
	binary.LittleEndian.PutUint16(b[3:], *controlHalf(unit.UpdateData, 4))
	binary.LittleEndian.PutUint16(b[5:], *controlHalf(unit.UpdateData, 8))
	return gameplayReportCoalesce(to, b)
}
func gameplayReportMana(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 || *controlByte(controlPlayer(unit), 2251) == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	b := gameplayReportID(69, unit, 5)
	binary.LittleEndian.PutUint16(b[3:], *controlHalf(unit.UpdateData, 4))
	return gameplayReportCoalesce(to, b)
}
func gameplayReportStats(to int, unit *server.Object, value byte) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return int(int8(uint32(unit.ObjClass)))
	}
	b := gameplayReportID(72, unit, 14)
	pl := controlPlayer(unit)
	binary.LittleEndian.PutUint16(b[3:], unit.HealthData.Max)
	binary.LittleEndian.PutUint16(b[5:], *controlHalf(unit.UpdateData, 8))
	binary.LittleEndian.PutUint16(b[7:], *controlHalf(unit.CObj(), 490))
	binary.LittleEndian.PutUint16(b[9:], *controlHalf(pl, 2235))
	binary.LittleEndian.PutUint16(b[11:], *controlHalf(pl, 2239))
	b[13] = value
	return int(int8(gameplayReportSend(to, b, false, 1)))
}
func gameplayReportArmor(to int, value uint32) int {
	b := []byte{73, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[1:], value)
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportDrop(to int, unit *server.Object) int {
	b := gameplayReportID(77, unit, 5)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportTimer(to int, value uint32) int {
	b := make([]byte, 13)
	b[0] = 211
	binary.LittleEndian.PutUint32(b[1:], value)
	binary.LittleEndian.PutUint32(b[5:], Sub_40A230())
	binary.LittleEndian.PutUint32(b[9:], GetServer().S().Frame())
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportEnchant(to int, unit *server.Object) int {
	b := gameplayReportID(90, unit, 7)
	binary.LittleEndian.PutUint32(b[3:], *gameplayReportWord(unit.CObj(), 340))
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportHidden(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&0x800000 != 0 {
		op := byte(56)
		if uint32(unit.ObjFlags)&0x1000000 != 0 {
			op = 55
		}
		gameplayReportSend(to, gameplayReportID(op, unit, 3), false, 1)
	}
	return 0
}
func gameplayReportHeight(to int, unit *server.Object) int {
	p := unit.CObj()
	height := *(*float32)(unsafe.Add(p, 104))
	if *gameplayReportWord(p, 20)&0x20 != 0 {
		b := []byte{159, 0, 0, byte(int64(height)), byte(int64(*(*float32)(unsafe.Add(p, 108)))), byte(int64(*(*float32)(unsafe.Add(p, 116))))}
		binary.LittleEndian.PutUint16(b[1:], uint16(unit.NetCode))
		return gameplayReportDirect(to, b)
	}
	op := byte(94)
	if height < 0 {
		op = 95
		height = -height
	}
	b := gameplayReportID(op, unit, 4)
	b[3] = byte(int64(height))
	return gameplayReportDirect(to, b)
}
func gameplayReportEarthquakeByte(to int, value byte) int {
	return gameplayReportDirect(to, []byte{151, value})
}
func gameplayReportEarthquake(pos *types.Pointf, amplitude int) int {
	players := &GetServer().S().Players
	for unit := players.FirstUnit(); unit != nil; unit = players.NextUnit(unit) {
		pl := controlPlayer(unit)
		dx := float64(pos.X) - float64(*(*float32)(unsafe.Add(pl, 3632)))
		dy := float64(pos.Y) - float64(*(*float32)(unsafe.Add(pl, 3636)))
		distance := dx*dx + dy*dy
		if distance < 90000 {
			gameplayReportEarthquakeByte(gameplayReportRecipient(unit), byte(int64((1-distance*0.000011111111)*float64(amplitude))))
		}
	}
	return 0
}
func gameplayReportAcquireCreature(to int, unit *server.Object) int {
	b := gameplayReportID(108, unit, 5)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	if controlFlags(0x8000000) {
		b[4] |= 0x80
	}
	gameplayReportSend(to, b, true, 1)
	return gameplayReportTotalHealth(to, unit)
}
func gameplayReportShield(to int, unit *server.Object) int {
	code := uint16(unit.NetCode)
	if controlFlags(0x8000000) {
		code |= 0x8000
	}
	b := []byte{109, byte(code), byte(code >> 8)}
	return gameplayReportSend(to, b, true, 1)
}
func gameplayReportMonitor(to int, unit *server.Object) int {
	b := gameplayReportID(219, unit, 5)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	gameplayReportSend(to, b, true, 1)
	return gameplayReportTotalHealth(to, unit)
}
func gameplayReportUnmonitor(to int, unit *server.Object) int {
	return gameplayReportSend(to, gameplayReportID(220, unit, 3), true, 1)
}
func gameplayReportSpeed(to int, unit *server.Object, flag byte, value uint32) int {
	b := gameplayReportID(104, unit, 8)
	binary.LittleEndian.PutUint32(b[3:], value)
	b[7] = flag
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportJournal(to int, record unsafe.Pointer, kind byte) int {
	b := make([]byte, 68)
	b[0], b[1] = 213, kind
	copy(b[2:66], alloc.GoString((*byte)(record)))
	if kind != 2 {
		binary.LittleEndian.PutUint16(b[66:], *controlHalf(record, 72))
	}
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportChapter(to int, chapter byte, done int) int {
	flag := byte(0)
	if done == 1 {
		flag = 1
	}
	return gameplayReportSend(to, []byte{214, chapter, flag}, false, 1)
}
func gameplayReportFlag(to int, a, b, c byte, value uint16) int {
	return gameplayReportSend(to, []byte{216, b, a, c, byte(value), byte(value >> 8)}, true, 1)
}
func gameplayReportBall(to int, flag byte, value uint16) int {
	return gameplayReportSend(to, []byte{217, flag, byte(value), byte(value >> 8)}, true, 1)
}
func gameplayReportSpellStat(to int, value uint32, flag byte) int {
	b := []byte{223, 0, 0, 0, 0, flag}
	binary.LittleEndian.PutUint32(b[1:], value)
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportSecondary(to int, unit *server.Object, flag byte) int {
	b := gameplayReportID(224, unit, 4)
	b[3] = flag
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportQuiver(to int, unit *server.Object) int {
	return gameplayReportSend(to, gameplayReportID(225, unit, 3), true, 0)
}
func gameplayReportInventoryLoaded(to int) int { return gameplayReportSend(to, []byte{113}, true, 0) }
func gameplayReportFriend(to int, unit *server.Object, value int) int {
	op := byte(53)
	if value == 1 {
		op = 52
	}
	return gameplayReportSend(to, gameplayReportID(op, unit, 3), true, 1)
}
func gameplayReportFriendReset(to int) int { return gameplayReportSend(to, []byte{54}, true, 1) }
func gameplayReportFade(to, fade, direction int) int {
	a, b := byte(0), byte(0)
	if fade != 0 {
		a = 1
	}
	if direction != 0 {
		b = 1
	}
	return gameplayReportSend(to, []byte{228, a, b}, true, 1)
}
func gameplayReportQuestStart(to int) int { return gameplayReportSend(to, []byte{240, 0}, true, 1) }
func gameplayReportQuestObject(to int, unit *server.Object) int {
	return gameplayReportSend(to, []byte{240, 1, byte(unit.NetCode), byte(unit.NetCode >> 8)}, true, 1)
}
func gameplayReportQuestLevel(to int, unit *server.Object) int {
	return gameplayReportSend(to, []byte{240, 4, *controlByte(unit.UpdateData, 320), byte(unit.NetCode), byte(unit.NetCode >> 8)}, false, 1)
}
func gameplayReportSilverKey(to int, unit *server.Object, value byte) int {
	return gameplayReportSend(to, []byte{240, 22, value, byte(unit.NetCode), byte(unit.NetCode >> 8)}, false, 1)
}
func gameplayReportGoldKey(to int, unit *server.Object, value byte) int {
	return gameplayReportSend(to, []byte{240, 23, value, byte(unit.NetCode), byte(unit.NetCode >> 8)}, false, 1)
}
func gameplayReportGauntlet(to int) int { return gameplayReportSend(to, []byte{240, 20}, false, 1) }
