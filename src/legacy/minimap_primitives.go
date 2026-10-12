package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

var minimapZoom uint32 = 0x456

func minimapZoomIn() {
	minimapZoom -= 10
	if int32(minimapZoom) < 500 {
		minimapZoom = 500
	}
}
func minimapZoomOut() {
	minimapZoom += 10
	if minimapZoom > 4000 {
		minimapZoom = 4000
	}
}
func minimapSetZoom(v int) int { minimapZoom = uint32(v); return v }
func minimapLine(a, b image.Point) int {
	r := GetClient().R2()
	r.AddPoint(a)
	r.AddPoint(b)
	return bool2int(r.DrawLineFromPoints(r.Data().Color2()))
}
func minimapColoredLine(x1, y1, x2, y2, color int) int {
	GetClient().R2().Data().SetColor2(noxcolor.RGBA5551(color))
	return minimapLine(image.Pt(x1, y1), image.Pt(x2, y2))
}
func minimapWall(p image.Point, dir byte, size int) int {
	color := int(memmap.Uint32(0x85B3FC, 956))
	if minimapZoom > 2000 {
		r := GetClient().R2()
		r.Data().SetColor2(noxcolor.RGBA5551(color))
		r.DrawPixel(p, r.Data().Color2())
		return 0
	}
	h := size / 2
	line := func(x1, y1, x2, y2 int) int { return minimapColoredLine(p.X+x1, p.Y+y1, p.X+x2, p.Y+y2, color) }
	switch dir {
	case 0:
		return line(0, size, size, 0)
	case 1:
		return line(0, 0, size, size)
	case 2:
		line(0, size, size, 0)
		return line(0, 0, size, size)
	case 3:
		line(0, 0, h, h)
		return line(0, size, size, 0)
	case 4:
		line(0, 0, size, size)
		return line(h, h, size, 0)
	case 5:
		line(0, size, size, 0)
		return line(h, h, size, size)
	case 6:
		line(0, 0, size, size)
		return line(0, size, h, h)
	case 7:
		line(0, 0, h, h)
		return line(h, h, size, 0)
	case 8:
		line(h, h, size, size)
		return line(h, h, size, 0)
	case 9:
		line(h, h, size, size)
		return line(h, h, 0, size)
	case 10:
		line(0, 0, h, h)
		return line(0, size, h, h)
	}
	return int(dir)
}
func minimapFlag(p image.Point) int {
	points := [...]image.Point{{-2, 4}, {-2, -4}, {2, -4}, {2, 0}, {-2, 0}}
	ret := 0
	for i := 1; i < len(points); i++ {
		ret = minimapLine(p.Add(points[i-1]), p.Add(points[i]))
	}
	return ret
}
func minimapCrown(p image.Point) int {
	points := [...]image.Point{{-4, 6}, {-6, -6}, {-2, 0}, {0, -6}, {2, 0}, {6, -6}, {4, 6}, {-4, 6}}
	ret := 0
	for i := 1; i < len(points); i++ {
		ret = minimapLine(p.Add(points[i-1]), p.Add(points[i]))
	}
	return ret
}
func minimapCircle(p image.Point) {
	r := GetClient().R2()
	r.DrawCircle(p.X, p.Y, 3, r.Data().Color2())
}

// The original shadow border has top, right and bottom strokes, with an open left.
func minimapOutline(x, y, w, h int) int {
	right, bottom := x+w-1, y+h-1
	minimapLine(image.Pt(x, y), image.Pt(right, y))
	minimapLine(image.Pt(right, y), image.Pt(right, bottom))
	return minimapLine(image.Pt(right, bottom), image.Pt(x, bottom))
}
func minimapPoint(x, y int, large bool) {
	size := 4
	if minimapZoom > 1200 {
		size = 2
		if minimapZoom < 1750 {
			size++
		}
	}
	if large {
		size += 2
	}
	r := GetClient().R2()
	r.DrawPoint(image.Pt(x, y), size, r.Data().Color2())
}
