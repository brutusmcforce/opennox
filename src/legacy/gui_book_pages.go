package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func bookInputBlocked() bool {
	p := LocalPlayerDrawable
	return p != nil && uiShopActive() == 0 && sub_47A260() == 0 && *(*byte)(unsafe.Pointer(&p.ObjFlags))&2 != 0
}
func bookPageComplete(guide bool) int {
	v := uint32(0)
	if guide {
		v = 1
	}
	*bookWord(1046868) = v
	*bookWord(1046872) = v
	return int(v)
}
func bookStartTurn(back, guide bool) {
	off, mode := uintptr(1046924), uint32(2)
	if back {
		off, mode = 1046928, 3
	}
	anim := (*ImageRef)(memmap.WordPtr(*bookWord(off))).Field24ptr()
	if guide {
		anim.OnEnd = bookImageEndKey(1)
	} else {
		anim.OnEnd = bookImageEndKey(0)
	}
	*bookWord(1046868) = mode
	anim.Field_3 = uint32(nox_xxx_bookGet_430B40_get_mouse_prev_seq())
}
func bookForward(event int) int {
	if *bookWord(1047520) == 1 {
		return 1
	}
	if event == 10 {
		bookToggle()
		return 1
	}
	if bookInputBlocked() {
		return 1
	}
	if event != 5 {
		return 0
	}
	bookHideWindow(bookWindow(*bookWord(1046944)), false)
	if *bookContents() != 0 {
		if *bookWord(1046936)+1 < *bookWord(1046940) {
			*bookWord(1046936)++
		} else {
			*bookWord(1046932) = 0
			*bookContents() = 0
			bookHideWindow(bookWindow(*bookWord(1046952)), false)
		}
		bookStartTurn(false, *bookWord(1046872) != 0)
		bookSound(788)
	} else if *bookWord(1046932) < *bookWord(1047508)-*bookWord(1047512)-1 {
		*bookWord(1046932)++
		bookStartTurn(false, *bookWord(1046872) != 0)
		bookSound(788)
	}
	return 1
}
func bookBackward(event int) int {
	if *bookWord(1047520) == 1 {
		return 1
	}
	if event == 10 {
		bookToggle()
		return 1
	}
	if bookInputBlocked() {
		return 1
	}
	if event != 5 {
		return 0
	}
	bookHideWindow(bookWindow(*bookWord(1046948)), false)
	if *bookContents() == 0 {
		if int32(*bookWord(1046932)) < int32(*bookWord(1047508))-int32(*bookWord(1047512)) {
			if int32(*bookWord(1046932)) <= 0 {
				*bookContents() = 1
				*bookWord(1046936) = *bookWord(1046940) - 1
				bookHideWindow(bookWindow(*bookWord(1046952)), true)
			} else {
				*bookWord(1046932)--
			}
		} else {
			*bookContents() = 1
			*bookWord(1046936) = 0
			bookHideWindow(bookWindow(*bookWord(1046952)), true)
		}
	} else {
		if *bookWord(1046936) == 0 {
			return 1
		}
		*bookWord(1046936)--
	}
	bookStartTurn(true, *bookWord(1046872) != 0)
	bookSound(788)
	return 1
}
func bookMoveToPage(page int) int {
	*bookContents() = 0
	*bookWord(1046932) = uint32(page)
	*bookWord(1046936) = 99
	bookStartTurn(false, *bookWord(1046872) != 0)
	bookSound(788)
	bookHideWindow(bookWindow(*bookWord(1046952)), false)
	bookHideWindow(bookWindow(*bookWord(1046944)), false)
	return bookHideWindow(bookWindow(*bookWord(1046948)), false)
}
func bookTab(w *gui.Window, event uint32) int {
	if *bookWord(1047520) == 1 || bookInputBlocked() {
		return 1
	}
	if event != 5 {
		if event == 6 || event == 7 {
			return 1
		}
		return 0
	}
	changed := false
	switch w.ID() {
	case 1310:
		if *bookWord(1046872) != 0 {
			bookPageComplete(false)
			if bookSort(bookClass(bookPlayer)) == 0 {
				bookPageComplete(true)
				bookSort(bookClass(bookPlayer))
				bookSound(925)
				return 1
			}
			changed = true
		}
		if *bookContents() == 0 || *bookWord(1046936) != 0 || changed {
			*bookContents() = 1
			*bookWord(1046936) = 0
			bookSound(788)
			bookStartTurn(true, false)
			bookHideWindow(bookWindow(*bookWord(1046952)), true)
			bookHideWindow(bookWindow(*bookWord(1046948)), false)
			bookRoot().Capture(false)
			Nox_xxx_bookSpellDnDclear_477660()
		}
	case 1320:
		if *bookWord(1046872) != 1 {
			bookPageComplete(true)
			if bookSort(bookClass(bookPlayer)) == 0 {
				bookPageComplete(false)
				bookSort(bookClass(bookPlayer))
				bookSound(925)
				return 1
			}
			changed = true
		}
		if *bookContents() != 0 && *bookWord(1046936) == 0 && !changed {
			return 1
		}
		*bookContents() = 1
		*bookWord(1046936) = 0
		bookSound(788)
		bookStartTurn(!changed, true)
		bookHideWindow(bookWindow(*bookWord(1046952)), true)
		bookHideWindow(bookWindow(*bookWord(1046948)), false)
		bookRoot().Capture(false)
		Nox_xxx_bookSpellDnDclear_477660()
	}
	return 1
}
func bookListEvents(w *gui.Window, event uint32, pos image.Point) int {
	if *bookWord(1047520) == 1 || bookInputBlocked() {
		return 1
	}
	if *bookContents() == 0 {
		if event >= 5 && event <= 8 {
			return 1
		}
		if event == 11 {
			bookToggle()
			return 1
		}
		return 0
	}
	*bookWord(1046656) = uint32(nox_xxx_guiFontHeightMB_43F320(nil) + 2)
	perColumn := int(int32(141 / *bookWord(1046656) - 1))
	origin := w.Off
	switch event {
	case 5:
		p := LocalPlayerDrawable
		if p != nil && *(*byte)(unsafe.Pointer(&p.ObjFlags))&2 != 0 {
			return 1
		}
		*bookWord(1047536) = uint32(pos.Y)
		*bookWord(1047532) = uint32(pos.X)
		rowY := pos.Y - origin.Y - 19
		if rowY < 0 || rowY/int(int32(*bookWord(1046656))) >= perColumn {
			*bookSelection() = 0xffffffff
			return 1
		}
		row := rowY / int(int32(*bookWord(1046656)))
		if pos.X-origin.X > 145 {
			row += perColumn
		}
		*bookSelection() = 2*uint32(perColumn)**bookWord(1046936) + uint32(row)
		w.Capture(true)
		if *bookSelection() >= *bookWord(1047508)-*bookWord(1047512) {
			*bookSelection() = 0xffffffff
			return 1
		}
		class := bookClass(bookPlayer)
		guide := *bookWord(1046868)
		id := *bookWord(1046960 + 4*uintptr(*bookSelection()))
		if class == 2 && guide == 1 {
			if *bookPlayerWord(bookPlayer, 4232, 0) == 0 && (!noxflags.HasGame(noxflags.GameFlag(0x2000)) || noxflags.HasGame(noxflags.GameModeQuest)) {
				w.Capture(false)
				bookMoveToPage(int(*bookSelection()))
				*bookWord(1047528) = 0
				*bookSelection() = 0xffffffff
				return 1
			}
			id += 74
		} else if guide != 0 {
			return 1
		}
		*bookWord(1047528) = id
		Nox_xxx_bookSaveSpellForDragDrop_477640(int(int32(id)), 1)
		bookSound(793)
		return 1
	case 6, 7:
		if int32(*bookSelection()) >= 0 {
			dx, dy := int(int32(*bookWord(1047532)))-pos.X, int(int32(*bookWord(1047536)))-pos.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			if dx >= 5 || dy >= 5 {
				quickbarDrop(*bookWord(1047528), 0, pos, nil)
			} else {
				bookMoveToPage(int(int32(*bookSelection())))
			}
		}
		if GetClient().Cli().GUI.Captured() != nil {
			w.Capture(false)
		}
		Nox_xxx_bookSpellDnDclear_477660()
		return 1
	case 8:
		return 1
	case 11:
		bookToggle()
		return 1
	}
	return 0
}
