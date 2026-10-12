package legacy

import (
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"github.com/brutusmcforce/opennox/server/noxscript"
)

var scriptBindingStringCount uint32
var scriptBindingMoverType uint32

func scriptBindingIntern(text unsafe.Pointer) int {
	i := scriptBindingStringCount
	if int32(i) < 1024 {
		*memmap.PtrPtr(0x973F18, uintptr(uint32(26664)+4*i)) = text
		i++
		scriptBindingStringCount = i
	}
	return int(int32(i) - 1)
}

func scriptBindingCString(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}

func scriptBindingObject(vm noxscript.VM) *server.Object {
	return GetServer().NoxScriptC().ScriptToObject(int(vm.PopI32()))
}

func scriptBindingJournal(vm noxscript.VM, edit bool) int {
	var flags uint16
	if edit {
		flags = uint16(vm.PopU32())
	}
	id := vm.PopU32()
	objectID := vm.PopI32()
	var obj *server.Object
	if objectID != 0 {
		obj = GetServer().NoxScriptC().ScriptToObject(int(objectID))
		if obj == nil {
			return 0
		}
	}
	name := scriptBindingCString(GetServer().S().NoxScriptVM.GetString(id))
	if objectID == 0 {
		if edit {
			journalUpdateAll(name, flags)
		} else {
			journalRemoveAll(name)
		}
	} else if edit {
		journalUnitUpdate(obj, name, flags)
	} else {
		journalUnitRemove(obj, name)
	}
	return 0
}

func scriptBindingRoamByte(vm noxscript.VM, group bool) int {
	value := byte(vm.PopU32())
	if group {
		s := GetServer().S()
		g := s.MapGroups.GroupByInd(int(vm.PopI32()))
		server.EachObjectRecursive(s, g, func(u *server.Object) bool { monsterControlByte(u, &value); return true })
	} else {
		monsterControlByte(scriptBindingObject(vm), &value)
	}
	return 0
}

func scriptBindingGiveXP(vm noxscript.VM) int {
	xp := vm.PopF32()
	if u := scriptBindingObject(vm); u != nil {
		unitGiveExperience(u, xp)
	}
	return 0
}

func scriptBindingHostState(vm noxscript.VM, trading bool) int {
	pl := GetServer().S().Players.ByInd(31)
	result := false
	if pl != nil {
		ud := pl.PlayerUnit.UpdateDataPlayer()
		if trading {
			result = ud.Trade70 != nil
		} else {
			result = ud.DialogWith != nil
		}
	}
	vm.PushBool(result)
	return 0
}

func scriptBindingSubclass(vm noxscript.VM, bit uint32) int {
	u := scriptBindingObject(vm)
	vm.PushBool(u != nil && uint32(u.ObjSubClass)&bit != 0)
	return 0
}

func scriptBindingOwnership(vm noxscript.VM, mode int) int {
	s := GetServer().S()
	// MakeEnemy does not require the local-player slot; the other three original
	// commands require that slot, and only act on its unit when one is present.
	var pl *server.Player
	if mode != 1 {
		pl = s.Players.ByInd(31)
	}
	u := scriptBindingObject(vm)
	if u == nil {
		return 0
	}
	switch mode {
	case 0:
		u.ObjSubClass |= 0x100
		if host := pl.PlayerUnit; host != nil {
			s.ObjSetOwner(host, u)
		}
	case 1:
		u.ObjSubClass &^= 0x100
		s.ObjClearOwner(u)
	case 2:
		if host := pl.PlayerUnit; host != nil {
			statePet(host, u)
		}
	case 3:
		if host := pl.PlayerUnit; host != nil {
			stateRemoveMonitors(host, u)
		}
	}
	return 0
}

func scriptBindingMover(u *server.Object, wp *server.Waypoint) {
	ud := u.UpdateDataMover()
	stateOn(u)
	u.VelVec = types.Pointf{}
	ud.Field_0 = 0
	ud.Field_2 = int32(wp.Index)
	GetServer().S().Objs.AddToUpdatable(u)
}

func scriptBindingMove(u *server.Object, wp *server.Waypoint) {
	if u.ObjFlags&0x8000 != 0 {
		return
	}
	if u.ObjClass&2 != 0 {
		ud := u.UpdateDataMonster()
		u.ClearActionStack()
		if a := u.MonsterPushAction(32); a != nil {
			a.Args[0] = 8
		}
		if wp.PointsCnt != 0 {
			if a := u.MonsterPushAction(10); a != nil {
				a.Args[0] = uintptr(unsafe.Pointer(wp))
				a.Args[2] = uintptr(byte(ud.Field333))
			}
		}
		if a := u.MonsterPushAction(8); a != nil {
			a.SetArgs(wp.PosVec, 0)
		}
	} else if uint32(u.TypeInd) == scriptBindingMoverType {
		scriptBindingMover(u, wp)
	} else {
		for it := GetServer().S().Objs.List; it != nil; it = it.Next() {
			if uint32(it.TypeInd) == scriptBindingMoverType && it.UpdateDataMover().Field_8 == u.Extent {
				scriptBindingMover(it, wp)
			}
		}
	}
}

func scriptBindingRoam(u *server.Object) {
	if u.ObjClass&2 == 0 || u.ObjFlags&0x8000 != 0 {
		return
	}
	ud := u.UpdateDataMonster()
	u.ClearActionStack()
	if a := u.MonsterPushAction(32); a != nil {
		a.Args[0] = 10
	}
	if a := u.MonsterPushAction(10); a != nil {
		a.Args[0] = 0
		a.Args[2] = uintptr(byte(ud.Field333))
	}
}

func scriptBindingHome(u *server.Object) {
	if u.ObjClass&2 == 0 || u.ObjFlags&0x8000 != 0 {
		return
	}
	ud := u.UpdateDataMonster()
	u.ClearActionStack()
	if a := u.MonsterPushAction(32); a != nil {
		a.Args[0] = 37
	}
	if a := u.MonsterPushAction(25); a != nil {
		off := 8 * uintptr(ud.Direction94)
		p := types.Ptf(float32(float64(memmap.Float32(0x587000, 194136+off))*10+float64(u.PosVec.X)), float32(float64(memmap.Float32(0x587000, 194140+off))*10+float64(u.PosVec.Y)))
		a.SetArgs(p)
	}
	if a := u.MonsterPushAction(37); a != nil {
		a.SetArgs(ud.Pos95, 0)
	}
}

func scriptBindingCallback(record, name unsafe.Pointer) int {
	r := objectXferStream{cryptfile.Global()}
	if int16(r.short(1)) > 1 {
		return 0
	}
	words := (*[2]uint32)(record)
	if r.read() {
		n := r.word(0)
		if n >= 1024 {
			return 0
		}
		var text [1024]byte
		r.cf.ReadWrite(text[:n])
		if n != 0 {
			value := scriptBindingCString(string(text[:n]))
			if objectXferEditor() {
				dst := unsafe.Slice((*byte)(name), len(value)+1)
				copy(dst, value)
				dst[len(value)] = 0
			} else {
				words[1] = uint32(GetServer().S().NoxScriptVM.ScriptIndexByName(value))
			}
		}
	} else {
		value := ""
		if objectXferEditor() {
			if name != nil {
				value = alloc.GoString((*byte)(name))
			}
		} else if words[1] != 0xffffffff {
			value = scriptBindingCString(GetServer().S().NoxScriptVM.ScriptNameByIndex(int(int32(words[1]))))
		}
		r.word(uint32(len(value)))
		r.cf.ReadWrite([]byte(value))
	}
	words[0] = r.word(words[0])
	return 1
}
