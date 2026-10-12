package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func creatureXferText(r objectXferStream, p unsafe.Pointer) {
	n := r.byte(byte(len(alloc.GoString((*byte)(p)))))
	r.raw(p, int(n))
	*(*byte)(unsafe.Add(p, int(n))) = 0
}
func creatureXferSpellName(r objectXferStream, p *uint32) {
	if r.read() {
		var name [256]byte
		n := r.byte(0)
		r.raw(unsafe.Pointer(&name[0]), int(n))
		id := spell.ParseID(alloc.GoString(&name[0]))
		if id < 0 {
			id = 0
		}
		*p = uint32(id)
	} else {
		name := spell.ID(int32(*p)).String()
		n := r.byte(byte(len(name)))
		r.cf.ReadWrite([]byte(name[:int(n)]))
	}
}
func creatureXferSpellTable(r objectXferStream, p unsafe.Pointer, version int) {
	if version < 34 {
		r.raw(unsafe.Add(p, 1488), 548)
		return
	}
	table := unsafe.Slice((*uint32)(unsafe.Add(p, 1488)), 137)
	if r.read() {
		clear(table)
		count := int32(r.word(0))
		for i := int32(0); i < count; i++ {
			var name [256]byte
			n := r.byte(0)
			r.raw(unsafe.Pointer(&name[0]), int(n))
			id := spell.ParseID(alloc.GoString(&name[0]))
			if id < 0 {
				id = 0
			}
			r.raw(unsafe.Add(p, 1488+4*int(id)), 4)
		}
	} else {
		count := uint32(0)
		for _, v := range table {
			if v != 0 {
				count++
			}
		}
		r.word(count)
		for id, v := range table {
			if v == 0 {
				continue
			}
			name := spell.ID(id).String()
			n := r.byte(byte(len(name)))
			r.cf.ReadWrite([]byte(name[:int(n)]))
			r.raw(unsafe.Add(p, 1488+4*id), 4)
		}
	}
}
func creatureXferHeader(r objectXferStream, u *server.Object, v int, npc bool, scriptBase unsafe.Pointer) {
	var direction [2]uint32
	if !r.read() {
		geometryIndexedDirection(int32(int16(u.Direction1)), (*[2]int32)(unsafe.Pointer(unsafe.Pointer(&direction[0]))))
	}
	r.raw(unsafe.Pointer(&direction[0]), 8)
	script := func(off, nameOff int) {
		p := unsafe.Add(u.UpdateData, off)
		if !npc && v < 3 {
			objectXferLegacyScript(p)
			return
		}
		var name unsafe.Pointer
		if scriptBase != nil {
			name = unsafe.Add(scriptBase, nameOff)
		}
		objectXferScript(p, name)
	}
	script(1232, 640)
	script(1264, 896)
	r.raw(unsafe.Add(u.UpdateData, 1220), 2)
	script(1224, 768)
	extra, last := 31, 52
	if npc {
		extra, last = 32, 50
	}
	if v >= extra {
		for _, pair := range [][2]int{{1240, 1024}, {1248, 1152}, {1256, 1280}, {1272, 1408}, {1280, 1536}, {1288, 1664}} {
			script(pair[0], pair[1])
		}
		if v >= last {
			script(1296, 1792)
		}
	}
	if r.read() {
		d := server.Dir16(int32(geometryDirectionAngle((*[2]uint32)(unsafe.Pointer(unsafe.Pointer(&direction[0]))))))
		u.Direction1 = d
		u.Direction2 = d
	}
	if npc || !r.read() || v >= 11 {
		r.word(0)
	}
}

func creatureXferStats(r objectXferStream, u *server.Object, v int, npc bool) {
	p := u.UpdateData
	at := func(off, n int) { r.raw(unsafe.Add(p, off), n) }
	word := func(off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
	start, statusWide, timingWide, actionNamed, more := 31, 51, 46, 34, 33
	if npc {
		start, statusWide, timingWide, actionNamed, more = 32, 49, 47, 35, 34
	}
	if v < start {
		return
	}
	at(1332, 1)
	if v < statusWide {
		*word(1440) = uint32(r.short(0))
	} else {
		at(1440, 4)
	}
	at(1352, 4)
	at(1336, 4)
	at(1344, 4)
	at(1312, 4)
	if npc {
		health := uint16(0)
		if u.HealthData != nil {
			health = u.HealthData.Cur
		}
		health = r.short(health)
		if u.HealthData != nil {
			u.HealthData.Cur = health
		}
	} else if v < 33 {
		r.cf.Seek(2, 1)
	}
	at(1304, 4)
	*word(1308) = *word(1304)
	if v < actionNamed {
		at(1360, 4)
	}
	creatureXferText(r, unsafe.Add(p, 1364))
	creatureXferSpellTable(r, p, v)
	for _, off := range []int{1448, 1456, 1464, 1472, 1480} {
		if v < timingWide {
			*(*uint16)(unsafe.Add(p, off)) = uint16(r.byte(0))
			*(*uint16)(unsafe.Add(p, off+2)) = uint16(r.byte(0))
		} else {
			at(off, 2)
			at(off+2, 2)
		}
		if v < more {
			r.word(0)
		}
	}
	if v >= more-1 {
		at(1316, 4)
	}
	if v >= more {
		at(2040, 4)
		if npc {
			at(1324, 1)
			at(1328, 4)
		}
		at(1320, 4)
		if v < 42 && r.short(0) == 0 {
			*(*byte)(unsafe.Add(p, 1445)) = 1
		}
		for _, off := range []int{2044, 2048, 2052} {
			if npc || v >= 53 {
				creatureXferSpellName(r, word(off))
			} else {
				at(off, 4)
			}
		}
	}
	if v >= actionNamed {
		if r.read() {
			var name [256]byte
			n := r.byte(0)
			r.raw(unsafe.Pointer(&name[0]), int(n))
			*word(1360) = uint32(unitActionIndex(alloc.GoString(&name[0]), true))
		} else {
			name := unitActionName(int32(*word(1360)), true)
			n := r.byte(byte(len(alloc.GoString(name))))
			r.raw(unsafe.Pointer(name), int(n))
		}
	}
}
