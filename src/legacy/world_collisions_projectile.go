package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func worldCollideSpellProjectile(a, b *server.Object, normal *types.Pointf) {
	d := a.UpdateData
	if b == nil {
		if normal != nil {
			collisionReflect(normal, &a.VelVec)
		}
		return
	}
	if b.ObjFlags&0x8020 != 0 {
		return
	}
	if b.ObjClass&4 != 0 {
		bd := b.UpdateData
		state := *controlByte(bd, 88)
		shield := state == 16
		skipSword := false
		if !shield && state == 1 && controlActionState(b) == 45 {
			class := *controlByte(controlPlayer(b), 3)
			if class == 1 || class == 2 {
				if uint32(gameex_flags)&16 == 0 {
					skipSword = true
				} else {
					shield = true
				}
			}
		}
		if shield && stateFront(&b.PosVec, int32(int16(b.Direction1)), &a.PrevPos)&1 != 0 {
			worldCollideSound(878, b)
			damageReflect(a, b)
			Nox_xxx_changeOwner_52BE40(a, b)
			return
		}
		if !skipSword && (state == 13 || state == 0 && uint32(gameex_flags)&4 != 0) && *equipmentWord(controlPlayer(b), 4)&0x400 != 0 && stateFront(&b.PosVec, int32(int16(b.Direction1)), &a.PrevPos)&1 != 0 {
			state := GetServer().S().Rand.Logic.IntClamp(18, 20)
			worldCollideSound(890, b)
			damageReflect(a, b)
			Nox_xxx_changeOwner_52BE40(a, b)
			Nox_xxx_playerSetState_4FA020(b, server.PlayerState(state))
			frames, _ := GetServer().S().PlayerAnimFrames(int(controlActionState(b)))
			*controlByte(bd, 236) = byte(frames - 1)
		}
		if controlInversion(b, a) != 0 {
			Nox_xxx_changeOwner_52BE40(a, b)
			return
		}
		if b.Buffs&(1<<27) != 0 && stateFront(&b.PosVec, int32(int16(b.Direction1)), &a.PosVec)&1 != 0 {
			Nox_xxx_changeOwner_52BE40(a, b)
			worldCollideSound(122, b)
			return
		}
	}
	if controlObject(d, 4) == b {
		// The C caller supplied only the target object in this argument record.
		args := server.SpellAcceptArg{Obj: b}
		GetServer().Nox_xxx_spellAccept4FD400(spell.ID(int32(*equipmentWord(d, 12))), controlObject(d, 8), controlObject(d, 0), a, &args, int(int32(*equipmentWord(d, 16))))
		GetServer().DelayedDelete(a)
	} else {
		damageReflect(a, b)
	}
}
