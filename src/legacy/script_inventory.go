package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"github.com/brutusmcforce/opennox/server/noxscript"
)

var scriptInventoryReserved int32
var scriptInventoryNotice uint32

func scriptInventoryCarry(u, it *server.Object) {
	glyph := memmap.PtrUint32(0x5D4594, 2386856)
	if *glyph == 0 {
		*glyph = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	if uiInventoryCapacity(int32(it.TypeInd), 1)-scriptInventoryReserved > 0 {
		return
	}
	var selected *server.Object
	price := int32(999999)
	for item := u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.ObjClass&0x10 != 0 || item.ObjFlags&0x100 != 0 || uint32(item.TypeInd) == *glyph || inventoryDroppable(item) {
			continue
		}
		if cost := shopPrice(1, nil, item); cost < price {
			price = cost
			selected = item
		}
	}
	if selected == nil {
		return
	}
	var pos types.Pointf
	inventoryRandomPlacement(50, &u.PosVec, &pos)
	inventoryDrop(u, selected, &pos)
	if scriptInventoryNotice == 0 {
		gameplayTextPrivate(u, alloc.InternCString("pickup.c:CarryingTooMuch"), 0)
		scriptInventoryNotice = 1
	}
}
func scriptInventoryStartup() {
	s := GetServer().S()
	var host *server.Object
	for host = s.Players.FirstUnit(); host != nil; host = s.Players.NextUnit(host) {
		if host.UpdateDataPlayer().Player.PlayerInd == 31 {
			break
		}
	}
	journalRemoveMask(host, 0xE)
	for it := host.InvFirstItem; it != nil; {
		next := it.InvNextItem
		if it.ObjClass&0x40 != 0 {
			GetServer().DelayedDelete(it)
		}
		it = next
	}
	*memmap.PtrUint32(0x5D4594, 2386832) = 1
}
func scriptInventoryHalberd(vm noxscript.VM) int {
	host := GetServer().S().Players.ByInd(31).PlayerUnit
	index := vm.PopI32()
	equipped := false
	for it := host.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x7800000 != 0 {
			equipped = it.ObjFlags&0x100 != 0
			GetServer().DelayedDelete(it)
			break
		}
	}
	name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, uintptr(247336+4*index))))
	item := controlRespawnItem(host, name, nil, 1, 1)
	if equipped {
		equipmentTryEquip(host, item)
	}
	return 0
}
