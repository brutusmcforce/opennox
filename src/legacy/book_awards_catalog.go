package legacy

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var bookEnchantN int32 = 29

func bookEnchantCountActive(u *server.Object) int8 {
	var count uint8
	for i := int32(0); i < bookEnchantN; i++ {
		id := *memmap.PtrInt32(0x587000, uintptr(66000+4*i))
		if u.HasEnchant(server.EnchantID(id)) {
			count++
		}
	}
	return int8(count)
}
func bookEnchantFirst() int32 {
	if bookEnchantN <= 0 {
		return -1
	}
	return *memmap.PtrInt32(0x587000, 66000)
}
func bookEnchantNext(id int32) int32 {
	if bookEnchantN <= 0 {
		return -1
	}
	for i := int32(0); i < bookEnchantN-1; i++ {
		if *memmap.PtrInt32(0x587000, uintptr(66000+4*i)) == id {
			return *memmap.PtrInt32(0x587000, uintptr(66004+4*i))
		}
	}
	return -1
}
func bookAbilityID(name string) int32 {
	for i := int32(0); ; i++ {
		p := (*byte)(memmap.LoadPtr(0x587000, uintptr(69736+4*i)))
		if p == nil {
			return 0
		}
		if alloc.GoString(p) == name {
			return i
		}
	}
}
func bookGuideID(name string) int32 {
	for i := int32(0); i < 41; i++ {
		if alloc.GoString(bookGuideName(i)) == name {
			return i
		}
	}
	return 0
}
func bookGuideName(id int32) *byte { return (*byte)(memmap.LoadPtr(0x587000, uintptr(70500+4*id))) }
func bookGuideWord(id, field int32) uint32 {
	return *memmap.PtrUint32(0x5D4594, uintptr(740076+28*id+field))
}
func bookGuideCreatureName(id int32) uint32 {
	if id <= 0 || id >= 41 || bookGuideWord(id, 4) == 0 {
		return 0
	}
	return bookGuideWord(id, 0)
}
func bookGuideCharmable(typ uint32) int32 {
	for i := int32(1); i < 41; i++ {
		if v := bookGuideWord(i, 4); v != 0 && v == typ {
			return i
		}
	}
	return 0
}
func bookGuideDescription(id int32) uint32 { return bookGuideWord(id, 8) }
func bookGuideFirst() int32                { return bookGuideNext(0) }
func bookGuideNext(id int32) int32 {
	for id++; id < 41; id++ {
		if bookGuideWord(id, 4) != 0 {
			return id
		}
	}
	return 0
}
func bookGuideImage(id int32) uint32 {
	if id <= 0 || id >= 41 {
		return 0
	}
	return bookGuideWord(id, 16)
}
func bookGuideCage(id int32) uint32 {
	if id <= 0 || id >= 41 {
		return 0
	}
	return bookGuideWord(id, 12)
}
func bookGuideSize(id int32) byte { return *memmap.PtrUint8(0x5D4594, uintptr(740100+28*id)) }
func bookGuideImagePrefix(off uintptr, tail string) string {
	var buf [16]byte
	copy(buf[:8], unsafe.Slice(memmap.PtrUint8(0x587000, off), 8))
	copy(buf[8:], tail)
	return alloc.GoStringS(buf[:])
}
func bookLoadGuides() int {
	for i := int32(1); i < 41; i++ {
		name := alloc.GoString(bookGuideName(i))
		typ := GetClient().Cli().Things.TypeByID(name)
		if typ == nil {
			return 0
		}
		row := unsafe.Slice(memmap.PtrUint32(0x5D4594, uintptr(740076+28*i)), 7)
		text := GetServer().S().Strings().GetStringInFile(strman.ID("creature:"+name), "ComGuide.c")
		row[0] = memmap.PtrWord(unsafe.Pointer(alloc.InternCString16(text)))
		if alloc.GoString(typ.Name) == "Bomber" {
			row[1] = 0
		} else {
			row[1] = uint32(GetClient().Cli().Things.IndByID(name))
		}
		text = GetServer().S().Strings().GetStringInFile(strman.ID("creature_desc:"+name), "ComGuide.c")
		row[2] = memmap.PtrWord(unsafe.Pointer(alloc.InternCString16(text)))
		row[3] = memmap.PtrWord(unsafe.Pointer(Nox_xxx_gLoadImg(bookGuideImagePrefix(71248, "Cage") + name).C()))
		row[4] = memmap.PtrWord(unsafe.Pointer(Nox_xxx_gLoadImg(bookGuideImagePrefix(71264, "k") + name).C()))
		row[5] = 0
		size := byte(4)
		if typ.ObjSubClass&1 != 0 {
			size = 1
		} else if typ.ObjSubClass&2 != 0 {
			size = 2
		}
		*(*byte)(unsafe.Add(unsafe.Pointer(&row[0]), 24)) = size
	}
	return 1
}
