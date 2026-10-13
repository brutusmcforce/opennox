package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func collisionRetained(u, v *server.Object) int32 {
	md := layout.As[server.MonsterUpdateData](u.UpdateData)
	for i := 0; i < int(md.Field543_0); i++ {
		if md.Field535[i] == v.NetCode {
			return 1
		}
	}
	return 0
}
func collisionEligible(a, b *server.Object) int32 {
	if collisionTrigger == 0 {
		types := &GetServer().S().Types
		collisionTrigger = uint32(types.IndByID("Trigger"))
		collisionPowder = uint32(types.IndByID("BlackPowder"))
		collisionHand = uint32(types.IndByID("TelekinesisHand"))
	}
	if a.Collide == collisionKey(collisionIdentityPentagram) && uint32(b.TypeInd) == collisionHand || b.Collide == collisionKey(collisionIdentityPentagram) && uint32(a.TypeInd) == collisionHand {
		return 0
	}
	ac, bc, af, bf := a.ObjClass, b.ObjClass, a.ObjFlags, b.ObjFlags
	if ac&0x80 != 0 && bf&9 != 0 || bc&0x80 != 0 && af&9 != 0 || af&0x60 != 0 || bf&0x60 != 0 || a == b || a.Collide == nil || b.Collide == nil {
		return 0
	}
	if ac&2 != 0 && af&0x4000 != 0 && bc&0xCC00 != 0 {
		return collisionRetained(a, b)
	}
	if bc&2 != 0 && bf&0x4000 != 0 && ac&0xCC00 != 0 {
		return collisionRetained(b, a)
	}
	if ac&0x2000 != 0 && uint32(b.TypeInd) == collisionPowder || bc&0x2000 != 0 && uint32(a.TypeInd) == collisionPowder {
		return 1
	}
	if af&8 != 0 && bf&8 != 0 || ac&0x80 != 0 && bf&8 != 0 || bc&0x80 != 0 && af&8 != 0 || uint32(a.TypeInd) != collisionTrigger && uint32(b.TypeInd) != collisionTrigger && (af&0x11 != 0 && bf&0x24000 != 0 || bf&0x11 != 0 && af&0x24000 != 0) || ac&0x80 != 0 && bc&0x80 != 0 || ac&0x4000 != 0 && bc&0x4000 != 0 || ac&0x8000 != 0 && bc&0x8000 != 0 || af&0x400 != 0 && Nox_xxx_unitsHaveSameTeam_4EC520(a, b) {
		return 0
	}
	return 1
}

// Arguments retain the spatial callback order: candidate first, scanning object second.
func collisionPair(candidate, u *server.Object) {
	if collisionEligible(u, candidate) == 0 {
		return
	}
	a, b := u, candidate
	switch {
	case a.ObjClass&0x4000 != 0:
		collisionElevator(a, b, 0)
	case b.ObjClass&0x4000 != 0:
		collisionElevator(b, a, 1)
	case a.ObjClass&0x8000 != 0:
		collisionShaft(a, b)
	case b.ObjClass&0x8000 != 0:
		collisionShaft(b, a)
	case a.ObjClass&0x80 != 0:
		if b.Shape.Kind == server.ShapeKindCircle {
			collisionGateCircle(a, b, 0)
		} else if b.Shape.Kind == server.ShapeKindBox {
			geometryGateBox(a, b, 0)
		}
	case a.Shape.Kind == server.ShapeKindCircle:
		if b.ObjClass&0x80 != 0 {
			collisionGateCircle(b, a, 1)
		} else if b.Shape.Kind == server.ShapeKindCircle {
			geometryCircleCircle(a, b)
		} else if b.Shape.Kind == server.ShapeKindBox {
			collisionCircleBox(a, b, 0)
		}
	case a.Shape.Kind == server.ShapeKindBox:
		if b.ObjClass&0x80 != 0 {
			geometryGateBox(b, a, 1)
		} else if b.Shape.Kind == server.ShapeKindCircle {
			collisionCircleBox(b, a, 1)
		} else if b.Shape.Kind == server.ShapeKindBox {
			geometryBoxBox(a, b)
		}
	}
}
func collisionScan(u *server.Object) {
	u.Pos24 = types.Pointf{}
	if u.ObjFlags&0x60 != 0 {
		return
	}
	GetServer().S().Map.EachObjInRect(types.Rectf{Min: u.CollideP1, Max: u.CollideP2}, func(v *server.Object) bool { collisionPair(v, u); return true })
	if u.ObjFlags&8 != 0 {
		return
	}
	if u.Shape.Kind == server.ShapeKindCircle {
		collisionCircleWalls(u)
	} else if u.Shape.Kind == server.ShapeKindBox {
		geometryBoxWalls(u)
	}
}
func collisionWallOpen(grid *[2]int32, u *server.Object) {
	s := GetServer().S()
	w := s.Walls.GetWallAtGrid(image.Pt(int(grid[0]), int(grid[1])))
	if w == nil || w.Flags4&4 == 0 || u.ObjClass&6 == 0 {
		return
	}
	ud := w.SecretData()
	if ud.State != server.SecretWallClosed || ud.Flags&2 == 0 {
		return
	}
	ud.State = server.SecretWallOpening
	ud.Progress = 0
	x, y := int32(ud.GridX)*23, int32(ud.GridY)*23
	pos := types.Pointf{float32(float64(x) + 11.5), float32(float64(y) + 11.5)}
	s.Audio.EventPos(sound.ByName(s.Walls.DefByInd(int(w.Tile1)).OpenSound()), pos, 0, 0)
}
func collisionObjectContains(u *server.Object, p *types.Pointf) int32 {
	switch u.Shape.Kind {
	case server.ShapeKindCircle:
		dx, dy := float64(u.NewPos.X)-float64(p.X), float64(u.NewPos.Y)-float64(p.Y)
		return int32(bool2int(!(dy*dy+dx*dx > float64(u.Shape.Circle.R2))))
	case server.ShapeKindBox:
		box := u.Shape.Box
		x, y := float64(u.NewPos.X), float64(u.NewPos.Y)
		px, py := float64(p.X), float64(p.Y)
		lx, ly := x+float64(box.LeftTop), y+float64(box.LeftBottom)
		if (ly-lx+px-py)*0.70709997 < 0 && (y+float64(box.LeftTop2)-(x+float64(box.LeftBottom2))+px-py)*0.70709997 > 0 && ((y+float64(box.RightBottom))+(x+float64(box.RightTop))-px-py)*0.70709997 > 0 && (ly+float64(float32(lx))-px-py)*0.70709997 < 0 {
			return 1
		}
	}
	return 0
}
