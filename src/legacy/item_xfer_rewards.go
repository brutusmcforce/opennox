package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func itemXferPositiveStart(u *server.Object, max int) (objectXferStream, int, uint32, bool) {
	r := objectXferStream{cryptfile.Global()}
	saved := u.Field34
	v := int(int16(r.short(uint16(max))))
	if v <= 0 || v > max {
		return r, v, saved, false
	}
	return r, v, saved, objectXferCommon(u, v) != 0
}
func itemXferSpellReward(u *server.Object) int {
	p := &u.UseDataSpellReward().Spell
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	if r.read() {
		if v < 31 {
			r.byte(0)
			a, b := r.byte(0), r.byte(0)
			if a >= 137 {
				a = 0
			}
			if b >= 137 {
				b = 0
			}
			*p = a
			if b != 0 {
				*p = b
			}
			if v == 10 {
				r.byte(0)
			}
		} else {
			readName := func() (byte, bool) {
				n := int(r.byte(0))
				if n >= 128 {
					return 0, false
				}
				var name [128]byte
				r.raw(unsafe.Pointer(&name[0]), n)
				return byte(spell.ParseID(alloc.GoStringS(name[:]))), true
			}
			if v >= 41 {
				id, valid := readName()
				if !valid {
					return 0
				}
				*p = id
			} else {
				if _, valid := readName(); !valid {
					return 0
				}
				a, valid := readName()
				if !valid {
					return 0
				}
				b, valid := readName()
				if !valid {
					return 0
				}
				*p = a
				if b != 0 {
					*p = b
				}
			}
		}
	} else {
		name := spell.ID(*p).String()
		n := int(r.byte(byte(len(name))))
		r.cf.ReadWrite([]byte(name)[:n])
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferAbilityReward(u *server.Object) int {
	p := &u.UseDataAbilityReward().Ability
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	name := server.Ability(*p).String()
	var buf [128]byte
	copy(buf[:], name)
	n := int(r.byte(byte(len(name))))
	if n >= 128 {
		return 0
	}
	r.raw(unsafe.Pointer(&buf[0]), n)
	buf[n] = 0
	// Use the same name resolver as the C callback, including invalid names.
	*p = byte(bookAbilityID(alloc.GoString(&buf[0])))
	return objectXferFinish(r, u, v, saved)
}
func itemXferFieldGuide(u *server.Object) int {
	data := u.UseDataFieldGuide()
	p := unsafe.Pointer(&data.CreatureBuf[0])
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	if r.read() {
		n := int(r.byte(0))
		if n >= 64 {
			return 0
		}
		r.raw(p, n)
		data.CreatureBuf[n] = 0
	} else {
		n := int(r.byte(byte(len(data.Creature()))))
		r.raw(p, n)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferGold(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	r.raw(u.InitData, 4)
	return objectXferFinish(r, u, v, saved)
}
func itemXferToxicCloud(u *server.Object) int {
	r, v, saved, ok := itemXferPositiveStart(u, 61)
	if !ok {
		return 0
	}
	r.raw(unsafe.Pointer(&u.UpdateDataToxicCloud().Life), 4)
	return objectXferFinish(r, u, v, saved)
}
