package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

type tileFillEntry struct {
	X, Y, Flags uint32
}

func tileFillStack() *[500]tileFillEntry {
	return (*[500]tileFillEntry)(memmap.PtrOff(0x973F18, 16200))
}

func pushTileFill(x, y, flags, key int32) {
	if x <= 0 || x >= 127 || y <= 0 || y >= 127 || flags&3 == 0 {
		return
	}
	rows := (*[128]*[128]worldTileCell)(unsafe.Pointer(worldTileGrid))
	cell := &rows[x][y]
	if !(flags&2 != 0 && int32(cell[6]) == key || flags&1 != 0 && int32(cell[1]) == key) {
		return
	}
	if flags&1 != 0 && y == 1 || flags&2 != 0 && x == 1 {
		return
	}
	count := uint32(dword_5d4594_2487248)
	stack := tileFillStack()
	// C scans using a signed count but tests capacity using the unsigned word.
	for i := int32(0); i < int32(count); i++ {
		p := &stack[i]
		if p.X == uint32(x) && p.Y == uint32(y) && p.Flags == uint32(flags) {
			return
		}
	}
	if count >= 500 {
		*memmap.PtrUint32(0x973F18, 22200) = 1
		return
	}
	dword_5d4594_2487248 = uint32(count + 1)
	stack[count] = tileFillEntry{X: uint32(x), Y: uint32(y), Flags: uint32(flags)}
}

func popTileFill(x, y, flags *uint32) bool {
	count := uint32(dword_5d4594_2487248)
	if int32(count) <= 0 {
		return false
	}
	count--
	dword_5d4594_2487248 = uint32(count)
	stack := tileFillStack()
	*x = stack[count].X
	// Output pointers can alias the count or stack. Reload after each write,
	// as C does, instead of snapshotting a triplet before writing outputs.
	*y = stack[uint32(dword_5d4594_2487248)].Y
	*flags = stack[uint32(dword_5d4594_2487248)].Flags
	return true
}
