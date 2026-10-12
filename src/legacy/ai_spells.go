package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/things"
	"github.com/brutusmcforce/opennox/libs/types"
	flags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func monsterCastBusy(u *server.Object) bool {
	a := monsterHeadSafe(u)
	return a >= 18 && a <= 20
}
func monsterSpellReady(u *server.Object, until uint32) bool {
	return u.Flags()&0x1000000 != 0 && monsterCanCast(u) && GetServer().S().Frame() >= until && !monsterCastBusy(u)
}
func monsterSpellCandidates(u *server.Object, mask uint32, out *[136]spell.ID) []spell.ID {
	permissions := (*[136]uint32)(unsafe.Pointer(&u.UpdateDataMonster().Field373))
	n := 0
	for i, v := range permissions {
		id := spell.ID(i + 1)
		if v&mask != 0 && GetServer().S().Spells.HasFlags(id, things.SpellFlags(16)) {
			out[n] = id
			n++
		}
	}
	return out[:n]
}
func monsterSpellEnchantActive(u *server.Object, id spell.ID) bool {
	for i := 0; i < 32; i++ {
		enc := server.EnchantID(i)
		if enc.Spell() == id {
			return u.HasEnchant(enc)
		}
	}
	return false
}
func monsterSummonSpell(id spell.ID) bool { return id >= 75 && id <= 114 }
func monsterSummonActive(u *server.Object) bool {
	for d := GetServer().S().Spells.Dur.List; d != nil; d = d.Next {
		if monsterSummonSpell(spell.ID(d.Spell)) && d.Caster16 == u {
			return true
		}
	}
	return false
}
func monsterSpellSchedule(u, target *server.Object, id spell.ID, lo, hi *uint16, until *uint32) {
	u.MonsterCast(id, target)
	// Read delay bounds after action callbacks, retaining the original update pointer.
	core := GetServer().S()
	*until = core.Frame() + uint32(core.Rand.Logic.IntClamp(int(*lo), int(*hi)))
}
func monsterSelectSpell(u, target *server.Object, mask uint32, lo, hi *uint16, until *uint32, rejectEnchant bool) bool {
	if !monsterSpellReady(u, *until) {
		return false
	}
	var buf [136]spell.ID
	list := monsterSpellCandidates(u, mask, &buf)
	if len(list) == 0 {
		return false
	}
	if rejectEnchant {
		for _, id := range list {
			if monsterSpellEnchantActive(u, id) {
				return false
			}
		}
	}
	id := list[GetServer().S().Rand.Logic.IntClamp(0, len(list)-1)]
	monsterSpellSchedule(u, target, id, lo, hi, until)
	return true
}
func monsterBuffSelf(u *server.Object) bool {
	d := u.UpdateDataMonster()
	return monsterSelectSpell(u, u, 0x10000000, &d.Field364_0, &d.Field364_2, &d.Field365, true)
}
func monsterCastOffensive(u, t *server.Object) bool {
	d := u.UpdateDataMonster()
	return monsterSelectSpell(u, t, 0x20000000, &d.Field366_0, &d.Field366_2, &d.Field367, false)
}
func monsterCastRelated(u *server.Object) bool {
	d := u.UpdateDataMonster()
	return monsterSelectSpell(u, u, 0x80000000, &d.Field370_0, &d.Field370_2, &d.Field371, true)
}
func monsterMagicMissile(t, u *server.Object) {
	if t.Class()&1 != 0 && t.SubClass()&2 != 0 && *(**server.Object)(unsafe.Add(t.UpdateData, 4)) == u {
		*memmap.PtrUint32(0x5D4594, 2489156) = 1
	}
}
func monsterCastInversion(u *server.Object) bool {
	d := u.UpdateDataMonster()
	if !monsterSpellReady(u, d.Field363) {
		return false
	}
	found := memmap.PtrUint32(0x5D4594, 2489156)
	*found = 0
	core := GetServer().S()
	radius := float32(core.Balance.Float("InversionRange") * .5)
	core.Map.EachMissileInCircle(u.PosVec, radius, func(t *server.Object) bool { monsterMagicMissile(t, u); return true })
	if *found == 0 {
		return false
	}
	var buf [136]spell.ID
	list := monsterSpellCandidates(u, 0x08000000, &buf)
	if len(list) == 0 {
		return false
	}
	id := list[core.Rand.Logic.IntClamp(0, len(list)-1)]
	monsterSpellSchedule(u, u, id, &d.Field362_0, &d.Field362_2, &d.Field363)
	return true
}
func monsterCastRelated2(u, t *server.Object) bool {
	d := u.UpdateDataMonster()
	if !monsterSpellReady(u, d.Field369) {
		return false
	}
	var buf [136]spell.ID
	list := monsterSpellCandidates(u, 0x40000000, &buf)
	if len(list) == 0 {
		return false
	}
	allSummons := true
	for _, id := range list {
		if !monsterSummonSpell(id) {
			allSummons = false
		}
	}
	for {
		id := list[GetServer().S().Rand.Logic.IntClamp(0, len(list)-1)]
		if monsterSummonSpell(id) && (monsterSummonActive(u) || !Nox_xxx_checkSummonedCreaturesLimit_500D70(u, int(id)-74)) {
			if allSummons {
				return false
			}
			continue
		}
		monsterSpellSchedule(u, t, id, &d.Field368_0, &d.Field368_2, &d.Field369)
		return true
	}
}
func monsterHealCandidate(t, u *server.Object) {
	core := GetServer().S()
	if t != u && !core.IsEnemyTo(u, t) && t.HealthData != nil && t.Flags()&0x8000 == 0 && core.CanInteract(u, t, 0) && t.HealthData.Cur < t.HealthData.Max>>1 {
		dword_5d4594_2489160 = uint32(uintptr(t.CObj()))
	}
}
func monsterHealSomeone(u *server.Object) bool {
	core := GetServer().S()
	d := u.UpdateDataMonster()
	if u.Flags()&0x1000000 == 0 || !monsterCanCast(u) || core.Frame()&31 != 0 || monsterCastBusy(u) {
		return false
	}
	status := d.StatusFlags
	if status&0x800 != 0 && u.HealthData.Max != 0 && u.HealthData.Cur < u.HealthData.Max>>1 {
		u.MonsterCast(41, u)
		return false
	}
	if status&0x1000 == 0 {
		return false
	}
	radius := float64(250)
	if flags.HasGame(4096) {
		radius = 640
	}
	dword_5d4594_2489160 = 0
	rect := types.Rectf{
		Min: types.Pointf{X: float32(float64(u.PosVec.X) - radius), Y: float32(float64(u.PosVec.Y) - radius)},
		Max: types.Pointf{X: float32(float64(u.PosVec.X) + radius), Y: float32(float64(u.PosVec.Y) + radius)},
	}
	core.Map.EachObjInRect(rect, func(t *server.Object) bool { monsterHealCandidate(t, u); return true })
	if dword_5d4594_2489160 == 0 {
		return false
	}
	u.MonsterCast(41, (*server.Object)(unsafe.Pointer(uintptr(dword_5d4594_2489160))))
	return true
}
func monsterCastSpell(id int, u *server.Object, args *[3]uint32) {
	d := u.UpdateDataMonster()
	if d.StatusFlags&0x20000 != 0 {
		controlMorphToPlayer((*server.Object)(unsafe.Pointer(u.CObj())))
	}
	monsterCalcDir(u, (*float32)(unsafe.Pointer(&args[1])))
	Nox_xxx_castSpellByUser_4FDD20(id, u, unsafe.Pointer(args))
	if d.StatusFlags&0x20000 != 0 {
		controlMorphFromPlayer((*server.Object)(unsafe.Pointer(u.CObj())))
	}
}
func monsterActionCast(u *server.Object, mode int) {
	d := u.UpdateDataMonster()
	head := d.AIStackHead()
	if d.Field120_2 != 0 {
		return
	}
	if uint32(d.Field120_1) == d.MonsterDef.MissileAttackFrame216 {
		if mode != 0 && mode != 1 {
			return
		}
		// The cast engine can pass these words through C and back into Go.
		args, freeArgs := alloc.New([3]uint32{})
		defer freeArgs()
		if mode == 1 {
			*args = [3]uint32{0, uint32(head.Args[2]), uint32(head.Args[3])}
		} else {
			args[0] = uint32(head.Args[2])
			monsterCastRecoil(u, head.ArgObj(2), (*types.Pointf)(unsafe.Pointer(&args[1])))
		}
		if !u.HasEnchant(29) {
			monsterCastSpell(int(head.Args[0]), u, args)
		}
	} else if d.Field120_1 == 1 {
		combatSound(u, 14)
	}
}
func monsterCastRecoil(u, t *server.Object, out *types.Pointf) {
	inv := 1 - float64(u.UpdateDataMonster().Field330)
	lo, hi := float32(inv), float32(inv+1)
	rng := GetServer().S().Rand.Logic
	r := rng.FloatClamp(float64(lo), float64(hi))
	out.X = float32(float64(t.PosVec.X) - r*float64(t.VelVec.X)*6)
	out.Y = float32(float64(t.PosVec.Y) - r*float64(t.VelVec.Y)*6)
	scale := float32(float64(lo)*.80000001 + .2)
	out.X = float32(rng.FloatClamp(-60, 60)*float64(scale) + float64(out.X))
	out.Y = float32(rng.FloatClamp(-60, 60)*float64(scale) + float64(out.Y))
}
