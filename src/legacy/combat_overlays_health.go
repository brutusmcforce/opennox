package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"strconv"
	"unsafe"
)

type combatHealth struct {
	Code           uint32
	Amount         int16
	Reserved       uint16
	Stamp          uint32
	Next, Previous *combatHealth
}

var combatHealthPool unsafe.Pointer
var combatHealthHead *combatHealth
var combatHealthFont unsafe.Pointer
var _ = [1]struct{}{}[20-unsafe.Sizeof(combatHealth{})]

func combatHealthClass() alloc.ClassT[combatHealth] {
	return alloc.AsClassT[combatHealth](combatHealthPool)
}
func combatHealthInit() bool {
	c := alloc.NewClassT("HealthChange", combatHealth{}, 32)
	combatHealthPool = c.UPtr()
	if combatHealthPool == nil {
		return false
	}
	combatHealthFont = GetClient().R2().GetFonts().FontPtrByName("numbers")
	return true
}
func combatHealthClear() { combatHealthClass().FreeAllObjects(); combatHealthHead = nil }
func combatHealthDestroy() {
	combatHealthClass().Free()
	combatHealthPool = nil
	combatHealthHead = nil
	combatHealthFont = nil
}
func combatHealthAdd(code uint32, amount int16) *combatHealth {
	p := combatHealthClass().NewObject()
	if p == nil {
		return nil
	}
	old := combatHealthHead
	p.Code, p.Amount, p.Stamp, p.Next = code, amount, uint32(GetClient().GetInputSeq()), old
	if old != nil {
		old.Previous = p
	}
	combatHealthHead = p
	return old
}
func combatHealthRemove(p *combatHealth) {
	if p.Previous != nil {
		p.Previous.Next = p.Next
	} else {
		combatHealthHead = p.Next
	}
	if p.Next != nil {
		p.Next.Previous = p.Previous
	}
	combatHealthClass().FreeObjectFirst(p)
}
func combatHealthDraw(v *noxrender.Viewport, dr *client.Drawable) {
	seq := uint32(GetClient().GetInputSeq())
	color := noxcolor.RGBA5551(nox_color_yellow_2589772)
	if dr.C() == *memmap.PtrPtr(0x852978, 8) {
		color = noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 940))
	}
	r := GetClient().R2()
	face := r.GetFonts().AsFont(combatHealthFont)
	for p := combatHealthHead; p != nil; {
		next := p.Next
		if seq-p.Stamp > 30 {
			combatHealthRemove(p)
		} else if p.Code == dr.NetCode32 {
			amount := int(p.Amount)
			if amount < 0 {
				amount = -amount
			}
			text := strconv.Itoa(amount)
			pos := v.ToScreenPos(dr.PosVec)
			pos.Y += int(2*(p.Stamp-seq)) - int(int16(dr.ZVal)) - int(int64(dr.ZSizeMax))
			pos.X -= r.GetStringSizeWrapped(face, text, 0).X / 2
			r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_black_2650656))
			for _, off := range []image.Point{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}} {
				r.DrawString(face, text, pos.Add(off))
			}
			value := color
			if p.Amount > 0 {
				value = noxcolor.RGBA5551(dword_8531A0_2572)
			}
			r.Data().SetTextColor(value)
			r.DrawString(face, text, pos)
		}
		p = next
	}
}
