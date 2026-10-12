package legacy

import (
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"strings"
	"unsafe"
)

func quickbarSaveRows(b *quickbarRecord, rows, slots int, class byte) int {
	if rows <= 0 {
		return 1
	}
	file := cryptfile.Global()
	for row := 0; row < rows; row++ {
		for slot := 0; slot < slots; slot++ {
			s := (*quickbarSlot)(unsafe.Add(unsafe.Pointer(b), 40*row+8*slot))
			var size [1]byte
			if file.ReadOnly() {
				file.ReadWrite(size[:])
				var name [256]byte
				file.ReadWrite(name[:int(size[0])])
				// The live spell parser is Go; the shared ability-name owner remains C.
				if class != 0 {
					text, _, _ := strings.Cut(string(name[:size[0]]), "\x00")
					s.ID = uint32(spell.ParseID(text))
				} else {
					s.ID = uint32(bookAbilityID(alloc.GoStringS(name[:size[0]])))
				}
			} else {
				var name string
				if class != 0 {
					name = spell.ID(int32(s.ID)).String()
				} else {
					name = server.Ability(s.ID).String()
				}
				size[0] = byte(len(name))
				file.ReadWrite(size[:])
				file.ReadWrite([]byte(name)[:int(size[0])])
			}
			file.ReadWrite(unsafe.Slice(quickbarFlagByte(s), 1))
		}
	}
	return 1
}
func quickbarSave() int {
	file := cryptfile.Global()
	if !file.ReadOnly() && *quickbarWord(1049688) == 1 {
		quickbarRestoreSlots()
		*quickbarWord(1049688) = 0
	}
	var class [1]byte
	if p := quickbarPlayer(); p != 0 {
		class[0] = byte(bookClass(p))
	} else {
		class[0] = *(*byte)(unsafe.Add(unsafe.Pointer(Nox_xxx_getHostInfoPtr_431770()), 66))
	}
	file.ReadWrite(class[:])
	quickbarSaveRows(quickbarMain(), 5, 5, class[0])
	if class[0] != 0 {
		quickbarSaveRows(quickbarAt(1047940), 3, 3, class[0])
		if quickbarMain().Directions[0] != nil {
			quickbarDirections(quickbarMain())
		}
	}
	return 1
}
