package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func motionWaypoint(raw uint32) *server.Waypoint {
	return (*server.Waypoint)(memmap.WordPtr(raw))
}
func motionWaypointAddress(w *server.Waypoint) uint32 { return memmap.PtrWord(unsafe.Pointer(w)) }
func motionMover(u *server.Object) {
	s := GetServer().S()
	d := u.UpdateDataMover()
	state := &d.Field_0
	if d.Field_8 == 0 {
		s.Objs.RemoveFromUpdatable(u)
		return
	}
	if d.Field_7 == nil {
		d.Field_7 = s.Objs.GetObjectByInd(int(d.Field_8))
		if d.Field_7 == nil {
			s.Objs.RemoveFromUpdatable(u)
			return
		}
	}
	if d.NextInd != 0 && d.Next == nil {
		d.Next = s.WPs.ByInd(int(d.NextInd))
	}
	if d.PrevInd != 0 && d.Prev == nil {
		d.Prev = s.WPs.ByInd(int(d.PrevInd))
	}
	target := d.Field_7
	if target.ObjFlags&4 == 0 || target.ObjFlags&0x20 != 0 {
		d.Field_7 = nil
		s.Objs.RemoveFromUpdatable(u)
		return
	}
	switch *state {
	case 0:
		if u.ObjFlags&0x1000000 != 0 {
			wp := s.WPs.ByInd(int(d.Field_2))
			if wp != nil {
				Nox_xxx_unitMove_4E7010(u, d.Field_7.PosVec)
				speed := float32(float64(int32(math.Float32bits(d.Field_1))) * .25)
				*state = 1
				d.Next = wp
				d.Prev = nil
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
			wp := d.Next
			dot := (float64(wp.PosVec.Y)-float64(u.PosVec.Y))*float64(u.VelVec.Y) + (float64(wp.PosVec.X)-float64(u.PosVec.X))*float64(u.VelVec.X)
			if dot <= 0 {
				if wp.PointsCnt == 1 {
					d.Prev = d.Next
					d.Next = wp.Points[0].Waypoint
				} else if wp.PointsCnt != 0 {
					for {
						i := s.Rand.Logic.IntClamp(0, int(wp.PointsCnt)-1)
						wp = d.Next
						next := wp.Points[i].Waypoint
						if next != d.Prev {
							d.Prev = d.Next
							d.Next = next
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
			Nox_xxx_unitMove_4E7010(d.Field_7, u.PosVec)
			wp := d.Next
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
