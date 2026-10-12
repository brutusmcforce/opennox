package legacy

import (
	"encoding/binary"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func interactionPickup(dr *client.Drawable) {
	if memmap.Uint32(0x5D4594, 1064928) == 0 {
		for i, name := range []string{"Gold", "QuestGoldPile", "QuestGoldChest"} {
			*memmap.PtrUint32(0x5D4594, 1064928+uintptr(4*i)) = uint32(GetClient().Cli().Things.IndByID(name))
		}
	}
	if dr == nil {
		return
	}
	typ := dr.TypeIDVal
	if typ == memmap.Uint32(0x5D4594, 1064928) || typ == memmap.Uint32(0x5D4594, 1064932) || typ == memmap.Uint32(0x5D4594, 1064936) || uiInventoryCapacity(int32(typ), 1) != 0 {
		uiInventoryItemRequest(115, dr)
		return
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(925, 100)
	text, free := alloc.CString16(GetClient().Cli().Strings().GetStringInFile("pickup.c:CarryingTooMuch", "gamewin.c"))
	defer free()
	interactionCentered(text)
}
func interactionReportSecondary(dr *client.Drawable) int {
	var msg [3]byte
	msg[0] = 224
	binary.LittleEndian.PutUint16(msg[1:], uint16(drawableUnitCode(dr)))
	return reliableClientSend(31, msg[:], nil, 1)
}

func nox_xxx_clientPickup_46C140(dr *nox_drawable) { interactionPickup(asDrawable(dr)) }

func nox_xxx_clientReportSecondaryWeapon_4BF010(dr int32) int32 {
	return int32(interactionReportSecondary((*client.Drawable)(memmap.WordPtr(uint32(dr)))))
}
