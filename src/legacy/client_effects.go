package legacy

import (
	"fmt"
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/types"
)

// Effect records overlay the drawable union with several packed layouts. Keep
// byte/word access explicit where those layouts overlap (including unaligned
// packet coordinates), while using Drawable fields for shared object state.
func effectByte(dr *client.Drawable, off int) *byte    { return (*byte)(effectAt(dr, off, 1)) }
func effectShort(dr *client.Drawable, off int) *uint16 { return (*uint16)(effectAt(dr, off, 2)) }

// effectAt translates a 386 drawable offset of a scalar run of n bytes. Union bytes (432-451)
// stay union-relative, like the effect accesses through &dr.Union.
func effectAt(dr *client.Drawable, off, n int) unsafe.Pointer {
	const union, unionSize = 432, 20
	if off >= union && off+n <= union+unionSize {
		return unsafe.Add(unsafe.Pointer(&dr.Union), off-union)
	}
	p := layoutClientDrawable.Off(uintptr(off))
	if n > 1 && layoutClientDrawable.Off(uintptr(off+n-1)) != p+uintptr(n-1) {
		panic(fmt.Errorf("effectAt: drawable offset %d+%d is not contiguous", off, n))
	}
	return unsafe.Add(dr.C(), p)
}
func effectMapped(off uintptr) *uint32 { return memmap.PtrUint32(0x5D4594, off) }
func effectRand(lo, hi int) int        { return GetServer().S().Rand.Other.Int(lo, hi) }
func effectDelete(dr *client.Drawable) int {
	GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	return 0
}
func effectMove(dr *client.Drawable, x, y int) {
	GetClient().Cli().Nox_xxx_updateSpritePosition_49AA90(dr, x, y)
}
func effectSpawn(typ int, pos image.Point) *client.Drawable {
	return GetClient().Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
}
func effectLink(dr *client.Drawable) { GetClient().Cli().Objs.List34Add(dr) }
func effectType(name string) uint32  { return uint32(GetClient().Cli().Things.IndByID(name)) }
func effectDistance(x, y int) int    { return int(screenDistance(int32(x), int32(y))) }
func effectAngle(x, y float32) int {
	p := [2]float32{x, y}
	return int(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&p))))))
}
func effectFloatInt(v float32) int { return int(floatToInt32(v)) }
func effectGlow(p image.Point, color uint32, radius, size int) {
	GetClient().R2().DrawGlow(p, noxcolor.RGBA5551(color), radius, int(int8(size)))
}
func effectColor(color uint32) { GetClient().R2().Data().SetColor2(noxcolor.RGBA5551(color)) }
func effectPoint(p image.Point, radius int) {
	r := GetClient().R2()
	r.DrawPoint(p, radius, r.Data().Color2())
}
func effectLine(a, b image.Point) {
	r := GetClient().R2()
	r.AddPoint(a)
	r.AddPoint(b)
	r.DrawLineFromPoints(r.Data().Color2())
}
func effectInside(vp *noxrender.Viewport, p image.Point, radius int) bool {
	return p.X-radius >= vp.Screen.Min.X && p.Y-radius >= vp.Screen.Min.Y && p.X+radius < vp.Screen.Max.X && p.Y+radius < vp.Screen.Max.Y
}
func effectScreen(vp *noxrender.Viewport, dr *client.Drawable) image.Point {
	p := vp.ToScreenPos(dr.PosVec)
	p.Y -= int(int16(dr.ZVal)) + int(int16(dr.ZVal2))
	return p
}

// sub_4B6770/sub_4B6880 differ in their intensity sampling but share lifetime
// and strict clipping. Byte-width glow intensity follows the renderer C ABI.
func effectSparkDraw(vp *noxrender.Viewport, dr *client.Drawable, core, glow uint32, random bool) int {
	u := (*client.DrawableUnionSpark)(unsafe.Pointer(&dr.Union))
	lifetime := int(int32(u.Expires - u.Born))
	left := int(int32(u.Expires - gameFrame()))
	if left == lifetime {
		left--
	}
	if left <= 0 {
		return effectDelete(dr)
	}
	p := effectScreen(vp, dr)
	inside := effectInside(vp, p, 10)
	if !random {
		// The moving-spark C viewport is unsigned for X clipping and the
		// bottom edge; keep that behavior at a partially clipped left edge.
		inside = uint32(p.X-10) >= uint32(vp.Screen.Min.X) && p.Y-10 >= vp.Screen.Min.Y && uint32(p.X+10) < uint32(vp.Screen.Max.X) && uint32(p.Y+10) < uint32(vp.Screen.Max.Y)
	}
	if inside {
		if random {
			radius := left * effectRand(0, 4) / lifetime
			if radius != 0 {
				effectGlow(p, glow, 2*radius+1, radius+1)
				effectColor(core)
				effectPoint(p, radius)
			}
		} else {
			radius := 4 * left / lifetime
			effectGlow(p, glow, 2*radius+1, 5*left/lifetime)
			effectColor(core)
			effectPoint(p, radius)
		}
	}
	return 1
}
func effectSparkBounce(dr *client.Drawable) int16 {
	if !nox_client_drawIsTick {
		return int16(dr.ZVal)
	}
	v := int(dr.VelZ)
	dr.ZVal += uint16(int16(v))
	result := int16(dr.ZVal)
	if result >= 0 {
		if gameFrame()&1 != 0 {
			dr.VelZ = int8(v - 1)
		}
	} else {
		dr.ZVal = uint16(-result)
		result = int16(26209 * v)
		dr.VelZ = int8(-9 * v / 10)
		if dr.VelZ < 2 {
			dr.ZVal = 0
			dr.VelZ = 0
		}
	}
	return result
}
func effectMovingSpark(vp *noxrender.Viewport, dr *client.Drawable, core, glow uint32) int {
	if nox_client_drawIsTick {
		u := (*client.DrawableUnionSpark)(unsafe.Pointer(&dr.Union))
		speed := int32(u.Speed)
		angle := uintptr(dr.Field_74_4)
		u.X += uint32(speed * *memmap.PtrInt32(0x587000, 192088+8*angle))
		y := speed**memmap.PtrInt32(0x587000, 192092+8*angle) + int32(u.Y)
		u.Y = uint32(y)
		// C shifts X unsigned and Y signed.
		effectMove(dr, int(u.X>>12), int(y>>12))
	}
	effectSparkBounce(dr)
	return effectSparkDraw(vp, dr, core, glow, false)
}
func effectMagicSparkle(vp *noxrender.Viewport, dr *client.Drawable) int {
	if effectRand(0, 10) >= 5 {
		return effectSparkDraw(vp, dr, uint32(nox_color_white_2523948), uint32(dword_5d4594_1313540), true)
	}
	return effectSparkDraw(vp, dr, uint32(dword_5d4594_1313540), uint32(dword_5d4594_1313536), true)
}
func effectPixieDust(vp *noxrender.Viewport, dr *client.Drawable) int {
	if effectRand(0, 10) >= 5 {
		return effectSparkDraw(vp, dr, uint32(nox_color_white_2523948), uint32(dword_5d4594_1313564), true)
	}
	return effectSparkDraw(vp, dr, uint32(dword_5d4594_1313564), *effectMapped(1313560), true)
}
func effectPixie(vp *noxrender.Viewport, dr *client.Drawable) int {
	lower := effectRand(0, 100) < 50
	z := int16(dr.ZVal)
	if !nox_client_drawIsTick {
		// don't move between the ticks
	} else if lower {
		if z > 0 {
			dr.ZVal--
		}
	} else if z < 35 {
		dr.ZVal++
	}
	p := effectScreen(vp, dr)
	if effectInside(vp, p, 10) {
		effectGlow(p, *effectMapped(1313560), 10, 4)
		prev := vp.ToScreenPos(image.Pt(int(dr.Field_8), int(dr.Field_9)))
		prev.Y -= int(int16(dr.ZVal))
		dx, dy := p.X-prev.X, p.Y-prev.Y
		if n := dx*dx + dy*dy; n > 400 {
			distance := int(math.Sqrt(float64(n)))
			prev = image.Pt(p.X-20*dx/distance, p.Y-20*dy/distance)
		}
		effectColor(uint32(dword_5d4594_1313564))
		effectLine(p, prev)
	}
	return 1
}
func effectOrb(vp *noxrender.Viewport, dr *client.Drawable, moving bool) int {
	if *effectMapped(1313660) == 0 {
		for i, name := range []string{"DrainManaOrb", "HealOrb", "CharmOrb", "WhiteOrb", "ManaBombOrb", "WhiteMoveOrb", "BlueMoveOrb"} {
			*effectMapped(1313660 + uintptr(4*i)) = effectType(name)
		}
	}
	var core, glow uint32
	switch dr.TypeIDVal {
	case *effectMapped(1313660), *effectMapped(1313684):
		core, glow = uint32(dword_5d4594_1313540), uint32(dword_5d4594_1313536)
	case *effectMapped(1313668):
		core, glow = *effectMapped(1313584), *effectMapped(1313580)
	case *effectMapped(1313672), *effectMapped(1313676), *effectMapped(1313680):
		core, glow = *effectMapped(1313592), *effectMapped(1313588)
	default:
		core, glow = uint32(dword_5d4594_1313532), *effectMapped(1313528)
	}
	u := (*client.DrawableUnionOrb)(unsafe.Pointer(&dr.Union))
	if moving && nox_client_drawIsTick {
		dx, dy := int(u.Target[0])-dr.PosVec.X, int(u.Target[1])-dr.PosVec.Y
		distance := effectDistance(dx, dy) + 1
		if distance <= 10 {
			return effectDelete(dr)
		}
		speed := int(u.Speed)
		effectMove(dr, dr.PosVec.X+dx*speed/distance, dr.PosVec.Y+dy*speed/distance)
	}
	p := vp.ToScreenPos(dr.PosVec)
	p.Y -= 22
	radius := int(u.Size)
	if !effectInside(vp, p, radius) {
		return 1
	}
	effectGlow(p, glow, radius, 5)
	effectColor(core)
	effectPoint(p, radius>>1)
	delta := image.Pt(int(dr.Field_8)-dr.PosVec.X, int(dr.Field_9)-dr.PosVec.Y)
	effectColor(core)
	r := GetClient().R2()
	r.AddPoint(p)
	r.AddPointRel(delta)
	r.DrawLineFromPoints(r.Data().Color2())
	period := u.Period
	if period == 0 || !nox_client_drawIsTick {
		return 1
	}
	u.Timer--
	if u.Timer != 0 {
		return 1
	}
	u.Timer = period
	u.Size--
	if u.Size != 0 {
		return 1
	}
	return effectDelete(dr)
}

func effectBlueRainSpark(vp *noxrender.Viewport, dr *client.Drawable) int {
	result := effectMovingSpark(vp, dr, uint32(nox_color_white_2523948), uint32(dword_5d4594_1313540))
	if result == 1 && byte(dr.VelZ) >= 5 && nox_client_drawIsTick {
		if *effectMapped(1313688) == 0 {
			*effectMapped(1313688) = effectType("WhiteSpark")
		}
		if child := effectSpawn(int(*effectMapped(1313688)), dr.PosVec); child != nil {
			effectInitSpark(child, dr.PosVec, 1611, 10, 96)
			child.ZVal = uint16(effectRand(5, 15))
			child.ZVal2 = 0
			child.VelZ = int8(effectRand(0, 8))
			effectLink(child)
		}
		return effectDelete(dr)
	}
	return result
}
func effectRainOrb(vp *noxrender.Viewport, dr *client.Drawable) int {
	if dword_5d4594_1313692 == 0 {
		dword_5d4594_1313692 = uint32(effectType("RainOrbWhite"))
		*effectMapped(1313696) = effectType("RainOrbBlue")
	}
	u := (*client.DrawableUnionRainOrb)(unsafe.Pointer(&dr.Union))
	z := int16(dr.ZVal)
	if z > 0 {
		p := vp.ToScreenPos(dr.PosVec)
		p.Y -= int(z)
		color := *effectMapped(1313588)
		if dr.TypeIDVal != uint32(dword_5d4594_1313692) {
			color = uint32(dword_5d4594_1313536)
		}
		effectGlow(p, color, int(u.Size), 5)
		delta := int(z) - int(int16(u.PrevZ))
		effectColor(color)
		r := GetClient().R2()
		r.AddPoint(p)
		r.AddPointRel(image.Pt(0, delta))
		r.DrawLineFromPoints(r.Data().Color2())
		if nox_client_drawIsTick {
			u.PrevZ = dr.ZVal
			dr.ZVal += uint16(int16(dr.VelZ))
		}
		effectColor(*effectMapped(1313592))
		effectPoint(p, int(u.Size)/3)
		return 1
	}
	if *effectMapped(1313700) == 0 {
		*effectMapped(1313700) = effectType("WhiteMoveOrb")
		*effectMapped(1313704) = effectType("BlueMoveOrb")
	}
	// X subtraction is unsigned in C, while Y is assigned to an int before the
	// float conversion. Preserve that asymmetry and the float32 call boundary.
	angle := uintptr(effectAngle(float32(uint32(dr.PosVec.X)-u.ToX), float32(int32(uint32(dr.PosVec.Y)-u.ToY))))
	x := effectFloatInt(float32(float64(*memmap.PtrFloat32(0x587000, 194136+8*angle))*150 + float64(int32(u.ToX))))
	y := effectFloatInt(float32(float64(*memmap.PtrFloat32(0x587000, 194140+8*angle))*150 + float64(int32(u.ToY))))
	coords := [4]uint16{uint16(x), uint16(y), uint16(dr.PosVec.X), uint16(dr.PosVec.Y + 20)}
	typ := *effectMapped(1313700)
	if dr.TypeIDVal != uint32(dword_5d4594_1313692) {
		typ = *effectMapped(1313704)
	}
	speed := byte(effectRand(6, 8))
	effectCreateOrb(int(typ), &coords, 0, 0, speed, 1)
	return effectDelete(dr)
}

func effectColoredSpark(vp *noxrender.Viewport, dr *client.Drawable, kind int) int {
	var core, glow uint32
	switch kind {
	case 0:
		core, glow = uint32(dword_5d4594_1313532), *effectMapped(1313528)
	case 1:
		core, glow = uint32(dword_5d4594_1313540), uint32(dword_5d4594_1313536)
	case 2:
		core, glow = *effectMapped(1313548), *effectMapped(1313544)
	case 3:
		core, glow = *effectMapped(1313584), *effectMapped(1313580)
	case 4:
		core, glow = uint32(dword_5d4594_1313532), *effectMapped(1313576)
	case 5:
		core, glow = *effectMapped(1313556), *effectMapped(1313552)
	case 6:
		core, glow = *effectMapped(1313572), *effectMapped(1313568)
	case 7:
		core, glow = uint32(nox_color_white_2523948), uint32(dword_5d4594_1313540)
	}
	return effectMovingSpark(vp, dr, core, glow)
}

func effectOrbitUpdate(dr *client.Drawable) int {
	age := int(int32(gameFrame() - dr.AnimStart))
	u := (*client.DrawableUnionOrb)(unsafe.Pointer(&dr.Union))
	target := image.Pt(int(u.Target[0]), int(u.Target[1]))
	dx, dy := dr.PosVec.X-target.X, dr.PosVec.Y-target.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	if age >= 60 || dx < 10 && dy < 10 {
		return effectDelete(dr)
	}
	radius := int(int16(u.Radius))
	radius -= age * radius / 60
	turn := (age << 8) / 120
	angle := int(u.Angle)
	if u.Speed != 0 {
		angle += turn
	} else {
		angle -= turn
	}
	angle = int(byte(angle))
	x := target.X + radius*int(*memmap.PtrInt32(0x587000, 192088+8*uintptr(angle)))/16
	y := target.Y + radius*int(*memmap.PtrInt32(0x587000, 192092+8*uintptr(angle)))/16
	effectMove(dr, x, y)
	dr.Field_8 = uint32(x)
	dr.Field_9 = uint32(y)
	return 1
}
