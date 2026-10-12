package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var spatialCursorOwner uint32

func spatialCursor(u *server.Object) *server.Object {
	pl := u.UpdateDataPlayer().Player
	mouse := types.Pointf{float32(pl.CursorVec.X), float32(pl.CursorVec.Y)}
	spatialCursorOwner = motionAddress(u)
	*memmap.PtrUint32(0x5D4594, 2491596) = 0
	*memmap.PtrUint32(0x5D4594, 2491600) = 0
	GetServer().S().Map.EachObjInCircle(mouse, 100, func(t *server.Object) bool { spatialCursorCandidate(t, &mouse); return true })
	return motionObject(memmap.Uint32(0x5D4594, 2491596))
}
func spatialCursorCandidate(u *server.Object, mouse *types.Pointf) {
	s := GetServer().S()
	if memmap.Uint32(0x5D4594, 2491604) == 0 {
		*memmap.PtrUint32(0x5D4594, 2491604) = uint32(s.Types.IndByID("Polyp"))
	}
	owner := motionObject(spatialCursorOwner)
	if u == owner || u.ObjFlags&0x8020 != 0 || u.HasEnchant(0) && !owner.HasEnchant(21) {
		return
	}
	if u.ObjClass&0x80000206 == 0 && uint32(u.TypeInd) != memmap.Uint32(0x5D4594, 2491604) {
		return
	}
	if !s.MapTraceVision(u, owner) {
		return
	}
	if u.ObjClass&4 != 0 {
		if u.NetCode == uint32(nox_player_netCode_85319C) && noxflags.HasEngine(noxflags.EngineNoRendering) {
			return
		}
		pl := s.Players.ByID(int(u.NetCode))
		if pl == nil || pl.Field3680&1 != 0 {
			return
		}
	}
	if u.ObjClass&0x200 != 0 && u.ObjFlags&0x4000 == 0 {
		return
	}
	// The compiled C spills projected Y before shape tests, and the score before
	// comparing against the previous float score. Preserve ties after rounding.
	y := float32(float64(u.PosVec.Y) - float64(u.ZVal))
	inside := false
	switch u.Shape.Kind {
	case server.ShapeKindCircle:
		r := float64(u.Shape.Circle.R)
		radius2 := floatToInt32(float32(r * r))
		if mouse.Y > y || mouse.Y < y {
			dx, dy := float64(mouse.X)-float64(u.PosVec.X), float64(mouse.Y)-float64(y)
			inside = !(float64(radius2) <= dy*dy+dx*dx)
		} else {
			inside = float64(mouse.X) > float64(u.PosVec.X)-r && float64(mouse.X) < float64(u.PosVec.X)+r
		}
	case server.ShapeKindBox:
		center := types.Pointf{u.PosVec.X, y}
		inside = collisionContains(&center, (*[11]float32)(unsafe.Pointer(&u.Shape)), mouse)
	}
	if !inside {
		return
	}
	score := float32(float64(u.ZVal) + float64(u.PosVec.Y))
	if score > memmap.Float32(0x5D4594, 2491600) {
		*memmap.PtrFloat32(0x5D4594, 2491600) = score
		*memmap.PtrUint32(0x5D4594, 2491596) = motionAddress(u)
	}
}
