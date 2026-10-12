package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func objectiveColor(header unsafe.Pointer) int32 {
	if header == nil {
		return 0
	}
	name := alloc.GoString(*(**byte)(header))
	for off := uintptr(205224); ; off += 8 {
		p := *memmap.PtrPtr(0x587000, off)
		if p == nil {
			return 0
		}
		if name == alloc.GoString((*byte)(p)) {
			return int32(*memmap.PtrUint32(0x587000, off+4))
		}
	}
}
func objectiveFlagID(u *server.Object) int32 {
	if u.ObjClass&0x10000000 == 0 {
		return 0
	}
	return objectiveColor(*(*unsafe.Pointer)(unsafe.Add(u.InitData, 4)))
}
func objectiveRememberOwner(u, t *server.Object) *server.Object {
	if t != nil {
		t = (*server.Object)(unsafe.Pointer(nox_xxx_findParentChainPlayer_4EC580(asObjectC(t))))
	}
	ud := u.UpdateData
	if t != nil && t.ObjClass&4 != 0 {
		*temporaryRefWord(ud, 0) = t
		*equipmentWord(ud, 4) = uint32(t.TeamVal.ID)
		*equipmentWord(ud, 16) = GetServer().S().Frame()
	} else {
		*temporaryRefWord(ud, 0) = nil
		*equipmentWord(ud, 4) = 0
	}
	return t
}
func objectivePickupBuffs(u *server.Object) uint32 {
	var out uint32
	for enc := 0; enc < 32; enc++ {
		out = uint32(bool2int(spellLifeHasBuff(u, int32(int8(enc)))))
		if out == 0 {
			continue
		}
		out = uint32(server.EnchantID(enc).Spell())
		if out == 0 {
			continue
		}
		out = uint32(bool2int(bool(nox_xxx_spellHasFlags_424A50(int(out), 0x80000))))
		if out != 0 {
			out = uint32(spellLifeBuffOff(u, int32(enc)))
		}
	}
	return out
}
func objectiveFlagUpdate(u *server.Object) int32 {
	ud := u.UpdateData
	stamp := *equipmentWord(ud, 8)
	if stamp == 0 {
		return 0
	}
	color := objectiveFlagID(u)
	team := u.TeamVal.ID
	core := GetServer().S()
	out := int32(3 * core.TickRate())
	if core.Frame()-*equipmentWord(ud, 8) > uint32(30*core.TickRate()) {
		inventorySound(305, u, 0, 0)
		*equipmentWord(ud, 8) = 0
		matchRosterFlagState(byte(team), 0, byte(color), 0)
		Nox_xxx_unitMove_4E7010(u, *(*types.Pointf)(ud))
		out = int32(gameplayTextInformationAll(8, unsafe.Pointer(&color)))
	}
	return out
}
func objectiveFollow(u, owner *server.Object) {
	radius := float64(*temporaryFloat(owner.CObj(), 176)) + float64(*temporaryFloat(u.CObj(), 176)) + 10
	dx, dy := movementDirectionVector(int32(int16(owner.Direction1)))
	end := types.Pointf{X: float32(radius*float64(dx) + float64(owner.PosVec.X)), Y: float32(radius*float64(dy) + float64(owner.PosVec.Y))}
	if GetServer().S().MapTraceRayAt(owner.PosVec, end, nil, nil, 5) {
		Nox_xxx_unitMove_4E7010(u, end)
	}
}
func objectiveCrownUpdate(u *server.Object) {
	ud := u.UpdateData
	pending := *temporaryRefWord(ud, 4)
	if pending != nil && pending.ObjFlags&0x8020 == 0 {
		inventoryCrownPickup(pending, u, 1)
		return
	}
	if prior := *temporaryRefWord(ud, 0); prior != nil && prior.ObjFlags&0x20 != 0 {
		*temporaryRefWord(ud, 0) = nil
	}
	if owner := u.ObjOwner; owner != nil {
		if owner.ObjFlags&0x8020 != 0 {
			GetServer().S().ObjClearOwner(u)
		} else {
			objectiveFollow(u, owner)
		}
	}
}
func objectiveCrownCollide(u, t *server.Object) uint32 {
	if t != nil && t.ObjFlags&0x8020 == 0 && t.ObjClass&4 != 0 {
		return uint32(inventoryCrownPickup(t, u, 1))
	}
	return uint32(uintptr(t.CObj()))
}
func objectiveClearBallTeam(u *server.Object) {
	objectiveRememberOwner(u, nil)
	teamRuntimeLeave(u.TeamPtr(), int(u.NetCode))
	Sub_4E8290(1, 0)
}
func objectiveBallUpdate(u *server.Object) {
	core := GetServer().S()
	ud := u.UpdateData
	*equipmentWord(u.CObj(), 112) = 1008981770
	if prior := *temporaryRefWord(ud, 0); prior != nil && prior.ObjFlags&0x20 != 0 {
		objectiveClearBallTeam(u)
	}
	// The C tick API narrows to uint32 before subtraction from the uint64 stamp.
	if uint64(uint32(PlatformTicks()))-*(*uint64)(unsafe.Add(ud, 8)) > 20000 {
		objectiveBallReset(u)
		return
	}
	owner := u.ObjOwner
	if owner == nil {
		u.ObjFlags &^= 0x40
		x, y := float64(u.VelVec.X), float64(u.VelVec.Y)
		if float64(*temporaryFloat(ud, 24)) > math.Sqrt(x*x+y*y) {
			*temporaryRefWord(ud, 0) = nil
		}
		return
	}
	if owner != *temporaryRefWord(ud, 0) || core.Frame()-*equipmentWord(ud, 16) <= *equipmentWord(ud, 20) {
		u.ObjFlags |= 0x40
		*(*uint64)(unsafe.Add(ud, 8)) = uint64(uint32(PlatformTicks()))
		owner = u.ObjOwner
		if owner.ObjFlags&0x8020 != 0 {
			core.ObjClearOwner(u)
			objectiveClearBallTeam(u)
		} else {
			objectiveFollow(u, owner)
		}
		return
	}
	u.ObjFlags &^= 0x40
	*equipmentWord(u.CObj(), 520) = 0
	dir := (int32(int16(owner.Direction1)) + int32(nox_common_randomInt_415FA0(-32, 32))) & 255
	dx, dy := movementDirectionVector(dir)
	origin := types.Pointf{X: float32(float64(u.PosVec.X) - float64(dx)*20), Y: float32(float64(u.PosVec.Y) - float64(dy)*20)}
	GetServer().ApplyForce(u, AsPointf(unsafe.Pointer(&origin)), float64(float32(30)))
	core.ObjClearOwner(u)
	Sub_4E8290(1, 0)
	inventorySound(926, u, 0, 0)
}
func objectiveBallReset(old *server.Object) int {
	core := GetServer().S()
	if dword_5d4594_527656 == 0 {
		dword_5d4594_527656 = uint32(core.Types.IndByID("GameBallStart"))
	}
	var n int
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if uint32(it.TypeInd) == uint32(dword_5d4594_527656) {
			n++
		}
	}
	if n == 0 {
		return 0
	}
	choice := int(int32(nox_common_randomInt_415FA0(0, int(int32(n-1)))))
	var start *server.Object
	for it := core.Objs.First(); it != nil; it = it.Next() {
		if uint32(it.TypeInd) == uint32(dword_5d4594_527656) {
			if choice == 0 {
				start = it
				break
			}
			choice--
		}
	}
	if start == nil {
		return 0
	}
	ball := core.NewObjectByTypeID("GameBall")
	if ball == nil {
		return 0
	}
	ud := ball.UpdateData
	*(*uint64)(unsafe.Add(ud, 8)) = uint64(uint32(PlatformTicks()))
	*equipmentWord(ud, 20) = uint32(floatToInt32(float32(core.Balance.Float("FlagballPossDuration"))))
	*temporaryFloat(ud, 24) = float32(core.Balance.Float("FlagballResetVel"))
	playerStateMark(ball, 1)
	GetServer().CreateObjectAt(ball, nil, types.Pointf{})
	core.ObjClearOwner(ball)
	objectiveRememberOwner(ball, nil)
	Sub_4E8290(0, 0)
	Nox_xxx_unitMove_4E7010(ball, start.PosVec)
	for _, off := range []int{80, 84, 88, 100} {
		*equipmentWord(ball.CObj(), off) = 0
	}
	if old != nil {
		for pl := core.Players.FirstUnit(); pl != nil; pl = core.Players.NextUnit(pl) {
			if *temporaryRefWord(unsafe.Pointer(pl.UpdateDataPlayer().Player), 3628) == old {
				nox_xxx_playerCameraUnlock_4E6040(asObjectC(pl))
				nox_xxx_playerCameraFollow_4E6060(asObjectC(pl), asObjectC(ball))
			}
		}
		GetServer().DelayedDelete(old)
	}
	return 1
}
func objectiveObelisk(u *server.Object) int32 {
	core := GetServer().S()
	energy := (*int32)(u.UpdateData)
	idle := true
	for pl := core.Players.FirstUnit(); pl != nil; pl = core.Players.NextUnit(pl) {
		if pl.ObjFlags&0x8000 != 0 {
			continue
		}
		ud := pl.UpdateDataPlayer()
		if u.TeamVal.Has() && !u.TeamVal.SameAs(&pl.TeamVal) {
			continue
		}
		dx := float64(u.PosVec.X) - float64(pl.PosVec.X)
		dy := float64(u.PosVec.Y) - float64(pl.PosVec.Y)
		if dx*dx+dy*dy >= 2500 || nox_xxx_mapCheck_537110(asObjectC(u), asObjectC(pl)) == 0 {
			continue
		}
		idle = false
		spent := false
		if *energy >= 1 {
			wand := *temporaryRefWord(pl.UpdateData, 104)
			rate := effectsRechargeRate(wand)
			if noxflags.HasGame(8192) && !noxflags.HasGame(4096) {
				rate = 1
			}
			if rate != 0 && wand != nil && wand.ObjClass&0x1000 != 0 && wand.ObjSubClass&0x47f0000 != 0 && int32(*equipmentWord(wand.UseData.Ptr, 112)) < 100 {
				if noxflags.HasGame(4096) {
					if core.Frame()%uint32(core.TickRate()>>1) == 0 {
						inventorySound(230, u, 0, 0)
					}
				} else {
					spent = true
					*energy--
				}
				if effectsRecharge(wand, rate) != 0 {
					data := unsafe.Slice((*byte)(wand.UseData.Ptr), 110)
					gameplayReportCharges(int(uint8(ud.Player.PlayerInd)), wand, byte(data[108]), byte(data[109]))
				}
			}
		}
		if *energy >= 1 && ud.ManaCur < ud.ManaMax {
			amount := int16(1)
			if noxflags.HasGame(4096) {
				var value float32
				switch *(*byte)(unsafe.Add(unsafe.Pointer(ud.Player), 2251)) {
				case 0:
					value = core.Players.Mult.Warrior.Mana
				case 1:
					value = core.Players.Mult.Wizard.Mana
				case 2:
					value = core.Players.Mult.Conjurer.Mana
				default:
					value = 1
				}
				amount = int16(floatToInt32(value))
			}
			ud.ManaCur += uint16(amount)
			addProtectionRecord(int32(ud.Player.ProtUnitManaCur), uint32(int32(amount)))
			if ud.ManaCur > ud.ManaMax {
				ud.ManaCur = ud.ManaMax
				setProtectionRecord(int32(ud.Player.ProtUnitManaCur), uint32(ud.ManaMax))
			}
			if noxflags.HasGame(4096) {
				if core.Frame()%uint32(core.TickRate()>>1) == 0 {
					inventorySound(230, u, 0, 0)
				}
			} else {
				*energy--
				spent = true
			}
		}
		if spent {
			if *energy%8 == 0 {
				u.NeedSync()
			}
			if core.Frame()-u.Field34 > uint32(int32(core.TickRate())>>1) {
				inventorySound(230, u, 0, 0)
				u.Field34 = core.Frame()
			}
		}
	}
	if !idle {
		return 0
	}
	half := uint32(core.TickRate() >> 1)
	out := int32(core.Frame() / half)
	if core.Frame()%half == 0 {
		out = *energy
		if out < 50 {
			if out%8 == 0 {
				u.NeedSync()
			}
			*energy++
		}
	}
	return out
}
