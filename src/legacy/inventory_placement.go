package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func inventoryShapeRadius(u *server.Object) float64 {
	shape := *(*uint32)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(172)) /* Shape.Kind */))
	if shape == 2 {
		return float64(*(*float32)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(176)) /* Shape.Circle.R */)))
	}
	if shape != 3 {
		return 0
	}
	w, h := float64(*(*float32)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(184)) /* Shape.Box.W */))), float64(*(*float32)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(188)) /* Shape.Box.H */)))
	if w <= h {
		return h * 0.5
	}
	return w * 0.5
}
func inventoryRandomPlacement(radius float32, origin, pos *types.Pointf) {
	step := float32(float64(radius) * 0.015625)
	angle := float32(GetServer().S().Rand.Logic.FloatClamp(float64(float32(-3.1415927)), float64(float32(3.1415927))))
	for i := 0; i < 64; i++ {
		next := float64(angle) + 1.8849558
		angle = float32(next)
		target := types.Pointf{X: float32(placementCos(next)*float64(radius) + float64(origin.X)), Y: float32(placementSin(float64(angle))*float64(radius) + float64(origin.Y))}
		if GetServer().S().MapTraceRayAt(*origin, target, nil, nil, 1) {
			*pos = target
			return
		}
		radius = float32(radius - step)
	}
	*pos = *origin
}
func inventoryTargetDrop(u, it *server.Object, pos *types.Pointf) int {
	origin := u.PosVec
	dx, dy := float64(pos.X)-float64(origin.X), float64(pos.Y)-float64(origin.Y)
	y := float32(dy)
	length := math.Sqrt(dy*float64(y) + dx*dx)
	storedLength := float32(length)
	target := *pos
	if length > 75 {
		target.X = float32(dx*75/float64(storedLength) + float64(origin.X))
		target.Y = float32(float64(y)*75/float64(storedLength) + float64(origin.Y))
	}
	if !GetServer().S().MapTraceRayAt(origin, target, nil, nil, 0) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(asObjectC(u), (*int8)(internCStr("drop.c:DropNotAllowed")), 0)
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if noxflags.HasGame(16) && uint32(it.TypeInd) == inventoryCache(1568248, "Crown") {
		return 0
	}
	return inventoryDrop(u, it, &target)
}
func inventoryForceDrop(u, it *server.Object) int {
	var p types.Pointf
	inventoryRandomPlacement(50, &u.PosVec, &p)
	return inventoryDrop(u, it, &p)
}

// The original walks a square spiral, advances it after each attempt, and
// falls back to the owner's position once an entire ring has failed.
func inventoryDropAll(u *server.Object) uint32 {
	var radius float32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		r := *(*float32)(unsafe.Add(it.CObj(), int(layoutServerObject.Off(176)) /* Shape.Circle.R */))
		if inventoryDropEligible(u, it) && float64(r) > float64(radius) {
			radius = r
		}
	}
	spacing := float32(float64(radius) + float64(radius) + 6)
	item := u.InvFirstItem
	if item == nil {
		return 0
	}
	x, y := int32(0), int32(-1)
	side, width, progress := int32(0), int32(3), int32(0)
	ringSucceeded := false
	rng := GetServer().S().Rand.Logic
	var result uint32
	for item != nil {
		nextItem := item.InvNextItem
		accepted := !inventoryDropEligible(u, item)
		if !accepted {
			limit := width - 1
			for {
				pos := types.Pointf{X: float32(float64(x)*float64(spacing) + float64(u.PosVec.X)), Y: float32(float64(u.PosVec.Y) - float64(y)*float64(spacing))}
				pos.X = float32(rng.FloatClamp(-3, 3) + float64(pos.X))
				pos.Y = float32(rng.FloatClamp(-3, 3) + float64(pos.Y))
				if GetServer().S().MapTraceRayAt(u.PosVec, pos, nil, nil, 1) {
					inventoryDrop(u, item, &pos)
					accepted = true
					ringSucceeded = true
				}
				result = uint32(limit)
				advance := true
				if progress != limit {
					progress++
				} else if side != 3 && limit != 0 {
					side++
					progress = 1
				} else {
					if !ringSucceeded {
						break
					}
					width += 2
					limit += 2
					x = 1 - width/2
					y = width / -2
					progress = 1
					side = 0
					ringSucceeded = false
					advance = false
				}
				if advance {
					switch side {
					case 0:
						x++
					case 1:
						y++
					case 2:
						x--
					case 3:
						y--
					}
				}
				if accepted {
					break
				}
			}
		}
		if !accepted {
			break
		}
		item = nextItem
		if item == nil {
			return 0
		}
	}
	for item != nil {
		next := item.InvNextItem
		result = uint32(bool2int(inventoryDropEligible(u, item)))
		if result != 0 {
			result = uint32(inventoryDrop(u, item, &u.PosVec))
		}
		item = next
	}
	return result
}

func inventoryChest(u, opener *server.Object) {
	if u == nil || opener == nil {
		return
	}
	count := 0
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		count++
	}
	if count == 0 {
		return
	}
	var dir types.Pointf
	switch {
	case u.ObjSubClass&0x100 != 0:
		dir = types.Pointf{X: -1, Y: -1}
	case u.ObjSubClass&0x200 != 0:
		dir = types.Pointf{X: 1, Y: -1}
	case u.ObjSubClass&0x400 != 0:
		dir = types.Pointf{X: 1, Y: 1}
	case u.ObjSubClass&0x800 != 0:
		dir = types.Pointf{X: -1, Y: 1}
	default:
		dir = types.Pointf{X: opener.PosVec.X - u.PosVec.X, Y: opener.PosVec.Y - u.PosVec.Y}
	}
	geometryNormalize((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&dir))))
	dist := inventoryShapeRadius(u) + 4 + 15
	var p [3]types.Pointf
	p[0] = types.Pointf{X: float32(dist*float64(dir.X) + float64(u.PosVec.X)), Y: float32(dist*float64(dir.Y) + float64(u.PosVec.Y))}
	p[1] = types.Pointf{X: float32(-float64(dir.Y)*30 + float64(p[0].X)), Y: float32(float64(dir.X)*30 + float64(p[0].Y))}
	p[2] = types.Pointf{X: float32(float64(dir.Y)*30 + float64(p[0].X)), Y: float32(-float64(dir.X)*30 + float64(p[0].Y))}
	distance := func(p types.Pointf) float64 {
		x, y := float64(p.X)-float64(opener.PosVec.X), float64(p.Y)-float64(opener.PosVec.Y)
		return y*y + x*x
	}
	d0, d1, d2 := float32(distance(p[0])), distance(p[1]), float32(distance(p[2]))
	if !(float64(d0) >= d1) {
		p[0], p[1] = p[1], p[0]
		old := d0
		d0 = float32(d1)
		d1 = float64(old)
	}
	if d1 < float64(d2) {
		p[1], p[2] = p[2], p[1]
		d1 = float64(d2)
	}
	if float64(d0) < d1 {
		p[0], p[1] = p[1], p[0]
	}
	for i := 2; i >= 0; i-- {
		if !GetServer().S().MapTraceRayAt(u.PosVec, p[i], nil, nil, 1) {
			p[i] = opener.PosVec
		}
	}
	index := 0
	for it := u.InvFirstItem; it != nil; {
		next := it.InvNextItem
		if it.Weight != 255 && it.ObjClass&2 == 0 {
			it.ObjFlags |= 0x40
			motionDeactivate(it)
			inventoryDrop(u, it, &p[index])
			index = (index + 1) % 3
		}
		it = next
	}
}
