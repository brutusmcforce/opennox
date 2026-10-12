package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
)

func chatBubbleDraw(v *noxrender.Viewport) {
	r := GetClient().R2()
	height := r.FontHeight(nil)
	half := height / 2
	Sub_437260()
	chatBubbleLayout(v)
	line := func(x1, y1, x2, y2 int) {
		r.AddPoint(image.Pt(x1, y1))
		r.AddPoint(image.Pt(x2, y2))
		r.DrawLineFromPoints(r.Data().Color2())
	}
	for b := *chatBubbleHead(); b != nil; b = b.Next {
		if b.Visible == 0 {
			continue
		}
		color := noxcolor.RGBA5551(nox_color_white_2523948)
		var name *uint16
		if dr := b.Drawable; dr != nil && dr.ObjClass&4 != 0 {
			team := teamRuntimeObject(int(dr.NetCode32))
			if p := GetServer().S().Players.ByID(int(b.Code)); p != nil {
				name = &p.NameFinal[0]
			}
			if team != nil {
				if t := GetServer().S().Teams.ByID(team.ID); t != nil {
					color = noxcolor.ToRGBA5551Color(GetServer().S().Teams.GetTeamColor(t))
				}
			}
		}
		x, y := int(b.X), int(b.Y)
		for pass := 0; pass < 2; pass++ {
			edge := *memmap.PtrUint32(0x852978, 4)
			if pass != 0 {
				edge = *memmap.PtrUint32(0x85B3FC, 956)
				x--
				y--
			}
			r.Data().SetColor2(noxcolor.RGBA5551(edge))
			right, bottom := x+int(b.Width), y+int(b.Height)
			line(x, y-half, x-half, y)
			line(x, y-half, right, y-half)
			line(right+half, y, right, y-half)
			line(right+half, y, right+half, bottom)
			line(right, bottom+half, right+half, bottom)
			if b.Arrow != 0 {
				center := x + int(b.Width)/2
				line(right, bottom+half, center+half, bottom+half)
				line(center, bottom+half+height, center+half, bottom+half)
				line(center, bottom+half+height, center-half, bottom+half)
				line(x, bottom+half, center-half, bottom+half)
			} else {
				line(x, bottom+half, right, bottom+half)
			}
			line(x, bottom+half, x-half, bottom)
			line(x-half, y, x-half, bottom)
		}
		r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
		r.Data().SetColor(noxcolor.RGBA5551(*memmap.PtrUint32(0x852978, 4)))
		r.DrawStringWrappedHL(nil, alloc.GoString16(&b.Text[0]), image.Rect(x, y, x+128, y))
		if name != nil {
			r.Data().SetTextColor(color)
			r.DrawStringWrappedHL(nil, alloc.GoString16(name), image.Rect(x, y-height-1, x+128, y-height-1))
		}
	}
	r.SetRectFullScreen()
}
