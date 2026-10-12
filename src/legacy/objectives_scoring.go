package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func objectiveTeamCount(t *server.Team) int { return teamRuntimeCount(t) }
func objectiveScore(u *server.Object) {
	gameplayReportChangeScore(u, 1)
	gameplayReportLesson(u)
}
func objectiveQuestScore(u *server.Object) {
	if dword_5d4594_2650652 != 0 && u != nil && u.UpdateData != nil {
		statisticsEvent(unsafe.Pointer(u.UpdateDataPlayer().Player), nil)
	}
}
func objectiveBallCollide(u, t *server.Object) {
	core := GetServer().S()
	ud := u.UpdateData
	var team *server.Team
	if t != nil {
		team = core.Teams.ByID(t.TeamVal.ID)
		if *temporaryRefWord(ud, 0) == t && team != nil && objectiveTeamCount(team) > 1 {
			last := memmap.PtrUint32(0x5d4594, 1568008)
			if core.Frame()-*last > 45 {
				gameplayTextPrivate(t, (*byte)(unsafe.Pointer(internCStr("objcoll.c:CantPickupBall"))), 0)
				*last = core.Frame()
			}
			inventorySound(928, u, 0, 0)
			return
		}
	}
	if u.ObjOwner != nil || t == nil || t.ObjClass&4 == 0 {
		inventorySound(928, u, 0, 0)
		return
	}
	for it := t.Field129; it != nil; it = it.Field128 {
		if it.TypeInd == u.TypeInd {
			return
		}
	}
	core.ObjSetOwner(t, u)
	if u.TeamVal.Has() {
		if team != nil {
			teamRuntimeSwitch(u.TeamPtr(), team, int(u.NetCode), 0)
		}
	} else {
		teamRuntimeJoin(t.TeamVal.ID, u.TeamPtr(), 1, int(u.NetCode), 0)
	}
	if team != nil {
		state := byte(2)
		if team.ColorInd == 2 {
			state = 4
		}
		Sub_4E8290(state, uint16(t.NetCode))
	}
	objectiveRememberOwner(u, t)
	inventorySound(927, u, 0, 0)
	u.ObjFlags |= 0x40
	objectivePickupBuffs(t)
}
func objectiveFlagCollide(u, t *server.Object) {
	if t == nil || t.ObjFlags&0x8000 != 0 {
		return
	}
	if noxflags.HasGame(32) {
		if t.ObjClass&4 != 0 {
			objectiveCTFPickup(u, t)
		}
		return
	}
	if noxflags.HasGame(64) {
		cache := memmap.PtrUint32(0x5d4594, 1567996)
		if *cache == 0 {
			*cache = uint32(GetServer().S().Types.IndByID("GameBall"))
		}
		if uint32(t.TypeInd) == *cache || t.ObjClass&4 != 0 {
			objectiveFlagBallScore(u, t)
		}
	}
}
func objectiveCTFPickup(u, t *server.Object) {
	core := GetServer().S()
	playerUD := t.UpdateData
	color := byte(objectiveFlagID(u))
	team := u.TeamVal.ID
	if t.TeamVal.SameAs(&u.TeamVal) {
		ud := u.UpdateData
		epsilon := *memmap.PtrFloat64(0x581450, 10160)
		if math.Abs(float64(u.PosVec.X)-float64(*temporaryFloat(ud, 0))) > epsilon || math.Abs(float64(u.PosVec.Y)-float64(*temporaryFloat(ud, 4))) > epsilon {
			Nox_xxx_unitMove_4E7010(u, *(*types.Pointf)(ud))
			netcode := t.NetCode
			gameplayTextInformationAll(4, unsafe.Pointer(&netcode))
			*equipmentWord(ud, 8) = 0
			matchRosterFlagState(byte(t.TeamVal.ID), 0, byte(color), 0)
			return
		}
		for flag := t.InvFirstItem; flag != nil; flag = flag.InvNextItem {
			if flag.ObjClass&0x10000000 == 0 {
				continue
			}
			limit := uint16(serverConfigScore(int16(32)))
			data := flag.UpdateData
			flagColor := objectiveFlagID(flag)
			flagTeam := flag.TeamVal.ID
			objectiveScore(t)
			if t.TeamVal.Has() {
				tm := core.Teams.ByID(t.TeamVal.ID)
				teamRuntimeLessons(tm, tm.Lessons+1)
				if dword_5d4594_2650652 != 0 && playerUD != nil {
					statisticsEvent(unsafe.Pointer((*server.PlayerUpdateData)(playerUD).Player), nil)
				}
			}
			inventoryRemove(t, flag)
			GetServer().CreateObjectAt(flag, nil, *(*types.Pointf)(data))
			Nox_xxx_unitRaise_4E46F0(flag, 0)
			playerStateMark(flag, 1)
			*equipmentWord(data, 8) = 0
			matchRosterFlagState(byte(flagTeam), 0, byte(flagColor), 0)
			inventoryMessage(5, t, uint32(flagColor))
			if limit > 0 {
				for tm := core.Teams.First(); tm != nil; tm = core.Teams.Next(tm) {
					if uint32(tm.Lessons) >= uint32(limit) {
						noxflags.SetGame(8)
						gameplayReportFlagWinner(tm, 0)
						break
					}
				}
				return
			}
		}
		return
	}
	ud := u.UpdateData
	if u.InvHolder != nil {
		return
	}
	tm := core.Teams.ByID(u.TeamVal.ID)
	if tm == nil || objectiveTeamCount(tm) == 0 {
		gameplayTextPrivate(t, (*byte)(unsafe.Pointer(internCStr("objcoll.c:FlagNoTeam"))), 0)
		return
	}
	for flag := t.InvFirstItem; flag != nil; flag = flag.InvNextItem {
		if flag.ObjClass&0x10000000 != 0 {
			inventoryForceDrop(t, flag)
			break
		}
	}
	nox_xxx_servFinalizeDelObject_4DADE0(asObjectC(u))
	flagColor := objectiveFlagID(u)
	inventoryInsert(t, u, 1)
	*equipmentWord(unsafe.Pointer((*server.PlayerUpdateData)(playerUD).Player), 4) |= 1
	gameplayReportEquipment(255, u)
	inventoryMessage(6, t, uint32(flagColor))
	playerStateUnmark(u, 1)
	*equipmentWord(ud, 8) = 0
	matchRosterFlagState(byte(team), 1, byte(color), uint16(t.NetCode))
	objectivePickupBuffs(t)
}
func objectivePointFX(id byte, u *server.Object) int16 {
	return int16(visibilityFXPoint(byte(id), u.PosVec))
}
func objectiveResetMotion(u *server.Object) {
	for _, off := range []int{80, 84, 88, 100} {
		*equipmentWord(u.CObj(), off) = 0
	}
}
func objectiveHomeBase(base, ball *server.Object) int16 {
	core := GetServer().S()
	ballID := core.Types.IndByID("GameBall")
	startID := core.Types.IndByID("GameBallStart")
	if ball == nil {
		return int16(startID)
	}
	if int(ball.TypeInd) != ballID {
		return int16(ball.TypeInd)
	}
	ud := ball.UpdateData
	var baseTeam, ownerTeam *server.Team
	if base.TeamVal.Has() {
		baseTeam = core.Teams.ByID(base.TeamVal.ID)
	}
	owner := *temporaryRefWord(ud, 0)
	if owner != nil && owner.TeamVal.Has() {
		ownerTeam = core.Teams.ByID(owner.TeamVal.ID)
	}
	if baseTeam == ownerTeam {
		objectiveScore(owner)
	}
	if baseTeam != nil {
		teamRuntimeLessons(baseTeam, baseTeam.Lessons+1)
		objectiveQuestScore(*temporaryRefWord(ud, 0))
		inventorySound(929, base, 0, 0)
		objectivePointFX(154, ball)
	}
	n := 0
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if int(it.TypeInd) == startID {
			n++
		}
	}
	choice := int(int32(nox_common_randomInt_415FA0(0, int(int32(n-1)))))
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if int(it.TypeInd) == startID {
			if choice == 0 {
				core.ObjClearOwner(ball)
				Nox_xxx_unitMove_4E7010(ball, it.PosVec)
				out := objectivePointFX(129, ball)
				objectiveResetMotion(ball)
				return out
			}
			choice--
		}
	}
	return 0
}
func objectiveFlagBallScore(flag, target *server.Object) int16 {
	core := GetServer().S()
	cache := memmap.PtrUint32(0x5d4594, 1567992)
	if *cache == 0 {
		*cache = uint32(core.Types.IndByID("GameBall"))
	}
	out := int16(*cache)
	ball := target
	if target.ObjClass&4 != 0 {
		out = int16(bool2int(stateOwns(target, 1567720, "GameBall")))
		if out == 0 {
			return out
		}
		ball = target.Field129
		if ball == nil {
			return out
		}
		out = int16(*cache)
		for uint32(ball.TypeInd) != *cache {
			ball = ball.Field128
			if ball == nil {
				return out
			}
		}
	}
	if ball == nil {
		return out
	}
	ud := ball.UpdateData
	if owner := *temporaryRefWord(ud, 0); owner != nil && owner.ObjFlags&0x20 != 0 {
		*temporaryRefWord(ud, 0) = nil
	}
	team := core.Teams.Next(core.Teams.ByID(flag.TeamVal.ID))
	if team == nil {
		team = core.Teams.First()
	}
	out = int16(memmap.PtrWord(team.C()))
	owner := *temporaryRefWord(ud, 0)
	if owner == nil {
		return out
	}
	out = int16(owner.TeamVal.ID)
	if owner.TeamVal.ID != team.ID() {
		return out
	}
	limit := uint16(serverConfigScore(int16(64)))
	objectiveScore(owner)
	teamRuntimeLessons(team, team.Lessons+1)
	objectiveQuestScore(owner)
	inventorySound(929, flag, 0, 0)
	inventoryMessage(9, owner, uint32(team.ID()))
	objectivePointFX(154, ball)
	if limit > 0 {
		for tm := core.Teams.First(); tm != nil; tm = core.Teams.Next(tm) {
			if uint32(tm.Lessons) >= uint32(limit) {
				noxflags.SetGame(8)
				gameplayReportFlagballWinner(tm)
				break
			}
		}
	}
	if dword_5d4594_1567988 == 0 {
		dword_5d4594_1567988 = uint32(core.Types.IndByID("GameBallStart"))
	}
	n := 0
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if uint32(it.TypeInd) == uint32(dword_5d4594_1567988) {
			n++
		}
	}
	choice := int(int32(nox_common_randomInt_415FA0(0, int(int32(n-1)))))
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if uint32(it.TypeInd) == uint32(dword_5d4594_1567988) {
			if choice == 0 {
				data := ball.UpdateData
				core.ObjClearOwner(ball)
				objectiveRememberOwner(ball, nil)
				teamRuntimeLeave(ball.TeamPtr(), int(ball.NetCode))
				resourceRestoreHP(ball)
				*(*uint64)(unsafe.Add(data, 8)) = uint64(uint32(PlatformTicks()))
				Nox_xxx_unitMove_4E7010(ball, it.PosVec)
				Sub_4E8290(0, 0)
				out = objectivePointFX(129, ball)
				objectiveResetMotion(ball)
				return out
			}
			choice--
		}
	}
	return 0
}
