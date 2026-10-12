package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func objectXferStart(u *server.Object, max int) (objectXferStream, int, uint32, bool) {
	r := objectXferStream{cryptfile.Global()}
	saved := u.Field34
	version := int(int16(r.short(uint16(max))))
	if version > max {
		return r, version, saved, false
	}
	return r, version, saved, objectXferCommon(u, version) != 0
}
func objectXferFinish(r objectXferStream, u *server.Object, version int, saved uint32) int {
	if r.read() && u.Field34 != 0 && objectXferInventory(uint16(version), u, int32(u.Field34)) == 0 {
		return 0
	}
	u.Field34 = saved
	return 1
}
func objectXferWord(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
func objectXferScriptName(u *server.Object, off int) unsafe.Pointer {
	if u.Field189 == nil {
		return nil
	}
	return unsafe.Add(u.Field189, off)
}

func objectXferPedestal(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	r.raw(u.CollideData, 4)
	return objectXferFinish(r, u, v, saved)
}
func objectXferReadable(u *server.Object) int {
	p := u.UseData.Ptr
	n := uint32(len(alloc.GoString((*byte)(p))) + 1)
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	if v >= 2 {
		n = r.word(n)
	}
	r.raw(p, int(n))
	if r.read() {
		*objectXferWord(p, 256) = 0
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferExit(u *server.Object) int {
	p := u.CollideData
	n := uint32(len(alloc.GoString((*byte)(p))) + 1)
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	if v >= 2 {
		n = r.word(n)
		r.raw(p, int(n))
	} else if r.read() {
		for off := 0; ; off++ {
			at := unsafe.Add(p, off)
			r.raw(at, 1)
			if *(*byte)(at) == 0 {
				break
			}
		}
	} else {
		r.raw(p, int(n))
	}
	if v >= 31 {
		r.raw(unsafe.Add(p, 80), 4)
		r.raw(unsafe.Add(p, 84), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferDoor(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.UpdateData
	direction := r.word(*objectXferWord(p, 12))
	state := r.word(uint32(*(*byte)(unsafe.Add(p, 1))))
	orientation := direction
	if v >= 41 {
		orientation = r.word(*objectXferWord(p, 4))
	}
	if r.read() {
		*objectXferWord(p, 12) = direction
		*(*uint16)(unsafe.Add(p, 40)) = uint16((int32(direction) << 8) / 32)
		*objectXferWord(p, 4) = orientation
		*objectXferWord(p, 8) = direction
		halfX := memmap.Int32(0x587000, 196184+8*uintptr(orientation)) / 2
		halfY := memmap.Int32(0x587000, 196188+8*uintptr(orientation)) / 2
		x := int32(int64((float64(halfX) + float64(u.PosVec.X)) * 0.043478262))
		y := int32(int64((float64(halfY) + float64(u.PosVec.Y)) * 0.043478262))
		worldDoorAttach(u.CObj(), int(x), int(y))
		*objectXferWord(p, 16) = uint32(x)
		*objectXferWord(p, 20) = uint32(y)
		*(*byte)(unsafe.Add(p, 1)) = byte(state)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferLegacyScript(p unsafe.Pointer) int {
	if !cryptfile.Global().ReadOnly() {
		return 0
	}
	prefabScriptScan(mapRoomRaw(unsafe.Pointer(nox_xxx_mapgenGetSomeFile_426A60())), mapRoomRaw(p))
	if noxflags.HasGame(0x400000) {
		return 1
	}
	*(*int32)(unsafe.Add(p, 4)) = -1
	return 0
}
func objectXferTrigger(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	p := u.UpdateData
	if r.read() {
		w, h := int32(r.word(0)), int32(r.word(0))
		u.Shape.Box.W = float32(w)
		u.Shape.Box.H = float32(h)
		if u.Shape.Box.W > 60 {
			u.Shape.Box.W = 60
		}
		if h > 60 {
			u.Shape.Box.H = 60
		}
	} else {
		r.word(uint32(int64(u.Shape.Box.W)))
		r.word(uint32(int64(u.Shape.Box.H)))
	}
	geometryShapeBox((*server.Shape)(unsafe.Pointer(unsafe.Pointer(&u.Shape))))
	if v < 41 {
		var scratch [3]byte
		for i := 0; i < 3; i++ {
			r.raw(unsafe.Pointer(&scratch[0]), 3)
		}
	} else {
		for off := 54; off < 60; off++ {
			r.raw(unsafe.Add(p, off), 1)
		}
	}
	r.raw(p, 4)
	if v >= 3 {
		objectXferScript(unsafe.Add(p, 20), objectXferScriptName(u, 256))
		objectXferScript(unsafe.Add(p, 28), objectXferScriptName(u, 384))
		if v >= 31 {
			objectXferScript(unsafe.Add(p, 12), objectXferScriptName(u, 512))
		}
	} else {
		objectXferLegacyScript(unsafe.Add(p, 20))
		objectXferLegacyScript(unsafe.Add(p, 28))
	}
	if r.read() && v < 31 {
		for i := 0; i < 4; i++ {
			n := r.byte(0)
			r.cf.Seek(int64(4*int(n)), 1)
		}
	}
	r.raw(unsafe.Add(p, 44), 4)
	r.raw(unsafe.Add(p, 48), 4)
	if r.read() {
		*(*byte)(unsafe.Add(p, 52)) = 0
		*(*byte)(unsafe.Add(p, 53)) = 0
	}
	if !r.read() || v >= 21 {
		r.raw(unsafe.Add(p, 52), 1)
		r.raw(unsafe.Add(p, 53), 1)
	}
	if v >= 61 {
		r.raw(unsafe.Add(p, 8), 1)
		r.raw(unsafe.Add(p, 9), 1)
		r.raw(unsafe.Pointer(&u.Field33), 4)
		if r.read() {
			stateAnimation(u, u.Field33)
		}
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferHole(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.CollideData
	if v < 42 {
		*objectXferWord(p, 24) = 0
	} else {
		r.raw(unsafe.Add(p, 24), 4)
	}
	if v < 41 {
		r.raw(unsafe.Add(p, 8), 8)
		*objectXferWord(p, 4) = 0xffffffff
		*objectXferWord(p, 0) = 0
		*objectXferWord(p, 16) = 0
		*(*uint16)(unsafe.Add(p, 20)) = 0
	} else {
		objectXferScript(p, objectXferScriptName(u, 128))
		r.raw(unsafe.Add(p, 8), 8)
		r.raw(unsafe.Add(p, 16), 4)
		r.raw(unsafe.Add(p, 20), 2)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferTransporter(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.UpdateData
	if r.read() {
		r.raw(unsafe.Add(p, 16), 4)
	} else {
		var id uint32
		if *(*unsafe.Pointer)(unsafe.Add(p, 12)) != nil {
			id = *objectXferWord(p, 16)
		}
		r.word(id)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferElevator(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	p := u.UpdateData
	r.raw(unsafe.Add(p, 8), 4)
	if v >= 41 {
		r.raw(unsafe.Add(p, 16), 4)
	}
	if v >= 61 {
		r.raw(unsafe.Add(p, 12), 1)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferShaft(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	r.raw(unsafe.Add(u.UpdateData, 8), 4)
	return objectXferFinish(r, u, v, saved)
}
func objectXferMover(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.UpdateData
	for _, off := range []int{4, 8, 32} {
		r.raw(unsafe.Add(p, off), 4)
	}
	if v >= 41 {
		r.raw(p, 1)
		if r.read() {
			r.raw(unsafe.Add(p, 16), 4)
			r.raw(unsafe.Add(p, 24), 4)
		} else {
			for _, off := range []int{12, 20} {
				var id uint32
				if q := *(*unsafe.Pointer)(unsafe.Add(p, off)); q != nil {
					id = *(*uint32)(q)
				}
				r.word(id)
			}
		}
	}
	if v >= 42 {
		r.raw(unsafe.Pointer(&u.SpeedBase), 4)
		r.raw(unsafe.Pointer(&u.SpeedCur), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferGlyph(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.InitData
	if v < 41 {
		r.word(0)
	}
	r.raw(unsafe.Pointer(&u.Direction1), 1)
	r.raw(unsafe.Add(p, 28), 4)
	r.raw(unsafe.Add(p, 32), 4)
	r.raw(unsafe.Add(p, 20), 1)
	count := int(*(*byte)(unsafe.Add(p, 20)))
	if r.read() && v < 31 {
		r.raw(p, 20)
	} else {
		for i := 0; i < count; i++ {
			field := objectXferWord(p, 4*i)
			if r.read() {
				n := int(r.byte(0))
				buf := make([]byte, n+1)
				r.raw(unsafe.Pointer(&buf[0]), n)
				*field = uint32(spell.ParseID(alloc.GoStringS(buf)))
			} else {
				name := spell.ID(*field).String()
				n := int(r.byte(byte(len(name))))
				buf := []byte(name)
				r.cf.ReadWrite(buf[:n])
			}
		}
	}
	if r.read() {
		u.Direction2 = u.Direction1
		*objectXferWord(p, 24) = 0
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferSentry(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	p := u.UpdateData
	r.raw(unsafe.Add(p, 4), 4)
	r.raw(unsafe.Add(p, 8), 4)
	if r.read() || noxflags.HasGame(0x200000) {
		*objectXferWord(p, 0) = *objectXferWord(p, 4)
	}
	if v >= 61 {
		r.raw(p, 4)
	}
	return objectXferFinish(r, u, v, saved)
}
