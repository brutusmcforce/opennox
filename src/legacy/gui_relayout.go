package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
)

// RelayoutQuickBar updates positions that the quick bar keeps for the GUI size (see nox_xxx_quickBarCreate_45E190).
// The windows are moved separately.
func RelayoutQuickBar(old, cur image.Point) {
	x, y := (cur.X-320)/2, cur.Y-74
	dword_5d4594_1047548 = uint32(x)
	dword_5d4594_1047552 = uint32(y)
	// drop area of the quick bar
	*memmap.PtrInt32(0x587000, 133656) = int32(x)
	*memmap.PtrInt32(0x587000, 133660) = int32(y - 17)
	*memmap.PtrInt32(0x587000, 133664) = int32(x + 320)
	*memmap.PtrInt32(0x587000, 133668) = int32(cur.Y)
	// position the trap button slides to; it's below the screen while it's hidden
	if int(int32(dword_5d4594_1049536)) > old.Y {
		dword_5d4594_1049536 = uint32(cur.Y + 1)
	} else {
		dword_5d4594_1049536 = uint32(int(int32(dword_5d4594_1049536)) + cur.Y - old.Y)
	}
}

// SpellBookListWindow returns the spell list window of the spell book.
func SpellBookListWindow() *gui.Window {
	return legacyGlobals.nox_win_unk1
}
