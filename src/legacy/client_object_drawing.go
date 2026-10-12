package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func objectDoorDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	data := (*spriteAnimationData)(dr.DrawData)
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames, int(dr.Field_74_4)))
	b := unsafe.Slice((*byte)(dr.C()), 512)
	if !noxflags.HasGame(noxflags.GameFlag(4096)) || b[432] != 1 {
		return 1
	}
	if *effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */) == 0 {
		*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, int(layoutClientDrawable.Off(440)) /* Union._ */))) = unsafe.Pointer(Nox_xxx_gLoadImg("DoorLockSilverSW").C())
		*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */))) = unsafe.Pointer(Nox_xxx_gLoadImg("DoorLockSilverSE").C())
		*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, int(layoutClientDrawable.Off(444)) /* Union._ */))) = unsafe.Pointer(Nox_xxx_gLoadImg("DoorLockGoldSW").C())
		*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, int(layoutClientDrawable.Off(448)) /* Union._ */))) = unsafe.Pointer(Nox_xxx_gLoadImg("DoorLockGoldSE").C())
	}
	p := vp.ToScreenPos(dr.PosVec).Sub(image.Pt(64, 79))
	light := GetClient().Sub469920(dr.PosVec)
	r := GetClient().R2()
	r.Data().SetMultiply14(1)
	r.SetColorMultAndIntensityRGB(byte(light[0]), byte(light[1]), byte(light[2]))
	silver := b[433] == 1
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
	r.DrawImageAt(r.GetBag().AsImage(noxrender.ImageHandle(*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, slot))))), p)
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
		*effectWord(tail, int(layoutClientDrawable.Off(432)) /* Union._ */), *effectWord(tail, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.X), uint32(dr.PosVec.Y)
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
	to := image.Pt(int(*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */)), int(*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */)))
	to = vp.ToScreenPos(to)
	to.Y -= int(int16(dr.ZVal)) + int(int16(dr.ZVal2)) + 4
	effectLine(from, to)
	d.SetAlphaEnabled(false)
	return 1
}
func objectShapeEdges(vp *noxrender.Viewport, dr *client.Drawable) [4]image.Point {
	base := vp.ToScreenPos(dr.PosVec)
	var p [4]image.Point
	for i, off := range []int{64, 88, 72, 80} {
		x := effectFloatInt(*(*float32)(unsafe.Add(dr.C(), off)))
		y := effectFloatInt(*(*float32)(unsafe.Add(dr.C(), off+4)))
		p[i] = base.Add(image.Pt(x, y))
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
	colors := unsafe.Slice((*byte)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)), 6)
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
