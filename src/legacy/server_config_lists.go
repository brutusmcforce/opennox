package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

type serverConfigAllowed struct {
	list legacyListNode
	name [26]uint16
}
type serverConfigBlocked struct {
	list    legacyListNode
	name    [26]uint16
	expires uint64
	address [24]byte
}

func serverConfigAllowedHead() *legacyListNode { return &serverConfigAllowedList }
func serverConfigBlockedHead() *legacyListNode { return &serverConfigBlockedList }

// ServerConfigAllowedNames returns the names on the allowed players list.
func ServerConfigAllowedNames() []string {
	var out []string
	for p := serverConfigAllowedFirst(); p != nil; p = serverConfigAllowedNext(p) {
		out = append(out, alloc.GoString16(&p.name[0]))
	}
	return out
}

// ServerConfigBlockedEntry is an entry of the blocked players list: a player name, or a serial when Address is not "0".
type ServerConfigBlockedEntry struct {
	Name, Address string
}

// ServerConfigBlocked returns the blocked players list.
func ServerConfigBlocked() []ServerConfigBlockedEntry {
	var out []ServerConfigBlockedEntry
	for p := serverConfigBlockedFirst(); p != nil; p = serverConfigBlockedNext(p) {
		out = append(out, ServerConfigBlockedEntry{Name: alloc.GoString16(&p.name[0]), Address: alloc.GoStringS(p.address[:])})
	}
	return out
}
func serverConfigAllowedFirst() *serverConfigAllowed {
	return (*serverConfigAllowed)(unsafe.Pointer(listNext(serverConfigAllowedHead())))
}
func serverConfigAllowedNext(p *serverConfigAllowed) *serverConfigAllowed {
	return (*serverConfigAllowed)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}
func serverConfigBlockedFirst() *serverConfigBlocked {
	return (*serverConfigBlocked)(unsafe.Pointer(listNext(serverConfigBlockedHead())))
}
func serverConfigBlockedNext(p *serverConfigBlocked) *serverConfigBlocked {
	return (*serverConfigBlocked)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}
func serverConfigAllowedAdd(name *uint16) unsafe.Pointer {
	p := legacyNew[serverConfigAllowed]()
	listInit(&p.list)
	alloc.StrCopyZero16P(p.name[:], name)
	listAppend(serverConfigAllowedHead(), &p.list)
	return unsafe.Pointer(serverPanelsAccessRefresh())
}
func serverConfigBlockedAdd(duration int32, name *uint16, address *byte) unsafe.Pointer {
	p := legacyNew[serverConfigBlocked]()
	listInit(&p.list)
	alloc.StrCopyZero16P(p.name[:], name)
	if address != nil {
		copy(p.address[:], alloc.GoString(address))
	}
	listAppend(serverConfigBlockedHead(), &p.list)
	if duration != 0 {
		p.expires = uint64(uint32(duration)*60000 + uint32(PlatformTicks()))
	}
	return unsafe.Pointer(serverPanelsAccessRefresh())
}
func serverConfigAllowedRemove(index int32) unsafe.Pointer {
	p := serverConfigAllowedFirst()
	for p != nil && index != 0 {
		p = serverConfigAllowedNext(p)
		index--
	}
	if p != nil {
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
	}
	return unsafe.Pointer(p)
}
func serverConfigBlockedRemove(index int32) {
	p := serverConfigBlockedFirst()
	for p != nil && index != 0 {
		p = serverConfigBlockedNext(p)
		index--
	}
	if p != nil {
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
	}
}
func serverConfigExpire() {
	index := int32(0)
	for p := serverConfigBlockedFirst(); p != nil; {
		next := serverConfigBlockedNext(p)
		if p.expires != 0 && uint64(uint32(PlatformTicks())) > p.expires {
			serverConfigBlockedRemove(index)
		} else {
			index++
		}
		p = next
	}
}
func serverConfigAdmissionInit() int32 {
	listClear(serverConfigAllowedHead())
	listClear(serverConfigBlockedHead())
	return serverConfigFileRead(userFilePath(alloc.GoString(memmap.PtrUint8(0x587000, 54280))))
}
func serverConfigAdmissionClose() unsafe.Pointer {
	serverConfigFileWrite(userFilePath("ban.txt"))
	for p := serverConfigAllowedFirst(); p != nil; {
		next := serverConfigAllowedNext(p)
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	result := serverConfigBlockedFirst()
	for p := result; p != nil; {
		next := serverConfigBlockedNext(p)
		listRemove(&p.list)
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	return unsafe.Pointer(result)
}
