package legacy

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func sustainedGlobal(i int) *uint32 {
	switch i {
	case 0:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487884))
	case 1:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487932))
	case 2:
		return (*uint32)(unsafe.Pointer(&nox_xxx_lightningOwner_5d4594_2487900))
	case 3:
		return (*uint32)(unsafe.Pointer(&nox_xxx_lightningTargetArrayIndex_5d4594_2487904))
	case 4:
		return (*uint32)(unsafe.Pointer(&nox_xxx_lightningTarget_5d4594_2487908))
	case 5:
		return (*uint32)(unsafe.Pointer(&nox_xxx_lightningClosestTargetDistance_5d4594_2487912))
	case 6:
		return (*uint32)(unsafe.Pointer(&nox_xxx_energyBoltTarget_5d4594_2487880))
	}
	panic("sustained spell global")
}
func sustainedPointer(v uint32) *server.Object { return (*server.Object)(memmap.WordPtr(v)) }
func sustainedFrame() uint32                   { return GetServer().S().Frame() }
func sustainedFPS() uint32                     { return uint32(GetServer().S().TickRate()) }
func sustainedGame(flag uint32) bool           { return noxflags.HasGame(noxflags.GameFlag(flag)) }
func sustainedAudio(id int, u *server.Object)  { GetServer().S().Audio.EventObj(sound.ID(id), u, 0, 0) }
func sustainedFX(code byte, pos types.Pointf) uint32 {
	return uint32(visibilityFXPoint(byte(code), pos))
}
func sustainedState(u *server.Object, state int) {
	Nox_xxx_playerSetState_4FA020(u, server.PlayerState(int(int32(state))))
}
func sustainedHurtRecently(u *server.Object) bool {
	return u.HealthData != nil && sustainedFrame()-u.Frame134 <= 1
}
func sustainedFront(a, b *server.Object) int32 {
	return stateFront(&a.PosVec, int32(int16(a.Direction1)), &b.PosVec)
}
func sustainedInteract(a, b *server.Object) bool { return GetServer().S().CanInteract(a, b, 0) }
func sustainedEnemy(a, b *server.Object) bool    { return GetServer().S().IsEnemyTo(a, b) }
func sustainedNew(name string) *server.Object {
	namePtr := internCStr(name)
	return GetServer().S().NewObjectByTypeID(GoStringP(unsafe.Pointer(namePtr)))
}
func sustainedDelete(u *server.Object)                      { nox_xxx_delayedDeleteObject_4E5CC0(asObjectC(u)) }
func sustainedStopRay(d *server.DurSpell, u *server.Object) { GetServer().S().NetStopRaySpell(d, u) }
func sustainedScalarInt(name string) int32                  { return floatToInt32(float32(spellEffectScalar(name))) }
func sustainedTableInt(name string, level uint32) int32 {
	return floatToInt32(float32(spellEffectTable(name, int32(level)-1)))
}
func sustainedShieldFX(u, source *server.Object) {
	var pos *types.Pointf
	if source != nil {
		pos = &source.PosVec
	}
	visibilityFXShield(u, pos)
}
func sustainedCharges(u, item *server.Object) {
	data := item.UseDataWand()
	gameplayReportCharges(int(controlPlayer(u).PlayerInd), item, data.Charges, data.MaxCharges)
}
func sustainedTagPacket(u, target *server.Object, mode byte) uint32 {
	var b [7]byte
	b[0] = 210
	binary.LittleEndian.PutUint16(b[1:], uint16(nox_xxx_netGetUnitCodeServ_578AC0(asObjectC(target))))
	binary.LittleEndian.PutUint16(b[3:], target.TypeInd)
	b[5] = mode
	b[6] = 1
	return uint32(reliableEnqueue(int(layout.As[server.Player](unsafe.Pointer(controlPlayer(u))).PlayerInd), b[:], nil, 1, 0))
}
func sustainedTagStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	if u == nil || u.ObjFlags&0x8020 != 0 || u.ObjClass&4 == 0 {
		return 1
	}
	t := d.Target48
	if t == nil || t.ObjFlags&0x8020 != 0 {
		return 1
	}
	d.Frame68 = sustainedFrame() + d.Level*uint32(sustainedScalarInt("TagDurationPerLevel"))
	nox_xxx_netMarkMinimapObject_417190(int(int32(layout.As[server.Player](unsafe.Pointer(controlPlayer(u))).PlayerInd)), asObjectC(t), 1)
	sustainedTagPacket(u, t, 1)
	return 0
}
func sustainedTagTick(p unsafe.Pointer) uint32 {
	u := (*server.DurSpell)(p).Target48
	if u == nil {
		return 1
	}
	return uint32(u.ObjFlags) >> 5 & 1
}
func sustainedTagCancel(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Caster16
	ret := controlRaw(u)
	if u != nil && u.ObjClass&4 != 0 {
		t := d.Target48
		ret = controlRaw(t)
		if t != nil {
			if t.ObjClass&4 == 0 {
				nox_xxx_netUnmarkMinimapObj_417300(int(int32(controlPlayer(u).PlayerInd)), asObjectC(t), 1)
			}
			ret = sustainedTagPacket(u, t, 2)
		}
	}
	return ret
}
func sustainedOvalShieldStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	duration := 20 * d.Level * sustainedFPS()
	if u == nil || u.ObjClass&4 == 0 {
		return 1
	}
	GetServer().S().Spells.Dur.CancelOffensiveFor(u)
	spellLifeApplyBuff(u, 27, int16(duration), int8(d.Level))
	d.Frame68 = duration + sustainedFrame()
	return 0
}
func sustainedOvalShieldTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	if u == nil || spellLifeHasBuff(u, 8) {
		return 1
	}
	if u.ObjClass&2 != 0 && spellLifeMoved(u, &d.Pos) != 0 {
		return 1
	}
	return uint32(bool2int(u.ObjFlags&0x8020 != 0))
}
func sustainedOvalShieldCancel(p unsafe.Pointer) uint32 {
	u := (*server.DurSpell)(p).Target48
	if u == nil {
		return 0
	}
	return spellLifeBuffOff(u, 27)
}
func sustainedMoonglowStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	dur := int16(sustainedScalarInt("MoonglowEnchantmentDuration"))
	u := d.Target48
	if u == nil || u.ObjFlags&0x8020 != 0 {
		return 1
	}
	if u.ObjClass&4 != 4 {
		spellLifeApplyBuff(u, 15, dur, int8(d.Level))
		return 1
	}
	light := sustainedNew("Moonglow")
	d.Field72 = int32(memmap.PtrWord(light.CObj())) // a pointer word
	if light == nil {
		return 1
	}
	pl := controlPlayer(u)
	pos := types.Ptf(2944, 2944)
	if *(*byte)(unsafe.Pointer(&layout.As[server.Player](unsafe.Pointer(pl)).Field3680))&0x10 != 0 {
		pos = types.Ptf(float32(int32(layout.As[server.Player](unsafe.Pointer(pl)).CursorVec.X)), float32(int32(layout.As[server.Player](unsafe.Pointer(pl)).CursorVec.Y)))
	}
	spellEffectCreate(light, u, pos)
	spellLifeApplyBuff(u, 1, dur, int8(d.Level))
	return 0
}
func sustainedMoonglowCancel(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 4 {
		return spellLifeBuffOff(u, 15)
	}
	if light := (*server.Object)(memmap.WordPtr(uint32(d.Field72))); light != nil {
		sustainedDelete(light)
	}
	d.Field72 = 0
	return spellLifeBuffOff(u, 1)
}
