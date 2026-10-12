package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

// The original distance estimate scales an eight-bit lookup table in powers of
// four. Preserve its approximation and the 32-bit wrap of squared coordinates.
func screenSqrt(v uint32) uint32 {
	if v < 4 {
		return uint32(*memmap.PtrUint8(0x587000, 155956+uintptr(v*64))) >> 7
	}
	if v < 16 {
		return uint32(*memmap.PtrUint8(0x587000, 155956+uintptr(v*16))) >> 6
	}
	if v < 64 {
		return uint32(*memmap.PtrUint8(0x587000, 155956+uintptr(v*4))) >> 5
	}
	if v < 256 {
		return uint32(*memmap.PtrUint8(0x587000, 155956+uintptr(v))) >> 4
	}
	shift := uint(2)
	for v>>shift >= 256 {
		shift += 2
	}
	value := uint32(*memmap.PtrUint8(0x587000, 155956+uintptr(v>>shift)))
	if shift < 8 {
		return value >> (4 - shift/2)
	}
	return value << (shift/2 - 4)
}
func screenDistance(x, y int32) uint32                  { a, b := uint32(x), uint32(y); return screenSqrt(a*a + b*b) }
func screenDistanceBetween(x1, y1, x2, y2 int32) uint32 { return screenDistance(x2-x1, y2-y1) }

func screenRopeLine(a, b image.Point) int {
	r := GetClient().R2()
	line := func(a, b image.Point) bool {
		r.AddPoint(a)
		r.AddPoint(b)
		return r.DrawLineFromPoints(r.Data().Color2())
	}
	effectColor(*effectMapped(1312492))
	line(a, b)
	dx, dy := b.X-a.X, b.Y-a.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	effectColor(*effectMapped(1312496))
	off := image.Pt(0, 1)
	if dx <= dy {
		off = image.Pt(1, 0)
	}
	line(a.Sub(off), b.Sub(off))
	if line(a.Add(off), b.Add(off)) {
		return 1
	}
	return 0
}
func screenCircle(x, y, radius int, color uint32) int {
	// Multiplication and logical shifts are intentionally unsigned, matching C.
	point := func(off uintptr) image.Point {
		return image.Pt(x+int(uint32(radius)**memmap.PtrUint32(0x587000, off)>>4), y+int(uint32(radius)**memmap.PtrUint32(0x587000, off+4)>>4))
	}
	first := point(192088)
	prev := first
	effectColor(color)
	GetClient().R2().Data().SetAlphaEnabled(true)
	for off := uintptr(192216); off < 194136; off += 128 {
		next := point(off)
		effectLine(prev, next)
		prev = next
	}
	effectLine(first, prev)
	GetClient().R2().Data().SetAlphaEnabled(false)
	return 1
}
