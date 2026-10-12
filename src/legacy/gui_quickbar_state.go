package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

type quickbarSlot struct{ ID, Flags uint32 }
type quickbarRecord struct {
	Rows       [5][5]quickbarSlot
	Selected   byte
	_          [3]byte
	Current    *[5]quickbarSlot
	Window     *gui.Window
	Slots      [5]*gui.Window
	Directions [5]*gui.Window
	_          uint32
}

// Some legacy transitions only compute a row address, without accessing it.
// Keep that arithmetic separate from indexing the five valid owned rows.
func quickbarRowPointer(b *quickbarRecord, row byte) *[5]quickbarSlot {
	return (*[5]quickbarSlot)(unsafe.Add(unsafe.Pointer(b), 40*uintptr(row)))
}

// The shared slot and record layouts remain visible to remaining C owners.
var _ [256 - int(unsafe.Sizeof(quickbarRecord{}))]byte
var _ [int(unsafe.Sizeof(quickbarRecord{})) - 256]byte
var _ [204 - int(unsafe.Offsetof(quickbarRecord{}.Current))]byte
var _ [int(unsafe.Offsetof(quickbarRecord{}.Current)) - 204]byte

func quickbarMain() *quickbarRecord {
	return (*quickbarRecord)(legacyGlobals.nox_xxx_aClosewoodengat_587000_133480)
}
func quickbarAt(off uintptr) *quickbarRecord { return (*quickbarRecord)(memmap.PtrOff(0x5D4594, off)) }
func quickbarPlayer() uint32                 { return uint32(dword_8531A0_2576) }
func quickbarWord(off uintptr) *uint32 {
	switch off {
	case 1047548:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047548))
	case 1047552:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047552))
	case 1047932:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047932))
	case 1047936:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047936))
	case 1049484:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049484))
	case 1049496:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049496))
	case 1049500:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049500))
	case 1049504:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049504))
	case 1049508:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049508))
	case 1049512:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049512))
	case 1049516:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049516))
	case 1049520:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049520))
	case 1049524:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049524))
	case 1049532:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049532))
	case 1049536:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049536))
	case 1049692:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049692))
	case 1049696:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1049696))
	}
	return memmap.PtrUint32(0x5D4594, off)
}
func quickbarByte(off uintptr) *byte          { return memmap.PtrUint8(0x5D4594, off) }
func quickbarWindow(off uintptr) *gui.Window  { return bookWindow(*quickbarWord(off)) }
func quickbarUserData(w *gui.Window) *uint32  { return (*uint32)(unsafe.Add(w.C(), 368)) }
func quickbarFlagByte(s *quickbarSlot) *byte  { return (*byte)(unsafe.Pointer(&s.Flags)) }
func quickbarPointer(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }
func quickbarLit(w *gui.Window, on bool) uint32 {
	if w != nil {
		if on {
			w.DrawData().Field0 |= 2
		} else {
			w.DrawData().Field0 &^= 2
		}
	}
	return quickbarPointer(w.C())
}
func quickbarDirection(b *quickbarRecord, slot int) uint32 {
	s := &b.Current[slot]
	return quickbarLit(b.Directions[slot], s.ID != 0 && s.Flags&1 != 0)
}
func quickbarDirections(b *quickbarRecord) uint32 {
	var ret uint32
	for i := 0; i < 5; i++ {
		ret = quickbarDirection(b, i)
	}
	return ret
}
func quickbarSelectRow(row int) uint32 {
	if row < 0 || row >= 5 || *quickbarWord(1049476) != 0 || *quickbarWord(1049496) != 0 {
		return uint32(row)
	}
	b := quickbarMain()
	b.Selected = byte(row)
	b.Current = &b.Rows[row]
	if quickbarPlayer() != 0 {
		bookSound(798)
	}
	return quickbarDirections(b)
}
func quickbarFirstEmpty(currentOnly bool) int {
	b := quickbarMain()
	row := int(b.Selected)
	for {
		for slot := 0; slot < 5; slot++ {
			if b.Rows[row][slot].ID == 0 {
				if !currentOnly {
					quickbarSelectRow(row)
				}
				return slot
			}
		}
		if currentOnly {
			return -1
		}
		row = (row + 1) % 5
		if row == int(b.Selected) {
			return -1
		}
	}
}
func quickbarContains(id uint32) int {
	b := quickbarMain()
	row := int(b.Selected)
	for {
		for _, s := range b.Rows[row] {
			if s.ID == id {
				return 1
			}
		}
		row = (row + 1) % 5
		if row == int(b.Selected) {
			return 0
		}
	}
}
func quickbarRemove(id uint32) uint32 {
	b := quickbarMain()
	start := int(b.Selected)
	row := start
	var ret uint32
	for {
		for slot := 0; slot < 5; slot++ {
			if b.Rows[row][slot].ID == id {
				b.Rows[row][slot].ID = 0
			}
		}
		ret = uint32(40 * (row + 1))
		row = (row + 1) % 5
		if row == start {
			return ret
		}
	}
}
func quickbarRestoreSlots() uint32 {
	b := quickbarMain()
	for slot := 0; slot < 5; slot++ {
		for row := 0; row < 5; row++ {
			off := uintptr(40*row + 8*slot)
			b.Rows[row][slot].ID = *quickbarWord(1047564 + off)
			*quickbarFlagByte(&b.Rows[row][slot]) = *quickbarByte(1047568 + off)
		}
	}
	return 232
}
