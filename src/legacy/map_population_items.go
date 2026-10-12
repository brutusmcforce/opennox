package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/common/memmap"
	"strings"
	"unsafe"
)

func mapPopulationSpellID(name uint32) uint32 {
	s := populationString(name)
	if len(s) >= 60 {
		return 0
	}
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - ('a' - 'A')
		}
	}
	id := spell.ParseID("SPELL_" + string(b))
	if id <= 0 {
		return 0
	}
	return uint32(id)
}
func mapPopulationAttach(owner, item uint32) uint32 {
	u := populationObject(item)
	o := populationObject(owner)
	u.Field125 = nil
	u.InvNextItem = o.InvFirstItem
	if o.InvFirstItem != nil {
		o.InvFirstItem.Field125 = u
	}
	o.InvFirstItem = u
	u.InvHolder = o
	return item
}
func mapPopulationInventory(cfg, owner, groups uint32) {
	if groups == 0 {
		return
	}
	roll := mapRoomRandomInt(1, 100)
	for {
		roll -= int32(*populationWord(groups, 0))
		if roll <= 0 {
			break
		}
		groups = *populationWord(groups, 2056)
		if groups == 0 {
			return
		}
	}
	var item uint32
	for i := uint32(0); i < *populationWord(groups, 2052); i++ {
		row := groups + 4 + 64*i
		name := row + 4
		switch *populationWord(row, 0) {
		case 0:
			item = mapRoomRaw(unsafe.Pointer(GetServer().S().NewObjectByTypeID(populationString(name))))
		case 3:
			item = mapPopulationItem(name, *populationWord(cfg, 1100), *populationWord(cfg, 1104))
		case 4:
			item = mapPopulationItem(name, *populationWord(cfg, 1108), *populationWord(cfg, 1112))
		case 5:
			item = mapPopulationSpellbook(cfg, name)
		}
		if item != 0 {
			mapPopulationAttach(owner, item)
		}
	}
}
func mapPopulationSpellbook(cfg, name uint32) uint32 {
	count := *populationWord(cfg, 1096)
	if count == 0 {
		return 0
	}
	u := GetServer().S().NewObjectByTypeID("SpellBook")
	if u == nil {
		return 0
	}
	var id byte
	if populationString(name) == "*" {
		id = byte(*populationWord(cfg, 548+4*int(mapRoomRandomInt(0, int32(count)-1))))
	} else {
		id = byte(mapPopulationSpellID(name))
		if id == 0 {
			GetServer().S().Objs.FreeObject(u)
			return 0
		}
	}
	mapPaintFinishBook(u, id)
	return mapRoomRaw(unsafe.Pointer(u))
}
func mapPopulationItem(name, definitions, count uint32) uint32 {
	id := populationString(name)
	def := definitions
	typ := id
	if id == "*" {
		index := mapRoomRandomInt(0, int32(count)-1)
		for i := int32(0); def != 0 && i != index; i++ {
			def = *populationWord(def, 152)
		}
		typ = populationString(def + 60)
	} else {
		for def != 0 && !strings.EqualFold(populationString(def), id) {
			def = *populationWord(def, 152)
		}
		if def != 0 {
			typ = populationString(def + 60)
		}
	}
	u := GetServer().S().NewObjectByTypeID(typ)
	if u == nil {
		return 0
	}
	if def != 0 {
		var attrs [5]uint32
		for i := 0; i < 4; i++ {
			n := *populationWord(def, 136+4*i)
			if n != 0 && mapRoomRandomInt(1, 100) <= int32(*memmap.PtrUint32(0x587000, 254688+uintptr(4*i))) {
				mod := populationString(*populationWord(def, 120+4*i) + 60*uint32(mapRoomRandomInt(0, int32(n)-1)))
				if !strings.EqualFold(mod, "none") {
					ind := GetServer().S().Modif.Nox_xxx_modifGetIdByName413290(mod)
					attrs[i] = mapRoomRaw(GetServer().S().Modif.Nox_xxx_modifGetDescById413330(ind).C())
				}
			}
		}
		if attrs[2] == 0 {
			attrs[3] = 0
		}
		stateAttributes(u, unsafe.Pointer(&attrs[0]))
	}
	return mapRoomRaw(unsafe.Pointer(u))
}
