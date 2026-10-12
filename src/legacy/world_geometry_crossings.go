package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"unsafe"
)

func geometryRectCrossings(line, rect *[4]float32, out *types.Pointf, capacity, extend int32) int32 {
	var count int32
	for edge := 0; edge < 4 && count < capacity; edge++ {
		p := (*types.Pointf)(unsafe.Add(unsafe.Pointer(out), uintptr(count)*8))
		var hit int32
		switch edge {
		case 0:
			hit = geometryCrossHorizontal(line, rect[0], rect[2], rect[1], p, extend)
		case 1:
			hit = geometryCrossVertical(line, rect[2], rect[1], rect[3], p, extend)
		case 2:
			hit = geometryCrossHorizontal(line, rect[0], rect[2], rect[3], p, extend)
		case 3:
			hit = geometryCrossVertical(line, rect[0], rect[1], rect[3], p, extend)
		}
		if hit != 0 {
			count++
		}
	}
	return count
}

// The original x87 code keeps the first ratio wide, but stores the second
// ratio as float32 before its bounds check and final coordinate calculation.
func geometryCrossHorizontal(line *[4]float32, lo, hi, y float32, out *types.Pointf, extend int32) int32 {
	denom := float64(line[1]) - float64(line[3])
	if denom == 0 {
		return 0
	}
	t := (float64(y) - float64(line[3])) / denom
	if t < 0 && extend == 0 || t > 1 {
		return 0
	}
	width := float64(lo) - float64(hi)
	if width == 0 {
		return 0
	}
	u := float32(((1-float64(t))*float64(line[2]) + t*float64(line[0]) - float64(hi)) / width)
	if u < 0 || u > 1 {
		return 0
	}
	out.Y = y
	out.X = float32((1-float64(u))*float64(hi) + float64(u)*float64(lo))
	return 1
}
func geometryCrossVertical(line *[4]float32, x, lo, hi float32, out *types.Pointf, extend int32) int32 {
	denom := float64(line[0]) - float64(line[2])
	if denom == 0 {
		return 0
	}
	t := (float64(x) - float64(line[2])) / denom
	if t < 0 && extend == 0 || t > 1 {
		return 0
	}
	height := float64(lo) - float64(hi)
	if height == 0 {
		return 0
	}
	u := float32(((1-float64(t))*float64(line[3]) + t*float64(line[1]) - float64(hi)) / height)
	if u < 0 || u > 1 {
		return 0
	}
	out.X = x
	out.Y = float32((1-float64(u))*float64(hi) + float64(u)*float64(lo))
	return 1
}
func geometryClipCenter(line, bounds, rect *[4]float32, out *types.Pointf) int32 {
	if bounds[0] < rect[0] || bounds[2] > rect[2] || bounds[1] < rect[1] || bounds[3] > rect[3] {
		var crossings [2]types.Pointf
		switch geometryRectCrossings(line, rect, &crossings[0], 2, 0) {
		case 1:
			which := 2
			if geometryRectFloat((*types.Pointf)(unsafe.Pointer(line)), rect) != 0 {
				which = 0
			}
			out.X = float32((float64(crossings[0].X) + float64(line[which])) * 0.5)
			out.Y = float32((float64(crossings[0].Y) + float64(line[which+1])) * 0.5)
		case 2:
			out.X = float32((float64(crossings[1].X) + float64(crossings[0].X)) * 0.5)
			out.Y = float32((float64(crossings[1].Y) + float64(crossings[0].Y)) * 0.5)
		default:
			return 0
		}
	} else {
		out.X = float32((float64(line[2]) + float64(line[0])) * 0.5)
		out.Y = float32((float64(line[3]) + float64(line[1])) * 0.5)
	}
	return 1
}
