package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"unsafe"
)

func serverOptionsTabOrder(id int) int {
	w := serverOptionsChild(uint(id))
	parent := w.Parent()
	var order []uint
	switch id {
	case 10161:
		order = []uint{10163, 10162}
	case 10162:
		order = []uint{10163, 10161}
	case 10163:
		order = []uint{10162, 10161}
	}
	for _, id := range order {
		c := serverOptionsChild(id)
		c.SetParent(nil)
		c.SetParent(parent)
	}
	w.SetParent(nil)
	return w.SetParent(parent)
}
func serverOptionsTab(tab int) int {
	general, access, players := serverOptionsChild(10163), serverOptionsChild(10161), serverOptionsChild(10162)
	setRaw := func(w *gui.Window, v uint32, disabled bool) {
		if disabled {
			w.DrawData().DisImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(v)))
		} else {
			w.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(v)))
		}
	}
	switch tab {
	case 1:
		setRaw(general, serverOptionsTabs3, false)
		setRaw(players, serverOptionsTabs2, false)
		if serverOptionsPlayersPanel != 0 {
			teamUIPlayersDestroy(true)
			serverOptionsPlayersPanel = 0
		}
		if serverOptionsGeneralPanel != 0 {
			serverPanelsGeneralClose()
			serverOptionsGeneralPanel = 0
		}
		serverOptionsAccessPanel = uint32(serverPanelsAccessOpen(serverOptionsWindow(1046492)))
		return serverOptionsTabOrder(10161)
	case 2:
		setRaw(general, serverOptionsTabs3, false)
		setRaw(access, serverOptionsTabs3, !noxflags.HasGame(1))
		if serverOptionsAccessPanel != 0 {
			serverPanelsAccessClose(true)
			serverOptionsAccessPanel = 0
		}
		if serverOptionsGeneralPanel != 0 {
			serverPanelsGeneralClose()
			serverOptionsGeneralPanel = 0
		}
		serverOptionsPlayersPanel = uint32(teamUIPlayersConstruct(serverOptionsWindow(1046492)))
		return serverOptionsTabOrder(10162)
	case 0:
		setRaw(access, serverOptionsTabs2, !noxflags.HasGame(1))
		setRaw(players, serverOptionsTabs2, false)
		if serverOptionsPlayersPanel != 0 {
			teamUIPlayersDestroy(true)
			serverOptionsPlayersPanel = 0
		} else if serverOptionsAccessPanel != 0 {
			serverPanelsAccessClose(true)
			serverOptionsAccessPanel = 0
		}
		serverOptionsGeneralPanel = uint32(serverPanelsGeneralOpen(serverOptionsWindow(1046492)))
		return serverOptionsTabOrder(10163)
	default:
		return tab - 2
	}
}
