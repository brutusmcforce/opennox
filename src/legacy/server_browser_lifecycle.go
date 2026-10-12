package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func browserMarkersEnable(enabled int) {
	if browserUI.resultCount != 0 {
		serverPanelsEnable(browserUI.world, 10070, uint(browserUI.resultCount)+10069, enabled)
	}
}
func browserListReset() {
	if browserUI.listMode == 1 {
		optionsSend(browserWindow(uint32(browserUI.gameList)), 16399, 0, 0)
	}
}
func browserAnimationFinish() int {
	anim := (*gui.Anim)(unsafe.Pointer(browserUI.animation))
	next := anim.Func13Ptr
	anim.Free()
	if !noxflags.HasGame(0x10000000) {
		browserClose()
	}
	if next != nil {
		gui.CallAnimationCallback(next)
	}
	return 1
}
func browserAnimationOut() int {
	anim := (*gui.Anim)(unsafe.Pointer(browserUI.animation))
	if anim.State() == gui.AnimOutDone {
		return browserAnimationFinish()
	}
	anim.SetState(gui.AnimOut)
	Sub_43BE40(2)
	Nox_xxx_clientPlaySoundSpecial_452D80(923, 100)
	return 1
}
func browserConnectionReset() int {
	GetClient().Cli().GUI.Focus(browserUI.world)
	if Sub_43BE30() == 0 || memmap.Uint32(0x5D4594, 815084) == 0 {
		Sub_44A400()
	}
	browserUI.connectionState = 0
	browserUI.transition = 0
	return 0
}
func browserClose() int {
	browserWindow(uint32(uintptr(browserUI.mapWindow))).Capture(false)
	sessionFilterClose()
	browserPopupClose()
	detail := browserWindow(uint32(uintptr(browserUI.detailPanel)))
	if detail != nil && detail.Parent() == nil {
		detail.StackPop()
		detail.Destroy()
		browserUI.detailPanel = nil
	}
	if browserUI.world != nil {
		browserUI.world.Destroy()
		browserUI.world = nil
	}
	browserConnectionReset()
	GetClient().Cli().GUI.Focus(nil)
	browserUI.resultCount = 0
	browserListClear(false)
	Sub_554D10()
	nox_client_setCursorType_477610(0)
	GetServer().SetUpdateFunc2(nil)
	return 1
}
func browserChooseCharacter() int {
	browserAnimationOut()
	(*gui.Anim)(unsafe.Pointer(browserUI.animation)).FncDoneOutPtr = animationCallbackKey(animationKeyBrowserHide)
	browserWindow(uint32(uintptr(browserUI.detailPanel))).StackPop()
	return uiWindowEnable(browserWindow(uint32(uintptr(browserUI.mapWindow))), 0)
}
func browserHideAfterChoice() int {
	if GetClient().GameGetStateCode() == 1700 {
		return browserAnimationFinish()
	}
	browserUI.world.SetHidden(true)
	browserWindow(uint32(uintptr(browserUI.detailPanel))).SetHidden(true)
	nox_client_setCursorType_477610(0)
	return 1
}
func browserShowList() int {
	w := browserWindow(uint32(browserUI.filter))
	if !browserHidden(w) {
		w.SetHidden(true)
		sessionFilterSave()
	}
	browserWindow(uint32(browserUI.gameList)).SetHidden(false)
	browserWindow(uint32(uintptr(browserUI.mapWindow))).SetHidden(true)
	browserWindow(uint32(browserUI.overview)).SetHidden(true)
	serverPanelsEnable(browserUI.world, 10006, 10007, 1)
	serverPanelsHide(browserUI.world, 10047, 10051, false)
	browserLabel("ListJoinServer")
	browserUI.listMode = 1
	return 0
}

func sub_4379C0() { browserListReset() }

func nox_client_guiXxx_43A9D0() int32 { return int32(browserClose()) }
