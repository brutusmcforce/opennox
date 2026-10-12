package legacy

import "unsafe"

func nox_xxx_gamedataGetFloat_419D40(k *int8) float64 {
	key := GoStringP(unsafe.Pointer(k))
	val := float64(GetServer().S().Balance.Float(key))
	return val
}

func nox_xxx_gamedataGetFloatTable_419D70(k *int8, i int) float64 {
	key := GoStringP(unsafe.Pointer(k))
	val := float64(GetServer().S().Balance.FloatInd(key, i))
	return val
}

func Nox_xxx_loadMonsterBin_517010() int {
	return monsterDefinitionLoad()
}
