package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"math"
)

func motionTrace(u *server.Object, target *uint32, normal *types.Pointf) int8 {
	points := [2]types.Pointf{u.PosVec, u.NewPos}
	previous, next := &points[0], &points[1]
	dx := float64(u.NewPos.X) - float64(u.PosVec.X)
	dyWide := float64(u.NewPos.Y) - float64(u.PosVec.Y)
	dy := dyWide
	length2 := dyWide*float64(dy) + float64(dx)*float64(dx)
	var hit *server.Object
	var objectNormal types.Pointf
	probe := func() bool {
		hit = spatialProbe(u, next, previous)
		if hit == nil {
			return false
		}
		objectNormal = types.Pointf{X: float32(float64(previous.X) - float64(hit.PosVec.X)), Y: float32(float64(previous.Y) - float64(hit.PosVec.Y))}
		return true
	}
	if length2 <= 36 {
		probe()
	} else {
		count := doubleToInt32(math.Sqrt(length2*.027777778)) + 1
		*next = *previous
		// The compiled C spills X to float32 before division, but retains Y wide.
		sx, sy := float32(float64(float32(dx))/float64(count)), float32(dy/float64(count))
		for i := int32(0); i < count; i++ {
			next.X = float32(float64(next.X) + float64(sx))
			next.Y = float32(float64(next.Y) + float64(sy))
			if probe() {
				break
			}
			*previous = *next
		}
	}
	start, end := u.PosVec, u.NewPos
	var wallPoint, wallNormal types.Pointf
	var grid image.Point
	wall := false
	if !GetServer().S().MapTraceRayAt(start, end, &wallPoint, &grid, 5) {
		*memmap.PtrUint32(0x5D4594, 2488612) = uint32(grid.X)
		*memmap.PtrUint32(0x5D4594, 2488616) = uint32(grid.Y)
		dword_5d4594_2488620 = 1
		ray := [4]float32{start.X, start.Y, wallPoint.X, wallPoint.Y}
		wall = spatialNormal(&[2]int32{int32(grid.X), int32(grid.Y)}, &ray, &wallNormal) != 0
		u.NewPos = u.PosVec
	}
	if hit != nil {
		dx, dy := float64(u.PosVec.X)-float64(hit.PosVec.X), float64(u.PosVec.Y)-float64(hit.PosVec.Y)
		wx, wy := float64(u.PosVec.X)-float64(wallPoint.X), float64(u.PosVec.Y)-float64(wallPoint.Y)
		if !wall || dy*dy+dx*dx < wy*wy+wx*wx {
			*normal = objectNormal
			*target = motionAddress(hit)
			return 1
		}
	}
	if wall {
		*normal = wallNormal
		*target = 0
		return 1
	}
	return 0
}
func motionProjectileDispatch(u *server.Object) {
	if memmap.Uint32(0x5D4594, 2488624) == 0 {
		*memmap.PtrUint32(0x5D4594, 2488624) = uint32(GetServer().S().Types.IndByID("SmallFist"))
		*memmap.PtrUint32(0x5D4594, 2488628) = uint32(GetServer().S().Types.IndByID("MediumFist"))
		*memmap.PtrUint32(0x5D4594, 2488632) = uint32(GetServer().S().Types.IndByID("LargeFist"))
	}
	if u.ObjFlags&0x60 != 0 {
		return
	}
	dword_5d4594_2488620 = 0
	var raw uint32
	var normal types.Pointf
	if motionTrace(u, &raw, &normal) == 0 {
		return
	}
	t := motionObject(raw)
	if t != nil {
		for _, off := range []uintptr{2488624, 2488628, 2488632} {
			if uint32(t.TypeInd) == memmap.Uint32(0x5D4594, off) {
				return
			}
		}
	}
	u.CallCollideWith(t, &normal)
	dword_5d4594_2488620 = 0
	if t != nil {
		normal.X = -normal.X
		normal.Y = -normal.Y
		t.CallCollideWith(u, &normal)
	}
}
