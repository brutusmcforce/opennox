package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func prefabGlobal(i int) *uint32 { return &prefabState[i] }
func prefabCompareNames(a, b uint32) int {
	return textCompareNarrow((*byte)(mapRoomPointer(a)), (*byte)(mapRoomPointer(b)))
}
func prefabTileDefinitionCount() int    { return int(worldTileDefinitionCount) }
func prefabSetBounds(bounds *[8]uint32) { prefabScriptBounds(bounds) }
func prefabInteresting(bounds *[4]int32, index uint32) uint32 {
	return prefabScriptPending(bounds, index)
}
func prefabAdjustScript(instance uint32, dx, dy int32) {
	prefabScriptObjectNames(int32(instance), dx, dy)
	offset := [2]int32{dx, dy}
	prefabScriptRewrite(alloc.GoString((*byte)(memmap.PtrOff(0x973F18, 30760))), offset)
}
func prefabCombineScripts(a, b, c string) {
	x, freeX := alloc.CString(a)
	defer freeX()
	y, freeY := alloc.CString(b)
	defer freeY()
	z, freeZ := alloc.CString(c)
	defer freeZ()
	Nox_script_readWriteZzz_541670(x, y, z)
}
func prefabTransferWallData(src, dst *server.Wall) {
	if dst.Flags4&4 != 0 {
		worldSecretRemove(dst.Data)
		dst.Data = nil
		dst.Flags4 &^= 4
	}
	if src.Flags4&12 != 0 {
		dst.Field10 = src.Field10
	}
	if src.Flags4&4 != 0 && src.Data != nil {
		dst.Flags4 |= 4
		dst.Data = src.Data
		src.Data = nil
		data := (*[8]uint32)(dst.Data)
		data[1], data[2], data[3] = uint32(dst.X5), uint32(dst.Y6), mapRoomRaw(dst.C())
		worldSecretInsert(dst.Data)
	}
	if src.Flags4&8 != 0 && dst.Flags4&8 == 0 {
		dst.Flags4 |= 8
		GetServer().S().Walls.AddBreakable(dst)
	}
}
