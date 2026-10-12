package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

func geometryObjectForce() float32 { return math.Float32frombits(uint32(collisionObjectForce)) }
func geometryWallForce() float32   { return math.Float32frombits(uint32(collisionWallForce)) }
func geometryHit(u, v *server.Object, normal *types.Pointf) {
	collisionAddHit(u, collisionObjectAddress(v), normal)
}
func geometryActivate(u *server.Object) { collisionActivate(u) }
func geometryWake(u *server.Object) {
	if u.ObjFlags&0x8000000 != 0 {
		geometryActivate(u)
		u.ObjFlags &^= 0x8000000
	}
}
func geometryGameBall(u *server.Object) int32 {
	p := memmap.PtrUint32(0x5D4594, 2491788)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID("GameBall"))
	}
	return int32(bool2int(uint32(u.TypeInd) == *p))
}
func geometryPointWall(u *server.Object, p *types.Pointf) int32 {
	dx, dy := float64(u.NewPos.X)-float64(p.X), float64(u.NewPos.Y)-float64(p.Y)
	distance64 := math.Sqrt(dy*dy + dx*dx)
	distance := float32(distance64)
	if distance64 == 0 {
		distance = 0.0099999998
	}
	if distance >= u.Shape.Circle.R {
		return 0
	}
	nx, ny := dx/float64(distance), dy/float64(distance)
	normal := types.Pointf{float32(nx), float32(ny)}
	inward := -ny*float64(u.VelVec.Y) - nx*float64(u.VelVec.X)
	if inward > 0 {
		ball := geometryGameBall(u)
		// The original helper call spills and reloads the two normal components.
		nx, ny = float64(normal.X), float64(normal.Y)
		if ball == 0 {
			speed := float64(float32(inward))
			u.VelVec.X = float32(speed*nx + float64(u.VelVec.X))
			u.VelVec.Y = float32(speed*ny + float64(u.VelVec.Y))
		}
	}
	force := (float64(u.Shape.Circle.R) - float64(distance)) * float64(geometryWallForce())
	u.Sub548600(types.Pointf{float32(force * nx), float32(force * ny)})
	geometryHit(u, nil, &normal)
	return 1
}

func geometryCircleCircle(a, b *server.Object) {
	dx, dy := float64(b.NewPos.X)-float64(a.NewPos.X), float64(b.NewPos.Y)-float64(a.NewPos.Y)
	delta := types.Pointf{float32(dx), float32(dy)}
	if dx == 0 && dy == 0 {
		index := GetServer().S().Rand.Logic.IntClamp(0, 3)
		delta.X = memmap.Float32(0x587000, uintptr(292784+8*index))
		delta.Y = memmap.Float32(0x587000, uintptr(292788+8*index))
		dx, dy = float64(delta.X), float64(delta.Y)
	}
	distance64 := math.Sqrt(dy*dy + dx*dx)
	distance := float32(distance64)
	if distance64 == 0 {
		distance = 0.0099999998
	}
	overlap64 := float64(a.Shape.Circle.R) + float64(b.Shape.Circle.R) - float64(distance)
	overlap := float32(overlap64)
	if !(overlap64 > 0) {
		return
	}
	if a.ObjClass&0x2204 != 0 && b.ObjClass&0x2204 != 0 && !GetServer().S().MapTraceRayAt(a.PosVec, b.PosVec, nil, nil, 0) {
		return
	}
	geometryHit(b, a, &delta)
	if a.ObjFlags&8 == 0 && b.ObjFlags&8 == 0 && (a.ObjClass&6 == 0 || b.ObjFlags&0x2000 == 0) {
		normal := types.Pointf{delta.X / distance, delta.Y / distance}
		vx, vy := float64(a.VelVec.X)-float64(b.VelVec.X), float64(a.VelVec.Y)-float64(b.VelVec.Y)
		mass := float64(a.Mass)
		if b.Mass <= a.Mass {
			mass = float64(b.Mass)
		}
		spring := float64(geometryObjectForce()) * float64(overlap)
		fx, fy := -spring*float64(normal.X), -spring*float64(normal.Y)
		if a.ObjClass&6 == 0 || b.ObjClass&6 == 0 {
			tangentX := -float64(normal.Y)
			tangentSpeed := tangentX*vx + vy*float64(normal.X)
			speed := float32(tangentSpeed)
			fx = float64(float32(fx - tangentSpeed*tangentX*mass*0.69999999))
			fy = float64(float32(fy - float64(speed)*mass*float64(normal.X)*0.69999999))
		}
		a.Sub548600(types.Pointf{float32(fx), float32(fy)})
	}
	geometryWake(a)
	geometryWake(b)
}
func geometryRotated(p types.Pointf) types.Pointf {
	return types.Pointf{float32((float64(p.X) + float64(p.Y)) * 0.70710677), float32((float64(p.Y) - float64(p.X)) * 0.70710677)}
}
func geometryRotatedBounds(u *server.Object, p types.Pointf) [4]float32 {
	w := float64(u.Shape.Box.W) * 0.5
	h := float32(float64(u.Shape.Box.H) * 0.5)
	return [4]float32{float32(float64(p.X) - w), p.Y - h, float32(w + float64(p.X)), h + p.Y}
}
func geometryBoxBox(a, b *server.Object) {
	centerA, centerB := geometryRotated(a.NewPos), geometryRotated(b.NewPos)
	boxA, boxB := geometryRotatedBounds(a, centerA), geometryRotatedBounds(b, centerB)
	if !(boxA[0] <= boxB[2] && boxA[1] <= boxB[3] && boxA[2] >= boxB[0] && boxA[3] >= boxB[1]) {
		return
	}
	delta := types.Pointf{b.NewPos.X - a.NewPos.X, b.NewPos.Y - a.NewPos.Y}
	geometryHit(b, a, &delta)
	if a.ObjFlags&8 == 0 && b.ObjFlags&8 == 0 && (a.ObjClass&6 == 0 || b.ObjFlags&0x2000 == 0) {
		loX := boxA[0]
		if boxA[0] <= boxB[0] {
			loX = boxB[0]
		}
		hiX := boxA[2]
		if boxA[2] >= boxB[2] {
			hiX = boxB[2]
		}
		loY := boxA[1]
		if boxA[1] <= boxB[1] {
			loY = boxB[1]
		}
		hiY := boxA[3]
		if boxA[3] >= boxB[3] {
			hiY = boxB[3]
		}
		width, height := float64(hiX)-float64(loX), float64(hiY)-float64(loY)
		var x float64
		var y float64
		if width >= height {
			y = height
			if centerA.Y < centerB.Y {
				y = -height
			}
		} else {
			x = float64(width)
			if centerA.X < centerB.X {
				x = -x
			}
		}
		fx, fy := x*float64(geometryObjectForce()), y*float64(geometryObjectForce())
		a.Sub548600(types.Pointf{float32((fx - fy) * 0.70710677), float32((fy + fx) * 0.70710677)})
	}
	geometryWake(a)
	geometryWake(b)
}
