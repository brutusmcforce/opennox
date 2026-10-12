package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"image"
	"sort"
	"unsafe"
)

var bookVector [2]float32

func bookContents() *uint32  { return (*uint32)(unsafe.Pointer(&nox_xxx_aNox_cfg_0_587000_132132)) }
func bookSelection() *uint32 { return (*uint32)(unsafe.Pointer(&nox_xxx_aNox_cfg_0_587000_132136)) }

// Named words remain shared with the quickbar until that owner is ported.
func bookWord(off uintptr) *uint32 {
	switch off {
	case 1046636:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046636))
	case 1046640:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046640))
	case 1046648:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046648))
	case 1046652:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046652))
	case 1046656:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046656))
	case 1046852:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046852))
	case 1046864:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046864))
	case 1046868:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046868))
	case 1046872:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046872))
	case 1046924:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046924))
	case 1046928:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046928))
	case 1046932:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046932))
	case 1046936:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046936))
	case 1046944:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046944))
	case 1046948:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046948))
	case 1046952:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046952))
	case 1046956:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1046956))
	case 1047512:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047512))
	case 1047516:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047516))
	case 1047520:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047520))
	case 1047524:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047524))
	case 1047528:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047528))
	case 1047532:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047532))
	case 1047536:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047536))
	case 1047540:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1047540))
	}
	return memmap.PtrUint32(0x5D4594, off)
}
func bookWindow(word uint32) *gui.Window { return AsWindowP(unsafe.Pointer(uintptr(word))) }
func bookRoot() *gui.Window              { return legacyGlobals.nox_win_unk1 }
func bookClass(p uint32) int             { return int(*(*byte)(unsafe.Pointer(uintptr(p) + 2251))) }
func bookPlayerWord(p uint32, off, id int) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(p) + uintptr(off+4*id)))
}
func bookSound(id sound.ID) { Nox_xxx_clientPlaySoundSpecial_452D80(id, 100) }
func bookHideWindow(w *gui.Window, hidden bool) int {
	if w == nil {
		return -2
	}
	w.SetHidden(hidden)
	return 0
}
func bookName(id int) *wchar2_t {
	if *bookWord(1046868) == 1 {
		return (*wchar2_t)(unsafe.Pointer(uintptr(uint32(bookGuideCreatureName(int32(id))))))
	}
	if bookClass(*bookWord(1047516)) != 0 {
		return nox_xxx_spellTitle_424930(id)
	}
	return nox_xxx_abilityGetName_0_425260(id)
}
func bookCompare(a, b uint32) int {
	x, y := bookName(int(a)), bookName(int(b))
	if x == nil || y == nil {
		return 0
	}
	// Keep the live locale/code-unit comparator; Go Unicode folding differs.
	return int(textCompareWide((*uint16)(unsafe.Pointer(x)), (*uint16)(unsafe.Pointer(y))))
}
func bookSetColor() int {
	*bookWord(1046880) = noxcolor.RGB5551Color(15, 15, 15).Color32()
	v := noxcolor.RGB5551Color(115, 100, 100).Color32()
	*bookWord(1046884) = v
	return int(v)
}
func bookToggle() {
	if uiWindowHidden(bookRoot()) != 0 {
		bookShow(0)
	} else {
		bookHide(0)
	}
}
func bookHide(reset int) int {
	if uiWindowHidden(bookRoot()) != 0 || *bookWord(1047520) == 1 {
		return 0
	}
	*bookWord(1046864) = 0
	GetClient().Cli().GUI.ValYYY = 1
	*bookWord(1046868) = 0
	if *bookWord(1046872) != 0 {
		*bookWord(1046868) = 1
	}
	g := GetClient().Cli().GUI
	w := bookRoot()
	if g.Captured() == w {
		w.Capture(false)
	}
	bookHideWindow(w, true)
	bookSound(787)
	if reset != 0 {
		nox_xxx_aNox_cfg_0_587000_132132 = 1
		*bookWord(1046936) = 0
		bookHideWindow(bookWindow(*bookWord(1046952)), true)
	}
	if int32(Nox_xxx_bookGetSpellDnDType_477670()) == 1 {
		Nox_xxx_bookSpellDnDclear_477660()
	}
	return 1
}
func bookShow(force int) {
	if nox_xxx_guiCursor_477600() != 0 {
		return
	}
	if Nox_xxx_playerAnimCheck_4372B0() != 0 && noxflags.HasGame(noxflags.GameModeCoop) {
		return
	}
	if dword_8531A0_2576 == 0 || bookSort(bookClass(uint32(dword_8531A0_2576))) != 0 {
		bookOpen(force)
		return
	}
	Nox_xxx_printCentered_445490(GetServer().S().Strings().GetStringInFile("EmptyBook", "guibook.c"))
	bookSound(925)
}
func bookOpen(force int) {
	if sessionQuitShown() != 0 || optionsVisible() != 0 {
		return
	}
	if force != 0 || Nox_xxx_get_57AF20() == 0 {
		*bookWord(1046864) = 1
		*bookWord(1046868) = 0
		if *bookWord(1046872) != 0 {
			*bookWord(1046868) = 1
		}
		bookRoot().ShowModal()
		bookSound(786)
	}
}
func bookSort(class int) int {
	*bookWord(1046656) = uint32(nox_xxx_guiFontHeightMB_43F320(nil) + 2)
	*bookWord(1047508) = 0
	*bookWord(1047512) = 0
	capacity := 2*(141 / *bookWord(1046656)) - 2
	table := unsafe.Slice(bookWord(1046960), 137)
	add := func(id int) { table[*bookWord(1047508)] = uint32(id); *bookWord(1047508)++ }
	all := noxflags.HasGame(noxflags.GameFlag(0x2000)) && !noxflags.HasGame(noxflags.GameModeQuest)
	p := *bookWord(1047516)
	if *bookWord(1046868) == 1 {
		for id := bookGuideFirst(); id != 0; id = bookGuideNext(int32(id)) {
			if (all || *bookPlayerWord(p, 4244, int(id)) != 0) && questRuntimeBookAllowed(int(id), 37) {
				add(int(id))
			}
		}
	} else if class != 0 {
		for id := nox_xxx_spellFirstValid_424AD0(); id != 0; id = nox_xxx_spellNextValid_424AF0(id) {
			if id == 34 || playerSpellClassCheck(int32(class), int32(id)) != 0 || (!all && *bookPlayerWord(p, 3696, id) == 0) {
				continue
			}
			if bool(nox_xxx_spellHasFlags_424A50(id, 0x15000)) {
				*bookWord(1047512)++
			}
			if !bool(nox_xxx_spellHasFlags_424A50(id, 0x2000)) {
				add(id)
			}
		}
	} else {
		for id := nox_xxx_bookFirstKnownAbil_425330(); id != 0; id = nox_xxx_bookNextKnownAbil_425350(id) {
			if all || *bookPlayerWord(p, 3696, int(id)) != 0 {
				add(int(id))
			}
		}
	}
	count := *bookWord(1047508) - *bookWord(1047512)
	*bookWord(1046940) = count/capacity + 1
	normal := table[:count]
	sort.SliceStable(normal, func(i, j int) bool { return bookCompare(normal[i], normal[j]) < 0 })
	if *bookWord(1047508) != 0 {
		return 1
	}
	return 0
}
func bookTemporaryShow(show int) int {
	result := int(*bookWord(1046864))
	if result == 0 {
		return result
	}
	if show != 0 {
		return bookHideWindow(bookRoot(), false)
	}
	result = int(sub_47A260())
	if result == 0 {
		result = bookHideWindow(bookRoot(), true)
	}
	return result
}
func bookIconPosition(p *image.Point) uintptr {
	if p == nil {
		return 0
	}
	w := bookWindow(*bookWord(1046952))
	if w == nil {
		return uintptr(uint32(0xfffffffe))
	}
	*p = w.GlobalPos()
	return 0
}
func bookPath(from, to image.Point) int {
	n := *bookWord(1046680)
	if int32(n) < 20 {
		n++
		*memmap.PtrFloat32(0x5D4594, 1046676+8*uintptr(n)) = float32(from.X)
		*memmap.PtrFloat32(0x5D4594, 1046680+8*uintptr(n)) = float32(from.Y)
		*memmap.PtrFloat32(0x5D4594, 1046684+8*uintptr(n)) = float32(to.X)
		*bookWord(1046680) = n
		*memmap.PtrFloat32(0x5D4594, 1046688+8*uintptr(n)) = float32(to.Y)
	}
	return int(n)
}
func bookStopAddition() {
	if *bookWord(1047520) == 0 {
		return
	}
	*bookWord(1047520) = 0
	bookHideWindow(bookWindow(*bookWord(1046956)), true)
	*bookWord(1046648) = 0
	quickbarSelectRow(int(*bookWord(1046612)))
	if noxflags.HasGame(noxflags.GameModeCoop) {
		runtimePauseStop()
		sub_413A00(0)
	}
}

func bookShown() int { return bool2int(bookRoot().GetFlags()&16 == 0) }
