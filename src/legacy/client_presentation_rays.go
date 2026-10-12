package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"image"
)

func presentationRays() []*client.Drawable {
	return presentationRayList[:]
}
func presentationTransientRays() []*client.Drawable {
	return presentationTransientRayList[:]
}
func presentationRayAdd(event *[7]byte) {
	if int32(*effectMapped(1304312)) >= 96 {
		return
	}
	if *effectMapped(1304352) == 0 {
		for i, name := range []string{"DynamicLightning", "DynamicChainLightning", "DynamicEnergyBolt", "OrbRay", "PlasmaRay", "DrainManaRay", "HealRay", "CharmRay", "DrainManaOrb", "HealOrb", "CharmOrb", "HarpoonRope"} {
			*effectMapped(1304352 + 4*uintptr(i)) = effectType(name)
		}
	}
	fromCode, toCode := binary.LittleEndian.Uint16(event[3:]), binary.LittleEndian.Uint16(event[5:])
	c := GetClient().Cli()
	from, to := c.Objs.ByNetCode(fromCode), c.Objs.ByNetCode(toCode)
	if from == nil || to == nil {
		return
	}
	var offset uintptr
	switch event[1] {
	case 1:
		offset = 1304368
	case 2:
		offset = 1304380
	case 3:
		offset = 1304356
	case 4:
		offset = 1304360
	case 5:
		offset = 1304372
	case 6:
		offset = 1304376
	case 7:
		offset = 1304396
	case 140:
		offset = 1304352
	default:
		return
	}
	slot := -1
	for i, dr := range presentationRays() {
		if dr == nil {
			slot = i
			break
		}
	}
	if slot < 0 {
		return
	}
	pos := image.Pt(from.PosVec.X+(to.PosVec.X-from.PosVec.X)/2, from.PosVec.Y+(to.PosVec.Y-from.PosVec.Y)/2)
	dr := effectSpawn(int(*effectMapped(offset)), pos)
	if dr == nil {
		return
	}
	*effectByte(dr, 432) = 1
	*effectWord(dr, 433) = uint32(event[2])
	*effectWord(dr, 437) = uint32(fromCode)
	*effectWord(dr, 441) = uint32(toCode)
	presentationRays()[slot] = dr
}
func presentationRayRemove(event *[7]byte) {
	from, to := uint32(binary.LittleEndian.Uint16(event[3:])), uint32(binary.LittleEndian.Uint16(event[5:]))
	for i, dr := range presentationRays() {
		if dr != nil && *effectWord(dr, 437) == from && *effectWord(dr, 441) == to {
			GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			presentationRays()[i] = nil
			return
		}
	}
}
func presentationRayClear() {
	for i, dr := range presentationRays() {
		if dr != nil {
			GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			presentationRays()[i] = nil
		}
	}
}
func presentationTransientClear() {
	for i := 0; i < int(int32(*effectMapped(1304308))); i++ {
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(presentationTransientRays()[i])
	}
	*effectMapped(1304308) = 0
	objectRenderBeamReset()
}
func presentationRayContains(dr *client.Drawable) bool {
	for _, p := range presentationRays() {
		if p == dr {
			return true
		}
	}
	for i := 0; i < int(int32(*effectMapped(1304308))); i++ {
		if presentationTransientRays()[i] == dr {
			return true
		}
	}
	return false
}
