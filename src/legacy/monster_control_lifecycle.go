package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

func monsterControlDeath(u *server.Object) uint32 {
	s := GetServer().S()
	ud := u.UpdateDataMonster()
	if noxflags.HasGame(4096) {
		spawnPolicyDeathRelease(u)
	}
	for p := s.Players.FirstUnit(); p != nil; p = s.Players.NextUnit(p) {
		if Nox_xxx_playerGetPossess_4DDF30(p) == u {
			Nox_xxx_playerObserveClear_4DDEF0(p)
		}
	}
	u.ClearActionStack()
	u.MonsterPushAction(ai.ACTION_DEAD)
	u.MonsterPushAction(ai.ACTION_DYING)
	if monsterIsZombie(u) {
		return 1
	}
	u.ObjFlags &^= 0x80
	collisionActivate(u)
	spellLifeClearBuffs(u)
	if ud.StatusFlags&0x80 != 0 {
		motionDecaySet(u, int32(s.TickRate()*uint32(s.Rand.Logic.IntClamp(10, 20))))
	} else if noxflags.HasGame(4096) {
		motionDecaySet(u, int32(s.TickRate()*uint32(s.Rand.Logic.IntClamp(5, 8))))
	}
	if owner := u.ObjOwner; owner != nil && owner.ObjClass&4 != 0 {
		pl := owner.UpdateDataPlayer().Player
		u.ObjSubClass &^= 0x80
		gameplayReportShield(int(pl.PlayerInd), u)
		s.Players.Nox_xxx_netUnmarkMinimapObj_417300(ntype.PlayerInd(pl.PlayerInd), u, 1)
	}
	u.ObjSubClass &^= 0x100
	controlTransferChildren(u)
	s.ObjClearOwner(u)
	if u.ObjSubClass&0x2000 == 0 {
		inventoryDropAll(u)
	}
	if !noxflags.HasGame(2048) && ud.Field547 == 2 && ud.Field546 == 2 && u.Obj130 != nil {
		killer := u.Obj130.FindOwnerChainPlayer()
		if killer.ObjClass&4 != 0 {
			Sub_4FC0B0(killer, 1)
		}
	}
	if !noxflags.HasGame(4096) {
		return 0
	}
	if u.Obj130 != nil {
		killer := u.Obj130.FindOwnerChainPlayer()
		if killer.ObjClass&4 != 0 {
			return questRuntimeIncrement(killer, 4664, 4)
		}
		return motionAddress(killer)
	}
	return 1
}
func monsterControlChapter() int {
	var player *server.Player
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		player = u.UpdateDataPlayer().Player
		if player.PlayerInd == 31 {
			break
		}
	}
	Set_nox_gameDisableMapDraw_5d4594_2650672(1)
	if player == nil {
		return 0
	}
	return gameplayReportChapter(int(player.PlayerInd), memmap.Uint8(0x5D4594, 2386828), int(memmap.Int32(0x5D4594, 2386832)))
}
