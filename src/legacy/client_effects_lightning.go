package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"math"
)

func effectPrepareLightning() int {
	for _, v := range [][4]int{
		{1316464, 255, 255, 255}, {1316488, 128, 128, 255}, {1316424, 128, 128, 255},
		{1316428, 64, 64, 255}, {1316516, 200, 200, 255}, {1316512, 128, 128, 255},
		{1316496, 255, 255, 255}, {1316468, 255, 255, 0}, {1316460, 30, 160, 30},
		{1316444, 60, 140, 60}, {1316504, 40, 225, 40}, {1316480, 150, 220, 150},
	} {
		*effectMapped(uintptr(v[0])) = noxcolor.RGB5551Color(byte(v[1]), byte(v[2]), byte(v[3])).Color32()
	}
	for i, name := range []string{"BlueSpark", "YellowSpark", "GreenSpark"} {
		*effectMapped(1316520 + uintptr(i*4)) = effectType(name)
	}
	return int(*effectMapped(1316528))
}
func effectPackedPoint(v uint32) image.Point { return image.Pt(int(int16(v)), int(int16(v>>16))) }
func effectPackPoint(p image.Point) uint32   { return uint32(uint16(p.X)) | uint32(uint16(p.Y))<<16 }
func effectLightningStep(from, to uint32) int {
	dword_5d4594_1316492++
	depth := int(int32(dword_5d4594_1316492))
	a, b := effectPackedPoint(from), effectPackedPoint(to)
	if dword_5d4594_1316492 != dword_5d4594_1316448 {
		denominator, numerator := 1, 1
		for i := 1; i < depth; i++ {
			denominator *= int(*memmap.PtrUint32(0x587000, 178212)) + 9
			numerator *= 10
		}
		jitter := int(int32(dword_5d4594_1316476))
		mid := image.Pt(numerator*effectRand(-jitter, jitter)/denominator+((a.X+b.X)>>1), numerator*effectRand(-jitter, jitter)/denominator+((a.Y+b.Y)>>1))
		packed := effectPackPoint(mid)
		effectLightningStep(from, packed)
		effectLightningStep(packed, to)
	} else {
		glow := *effectMapped(1316508) != 0
		radius := 32
		color := uint32(dword_5d4594_1316472)
		size := 3
		if glow {
			radius = int(*memmap.PtrUint8(0x5D4594, 1316420)) + 48
			color = *effectMapped(1316440)
			size = 12
		}
		nox_draw_set54RGB32_434040(int(color))
		r := GetClient().R2()
		r.Data().SetField262(size)
		r.AddPoint(a)
		r.AddPoint(b)
		r.DrawParticles49ED80(radius)
		effectColor(uint32(dword_5d4594_1316472))
		effectLine(a, b)
		if glow {
			dx, dy := a.X-b.X, a.Y-b.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			for _, sign := range []int{1, -1} {
				for i := 1; i < int(*memmap.PtrUint8(0x5D4594, 1316420)>>1); i++ {
					r.ClearPoints()
					offset := image.Pt(0, sign*i)
					if dx <= dy {
						offset = image.Pt(sign*i, 0)
					}
					effectLine(b.Add(offset), a.Add(offset))
				}
			}
		}
	}
	dword_5d4594_1316492--
	return int(int32(dword_5d4594_1316492))
}
func effectLightningPasses(a, b image.Point, mode int, coords *[4]int16, unused, outer, inner int) int {
	dx, dy := a.X-b.X, a.Y-b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	length := int(int64(math.Sqrt(float64(dx*dx + dy*dy))))
	steps := uint32(nox_xxx_lightningSteps_587000_178216)
	if length >= 512 {
		dword_5d4594_1316476 = uint32(*memmap.PtrUint32(0x587000, 178204))
		dword_5d4594_1316448 = uint32(steps)
	} else {
		lo, hi := *memmap.PtrUint32(0x587000, 178208), *memmap.PtrUint32(0x587000, 178204)
		dword_5d4594_1316476 = uint32(lo + uint32(length)*(hi-lo)/512)
		var decrease uint32
		switch {
		case length < 64:
			decrease = 3
		case length < 128:
			decrease = 2
		case length < 256:
			decrease = 1
		}
		depth := steps - decrease
		if length < 256 && int32(depth) < 1 {
			depth = 1
		}
		dword_5d4594_1316448 = uint32(depth)
	}
	*effectMapped(1316532) = uint32(mode)
	if mode == 1 || mode == 3 {
		*memmap.PtrUint16(0x5D4594, 1316432) = uint16(coords[0])
		*memmap.PtrUint16(0x5D4594, 1316434) = uint16(coords[2])
		if dx <= dy {
			*effectMapped(1316500) = uint32(bool2int(a.Y >= b.Y) + 2)
		} else {
			*effectMapped(1316500) = uint32(bool2int(a.X >= b.X))
		}
	}
	from, to := effectPackPoint(a), effectPackPoint(b)
	if outer != 0 {
		dword_5d4594_1316492 = 1
		dword_5d4594_1316472 = dword_5d4594_1316456
		*effectMapped(1316508) = 0
		effectLightningStep(from, to)
		dword_5d4594_1316492 = 1
		dword_5d4594_1316472 = dword_5d4594_1316452
		effectLightningStep(from, to)
	}
	if inner == 0 {
		return 0
	}
	dword_5d4594_1316492 = 1
	dword_5d4594_1316472 = dword_5d4594_1316436
	*effectMapped(1316440) = uint32(dword_5d4594_1316484)
	*effectMapped(1316508) = 1
	return effectLightningStep(from, to)
}
func effectRayEndpoints(dr *client.Drawable, stripStatic bool) (a, b image.Point, ok bool) {
	if *effectByte(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) == 0 {
		return image.Pt(int(*effectShort(dr, int(layoutClientDrawable.Off(437)) /* Union._ */)), int(*effectShort(dr, int(layoutClientDrawable.Off(439)) /* Union._ */))), image.Pt(int(*effectShort(dr, int(layoutClientDrawable.Off(441)) /* Union._ */)), int(*effectShort(dr, int(layoutClientDrawable.Off(443)) /* Union._ */))), true
	}
	lookup := func(code uint32) *client.Drawable {
		if code&0x8000 != 0 {
			if stripStatic {
				code &= 0x7fff
			}
			return GetClient().Cli().Objs.ByNetCodeStatic(int(code))
		}
		if stripStatic {
			code &= 0x7fff
		}
		return GetClient().Cli().Objs.ByNetCodeDynamic(int(code))
	}
	from, to := lookup(*effectWord(dr, int(layoutClientDrawable.Off(437)) /* Union._ */)), lookup(*effectWord(dr, int(layoutClientDrawable.Off(441)) /* Union._ */))
	if from == nil || to == nil {
		return a, b, false
	}
	return from.PosVec, to.PosVec, true
}

// Lightning and chain lightning share the repaired endpoint construction.
func effectLightningDraw(vp *noxrender.Viewport, dr *client.Drawable, kind int) int {
	if kind == 2 && *effectByte(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) == 0 && *effectWord(dr, int(layoutClientDrawable.Off(433)) /* Union._ */) != 0 {
		*effectWord(dr, int(layoutClientDrawable.Off(433)) /* Union._ */)--
		if *effectWord(dr, int(layoutClientDrawable.Off(433)) /* Union._ */) == 0 {
			return effectDelete(dr)
		}
	}
	a, b, ok := effectRayEndpoints(dr, false)
	if !ok {
		return 1
	}
	sa, sb := vp.ToScreenPos(a), vp.ToScreenPos(b)
	sa.Y -= 20
	sb.Y -= 20
	particle := uintptr(1316520)
	outer := 1
	switch kind {
	case 1:
		*memmap.PtrUint8(0x5D4594, 1316420) = byte(2 * (int(int8(*effectByte(dr, int(layoutClientDrawable.Off(433)) /* Union._ */))) + 127))
		dword_5d4594_1316436 = uint32(*effectMapped(1316496))
		dword_5d4594_1316484 = uint32(*effectMapped(1316468))
		outer = 0
		particle = 1316524
	case 2:
		dword_5d4594_1316452 = uint32(*effectMapped(1316444))
		dword_5d4594_1316436 = uint32(*effectMapped(1316504))
		dword_5d4594_1316456 = uint32(*effectMapped(1316460))
		dword_5d4594_1316484 = uint32(*effectMapped(1316480))
		*memmap.PtrUint8(0x5D4594, 1316420) = 1
		particle = 1316528
	default:
		dword_5d4594_1316452 = uint32(*effectMapped(1316428))
		dword_5d4594_1316436 = uint32(*effectMapped(1316464))
		dword_5d4594_1316456 = uint32(*effectMapped(1316424))
		dword_5d4594_1316484 = uint32(*effectMapped(1316488))
		*memmap.PtrUint8(0x5D4594, 1316420) = 1
	}
	effectLightningPasses(sa, sb, 2, nil, outer, outer, 1)
	if !noxflags.HasGame(noxflags.GamePause) {
		effectLightningParticles(int(*effectMapped(particle)), a, b)
	}
	return 1
}
