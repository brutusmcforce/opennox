package legacy

import (
	"image"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

// Journal entries remain C allocated for the shared player/save layout. Direct
// calloc preserves the original nil-on-failure behavior. List logic is Go.
func journalAdd(p *server.Player, name string, flags uint16) *server.PlayerJournal {
	n := legacyNew[server.PlayerJournal]()
	if n == nil {
		return nil
	}
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	copy(n.EntryBuf[:63], name)
	n.Field3 = flags
	n.Next = p.Journal
	if p.Journal != nil {
		p.Journal.Prev = n
	}
	p.Journal = n
	return n
}
func journalFind(p *server.Player, name string) *server.PlayerJournal {
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	for n := p.Journal; n != nil; n = n.Next {
		if alloc.GoString(&n.EntryBuf[0]) == name {
			return n
		}
	}
	return nil
}
func journalUnlink(p *server.Player, n *server.PlayerJournal) {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
	if p.Journal == n {
		p.Journal = n.Next
	}
	legacyFree(unsafe.Pointer(n))
}
func journalRemove(p *server.Player, name string) int {
	n := journalFind(p, name)
	if n == nil {
		return 0
	}
	journalUnlink(p, n)
	return 1
}
func journalUpdate(p *server.Player, name string, flags uint16) *server.PlayerJournal {
	n := journalFind(p, name)
	if n != nil {
		n.Field3 = flags
	}
	return n
}
func journalUnitAdd(u *server.Object, name string, flags uint16) {
	p := u.UpdateDataPlayer().Player
	n := journalAdd(p, name, flags)
	if n == nil {
		return
	}
	if p.PlayerInd == 31 {
		journalMeasure()
	} else {
		gameplayReportJournal(int(p.PlayerInd), alloc.GoStringS(n.EntryBuf[:]), n.Field3, 1)
	}
}
func journalUnitRemove(u *server.Object, name string) {
	p := u.UpdateDataPlayer().Player
	if journalRemove(p, name) == 0 {
		return
	}
	if p.PlayerInd == 31 {
		journalMeasure()
	} else {
		gameplayReportJournal(int(p.PlayerInd), name, 0, 2)
	}
}
func journalRemoveAll(name string) int {
	for u := GetServer().S().Players.FirstUnit(); u != nil; u = GetServer().S().Players.NextUnit(u) {
		journalUnitRemove(u, name)
	}
	return 0
}
func journalUnitUpdate(u *server.Object, name string, flags uint16) uint32 {
	p := u.UpdateDataPlayer().Player
	n := journalUpdate(p, name, flags)
	if n == nil {
		return 0
	}
	if p.PlayerInd != 31 {
		return uint32(gameplayReportJournal(int(p.PlayerInd), alloc.GoStringS(n.EntryBuf[:]), n.Field3, 3))
	}
	return memmap.PtrWord(unsafe.Pointer(n))
}
func journalUpdateAll(name string, flags uint16) int {
	for u := GetServer().S().Players.FirstUnit(); u != nil; u = GetServer().S().Players.NextUnit(u) {
		journalUnitUpdate(u, name, flags)
	}
	return 0
}
func journalRemoveMask(u *server.Object, mask uint16) int {
	p := u.UpdateDataPlayer().Player
	for n := p.Journal; n != nil; {
		next := n.Next
		if n.Field3&mask != 0 {
			journalUnlink(p, n)
		}
		n = next
	}
	return 0
}
func journalText(n *server.PlayerJournal) string {
	sm := GetServer().S().Strings()
	prefix := ""
	id := ""
	switch n.Field3 {
	case 2:
		id = "Journal:QuestLabel"
	case 4:
		id = "Journal:CompletedLabel"
	case 8:
		id = "Journal:HintLabel"
	}
	if id != "" {
		prefix = sm.GetStringInFile(strman.ID(id), "GUIJourn.c")
	}
	return prefix + " " + sm.GetStringInFile(strman.ID("Journal:"+alloc.GoString(&n.EntryBuf[0])), "GUIJourn.c")
}
func journalLocalPlayer() *server.Player {
	return dword_8531A0_2576
}
func journalMeasure() {
	p := journalLocalPlayer()
	if p == nil {
		return
	}
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	gap := r.FontHeight(font)
	height := -gap
	for n := p.Journal; n != nil; n = n.Next {
		height += gap + r.GetStringSizeWrapped(font, journalText(n), 240).Y
	}
	if height < 0 {
		height = 0
	}
	*memmap.PtrUint32(0x5D4594, 1064848) = uint32(height)
}
func journalDraw(x, y, scroll int) {
	p := journalLocalPlayer()
	if p == nil {
		return
	}
	nox_client_drawSetColor_434460(int(nox_color_black_2650656))
	nox_client_drawRectFilledOpaque_49CE30(x, y, 260, 150)
	n := p.Journal
	if n == nil {
		return
	}
	for n.Next != nil {
		n = n.Next
	}
	r := GetClient().R2()
	font := r.GetFonts().AsFont(nil)
	gap := r.FontHeight(font)
	top := y - scroll
	for ; n != nil; n = n.Prev {
		color := uint32(dword_8531A0_2572)
		switch n.Field3 {
		case 1:
			color = uint32(nox_color_white_2523948)
		case 2:
			color = memmap.Uint32(0x85B3FC, 940)
		case 4:
			color = memmap.Uint32(0x85B3FC, 956)
		case 8:
			color = uint32(nox_color_yellow_2589772)
		}
		text := journalText(n)
		height := r.GetStringSizeWrapped(font, text, 240).Y
		bottom := top + height
		if bottom > y {
			nox_xxx_drawSetTextColor_434390(int(color))
			r.DrawStringWrapped(font, text, image.Rect(x+10, top, x+250, top))
		}
		top = bottom + gap
		if top > y+150 {
			break
		}
	}
}
