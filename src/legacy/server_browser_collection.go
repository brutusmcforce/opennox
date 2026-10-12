package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unicode/utf16"
	"unsafe"
)

// Keep one owned snapshot when the selected record is about to be replaced or
// released. C abandoned every detached list; selection needs only this record.
var browserSelectedSnapshot unsafe.Pointer

func browserKeepSelection(p unsafe.Pointer) {
	if browserUI.selected != p {
		return
	}
	if browserSelectedSnapshot == nil {
		browserSelectedSnapshot, _ = alloc.Malloc(169)
	}
	copy(unsafe.Slice((*byte)(browserSelectedSnapshot), 169), unsafe.Slice((*byte)(p), 169))
	browserUI.selected = browserSelectedSnapshot
}
func browserListHead() *legacyListNode {
	if browserUI.servers == nil {
		browserUI.servers, _ = alloc.New(legacyListNode{})
	}
	return browserUI.servers
}
func browserListClear(windows bool) unsafe.Pointer {
	head := browserListHead()
	if memmap.Uint32(0x5D4594, 1305808) == 0 {
		listClear(head)
	}
	first := listNext(head)
	for p := first; p != nil; {
		next := listNext(p)
		browserKeepSelection(unsafe.Pointer(p))
		listRemove(p)
		if windows {
			w := (*gui.Window)(unsafe.Pointer(uintptr((*Nox_gui_server_ent_t)(unsafe.Pointer(p)).Field_7)))
			w.Destroy()
		}
		alloc.Free(p)
		p = next
	}
	*memmap.PtrUint32(0x5D4594, 1305808) = 1
	return unsafe.Pointer(first)
}
func browserASCIICompare(a, b string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}
func browserWideCompare(a, b string) int {
	x, y := utf16.Encode([]rune(a)), utf16.Encode([]rune(b))
	n := len(x)
	if len(y) < n {
		n = len(y)
	}
	for i := 0; i < n; i++ {
		if x[i] < y[i] {
			return -1
		}
		if x[i] > y[i] {
			return 1
		}
	}
	if len(x) < len(y) {
		return -1
	}
	if len(x) > len(y) {
		return 1
	}
	return 0
}
func browserListInsert(p *legacyListNode) int {
	rec := (*Nox_gui_server_ent_t)(unsafe.Pointer(p))
	head := browserListHead()
	mode := uint32(browserUI.sort)
	switch mode {
	case 2:
		rec.Sort_key = int32(rec.PlayersVal)
	case 3:
		rec.Sort_key = 32 - int32(rec.PlayersVal)
	case 6:
		rec.Sort_key = rec.PingVal
	case 7:
		rec.Sort_key = 1000 - rec.PingVal
	case 8:
		rec.Sort_key = int32(rec.StatusVal & 0x30)
	case 9:
		rec.Sort_key = 48 - int32(rec.StatusVal&0x30)
	}
	if mode == 2 || mode == 3 || mode >= 6 && mode <= 9 {
		return listAscending(head, p)
	}
	index := 0
	for it := listNext(head); it != nil; it = listNext(it) {
		other := (*Nox_gui_server_ent_t)(unsafe.Pointer(it))
		cmp := 0
		if mode < 2 {
			// Original comparison follows the record's C string through its terminator.
			cmp = browserASCIICompare(alloc.GoString((*byte)(unsafe.Add(unsafe.Pointer(p), 120))), alloc.GoString((*byte)(unsafe.Add(unsafe.Pointer(it), 120))))
		} else {
			cmp = browserWideCompare(browserModeName(uint16(rec.Flags())), browserModeName(uint16(other.Flags())))
		}
		if mode == 1 || mode == 5 {
			cmp = -cmp
		}
		if cmp <= 0 {
			listAppend(it, p)
			return index
		}
		index++
	}
	listAppend(head, p)
	return index
}
func browserListAdd(record unsafe.Pointer) int {
	if browserUI.sort > 9 {
		return 0
	} // No valid UI path uses other modes.
	p, _ := alloc.Malloc(169)
	copy(unsafe.Slice((*byte)(p), 169), unsafe.Slice((*byte)(record), 169))
	return browserListInsert((*legacyListNode)(p))
}
func browserListRender() {
	for p := listNext(browserListHead()); p != nil; p = listNext(p) {
		browserRow((*Nox_gui_server_ent_t)(unsafe.Pointer(p)))
	}
}
func browserListResort() {
	head := browserListHead()
	var nodes []*legacyListNode
	for p := listNext(head); p != nil; p = listNext(p) {
		browserKeepSelection(unsafe.Pointer(p))
		nodes = append(nodes, p)
	}
	listClear(head)
	for _, p := range nodes {
		browserListInsert(p)
	}
	browserListRender()
}
func browserListAt(index int32) unsafe.Pointer {
	if index < 0 {
		return nil
	}
	for p := listNext(browserListHead()); p != nil; p = listNext(p) {
		if index == 0 {
			return unsafe.Pointer(p)
		}
		index--
	}
	return nil
}
func browserListID(id int32) unsafe.Pointer {
	for p := listNext(browserListHead()); p != nil; p = listNext(p) {
		if (*Nox_gui_server_ent_t)(unsafe.Pointer(p)).Field_9 == uint32(id) {
			return unsafe.Pointer(p)
		}
	}
	return nil
}
func browserListMissing(addr string, port int16) int {
	for p := listNext(browserListHead()); p != nil; p = listNext(p) {
		rec := (*Nox_gui_server_ent_t)(unsafe.Pointer(p))
		if rec.Addr() == addr && int(port) == rec.Port() {
			return 0
		}
	}
	return 1
}
