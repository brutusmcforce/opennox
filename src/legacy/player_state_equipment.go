package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func playerStateEquip(command byte, code, mask uint32, mods *[4]byte) unsafe.Pointer {
	s := GetServer().S()
	pl := s.Players.ByID(int(code))
	if pl == nil {
		return nil
	}
	slots := pl.Armor[:]
	if command == 80 || command == 81 {
		pl.WeaponEquip |= mask
		slots = pl.Weapon[:]
	} else {
		pl.ArmorEquip |= mask
	}
	for i := range slots {
		slot := &slots[i]
		if slot.Field0 != 0 {
			continue
		}
		slot.Field0 = mask
		var result unsafe.Pointer
		for j, id := range mods {
			result = unsafe.Pointer(s.Modif.Nox_xxx_modifGetDescById413330(int(id)))
			slot.Field4[j] = result
		}
		return result
	}
	return pl.C()
}
func playerStateUnequip(command byte, code, mask uint32) unsafe.Pointer {
	pl := GetServer().S().Players.ByID(int(code))
	if pl == nil {
		return nil
	}
	slots := pl.Armor[:]
	if command == 84 {
		pl.WeaponEquip &^= mask
		slots = pl.Weapon[:]
	} else {
		pl.ArmorEquip &^= mask
	}
	for i := range slots {
		if slots[i].Field0 == mask {
			slots[i].Field0 = 0
			break
		}
	}
	return pl.C()
}
func playerStateRespawn(pl *server.Player, mask byte) {
	if pl == nil {
		return
	}
	if !noxflags.HasGame(1) {
		pl.WeaponEquip = 0
		pl.ArmorEquip = 0
	}
	// The sixth word in each display slot belongs to another owner.
	for i := range pl.Weapon {
		pl.Weapon[i].Field0 = 0
		pl.Weapon[i].Field4 = [4]unsafe.Pointer{}
	}
	for i := range pl.Armor {
		pl.Armor[i].Field0 = 0
		pl.Armor[i].Field4 = [4]unsafe.Pointer{}
	}
	s := GetServer().S()
	byName := func(name string) *server.ModifierEff {
		return s.Modif.Nox_xxx_modifGetDescById413330(s.Modif.Nox_xxx_modifGetIdByName413290(name))
	}
	color := byName("UserColor1")
	if color == nil {
		return
	}
	id := func(m *server.ModifierEff) byte { return byte(m.Index()) }
	colorID := uint32(color.Index())
	colorAt := func(offset uintptr) byte {
		index := *(*byte)(unsafe.Add(pl.C(), offset))
		return id(s.Modif.Nox_xxx_modifGetDescById413330(int(colorID + uint32(index))))
	}
	equip := func(command byte, bit uint32, mods [4]byte) { playerStateEquip(command, pl.NetCodeVal, bit, &mods) }
	class := *(*byte)(unsafe.Add(pl.C(), 2251))
	if class != 0 || noxflags.HasGame(2048) {
		mods := [4]byte{255, colorAt(2269), colorAt(2270), 255}
		if mask&1 != 0 {
			equip(82, 1024, mods)
		}
	}
	mods := [4]byte{255, colorAt(2268), 255, 255}
	if mask&2 != 0 {
		equip(82, 4, mods)
	}
	mods = [4]byte{colorAt(2272), colorAt(2271), 255, 255}
	if mask&4 != 0 {
		equip(82, 1, mods)
	}
	mods = [4]byte{255, 255, 255, 255}
	switch class {
	case 1:
		if noxflags.HasGame(2048) {
			if mask&8 != 0 {
				mods[0] = id(byName("ArmorQuality1"))
				equip(80, 0x8000, mods)
			}
		} else if noxflags.HasGame(4096) {
			mods[2] = id(byName("Replenishment1"))
			equip(80, 0x10000, mods)
		} else if mask&16 != 0 {
			equip(79, 0x4000, mods)
		}
	case 0:
		if noxflags.HasGame(2048) {
			if mask&32 != 0 {
				mods[0] = id(byName("ArmorQuality1"))
				// The original overwrites the quality byte with the material ID.
				mods[0] = id(byName("Material1"))
				equip(80, 256, mods)
			}
		} else if noxflags.HasGame(4096) {
			equip(80, 256, mods)
		} else {
			if mask&64 != 0 {
				equip(80, 512, mods)
			}
			if mask&128 != 0 {
				equip(79, 0x1000000, mods)
			}
		}
	case 2:
		if noxflags.HasGame(2048) {
			if mask&8 != 0 {
				mods[0] = id(byName("ArmorQuality1"))
				equip(80, 0x8000, mods)
			}
		} else if noxflags.HasGame(4096) {
			equip(80, 4, mods)
		}
	}
}
