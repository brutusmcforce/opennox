package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

var monsterScriptCacheCold uint32 = 1

func monsterCacheNode(raw uint32) *[3]uint32 { return (*[3]uint32)(memmap.WordPtr(raw)) }
func monsterCacheList(active bool) *[2]uint32 {
	off := uintptr(2386620)
	if active {
		off = 2386820
	}
	return (*[2]uint32)(memmap.PtrOff(0x5D4594, off))
}
func monsterCachePrepend(list *[2]uint32, raw uint32) uint32 {
	p := monsterCacheNode(raw)
	p[1] = 0
	p[2] = list[0]
	if list[0] != 0 {
		monsterCacheNode(list[0])[1] = raw
	} else {
		list[1] = raw
	}
	list[0] = raw
	return raw
}
func monsterCacheUnlink(list *[2]uint32, raw uint32) uint32 {
	p := monsterCacheNode(raw)
	if p[2] != 0 {
		monsterCacheNode(p[2])[1] = p[1]
	} else {
		list[1] = p[1]
	}
	if p[1] != 0 {
		monsterCacheNode(p[1])[2] = p[2]
		return raw
	}
	list[0] = p[2]
	return p[2]
}
func monsterCacheReset() uint32 {
	*monsterCacheList(false) = [2]uint32{}
	*monsterCacheList(true) = [2]uint32{}
	var result uint32
	for i := uintptr(0); i < 16; i++ {
		result = monsterCachePrepend(monsterCacheList(false), memmap.PtrWord(memmap.PtrOff(0x5D4594, 2386628+12*i)))
	}
	monsterScriptCacheCold = 0
	return result
}
func monsterCacheFind(id int32) *server.Object {
	if monsterScriptCacheCold != 0 {
		monsterCacheReset()
	}
	list := monsterCacheList(true)
	for p := list[0]; p != 0; p = monsterCacheNode(p)[2] {
		u := motionObject(monsterCacheNode(p)[0])
		if u.ObjFlags&0x20 == 0 && uint32(u.ScriptIDVal) == uint32(id) {
			monsterCacheUnlink(list, p)
			monsterCachePrepend(list, p)
			return u
		}
	}
	return nil
}
func monsterCachePrepare(u *server.Object) uint32 {
	free, active := monsterCacheList(false), monsterCacheList(true)
	node := free[0]
	if node != 0 {
		free[0] = monsterCacheNode(node)[2] // C leaves the free predecessor/tail words stale.
	} else {
		node = active[1]
		monsterCacheUnlink(active, node)
	}
	monsterCacheNode(node)[0] = motionAddress(u)
	return monsterCachePrepend(active, node)
}
func monsterCacheRemove(u *server.Object) uint32 {
	if monsterScriptCacheCold != 0 {
		return monsterScriptCacheCold
	}
	list := monsterCacheList(true)
	if list[0] == 0 {
		return 0
	}
	raw := motionAddress(u)
	for p := list[0]; p != 0; p = monsterCacheNode(p)[2] {
		if monsterCacheNode(p)[0] == raw {
			monsterCacheUnlink(list, p)
			return monsterCachePrepend(monsterCacheList(false), p)
		}
	}
	return raw
}
func monsterCacheClear() uint32 {
	result := monsterScriptCacheCold
	if result != 0 {
		return result
	}
	list := monsterCacheList(true)
	for p := list[0]; p != 0; {
		next := monsterCacheNode(p)[2]
		monsterCacheUnlink(list, p)
		result = monsterCachePrepend(monsterCacheList(false), p)
		p = next
	}
	return result
}
