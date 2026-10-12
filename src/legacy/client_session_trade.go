package legacy

import (
	"encoding/binary"
	"unsafe"
)

func clientSessionTrade(data []byte) int {
	if len(data) < 2 {
		return -1
	}
	size := 0
	switch data[1] {
	case 1, 2, 7:
		size = 2
	case 3:
		size = 3
	case 4:
		size = 15
	case 5, 9, 27:
		size = 4
	case 6:
		size = 14
	case 8:
		size = 18
	case 12:
		size = 52
	case 13:
		size = 86
	case 29, 31:
		size = 8
	default:
		return -1
	}
	if len(data) < size {
		return -1
	}
	if !Nox_client_isConnected() {
		return size
	}
	word := func(off int) uint32 { return uint32(binary.LittleEndian.Uint16(data[off:])) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	p := unsafe.Pointer(&data[0])
	switch data[1] {
	case 1:
		uiTradeFinish()
	case 2:
		uiShopClose()
	case 3:
		uiTradeAcceptance(p)
	case 4:
		uiTradeAdd(p)
	case 5:
		uiTradeRemove(p)
	case 6:
		uiTradeMoney(p)
	case 7:
		uiTradePrepare()
	case 8:
		uiShopAdd(word(2), word(4), dword(6), uint16(dword(10)), unsafe.Add(p, 14))
	case 9:
		uiShopRemove(word(2))
	case 12:
		uiTradeStart(p)
	case 13:
		// Terminate the fixed wire fields before passing them to string owners.
		name := make([]uint16, 26)
		for i := 0; i < 25; i++ {
			name[i] = uint16(word(4 + 2*i))
		}
		end := 54
		for end < 86 && data[end] != 0 {
			end++
		}
		uiShopStart(&name[0], string(data[54:end]), word(2))
	case 27:
		uiShopGoldWarning(word(2))
	case 29:
		uiShopSellShow(word(2), dword(4))
	case 31:
		uiShopRepairShow(word(2), dword(4))
	}
	return size
}
