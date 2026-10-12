package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func effectDispatchRay(packet *[9]byte) uint32 {
	count := *effectMapped(1304308)
	if int32(count) >= 96 {
		return count
	}
	if *effectMapped(1304316) == 0 {
		for _, v := range []struct {
			offset uintptr
			name   string
		}{
			{1304316, "DynamicLightning"}, {1304320, "DynamicChainLightning"}, {1304324, "DynamicEnergyBolt"},
			{1304348, "GreenZap"}, {1304332, "PlasmaRay"}, {1304336, "DrainManaOrb"}, {1304340, "HealOrb"}, {1304344, "CharmOrb"},
		} {
			*effectMapped(v.offset) = effectType(v.name)
		}
		dword_5d4594_1304328 = uint32(effectType("OrbRay"))
	}
	var coords [4]uint16
	for i := range coords {
		coords[i] = binary.LittleEndian.Uint16(packet[1+2*i:])
	}
	midpoint := image.Pt(int(coords[0])+(int(coords[2])-int(coords[0]))/2, int(coords[1])+(int(coords[3])-int(coords[1]))/2)
	var typ uint32
	switch packet[0] {
	case 0x7d:
		typ = *effectMapped(1304332)
	case 0x8c:
		typ = *effectMapped(1304316)
	case 0x8d:
		typ = *effectMapped(1304324)
	case 0x8e:
		typ = *effectMapped(1304320)
	case 0x8f, 0x90, 0x91:
		speed := byte(effectRand(6, 12))
		typ = uint32(dword_5d4594_1304328)
		offset := uintptr(1304336)
		if packet[0] == 0x90 {
			offset = 1304344
		} else if packet[0] == 0x91 {
			offset = 1304340
		}
		emit := func(p *[4]uint16) {
			if effectRand(0, 100) < 50 {
				// Original call evaluates the Y random sample before X.
				y := effectRand(-20, 20)
				x := effectRand(-20, 20)
				effectCreateOrb(int(*effectMapped(offset)), p, x, y, speed, 0)
			}
		}
		emit(&coords)
		if packet[0] == 0x90 {
			reverse := [4]uint16{coords[2], coords[3], coords[0], coords[1]}
			emit(&reverse)
		}
	default:
		return uint32(int(packet[0]) - 125)
	}
	dr := effectSpawn(int(typ), midpoint)
	if dr == nil {
		return 0
	}
	*effectByte(dr, 432) = 0
	*effectWord(dr, 437) = binary.LittleEndian.Uint32(packet[1:])
	*effectWord(dr, 441) = binary.LittleEndian.Uint32(packet[5:])
	count = *effectMapped(1304308)
	*(**client.Drawable)(memmap.PtrOff(0x5D4594, 1303540+uintptr(4*count))) = dr
	*effectMapped(1304308) = count + 1
	return uint32(uintptr(unsafe.Pointer(dr)))
}
