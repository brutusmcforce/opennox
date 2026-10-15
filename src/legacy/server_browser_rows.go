package legacy

import (
	"fmt"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
)

// windowArg returns the window passed in an event argument.
func windowArg(a uintptr) *gui.Window { return (*gui.Window)(unsafe.Pointer(a)) }
func browserNarrowText(s string) string {
	// nox_swprintf's %S widens each unsigned byte, rather than decoding UTF-8.
	out := make([]rune, len(s))
	for i := range out {
		out[i] = rune(s[i])
	}
	return string(out)
}
func browserText(w *gui.Window, s string, color int) {
	p, free := alloc.CString16(s)
	defer free()
	optionsSend(w, 16397, uintptr(unsafe.Pointer(p)), uintptr(color))
}
func browserCappedString(p unsafe.Pointer, limit int) string {
	out := make([]byte, 0, limit)
	for i := 0; i < limit; i++ {
		v := *(*byte)(unsafe.Add(p, i))
		if v == 0 {
			break
		}
		out = append(out, v)
	}
	return string(out)
}
func browserMarkerImages(rec *Nox_gui_server_ent_t) {
	band := 0
	for ; band < 3; band++ {
		if uint32(rec.PingVal) <= memmap.Uint32(0x587000, 87484+4*uintptr(band)) {
			break
		}
	}
	if band > 2 {
		band = 2
	}
	i := 4 * band // large icons, then small ones in the world view
	if browserUI.region == -1 {
		i += 2
	}
	draw := rec.Field_7.DrawData()
	draw.BgImageHnd = browserGameIcons[i]
	draw.HlImageHnd = browserGameIcons[i+1]
	draw.DisImageHnd = draw.BgImageHnd
}
func browserRow(rec *Nox_gui_server_ent_t) {
	name := rec.ServerName()
	if browserUI.listMode != 0 {
		optionsSend(browserUI.gameList, 16397, uintptr(memmap.PtrOff(0x587000, 91164)), 4)
		first := byte(0)
		if len(name) > 0 {
			first = name[0]
		}
		if first == memmap.Uint8(0x5D4594, 815120) {
			name = fmt.Sprintf("%s:%d", rec.Addr(), rec.Port())
		}
		text, free := alloc.CString16(browserNarrowText(name))
		browserTrimName(text, 100)
		optionsSend(browserUI.playersColumn, 16397, uintptr(unsafe.Pointer(text)), 4)
		free()
		browserText(browserUI.modeColumn, fmt.Sprintf("%d/%d", rec.PlayersVal, rec.MaxPlayersVal), 4)
		mode := browserModeName(uint16(rec.Flags()))
		if uint16(rec.Flags())&0x1000 != 0 {
			mode = fmt.Sprintf("%s %d", mode, rec.QuestLevel())
			// The original shared quest label buffer remains observable to its owner.
			alloc.StrCopy16(unsafe.Slice(memmap.PtrUint16(0x5D4594, 814772), 128), mode)
			optionsSend(browserUI.mapColumn, 16397, uintptr(memmap.PtrOff(0x5D4594, 814772)), 4)
		} else {
			browserText(browserUI.mapColumn, mode, 4)
		}
		ping := fmt.Sprint(rec.PingVal)
		if rec.PingVal == 9999 {
			ping = "--"
		}
		browserText(browserUI.pingColumn, ping, 4)
		status := ""
		if rec.StatusVal&0x20 != 0 {
			status = serverPanelsText("noxworld.c", "Noxworld.wnd:private")
		}
		if rec.StatusVal&0x10 != 0 {
			if status != "" {
				status += "+"
			}
			status += serverPanelsText("noxworld.c", "Noxworld.wnd:closed")
		}
		if status == "" {
			key := "Full"
			if rec.PlayersVal < rec.MaxPlayersVal {
				key = "Open"
			}
			status = serverPanelsText("noxworld.c", key)
		}
		browserText(browserUI.statusColumn, status, 4)
		return
	}
	region := int32(browserUI.region)
	parent := browserUI.mapWindow
	size, offset, flags := 20, 10, gui.StatusFlags(1192)
	if region == -1 {
		region = int32(browserRegion(int32(rec.Field_11_0), int32(rec.Field_11_2)))
		parent = browserUI.overview.ChildByID(uint(10054 + region))
		size, offset, flags = 10, 5, 1185
	}
	x := uint16(rec.Field_11_0) - memmap.Uint16(0x587000, 87528+8*uintptr(uint32(region)))
	y := uint16(rec.Field_11_2) - memmap.Uint16(0x587000, 87530+8*uintptr(uint32(region)))
	if browserUI.region == -1 {
		x >>= 1
		y >>= 1
	}
	rec.Field_11_0, rec.Field_11_2 = int16(x), int16(y)
	draw := gui.WindowData{Style: 257, Window: browserUI.world}
	w := NewButtonOrCheckbox(parent, flags, int(int16(x))-offset, int(int16(y))-offset, size, size, &draw)
	rec.Field_7 = w
	browserMarkerImages(rec)
	if name == "" {
		name = fmt.Sprintf("%s:%d", rec.Addr(), rec.Port())
	}
	name = browserNarrowText(name)
	tooltip := fmt.Sprintf("%s %dms", name, rec.PingVal)
	if rec.PingVal == 9999 {
		tooltip = name + " -- ms"
	}
	w.DrawData().SetTooltip(GetServer().S().Strings(), tooltip)
	w.SetFunc94(browserEvent)
	w.SetID(uint(rec.Field_9 + 10070))
}
func browserDetails(record unsafe.Pointer) {
	rec := (*Nox_gui_server_ent_t)(record)
	w := browserUI.detailList
	optionsSend(w, 16399, 0, 0)
	add := func(s string, color int) { browserText(w, s, color) }
	title := func(key string) { add(serverPanelsText("noxworld.c", key), 14) }
	separator := func(off uintptr) { optionsSend(w, 16397, uintptr(memmap.PtrOff(0x587000, off)), ^uintptr(0)) }
	title("Name")
	name := browserCappedString(unsafe.Pointer(&(*Nox_gui_server_ent_t)(record).ServerNameBuf[0]), 255)
	if name == "" {
		name = fmt.Sprintf("%s:%d", rec.Addr(), rec.Port())
	}
	add(browserNarrowText(name), -1)
	separator(89396)
	title("Ping")
	ping := fmt.Sprint(rec.PingVal)
	if rec.PingVal == 9999 {
		ping = "--"
	}
	add(ping, -1)
	separator(89464)
	title("GameType")
	add(browserModeName(uint16(rec.Flags())), -1)
	if uint16(rec.Flags())&0x1000 != 0 {
		separator(89520)
		title("Stage")
		add(fmt.Sprint(rec.QuestLevel()), -1)
	}
	separator(89580)
	title("Map")
	add(browserNarrowText((*Nox_gui_server_ent_t)(record).MapName()), -1)
	separator(89636)
	if uint16(rec.Flags())&0xc000 != 0 {
		key := "Clan"
		if uint16(rec.Flags())&0x4000 != 0 {
			key = "Individual"
		}
		add(serverPanelsText("noxworld.c", key), 6)
		add(serverPanelsText("noxworld.c", "Ladder"), 6)
	}
	separator(89788)
	title("Occupancy")
	add(fmt.Sprintf("%d/%d\n", rec.PlayersVal, rec.MaxPlayersVal), -1)
	if uint16(rec.Flags())&0x2000 == 0 {
		return
	}
	separator(89860)
	title("Resolution")
	add(Get_video_mode_string(int(rec.Field_25_2&0x7f)), -1)
	separator(89916)
	title("DisabledSpells")
	count := 0
	for id := 1; id <= 136; id++ {
		sp := GetServer().S().Spells.DefByInd(spell.ID(id))
		if sp.IsValid() && uint32(GetServer().S().Spells.Flags(spell.ID(id)))&0x7000000 != 0 && rec.Field_33_3[id/8]&(1<<uint(id%8)) == 0 {
			s, _ := Nox_xxx_spellTitle_424930(id)
			add(s, 4)
			count++
		}
	}
	if count == 0 {
		add(serverPanelsText("noxworld.c", "None"), 4)
	}
	separator(90024)
	title("DisabledWeapons")
	count = 0
	for i := uint(0); i < 27; i++ {
		if rec.Field_38_3[i/8]&(1<<(i%8)) == 0 {
			p := unsafe.Pointer(runtimeEquipmentLabel(false, uint32(1)<<i))
			if p != nil {
				add(alloc.GoString16((*uint16)(p)), -1)
				count++
			}
		}
	}
	if count == 0 {
		add(serverPanelsText("noxworld.c", "None"), 4)
	}
	separator(90132)
	title("DisabledArmor")
	count = 0
	for i := uint(0); i < 26; i++ {
		if rec.Field_39_3[i/8]&(1<<(i%8)) == 0 {
			p := unsafe.Pointer(runtimeEquipmentLabel(true, uint32(1)<<i))
			if p != nil {
				add(alloc.GoString16((*uint16)(p)), -1)
				count++
			}
		}
	}
	if count == 0 {
		add(serverPanelsText("noxworld.c", "None"), 4)
	}
}
