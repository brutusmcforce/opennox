package legacy

import (
	"encoding/binary"
	"image"
	"net/netip"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func browserLabel(key string) {
	// Static labels retain their text pointer; the typed event interns this fixed
	// localized string for the widget's lifetime.
	browserUI.label.Func94(&gui.StaticTextSetText{Str: browserString(key)})
}
func browserShowRegion() {
	w := browserUI.mapWindow
	w.SetHidden(false)
	browserUI.overview.SetHidden(true)
	if r := browserUI.region; r >= 0 && r < int32(len(browserRegionMaps)) {
		w.DrawData().BgImageHnd = browserRegionMaps[r]
	}
	key := "JoinServer"
	if browserUI.creating == 1 {
		key = "CreateMsg"
	}
	browserLabel(key)
	sub_49FDB0(int32(browserUI.region))
}
func browserEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	code := ev.EventCode()
	a, b := ev.EventArgsC()
	switch code {
	case 16403, 16412:
		browserColumnsSend(code, a, 0)
		return nil
	case 16400:
		if windowArg(a).ID() == 10061 {
			p := GetClient().GetMousePos()
			browserUI.selected = browserPopupAt(int32(b))
			point := [2]uint32{uint32(p.X), uint32(p.Y)}
			browserInfoPopup(&point, browserUI.selected)
		}
		return nil
	case 23:
		return gui.RawEventResp(1)
	case 16384:
		id := windowArg(a).ID()
		if id >= 10043 && id <= 10044 {
			browserColumnsSend(code, a, 0)
		}
		return nil
	}
	// Resource construction events carry numeric child IDs, not window pointers.
	if code != 16391 {
		return nil
	}
	id := windowArg(a).ID()
	if id != 10043 && id != 10044 && id != 10035 && id != 10036 {
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	}
	root := browserUI.world
	if id >= 10070 {
		p := GetClient().GetMousePos()
		point := [2]uint32{uint32(p.X - 216), uint32(p.Y - 27)}
		head := browserListHead()
		if browserCount(&point, head) >= 2 {
			browserUI.popupWindow = browserPopup(root, &point, head)
		} else {
			browserUI.selected = browserListID(int32(id - 10070))
			point = [2]uint32{uint32(p.X), uint32(p.Y)}
			browserInfoPopup(&point, browserUI.selected)
		}
		Nox_xxx_cursorSetTooltip_4776B0(alloc.GoString16(memmap.PtrUint16(0x5D4594, 815112)))
	}
	switch id {
	case 10007:
		if browserUI.creating != 0 {
			Sub_4373A0()
		}
	case 10010:
		Sub_4373A0()
	case 10047, 10048, 10049, 10050, 10051:
		browserListReset()
		nox_wol_servers_sortBtnHandler_4A0290(int32(id))
		browserListResort()
	case 10054, 10055, 10056, 10057:
		browserUI.region = int32(id - 10054)
		browserShowRegion()
		Nox_client_refreshServerList_4378B0()
		x, y := uint32(408), uint32(239)
		if noxflags.HasGame(0x2000000) {
			x = memmap.Uint32(0x5D4594, 1308732) + 216
			y = memmap.Uint32(0x5D4594, 1308736) + 27
		}
		browserCreateAt(x, y)
	case 10006:
		browserUI.creating = 0
		Nox_client_refreshServerList_4378B0()
	case 10004:
		browserUI.creating = 0
		browserShowList()
		Nox_client_refreshServerList_4378B0()
	case 10005:
		browserUI.creating = 0
		browserUI.gameList.SetHidden(true)
		browserUI.mapWindow.SetHidden(true)
		browserUI.overview.SetHidden(true)
		browserUI.filter.SetHidden(false)
		serverPanelsEnable(root, 10006, 10007, 0)
		serverPanelsHide(root, 10047, 10051, true)
		browserLabel("FilterMsg")
	case 10003:
		browserUI.creating = 1
		noxflags.UnsetGame(0x10000)
		GetClient().ChangeMousePos(image.Pt(408, 239), true)
		browserMarkersEnable(0)
		if !browserHidden(browserUI.filter) {
			sessionFilterSave()
		}
		questRuntimeSetWord(1556160, 1)
		nox_xxx_cliShowHideTubes_470AA0(1)
		browserUI.region = 0
		root.ChildByID(10020).Func93(&gui.RawEvent{Event: 5, Arg1: 15663512})
	case 10002:
		if noxflags.HasGame(0x1000000) {
			return nil
		}
		browserUI.creating = 1
		noxflags.SetGame(0x10000)
		nox_xxx_cliShowHideTubes_470AA0(0)
		GetClient().ChangeMousePos(image.Pt(408, 239), true)
		browserMarkersEnable(0)
		if !browserHidden(browserUI.filter) {
			sessionFilterSave()
		}
		if noxflags.HasGame(0x2000000) || Sub_4D6F30() != 0 {
			return nil
		}
		browserUI.region = 0
		root.ChildByID(10020).Func93(&gui.RawEvent{Event: 5, Arg1: 15663512})
	case 4001:
		switch browserUI.connectionState {
		case 6:
			text := (*uint16)(unsafe.Pointer(uintptr(Sub_449E60(4))))
			var buf [22]byte
			// The original stack buffer left header padding undefined; initialize it.
			for i := 0; i < 8; i++ {
				v := *(*uint16)(unsafe.Add(unsafe.Pointer(text), i*2))
				if v == 0 {
					break
				}
				binary.LittleEndian.PutUint16(buf[4+i*2:], v)
			}
			SendXXX_5550D0(netip.AddrPortFrom(int2ip(uint32(nox_client_getServerAddr_43B300())), uint16(nox_client_getServerPort_43B320())), buf[:])
			browserUI.connectionState = 3
			browserUI.connectionDeadline = uint64(uint32(PlatformTicks()) + 20000)
			Sub_449EA0(0)
		case 10:
			Sub_449E60(4)
			Sub_449E30(browserString("Finding"))
			browserUI.connectionState = 11
			Sub_449EA0(0)
		case 1:
			browserConnectionReset()
			nox_game_showGameSel_4379F0()
		default:
			if browserUI.retry != 0 {
				browserUI.refreshDeadline = uint64(uint32(PlatformTicks()) + 1000)
			}
		}
	case 4002:
		browserConnectionReset()
		nox_game_showGameSel_4379F0()
	case 10001:
		if browserUI.transition != 0 {
			return nil
		}
		if sub_43B340()&0x1000 != 0 {
			questRuntimeSetWord(1556164, 1)
			nox_xxx_cliShowHideTubes_470AA0(1)
		}
		browserChooseCharacter()
		count := 0
		if questRuntimeWord(1556164) != 0 {
			count = Nox_client_countPlayerFiles04_4DC7D0()
		} else {
			count = int(sessionCharacterCount())
		}
		if count != 0 {
			Sub_4A7A70(1)
			Nox_game_showSelChar_4A4DB0()
		} else {
			Sub_4A7A70(0)
			characterShowClass()
		}
		browserUI.transition = 1
		browserPopupClose()
	}
	return nil
}
