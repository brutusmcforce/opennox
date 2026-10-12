package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
)

var animationCallbackStorage [25]byte

const (
	animationKeyShowSelClass  = 0
	animationKeyShowSelColor  = 1
	animationKeyClassStart    = 2
	animationKeyClassDone     = 3
	animationKeyColorStart    = 4
	animationKeyColorDone     = 5
	animationKeyShowOptions   = 6
	animationKeyOptionsDone   = 7
	animationKeyBindingsShow  = 8
	animationKeyBindingsBack  = 9
	animationKeyBindingsDone  = 10
	animationKeyOptionsStart  = 11
	animationKeyOptionsFinish = 12
	animationKeyBrowserHide   = 13
	animationKeyBrowserDone   = 14
	animationKeyBrowserOut    = 15
	animationKeyBrowserShow   = 16
	animationKeyMainMenuStart = 17
	animationKeyMainMenuDone  = 18
	animationKeyDrawGeneral   = 19
	animationKeyShowSelChar   = 20
	animationKeySwitchStates  = 21
	animationKeyShowMainMenu  = 22
	animationKeySelCharStart  = 23
	animationKeySelCharDone   = 24
)

func animationCallbackKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&animationCallbackStorage[id])
}

func init() {
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyShowSelClass), func() int { return characterShowClass() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyShowSelColor), func() int { return characterShowColor() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyClassStart), func() int { return characterClassStart() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyClassDone), func() int { return characterClassDone() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyColorStart), func() int { return characterColorStart() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyColorDone), func() int { return characterColorDone() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyShowOptions), func() int { return optionsMenu.construct() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyOptionsDone), func() int { return optionsMenuDone() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBindingsShow), func() int { return bindingMenu.construct() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBindingsBack), func() int { return bindingMenuBack() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBindingsDone), func() int { return bindingMenuDone() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyOptionsStart), func() int { return Sub_4AA9C0() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyOptionsFinish), func() int { return Sub_4AAA10() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBrowserHide), func() int { return browserHideAfterChoice() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBrowserDone), func() int { return browserAnimationFinish() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBrowserOut), func() int { return browserAnimationOut() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyBrowserShow), func() int { return browserShow() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyMainMenuStart), func() int { return WinMainMenuAnimOutStartFnc() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyMainMenuDone), func() int { return WinMainMenuAnimOutDoneFnc() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyDrawGeneral), func() int { return nox_client_drawGeneralCallback_4A2200() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyShowSelChar), func() int { return Nox_game_showSelChar_4A4DB0() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeySwitchStates), func() int { return bool2int(GetClient().GameStateSwitch()) })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeyShowMainMenu), func() int { return Nox_game_showMainMenu_4A1C00() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeySelCharStart), func() int { return Sub_4A50A0() })
	gui.RegisterAnimationCallbackGo(animationCallbackKey(animationKeySelCharDone), func() int { return Sub_4A50D0() })
}
