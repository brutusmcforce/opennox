package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func playerDeathTeam(u *server.Object) *server.Team {
	if u == nil || !u.TeamVal.Has() {
		return nil
	}
	return GetServer().S().Teams.ByID(u.TeamVal.ID)
}
func playerDeathStats(actor, target *server.Object) {
	if Get_dword_5d4594_2650652() == 0 || actor == nil || target == nil || actor.UpdateData == nil || target.UpdateData == nil {
		return
	}
	statisticsEvent(*(*unsafe.Pointer)(unsafe.Add(actor.UpdateData, 276)), *(*unsafe.Pointer)(unsafe.Add(target.UpdateData, 276)))
}
func playerDeathTeamDelta(team *server.Team, delta int) {
	if team != nil {
		teamRuntimeLessons(team, team.Lessons+delta)
	}
}
func playerDeathElimination(victim, killer *server.Object) {
	vt, kt := playerDeathTeam(victim), playerDeathTeam(killer)
	if killer == victim {
		gameplayReportSubtractLessons(victim, 1)
		gameplayReportEliminationDeath(victim)
		gameplayReportLesson(victim)
		playerDeathTeamDelta(vt, 1)
		playerDeathStats(killer, killer)
		return
	}
	if killer != nil && killer.ObjClass&4 != 0 {
		if kt != nil && kt == vt {
			gameplayReportSubtractLessons(killer, 1)
			gameplayReportLesson(killer)
			playerDeathStats(killer, killer)
		} else {
			gameplayReportChangeScore(killer, 1)
			gameplayReportLesson(killer)
			playerDeathStats(killer, victim)
		}
	} else if killer == nil {
		playerDeathStats(victim, victim)
	}
	gameplayReportEliminationDeath(victim)
	gameplayReportLesson(victim)
	playerDeathTeamDelta(vt, 1)
}
func playerDeathArena(victim, killer, assist *server.Object, info *server.Player) {
	vt, kt := playerDeathTeam(victim), playerDeathTeam(killer)
	var at *server.Team
	if info != nil && assist != nil {
		at = playerDeathTeam(assist)
	}
	if killer == victim || (killer == nil && assist == nil) {
		gameplayReportSubtractLessons(victim, 1)
		gameplayReportLesson(victim)
		playerDeathTeamDelta(vt, -1)
		playerDeathStats(killer, killer)
	} else {
		if killer != nil && killer.ObjClass&4 != 0 {
			if kt != nil && kt == vt {
				gameplayReportSubtractLessons(killer, 1)
				gameplayReportLesson(killer)
				playerDeathTeamDelta(kt, -1)
				playerDeathStats(killer, killer)
			} else {
				gameplayReportChangeScore(killer, 1)
				gameplayReportLesson(killer)
				playerDeathTeamDelta(kt, 1)
				playerDeathStats(killer, victim)
			}
		}
		gameplayReportEliminationDeath(victim)
		gameplayReportLesson(victim)
	}
	if assist == nil || (kt != nil && kt == at) || (at != nil && vt != nil && vt == at) {
		return
	}
	gameplayReportChangeScore(assist, 1)
	gameplayReportLesson(assist)
	playerDeathTeamDelta(at, 1)
	if info != nil {
		playerDeathStats(assist, victim)
	}
}
func playerDeathKotr(victim, killer *server.Object) {
	vt, kt := playerDeathTeam(victim), playerDeathTeam(killer)
	if killer == nil || killer.ObjClass&4 == 0 {
		return
	}
	finish := func() { gameplayReportEliminationDeath(victim); gameplayReportLesson(victim) }
	owns := func(u *server.Object) bool { return stateOwns(u, 1567716, "Crown") }
	if killer == victim || (vt != nil && vt == kt) {
		if owns(killer) {
			gameplayReportSubtractLessons(killer, 1)
			gameplayReportLesson(killer)
			playerDeathTeamDelta(kt, -1)
			playerDeathStats(killer, killer)
		}
		finish()
		return
	}
	kc, vc := owns(killer), owns(victim)
	if kt == nil {
		if kc || vc {
			key := "KotRPawnKillsKingPoints"
			if kc {
				key = "KotRKingKillsPawnPoints"
			}
			points := floatToInt32(float32(GetServer().S().Balance.Float(key)))
			gameplayReportChangeScore(killer, int(points))
			gameplayReportLesson(killer)
			playerDeathStats(killer, victim)
			if !noxflags.HasGamePlay(4) && vc {
				itemDropCrowns(victim, memmap.PtrWord(killer.CObj()))
			}
		}
	} else if kc || vc {
		key := "KotRPawnKillsKingPoints"
		if kc {
			key = "KotRKingKillsPawnPoints"
			if vc {
				key = "KotRKingKillsKingPoints"
			}
		}
		points := floatToInt32(float32(GetServer().S().Balance.Float(key)))
		gameplayReportChangeScore(killer, int(points))
		playerDeathTeamDelta(kt, int(points))
		gameplayReportLesson(killer)
		playerDeathStats(killer, victim)
	}
	finish()
}
