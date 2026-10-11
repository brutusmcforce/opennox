package legacy

/*
#include <stdint.h>
#include "client__gui__window.h"
extern uint32_t dword_5d4594_1047548;
extern uint32_t dword_5d4594_1047552;
extern uint32_t dword_5d4594_1049536;
extern nox_window* nox_win_unk1;
*/
import "C"
import (
	"image"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
)

// RelayoutQuickBar updates positions that the quick bar keeps for the GUI size (see nox_xxx_quickBarCreate_45E190).
// The windows are moved separately.
func RelayoutQuickBar(old, cur image.Point) {
	x, y := (cur.X-320)/2, cur.Y-74
	C.dword_5d4594_1047548 = C.uint32_t(x)
	C.dword_5d4594_1047552 = C.uint32_t(y)
	// drop area of the quick bar
	*memmap.PtrInt32(0x587000, 133656) = int32(x)
	*memmap.PtrInt32(0x587000, 133660) = int32(y - 17)
	*memmap.PtrInt32(0x587000, 133664) = int32(x + 320)
	*memmap.PtrInt32(0x587000, 133668) = int32(cur.Y)
	// position the trap button slides to; it's below the screen while it's hidden
	if int(int32(C.dword_5d4594_1049536)) > old.Y {
		C.dword_5d4594_1049536 = C.uint32_t(cur.Y + 1)
	} else {
		C.dword_5d4594_1049536 = C.uint32_t(int(int32(C.dword_5d4594_1049536)) + cur.Y - old.Y)
	}
}

// SpellBookListWindow returns the spell list window of the spell book.
func SpellBookListWindow() *gui.Window {
	return asWindow(C.nox_win_unk1)
}
