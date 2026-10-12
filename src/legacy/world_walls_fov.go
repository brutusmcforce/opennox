package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
)

// Restrict the wall's horizontal span using the real field-of-view scanlines.
// Some corner shapes select a different sprite when only one segment is visible.
func worldWallVisibleSpan(pos image.Point, dir, width int, lo, hi, sprite *int) bool {
	segment := func(base uintptr) (image.Point, image.Point) {
		off := base + uintptr(16*dir)
		a := image.Pt(int(memmap.Int32(0x587000, off)), int(memmap.Int32(0x587000, off+4)))
		b := image.Pt(int(memmap.Int32(0x587000, off+8)), int(memmap.Int32(0x587000, off+12)))
		return pos.Add(a), pos.Add(b)
	}
	clip := func(a, b image.Point, left, right *int) bool { return GetClient().Sub4C42A0(a, b, left, right) != 0 }
	a, b := segment(85440)
	switch dir {
	case 7, 9:
		first := clip(a, b, lo, hi)
		if !first {
			*lo = b.X
		}
		a, b = segment(85504)
		if !clip(a, b, lo, hi) {
			if !first {
				return false
			}
			if *hi > a.X {
				*hi = a.X
			}
		}
	case 8, 10:
		firstLo, firstHi := a.X, b.X
		first := clip(a, b, &firstLo, &firstHi) && firstHi-firstLo >= 3
		a, b = segment(85504)
		*lo, *hi = a.X, b.X
		second := clip(a, b, lo, hi) && *hi-*lo >= 3
		if first {
			if second {
				*lo = min(*lo, firstLo)
				if *lo <= a.X {
					*lo = 0
				}
				*hi = max(*hi, firstHi)
				if *hi >= b.X {
					*hi = width
				}
			} else {
				*lo, *hi = firstLo, firstHi
				if dir != 8 {
					*sprite = 1
					if *lo == a.X {
						*lo = 0
					}
				} else {
					*sprite = 0
					if firstHi == b.X {
						*hi = width
					}
				}
			}
		} else {
			if !second {
				return false
			}
			*sprite = 13
			if dir != 8 {
				*sprite = 14
			}
			if *hi == b.X {
				*hi = width
			}
			if *lo == a.X {
				*lo = 0
			}
		}
	default:
		if !clip(a, b, lo, hi) {
			return false
		}
	}
	return true
}
