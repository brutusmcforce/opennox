package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/strman"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func bookAwardError(u *server.Object, key, file string) {
	Nox_xxx_netSendLineMessage_4D9EB0(u, GetServer().S().Strings().GetStringInFile(strman.ID(key), file))
}
func bookAwardPrivate(u *server.Object, key string) {
	gameplayTextPrivate(u, alloc.InternCString(key), 0)
}
func bookAwardClientMessage(id int32, success bool) {
	off := int32(216380)
	file := "plyrspel.c"
	if success {
		off = 217092
		file = "Ability.c"
	}
	key := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, uintptr(off+4*id))))
	Nox_xxx_printCentered_445490(GetServer().S().Strings().GetStringInFile(strman.ID(key), file))
}
func bookReportSpell(u *server.Object, id, notify, auto int32) uint32 {
	if uint32(u.ObjClass)&4 == 0 {
		return memmap.PtrWord(u.CObj())
	}
	p := u.UpdateDataPlayer().Player
	data := [4]byte{111, byte(id), byte(*equipmentWord(unsafe.Pointer(p), 3696+4*int(id))), byte(notify)}
	if auto != 0 {
		data[3] |= 128
	}
	return uint32(reliableEnqueue(int(p.PlayerInd), data[:], nil, 1, 1))
}
func bookReportAbility(u *server.Object, id, auto int32) uint32 {
	if uint32(u.ObjClass)&4 == 0 {
		return memmap.PtrWord(u.CObj())
	}
	p := u.UpdateDataPlayer().Player
	data := [3]byte{205, byte(id), byte(*equipmentWord(unsafe.Pointer(p), 3696+4*int(id)))}
	if auto != 0 {
		data[2] |= 128
	}
	return uint32(reliableEnqueue(int(p.PlayerInd), data[:], nil, 1, 1))
}
func bookReportGuide(u *server.Object, id, notify, auto int32) uint32 {
	if uint32(u.ObjClass)&4 == 0 {
		return memmap.PtrWord(u.CObj())
	}
	data := [3]byte{209, byte(id), byte(notify)}
	if auto != 0 {
		data[2] |= 128
	}
	return uint32(reliableEnqueue(int(u.UpdateDataPlayer().Player.PlayerInd), data[:], nil, 1, 1))
}
func bookAwardAbility(u *server.Object, id, notify int32) int {
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	if id <= 0 || id >= 6 {
		bookAwardError(u, "AwardAbilityError", "Ability.c")
		return 0
	}
	p := u.UpdateDataPlayer().Player
	if p.SpellLvl[id] != 0 {
		bookAwardPrivate(u, "use.c:HadAbility")
		return 0
	}
	p.SpellLvl[id] = 5
	Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(p.Prot4636, int(id), int(p.SpellLvl[id]))
	bookReportAbility(u, id, notify)
	core := GetServer().S()
	if noxflags.HasGame(4096) {
		controlRewardNotify(u, 2, u, byte(id))
		if !core.Players.CheckXxx(u) {
			for it := core.Players.FirstUnit(); it != nil; it = core.Players.NextUnit(it) {
				if it != u {
					controlRewardNotify(it, 2, u, byte(id))
				}
			}
		}
	}
	return 1
}
func bookAwardGuide(u *server.Object, id, notify int32) int {
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	if id <= 0 || id >= 41 {
		bookAwardError(u, "AwardGuideError", "PlyrGide.c")
		return 0
	}
	p := u.UpdateDataPlayer().Player
	if p.BeastScrollLvl[id] != 0 {
		return 0
	}
	p.BeastScrollLvl[id] = 1
	Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(p.Prot4640, int(id), int(p.BeastScrollLvl[id]))
	core := GetServer().S()
	if notify != 0 {
		core.Audio.EventObj(227, u, 0, 0)
		controlRewardNotify(u, 1, u, byte(id))
	}
	for off := uintptr(216292); ; off += 4 {
		family := memmap.LoadPtr(0x587000, off)
		if family == nil {
			break
		}
		if *(*int32)(family) != id {
			continue
		}
		for entry := unsafe.Add(family, 4); ; entry = unsafe.Add(entry, 4) {
			child := *(*int32)(entry)
			if child == 0 {
				break
			}
			p.BeastScrollLvl[child] = 1
			Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(p.Prot4640, int(child), int(p.BeastScrollLvl[child]))
		}
	}
	if notify != 0 {
		for pl := core.Players.First(); pl != nil; pl = core.Players.Next(pl) {
			if it := pl.PlayerUnit; it != nil && it != u {
				controlRewardNotify(it, 1, u, byte(id))
			}
		}
	}
	bookReportGuide(u, id, notify, 0)
	return 1
}
func bookQuestSingleLevel(id int32) bool {
	switch id {
	case 19, 34, 45, 46, 47, 48, 49, 117, 118, 119, 120, 121, 122, 123, 124, 125, 134:
		return true
	}
	return false
}
func bookAwardSpell(u *server.Object, id, notify, auto, force int32) int {
	if uint32(u.ObjClass)&4 == 0 {
		return 0
	}
	if id <= 0 || id >= 137 {
		bookAwardError(u, "AwardSpellError", "plyrspel.c")
		return 0
	}
	p := u.UpdateDataPlayer().Player
	if noxflags.HasGame(6144) && p.SpellLvl[id] == 3 || p.SpellLvl[id] == 5 {
		bookAwardError(u, "MaxSpellLevel", "plyrspel.c")
		return 0
	}
	if noxflags.HasGame(4096) && bookQuestSingleLevel(id) && p.SpellLvl[id] != 0 {
		bookAwardError(u, "MaxSpellLevel", "plyrspel.c")
		return 0
	}
	p.SpellLvl[id]++
	if p.SpellLvl[id] > 5 {
		p.SpellLvl[id] = 5
	}
	if noxflags.HasGame(4096) && p.SpellLvl[id] > 3 {
		p.SpellLvl[id] = 3
	}
	if force != 0 {
		p.SpellLvl[id] = uint32(force)
	}
	Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(p.Prot4636, int(id), int(p.SpellLvl[id]))
	core := GetServer().S()
	sf := uint32(core.Spells.Flags(spell.ID(id)))
	family := uint32(0)
	if sf&0x1000 != 0 {
		family = 0x2000
	} else if sf&0x4000 != 0 {
		family = 0x8000
	} else if sf&0x10000 != 0 {
		family = 0x20000
	}
	if family != 0 {
		for i := int32(1); i < 137; i++ {
			if uint32(core.Spells.Flags(spell.ID(i)))&family == 0 || !core.Spells.DefByInd(spell.ID(i)).IsValid() {
				continue
			}
			if force != 0 {
				p.SpellLvl[i] = uint32(force)
			} else {
				p.SpellLvl[i]++
			}
			if p.SpellLvl[i] > 5 {
				p.SpellLvl[i] = 5
			}
			// Preserve the original-ID cap and bookkeeping in historical family awards.
			if noxflags.HasGame(4096) && p.SpellLvl[id] > 3 {
				p.SpellLvl[id] = 3
			}
			Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(p.Prot4636, int(id), int(p.SpellLvl[i]))
		}
	}
	if notify != 0 {
		core.Audio.EventObj(226, u, 0, 0)
		show := !(noxflags.HasGame(2048) && (id == 34 || uint32(core.Spells.Flags(spell.ID(id)))&0x15000 != 0))
		if (!noxflags.HasGame(4096) || p.Field4792 != 0) && show {
			controlRewardNotify(u, 0, u, byte(id))
			if !core.Players.CheckXxx(u) {
				for pl := core.Players.First(); pl != nil; pl = core.Players.Next(pl) {
					if it := pl.PlayerUnit; it != nil && it != u {
						controlRewardNotify(it, 0, u, byte(id))
					}
				}
			}
		}
	}
	if noxflags.HasGame(2048) && notify == 1 && auto == 1 {
		if session := u.UpdateDataPlayer().Trade70; session != nil {
			shopExit(session)
		}
	}
	bookReportSpell(u, id, notify, auto)
	return 1
}
