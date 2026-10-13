package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var motionDecayHead uint32
var motionSentryHead uint32

func motionObject(raw uint32) *server.Object {
	return (*server.Object)(memmap.WordPtr(raw))
}
func motionAddress(u *server.Object) uint32 { return memmap.PtrWord(unsafe.Pointer(u)) }

func motionDecaySet(u *server.Object, delay int32) uint32 {
	if u.ObjFlags&0x10000 != 0 {
		return uint32(u.ObjFlags)
	}
	if u.ObjFlags&0x400000 != 0 {
		motionDecayRemove(u)
	}
	deadline := GetServer().S().Frame() + uint32(delay)
	u.Field34 = deadline
	var prev *server.Object
	cur := motionObject(motionDecayHead)
	for cur != nil && deadline >= cur.Field34 {
		prev, cur = cur, motionObject(cur.Field117)
	}
	if prev != nil {
		prev.Field117 = motionAddress(u)
	} else {
		motionDecayHead = motionAddress(u)
	}
	u.Field117 = motionAddress(cur)
	u.ObjFlags |= 0x400000
	return uint32(u.ObjFlags)
}
func motionDecayRemove(u *server.Object) uint32 {
	result := uint32(u.ObjFlags)
	if result&0x400000 == 0 {
		return result
	}
	u.ObjFlags &^= 0x400000
	result = motionDecayHead
	var prev *server.Object
	for result != 0 && result != motionAddress(u) {
		prev = motionObject(result)
		result = prev.Field117
	}
	if result != 0 {
		if prev != nil {
			result = u.Field117
			prev.Field117 = result
		} else {
			motionDecayHead = u.Field117
		}
	}
	return result
}
func motionDecayTick() {
	for u := motionObject(motionDecayHead); u != nil; {
		next := motionObject(u.Field117)
		if u.InvHolder != nil {
			motionDecayRemove(u)
		} else {
			if u.Field34 > GetServer().S().Frame() {
				return
			}
			motionDecayRemove(u)
			u.Field5 |= 0x80
			GetServer().DelayedDelete(u)
		}
		u = next
	}
}
func motionDecayClear() uint32 {
	for u := motionObject(motionDecayHead); u != nil; {
		next := motionObject(u.Field117)
		motionDecayRemove(u)
		u = next
	}
	motionDecayHead = 0
	return 0
}
func motionSentryRemove(u *server.Object) uint32 {
	if u.ObjFlags&0x80000000 != 0 {
		if u.Field125 != nil {
			u.Field125.InvNextItem = u.InvNextItem
		} else {
			motionSentryHead = motionAddress(u.InvNextItem)
		}
		if u.InvNextItem != nil {
			u.InvNextItem.Field125 = u.Field125
		}
	}
	u.ObjFlags &^= 0x80000000
	return motionAddress(u)
}
func motionActivate(u *server.Object) int8 {
	if u.ObjClass&1 != 0 {
		return int8(motionAddress(u))
	}
	return collisionActivate(u)
}
func motionDeactivate(u *server.Object) {
	if u.ObjClass&1 == 0 && u.Field116&1 != 0 {
		collisionRemoveActive(u)
	}
}

// MotionSentryForget removes a deleted object from the sentry globe list. A globe normally unlinks itself in
// its next update; one deleted before that would stay in the list and be read after its memory is reused.
func MotionSentryForget(u *server.Object) {
	for p := motionObject(motionSentryHead); p != nil; p = p.InvNextItem {
		if p == u {
			motionSentryRemove(u)
			return
		}
	}
}
