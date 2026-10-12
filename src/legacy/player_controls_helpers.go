package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func controlByte(p unsafe.Pointer, n int) *byte          { return (*byte)(unsafe.Add(p, n)) }
func controlHalf(p unsafe.Pointer, n int) *uint16        { return (*uint16)(unsafe.Add(p, n)) }
func controlPtr(p unsafe.Pointer, n int) *unsafe.Pointer { return (*unsafe.Pointer)(unsafe.Add(p, n)) }
func controlObject(p unsafe.Pointer, n int) *server.Object {
	return (*server.Object)(*controlPtr(p, n))
}
func controlPlayer(u *server.Object) unsafe.Pointer { return *controlPtr(u.UpdateData, 276) }
func controlFlags(mask uint32) bool                 { return noxflags.HasGame(noxflags.GameFlag(mask)) }
func controlRaw(u *server.Object) uint32            { return memmap.PtrWord(u.CObj()) }
func controlRemoveGlyphs(u *server.Object) int32 {
	typ := stateType(1565600, "Glyph")
	for it := GetServer().S().Objs.List; it != nil; it = it.ObjNext {
		if it.HasOwner(u) && uint32(it.TypeInd) == typ && it.ObjFlags&0x20 == 0 {
			visibilityFXPoint(129, it.PosVec)
			GetServer().DelayedDelete(it)
		}
	}
	return 0
}
func controlRemoveCreatures(u *server.Object) {
	for it := u.Field129; it != nil; {
		next := it.Field128
		if server.Nox_xxx_creatureIsMonitored_500CC0(u, it) {
			for item := it.InvFirstItem; item != nil; {
				n := item.InvNextItem
				GetServer().DelayedDelete(item)
				item = n
			}
			visibilityFXPoint(129, it.PosVec)
			GetServer().DelayedDelete(it)
		}
		it = next
	}
}
func controlFindBall() *server.Object {
	if dword_5d4594_1565616 == 0 {
		dword_5d4594_1565616 = uint32(GetServer().S().Types.IndByID("GameBall"))
	}
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if uint32(u.TypeInd) == uint32(dword_5d4594_1565616) {
			return u
		}
	}
	return nil
}
func controlNextObserver(pl unsafe.Pointer) *server.Object {
	if dword_5d4594_1565616 == 0 {
		dword_5d4594_1565616 = uint32(GetServer().S().Types.IndByID("GameBall"))
	}
	players := &GetServer().S().Players
	current := controlObject(pl, 3628)
	var u *server.Object
	if current != nil && current.ObjClass&4 != 0 {
		u = players.NextUnit(current)
	} else if current == nil && controlFlags(64) {
		u = controlFindBall()
		if u == nil {
			u = players.FirstUnit()
		}
	} else {
		u = players.FirstUnit()
	}
	eligible := func(u *server.Object) bool {
		p := players.ByID(int(u.NetCode))
		return u.ObjFlags&0x20 == 0 && *controlByte(unsafe.Pointer(p), int(layoutServerPlayer.Off(3680)) /* Field3680 */)&1 == 0
	}
	if u != nil && u.ObjClass&4 != 0 {
		for u != nil && !eligible(u) {
			u = players.NextUnit(u)
		}
	}
	if u != nil {
		return u
	}
	if u = controlFindBall(); u != nil {
		return u
	}
	for u = players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if eligible(u) {
			break
		}
	}
	return u
}
func controlSlave(u *server.Object, next bool) *server.Object {
	if u == nil {
		return nil
	}
	if next {
		if u.ObjOwner == nil {
			return nil
		}
		u = u.Field128
	} else {
		u = u.Field129
	}
	for ; u != nil; u = u.Field128 {
		if u.ObjClass&2 != 0 && *controlByte(u.UpdateData, 1440)&0x80 != 0 {
			return u
		}
	}
	return nil
}
func controlObserverSlave(pl unsafe.Pointer) *server.Object {
	u := controlObject(pl, 3628)
	if u != nil {
		u = controlSlave(u, true)
	} else {
		u = controlSlave(controlObject(pl, 2056), false)
	}
	for ; u != nil; u = controlSlave(u, true) {
		if u.ObjFlags&0x8020 == 0 {
			return u
		}
	}
	for u = controlSlave(controlObject(pl, 2056), false); u != nil; u = controlSlave(u, true) {
		if u.ObjFlags&0x8020 == 0 {
			return u
		}
	}
	return nil
}
func controlRemoveChildren(u *server.Object) {
	if u == nil {
		return
	}
	for it := u.Field129; it != nil; {
		next := it.Field128
		it.ObjOwner = nil
		it.Field128 = nil
		it = next
	}
	u.Field129 = nil
}
func controlTransferChildren(u *server.Object) {
	if u == nil {
		return
	}
	for it := u.Field129; it != nil; {
		next := it.Field128
		GetServer().S().ObjSetOwner(u.ObjOwner, it)
		it = next
	}
}
func controlBoltDamage(level int32, def unsafe.Pointer) float64 {
	typ := stateType(1568264, "ArcherBolt")
	base := float64(*controlHalf(def, 72))
	if controlFlags(2048) && *equipmentWord(def, 4) == typ {
		base = float64(nox_xxx_gamedataGetFloat_419D40(internCStr("BoltSoloDamageMin")))
	}
	return float64(level-int32(*controlHalf(def, 60)))*float64(*(*float32)(unsafe.Add(def, 64))) + base
}
func controlRespawnFlags() int8 {
	s := GetServer().S()
	var mask byte
	for i, bit := range []uint32{0x400, 4, 1, 0x8000, 0x4000, 0x100, 0x200, 0x1000000} {
		var typ int
		if i == 3 || i == 5 || i == 6 {
			typ = int(s.Weapons.Sub_415840(bit))
		} else {
			typ = int(s.Armor.Sub_415CD0(bit))
		}
		if s.Types.ByInd(typ).Allowed() {
			mask |= 1 << i
		}
	}
	return int8(mask)
}
func controlGlyphCount(u *server.Object) int32 {
	typ := stateType(1568268, "Glyph")
	var n int32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if uint32(it.TypeInd) == typ {
			n++
		}
	}
	return n
}
func controlEquippedByCode(u *server.Object, code uint32) *server.Object {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.NetCode == code {
			return it
		}
	}
	return nil
}
func controlClearWaypoints(u *server.Object) {
	for i := 0; i < 3; i++ {
		ptr := controlPtr(u.UpdateData, 168+4*i)
		if *ptr != nil {
			GetServer().DelayedDelete((*server.Object)(*ptr))
		}
		*ptr = nil
	}
	*controlByte(u.UpdateData, 180) = 0
	*controlByte(u.UpdateData, 181) = 0
}
func controlSetWaypoint(u *server.Object, x, y uint32) {
	data := u.UpdateData
	if *controlByte(controlPlayer(u), 3680)&3 != 0 {
		return
	}
	ptr := controlPtr(data, 168+4*int(*controlByte(data, 180)))
	pos := *(*types.Pointf)(unsafe.Pointer(&[2]uint32{x, y}))
	if *ptr != nil {
		Nox_xxx_unitMove_4E7010((*server.Object)(*ptr), pos)
	} else {
		it := GetServer().S().NewObjectByTypeID("PlayerWaypoint")
		*ptr = it.CObj()
		GetServer().CreateObjectAt(it, u, pos)
	}
}
func controlConfusedDirection(u *server.Object) int32 {
	n := int32((GetServer().S().Frame() + u.NetCode) % 40)
	if n > 20 {
		n = 40 - n
	}
	v := (int32(u.BuffsPower[3])+3)*(n-10) + int32(int16(u.Direction2))
	for v < 0 {
		v += 256
	}
	for v >= 256 {
		v -= 256
	}
	return v
}
func controlStartEligible(u *server.Object, team int32) bool {
	return u.ObjFlags&0x1000000 != 0 && (team == 0 || !u.TeamVal.Has() || Nox_xxx_teamCompare2_419180(&u.TeamVal, server.TeamID(byte(team))) != 0)
}
func controlSubStamina(u *server.Object, amount int32) int32 {
	if u.ObjClass&6 == 0 {
		return 1
	}
	off := 1128
	if u.ObjClass&4 != 0 {
		off = 91
	}
	p := controlByte(u.UpdateData, off)
	if int32(*p) < amount {
		return 0
	}
	*p -= byte(amount)
	if off == 91 {
		gameplayReportStamina(int(*controlByte(controlPlayer(u), 2064)), u)
	}
	return 1
}
func controlAdjustStamina(u *server.Object, amount int8) {
	if u.ObjClass&4 != 0 {
		*controlByte(u.UpdateData, 91) -= byte(amount)
		gameplayReportStamina(int(*controlByte(controlPlayer(u), 2064)), u)
	}
}
func controlWeaponStamina(mask uint32) int32 {
	for _, v := range [][2]uint32{{0x200, 70}, {0x4000, 100}, {0x800, 50}, {0x100, 45}, {0x1000, 75}, {0x2000, 100}, {0x7ff8000, 45}, {0x400, 75}} {
		if mask&v[0] != 0 {
			return int32(v[1])
		}
	}
	return 10
}
func controlHasWaypoint(u *server.Object) bool {
	return *controlPtr(u.UpdateData, 168+4*int(*controlByte(u.UpdateData, 181))) != nil
}
func controlCanMove(u *server.Object) bool {
	d := u.UpdateData
	if u.Buffs&(1<<25|1<<5) != 0 || controlFlags(4096) && *equipmentWord(d, 280) != 0 {
		return false
	}
	if *controlByte(d, 88) == 1 {
		it := controlObject(d, 104)
		if it != nil && it.ObjClass&0x1000000 != 0 && it.ObjSubClass&8 != 0 {
			return false
		}
	}
	return true
}
func controlCanAttack(u *server.Object) bool {
	return u.Buffs&(1<<25) == 0 && *controlByte(u.UpdateData, 88) != 23
}
func controlAimsAtEnemy(u *server.Object) bool {
	if u == nil {
		return false
	}
	target := controlObject(u.UpdateData, 288)
	return target == nil || GetServer().S().IsEnemyTo(u, target) || controlFlags(4096)
}
func controlWeaponAnimation(mask uint32) int32 {
	for i := uint(2); i < 27; i++ {
		if mask&(1<<i) != 0 {
			return int32(*memmap.PtrUint32(0x587000, 215824+uintptr(4*i)))
		}
	}
	return 0
}
func controlActionState(u *server.Object) int32 {
	d := u.UpdateData
	pl := controlPlayer(u)
	weapon := *equipmentWord(pl, 4)
	switch *controlByte(d, 88) {
	case 0:
		return 4
	case 1, 14, 22:
		if GetServer().S().Abils.IsActive(u, server.Ability(2)) && GetServer().S().Abils.IsActiveVal(u, server.Ability(2)) {
			return 46
		}
		if GetServer().S().Abils.IsActive(u, server.Ability(1)) {
			return 45
		}
		if weapon&0x47f0000 != 0 {
			return int32((^byte(*equipmentWord(controlObject(d, 104).UseData.Ptr, 96)) & 2) | 29)
		}
		if (weapon == 0 || weapon == 1) && *controlByte(pl, 8) != 0 {
			return int32(*controlByte(pl, 8))
		}
		return controlWeaponAnimation(weapon)
	case 2, 10:
		return 21
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 6
	case 12:
		return 3
	case 13:
		if weapon&0x400 != 0 {
			return 38
		}
		return 0
	case 15, 16, 17:
		return 40
	case 18:
		return 48
	case 19:
		return 49
	case 20:
		return 47
	case 21:
		return 30
	case 23:
		return 50
	case 24:
		return 19
	case 25:
		return 20
	case 26:
		return 15
	case 27, 28, 29:
		return 16
	case 30:
		return 52
	case 32:
		return 54
	}
	return 0
}
func controlBotState(u *server.Object) int8 {
	d := *controlPtr(u.UpdateData, 292)
	index := int8(*controlByte(d, 544))
	if index == -1 {
		return 13
	}
	switch *equipmentWord(d, 24*(int(index)+23)) {
	case 7, 8, 10, 13, 29:
		if *equipmentWord(d, 1440)&0x4000 != 0 {
			return 5
		}
		return 0
	case 9:
		return 0
	case 16, 17:
		return 1
	case 18, 19, 20:
		return 2
	case 21, 23:
		return 16
	case 22:
		return 17
	case 24:
		return 5
	case 30:
		return 3
	case 31:
		return 4
	}
	return 13
}
