package server

import (
	"testing"
	"unsafe"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/object"
)

// Update effects (Regeneration, ContinualReplenishment) get the item, as in the original, and find the holder themselves.
func TestItemsApplyUpdateEffect(t *testing.T) {
	key := unsafe.Pointer(new(byte))
	var got []*Object
	RegisterModifierEffect3(key, func(_ *ModifierEff, a, _ *Object) int32 {
		got = append(got, a)
		return 0
	})
	t.Cleanup(func() { delete(modifierEffects3, key) })
	mod := &ModifierEff{}
	mod.Update100.Fnc = key

	var s Server
	staff := &Object{ObjClass: object.ClassWand, ObjFlags: object.FlagEquipped}
	unequipped := &Object{ObjClass: object.ClassWand}
	for _, it := range []*Object{staff, unequipped} {
		it.InitData = unsafe.Pointer(&ModifierInitData{Modifiers: [4]*ModifierEff{2: mod}})
	}
	player := &Object{ObjClass: object.ClassPlayer, InvFirstItem: staff}
	staff.InvNextItem = unequipped
	staff.InvHolder, unequipped.InvHolder = player, player
	s.ItemsApplyUpdateEffect(player)
	must.Eq(t, []*Object{staff}, got)
}

// As the original: the object itself counts, so a player is its own player even if it has an owner.
func TestFindOwnerChainPlayer(t *testing.T) {
	player := &Object{ObjClass: object.ClassPlayer}
	other := &Object{ObjClass: object.ClassPlayer}
	player.ObjOwner = other
	must.EqOp(t, player, player.FindOwnerChainPlayer())
	missile := &Object{ObjClass: object.ClassMissile}
	summon := &Object{ObjClass: object.ClassMonster, ObjOwner: player}
	missile.ObjOwner = summon
	must.EqOp(t, player, missile.FindOwnerChainPlayer())
	lone := &Object{ObjClass: object.ClassMissile}
	must.EqOp(t, lone, lone.FindOwnerChainPlayer())
}

// Scripts read these from monsters; anything else returns 0 instead of panicking.
func TestMonsterAIGetters(t *testing.T) {
	var ud MonsterUpdateData
	ud.RetreatLevel = 0.25
	mon := &Object{ObjClass: object.ClassMonster, UpdateData: unsafe.Pointer(&ud)}
	must.EqOp(t, float32(0.25), mon.RetreatLevel())
	item := &Object{ObjClass: object.ClassFood}
	must.EqOp(t, float32(0), item.RetreatLevel())
	must.EqOp(t, float32(0), item.Aggression())
	must.EqOp(t, float32(0), item.RegroupLevel())
}
