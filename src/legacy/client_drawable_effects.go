package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

var drawableSummonSpark uint32
var drawableShieldFound uint32

func drawableSummonStart(owner int16, pos [2]uint16, typ uint16, dir byte, value int16) uint32 {
	outer := effectSpawn(int(effectType("SummonEffect")), image.Pt(int(pos[0]), int(pos[1])))
	if outer == nil {
		return 0
	}
	inner := GetClient().Nox_new_drawable_for_thing(int(typ))
	if inner == nil {
		return 0
	}
	inner.PosVec = image.Pt(int(pos[0]), int(pos[1]))
	inner.AnimDir = byte(geometryDirection4Index(int32(dir)))
	inner.AnimInd = 8
	*(**client.Drawable)(unsafe.Add(outer.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)) = inner
	*(*uint32)(unsafe.Add(outer.C(), int(layoutClientDrawable.Off(436)) /* Union._ */)) = uint32(uint16(owner))<<16 | uint32(uint16(value))
	frame := GetServer().S().Frame()
	outer.AnimStart = frame
	return frame
}
func drawableSummonStop(owner int16) {
	typ := memmap.PtrUint32(0x5D4594, 1313744)
	if *typ == 0 {
		*typ = uint32(effectType("SummonEffect"))
	}
	if drawableSummonSpark == 0 {
		drawableSummonSpark = uint32(effectType("BlueSpark"))
	}
	for dr := GetClient().Cli().Objs.List1; dr != nil; dr = dr.NextPtr {
		if dr.TypeIDVal != *typ || *(*uint16)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(438)) /* Union._ */)) != uint16(owner) {
			continue
		}
		effectCreatePointSparks(int(drawableSummonSpark), 50, 1000, 30, dr.PosVec.X, dr.PosVec.Y)
		GetClient().Nox_xxx_spriteDelete_45A4B0(*(**client.Drawable)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)))
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return
	}
}
func drawableShieldLoad() int {
	last := 0
	for i, name := range []string{"SphericalShieldNW", "SphericalShieldN", "SphericalShieldNE", "SphericalShieldW", "", "SphericalShieldE", "SphericalShieldSW", "SphericalShieldS", "SphericalShieldSE"} {
		if i == 4 {
			continue
		}
		last = int(effectType(name))
		*memmap.PtrUint32(0x5D4594, 1313748+4*uintptr(i)) = uint32(last)
	}
	*memmap.PtrUint32(0x5D4594, 1313764) = 0
	*memmap.PtrUint32(0x5D4594, 1313784) = 1
	return last
}
func drawableShieldScan(dr *client.Drawable, code uint32) {
	for i := 0; i < 9; i++ {
		if dr.TypeIDVal == *memmap.PtrUint32(0x5D4594, 1313748+4*uintptr(i)) && *(*uint32)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)) == code {
			drawableShieldFound = 1
		}
	}
}
func drawableShield(code uint32, dir int) uintptr {
	if *memmap.PtrUint32(0x5D4594, 1313784) == 0 {
		drawableShieldLoad()
	}
	// Preserve legacy lookups, including its misspelled north name.
	switch dir {
	case 0, 2:
		effectType("SphericalShieldNW")
	case 1:
		effectType("ShpericalShieldN")
	case 3:
		effectType("SphericalShieldW")
	case 5:
		effectType("SphericalShieldE")
	case 6:
		effectType("SphericalShieldSW")
	case 7:
		effectType("SphericalShieldS")
	case 8:
		effectType("SphericalShieldSE")
	}
	target := GetClient().Cli().Objs.ByNetCode(uint16(code))
	if target == nil {
		return 0
	}
	pos := target.PosVec
	drawableShieldFound = 0
	GetClient().Cli().Objs.EachInRect(image.Rect(pos.X-10, pos.Y-10, pos.X+10, pos.Y+10), func(dr *client.Drawable) { drawableShieldScan(dr, code) })
	if drawableShieldFound == 1 {
		return 1
	}
	dr := effectSpawn(int(*memmap.PtrUint32(0x5D4594, 1313748+4*uintptr(dir))), pos.Add(image.Pt(0, 3)))
	if dr != nil {
		*(*uint32)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(432)) /* Union._ */)) = code
	}
	return uintptr(dr.C())
}
func drawableEffectTypes() int {
	if v := effectType("Spark"); v != 0 {
		*memmap.PtrUint32(0x5D4594, 1200772) = uint32(v)
	} else {
		*memmap.PtrUint32(0x5D4594, 1200772) = 0
		return 0
	}
	clientGameBlueSpark = effectType("BlueSpark")
	if clientGameBlueSpark == 0 {
		return 0
	}
	for i, name := range []string{"YellowSpark", "CyanSpark", "GreenSpark", "Puff"} {
		v := effectType(name)
		*memmap.PtrUint32(0x5D4594, 1200780+4*uintptr(i)) = uint32(v)
		if v == 0 {
			return 0
		}
	}
	for i := 0; i < 5; i++ {
		name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 161216+4*uintptr(i))))
		v := effectType(name)
		*memmap.PtrUint32(0x5D4594, 1200812+4*uintptr(i)) = uint32(v)
		if v == 0 {
			return 0
		}
	}
	clientGameVioletSpark = effectType("VioletSpark")
	return bool2int(clientGameVioletSpark != 0)
}

func sub_4B7EE0(owner int16) { drawableSummonStop(int16(owner)) }

func nox_xxx_fxShield_4B8090(code uint32, dir int32) *uint32 {
	return (*uint32)(unsafe.Pointer(drawableShield(uint32(code), int(dir))))
}
