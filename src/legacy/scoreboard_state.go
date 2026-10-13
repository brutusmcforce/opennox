package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// These tables remain shared with existing callers; fixed-width fields preserve
// the original 80-byte player and 56-byte team records on every Go architecture.
type scoreboardPlayer struct {
	Name                 [26]uint16
	Team                 int32
	Class                byte
	_                    [3]byte
	NetCode, Score, Ping uint32
	Objective            int32
	Flags                uint32
}
type scoreboardTeam struct {
	Name  [22]uint16
	ID    uint32
	Score int32
	Color byte
	_     [3]byte
}

var scoreboardData = struct {
	dirty, requested, localRow, width, height, mode, localCode *uint32
	parent, rank, time, limit                                  **gui.Window
}{
	(*uint32)(unsafe.Pointer(&dword_587000_145664)),
	(*uint32)(unsafe.Pointer(&dword_587000_145668)),
	(*uint32)(unsafe.Pointer(&dword_587000_145672)),
	(*uint32)(unsafe.Pointer(&dword_5d4594_1090040)),
	(*uint32)(unsafe.Pointer(&dword_5d4594_1090044)),
	(*uint32)(unsafe.Pointer(&dword_5d4594_1090120)),
	(*uint32)(unsafe.Pointer(&nox_player_netCode_85319C)),
	(**gui.Window)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1090048)),
	&legacyGlobals.dword_5d4594_1090100,
	&dword_5d4594_1090108,
	&dword_5d4594_1090112,
}

func scoreboardPlayers() []scoreboardPlayer {
	return unsafe.Slice((*scoreboardPlayer)(memmap.PtrOff(0x5D4594, 1084132)), 32)
}
func scoreboardTeams() []scoreboardTeam {
	return unsafe.Slice((*scoreboardTeam)(memmap.PtrOff(0x5D4594, 1087204)), 9)
}
func scoreboardCount(off uintptr) *byte { return memmap.PtrUint8(0x5D4594, off) }
func scoreboardColumn(side, column int) *gui.Window {
	return scoreboardColumns[column][side]
}

func scoreboardHidden(w *gui.Window) bool { return w == nil || w.GetFlags().IsHidden() }
func scoreboardVisible() bool {
	return *scoreboardData.mode != 0 && !scoreboardHidden(*scoreboardData.parent)
}
func scoreboardEnabled() bool { return *scoreboardData.parent != nil && *scoreboardData.mode != 0 }
func scoreboardOpen() {
	w := *scoreboardData.parent
	if w == nil {
		return
	}
	if scoreboardHidden(w) {
		w.Show()
	}
	*scoreboardData.mode = 0
	Sub_4703F0()
}
func scoreboardSetFlag(action, team byte, code uint16) byte {
	result := team
	if team == 1 {
		result = action
		if action == 0 || action == 2 {
			*memmap.PtrUint16(0x5D4594, 1090128) = 0
		} else if action == 1 {
			result = byte(code)
			*memmap.PtrUint16(0x5D4594, 1090128) = code
		}
	} else if team == 2 {
		result = action
		if action == 0 || action == 2 {
			*memmap.PtrUint16(0x5D4594, 1090130) = 0
		} else if action == 1 {
			*memmap.PtrUint16(0x5D4594, 1090130) = code
		}
	}
	return result
}
func scoreboardSetBall(action byte, code uint16) byte {
	result := action
	if action == 0 || action == 1 {
		*memmap.PtrUint16(0x5D4594, 1090132) = 0
	} else if action == 2 || action == 4 {
		result = byte(code)
		*memmap.PtrUint16(0x5D4594, 1090132) = code
	}
	return result
}
func scoreboardClearFlags() int {
	*memmap.PtrUint16(0x5D4594, 1090128) = 0
	*memmap.PtrUint16(0x5D4594, 1090130) = 0
	*memmap.PtrUint16(0x5D4594, 1090132) = 0
	return 0
}
func scoreboardHasTeam(id uint32) bool {
	for _, r := range scoreboardTeams()[:int(*scoreboardCount(1090116))] {
		if r.ID == id {
			return true
		}
	}
	return false
}
func scoreboardHasPlayer(id uint32) bool {
	for _, r := range scoreboardPlayers()[:int(*scoreboardCount(1090117))] {
		if r.NetCode == id {
			return true
		}
	}
	return false
}
func scoreboardTeamIndex(id uint32) byte {
	for i, r := range scoreboardTeams()[:int(*scoreboardCount(1090116))] {
		if r.ID == id {
			return byte(i)
		}
	}
	return 0
}
func scoreboardTeamColor(index byte) byte {
	return memmap.Uint8(0x587000, 145584+uintptr(scoreboardTeams()[index].Color%10)*8)
}
func scoreboardLocalRank() byte {
	rows := scoreboardPlayers()[:int(*scoreboardCount(1090117))]
	var score uint32
	found := false
	for _, r := range rows {
		if r.NetCode == *scoreboardData.localCode {
			score = r.Score
			found = true
			break
		}
	}
	if !found {
		return 0
	}
	rank := byte(1)
	low := noxflags.HasGame(noxflags.GameModeElimination)
	for _, r := range rows {
		if (!low && r.Score > score) || (low && r.Score < score) {
			rank++
		}
	}
	return rank
}
func scoreboardTeamRank(score uint32) byte {
	rank := byte(1)
	low := noxflags.HasGame(noxflags.GameModeElimination)
	for _, r := range scoreboardTeams()[:int(*scoreboardCount(1090116))] {
		v := uint32(r.Score)
		if (!low && v > score) || (low && v < score) {
			rank++
		}
	}
	return rank
}
func scoreboardObjective(p *server.Player) int {
	code := p.NetCodeVal
	if noxflags.HasGame(noxflags.GameModeCTF) {
		if code == uint32(memmap.Uint16(0x5D4594, 1090128)) {
			return 2
		}
		if code == uint32(memmap.Uint16(0x5D4594, 1090130)) {
			return 3
		}
	} else if noxflags.HasGame(noxflags.GameModeFlagBall) {
		if code == uint32(memmap.Uint16(0x5D4594, 1090132)) {
			return 4
		}
	} else if noxflags.HasGame(noxflags.GameModeKOTR) {
		if dr := GetClient().Cli().Objs.ByNetCodeDynamic(int(code)); dr != nil && dr.HasEnchant(30) {
			return 1
		}
	}
	return 0
}
func scoreboardText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "guirank.c")
}
func scoreboardLoadClasses() *uint16 {
	var result *uint16
	for i := uintptr(0); i < 3; i++ {
		id := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 145676+4*i)))
		result = alloc.InternCString16(scoreboardText(id))
		scoreboardClassNames[i] = result
	}
	return result
}
func scoreboardClipName(p *uint16) *uint16 {
	width := int32(memmap.Uint32(0x5D4594, 1084036))
	if p == nil || width == 0 {
		return nil
	}
	s := alloc.GoString16(p)
	r := GetClient().R2()
	size := r.GetStringSizeWrapped(nil, s, 0).X
	if size == 0 {
		return nil
	}
	if int32(size)+5 > width {
		n := 0
		for *(*uint16)(unsafe.Add(unsafe.Pointer(p), n*2)) != 0 {
			n++
		}
		for {
			*(*uint16)(unsafe.Add(unsafe.Pointer(p), n*2)) = 0
			n--
			if int32(r.GetStringSizeWrapped(nil, alloc.GoString16(p), 0).X)+5 <= width {
				break
			}
		}
	}
	return p
}

func scoreboardScreenWidth() int   { return int(nox_win_width) }
func scoreboardYellow() uint32     { return Get_nox_color_yellow_2589772() }
func scoreboardWhite() uint32      { return Get_nox_color_white_2523948() }
func scoreboardTitleColor() uint32 { return Get_dword_8531A0_2572() }

// Compile-time layout checks for the tables shared with legacy storage.
var (
	_ [4]byte  = [unsafe.Alignof(scoreboardPlayer{})]byte{}
	_ [4]byte  = [unsafe.Alignof(scoreboardTeam{})]byte{}
	_ [80]byte = [unsafe.Sizeof(scoreboardPlayer{})]byte{}
	_ [56]byte = [unsafe.Sizeof(scoreboardTeam{})]byte{}
	_ [52]byte = [unsafe.Offsetof(scoreboardPlayer{}.Team)]byte{}
	_ [56]byte = [unsafe.Offsetof(scoreboardPlayer{}.Class)]byte{}
	_ [60]byte = [unsafe.Offsetof(scoreboardPlayer{}.NetCode)]byte{}
	_ [64]byte = [unsafe.Offsetof(scoreboardPlayer{}.Score)]byte{}
	_ [68]byte = [unsafe.Offsetof(scoreboardPlayer{}.Ping)]byte{}
	_ [72]byte = [unsafe.Offsetof(scoreboardPlayer{}.Objective)]byte{}
	_ [76]byte = [unsafe.Offsetof(scoreboardPlayer{}.Flags)]byte{}
	_ [44]byte = [unsafe.Offsetof(scoreboardTeam{}.ID)]byte{}
	_ [48]byte = [unsafe.Offsetof(scoreboardTeam{}.Score)]byte{}
	_ [52]byte = [unsafe.Offsetof(scoreboardTeam{}.Color)]byte{}
)
