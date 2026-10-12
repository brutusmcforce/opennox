package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func screenHarpoonRope(vp *noxrender.Viewport, dr *client.Drawable) int {
	var from, to image.Point
	if *effectByte(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) == 0 {
		from = vp.ToScreenPos(image.Pt(int(*effectShort(dr, int(layoutClientDrawable.Off(437)) /* Union._ */)), int(*effectShort(dr, int(layoutClientDrawable.Off(439)) /* Union._ */))))
		to = vp.ToScreenPos(image.Pt(int(*effectShort(dr, int(layoutClientDrawable.Off(441)) /* Union._ */)), int(*effectShort(dr, int(layoutClientDrawable.Off(443)) /* Union._ */))))
		from.Y -= 20
		to.Y -= 20
	} else {
		lookup := func(code uint32) *client.Drawable {
			if code&0x8000 != 0 {
				return GetClient().Cli().Objs.ByNetCodeStatic(int(code & 0x7fff))
			}
			return GetClient().Cli().Objs.ByNetCodeDynamic(int(code & 0x7fff))
		}
		a, b := lookup(*effectWord(dr, int(layoutClientDrawable.Off(437)) /* Union._ */)), lookup(*effectWord(dr, int(layoutClientDrawable.Off(441)) /* Union._ */))
		if a == nil || b == nil {
			return 1
		}
		from = vp.ToScreenPos(a.PosVec)
		to = effectScreen(vp, b)
		dir := uintptr(*effectByte(a, int(layoutClientDrawable.Off(297)) /* AnimDir */))
		from.X += int(*memmap.PtrInt32(0x587000, 175864+8*dir))
		from.Y += int(*memmap.PtrInt32(0x587000, 175868+8*dir))
		to.Y -= 8
	}
	*effectMapped(1312492) = uint32(noxcolor.RGB5551Color(144, 104, 64).Color32())
	*effectMapped(1312496) = uint32(noxcolor.RGB5551Color(24, 16, 0).Color32())
	screenRopeLine(from, to)
	return 1
}
func screenMaiden(vp *noxrender.Viewport, dr *client.Drawable) int {
	r := GetClient().R2().Data()
	s := GetServer().S()
	if !noxflags.HasGame(noxflags.GameFlag22) {
		npc := s.NPCs.ByID(int(dr.NetCode32))
		if npc == nil {
			return 1
		}
		for i, color := range npc.Color8 {
			r.SetMaterial(i+1, noxcolor.RGBA5551(color))
		}
	} else {
		for obj := s.Objs.First(); obj != nil; obj = obj.Next() {
			if uint32(obj.Extent) != dr.NetCode32 {
				continue
			}
			colors := unsafe.Slice((*byte)(unsafe.Add(obj.UpdateData, 2076)), 18)
			for i := 0; i < 6; i++ {
				r.SetMaterialRGB(i+1, int(colors[3*i]), int(colors[3*i+1]), int(colors[3*i+2]))
			}
			break
		}
	}
	return Nox_thing_monster_draw(vp, dr)
}
func screenUndead(vp *noxrender.Viewport, dr *client.Drawable) int {
	if *effectMapped(1313728) == 0 {
		*effectMapped(1313736) = effectType("WhiteOrb")
		*effectMapped(1313732) = uint32(noxcolor.RGB5551Color(100, 100, 255).Color32())
		*effectMapped(1313728) = 1
	}
	if gameFrame()-*effectWord(dr, int(layoutClientDrawable.Off(316)) /* AnimStart */) > 70 {
		return effectDelete(dr)
	}
	if nox_client_drawIsTick && effectRand(0, 100) > 85 {
		coords := [4]uint16{*effectShort(dr, int(layoutClientDrawable.Off(324)) /* Field_81 */), *effectShort(dr, int(layoutClientDrawable.Off(328)) /* Field_82 */)}
		coords[2] = uint16(dr.PosVec.X) + uint16(effectRand(-5, 5))
		coords[3] = uint16(dr.PosVec.Y) + uint16(effectRand(-5, 5))
		speed := byte(effectRand(6, 10))
		dy := effectRand(-5, 5)
		effectCreateOrb(int(*effectMapped(1313736)), &coords, 0, dy, speed, 0)
	}
	distance := int(screenDistanceBetween(int32(dr.PosVec.X), int32(dr.PosVec.Y), int32(*effectWord(dr, int(layoutClientDrawable.Off(324)) /* Field_81 */)), int32(*effectWord(dr, int(layoutClientDrawable.Off(328)) /* Field_82 */))))
	radius := 8 - distance/40
	if radius < 0 {
		radius = 1
	}
	effectGlow(effectScreen(vp, dr), *effectMapped(1313732), radius, 12)
	return 1
}
func screenWaypoint(vp *noxrender.Viewport, dr *client.Drawable) int {
	color := *memmap.PtrUint32(0x85B3FC, 940)
	pos := dr.PosVec.Sub(vp.World.Min)
	screenCircle(pos.X, pos.Y, 10, color)
	angle := int(byte(2 * byte(gameFrame())))
	GetClient().R2().Data().SetAlphaEnabled(true)
	effectColor(color)
	point := func(i int) image.Point {
		return pos.Add(image.Pt(10*int(*memmap.PtrInt32(0x587000, 192088+8*uintptr(i)))/16, 10*int(*memmap.PtrInt32(0x587000, 192092+8*uintptr(i)))/16))
	}
	for i := 0; i < 5; i++ {
		next := (angle + 102) % 256
		effectLine(point(angle), point(next))
		angle = next
	}
	GetClient().R2().Data().SetAlphaEnabled(false)
	return 1
}
