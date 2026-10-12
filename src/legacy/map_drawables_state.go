package legacy

import (
	"encoding/binary"
	"image"
	"io"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
)

type mapDrawableReader struct{ count *uint32 }

func (r mapDrawableReader) bytes(p []byte) {
	cryptfile.Global().ReadWrite(p)
	*r.count += uint32(len(p))
}
func (r mapDrawableReader) u8() byte { var p [1]byte; r.bytes(p[:]); return p[0] }
func (r mapDrawableReader) u16() uint16 {
	var p [2]byte
	r.bytes(p[:])
	return binary.LittleEndian.Uint16(p[:])
}
func (r mapDrawableReader) u32() uint32 {
	var p [4]byte
	r.bytes(p[:])
	return binary.LittleEndian.Uint32(p[:])
}
func (r mapDrawableReader) skip(n int32) {
	cryptfile.Global().Seek(int64(n), io.SeekCurrent)
	*r.count += uint32(n)
}
func (r mapDrawableReader) field(d *client.Drawable, off, n int) {
	r.bytes(unsafe.Slice((*byte)(unsafe.Add(d.C(), off)), n))
}
func mapDrawableTeam(d *client.Drawable) {
	teamRuntimeJoin(d.TeamVal.ID, d.TeamPtr(), 0, int(*effectWord(d, 128)), 0)
}
func mapDrawableOld(typ int, inner, outer int16, r mapDrawableReader) *client.Drawable {
	id, flags := r.u32(), r.u32()
	x, y := r.u32(), r.u32()
	if outer >= 40 && inner >= 4 {
		x = uint32(floatToInt32(math.Float32frombits(x)))
		y = uint32(floatToInt32(math.Float32frombits(y)))
	}
	if outer >= 10 {
		n := r.u8()
		var name [256]byte
		r.bytes(name[:int(n)])
	}
	var team byte
	if outer >= 20 {
		team = r.u8()
	}
	if outer >= 30 {
		r.u8()
	}
	var extra uint32
	if outer >= 40 {
		r.skip(4)
		if inner >= 2 {
			var n uint16
			if inner < 5 {
				n = uint16(r.u32())
			} else {
				n = r.u16()
			}
			// The legacy layout narrows the multiplied byte count to 16 bits.
			r.skip(int32(uint16(4 * uint32(n))))
		}
		if inner >= 3 {
			extra = r.u32()
		}
	}
	d := GetClient().Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(int(int32(x)), int(int32(y))))
	if d == nil {
		return nil
	}
	*effectWord(d, 120) = *effectWord(d, 120)&0xEEBF7E9D | flags
	*effectWord(d, 280) = *effectWord(d, 280)&0xFFFFFFA1 | extra
	*effectWord(d, 128) = id
	d.TeamVal.ID = server.TeamID(team)
	if !noxflags.HasGame(noxflags.GameHost) && team != 0 && (uint32(d.Class())&0x10000000 == 0 || !noxflags.HasGame(noxflags.GameFlag(128))) {
		mapDrawableTeam(d)
	}
	return d
}
func mapDrawableBase(typ int, outer int16, r mapDrawableReader) *client.Drawable {
	if outer < 40 {
		return mapDrawableOld(typ, 1, outer, r)
	}
	inner := int16(r.u16())
	if inner < 61 {
		return mapDrawableOld(typ, inner, outer, r)
	}
	id := r.u32()
	r.u32()
	x := floatToInt32(math.Float32frombits(r.u32()))
	y := floatToInt32(math.Float32frombits(r.u32()))
	d := GetClient().Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(int(x), int(y)))
	if d == nil {
		return nil
	}
	*effectWord(d, 128) = id
	if r.u8() != 0 {
		flags := r.u32()
		*effectWord(d, 120) = *effectWord(d, 120)&0xEEBF7E9D | flags
		r.skip(int32(r.u8()))
		d.TeamVal.ID = server.TeamID(r.u8())
		r.u8()
		r.skip(4 * int32(r.u16()))
		extra := r.u32()
		*effectWord(d, 280) = *effectWord(d, 280)&0xFFFFFFA1 | extra
		if inner >= 63 {
			r.skip(2)
			n := r.u32()
			r.skip(int32(n))
			r.skip(4)
		}
		if inner >= 64 {
			r.skip(4)
		}
	}
	if !noxflags.HasGame(noxflags.GameHost) {
		p := memmap.PtrUint32(0x5D4594, 1309788)
		if *p == 0 {
			*p = uint32(GetClient().Cli().Things.IndByID("FlagMarker"))
		}
		if d.TeamVal.ID != 0 && d.TypeIDVal != *p {
			mapDrawableTeam(d)
		}
	}
	return d
}
