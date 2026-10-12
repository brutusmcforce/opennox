package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func extensionPlayerName(record unsafe.Pointer, teamOut *int32) int {
	if record == nil || uintptr(record) == 0xfffffffe {
		return 0
	}
	query := alloc.GoString((*byte)(unsafe.Add(record, 2)))
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		// The historical formatter narrows UTF-16 units to their low byte.
		var buf [28]byte
		n := 0
		for _, ch := range pl.NameFinal {
			if ch == 0 || byte(ch) == 0 {
				break
			}
			buf[n] = byte(ch)
			n++
		}
		if string(buf[:n]) != query {
			continue
		}
		team := teamRuntimeObject(int(pl.NetCodeVal))
		if team == nil {
			return 0
		}
		if teamOut != nil {
			*teamOut = int32(uintptr(team.C()))
		}
		*(*byte)(unsafe.Add(record, 1)) = byte(team.ID)
		*(*byte)(record) = byte(pl.PlayerClass())
		return 1
	}
	return 0
}

func extensionWeaponRoll(u *server.Object, direction int8) int {
	ud := u.UpdateData
	pl := u.UpdateDataPlayer().Player
	if *(*byte)(unsafe.Add(pl.C(), 3680))&3 != 0 || *(*byte)(unsafe.Add(ud, 88)) == 1 {
		return 0
	}
	current := *(**server.Object)(unsafe.Add(ud, 104))
	next := func(it *server.Object) *server.Object { return it.InvNextItem }
	// Without a current weapon, both directions search the forward inventory.
	start := u.InvFirstItem
	if current != nil {
		if direction == 0 {
			next = func(it *server.Object) *server.Object { return it.Field125 }
		}
		start = next(current)
	}
	for it := start; it != nil; it = next(it) {
		bits := equipmentWeaponBits(it)
		if bits == 0 || bits == 2 || !Nox_xxx_playerClassCanUseItem_57B3D0(it, pl.PlayerClass()) || !equipmentCheckStrength(u, it) {
			continue
		}
		if current != nil && equipmentTryDequip(u, current) == 0 {
			return 0
		}
		return equipmentTryEquip(u, it)
	}
	return 0
}

func extensionDropTrap(u *server.Object) int {
	if u == nil {
		return 0
	}
	pl := u.UpdateDataPlayer().Player
	pos := types.Pointf{X: *(*float32)(unsafe.Add(pl.C(), 3632)), Y: *(*float32)(unsafe.Add(pl.C(), 3636))}
	if *(*byte)(unsafe.Add(pl.C(), 3680))&3 != 0 || *(*byte)(unsafe.Add(u.UpdateData, 88)) == 1 {
		return 0
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if byte(uint32(it.ObjClass)>>16) == 17 {
			inventoryTargetDrop(u, it, &pos)
			return 1
		}
	}
	return 0
}
