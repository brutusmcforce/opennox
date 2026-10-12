package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

var monsterPendingClass *alloc.Class
var monsterPendingHead uint32

type monsterPendingNode struct{ Owner, Unit, Next uint32 }

func monsterPendingInit() int {
	monsterPendingHead = 0
	monsterPendingClass = alloc.NewClass("PendingOwn", 12, 512)
	return bool2int(monsterPendingClass != nil)
}
func monsterPendingFree() uint32 {
	monsterPendingClass.Free()
	monsterPendingClass = nil
	monsterPendingHead = 0
	return 0
}
func monsterPendingClear() { monsterPendingClass.FreeAllObjects(); monsterPendingHead = 0 }
func monsterPendingAdd(owner, unit uint32) uint32 {
	p := (*monsterPendingNode)(monsterPendingClass.NewObject())
	if p == nil {
		return 0
	}
	*p = monsterPendingNode{owner, unit, monsterPendingHead}
	monsterPendingHead = uint32(uintptr(unsafe.Pointer(p)))
	return monsterPendingHead
}
func monsterPendingResolve() {
	for raw := monsterPendingHead; raw != 0; {
		p := (*monsterPendingNode)(unsafe.Pointer(uintptr(raw)))
		owner := objectLookupByScriptID(p.Owner)
		unit := objectLookupByScriptID(p.Unit)
		if owner != nil && unit != nil {
			GetServer().S().ObjSetOwner(owner, unit)
		}
		raw = p.Next
	}
	monsterPendingClear()
}
