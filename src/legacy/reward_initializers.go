package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func rewardInitSpark(u *server.Object) unsafe.Pointer {
	*equipmentWord(u.UpdateData, 0) = 32
	*equipmentWord(u.UpdateData, 4) = 32
	return u.UpdateData
}
func rewardInitFrog(u *server.Object) int32 {
	data := unsafe.Slice((*byte)(u.UpdateData), 3)
	data[0] = byte(rewardRoll(55, 60))
	data[1] = 1
	data[2] = 0
	dir := rewardRoll(0, 255)
	u.Direction2 = server.Dir16(uint16(dir))
	return dir
}
func rewardInitBoulder(u *server.Object) *server.Object {
	u.Pos39 = u.PosVec
	return u
}

func rewardInitBreakable(u *server.Object) {
	if u.Field5&14 == 0 {
		u.SetXStatus(2)
	}
}
func rewardInitDirection(u *server.Object, name bool) int32 {
	dir := int32(int32(geometryDirectionAngle((*[2]uint32)(unsafe.Pointer(u.InitData)))))
	u.Direction1 = server.Dir16(uint16(dir))
	u.Direction2 = u.Direction1
	if !name {
		return dir
	}
	typ := GetServer().S().Types.IndByID(alloc.GoString((*byte)(unsafe.Add(u.UpdateData, 16))))
	*equipmentWord(u.UpdateData, 12) = uint32(typ)
	return int32(typ)
}
func rewardInitGold(u *server.Object) int32 {
	if *equipmentWord(u.InitData, 0) != 0 {
		return int32(memmap.PtrWord(u.CObj()))
	}
	players := &GetServer().S().Players
	var sum float32
	var count int32
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if pl.PlayerUnit != nil {
			sum += pl.PlayerUnit.Experience
		}
		count++
	}
	average := float32(float64(sum) / float64(count))
	lo := effectsTruncWord(float64(average) * math.Float64frombits(uint64(qword_581450_10256)))
	hi := effectsTruncWord(float64(average) * memmap.Float64(0x581450, 10264))
	extra := rewardRoll(lo, hi) - effectsTruncWord(float64(average)*memmap.Float64(0x581450, 10248))
	result := rewardRoll(15, 30)
	*equipmentWord(u.InitData, 0) = uint32(result + extra)
	return result
}
func rewardInitGenerator(u *server.Object) int32 {
	stage := int32(memmap.Uint32(0x5d4594, 2388660))
	level := *(*byte)(unsafe.Add(u.UpdateData, 83+int(stage)))
	names := [4]string{"GeneratorMaxActiveCreaturesHigh", "GeneratorMaxActiveCreaturesNormal", "GeneratorMaxActiveCreaturesLow", "GeneratorMaxActiveCreaturesSingular"}
	if level < 4 {
		*(*byte)(unsafe.Add(u.UpdateData, 87)) = byte(effectsTruncWord(GetServer().S().Balance.Float(names[level])))
	}
	result := int32(u.ObjSubClass)
	for i, dir := range []int{0, 2, 8, 6} {
		if u.ObjSubClass&(1<<uint(i)) != 0 {
			result = int32(int32(geometryDirection4Angle(int32(dir))))
			u.Direction1 = server.Dir16(uint16(result))
			break
		}
	}
	u.Direction2 = u.Direction1
	return result
}
func rewardPlaceAnkh() {
	if dword_5d4594_1568280 == 0 {
		dword_5d4594_1568280 = uint32(GetServer().S().Types.IndByID("RewardMarker"))
		*memmap.PtrUint32(0x5d4594, 1568284) = uint32(GetServer().S().Types.IndByID("RewardMarkerPlus"))
	}
	eligible := func(u *server.Object) bool {
		return (uint32(u.TypeInd) == uint32(dword_5d4594_1568280) || uint32(u.TypeInd) == memmap.Uint32(0x5d4594, 1568284)) && *(*byte)(u.InitData)&0x80 != 0
	}
	var count int32
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if eligible(u) {
			count++
		}
	}
	pick := rewardRoll(0, count-1)
	var i int32
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if eligible(u) {
			if i == pick {
				if it := GetServer().S().NewObjectByTypeID("Ankh"); it != nil {
					GetServer().CreateObjectAt(it, nil, u.PosVec)
					GetServer().DelayedDelete(u)
					return
				}
			} else {
				i++
			}
		}
	}
}
func rewardSelectMarkers() {
	stage := int32(rewardWord(202028))
	players := int32(questRuntimeCount())
	if dword_5d4594_1568288 == 0 {
		dword_5d4594_1568288 = uint32(GetServer().S().Types.IndByID("RewardMarker"))
		*memmap.PtrUint32(0x5d4594, 1568292) = uint32(GetServer().S().Types.IndByID("RewardMarkerPlus"))
		*memmap.PtrUint32(0x5d4594, 1568296) = uint32(GetServer().S().Types.IndByID("RedPotion"))
	}
	var fraction float32
	if stage == 1 {
		fraction = .5
	} else {
		switch players {
		case 1, 2:
			fraction = .40000001
		case 3, 4:
			fraction = .69999999
		case 5, 6:
			fraction = 1
		}
	}
	var markers, potions []*server.Object
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if uint32(u.TypeInd) == uint32(dword_5d4594_1568288) {
			if *(*byte)(unsafe.Add(u.InitData, 216))&1 == 0 {
				markers = append(markers, u)
			}
		} else if uint32(u.TypeInd) == memmap.Uint32(0x5d4594, 1568296) {
			potions = append(potions, u)
		}
	}
	if len(markers) == 0 && len(potions) == 0 {
		return
	}
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if uint32(u.TypeInd) == uint32(dword_5d4594_1568288) {
			flags := (*byte)(unsafe.Add(u.InitData, 216))
			if *flags&1 != 0 {
				*flags |= 0x80
			}
		} else if uint32(u.TypeInd) == memmap.Uint32(0x5d4594, 1568292) {
			*(*byte)(unsafe.Add(u.InitData, 216)) |= 0x80
		}
	}
	shuffle := func(items []*server.Object) {
		for i := len(items) - 1; i > 0; i-- {
			j := int(rewardRoll(0, int32(i)))
			items[j], items[i] = items[i], items[j]
		}
	}
	if len(markers) != 0 {
		n := effectsTruncWord(float64(len(markers))*float64(fraction) + .5)
		shuffle(markers)
		for i := int32(0); i < n; i++ {
			*(*byte)(unsafe.Add(markers[i].InitData, 216)) |= 0x80
		}
	}
	if len(potions) != 0 {
		n := uint32(effectsTruncWord(float64(len(potions))*float64(fraction) + .5))
		shuffle(potions)
		for i, u := range potions {
			if uint32(i) >= n {
				GetServer().DelayedDelete(u)
			}
		}
	}
}
