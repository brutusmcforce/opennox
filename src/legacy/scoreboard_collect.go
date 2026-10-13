package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

func scoreboardObserver(p *server.Player) bool { return p.Field3680&1 != 0 && p.Field3680&0x20 == 0 }

// Copy through the terminator without clearing the remaining shared record.
// Player names can occupy 27 UTF-16 units: the original row writes the team
// field after clipping, including when that field overlaps the name terminator.
func scoreboardCopyName(dst, src *uint16) {
	for i := uintptr(0); ; i += 2 {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(src), i))
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), i)) = v
		if v == 0 {
			return
		}
	}
}
func scoreboardCollectPlayer(p *server.Player, score uint32) {
	r := &scoreboardPlayers()[*scoreboardCount(1090117)]
	r.NetCode = p.NetCodeVal
	r.Score = score
	r.Ping = uint32(*(*uint16)(unsafe.Pointer(&p.Field2148)))
	r.Flags = p.Field3680
	r.Objective = int32(scoreboardObjective(p))
	scoreboardCopyName(&r.Name[0], &p.NameFinal[0])
	scoreboardClipName(&r.Name[0])
	r.Class = p.InfoBytes[66]
	r.Team = -1
	if ot := objectRenderTeam(int(p.NetCodeVal)); ot != nil && ot.Has() {
		if tm := GetServer().S().Teams.ByID(ot.ID); tm != nil {
			r.Team = int32(tm.IDVal)
		}
	}
	*scoreboardCount(1090117)++
}
func scoreboardCollectTeams(low bool) {
	teams := &GetServer().S().Teams
	*scoreboardCount(1090116) = 0
	previous := int32(0x7fffffff)
	if low {
		previous = -0x80000000
	}
	for int(*scoreboardCount(1090116)) < teams.Count() {
		best := int32(-0x80000000)
		if low {
			best = 0x7fffffff
		}
		var selected *server.Team
		for tm := teams.First(); tm != nil; tm = teams.Next(tm) {
			score := int32(tm.Lessons)
			if scoreboardHasTeam(uint32(tm.IDVal)) {
				continue
			}
			if (!low && score <= previous && score > best) || (low && score >= previous && score < best) {
				selected = tm
				best = score
			}
		}
		// Valid game scores exclude the original selection sentinels.
		r := &scoreboardTeams()[*scoreboardCount(1090116)]
		r.Score = int32(selected.Lessons)
		r.ID = uint32(selected.IDVal)
		r.Color = byte(selected.ColorInd)
		scoreboardCopyName(&r.Name[0], (*uint16)(unsafe.Pointer(selected)))
		scoreboardClipName(&r.Name[0])
		*scoreboardCount(1090116)++
		previous = best
	}
}
func scoreboardCollectMode(quest, low bool) {
	players := &GetServer().S().Players
	if !quest {
		scoreboardCollectTeams(low)
	}
	*scoreboardCount(1090117) = 0
	*scoreboardCount(1090118) = 0
	count := players.Count()
	headless := noxflags.HasEngine(noxflags.EngineNoRendering)
	if noxflags.HasGame(noxflags.GameHost) && headless {
		count--
	}
	for p := players.First(); p != nil; p = players.Next(p) {
		if quest {
			if scoreboardObserver(p) {
				p.Field2108 |= 0x80000000
			} else if p.Field2108 == 0 {
				p.Field2108 = 0x08000000
			}
		} else if scoreboardObserver(p) {
			if low {
				p.Field2140 += 65535
			} else {
				p.Lessons -= 65535
			}
		}
	}
	previous := uint32(0x7fffffff)
	if quest {
		previous = 0
	} else if low {
		previous = 0xffffffff
	}
	for int(*scoreboardCount(1090117)) < count {
		best := uint32(0x80000000)
		if quest {
			best = 0xffffffff
		} else if low {
			best = 0x7fffffff
		}
		var selected *server.Player
		for p := players.First(); p != nil; p = players.Next(p) {
			if (headless && p.PlayerInd == 31) || scoreboardHasPlayer(p.NetCodeVal) {
				continue
			}
			score := uint32(p.Lessons)
			if quest {
				score = p.Field2108
			} else if low {
				score = p.Field2140
			}
			take := false
			if quest {
				take = score >= previous && score < best
			} else if low {
				take = int32(score) >= int32(previous) && int32(score) < int32(best)
			} else {
				take = int32(score) <= int32(previous) && int32(score) > int32(best)
			}
			if take {
				selected = p
				best = score
			}
		}
		score := best
		if scoreboardObserver(selected) {
			if quest {
				score += 0x80000000
			} else if low {
				score -= 65535
			} else {
				score += 65535
			}
		} else if quest && score == 0x08000000 {
			score = 0
		} else {
			*scoreboardCount(1090118)++
		}
		scoreboardCollectPlayer(selected, score)
		previous = best
	}
	for p := players.First(); p != nil; p = players.Next(p) {
		if quest {
			if scoreboardObserver(p) {
				p.Field2108 &= 0x7fffffff
			} else if p.Field2108 == 0x08000000 {
				p.Field2108 = 0
			}
		} else if scoreboardObserver(p) {
			if low {
				p.Field2140 -= 65535
			} else {
				p.Lessons += 65535
			}
		}
	}
}
func scoreboardCollect() {
	scoreboardCollectMode(*scoreboardData.mode == 5, noxflags.HasGame(noxflags.GameModeElimination))
}
func scoreboardCollectOrdinary() { scoreboardCollectMode(false, false) }
