package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func browserWireList(list, slider, up, down *gui.Window, height int) {
	slider.DrawData().Window = list
	up.DrawData().Window = list
	down.DrawData().Window = list
	d := uiListData(list)
	d.Field_9 = slider.C()
	d.Field_7 = up.C()
	d.Field_8 = down.C()
	slider.Field100Ptr.SizeVal = image.Pt(16, height)
}
func browserShow() int {
	questRuntimeSetWord(1556160, 0)
	questRuntimeSetWord(1556164, 0)
	if !noxflags.HasGame(0x2000000) && Sub_4D6F30() == 0 {
		browserUI.creating = 0
		browserUI.region = -1
	}
	*memmap.PtrUint64(0x5D4594, 815076) = 0
	*memmap.PtrUint32(0x5D4594, 815084) = 0
	GetClient().GameAddStateCode(10000)
	Sub_4A24C0(1)
	Sub_4A1BE0(0)
	if browserUI.region != -1 {
		browserMapPolygons()
	}
	root := browserUI.world
	if root != nil {
		browserUI.transition = 0
		anim := (*gui.Anim)(unsafe.Pointer(browserUI.animation))
		anim.SetState(gui.AnimIn)
		anim.FncDoneOutPtr = animationCallbackKey(animationKeyBrowserDone)
		Sub_43BE40(3)
		Nox_xxx_clientPlaySoundSpecial_452D80(922, 100)
		root.SetHidden(false)
		browserUI.detailPanel.SetHidden(false)
		uiWindowEnable(browserUI.mapWindow, 1)
		if browserUI.creating != 0 {
			browserMarkersEnable(0)
		}
		serverPanelsHide(root, 10047, 10051, browserUI.listMode == 0)
		if questRuntimeWord(1556104) == 2 {
			for i := 0; i < 2; i++ {
				optionsSend(root, 16391, uintptr(root.ChildByID(10010).C()), 0)
			}
		}
		return 1
	}
	root = Nox_new_window_from_file("noxworld.wnd", browserEvent)
	browserUI.world = root
	if root == nil {
		return 0
	}
	browserMapPolygonsClear()
	root.ShowModal()
	root.SetAllFuncs(browserMapInput, nil, nil)
	anim := Nox_gui_makeAnimation_43C5B0(root, 0, 0, 0, -480, 0, 20, 0, -40)
	browserUI.animation = anim
	if anim == nil {
		return 0
	}
	anim.StateID = 10000
	anim.Func12Ptr = animationCallbackKey(animationKeyBrowserOut)
	anim.FncDoneOutPtr = animationCallbackKey(animationKeyBrowserDone)
	mapWin, overview := root.ChildByID(10020), root.ChildByID(10021)
	browserUI.mapWindow = mapWin
	browserUI.overview = overview
	mapWin.SetFunc93(browserMapInput)
	overview.SetFunc93(browserMapInput)
	overview.SetFunc94(browserEvent)
	browserUI.label = root.ChildByID(10011)
	main := root.ChildByID(10037)
	browserUI.gameList = main
	details := root.ChildByID(10034)
	browserUI.detailList = details
	panel := root.ChildByID(10033)
	browserUI.detailPanel = panel
	browserUI.filter = sessionFilterOpen(root)
	browserUI.joinButton = root.ChildByID(10001)
	for i, p := range []**gui.Window{&browserUI.playersColumn, &browserUI.modeColumn, &browserUI.mapColumn, &browserUI.pingColumn, &browserUI.statusColumn} {
		*p = root.ChildByID(uint(10038 + i))
	}
	if noxflags.HasGame(0x1000000) {
		uiWindowEnable(root.ChildByID(10002), 0)
	}
	root.ChildByID(10003).DrawData().Field0 &^= 4
	root.ChildByID(10046).SetDraw(browserMouseDraw)
	main.SetFunc94(browserEvent)
	main.SetHidden(true)
	panel.SetHidden(true)
	root.ChildByID(10001).DrawData().Window = root
	browserUI.region = 0
	overview.SetHidden(true)
	browserLabel("JoinServer")
	serverPanelsHide(mapWin, 10620, 10631, true)
	browserMapPolygons()
	for i, name := range []string{"NWGameIconLargeGreen", "NWGameIconLargeGreenLit", "NWGameIconSmallGreen", "NWGameIconSmallGreenLit", "NWGameIconLargeYellow", "NWGameIconLargeYellowLit", "NWGameIconSmallYellow", "NWGameIconSmallYellowLit", "NWGameIconLargeRed", "NWGameIconLargeRedLit", "NWGameIconSmallRed", "NWGameIconSmallRedLit"} {
		browserGameIcons[i] = Nox_xxx_gLoadImg(name).C()
	}
	for i, name := range []string{"NWMapULLg", "NWMapURLg", "NWMapLLLg", "NWMapLRLg"} {
		browserRegionMaps[i] = Nox_xxx_gLoadImg(name).C()
	}
	slider := main.ChildByID(10053)
	browserWireList(main, slider, main.ChildByID(10043), main.ChildByID(10044), 12)
	slider.Field100Ptr.DrawData().ImgPtVal = image.Pt(0, -15)
	browserWireList(details, panel.ChildByID(10032), panel.ChildByID(10035), panel.ChildByID(10036), 10)
	browserColumnsInit()
	serverConfigModeStore(uint32(GetServer().S().ServerPort()))
	browserUI.creating = 0
	browserUI.transition = 0
	browserUI.hosting = 0
	browserUI.connectionState = 0
	browserUI.refreshDeadline = uint64(uint32(PlatformTicks()) + 1000)
	if browserUI.listMode == 1 {
		browserShowList()
	}
	Nox_xxx_createSocketLocal(0)
	browserUI.resultCount = 0
	Nox_xxx_loadModifyers_4158C0()
	Nox_xxx_loadLook_415D50()
	// video_highres.go already enabled these limits in every build profile.
	Sub_430C30_set_video_max(3840, 2160)
	nox_client_setCursorType_477610(0)
	if browserUI.pendingKicked != 0 {
		browserNotice(false)
	} else if browserUI.pendingTimeout != 0 {
		browserNotice(true)
	}
	GetServer().SetUpdateFunc2(func() bool { return browserTick() != 0 })
	if Sub_44A4A0() != 0 {
		Sub_44A4B0()
	}
	mapWin.SetDraw(browserMapDraw)
	for _, id := range []uint{10054, 10055, 10056, 10057} {
		root.ChildByID(id).SetDraw(browserMapDraw)
	}
	if browserUI.listMode == 0 {
		serverPanelsHide(root, 10047, 10051, true)
	}
	if questRuntimeWord(1556104) == 1 {
		optionsSend(root, 16391, uintptr(root.ChildByID(10002).C()), 0)
		questRuntimeSetWord(1556160, 1)
		mapWin.Func93(&gui.RawEvent{Event: 5, Arg1: 15663512})
	}
	return 1
}
func nox_game_showGameSel_4379F0() int { return browserShow() }
