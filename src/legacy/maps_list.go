package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

type Nox_map_list_item struct {
	list      legacyListNode // 0, 0
	Name      [12]byte       // 3, 12
	Field_6   int32          // 6, 24
	Field_7   uint32         // 7, 28
	Field_8_0 byte           // 8, 32
	Field_8_1 byte           // 8, 33
	Field_8_2 uint16         // 8, 34
}

// MapName returns the map name of the item (see mapNameToken).
func (p *Nox_map_list_item) MapName() string {
	return mapNameResolve(alloc.GoString(&p.Name[0]))
}

// SetMapName sets the map name of the item.
func (p *Nox_map_list_item) SetMapName(name string) {
	alloc.StrCopyZero(p.Name[:], mapNameToken(name))
}

func Get_nox_common_maplist() unsafe.Pointer {
	return unsafe.Pointer(&legacyGlobals.nox_common_maplist)
}

func Nox_common_maplist_add_4D0760(p *Nox_map_list_item) {
	mapCatalogAdd(p)
}

func Sub_425770(p *Nox_map_list_item) {
	listInit((*legacyListNode)(unsafe.Pointer(p)))
}
