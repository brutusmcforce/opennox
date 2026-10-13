package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func objectDoorDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	data := (*spriteAnimationData)(dr.DrawData)
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames, int(dr.Field_74_4)))
	door := layout.As[client.DrawableUnionDoor](unsafe.Pointer(&dr.Union))
	if !noxflags.HasGame(noxflags.GameFlag(4096)) || door.Field_108_0 != 1 {
		return 1
	}
	if door.Field_109 == nil {
		door.Field_110 = noxrender.ImageHandle(Nox_xxx_gLoadImg("DoorLockSilverSW").C())
		door.Field_109 = noxrender.ImageHandle(Nox_xxx_gLoadImg("DoorLockSilverSE").C())
		door.Field_111 = noxrender.ImageHandle(Nox_xxx_gLoadImg("DoorLockGoldSW").C())
		door.Field_112 = noxrender.ImageHandle(Nox_xxx_gLoadImg("DoorLockGoldSE").C())
	}
	p := vp.ToScreenPos(dr.PosVec).Sub(image.Pt(64, 79))
	light := GetClient().Sub469920(dr.PosVec)
	r := GetClient().R2()
	r.Data().SetMultiply14(1)
	r.SetColorMultAndIntensityRGB(byte(light[0]), byte(light[1]), byte(light[2]))
	silver := door.Field_108_1 == 1
	slot := 436
	switch dr.Field_74_4 {
	case 0:
		p = p.Add(image.Pt(-15, -20))
		slot = 440
		if !silver {
			slot = 444
		}
	case 8:
		p = p.Add(image.Pt(15, -20))
		if !silver {
			slot = 448
		}
	case 16:
		p = p.Add(image.Pt(8, 2))
		slot = 440
		if !silver {
			slot = 444
		}
	default:
		p = p.Add(image.Pt(-8, 2))
		if !silver {
			slot = 448
		}
	}
	images := [4]noxrender.ImageHandle{door.Field_109, door.Field_110, door.Field_111, door.Field_112} // slots 436-448
	r.DrawImageAt(r.GetBag().AsImage(images[(slot-436)/4]), p)
	return 1
}
func objectArrowDraw(vp *noxrender.Viewport, dr *client.Drawable, weak bool) int {
	off := uintptr(1313720)
	name := "ArrowTailLink"
	if weak {
		off = 1313724
		name = "WeakArrowTailLink"
	}
	typ := effectMapped(off)
	if *typ == 0 {
		*typ = effectType(name)
	}
	from := image.Pt(int(dr.Field_81), int(dr.Field_82))
	delta := dr.PosVec.Sub(from)
	if nox_client_drawIsTick && delta.X*delta.X+delta.Y*delta.Y > 200 {
		tail := effectSpawn(int(*typ), from)
		if tail == nil {
			return spriteSlaveDraw(vp, dr)
		}
		u := (*client.DrawableUnionTail)(unsafe.Pointer(&tail.Union))
		u.X, u.Y = uint32(dr.PosVec.X), uint32(dr.PosVec.Y)
		effectLink(tail)
		dr.Field_81, dr.Field_82 = uint32(dr.PosVec.X), uint32(dr.PosVec.Y)
		GetClient().Cli().Objs.TransparentDecay(tail, int(gameFPS()/3))
	}
	return spriteSlaveDraw(vp, dr)
}
func objectArrowTailDraw(vp *noxrender.Viewport, dr *client.Drawable, weak bool) int {
	remaining := dr.Deadline - gameFrame()
	if remaining == 0 {
		return 1
	}
	// The C comparison and shift/division are unsigned because Deadline is uint32.
	index := int((remaining << 6) / (gameFPS() / 3))
	if index >= 64 {
		index = 63
	}
	off := uintptr(1313012)
	if weak {
		off = 1313268
	}
	effectColor(*effectMapped(off + uintptr(index*4)))
	d := GetClient().R2().Data()
	d.SetAlphaEnabled(true)
	d.SetAlpha(128)
	from := effectScreen(vp, dr).Sub(image.Pt(0, 4))
	u := (*client.DrawableUnionTail)(unsafe.Pointer(&dr.Union))
	to := image.Pt(int(u.X), int(u.Y))
	to = vp.ToScreenPos(to)
	to.Y -= int(int16(dr.ZVal)) + int(int16(dr.ZVal2)) + 4
	effectLine(from, to)
	d.SetAlphaEnabled(false)
	return 1
}
func objectShapeEdges(vp *noxrender.Viewport, dr *client.Drawable) [4]image.Point {
	base := vp.ToScreenPos(dr.PosVec)
	var p [4]image.Point
	box := &dr.Shape.Box
	for i, c := range [4][2]float32{
		{box.LeftTop, box.LeftBottom},     // 64
		{box.RightBottom2, box.RightTop2}, // 88
		{box.LeftBottom2, box.LeftTop2},   // 72
		{box.RightTop, box.RightBottom},   // 80
	} {
		p[i] = base.Add(image.Pt(effectFloatInt(c[0]), effectFloatInt(c[1])))
	}
	return p
}
func objectDrawEdges(p [4]image.Point, dx int) {
	shift := image.Pt(dx, 0)
	effectLine(p[0].Add(shift), p[3].Add(shift))
	effectLine(p[1].Add(shift), p[3].Add(shift))
	effectLine(p[0].Add(shift), p[2].Add(shift))
	effectLine(p[1].Add(shift), p[2].Add(shift))
}
func objectPressureDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	colors := &(*client.DrawableUnionPressurePlate)(unsafe.Pointer(&dr.Union)).Colors
	if colors[0]|colors[1]|colors[2]|colors[3]|colors[4]|colors[5] == 0 {
		return 1
	}
	p := objectShapeEdges(vp, dr)
	for i := 0; i < 2; i++ {
		nox_set_color_rgb_434430(int(colors[3*i]), int(colors[3*i+1]), int(colors[3*i+2]))
		objectDrawEdges(p, i)
	}
	return 1
}
func objectTriggerDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	effectColor(uint32(nox_color_black_2650656))
	d := GetClient().R2().Data()
	d.SetAlphaEnabled(true)
	p := objectShapeEdges(vp, dr)
	p[2], p[3] = p[3], p[2]
	objectDrawEdges(p, 0)
	d.SetAlphaEnabled(false)
	return 1
}
func objectPowderDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	p := vp.ToScreenPos(dr.PosVec)
	effectColor(memmap.Uint32(0x85B3FC, 956))
	r := GetClient().R2()
	cl := r.Data().Color2()
	r.DrawRectFilledOpaque(p.X-1, p.Y-1, 3, 3, cl)
	for _, delta := range []image.Point{{-5, 0}, {0, 7}, {8, -2}} {
		q := p.Add(delta)
		r.DrawRectFilledOpaque(q.X, q.Y, 1, 1, cl)
	}
	return 1
}
