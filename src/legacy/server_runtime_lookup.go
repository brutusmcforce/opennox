package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func runtimeLookupString(id, file string) *uint16 {
	return alloc.InternCString16(GetServer().S().Strings().GetStringInFile(strman.ID(id), file))
}
func runtimeEquipmentBase(armor bool) uintptr {
	if armor {
		return 35496
	}
	return 33392
}
func runtimeEquipmentLoad(armor bool) {
	off, ready, file := runtimeEquipmentBase(armor), uintptr(371248), `C:\NoxPost\src\common\Object\WeapLook.c`
	if armor {
		ready, file = 371256, `C:\NoxPost\src\common\Object\ArmrLook.c`
	}
	flag := memmap.PtrUint32(0x5D4594, ready)
	if *flag != 0 {
		return
	}
	for ; memmap.Uint32(0x587000, off+4) != 0; off += 12 {
		id := alloc.GoString((*byte)(unsafe.Pointer(uintptr(memmap.Uint32(0x587000, off+4)))))
		*memmap.PtrUint32(0x587000, off) = uint32(uintptr(unsafe.Pointer(runtimeLookupString(id, file))))
	}
	*flag = 1
}

// The hosted C locale folds ASCII letters; Unicode case folding would broaden
// the accepted equipment names. Compare UTF-16 units through their terminator.
func runtimeEquipmentNameEqual(a, b *uint16) bool {
	for i := uintptr(0); ; i += 2 {
		x, y := *(*uint16)(unsafe.Add(unsafe.Pointer(a), i)), *(*uint16)(unsafe.Add(unsafe.Pointer(b), i))
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
		if x == 0 {
			return true
		}
	}
}
func runtimeEquipmentMask(armor bool, name *uint16) uint32 {
	for off := runtimeEquipmentBase(armor); memmap.Uint32(0x587000, off) != 0; off += 12 {
		p := (*uint16)(unsafe.Pointer(uintptr(memmap.Uint32(0x587000, off))))
		if runtimeEquipmentNameEqual(name, p) {
			return memmap.Uint32(0x587000, off+8)
		}
	}
	return 0
}
func runtimeEquipmentLabel(armor bool, mask uint32) *uint16 {
	for off := runtimeEquipmentBase(armor); memmap.Uint32(0x587000, off) != 0; off += 12 {
		if memmap.Uint32(0x587000, off+8) == mask {
			return (*uint16)(unsafe.Pointer(uintptr(memmap.Uint32(0x587000, off))))
		}
	}
	return nil
}
func runtimeModifierRow(mask byte) uintptr {
	for off := uintptr(27332); off < 27452; off += 20 {
		if int(int8(mask)) == int(memmap.Uint8(0x587000, off)) {
			return off
		}
	}
	return 0
}
func runtimeModifierIcon(mask byte) uint32 {
	ready := memmap.PtrUint32(0x5D4594, 251624)
	if *ready == 0 {
		for off := uintptr(27332); off < 27452; off += 20 {
			id := alloc.GoString((*byte)(unsafe.Pointer(uintptr(memmap.Uint32(0x587000, off+4)))))
			*memmap.PtrUint32(0x587000, off+8) = uint32(uintptr(Nox_xxx_gLoadImg(id).C()))
		}
		*ready = 1
	}
	if off := runtimeModifierRow(mask); off != 0 {
		return memmap.Uint32(0x587000, off+8)
	}
	return 0
}
func runtimeModifierLabel(mask byte) *uint16 {
	off := runtimeModifierRow(mask)
	if off == 0 {
		return nil
	}
	id := alloc.GoString((*byte)(unsafe.Pointer(uintptr(memmap.Uint32(0x587000, off+12)))))
	return runtimeLookupString(id, `C:\NoxPost\src\common\Object\Modifier.c`)
}
func runtimeMaterial(u *server.Object) bool {
	cache := memmap.PtrUint32(0x5D4594, 251620)
	if *cache == 0 {
		m := &GetServer().S().Modif
		*cache = uint32(uintptr(m.Nox_xxx_modifGetDescById413330(m.Nox_xxx_modifGetIdByName413290("Material7")).C()))
	}
	return uint32(u.ObjClass)&0x13001000 != 0 && *(*uint32)(unsafe.Add(u.InitData, 4)) == *cache
}
func runtimeArmorConductivity(u *server.Object) float64 {
	if u.Class().Has(object.ClassArmor) {
		if m := GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd)); m != nil {
			return float64(m.DamageCoeffOrArmor64)
		}
	}
	return 0
}
