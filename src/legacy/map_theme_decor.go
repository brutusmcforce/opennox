package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func mapThemeWallFloor(p, f uint32) uint32 {
	row := mapThemeAlloc(1, 128)
	if row == 0 || !mapThemeNext(f) {
		return 0
	}
	mapThemeCopy(row, mapThemeText())
	if !mapThemeNext(f) {
		return 0
	}
	mapThemeCopy(row+60, mapThemeText())
	*populationWord(row, 124) = *populationWord(p, 84)
	*populationWord(p, 84) = row
	*populationWord(p, 88)++
	return row
}
func mapThemeEdging(p uint32, kind int32, f uint32) uint32 {
	row := mapThemeAlloc(1, 128)
	if row == 0 {
		return 0
	}
	*populationWord(row, 0) = uint32(kind)
	if !mapThemeNext(f) {
		return 0
	}
	mapThemeCopy(row+4, mapThemeText())
	if !mapThemeNext(f) {
		return 0
	}
	mapThemeCopy(row+64, mapThemeText())
	mapThemeAppend(populationWord(p, 120), row, 124)
	return 1
}
func mapThemeDecorCopy(cfg, p, f uint32) uint32 {
	if !mapThemeNext(f) {
		return 0
	}
	kind := *populationWord(p, 60)
	source := p
	if kind <= 2 {
		source = *populationWord(cfg, 88+32*int(kind))
	}
	name := mapThemeText()
	for source != 0 && populationString(source) != name {
		source = *populationWord(source, 220)
	}
	if source == 0 {
		source = *populationWord(cfg, 152)
		for source != 0 && populationString(source) != name {
			source = *populationWord(source, 220)
		}
	}
	if source == 0 {
		return 0
	}
	for q := *populationWord(source, 84); q != 0; q = *populationWord(q, 124) {
		row := mapThemeAlloc(1, 128)
		if row == 0 {
			return 0
		}
		copy(unsafe.Slice(mapThemeByte(row, 0), 128), unsafe.Slice(mapThemeByte(q, 0), 128))
		*populationWord(row, 124) = *populationWord(p, 84)
		*populationWord(p, 84) = row
		*populationWord(p, 88)++
	}
	tail := *populationWord(p, 92)
	if tail != 0 {
		for *populationWord(tail, 20) != 0 {
			tail = *populationWord(tail, 20)
		}
	}
	for q := *populationWord(source, 92); q != 0; {
		row := mapThemeAlloc(1, 24)
		if row == 0 {
			return 0
		}
		copy(unsafe.Slice(mapThemeByte(row, 0), 24), unsafe.Slice(mapThemeByte(q, 0), 24))
		*populationWord(row, 16) = 1
		*populationWord(row, 20) = 0
		count := uint32(1)
		if tail != 0 {
			*populationWord(tail, 20) = row
			count = *populationWord(p, 96) + 1
		} else {
			*populationWord(p, 92) = row
		}
		*populationWord(p, 96) = count
		q = *populationWord(q, 20)
		tail = row
	}
	return 1
}
func mapThemeDecor(cfg, f uint32) uint32 {
	p := mapThemeAlloc(1, 224)
	if p == 0 {
		return 0
	}
	*populationWord(p, 72) = 1000
	*populationWord(p, 80) = 999999
	if !mapThemeNext(f) || mapThemeRead(f, mapThemeByte(p, 0)) == 0 {
		return 0
	}
	kind := -1
	for i, s := range []string{"room", "hall", "template", "backdrop"} {
		if mapThemeLower(mapThemeText()) == s {
			kind = i
			break
		}
	}
	if kind < 0 {
		return 0
	}
	*populationWord(p, 60) = uint32(kind)
	var wall uint32
	for mapThemeNext(f) {
		key := mapThemeLower(mapThemeText())
		switch key {
		case "end":
			off := 88 + 32*kind
			*populationWord(p, 220) = *populationWord(cfg, off)
			*populationWord(cfg, off) = p
			*populationWord(cfg, off+4)++
			return 1
		case "wall_floor":
			wall = mapThemeWallFloor(p, f)
			if wall == 0 {
				return 0
			}
		case "must_occur":
			*mapThemeByte(p, 67) = 1
		case "copy":
			if mapThemeDecorCopy(cfg, p, f) == 0 {
				return 0
			}
		case "decor_set":
			if mapThemeDecorSet(p, f) == 0 {
				return 0
			}
		case "occur_constraint", "occur_limit", "frequency", "room_size_constraint", "door", "double_door":
			op := map[string]int{"occur_constraint": 23, "occur_limit": 24, "frequency": 25, "room_size_constraint": 26, "door": 27, "double_door": 28}[key]
			if mapThemeProperty(p, f, op) == 0 {
				return 0
			}
		default:
			found := false
			// Each edging read replaces the shared token. The original table
			// walk compares subsequent entries against that updated token.
			for i := 0; ; i++ {
				name := *memmap.PtrUint32(0x587000, 253200+uintptr(4*i))
				if name == 0 {
					break
				}
				if mapThemeLower(populationString(name)) == mapThemeLower(mapThemeText()) {
					if wall == 0 || mapThemeEdging(wall, int32(i), f) == 0 {
						return 0
					}
					found = true
				}
			}
			if !found {
				return 0
			}
		}
	}
	return 0
}
