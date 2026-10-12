package legacy

import (
	"encoding/binary"
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/netlist"
)

func interactionPlayerBlocked() bool {
	p := Get_dword_8531A0_2576()
	return p != nil && p.Field3680&3 != 0
}
func interactionTalk(dr *client.Drawable) {
	if dr == nil || interactionPlayerBlocked() || uiShopActive() == 1 || sessionQuitShown() == 1 {
		return
	}
	var msg [4]byte
	binary.LittleEndian.PutUint16(msg[:], 464)
	binary.LittleEndian.PutUint16(msg[2:], uint16(dr.NetCode32))
	GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])
}
func interactionUse(dr *client.Drawable) {
	if dr == nil || interactionPlayerBlocked() {
		return
	}
	uiInventoryItemRequest(123, dr)
}
func interactionTrade(dr *client.Drawable) {
	if dr == nil || interactionPlayerBlocked() || sub_47A260() == 1 || sessionQuitShown() == 1 {
		return
	}
	var msg [4]byte
	binary.LittleEndian.PutUint16(msg[:], 5577)
	binary.LittleEndian.PutUint16(msg[2:], uint16(nox_xxx_netGetUnitCodeCli_578B00(int32(uintptr(dr.C())))))
	GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])
}
func interactionMouseMode(v int32) int32 {
	switch v {
	case 1:
		interactionPrimaryKey = 1
		interactionCursorMode = 9
		return 0
	case 2:
		interactionPrimaryKey = 2
		interactionCursorMode = 13
		return 0
	default:
		interactionPrimaryKey = 0
		interactionCursorMode = 5
		return v - 2
	}
}
func interactionInitTime() uint64 {
	v := uint64(uint32(PlatformTicks()))
	*memmap.PtrUint64(0x5D4594, 811908) = v
	return v
}
func interactionHasBuff(dr *client.Drawable, bit uint8) bool {
	// The original x86 shift masks its count, including the signed-char ABI.
	return dr != nil && dr.Buffs&(uint32(1)<<(bit&31)) != 0
}
func interactionTextState(text *uint16, value uint32) *uint16 {
	dst := memmap.PtrUint16(0x5D4594, 811376)
	for i := uintptr(0); ; i += 2 {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(text), i))
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), i)) = v
		if v == 0 {
			break
		}
	}
	*memmap.PtrUint32(0x5D4594, 811060) = value
	return dst
}
func interactionObserverToggle() int {
	p := Get_dword_8531A0_2576()
	show := 1
	if p != nil && p.Field3680&1 != 0 {
		show = 0
	}
	return int(nox_xxx_showObserverWindow_48CA70(int32(show)))
}
func interactionToggleDrawing() uint32 {
	v := 1 - uint32(interactionDrawToggle)
	interactionDrawToggle = uint32(v)
	return v
}
func interactionFrameGate() int {
	if serverOptionsRoot != 0 || sessionQuitShown() != 0 || sub_49C810() != 0 || sessionMOTDShown() != 0 || scoreboardEnabled() || Nox_gui_console_flagXxx_451410() != 0 {
		*memmap.PtrUint32(0x5D4594, 811920) = GetServer().S().Frame()
		return 1
	}
	frame := GetServer().S().Frame()
	if frame != 2 {
		return bool2int(frame-memmap.Uint32(0x5D4594, 811920) == 1)
	}
	*memmap.PtrUint32(0x5D4594, 811920) = frame
	return 1
}
func interactionHUDVisibility() {
	flag := Nox_client_getRenderGUI()
	if memmap.Uint32(0x5D4594, 811064) == uint32(flag) || noxflags.HasEngine(noxflags.EngineNoRendering) {
		return
	}
	*memmap.PtrUint32(0x5D4594, 811064) = uint32(flag)
	sub_4721A0(flag)
	quickbarVisible(flag != 0)
	nox_window_set_visible_unk5(flag)
	bookTemporaryShow(flag)
	teamUIHUDShow(false, flag)
	teamUIHUDShow(true, flag)
	Sub_4706C0(flag)
	if flag == 0 {
		Sub_478000()
	}
}
func interactionPlayerAnimation() int {
	dr := *(**client.Drawable)(memmap.PtrOff(0x852978, 8))
	return bool2int(dr == nil || dr.AnimInd == 1 || dr.AnimInd == 2 || dr.AnimInd == 51)
}
func interactionIsObserver() int {
	p := Get_dword_8531A0_2576()
	// The C predicate reads all four bytes, including the legacy padding.
	return bool2int(p != nil && *(*uint32)(unsafe.Pointer(&p.Active)) == 1 && p.Field3680&3 != 0)
}

// Legacy-shaped Go adapters preserve numeric conventions for existing callers.
// Only entries with live C callers retain an export.
func nox_xxx_clientTalk_42E7B0(dr *nox_drawable) { interactionTalk(asDrawable(dr)) }

func nox_xxx_clientCollideOrUse_42E810(dr *nox_drawable) { interactionUse(asDrawable(dr)) }

func nox_xxx_clientTrade_42E850(dr *nox_drawable) { interactionTrade(asDrawable(dr)) }

func sub_430AA0(v int32) int32 { return int32(interactionMouseMode(int32(v))) }

func nox_client_mousePriKey_430AF0() int32 { return int32(interactionPrimaryKey) }

func nox_xxx_cursor_430B00() int32 { return int32(interactionCursorMode) }

func nox_client_setMousePos_430B10(x, y int32) {
	GetClient().ChangeMousePos(image.Pt(int(x), int(y)), true)
}

func nox_xxx_initTime_435570() int64 { return int64(interactionInitTime()) }

func nox_client_drawable_testBuff_4356C0(dr *nox_drawable, bit int8) bool {
	return bool(interactionHasBuff(asDrawable(dr), uint8(bit)))
}

func sub_435700(text *uint16, v int32) *uint16 {
	return (*uint16)(unsafe.Pointer(interactionTextState((*uint16)(unsafe.Pointer(text)), uint32(v))))
}

func nox_xxx_cliToggleObsWindow_4357A0() int32 { return int32(interactionObserverToggle()) }

func sub_435F60() int32 { return int32(interactionToggleDrawing()) }

func sub_437100() { interactionHUDVisibility() }

func nox_xxx_playerAnimCheck_4372B0() int32 { return int32(interactionPlayerAnimation()) }

func nox_xxx_clientIsObserver_4372E0() int32 { return int32(interactionIsObserver()) }
