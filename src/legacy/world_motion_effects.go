package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func motionRadial(p *types.Pointf, radius float32, callback func(*server.Object)) uint32 {
	if p != nil {
		rect := types.Rectf{Min: types.Pointf{X: float32(float64(p.X) - float64(radius)), Y: float32(float64(p.Y) - float64(radius))}, Max: types.Pointf{X: float32(float64(p.X) + float64(radius)), Y: float32(float64(p.Y) + float64(radius))}}
		GetServer().S().Map.EachObjInRect(rect, func(u *server.Object) bool {
			motionRadialCandidate(u, p, radius, callback)
			return true
		})
	}
	return memmap.PtrWord(unsafe.Pointer(p))
}
func motionScorchInit() int32 {
	var result uint32
	for _, off := range []uintptr{276824, 276832, 276840} {
		name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, off)))
		result = uint32(GetServer().S().Types.IndByID(name))
		*memmap.PtrUint32(0x587000, off+4) = result
	}
	*memmap.PtrUint32(0x5D4594, 2488636) = 1
	return int32(result)
}
func motionScorch(p *types.Pointf, size int32) {
	s := GetServer().S()
	if memmap.Uint32(0x5D4594, 2488636) == 0 {
		motionScorchInit()
	}
	if size < 0 || size > 2 {
		return
	}
	index := s.Rand.Logic.IntClamp(0, 0)
	id := memmap.Uint32(0x587000, uintptr(276828+8*size+8*int32(index)))
	if u := s.NewObjectByTypeInd(int(id)); u != nil {
		GetServer().CreateObjectAt(u, nil, *p)
		lo, hi := 10, 20
		if noxflags.HasGame(noxflags.GameModeQuest) {
			lo, hi = 5, 8
		}
		lifetime := s.Rand.Logic.IntClamp(lo, hi)
		motionDecaySet(u, int32(uint32(s.TickRate())*uint32(lifetime)))
	}
}

func motionRadialCandidate(u *server.Object, p *types.Pointf, radius float32, callback func(*server.Object)) {
	if u == nil || p == nil {
		return
	}
	var normal types.Pointf
	var overlap float64
	if u.Shape.Kind == server.ShapeKindBox {
		overlap = collisionBoxDistance(p, radius, u, &normal)
	} else {
		dx := float64(p.X) - float64(u.PosVec.X)
		dyWide := float64(p.Y) - float64(u.PosVec.Y)
		dy := dyWide
		distance := math.Sqrt(dyWide*float64(dy) + float64(dx)*float64(dx))
		if u.Shape.Kind == server.ShapeKindCircle {
			distance -= float64(u.Shape.Circle.R)
		}
		overlap = float64(radius) - distance
	}
	if overlap > 0 {
		callback(u)
	}
}
