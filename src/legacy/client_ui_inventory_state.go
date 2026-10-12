package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func sub_4673F0(a, b int) int {
	*memmap.PtrUint32(0x5D4594, 1062580) = uint32(a)
	*memmap.PtrUint32(0x5D4594, 1062584) = uint32(b)
	return a
}

func sub_467410(v int) int { *memmap.PtrUint32(0x5D4594, 1062540) = uint32(v); return v }

func sub_467420(v int8) int8 { *memmap.PtrUint8(0x5D4594, 1062536) = byte(v); return v }

func sub_467430() byte { return memmap.Uint8(0x5D4594, 1062536) }

func sub_467440(v int) int { *memmap.PtrUint32(0x5D4594, 1062544) = uint32(v); return v }

func sub_467450(v int) int { *memmap.PtrUint32(0x5D4594, 1062548) = uint32(v); return v }

func sub_467470(index int, v float32) int {
	i := uint8(index)
	*memmap.PtrFloat32(0x5D4594, 1063100+uintptr(i)*4) = v
	return int(i)
}

func sub_467490(v int) int { dword_5d4594_1062552 = uint32(v); return v }

func sub_4674A0() int { return int(dword_5d4594_1062552) }

func nox_window_set_visible_unk5(v int) {
	uiMeterHide((*gui.Window)(unsafe.Pointer(legacyGlobals.nox_win_unk5)), v == 0)
}
func uiInventoryUsePotion(typ uint32) {
	if memmap.Uint32(0x5D4594, 1096672) == 1 || noxflags.HasGame(noxflags.GamePause) {
		return
	}
	if cell := uiInventoryFindType(typ); cell != nil {
		cell.Drawable.NetCode32 = cell.Codes[0]
		uiInventoryUse(cell.Drawable)
	}
}

func sub_467590() int {
	if p := uiMeterPlayer(); p != nil {
		return int(*(*int8)(unsafe.Add(p, 3684)))
	}
	return 1
}
func uiInventoryMode() int { return int(dword_5d4594_1049864) }

func uiInventoryItemHealth(code int, current, maximum int16) int16 {
	if found := uiInventoryFindCode(uint32(code)); found != nil {
		dr := found.Cell.Drawable
		*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(292)) /* Field_73_1 */)) = uint16(current)
		*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(294)) /* Field_73_2 */)) = uint16(maximum)
		return int16(memmap.PtrWord(dr.C()))
	}
	dr := uiInventoryDragged()
	if dr != nil && dr.NetCode32 == uint32(code) {
		*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(292)) /* Field_73_1 */)) = uint16(current)
		*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(294)) /* Field_73_2 */)) = uint16(maximum)
		return maximum
	}
	return int16(memmap.PtrWord(unsafe.Pointer(dr)))
}

func sub_467680() {
	if dword_5d4594_1049864 == 6 {
		dword_5d4594_1049864 = 0
	}
}

func sub_467740(v int) int { dword_5d4594_1062488 = uint32(v); return v }

func sub_4678C0() int { return int(dword_5d4594_1062488) }

func sub_467930(code, current, maximum int) unsafe.Pointer {
	if code == 0 {
		return nil
	}
	found := uiInventoryFindCode(uint32(code))
	if found == nil {
		return nil
	}
	cell := found.Cell
	dr := cell.Drawable
	*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(448)) /* Union._ */)) = uint16(current)
	*(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(450)) /* Union._ */)) = uint16(maximum)
	if cell.Equipped == 1 {
		return unsafe.Pointer(uintptr(sub_470D90(current, maximum))) // an integer result, not a pointer
	}
	return unsafe.Pointer(cell)
}
