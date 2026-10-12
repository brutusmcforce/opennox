package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"math"
	"unsafe"
)

func updateType(off uintptr, name string) uint32 {
	p := effectMapped(off)
	if *p == 0 {
		*p = effectType(name)
	}
	return *p
}
func updateRadial(pos image.Point, angle, radius int) image.Point {
	return pos.Add(image.Pt(radius*int(*memmap.PtrInt32(0x587000, 192088+uintptr(8*angle)))/16, radius*int(*memmap.PtrInt32(0x587000, 192092+uintptr(8*angle)))/16))
}
func updateCloud(dr *client.Drawable, count, radius int) int {
	if *effectMapped(1523016) == 0 {
		*effectMapped(1523016) = effectType("GreenPuff")
		*effectMapped(1523020) = effectType("GreenSmoke")
	}
	for count > 0 {
		angle := effectRand(0, 255)
		r := effectRand(0, radius)
		pos := updateRadial(dr.PosVec, angle, r)
		smoke := effectRand(0, 10) < 3
		typ := *effectMapped(1523016)
		if smoke {
			typ = *effectMapped(1523020)
		}
		if child := effectSpawn(int(typ), pos); child != nil {
			child.ZVal = 0
			effectLink(child)
			*effectByte(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = byte(effectRand(1, 3))
			child.Field_115 = drawableUpdateIdentity(updateID_sub_4CE340)
			GetClient().Cli().Objs.TransparentDecay(child, effectRand(10, 32))
			GetClient().Cli().Objs.List5Add(child)
			GetClient().Cli().Objs.List6Add(child)
		}
		count--
	}
	return count
}
func updateCloudRise(dr *client.Drawable) int { dr.ZVal += uint16(*effectByte(dr, int(layoutClientDrawable.Off(432)) /* Union._ */)); return 1 }
func updateCloudFrame(dr *client.Drawable, radius int) int {
	if gameFrame()&1 != 0 {
		updateCloud(dr, 1, radius)
	}
	return 1
}

// The original helper's return was the last allocation, or the nonpositive count.
func updateDeathBallSparks(dr *client.Drawable, count int) uint32 {
	typ := updateType(1523008, "DeathBallSpark")
	out := uint32(count)
	for ; count > 0; count-- {
		child := effectSpawn(int(typ), dr.PosVec)
		out = memmap.PtrWord(unsafe.Pointer(child))
		if child == nil {
			continue
		}
		*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(dr.PosVec.X) << 12
		*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.Y) << 12
		child.Field_74_4 = byte(effectRand(0, 255))
		*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = uint32(effectRand(1000, 3000))
		*effectWord(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(10, 40))
		*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
		child.ZVal = 22
		child.VelZ = int8(effectRand(0, 4))
		effectLink(child)
	}
	return out
}
func updateDeathBallCharge(dr *client.Drawable) int {
	typ := updateType(1523012, "CharmOrb")
	coords := [4]uint16{uint16(dr.PosVec.X), uint16(dr.PosVec.Y)}
	for i := 0; i < 10; i++ {
		angle := effectRand(0, 255)
		radius := effectRand(2, 8)
		coords[2] = coords[0] + uint16(radius*int(*memmap.PtrInt16(0x587000, 192088+uintptr(8*angle))))
		coords[3] = coords[1] + uint16(radius*int(*memmap.PtrInt16(0x587000, 192092+uintptr(8*angle))))
		if effectRand(0, 100) < 50 {
			speed := byte(effectRand(6, 10))
			y := effectRand(-20, 20)
			x := effectRand(-20, 20)
			effectCreateOrb(int(typ), &coords, x, y, speed, 0)
		}
	}
	return 1
}
func updateFireball(dr *client.Drawable, speed int) {
	typ := updateType(1522964, "Spark")
	dx, dy := dr.PosVec.X-int(dr.Field_8), dr.PosVec.Y-int(dr.Field_9)
	ax, ay := dx, dy
	if ax < 0 {
		ax = -ax
	}
	if ay < 0 {
		ay = -ay
	}
	for n := (ax + ay) / 7; n > 0; n-- {
		fraction := effectRand(0, 100)
		pos := image.Pt(int(dr.Field_8)+dx*fraction/100, int(dr.Field_9)+dy*fraction/100)
		if child := effectSpawn(int(typ), pos); child != nil {
			*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(pos.X) << 12
			*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(pos.Y) << 12
			jitter := effectRand(-25, 25)
			child.Field_74_4 = byte(effectAngle(float32(-dx), float32(-dy)) + jitter)
			*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = uint32(speed * effectRand(100, 300))
			*effectWord(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(30, 45))
			*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
			child.ZVal, child.ZVal2 = 28, 0
			child.VelZ = int8(effectRand(-2, 4))
			effectLink(child)
		}
	}
}
func updateFireballFrame(dr *client.Drawable, speed int) int {
	if dr.Field_120 == 0 && !noxflags.HasGame(noxflags.GamePause) {
		updateFireball(dr, speed)
	}
	return 1
}
func updateManaBomb(dr *client.Drawable) int {
	radius := int(GetServer().S().Balance.Float("ManaBombOutRadius"))
	if dword_5d4594_1522956 == 0 {
		dword_5d4594_1522956 = uint32(effectType("ManaBombOrb"))
		*effectMapped(1522960) = effectType("VioletSpark")
	}
	for i := 0; i < 20; i++ {
		r := radius/4 + effectRand(0, radius)
		if r > radius {
			r = radius
		}
		angle := effectRand(0, 255)
		if child := effectSpawn(int(*effectMapped(1522960)), updateRadial(dr.PosVec, angle, r)); child != nil {
			*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(child.PosVec.X) << 12
			*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(child.PosVec.Y) << 12
			child.Field_74_4 = 0
			*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = 0
			*effectWord(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(10, 30))
			*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
			child.ZVal = 0
			child.VelZ = int8(effectRand(2, 8))
			effectLink(child)
		}
	}
	if gameFrame()&1 != 0 && gameFrame()-dr.Field_80 < 10 {
		coords := [4]int16{int16(dr.PosVec.X), int16(dr.PosVec.Y)}
		for angle := int(gameFrame() % 51); angle < 256; angle += 51 {
			coords[2] = coords[0] + int16(radius/16*int(*memmap.PtrInt16(0x587000, 192088+uintptr(8*angle))))
			coords[3] = coords[1] + int16(radius/16*int(*memmap.PtrInt16(0x587000, 192092+uintptr(8*angle))))
			effectCreateOrbit(int(dword_5d4594_1522956), &coords, int16(angle), 0, 0)
			effectCreateOrbit(int(dword_5d4594_1522956), &coords, int16(angle), 1, 0)
		}
	}
	return 1
}
func updateMagicMissile(dr *client.Drawable) int {
	if *effectMapped(1522984) == 0 {
		*effectMapped(1522984) = effectType("Spark")
		*effectMapped(1522988) = effectType("Puff")
		*effectMapped(1522992) = effectType("MagicMissileTailLink")
	}
	dx, dy := dr.PosVec.X-int(dr.Field_8), dr.PosVec.Y-int(dr.Field_9)
	densityWord := *memmap.PtrUint32(0x587000, 190108)
	density := int(densityWord)
	// The C entry gate is unsigned, but its loop termination is signed.
	if densityWord != 0 && density < 0 {
		density = 1
	}
	for i := 0; i < density; i++ {
		x := int(dr.Field_8) + dx*effectRand(0, 100)/100
		y := int(dr.Field_9) + dy*effectRand(0, 100)/100
		if child := effectSpawn(int(*effectMapped(1522984)), image.Pt(x, y)); child != nil {
			*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(x) << 12
			*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(y) << 12
			child.Field_74_4 = byte(effectRand(0, 255))
			*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = 0
			*effectWord(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(3, 10))
			*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
			child.ZVal = 20
			child.VelZ = int8(effectRand(0, 6))
			effectLink(child)
		}
	}
	updateTailLink(dr, *effectMapped(1522992), int(gameFPS()/3))
	return 1
}

// Advance the anchor only after ownership and decay can be assigned to a tail.
func updateTailLink(dr *client.Drawable, typ uint32, lifetime int) {
	anchor := image.Pt(int(*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */)), int(*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */)))
	dx, dy := dr.PosVec.X-anchor.X, dr.PosVec.Y-anchor.Y
	if uint32(dx*dx+dy*dy) <= 200 {
		return
	}
	if child := effectSpawn(int(typ), anchor); child != nil {
		*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(dr.PosVec.X)
		*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.Y)
		effectLink(child)
		*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(dr.PosVec.X)
		*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.Y)
		GetClient().Cli().Objs.TransparentDecay(child, lifetime)
	}
}
func updateMagicTrail(dr *client.Drawable) {
	dx, dy := dr.PosVec.X-int(*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */)), dr.PosVec.Y-int(*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */))
	if uint32(dx*dx+dy*dy) > 200 {
		updateTailLink(dr, updateType(1523000, "MagicTailLink"), int(gameFPS()))
	}
	updateTrailSparks(dr, true)
}
func updateTrailSparks(dr *client.Drawable, magic bool) int {
	off, name := uintptr(1522996), "BlueSpark"
	count, jitter, lo, hi := 5, 3, 2, 10
	if magic {
		off = 1523004
		count, jitter, lo, hi = 4, 8, 10, 20
	}
	typ := updateType(off, name)
	dx, dy := dr.PosVec.X-int(dr.Field_8), dr.PosVec.Y-int(dr.Field_9)
	sx, sy := 0, 0
	for i := 0; i < count; i++ {
		x := int(dr.Field_8) + sx/count + effectRand(-jitter, jitter)
		y := int(dr.Field_9) + sy/count + effectRand(-jitter, jitter)
		if child := effectSpawn(int(typ), image.Pt(x, y)); child != nil {
			pos := image.Pt(x, y)
			if magic {
				child.DrawFuncPtr = unsafe.Pointer(drawableDrawKey(drawKey_nox_thing_magic_sparkle_draw))
				child.SetLightColor(128, 128, 255)
				pos = dr.PosVec
			} else {
				child.DrawFuncPtr = unsafe.Pointer(drawableDrawKey(drawKey_nox_thing_pixie_dust_draw))
				child.SetLightColor(255, 200, 75)
			}
			*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(pos.X) << 12
			*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(pos.Y) << 12
			child.Field_74_4 = 0
			*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = 0
			*effectWord(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(lo, hi))
			*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
			child.ZVal, child.ZVal2 = dr.ZVal, dr.ZVal2
			child.VelZ = 0
			effectLink(child)
		}
		sx += dx
		sy += dy
	}
	return 1
}
func updateTeleportWake(dr *client.Drawable) int {
	typ := updateType(1522980, "BlueSpark")
	x := dr.PosVec.X + effectRand(-5, 5)
	y := dr.PosVec.Y + effectRand(-5, 5)
	if child := effectSpawn(int(typ), image.Pt(x, y)); child != nil {
		effectInitSpark(child, child.PosVec, 100, 10, 32)
		child.ZVal = 0
		child.VelZ = int8(effectRand(3, 8))
		effectLink(child)
	}
	return 1
}
func updateVortex(dr *client.Drawable) int {
	if *effectMapped(1522952) == 0 {
		*effectMapped(1522952) = effectType("WhiteVortexOrb")
		*effectMapped(1522944) = noxcolor.RGB5551Color(200, 200, 200).Color32()
		*effectMapped(1522948) = noxcolor.RGB5551Color(255, 255, 255).Color32()
	}
	angle := effectRand(0, 255)
	if child := effectSpawn(int(*effectMapped(1522952)), updateRadial(dr.PosVec, angle, 50)); child != nil {
		*effectByte(child, int(layoutClientDrawable.Off(448)) /* Union._ */) = byte(angle)
		dr.ZVal = 0 // The source clears the parent height, not the new orb's height.
		*effectByte(child, int(layoutClientDrawable.Off(449)) /* Union._ */) = byte(effectRand(2, 3))
		if effectRand(0, 100) > 50 {
			*effectByte(child, int(layoutClientDrawable.Off(449)) /* Union._ */) = -*effectByte(child, int(layoutClientDrawable.Off(449)) /* Union._ */)
		}
		*effectByte(child, int(layoutClientDrawable.Off(451)) /* Union._ */) = 1
		*effectByte(child, int(layoutClientDrawable.Off(450)) /* Union._ */) = 50
		*effectWord(child, int(layoutClientDrawable.Off(440)) /* Union._ */) = uint32(dr.PosVec.X)
		*effectWord(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = uint32(dr.PosVec.Y)
		*effectWord(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = *effectMapped(1522944)
		*effectWord(child, int(layoutClientDrawable.Off(436)) /* Union._ */) = *effectMapped(1522948)
		effectLink(child)
		GetClient().Cli().Objs.List6Add(child)
	}
	return 1
}
func updateTransfer(typ int, vp *noxrender.Viewport, dr *client.Drawable, reverse, unsignedBounds bool) {
	if effectRand(0, 100) >= 50 {
		return
	}
	from, to, ok := effectRayEndpoints(dr, true)
	if !ok {
		return
	}
	if reverse {
		from, to = to, from
	}
	x := from.X + effectRand(-20, 20)
	y := from.Y + effectRand(-20, 20)
	view := unsafe.Slice((*int32)(vp.C()), 10)
	z := int(int16(dr.ZVal))
	screenX := x + int(view[0]) - int(view[4])
	screenY := y + int(view[1]) - int(view[5]) - z
	if screenX < 0 {
		x = int(view[4]+view[0]) + 1
	}
	if screenY < 0 {
		y = int(view[1]+view[5]) - z + 1
	}
	// Charm/heal retain the unsigned viewport comparison in their C callers.
	right, bottom := screenX >= int(view[8]), screenY >= int(view[9])
	if unsignedBounds {
		right = uint32(screenX) >= uint32(view[8])
		bottom = uint32(screenY) >= uint32(view[9])
	}
	if right {
		x = int(view[2]+view[4]) - 1
	}
	if bottom {
		y = int(view[3]+view[5]) - z - 1
	}
	if child := effectSpawn(typ, image.Pt(x, y)); child != nil {
		speed := effectRand(6, 12)
		*effectShort(child, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint16(to.X)
		*effectShort(child, int(layoutClientDrawable.Off(434)) /* Union._ */) = uint16(to.Y)
		*effectByte(child, int(layoutClientDrawable.Off(443)) /* Union._ */) = byte(speed)
		*effectByte(child, int(layoutClientDrawable.Off(444)) /* Union._ */) = byte(effectRand(3, 10))
	}
}
func updateHealDrain(vp *noxrender.Viewport, dr *client.Drawable, heal bool) int {
	off, name := uintptr(1522976), "DrainManaOrb"
	if heal {
		off, name = 1522972, "HealOrb"
	}
	updateTransfer(int(updateType(off, name)), vp, dr, true, heal)
	return 1
}
func updateCharm(vp *noxrender.Viewport, dr *client.Drawable) int {
	if dword_5d4594_1522968 == 0 {
		dword_5d4594_1522968 = uint32(effectType("CharmOrb"))
	}
	updateTransfer(int(dword_5d4594_1522968), vp, dr, true, true)
	updateTransfer(int(dword_5d4594_1522968), vp, dr, false, true)
	return 1
}
func updateHeight(dr *client.Drawable, bounce bool) int {
	h := (*float32)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(436)) /* Union._ */))
	v := (*float32)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(440)) /* Union._ */))
	restitution := math.Float32frombits(*effectWord(dr, int(layoutClientDrawable.Off(444)) /* Union._ */))
	for frame := *effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */); frame < gameFrame(); frame++ {
		if !bounce {
			if *h > 0 {
				*h = float32(float64(*h) + float64(*v))
				*v = float32(float64(*v) - 1)
			}
			if *h <= 0 {
				*h, *v = 0, 0
			}
		} else {
			next := float64(*h) + float64(*v)
			*h = float32(next)
			if next >= 0 {
				*v = float32(float64(*v) - 0.5)
			} else {
				value := -float64(*v) * float64(restitution) * 0.1
				*h = 0
				*v = float32(value)
				if value < 2 {
					*h, *v = 0, 0
				}
			}
		}
	}
	dr.ZVal = uint16(int64(*h))
	dr.VelZ = int8(int64(*v))
	*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = gameFrame()
	return 1
}
