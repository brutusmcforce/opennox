package legacy

import (
	"math"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func spellLifeWord(p unsafe.Pointer, off int) *uint32 { return rawAt[uint32](p, off) }
func spellLifeManaCost(u *server.Object, id, mode int32) int32 {
	if id >= 75 && id <= 114 {
		return spellEffectSummonCost(id, u)
	}
	return int32(GetServer().S().Spells.ManaCost(spell.ID(id), int(mode)))
}
func spellLifeCheckMana(u *server.Object, list unsafe.Pointer, n int32) int32 {
	if u == nil || list == nil || n == 0 {
		return 0
	}
	if noxflags.HasEngine(noxflags.EngineGodMode) || u.ObjClass&2 != 0 {
		return 1
	}
	mana := int32(uint16(resourceGetMana(u)))
	for i := int32(0); i < n; i++ {
		cost := spellLifeManaCost(u, int32(*spellLifeWord(list, int(i)*4)), 2)
		if cost > mana {
			return 0
		}
		mana -= cost
	}
	return 1
}
func spellLifeSpendMana(u *server.Object, id, mode int32) int32 {
	if u.ObjClass&4 == 0 || id == 0 {
		return -1
	}
	if noxflags.HasEngine(noxflags.EngineGodMode) {
		return 0
	}
	cost := spellLifeManaCost(u, id, mode)
	if int32(u.UpdateDataPlayer().ManaCur) >= cost {
		resourceSubMana(u, cost)
		return cost
	}
	layout.As[server.PlayerUpdateData](u.UpdateData).Field20_0 = uint16(GetServer().S().Spells.ManaCost(spell.ID(id), 1))
	layout.As[server.PlayerUpdateData](u.UpdateData).Field20_1 = uint16(GetServer().S().TickRate())
	return -1
}
func spellLifeRefundMana(u *server.Object, amount int16) uint16 { return resourceAddMana(u, amount) }
func spellLifeCheckClass(u *server.Object, id int32) int32 {
	parent := u.FindOwnerChainPlayer()
	if !GetServer().S().Spells.DefByInd(spell.ID(id)).IsEnabled() {
		return 10
	}
	if u.ObjClass&4 != 0 {
		return playerSpellClassCheck(int32(controlPlayer(u).InfoBytes[66]), int32(id))
	}
	v := -int32(bool2int(server.Sub_57AEE0(spell.ID(int(int32(id))), parent)))
	v = int32(uint32(v)&0xffffff00 | uint32(byte(v)&0xf6))
	return v + 10
}
func spellLifePower(id int32, u *server.Object) int32 {
	typ := stateType(1569720, "ImaginaryCaster")
	if uint32(u.TypeInd) == typ {
		return 1
	}
	if controlFlags(1392) {
		return 3
	}
	if u.ObjClass&4 != 0 {
		return int32(controlPlayer(u).SpellLvl[int(id)])
	}
	if u.ObjClass&2 != 0 {
		return int32(layout.As[server.MonsterUpdateData](u.UpdateData).Field510)
	}
	return 3
}
func spellLifeMoved(u *server.Object, p *types.Pointf) int32 {
	dx := float32(math.Abs(float64(p.X) - float64(u.PosVec.X)))
	dy := math.Abs(float64(p.Y) - float64(u.PosVec.Y))
	if dx >= 5 || dy >= 5 {
		return 1
	}
	return 0
}
