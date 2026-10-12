package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

var geometryCircleWallSpans = [...]struct {
	positive         bool
	lo, hi           float32
	negative         bool
	otherLo, otherHi float32
}{
	{false, 0, 0, true, 0, 23}, {true, 0, 23, false, 0, 0}, {true, 0, 23, true, 0, 23},
	{true, 0, 11.5, true, 0, 23}, {true, 0, 23, true, 11.5, 23}, {true, 11.5, 23, true, 0, 23},
	{true, 0, 23, true, 0, 11.5}, {true, 0, 11.5, true, 11.5, 23}, {true, 11.5, 23, true, 11.5, 23},
	{true, 11.5, 23, true, 0, 11.5}, {true, 0, 11.5, true, 0, 11.5},
}

func geometryCircleWall(grid *[2]int32, u *server.Object) int32 {
	s := GetServer().S()
	flags := byte(64)
	if u.ObjFlags&0x4000 != 0 && u.Shape.Kind == server.ShapeKindCircle && !(u.Shape.Circle.R > 9) {
		flags = 0
	}
	lookup := func(x, y int32) int8 { return s.Sub_57B500(image.Pt(int(x), int(y)), flags) }
	kind := uint8(lookup(grid[0], grid[1]))
	if kind == 255 {
		return 0
	}
	span := geometryCircleWallSpans[kind]
	mask := memmap.Uint8(0x587000, 292496+uintptr(kind))
	origin := types.Pointf{float32(23 * grid[0]), float32(23 * grid[1])}
	var point types.Pointf
	var lo, hi, result int32
	respond := func() {
		if u.ObjClass&4 != 0 && u.UpdateData != nil {
			var p unsafe.Pointer
			if wl := s.Walls.GetWallAtGrid(image.Pt(int(grid[0]), int(grid[1]))); wl != nil {
				p = wl.C()
			}
			*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 296)) = p
		}
		if geometryPointWall(u, &point) != 0 {
			result = 1
		}
	}
	if span.positive {
		if mask&2 != 0 && lookup(grid[0]-1, grid[1]-1) == -1 {
			lo = 1
		}
		if mask&4 != 0 && lookup(grid[0]+1, grid[1]+1) == -1 {
			hi = 1
		}
		if kind == 7 || kind == 10 {
			hi = 1
		} else if kind == 8 || kind == 9 {
			lo = 1
		}
		if geometryProjectPositive(&origin, span.lo, span.hi, lo, hi, &u.NewPos, &point) != 0 {
			respond()
		}
	}
	if span.negative {
		if mask&8 != 0 && lookup(grid[0]-1, grid[1]+1) == -1 {
			lo = 1
		}
		if mask&1 != 0 && lookup(grid[0]+1, grid[1]-1) == -1 {
			hi = 1
		}
		end := hi
		if kind == 7 || kind == 8 {
			lo = 1
		} else if kind == 9 || kind == 10 {
			end = 1
		}
		if geometryProjectNegative(&origin, span.otherLo, span.otherHi, lo, end, &u.NewPos, &point) != 0 {
			respond()
		}
	}
	return result
}
func geometryBoxWalls(u *server.Object) {
	if u.ObjClass&0x400000 != 0 {
		return
	}
	x1 := floatToInt32(float32(float64(u.CollideP1.X) * 0.043478262))
	y1 := floatToInt32(float32(float64(u.CollideP1.Y) * 0.043478262))
	x2 := floatToInt32(float32(float64(u.CollideP2.X) * 0.043478262))
	y2 := floatToInt32(float32(float64(u.CollideP2.Y) * 0.043478262))
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			grid := [2]int32{x, y}
			if geometryBoxWall(&grid, u) != 0 {
				collisionWallOpen(&grid, u)
			}
		}
	}
}
func geometryBoxWall(grid *[2]int32, u *server.Object) int32 {
	kind := uint8(GetServer().S().Sub_57B500(image.Pt(int(grid[0]), int(grid[1])), 64))
	if kind == 255 {
		return 0
	}
	data := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 292520+24*uintptr(kind))), 24)
	f := (*[6]float32)(unsafe.Pointer(&data[0]))
	center, previous := geometryRotated(u.NewPos), geometryRotated(u.PrevPos)
	bounds := geometryRotatedBounds(u, center)
	x, y := float64(23*grid[0]), float64(23*grid[1])
	origin := types.Pointf{float32((y + x) * 0.70710677), float32((y - x) * 0.70710677)}
	mask := uint8(geometryQuadrant(&center, &origin))
	var result int32
	if data[0] != 0 && mask&data[1] == 0 {
		edge := types.Pointf{float32(float64(origin.X) + 16.263456), float32(float64(origin.Y) - 16.263456 + float64(f[1]))}
		if geometryWallVertical(u, &center, &previous, &bounds, &edge, f[2]) != 0 {
			result = 1
		}
	}
	if data[12] != 0 && mask&data[13] == 0 {
		edge := origin
		edge.X = origin.X + f[4]
		if geometryWallHorizontal(u, &center, &previous, &bounds, &edge, f[5]) != 0 {
			result = 1
		}
	}
	return result
}
