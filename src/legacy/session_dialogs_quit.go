package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

var sessionQuitRoot *gui.Window

func sessionQuitShown() int {
	return bool2int(sessionQuitRoot != nil && sessionQuitRoot.Flags&gui.StatusHidden == 0)
}
func sessionQuitCapture() int {
	if sessionQuitRoot.Capture(true) {
		return 0
	}
	return -4
}
func sessionQuitHide() {
	if sessionQuitRoot != nil && sessionQuitRoot.Flags&gui.StatusHidden == 0 {
		sessionQuitToggle()
	}
}
func sessionQuitColors() *gui.Window {
	if sessionQuitRoot != nil {
		sessionQuitRoot.DrawData().BgColorVal = uint32(nox_color_black_2650656)
	}
	var last *gui.Window
	for id := uint(9001); id <= 9006; id++ {
		last = sessionQuitRoot.ChildByID(id)
		if last != nil {
			last.DrawData().TextColorVal = uint32(nox_color_orange_2614256)
		}
	}
	return last
}
func sessionQuitText(key string) string { return serverPanelsText("guiquit.c", key) }
func sessionQuitToggle() {
	w := sessionQuitRoot
	if w == nil {
		return
	}
	if w.Flags&gui.StatusHidden == 0 {
		GetClient().Cli().GUI.Focus(nil)
		w.Capture(false)
		w.Hide()
		w.Flags &^= gui.StatusEnabled
		Sub_413A00(0)
		return
	}
	player := LocalPlayerDrawable.C()
	solo := noxflags.HasGame(noxflags.GameModeCoop)
	if player != nil && solo {
		state := *(*uint32)(unsafe.Add(player, 276))
		if state == 1 || state == 2 || state == 51 {
			return
		}
	}
	if *bookWord(1047520) == 1 || noxflags.HasGame(noxflags.GamePause) {
		return
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
	w.ShowModal()
	w.Flags |= gui.StatusEnabled
	w.Capture(true)
	if solo {
		serverOptionsSetText(w.ChildByID(9003), 16385, sessionQuitText("SoloSaveLabel"), 0)
		serverPanelsHide(w, 9001, 9002, false)
		serverPanelsHide(w, 9007, 9009, true)
		w.ChildByID(9004).SetPos(w.ChildByID(9009).Off)
		Sub_413A00(1)
		uiWindowResize(w, 220, 285)
	} else {
		serverOptionsSetText(w.ChildByID(9003), 16385, sessionQuitText("MultiplayerSaveLabel"), 0)
		serverPanelsHide(w, 9001, 9002, true)
		serverPanelsHide(w, 9007, 9008, false)
		uiWindowEnable(w.ChildByID(9007), 1)
		if vote := w.ChildByID(9009); vote != nil {
			vote.Show()
			optionsSend(vote, 16385, uintptr(memmap.PtrOff(0x5D4594, 825772)), 0)
			uiWindowEnable(vote, bool2int(!noxflags.HasGame(49152) && uint8(GetServer().S().Teams.Count()) != 0))
			w.ChildByID(9004).SetPos(vote.Off.Add(image.Pt(0, 45)))
		}
		uiWindowResize(w, 220, 330)
		if noxflags.HasGame(noxflags.GameModeQuest) {
			uiWindowEnable(w.ChildByID(9007), 0)
			uiWindowEnable(w.ChildByID(9003), 0)
		}
		if noxflags.HasEngine(noxflags.EngineNoRendering) {
			for _, id := range []uint{9007, 9005, 9003} {
				uiWindowEnable(w.ChildByID(id), 0)
			}
		}
	}
}

// SessionQuitEvent is the real quit-menu callback, used directly by the Go loader.
func SessionQuitEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() != 16391 {
		return nil
	}
	a, _ := ev.EventArgsC()
	button := (*gui.Window)(unsafe.Pointer(a))
	id := button.ID()
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	switch id {
	case 9001:
		sessionQuitToggle()
		Sub_413A00(1)
		if !noxflags.HasGame(noxflags.GameModeCoop) || Nox_xxx_playerAnimCheck_4372B0() != 0 {
			Sub_445B40()
		} else {
			Nox_xxx_dialogMsgBoxCreate_449A10(nil, sessionQuitText("SelChar.c:LoadLabel"), sessionQuitText("GUIQuit.c:ReallyLoadMessage"), 56, func() { Sub_445B40() }, func() { Sub_413A00(0) })
		}
	case 9002:
		p := LocalPlayerDrawable.C()
		if *(*uint32)(unsafe.Add(p, 120))&0x8000 == 0 {
			sessionQuitToggle()
			if noxflags.HasGame(noxflags.GameModeCoop) {
				Nox_setSaveFileName_4DB130("AUTOSAVE")
				Sub_4DB170(true, nil, 0)
			}
		} else {
			uiWindowEnable(w.ChildByID(id), 0)
		}
	case 9003:
		sessionQuitToggle()
		if noxflags.HasGame(noxflags.GameModeCoop) {
			Nox_savegame_sub_46D580()
		} else {
			playerFileSaveRequest()
		}
		if Sub_43C6E0() == 0 {
			Sub_43CF70()
		}
	case 9004:
		sessionQuitRoot.Capture(false)
		Nox_xxx_dialogMsgBoxCreate_449A10(sessionQuitRoot, sessionQuitText("GUIQuit.c:ReallyQuitTitle"), sessionQuitText("GUIQuit.c:ReallyQuitMessage"), 56, func() { Nox_client_quit_4460C0(); sessionQuitToggle() }, func() { sessionQuitCapture() })
	case 9005:
		sessionQuitToggle()
		optionsShow()
	case 9006:
		sessionQuitToggle()
	case 9007:
		if noxflags.HasGame(noxflags.GameHost) {
			consoleCommandRemote(GetServer().S().Players.ByID(int(nox_player_netCode_85319C)), 0, "")
		} else {
			Nox_xxx_netServerCmd_440950(0, "")
		}
		sessionQuitToggle()
	case 9008:
		sessionQuitToggle()
		serverOptionsConstruct()
	case 9009:
		sessionQuitToggle()
		topic := uint32(0)
		if noxflags.HasGame(noxflags.GameModeQuest) {
			topic = 4
		}
		voteGUIShow(topic)
	}
	button.DrawData().Field0 &^= 2
	return nil
}
