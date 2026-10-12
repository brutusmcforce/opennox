package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// Bot update data remains C-owned because the object lifecycle frees it in C.
// Morphing temporarily exchanges player and monster records without copying them.
func controlBotCreate(u *server.Object) uint32 {
	d := u.UpdateData
	ptr := controlPtr(d, 292)
	result := controlRaw(u)
	if *ptr == nil {
		*ptr = legacyCalloc(1, 2200)
		result = memmap.PtrWord(*ptr)
	}
	b := *ptr
	if b == nil {
		return result
	}
	clear(unsafe.Slice((*byte)(b), 2200))
	*controlPtr(b, 2180) = d
	*controlPtr(b, 484) = unsafe.Pointer(monsterDefinitionByType(uint32(GetServer().S().Types.IndByID("NPC"))))
	for _, v := range [][2]uint32{{1336, 1048576000}, {1344, 1061997773}, {1440, 186376}, {552, 5}, {1360, 38}, {1308, 1056964608}, {1304, 1062501089}, {1312, 1125515264}, {1316, 1106247680}, {1320, 1065353216}, {1328, 1056964608}, {1352, 1065353216}, {2040, 3}, {2096, 0xffffffff}, {2100, 0xffffffff}, {0, 0xdeadface}} {
		*equipmentWord(b, int(v[0])) = v[1]
	}
	for _, v := range [][2]int{{1340, 1}, {1348, 1}, {1332, 255}, {1324, 30}} {
		*controlByte(b, v[0]) = byte(v[1])
	}
	for i := 1228; i <= 1300; i += 8 {
		*equipmentWord(b, i) = 0xffffffff
	}
	*equipmentWord(b, 376) = uint32(int32(int16(u.Direction1)))
	*equipmentWord(b, 380) = *equipmentWord(u.CObj(), int(layoutServerObject.Off(56)) /* PosVec.X */)
	*equipmentWord(b, 384) = *equipmentWord(u.CObj(), int(layoutServerObject.Off(60)) /* PosVec.Y */)
	class := *controlByte(controlPlayer(u), 2251)
	if class == 0 {
		return 0
	}
	if class > 2 {
		return uint32(class) - 2
	}
	fps := GetServer().S().TickRate()
	f := uint32(fps)
	*equipmentWord(b, 1356) = 1112014848
	*equipmentWord(b, 1640) = 0x8000000
	*equipmentWord(b, 1440) |= 0x20
	*controlHalf(b, 1450) = uint16(f >> 1)
	*controlHalf(b, 1464) = uint16(3 * f)
	*controlHalf(b, 1466) = uint16(30 * f)
	*controlHalf(b, 1474) = uint16(2 * f)
	*controlHalf(b, 1482) = uint16(6 * f)
	*equipmentWord(b, 1504) = 0x80000000
	if class == 2 {
		for _, i := range []int{430} {
			*equipmentWord(b, 4*i) = 0x10000000
		}
		for _, i := range []int{432, 446} {
			*equipmentWord(b, 4*i) = 0x20000000
		}
		for _, i := range []int{401, 424, 456, 455, 464} {
			*equipmentWord(b, 4*i) = 0x40000000
		}
		*controlHalf(b, 1458) = uint16(6 * f)
		return f
	}
	for _, i := range []int{423, 408, 411} {
		*equipmentWord(b, 4*i) = 0x10000000
	}
	for _, i := range []int{384, 405} {
		*equipmentWord(b, 4*i) = 0x20000000
	}
	*controlHalf(b, 1458) = uint16(f)
	n := GetServer().S().Rand.Logic.IntClamp(0, 100)
	if n < 33 {
		*equipmentWord(b, 1596) = 0x40000000
	} else if n < 66 {
		*equipmentWord(b, 1552) = 0x40000000
	} else {
		*equipmentWord(b, 1660) = 0x40000000
		*equipmentWord(b, 1688) = 0x40000000
	}
	return 6 * f
}
func controlMorphFromPlayer(u *server.Object) int8 {
	if u.ObjClass&4 != 0 {
		u.ObjClass = u.ObjClass&^4 | 2
		u.UpdateData = *controlPtr(u.UpdateData, 292)
		u.ObjSubClass = 16
	}
	return int8(u.ObjClass)
}
func controlMorphToPlayer(u *server.Object) int8 {
	if u.ObjClass&2 != 0 {
		u.ObjClass = u.ObjClass&^2 | 4
		u.UpdateData = *controlPtr(u.UpdateData, 2180)
		u.ObjSubClass = 0
	}
	return int8(u.ObjClass)
}
func controlBotUpdate(u *server.Object) uint32 {
	d := u.UpdateData
	if *controlPtr(d, 292) == nil {
		controlBotCreate(u)
	}
	b := *controlPtr(d, 292)
	if b == nil {
		*controlPtr(u.CObj(), int(layoutServerObject.Off(744)) /* Update */) = updateIdentityKey(updateIDPlayer)
		return 0
	}
	result := controlRespawnBot(u)
	if result != 0 {
		return uint32(result)
	}
	*equipmentWord(b, 1440) |= 0x100
	controlMorphFromPlayer(u)
	Nox_xxx_unitUpdateMonster_50A5C0(u)
	controlMorphToPlayer(u)
	*controlByte(d, 88) = byte(controlBotState(u))
	*controlByte(d, 236) = *controlByte(b, 481)
	pl := *controlPtr(d, 276)
	*equipmentWord(pl, 3632) = *equipmentWord(u.CObj(), int(layoutServerObject.Off(56)) /* PosVec.X */)
	result = int32(*equipmentWord(u.CObj(), int(layoutServerObject.Off(60)) /* PosVec.Y */))
	*equipmentWord(pl, 3636) = uint32(result)
	return uint32(result)
}
