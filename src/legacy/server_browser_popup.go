package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func browserPopup(parent *gui.Window, point *[2]uint32, head *legacyListNode) *gui.Window {
	browserUI.popupCount = 0
	for it := listNext(head); it != nil; it = listNext(it) {
		if browserHit(point, unsafe.Pointer(it)) {
			*memmap.PtrUint32(0x5D4594, 1307316+4*uintptr(browserUI.popupCount)) = uint32(uintptr(unsafe.Pointer(it)))
			browserUI.popupCount++
		}
	}
	if browserUI.popupCount == 0 {
		return browserWindow(uint32(browserUI.popup))
	}
	// The retired raw slot has no producer. Preserve its nil parser callback;
	// inheriting the native parent handler would change existing popup behavior.
	if *(*unsafe.Pointer)(unsafe.Add(parent.C(), 376)) != nil {
		panic("unregistered raw popup callback")
	}
	w := Nox_new_window_from_file("proxlist.wnd", nil)
	browserUI.popup = uint32(uintptr(w.C()))
	var pos [2]uint32
	browserPopupClamp(int32(point[0]+216), int32(point[1]+27), &pos)
	w.SetPos(image.Pt(int(int32(pos[0])), int(int32(pos[1]))))
	w.DrawData().Window = parent
	slider, up, down, list := w.ChildByID(10064), w.ChildByID(10062), w.ChildByID(10063), w.ChildByID(10061)
	en := noxrender.ImageHandle(unsafe.Pointer(uintptr(uiMeterLoadImage("UISlider"))))
	lit := noxrender.ImageHandle(unsafe.Pointer(uintptr(uiMeterLoadImage("UISliderLit"))))
	gui.ButtonSetImage(slider, nil, nil, en, lit, lit)
	slider.DrawData().Window = list
	up.DrawData().Window = list
	down.DrawData().Window = list
	data := uiListData(list)
	data.Field_9 = slider.C()
	data.Field_7 = up.C()
	data.Field_8 = down.C()
	slider.Field100Ptr.SizeVal = image.Pt(16, 10)
	for i := uint32(0); i < uint32(browserUI.popupCount); i++ {
		p := unsafe.Pointer(uintptr(memmap.Uint32(0x5D4594, 1307316+4*uintptr(i))))
		name := browserCappedString(unsafe.Add(p, 120), 15)
		if name == "" {
			name = fmt.Sprintf("%s:%d", alloc.GoString((*byte)(unsafe.Add(p, 12))), *(*uint16)(unsafe.Add(p, 109)))
		}
		browserText(list, fmt.Sprintf("%s   %dms", browserNarrowText(name), *(*int32)(unsafe.Add(p, 96))), -1)
	}
	return w
}
func browserPopupClose() int {
	w := browserWindow(uint32(browserUI.popup))
	if w == nil {
		return 0
	}
	w.Destroy()
	browserUI.popup = 0
	return 0
}
func browserPopupAt(index int32) unsafe.Pointer {
	if index >= int32(browserUI.popupCount) {
		return nil
	}
	return unsafe.Pointer(uintptr(memmap.Uint32(0x5D4594, uintptr(uint32(1307316+4*index)))))
}

func sub_4A2610(parent int32, point *uint32, head *int32) int32 {
	return int32(uintptr(browserPopup(browserWindow(uint32(parent)), (*[2]uint32)(unsafe.Pointer(point)), (*legacyListNode)(unsafe.Pointer(head))).C()))
}

func sub_4A2890() int32 { return int32(browserPopupClose()) }

func sub_4A28B0() int32 {
	if browserUI.popup != 0 {
		return 1
	}
	return 0
}

func sub_4A28C0(index int32) int32 { return int32(uintptr(browserPopupAt(int32(index)))) }
