package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"unsafe"
)

func quickbarQueueTarget(id uint32) int {
	if *quickbarWord(1047916) != 0 {
		return 0
	}
	*quickbarWord(1047916) = id
	*quickbarByte(1047920) = 0
	InputSetKeyTimeoutLegacy(5)
	*quickbarWord(1047928) = 0
	*quickbarWord(1047924) = 0
	return 1
}
func quickbarQueueInstant(id uint32) int {
	if *quickbarWord(1047916) != 0 {
		return 0
	}
	*quickbarWord(1047916) = id
	*quickbarByte(1047920) = 0
	InputSetKeyTimeoutLegacy(5)
	*quickbarWord(1047924) = 1
	return 1
}
func quickbarLastButton(slot int) {
	*memmap.PtrUint32(0x587000, 133484) = uint32(slot)
	*quickbarWord(1049540) = GetServer().S().Frame()
}
func quickbarSendAbility(id uint32) int {
	if cur := memmap.Uint32(0x5D4594, 1096672); cur != 0 {
		return int(cur)
	}
	if id == 0 {
		return 0
	}
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(31), netlist.Kind0, []byte{122, byte(id)}))
}
func quickbarSendSpells(ids *uint32, count int, flag byte) int {
	if cur := memmap.Uint32(0x5D4594, 1096672); cur != 0 {
		return int(cur)
	}
	if *ids == 0 {
		return 0
	}
	var msg [22]byte
	msg[0] = 121
	msg[21] = flag
	for i := 0; i < 5 && i < count; i++ {
		binary.LittleEndian.PutUint32(msg[1+4*i:], *(*uint32)(unsafe.Add(unsafe.Pointer(ids), 4*i)))
	}
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(31), netlist.Kind0, msg[:]))
}
func quickbarInvokeSlot(slot int) {
	p := quickbarPlayer()
	if p == 0 || Nox_xxx_playerAnimCheck_4372B0() != 0 {
		return
	}
	s := &quickbarMain().Current[slot]
	if bookClass(p) != 0 {
		quickbarSendSpells(&s.ID, 1, byte(s.Flags)&1)
	} else {
		quickbarSendAbility(s.ID)
	}
	quickbarLastButton(slot)
}
func quickbarSendPending() {
	if *quickbarWord(1047916) != 0 {
		quickbarSendSpells(quickbarWord(1047916), 1, *quickbarByte(1047924))
		*quickbarWord(1047916) = 0
	}
}
func quickbarSpellCursor(id uint32, flags byte) int {
	if id == 0 {
		*quickbarWord(1047928) = 0
		return 0
	}
	bookSound(766)
	if flags&1 == 0 || bool(nox_xxx_spellHasFlags_424A50(int(id), 0x2000)) {
		*quickbarWord(1047556) = id
		*quickbarWord(1047928) = 1
		return 0
	}
	quickbarQueueInstant(id)
	return 1
}
func quickbarAbilityCursor(id uint32) uint32 {
	if id == 0 {
		*quickbarWord(1047932) = 0
		*quickbarWord(1047936) = 0
		return 0
	}
	off := uintptr(1047764 + 24*id)
	instant := *quickbarWord(off + 4)
	if instant != 0 {
		*quickbarWord(1047932) = 0
		*quickbarWord(1047936) = 0
		quickbarSendAbility(id)
	} else {
		*quickbarWord(1047932) = 1
		*quickbarWord(1047936) = *quickbarWord(off)
	}
	bookSound(766)
	return instant
}
func quickbarSendPendingAbility() int {
	id := *quickbarWord(1047936)
	if id == 0 {
		return 0
	}
	ret := quickbarSendAbility(id)
	*quickbarWord(1047936) = 0
	*quickbarWord(1047932) = 0
	return ret
}
func quickbarResetFlash() int {
	clear(unsafe.Slice(quickbarByte(1049544), 136))
	*quickbarByte(1049680) = 0
	return 0
}
func quickbarSetFlash(id int, value byte) {
	if id >= 0 && id < 140 {
		*quickbarByte(1049544 + uintptr(id)) = value
	}
}
