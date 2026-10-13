package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/strman"
)

func quickbarText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "guispell.c")
}

func quickbarSlotTooltip(b *quickbarRecord, slot int) {
	b.Slots[slot].DrawData().SetTooltip(GetClient().Cli().Strings(), GoWString(nox_xxx_spellTitle_424930(int(b.Current[slot].ID))))
}

func quickbarSetSlot(b *quickbarRecord, id uint32, slot int) {
	b.Current[slot].ID = id
	if b.Slots[slot] != nil {
		quickbarSlotTooltip(b, slot)
	}
	if quickbarPlayer() != nil {
		bookSound(794)
	}
}

func quickbarBookSlot(kind uintptr, id uint32, slot int) uintptr {
	b := quickbarMain()
	quickbarSetSlot(b, id, slot)
	if kind != 2 {
		return kind
	}
	if id != 0 && bool(nox_xxx_spellHasFlags_424A50(int(id), 0x600)) {
		*quickbarFlagByte(&b.Current[slot]) = 1
		return 1
	}
	*quickbarFlagByte(&b.Current[slot]) = 0
	return uintptr(unsafe.Pointer(b))
}

func quickbarPoint(b *quickbarRecord, p image.Point) int {
	for i, w := range b.Slots {
		// The legacy point helper includes edges and treats a nil window as
		// a zero-sized rectangle at the origin.
		pos, size := image.Point{}, image.Point{}
		if w != nil {
			pos, size = uiWindowPosition(w), w.SizeVal
		}
		if p.X >= pos.X && p.X <= pos.X+size.X && p.Y >= pos.Y && p.Y <= pos.Y+size.Y {
			return i
		}
	}
	return -1
}

func quickbarPut(b *quickbarRecord, id uint32, p image.Point) int {
	i := quickbarPoint(b, p)
	if i < 0 {
		return 0
	}
	if b != quickbarAt(1047940) {
		quickbarSetSlot(b, id, i)
		return 1
	}
	message := "RestrictedTrapSpell"
	if GetServer().S().Spells.CanUseInTraps(spell.ID(id)) {
		duplicate := false
		for j := 0; j < 3; j++ {
			if b.Current[j].ID == id {
				duplicate = true
				break
			}
		}
		if !duplicate {
			quickbarSetSlot(b, id, i)
			return 1
		}
		message = "OneSpellPerTrap"
	}
	Nox_xxx_printCentered_445490(quickbarText(message))
	bookSound(925)
	return 0
}

func quickbarDrop(id uint32, flag byte, p image.Point, source *quickbarRecord) uint32 {
	oldID, oldFlag := uint16(137), byte(0)
	trap := quickbarAt(1047940)
	if uiWindowHidden(trap.Window) != 0 || quickbarPut(trap, id, p) == 0 {
		if source != trap {
			first := 4
			if *quickbarWord(1049476) != 0 {
				first = 0
			}
			for row := first; row < 5; row++ {
				b := quickbarAt(1048196 + uintptr(row)*256)
				i := quickbarPoint(b, p)
				if i >= 0 {
					oldID, oldFlag = uint16(b.Current[i].ID), byte(b.Current[i].Flags)
				}
				if quickbarPut(b, id, p) == 0 {
					continue
				}
				if Nox_xxx_bookGetSpellDnDType_477670() == 1 {
					flag = 0
					if bool(nox_xxx_spellHasFlags_424A50(int(id), 0x600)) {
						flag = 1
					}
				}
				*quickbarFlagByte(&b.Current[i]) = flag
				quickbarDirections(b)
				break
			}
		}
	}
	return uint32(oldID)<<16 | uint32(oldFlag)
}

func quickbarSwap(b *quickbarRecord, first, second int) {
	b.Current[first], b.Current[second] = b.Current[second], b.Current[first]
	quickbarSlotTooltip(b, first)
	quickbarSlotTooltip(b, second)
}

func quickbarSwapRows(first, second *[5]quickbarSlot, i, j int) {
	first[i], second[j] = second[j], first[i]
	w := quickbarMain().Slots[j]
	w.DrawData().SetTooltip(GetClient().Cli().Strings(), GoWString(nox_xxx_spellTitle_424930(int(second[j].ID))))
}

func quickbarBuildTrap() {
	b := quickbarAt(1047940)
	var ids [5]uint32
	n := 0
	for i := 0; i < 3; i++ {
		if id := b.Current[i].ID; id != 0 {
			ids[n] = id
			n++
		}
	}
	if n == 0 {
		if p := quickbarPlayer(); p != nil && *bookPlayerWord(unsafe.Pointer(p), 3832, 0) != 0 {
			Nox_xxx_printCentered_445490(quickbarText("TrapError"))
			bookSound(925)
		}
		return
	}
	*quickbarByte(1049488) = 0
	ids[n] = 34
	quickbarSendSpells(&ids[0], n+1, *quickbarByte(1047924))
	*quickbarWord(1047916) = 0
	*quickbarWord(1049480) = 0
}

func quickbarDirectionLit(b *quickbarRecord, i int) int {
	return int((b.Directions[i].DrawData().Field0 >> 1) & 1)
}

func quickbarDrawOne(*gui.Window, *gui.WindowData) int { return 1 }

func quickbarSlotPosition(slot int, out *[2]int32) {
	if out != nil && slot >= 0 && slot < 5 {
		w := quickbarMain().Slots[slot]
		if w != nil {
			p := uiWindowPosition(w)
			out[0], out[1] = int32(p.X), int32(p.Y)
		}
	}
}
