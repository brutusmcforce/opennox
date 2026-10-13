package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
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
func objectXferWord(p unsafe.Pointer, off int) *uint32 { return rawAt[uint32](p, off) }
// objectXferScriptName returns the editor-mode name slot of a script callback (see server.ObjectEditorData).
func objectXferScriptName(u *server.Object, slot int) unsafe.Pointer {
	return u.EditorData().ScriptName(slot)
}

func objectXferPedestal(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	r.raw(unsafe.Pointer(&u.CollideDataAwardSpell().Spell), 4)
	return objectXferFinish(r, u, v, saved)
}
func objectXferReadable(u *server.Object) int {
	data := u.UseDataRead()
	p := unsafe.Pointer(&data.TextBuf[0])
	n := uint32(len(data.Text()) + 1)
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	if v >= 2 {
		n = r.word(n)
	}
	r.raw(p, int(n))
	if r.read() {
		data.LastRead = 0
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferExit(u *server.Object) int {
	// The map name is copied with its actual length, which is not bounded by the field.
	p := unsafe.Pointer(&u.CollideDataExit().MapName[0])
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
	p := u.UpdateDataDoor()
	direction := r.word(p.Dir)
	state := r.word(uint32(p.KeyType))
	orientation := direction
	if v >= 41 {
		orientation = r.word(p.BaseDir)
	}
	if r.read() {
		p.Dir = direction
		p.Angle = int16(uint16((int32(direction) << 8) / 32))
		p.BaseDir = orientation
		p.PrevDir = direction
		halfX := memmap.Int32(0x587000, 196184+8*uintptr(orientation)) / 2
		halfY := memmap.Int32(0x587000, 196188+8*uintptr(orientation)) / 2
		x := int32(int64((float64(halfX) + float64(u.PosVec.X)) * 0.043478262))
		y := int32(int64((float64(halfY) + float64(u.PosVec.Y)) * 0.043478262))
		worldDoorAttach(u.CObj(), int(x), int(y))
		p.GridX = uint32(x)
		p.GridY = uint32(y)
		p.KeyType = byte(state)
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
	p := u.UpdateDataTrigger()
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
		for i := range p.Colors {
			r.raw(unsafe.Pointer(&p.Colors[i]), 1)
		}
	}
	r.raw(unsafe.Pointer(&p.Flags), 4)
	if v >= 3 {
		objectXferScript(unsafe.Pointer(&p.ScriptActivate), objectXferScriptName(u, server.ScriptNameTriggerActivate))
		objectXferScript(unsafe.Pointer(&p.ScriptDeactivate), objectXferScriptName(u, server.ScriptNameTriggerDeactivate))
		if v >= 31 {
			objectXferScript(unsafe.Pointer(&p.ScriptCollide), objectXferScriptName(u, server.ScriptNameTriggerCollide))
		}
	} else {
		objectXferLegacyScript(unsafe.Pointer(&p.ScriptActivate))
		objectXferLegacyScript(unsafe.Pointer(&p.ScriptDeactivate))
	}
	if r.read() && v < 31 {
		for i := 0; i < 4; i++ {
			n := r.byte(0)
			r.cf.Seek(int64(4*int(n)), 1)
		}
	}
	r.raw(unsafe.Pointer(&p.AllowClass), 4)
	r.raw(unsafe.Pointer(&p.IgnoreClass), 4)
	if r.read() {
		p.AllowTeam = 0
		p.IgnoreTeam = 0
	}
	if !r.read() || v >= 21 {
		r.raw(unsafe.Pointer(&p.AllowTeam), 1)
		r.raw(unsafe.Pointer(&p.IgnoreTeam), 1)
	}
	if v >= 61 {
		r.raw(unsafe.Pointer(&p.State), 1)
		r.raw(unsafe.Pointer(&p.Field9), 1)
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
	d := u.CollideDataHole()
	if v < 42 {
		d.Triggered = 0
	} else {
		r.raw(unsafe.Pointer(&d.Triggered), 4)
	}
	if v < 41 {
		r.raw(unsafe.Pointer(&d.Dest), 8)
		d.Script.Func = -1
		d.Script.Flags = 0
		d.Frame = 0
		d.Delay = 0
	} else {
		objectXferScript(unsafe.Pointer(&d.Script), objectXferScriptName(u, server.ScriptNameHole))
		r.raw(unsafe.Pointer(&d.Dest), 8)
		r.raw(unsafe.Pointer(&d.Frame), 4)
		r.raw(unsafe.Pointer(&d.Delay), 2)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferTransporter(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.UpdateDataTransporter()
	if r.read() {
		r.raw(unsafe.Pointer(&p.DestExtent), 4)
	} else {
		var id uint32
		if p.Destination != nil {
			id = p.DestExtent
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
	p := u.UpdateDataElevator()
	r.raw(unsafe.Pointer(&p.PartnerExtent), 4)
	if v >= 41 {
		r.raw(unsafe.Pointer(&p.Field_4), 4)
	}
	if v >= 61 {
		r.raw(unsafe.Pointer(&p.Field_3), 1)
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferShaft(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	r.raw(unsafe.Pointer(&u.UpdateDataElevatorShaft().Field_2), 4)
	return objectXferFinish(r, u, v, saved)
}
func objectXferMover(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	p := u.UpdateDataMover()
	r.raw(unsafe.Pointer(&p.Field_1), 4)
	r.raw(unsafe.Pointer(&p.Field_2), 4)
	r.raw(unsafe.Pointer(&p.Field_8), 4)
	if v >= 41 {
		r.raw(unsafe.Pointer(&p.Field_0), 1)
		if r.read() {
			r.raw(unsafe.Pointer(&p.NextInd), 4)
			r.raw(unsafe.Pointer(&p.PrevInd), 4)
		} else {
			for _, wp := range []*server.Waypoint{p.Next, p.Prev} {
				var id uint32
				if wp != nil {
					id = wp.Index
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
	p := u.InitDataGlyph()
	if v < 41 {
		r.word(0)
	}
	r.raw(unsafe.Pointer(&u.Direction1), 1)
	r.raw(unsafe.Pointer(&p.SpellArg.Pos.X), 4)
	r.raw(unsafe.Pointer(&p.SpellArg.Pos.Y), 4)
	r.raw(unsafe.Pointer(&p.SpellsCnt), 1)
	count := int(byte(p.SpellsCnt))
	if r.read() && v < 31 {
		r.raw(unsafe.Pointer(&p.Spells[0]), 20)
	} else {
		for i := 0; i < count; i++ {
			field := &p.Spells[i]
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
		p.SpellArg.Obj = nil
	}
	return objectXferFinish(r, u, v, saved)
}
func objectXferSentry(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	p := u.UpdateDataSentryGlobe()
	r.raw(unsafe.Pointer(&p.BaseAngle), 4)
	r.raw(unsafe.Pointer(&p.Speed), 4)
	if r.read() || noxflags.HasGame(0x200000) {
		p.Angle = p.BaseAngle
	}
	if v >= 61 {
		r.raw(unsafe.Pointer(&p.Angle), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
