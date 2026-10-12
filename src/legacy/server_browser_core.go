package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"math"
	"unsafe"
)

func browserModeName(flags uint16) string {
	name := "Arena"
	switch {
	case flags&0x1000 != 0:
		name = "Quest"
	case flags&0x20 != 0:
		name = "CTF"
	case flags&0x400 != 0:
		name = "Highlander"
	case flags&0x10 != 0:
		name = "KotR"
	case flags&0x40 != 0:
		name = "Flagball"
	case flags&0x80 != 0:
		name = "Chat"
	}
	return GetServer().S().Strings().GetStringInFile(strman.ID(name), "noxworld.c")
}
func browserHit(point *[2]uint32, record unsafe.Pointer) bool {
	// C converts the unsigned subtraction to double, including its wraparound.
	x := float64(uint32(int32(*(*int16)(unsafe.Add(record, 44)))) - point[0])
	y := float64(uint32(int32(*(*int16)(unsafe.Add(record, 46)))) - point[1])
	return math.Sqrt(y*y+x*x) <= memmap.Float64(0x581450, 9720)
}
func browserCount(point *[2]uint32, head *legacyListNode) int {
	n := 0
	for it := listNext(head); it != nil; it = listNext(it) {
		if browserHit(point, unsafe.Pointer(it)) {
			n++
		}
	}
	return n
}
func browserPopupClamp(x, y int32, out *[2]uint32) {
	px, py := x-100, y-20
	if px+200 > 600 {
		px = 400
	}
	if py+200 > 451 {
		py = 251
	}
	if uint32(py) < 27 {
		py = 27
	}
	if px < 216 {
		px = 216
	}
	out[0], out[1] = uint32(px), uint32(py)
}
func browserRegion(x, y int32) int {
	for i := uintptr(0); i < 4; i++ {
		off := 87528 + 8*i
		if x > int32(memmap.Int16(0x587000, off)) && x < int32(memmap.Int16(0x587000, off+4)) && y > int32(memmap.Int16(0x587000, off+2)) && y < int32(memmap.Int16(0x587000, off+6)) {
			return int(i)
		}
	}
	return 0
}
func browserTrimName(text *uint16, width byte) *uint16 {
	end := alloc.StrLen(text)
	for {
		measured := GetClient().R2().GetStringSizeWrapped(nil, alloc.GoString16(text), 0).X
		*(*uint16)(unsafe.Add(unsafe.Pointer(text), 2*end)) = 0
		end--
		if measured+5 <= int(width) {
			return text
		}
	}
}
func browserFormatEndpoint(addr string, port uint16, dst *byte) int {
	s := fmt.Sprintf("%s:%d", addr, port)
	out := unsafe.Slice(dst, len(s)+1)
	copy(out, s)
	out[len(s)] = 0
	return len(s)
}

// Browser compatibility helpers call the native owners.

func nox_gui_wol_gameModeString_43BCB0(v int16) *wchar2_t {
	return (*wchar2_t)(unsafe.Pointer(alloc.InternCString16(browserModeName(uint16(v)))))
}

func sub_4A25C0(point *uint32, head *int32) int32 {
	return int32(browserCount((*[2]uint32)(unsafe.Pointer(point)), (*legacyListNode)(unsafe.Pointer(head))))
}

func sub_4A2830(x, y int32, out *uint32) *uint32 {
	browserPopupClamp(int32(x), int32(y), (*[2]uint32)(unsafe.Pointer(out)))
	return out
}

func sub_438DD0(x, y uint32) int32 {
	if browserUI.region == -1 {
		if x > 216 && x < 600 && y > 27 && y < 451 {
			return 1
		}
	} else if x > 226 && x < 590 && y > 37 && y < 441 {
		return 1
	}
	return 0
}

func sub_43AF40() int32 { return int32(browserUI.creating) }

func sub_43AF80() int32 { return int32(browserUI.connectionState) }

func sub_43AF90(v int32) int32 { browserUI.connectionState = uint32(v); return v }

func nox_client_setConnError_43AFA0(v int32) {
	browserUI.connectionError = uint32(v)
	browserUI.connectionState = 2
}

func nox_client_getServerAddr_43B300() uint32 {
	if browserUI.hasSelection == 0 {
		return 0
	}
	// Preserve the established short/octal/hex IPv4 forms at this boundary.
	return browserIPv4(GoStringP(unsafe.Add(browserUI.selected, 12)))
}

func nox_client_getServerPort_43B320() int32 {
	if browserUI.hasSelection == 0 {
		return 0
	}
	return int32(memmap.Uint32(0x5D4594, 814604))
}

func sub_43B340() int32 {
	if browserUI.hasSelection == 0 {
		return 0
	}
	return int32(*(*uint16)(unsafe.Add(browserUI.selected, 163)))
}

func sub_43B6D0() int32 { return int32(browserUI.transition) }

func sub_43BC10(text *wchar2_t, width uint8) *uint16 {
	return (*uint16)(unsafe.Pointer(browserTrimName((*uint16)(unsafe.Pointer(text)), byte(width))))
}

func nox_sprintAddrPort_43BC80(addr *int8, port uint16, dst *int8) int32 {
	return int32(browserFormatEndpoint(GoStringP(unsafe.Pointer(addr)), uint16(port), (*byte)(unsafe.Pointer(dst))))
}

func nox_wol_servers_sortBtnHandler_4A0290(id int32) {
	if id < 10047 || id > 10051 {
		return
	}
	v := uint32(id-10047) * 2
	if uint32(browserUI.sort) == v {
		v++
	}
	browserUI.sort = uint32(v)
}
