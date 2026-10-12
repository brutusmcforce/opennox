package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"math"
)

func worldAbsScratch(v float32) float64 {
	bits := math.Float32bits(v)
	// The original x87 load/store quiets signaling NaNs before the mask.
	if bits&0x7f800000 == 0x7f800000 && bits&0x007fffff != 0 {
		bits |= 0x00400000
	}
	*memmap.PtrUint32(0x5D4594, 527672) = bits
	*memmap.PtrUint32(0x5D4594, 527672) &= 0x7fffffff
	return float64(memmap.Float32(0x5D4594, 527672))
}

func worldRoundScratch(v float32) uint32 {
	if v < 0 {
		return 0
	}
	*memmap.PtrFloat32(0x5D4594, 527676) = float32(float64(v) + 8388608)
	result := memmap.Uint32(0x5D4594, 527676) & 0x7fffff
	*memmap.PtrUint32(0x5D4594, 527680) = result
	return result
}

func worldInitTrigTables() int8 {
	flag := memmap.PtrUint8(0x5D4594, 371240)
	if *flag != 0 {
		return int8(*flag)
	}
	*flag = 1
	for i := 0; i < 4096; i++ {
		value := int64(math.Sin(float64(i)*0.0015339808) * memmap.Float64(0x581450, 7200))
		*memmap.PtrUint32(0x85B3FC, 12260+uintptr(i)*4) = uint32(value)
	}
	var last int64
	for i := 0; i < 8192; i++ {
		last = int64(math.Acos(float64(i)*0.00024414062-1) * memmap.Float64(0x581450, 7184))
		*memmap.PtrUint32(0x5D4594, 338472+uintptr(i)*4) = uint32(last)
	}
	return int8(last)
}
