package legacy

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func prefabScriptName(name string, instance, x, y int32, objectOnly bool) string {
	if end := strings.IndexByte(name, 0); end >= 0 {
		name = name[:end]
	}
	var out string
	if objectOnly {
		out = fmt.Sprintf("%s%%%d", name, instance)
	} else {
		out = fmt.Sprintf("%s%%%d%%%d%%%d", name, instance, x, y)
	}
	if len(out) >= 256 {
		out = "ERROR_NAME_TOO_LONG!"
	}
	return out
}
func prefabScriptSetCallback(u *server.Object, event int, name string) {
	if u.Field189 == nil {
		return
	}
	off, field := -1, -1
	data := u.UpdateData
	if event == 14 {
		off = 0
	} else if u.Class().Has(object.ClassTrigger) {
		switch event {
		case 0:
			off, field = 512, 16
		case 1:
			off, field = 256, 24
		case 2:
			off, field = 384, 32
		}
	} else if u.Class().Has(object.ClassMonster) {
		switch event {
		case 3:
			off, field = 640, 1236
		case 4:
			off, field = 768, 1228
		case 5:
			off, field = 896, 1268
		case 6:
			off, field = 1024, 1244
		case 7:
			off, field = 1152, 1252
		case 8:
			off, field = 1280, 1260
		case 9:
			off, field = 1408, 1276
		case 10:
			off, field = 1536, 1284
		case 11:
			off, field = 1664, 1292
		case 13:
			off, field = 1792, 1300
		}
	} else if u.Class().Has(object.ClassHole) {
		if event == 12 {
			off, field = 128, 4
			data = u.CollideData
		}
	} else if u.Class().Has(object.ClassMonsterGenerator) {
		switch event {
		case 15:
			off, field = 1920, 52
		case 16:
			off, field = 2048, 60
		case 17:
			off, field = 2304, 68
		case 18:
			off, field = 2176, 76
		}
	}
	if off < 0 {
		return
	}
	if noxflags.HasGame(noxflags.GameFlag(0x600000)) {
		// Preserve the legacy storage layout and strcpy's used-prefix mutation.
		dst := unsafe.Slice((*byte)(unsafe.Add(u.Field189, off)), len(name)+1)
		copy(dst, name)
		dst[len(name)] = 0
	} else {
		index := GetServer().S().NoxScriptVM.ScriptIndexByName(name)
		if event == 14 {
			*(*int32)(unsafe.Add(u.CObj(), 768)) = int32(index)
		} else {
			*(*int32)(unsafe.Add(data, field)) = int32(index)
		}
	}
}
func prefabScriptObjectNames(instance, x, y int32) {
	s := GetServer().S()
	for p := *prefabGlobal(prefabObjects); p != 0; p = *prefabWord(p, 4) {
		u := (*server.Object)(mapRoomPointer(*prefabWord(p, 0)))
		if uint32(u.ObjFlags)&0x80000000 == 0 {
			continue
		}
		if u.IDPtr != nil {
			name := prefabScriptName(u.ID(), instance, 0, 0, true)
			p := legacyRealloc(u.IDPtr, uintptr(len(name)+1))
			if p != nil {
				u.IDPtr = p
				dst := unsafe.Slice((*byte)(p), len(name)+1)
				copy(dst, name)
				dst[len(name)] = 0
			}
		}
		events := []int{14}
		typ := s.Types.ByInd(int(u.TypeInd))
		if typ != nil {
			switch typ.Xfer {
			case xferIdentityKey(xferIDTrigger):
				events = append(events, 1, 2, 0)
			case xferIdentityKey(xferIDMonster):
				events = append(events, 3, 5, 4, 6, 7, 8, 9, 10, 11)
			case xferIdentityKey(xferIDHole):
				events = append(events, 12)
			case xferIdentityKey(xferIDMonsterGenerator):
				events = append(events, 15, 16, 18, 17)
			}
		}
		for _, event := range events {
			name, ok := s.NoxScriptVM.Nox_script_objCallbackName_508CB0(u, event)
			if ok && name != "" {
				prefabScriptSetCallback(u, event, prefabScriptName(name, instance, x, y, false))
			}
		}
		u.ObjFlags &^= object.Flags(0x80000000)
	}
	for wp := s.WPs.First(); wp != nil; wp = wp.Next() {
		if wp.Flags&0x80000000 == 0 {
			continue
		}
		if name := wp.ID(); name != "" {
			alloc.StrCopyZero(wp.NameBuf[:], prefabScriptName(name, instance, 0, 0, true))
		}
		wp.Flags &^= 0x80000000
	}
}
func prefabScriptPending(bounds *[4]int32, index uint32) uint32 {
	s := GetServer().S()
	for u := s.Objs.Pending; u != nil; u = u.Next() {
		u.ScriptIDVal = int(u.Extent)
		u.Extent = index
		index++
	}
	if s.Objs.Pending == nil {
		return index
	}
	grid := (*[2]uint32)(memmap.PtrOff(0x5D4594, 739980))
	dx, dy := uint32(bounds[0])-23*grid[0], uint32(bounds[1])-23*grid[1]
	word := func(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
	ref := func(p unsafe.Pointer, off int, storePointer bool) {
		target := s.Objs.PendingByScriptID(int(*word(p, off)))
		value, ptr := uint32(0), uint32(0)
		if target != nil {
			value = target.Extent
			ptr = mapRoomRaw(target.CObj())
		}
		*word(p, off) = value
		if storePointer {
			*word(p, off-4) = ptr
		}
	}
	shift := func(p unsafe.Pointer, off int) {
		f := (*[2]float32)(unsafe.Add(p, off))
		f[0] = float32(float64(int32(dx)) + float64(f[0]))
		f[1] = float32(float64(int32(dy)) + float64(f[1]))
	}
	for u := s.Objs.Pending; u != nil; u = u.Next() {
		typ := s.Types.ByInd(int(u.TypeInd))
		if typ == nil {
			continue
		}
		switch typ.Xfer {
		case xferIdentityKey(xferIDElevator), xferIdentityKey(xferIDElevatorShaft):
			ref(u.UpdateData, 8, true)
		case xferIdentityKey(xferIDTransporter):
			ref(u.UpdateData, 16, true)
		case xferIdentityKey(xferIDHole):
			*word(u.CollideData, 8) += dx
			*word(u.CollideData, 12) += dy
		case xferIdentityKey(xferIDExit):
			shift(u.CollideData, 80)
		case xferIdentityKey(xferIDMover):
			wp := s.WPs.PendingByIndTmp(*word(u.UpdateData, 8))
			*word(u.UpdateData, 8) = 0
			if wp != nil {
				*word(u.UpdateData, 8) = wp.Index
			}
			ref(u.UpdateData, 32, false)
		case xferIdentityKey(xferIDGlyph):
			shift(u.InitData, 28)
		}
	}
	return index
}
