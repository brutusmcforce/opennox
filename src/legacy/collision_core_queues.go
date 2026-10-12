package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var collisionHitClass unsafe.Pointer
var collisionHitHead uint32
var collisionBuckets [256]uint32
var collisionActiveHead, collisionActiveTail uint32
var collisionAngleHead uint32
var collisionTrigger, collisionPowder, collisionHand uint32
var collisionSmallFist, collisionMediumFist, collisionLargeFist, collisionMeteor, collisionTypesReady uint32
var collisionObjectForce uint32 = 0x42200000
var collisionWallForce uint32 = 0x42c80000

// Hit records stay in the fixed C-backed class. Links and sentinel targets retain
// their 32-bit representation across callbacks and the remaining C callers.
type collisionHitRecord struct {
	BucketNext, Next, A, B uint32
	Normal                 types.Pointf
	Bucket                 int32
}

// Compile-time checks for the fixed allocation record layout.





func collisionHitAt(p uint32) *collisionHitRecord {
	return (*collisionHitRecord)(memmap.WordPtr(p))
}
func collisionObjectAt(p uint32) *server.Object      { return (*server.Object)(memmap.WordPtr(p)) }
func collisionObjectAddress(u *server.Object) uint32 { return memmap.PtrWord(u.CObj()) }
func collisionResetHits() {
	if collisionHitClass == nil {
		collisionHitClass = alloc.NewClass("Hit", 28, 1024).UPtr()
		clear(collisionBuckets[:])
	}
	for p := collisionHitHead; p != 0; p = collisionHitAt(p).Next {
		collisionBuckets[collisionHitAt(p).Bucket] = 0
	}
	alloc.AsClass(collisionHitClass).FreeAllObjects()
	collisionHitHead = 0
}
func collisionAddHit(u *server.Object, target uint32, normal *types.Pointf) {
	code := u.NetCode
	if target > 6 {
		code += collisionObjectAt(target).NetCode
	}
	bucket := int32(code) % 256
	for p := collisionBuckets[bucket]; p != 0; p = collisionHitAt(p).BucketNext {
		h := collisionHitAt(p)
		if h.A == collisionObjectAddress(u) && h.B == target || h.A == target && h.B == collisionObjectAddress(u) {
			return
		}
	}
	p := alloc.AsClass(collisionHitClass).NewObject()
	if p == nil {
		return
	}
	h := (*collisionHitRecord)(p)
	*h = collisionHitRecord{collisionBuckets[bucket], collisionHitHead, collisionObjectAddress(u), target, *normal, bucket}
	collisionBuckets[bucket] = memmap.PtrWord(p)
	collisionHitHead = memmap.PtrWord(p)
}
func collisionDispatch() {
	for p := collisionHitHead; p != 0; p = collisionHitAt(p).Next {
		h := collisionHitAt(p)
		if h.B > 6 || h.B == 0 {
			a := collisionObjectAt(h.A)
			a.CallCollideWith(collisionObjectAt(h.B), &h.Normal)
			if h.B != 0 {
				spellLifeCollide(collisionObjectAt(h.A), collisionObjectAt(h.B))
			}
		}
		if h.B == 6 {
			collisionObjectAt(h.A).CallDamage(nil, nil, 2, 12)
			collisionActivate(collisionObjectAt(h.A))
		} else if h.B != 0 {
			normal := types.Pointf{-h.Normal.X, -h.Normal.Y}
			b := collisionObjectAt(h.B)
			b.CallCollideWith(collisionObjectAt(h.A), &normal)
			spellLifeCollide(collisionObjectAt(h.B), collisionObjectAt(h.A))
			if collisionObjectAt(h.A).ObjFlags&8 != 0 {
				collisionActivate(collisionObjectAt(h.B))
			} else if collisionObjectAt(h.B).ObjFlags&8 != 0 {
				collisionActivate(collisionObjectAt(h.A))
			}
		}
	}
}
func collisionActivate(u *server.Object) int8 {
	result := memmap.PtrWord(u.Update)
	if result == 0 {
		result = memmap.PtrWord(u.Collide)
		if result == 0 || u.ObjFlags&0x40 != 0 {
			return int8(result)
		}
	}
	cls := u.ObjClass
	allow := cls&0x400000 == 0 && u.ObjFlags&8 == 0
	if !allow {
		spike := GetServer().S().Types.IndByID("Spike")
		result = uint32(GetServer().S().Types.IndByID("PeriodicSpike"))
		cls = u.ObjClass
		allow = cls&0xE080 != 0 || u.Collide == collisionKey(collisionIdentityFist) || u.Collide == collisionKey(collisionIdentityUndeadKiller) || uint32(u.TypeInd) == uint32(spike) || uint32(u.TypeInd) == result
	}
	if allow && u.ObjFlags&4 != 0 {
		if cls&0x2008 != 0 {
			GetServer().S().AI.Paths.Sub_50B500()
		}
		result = uint32(uint8(u.Field116))
		if result&1 == 0 {
			if collisionActiveTail != 0 {
				collisionObjectAt(collisionActiveTail).Field115 = collisionObjectAddress(u)
			} else {
				collisionActiveHead = collisionObjectAddress(u)
			}
			collisionActiveTail = collisionObjectAddress(u)
			result = uint32(uint8(u.Field116)) | 1
			u.Field115 = 0
			u.Field116 = u.Field116&^0xff | result
		}
	}
	return int8(result)
}
func collisionRemoveActive(u *server.Object) {
	if u.Field116&1 == 0 {
		return
	}
	var previous uint32
	for p := collisionActiveHead; p != 0; p = collisionObjectAt(p).Field115 {
		if p != collisionObjectAddress(u) {
			previous = p
			continue
		}
		if previous != 0 {
			collisionObjectAt(previous).Field115 = u.Field115
		} else {
			collisionActiveHead = u.Field115
		}
		if collisionActiveTail == p {
			collisionActiveTail = previous
		}
		u.Field115 = 0xffffffff
		u.Field116 &^= 1
		return
	}
}
func collisionPopActive() *server.Object {
	u := collisionObjectAt(collisionActiveHead)
	collisionActiveHead = u.Field115
	if collisionActiveHead == 0 {
		collisionActiveTail = 0
	}
	u.Field115 = 0xffffffff
	u.Field116 &^= 1
	return u
}
func collisionNextActive(u *server.Object) uint32 {
	if u == nil {
		return 0
	}
	return u.Field115
}
