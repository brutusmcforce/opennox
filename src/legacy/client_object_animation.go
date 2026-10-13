package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
)

func objectGeneratorCountdown(dr *client.Drawable) unsafe.Pointer {
	data := (*spriteConditionalData)(dr.DrawData)
	(*client.DrawableUnionGenerator)(unsafe.Pointer(&dr.Union)).Remaining = uint32(data.Count[3]) * (uint32(data.Delay[3]) + 1)
	return dr.DrawData
}
func objectGeneratorDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	data := (*spriteConditionalData)(dr.DrawData)
	flags := dr.Flags70Val
	state := 0
	if flags&0x100 != 0 {
		state = 1
	} else if flags&0x200 != 0 {
		state = 2
	} else if flags&0xc00 != 0 {
		state = 3
	}
	count, period := int(data.Count[state]), uint32(data.Delay[state])+1
	var index int
	switch data.Kind[state] {
	case client.AnimOneShot:
		// A one-shot generator is drawn during its destruction countdown or in
		// completed state. Those lifecycle states replace this provisional index.
		index = int(uintptr(unsafe.Pointer(&data.Delay[state])))
	case client.AnimLoop:
		index = int(int32((gameFrame() + dr.NetCode32) / period))
		if index >= count {
			index %= count
		}
	case client.AnimRandom:
		index = effectRand(0, count-1)
	case client.AnimSlave:
		index = int(int32(dr.AnimFrameSlave))
	default:
		return 0
	}
	if flags&0x800 != 0 {
		index = count - 1
		*(*uint32)(unsafe.Pointer(&dr.ObjClass)) &= 0xfff7ffff
		*(*uint32)(unsafe.Pointer(&dr.ObjFlags)) &= 0xdfffffff
	}
	u := (*client.DrawableUnionGenerator)(unsafe.Pointer(&dr.Union))
	if remaining := u.Remaining; remaining != 0 {
		index = int(int32((period*uint32(count) - remaining) / period))
		if index >= count {
			index = count - 1
		}
		if index < 0 {
			index = 0
		}
		u.Remaining = remaining - 1
		if remaining == 1 {
			dr.Flags70Val = dr.Flags70Val&^0x400 | 0x800
		}
	}
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames[state], index))
	if dr.Flags70Val&0xc00 == 0 {
		index = int(int32((gameFrame() + dr.NetCode32) / (uint32(data.Delay[4]) + 1)))
		if index >= int(data.Count[4]) {
			index %= int(data.Count[4])
		}
		Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames[4], index))
	}
	if dr.Flags70Val&0x800 != 0 {
		*(*uint32)(unsafe.Pointer(&dr.ObjFlags)) |= 1
	}
	return 1
}
func objectGlyphDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	r := GetClient().R2()
	d := r.Data()
	alpha := byte(255)
	local := LocalPlayerDrawable
	if noxflags.HasGame(noxflags.GameFlag(2)) && local != nil && *(*uint32)(unsafe.Pointer(&dr.ObjFlags))&0x40000000 == 0 {
		player := local
		if player.HasEnchant(21) {
			d.SetColorize17(1)
			cl := noxcolor.RGBA5551(dword_8531A0_2572).ColorNRGBA()
			r.SetColorMultAndIntensityRGB(cl.R, cl.G, cl.B)
		} else {
			delta := dr.PosVec.Sub(player.PosVec)
			dist := delta.X*delta.X + delta.Y*delta.Y
			if dist >= 22500 {
				return 1
			}
			alpha = byte(200 - 200*dist/22500)
		}
	}
	d.SetAlphaEnabled(true)
	d.SetAlpha(alpha)
	ret := spriteAnimateDraw(vp, dr)
	d.SetAlphaEnabled(false)
	d.SetColorize17(0)
	return ret
}
func objectShieldDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	code := (*client.DrawableUnionShield)(unsafe.Pointer(&dr.Union)).Owner
	var owner *client.Drawable
	if code&0x8000 != 0 {
		owner = GetClient().Cli().Objs.ByNetCodeStatic(int(code & 0x7fff))
	} else {
		owner = GetClient().Cli().Objs.ByNetCodeDynamic(int(code & 0x7fff))
	}
	if owner == nil {
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	effectMove(dr, owner.PosVec.X, owner.PosVec.Y+3)
	return spriteAnimateDraw(vp, dr)
}
func objectSummonDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	pos, frame := dr.PosVec, dr.AnimFrameSlave
	data := (*spriteAnimationData)(dr.DrawData)
	if drawableSummonSpark == 0 {
		drawableSummonSpark = uint32(effectType("BlueSpark"))
	}
	age := gameFrame() - dr.AnimStart
	summon := layout.As[client.DrawableUnionSummon](unsafe.Pointer(&dr.Union)) // see drawableSummonStart
	duration, child := uint32(summon.Duration), summon.Child
	if age >= duration {
		GetClient().Nox_xxx_spriteDelete_45A4B0(child)
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	if age >= duration-1 {
		effectCreatePointSparks(int(drawableSummonSpark), 50, 1000, 30, pos.X, pos.Y)
	}
	spriteAnimateDraw(vp, dr)
	var n uint32
	for off := uintptr(192092); off < 194140; off += 80 {
		if n >= uint32(data.Count) {
			n = 0
		}
		index := (n + gameFrame() + dr.NetCode32) / (uint32(data.Delay) + 1)
		if index >= uint32(data.Count) {
			index %= uint32(data.Count)
		}
		dr.AnimFrameSlave = index
		data.Kind = client.AnimSlave
		dr.PosVec = image.Pt(pos.X+2*int(memmap.Int32(0x587000, off-4)), pos.Y+2*int(memmap.Int32(0x587000, off)))
		if dr.PosVec.X >= 0 && dr.PosVec.X < 5888 && dr.PosVec.Y >= 0 && dr.PosVec.Y < 5888 {
			spriteAnimateDraw(vp, dr)
		}
		n++
	}
	dr.PosVec, dr.AnimFrameSlave = pos, frame
	data.Kind = client.AnimLoop
	d := GetClient().R2().Data()
	d.SetAlphaEnabled(true)
	d.SetAlpha(byte(int64((float64(gameFrame()) - float64(dr.AnimStart)) / float64(duration) * 255)))
	child.CallDraw(vp)
	d.SetAlphaEnabled(false)
	return 1
}
