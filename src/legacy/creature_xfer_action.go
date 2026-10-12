package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func creatureXferNetReference(r objectXferStream, p *uint32) {
	id := uint32(0)
	if u := objectLookupByNetCode(*p); u != nil {
		id = uint32(u.ScriptIDVal)
	}
	id = r.word(id)
	if r.read() {
		*p = id
	}
}
func creatureXferAction(u *server.Object) int {
	r := objectXferStream{cryptfile.Global()}
	p := u.UpdateData
	version := int16(r.short(4))
	if version > 4 {
		return 0
	}
	if version >= 2 {
		present := byte(0)
		if noxflags.HasGame(1) && !noxflags.HasGame(0x400000) {
			present = 1
		}
		if r.byte(present) == 0 {
			return 1
		}
	}
	frame := GetServer().S().Frame()
	savedFrame := r.word(frame)
	delta := frame - savedFrame
	residue := r.word(0)
	at := func(off, n int) { r.raw(unsafe.Add(p, off), n) }
	word := func(off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
	byteAt := func(off int) *byte { return (*byte)(unsafe.Add(p, off)) }
	adjust := func(off int) { creatureXferAdjust(word(off), delta) }
	at(8, 4)
	at(12, 8*int(*word(8)))
	at(268, 4)
	at(272, 8)
	at(280, 4)
	at(284, 1)
	if r.read() {
		adjust(280)
	}
	at(296, 4)
	residue = 0
	for i := int32(0); i < int32(*word(296)); i++ {
		dest := (**server.Waypoint)(unsafe.Add(p, 300+4*i))
		if r.read() {
			*dest = GetServer().S().WPs.PendingByInd(int(r.word(0)))
		} else {
			r.raw(unsafe.Pointer(*dest), 4)
		}
		residue++
	}
	at(364, 4)
	at(368, 8)
	at(376, 4)
	at(380, 8)
	var name [256]byte
	initial := sound.ID(*word(388)).String()
	copy(name[:], initial)
	n := r.byte(byte(len(initial)))
	r.raw(unsafe.Pointer(&name[0]), int(n))
	name[n] = 0
	*word(388) = uint32(sound.ByName(alloc.GoString(&name[0])))
	at(396, 8)
	at(404, 4)
	if r.read() {
		adjust(404)
	}
	at(481, 1)
	at(482, 1)
	at(483, 1)
	if version < 3 {
		residue = r.word(residue)
	}
	at(496, 4)
	at(500, 8)
	if r.read() {
		adjust(496)
	}
	r.word(residue)
	at(536, 4)
	at(540, 4)
	adjust(536)
	adjust(540)
	at(544, 1)
	for i := 0; i <= int(int8(*byteAt(544))); i++ {
		creatureXferArgument(unsafe.Add(p, 552+24*i), delta)
	}
	at(1129, 1)
	for i := 0; i < int(*byteAt(1129)); i++ {
		dest := unsafe.Add(p, 1132+4*i)
		if r.read() {
			r.raw(dest, 4)
		} else {
			cache := memmap.PtrUint32(0x5d4594, 2487688)
			if *cache == 0 {
				*cache = uint32(GetServer().S().Types.IndByID("NewPlayer"))
			}
			ref := *(**server.Object)(dest)
			r.raw(unsafe.Pointer(&ref.ScriptIDVal), 4)
		}
	}
	if r.read() {
		at(1196, 4)
	} else {
		id := uint32(0)
		if ref := *(**server.Object)(unsafe.Add(p, 1196)); ref != nil {
			id = uint32(ref.ScriptIDVal)
		}
		r.word(id)
	}
	at(1204, 4)
	adjust(1204)
	at(2096, 4)
	at(2100, 4)
	at(2104, 1)
	at(2105, 1)
	text := (*byte)(unsafe.Add(p, 2106))
	n = r.byte(byte(len(alloc.GoString(text))))
	r.raw(unsafe.Pointer(text), int(n))
	*byteAt(2106 + int(n)) = 0
	if version < 4 {
		return 1
	}
	at(4, 4)
	at(288, 4)
	at(292, 4)
	creatureXferNetReference(r, word(392))
	at(492, 4)
	for _, off := range []int{508, 512, 516, 520, 528, 532} {
		at(off, 4)
		adjust(off)
	}
	at(524, 4)
	at(548, 4)
	adjust(548)
	at(1128, 1)
	creatureXferNetReference(r, word(1200))
	at(1208, 4)
	adjust(1208)
	at(1212, 4)
	adjust(1212)
	id := uint32(0)
	if ref := *(**server.Object)(unsafe.Add(p, 1216)); ref != nil {
		id = uint32(ref.ScriptIDVal)
	}
	id = r.word(id)
	if r.read() {
		*word(1216) = id
	}
	at(2172, 1)
	for i := 0; i < int(*byteAt(2172)); i++ {
		creatureXferNetReference(r, word(2140+4*i))
	}
	return 1
}
