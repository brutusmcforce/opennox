package legacy

import "unsafe"

func quickbarCopySlot(dst, src *quickbarSlot) {
	dst.ID = src.ID
	*quickbarFlagByte(dst) = byte(src.Flags)
}
func quickbarCloseExpanded() int {
	if *quickbarWord(1049476) != 1 {
		return 0
	}
	b := quickbarMain()
	for row := 0; row < 4; row++ {
		panel := quickbarAt(1048196 + uintptr(256*row))
		for slot := 0; slot < 5; slot++ {
			quickbarCopySlot(&b.Rows[row][slot], &panel.Rows[0][slot])
			bookHideWindow(panel.Directions[slot], true)
		}
		bookHideWindow(panel.Window, true)
	}
	// The prerequisite C correction restores the row saved by expansion.
	b.Selected = *quickbarByte(1047912)
	b.Current = quickbarRowPointer(b, b.Selected)
	quickbarDirections(b)
	bookHideWindow(quickbarWindow(1049512), false)
	bookSound(800)
	*quickbarWord(1049476) = 0
	return 1
}
func quickbarExpand() {
	if *quickbarWord(1049476) != 0 {
		quickbarCloseExpanded()
		return
	}
	if *quickbarWord(1049496) != 0 || Nox_xxx_get_57AF20() != 0 {
		return
	}
	if *quickbarWord(1049484) == 1 {
		quickbarCloseTrap()
	}
	b := quickbarMain()
	for row := 0; row < 4; row++ {
		panel := quickbarAt(1048196 + uintptr(256*row))
		for slot := 0; slot < 5; slot++ {
			quickbarCopySlot(&panel.Rows[0][slot], &b.Rows[row][slot])
			bookHideWindow(panel.Directions[slot], false)
		}
		bookHideWindow(panel.Window, false)
		quickbarDirections(panel)
	}
	*quickbarWord(1047912) = uint32(b.Selected)
	b.Selected = 4
	b.Current = &b.Rows[4]
	quickbarDirections(b)
	bookHideWindow(quickbarWindow(1049512), true)
	bookSound(799)
	*quickbarWord(1049476) = 1
}
func quickbarClearSlots() uint32 {
	quickbarCloseExpanded()
	b := quickbarMain()
	for row := 0; row < 5; row++ {
		panel := quickbarAt(1048196 + uintptr(256*row))
		for slot := 0; slot < 5; slot++ {
			b.Rows[row][slot].ID = 0
			*quickbarFlagByte(&b.Rows[row][slot]) = 0
			panel.Rows[0][slot].ID = 0
			*quickbarFlagByte(&panel.Rows[0][slot]) = 0
		}
	}
	trap := quickbarAt(1047940)
	for row := 0; row < 3; row++ {
		for slot := 0; slot < 3; slot++ {
			trap.Rows[row][slot].ID = 0
			*quickbarFlagByte(&trap.Rows[row][slot]) = 0
		}
	}
	if b.Directions[0] != nil {
		return quickbarDirections(b)
	}
	return quickbarPointer(unsafe.Pointer(b))
}
func quickbarClearAbilities() uint32 {
	for id := 1; id <= 5; id++ {
		*quickbarWord(1047764 + uintptr(24*id) + 16) = 0
	}
	return quickbarClearSlots()
}
func quickbarTrapSelect(row int) int {
	if row < 0 || row >= 3 {
		return row
	}
	b := quickbarAt(1047940)
	b.Selected = byte(row)
	b.Current = &b.Rows[row]
	return 5 * int(b.Selected)
}
func quickbarTrapNext() {
	row := *quickbarByte(1048140)
	if row == 2 {
		row = 0
	} else {
		row++
	}
	b := quickbarAt(1047940)
	b.Selected = row
	b.Current = quickbarRowPointer(b, row)
	bookSound(798)
}
func quickbarTrapPrevious() {
	row := *quickbarByte(1048140)
	if row == 0 {
		row = 2
	} else {
		row--
	}
	b := quickbarAt(1047940)
	b.Selected = row
	b.Current = quickbarRowPointer(b, row)
	bookSound(798)
}
func quickbarMoveRow(delta int) uint32 {
	if v := *quickbarWord(1049476); v != 0 {
		return v
	}
	if v := *quickbarWord(1049496); v != 0 {
		return v
	}
	if v := Nox_xxx_get_57AF20(); v != 0 {
		return uint32(v)
	}
	quickbarLastButton(-1)
	row := int(quickbarMain().Selected) + delta
	if delta > 0 && row > 4 {
		row = 0
	} else if delta < 0 && row < 0 {
		row = 4
	}
	return quickbarSelectRow(row)
}
func quickbarTimedRow() {
	if *quickbarWord(1049476) != 0 || *quickbarWord(1049496) != 0 {
		return
	}
	var row uint32
	if InputKeyCheckTimeoutLegacy(7, GetServer().S().TickRate()>>1) {
		*quickbarWord(1049712) = 0
	} else {
		*quickbarWord(1049712)++
		row = *quickbarWord(1049712)
		if int32(row) >= 5 {
			quickbarMain().Selected = 4
			return
		}
	}
	b := quickbarMain()
	b.Selected = byte(row)
	b.Current = quickbarRowPointer(b, b.Selected)
	InputSetKeyTimeoutLegacy(7)
	bookSound(798)
	quickbarDirections(b)
}
func quickbarCloseTrap() {
	if *quickbarWord(1049484) == 0 {
		return
	}
	bookHideWindow(quickbarWindow(1048148), true)
	bookHideWindow(quickbarWindow(1049512), false)
	quickbarLit(quickbarWindow(1049500), false)
	bookSound(797)
	*quickbarWord(1049484) = 0
}
func quickbarToggleTrap() {
	if *quickbarWord(1049484) == 1 {
		quickbarCloseTrap()
		return
	}
	if *quickbarWord(1049476) == 1 {
		quickbarCloseExpanded()
	}
	bookHideWindow(quickbarWindow(1048148), false)
	bookHideWindow(quickbarWindow(1049512), true)
	quickbarLit(quickbarWindow(1049500), true)
	bookSound(796)
	*quickbarWord(1049484) = 1
}
func quickbarCancelCapture() uint32 {
	ret := *quickbarWord(1049532)
	if ret == 0 && *quickbarWord(1047928) == 0 && *quickbarWord(1047932) == 0 {
		return ret
	}
	quickbarWindow(1049532).Capture(false)
	*quickbarWord(1049532) = 0
	*quickbarWord(1047928) = 0
	*quickbarWord(1047932) = 0
	return 1
}
func quickbarRememberMouseSequence() {
	if *quickbarWord(1049476) == 0 {
		*quickbarWord(1049492) = uint32(nox_xxx_bookGet_430B40_get_mouse_prev_seq())
	}
}
