package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/server"
	"io"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func mapDrawableByte(d *client.Drawable, off int) *byte    { return (*byte)(unsafe.Add(d.C(), off)) }
func mapDrawableShort(d *client.Drawable, off int) *uint16 { return (*uint16)(unsafe.Add(d.C(), off)) }
func mapDrawableBox(d *client.Drawable) {
	geometryShapeBox((*server.Shape)(unsafe.Pointer(unsafe.Add(d.C(), 44))))
}
func mapDrawableTyped(kind, typ int) int {
	var count uint32
	r := mapDrawableReader{&count}
	outer := int16(r.u16())
	d := mapDrawableBase(typ, outer, r)
	if d == nil {
		GetClient().Cli().Nox_xxx_spriteLoadError_4356E0()
		return 0
	}
	switch kind {
	case 1:
		if uint32(d.Class())&0x80 != 0 {
			direction, lock := r.u32(), r.u32()
			*mapDrawableByte(d, 433) = byte(lock)
			*mapDrawableByte(d, 432) = byte(bool2int(byte(lock) != 0))
			index := direction
			if outer >= 41 {
				index = r.u32()
			}
			*mapDrawableByte(d, 299) = byte(direction)
			dx := memmap.Int32(0x587000, uintptr(uint32(196184)+8*index)) / 2
			dy := memmap.Int32(0x587000, uintptr(uint32(196188)+8*index)) / 2
			// C stores X in a signed local; the inline Y expression divides unsigned.
			x := int32(uint32(d.PosVec.X) + uint32(dx))
			y := uint32(d.PosVec.Y) + uint32(dy)
			worldDrawableAttach(d, int(x/23), int(y/23))
		} else if uint32(d.Class())&0x200 != 0 {
			w, h := float32(int32(r.u32())), float32(int32(r.u32()))
			if w > 60 {
				w = 60
			}
			if h > 60 {
				h = 60
			}
			*effectWord(d, 56) = math.Float32bits(w)
			*effectWord(d, 60) = math.Float32bits(h)
			mapDrawableBox(d)
		}
	case 2:
		for _, f := range [][2]int{{136, 4}, {140, 4}, {144, 4}, {148, 4}, {152, 12}, {164, 2}, {166, 2}, {168, 4}} {
			r.field(d, f[0], f[1])
		}
		if outer >= 2 {
			for _, f := range [][2]int{{176, 2}, {178, 48}, {226, 16}, {242, 16}, {258, 2}, {260, 2}, {262, 2}, {264, 4}, {270, 2}, {272, 2}, {274, 1}} {
				r.field(d, f[0], f[1])
			}
			if outer > 40 {
				if outer >= 42 {
					r.field(d, 172, 4)
				} else {
					*effectWord(d, 172) = uint32(r.u8())
				}
			}
		} else {
			for _, off := range []int{176, 258, 260, 262, 270} {
				*mapDrawableShort(d, off) = 0
			}
			*effectWord(d, 264) = 0
			*mapDrawableByte(d, 274) = 128
			if cryptfile.Global().ReadOnly() && !(math.Float32frombits(*effectWord(d, 140)) <= 63 && float64(int32(*effectWord(d, 148)))*memmap.Float64(0x581450, 9752) <= memmap.Float64(0x581450, 9744)) {
				particleLightIntensity(unsafe.Add(d.C(), 136), 63, false)
			}
		}
		if cryptfile.Global().ReadOnly() {
			for _, off := range []int{432, 433, 434} {
				*mapDrawableByte(d, off) = 0
			}
			for i := 0; i < 16; i++ {
				if *mapDrawableByte(d, 178+3*i) != 0 || *mapDrawableByte(d, 179+3*i) != 0 || *mapDrawableByte(d, 180+3*i) != 0 {
					(*mapDrawableByte(d, 432))++
				}
				if *mapDrawableByte(d, 226+i) != 0 {
					(*mapDrawableByte(d, 433))++
				}
				if *mapDrawableByte(d, 242+i) != 0 {
					(*mapDrawableByte(d, 434))++
				}
			}
			a := float64(*mapDrawableShort(d, 164)) * memmap.Float64(0x581450, 9752) * memmap.Float64(0x581450, 9736)
			*mapDrawableShort(d, 268) = uint16(int64(a))
		}
	case 3:
		// These two readers activate before reading their trailing fields.
		*effectWord(d, 288) = 0
		d.SetActive()
		for i := 0; i < 4; i++ {
			n := r.u8()
			var name [256]byte
			r.bytes(name[:int(n)])
			id := GetServer().S().Modif.Nox_xxx_modifGetIdByName413290(alloc.GoString(&name[0]))
			mod := GetServer().S().Modif.Nox_xxx_modifGetDescById413330(id)
			*effectWord(d, 432+4*i) = uint32(uintptr(unsafe.Pointer(mod)))
			*mapDrawableShort(d, 448) = 65535
			*mapDrawableShort(d, 450) = 65535
		}
		return int(count)
	case 4:
		*effectWord(d, 288) = 0
		d.SetActive()
		*effectWord(d, 56) = math.Float32bits(float32(int32(r.u32())))
		*effectWord(d, 60) = math.Float32bits(float32(int32(r.u32())))
		mapDrawableBox(d)
		copy(unsafe.Slice(mapDrawableByte(d, 432), 6), []byte{90, 90, 90, 10, 10, 10})
		if outer >= 41 {
			for off := 432; off < 438; off++ {
				r.field(d, off, 1)
			}
		}
		return int(count)
	case 5:
		if outer >= 61 {
			r.u32()
			if r.u8() == 1 {
				GetClient().Cli().Objs.MinimapAdd(d, 1)
			}
		}
	}
	*effectWord(d, 288) = 0
	d.SetActive()
	return int(count)
}
func mapDrawableDispatch(typ uint16) int {
	first := memmap.PtrUint32(0x5D4594, 1309792)
	if *first == 0 {
		for i, name := range []string{"ColorLight", "ColorLightMovable", "TeamBase", "PressurePlate"} {
			*memmap.PtrUint32(0x5D4594, uintptr(1309792+4*i)) = uint32(GetClient().Cli().Things.IndByID(name))
		}
	}
	t := uint32(typ)
	switch {
	case t == *first || t == memmap.Uint32(0x5D4594, 1309796):
		return mapDrawableTyped(2, int(typ))
	case t == memmap.Uint32(0x5D4594, 1309800):
		return mapDrawableTyped(3, int(typ))
	case t == memmap.Uint32(0x5D4594, 1309804):
		return mapDrawableTyped(4, int(typ))
	}
	obj := GetClient().Cli().Things.TypeByInd(int(typ))
	if obj == nil {
		return 0
	}
	if uint32(obj.ObjClass)&0x400000 != 0 {
		if uint32(obj.ObjSubClass)&0x18 != 0 {
			return mapDrawableTyped(5, int(typ))
		}
		return mapDrawableTyped(1, int(typ))
	}
	return 0
}
func mapDrawableSection() int {
	f := cryptfile.Global()
	version := []byte{1, 0}
	f.ReadWrite(version)
	if int16(binary.LittleEndian.Uint16(version)) > 1 || !f.ReadOnly() {
		return 0
	}
	var code [2]byte
	f.ReadWrite(code[:])
	for binary.LittleEndian.Uint16(code[:]) != 0 {
		var length [4]byte
		f.ReadMaybeAlign(length[:])
		remaining := int32(binary.LittleEndian.Uint32(length[:]))
		typ := uint16(Nox_xxx_objectTOCgetTT(binary.LittleEndian.Uint16(code[:])))
		if typ == 0 {
			return 0
		}
		obj := GetClient().Cli().Things.TypeByInd(int(typ))
		if obj != nil && uint32(obj.ObjClass)&0x20400000 != 0 {
			remaining -= int32(mapDrawableDispatch(typ))
		}
		if remaining > 0 {
			f.Seek(int64(remaining), io.SeekCurrent)
		}
		f.ReadWrite(code[:])
	}
	return 1
}
