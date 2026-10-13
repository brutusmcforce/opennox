package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func gameplayReportModifiers(unit *server.Object) (out [4]byte, any bool) {
	for i := range out {
		out[i] = 255
		if p := itemModifiers(unit)[i]; p != nil {
			out[i] = byte(p.Ind4)
			any = true
		}
	}
	return
}
func gameplayReportEquipment(to int, unit *server.Object) int {
	class := uint32(unit.ObjClass)
	weapon := class&0x11001000 != 0
	if !weapon && class&0x2000000 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	mods, modified := gameplayReportModifiers(unit)
	opcode := byte(79)
	if weapon {
		opcode = 80
	}
	if modified {
		opcode = 82
		if weapon {
			opcode = 81
		}
	}
	holder := unit.InvHolder
	code := uint16(holder.NetCode)
	if uint32(holder.ObjClass)&4 != 0 {
		code |= 0x8000
	}
	b := []byte{opcode, byte(code), byte(code >> 8), 0, 0, 0, 0}
	var flags uint32
	if weapon {
		flags = GetServer().S().Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(unit)
	} else {
		flags = GetServer().S().Armor.Nox_xxx_unitArmorInventoryEquipFlags_415C70(unit)
	}
	binary.LittleEndian.PutUint32(b[3:], flags)
	if modified {
		b = append(b, mods[:]...)
	}
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportDequipFlags(to int, unit *server.Object) int {
	class := uint32(unit.ObjClass)
	weapon := class&0x11001000 != 0
	if !weapon && class&0x2000000 == 0 {
		return gameplayReportPtr(unit.CObj())
	}
	opcode := byte(83)
	var flags uint32
	if weapon {
		opcode = 84
		flags = GetServer().S().Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(unit)
	} else {
		flags = GetServer().S().Armor.Nox_xxx_unitArmorInventoryEquipFlags_415C70(unit)
	}
	code := uint16(unit.InvHolder.NetCode)
	b := []byte{opcode, byte(code), byte(code >> 8), 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[3:], flags)
	return gameplayReportSend(to, b, true, 0)
}
func gameplayReportEquip(to int, unit *server.Object, notify int) int {
	gameplayReportSend(to, gameplayReportID(96, unit, 3), true, 0)
	if notify != 0 {
		return gameplayReportEquipment(255, unit)
	}
	return 0
}
func gameplayReportPickup(to int, unit *server.Object) int {
	if uint32(unit.ObjClass)&0x13001000 != 0 {
		return gameplayReportModifiablePickup(to, unit)
	}
	b := gameplayReportID(75, unit, 5)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	gameplayReportSend(to, b, true, 0)
	return gameplayReportItemHealth(to, unit)
}
func gameplayReportModifiablePickup(to int, unit *server.Object) int {
	b := gameplayReportID(76, unit, 9)
	binary.LittleEndian.PutUint16(b[3:], unit.TypeInd)
	mods, _ := gameplayReportModifiers(unit)
	copy(b[5:], mods[:])
	gameplayReportSend(to, b, true, 0)
	return gameplayReportItemHealth(to, unit)
}
func gameplayReportTeamBase(to int, unit *server.Object) int {
	cache := memmap.PtrUint32(0x5D4594, 1556320)
	if *cache == 0 {
		*cache = uint32(GetServer().S().Types.IndByID("TeamBase"))
	}
	if uint32(unit.ObjClass)&0x13001000 == 0 && uint32(unit.TypeInd) != *cache {
		return int(*cache)
	}
	b := gameplayReportID(103, unit, 7)
	mods, _ := gameplayReportModifiers(unit)
	copy(b[3:], mods[:])
	return gameplayReportSend(to, b, false, 1)
}
func gameplayReportNPC(to int, unit *server.Object) int {
	code := uint16(unit.NetCode)
	if unit.Poison540 != 0 {
		code |= 0x8000
	}
	b := make([]byte, 21)
	b[0] = 105
	binary.LittleEndian.PutUint16(b[1:], code)
	copy(b[3:], unsafe.Slice(&layout.As[server.MonsterUpdateData](unit.UpdateData).Color[0].R, 18))
	result := gameplayReportSend(to, b, true, 1)
	for it := unit.InvFirstItem; it != nil; it = it.InvNextItem {
		if uint32(it.ObjFlags)&0x100 != 0 {
			result = gameplayReportEquipment(to, it)
		}
	}
	return result
}
