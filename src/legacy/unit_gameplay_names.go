package legacy

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"strings"
	"unsafe"
)

func unitNPCName(u *server.Object) *uint16 {
	name := alloc.GoString((*byte)(u.IDPtr))
	if u.IDPtr == nil {
		name = GetServer().S().Types.ByInd(int(u.TypeInd)).ID()
	}
	name, _, _ = strings.Cut(name, "\x00")
	if i := strings.LastIndexByte(name, ':'); i >= 0 {
		name = name[i+1:]
	}
	key := "NPC:" + strings.ReplaceAll(name, "_", "")
	alloc.StrCopy(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1563460)), 512), key)
	return alloc.InternCString16(GetServer().S().Strings().GetStringInFile(strman.ID(key), "C:\\NoxPost\\src\\Server\\DBase\\objdb.c"))
}
func unitItemName(u *server.Object) *uint16 {
	dst := unsafe.Slice(memmap.PtrUint16(0x5D4594, 1565660), 1036)
	core := GetServer().S()
	lookup := func(key string) *uint16 {
		return alloc.InternCString16(core.Strings().GetStringInFile(strman.ID(key), "C:\\NoxPost\\src\\Server\\Object\\objutil.c"))
	}
	if u.ObjClass&0x13001000 == 0 {
		alloc.StrCopy16P(dst, lookup("NoDescription"))
		return &dst[0]
	}
	var base *server.Modifier
	if u.ObjClass&0x11001000 != 0 {
		base = core.Modif.Nox_xxx_getProjectileClassById413250(int(u.TypeInd))
	} else {
		base = core.Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd))
	}
	if base == nil {
		textFormatBuffer(dst, lookup("NoInfo"), textFormatPointer(unsafe.Pointer(internCStr(core.Types.ByInd(int(u.TypeInd)).ID()))))
		return &dst[0]
	}
	var out []uint16
	appendText := func(p *uint16) { out = append(out, unsafe.Slice(p, alloc.StrLen(p))...) }
	appendText(memmap.PtrUint16(0x5D4594, 1567732))
	mods := (*[4]unsafe.Pointer)(u.InitData)
	for i := 0; i < 2; i++ {
		if mods[i] != nil {
			if p := *(**uint16)(unsafe.Add(mods[i], 8)); p != nil {
				appendText(p)
				out = append(out, ' ')
			}
		}
	}
	if base.Desc8 != nil {
		appendText(base.Desc8)
	}
	for i := 2; i < 4; i++ {
		if mods[i] != nil {
			off := 8
			if i == 3 {
				off = 12
			}
			if p := *(**uint16)(unsafe.Add(mods[i], off)); p != nil {
				out = append(out, ' ')
				appendText(p)
			}
		}
	}
	alloc.StrCopy16S(dst, append(out, 0))
	return &dst[0]
}
