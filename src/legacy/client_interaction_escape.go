package legacy

import "github.com/brutusmcforce/opennox/common/memmap"

func interactionEscape() {
	if memmap.Uint32(0x5D4594, 1096672) != 0 || Nox_video_inFadeTransition_44E0D0() != 0 {
		return
	}
	handled := quickbarCancelCapture() != 0
	// Chat closes even when the preceding capture controller handled Escape.
	chat := interactionChatClose()
	if chat != 0 || handled {
		return
	}
	if *bookWord(1047520) == 1 {
		bookFinishAddition()
		return
	}
	if uiAmountCancel() != 0 {
		handled = true
	}
	if quickbarCloseExpanded() != 0 {
		handled = true
	}
	if uiInventoryCloseIdentify() != 0 {
		handled = true
	}
	dialog := Sub_44A4E0()
	if dialog != 0 || handled {
		return
	}
	switch uiShopMode() {
	case 2, 3, 4:
		uiShopSetMode(1)
		return
	}
	if uiShopCancelRequest() != 0 || interactionConversationCancel() != 0 {
		return
	}
	if uiInventoryCloseWindow() != 0 {
		handled = true
	}
	if bookHide(0) != 0 {
		handled = true
	}
	if Nox_gui_console_Hide_4512B0() != 0 {
		handled = true
	}
	if sessionMOTDClose() != 0 {
		handled = true
	}
	if serverOptionsTryClose() != 0 {
		handled = true
	}
	if voteGUIHide() != 0 {
		handled = true
	}
	if optionsClose(1) != 0 {
		handled = true
	}
	if bindingClose(1) != 0 {
		handled = true
	}
	// Save closure remains last and still runs after any earlier window closed.
	saved := Sub_46D6F0()
	if saved != 0 || handled {
		return
	}
	if sessionSaveAllowed() != 0 {
		sessionQuitToggle()
	} else {
		Nox_xxx_clientPlaySoundSpecial_452D80(231, 100)
	}
}

func nox_xxx_consoleEsc_49B7A0() { interactionEscape() }
