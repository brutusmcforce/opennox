package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
)

func tileAtPoint(p types.Pointf) int32 {
	// Match C's PC53 expressions and explicit float32 spills before conversion.
	x, y := float64(p.X)+11.5, float64(p.Y)+11.5
	i := floatToInt32(float32(x * 0.021739131))
	j := floatToInt32(float32(y * 0.021739131))
	if i <= 1 || i >= 127 || j <= 1 || j >= 127 {
		return -1
	}
	u, v := floatToInt32(float32(x))%46, floatToInt32(float32(y))%46
	rows := (*[128]*[128]worldTileCell)(unsafe.Pointer(worldTileGrid))
	var fallback int32
	var head unsafe.Pointer
	var local [2]int32
	if u <= v {
		if 46-u <= v {
			cell := &rows[i][j]
			fallback, head = int32(cell[6]), memmap.WordPtr(cell[10])
			local = [2]int32{u, v - 23}
		} else {
			cell := &rows[i-1][j]
			fallback, head = int32(cell[1]), memmap.WordPtr(cell[5])
			local = [2]int32{u + 23, v}
		}
	} else if 46-u <= v {
		cell := &rows[i][j]
		fallback, head = int32(cell[1]), memmap.WordPtr(cell[5])
		local = [2]int32{u - 23, v}
	} else {
		cell := &rows[i][j-1]
		fallback, head = int32(cell[6]), memmap.WordPtr(cell[10])
		local = [2]int32{u, v + 23}
	}
	if head == nil {
		return fallback
	}
	return findSubtileAt((*[5]uint32)(head), &local, fallback)
}
