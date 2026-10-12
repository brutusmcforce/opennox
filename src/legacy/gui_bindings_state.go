package legacy

import (
	"fmt"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// bindingSlot returns the variable that holds a window of the key binding editors.
func bindingSlot(off int) **gui.Window {
	switch off {
	case 1321224:
		return &dword_5d4594_1321224
	case 1321228:
		return &dword_5d4594_1321228
	case 1321232:
		return &dword_5d4594_1321232
	case 1321236:
		return &legacyGlobals.dword_5d4594_1321236
	case 1321240:
		return &legacyGlobals.dword_5d4594_1321240
	case 1321244:
		return &legacyGlobals.dword_5d4594_1321244
	case 1321248:
		return &legacyGlobals.dword_5d4594_1321248
	case 1321252:
		return &dword_5d4594_1321252
	case 1522604:
		return &dword_5d4594_1522604
	case 1522612:
		return &dword_5d4594_1522612
	case 1522616:
		return &legacyGlobals.dword_5d4594_1522616
	case 1522620:
		return &legacyGlobals.dword_5d4594_1522620
	case 1522624:
		return &legacyGlobals.dword_5d4594_1522624
	case 1522628:
		return &legacyGlobals.dword_5d4594_1522628
	case 1522632:
		return &dword_5d4594_1522632
	}
	panic(off)
}

// bindingEditor selects the two existing UI owners while sharing their row logic.
type bindingEditor bool

const (
	bindingInGame bindingEditor = false
	bindingMenu   bindingEditor = true
)

func (e bindingEditor) offsets() (root, modal, base, selected, buffer int) {
	if e {
		return 1522604, 1522612, 1522616, 1522632, 1522636
	}
	return 1321228, 1321232, 1321236, 1321252, 1321256
}
func bindingWindow(off int) *gui.Window {
	return *bindingSlot(off)
}
func bindingStore(off int, w *gui.Window)              { *bindingSlot(off) = w }
func bindingList(w *gui.Window) *gui.ScrollListBoxData { return (*gui.ScrollListBoxData)(w.WidgetData) }
func bindingSend(w *gui.Window, event int, a, b uintptr) int {
	return gui.EventRespInt(w.Func94(gui.AsWindowEvent(event, a, b)))
}
func bindingSetText(w *gui.Window, event int, text string, row int) {
	bindingSend(w, event, uintptr(unsafe.Pointer(alloc.InternCString16(text))), uintptr(row))
}
func bindingGetText(w *gui.Window, row int) string {
	return GoWStringP(unsafe.Pointer(uintptr(bindingSend(w, 16406, uintptr(row), 0))))
}
func (e bindingEditor) assign(key uint32, mouse bool) int {
	k := keybind.Key(key)
	if !mouse && !k.IsValid() {
		return 0
	}
	_, _, base, selected, _ := e.offsets()
	dst := bindingWindow(selected)
	if dst == nil {
		return 1
	}
	title := ""
	if k.IsValid() {
		title = k.Title(GetClient().Strings())
	}
	for _, off := range []int{base + 8, base + 12} {
		w := bindingWindow(off)
		count := bindingList(w).Field_11_0
		// C enters on an unsigned count, then checks a signed count in its do/while.
		n := int(int16(count))
		if count != 0 && n < 1 {
			n = 1
		}
		for row := 0; row < n; row++ {
			if bindingGetText(w, row) == title {
				bindingSetText(w, 16407, " ", row)
			}
		}
	}
	bindingSetText(dst, 16407, title, int(int32(bindingList(dst).Field_12)))
	bindingSend(dst, 16403, ^uintptr(0), 0)
	*bindingSlot(selected) = nil
	return 1
}
func (e bindingEditor) apply() uint32 {
	_, _, base, _, _ := e.offsets()
	titles := bindingWindow(base + 4)
	GetClient().GetCtrlEvent().Reset()
	bind := func(key, title string) {
		k := Nox_xxx_keybind_nameByTitle_42E960(key)
		if k == 0 {
			return
		}
		event := Nox_xxx_bindevent_bindNameByTitle_42EA40(title)
		name := "(null)"
		if event != nil {
			name = event.Name
		}
		Nox_client_parseConfigHotkeysLine_42CF50(fmt.Sprintf("%s = %s", k.String(), name))
	}
	for row := 0; row < int(int16(bindingList(titles).Field_11_0)); row++ {
		title := bindingGetText(titles, row)
		bind(bindingGetText(bindingWindow(base+12), row), title)
		bind(bindingGetText(bindingWindow(base+8), row), title)
	}
	sm := GetClient().Strings()
	escape := sm.GetString("keybind:Esc")
	if Nox_xxx_keybind_nameByTitle_42E960(escape) == 0 {
		return 0
	}
	bind(escape, sm.GetString("bindevent:ToggleQuitMenu"))
	return 1
}
func (e bindingEditor) modalEvent(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	_, modal, _, selected, _ := e.offsets()
	close := func() {
		GetClient().Cli().GUI.Focus(nil)
		bindingWindow(modal).StackPop()
		bindingWindow(modal).SetHidden(true)
	}
	key := uint32(0)
	switch ev.EventCode() {
	case 6, 7:
		key = 0x10000
	case 10, 11:
		key = 0x10002
	case 14, 15:
		key = 0x10001
	case 19:
		key = 0x10003
	case 20:
		key = 0x10004
	case 21:
		if a == 1 {
			if b == 2 {
				close()
				if w := bindingWindow(selected); w != nil {
					bindingSend(w, 16403, ^uintptr(0), 0)
				}
			}
			return gui.RawEventResp(1)
		}
		if b != 1 || e.assign(uint32(a), false) == 0 {
			return nil
		}
		key = 0x10000
	default:
		return nil
	}
	e.assign(key, true)
	close()
	return gui.RawEventResp(1)
}
func bindingFilter(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 19 || ev.EventCode() == 20 {
		return nil
	}
	return uiListSingleInput(w, ev)
}
