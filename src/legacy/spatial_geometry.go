package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"image"
)

func spatialQuadrant(p *types.Pointf) int8 {
	x, y := uint8(floatToInt32(p.X)%23), uint8(floatToInt32(p.Y)%23)
	if x >= 12 {
		if y >= 12 {
			return 4
		}
		return 1
	}
	if y >= 12 {
		return 8
	}
	return 2
}
func spatialTriangle(p *types.Pointf, grid [2]int32) int32 {
	x := floatToInt32(float32(float64(p.X) - float64(grid[0]*23)))
	y := floatToInt32(float32(float64(p.Y) - float64(grid[1]*23)))
	// The C byte assignment intentionally retains the upper 24 bits.
	b := int32(0)
	if 22-x <= y {
		b = 1
	}
	v := (y & ^int32(255)) | b
	if x <= y {
		return ((v - 1) &^ int32(2)) + 3
	}
	return v + 1
}
func spatialNormal(grid *[2]int32, ray *[4]float32, normal *types.Pointf) int32 {
	start, end := types.Pointf{ray[0], ray[1]}, types.Pointf{ray[2], ray[3]}
	tri, quad := int8(spatialTriangle(&start, *grid)), spatialQuadrant(&end)
	// Four diagonal outward normals, in clockwise order from northwest.
	normals := [4]types.Pointf{{-.70709997, -.70709997}, {.70709997, -.70709997}, {.70709997, .70709997}, {-.70709997, .70709997}}
	face := func(t int8) int {
		switch t {
		case 0:
			if quad&2 != 0 {
				return 3
			}
			return 0
		case 1:
			if quad&1 != 0 {
				return 0
			}
			return 1
		case 2:
			if quad&1 != 0 {
				return 2
			}
			return 1
		default:
			if quad&4 != 0 {
				return 3
			}
			return 2
		}
	}
	n := 0
	switch GetServer().S().Sub_57B500(image.Pt(int(grid[0]), int(grid[1])), 64) {
	case 0:
		if tri != 0 && tri != 1 {
			n = 2
		}
	case 1:
		n = 3
		if tri == 1 || tri == 2 {
			n = 1
		}
	case 2:
		if tri < 0 || tri > 3 {
			return 1
		}
		n = face(tri)
	case 3:
		n = 2
		if tri == 0 || tri == 1 {
			n = face(tri)
		}
	case 4:
		n = 3
		if tri == 1 || tri == 2 {
			n = face(tri)
		}
	case 5:
		n = 0
		if tri == 2 || tri == 3 {
			n = face(tri)
		}
	case 6:
		n = 1
		if tri == 0 || tri == 3 {
			n = face(tri)
		}
	case 7:
		n = 3
		if quad&1 != 0 {
			n = 2
		}
		if tri == 1 {
			n = face(tri)
		}
	case 8:
		n = 3
		if quad&1 != 0 {
			n = 0
		}
		if tri == 2 {
			n = face(tri)
		}
	case 9:
		n = 0
		if quad&4 != 0 {
			n = 1
		}
		if tri == 3 {
			n = face(tri)
		}
	case 10:
		n = 2
		if quad&2 != 0 {
			n = 1
		}
		if tri == 0 {
			n = face(tri)
		}
	default:
		return 0
	}
	*normal = normals[n]
	return 1
}
