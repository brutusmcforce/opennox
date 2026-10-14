package legacy

import (
	"bytes"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func sessionResetPlayers() {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if u := pl.PlayerUnit; u != nil {
			dword_5d4594_2649712 &^= uint32(uint32(1) << uint32(pl.PlayerInd&31))
			pl.Field3676 = 2
			controlDefaultItems(u, 1, 0)
			pl.Field2140 = 0
			pl.Lessons = 0
		}
	}
	s.MapSend.Sub_51A100()
	noxflags.UnsetGame(0x20000)
	matchRosterSettings()
	serverConfigUpdatedClear()
}
func sessionSaveAllowed() int32 {
	if !noxflags.HasGame(2048) {
		return 1
	}
	pl := GetServer().S().Players.ByInd(31)
	if pl == nil || pl.PlayerUnit == nil || Sub_4DCC90() != 0 || Get_dword_5d4594_251744() != 0 || GetServer().S().Frame()-memmap.Uint32(0x5D4594, 1563068) < 30 || memmap.Uint32(0x5D4594, 1096672) != 0 || pl.PlayerUnit.ObjFlags&0x4000 != 0 {
		return 0
	}
	return int32(bool2int(Sub_4DCC10(pl.PlayerUnit) != 0))
}
func sessionPlayerReady(index int32) {
	s := GetServer().S()
	pl := s.Players.ByInd(ntype.PlayerInd(index))
	if pl == nil {
		return
	}
	playerStateAddStatus(pl, 16)
	if pl.PlayerUnit != nil {
		// The client's friend list holds net codes, which a new map reuses for other objects: without a reset
		// enemies could get a friendly selection circle after a map change, and friends a hostile one (opennox#535).
		gameplayReportFriendReset(int(index))
		statePlayerVisibility(index)
	}
	if noxflags.HasGame(8192) && !noxflags.HasGame(128) {
		if u := pl.PlayerUnit; u != nil {
			spellLifeBuffOff(u, 23)
			spellLifeApplyBuff(u, 23, int16(5*uint32(uint16(s.TickRate()))), 5)
		}
		for u := s.Objs.First(); u != nil; u = u.Next() {
			if u.ObjClass&0x10000000 != 0 {
				s.Players.Nox_xxx_netMarkMinimapObject_417190(ntype.PlayerInd(index), u, 1)
			}
		}
	}
	if pl.Field3680&1 != 0 {
		pl.Pos3632Vec = pl.PlayerUnit.PosVec
		if noxflags.HasGame(512) {
			controlLeaveObserver((*server.Player)(pl.C()))
		}
	}
	GetServer().Nox_xxx_wall_4DF1E0(int(index))
	if noxflags.HasGame(4096) && GetServer().GetFlag3592() {
		gameplayReportGauntlet(int(index))
	}
	if index == 31 && noxflags.HasGame(128) {
		if Get_nox_server_sanctuaryHelp_54276() == 1 {
			nox_xxx_cliShowHelpGui_49C560()
			Nox_xxx_netStatsMultiplier_4D9C20(pl.PlayerUnit)
			return
		}
		serverOptionsConstruct()
	}
	Nox_xxx_netStatsMultiplier_4D9C20(pl.PlayerUnit)
}
func sessionPlayerIncoming(index int32) uint32 {
	s := GetServer().S()
	pl := s.Players.ByInd(ntype.PlayerInd(index))
	if pl == nil {
		panic("session arrival without player record")
	}
	if noxflags.HasGame(4096) {
		if index != 31 && pl.PlayerUnit != nil {
			pl.PlayerUnit.UpdateDataPlayer().Field138 = 1
		}
		gameplayReportQuestStart(int(index))
		if pl.PlayerUnit != nil {
			questRuntimeReset(pl.PlayerUnit)
		}
	}
	resetNetworkAliases(&pl.NetData16)
	u := pl.PlayerUnit
	dword_5d4594_2649712 |= uint32(uint32(1) << uint32(index&31))
	pos := u.PosVec
	matchRosterSendPlayers(int(index))
	pl.Field4700 = 0
	u.CallInitWithArg(nil)
	pl.Field3676 = 3
	if !noxflags.HasGame(512) {
		pl.Pos3632Vec = pos
	}
	if Get_nox_server_sendMotd_108752() != 0 && noxflags.HasGame(8192) && !noxflags.HasGame(4096) {
		Nox_xxx_playerSendMOTD_4DD140(ntype.PlayerInd(index))
	}
	for other := s.Players.First(); other != nil; other = s.Players.Next(other) {
		if other.PlayerUnit == nil || other == pl {
			continue
		}
		s.Players.Nox_xxx_netMarkMinimapObject_417190(ntype.PlayerInd(index), other.PlayerUnit, 1)
		s.Players.Nox_xxx_netMarkMinimapObject_417190(ntype.PlayerInd(other.PlayerInd), pl.PlayerUnit, 1)
		matchRosterSimpleObject(int(other.PlayerInd), pl.PlayerUnit)
		if noxflags.HasGame(4096) {
			gameplayReportTeam(int(other.PlayerInd), pl.PlayerUnit)
			gameplayReportTeam(int(index), other.PlayerUnit)
		}
	}
	matchRosterMinimap(int(index))
	matchRosterTeamRoster(int(index))
	if noxflags.HasGame(1024) && playerStateAdmission(pl) == 0 {
		playerStateAddStatus(pl, 256)
	}
	statePlayerVisibility(index)
	if noxflags.HasGame(64) {
		p := matchRosterFlagBase()
		gameplayReportBall(int(index), *(*byte)(p), *(*uint16)(unsafe.Add(p, 2)))
	} else if noxflags.HasGame(32) {
		for tm := s.Teams.First(); tm != nil; tm = s.Teams.Next(tm) {
			p := matchRosterFlagRecord(byte(tm.ID()))
			gameplayReportFlag(int(index), *(*byte)(p), *(*byte)(unsafe.Add(p, 2)), *(*byte)(unsafe.Add(p, 1)), *(*uint16)(unsafe.Add(p, 4)))
		}
	}
	playerStateAllStatus(pl)
	if serverConfigFlagsQuery(8192) != 0 {
		matchRosterPlayerIDs(pl)
	}
	if noxflags.HasGame(4096) {
		for unit := s.Players.FirstUnit(); unit != nil; unit = s.Players.NextUnit(unit) {
			if controlPlayer(unit).Field4792 == 1 {
				gameplayReportQuestObject(int(index), unit)
			}
		}
		questRuntimeDepartureStamp(int(index))
	}
	return u.NetCode
}
func sessionPlayerDeparture(index int32) {
	s := GetServer().S()
	pl := s.Players.ByInd(ntype.PlayerInd(index))
	GetServer().Nox_script_event_playerLeave(pl)
	if sessionRosterContains(index) != 0 {
		sessionRosterRemove(index)
	}
	if pl.Field2068 != 0 {
		if group := playerGroupFind(pl.Field2068); group != nil {
			playerGroupRemoveMember(group, index)
		}
	}
	data := pl.PlayerUnit.UpdateDataPlayer()
	if shop := (*shopSession)(unsafe.Pointer(data.Trade70)); shop != nil {
		shopCancel(shop)
	}
	data.Trade70 = nil
	shopPlayerCleanup(int(pl.PlayerInd))
	Sub_4FF990(uint32(1) << uint32(pl.PlayerInd&31))
	if sessionServerBuild || !noxflags.HasGame(2) {
		pl.Active = 0
	}
	for _, word := range []int{1146, 1148, 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159, 1147, 1160, 1161} {
		p := &(*[16]uint32)(unsafe.Pointer(&pl.ProtUnitHPCur))[word-1146]
		if deleteProtectionRecord(*p) {
			*p = 0
		}
	}
	code := uint16(s.GetUnitNetCode(pl.PlayerUnit))
	msg := [3]byte{46, byte(code), byte(code >> 8)}
	reliableEnqueue(int(index)|0x80, msg[:], nil, 0, 0)
	GetServer().DelayedDelete(pl.PlayerUnit)
	pl.PlayerUnit = nil
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		ud := u.UpdateDataPlayer()
		ud.QuestLevels[index] = 0
		ud.SilverKeys[index] = 0
		ud.GoldKeys[index] = 0
		ud.Field81[index] = 0
	}
	if playerStateMultiple() != 0 {
		if noxflags.HasGame(1024) && Nox_xxx_serverIsClosing_446180() == 0 && playerStateCompetitors() == 1 {
			matchRosterWinner(true)
		}
	} else {
		serverConfigTimerSet(0)
		playerStateLessons(1)
		GetServer().TeamsResetYyy()
		playerStateReset()
	}
	reliableRemoveRecipient(byte(index))
	reliableResetRate(int(index))
	ReliableResetSequence(ntype.PlayerInd(index))
	if !noxflags.HasGame(4096) {
		return
	}
	if worldQuestExitReady() {
		GetServer().SwitchMap(alloc.GoString(worldQuestPending()))
	} else if worldQuestMaybeWarp() {
		questRuntimeResetAll()
		questRuntimeSetWord(1556124, 1)
		questRuntimeSetStage(uint32(Nox_server_questNextStageThreshold_4D74F0(int(questRuntimeStage()))) - 1)
		GetServer().SwitchMap(alloc.GoString(worldQuestPending()))
	} else {
		for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
			if u.UpdateDataPlayer().Field78 != nil {
				worldQuestCountdown()
				break
			}
		}
	}
}
func sessionBroadcastSettings() {
	var data [58]byte
	serverOptionsRead(data[:])
	previous := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1563214), 58)
	if bytes.Equal(data[:], previous) {
		return
	}
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if pl.PlayerInd != 31 {
			matchRosterGUISettings(1, unsafe.Pointer(&data[0]), int(pl.PlayerInd))
		}
	}
	copy(previous, data[:])
}
