package legacy

import (
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
)

func effectPlasmaSegment(from, to image.Point, slot int) int8 {
	vp := drawViewport()
	if *effectMapped(1316416) == 0 {
		*effectMapped(1316416) = effectType("CyanSpark")
	}
	off := uintptr(28 * (int(dword_5d4594_1316412) + 30*slot))
	*effectMapped(1313884 + off) = uint32(from.X)
	*effectMapped(1313888 + off) = uint32(from.Y)
	*effectMapped(1313912 + off) = uint32(to.X)
	*effectMapped(1313916 + off) = uint32(to.Y)
	dword_5d4594_1316412++
	result := int8(byte(gameFrame()))
	if gameFrame()&4 != 0 {
		result = int8(effectRand(0, 10))
		if result > 5 {
			p := vp.ToWorldPos(to)
			result = int8(effectCreateEnergySparks(p.X, p.Y, 8, int(*effectMapped(1316416))))
		}
	}
	return result
}
func effectPlasmaSetup(angle int, from, to image.Point) {
	dx, dy := to.X-from.X, to.Y-from.Y
	distance := effectDistance(dx, dy)
	segments := distance/40 + 1
	if distance/40+2 >= 30 {
		segments = 28
	}
	dword_5d4594_1316408 = uint32(segments)
	cosine := float64(*memmap.PtrFloat32(0x587000, 194136+8*uintptr(angle)))
	sine := float64(*memmap.PtrFloat32(0x587000, 194140+8*uintptr(angle)))
	fy := float32(dy)
	dword_5d4594_1313880 = uint32(math.Float32bits(fy))
	length := float32(math.Sqrt(float64(dy)*float64(fy)+float64(dx)*float64(dx)) + 0.0099999998)
	*memmap.PtrFloat32(0x5D4594, 1313876) = float32(float64(dx) / float64(length))
	normalizedY := float64(fy) / float64(length)
	dword_5d4594_1313880 = uint32(math.Float32bits(float32(normalizedY)))
	dot := normalizedY*sine + float64(*memmap.PtrFloat32(0x5D4594, 1313876))*cosine
	if dot < 0 {
		dot *= 0.2
	}
	scale := (1 - dot) * float64(distance) * 2.3
	*memmap.PtrFloat32(0x5D4594, 1313868) = float32(cosine * scale)
	*memmap.PtrFloat32(0x5D4594, 1313872) = float32(sine * scale)
	tangent := image.Pt(effectFloatInt(*memmap.PtrFloat32(0x5D4594, 1313868)), effectFloatInt(*memmap.PtrFloat32(0x5D4594, 1313872)))
	points := [4]image.Point{from, to, tangent, image.Pt(dx, dy)}
	for slot := 0; slot < 3; slot++ {
		phase := memmap.PtrFloat32(0x5D4594, 1313856+uintptr(slot*4))
		if slot == 0 {
			*phase = float32(float64(*phase) + 0.2)
		} else {
			*phase = float32(float64(*phase) + 0.25)
		}
		if *phase >= 1 {
			for n := segments + 1; n > 0; n-- {
				off := uintptr(1313872 + 28*(n+30*slot))
				// Copy only the five animation words; segment endpoints stay separate.
				for _, delta := range []uintptr{0, 4, 8, 12, 16} {
					*effectMapped(off + 20 + delta) = *effectMapped(off - 8 + delta)
				}
			}
			*phase = 0
			*effectMapped(1313908 + uintptr(840*slot)) = 0
		}
		dword_5d4594_1316412 = 0
		effectCurveSegments(points, segments, *phase, func(a, b image.Point) { effectPlasmaSegment(a, b, slot) })
	}
}
func effectPlasma(angle int, from, to image.Point) int {
	if *effectMapped(1316404) == 0 {
		blue := noxcolor.RGB5551Color(40, 180, 255).Color32()
		*effectMapped(1313828) = blue
		*effectMapped(1313832) = noxcolor.RGB5551Color(255, 255, 255).Color32()
		*effectMapped(1313836) = blue
		*effectMapped(1313844) = 8
		*effectMapped(1313848) = 12
		*effectMapped(1313852) = 8
		for i := 0; i < 3; i++ {
			*effectMapped(1313856 + uintptr(4*i)) = 0
		}
		*memmap.PtrUint8(0x5D4594, 1313840) = 16
		*memmap.PtrUint8(0x5D4594, 1313841) = 24
		*memmap.PtrUint8(0x5D4594, 1313842) = 16
		for i := 0; i < 90; i++ {
			for j := 0; j < 5; j++ {
				*effectMapped(1313892 + uintptr(28*i+4*j)) = 0
			}
		}
		*effectMapped(1316404) = 1
	}
	effectPlasmaSetup(angle, from, to)
	segments := int(int32(dword_5d4594_1316408))
	// the plasma wobble only advances on game ticks
	for slot := 0; slot < 3 && nox_client_drawIsTick; slot++ {
		for n := 0; n <= segments; n++ {
			off := uintptr(28 * (30*slot + n))
			heading := int(int32(*effectMapped(1313896 + off) + *effectMapped(1313900 + off)))
			if heading >= 256 {
				heading -= 256
			} else if heading < 0 {
				heading += 255
			}
			if n == 0 {
				heading = angle
			}
			*effectMapped(1313896 + off) = uint32(heading)
			velocity := int(int32(*effectMapped(1313904 + off) + *effectMapped(1313900 + off)))
			limit := 25
			if slot == 1 {
				limit = 10
			}
			if velocity > limit {
				velocity = limit
			} else if velocity < -limit {
				velocity = -limit
			}
			*effectMapped(1313900 + off) = uint32(velocity)
			remaining := int(int32(*effectMapped(1313908 + off) - 1))
			*effectMapped(1313908 + off) = uint32(remaining)
			if remaining <= 0 {
				*effectMapped(1313908 + off) = uint32(effectRand(10, 90))
				*effectMapped(1313904 + off) = uint32(effectRand(4, 8))
				if effectRand(0, 1) != 0 {
					*effectMapped(1313904 + off) = -*effectMapped(1313904 + off)
				}
				if slot == 1 {
					*effectMapped(1313892 + off) = uint32(effectRand(40, 50))
				} else {
					*effectMapped(1313892 + off) = uint32(effectRand(80, 110))
					if n < 4 && effectRand(0, 100) > 90 {
						*effectMapped(1313892 + off) = uint32(effectRand(150, 200))
						*effectMapped(1313908 + off) *= 2
					}
				}
			}
		}
	}
	for slot := 0; slot < 3; slot++ {
		for n := 0; n < segments; n++ {
			off := uintptr(28 * (30*slot + n))
			radius := float32(int32(*effectMapped(1313892 + off)))
			tangent := func(heading uint32) image.Point {
				x := radius * *memmap.PtrFloat32(0x587000, 194136+8*uintptr(heading))
				y := radius * *memmap.PtrFloat32(0x587000, 194140+8*uintptr(heading))
				return image.Pt(effectFloatInt(x), effectFloatInt(y))
			}
			a := image.Pt(int(*effectMapped(1313884 + off)), int(*effectMapped(1313888 + off)))
			b := image.Pt(int(*effectMapped(1313912 + off)), int(*effectMapped(1313916 + off)))
			ta, tb := tangent(*effectMapped(1313896 + off)), tangent(*effectMapped(1313924 + off))
			color := int(*effectMapped(1313828 + uintptr(4*slot)))
			effectCurveColor(color)
			effectCurveGlow(1, color, int(*effectMapped(1313844 + uintptr(4*slot))), int8(*memmap.PtrUint8(0x5D4594, 1313840+uintptr(slot))))
			effectCurveRaster([4]image.Point{a, b, ta, tb}, 8, bool2int(slot == 1))
		}
	}
	return int(int32(dword_5d4594_1316408))
}
func effectPlasmaDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	mouse := GetClient().GetMousePos()
	a, b, ok := effectRayEndpoints(dr, true)
	if !ok {
		return 1
	}
	a = vp.ToScreenPos(a)
	b = vp.ToScreenPos(b)
	a.Y -= 20
	b.Y -= 20
	u := (*client.DrawableUnionRay)(unsafe.Pointer(&dr.Union))
	angle := int(u.Arg[0])
	if u.ByCode != 0 {
		angle = int(byte(effectAngle(float32(float64(mouse.X)-float64(a.X)), float32(float64(mouse.Y)-float64(a.Y)))))
	}
	effectPlasma(angle, a, b)
	return 1
}
