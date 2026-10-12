package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

type serverOptionsMode struct {
	key, title string
	mode       uint16
	hidden     bool
}

var serverOptionsModes = []serverOptionsMode{
	{key: "CTF", mode: 0x20}, {key: "Arena", mode: 0x100}, {key: "Highlander", mode: 0x400}, {key: "KotR", mode: 0x10}, {key: "Flagball", mode: 0x40}, {key: "Quest", mode: 0x1000, hidden: true}, {key: "Noxworld.c:Chat", mode: 0x80},
}
var serverOptionsModesLoaded bool

func serverOptionsLoadModes() {
	if serverOptionsModesLoaded {
		return
	}
	for i := range serverOptionsModes {
		p := &serverOptionsModes[i]
		p.title = serverOptionsText(p.key)
	}
	serverOptionsModesLoaded = true
}
func serverOptionsModeName(mode uint16) string {
	serverOptionsLoadModes()
	mode &= 0x17f0
	for _, p := range serverOptionsModes {
		if p.mode == mode {
			return p.title
		}
	}
	return serverOptionsModes[1].title
}
func serverOptionsModeFromName(title string) int {
	for _, p := range serverOptionsModes {
		if p.title == title {
			return int(p.mode)
		}
	}
	return 0
}
func serverOptionsSelectedMode() int {
	return serverOptionsModeFromName(serverOptionsChild(10119).DrawData().Text())
}
func serverOptionsPopulate() {
	w := serverOptionsChild(10120)
	teamUIEvent(w, 16399, 0, 0)
	for _, p := range serverOptionsModes {
		if !p.hidden {
			serverOptionsSetText(w, 16397, p.title, -1)
		}
	}
}
func serverOptionsMeasure() {
	w := serverOptionsChild(10120)
	r := GetClient().R2()
	max, n := 0, 0
	for _, p := range serverOptionsModes {
		if p.hidden {
			continue
		}
		n++
		if x := r.GetStringSizeWrapped(w.DrawData().Font(), p.title, 0).X; x > max {
			max = x
		}
	}
	height := n*(r.FontHeight(w.DrawData().Font())+1) + 2
	w.EndPos.Y = w.Off.Y + height
	w.SizeVal.Y = height
	w.SizeVal.X = max + 7
	w.Off.X = w.EndPos.X - w.SizeVal.X
}
func serverOptionsTooltip(damage bool, flags byte) int {
	key := "AutoAssign"
	if damage {
		key = "TeamDamage"
	}
	if flags&4 != 0 {
		key += "OnTT"
	} else {
		key += "OffTT"
	}
	Nox_xxx_cursorSetTooltip_4776B0(serverOptionsText(key))
	return 1
}
func serverOptionsTeamCount() int {
	if serverOptionsRoot == 0 {
		return 0
	}
	n := byte(GetServer().S().Teams.Count())
	if teamUISettingsLocked() && noxflags.HasGame(128) || noxflags.HasGame(0x8000) {
		n = byte(teamRuntimeGroupCount())
	}
	p := serverOptionsStoreText(1046364, 64, serverOptionsFormat("NumTeamsMsg", int(n)))
	return teamUIEvent(serverOptionsChild(10110), 16385, p, ^uintptr(0))
}
func serverOptionsName(name string) int {
	if len(name) > 15 {
		name = name[:15]
	}
	return serverOptionsSetText(serverOptionsWindow(1046512), 16414, name, 0)
}

// Keep the shared formatting buffers while unported GUI/rule owners can observe them.
func serverOptionsBufferText(off uintptr) string {
	return alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, off)))
}
