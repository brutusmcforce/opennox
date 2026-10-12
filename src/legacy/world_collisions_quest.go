package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func worldQuestPending() *byte { return memmap.PtrUint8(0x5D4594, 1567844) }
func worldQuestCountdown() int32 {
	seconds := floatToInt32(float32(GetServer().S().Balance.Float("QuestExitTimerStart")))
	if memmap.Uint32(0x587000, 4660) != 0 {
		seconds = int32(memmap.Uint32(0x5D4594, 3468)-uint32(PlatformTicks())) / 1000
	}
	total, ready := 0, 0
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if *equipmentWord(controlPlayer(u), 4792) == 1 {
			total++
			if *controlPtr(u.UpdateData, 312) != nil {
				ready++
			}
		}
	}
	if total == 0 {
		return int32(serverConfigTimerSet(int32(0)))
	}
	reduction := floatToInt32(float32(float64(ready) / float64(total) * float64(seconds)))
	next := seconds - reduction
	reduced := next < seconds
	if reduced {
		seconds = next
	}
	if reduced || !GetServer().GetFlag3592() {
		GetServer().ServStartCountdown(int(seconds), "objcoll.c:ExitCountdown")
		return int32(gameplayReportGauntlet(255))
	}
	return 1
}
func worldQuestIgnoreHost(u *server.Object) bool {
	return noxflags.HasGame(noxflags.GameHost) && noxflags.HasEngine(noxflags.EngineNoRendering) && *controlByte(controlPlayer(u), 2064) == 31
}
func worldQuestMaybeWarp() bool {
	stage := uint32(Nox_game_getQuestStage_4E3CC0())
	threshold := uint32(Nox_server_questNextStageThreshold_4D74F0(int(stage)))
	count := 0
	allow := Nox_server_questAllowDefault()
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if worldQuestIgnoreHost(u) {
			continue
		}
		pl := controlPlayer(u)
		if *equipmentWord(pl, 4792) == 0 {
			continue
		}
		count++
		if *controlPtr(u.UpdateData, 316) == nil {
			return false
		}
		if *equipmentWord(pl, 4696) >= threshold {
			allow = true
		}
	}
	return count != 0 && allow
}
func worldQuestExitReady() bool {
	count := 0
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if worldQuestIgnoreHost(u) || *equipmentWord(controlPlayer(u), 4792) == 0 {
			continue
		}
		count++
		if *controlPtr(u.UpdateData, 312) == nil {
			return false
		}
	}
	return count != 0
}
func worldCollideExit(a, b *server.Object) {
	core := GetServer().S()
	if dword_5d4594_1567960 == 0 {
		dword_5d4594_1567960 = uint32(core.Types.IndByID("Glyph"))
	}
	if b == nil || b.ObjClass&4 == 0 {
		return
	}
	d := b.UpdateData
	mapName := (*byte)(a.CollideData)
	quest := noxflags.HasGame(noxflags.GameModeQuest)
	if quest && (*controlPtr(d, 312) != nil || *controlPtr(d, 316) != nil) {
		return
	}
	if a.ObjSubClass&2 != 0 && questRuntimeWord(1556120) == 0 {
		return
	}
	if Sub_4DCC90() == 1 || Sub_4DCC10(b) == 0 || noxflags.HasGame(noxflags.GamePause) {
		return
	}
	if *mapName == 0 && !quest {
		return
	}
	pl := b.UpdateDataPlayer().Player
	if *controlByte(pl.C(), 2251) == 1 {
		for it := b.Field129; it != nil; it = it.Field128 {
			if uint32(it.TypeInd) == uint32(dword_5d4594_1567960) && it.InvHolder == nil {
				GetServer().DelayedDelete(it)
				if count := controlByte(d, 244); *count != 0 {
					*count -= 1
				}
			}
		}
	}
	Sub_4DCBF0(1)
	if noxflags.HasGame(noxflags.GameModeCoop) {
		Nox_setSaveFileName_4DB130("WORKING")
		Sub_4DB170(true, a.CObj(), 0)
		return
	}
	ready := true
	if quest {
		mapName = mapQuestChoose()
		for ability := 1; ability < 6; ability++ {
			if core.Abils.IsActive(b, server.Ability(ability)) {
				Sub_4FC300(b, int(int32(ability)))
			}
		}
		if a.ObjSubClass&1 != 0 {
			stage := uint32(Nox_game_getQuestStage_4E3CC0()) + 1
			questRuntimeStageComplete(b)
			pl = (*server.Player)(*controlPtr(d, 276))
			if highest := equipmentWord(pl.C(), 4696); *highest < stage {
				*highest = stage
				questRuntimeHighestMessage(int(uint8(pl.PlayerInd)), uint16(stage))
			}
			*controlPtr(d, 312) = a.CObj()
			*controlPtr(d, 316) = nil
			Nox_xxx_playerSetState_4FA020(b, 13)
			Nox_xxx_playerGoObserver_4E6860((*server.Player)(*controlPtr(d, 276)), 0, 0)
			gameplayTextInformationAll(18, unsafe.Pointer(&b.NetCode))
			ready = worldQuestExitReady()
		} else if a.ObjSubClass&2 != 0 {
			*controlPtr(d, 312) = nil
			*controlPtr(d, 316) = a.CObj()
			questRuntimeSetWord(1556108, core.Frame())
			Nox_xxx_playerSetState_4FA020(b, 13)
			Nox_xxx_playerGoObserver_4E6860((*server.Player)(*controlPtr(d, 276)), 0, 0)
			for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
				if u != b {
					gameplayTextInformation(int(*controlByte(controlPlayer(u), 2064)), 19, unsafe.Pointer(&b.NetCode))
				}
			}
			worldCollideMessage(b, "objcoll.c:PlayerEntersWarp")
			ready = worldQuestMaybeWarp()
		}
		worldCollideSound(1003, a)
		if a.ObjSubClass&2 == 0 && !ready {
			worldQuestCountdown()
		}
		worldCopyNarrow(unsafe.Pointer(worldQuestPending()), unsafe.Pointer(mapName))
		if !ready {
			return
		}
	}
	if mapName != nil && *mapName != 0 {
		if a.ObjSubClass&2 != 0 {
			stage := uint32(Nox_game_getQuestStage_4E3CC0())
			threshold := uint32(Nox_server_questNextStageThreshold_4D74F0(int(stage)))
			Nox_game_setQuestStage_4E3CD0(int(threshold - 1))
			questRuntimeSetWord(1556124, 1)
			questRuntimeResetAll()
		}
		GetServer().SwitchMap(alloc.GoString(mapName))
	}
}
func worldCollideSoulGate(a, b *server.Object) {
	d := a.CollideData
	core := GetServer().S()
	if !noxflags.HasGame(noxflags.GameModeQuest) || b == nil || b.ObjClass&4 == 0 {
		return
	}
	questRuntimeGateSet(0)
	first := true
	for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
		if *equipmentWord(controlPlayer(u), 4792) == 1 && *controlPtr(u.UpdateData, 308) != nil {
			first = false
		}
	}
	if first {
		questRuntimeSetSoulFrame(core.Frame())
	}
	if *controlPtr(b.UpdateData, 308) != a.CObj() || core.Frame()-*equipmentWord(d, 0) > uint32(core.TickRate()) {
		worldCollideSound(1005, a)
		visibilityFXPoint(130, a.PosVec)
		worldCollideMessage(b, "objcoll.c:SoulGateCollide")
	}
	*controlPtr(b.UpdateData, 308) = a.CObj()
	*equipmentWord(d, 0) = core.Frame()
}
