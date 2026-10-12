package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"image"
)

// Update the image interval cache without changing the renderer's clip rectangle.
func worldWallImageInterval(a, b int) bool {
	lo, hi := min(a, b), max(a, b)
	clip := objectRenderClipData().ClipRect()
	if lo < clip.Min.X {
		lo = clip.Min.X
	} else if lo >= clip.Max.X {
		return false
	}
	if hi >= clip.Max.X {
		hi = clip.Max.X
	} else if hi < clip.Min.X {
		return false
	}
	if lo == hi {
		return false
	}
	if lo != clip.Min.X || hi != clip.Max.X {
		*memmap.PtrUint32(0x973F18, 52) = uint32(lo)
		*memmap.PtrUint32(0x973F18, 12) = uint32(hi)
		dword_5d4594_3799452 = 1
	}
	return true
}
func worldWallDraw(vp *noxrender.Viewport, w *server.Wall) {
	if w == nil || w.Flags4&1 == 0 {
		return
	}
	c := GetClient()
	r := c.R2()
	data := r.Data()
	width, height := int(nox_win_width), int(nox_win_height)
	origin := vp.ToScreenPos(image.Pt(23*int(w.X5), 23*int(w.Y6)))
	dir := int(memmap.Int32(0x587000, 149364+4*uintptr(w.Field3)))
	if dir == -1 {
		dir = int(w.Dir0)
	}
	sprite := dir
	if w.Flags4&0x40 != 0 {
		if dir == 0 {
			sprite = 11
		} else if dir == 1 {
			sprite = 12
		}
	}
	lo, hi := 0, width
	if memmap.Uint32(0x587000, 80808) != 0 {
		if !worldWallVisibleSpan(origin, dir, width, &lo, &hi, &sprite) {
			return
		}
		if lo >= hi {
			w.Field3 = 0
			w.Flags4 &^= 3
			return
		}
	}
	def := GetServer().S().Walls.DefByInd(int(w.Tile1))
	front := w.Flags4&2 != 0
	highFront, highFloors, translucent := nox_client_highResFrontWalls_80820 != 0, nox_client_highResFloors_154952 != 0, nox_client_translucentFrontWalls_805844 != 0
	options := 0
	layer := int((w.Flags4 >> 2) & 2)
	if front {
		if memmap.Uint32(0x5D4594, 805848) != 0 && translucent {
			if !highFront && highFloors {
				options = 4
			} else {
				options = 8
			}
		}
		if !highFront {
			options |= 4
		}
		layer = int((w.Flags4&8 | 4) >> 2)
	}
	if front && translucent && def.Flags32&4 == 0 {
		options |= 2
	} else {
		options |= 1
	}
	offset := def.DrawOffset(sprite, int(w.Field2), layer)
	img := def.Sprite(sprite, int(w.Field2), layer)
	x, y := 23*int(w.X5), 23*int(w.Y6)
	if memmap.Uint32(0x587000, 80816) != 0 {
		var a, b image.Point
		switch dir {
		case 0, 3:
			a, b = image.Pt(x, y+23), image.Pt(x+23, y)
		case 1, 4:
			a, b = image.Pt(x, y), image.Pt(x+23, y+23)
		case 7:
			a, b = image.Pt(x, y), image.Pt(x+23, y)
		case 8:
			a, b = image.Pt(x+11, y+11), image.Pt(x-23, y-23)
		case 10:
			a, b = image.Pt(x, y+23), image.Pt(x+11, y+12)
		default:
			a, b = image.Pt(x, y+23), image.Pt(x+23, y+23)
		}
		first := c.Sub469920(a)
		// The real light sampler reuses its scratch buffer on the next call.
		color := [3]uint32{first[0], first[1], first[2]}
		second := c.Sub469920(b)
		p := image.Pt(origin.X+offset.X-51, origin.Y-offset.Y-73)
		data.SetMultiply14(1)
		r.SetColorMultAndIntensityRGB(byte(color[0]), byte(color[1]), byte(color[2]))
		if options&2 == 0 {
			wallEdgeDraw(noxrender.ImageHandle(img), p, &color, (*[3]uint32)(second), height, lo, hi, 0)
		} else if worldWallImageInterval(lo, hi) {
			data.SetAlphaEnabled(true)
			data.SetAlpha(128)
			r.DrawImageAt(r.GetBag().AsImage(noxrender.ImageHandle(img)), p)
			data.SetAlphaEnabled(false)
		}
	} else {
		color := c.Sub469920(image.Pt(x+11, y+11))
		p := image.Pt(origin.X+offset.X-50, origin.Y-offset.Y-72)
		data.SetMultiply14(1)
		r.SetColorMultAndIntensityRGB(byte(color[0]), byte(color[1]), byte(color[2]))
		if worldWallImageInterval(lo, hi) {
			if options&2 != 0 {
				data.SetAlphaEnabled(true)
				data.SetAlpha(128)
			}
			r.DrawImageAt(r.GetBag().AsImage(noxrender.ImageHandle(img)), p)
			if options&2 != 0 {
				data.SetAlphaEnabled(false)
			}
		}
	}
	data.SetMultiply14(0)
	w.Field3 = 0
	w.Flags4 &^= 3
	w.Field12 = 1
}
