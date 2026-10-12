package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func serverPanelsByte(p *byte, mask byte, on int) *byte {
	if on != 0 {
		*p |= mask
	} else {
		*p &^= mask
	}
	return p
}
func serverPanelsWordMask(p *uint32, mask uint32, on int) *uint32 {
	if on != 0 {
		*p |= mask
	} else {
		*p &^= mask
	}
	return p
}
func serverPanelsWeaponPointer() *uint32        { return memmap.PtrUint32(0x5D4594, 1045452) }
func serverPanelsWeaponStore(p *uint32) *uint32 { *serverPanelsWeaponPointer() = *p; return p }
func serverPanelsArmorStore(v uint32) uint32    { *memmap.PtrUint32(0x5D4594, 1045456) = v; return v }
func serverPanelsArmorLoad() uint32             { return *memmap.PtrUint32(0x5D4594, 1045456) }
func serverPanelsSpellPointer() *uint32         { return memmap.PtrUint32(0x5D4594, 1045488) }
func serverPanelsSpellStore(p *uint32) {
	copy(unsafe.Slice(serverPanelsSpellPointer(), 5), unsafe.Slice(p, 5))
}
func serverPanelsWeaponQuery(mask uint32) bool {
	// C selects one byte; nonpositive inputs use the preceding byte.
	value := int32(mask)
	offset := uintptr(0)
	if value > 0 {
		for {
			next := value >> 8
			if next > 0 {
				value = next
			}
			offset++
			if next <= 0 {
				break
			}
		}
	}
	return byte(value)&*memmap.PtrUint8(0x5D4594, 1045451+offset) != 0
}
func serverPanelsArmorQuery(mask uint32) bool { return mask&serverPanelsArmorLoad() != 0 }
func serverPanelsClass(index byte) bool {
	mask := uint32(1) << uint(index&31)
	if *serverPanelsWord(1045460) != 0 {
		return serverPanelsArmorQuery(mask)
	}
	return serverPanelsWeaponQuery(mask)
}
func serverPanelsSpellBit(p *uint32, index, on int) uint32 {
	group := byte(index / 32)
	word := (*uint32)(unsafe.Add(unsafe.Pointer(p), uintptr(group)*4))
	mask := uint32(1) << uint(uint32(index)&31)
	if on != 0 {
		*word |= mask
		return mask
	}
	*word &^= mask
	return ^mask
}
func serverPanelsSpellQuery(p *uint32, index int) bool {
	word := *(*uint32)(unsafe.Add(unsafe.Pointer(p), uintptr(byte(index/32))*4))
	return word&(uint32(1)<<uint(uint32(index)&31)) != 0
}
func serverPanelsSpellSnapshot(p *uint32) uint32 {
	out := unsafe.Slice(p, 5)
	for i := range out {
		out[i] = 0xffffffff
	}
	sp := &GetServer().S().Spells
	var result uint32
	for index := 1; index <= 136; index++ {
		def := sp.DefByInd(spell.ID(index))
		result = uint32(bool2int(def.IsValid()))
		if result == 0 {
			continue
		}
		result = uint32(sp.Flags(spell.ID(index)))
		if result&0x7000000 == 0 {
			continue
		}
		result = uint32(bool2int(def.IsEnabled()))
		if result != 0 {
			continue
		}
		word := &out[index/32]
		*word &^= uint32(1) << uint(index%32)
		result = uint32(uintptr(unsafe.Pointer(word)))
	}
	return result
}
func serverPanelsSpellApply(p *uint32) int {
	sp := &GetServer().S().Spells
	result := 0
	for index := 1; index <= 136; index++ {
		result = bool2int(sp.Enable(spell.ID(index), serverPanelsSpellQuery(p, index)))
	}
	return result
}
func serverPanelsWeaponSnapshot(p *uint32) int {
	*p = 0xffffffff
	s := GetServer().S()
	result := 0
	for i := uint(0); i < 27; i++ {
		bit := uint32(1) << i
		result = int(s.Weapons.Sub_415840(bit))
		if result == 0 {
			continue
		}
		result = bool2int(s.Types.ByInd(result).Allowed())
		if result == 0 {
			*p &^= bit
			result = int(^byte(1 << uint(i&7)))
		}
	}
	return result
}
func serverPanelsArmorSnapshot() uint32 {
	s := GetServer().S()
	out := uint32(0xffffffff)
	for i := uint(0); i < 26; i++ {
		bit := uint32(1) << i
		if typ := int(s.Armor.Sub_415CD0(bit)); typ != 0 && !s.Types.ByInd(typ).Allowed() {
			out &^= bit
		}
	}
	return out
}
