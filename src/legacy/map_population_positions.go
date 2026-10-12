package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"math"
	"unsafe"
)

// These named C globals are separate from the same-offset blob words. Preserve
// that distinction even though the original layout suggests a contiguous table.
func populationPositionWord(index int) *uint32 {
	if index >= 2 && index <= 5 {
		return populationGlobal(index + 4)
	}
	return populationBlob(uintptr(2487612 + 4*index))
}
func mapPopulationPrefabPositions(bits float32) {
	cfg := math.Float32bits(bits)
	half := float32(float64(*populationFloat(cfg, 64)) * 0.5)
	*populationBlob(2487608) = 0
	n := int32(*populationWord(cfg, 84))
	if n >= 5 {
		n = 5
	}
	*populationGlobal(10) = uint32(n)
	if n == 1 {
		*populationBlob(2487588) = 0
	} else if n == 2 {
		if mapRoomRandomInt(0, 100) >= 50 {
			*populationBlob(2487588) = 1
			*populationBlob(2487592) = 0
		} else {
			*populationBlob(2487588) = 0
			*populationBlob(2487592) = 1
		}
	} else {
		for i := int32(0); i < n; i++ {
			*populationBlob(uintptr(2487588 + 4*i)) = uint32(i)
		}
		for i := int32(0); i < n; i++ {
			j := mapRoomRandomInt(0, n-2)
			a, b := populationBlob(uintptr(2487588+4*j)), populationBlob(uintptr(2487592+4*j))
			*a, *b = *b, *a
		}
	}
	set := func(i int, x float32) { *populationPositionWord(i) = math.Float32bits(x) }
	random := func() float32 { return float32(mapRoomRandomFloat(-half, half)) }
	switch n {
	case 1:
		set(0, random())
		set(1, random())
	case 2:
		if mapRoomRandomInt(0, 100) >= 50 {
			set(0, -half)
			set(1, random())
			set(2, half)
			set(3, random())
		} else {
			set(0, random())
			set(1, -half)
			set(2, random())
			set(3, half)
		}
	case 3:
		roll := mapRoomRandomInt(0, 100)
		switch {
		case roll < 25:
			set(0, -half)
			set(1, -half)
			set(3, -half)
			set(2, half)
			set(4, 0)
			set(5, half)
		case roll < 50:
			set(0, -half)
			set(5, -half)
			set(1, half)
			set(2, half)
			set(3, half)
			set(4, 0)
		case roll < 75:
			set(1, -half)
			set(4, -half)
			set(0, half)
			set(2, half)
			set(3, half)
			set(5, 0)
		default:
			set(0, -half)
			set(1, -half)
			set(2, -half)
			set(3, half)
			set(4, half)
			set(5, 0)
		}
	case 4, 5:
		*populationBlob(2487644) = 0
		set(0, -half)
		set(1, -half)
		set(3, -half)
		set(4, -half)
		set(2, half)
		set(5, half)
		set(6, half)
		set(7, half)
		*populationBlob(2487648) = 0
	}
}
func populationNextPosition(r *mapRoom, out *types.Pointf) int64 {
	index := *populationBlob(uintptr(2487588 + 4**populationBlob(2487608)))
	out.X = math.Float32frombits(*populationBlob(uintptr(2487612 + 8*index)))
	out.Y = math.Float32frombits(*populationBlob(uintptr(2487616 + 8*index)))
	*populationBlob(2487608)++
	out.X = float32(float64(out.X) - float64(r.Size.X)*0.5)
	out.Y = float32(float64(out.Y) - float64(r.Size.Y)*0.5)
	out.X = float32(float64(int32(int64(float64(out.X)*0.030743772))) * 32.526913)
	ret := int64(float64(out.Y) * 0.030743772)
	out.Y = float32(float64(int32(ret)) * 32.526913)
	return ret
}
func mapPopulationNextPosition(room, out uint32) int64 {
	return populationNextPosition(populationRoom(room), (*types.Pointf)(unsafe.Pointer(mapRoomPointer(out))))
}
