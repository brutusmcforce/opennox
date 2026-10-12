package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"image"
	"unsafe"
)

// These helpers intentionally use the current C render-data pointer: remaining
// C UI callers save and restore that state, which can be switched independently.
func objectRenderClipData() *noxrender.RenderData {
	return (*noxrender.RenderData)(unsafe.Pointer(legacyGlobals.nox_draw_curDrawData_3799572))
}
func objectRenderSaveClip() {
	if dword_5d4594_1305748 != 0 {
		return
	}
	d := objectRenderClipData()
	clip, rect := d.ClipRect(), d.ClipRect2()
	*effectMapped(1305772) = *(*uint32)(d.C())
	for i, v := range []int{clip.Min.X, clip.Min.Y, clip.Max.X, clip.Max.Y} {
		*effectMapped(1305756 + 4*uintptr(i)) = uint32(v)
	}
	for i, v := range []int{rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y} {
		*effectMapped(1305732 + 4*uintptr(i)) = uint32(v)
	}
	dword_5d4594_1305748 = 1
}
func objectRenderRestoreClip() int {
	if dword_5d4594_1305748 == 0 {
		return 0
	}
	d := objectRenderClipData()
	*(*uint32)(d.C()) = *effectMapped(1305772)
	rect := func(off uintptr) image.Rectangle {
		return image.Rectangle{Min: image.Pt(int(*effectMapped(off)), int(*effectMapped(off + 4))), Max: image.Pt(int(*effectMapped(off + 8)), int(*effectMapped(off + 12)))}
	}
	d.SetClipRect(rect(1305756))
	d.SetClipRect2(rect(1305732))
	dword_5d4594_1305748 = 0
	return int(*effectMapped(1305740))
}
