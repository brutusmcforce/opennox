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
