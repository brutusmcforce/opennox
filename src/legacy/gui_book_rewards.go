package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func bookFind(id, limit int) int {
	for i := 0; i < limit; i++ {
		if *bookWord(1046960 + 4*uintptr(i)) == uint32(id) {
			return i
		}
	}
	return limit
}
func bookRewardPresent(kind, id, auto int) {
	presentationBookReward(kind, id, auto)
}
func bookSpellFamily(id int) int {
	for _, pair := range [][2]int{{0x1000, 0x2000}, {0x4000, 0x8000}, {0x10000, 0x20000}} {
		if bool(nox_xxx_spellHasFlags_424A50(id, pair[0])) {
			return pair[1]
		}
	}
	return 0
}
func bookGuideFamily(id int, visit func(int)) {
	for off := uintptr(132124); ; off += 4 {
		p := memmap.Uint32(0x587000, off)
		if p == 0 {
			return
		}
		words := (*uint32)(unsafe.Pointer(uintptr(p)))
		if *words != uint32(id) {
			continue
		}
		for i := uintptr(4); ; i += 4 {
			n := *(*uint32)(unsafe.Add(unsafe.Pointer(words), i))
			if n == 0 {
				break
			}
			visit(int(n))
		}
	}
}
func bookSpellReward(id, rank, notify, auto int) {
	p := uint32(dword_8531A0_2576)
	if p == 0 {
		return
	}
	if playerSpellClassCheck(int32(bookClass(p)), int32(id)) == 9 {
		return
	}
	*bookPlayerWord(p, 3696, id) = uint32(rank)
	if family := bookSpellFamily(id); family != 0 {
		for n := 1; n < 137; n++ {
			if bool(nox_xxx_spellHasFlags_424A50(n, family)) && bool(nox_xxx_spellIsValid_424B50(n)) {
				*bookPlayerWord(p, 3696, n) = uint32(rank)
			}
		}
	}
	bookPageComplete(false)
	bookSort(bookClass(p))
	if id == 34 {
		show := 0
		if *bookPlayerWord(p, 3832, 0) == 1 && notify != 0 {
			show = 1
		}
		quickbarAddTrap(show)
	} else if notify != 0 {
		if index := bookFind(id, 137); index != 137 {
			bookHide(0)
			bookMoveToPage(index)
			bookOpen(0)
			bookRewardPresent(2, id, auto)
		}
	}
}
func bookGuideReward(id, notify int) {
	p := uint32(dword_8531A0_2576)
	if p == 0 {
		return
	}
	*bookPlayerWord(p, 4244, id) = 1
	bookGuideFamily(id, func(n int) { *bookPlayerWord(p, 4244, n) = 1 })
	bookPageComplete(true)
	bookSort(bookClass(p))
	if notify != 0 {
		if index := bookFind(id, 41); index != 41 {
			bookHide(0)
			bookMoveToPage(index)
			bookOpen(0)
			bookRewardPresent(4, id, 0)
		}
	}
}
func bookSetForward(kind uintptr, id int, pos image.Point) uintptr {
	bookShow(1)
	bookRoot().SetPos(pos)
	limit := 0
	switch kind {
	case 2, 4:
		limit = 137
	case 3:
		limit = 6
	default:
		return kind
	}
	index := bookFind(id, limit)
	if index != limit {
		return uintptr(bookMoveToPage(index))
	}
	return uintptr(unsafe.Pointer(bookWord(1046960 + 4*uintptr(index))))
}
func bookAbilityReward(id int, notify uintptr, auto int) {
	p := nox_common_playerInfoGetByID_417040(int(nox_player_netCode_85319C))
	if p == nil {
		return
	}
	bookPageComplete(false)
	bookSort(bookClass(uint32(uintptr(unsafe.Pointer(p)))))
	if notify != 0 {
		if index := bookFind(id, 6); index != 6 {
			bookHide(0)
			bookMoveToPage(index)
			bookOpen(0)
			if auto != 0 {
				bookRewardPresent(3, id, auto)
			}
		}
	}
}
func bookRemoveSpell(id int) int {
	p := uint32(dword_8531A0_2576)
	result := bookHide(1)
	if p == 0 {
		return result
	}
	*bookPlayerWord(p, 3696, id) = 0
	quickbarRemove(uint32(id))
	if family := bookSpellFamily(id); family != 0 {
		for n := 1; n < 137; n++ {
			if bool(nox_xxx_spellHasFlags_424A50(n, family)) && bool(nox_xxx_spellIsValid_424B50(n)) {
				*bookPlayerWord(p, 3696, n) = 0
				quickbarRemove(uint32(n))
			}
		}
	}
	return bookSort(bookClass(p))
}
func bookRemoveGuide(id int) int {
	p := uint32(dword_8531A0_2576)
	result := bookHide(1)
	if p == 0 {
		return result
	}
	*bookPlayerWord(p, 4244, id) = 0
	quickbarRemove(uint32(id + 74))
	bookGuideFamily(id, func(n int) { *bookPlayerWord(p, 4244, n) = 0; quickbarRemove(uint32(n + 74)) })
	return bookSort(bookClass(p))
}
func bookRemoveAbility(id int) uintptr {
	bookHide(1)
	quickbarAbilityReward(id, 0, 0)
	if p := *bookWord(1047516); p != 0 {
		*bookPlayerWord(p, 3696, id) = 0
	}
	quickbarRemove(uint32(id))
	p := nox_common_playerInfoGetByID_417040(int(nox_player_netCode_85319C))
	if p == nil {
		return 0
	}
	return uintptr(bookSort(bookClass(uint32(uintptr(unsafe.Pointer(p))))))
}
