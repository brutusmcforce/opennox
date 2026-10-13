package legacy

import (
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func mapPaintSelectObject(name *byte) uint32 {
	if name == nil {
		return 0
	}
	id := GoStringP(unsafe.Pointer(name))
	if strings.EqualFold(id, "NONE") {
		*mapPaintGlobal(paintObjectType) = 0
	} else {
		*mapPaintGlobal(paintObjectType) = uint32(GetServer().S().Types.IndByID(id))
	}
	return 1
}
func mapPaintPlaceObject(pos *types.Pointf) *server.Object {
	ind := *mapPaintGlobal(paintObjectType)
	if ind == 0 {
		return nil
	}
	var world types.Pointf
	if mapPaintTransform(pos, &world) == 0 {
		return nil
	}
	u := GetServer().S().NewObjectByTypeInd(int(ind))
	if u != nil {
		return mapPaintMoveObject(u, pos)
	}
	return nil
}
func mapPaintMoveObject(u *server.Object, pos *types.Pointf) *server.Object {
	var world types.Pointf
	if mapPaintTransform(pos, &world) == 0 {
		return u
	}
	if u == nil {
		return nil
	}
	u.Extent = *mapPaintGlobal(paintObjectCounter)
	*mapPaintGlobal(paintObjectCounter)++
	u.PosVec = world
	u.ObjFlags |= 0x1000000
	if GetServer().S().Types.ByInd(int(u.TypeInd)).Xfer == xferIdentityKey(xferIDDoor) {
		u.PosVec.X = float32(int32(23 * int64(float64(world.X)*0.043478262+0.5)))
		u.PosVec.Y = float32(int32(23 * int64(float64(world.Y)*0.043478262+0.5)))
	}
	GetServer().CreateObjectAt(u, nil, u.PosVec)
	GetServer().S().Objs.ObjectsClearPending()
	return u
}
func mapPaintOrientObject(u *server.Object, dir int32) uint32 {
	if u == nil {
		return 0
	}
	if u.ObjClass&2 != 0 {
		angle := uint32(geometryDirection4Angle(int32(mapPaintDirection(dir))))
		layout.As[server.MonsterUpdateData](u.UpdateData).Direction94 = angle
		u.Direction1 = server.Dir16(angle)
		return 1
	}
	if u.ObjClass&0x80 == 0 {
		return 0
	}
	angle := uint32(0)
	switch dir {
	case 1:
		angle = 0
	case 3:
		angle = 24
	case 5:
		angle = 8
	case 7:
		angle = 16
	default:
		return 0
	}
	data := u.UpdateDataDoor()
	data.Dir = angle
	data.BaseDir = angle
	data.PrevDir = angle
	return 1
}
func mapPaintFinishBook(u *server.Object, value byte) uint32 {
	if u == nil || GetServer().S().Types.ByInd(int(u.TypeInd)).Xfer != xferIdentityKey(xferIDSpellReward) {
		return 0
	}
	u.UseDataSpellReward().Spell = value
	return 1
}
func mapPaintDoor(r *mapRoom, pos *types.Pointf, span int32, vertical, double bool) uint32 {
	off := 100
	if double {
		off = 160
	}
	name := unsafe.Add(r.Decoration, off)
	if *(*byte)(name) == 0 {
		return 0
	}
	if vertical {
		mapPaintWallLineY(pos, span+1)
	} else {
		mapPaintWallLineX(pos, span+1)
	}
	p := *pos
	axis := &p.X
	if vertical {
		axis = &p.Y
	}
	*axis = float32(float64(span/2)*32.526913 + float64(*axis))
	mapPaintEraseWall(&p)
	delta := 16.263456
	if double {
		*axis = float32(float64(*axis) + 32.526913)
		mapPaintEraseWall(&p)
		delta = 48.790367
	}
	*axis = float32(float64(*axis) - delta)
	mapPaintSelectObject((*byte)(name))
	dir := int32(5)
	if vertical {
		dir = 7
	}
	if u := mapPaintPlaceObject(&p); u != nil {
		mapPaintOrientObject(u, dir)
	}
	if double {
		*axis = float32(float64(*axis) + 65.053825)
		dir = 3
		if vertical {
			dir = 1
		}
		if u := mapPaintPlaceObject(&p); u != nil {
			mapPaintOrientObject(u, dir)
		}
	}
	return 1
}
func mapPaintDoorX(r *mapRoom, p *types.Pointf, n int32) uint32 {
	return mapPaintDoor(r, p, n, false, false)
}
func mapPaintDoubleDoorX(r *mapRoom, p *types.Pointf, n int32) uint32 {
	return mapPaintDoor(r, p, n, false, true)
}
func mapPaintDoorY(r *mapRoom, p *types.Pointf, n int32) uint32 {
	return mapPaintDoor(r, p, n, true, false)
}
func mapPaintDoubleDoorY(r *mapRoom, p *types.Pointf, n int32) uint32 {
	return mapPaintDoor(r, p, n, true, true)
}
func mapPaintRoomDoors(config unsafe.Pointer, r *mapRoom) {
	for dir := 0; dir < 4; dir++ {
		for i := 0; i < int(r.Counts[dir]); i++ {
			mapPaintJoinDoors(config, r, r.Neighbors[dir][i], int32(dir))
		}
	}
}
func mapPaintJoinDoors(_ unsafe.Pointer, r, t *mapRoom, dir int32) {
	var p types.Pointf
	switch dir {
	case 0, 1:
		p.X = max(r.Min.X, t.Min.X)
		p.Y = r.Min.Y
		if dir == 1 {
			p.Y = r.Max.Y
		}
		n := min(r.Width, t.Width)
		if n >= 2 && (n == 2 || mapPaintDoubleDoorX(r, &p, n) == 0) {
			mapPaintDoorX(r, &p, n)
		}
	case 2, 3:
		p.X = r.Max.X
		if dir == 3 {
			p.X = r.Min.X
		}
		p.Y = max(r.Min.Y, t.Min.Y)
		n := min(r.Height, t.Height)
		if n >= 2 && (n == 2 || mapPaintDoubleDoorY(r, &p, n) == 0) {
			mapPaintDoorY(r, &p, n)
		}
	}
}
