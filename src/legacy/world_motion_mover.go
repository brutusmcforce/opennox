package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func motionWaypoint(raw uint32) *server.Waypoint {
	return (*server.Waypoint)(memmap.WordPtr(raw))
}
func motionWaypointAddress(w *server.Waypoint) uint32 { return memmap.PtrWord(unsafe.Pointer(w)) }
func motionMover(u *server.Object) {
	s := GetServer().S()
	d := (*[9]uint32)(u.UpdateData)
	state := (*byte)(u.UpdateData)
	if d[8] == 0 {
		s.Objs.RemoveFromUpdatable(u)
		return
	}
	if d[7] == 0 {
		d[7] = motionAddress(s.Objs.GetObjectByInd(int(d[8])))
		if d[7] == 0 {
			s.Objs.RemoveFromUpdatable(u)
			return
		}
	}
	if d[4] != 0 && d[3] == 0 {
		d[3] = motionWaypointAddress(s.WPs.ByInd(int(d[4])))
	}
	if d[6] != 0 && d[5] == 0 {
		d[5] = motionWaypointAddress(s.WPs.ByInd(int(d[6])))
	}
	target := motionObject(d[7])
	if target.ObjFlags&4 == 0 || target.ObjFlags&0x20 != 0 {
		d[7] = 0
		s.Objs.RemoveFromUpdatable(u)
		return
	}
	switch *state {
	case 0:
		if u.ObjFlags&0x1000000 != 0 {
			wp := s.WPs.ByInd(int(d[2]))
			if wp != nil {
				Nox_xxx_unitMove_4E7010(u, motionObject(d[7]).PosVec)
				speed := float32(float64(int32(d[1])) * .25)
				*state = 1
				d[3] = motionWaypointAddress(wp)
				d[5] = 0
				u.SpeedCur, u.SpeedBase = speed, speed
			} else {
				*state = 3
			}
		}
	case 1:
		if u.ObjFlags&0x1000000 == 0 {
			*state = 2
			return
		}
		if u.VelVec.X != 0 && u.VelVec.Y != 0 {
			wp := motionWaypoint(d[3])
			dot := (float64(wp.PosVec.Y)-float64(u.PosVec.Y))*float64(u.VelVec.Y) + (float64(wp.PosVec.X)-float64(u.PosVec.X))*float64(u.VelVec.X)
			if dot <= 0 {
				if wp.PointsCnt == 1 {
					d[5] = d[3]
					d[3] = motionWaypointAddress(wp.Points[0].Waypoint)
				} else if wp.PointsCnt != 0 {
					for {
						i := s.Rand.Logic.IntClamp(0, int(wp.PointsCnt)-1)
						wp = motionWaypoint(d[3])
						next := motionWaypointAddress(wp.Points[i].Waypoint)
						if next != d[5] {
							d[5] = d[3]
							d[3] = next
							break
						}
					}
				} else {
					*state = 3
					Nox_xxx_unitMove_4E7010(target, wp.PosVec)
				}
			}
		}
		if *state == 1 {
			Nox_xxx_unitMove_4E7010(motionObject(d[7]), u.PosVec)
			wp := motionWaypoint(d[3])
			dx := float64(wp.PosVec.X) - float64(u.PosVec.X)
			dyWide := float64(wp.PosVec.Y) - float64(u.PosVec.Y)
			dy := dyWide
			den := float32(math.Sqrt(dyWide*float64(dy)+dx*dx) + .1)
			u.VelVec.X = float32(dx * float64(u.SpeedCur) / float64(den))
			if u.VelVec.X == 0 {
				u.VelVec.X = math.Float32frombits(0x00800000)
			}
			u.VelVec.Y = float32(float64(dy) * float64(u.SpeedCur) / float64(den))
			if u.VelVec.Y == 0 {
				u.VelVec.Y = math.Float32frombits(0x00800000)
			}
		}
	case 2:
		if u.ObjFlags&0x1000000 != 0 {
			Nox_xxx_unitMove_4E7010(u, target.PosVec)
			*state = 1
		}
	case 3:
		s.Objs.RemoveFromUpdatable(u)
	}
}
