package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// C SpawnClass: object +0, next +4, prev +8.
type spawnPolicyNode struct {
	Object     *server.Object
	Next, Prev *spawnPolicyNode
}

// C MonsterListClass: object +0, average +4, class mask +8, 32 distances
// starting +12, then next +140 and prev +144. It is allocated here and used
// by the separate visible-cull half of this batch.
type spawnPolicyMonsterListNode struct {
	Object     *server.Object
	Average    float32
	ClassMask  uint32
	Distance   [32]float32
	Next, Prev *spawnPolicyMonsterListNode
}




func spawnPolicySpawnClass() *alloc.Class {
	return alloc.AsClass(unsafe.Pointer(legacyGlobals.nox_alloc_spawn_2386216))
}
func spawnPolicyMonsterListClass() *alloc.Class {
	return alloc.AsClass(unsafe.Pointer(legacyGlobals.nox_alloc_monsterList_2386220))
}
func spawnPolicyHead() *spawnPolicyNode {
	return (*spawnPolicyNode)(memmap.WordPtr(dword_5d4594_2386212))
}
func spawnPolicySetHead(p *spawnPolicyNode) {
	dword_5d4594_2386212 = memmap.PtrWord(unsafe.Pointer(p))
}

// spawnPolicyInit is 50D780. It intentionally leaves an already-created
// SpawnClass installed if allocation of MonsterListClass fails, matching C.
func spawnPolicyInit() int {
	spawn := alloc.NewClass("SpawnClass", unsafe.Sizeof(spawnPolicyNode{}), 96)
	legacyGlobals.nox_alloc_spawn_2386216 = spawn.UPtr()
	if spawn == nil {
		return 0
	}
	dword_5d4594_2386212 = 0

	list := alloc.NewClass("MonsterListClass", unsafe.Sizeof(spawnPolicyMonsterListNode{}), 96)
	legacyGlobals.nox_alloc_monsterList_2386220 = list.UPtr()
	if list == nil {
		return 0
	}
	dword_5d4594_2386224 = 0
	dword_5d4594_2386228 = 0
	return 1
}

// spawnPolicyReset is 50D7E0: it returns active records to both owned classes
// and clears their heads/count, while retaining the classes themselves.
func spawnPolicyReset() {
	spawnPolicySpawnClass().FreeAllObjects()
	dword_5d4594_2386212 = 0
	spawnPolicyMonsterListClass().FreeAllObjects()
	dword_5d4594_2386224 = 0
	dword_5d4594_2386228 = 0
}

// spawnPolicyFree releases both classes and clears their globals.
func spawnPolicyFree() {
	spawnPolicySpawnClass().Free()
	legacyGlobals.nox_alloc_spawn_2386216 = nil
	dword_5d4594_2386212 = 0
	spawnPolicyMonsterListClass().Free()
	legacyGlobals.nox_alloc_monsterList_2386220 = nil
	dword_5d4594_2386224 = 0
	dword_5d4594_2386228 = 0
}

func spawnPolicyLink(n *spawnPolicyNode) {
	n.Prev = nil
	n.Next = spawnPolicyHead()
	if n.Next != nil {
		n.Next.Prev = n
	}
	spawnPolicySetHead(n)
}

func spawnPolicyUnlink(n *spawnPolicyNode) {
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		spawnPolicySetHead(n.Next)
	}
}

func spawnPolicyUpdateWord(u *server.Object, off uintptr) *uint32 {
	return (*uint32)(unsafe.Add(u.UpdateData, off))
}

// spawnPolicyRegister is 50E030. It retains C inventory insertion: that engine
// owns Inventory links and the historical null-Glyph behavior.
func spawnPolicyRegister(generator, child *server.Object) int {
	childData := child.UpdateData
	if *(*unsafe.Pointer)(unsafe.Add(childData, 2196)) != nil {
		return 1
	}
	n := (*spawnPolicyNode)(spawnPolicySpawnClass().NewObject())
	if n == nil {
		return 0
	}
	spawnPolicyLink(n)
	n.Object = child
	*(*byte)(unsafe.Add(generator.UpdateData, 86))++
	*(*unsafe.Pointer)(unsafe.Add(childData, 2196)) = unsafe.Pointer(n)
	*(*unsafe.Pointer)(unsafe.Add(childData, 2192)) = unsafe.Pointer(generator)

	if child.ObjSubClass&object.SubClass(0x2000) != 0 {
		glyph := GetServer().S().NewObjectByTypeID("Glyph")
		if glyph != nil {
			gd := glyph.InitData
			for i := uintptr(0); i < 3; i++ {
				*(*uint32)(unsafe.Add(gd, 4*i)) = *spawnPolicyUpdateWord(child, 2044+4*i)
			}
			*(*uint32)(unsafe.Add(gd, 24)) = 0
			*(*uint32)(unsafe.Add(gd, 28)) = *(*uint32)(unsafe.Pointer(&child.PosVec.X))
			*(*uint32)(unsafe.Add(gd, 32)) = *(*uint32)(unsafe.Pointer(&child.PosVec.Y))
			var cnt byte
			for i := uintptr(0); i < 3; i++ {
				if *(*uint32)(unsafe.Add(gd, 4*i)) != 0 {
					cnt++
				}
			}
			*(*byte)(unsafe.Add(gd, 20)) = cnt
		}
		inventoryInsert(child, glyph, 1)
	}
	return 1
}

// spawnPolicyRelease is 50E140. It deliberately decrements the owning
// generator's byte before removing/freeing the SpawnClass node.
func spawnPolicyRelease(child *server.Object) {
	if child == nil {
		return
	}
	data := child.UpdateData
	if owner := *(**server.Object)(unsafe.Add(data, 2192)); owner != nil {
		p := (*byte)(unsafe.Add(owner.UpdateData, 86))
		*p = *p - 1
		*(*unsafe.Pointer)(unsafe.Add(data, 2192)) = nil
	}
	if p := *(**spawnPolicyNode)(unsafe.Add(data, 2196)); p != nil {
		spawnPolicyUnlink(p)
		spawnPolicySpawnClass().FreeObjectFirst(unsafe.Pointer(p))
		*(*unsafe.Pointer)(unsafe.Add(data, 2196)) = nil
	}
}

// spawnPolicyDeathRelease is 50E1E0.
func spawnPolicyDeathRelease(u *server.Object) {
	if u != nil && !monsterIsZombie(u) {
		spawnPolicyRelease(u)
	}
}

// spawnPolicyGlyphRelease is 50E210. The cache is blob-backed; the spawn/list
// heads above are standalone C vardefs and must not be accessed through memmap.
func spawnPolicyGlyphRelease(u *server.Object) {
	glyph := memmap.PtrUint32(0x5D4594, 2386360)
	if *glyph == 0 {
		*glyph = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	if u.ObjClass&object.ClassMonster != 0 && u.ObjSubClass&object.SubClass(0x2000) != 0 && *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 2196)) != nil {
		for it := u.InvFirstItem; it != nil; {
			next := it.InvNextItem // C obtains next before delayed delete.
			if uint32(it.TypeInd) == *glyph {
				GetServer().DelayedDelete(it)
			}
			it = next
		}
	}
	spawnPolicyRelease(u)
}

// The retained object and monster-death owners call these two C entrypoints.
//
