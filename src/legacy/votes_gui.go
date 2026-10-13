package legacy

import (
	"encoding/binary"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

var voteTopic uint32
var voteWindow, votePlayersWindow, voteTopicsWindow *gui.Window
var voteNameCount, votePreviousNameCount, voteChoice, votePreviousChoice uint32
var voteNames, votePreviousNames [32][28]uint16

func voteGUIInit() int {
	voteWindow = Nox_new_window_from_file("GuiKick.wnd", voteGUIProc)
	if voteWindow == nil {
		return 0
	}
	votePlayersWindow = voteWindow.ChildByID(4320)
	voteTopicsWindow = voteWindow.ChildByID(4321)
	voteWindow.SetPos(image.Pt((int(nox_win_width)-voteWindow.SizeVal.X)/2, voteWindow.Offs().Y))
	voteWindow.SetHidden(true)
	voteNameCount = 0
	votePreviousNameCount = 0
	voteChoice = 0
	votePreviousChoice = 0
	return 1
}
func voteGUIHide() int {
	if voteWindow == nil || voteWindow.GetFlags().IsHidden() {
		return 0
	}
	voteWindow.StackPop()
	voteWindow.SetHidden(true)
	return 1
}
func voteGUIText(key string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(key), "GUIVote.c")
}
func voteGUIAdd(w *gui.Window, text string, color uint32) {
	serverOptionsSetText(w, 16397, text, int(color))
}
func voteGUIShow(topic uint32) {
	teamUIEvent(votePlayersWindow, 16399, 0, 0)
	teamUIEvent(voteTopicsWindow, 16399, 0, 0)
	voteTopic = topic
	switch topic {
	case 4:
		if pl := Get_dword_8531A0_2576(); pl != nil && pl.Field4792 == 0 {
			Nox_xxx_dialogMsgBoxCreate_449A10(nil, voteGUIText("guiquit.c:Vote"), voteGUIText("GUIVote.c:NotAllowedVote"), 33, nil, nil)
			return
		}
		votePlayersWindow.SetHidden(true)
		voteTopicsWindow.SetHidden(false)
		serverOptionsSetText(voteWindow.ChildByID(4301), 16385, voteGUIText("SelectVoteTopic"), 0)
		voteGUIAdd(voteTopicsWindow, voteGUIText("VoteTopicLabel")+" "+voteGUIText("VoteResetServer"), 4)
		voteGUIAdd(voteTopicsWindow, voteGUIText("VoteTopicLabel")+" "+voteGUIText("VoteKickPlayer"), 4)
	case 2:
		votePlayersWindow.SetHidden(true)
		voteTopicsWindow.SetHidden(false)
		serverOptionsSetText(voteWindow.ChildByID(4301), 16385, voteGUIText("Vote:ResetQuest"), 0)
		voteGUIAdd(voteTopicsWindow, voteGUIText("WindowDir:Yes"), 4)
		voteGUIAdd(voteTopicsWindow, voteGUIText("WindowDir:No"), 4)
		selected := uintptr(1)
		if voteChoice == 1 {
			selected = 0
		}
		teamUIEvent(voteTopicsWindow, 16403, selected, 0)
	case 0, 1, 3:
		votePlayersWindow.SetHidden(false)
		voteTopicsWindow.SetHidden(true)
		serverOptionsSetText(voteWindow.ChildByID(4301), 16385, voteGUIText("VoteKickPlayer"), 0)
		s := GetServer().S()
		local := Get_dword_8531A0_2576()
		index := 0
		var teamID byte
		color := uint32(4)
		if s.Teams.Count() != 0 {
			member := teamRuntimeObject(ClientPlayerNetCode())
			if member == nil {
				return
			}
			team := s.Teams.ByID(member.ID)
			if team == nil {
				voteWindow.ShowModal()
				return
			}
			teamID = byte(team.ID())
			color = memmap.Uint32(0x587000, 156400+8*uintptr(teamID%10))
		}
		for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
			if pl == local {
				continue
			}
			if s.Teams.Count() != 0 && !teamRuntimeContains(teamRuntimeObject(int(pl.NetCodeVal)), server.TeamID(teamID)) {
				continue
			}
			teamUIEvent(votePlayersWindow, 16397, uintptr(unsafe.Pointer(&pl.NameFinal[0])), uintptr(color))
			for j := uint32(0); j < voteNameCount && j < 32; j++ {
				if voteNamesEqual(&voteNames[j][0], &pl.NameFinal[0]) {
					teamUIEvent(votePlayersWindow, 16405, uintptr(index), 0)
				}
			}
			index++
		}
	default:
		return
	}
	voteWindow.ShowModal()
}
func voteSendName(name *uint16, withdraw bool) {
	if name == nil {
		return
	}
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		if !voteNamesEqual(&p.NameFinal[0], name) {
			continue
		}
		var msg [52]byte
		msg[0] = 238
		if withdraw {
			msg[1] = 2
		}
		for i := 0; i < 24; i++ {
			v := *(*uint16)(unsafe.Add(unsafe.Pointer(name), 2*i))
			if v == 0 {
				break
			}
			binary.LittleEndian.PutUint16(msg[2+2*i:], v)
		}
		GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])
		return
	}
}
func voteGUISelection() int {
	votePreviousNameCount = voteNameCount
	votePreviousNames = voteNames
	voteNameCount = 0
	selected := unsafe.Pointer(uintptr(teamUIEvent(votePlayersWindow, 16404, 0, 0))) // multi-select: the selected indexes
	if selected != nil {
		for i := 0; i < 32; i++ {
			index := *(*int32)(unsafe.Add(selected, 4*i))
			if index == -1 {
				break
			}
			text := (*uint16)(unsafe.Pointer(uintptr(teamUIEvent(votePlayersWindow, 16406, uintptr(index), 0))))
			if text == nil {
				continue
			}
			dst := &voteNames[voteNameCount]
			clear(dst[:])
			for j := 0; j < len(dst)-1; j++ {
				dst[j] = *(*uint16)(unsafe.Add(unsafe.Pointer(text), 2*j))
				if dst[j] == 0 {
					break
				}
			}
			voteNameCount++
		}
	}
	for i := uint32(0); i < votePreviousNameCount && i < 32; i++ {
		found := false
		for j := uint32(0); j < voteNameCount; j++ {
			if voteNamesEqual(&votePreviousNames[i][0], &voteNames[j][0]) {
				found = true
				break
			}
		}
		if !found {
			voteSendName(&votePreviousNames[i][0], true)
		}
	}
	for i := uint32(0); i < voteNameCount; i++ {
		found := false
		for j := uint32(0); j < votePreviousNameCount && j < 32; j++ {
			if voteNamesEqual(&voteNames[i][0], &votePreviousNames[j][0]) {
				found = true
				break
			}
		}
		if !found {
			voteSendName(&voteNames[i][0], false)
		}
	}
	return int(voteNameCount)
}
func voteGUIResetChoice() int {
	result := 0
	if teamUIEvent(voteTopicsWindow, 16404, 0, 0) != 0 {
		voteChoice = 0
		if votePreviousChoice == 1 {
			result = reliableClientSend(31, []byte{238, 5}, nil, 1)
			votePreviousChoice = voteChoice
			return result
		}
	} else {
		result = 1
		voteChoice = 1
		if votePreviousChoice == 0 {
			result = reliableClientSend(31, []byte{238, 4}, nil, 1)
			votePreviousChoice = voteChoice
			return result
		}
	}
	votePreviousChoice = uint32(result)
	return result
}
func voteGUIConfirmTopic() {
	switch teamUIEvent(voteTopicsWindow, 16404, 0, 0) {
	case 0:
		voteGUIShow(2)
	case 1:
		voteGUIShow(3)
	}
}
func voteGUIProc(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() != 16391 {
		return gui.RawEventResp(0)
	}
	a, _ := ev.EventArgsC()
	child := (*gui.Window)(unsafe.Pointer(a))
	switch child.ID() {
	case 4311:
		if voteTopic == 4 {
			voteGUIConfirmTopic()
		} else if voteTopic == 2 {
			voteGUIResetChoice()
			voteGUIHide()
		} else if voteTopic == 0 || voteTopic == 1 || voteTopic == 3 {
			voteGUISelection()
			voteGUIHide()
		}
	case 4312:
		voteGUIHide()
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
	return gui.RawEventResp(1)
}

func voteGUIClose() int {
	voteWindow.StackPop()
	voteWindow.Capture(false)
	voteWindow.Destroy()
	voteWindow, votePlayersWindow, voteTopicsWindow = nil, nil, nil
	voteNameCount, votePreviousNameCount, voteChoice, votePreviousChoice = 0, 0, 0, 0
	return 0
}
func voteGUIReset() int { voteChoice, votePreviousChoice = 0, 0; return 0 }
