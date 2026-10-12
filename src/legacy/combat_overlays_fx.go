package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

type combatFX struct {
	Kind, Field4                               uint32
	History                                    [6]image.Point
	Count                                      byte
	Reserved                                   [3]byte
	Owner                                      *client.Drawable
	Next, Previous, GlobalNext, GlobalPrevious *combatFX
}

var _ = [1]struct{}{}[80-unsafe.Sizeof(combatFX{})]
var _ = [1]struct{}{}[64-unsafe.Offsetof(combatFX{}.Next)]

func combatFXHead() **combatFX { return &combatFXList }

// ResetCombatFX forgets all combat effects (their memory is freed with the DrawableFX class).
func ResetCombatFX() { combatFXList = nil }
func combatFXAttach(p *combatFX, dr *client.Drawable) {
	if p == nil || dr == nil {
		return
	}
	p.Owner = dr
	p.GlobalNext = *combatFXHead()
	p.GlobalPrevious = nil
	if p.GlobalNext != nil {
		p.GlobalNext.GlobalPrevious = p
	}
	*combatFXHead() = p
	p.Next = (*combatFX)(dr.Field_114.C())
	p.Previous = nil
	if p.Next != nil {
		p.Next.Previous = p
	}
	dr.Field_114 = (*client.DrawableFX)(unsafe.Pointer(p))
}
func combatFXDetach(p *combatFX) {
	if p == nil {
		return
	}
	if p.GlobalNext != nil {
		p.GlobalNext.GlobalPrevious = p.GlobalPrevious
	}
	if p.GlobalPrevious != nil {
		p.GlobalPrevious.GlobalNext = p.GlobalNext
	} else {
		*combatFXHead() = p.GlobalNext
	}
	if p.Next != nil {
		p.Next.Previous = p.Previous
	}
	if p.Previous != nil {
		p.Previous.Next = p.Next
	} else {
		p.Owner.Field_114 = (*client.DrawableFX)(unsafe.Pointer(p.Next))
	}
}
func combatFXDraw(v *noxrender.Viewport, dr *client.Drawable) {
	for p := (*combatFX)(dr.Field_114.C()); p != nil; {
		next := p.Next
		switch p.Kind {
		case 1:
			combatFXTrail(v, dr, p, false)
		case 2:
			combatFXTrail(v, dr, p, true)
		}
		p = next
	}
}
func combatFXTrail(v *noxrender.Viewport, dr *client.Drawable, p *combatFX, line bool) {
	// Preserve the original axis-equality stop test and its extra history slot.
	stop := true
	for i := 0; i < int(p.Count); i++ {
		a, b := p.History[i], p.History[i+1]
		if a.X != b.X && a.Y != b.Y {
			stop = false
			break
		}
	}
	if stop && dr.PosVec == image.Pt(int(dr.Field_8), int(dr.Field_9)) {
		p.Count = 0
		return
	}
	r := GetClient().R2()
	pos := dr.PosVec
	alpha := byte(255)
	if line {
		off := uintptr(dr.AnimFrameSlave) * 64
		dir := image.Pt(int(int32(float32(float64(memmap.Float32(0x587000, 194136+off))*-12))), int(int32(float32(float64(memmap.Float32(0x587000, 194140+off))*-12))))
		offset := dir.Sub(image.Pt(0, int(int16(dr.ZVal))+int(int16(dr.ZVal2))+10))
		last := v.ToScreenPos(pos).Add(offset)
		for i := 0; i < int(p.Count); i++ {
			alpha -= 63
			r.Data().SetAlphaEnabled(true)
			r.Data().SetAlpha(alpha)
			cur := v.ToScreenPos(p.History[i]).Add(offset)
			color := noxcolor.RGB5551Color(200, 200, 200)
			r.Data().SetColor2(color)
			r.AddPoint(last)
			r.AddPoint(cur)
			r.DrawLineFromPoints(r.Data().Color2())
			last = cur
		}
	} else {
		for i := 0; i < int(p.Count); i++ {
			alpha -= 42
			r.Data().SetAlphaEnabled(true)
			r.Data().SetAlpha(alpha)
			dr.PosVec = p.History[i]
			dr.CallDraw(v)
		}
		dr.PosVec = pos
	}
	if nox_client_drawIsTick {
		// the position history is only updated on game ticks
		copy(p.History[1:int(p.Count)+1], p.History[:int(p.Count)])
		p.History[0] = pos
		if p.Count != 5 {
			p.Count++
		}
	}
	r.Data().SetAlphaEnabled(false)
}
