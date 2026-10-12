package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

type combatAlly struct {
	Code             uint32
	Flag             byte
	Reserved         byte
	Current, Maximum uint16
	Padding          uint16
	Active           uint32
}

func combatAllies() []combatAlly {
	return unsafe.Slice((*combatAlly)(memmap.PtrOff(0x5D4594, 1200916)), 32)
}
func combatAllyFree() *combatAlly {
	for i := range combatAllies() {
		a := &combatAllies()[i]
		if a.Active == 0 {
			return a
		}
	}
	return nil
}
func combatAllyLookup(code uint32) *combatAlly {
	for i := range combatAllies() {
		a := &combatAllies()[i]
		if a.Active != 0 && a.Code == code {
			return a
		}
	}
	return nil
}
func combatAllyAdd(code uint32, current, maximum uint16) bool {
	for i := range combatAllies() {
		a := &combatAllies()[i]
		if a.Active == 1 && a.Code == code {
			return true
		}
	}
	a := combatAllyFree()
	if a == nil {
		return false
	}
	a.Code, a.Flag, a.Current, a.Maximum, a.Active = code, 0, current, maximum, 1
	return true
}
func combatAllyRemove(code uint32) bool {
	a := combatAllyLookup(code)
	if a == nil {
		return false
	}
	a.Active = 0
	return true
}
func combatAllyFlag(code uint32, flag byte) bool {
	a := combatAllyLookup(code)
	if a == nil {
		return false
	}
	a.Flag = flag
	return true
}
func combatAllyPair(code uint32, current, maximum uint16) bool {
	a := combatAllyLookup(code)
	if a == nil {
		return false
	}
	a.Current, a.Maximum = current, maximum
	return true
}
func combatAllyFirst(code uint32, current uint16) bool {
	a := combatAllyLookup(code)
	if a == nil {
		return false
	}
	a.Current = current
	return true
}
func combatAllyRead(code uint32, current, maximum *uint16, flag *byte) bool {
	a := combatAllyLookup(code)
	if a == nil {
		return false
	}
	*current, *maximum, *flag = a.Current, a.Maximum, a.Flag
	return true
}
func combatAllyClear() {
	for i := range combatAllies() {
		a := &combatAllies()[i]
		a.Code, a.Flag, a.Current, a.Maximum, a.Active = 0, 0, 0, 0, 0
	}
}

type combatFriend struct {
	Code uint32
	Next *combatFriend
}

var combatFriendPool unsafe.Pointer
var combatFriendHead *combatFriend

func combatFriendClass() alloc.ClassT[combatFriend] {
	return alloc.AsClassT[combatFriend](combatFriendPool)
}
func combatFriendInit() bool {
	c := alloc.NewClassT("FriendListClass", combatFriend{}, 128)
	combatFriendPool = c.UPtr()
	return combatFriendPool != nil
}
func combatFriendClear() { combatFriendClass().FreeAllObjects(); combatFriendHead = nil }
func combatFriendDestroy() {
	combatFriendClass().Free()
	combatFriendPool = nil
	combatFriendHead = nil
}
func combatFriendAdd(code uint32) *combatFriend {
	p := combatFriendClass().NewObject()
	if p != nil {
		p.Code, p.Next = code, combatFriendHead
		combatFriendHead = p
	}
	return p
}
func combatFriendRemove(code uint32) {
	link := &combatFriendHead
	for p := *link; p != nil; p = *link {
		if p.Code == code {
			*link = p.Next
			combatFriendClass().FreeObjectFirst(p)
			return
		}
		link = &p.Next
	}
}
func combatFriendHas(code uint32) bool {
	for p := combatFriendHead; p != nil; p = p.Next {
		if p.Code == code {
			return true
		}
	}
	return false
}
