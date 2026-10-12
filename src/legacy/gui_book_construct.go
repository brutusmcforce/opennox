package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

func bookEvent(fn func(*gui.Window, uint32, uint32) int) gui.WindowFunc {
	return func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, _ := ev.EventArgsC()
		return gui.RawEventResp(fn(w, uint32(ev.EventCode()), uint32(a)))
	}
}
func bookPoint(v uint32) image.Point { return image.Pt(int(uint16(v)), int(v>>16)) }
func bookInit() int {
	bookPlayer = dword_8531A0_2576.C()
	for _, r := range []struct {
		name string
		off  uintptr
	}{
		{"ArrowNW", 1046888}, {"ArrowN", 1046892}, {"ArrowNE", 1046896}, {"ArrowW", 1046900},
		{"ArrowE", 1046908}, {"ArrowSW", 1046912}, {"ArrowS", 1046916}, {"ArrowSE", 1046920},
		{"BookOfKnowledge", 1046856}, {"GuideTabLit", 1046660}, {"SpellTabLit", 1046644},
	} {
		v := uiMeterLoadImage(r.name)
		*bookWord(r.off) = v
		if v == 0 {
			return 0
		}
	}
	for i, name := range []string{"BookPageForward", "BookPageBackward"} {
		ref := Nox_xxx_gLoadAnim(name)
		*bookWord(1046924 + 4*uintptr(i)) = memmap.PtrWord(ref.C())
		if ref == nil {
			return 0
		}
		ref.Field24ptr().OnEnd = bookImageEndKey(0)
	}
	g := GetClient().Cli().GUI
	drawOne := func(*gui.Window, *gui.WindowData) int { return 1 }
	root := g.NewWindowRaw(nil, 1196, 5, int(nox_win_height)-323, 285, 168, nil)
	legacyGlobals.nox_win_unk1 = root
	root.SetAllFuncs(bookEvent(func(w *gui.Window, e, a uint32) int { return bookListEvents(w, e, bookPoint(a)) }), func(w *gui.Window, _ *gui.WindowData) int { return bookDrawList(w) }, nil)
	bookHideWindow(root, true)
	for _, tab := range []struct{ x, y, id int }{{257, 15, 1320}, {253, 61, 1310}} {
		w := g.NewWindowRaw(root, 8, tab.x, tab.y, 27, 40, nil)
		if w == nil {
			return 0
		}
		w.SetAllFuncs(bookEvent(func(w *gui.Window, e, _ uint32) int { return bookTab(w, e) }), drawOne, bookTooltipKey(4))
		w.SetID(uint(tab.id))
	}
	back := g.NewWindowRaw(root, 136, 24, 138, 20, 20, nil)
	*bookWord(1046944) = memmap.PtrWord(back.C())
	back.SetAllFuncs(bookEvent(func(_ *gui.Window, e, _ uint32) int { return bookBackward(int(e)) }), nil, nil)
	back.DrawData().BgImageHnd = noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("ArrowW")))
	forward := g.NewWindowRaw(root, 136, 233, 138, 20, 20, nil)
	*bookWord(1046948) = memmap.PtrWord(forward.C())
	forward.SetAllFuncs(bookEvent(func(_ *gui.Window, e, _ uint32) int { return bookForward(int(e)) }), nil, nil)
	forward.DrawData().BgImageHnd = noxrender.ImageHandle(memmap.WordPtr(uiMeterLoadImage("ArrowE")))
	icon := g.NewWindowRaw(root, 8, 63, 19, 30, 30, nil)
	*bookWord(1046952) = memmap.PtrWord(icon.C())
	icon.SetAllFuncs(bookEvent(func(w *gui.Window, e, a uint32) int { return bookIconEvents(w, int(e), bookPoint(a)) }), func(w *gui.Window, _ *gui.WindowData) int { return bookDrawIcon(w) }, nil)
	bookHideWindow(icon, true)
	moving := g.NewWindowRaw(nil, 40, 0, 0, 30, 30, nil)
	*bookWord(1046956) = memmap.PtrWord(moving.C())
	moving.SetAllFuncs(nil, func(w *gui.Window, _ *gui.WindowData) int { return bookDrawAddition(w) }, nil)
	bookHideWindow(moving, true)
	bookPageComplete(false)
	*bookContents() = 1
	*bookWord(1046936) = 0
	return 1
}
