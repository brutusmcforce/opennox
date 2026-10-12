package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/legacy/timer"
)

type optionsEditor bool

const (
	optionsInGame optionsEditor = false
	optionsMenu   optionsEditor = true
)

func optionsWord(off int) *uint32 {
	switch off {
	case 1309720:
		return (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1309720))
	case 1309728:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309728))
	case 1309732:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309732))
	case 1309736:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309736))
	case 1309820:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309820))
	case 1309824:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309824))
	case 1309828:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309828))
	case 1309832:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309832))
	case 1309836:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1309836))
	case 172880:
		return (*uint32)(unsafe.Pointer(&nox_xxx_normalWndBits_587000_172880))
	}
	panic(off)
}
func optionsWindow(off int) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*optionsWord(off))))
}
func optionsStore(off int, w *gui.Window) { *optionsWord(off) = uint32(uintptr(w.C())) }
func (e optionsEditor) rootOffset() int {
	if e {
		return 1309720
	}
	return 1309820
}
func (e optionsEditor) buttonOffset(ch int) int {
	if e {
		return [...]int{1309728, 1309732, 1309736}[ch]
	}
	return [...]int{1309828, 1309836, 1309832}[ch]
}
func optionsTimer(ch int) *timer.Timer {
	switch ch {
	case 0:
		return (*timer.Timer)(Get_dword_587000_127004())
	case 1:
		return (*timer.Timer)(Get_dword_587000_122852())
	case 2:
		return (*timer.Timer)(Get_dword_587000_93164())
	}
	panic(ch)
}
func optionsEnabled(ch int) int {
	switch ch {
	case 0:
		return Sub_453070()
	case 1:
		return Sub_44D990()
	case 2:
		return Sub_43DC30()
	}
	panic(ch)
}
func optionsToggle(ch int) {
	switch ch {
	case 0:
		if optionsEnabled(ch) == 1 {
			Sub_453050()
		} else {
			*audioEventPlayback = 1
		}
	case 1:
		if optionsEnabled(ch) == 1 {
			Sub_44D960()
		} else {
			audioEventDialogEnable()
		}
	case 2:
		if optionsEnabled(ch) == 1 {
			Sub_43DC00()
		} else {
			audioEventMusicEnable()
			t := optionsTimer(ch)
			// Enabling music resets its target from the current volume.
			t.SetRaw(t.Current >> 16)
		}
	}
}
func optionsPreview() {
	if Dialogs.Sub_44D930() {
		return
	}
	cursor := memmap.PtrUint32(0x5D4594, 1309744)
	off := uint32(172892) + 4*(*cursor)
	name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, uintptr(off))))
	(*cursor)++
	v, _ := GetServer().S().Strings().GetVariantInFile(strman.ID(name), "C:\\NoxPost\\src\\client\\shell\\Options.c")
	*cursor %= 3
	if v.Str2 != "" {
		Dialogs.PlayFile(v.Str2, 100)
	}
}
func optionsSend(w *gui.Window, event int, a, b uintptr) int {
	return gui.EventRespInt(w.Func94(gui.AsWindowEvent(event, a, b)))
}
func (e optionsEditor) volume(w *gui.Window, event, ch, value int) {
	if event == 16396 && ch == 2 {
		return
	}
	// Even a suppressed drag updates the timer target; release ignores music.
	optionsTimer(ch).SetRaw(uint32(value))
	if event == 16393 && ch == 0 {
		current := GetClient().Cli().GUI.WinYYY
		if current != nil && current.Parent() == w {
			return
		}
	}
	enabled := optionsEnabled(ch)
	if value != 0 {
		if enabled == 0 {
			optionsWindow(e.buttonOffset(ch)).Func93(gui.AsWindowEvent(21, 28, 2))
		}
		if ch == 0 {
			Nox_xxx_clientPlaySoundSpecial_452D80(768, 100)
		} else if ch == 1 {
			optionsPreview()
		}
	} else if enabled == 1 {
		optionsWindow(e.buttonOffset(ch)).Func93(gui.AsWindowEvent(21, 28, 2))
	}
}
func optionsSensitivity(value int) float32 {
	// C retains precision through division/subtraction and rounds at the powf call.
	exponent := float32(float64(float32(value))/50 - 1)
	return float32(math.Pow(10, float64(exponent)))
}
func optionsViewportID(cut int) uint {
	switch {
	case cut <= 69:
		return 311
	case cut <= 79:
		return 312
	case cut <= 89:
		return 313
	default:
		return 314
	}
}
