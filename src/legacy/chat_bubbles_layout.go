package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func chatBubbleLayout(v *noxrender.Viewport) {
	r := GetClient().R2()
	pad := int32(r.FontHeight(nil))
	ordinal := uint32(0)
	for b := *chatBubbleHead(); b != nil; {
		next := b.Next
		b.Ordinal = ordinal
		ordinal++
		code := int(b.Code & 0x7fff)
		if b.Code&0x8000 != 0 {
			b.Drawable = GetClient().Cli().Objs.ByNetCodeStatic(code)
		} else {
			b.Drawable = GetClient().Cli().Objs.ByNetCodeDynamic(code)
		}
		if dr := b.Drawable; dr != nil {
			b.Position = uint32(uint16(dr.PosVec.X)) | uint32(uint16(dr.PosVec.Y))<<16
		}
		// the screen position follows the zoom and UI scale of the current drawing pass
		sp := worldToDrawScreen(image.Pt(int(uint16(b.Position)), int(uint16(b.Position>>16))))
		b.X, b.Y = int32(sp.X), int32(sp.Y)
		size := r.GetStringSizeWrapped(nil, alloc.GoString16(&b.Text[0]), 128)
		b.Width, b.Height = int32(min(size.X, 128)), int32(size.Y)
		b.X += b.Width / -2
		b.Y += -64 - b.Height
		b.Visible, b.Arrow = 1, 1
		if noxflags.HasGame(2048) {
			x, y := uint32(b.X), uint32(b.Y-64)
			if x < uint32(v.Screen.Min.X) || x > uint32(v.Screen.Max.X) || y < uint32(v.Screen.Min.Y) || y > uint32(v.Screen.Max.Y) {
				b.Visible, b.Arrow = 0, 0
			}
		}
		if b.Visible != 0 {
			minX, maxX := int32(v.Screen.Min.X)+pad, int32(v.Screen.Max.X)-b.Width-pad
			if b.X < minX {
				b.X, b.Arrow = minX, 0
			} else if b.X > maxX {
				b.X, b.Arrow = maxX, 0
			}
			minY, maxY := int32(v.Screen.Min.Y)+2*pad+2, int32(v.Screen.Max.Y)-b.Height-pad
			if b.Y < minY {
				b.Y, b.Arrow = minY, 0
			} else if b.Y > maxY {
				b.Y, b.Arrow = maxY, 0
			}
			rect := [4]int32{b.X, b.Y, b.X + b.Width, b.Y + b.Height}
			chatBubblePlace(&rect, &b.Arrow)
			b.X, b.Y = rect[0], rect[1]
		}
		if GetServer().S().Frame() > b.Expiry {
			chatBubbleUnlink(b)
		}
		b = next
	}
	for b := *chatBubbleHead(); b != nil; b = b.Next {
		chatBubbleArrange(b)
	}
}
