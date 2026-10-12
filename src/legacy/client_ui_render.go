package legacy

import (
	"image"
	"unsafe"
)

func uiRenderIntersect(out, a, b *[4]int32) *[4]int32 {
	x1, x2 := max(a[0], b[0]), min(a[2], b[2])
	if x1 >= x2 {
		return nil
	}
	y1, y2 := max(a[1], b[1]), min(a[3], b[3])
	if y1 >= y2 {
		return nil
	}
	*out = [4]int32{x1, y1, x2, y2}
	return out
}
func uiRenderIntersectRects(a, b image.Rectangle) (image.Rectangle, bool) {
	aa := [4]int32{int32(a.Min.X), int32(a.Min.Y), int32(a.Max.X), int32(a.Max.Y)}
	bb := [4]int32{int32(b.Min.X), int32(b.Min.Y), int32(b.Max.X), int32(b.Max.Y)}
	var out [4]int32
	if uiRenderIntersect(&out, &aa, &bb) == nil {
		return image.Rectangle{}, false
	}
	return image.Rectangle{Min: image.Pt(int(out[0]), int(out[1])), Max: image.Pt(int(out[2]), int(out[3]))}, true
}
func uiRenderFlag(v uint32) uint32 {
	p := (*uint32)(objectRenderClipData().C())
	old := *p
	*p = v
	return old
}
func uiRenderCopyRect(x, y, w, h int) int {
	d := objectRenderClipData()
	rect, ok := uiRenderIntersectRects(image.Rectangle{Min: image.Pt(x, y), Max: image.Pt(x+w, y+h)}, d.Rect3())
	if !ok {
		return 0
	}
	d.SetClipRect(rect)
	rect.Max.X--
	rect.Max.Y--
	d.SetClipRect2(rect)
	// noxCopyRect returns scalar true despite the old pointer return declaration.
	return 1
}
func uiRenderNarrowClip(left, right int) int {
	r := objectRenderClipData().ClipRect()
	left = max(left, r.Min.X)
	right = min(right, r.Max.X)
	return uiRenderCopyRect(left, r.Min.Y, right-left, r.Dy())
}
func uiRenderBounds(x1, y1, x2, y2 int) int {
	dword_5d4594_3807140 = uint32(x1)
	dword_5d4594_3807136 = uint32(y1)
	dword_5d4594_3807116 = uint32(x2)
	dword_5d4594_3807152 = uint32(y2)
	return y2
}
func uiRenderFill(p unsafe.Pointer, color uint32, size int32) {
	if size <= 0 {
		return
	}
	n := int(size) / 4
	for i := 0; i < n; i++ {
		*(*uint32)(unsafe.Add(p, 4*i)) = color
	}
	if size&3 >= 2 {
		*(*uint16)(unsafe.Add(p, 4*n)) = uint16(color)
	}
}
func uiRenderBorder(x, y, w, h, color, width int) {
	if w == 0 || h == 0 {
		return
	}
	d := objectRenderClipData()
	if *(*uint32)(d.C()) != 0 {
		_, ok := uiRenderIntersectRects(image.Rectangle{Min: image.Pt(x, y), Max: image.Pt(x+w, y+h)}, d.ClipRect())
		if !ok {
			return
		}
	}
	nox_draw_set54RGB32_434040(color)
	r := GetClient().R2()
	r.Data().SetField262(width)
	line := func(x1, y1, x2, y2 int) {
		r.AddPoint(image.Pt(x1, y1))
		r.AddPoint(image.Pt(x2, y2))
		r.DrawParticles49ED80(64)
	}
	line(x, y, x+w-1, y)
	line(x+w, y, x+w, y+h-1)
	line(x+w, y+h, x, y+h)
	line(x, y+h-1, x, y+1)
}
