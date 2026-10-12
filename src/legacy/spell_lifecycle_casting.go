package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func spellLifeForbiddenObjective(id int32) bool {
	return GetServer().S().Spells.HasFlags(spell.ID(id), 0x80000)
}
func spellLifeHasTeam(u *server.Object) bool {
	return u.TeamVal.Has()
}
func spellLifeOwnsType(u *server.Object, id uint32) bool {
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) == id {
			return true
		}
	}
	return false
}
func spellLifeCantCast(u *server.Object, id, queued int32) int32 {
	if queued == 0 {
		if controlFlags(16) {
			typ := stateType(1569704, "Crown")
			if spellLifeForbiddenObjective(id) && spellLifeOwnsType(u, typ) && spellLifeHasTeam(u) {
				return 17
			}
		} else if controlFlags(64) {
			typ := stateType(1569708, "GameBall")
			if spellLifeForbiddenObjective(id) && spellLifeOwnsType(u, typ) {
				return 16
			}
		} else if controlFlags(32) && spellLifeForbiddenObjective(id) {
			for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
				if it.ObjClass&0x10000000 != 0 {
					return 13
				}
			}
		}
	}
	if spellLifeHasBuff(u, 29) {
		return 14
	}
	if !questRuntimeBookAllowed(int(id), 111) {
		return 10
	}
	count := func(off uintptr) int32 {
		return int32(Nox_xxx_unitIsUnitTT_4E7C80(u, int(int32(*memmap.PtrUint32(0x5d4594, off)))))
	}
	switch id {
	case 29:
		if count(1569692) > 0 || count(1569688) > 0 || count(1569684) > 0 {
			return 3
		}
	case 31:
		if count(1569696) > 0 {
			return 3
		}
	case 52:
		if count(1569700) > 0 {
			return 3
		}
	case 50, 58:
		off, name := uintptr(1569680), "MagicMissileCount"
		if id == 58 {
			off, name = 1569676, "PixieCount"
		}
		n := count(off)
		power := spellLifePower(id, u) - 1
		limit := int32(int64(float64(nox_xxx_gamedataGetFloatTable_419D70(internCStr(name), int(int32(power))))))
		if n >= limit {
			return 3
		}
	}
	return 0
}
func spellLifeCaptureAllowed(id int32, u *server.Object) int32 {
	ball := *memmap.PtrUint32(0x5d4594, 1569712)
	if ball == 0 {
		ball = stateType(1569712, "GameBall")
		*memmap.PtrUint32(0x5d4594, 1569716) = uint32(GetServer().S().Types.IndByID("Crown"))
	}
	if u == nil {
		return 0
	}
	if spellLifeForbiddenObjective(id) && u.ObjClass&4 != 0 {
		if controlFlags(32) {
			if *controlByte(controlPlayer(u), 4)&1 != 0 {
				return 0
			}
		} else if controlFlags(64) {
			if spellLifeOwnsType(u, ball) {
				return 0
			}
		} else if controlFlags(16) {
			if spellLifeOwnsType(u, *memmap.PtrUint32(0x5d4594, 1569716)) && spellLifeHasTeam(u) {
				return 0
			}
		}
	}
	return 1
}
func spellLifePhoneme(code int32, phon int8) int32 {
	player := false
	if controlFlags(1) {
		player = Nox_server_getObjectFromNetCode_4ECCB0(int(code)).ObjClass&4 != 0
	} else {
		player = GetClient().Cli().Objs.ByNetCodeDynamic(int(code)).Class()&4 != 0
	}
	female := false
	if player {
		p := GetServer().S().Players.ByID(int(code))
		if p == nil {
			return 0
		}
		female = *controlByte(unsafe.Pointer(p), int(layoutServerPlayer.Off(2252)) /* info */) != 0
	}
	var sound int32
	switch phon {
	case 0:
		sound = 193
	case 1:
		sound = 186
	case 2:
		sound = 187
	case 3:
		sound = 192
	case 5:
		sound = 188
	case 6:
		sound = 191
	case 7:
		sound = 190
	case 8:
		sound = 189
	default:
		return 0
	}
	if female {
		sound += 8
	}
	return sound
}
func spellLifeBroadcastPhoneme(u *server.Object, phon int8) int32 {
	s := GetServer().S()
	for p := s.Players.FirstUnit(); p != nil; p = s.Players.NextUnit(p) {
		if p != u {
			aud := spellLifePhoneme(int32(u.NetCode), phon)
			netCode := uint32(p.NetCode)
			GetServer().S().Audio.EventObj(sound.ID(aud), u, 2, netCode)
		}
	}
	return 0
}
func spellLifeCreateFly(u, target *server.Object, id int32) *server.Object {
	s := GetServer().S()
	radius := float32(float64(*(*float32)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(176)) /* Shape.Circle.R */))) + 4)
	if target == nil {
		var pos *types.Pointf
		if u.ObjClass&4 != 0 {
			pl := controlPlayer(u)
			pos = &types.Pointf{X: float32(int32(*spellLifeWord(pl, 2284))), Y: float32(int32(*spellLifeWord(pl, 2288)))}
		}
		target = s.Nox_xxx_spellFlySearchTarget(pos, u, s.Spells.DefByInd(spell.ID(id)).Def.Flags, 600, 0, u)
	}
	dir := int16(u.Direction1)
	dx := *memmap.PtrFloat32(0x587000, uintptr(194136+8*int(dir)))
	dy := *memmap.PtrFloat32(0x587000, uintptr(194140+8*int(dir)))
	// The x87 source keeps this intermediate at 53-bit precision through
	// the inherited velocity addition; storing float32 early changes one ULP.
	y := float64(radius)*float64(dy) + float64(u.PosVec.Y)
	pos := types.Pointf{X: float32(float64(radius)*float64(dx) + float64(u.PosVec.X) + float64(u.VelVec.X)), Y: float32(y + float64(u.VelVec.Y))}
	if !s.MapTraceRayAt(u.PosVec, pos, nil, nil, 5) {
		return nil
	}
	out := s.NewObjectByTypeID("Magic")
	if out == nil {
		return nil
	}
	*spellLifeWord(out.UpdateData, 16) = uint32(spellLifePower(id, u))
	GetServer().CreateObjectAt(out, u, pos)
	out.Direction1 = u.Direction1
	out.Direction2 = u.Direction1
	*controlPtr(out.UpdateData, 0) = u.CObj()
	*controlPtr(out.UpdateData, 4) = target.CObj()
	*controlPtr(out.UpdateData, 8) = u.CObj()
	*spellLifeWord(out.UpdateData, 12) = uint32(id)
	// Preserve the direction helper's side effects before velocity stores.
	var indexed [2]int32
	geometryIndexedDirection(int32(u.Direction1), &indexed)
	out.VelVec.X = float32(float64(dx) * float64(*(*float32)(unsafe.Add(out.CObj(), int(layoutServerObject.Off(544)) /* SpeedCur */))))
	out.VelVec.Y = float32(float64(dy) * float64(*(*float32)(unsafe.Add(out.CObj(), int(layoutServerObject.Off(544)) /* SpeedCur */))))
	out.VelVec.X = float32(float64(out.VelVec.X) + float64(u.VelVec.X))
	out.VelVec.Y = float32(float64(out.VelVec.Y) + float64(u.VelVec.Y))
	if spellLifeHasBuff(u, 21) {
		spellLifeApplyBuff(out, 21, int16(spellLifeBuffTimer(u, 21)), spellLifeBuffPower(u, 21))
	}
	s.Audio.EventObj(s.Spells.DefByInd(spell.ID(id)).GetAudio(0), u, 0, 0)
	return out
}
func spellLifeCollide(u, target *server.Object) {
	if spellLifeHasBuff(u, 22) && target.ObjFlags&0x8008 == 0 && target.ObjClass&6 != 0 && GetServer().S().IsEnemyTo(u, target) {
		power := int32(u.BuffsPower[22]) - 1
		GetServer().S().Audio.EventObj(135, u, 0, 0)
		spellLifeBuffOff(u, 22)
		damage := float32(nox_xxx_gamedataGetFloatTable_419D70(internCStr("ShockDamage"), int(int32(power))))
		target.CallDamage(u, u, int(floatToInt32(float32(damage))), 9)
	}
	if target.ObjClass&0x20006 != 0 && target.ObjFlags&0x8020 == 0 && !itemOwnerSameTeam(target, u) {
		spellLifeBuffOff(u, 0)
	}
	if u.ObjClass&4 != 0 && target.ObjClass&0x20000 != 0 && target.ObjFlags&0x8020 == 0 {
		spellLifeBuffOff(u, 0)
	}
}
