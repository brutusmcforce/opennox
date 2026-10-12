package legacy

import (
	"strings"
	"unsafe"
)

func mapThemeEquipmentFree(p uint32) {
	for i := 0; i < 4; i++ {
		if *populationWord(p, 136+4*i) != 0 {
			mapThemeFree(*populationWord(p, 120+4*i))
		}
	}
	mapThemeFree(p)
}
func mapThemeAttributes(p, f uint32) uint32 {
	var slots [4][]string
	for {
		if !mapThemeNext(f) {
			return 0
		}
		key := mapThemeLower(mapThemeText())
		if key == "end" {
			break
		}
		slot := -1
		switch key {
		case "quality", "effectiveness":
			slot = 0
		case "material":
			slot = 1
		case "primary_enchantment":
			slot = 2
		case "secondary_enchantment":
			slot = 3
		}
		if slot < 0 {
			return 0
		}
		for {
			if !mapThemeNext(f) {
				return 0
			}
			if mapThemeLower(mapThemeText()) == "end" {
				break
			}
			slots[slot] = append(slots[slot], mapThemeText())
		}
	}
	for i, changes := range slots {
		values := changes
		if template := *mapThemeTemplate(); template != 0 {
			values = nil
			for j := int32(0); j < int32(*populationWord(template, 136+4*i)); j++ {
				values = append(values, populationString(*populationWord(template, 120+4*i)+uint32(60*j)))
			}
			for _, value := range changes {
				remove := strings.HasPrefix(value, "-")
				match := value
				if remove {
					match = match[1:]
				}
				index := len(values)
				for j, old := range values {
					if mapThemeLower(old) == mapThemeLower(match) {
						index = j
						break
					}
				}
				if remove {
					if index < len(values) {
						values = append(values[:index], values[index+1:]...)
					}
				} else if index == len(values) {
					values = append(values, value)
				}
			}
		}
		if len(values) != 0 {
			a := mapThemeAlloc(uint32(len(values)), 60)
			*populationWord(p, 120+4*i) = a
			if a == 0 {
				return 0
			}
			for j, value := range values {
				mapThemeCopy(a+uint32(60*j), value)
			}
		}
		*populationWord(p, 136+4*i) = uint32(len(values))
	}
	return 1
}
func mapThemeEquipment(cfg, f uint32, armor bool) uint32 {
	*mapThemeTemplate() = 0
	kind, off := "weapon", 1100
	if armor {
		kind, off = "armor", 1108
	}
	for {
		if !mapThemeNext(f) {
			return 0
		}
		key := mapThemeLower(mapThemeText())
		if key == "end" {
			break
		}
		if key != kind {
			continue
		}
		p := mapThemeAlloc(1, 156)
		if p == 0 || !mapThemeNext(f) {
			return 0
		}
		mapThemeCopy(p+60, mapThemeText())
		if mapThemeLower(populationString(p+60)) != "template" {
			if !mapThemeNext(f) {
				return 0
			}
			mapThemeCopy(p, mapThemeText())
		}
		if mapThemeAttributes(p, f) == 0 {
			return 0
		}
		if mapThemeLower(populationString(p+60)) == "template" {
			*mapThemeTemplate() = p
		} else {
			*populationWord(p, 152) = 0
			mapThemeAppend(populationWord(cfg, off), p, 152)
			*populationWord(cfg, off+4)++
		}
	}
	if p := *mapThemeTemplate(); p != 0 {
		mapThemeEquipmentFree(p)
	}
	*mapThemeTemplate() = 0
	return 1
}
func mapThemeSpells(cfg, f uint32) uint32 {
	for {
		if !mapThemeNext(f) {
			return 0
		}
		if mapThemeLower(mapThemeText()) == "end" || *populationWord(cfg, 1096) >= 137 {
			return 1
		}
		if id := mapPopulationSpellID(mapRoomRaw(unsafe.Pointer(mapThemeToken()))); id != 0 {
			count := populationWord(cfg, 1096)
			*populationWord(cfg, 548+int(*count)*4) = id
			*count++
		}
	}
}
