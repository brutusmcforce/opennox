package legacy

import "github.com/brutusmcforce/opennox/client/gui"

func quickbarDestroyWindow(w *gui.Window) int {
	if w == nil {
		return -2
	}
	w.Destroy()
	return 0
}
func quickbarDestroy() int {
	for _, off := range []uintptr{1049500, 1049504, 1049520, 1049508, 1049512, 1049516} {
		quickbarDestroyWindow(quickbarWindow(off))
		*quickbarWord(off) = 0
	}
	for row := 0; row < 5; row++ {
		b := quickbarAt(1048196 + uintptr(256*row))
		quickbarDestroyWindow(b.Window)
		b.Window = nil
		for slot := 0; slot < 5; slot++ {
			quickbarDestroyWindow(b.Slots[slot])
			b.Slots[slot] = nil
			quickbarDestroyWindow(b.Directions[slot])
			b.Directions[slot] = nil
		}
	}
	trap := quickbarAt(1047940)
	quickbarDestroyWindow(trap.Window)
	trap.Window = nil
	var ret int
	for slot := 0; slot < 3; slot++ {
		ret = quickbarDestroyWindow(trap.Slots[slot])
		trap.Slots[slot] = nil
	}
	*quickbarWord(1049532) = 0
	*quickbarWord(1047928) = 0
	*quickbarWord(1047932) = 0
	return ret
}
func quickbarVisible(show bool) int {
	b := quickbarMain()
	if !show && *quickbarWord(1049476) != 0 {
		quickbarExpand()
	}
	for _, off := range []uintptr{1049500, 1049504, 1049520, 1049508, 1049512} {
		bookHideWindow(quickbarWindow(off), !show)
	}
	for slot := 0; slot < 5; slot++ {
		bookHideWindow(b.Slots[slot], !show)
		bookHideWindow(b.Directions[slot], !show)
		if show && *quickbarWord(1049484) != 0 {
			trap := quickbarAt(1047940)
			bookHideWindow(trap.Slots[slot], false)
			bookHideWindow(trap.Directions[slot], false)
			bookHideWindow(quickbarWindow(1049512), true)
		}
	}
	if show {
		if *quickbarWord(1049484) != 0 {
			bookHideWindow(quickbarWindow(1048148), false)
		}
		return bookHideWindow(b.Window, false)
	}
	bookHideWindow(b.Window, true)
	return bookHideWindow(quickbarWindow(1048148), true)
}
func quickbarPrepare() int {
	if *quickbarWord(1049508) != 0 {
		quickbarDestroy()
	}
	if quickbarCreate() == 0 {
		return 0
	}
	quickbarVisible(Nox_client_getRenderGUI() != 0)
	return 1
}
