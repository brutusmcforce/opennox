package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"image"
	"unsafe"
)

func bookDrawImage(handle uint32, pos image.Point) {
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(noxrender.ImageHandle(memmap.WordPtr(handle))), pos)
}
func bookSpellImage(id int) uint32 { return memmap.PtrWord(nox_xxx_spellIcon_424A90(id)) }
func bookAbilityImage(id int) uint32 {
	return memmap.PtrWord(unsafe.Pointer(nox_xxx_spellGetAbilityIcon_425310(id, 0)))
}
func bookDrawIcon(w *gui.Window) int {
	mode, p := *bookWord(1046868), bookPlayer
	var img uint32
	if mode != 0 {
		if mode != 1 || bookClass(p) != 2 || (*bookPlayerWord(p, 4232, 0) == 0 && (!noxflags.HasGame(noxflags.GameFlag(0x2000)) || noxflags.HasGame(noxflags.GameModeQuest))) {
			return 1
		}
		img = bookSpellImage(int(*bookWord(1046960 + 4*uintptr(*bookWord(1046932)))) + 74)
	} else if bookClass(p) != 0 {
		img = bookSpellImage(int(*bookWord(1046960 + 4*uintptr(*bookWord(1046932)))))
	} else {
		img = bookAbilityImage(int(*bookWord(1046960 + 4*uintptr(*bookWord(1046932)))))
	}
	if img != 0 {
		bookDrawImage(img, w.GlobalPos())
	}
	return 1
}
func bookIconTooltip(w *gui.Window) int {
	var key string
	switch w.ID() {
	case 1310:
		if bookClass(bookPlayer) == 0 {
			key = "ToolTipAbilityTab"
		} else {
			key = "ToolTipSpellTab"
		}
	case 1320:
		key = "ToolTipGuideTab"
	}
	if key != "" {
		Nox_xxx_cursorSetTooltip_4776B0(GetServer().S().Strings().GetStringInFile(strman.ID(key), "guibook.c"))
	}
	return 1
}
func bookIconEvents(w *gui.Window, event int, pos image.Point) int {
	dr := LocalPlayerDrawable
	if dr != nil && *(*byte)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(120)) /* ObjFlags */))&2 != 0 {
		return 1
	}
	p, mode := bookPlayer, *bookWord(1046868)
	class := bookClass(p)
	// This quest guard intentionally differs from icon drawing (frozen C behavior).
	if mode == 1 && (class == 1 || class == 0 || (!noxflags.HasGame(noxflags.GameFlag(0x2000)) && !noxflags.HasGame(noxflags.GameModeQuest) && *bookPlayerWord(p, 4232, 0) == 0)) {
		return 0
	}
	clear := func() { *bookWord(1047540) = 0; w.Capture(false); Nox_xxx_bookSpellDnDclear_477660() }
	switch event {
	case 5:
		if GetClient().Cli().GUI.Captured() != nil || *bookWord(1047540) != 0 {
			return 1
		}
		id := *bookWord(1046960 + 4*uintptr(*bookWord(1046932)))
		if class == 2 && mode == 1 {
			id += 74
		} else if mode != 0 {
			return 1
		}
		*bookWord(1047540) = id
		if class != 0 && mode == 0 && bool(nox_xxx_spellHasFlags_424A50(int(id), 0x15000)) {
			*bookWord(1047540) = 0
			return 1
		}
		Nox_xxx_bookSaveSpellForDragDrop_477640(int(int32(id)), 1)
		w.Capture(true)
		bookSound(793)
		return 1
	case 6, 7:
		if *bookWord(1047540) == 0 {
			return 1
		}
		origin := uiWindowPosition(w)
		inside := pos.X >= origin.X && pos.X <= origin.X+w.SizeVal.X && pos.Y >= origin.Y && pos.Y <= origin.Y+w.SizeVal.Y
		if !inside {
			if *quickbarWord(1047928) != 0 {
				quickbarQueueTarget(*bookWord(1047540))
			} else if *quickbarWord(1047932) != 0 {
				quickbarSendPendingAbility()
			} else {
				quickbarDrop(*bookWord(1047540), 0, pos, nil)
			}
			clear()
			return 1
		}
		if class != 0 {
			if *quickbarWord(1047928) != 0 {
				quickbarSpellCursor(0, 0)
				clear()
				return 1
			}
			if !bool(nox_xxx_spellHasFlags_424A50(int(*bookWord(1047540)), 0x600)) {
				quickbarSpellCursor(*bookWord(1047540), 0)
				Nox_xxx_bookSpellDnDclear_477660()
				return 1
			}
			if quickbarSpellCursor(*bookWord(1047540), 1) != 0 {
				clear()
				return 1
			}
		} else {
			if *quickbarWord(1047932) != 0 {
				quickbarAbilityCursor(0)
				clear()
				return 1
			}
			if quickbarAbilityCursor(*bookWord(1047540)) != 0 {
				clear()
				return 1
			}
		}
		Nox_xxx_bookSpellDnDclear_477660()
		return 1
	case 8:
		return 1
	}
	return 0
}
