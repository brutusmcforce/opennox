package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// These records remain in the shared C-owned pools: the retained trade engine
// also reads and allocates them. Only the policy and list operations move to Go.
type shopSession = server.TradeSession
type shopItem = server.TradeItem

func shopHead() *shopSession     { return (*shopSession)(memmap.WordPtr(dword_5d4594_2386500)) }
func shopSetHead(s *shopSession) { dword_5d4594_2386500 = memmap.PtrWord(unsafe.Pointer(s)) }
func shopCached() []*shopSession {
	return shopCachedList[:]
}
func shopGold(u *server.Object) *uint32 { return (*uint32)(u.InitData) }
func shopFreeItem(n *shopItem) {
	alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).FreeObjectFirst(unsafe.Pointer(n))
}
func shopInit() int {
	legacyGlobals.nox_alloc_tradeSession_2386492 = alloc.NewClass("TradeSessions", 64, 64).UPtr()
	if legacyGlobals.nox_alloc_tradeSession_2386492 == nil {
		return 0
	}
	legacyGlobals.nox_alloc_tradeItems_2386496 = alloc.NewClass("TradeItems", 16, 500).UPtr()
	if legacyGlobals.nox_alloc_tradeItems_2386496 == nil {
		shopFree()
		return 0
	}
	clear(shopCached())
	shopSetHead(nil)
	return 1
}
func shopFree() int {
	if legacyGlobals.nox_alloc_tradeSession_2386492 != nil {
		alloc.AsClass(legacyGlobals.nox_alloc_tradeSession_2386492).Free()
	}
	legacyGlobals.nox_alloc_tradeSession_2386492 = nil
	if legacyGlobals.nox_alloc_tradeItems_2386496 != nil {
		alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).Free()
	}
	legacyGlobals.nox_alloc_tradeItems_2386496 = nil
	clear(shopCached())
	shopSetHead(nil)
	return 0
}
func shopReset() int {
	for s := shopHead(); s != nil; s = s.Next {
		if s.Kind != 0 {
			for it := s.Stock; it != nil; it = it.Next {
				GetServer().S().Objs.FreeObject(it.Object)
			}
			s.Stock = nil
		}
	}
	// Preserve the original reset's gold lifetime; ordinary destruction below
	// releases gold objects, whereas this bulk reset only releases stock.
	alloc.AsClass(legacyGlobals.nox_alloc_tradeSession_2386492).FreeAllObjects()
	alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).FreeAllObjects()
	clear(shopCached())
	shopSetHead(nil)
	return 0
}
func shopCreate() *shopSession {
	s := (*shopSession)(alloc.AsClass(legacyGlobals.nox_alloc_tradeSession_2386492).NewObject())
	if s == nil {
		return nil
	}
	s.Gold[0] = GetServer().S().NewObjectByTypeID("Gold")
	s.Gold[1] = GetServer().S().NewObjectByTypeID("Gold")
	s.Next = shopHead()
	if s.Next != nil {
		s.Next.Prev = s
	}
	shopSetHead(s)
	return s
}
func shopAdd(s *shopSession, u *server.Object) *shopItem {
	n := (*shopItem)(alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).NewObject())
	if n == nil {
		return nil
	}
	n.Object, n.Value = u, uint32(shopPrice(1, s, u))
	key := shopStockKey(n)
	if s.Stock == nil || key <= shopStockKey(s.Stock) {
		n.Next = s.Stock
		if n.Next != nil {
			n.Next.Prev = n
		}
		s.Stock = n
		return n
	}
	prev := s.Stock
	for prev.Next != nil && key > shopStockKey(prev.Next) {
		prev = prev.Next
	}
	n.Next, n.Prev = prev.Next, prev
	if n.Next != nil {
		n.Next.Prev = n
	}
	prev.Next = n
	return n
}
func shopFreeList(n *shopItem) uint32 {
	for n != nil {
		next := n.Next
		alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).FreeObjectFirst(unsafe.Pointer(n))
		n = next
	}
	return 0
}
func shopFind(n *shopItem, code uint32) *shopItem {
	for ; n != nil; n = n.Next {
		if n.Object.NetCode == code {
			return n
		}
	}
	return nil
}
func shopDestroy(s *shopSession) {
	for n := s.Stock; n != nil; {
		next := n.Next
		GetServer().S().Objs.FreeObject(n.Object)
		alloc.AsClass(legacyGlobals.nox_alloc_tradeItems_2386496).FreeObjectFirst(unsafe.Pointer(n))
		n = next
	}
	GetServer().S().Objs.FreeObject(s.Gold[0])
	GetServer().S().Objs.FreeObject(s.Gold[1])
	shopFreeList(s.Offers[0])
	shopFreeList(s.Offers[1])
	if s.Next != nil {
		s.Next.Prev = s.Prev
	}
	if s.Prev != nil {
		s.Prev.Next = s.Next
	}
	if shopHead() == s {
		shopSetHead(s.Next)
	}
	alloc.AsClass(legacyGlobals.nox_alloc_tradeSession_2386492).FreeObjectFirst(unsafe.Pointer(s))
}
func shopDetach(s *shopSession, u *server.Object) uint32 {
	s.Active = 0
	if u.ObjClass&4 == 0 {
		return memmap.PtrWord(u.CObj())
	}
	data := u.UpdateData
	ptr := (**shopSession)(unsafe.Add(data, 280))
	if *ptr == s {
		*ptr = nil
	}
	return memmap.PtrWord(data)
}
func shopExit(s *shopSession) {
	shopDetach(s, s.Units[0])
	shopDetach(s, s.Units[1])
	player := s.Units[0]
	if player.ObjClass&4 == 0 {
		player = s.Units[1]
	}
	stateUnfreeze(player, 0)
	shopSendShort(player, 713, 1)
	if noxflags.HasGame(noxflags.GameModeQuest) {
		shopCached()[uint8(player.UpdateDataPlayer().Player.PlayerInd)] = s
	} else {
		shopDestroy(s)
	}
}
func shopCancel(s *shopSession) {
	if s.Kind != 0 {
		shopExit(s)
	} else {
		shopCancelTrade(s)
	}
}
func shopPlayerCleanup(ind int) {
	if s := shopCached()[ind]; s != nil {
		shopDestroy(s)
	}
	shopCached()[ind] = nil
}
func shopAddGold(u *server.Object, value uint32) uint32 {
	return resourceAddGold(u, value)
}
func shopSubGold(u *server.Object, value uint32) uint32 {
	return resourceSubGold(u, value)
}
func shopGetGold(u *server.Object) uint32 {
	return resourceGetGold(u)
}
func shopPut(u, item *server.Object) {
	inventoryInsert(u, item, 1)
}
