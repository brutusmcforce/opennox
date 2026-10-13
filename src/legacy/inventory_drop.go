package legacy

import (
	"encoding/binary"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"

	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

func inventoryCache(off uintptr, name string) uint32 {
	p := memmap.PtrUint32(0x5D4594, off)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID(name))
	}
	return *p
}
func inventorySound(id int, u *server.Object, a, b int) {
	nox_xxx_aud_501960(int32(int32(id)), asObjectC(u), int(int32(a)), int32(int32(b)))
}
func inventoryMessage(kind int, u *server.Object, value uint32) {
	var data [10]byte
	binary.LittleEndian.PutUint32(data[2:], u.NetCode)
	binary.LittleEndian.PutUint32(data[6:], value)
	gameplayTextInformationAll(int(int32(kind)), unsafe.Pointer(&data[0]))
}
func inventoryDefaultDrop(u, it *server.Object, pos *types.Pointf) int {
	if it.InvHolder != u {
		return 0
	}
	if u.ObjClass&4 != 0 && inventoryDroppable(it) && equipmentDropPolicy(it, 1) != 0 {
		if u.ObjFlags&0x8020 == 0 {
			gameplayTextPrivate(u, (*byte)(unsafe.Pointer(internCStr("drop.c:CantDropThat"))), 0)
			inventorySound(925, u, 2, int(u.NetCode))
		}
		return 0
	}
	inventoryRemove(u, it)
	GetServer().CreateObjectAt(it, nil, *pos)
	if it.ObjClass&0x1000000 != 0 && nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) == 4 {
		for ammo := u.InvFirstItem; ammo != nil; ammo = ammo.InvNextItem {
			if ammo.ObjClass&0x1000000 != 0 && nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(ammo)) == 2 && ammo.ObjFlags&0x100 == 0 && ammo.UseDataAmmo().Infinite != 0 {
				inventoryRemove(u, ammo)
				GetServer().DelayedDelete(ammo)
				break
			}
		}
	}
	if it.ObjClass&0x10000000 != 0 {
		team := it.TeamVal.ID
		value := objectiveFlagID(it)
		inventoryMessage(7, u, uint32(value))
		playerStateMark(it, 1)
		it.UpdateDataFlag().DropFrame = GetServer().S().Frame()
		matchRosterFlagState(byte(team), 2, byte(value), 0)
	}
	glyph := inventoryCache(1568252, "Glyph")
	if !noxflags.HasGame(2048|4096) && it.ObjFlags&0x80000 == 0 && it.ObjClass&0x10000000 == 0 && uint32(it.TypeInd) != glyph {
		Nox_xxx_unitSetDecayTime_511660(it, int(10*GetServer().S().TickRate()))
	}
	Nox_xxx_unitRaise_4E46F0(it, 0)
	if it.ObjClass&64 != 0 {
		inventorySound(821, it, 0, 0)
	}
	if it.ObjClass&2 != 0 {
		layout.As[server.MonsterUpdateData](it.UpdateData).AIAction340 = 15
		*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](it.UpdateData).StatusFlags)) |= 0x100
	}
	if *memmap.PtrUint32(0x5D4594, 1568256) == 0 {
		*memmap.PtrUint32(0x5D4594, 1568256) = uint32(GetServer().S().Types.IndByID("Torch"))
		*memmap.PtrUint32(0x5D4594, 1568244) = uint32(GetServer().S().Types.IndByID("Lantern"))
	}
	if uint32(it.TypeInd) == *memmap.PtrUint32(0x5D4594, 1568256) || uint32(it.TypeInd) == *memmap.PtrUint32(0x5D4594, 1568244) {
		Nox_xxx_spellBuffOff_4FF5B0(u, 15)
	}
	return 1
}
func inventoryTrapDrop(u, it *server.Object, pos *types.Pointf) int {
	if worldTileWater(*pos) {
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	GetServer().S().ObjSetOwner(u, it)
	return 1
}
func inventoryGlyphDrop(u, it *server.Object, pos *types.Pointf) int {
	if inventoryTrapDrop(u, it, pos) == 0 {
		return 0
	}
	it.InitDataGlyph().SpellArg.Pos = *pos
	dir := types.Pointf{X: u.PosVec.X - pos.X, Y: u.PosVec.Y - pos.Y}
	angle := uint16(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&dir))))))
	it.Direction1, it.Direction2 = server.Dir16(angle), server.Dir16(angle)
	inventorySound(825, it, 0, 0)
	return 1
}
func inventoryDrop(u, it *server.Object, pos *types.Pointf) int {
	if it == nil {
		return 0
	}
	if noxflags.HasGame(0x2000) && !noxflags.HasGame(4096) && it.ObjClass&0x3001010 != 0 {
		it.ObjFlags |= 0x40
		motionDeactivate(it)
	}
	if it.Drop.Ptr != nil {
		if fn := inventoryNativeDrops[it.Drop.Ptr]; fn != nil {
			return fn(u, it, pos)
		}
		panic("unregistered inventory drop callback")
	}
	return inventoryDefaultDrop(u, it, pos)
}
func inventoryPotionDrop(u, it *server.Object, pos *types.Pointf) int {
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	inventorySound(833, it, 0, 0)
	if !noxflags.HasGame(2048 | 4096) {
		Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
	}
	return 1
}
func inventoryFoodSound(off uintptr, u, it *server.Object) {
	for ; *memmap.PtrUint16(0x587000, off+6) != 0; off += 8 {
		if uint32(it.ObjSubClass)&*memmap.PtrUint32(0x587000, off) != 0 || it.Material&*memmap.PtrUint16(0x587000, off+4) != 0 {
			inventorySound(int(*memmap.PtrUint16(0x587000, off+6)), u, 0, 0)
			return
		}
	}
}
func inventoryFoodDrop(u, it *server.Object, pos *types.Pointf) int {
	if u == nil || it == nil || pos == nil {
		return 0
	}
	rv := inventoryDefaultDrop(u, it, pos)
	if rv != 0 {
		if !noxflags.HasGame(2048) {
			Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
		}
		inventoryFoodSound(205704, u, it)
	}
	return rv
}
func inventoryEquipmentDrop(u, it *server.Object, pos *types.Pointf, armor bool) int {
	if inventoryDefaultDrop(u, it, pos) != 1 {
		return 0
	}
	if armor {
		equipmentArmorDropSound(it)
	} else {
		equipmentDropSound(it)
	}
	if !noxflags.HasGame(2048|4096) && int32(serverConfigFlagsQuery(int32(2))) != 0 {
		Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
	}
	return 1
}
