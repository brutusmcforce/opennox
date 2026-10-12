package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/strman"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func clientSessionQuest(data []byte) int {
	if len(data) < 2 {
		return -1
	}
	size := 0
	switch data[1] {
	case 0, 20:
		size = 2
	case 1, 5, 6, 7, 8, 9, 10, 15, 17, 18, 19, 29:
		size = 4
	case 2:
		size = 14
	case 4, 22, 23, 30, 31, 32:
		size = 5
	case 11:
		size = 16
	case 12:
		size = 90
	case 13, 14:
		size = 69
	case 16:
		size = 12
	case 21:
		size = 8
	case 24, 28:
		size = 3
	case 25:
		size = 7
	case 26:
		size = 6
	case 33:
		size = 52
	default:
		return -1
	}
	if len(data) < size {
		return -1
	}
	if data[1] == 33 && data[51] >= 5 {
		if Nox_client_isConnected() {
			return -1
		}
		return size
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	c := GetClient()
	on := Nox_client_isConnected()
	// These cache updates and the FireBoom sprite also occur while disconnected.
	if data[1] == 16 {
		p := memmap.PtrUint32(0x5D4594, 1200904)
		if *p == 0 {
			*p = uint32(c.Cli().Things.IndByID("GreenZap"))
		}
	}
	if data[1] == 25 {
		spark, boom := memmap.PtrUint32(0x5D4594, 1200908), memmap.PtrUint32(0x5D4594, 1200912)
		if *spark == 0 {
			*spark = uint32(c.Cli().Things.IndByID("GreenSpark"))
			*boom = uint32(c.Cli().Things.IndByID("FireBoom"))
		}
		if dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(int(*boom), image.Pt(int(int16(word(2))), int(int16(word(4))))); dr != nil {
			c.Cli().Objs.List34Add(dr)
		}
	}
	if data[1] >= 12 && data[1] <= 14 {
		defer Set_nox_client_gui_flag_1556112(0)
	}
	if !on {
		return size
	}
	switch data[1] {
	case 0:
		if noxflags.HasGame(1) {
			quickbarClearAbilities()
			quickbarPrepare()
			path := alloc.GoString(memmap.PtrUint8(0x85B3FC, 10984))
			Sub_41CC00(path)
			playerFileClientLoad(path)
		}
		nox_xxx_cliShowHideTubes_470AA0(1)
		voteGUIReset()
		nox_xxx_cliPrepareGameplay2_4721D0()
		scoreboardOpen()
	case 1:
		audioEventPlay(1008, 100, 0, 0)
		if p := GetServer().S().Players.ByID(int(word(2))); p != nil && !noxflags.HasGame(1) {
			*(*uint32)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(4792)) /* Field4792 */)) = 1
		}
	case 2:
		interactionGameOverShow((*uint16)(unsafe.Pointer(&data[0])))
	case 4:
		if p := GetServer().S().Players.ByID(int(word(3))); p != nil {
			*(*byte)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(4816)) /* data4800 */)) = data[2]
		}
		if int(word(3)) == ClientPlayerNetCode() {
			*memmap.PtrUint32(0x5D4594, 1050012) = uint32(data[2])
		}
	case 5, 6, 7, 8, 9, 10, 11:
	case 12:
		briefingWinReport(unsafe.Pointer(&data[0]))
	case 13:
		briefingSelection(unsafe.Pointer(&data[0]), int(data[4]&1), false)
	case 14:
		briefingSelection(unsafe.Pointer(&data[0]), 1, true)
	case 15:
		if dr := c.Cli().Objs.ByNetCodeStatic(int(word(2))); dr != nil {
			*(*byte)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)) = 0
		}
	case 16:
		if dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(int(memmap.Uint32(0x5D4594, 1200904)), image.Pt(int(word(6)), int(word(8)))); dr != nil {
			b := unsafe.Slice((*byte)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)), 13)
			b[0] = 0
			binary.LittleEndian.PutUint32(b[5:], dword(2))
			binary.LittleEndian.PutUint32(b[9:], dword(6))
			binary.LittleEndian.PutUint32(b[1:], uint32(word(10)))
			c.Cli().Objs.TransparentDecay(dr, int(word(10)))
		}
	case 17:
		bookRemoveSpell(int(word(2)))
	case 18:
		bookRemoveAbility(int(word(2)))
	case 19:
		bookRemoveGuide(int(word(2)))
	case 20:
		if !scoreboardVisible() {
			scoreboardOpen()
		}
	case 21:
		if p := GetServer().S().Players.ByID(int(word(6))); p != nil {
			*(*uint32)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(4820)) /* data4800 */)) = dword(2)
		}
	case 22, 23:
		if p := GetServer().S().Players.ByID(int(word(3))); p != nil {
			*(*byte)(unsafe.Add(p.C(), 4824+int(data[1]-22))) = data[2]
		}
	case 24:
		interactionKeyUpdate(uintptr(data[2]))
	case 25:
		effectSparkBurst(image.Pt(int(int16(word(2))), int(int16(word(4)))), int(memmap.Uint32(0x5D4594, 1200908)), data[6])
	case 26:
		effectCreatePointSparks(int(memmap.Uint32(0x5D4594, 1200788)), 25, 500, 25, int(int16(word(2))), int(int16(word(4))))
	case 28:
		onlineSessionBriefing(uint32(data[2]))
	case 29:
		uiInventorySetWindowLevel(int(word(2)))
	case 30, 31, 32:
		p := GetServer().S().Players.ByID(int(word(3)))
		if p == nil {
			break
		}
		var title, key string
		switch data[1] {
		case 30:
			var ok bool
			title, ok = Nox_xxx_spellTitle_424930(int(data[2]))
			if !ok {
				title = "(null)"
			}
			key = "plyrspel.c:AwardSpell"
		case 31:
			ptr := bookGuideCreatureName(int32(data[2]))
			if ptr == 0 {
				title = "(null)"
			} else {
				title = alloc.GoString16((*uint16)(memmap.WordPtr(ptr)))
			}
			key = "PlyrGide.c:AwardGuide"
		case 32:
			title = Nox_xxx_abilityGetName_0_425260(int(data[2]))
			key = "ComAblty.c:AwardAbility"
		}
		if p == Get_dword_8531A0_2576() {
			Nox_xxx_printCentered_445490(clientGameProgressText(key, title))
		} else {
			Nox_xxx_printCentered_445490(clientGameProgressText(key+"ToOther", alloc.GoString16((*uint16)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(4704)) /* NameFinal */))), title))
		}
	case 33:
		offsets := [5]uintptr{160948, 160988, 161028, 161068, 161112}
		if int(data[51]) >= len(offsets) {
			return -1
		}
		key := alloc.GoString(memmap.PtrUint8(0x587000, offsets[data[51]]))
		end := 2
		for end < 51 && data[end] != 0 {
			end++
		}
		sm := GetServer().S().Strings()
		title := sm.GetStringInFile(strman.ID(key), "cdecode.c")
		Nox_xxx_printCentered_445490(clientGameProgressText(string(data[2:end]), title))
	}
	return size
}
