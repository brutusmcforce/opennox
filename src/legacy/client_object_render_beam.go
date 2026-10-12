package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func objectRenderBeamColors() int {
	*effectMapped(1321536) = uint32(noxcolor.RGB5551Color(255, 200, 255).Color32())
	*effectMapped(1321796) = uint32(noxcolor.RGB5551Color(255, 0, 255).Color32())
	value := uint32(noxcolor.RGB5551Color(100, 40, 100).Color32())
	*effectMapped(1321532) = value
	return int(value)
}
func objectRenderBeamAppend(packet unsafe.Pointer) int {
	count := int32(dword_5d4594_1321800)
	if count < 32 {
		count++
		*effectMapped(1321532 + uintptr(count*8)) = *(*uint32)(unsafe.Add(packet, 1))
		*effectMapped(1321536 + uintptr(count*8)) = *(*uint32)(unsafe.Add(packet, 5))
		dword_5d4594_1321800 = uint32(count)
	}
	return int(count)
}
func objectRenderBeamReset() { dword_5d4594_1321800 = 0 }
func objectRenderBeamDraw(vp *noxrender.Viewport) int {
	if *memmap.PtrPtr(0x852978, 8) == nil {
		return 0
	}
	if dword_5d4594_1321800 == 0 {
		return 0
	}
	sight := &GetClient().Cli().Sight
	var i int32
	for {
		coords := (*[4]uint16)(memmap.PtrOff(0x5D4594, 1321540+uintptr(i*8)))
		from := vp.ToScreenPos(image.Pt(int(coords[0]), int(coords[1])))
		to := from.Add(image.Pt(int(coords[2])-int(coords[0]), int(coords[3])-int(coords[1])))
		visible := sight.Sub_4992B0(from.X, from.Y)
		if from.X <= 0 || from.X >= vp.Size.X-1 || from.Y <= 0 || from.Y >= vp.Size.Y-1 {
			visible = 0
		}
		n := sight.Sub_498C20(&from, &to, int32(uintptr(vp.C())))
		if n != 0 {
			previous := from
			for j := 0; j < int(n); j++ {
				next := sight.Sub_499290(j)
				if visible != 0 {
					objectRenderBeamLine(previous, next)
				}
				previous = next
				visible = 1 - visible
			}
			if visible != 0 {
				objectRenderBeamLine(previous, to)
			}
		} else if visible != 0 {
			objectRenderBeamLine(from, to)
		}
		i++
		if i >= int32(dword_5d4594_1321800) {
			return int(i)
		}
	}
}
func objectRenderBeamLine(a, b image.Point) int {
	r := GetClient().R2()
	delta := b.Sub(a)
	relative := func(p image.Point) bool {
		r.AddPoint(p)
		r.AddPointRel(delta)
		return r.DrawLineFromPoints(r.Data().Color2())
	}
	effectColor(*effectMapped(1321532))
	r.Data().SetAlphaEnabled(true)
	relative(a)
	r.Data().SetAlphaEnabled(false)
	a.Y -= 22
	b.Y -= 22
	effectColor(*effectMapped(1321536))
	relative(a)
	effectColor(*effectMapped(1321796))
	dx, dy := delta.X, delta.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	off := image.Pt(0, 1)
	if dx <= dy {
		off = image.Pt(1, 0)
	}
	effectLine(a.Sub(off), b.Sub(off))
	r.AddPoint(a.Add(off))
	r.AddPoint(b.Add(off))
	return bool2int(r.DrawLineFromPoints(r.Data().Color2()))
}
