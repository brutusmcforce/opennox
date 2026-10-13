package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

func serverPanelsSpellOpen(parent *gui.Window) int {
	w := Nox_new_window_from_file("spelllst.wnd", serverPanelsSpellProc)
	*serverPanelsWord(1045484) = memmap.PtrWord(w.C())
	if w == nil {
		return 0
	}
	w.SetDraw(serverPanelsBackground)
	serverPanelsOwner(w, parent)
	list, classes := w.ChildByID(1110), w.ChildByID(1112)
	*serverPanelsWord(1045480) = memmap.PtrWord(list.C())
	*serverPanelsWord(1045508) = memmap.PtrWord(classes.C())
	serverPanelsSpellLayout()
	teamUIEvent(list, 16399, 0, 0)
	teamUIEvent(classes, 16399, 0, 0)
	sp := &GetServer().S().Spells
	titleBuf, free := alloc.Make([]uint16{}, 64)
	defer free()
	for i := 1; i < server.SpellsMax; i++ {
		if !sp.DefByInd(spell.ID(i)).IsValid() {
			continue
		}
		flags := uint32(sp.Flags(spell.ID(i)))
		if flags&0x15000 != 0 {
			continue
		}
		class := ""
		if flags&0x1000000 != 0 || flags&0x6000000 == 0x6000000 {
			class = serverPanelsText("spelllst.c", "Common")
		} else {
			if flags&0x6000000 == 0 {
				continue
			}
			if flags&0x2000000 != 0 {
				class += serverPanelsText("spelllst.c", "SpellWizard")
			}
			if flags&0x4000000 != 0 {
				class += serverPanelsText("spelllst.c", "SpellConjurer")
			}
		}
		serverOptionsSetText(classes, 16397, class, -1)
		title, _ := Nox_xxx_spellTitle_424930(i)
		alloc.StrCopy16(titleBuf[:63], title)
		teamUIEvent(list, 16397, uintptr(unsafe.Pointer(&titleBuf[0])), ^uintptr(0))
	}
	for _, pair := range []struct {
		code int
		id   uint
	}{{16408, 1113}, {16409, 1114}} {
		p := serverOptionsPtr(w.ChildByID(pair.id))
		teamUIEvent(list, pair.code, p, 0)
		teamUIEvent(classes, pair.code, p, 0)
	}
	serverPanelsSpellSnapshot(serverPanelsSpellPointer())
	serverPanelsSpellRefresh()
	if !noxflags.HasGame(1) || noxflags.HasGame(49152) {
		serverPanelsEnable(w, 1115, 1133, 0)
	}
	return int(memmap.PtrWord(w.C())) // stored as a pointer word
}
func serverPanelsSpellLayout() *gui.Window {
	list, classes := serverPanelsWindow(1045480), serverPanelsWindow(1045508)
	height := GetClient().R2().FontHeight(list.DrawData().Font()) + 1
	uiWindowResize(list, list.SizeVal.X, 15*height+2)
	uiWindowResize(classes, classes.SizeVal.X, 15*height+2)
	y := list.Off.Y + height + 2
	var last *gui.Window
	for id := uint(1120); id <= 1133; id++ {
		last = serverPanelsWindow(1045484).ChildByID(id)
		last.Off.Y = y
		y += height
	}
	return last
}
func serverPanelsSpellRefresh() int {
	list := uiListData(serverPanelsWindow(1045480))
	index := uiListIndex(list)
	rows := uiListRows(list)
	result := uint32(0)
	for id := uint(1120); id <= 1133; id++ {
		child := serverPanelsWindow(1045484).ChildByID(id)
		title := &rows[index].Text[0]
		index++
		var spellID spell.ID
		if *title != 0 {
			spellID = GetServer().S().Spells.ByTitle(alloc.GoString16(title))
			child.SetHidden(false)
		} else {
			child.SetHidden(true)
		}
		result = child.DrawData().Field0
		if spellID != 0 && serverPanelsSpellQuery(serverPanelsSpellPointer(), int(spellID)) {
			result |= 4
		} else {
			result &^= 4
		}
		child.DrawData().Field0 = result
	}
	return int(int8(byte(result)))
}
func serverPanelsSpellProc(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
	a, b := e.EventArgsC()
	return gui.RawEventResp(serverPanelsSpellEvent(w, e.EventCode(), a, int(b)))
}
func serverPanelsSpellEvent(_ *gui.Window, event int, arg uintptr, _ int) int {
	w := serverPanelsWindow(1045484)
	list, classes := serverPanelsWindow(1045480), serverPanelsWindow(1045508)
	child := (*gui.Window)(unsafe.Pointer(arg))
	scroll := func() {
		teamUIEvent(list, 16384, arg, 0)
		teamUIEvent(classes, 16384, arg, 0)
		serverPanelsSpellRefresh()
	}
	if event == 16384 {
		if child == w.ChildByID(1113) || child == w.ChildByID(1114) {
			scroll()
		}
		return 0
	}
	if event != 16391 {
		return 0
	}
	id := child.ID()
	switch {
	case id == 1113 || id == 1114:
		scroll()
		return 0
	case id == 1115 || id == 1116:
		d := uiListData(list)
		rows := uiListRows(d)
		settings := serverOptionsCurrent()
		for i := 0; i < int(int16(d.Count)); i++ {
			index := GetServer().S().Spells.ByTitle(alloc.GoString16(&rows[i].Text[0]))
			if index == 0 {
				continue
			}
			if id == 1115 {
				if (!noxflags.HasGame(64) && settings[52]&0x40 == 0) || index != 132 {
					serverPanelsSpellBit(serverPanelsSpellPointer(), int(index), 1)
				}
			} else {
				serverPanelsSpellBit(serverPanelsSpellPointer(), int(index), 0)
			}
		}
		if Get_dword_5d4594_2650652() != 0 {
			var rules server.Settings2
			ruleLoad(&rules, "user.rul", nil, 4, 6128)
			out := unsafe.Slice(serverPanelsSpellPointer(), 5)
			for i := range out {
				out[i] &= rules.Field24.Vals[i]
			}
		}
		serverPanelsSpellRefresh()
		serverOptionsDirty(1)
	case id >= 1120 && id <= 1133:
		d := uiListData(list)
		row := uiListIndex(d) + int(id) - 1120
		index := GetServer().S().Spells.ByTitle(alloc.GoString16(&uiListRows(d)[row].Text[0]))
		if index == 0 {
			serverOptionsDirty(1)
			break
		}
		allowed := true
		if Get_dword_5d4594_2650652() != 0 {
			var rules server.Settings2
			ruleLoad(&rules, "user.rul", nil, 4, 6128)
			allowed = serverPanelsSpellQuery(&rules.Field24.Vals[0], int(index))
		}
		reason := ""
		if !allowed {
			reason = "NotInternet"
		} else if index == 132 && (noxflags.HasGame(64) || serverOptionsCurrent()[52]&0x40 != 0) {
			reason = "plyrspel.c:Illegal"
		}
		if reason != "" {
			child.DrawData().Field0 ^= 4
			Nox_xxx_dialogMsgBoxCreate_449A10(w, serverPanelsText("spelllst.c", "Notice"), serverPanelsText("spelllst.c", reason), 33, nil, nil)
			Sub_44A360(1)
		} else {
			serverPanelsSpellBit(serverPanelsSpellPointer(), int(index), bool2int(child.DrawData().Field0&4 == 0))
			serverOptionsDirty(1)
		}
	}
	Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
	return 0
}
