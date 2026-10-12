package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

type summonRecord struct {
	Code, Type, Active uint32
	X, Y               int32
	Size               byte
	_                  [3]byte
	Index, Flash       uint32
}

func summonWord(off uintptr) *uint32 {
	switch off {
	case 1320988:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1320988))
	case 1320992:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1320992))
	case 1321024:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321024))
	case 1321032:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321032))
	case 1321036:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321036))
	case 1321040:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321040))
	case 1321044:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321044))
	case 1321196:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321196))
	case 1321204:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321204))
	case 1321208:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1321208))
	}
	return memmap.PtrUint32(0x5D4594, off)
}
func summonMenuWidth() *uint32 {
	return (*uint32)(unsafe.Pointer(&nox_xxx_screenWidth_587000_184452))
}
func summonAt(i int32) *summonRecord {
	return (*summonRecord)(memmap.PtrOff(0x5D4594, 1321052+uintptr(i)*32))
}
func summonAddress(r *summonRecord) uint32 { return memmap.PtrWord(unsafe.Pointer(r)) }
func summonGrid(x, y int32) *uint32        { return summonWord(1321180 + uintptr(y+2*x)*4) }
func summonClearGrid() uint32 {
	for i := uintptr(1321180); i < 1321196; i += 4 {
		*summonWord(i) = 0
	}
	return memmap.PtrWord(unsafe.Pointer(summonWord(1321200)))
}
func summonFirst() *summonRecord {
	for i := int32(0); i < 4; i++ {
		r := summonAt(i)
		if r.Active != 0 {
			return r
		}
	}
	return nil
}
func summonNext(r *summonRecord) *summonRecord {
	for i := int32(r.Index + 1); i < 4; i++ {
		r := summonAt(i)
		if r.Active != 0 {
			return r
		}
	}
	return nil
}
func summonFind(code uint32) *summonRecord {
	for i := int32(0); i < 4; i++ {
		r := summonAt(i)
		if r.Active != 0 && r.Code == code {
			return r
		}
	}
	return nil
}
func summonAllocate() *summonRecord {
	for i := int32(0); i < 4; i++ {
		r := summonAt(i)
		if r.Active == 0 {
			*r = summonRecord{Active: 1, Index: uint32(i)}
			return r
		}
	}
	return nil
}
func summonDeactivate(r *summonRecord) uint32 { r.Active = 0; return summonAddress(r) }
func summonFootprint(size int32) (w, h int32) {
	switch size {
	case 1:
		return 1, 1
	case 2:
		return 1, 2
	case 4:
		return 2, 2
	}
	return size, size
}
func summonAvailable(pos *[2]int32, size int32) int {
	w, h := summonFootprint(size)
	for y := pos[1]; y < pos[1]+h; y++ {
		for x := pos[0]; x < pos[0]+w; x++ {
			if *summonGrid(x, y) != 0 {
				return 0
			}
		}
	}
	return 1
}
func summonPaint(pos *[2]int32, size int32, value uint32) int32 {
	w, h := summonFootprint(size)
	result := h
	for y := pos[1]; y < pos[1]+h; y++ {
		for x := pos[0]; x < pos[0]+w; x++ {
			*summonGrid(x, y) = value
		}
		result = pos[1] + h
	}
	return result
}
func summonPlace(r *summonRecord) int32 {
	pos := (*[2]int32)(unsafe.Pointer(&r.X))
	for i := int32(0); ; {
		r.X = *memmap.PtrInt32(0x587000, 184456+uintptr(i)*8)
		r.Y = *memmap.PtrInt32(0x587000, 184460+uintptr(i)*8)
		if summonAvailable(pos, int32(r.Size)) != 0 {
			return summonPaint(pos, int32(r.Size), summonAddress(r))
		}
		i += int32(r.Size)
		if i >= 4 {
			return int32(r.Size)
		}
	}
}
func summonLayout() uint32 {
	summonClearGrid()
	for _, size := range []byte{1, 2, 4} {
		for r := summonFirst(); r != nil; r = summonNext(r) {
			if r.Size == size {
				summonPlace(r)
			}
		}
	}
	return 0
}
func summonGet(pos *[2]int32) uint32 {
	if pos == nil || pos[0] < 0 || pos[0] >= 2 || pos[1] < 0 || pos[1] >= 2 {
		return 0
	}
	return *summonGrid(pos[0], pos[1])
}
func summonMobile(r *summonRecord) int {
	p := summonWord(1321208)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID("CarnivorousPlant"))
	}
	if r.Type != *p {
		return 1
	}
	return 0
}
func summonAnyMobile() int {
	p := summonWord(1321208)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID("CarnivorousPlant"))
	}
	for r := summonFirst(); r != nil; r = summonNext(r) {
		if summonMobile(r) != 0 {
			return 1
		}
	}
	return 0
}
func summonClass(id int) byte {
	v := GetClient().Cli().Things.TypeByInd(id).ObjSubClass
	if v&1 != 0 {
		return 1
	}
	if v&2 != 0 {
		return 2
	}
	return 4
}
