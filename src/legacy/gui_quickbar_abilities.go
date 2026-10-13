package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
)

func quickbarAbilityState(id, state uint32) uintptr {
	frame := GetServer().S().Frame()
	for n := 1; n <= 5; n++ {
		off := uintptr(1047764 + 24*n)
		if *quickbarWord(off) == id {
			*quickbarWord(off + 20) = 0
			if state == 0 {
				*quickbarWord(off + 20) = frame
			}
			*quickbarWord(off + 8) = state
		}
	}
	return uintptr(unsafe.Pointer(quickbarWord(1047928)))
}
func quickbarResetAbility(id byte) uintptr {
	if id != 6 {
		return quickbarAbilityState(*quickbarWord(1047764 + 24*uintptr(id)), 1)
	}
	var ret uintptr
	for n := 1; n <= 5; n++ {
		ret = quickbarAbilityState(*quickbarWord(1047764 + uintptr(24*n)), 1)
	}
	return ret
}
func quickbarAbilityFlags(id, on uint32) uintptr {
	mask := uint32(1) << (id & 31)
	for n := 1; n <= 5; n++ {
		off := uintptr(1047764 + 24*n)
		if *quickbarWord(off) == id {
			if on != 0 {
				*quickbarWord(off + 12) |= mask
			} else {
				*quickbarWord(off + 12) &^= mask
			}
		}
	}
	return uintptr(unsafe.Pointer(quickbarWord(1047920)))
}
func quickbarAbilityAvailable(id uint32) int {
	for n := 1; n <= 5; n++ {
		off := uintptr(1047764 + 24*n)
		if *quickbarWord(off) == id {
			return bool2int(*quickbarWord(off + 12)&(uint32(1)<<(id&31)) != 0)
		}
	}
	return 0
}
func quickbarAbilityReward(id, rank int, notify uintptr) {
	if id < 1 || id >= 6 {
		return
	}
	for n := 1; n <= 5; n++ {
		off := uintptr(1047764 + 24*n)
		if *quickbarWord(off) == uint32(id) && *quickbarWord(off + 16) != uint32(rank) {
			if p := quickbarPlayer(); p != nil && noxflags.HasGame(noxflags.GameFlag(2)) {
				*bookPlayerWord(unsafe.Pointer(p), 3696, id) = uint32(rank)
			}
			*quickbarWord(off + 16) = uint32(rank)
			if rank != 0 {
				bookAbilityReward(id, notify, int(notify))
			}
		}
	}
}
