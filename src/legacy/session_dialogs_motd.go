package legacy

import (
	"io"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

var sessionMOTDRoot, sessionMOTDList *gui.Window
var sessionMOTDFile *byte

func sessionMOTDRead(slot int) int {
	sessionMOTDFile = nil
	size := memmap.PtrUint32(0x5D4594, 826040+uintptr(slot)*4)
	*size = 0
	f, err := ifs.Open("motd.txt")
	if err != nil {
		return 0
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0
	}
	// The supported client has a 32-bit address space. Reject an unrepresentable
	// allocation instead of the old wrapped allocation or nil dereference.
	if info.Size() < 0 || uint64(info.Size()) >= uint64(^uint32(0)>>1) {
		return 0
	}
	n := int(info.Size())
	p, release := alloc.Calloc(n+1, 1)
	if p == nil {
		release()
		return 0
	}
	sessionMOTDFile = (*byte)(p)
	io.ReadFull(f, unsafe.Slice(sessionMOTDFile, n))
	*size = uint32(n + 1)
	return 0
}
func sessionMOTDFree(slot int) uintptr {
	var result uintptr
	if sessionMOTDFile != nil {
		alloc.Free(sessionMOTDFile)
		result = uintptr(uint32(slot))
	}
	sessionMOTDFile = nil
	*memmap.PtrUint32(0x5D4594, 826040+uintptr(uint32(slot))*4) = 0
	return result
}
func sessionMOTDOpen() *gui.Window {
	w := Nox_new_window_from_file("motd.wnd", sessionMOTDEvent)
	sessionMOTDRoot = w
	sessionMOTDList = w.ChildByID(4203)
	normal, lit := Nox_xxx_gLoadImg("UISlider"), Nox_xxx_gLoadImg("UISliderLit")
	slider, up, down := w.ChildByID(4204), w.ChildByID(4205), w.ChildByID(4206)
	browserWireList(sessionMOTDList, slider, up, down, 10)
	gui.ButtonSetImage(slider, nil, nil, normal.C(), lit.C(), lit.C())
	return w
}
func sessionMOTDEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16391 {
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		sessionMOTDClose()
	}
	return nil
}
func sessionMOTDNext(s string) (string, int) {
	i := strings.IndexAny(s, "\r\n")
	if i < 0 {
		return s, -1
	}
	next := i + 1
	if s[i] == '\r' && next < len(s) && s[next] == '\n' {
		next++
	}
	return s[:i], next
}
func sessionMOTDLine(src, dst *byte) *byte {
	line, next := sessionMOTDNext(alloc.GoString(src))
	out := unsafe.Slice(dst, len(line)+1)
	copy(out, line)
	out[len(line)] = 0
	if next < 0 {
		return nil
	}
	return (*byte)(unsafe.Add(unsafe.Pointer(src), next))
}
func sessionMOTDAddText(text string) uintptr {
	wide, free := alloc.Make([]uint16{}, len(text)+1)
	defer free()
	for i := 0; i < len(text); i++ {
		wide[i] = uint16(text[i])
	}
	return uintptr(optionsSend(sessionMOTDList, 16397, uintptr(unsafe.Pointer(&wide[0])), ^uintptr(0)))
}
func sessionMOTDAdd(p unsafe.Pointer) uintptr {
	text := alloc.GoString((*byte)(p))
	if text == "" {
		return uintptr(p)
	}
	return sessionMOTDAddText(text)
}
func sessionMOTDClose() int {
	if sessionMOTDRoot == nil || sessionMOTDRoot.GetFlags().IsHidden() {
		return 0
	}
	GetClient().Cli().GUI.Focus(nil)
	sessionMOTDRoot.Hide()
	sessionMOTDRoot.Flags &^= gui.StatusEnabled
	sessionMOTDList.Flags &^= gui.StatusEnabled
	optionsSend(sessionMOTDList, 16399, 0, 0)
	return 1
}
func sessionMOTDShown() int {
	return bool2int(sessionMOTDRoot != nil && sessionMOTDRoot.Flags&gui.StatusHidden == 0)
}
func sessionMOTDShow() {
	if sessionQuitShown() != 0 {
		return
	}
	if questRuntimeWord(1556160) != 0 && noxflags.HasGame(128) {
		return
	}
	if !noxflags.HasEngine(noxflags.EngineNoRendering) && sessionMOTDRoot != nil && sessionMOTDRoot.Flags&gui.StatusHidden != 0 && Sub_44A4A0() == 0 && sub_49C810() == 0 {
		sessionMOTDRoot.ShowModal()
		// The legacy child lookup with a nil parent (ID4100) never finds a window.
		sessionMOTDRoot.Flags |= gui.StatusEnabled
		sessionMOTDList.Flags |= gui.StatusEnabled
		p := (*byte)(*memmap.PtrPtr(0x5D4594, 826060))
		if p != nil {
			text := alloc.GoString(p)
			for {
				line, next := sessionMOTDNext(text)
				if next < 0 {
					if line != "" {
						sessionMOTDAddText(line)
					}
					break
				}
				if line == "" {
					line = " "
				}
				sessionMOTDAddText(line)
				text = text[next:]
			}
		}
		sessionMOTDRoot.ChildByID(4202).Focus()
	}
	*memmap.PtrUint32(0x5D4594, 826068) = 0
}
