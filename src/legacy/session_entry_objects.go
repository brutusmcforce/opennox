package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var sessionShadowHead *server.Object

func sessionShadowAdd(u *server.Object) *server.Object {
	if u.ObjFlags&0x410000 == 0 {
		u.Field118 = 0
		u.ObjFlags |= 0x10000
		u.Field117 = memmap.PtrWord(unsafe.Pointer(sessionShadowHead))
		if sessionShadowHead != nil {
			sessionShadowHead.Field118 = memmap.PtrWord(u.CObj())
		}
		sessionShadowHead = u
	}
	return u
}
func sessionShadowRemove(u *server.Object) *server.Object {
	if u.ObjFlags&0x10000 != 0 {
		u.ObjFlags &^= 0x10000
		next, prev := (*server.Object)(memmap.WordPtr(u.Field117)), (*server.Object)(memmap.WordPtr(u.Field118))
		if prev != nil {
			prev.Field117 = u.Field117
		} else {
			sessionShadowHead = next
		}
		if next != nil {
			next.Field118 = u.Field118
		}
	}
	return u
}
func sessionClearCrowns() {
	s := GetServer().S()
	cache := memmap.PtrUint32(0x5D4594, 1523076)
	if *cache == 0 {
		*cache = uint32(s.Types.IndByID("Crown"))
	}
	for u := s.Objs.First(); u != nil; {
		next := u.Next()
		if uint32(u.TypeInd) == *cache {
			GetServer().DelayedDelete(u)
			itemRespawnRemove(u)
		}
		u = next
	}
}
func sessionClearTiles() {
	count := int(worldTileGridCapacity)
	grid := unsafe.Slice((*unsafe.Pointer)(unsafe.Pointer(worldTileGrid)), count)
	for y := 0; y < count; y++ {
		for x := 0; x < count; x++ {
			p := unsafe.Add(grid[x], 44*y)
			*(*byte)(p) = 0
			first, second := (*[5]uint32)(unsafe.Add(p, 4)), (*[5]uint32)(unsafe.Add(p, 24))
			first[0] = 255
			second[0] = 255
			mapPaintSubtileClear(first)
			mapPaintSubtileClear(second)
		}
	}
}
