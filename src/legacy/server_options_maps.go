package legacy

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/server"
)

func serverOptionsMapList(mode int, current string, update bool) {
	w := serverOptionsWindow(1046496)
	teamUIEvent(w, 16399, 0, 0)
	*serverOptionsWord(1046552) = uint32(mode)
	// The list from the window file holds 256 maps and drops the first ones for more (opennox#651).
	n := 0
	for it := mapCatalogFirst(); it != nil; it = mapCatalogNext(it) {
		n++
	}
	uiListGrow(w, n)
	selected, index := -1, 0
	for it := mapCatalogFirst(); it != nil; it = mapCatalogNext(it) {
		if it.Field_6 == 0 || int(sessionMapFlags(it))&mode == 0 {
			continue
		}
		var base, custom server.Settings2
		name := alloc.GoString(&it.Name[0])
		alloc.StrCopyZero(base.Field0[:], name)
		custom = base
		ruleLoad(&base, "", nil, 1, uint16(mode))
		ruleLoad(&custom, "user.rul", nil, 3, uint16(mode))
		color := -1
		if !bytes.Equal(unsafe.Slice((*byte)(unsafe.Pointer(&base.Field24)), 28), unsafe.Slice((*byte)(unsafe.Pointer(&custom.Field24)), 28)) {
			color = 6
		}
		serverOptionsSetText(w, 16397, serverOptionsFormat("RecPlayers", name, int(it.Field_8_0), int(it.Field_8_1)), color)
		if mapASCIIEqual(current, name) {
			selected = index
			teamUIEvent(w, 16403, uintptr(index), 0)
			teamUIEvent(w, 16412, uintptr(index), 0)
		}
		index++
	}
	if selected < 0 {
		teamUIEvent(w, 16403, 0, 0)
		teamUIEvent(w, 16412, 0, 0)
	}
	if !update {
		return
	}
	data := serverOptionsCurrent()
	selected = teamUIEvent(w, 16404, 0, 0)
	if selected < 0 {
		data[0] = 0
	} else {
		alloc.StrCopyZero(data[:9], serverOptionsMapToken(serverOptionsGetText(w, 16406, selected)))
	}
	ruleLoad((*server.Settings2)(unsafe.Pointer(&data[0])), "user.rul", nil, 7, uint16(mode))
	serverOptionsSettingsLabels(data)
}
func serverOptionsResetMap() int {
	data := serverOptionsCurrent()
	w := serverOptionsWindow(1046496)
	selected := teamUIEvent(w, 16404, 0, 0)
	text := serverOptionsGetText(w, 16406, selected)
	units := alloc.InternCString16(text)
	buf := make([]uint16, 256)
	alloc.StrCopy16P(buf[:255], units)
	teamUIEvent(w, 16398, uintptr(selected), 0)
	teamUIEvent(w, 16402, uintptr(selected), 0)
	serverOptionsSetText(w, 16397, alloc.GoString16(&buf[0]), -1)
	teamUIEvent(w, 16403, uintptr(selected), 0)
	ifs.Remove("maps\\" + alloc.GoString(&data[0]) + "\\user.rul")
	ruleLoad((*server.Settings2)(unsafe.Pointer(&data[0])), "user.rul", nil, 5, binary.LittleEndian.Uint16(data[52:]))
	serverOptionsSettingsLabels(data)
	return serverOptionsDirty(1)
}
