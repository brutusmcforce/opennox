package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

func spatialPredict(a, b *server.Object, speed float32, out *types.Pointf) *server.Object {
	dx, dy := float64(b.PosVec.X)-float64(a.PosVec.X), float64(b.PosVec.Y)-float64(a.PosVec.Y)
	scale := math.Sqrt(dy*dy+dx*dx) / float64(speed)
	out.X = float32(scale*float64(b.VelVec.X) + float64(b.PosVec.X))
	out.Y = float32(scale*float64(b.VelVec.Y) + float64(b.PosVec.Y))
	return b
}
func spatialEligible(a, b *server.Object) bool {
	if b.ObjClass&1 != 0 {
		return false
	}
	af, bf := a.ObjFlags, b.ObjFlags
	if af&0x20 != 0 || bf&0x60 != 0 || a.Collide == nil || b.Collide == nil {
		return false
	}
	if bf&0x80 != 0 {
		return true
	}
	if af&0x11 != 0 && bf&0x4000 != 0 || bf&0x11 != 0 && af&0x4000 != 0 {
		return false
	}
	if (af|bf)&0x400 != 0 && b.TeamVal.SameAs(&a.TeamVal) {
		return false
	}
	if owner := a.ObjOwner; owner != nil && a.ObjClass&1 != 0 && a.ObjSubClass&2 == 0 && owner.ObjClass&2 != 0 && b.ObjClass&2 != 0 {
		if !GetServer().S().IsEnemyTo(owner, b) || b.TeamVal.SameAs(&owner.TeamVal) {
			return false
		}
	}
	return true
}
func spatialTeamEligible(a, b *server.Object) bool {
	return spatialEligible(b, a) && (!a.TeamVal.SameAs(&b.TeamVal) || noxflags.HasGamePlay(1))
}
func spatialCandidate(a, b *server.Object, next, prev *types.Pointf, hit **server.Object) {
	if int8(b.ObjClass) >= 0 {
		if spatialEligible(a, b) && collisionObjectContains(b, next) != 0 {
			*hit = b
		}
		return
	}
	if b.ObjSubClass&4 != 0 {
		return
	}
	dir := b.UpdateDataDoor().Dir
	ray := [4]float32{prev.X, prev.Y, next.X, next.Y}
	door := [4]float32{b.PosVec.X, b.PosVec.Y, float32(float64(memmap.Int32(0x587000, uintptr(196184+8*dir))) + float64(b.PosVec.X)), float32(float64(memmap.Int32(0x587000, uintptr(196188+8*dir))) + float64(b.PosVec.Y))}
	if server.LineTraceXxx(types.Rectf{Min: types.Pointf{ray[0], ray[1]}, Max: types.Pointf{ray[2], ray[3]}}, types.Rectf{Min: types.Pointf{door[0], door[1]}, Max: types.Pointf{door[2], door[3]}}) {
		*next = *prev
		*hit = b
	}
}
func spatialProbe(a *server.Object, next, prev *types.Pointf) *server.Object {
	var hit *server.Object
	GetServer().S().Map.Sub517B70(*next, func(b *server.Object) { spatialCandidate(a, b, next, prev, &hit) })
	return hit
}
func spatialRay(ray *[4]float32) bool {
	return GetServer().S().MapTraceRay(types.Pointf{ray[0], ray[1]}, types.Pointf{ray[2], ray[3]}, 9)
}
