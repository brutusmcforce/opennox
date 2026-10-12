package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func quickbarCorpse() bool {
	p := LocalPlayerDrawable
	return p != nil && *(*byte)(unsafe.Add(p.C(), 120))&2 != 0
}

func quickbarSlotEvent(w *gui.Window, event, arg uint32) int {
	b := (*quickbarRecord)(unsafe.Pointer(uintptr(*quickbarUserData(w))))
	if quickbarCorpse() {
		return 1
	}
	i := 0
	for i < 5 && b.Slots[i] != w {
		i++
	}
	if i == 5 {
		return 0
	}
	s := &b.Current[i]
	pos := bookPoint(arg)
	switch event {
	case 5:
		if s.ID == 0 || *quickbarWord(1047928) != 0 || *quickbarWord(1047932) != 0 || Nox_xxx_get_57AF20() != 0 ||
			bookClass(quickbarPlayer()) != 0 && !bool(nox_xxx_spellIsEnabled_424B70(int(s.ID))) {
			return 1
		}
		w.Capture(true)
		Nox_xxx_bookSaveSpellForDragDrop_477640(int(s.ID), 2)
		*quickbarWord(1049692) = quickbarPointer(unsafe.Pointer(b))
		*quickbarWord(1049696) = quickbarPointer(unsafe.Pointer(b.Current))
		bookSound(793)
		return 1
	case 6, 7:
		if GetClient().Cli().GUI.Captured() == nil {
			return 1
		}
		defer func() {
			if *quickbarWord(1047928) == 0 && *quickbarWord(1047932) == 0 {
				w.Capture(false)
				*quickbarWord(1049532) = 0
			}
			Nox_xxx_bookSpellDnDclear_477660()
		}()
		if Nox_xxx_get_57AF20() != 0 {
			return 1
		}
		j := quickbarPoint(b, pos)
		if *quickbarWord(1047932) != 0 {
			if j < 0 {
				quickbarSendPendingAbility()
			}
			*quickbarWord(1047932) = 0
			return 1
		}
		if *quickbarWord(1047928) != 0 {
			if j < 0 {
				quickbarQueueTarget(*quickbarWord(1047556))
			}
			*quickbarWord(1047928) = 0
			return 1
		}
		if j < 0 {
			old := quickbarDrop(s.ID, byte(s.Flags), pos, b)
			if old>>16 == 137 {
				point := [2]int32{int32(pos.X), int32(pos.Y)}
				inside := geometryRectInt(&point, (*[4]int32)(unsafe.Pointer(memmap.PtrOff(0x587000, 133656)))) != 0
				previous := *quickbarWord(1049696)
				if inside || previous != 0 && previous != quickbarPointer(unsafe.Pointer(b.Current)) {
					return 1
				}
				b.Current[i].ID = 0
				*quickbarFlagByte(&b.Current[i]) = 0
			} else {
				b.Current[i].ID = old >> 16
				*quickbarFlagByte(&b.Current[i]) = byte(old)
			}
		} else {
			if *quickbarWord(1049692) == quickbarPointer(unsafe.Pointer(b)) && b == quickbarMain() {
				previous := *quickbarWord(1049696)
				if previous != 0 && previous != quickbarPointer(unsafe.Pointer(b.Current)) {
					quickbarSwapRows((*[5]quickbarSlot)(unsafe.Pointer(uintptr(previous))), b.Current, i, j)
					bookSound(794)
					quickbarDirections(b)
					return 1
				}
			}
			if j == i {
				*quickbarWord(1049532) = quickbarPointer(w.C())
				if quickbarPlayer() == 0 || b == quickbarAt(1047940) {
					w.Capture(false)
					*quickbarWord(1049532) = 0
					return 1
				}
				if bookClass(quickbarPlayer()) != 0 {
					if bool(nox_xxx_spellIsEnabled_424B70(int(s.ID))) {
						quickbarSpellCursor(s.ID, byte(s.Flags))
						quickbarLastButton(i)
					}
				} else {
					quickbarAbilityCursor(s.ID)
					quickbarLastButton(i)
				}
				return 1
			}
			quickbarSwap(b, i, j)
			bookSound(794)
		}
		if b != quickbarAt(1047940) {
			quickbarDirections(b)
		}
		return 1
	case 8, 12, 16:
		return 0
	default:
		return 1
	}
}

func quickbarBookEvent(_ *gui.Window, event, _ uint32) int {
	switch event {
	case 5, 6:
		return 1
	case 7:
		bookToggle()
		return 1
	}
	return 0
}

func quickbarDirectionEvent(w *gui.Window, event, _ uint32) int {
	data := *quickbarUserData(w)
	i := int(uint16(data))
	b := quickbarAt(1048196 + uintptr(data>>16)*256)
	s := &b.Current[i]
	if s.ID == 0 {
		return 0
	}
	if event != 5 {
		if event == 6 || event == 7 {
			return 1
		}
		return 0
	}
	if bool(nox_xxx_spellHasFlags_424A50(int(s.ID), 0x200400)) {
		bookSound(925)
	} else {
		*quickbarFlagByte(s) ^= 1
		quickbarDirection(b, i)
		bookSound(921)
	}
	return 1
}

func quickbarRowEvent(w *gui.Window, event, _ uint32) int {
	action := *quickbarUserData(w)
	if quickbarCorpse() || optionsVisible() != 0 {
		return 1
	}
	if event == 5 {
		switch action {
		case 0, 1, 2:
			quickbarWindow(1049508).DrawData().BgImageHnd = nil
			switch action {
			case 0:
				quickbarMoveRow(-1)
			case 1:
				quickbarMoveRow(1)
			case 2:
				quickbarExpand()
			}
		case 3:
			quickbarTrapPrevious()
		case 4:
			quickbarTrapNext()
		}
		*quickbarWord(1049700), *quickbarWord(1049704) = 2, action
	} else if event != 6 && event != 7 {
		return 0
	}
	return 1
}

func quickbarTrapTargets() int {
	b := quickbarAt(1047940)
	for i := 0; i < 3; i++ {
		if id := b.Current[i].ID; id != 0 && bool(nox_xxx_spellHasFlags_424A50(int(id), 8)) {
			return 1
		}
	}
	return 0
}

func quickbarTrapButtonEvent(w *gui.Window, event, arg uint32) int {
	if event == 5 || event == 6 {
		return 1
	}
	if event != 7 {
		quickbarLit(quickbarWindow(1049504), false)
		return 0
	}
	if *quickbarWord(1047928) != 0 {
		p := bookPoint(arg)
		if !uiWindowPointIn(w, int32(p.X), int32(p.Y)) {
			quickbarBuildTrap()
		}
		w.Capture(false)
		*quickbarWord(1049532), *quickbarWord(1047928) = 0, 0
		return 1
	}
	quickbarLit(quickbarWindow(1049504), true)
	if p := quickbarPlayer(); p != 0 {
		if bookClass(p) == 1 && quickbarTrapTargets() != 0 {
			w.Capture(true)
			*quickbarWord(1049532) = quickbarPointer(w.C())
			*quickbarWord(1047928) = 1
		} else {
			quickbarBuildTrap()
		}
	}
	return 1
}

func quickbarGenericEvent(_ *gui.Window, event, _ uint32) int {
	if event == 8 || event == 12 || event == 16 {
		return 0
	}
	return 1
}

func quickbarModifier() {
	b := quickbarMain()
	if b == nil {
		return
	}
	same := *quickbarWord(1049492) == uint32(nox_xxx_bookGet_430B40_get_mouse_prev_seq())
	if *quickbarWord(1049496) != 0 {
		if same {
			return
		}
		*quickbarWord(1049496) = 0
	} else {
		if !same {
			return
		}
		*quickbarWord(1049496) = 1
	}
	for i := 0; i < 5; i++ {
		s := &b.Current[i]
		if s.ID != 0 && !bool(nox_xxx_spellHasFlags_424A50(int(s.ID), 0x200400)) {
			*quickbarFlagByte(s) ^= 1
			quickbarDirection(b, i)
		}
	}
	quickbarDirections(b)
	bookSound(921)
}

func quickbarTrapEvent(w *gui.Window, event, _ uint32) int {
	if event == 5 || event == 6 {
		return 1
	}
	if event == 7 {
		if uiWindowHidden(w) == 0 {
			quickbarToggleTrap()
		}
		return 1
	}
	return 0
}
