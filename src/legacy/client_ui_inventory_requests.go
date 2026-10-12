package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/internal/netlist"
)

func uiInventoryItemRequest(op byte, dr *client.Drawable) int {
	var msg [3]byte
	msg[0] = op
	binary.LittleEndian.PutUint16(msg[1:], uint16(drawableUnitCode(dr)))
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:]))
}
func uiInventoryEquipRequest(dr *client.Drawable) int  { return uiInventoryItemRequest(117, dr) }
func uiInventoryDequipRequest(dr *client.Drawable) int { return uiInventoryItemRequest(118, dr) }
func uiInventoryUse(dr *client.Drawable) {
	if dr != nil {
		uiInventoryItemRequest(116, dr)
	}
}

func uiInventoryTrade(op byte, code uint16) int {
	msg := [4]byte{201, op, byte(code), byte(code >> 8)}
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:]))
}

func nox_xxx_send2ServInvenFail_461630(v int16) int32 {
	msg := [3]byte{241, byte(v), byte(uint16(v) >> 8)}
	return int32(reliableClientSend(31, msg[:], nil, 0))
}

func nox_xxx_clientDrop_465BE0(pos *[2]int32) int32 {
	dr := uiInventoryDragged()
	if dr == nil {
		return 0
	}
	msg := [7]byte{114}
	binary.LittleEndian.PutUint16(msg[1:], uint16(drawableUnitCode(dr)))
	binary.LittleEndian.PutUint16(msg[3:], uint16(pos[0]))
	binary.LittleEndian.PutUint16(msg[5:], uint16(pos[1]))
	return int32(bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])))
}

func nox_xxx_clientKeyEquip_465C30(col, row int32) int32 {
	uiInventorySetClick(int(col), int(row))
	uiInventoryDragCopy()
	uiInventoryEquipRequest(uiInventoryDragged())
	return int32(uiInventoryPlace(uiInventoryDragged(), int(col), int(row)))
}
func uiInventorySetDragged(dr *client.Drawable) {
	uiInventoryDraggedItem = dr
}
