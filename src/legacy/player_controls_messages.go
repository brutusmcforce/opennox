package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func controlRewardNotify(u *server.Object, kind int32, target *server.Object, selector byte) int32 {
	if u == nil || u.ObjClass&4 == 0 {
		return int32(controlRaw(u))
	}
	if kind < 0 || kind > 2 {
		return kind - 2
	}
	msg := [5]byte{0xf0, byte(30 + kind), selector}
	binary.LittleEndian.PutUint16(msg[3:], uint16(target.NetCode))
	return int32(reliableEnqueue(int(controlPlayer(u).PlayerInd), msg[:], nil, 1, 0))
}
func controlLockedDoor(u *server.Object, key *int8, selector byte) {
	if u == nil || u.ObjClass&4 == 0 || key == nil {
		return
	}
	s := GoStringP(unsafe.Pointer((*int8)(key)))
	if len(s) == 0 || len(s) > 48 {
		return
	}
	// The fixed wire format includes unused bytes; keep their padding zero.
	msg := [52]byte{0xf0, 33}
	copy(msg[2:], s)
	msg[51] = selector
	reliableEnqueue(int(controlPlayer(u).PlayerInd), msg[:], nil, 1, 0)
}
func controlRespawnNotify(u *server.Object, flag byte) int32 {
	msg := [9]byte{0xe9}
	binary.LittleEndian.PutUint16(msg[1:], uint16(u.NetCode))
	binary.LittleEndian.PutUint32(msg[3:], GetServer().S().Frame())
	msg[7] = byte(controlRespawnFlags())
	msg[8] = flag
	return int32(reliableEnqueue(255, msg[:], nil, 0, 1))
}
func controlGuideLevel(u, target *server.Object) int32 {
	if u == nil || target == nil || u.ObjClass&4 == 0 || target.ObjClass&2 == 0 {
		return 0
	}
	name := nox_xxx_getUnitName_4E39D0((*nox_object_t)(target.CObj()))
	id := bookGuideID(GoStringP(unsafe.Pointer(name)))
	if id == 0 {
		return 0
	}
	return int32(controlPlayer(u).BeastScrollLvl[int(id)])
}
func controlGuideDamage(u, target *server.Object, damage *int32) int32 {
	level := controlGuideLevel(u, target)
	if level == 0 {
		return 0
	}
	v := float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("FieldGuideDamageBonus")))*float64(*damage) + 0.5)
	*damage = floatToInt32(v)
	return *damage
}
func controlInversion(u, target *server.Object) int32 {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 == 0 || it.ObjClass&0x13001000 == 0 {
			continue
		}
		for i := 2; i < 4; i++ {
			mod := it.InitDataModifier().Modifiers[i]
			if mod == nil {
				continue
			}
			fn := mod.DefendCollide88.Fnc
			if fn == nil || fn != modifierKey(modifierIDInversionEffect) {
				continue
			}
			var result int32
			owner := target.FindOwnerChainPlayer()
			server.CallModifierEffect6(fn, (*server.ModifierEff)(mod), it, u, target, owner, unsafe.Pointer(&result))
			if result == 1 {
				return 1
			}
		}
	}
	return 0
}
func controlScheduledSpell(u, target *server.Object, back bool) int32 {
	d := u.UpdateDataPlayer()
	count := (*byte)(unsafe.Pointer(&d.TrapSpellsCnt))
	if *count == 0 {
		return 0
	}
	id := d.TrapSpells[0]
	if back {
		id = d.TrapSpells[*count-1]
	}
	reason := int32(spellLifeCantCast(u, int32(id), 0))
	args := server.SpellAcceptArg{Obj: target, Pos: types.Pointf{X: float32(int32(d.Field55)), Y: float32(int32(d.Field56))}}
	if reason != 0 {
		nox_xxx_netInformTextMsg_4DA0F0(int32(controlPlayer(u).PlayerInd), 0, &reason)
		nox_xxx_aud_501960(231, (*nox_object_t)(u.CObj()), 0, 0)
	} else {
		nox_xxx_castSpellByUser_4FDD20(int(int32(id)), (*nox_object_t)(u.CObj()), unsafe.Pointer(&args))
	}
	// Front removal shifts and clears; back removal only decrements the count.
	if !back {
		for i := 1; i < int(*count); i++ {
			d.TrapSpells[i-1] = d.TrapSpells[i]
		}
		d.TrapSpells[*count-1] = 0
	}
	*count--
	return 1
}
