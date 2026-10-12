package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

const tooltipSource = `C:\NoxPost\src\client\Gui\ToolTip.c`

func tooltipString(id string) *uint16 {
	return alloc.InternCString16(GetServer().S().Strings().GetStringInFile(strman.ID(id), tooltipSource))
}
func uiCursorTooltip(text *uint16) {
	dst := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1096676)), 256)
	if text == nil {
		dst[0] = 0
		return
	}
	alloc.StrCopyZero16P(dst, text)
}
func uiItemTooltip(dr *client.Drawable) *uint16 {
	dst := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1317000)), 1024)
	dst[0] = 0
	alloc.StrCopyZero16P(dst, (*uint16)(memmap.PtrOff(0x5D4594, 1319048)))
	if dr == nil {
		return &dst[0]
	}
	appendText := func(text *uint16) {
		n := alloc.StrLenS(dst)
		alloc.StrCopyZero16P(dst[n:], text)
	}
	space := alloc.InternCString16(" ")
	class, sub := uint32(dr.Class()), uint32(dr.SubClass())
	item := (*client.DrawableUnionItem)(unsafe.Pointer(&dr.Union))
	if class&0x13001000 == 0 {
		if class&0x100 == 0 || sub&7 == 0 {
			if p := GetClient().Cli().Things.TypeByInd(int(dr.TypeIDVal)); p != nil && p.PrettyName != nil {
				return p.PrettyName
			}
			return &dst[0]
		}
		metadata := item.Field_108
		kind, sentinel := byte(4), uint32(6)
		if sub&1 != 0 {
			kind, sentinel = 1, 137
		} else if sub&2 != 0 {
			kind, sentinel = 2, 41
		}
		if metadata == 0 {
			var request [4]byte
			request[0], request[3] = 0xe2, kind
			binary.LittleEndian.PutUint16(request[1:3], uint16(nox_xxx_netGetUnitCodeCli_578B00(int32(memmap.PtrWord(unsafe.Pointer(dr))))))
			item.Field_108 = sentinel
			Nox_xxx_netClientSend2_4E53C0(31, unsafe.Pointer(&request[0]), 4, 0, 1)
			return &dst[0]
		}
		if metadata == sentinel {
			return &dst[0]
		}
		lang := GetServer().S().Strings().Lang()
		var title *uint16
		switch kind {
		case 1:
			title = (*uint16)(unsafe.Pointer(nox_xxx_spellTitle_424930(int(metadata))))
		case 2:
			title = (*uint16)(memmap.WordPtr(uint32(bookGuideCreatureName(int32(metadata)))))
		case 4:
			title = alloc.InternCString16(Nox_xxx_abilityGetName_0_425260(int(metadata)))
		}
		key := "BookOf"
		prefix := lang != 6
		if kind == 2 {
			key = "LoreScroll"
			prefix = lang == 3 || lang == 5
		}
		if prefix {
			appendText(tooltipString(key))
			appendText(space)
			appendText(title)
		} else {
			// C's lore prefix is formatted with %s, whose nil case is "(null)".
			if kind == 2 && title == nil {
				title = alloc.InternCString16("(null)")
			}
			appendText(title)
			appendText(space)
			appendText(tooltipString(key))
		}
		return &dst[0]
	}
	var def *server.Modifier
	if class&0x11001000 != 0 {
		def = GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
	} else {
		def = GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
	}
	if def == nil {
		textFormatBuffer(dst, tooltipString("NoArmsInfo"), textFormatPointer(unsafe.Pointer(nox_get_thing_name(int(dr.TypeIDVal)))))
		return &dst[0]
	}
	var mods [4]*uint16
	if class&0x1000000 == 0 || sub&0x7800000 == 0 {
		for i, word := range [4]uint32{item.Field_108, item.Field_109, item.Field_110, item.Field_111} {
			if word == 0 {
				continue
			}
			mod := (*server.ModifierEff)(memmap.WordPtr(word))
			if i == 3 {
				mods[i] = mod.SecondaryDescPtr()
			} else {
				mods[i] = mod.DescPtr()
			}
		}
	}
	before := func(i int) {
		if mods[i] != nil {
			appendText(mods[i])
			appendText(space)
		}
	}
	after := func(i int) {
		if mods[i] != nil {
			appendText(space)
			appendText(mods[i])
		}
	}
	switch GetServer().S().Strings().Lang() {
	case 2:
		appendText(def.Desc8)
		after(2)
		after(3)
		before(1)
		after(0)
	case 3:
		before(0)
		appendText(def.Desc8)
		after(1)
		after(2)
		after(3)
	case 5:
		appendText(def.Desc8)
		after(1)
		after(0)
		after(2)
		after(3)
	case 6:
		before(0)
		before(2)
		before(3)
		after(1)
		appendText(def.Desc8)
	default:
		before(0)
		before(1)
		appendText(def.Desc8)
		after(2)
		after(3)
	}
	return &dst[0]
}
