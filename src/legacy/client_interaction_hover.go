package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func interactionHoverCode() uint32 {
	return drawableUnitCode((*client.Drawable)(interactionHoverDrawable))
}
func interactionHoverEnumerate() {
	if memmap.Uint32(0x5D4594, 1096632) == 0 {
		*memmap.PtrUint32(0x5D4594, 1096632) = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	p := Sub_473970(GetClient().GetMousePos())
	interactionHoverDrawable = nil
	interactionUsableCursorDrawable = nil
	*memmap.PtrUint32(0x5D4594, 1096628) = 0
	interactionHoverDepth = 0
	rect := image.Rectangle{Min: image.Pt(p.X-96, p.Y-96), Max: image.Pt(p.X+96, p.Y+96)}
	GetClient().Cli().Objs.EachInRect(rect, func(dr *client.Drawable) { interactionHover(dr, p) })
}
func interactionHover(dr *client.Drawable, point image.Point) {
	if memmap.Uint32(0x5D4594, 1096648) == 0 {
		*memmap.PtrUint32(0x5D4594, 1096648) = uint32(GetClient().Cli().Things.IndByID("Polyp"))
	}
	player := *(**client.Drawable)(memmap.PtrOff(0x852978, 8))
	if dr == player {
		return
	}
	if uint32(dr.ObjFlags)&0x8000 != 0 || interactionHasBuff(dr, 0) && !interactionHasBuff(player, 21) {
		return
	}
	class, sub := uint32(dr.ObjClass), uint32(dr.ObjSubClass)
	if class&2 != 0 && sub&0x4000 != 0 {
		return
	}
	if class&0x80400206 == 0 && dr.TypeIDVal != memmap.Uint32(0x5D4594, 1096648) {
		return
	}
	if !GetClient().Nox_xxx_client_4984B0_drawable(dr) {
		return
	}
	if class&4 != 0 {
		p := GetServer().S().Players.ByID(int(dr.NetCode32))
		if p == nil || p.Field3680&1 != 0 {
			return
		}
	}
	if class&0x400000 != 0 && sub&0x80 == 0 || class&2 != 0 && dr.AnimInd == 10 {
		return
	}
	x, y := int32(dr.PosVec.X), int32(dr.PosVec.Y)
	z := int16(dr.ZVal)
	top := floatToInt32(float32(float64(y) - float64(dr.ZSizeMax) - float64(z)))
	bottom := floatToInt32(float32(float64(y) - float64(dr.ZSizeMin) - float64(z)))
	px, py := int32(point.X), int32(point.Y)
	hit := false
	switch dr.Shape.Kind {
	case server.ShapeKindCircle:
		radius := floatToInt32(dr.Shape.Circle.R)
		square := floatToInt32(float32(float64(dr.Shape.Circle.R) * float64(dr.Shape.Circle.R)))
		center := bottom
		if py <= bottom {
			center = top
			if py >= top {
				if px <= x-radius || px >= x+radius {
					return
				}
				hit = true
			}
		}
		if !hit {
			dx, dy := uint32(px-x), uint32(py-center)
			if int32(dx*dx+dy*dy) >= square {
				return
			}
			hit = true
		}
	case server.ShapeKindBox:
		pos := types.Pointf{X: float32(x), Y: float32(bottom)}
		pt := types.Pointf{X: float32(px), Y: float32(py)}
		shape := (*[11]float32)(unsafe.Pointer(&dr.Shape))
		hit = collisionContains(&pos, shape, &pt)
		if !hit {
			pos.Y = float32(top)
			hit = collisionContains(&pos, shape, &pt)
		}
		if !hit {
			left := x + floatToInt32(dr.Shape.Box.LeftBottom2)
			lo := top + floatToInt32(dr.Shape.Box.LeftTop2)
			hi := bottom + floatToInt32(dr.Shape.Box.LeftTop2)
			hit = uint32(px) > uint32(left) && uint32(px) < uint32(x) && py > lo && py < hi
		}
		if !hit {
			right := x + floatToInt32(dr.Shape.Box.RightTop)
			lo := top + floatToInt32(dr.Shape.Box.RightBottom)
			hi := bottom + floatToInt32(dr.Shape.Box.RightBottom)
			if px < x || px >= right || py <= lo || py >= hi {
				return
			}
			hit = true
		}
	default:
		return
	}
	if !hit {
		return
	}
	depth := floatToInt32(float32(float64(z) + float64(y) + float64(dr.ZSizeMin)))
	if depth > memmap.Int32(0x5D4594, 1096628) {
		*memmap.PtrUint32(0x5D4594, 1096628) = uint32(depth)
		interactionHoverDrawable = dr.C()
	}
	if dr != *(**client.Drawable)(memmap.PtrOff(0x852978, 8)) && depth > int32(interactionHoverDepth) && glyphSelectionAllowed(dr) != 0 {
		p := Get_dword_8531A0_2576()
		if p != nil && p.PlayerClass() == 1 && dr.TypeIDVal == memmap.Uint32(0x5D4594, 1096632) {
			if interactionUsableCursorDrawable == nil {
				interactionUsableCursorDrawable = dr.C()
				interactionHoverDepth = 0
			}
		} else {
			interactionHoverDepth = uint32(depth)
			interactionUsableCursorDrawable = dr.C()
		}
	}
}

func nox_xxx_packetGetMarshall_476F40() uint32 { return uint32(interactionHoverCode()) }

func nox_xxx_clientEnumHover_476FA0() { interactionHoverEnumerate() }

func nox_xxx_clientOnCursorHover_477050(dr, p int32) {
	point := (*[2]int32)(unsafe.Pointer(uintptr(uint32(p))))
	interactionHover((*client.Drawable)(unsafe.Pointer(uintptr(uint32(dr)))), image.Pt(int(point[0]), int(point[1])))
}

func nox_xxx_guiCursor_477600() int32 { return int32(memmap.Uint32(0x5D4594, 1096672)) }
