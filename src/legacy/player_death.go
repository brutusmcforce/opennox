package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func playerDeathWord(p unsafe.Pointer, off uintptr) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
func playerDeathNotify(u *server.Object) {
	b := [3]byte{232, byte(u.NetCode), byte(u.NetCode >> 8)}
	reliableEnqueue(255, b[:], nil, 0, 1)
}
func playerDeath(u *server.Object) {
	s := GetServer().S()
	ud := u.UpdateDataPlayer()
	pl := ud.Player
	if *memmap.PtrUint32(0x5D4594, 2491688) == 0 {
		*memmap.PtrUint32(0x5D4594, 2491688) = uint32(s.Types.IndByID("AnkhTradable"))
	}
	if noxflags.HasGame(2048) {
		Sub_4DB170(false, nil, 0)
	}
	killer := u.Obj130
	if killer != nil {
		killer = killer.FindOwnerChainPlayer()
	}
	var assist *server.Object
	var info *server.Player
	if pl.Field3600 != 0 && s.Frame()-*playerDeathWord(pl.C(), layoutServerPlayer.Off(3608) /* field3608 */) < 10*s.TickRate() {
		info = s.Players.ByInd(ntype.PlayerInd(pl.Field3604))
		if info != nil {
			if *playerDeathWord(info.C(), layoutServerPlayer.Off(2092) /* Active */) != 0 && info.PlayerUnit != nil {
				assist = objectLookupByNetCode(info.NetCodeVal)
			} else {
				info = nil
			}
			if assist == killer || assist == u {
				assist = nil
			}
		}
	}
	if noxflags.HasGame(noxflags.GameOnline) {
		var b [11]byte
		if killer != nil && killer.ObjClass&4 != 0 {
			binary.LittleEndian.PutUint16(b[2:], uint16(killer.NetCode))
		}
		source := u.Obj130
		style, weapon := byte(ud.Field76), uint16(ud.Field75)
		if source != nil {
			if source.ObjClass&2 != 0 {
				style = 1
				weapon = source.TypeInd
			} else if source.ObjClass&4 == 0 && source.ObjOwner != nil && source.ObjOwner.ObjClass&4 == 0 && source.ObjOwner.ObjClass&2 != 0 {
				style = 1
				weapon = source.ObjOwner.TypeInd
			}
		}
		if assist != nil && assist.ObjClass&4 != 0 {
			binary.LittleEndian.PutUint16(b[4:], uint16(assist.NetCode))
		}
		binary.LittleEndian.PutUint16(b[6:], uint16(u.NetCode))
		binary.LittleEndian.PutUint16(b[8:], weapon)
		b[10] = style
		gameplayTextInformationAll(14, unsafe.Pointer(&b[0]))
		if style == 2 && weapon == 2 {
			Sub_4FC0B0(killer, 1)
		}
		ud.Field76 = 0
	}
	audio := sound.ID(321)
	if u.Field131 == 16 {
		audio = 299
	} else if *(*byte)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(2252)) /* info */)) != 0 {
		audio = 331
	}
	s.Audio.EventObj(audio, u, 0, 0)
	u.ObjFlags |= 0x8000
	Nox_xxx_playerSetState_4FA020(u, 3)
	ud.Field47_0 = 0
	clear(ud.TrapSpells[:])
	ud.TrapSpellsCnt &= 0xffffff00
	ud.SpellCastStart = 0
	if playerStateMultiple() != 0 {
		if noxflags.HasGame(256) {
			playerDeathArena(u, killer, assist, info)
		} else if noxflags.HasGame(16) {
			playerDeathKotr(u, killer)
		} else if noxflags.HasGame(1024) {
			playerDeathElimination(u, killer)
		}
	}
	if noxflags.HasGame(1024) {
		limit := serverConfigScore(1024)
		if limit != 0 && pl.Field2140 >= uint32(uint16(limit)) {
			stateRemoveSpawned(u)
		}
	}
	u.ObjFlags |= 0x10
	sessionShadowRemove(u)
	if !noxflags.HasGame(4096) {
		inventoryDropAll(u)
	}
	playerDeathNotify(u)
	ud.ManaCur = 0
	addProtectionRecord(int32(pl.ProtUnitManaCur), 0)
	stateBuffs(u, 0)
	Nox_xxx_playerCancelAbils_4FC180(u)
	pl.Field3600 = 0
	spellLifeCancelPlayer(u)
	spellLifeClearBuffs(u)
	if ud.Trade70 != nil {
		shopCancel((*shopSession)(unsafe.Pointer(ud.Trade70)))
	}
	ud.Trade70 = nil
	if noxflags.HasGame(4096) {
		if ud.Field80 != 0 {
			ud.Field80--
			questRuntimeIncrement(u, 4660, 2)
		} else {
			ud.Field137 = s.Frame()
			var b [14]byte
			b[0] = 240
			b[1] = 2
			for i, off := range []uintptr{4668, 4672, 4664, 4688} {
				binary.LittleEndian.PutUint16(b[2+2*i:], *(*uint16)(unsafe.Add(pl.C(), off)))
			}
			reliableEnqueue(int(pl.PlayerIndex()), b[:], nil, 1, 0)
			questRuntimeReset(u)
			questDeathPenalty(u)
			ud.Field80 = uint32(floatToInt32(float32(s.Balance.Float("QuestGameStartingExtraLives"))))
			*(*byte)(unsafe.Add(unsafe.Pointer(ud), 452+uintptr(pl.PlayerIndex()))) = byte(ud.Field80)
		}
	}
}
