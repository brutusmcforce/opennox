package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"math"
	"unsafe"
)

func resourceClientParser(kind string, typ *client.ObjectType, f *binfile.MemFile, input *byte) bool {
	s := alloc.GoString(input)
	switch kind {
	case "direction", "penumbra":
		deg, _, ok := resourceScanInt(s)
		limit := int32(360)
		if kind == "penumbra" {
			limit = 180
		}
		if !ok || deg < 0 || deg >= limit {
			return false
		}
		v := float64(deg)*memmap.Float64(0x581450, 9560)*math.Float64frombits(uint64(qword_581450_9552)) + math.Float64frombits(uint64(qword_581450_9544))
		n := uint16(effectsTruncWord(v))
		if kind == "direction" {
			typ.LightDir = n
			typ.Field_10 = 0
		} else {
			typ.LightPenumbra = n
		}
		return true
	case "update":
		tok := resourceTokenizer{b: unsafe.Slice(input, len(s)+1)}
		name, ok := tok.next(" \t\n\r")
		if !ok {
			return false
		}
		for off := uintptr(175072); ; off += 8 {
			p := (*byte)(*memmap.PtrPtr(0x587000, off))
			if p == nil {
				return false
			}
			if alloc.GoString(p) == name {
				typ.ClientUpdate = *memmap.PtrPtr(0x587000, off+4)
				return true
			}
		}
	case "image":
		id := int32(f.ReadU32())
		var imageType byte
		var name string
		if id == -1 {
			imageType = f.ReadU8()
			n := int(f.ReadU8())
			raw := make([]byte, n)
			_, _ = f.Read(raw)
			name = alloc.GoStringS(raw)
		}
		typ.PrettyImage = uint32(uintptr(GetClient().R2().GetBag().ImageRef(int(id), imageType, name).C()))
		return true
	}
	panic(kind)
}
func resourceClientField(kind string) client.ThingFieldFunc {
	return func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		StrNCopyBytes(buf, str)
		if !resourceClientParser(kind, typ, f, &buf[0]) {
			return fmt.Errorf("failed to parse %q", str)
		}
		return nil
	}
}
func resourceLinkCatalogs() int {
	for off := uintptr(208180); ; off += 20 {
		p := (*byte)(*memmap.PtrPtr(0x587000, off))
		if p == nil {
			break
		}
		name := alloc.GoString(p)
		if len(name) > 0 && name[0] == '#' {
			name = name[1:]
		}
		*memmap.PtrUint32(0x587000, off+4) = uint32(GetServer().S().Types.IndByID(name))
	}
	for _, start := range []uintptr{210712, 210856, 211000, 209344} {
		for off := start; ; off += 24 {
			p := (*byte)(*memmap.PtrPtr(0x587000, off))
			if p == nil {
				break
			}
			mods := &GetServer().S().Modif
			id := mods.Nox_xxx_modifGetIdByName413290(alloc.GoString(p))
			*memmap.PtrPtr(0x587000, off-4) = unsafe.Pointer(mods.Nox_xxx_modifGetDescById413330(id))
		}
	}
	return 0
}
