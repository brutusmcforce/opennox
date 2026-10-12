package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

type objectXferStream struct{ cf *cryptfile.CryptFile }

func (r objectXferStream) read() bool                  { return r.cf.ReadOnly() }
func (r objectXferStream) raw(p unsafe.Pointer, n int) { r.cf.ReadWrite(unsafe.Slice((*byte)(p), n)) }
func (r objectXferStream) byte(v byte) byte            { r.raw(unsafe.Pointer(&v), 1); return v }
func (r objectXferStream) short(v uint16) uint16 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	r.cf.ReadWrite(b[:])
	return binary.LittleEndian.Uint16(b[:])
}
func (r objectXferStream) word(v uint32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	r.cf.ReadWrite(b[:])
	return binary.LittleEndian.Uint32(b[:])
}
func objectXferEditor() bool { return noxflags.HasGame(0x600000) }
func objectXferScript(p, name unsafe.Pointer) int {
	return scriptBindingCallback(p, name)
}
func (r objectXferStream) name(u *server.Object) bool {
	n := r.byte(byte(len(alloc.GoString((*byte)(u.IDPtr)))))
	if r.read() && n != 0 {
		u.IDPtr = legacyCalloc(1, uintptr(n)+1)
		if u.IDPtr == nil {
			return false
		}
	}
	r.raw(u.IDPtr, int(n))
	if u.IDPtr != nil {
		*(*byte)(unsafe.Add(u.IDPtr, uintptr(n))) = 0
	}
	return true
}
func (r objectXferStream) flags(u *server.Object) {
	old := uint32(u.ObjFlags)
	saved := r.word(old & 0x11408162)
	u.ObjFlags = object.Flags((uint32(u.ObjFlags) & 0xEEBF7E9D) | (old & 0x40) | saved)
	if r.read() {
		if saved&0x1000000 != 0 {
			stateOn(u)
		} else {
			stateOff(u)
		}
	}
}
func (r objectXferStream) position(u *server.Object, integer bool) {
	if r.read() && integer {
		// The stream checksum depends on call boundaries: C reads this pair together.
		var xy [2]int32
		r.raw(unsafe.Pointer(&xy[0]), 8)
		u.PosVec.X = float32(xy[0])
		u.PosVec.Y = float32(xy[1])
	} else {
		r.raw(unsafe.Pointer(&u.PosVec.X), 4)
		r.raw(unsafe.Pointer(&u.PosVec.Y), 4)
	}
	if r.read() {
		u.NewPos = u.PosVec
	}
}
func (r objectXferStream) scriptID(u *server.Object) {
	r.raw(unsafe.Pointer(&u.ScriptIDVal), 4)
	if r.read() && u.ScriptIDVal == 0 && !objectXferEditor() {
		u.ScriptIDVal = int(GetServer().S().Objs.NextObjectScriptID())
	}
}
func (r objectXferStream) inventoryCount(u *server.Object) {
	var count byte
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		count++
	}
	count = r.byte(count)
	if r.read() {
		u.Field34 = uint32(count)
	}
}
func (r objectXferStream) owned(u *server.Object, wide bool) {
	var count uint32
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.ObjFlags)&0x20 == 0 && Sub_4E3B80(int(it.TypeInd)) {
			count++
		}
	}
	if wide {
		count = r.word(uint32(uint16(count)))
	} else {
		count = uint32(r.short(uint16(count)))
	}
	if r.read() {
		for i := 0; i < int(uint16(count)); i++ {
			id := r.word(0)
			if !objectXferEditor() {
				monsterPendingAdd(uint32(u.ScriptIDVal), uint32(id))
			}
		}
	} else {
		for it := u.Field129; it != nil; it = it.Field128 {
			if uint32(it.ObjFlags)&0x20 == 0 && Sub_4E3B80(int(it.TypeInd)) {
				r.raw(unsafe.Pointer(&it.ScriptIDVal), 4)
			}
		}
	}
}
func (r objectXferStream) extra(u *server.Object) {
	v := r.word(u.Field5 & 0x5e)
	u.UnsetXStatus(0x5e)
	u.SetXStatus(v)
}
func objectXferOld(u *server.Object, inner, outer int, r objectXferStream) int {
	if r.read() {
		u.Field34 = 0
	}
	r.raw(unsafe.Pointer(&u.Extent), 4)
	r.flags(u)
	r.position(u, outer < 40 || inner < 4)
	if outer >= 10 && !r.name(u) {
		return 0
	}
	if outer >= 20 {
		r.raw(unsafe.Add(u.CObj(), 52), 1)
	}
	if outer >= 30 {
		r.inventoryCount(u)
	}
	if outer >= 40 {
		r.scriptID(u)
		if inner >= 2 {
			r.owned(u, inner < 5)
		}
		if inner >= 3 {
			r.extra(u)
		}
	}
	return 1
}
func objectXferCommon(u *server.Object, outer int) int {
	r := objectXferStream{cryptfile.Global()}
	originalLifetime := u.Field34
	inner := 0
	if outer >= 40 || !r.read() {
		inner = int(int16(r.short(64)))
		if inner > 64 {
			return 0
		}
	}
	if outer < 40 || inner < 61 {
		return objectXferOld(u, inner, outer, r)
	}
	if r.read() {
		u.Field34 = 0
	}
	r.raw(unsafe.Pointer(&u.Extent), 4)
	r.scriptID(u)
	r.position(u, false)
	present := r.byte(byte(GetServer().S().Sub_4F40A0(u)))
	if present == 0 {
		return 1
	}
	r.flags(u)
	if !r.name(u) {
		return 0
	}
	r.raw(unsafe.Add(u.CObj(), 52), 1)
	r.inventoryCount(u)
	r.owned(u, false)
	r.extra(u)
	if inner >= 63 && objectXferScript(unsafe.Pointer(&u.ScriptPickup), u.Field189) == 0 {
		return 0
	}
	if inner >= 64 {
		life := int32(r.word(originalLifetime - GetServer().S().Frame()))
		if life > 0 && r.read() && uint32(u.ObjFlags)&0x400000 != 0 {
			u.Field32 = uint32(life)
		}
	}
	return 1
}
