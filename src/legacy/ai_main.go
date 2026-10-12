package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func monsterIdleAudio(u *server.Object) {
	ud, core := u.UpdateDataMonster(), GetServer().S()
	if !u.Class().Has(object.ClassMonster) {
		return
	}
	if !monsterAttackAtWill(u) && ud.CurrentEnemy == nil && float64(math.Float32frombits(ud.Field131)) <= 300 {
		a := monsterHeadSafe(u)
		if (a == 0 || a == 4) && core.Frame() >= ud.Field132 {
			ud.Field132 = core.Frame() + uint32(nox_common_randomInt_415FA0(int(20*core.TickRate()), int(60*core.TickRate())))
			combatSound(u, 4)
		}
	}
}
func monsterDangerous(t, u *server.Object) int16 {
	cloud := memmap.PtrUint32(0x5D4594, 2489468)
	small := memmap.PtrUint32(0x5D4594, 2489472)
	if *cloud == 0 {
		*cloud = uint32(GetServer().S().Types.IndByID("ToxicCloud"))
		*small = uint32(GetServer().S().Types.IndByID("SmallToxicCloud"))
	}
	var value uint32
	if t.Class()&0x2000 != 0 {
		value = (uint32(u.SubClass()) >> 10) & 1
		if value == 0 {
			dword_5d4594_2489460 = 0
		}
	} else {
		value = uint32(t.TypeInd)
		if value == *cloud || value == *small {
			value = uint32(u.SubClass())
			if value&0x200 == 0 {
				dword_5d4594_2489460 = 0
			}
		} else if t.Class()&0x10000 != 0 {
			dword_5d4594_2489460 = 0
		}
	}
	return int16(value)
}
func monsterUnwindAttacks(u *server.Object) {
	for {
		switch u.UpdateDataMonster().AIStackHead().Type() {
		case 16, 17, 18, 19, 20, 25, 26, 27, 28:
			u.MonsterPopAction()
		default:
			return
		}
	}
}
func monsterDodge(u *server.Object) bool {
	core := GetServer().S()
	vx, vy := movementDirectionVector(int32(int16(u.Direction1)))
	sideX, sideY := -vy, vx
	origin := u.PosVec
	for attempt := 0; attempt < 5; attempt++ {
		distance := core.Rand.Logic.FloatClamp(2, 3) * float64(u.Shape.Circle.R)
		mag := float32(distance)
		if distance > 15 {
			mag = 15
		}
		if nox_common_randomInt_415FA0(0, 100) < 50 {
			mag = -mag
		}
		end := types.Pointf{X: float32(float64(mag)*float64(sideX) + float64(u.PosVec.X)), Y: float32(float64(mag)*float64(sideY) + float64(u.PosVec.Y))}
		if !core.MapTraceRayAt(origin, end, nil, nil, 1) || !core.MapTraceObstacles(u, origin, end) || tileAtPoint(end) == 6 {
			continue
		}
		monsterUnwindAttacks(u)
		u.MonsterPushAction(ai.ActionType(41), core.Frame()+core.TickRate())
		u.MonsterPushAction(ai.ACTION_DODGE, end, 0)
		return true
	}
	return false
}
func monsterShieldThreat(u *server.Object) *server.Object {
	monsterShieldThreatObj = nil
	*memmap.PtrUint32(0x5D4594, 2487988) = 1315859240
	GetServer().S().Map.EachMissileInCircle(u.PosVec, 100, func(t *server.Object) bool { monsterShieldCandidate(t, u); return true })
	return monsterShieldThreatObj
}
func monsterShieldCandidate(t, u *server.Object) {
	if !spatialTeamEligible(u, t) {
		return
	}
	point := t.PrevPos
	if stateFront(&u.PosVec, int32(int16(u.Direction1)), &point)&1 == 0 {
		return
	}
	vx, vy := movementDirectionVector(int32(int16(u.Direction1)))
	if !(float64(t.VelVec.X)*float64(vx)+float64(t.VelVec.Y)*float64(vy) < 0) {
		return
	}
	// x87 keeps both deltas and distance full precision through the angular test.
	dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
	cross := float64(vx)*float64(dy) - float64(vy)*float64(dx)
	if !(math.Abs(cross) < 20) {
		return
	}
	speed := math.Sqrt(float64(t.VelVec.X)*float64(t.VelVec.X) + float64(t.VelVec.Y)*float64(t.VelVec.Y))
	// Normalized X is the one early float32 spill in this calculation.
	nx := float32(float64(t.VelVec.X) / speed)
	distance := math.Sqrt(dy*dy + dx*dx)
	dot := float64(t.VelVec.Y)/speed*-(float64(dy)/float64(distance)) + float64(nx)*-(float64(dx)/float64(distance))
	if dot > .69999999 && GetServer().S().CanInteract(u, t, 0) && float64(float32(distance)) < float64(*memmap.PtrFloat32(0x5D4594, 2487988)) {
		monsterShieldThreatObj = t
		*memmap.PtrFloat32(0x5D4594, 2487988) = float32(distance)
	}
}
func monsterMainAI(u *server.Object) {
	ud, core := u.UpdateDataMonster(), GetServer().S()
	head := ud.AIStackHead() // keep this captured pointer across pushes, as C does.
	if (head.Type() == 0 || head.Type() == 4) && (uint32(byte(core.Frame()))+uint32(byte(u.NetCode))-uint32(byte(ud.Field137)))&15 != 0 {
		return
	}
	for i := int(ud.AIStackInd) - 1; i >= 0; i-- {
		if ud.AIStack[i].Action == 61 {
			return
		}
	}
	if u.ObjFlags&0x8000 != 0 {
		return
	}
	set := Nox_xxx_monsterGetSoundSet_424300(u)
	play := func(slot int) {
		if set != nil {
			id := *(*uint32)(unsafe.Add(set, slot*4))
			core.Audio.EventObj(sound.ID(id), u, 0, 0)
		}
	}
	if noxflags.HasGame(2048) && nox_xxx_guiCursor_477600() == 0 && u.Field5&16 != 0 {
		if host := core.Players.HostUnit(); host != nil && host.ObjFlags&2 == 0 && !ud.HasAction(2) {
			p := host.UpdateDataPlayer()
			mouse := p.Player.CursorVec
			dx, dy := float64(mouse.X)-float64(u.PosVec.X), float64(mouse.Y)-float64(u.PosVec.Y)
			if dy*dy+dx*dx < 100 && Nox_xxx_findObjectAtCursor_54AF40(host) == u {
				u.PrevPos = u.PosVec
				u.VelVec = types.Pointf{}
				u.ForceVec = types.Pointf{}
				u.Pos24 = types.Pointf{}
				u.MonsterPushAction(ai.ActionType(71))
				u.MonsterPushAction(ai.ActionType(2), core.TickRate())
				u.MonsterPushAction(ai.ActionType(67))
				u.MonsterPushAction(ai.ActionType(2), 999999)
				u.MonsterPushAction(ai.ActionType(26), host)
				if p.DialogWith == nil && p.Trade70 == nil {
					play(1)
				}
				return
			}
		}
	}
	if u.HasEnchant(3) && !ud.HasAction(36) {
		head = u.MonsterPushAction(ai.ActionType(62), 3)
		u.MonsterPushAction(ai.ActionType(36))
	}
	if !u.HasEnchant(29) && monsterCastInversion(u) {
		return
	}
	if u.HasEnchant(11) && monsterMoving(u) && !ud.HasAction(24) {
		u.MonsterPushAction(ai.ActionType(62), 11)
		u.MonsterPushAction(ai.ActionType(24), u.PosVec, 0)
		play(12)
		return
	}
	if byte(core.Frame())&15 == 0 && (monsterAggressionMid(u) || monsterAttackAtWill(u)) && ud.HasAction(4) && monsterHeadSafe(u) != 4 && !ud.HasAction(15) {
		t := core.EnemyAggroYyy(u, 100)
		if t != nil && t.Class()&6 != 0 {
			ud.StatusFlags |= 0x200
			u.Obj130 = ud.CurrentEnemy
			u.Field131 = 11
			u.Frame134 = core.Frame()
		}
	}
	if !monsterAggressionRetreat(u) && monsterMoving(u) && !monsterCastBusy(u) && !u.HasEnchant(3) && !monsterMoveAttempt(u) {
		if enemy := ud.CurrentEnemy; enemy != nil {
			dist := float64(stateDistance(u, enemy))
			if dist < float64(ud.FleeRange) {
				if ud.StatusFlags&0x20 != 0 && ud.Field376 != 0 && !u.HasEnchant(29) && core.Frame() >= ud.Field371 && float64(ud.FleeRange)*.5 > float64(float32(dist)) {
					// The retained spell engine passes this buffer back through Go.
					// Keep it C-owned across that callback boundary.
					args, freeArgs := alloc.New([3]uint32{})
					*args = [3]uint32{memmap.PtrWord(u.CObj()), math.Float32bits(u.PosVec.X), math.Float32bits(u.PosVec.Y)}
					monsterCastSpell(4, u, args)
					freeArgs()
					ud.Field371 = core.Frame() + uint32(nox_common_randomInt_415FA0(int(ud.Field370_0), int(ud.Field370_2)))
					return
				}
				if !ud.HasAction(24) && ud.FleeRange != 0 {
					u.MonsterPushAction(ai.ActionType(28), int(int16(u.Direction1))+128)
					u.MonsterPushAction(ai.ActionType(68))
					u.MonsterPushAction(ai.ActionType(63), float32(float64(ud.FleeRange)+30))
					u.MonsterPushAction(ai.ActionType(24), ud.CurrentEnemy.PosVec, 0)
					if nox_common_randomInt_415FA0(0, 1) != 0 {
						play(12)
					}
					return
				}
			}
		}
	}
	if u.HealthData.Max != 0 && monsterMoving(u) && !monsterMoveAttempt(u) && !ud.HasAction(24) && !ud.HasAction(6) && !ud.HasAction(14) {
		fraction := float32(float64(u.HealthData.Cur) / float64(u.HealthData.Max))
		muted := ud.StatusFlags&0x20 != 0 && u.HasEnchant(29)
		if float64(fraction) <= float64(ud.RetreatLevel) || muted {
			monsterUnwindAttacks(u)
			u.MonsterPushAction(ai.ActionType(68))
			if (ud.StatusFlags&0x80 != 0 || u.SubClass()&0x80 != 0) && noxflags.HasGame(2048) {
				u.MonsterPushAction(ai.ActionType(14))
			} else {
				u.MonsterPushAction(ai.ActionType(6))
			}
			play(13)
			GetServer().NoxScriptC().ScriptCallback(&ud.ScriptRetreat, nil, u, server.ScriptEventType(4))
			return
		}
	}
	if u.SubClass()&0x10 != 0 && ud.Field514&0x400 != 0 && head.Type() != 16 && head.Type() != 17 && monsterShieldThreat(u) != nil {
		if head.Type() != 1 && head.Type() != 23 {
			monsterUnwindAttacks(u)
			u.MonsterPushAction(ai.ActionType(1), core.Frame()+core.TickRate())
		}
		return
	}
	if monsterHasShield(u) && head.Type() != 16 && head.Type() != 17 && monsterShieldThreat(u) != nil {
		if !ud.HasAction(21) {
			monsterUnwindAttacks(u)
			u.MonsterPushAction(ai.ActionType(21), core.Frame()+(core.TickRate()>>1))
		}
		return
	}
	if noxflags.HasGame(2048) && !monsterAggressionRetreat(u) && !u.HasEnchant(3) && ud.MonsterDef.StatusFlags92&8 != 0 && !ud.HasAction(9) && monsterShieldThreat(u) != nil && monsterDodge(u) {
		return
	}
	switch head.Type() {
	case 7, 37, 8, 10, 24:
		dx, dy := float64(math.Float32frombits(ud.Field125))-float64(u.PosVec.X), float64(math.Float32frombits(ud.Field126))-float64(u.PosVec.Y)
		if dy*dy+dx*dx > 225 {
			ud.Field124 = core.Frame()
			ud.Field125 = math.Float32bits(u.PosVec.X)
			ud.Field126 = math.Float32bits(u.PosVec.Y)
		} else if uint32(core.Frame()-ud.Field124) > uint32(int32(core.TickRate())>>1) {
			unitDebug(1, GetServer().S().Frame(), u.NetCode, GetServer().S().Types.ByInd(int(u.TypeInd)).ID())
			ud.StatusFlags |= 0x200000
			if ud.HasAction(6) || ud.HasAction(14) || ud.HasAction(24) {
				ud.Field127 = core.Frame()
			}
			if ud.HasAction(15) {
				monsterDodge(u)
			} else if nox_common_randomInt_415FA0(0, 100) >= 33 || !monsterDodge(u) {
				if st := u.MonsterPushAction(ai.ActionType(1)); st != nil {
					st.SetArgs(core.Frame() + uint32(nox_common_randomInt_415FA0(int(core.TickRate()>>1), int(2*core.TickRate()))))
				}
			}
			ud.Field124 = core.Frame()
			ud.Field125 = math.Float32bits(u.PosVec.X)
			ud.Field126 = math.Float32bits(u.PosVec.Y)
			return
		}
	}
	if !monsterAggressionRetreat(u) && monsterHasMissingHealth(u) && byte(core.Frame())&15 == 0 {
		if t := lifecycleFoodSearch(u, 75, false); t != nil {
			Nox_xxx_inventoryServPlace_4F36F0(u, t, 1, 1)
			if t.SubClass()&0x90 != 0 {
				effectsUse(u, t)
			}
		}
	}
	if ud.StatusFlags&0x20000 != 0 {
		data := (*server.PlayerUpdateData)(memmap.WordPtr(ud.Field545))
		p := data.Player
		if *(*byte)(unsafe.Add(unsafe.Pointer(p), int(layoutServerPlayer.Off(2251)) /* info */)) == 0 && *(*uint32)(unsafe.Add(unsafe.Pointer(p), 4)) == 0 && byte(core.Frame())&15 == 0 {
			if t := lifecycleFoodSearch(u, 75, true); t != nil {
				controlMorphToPlayer(u)
				Nox_xxx_inventoryServPlace_4F36F0(u, t, 1, 1)
				controlMorphFromPlayer(u)
			}
		}
	}
}
