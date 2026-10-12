package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func playerCorpseCache() {
	s := GetServer().S()
	parts := [...]string{"Skull", "RibCage", "Pelvis", "LeftLowerLeg", "LeftUpperLeg", "LeftLowerArm", "LeftUpperArm", "RightLowerLeg", "RightUpperLeg", "RightLowerArm", "RightUpperArm"}
	suffixIndex := uintptr(0)
	for direction := uintptr(0); direction < 9; direction++ {
		if direction == 4 {
			continue
		}
		suffix := alloc.GoString((*byte)(memmap.PtrOff(0x587000, 281216+4*suffixIndex)))
		suffixIndex++
		for part, name := range parts {
			*memmap.PtrUint32(0x5D4594, 2488740+44*direction+4*uintptr(part)) = uint32(s.Types.IndByID("Corpse" + name + suffix))
		}
	}
	*memmap.PtrUint32(0x5D4594, 2488736) = 1
}
func playerCorpseCreate(pos types.Pointf, angle int32) {
	if *memmap.PtrUint32(0x5D4594, 2488736) == 0 {
		playerCorpseCache()
	}
	s := GetServer().S()
	direction := uintptr(uint32(geometryDirection4Index(angle)))
	for part := uintptr(0); part < 11; part++ {
		typ := *memmap.PtrUint32(0x5D4594, 2488740+44*direction+4*part)
		u := s.NewObjectByTypeInd(int(typ))
		if u == nil {
			break
		}
		if Get_dword_5d4594_2650652() != 0 && noxflags.HasGame(noxflags.GameOnline) {
			u.ObjFlags |= 0x40
		}
		x := *memmap.PtrFloat32(0x587000, 280376+88*direction+8*part) + pos.X
		y := *memmap.PtrFloat32(0x587000, 280380+88*direction+8*part) + pos.Y
		GetServer().CreateObjectAt(u, nil, types.Pointf{X: x, Y: y})
		delay := s.TickRate() * uint32(s.Rand.Logic.IntClamp(10, 20))
		motionDecaySet(u, int32(delay))
	}
}
