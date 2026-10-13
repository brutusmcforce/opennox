package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

func pathDebug(s string) {
	if noxflags.HasEngine(noxflags.EngineShowAI) {
		ai.Log.Printf("%s", s)
	}
}
func pathTakeStatus() uint32                                  { p := memmap.PtrUint32(0x5D4594, 2490500); v := *p; *p = 0; return v }
func pathWaypointWords(ud *server.MonsterUpdateData) []uint32 { return unsafe.Slice(&ud.Field75, 16) }

// Output includes one compatibility slot beyond capacity. The original builder
// writes that slot before checking its limit; the movement owner aliases it
// with Field91 and resets Field91 immediately after construction.
func pathBuildGraph(start, end *server.Waypoint, output []uint32) int {
	flag := memmap.PtrUint32(0x5D4594, 2490500)
	if waypointEnabledMask(start, 128) && waypointEnabledMask(end, 128) {
		dword_5d4594_2490504++
		epoch := uint32(dword_5d4594_2490504)
		start.Field15, start.Field16, start.Field14 = 0, 0, epoch
		for frontier := start; frontier != nil; {
			var next *server.Waypoint
			for cur := frontier; cur != nil; cur = roamWaypoint(cur.Field16) {
				if cur == end {
					*flag = 0
					scratch := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2489476), 256)
					count := 0
					for node := cur; node != nil && count < len(scratch); node = roamWaypoint(node.Field15) {
						scratch[count] = roamWaypointWord(node)
						count++
					}
					if count == len(scratch) {
						pathDebug("BuildWaypointPath: Node list exceeded internal buffer.\n")
					}
					copied := 0
					for copied < count {
						output[copied] = scratch[count-1-copied]
						if copied == len(output)-1 {
							break
						}
						copied++
					}
					if copied != count {
						pathDebug("BuildWaypointPath: Node list too long.\n")
						*flag = 1
					}
					return copied
				}
				for i := 0; i < int(cur.PointsCnt); i++ {
					node := cur.Points[i].Waypoint
					if node.Field14 != epoch && waypointEnabledMask(node, 128) {
						node.Field15 = roamWaypointWord(cur)
						node.Field16 = roamWaypointWord(next)
						next = node
						node.Field14 = epoch
					}
				}
			}
			frontier = next
		}
	}
	*flag = 2
	return 0
}
func pathBuildBetween(u *server.Object, target *types.Pointf) int {
	ud := u.UpdateDataMonster()
	start := GetServer().Sub_50CB20(u, &u.PosVec)
	end := GetServer().Sub_50CB20(u, target)
	if start == nil || end == nil || start == end {
		return 0
	}
	return pathBuildGraph(start, end, unsafe.Slice(&ud.Field75, 17))
}
func pathSetWaypoints(u *server.Object, target *types.Pointf) int {
	ud := u.UpdateDataMonster()
	ud.Field92 = math.Float32bits(target.X)
	ud.Field93 = math.Float32bits(target.Y)
	count := pathBuildBetween(u, target)
	ud.Field74 = uint32(count)
	ud.Field91 = 0
	return count
}
func pathFollowWaypoints(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	if ud.Field74 != 0 && ud.Field2 == 0 {
		index := ud.Field91
		wp := roamWaypoint(pathWaypointWords(ud)[index])
		dx, dy := float64(wp.PosVec.X)-float64(u.PosVec.X), float64(wp.PosVec.Y)-float64(u.PosVec.Y)
		if dy*dy+dx*dx < 64 {
			if index == ud.Field74-1 {
				ud.Field74 = 0
				return true
			}
			ud.Field91 = index + 1
			wp = roamWaypoint(pathWaypointWords(ud)[index+1])
		}
		GetServer().Nox_xxx_creatureSetDetailedPath_50D220(u, &wp.PosVec)
		// Successful detailed-path generation appends at least its target point.
		if byte(ud.Field71) == 0 {
			ud.Path[ud.Field2-1] = wp.PosVec
		}
	}
	return pathActuallyMove(u) && byte(ud.Field71) == 2
}
func pathActuallyMove(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	if ud.Field2 == 0 {
		return false
	}
	index := int32(ud.Field67)
	if index >= int32(ud.Field2) {
		ud.Field2 = 0
		return false
	}
	origin := u.PosVec
	chosen, near := -1, -1
	best := float32(10000000)
	for i := int(index); i < int(int32(ud.Field2)); i++ {
		p := ud.Path[i]
		if !GetServer().S().MapTraceRayAt(origin, p, nil, nil, 132) {
			continue
		}
		dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
		dist := dy*dy + dx*dx
		if dist > 64 {
			if chosen < 0 || float64(best) > dist {
				best = float32(dist)
				chosen = i
			}
		} else {
			if uint32(i) == ud.Field2-1 {
				ud.Field2 = 0
				return true
			}
			near = i
		}
	}
	if chosen < 0 {
		chosen = near
	}
	if chosen < 0 {
		ud.Field2 = 0
		return false
	}
	*memmap.PtrUint32(0x5D4594, 2386204) = uint32(chosen)
	ud.Field67 = uint32(chosen)
	p := ud.Path[chosen]
	dx, dy := float64(p.X)-float64(u.PosVec.X), float64(p.Y)-float64(u.PosVec.Y)
	// x87 retains both deltas for length, then reloads X as float32 and Y as
	// float64 after direction conversion. The denominator is a float32 spill.
	denom := float32(math.Sqrt(dy*dy+dx*dx) + memmap.Float64(0x581450, 10288))
	var direction types.Pointf
	if chosen <= 0 {
		direction = ud.Path[1].Sub(ud.Path[0])
	} else {
		direction = ud.Path[chosen].Sub(ud.Path[chosen-1])
	}
	angle := int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&direction)))))
	u.Direction1, u.Direction2 = server.Dir16(angle), server.Dir16(angle)
	speed := float64(u.SpeedCur)
	if ud.StatusFlags&0x4000 != 0 {
		speed *= float64(ud.MonsterDef.RunMultiplier96)
	}
	u.ForceVec.X = float32(speed * float64(float32(dx)) / float64(denom))
	u.ForceVec.Y = float32(speed * dy / float64(denom))
	return false
}
func pathSetMove(u *server.Object) bool {
	ud := u.UpdateDataMonster()
	pos := ud.AIStackHead().ArgPos(0)
	target := &pos
	dx, dy := float64(target.X)-float64(u.PosVec.X), float64(target.Y)-float64(u.PosVec.Y)
	if math.Sqrt(dy*dy+dx*dx)+memmap.Float64(0x581450, 10288) <= 8 {
		return true
	}
	core := GetServer().S()
	flags := server.MapTraceFlags((uint32(u.ObjFlags) >> 12) & 4)
	if !core.MapTraceRayAt(u.PosVec, *target, nil, nil, flags) {
		if ud.Field2 == 0 {
			dx, dy := float64(math.Float32frombits(ud.Field92))-float64(target.X), float64(math.Float32frombits(ud.Field93))-float64(target.Y)
			if ud.Field74 == 0 || uint32(core.Frame()-ud.Field70) > 10 && dy*dy+dx*dx > 10000 {
				pathSetWaypoints(u, target)
			}
			if ud.Field74 != 0 {
				if pathFollowWaypoints(u) {
					ud.Field2, ud.Field74 = 0, 0
					return true
				}
				return false
			}
			pathDebug(" ** Waypoint path failed, doing detailed path\n")
			GetServer().Nox_xxx_creatureSetDetailedPath_50D220(u, target)
		}
	} else {
		dx, dy := float64(ud.Field68.X)-float64(target.X), float64(ud.Field68.Y)-float64(target.Y)
		if ud.Field2 == 0 || uint32(core.Frame()-ud.Field70) > 10 && dy*dy+dx*dx > 2500 {
			GetServer().Nox_xxx_creatureSetDetailedPath_50D220(u, target)
		}
	}
	if ud.Field74 != 0 {
		if pathFollowWaypoints(u) {
			ud.Field2, ud.Field74 = 0, 0
			return true
		}
		return false
	}
	if ud.Field2 != 0 && pathActuallyMove(u) {
		ud.Field2 = 0
		return true
	}
	return false
}
