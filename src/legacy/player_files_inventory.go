package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func playerFileInventory(u *server.Object) int {
	data := u.UpdateData
	core := GetServer().S()
	r := playerFileStream()
	cached := memmap.PtrUint32(0x5D4594, 527704)
	if *cached == 0 {
		*cached = uint32(core.Types.IndByID("Glyph"))
	}
	if r.read() {
		controlLevelFromXP(u)
	}
	version := int16(r.short(3))
	if version > 3 {
		return 0
	}
	present := playerFilePresent(r, noxflags.HasGame(8192) && !noxflags.HasGame(4096))
	if present {
		if !noxflags.HasGame(2048 | 4096) {
			return 0
		}
		p := u.UpdateDataPlayer().Player
		gold := r.word(p.GoldVal)
		if r.read() {
			resourceSubGold(u, resourceGetGold(u))
			resourceAddGold(u, gold)
			for it := u.InvFirstItem; it != nil; {
				next := it.InvNextItem
				GetServer().DelayedDelete(it)
				it = next
			}
			count := int32(r.word(0))
			if noxflags.HasGame(4096) && count > 2560 {
				return 0
			}
			for i := int32(0); i < count; i++ {
				name := playerFileName(r, "")
				it := core.NewObjectByTypeID(name)
				if it == nil {
					return 0
				}
				if err := it.CallXfer(nil); err != nil {
					return 0
				}
				it.PosVec.X = 2944
				it.PosVec.Y = 2944
				if noxflags.HasGame(4096) && !questEligibilityItem(it) {
					return 0
				}
				objectXferPlace(it, u.CObj(), nil)
				GetServer().ObjectsAddPending()
				if !Nox_xxx_inventoryServPlace_4F36F0(u, it, 1, 1) {
					if !noxflags.HasGame(4096) {
						return 0
					}
					GetServer().DelayedDelete(it)
				}
				if uint32(it.ObjFlags)&0x20 == 0 && uint32(it.ObjFlags)&0x100 != 0 {
					equipmentTryDequip(u, it)
				}
			}
			equipped := r.byte(0)
			for i := 0; i < int(equipped); i++ {
				id := r.word(0)
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					if uint32(it.ScriptIDVal) == id {
						equipmentTryEquip(u, it)
					}
				}
			}
			if noxflags.HasGame(2048) {
				sub_467750(0, 0)
				sub_467740(0)
			}
			id := r.word(0)
			if id != 0 {
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					if uint32(it.ScriptIDVal) == id {
						gameplayReportSecondary(int(uint8(p.PlayerInd)), it, 0)
						break
					}
				}
			}
			if version >= 2 {
				id = r.word(0)
				if id != 0 {
					for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
						if uint32(it.ScriptIDVal) == id {
							gameplayReportQuiver(int(uint8(p.PlayerInd)), it)
							break
						}
					}
				}
			}
			if noxflags.HasGame(4096) {
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					it.ScriptIDVal = int(core.Objs.NextObjectScriptID())
					it.Extent = it.NetCode
				}
			}
		} else {
			count := uint32(0)
			for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
				if !noxflags.HasGame(4096) || playerFileInventoryAllowed(it) {
					count++
				}
			}
			r.word(count)
			gridCount := 0
			if noxflags.HasGame(2048) {
				gridCount = playerFileInventoryCount()
			}
			writeItem := func(it *server.Object) bool {
				playerFileName(r, core.Types.ByInd(int(it.TypeInd)).ID())
				return it.CallXfer(nil) == nil
			}
			if count != uint32(gridCount) || !noxflags.HasGame(2048) {
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					if !noxflags.HasGame(4096) || playerFileInventoryAllowed(it) {
						if !writeItem(it) {
							return 0
						}
					}
				}
			} else {
				for row := 0; row < 20; row++ {
					for col := 0; col < 4; col++ {
						cell := &uiInventoryGrid()[row+21*col]
						for i := 0; i < int(cell.Count); i++ {
							it := controlEquippedByCode(u, cell.Codes[i])
							if it == nil || !writeItem(it) {
								return 0
							}
						}
					}
				}
			}
			equipped := byte(0)
			for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
				if uint32(it.ObjFlags)&0x100 != 0 {
					equipped++
				}
			}
			r.byte(equipped)
			for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
				if uint32(it.ObjFlags)&0x100 != 0 {
					r.raw(unsafe.Pointer(&it.ScriptIDVal), 4)
				}
			}
			selectedID := func(code uint32) uint32 {
				if code != 0 {
					for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
						if it.NetCode == code {
							return uint32(it.ScriptIDVal)
						}
					}
				}
				return 0
			}
			r.word(selectedID(uint32(sub_4678B0())))
			r.word(selectedID(uint32(sub_4678C0())))
		}
		tail := (*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](data).CurTraps))
		if version < 3 {
			*tail = 0
		} else {
			r.raw(unsafe.Pointer(tail), 1)
		}
		if r.read() && noxflags.HasGame(4096) {
			*tail = 0
		}
	}
	if noxflags.HasGame(4096) && !questEligibilityInventoryLimit(u) {
		return 0
	}
	gameplayReportInventoryLoaded(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)))
	return 1
}
