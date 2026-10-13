package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/types"
)

func mapPopulationProgress(flag byte) {
	if noxflags.HasGame(noxflags.GameFlag22) {
		return
	}
	now := uint32(PlatformTicks())
	*populationGlobal(1) = now
	if *populationGlobal(2) > now {
		*populationGlobal(2) = 0
	}
	nox_input_pollEvents_4453A0()
	if now-*populationGlobal(2) > *memmap.PtrUint32(0x587000, 254948) {
		seq := *populationBlob(2487572)
		data := [3]byte{flag, byte(seq), byte(seq >> 8)}
		*populationBlob(2487572)++
		clientMapProgress(data[:])
		*populationGlobal(2) = now
	}
}
func mapPopulationExitName(object, name uint32) uint32 {
	if object == 0 || name == 0 {
		return 0
	}
	u := populationObject(object)
	typ := GetServer().S().Types.ByInd(int(u.TypeInd))
	if typ == nil || typ.Xfer != xferIdentityKey(xferIDExit) {
		return 0
	}
	dst := u.CollideDataExit().MapName[:]
	s := populationString(name)
	for i := range dst {
		if i < len(s) {
			dst[i] = s[i]
		} else {
			dst[i] = 0
		}
	}
	return 1
}
func mapPopulationRoomExit(room, definition uint32) uint32 {
	r := populationRoom(room)
	dir := int32(*populationWord(definition, 60))
	if r.Counts[dir] != 0 {
		return 0
	}
	var pos types.Pointf
	switch dir {
	case 0:
		pos.X = float32(float64(r.Size.X)*0.5 + float64(r.Min.X) + 1)
		pos.Y = float32(float64(r.Min.Y) + 10)
	case 1:
		pos.X = float32(float64(r.Size.X)*0.5 + float64(r.Min.X) + 1)
		pos.Y = float32(float64(r.Max.Y) - 10)
	case 2:
		pos.X = float32(float64(r.Max.X) - 10)
		pos.Y = float32(float64(r.Size.Y)*0.5 + float64(r.Min.Y) + 1)
	case 3:
		pos.X = float32(float64(r.Min.X) + 10)
		pos.Y = float32(float64(r.Size.Y)*0.5 + float64(r.Min.Y) + 1)
	}
	mapPaintSelectObject((*byte)(mapRoomPointer(definition)))
	u := mapPaintPlaceObject(&pos)
	if u != nil {
		var grid [2]int32
		mapRoomRound(&pos, &grid)
		var w, h float64
		switch dir {
		case 0:
			w, h = 3, 2
			grid[0]--
		case 1:
			w, h = 3, 2
			grid[0]--
			grid[1] -= 2
		case 2:
			w, h = 2, 3
			grid[0] -= 2
			grid[1]--
		case 3:
			w, h = 2, 3
			grid[1]--
		}
		pos = types.Pointf{X: float32(float64(grid[0]) * 32.526913), Y: float32(float64(grid[1]) * 32.526913)}
		mapRoomAddExclusion(r, &pos, float32(w*32.526913), float32(h*32.526913))
	}
	return mapRoomRaw(unsafe.Pointer(u))
}
func mapPopulationExit(cfg uint32) uint32 {
	if *populationWord(cfg, 472) == 0 {
		return 1
	}
	for r := mapPopulationFarthest(); r != 0; r = *populationWord(r, 68) {
		if populationRoom(r).Kind != 1 {
			continue
		}
		for i := int32(0); i < int32(*populationWord(cfg, 472)); i++ {
			if u := mapPopulationRoomExit(r, cfg+216+64*uint32(i)); u != 0 {
				mapPopulationExitName(u, cfg+476)
				return 1
			}
		}
	}
	return 0
}
func mapPopulationFinish(cfg uint32) {
	mapPopulationProgress(157)
	if mapPopulationExit(cfg) == 0 {
		name, free := alloc.CString("NoExit")
		defer free()
		file, freeFile := alloc.CString("C:\\NoxPost\\src\\Server\\MapGen\\Generate\\populate.c")
		defer freeFile()
		for line := 848; line <= 850; line++ {
			text := nox_strman_loadString_40F1D0((*int8)(unsafe.Pointer(name)), nil, (*int8)(unsafe.Pointer(file)), int(int32(line)))
			textFormatAll(0, (*uint16)(unsafe.Pointer(text)))
		}
	}
	for p := mapPopulationStart(); p != 0; p = *populationWord(p, 64) {
		mapPopulationProgress(157)
		r := populationRoom(p)
		if r.Decoration != nil && r.Flags&2 == 0 {
			mapPopulationRoom(cfg, p)
		}
		if *populationWord(cfg, 60) != 0 {
			for e := r.Exclusions; e != nil; e = e.Next {
				name := "CrystalRed"
				if e.Kind != 0 {
					name = "CrystalBlue"
				}
				str, free := alloc.CString(name)
				nox_xxx_tileGetDefByName_51D4D0((*int8)(unsafe.Pointer(str)))
				free()
				mapPaintRect(mapRoomPointer(cfg), &e.Min, int32(int64((float64(e.Max.X)-float64(e.Min.X)+0.5)*0.030743772)), int32(int64((float64(e.Max.Y)-float64(e.Min.Y)+0.5)*0.030743772)))
			}
		}
	}
	r := populationRoom(mapPopulationStart())
	pos := types.Pointf{X: float32((float64(r.Max.X) + float64(r.Min.X)) * 0.5), Y: float32((float64(r.Max.Y) + float64(r.Min.Y)) * 0.5)}
	name, free := alloc.CString("PlayerStart")
	mapPaintSelectObject((*byte)(unsafe.Pointer(name)))
	free()
	mapPaintPlaceObject(&pos)
	mapMetadataSetAmbient(*(*[3]uint32)(unsafe.Pointer(populationWord(cfg, 536))))
	mapPopulationMetadataFree()
}
