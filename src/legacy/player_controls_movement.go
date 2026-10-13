package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func controlFindStart(out *types.Pointf, u *server.Object) {
	s := GetServer().S()
	if dword_5d4594_1568868 == 0 {
		dword_5d4594_1568868 = uint32(s.Types.IndByID("PlayerStart"))
	}
	if u == nil {
		return
	}
	team := int32(0)
	if u.TeamVal.Has() {
		team = int32(*(*byte)(unsafe.Pointer(&u.TeamVal.ID)))
		nox_xxx_getTeamByID_418AB0(int(team))
	}
	var last, best *server.Object
	var choices []*server.Object
	for it := s.Objs.List; it != nil; it = it.ObjNext {
		if uint32(it.TypeInd) == uint32(dword_5d4594_1568868) {
			last = it
			if controlStartEligible(it, team) {
				choices = append(choices, it)
			}
		}
	}
	if len(choices) == 0 {
		if last == nil {
			*out = types.Pointf{X: 2000, Y: 2000}
		} else {
			*out = last.PosVec
		}
		return
	}
	// Pick the start farthest from its nearest enemy. With no enemies, choose
	// randomly among eligible starts, preserving traversal order and RNG use.
	noEnemies := true
	maxDist := float32(0)
	for _, it := range choices {
		minDist := float32(10000000)
		for p := s.Players.FirstUnit(); p != nil; p = s.Players.NextUnit(p) {
			if p != u && GetServer().S().IsEnemyTo(u, p) {
				dx := float64(it.PosVec.X) - float64(p.PosVec.X)
				dy := float64(it.PosVec.Y) - float64(p.PosVec.Y)
				dist := dy*dy + dx*dx
				if dist < float64(minDist) {
					minDist = float32(dist)
				}
				noEnemies = false
			}
		}
		if minDist > maxDist {
			maxDist = minDist
			best = it
		}
	}
	if noEnemies || best == nil {
		best = choices[s.Rand.Logic.IntClamp(0, len(choices)-1)]
	}
	*out = best.PosVec
}
func controlDropBall(u *server.Object) int32 {
	typ := stateType(1568872, "GameBall")
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) != typ {
			continue
		}
		it.ObjFlags &^= 0x40
		GetServer().ApplyForce(it, u.PosVec, float64(float32(100)))
		*(*unsafe.Pointer)(unsafe.Pointer(&it.Obj130)) = nil
		GetServer().S().ObjSetOwner(nil, it)
		GetServer().S().Audio.EventObj(sound.ID(926), u, 0, 0)
		sub_4E8290(int8(1), int16(0))
		return 1
	}
	return 0
}
func controlNearStart(u *server.Object, out *types.Pointf) int32 {
	*out = u.PosVec
	var result int32
	for i := 0; i < 32; i++ {
		inventoryRandomPlacement(60, &u.PosVec, out)
		result = int32(bool2int(worldTileWater(*out)))
		if result == 0 {
			break
		}
	}
	return int32(result)
}
func controlWalkWaypoint(u *server.Object) int32 {
	ptr := &layout.As[server.PlayerUpdateData](u.UpdateData).Waypoints[*controlWaypointByte(u.UpdateData, 1)]
	it := *ptr
	if it == nil {
		return 0
	}
	dx := float32(float64(it.PosVec.X) - float64(u.PosVec.X))
	dy := float64(it.PosVec.Y) - float64(u.PosVec.Y)
	pos := types.Pointf{X: dx, Y: float32(dy)}
	if dy*float64(pos.Y)+float64(dx)*float64(dx) < 100 {
		GetServer().DelayedDelete(it)
		*ptr = nil
		return 0
	}
	u.Direction2 = server.Dir16(geometryVectorAngle(&pos))
	if u.Buffs&(1<<3) != 0 {
		u.Direction2 = server.Dir16(controlConfusedDirection(u))
	}
	off := uintptr(int(int16(u.Direction2)) * 8)
	for i := 0; i < 2; i++ {
		f := (*float32)(unsafe.Add(u.CObj(), 88+4*i))
		*f = float32(float64(*memmap.PtrFloat32(0x587000, 194136+off+uintptr(i*4)))*float64(u.SpeedCur) + float64(*f))
	}
	return 1
}
func controlInputAttack(u *server.Object) {
	if !controlAimsAtEnemy(u) {
		return
	}
	d := u.UpdateData
	weapon := controlPlayer(u).WeaponEquip
	setAttack := func() { u.Field34 = GetServer().S().Frame(); layout.As[server.PlayerUpdateData](d).Field59_0 = 0 }
	if weapon == 0 {
		if layout.As[server.PlayerUpdateData](d).State != 1 {
			Nox_xxx_playerSetState_4FA020(u, 1)
		}
		return
	}
	if weapon&0x47f0000 != 0 && controlActionState(u) != 29 {
		it := layout.As[server.PlayerUpdateData](d).Field26
		data := it.UseDataWand()
		if data.Charges != 0 || data.MaxCharges == 0 {
			setAttack()
			Nox_xxx_playerSetState_4FA020(u, 1)
			effectsUse(u, it)
		} else if controlSubStamina(u, 45) != 0 {
			data.Flags |= 2
			setAttack()
			Nox_xxx_playerSetState_4FA020(u, 1)
		}
	} else if *(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](d).State)) != 1 {
		cost := controlWeaponStamina(weapon)
		if controlSubStamina(u, cost) != 0 {
			setAttack()
			if !Nox_xxx_playerSetState_4FA020(u, 1) {
				controlAdjustStamina(u, -int8(cost))
			}
		}
	}
	spellLifeBuffOff(u, int32(0))
	spellLifeBuffOff(u, int32(23))
	GetServer().S().Spells.Dur.CancelFor(spell.ID(67), u)
}
func controlFollowEnemy(u *server.Object) int32 {
	if u == nil {
		return 0
	}
	it := u.Obj130
	if it == nil {
		return 0
	}
	it = it.FindOwnerChainPlayer()
	if it.ObjFlags&0x20 != 0 || it.ObjClass&2 != 0 || it.ObjClass&4 != 0 && *(*byte)(unsafe.Pointer(&controlPlayer(it).Field3680))&1 != 0 {
		return 0
	}
	Nox_xxx_playerCameraFollow_4E6060(u, it)
	return 1
}
