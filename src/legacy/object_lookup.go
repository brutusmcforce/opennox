package legacy

import (
	"strings"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/server"
)

// The cache is owned by Go. Only the C-owned object pointers cross the two
// remaining C entry points; list nodes never leave this package.
type netCodeCacheNode struct {
	value                  *server.Object
	towardHead, towardTail *netCodeCacheNode
}
type netCodeCacheList struct{ first, last *netCodeCacheNode }
type netCodeCacheStorage struct {
	free  netCodeCacheList
	nodes [16]netCodeCacheNode
	used  netCodeCacheList
}

var netCodeCacheState netCodeCacheStorage
var netCodeCacheNeedInit uint32 = 1

func (l *netCodeCacheList) prepend(n *netCodeCacheNode) {
	n.towardHead = nil
	n.towardTail = l.first
	if l.first != nil {
		l.first.towardHead = n
	} else {
		l.last = n
	}
	l.first = n
}
func (l *netCodeCacheList) remove(n *netCodeCacheNode) {
	if n.towardTail != nil {
		n.towardTail.towardHead = n.towardHead
	} else {
		l.last = n.towardHead
	}
	if n.towardHead != nil {
		n.towardHead.towardTail = n.towardTail
	} else {
		l.first = n.towardTail
	}
}
func netCodeCacheInit() {
	c := &netCodeCacheState
	c.free = netCodeCacheList{}
	c.used = netCodeCacheList{}
	// Initialization intentionally retains values in unused nodes.
	for i := range c.nodes {
		c.free.prepend(&c.nodes[i])
	}
	netCodeCacheNeedInit = 0
}
func netCodeCacheLookup(code uint32) *server.Object {
	if netCodeCacheNeedInit != 0 {
		netCodeCacheInit()
	}
	c := &netCodeCacheState
	for n := c.used.first; n != nil; n = n.towardTail {
		if n.value.NetCode == code {
			c.used.remove(n)
			c.used.prepend(n)
			return n.value
		}
	}
	return nil
}
func netCodeCacheNextUnused() *netCodeCacheNode {
	c := &netCodeCacheState
	n := c.free.first
	if n != nil {
		c.free.first = n.towardTail
	}
	// Match the legacy pop: the new head's backlink and free tail remain stale.
	return n
}
func netCodeCacheAdd(u *server.Object) {
	c := &netCodeCacheState
	n := netCodeCacheNextUnused()
	if n == nil {
		n = c.used.last
		c.used.remove(n)
	}
	n.value = u
	c.used.prepend(n)
}
func netCodeCacheInvalidate(u *server.Object) {
	if netCodeCacheNeedInit != 0 {
		return
	}
	c := &netCodeCacheState
	for n := c.used.first; n != nil; n = n.towardTail {
		if n.value == u {
			c.used.remove(n)
			c.free.prepend(n)
			return
		}
	}
}
func netCodeCacheFlush() {
	if netCodeCacheNeedInit != 0 {
		return
	}
	c := &netCodeCacheState
	for n := c.used.first; n != nil; {
		next := n.towardTail
		c.used.remove(n)
		c.free.prepend(n)
		n = next
	}
}

func objectLookupByNetCode(code uint32) *server.Object {
	if u := netCodeCacheLookup(code); u != nil {
		return u
	}
	objs := &GetServer().S().Objs
	for u := objs.List; u != nil; u = u.ObjNext {
		if !u.ObjFlags.Has(object.FlagDestroyed) && u.NetCode == code {
			netCodeCacheAdd(u)
			return u
		}
		for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
			if !it.ObjFlags.Has(object.FlagDestroyed) && it.NetCode == code {
				netCodeCacheAdd(it)
				return it
			}
		}
	}
	for u := objs.Pending; u != nil; u = u.ObjNext {
		if !u.ObjFlags.Has(object.FlagDestroyed) && u.NetCode == code {
			netCodeCacheAdd(u)
			return u
		}
	}
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if u := pl.PlayerUnit; u != nil && !u.ObjFlags.Has(object.FlagDestroyed) && u.NetCode == code {
			return u
		}
	}
	return nil
}
func objectLookupByScriptID(id uint32) *server.Object {
	objs := &GetServer().S().Objs
	for u := objs.List; u != nil; u = u.ObjNext {
		if !u.ObjFlags.Has(object.FlagDestroyed) && uint32(u.ScriptIDVal) == id {
			return u
		}
		for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
			if !it.ObjFlags.Has(object.FlagDestroyed) && uint32(it.ScriptIDVal) == id {
				return it
			}
		}
	}
	for u := objs.Pending; u != nil; u = u.ObjNext {
		if !u.ObjFlags.Has(object.FlagDestroyed) && uint32(u.ScriptIDVal) == id {
			return u
		}
	}
	for u := objs.MissileList; u != nil; u = u.ObjNext {
		if !u.ObjFlags.Has(object.FlagDestroyed) && uint32(u.ScriptIDVal) == id {
			return u
		}
	}
	return nil
}

// Unqualified names remove only the first prefix. In particular, empty IDs and
// names containing multiple colons differ from Object.FindByID's semantics.
func objectLookupNameMatch(u *server.Object, name string, exact bool) bool {
	if u.IDPtr == nil {
		return false
	}
	id := u.ID()
	if !exact {
		if i := strings.IndexByte(id, ':'); i >= 0 {
			id = id[i+1:]
		}
	}
	return id == name
}
func objectLookupNameAt(u *server.Object, name string, exact bool) *server.Object {
	if objectLookupNameMatch(u, name, exact) {
		return u
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if objectLookupNameMatch(it, name, exact) {
			return it
		}
	}
	return nil
}
func objectLookupByName(name string) *server.Object {
	// The former Go wrapper passed a C string, so embedded NUL ends the query.
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	exact := strings.IndexByte(name, ':') >= 0
	objs := &GetServer().S().Objs
	for u := objs.List; u != nil; u = u.ObjNext {
		if v := objectLookupNameAt(u, name, exact); v != nil {
			return v
		}
	}
	for u := objs.Pending; u != nil; u = u.ObjNext {
		if v := objectLookupNameAt(u, name, exact); v != nil {
			return v
		}
	}
	return nil
}
