package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

func collisionBoxDistance(p *types.Pointf, radius float32, u *server.Object, normal *types.Pointf) float64 {
	const scale = 0.70709997
	b := u.Shape.Box
	x, y := float64(u.NewPos.X), float64(u.NewPos.Y)
	px, py := float64(p.X), float64(p.Y)
	x1, y1 := x+float64(b.LeftTop), y+float64(b.LeftBottom)
	x2, y2 := x+float64(b.LeftBottom2), y+float64(b.LeftTop2)
	// The compiled C spills this corner Y and the four signed projections.
	x3, y3 := x+float64(b.RightTop), float64(float32(y+float64(b.RightBottom)))
	x4, y4 := x+float64(b.RightBottom2), y+float64(b.RightTop2)
	a := float32((y1 - x1 + px - py) * scale)
	c := float32((y2 - x2 + px - py) * scale)
	d := float32((y3 + x3 - px - py) * scale)
	e := float32((y1 + x1 - px - py) * scale)
	region := 0
	if a <= 0 {
		if c < 0 {
			region = 2
		}
	} else {
		region = 1
	}
	if d >= 0 {
		if !(e <= 0) {
			region |= 4
		}
	} else {
		region |= 8
	}
	switch region {
	case 0:
		dx, dy := px-x, py-y
		if dx == 0 && dy == 0 {
			index := GetServer().S().Rand.Logic.IntClamp(0, 3)
			dx = float64(memmap.Float32(0x587000, uintptr(289928+8*index)))
			dy = float64(memmap.Float32(0x587000, uintptr(289932+8*index)))
		}
		distance := math.Sqrt(dy*dy + dx*dx)
		if distance == 0 {
			distance = 0.1
		}
		*normal = types.Pointf{float32(dx / distance), float32(dy / distance)}
		return distance
	case 1:
		*normal = types.Pointf{scale, -scale}
		return float64(radius) - float64(a)
	case 2:
		*normal = types.Pointf{-scale, scale}
		return float64(radius) + float64(c)
	case 4:
		*normal = types.Pointf{-scale, -scale}
		return float64(radius) - float64(e)
	case 8:
		*normal = types.Pointf{scale, scale}
		return float64(radius) + float64(d)
	}
	var dx, dy float64
	switch region {
	case 5:
		dx, dy = px-x1, py-y1
	case 6:
		dx, dy = px-x2, py-y2
	case 9:
		dx, dy = px-x3, py-y3
	case 10:
		dx, dy = px-x4, py-y4
	default:
		return -1
	}
	wide := math.Sqrt(dy*dy + dx*dx)
	distance := float32(wide)
	if wide == 0 {
		distance = 0.1
	}
	result := float64(radius) - float64(distance)
	if result >= 0 {
		*normal = types.Pointf{float32(dx / float64(distance)), float32(dy / float64(distance))}
	}
	return result
}
func collisionCircleBox(a, b *server.Object, mode int32) {
	var normal types.Pointf
	distance := collisionBoxDistance(&a.NewPos, a.Shape.Circle.R, b, &normal)
	if !(distance >= 0) {
		return
	}
	if a.ObjClass&0x2204 != 0 && b.ObjClass&0x2204 != 0 && !GetServer().S().MapTraceRayAt(a.PosVec, b.PosVec, nil, nil, 0) {
		return
	}
	geometryHit(a, b, &normal)
	if a.ObjFlags&8 == 0 && b.ObjFlags&8 == 0 && (mode != 0 || a.ObjClass&6 == 0 || b.ObjFlags&0x2000 == 0) {
		spring := float64(geometryObjectForce()) * float64(float32(distance))
		fx, fy := float32(float64(normal.X)*spring), float32(float64(normal.Y)*spring)
		vx, vy := float64(a.VelVec.X)-float64(b.VelVec.X), float64(a.VelVec.Y)-float64(b.VelVec.Y)
		tx, ty := -float64(normal.Y), float64(normal.X)
		speed := float32(tx*vx + vy*ty)
		mass := a.Mass
		if b.Mass <= a.Mass {
			mass = b.Mass
		}
		friction := float64(mass) * float64(speed)
		fx = float32(float64(fx) - friction*tx*0.69999999)
		fy = float32(float64(fy) - friction*ty*0.69999999)
		if mode != 0 {
			b.Sub548600(types.Pointf{-fx, -fy})
		} else {
			a.Sub548600(types.Pointf{fx, fy})
		}
	}
	geometryWake(a)
	geometryWake(b)
}
