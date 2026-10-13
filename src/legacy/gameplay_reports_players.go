package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

func gameplayReportInterestingID(unit *server.Object) int {
	b := gameplayReportID(210, unit, 7)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	b[5], b[6] = 2, 2
	players := &GetServer().S().Players
	for it := players.FirstUnit(); it != nil; it = players.NextUnit(it) {
		gameplayReportSend(gameplayReportRecipient(it), b, false, 1)
	}
	return 0
}
func gameplayReportReset(unit *server.Object) int {
	result := gameplayReportPtr(unit.CObj())
	if uint32(unit.ObjClass)&4 != 0 {
		ud := unit.UpdateData
		layout.As[server.PlayerUpdateData](ud).Field62 = 0
		layout.As[server.PlayerUpdateData](ud).Field63 = 0
		layout.As[server.PlayerUpdateData](ud).Field64 = GetServer().S().Frame()
		if layout.As[server.PlayerUpdateData](ud).IsCamping != 0 {
			result = gameplayReportInterestingID(unit)
		}
		layout.As[server.PlayerUpdateData](ud).IsCamping = 0
	}
	return result
}
func gameplayReportResetAll() int {
	players := &GetServer().S().Players
	for unit := players.FirstUnit(); unit != nil; unit = players.NextUnit(unit) {
		gameplayReportReset(unit)
	}
	return 0
}
func gameplayReportPlayerHealthToTeam(to int) int {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(to))
	if pl == nil {
		return 0
	}
	unit := pl.PlayerUnit
	if unit == nil || unit.HealthData == nil {
		return gameplayReportPtr(unit.CObj())
	}
	hp := unit.HealthData.Cur
	gameplayReportSend(to, []byte{67, byte(hp), byte(hp >> 8)}, true, 1)
	if controlFlags(4096) {
		return gameplayReportTeam(to|0x80, unit)
	}
	return 0
}
func gameplayReportWinner(op byte, id uint16, value byte) int {
	b := []byte{op, byte(id), byte(id >> 8), value, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[4:], GetServer().S().Frame())
	return gameplayReportSend(255, b, true, 1)
}
func gameplayReportDMWinner(unit *server.Object, value byte) int {
	if unit != nil && uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	return gameplayReportWinner(88, gameplayReportCode(unit), value)
}
func gameplayReportDMTeamWinner(team *server.Team, value byte) int {
	id := uint16(0)
	if team != nil {
		id = uint16(team.IDVal)
	}
	return gameplayReportWinner(89, id, value)
}
func gameplayReportFlagballWinner(team *server.Team) int {
	return gameplayReportWinner(86, uint16(team.IDVal), 0)
}
func gameplayReportFlagWinner(team *server.Team, value byte) int {
	id := uint16(65535)
	if team != nil {
		id = uint16(team.IDVal)
	}
	return gameplayReportWinner(87, id, value)
}
func gameplayReportScavenger(unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	b := gameplayReportID(85, unit, 7)
	pl := controlPlayer(unit)
	binary.LittleEndian.PutUint16(b[3:], *(*uint16)(unsafe.Pointer(&pl.Field2152)))
	binary.LittleEndian.PutUint16(b[5:], *(*uint16)(unsafe.Pointer(&pl.Field2156)))
	return gameplayReportSend(255, b, false, 1)
}
func gameplayReportEliminationDeath(unit *server.Object) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return 0
	}
	controlPlayer(unit).Field2140++
	if !controlFlags(1024) {
		return 0
	}
	core := GetServer().S()
	elapsed := core.Frame()-*memmap.PtrUint32(0x5D4594, 3520) > 20*uint32(core.TickRate())
	if elapsed && *memmap.PtrUint32(0x5D4594, 3536) == 0 {
		for pl := core.Players.First(); pl != nil; pl = core.Players.Next(pl) {
			if *(*byte)(unsafe.Pointer(&pl.Field3680))&1 != 0 {
				Nox_xxx_netNeedTimestampStatus_4174F0(pl, 256)
			}
		}
		*memmap.PtrUint32(0x5D4594, 3536) = 1
	}
	if controlFlags(0x4000000) || GetServer().GetFlag3592() || !elapsed {
		return 0
	}
	limit := int(int32(*memmap.PtrUint32(0x5D4594, 3476)))
	offset := uintptr(198928)
	if !noxflags.HasGamePlay(noxflags.GameplayFlag4) {
		if Sub_40A770() >= limit {
			return 0
		}
	} else {
		if int(uint8(core.Teams.Count())) >= limit {
			return 0
		}
		found := false
		for tm := core.Teams.First(); tm != nil; tm = core.Teams.Next(tm) {
			if Nox_xxx_countNonEliminatedPlayersInTeam_40A830(tm) == 1 {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
		offset = 198872
	}
	seconds := floatToInt32(float32(core.Balance.Float("SuddenDeathCountdown")))
	GetServer().ServStartCountdown(int(seconds), strman.ID(alloc.GoString((*byte)(memmap.PtrOff(0x587000, offset)))))
	return 0
}
func gameplayReportChangeScore(unit *server.Object, value int) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	pl := controlPlayer(unit)
	*(*uint32)(unsafe.Pointer(&pl.Lessons)) += uint32(value)
	return gameplayReportPtr(unsafe.Pointer(pl))
}
func gameplayReportSubtractLessons(unit *server.Object, value int) int {
	if uint32(unit.ObjClass)&4 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	pl := controlPlayer(unit)
	*(*uint32)(unsafe.Pointer(&pl.Lessons)) -= uint32(value)
	return gameplayReportPtr(unsafe.Pointer(pl))
}
func gameplayReportLesson(unit *server.Object) int {
	b := []byte{78, byte(unit.NetCode), byte(unit.NetCode >> 8), 0, 0, 0, 0, 0, 0, 0, 0}
	pl := controlPlayer(unit)
	binary.LittleEndian.PutUint32(b[3:], *(*uint32)(unsafe.Pointer(&pl.Lessons)))
	binary.LittleEndian.PutUint32(b[7:], pl.Field2140)
	return gameplayReportSend(255, b, true, 1)
}
func gameplayReportAnything(unit *server.Object) int {
	if unit == nil || uint32(unit.ObjClass)&4 == 0 {
		return 0
	}
	core := GetServer().S()
	ud := unit.UpdateData
	pl := controlPlayer(unit)
	to := gameplayReportRecipient(unit)
	armor := layout.As[server.PlayerUpdateData](ud).Field57
	if *(*float32)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).Field58)) != *(*float32)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).Field57)) {
		gameplayReportArmor(to, armor)
		layout.As[server.PlayerUpdateData](ud).Field58 = armor
	}
	if pl.Field2168 != pl.GoldVal {
		gameplayReportPlayerStat(to, unit)
		pl.Field2168 = pl.GoldVal
	}
	if pl.Field2172 != *(*byte)(unsafe.Pointer(&unit.Field110)) {
		gameplayReportObjectByte(to, unit)
		pl.Field2172 = *(*byte)(unsafe.Pointer(&unit.Field110))
	}
	if controlFlags(4096) {
		for i := 0; i < 32; i++ {
			player := core.Players.ByInd(ntype.PlayerInd(i))
			pud := layout.As[server.PlayerUpdateData](ud)
			if player != nil && player.PlayerUnit != nil && pud.Field80 != uint32(pud.QuestLevels[i]) {
				gameplayReportQuestLevel(i, unit)
				pud.QuestLevels[i] = byte(pud.Field80)
			}
		}
		for _, key := range []struct {
			name   string
			cache  uintptr
			keys   func(*server.PlayerUpdateData) *[32]byte
			report func(int, *server.Object, byte) int
		}{{"SilverKey", 1556324, func(d *server.PlayerUpdateData) *[32]byte { return &d.SilverKeys }, gameplayReportSilverKey}, {"GoldKey", 1556328, func(d *server.PlayerUpdateData) *[32]byte { return &d.GoldKeys }, gameplayReportGoldKey}} {
			for i := 0; i < 32; i++ {
				cache := memmap.PtrUint32(0x5D4594, key.cache)
				if *cache == 0 {
					*cache = uint32(core.Types.IndByID(key.name))
				}
				if core.Players.ByInd(ntype.PlayerInd(i)) == nil {
					continue
				}
				value := byte(0)
				for it := unit.InvFirstItem; it != nil; it = it.InvNextItem {
					if uint32(it.TypeInd) == *cache {
						value = 1
						break
					}
				}
				if keys := key.keys(layout.As[server.PlayerUpdateData](ud)); value != keys[i] {
					key.report(i, unit, value)
					keys[i] = value
				}
			}
		}
	}
	if pl.Field2184 != 0 {
		gameplayReportTotalHealth(to, unit)
		gameplayReportTotalMana(to, unit)
		gameplayReportStats(to, unit, pl.Level)
		pl.Field2184 = 0
	}
	if hp := unit.HealthData; hp != nil && hp.Cur != layout.As[server.PlayerUpdateData](ud).Field2_1 {
		gameplayReportPlayerHealthToTeam(to)
		layout.As[server.PlayerUpdateData](ud).Field2_1 = hp.Cur
	}
	if layout.As[server.PlayerUpdateData](ud).ManaCur != layout.As[server.PlayerUpdateData](ud).ManaPrev {
		gameplayReportMana(to, unit)
		layout.As[server.PlayerUpdateData](ud).ManaPrev = layout.As[server.PlayerUpdateData](ud).ManaCur
	}
	return 0
}
