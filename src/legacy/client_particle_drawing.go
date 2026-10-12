package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

func particleRect(pos image.Point, size int) {
	r := GetClient().R2()
	r.DrawRectFilledOpaque(pos.X, pos.Y, size, size, r.Data().Color2())
}
func particleMagicDraw(vp *noxrender.Viewport, dr *client.Drawable, missile bool) int {
	if dword_5d4594_1313804 == 0 {
		*effectMapped(1313808) = particleRGB(0, 200, 255)
		*effectMapped(1313812) = particleRGB(255, 255, 50)
		dword_5d4594_1313804 = 1
	}
	pos := effectScreen(vp, dr)
	if !effectInside(vp, pos, 10) {
		return 1
	}
	size := effectRand(1, 4)
	color := *effectMapped(1313808)
	if missile {
		color = *effectMapped(1313812)
	}
	effectGlow(pos, color, 2*size+1, (size>>1)+3)
	effectColor(uint32(nox_color_white_2523948))
	particleRect(pos.Sub(image.Pt(size>>1, size>>1)), size)
	light := unsafe.Add(dr.C(), 136)
	if missile {
		particleLightColor(light, 255, 180, 50)
	} else {
		particleLightColor(light, 200, 200, 255)
	}
	intensity := float32(GetServer().S().Rand.Other.Float(0, 100))
	particleLightIntensity(light, intensity, true)
	return 1
}
func particleTailDraw(vp *noxrender.Viewport, dr *client.Drawable, missile bool) int {
	from := effectScreen(vp, dr)
	to := vp.ToScreenPos(image.Pt(int(*effectWord(dr, 432)), int(*effectWord(dr, 436))))
	to.Y -= int(int16(dr.ZVal)) + int(int16(dr.ZVal2))
	remaining := int32(dr.Deadline - gameFrame())
	if remaining <= 0 {
		return 1
	}
	period := gameFPS()
	off := uintptr(1312500)
	var index int
	if missile {
		period /= 3
		off = 1312756
		index = int(remaining<<6) / int(period)
	} else {
		index = int(uint32(remaining<<6) / period)
	}
	if index >= 64 {
		index = 63
	}
	color := *effectMapped(off + uintptr(4*index))
	light := unsafe.Add(dr.C(), 136)
	if missile {
		particleLightColor(light, 255, 128, 50)
	} else {
		particleLightColor(light, 128, 128, 255)
	}
	intensity := float32(float64(remaining) * 20 / float64(int32(period)))
	particleLightIntensity(light, intensity, true)
	effectColor(color)
	data := GetClient().R2().Data()
	data.SetAlphaEnabled(true)
	data.SetAlpha(128)
	effectLine(from, to)
	data.SetAlphaEnabled(false)
	return 1
}
func particleBubbleDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	size, phase := effectByte(dr, 440), effectByte(dr, 441)
	timer, period := effectByte(dr, 442), effectByte(dr, 443)
	if *phase != 3 && dr.Deadline != 0 && dr.Deadline <= gameFrame() {
		*phase = 3
		*timer = 4
		*period = 4
		GetClient().Cli().Objs.TransparentDecay(dr, int(gameFPS()))
	}
	if *phase == 3 && *size == 0 {
		effectDelete(dr)
		return 0
	}
	pos := vp.ToScreenPos(dr.PosVec)
	pos.Y -= int(int16(dr.ZVal))
	effectGlow(pos, *effectWord(dr, 432), int(*size), int(*size)+3)
	effectColor(*effectWord(dr, 436))
	effectPoint(pos, int(*size)>>1)
	if nox_client_drawIsTick && gameFrame()&3 != 0 {
		dr.ZVal += uint16(int8(*effectByte(dr, 446)))
	}
	if *timer != 0 && nox_client_drawIsTick {
		*timer--
		if *timer == 0 {
			switch *phase {
			case 1:
				*size++
				if *size >= 5 {
					*phase = 2
				}
			case 2:
				*size--
				if *size == 0 {
					*phase = 1
				}
			default:
				*size--
				if *size == 0 {
					effectDelete(dr)
					return 0
				}
			}
			*timer = *period
		}
	}
	flip := effectByte(dr, 445)
	if *flip != 0 {
		*flip--
		if *flip == 0 {
			*effectByte(dr, 446) = -*effectByte(dr, 446)
			*flip = *effectByte(dr, 444)
		}
	}
	if int16(dr.ZVal) < 0 {
		effectDelete(dr)
		return 0
	}
	return 1
}
func particleBlueRain(vp *noxrender.Viewport, dr *client.Drawable) int {
	if noxflags.HasGame(noxflags.GameFlag(0x200000)) || !nox_client_drawIsTick {
		return 1
	}
	typ := updateType(1313716, "BlueRainSpark")
	for i := 0; i < 2; i++ {
		x := dr.PosVec.X + effectRand(-10, 10)
		y := dr.PosVec.Y + effectRand(-10, 10)
		if child := effectSpawn(int(typ), image.Pt(x, y)); child != nil {
			*effectWord(child, 432) = uint32(x) << 12
			*effectWord(child, 436) = uint32(y) << 12
			child.Field_74_4 = 0
			*effectWord(child, 440) = 0
			*effectWord(child, 448) = gameFrame() + uint32(effectRand(90, 120))
			*effectWord(child, 444) = gameFrame()
			child.ZVal2 = 0
			child.VelZ = -5
			child.ZVal = uint16(y - vp.World.Min.Y)
			effectLink(child)
		}
	}
	return 1
}
func particleFallingSparks(typ int, vp *noxrender.Viewport, dr *client.Drawable) *client.Drawable {
	if !nox_client_drawIsTick {
		return nil
	}
	var last *client.Drawable
	for i := 0; i < 2; i++ {
		x := dr.PosVec.X + effectRand(-15, 15)
		y := dr.PosVec.Y + effectRand(-15, 15)
		height := uint16(y - vp.World.Min.Y)
		velocity := int8(-effectRand(8, 12))
		last = effectCreateRainOrb(typ, image.Pt(x, y), dr.PosVec, height, velocity)
	}
	return last
}
func particleLevelUp(vp *noxrender.Viewport, dr *client.Drawable, blue bool) int {
	off, name := uintptr(1313708), "RainOrbWhite"
	if blue {
		off, name = 1313712, "RainOrbBlue"
	}
	particleFallingSparks(int(updateType(off, name)), vp, dr)
	return 1
}
func particleSpiderSpit(vp *noxrender.Viewport, dr *client.Drawable) int {
	z := int(int16(dr.ZVal))
	pos := vp.ToScreenPos(dr.PosVec)
	pos.Y -= z
	previous := image.Pt(int(dr.Field_8), int(dr.Field_9))
	old := vp.ToScreenPos(previous)
	old.Y -= z
	dx, dy := dr.PosVec.X-previous.X, dr.PosVec.Y-previous.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	effectColor(*memmap.PtrUint32(0x85B3FC, 956))
	if dx <= dy {
		effectLine(pos.Add(image.Pt(1, 0)), old.Sub(image.Pt(1, 0)))
		effectLine(pos.Sub(image.Pt(1, 0)), old.Sub(image.Pt(1, 0)))
	} else {
		effectLine(pos.Add(image.Pt(0, 1)), old.Add(image.Pt(0, 1)))
		effectLine(pos.Add(image.Pt(0, 1)), old.Sub(image.Pt(0, 1)))
	}
	particleRect(pos.Sub(image.Pt(1, 1)), 4)
	effectColor(uint32(nox_color_white_2523948))
	effectLine(pos, old)
	particleRect(pos, 2)
	return 1
}
func particleVortexDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	if *effectMapped(1313820) == 0 {
		dword_5d4594_1313816 = uint32(particleRGB(170, 170, 170))
		*effectMapped(1313820) = 1
	}
	center := image.Pt(int(*effectWord(dr, 440)), int(*effectWord(dr, 444)))
	radius := int(*effectByte(dr, 450))
	angle := int(*effectByte(dr, 448))
	world := updateRadial(center, angle, radius)
	pos := vp.ToScreenPos(world)
	pos.Y -= int(int16(dr.ZVal))
	if pos.X <= vp.Screen.Min.X || pos.X >= vp.Screen.Max.X || pos.Y <= vp.Screen.Min.Y || pos.Y >= vp.Screen.Max.Y {
		effectDelete(dr)
		return 0
	}
	if world.Y >= center.Y {
		effectGlow(pos, *effectWord(dr, 432), 3, 5)
		effectColor(*effectWord(dr, 436))
		effectPoint(pos, 3)
	} else {
		effectGlow(pos, uint32(dword_5d4594_1313816), 2, 4)
		effectColor(uint32(dword_5d4594_1313816))
		effectPoint(pos, 2)
	}
	speed := int(int8(*effectByte(dr, 449)))
	previous := angle - 2*speed
	if speed > 0 {
		if previous < 0 {
			previous += 256
		}
	} else if previous >= 256 {
		previous -= 256
	}
	// The source uses unsigned products/division for the trailing endpoint.
	oldWorld := image.Pt(int(uint32(center.X)+uint32(radius)**memmap.PtrUint32(0x587000, 192088+uintptr(8*previous))/16), int(uint32(center.Y)+uint32(radius)**memmap.PtrUint32(0x587000, 192092+uintptr(8*previous))/16))
	old := vp.ToScreenPos(oldWorld)
	old.Y -= int(int16(dr.ZVal))
	if oldWorld.Y >= center.Y {
		effectColor(*effectWord(dr, 436))
	} else {
		effectColor(uint32(dword_5d4594_1313816))
	}
	effectLine(pos, old)
	if nox_client_drawIsTick {
		*effectByte(dr, 448) += *effectByte(dr, 449)
		dr.ZVal += uint16(*effectByte(dr, 451))
	}
	shrink := effectFloatInt(float32(float64(int16(dr.ZVal)) * 0.0024999999 * 50))
	if 50-shrink <= 0 {
		effectDelete(dr)
		return 0
	}
	*effectByte(dr, 450) = byte(50 - shrink)
	return 1
}
