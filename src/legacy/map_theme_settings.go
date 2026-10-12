package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func mapThemeAlgorithm(cfg, f uint32) uint32 {
	var value [60]byte
	fields := map[string]int{"midhalllength": 4, "halllengthvariance": 8, "midhallwidth": 12, "hallwidthvariance": 16, "halllimit": 20, "hallbranchrate": 24, "hallroomrate": 28, "midroomsize": 32, "roomvariance": 36, "irregularroomrate": 40, "emptyroomrate": 44, "mergerate": 48, "adjacentportalrate": 52, "recursionlimit": 72, "seed": 76}
	for mapThemeNext(f) {
		if mapThemeLower(mapThemeText()) == "end" {
			return 1
		}
		if mapThemeRead(f, &value[0]) == 0 {
			return 0
		}
		key := mapThemeLower(mapThemeText())
		s := alloc.GoString(&value[0])
		if off, ok := fields[key]; ok {
			*populationWord(cfg, off) = uint32(mapThemeInt(s))
			continue
		}
		switch key {
		case "mapsize":
			*populationFloat(cfg, 64) = float32(mapThemeFloat(s))
		case "usedoors", "debug":
			off := 56
			if key == "debug" {
				off = 60
			}
			*populationWord(cfg, off) = 0
			if mapThemeLower(s) == "true" {
				*populationWord(cfg, off) = 1
			}
		case "skeleton":
			i := mapThemeTable(253244, s)
			if i < 0 {
				i = 0
			}
			*populationWord(cfg, 0) = uint32(i)
		}
	}
	return 0
}
func mapThemeExit(cfg, f uint32) uint32 {
	for mapThemeNext(f) {
		switch mapThemeLower(mapThemeText()) {
		case "end":
			return 1
		case "linkdata":
			if !mapThemeNext(f) {
				return 0
			}
			mapThemeCopy(cfg+476, mapThemeText())
		case "object":
			if !mapThemeNext(f) {
				return 0
			}
			p := cfg + 216 + 64**populationWord(cfg, 472)
			mapThemeCopy(p, mapThemeText())
			if !mapThemeNext(f) {
				return 0
			}
			for i, s := range []string{"north", "south", "east", "west"} {
				if mapThemeLower(mapThemeText()) == s {
					*populationWord(p, 60) = uint32(i)
					break
				}
			}
			*populationWord(cfg, 472)++
		default:
			return 0
		}
	}
	return 0
}
func mapThemeProperty(p, f uint32, op int) uint32 {
	if op == 27 {
		return mapThemeRead(f, mapThemeByte(p, 100))
	}
	if op == 28 {
		return mapThemeRead(f, mapThemeByte(p, 160))
	}
	if op == 23 {
		*mapThemeByte(p, 64) = 0
	}
	if !mapThemeNext(f) {
		return 0
	}
	s := mapThemeText()
	switch op {
	case 23:
		if mapThemeLower(s) == "none" {
			return 1
		}
		// strtok mutates only the delimiter terminating each returned token. Leading
		// and repeated separators are skipped without being overwritten.
		b := unsafe.Slice(mapThemeToken(), len(s)+1)
		for i := 0; i < len(s); {
			for i < len(s) && b[i] == '+' {
				i++
			}
			if i == len(s) {
				break
			}
			start := i
			for i < len(s) && b[i] != '+' {
				i++
			}
			end := i
			if i < len(s) {
				b[i] = 0
				i++
			}
			index := mapThemeTable(253144, string(b[start:end]))
			if index < 0 {
				return 0
			}
			*mapThemeByte(p, 64) |= byte(1 << index)
		}
	case 24:
		*mapThemeByte(p, 65) = byte(mapThemeInt(s))
	case 25:
		for i, name := range []string{"common", "uncommon", "rare", "very_rare", "hardly_ever"} {
			if mapThemeLower(s) == name {
				*populationWord(p, 72) = []uint32{1000, 500, 100, 10, 1}[i]
				break
			}
		}
	case 26:
		*populationWord(p, 76) = mapThemeBound(s, 0)
		if !mapThemeNext(f) {
			return 0
		}
		*populationWord(p, 80) = mapThemeBound(mapThemeText(), 999999)
	}
	return 1
}
func mapThemeBound(s string, wild uint32) uint32 {
	if s == "*" {
		return wild
	}
	return uint32(mapThemeInt(s))
}
func mapThemeValidate(p uint32) uint32 {
	for i := 0; i < 6; i++ {
		*populationWord(p, 8+4*i) = 0
	}
	for d := *populationWord(p, 0); d != 0; d = *populationWord(d, 220) {
		mask := *mapThemeByte(d, 64)
		for i := 0; i < 6; i++ {
			if mask == 0 || mask&(1<<i) != 0 {
				*populationWord(p, 8+4*i) += *populationWord(d, 72)
			}
		}
	}
	for i := 0; i < 6; i++ {
		if *populationWord(p, 8+4*i) == 0 {
			return 0
		}
	}
	return 1
}
func mapThemeNewPrefab(cfg, f uint32) uint32 {
	if !mapThemeNext(f) {
		return 0
	}
	p := mapThemeAlloc(1, 160)
	if p == 0 {
		return 0
	}
	mapThemeCopy(p, mapThemeText())
	*populationWord(p, 156) = *populationWord(cfg, 80)
	*populationWord(cfg, 80) = p
	*populationWord(cfg, 84)++
	return p
}
func mapThemePrefabs(cfg, f uint32) uint32 {
	var p uint32
	required := false
	for mapThemeNext(f) {
		switch mapThemeLower(mapThemeText()) {
		case "end":
			return 1
		case "must_occur":
			required = true
		case "areamap":
			p = mapThemeNewPrefab(cfg, f)
			if p == 0 {
				return 0
			}
			if required {
				*populationWord(p, 72) = 1
				required = false
			}
		case "foreach":
			if p == 0 {
				return 0
			}
			r := mapThemeForeach(f)
			if r == 0 {
				return 0
			}
			*populationWord(r, 8) = *populationWord(p, 152)
			*populationWord(p, 152) = r
		}
	}
	return 0
}
func mapThemeCleanup(cfg uint32) uint32 {
	for _, off := range []int{88, 120, 1100, 1108, 80} {
		next := 220
		if off >= 1100 {
			next = 152
		}
		if off == 80 {
			next = 156
		}
		for p := *populationWord(cfg, off); p != 0; {
			q := *populationWord(p, next)
			if off >= 1100 {
				mapThemeEquipmentFree(p)
			} else {
				mapThemeFree(p)
			}
			p = q
		}
	}
	return 0
}
