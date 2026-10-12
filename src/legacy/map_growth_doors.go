package legacy

import "github.com/brutusmcforce/opennox/libs/types"

func mapGrowthDoors() *mapRoom {
	var previous *mapRoom
	var run, start int32
	for r := mapRoomHead(); r != nil; r = r.Next {
		if r.Kind != 1 {
			continue
		}
		mapPopulationProgress(156)
		for dir := int32(0); dir < 4; dir++ {
			if r.Counts[dir] == 0 && mapRoomRandomInt(1, 100) <= int32(*populationBlob(1549848)) {
				pos := r.Grid
				length := r.Width
				switch dir {
				case 0:
					pos[1]--
				case 1:
					pos[1] += r.Height
				case 2:
					pos[0] += r.Width
					length = r.Height
				case 3:
					pos[0]--
					length = r.Height
				}
				for i := int32(0); i < length; i++ {
					next := mapRoomAt(&pos)
					if next != nil && next == previous && next.Flags&2 != 0 {
						run++
					}
					if next != previous || i == length-1 {
						if previous != nil && run >= 3 {
							mapGrowthDoorOpening(r, previous, dir, (start+i)/2, run)
						}
						start = i
						run = 1
						previous = next
						if next != nil && next.Kind == 1 && (dir == 1 || dir == 3) {
							previous = nil
						}
					}
					if dir < 2 {
						pos[0]++
					} else {
						pos[1]++
					}
				}
			}
			previous = nil
		}
	}
	return nil
}

// Door pairs share the same edge scan and exclusion/waypoint bookkeeping. Keep
// double intermediates until each original float store, including the pair center.
func mapGrowthDoorOpening(r, n *mapRoom, dir, mid, run int32) {
	var p types.Pointf
	tangent, normal := &p.X, &p.Y
	orientation, secondOrientation := int32(5), int32(3)
	far := dir == 1 || dir == 2
	if dir < 2 {
		p.X = float32(float64(mid)*32.526913 + float64(r.Min.X))
		p.Y = r.Min.Y
		if far {
			p.Y = r.Max.Y
		}
	} else {
		tangent, normal = &p.Y, &p.X
		orientation, secondOrientation = 7, 1
		p.Y = float32(float64(mid)*32.526913 + float64(r.Min.Y))
		p.X = r.Min.X
		if far {
			p.X = r.Max.X
		}
	}
	place := func(pair bool, orientation int32) {
		id := "ArchedDoor"
		if pair {
			id = "ArchedHalfDoor"
		}
		*mapPaintGlobal(paintObjectType) = uint32(GetServer().S().Types.IndByID(id))
		if u := mapPaintPlaceObject(&p); u != nil {
			mapPaintOrientObject(u, orientation)
		}
	}
	exclude := func(a types.Pointf) {
		axis := &a.Y
		if dir >= 2 {
			axis = &a.X
		}
		if far {
			*axis = float32(float64(*axis) - 32.526913)
		}
		mapRoomAddExclusion(r, &a, 32.526913, 32.526913)
		if mapRoomIsHall(n) != 0 {
			delta := -32.526913
			if far {
				delta = 32.526913
			}
			*axis = float32(float64(*axis) + delta)
			mapRoomAddExclusion(n, &a, 32.526913, 32.526913)
		}
	}
	mapPaintEraseWall(&p)
	*tangent = float32(float64(*tangent) - 16.263456)
	place(run >= 4, orientation)
	exclude(p)
	center := float64(*tangent) + 16.263456
	if run >= 4 {
		*tangent = float32(center + 32.526913)
		mapPaintEraseWall(&p)
		*tangent = float32(float64(*tangent) + 16.263456)
		place(true, secondOrientation)
		center = float64(float32(float64(*tangent) - 32.526913))
		a := p
		if dir < 2 {
			a.X = float32(center)
		} else {
			a.Y = float32(center)
		}
		exclude(a)
	}
	var a, b types.Pointf
	if dir < 2 {
		a = types.Ptf(float32(center), float32(float64(*normal)-32.526913))
		b = types.Ptf(float32(center), float32(float64(*normal)+32.526913))
	} else {
		a = types.Ptf(float32(float64(*normal)-32.526913), float32(center))
		b = types.Ptf(float32(float64(*normal)+32.526913), float32(center))
	}
	populationEnsureWaypoint(&a)
	populationEnsureWaypoint(&b)
	populationWaypointConnect(&a, &b)
	populationWaypointConnect(&b, &a)
	if far {
		mapRoomRememberPoint(r, &a)
		if n.Kind == 1 {
			mapRoomRememberPoint(n, &b)
		}
	} else {
		mapRoomRememberPoint(r, &b)
		if n.Kind == 1 {
			mapRoomRememberPoint(n, &a)
		}
	}
	mapRoomConnect(r, n, dir)
	mapRoomConnect(n, r, mapRoomOpposite(dir))
}
