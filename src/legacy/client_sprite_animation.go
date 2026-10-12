package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"unsafe"
)

type spriteAnimationData struct {
	Size         uint32
	Frames       *noxrender.ImageHandle
	Count, Delay byte
	Reserved     uint16
	Kind         client.AnimKind
}
type spriteConditionalData struct {
	Size     uint32
	Frames   [5]*noxrender.ImageHandle
	Count    [5]byte
	Delay    [5]byte
	Reserved uint16
	Kind     [5]client.AnimKind
}

var _ = [1]struct{}{}[16-unsafe.Sizeof(spriteAnimationData{})]
var _ = [1]struct{}{}[56-unsafe.Sizeof(spriteConditionalData{})]

func spriteFrameImage(frames *noxrender.ImageHandle, index int) noxrender.ImageHandle {
	return *(*noxrender.ImageHandle)(unsafe.Add(unsafe.Pointer(frames), index*4))
}
func spriteAnimateDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	data := (*spriteAnimationData)(dr.DrawData)
	count := int(data.Count)
	period := uint32(data.Delay) + 1
	var index int
	switch data.Kind {
	case client.AnimOneShot, client.AnimOneShotRemove, client.AnimLoopAndFade:
		if data.Kind == client.AnimLoopAndFade {
			GetClient().R2().Data().SetAlphaEnabled(true)
		}
		index = int(int32((gameFrame() - dr.AnimStart) / period))
		end := count
		if data.Kind == client.AnimLoopAndFade {
			end *= 2
		}
		if index >= end {
			if data.Kind == client.AnimOneShot {
				index = count - 1
			} else {
				// Preserve the fade callback's early-delete render state.
				GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
				return 0
			}
		}
		if data.Kind == client.AnimLoopAndFade {
			GetClient().R2().Data().SetAlpha(byte(200 - 200*index/end))
			if index >= count {
				index %= count
			}
		}
	case client.AnimLoop:
		flags := *effectWord(dr, 120)
		class := *effectWord(dr, 112)
		if flags&0x1000000 != 0 || class&0x10000000 != 0 && noxflags.HasGame(noxflags.GameFlag(32)) {
			index = int(int32((gameFrame() + *effectWord(dr, 128)) / period))
			if index >= count {
				index %= count
			}
		} else if class&0x10000000 != 0 {
			return 1
		}
	case client.AnimRandom:
		index = effectRand(0, count-1)
	case client.AnimSlave:
		index = int(int32(dr.AnimFrameSlave))
	default:
		return 1
	}
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames, index))
	if data.Kind == client.AnimLoopAndFade {
		GetClient().R2().Data().SetAlphaEnabled(false)
	}
	return 1
}
func spriteConditionalDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	data := (*spriteConditionalData)(dr.DrawData)
	state := 1
	if *effectWord(dr, 120)&0x1000000 != 0 {
		state = 0
	}
	count := int(data.Count[state])
	var index int
	switch data.Kind[state] {
	case client.AnimLoop:
		index = int(int32((gameFrame() + *effectWord(dr, 128)) / (uint32(data.Delay[state]) + 1)))
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
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(data.Frames[state], index))
	return 1
}
func spriteStaticDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	if *effectWord(dr, 112)&0x40000 == 0 || *effectWord(dr, 120)&0x1000000 != 0 {
		img := *(*noxrender.ImageHandle)(unsafe.Add(dr.DrawData, 4))
		Nox_xxx_drawObject_4C4770_draw(vp, dr, img)
	}
	return 1
}
func spriteSlaveDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	frames := *(**noxrender.ImageHandle)(unsafe.Add(dr.DrawData, 4))
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(frames, int(int32(dr.AnimFrameSlave))))
	return 1
}
func spriteBoulderDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	oldX, oldY := effectWord(dr, 432), effectWord(dr, 436)
	index, bank := effectWord(dr, 440), effectWord(dr, 444)
	x, y := uint32(dr.PosVec.X), uint32(dr.PosVec.Y)
	if *oldX == 0 && *oldY == 0 {
		*oldX, *oldY = x, y
	}
	dx, dy := int32(x-*oldX), int32(y-*oldY)
	if dx*dx+dy*dy >= 100 {
		if dx <= 0 && dy > 0 || dx > 0 && dy <= 0 {
			*bank = 0
		} else {
			*bank = 16
		}
		if dy > 0 {
			*index++
			if *index >= 16 {
				*index = 0
			}
		} else if *index != 0 {
			*index--
		} else {
			*index = 15
		}
		*oldX, *oldY = x, y
	}
	frames := *(**noxrender.ImageHandle)(unsafe.Add(dr.DrawData, 4))
	Nox_xxx_drawObject_4C4770_draw(vp, dr, spriteFrameImage(frames, int(int32(*index+*bank))))
	return 1
}
