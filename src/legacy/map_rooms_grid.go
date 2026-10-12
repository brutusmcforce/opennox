package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"math"
	"unsafe"
)

func mapRoomRound(pos *types.Pointf, out *[2]int32) int64 {
	shift := 0.5
	if pos.X < 0 {
		shift = -0.5
	}
	out[0] = int32(int64(float64(pos.X)*0.030743772 + shift))
	shift = 0.5
	if pos.Y < 0 {
		shift = -0.5
	}
	ret := int64(float64(pos.Y)*0.030743772 + shift)
	out[1] = int32(ret)
	return ret
}
func mapRoomGridRows() **mapRoomCellData {
	return (**mapRoomCellData)(mapRoomPointer(*mapRoomGlobalWord(0)))
}
func mapRoomGridRow(x int32) *mapRoomCellData {
	return *(**mapRoomCellData)(unsafe.Add(unsafe.Pointer(mapRoomGridRows()), int(x)*4))
}
func mapRoomGridCellAt(x, y int32) *mapRoomCellData {
	return (*mapRoomCellData)(unsafe.Add(unsafe.Pointer(mapRoomGridRow(x)), int(y)*20))
}
func mapRoomCell(pos *[2]int32) *mapRoomCellData {
	x := int32(*mapRoomGlobalWord(1)) + pos[0]
	y := int32(*mapRoomGlobalWord(1)) + pos[1]
	size := int32(*mapRoomGlobalWord(2))
	if x < 0 || x >= size || y < 0 || y >= size {
		return nil
	}
	return mapRoomGridCellAt(x, y)
}
func mapRoomGridInit(config unsafe.Pointer) uint32 {
	radius := uint32(*mapRoomWord(config, 68))
	size := 2*radius + 1
	*mapRoomGlobalWord(2) = size
	*mapRoomGlobalWord(1) = radius
	table := mapRoomCalloc(size, 4)
	*mapRoomGlobalWord(0) = mapRoomRaw(table)
	if table == nil {
		return 0
	}
	if size > 0 {
		var x int32
		for {
			p := mapRoomCalloc(*mapRoomGlobalWord(2), 20)
			*mapRoomRef(table, int(x)*4) = p
			if p == nil {
				return 0
			}
			x++
			if x >= int32(*mapRoomGlobalWord(2)) {
				break
			}
		}
	}
	for y := int32(0); y < int32(*mapRoomGlobalWord(2)); y++ {
		for x := int32(0); x < int32(*mapRoomGlobalWord(2)); x++ {
			cell := mapRoomGridCellAt(x, y)
			cell.X = x - int32(*mapRoomGlobalWord(1))
			cell.Y = y - int32(*mapRoomGlobalWord(1))
			cell.Room = nil
		}
	}
	return 1
}
func mapRoomGridFree() {
	for x := int32(0); x < int32(*mapRoomGlobalWord(2)); x++ {
		mapRoomRelease(unsafe.Pointer(mapRoomGridRow(x)))
	}
	mapRoomRelease(mapRoomPointer(*mapRoomGlobalWord(0)))
}
func mapRoomOccupy(r *mapRoom) int32 {
	var pos [2]int32
	mapRoomRound(&r.Pos, &pos)
	for y := int32(0); y < r.Height; y++ {
		for x := int32(0); x < r.Width; x++ {
			at := [2]int32{pos[0] + x, pos[1] + y}
			if cell := mapRoomCell(&at); cell != nil {
				cell.Room = r
			}
		}
	}
	return r.Height
}
func mapRoomVacate(r *mapRoom) int32 {
	var pos [2]int32
	mapRoomRound(&r.Pos, &pos)
	for y := int32(0); y < r.Height; y++ {
		for x := int32(0); x < r.Width; x++ {
			at := [2]int32{pos[0] + x, pos[1] + y}
			if cell := mapRoomCell(&at); cell != nil {
				cell.Room = nil
			}
		}
	}
	return r.Height
}
func mapRoomOverlap(r *mapRoom) *mapRoom {
	var pos [2]int32
	mapRoomRound(&r.Pos, &pos)
	for y := int32(0); y < r.Height; y++ {
		for x := int32(0); x < r.Width; x++ {
			at := [2]int32{pos[0] + x, pos[1] + y}
			if cell := mapRoomCell(&at); cell != nil && cell.Room != nil {
				return cell.Room
			}
		}
	}
	return nil
}
func mapRoomAt(pos *[2]int32) *mapRoom {
	if cell := mapRoomCell(pos); cell != nil {
		return cell.Room
	}
	return nil
}
func mapRoomResolveOverlap(r, other *mapRoom) uint32 {
	for attempt := 0; attempt < 25; attempt++ {
		distance := [4]int32{other.Grid[1] + other.Height - r.Grid[1], r.Height + r.Grid[1] - other.Grid[1], other.Grid[0] + other.Width - r.Grid[0], r.Width + r.Grid[0] - other.Grid[0]}
		best, side := int32(99999), -1
		for i, v := range distance {
			if v < best {
				best = v
				side = i
			}
		}
		var pos types.Pointf
		switch side {
		case 0:
			pos = types.Ptf(r.Pos.X, float32(float64(distance[0])*32.526913+float64(r.Pos.Y)))
		case 1:
			pos = types.Ptf(r.Pos.X, float32(float64(r.Pos.Y)-float64(distance[1])*32.526913))
		case 2:
			pos = types.Ptf(float32(float64(distance[2])*32.526913+float64(r.Pos.X)), r.Pos.Y)
		case 3:
			pos = types.Ptf(float32(float64(r.Pos.X)-float64(distance[3])*32.526913), r.Pos.Y)
		}
		mapRoomSetPos(r, &pos)
		other = mapRoomOverlap(r)
		if other == nil {
			return 1
		}
	}
	return 0
}
func mapRoomWithinBounds(config unsafe.Pointer, r *mapRoom) uint32 {
	y := float64(r.Pos.Y)
	for iy := int32(0); iy < r.Height; iy++ {
		x := float64(r.Pos.X)
		for ix := int32(0); ix < r.Width; ix++ {
			limit := float64(*mapRoomFloatWord(config, 64))
			if !(limit >= math.Abs(x) && limit >= math.Abs(y)) {
				return 0
			}
			x += 32.526913
		}
		y += 32.526913
	}
	return 1
}
func mapRoomCanPlace(config unsafe.Pointer, r *mapRoom) uint32 {
	if mapRoomWithinBounds(config, r) == 0 {
		return 0
	}
	return uint32(bool2int(mapRoomOverlap(r) == nil))
}
func mapRoomUpdateBounds(r *mapRoom) *mapRoom {
	right := float64(r.Size.X) + float64(r.Pos.X)
	r.Min = r.Pos
	r.Max.X = float32(right)
	r.Max.Y = float32(float64(r.Size.Y) + float64(r.Pos.Y))
	return r
}
func mapRoomSetPos(r *mapRoom, pos *types.Pointf) *mapRoom {
	r.Pos = *pos
	mapRoomRound(pos, &r.Grid)
	return mapRoomUpdateBounds(r)
}
