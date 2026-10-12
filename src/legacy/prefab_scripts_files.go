package legacy

import (
	"io"
	"os"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func prefabScriptRewrite(path string, offset [2]int32) uint32 {
	instance := int32(*prefabGlobal(prefabInstance))
	prefabScriptObjectNames(instance, offset[0], offset[1])
	backup := alloc.GoString((*byte)(memmap.PtrOff(0x973F18, 42152))) + alloc.GoString((*byte)(memmap.PtrOff(0x587000, 282608)))
	if ifs.Copy(path, backup) != nil {
		return 0
	}
	input, err := ifs.Open(backup)
	if err != nil {
		return 0
	}
	defer input.Close()
	output, err := ifs.OpenFile(path, os.O_RDWR)
	if err != nil {
		return 0
	}
	defer output.Close()
	in, out := binfile.NewFile(input), binfile.NewFile(output)
	if !prefabScriptCopyBytes(in, out, 12) {
		return 0
	}
	n := prefabScriptCopyWord(in, out)
	for i := int32(0); i < int32(n); i++ {
		if !prefabScriptCopyString(in, out) {
			return 0
		}
	}
	prefabScriptCopyWord(in, out)
	n = prefabScriptCopyWord(in, out)
	for i := int32(0); i < int32(n); i++ {
		prefabScriptCopyWord(in, out)
		size := prefabScriptReadInt(in)
		if size >= 2048 {
			return 0
		}
		name := make([]byte, int(size))
		if _, err := io.ReadFull(in, name); err != nil {
			return 0
		}
		if i > 1 {
			name = []byte(prefabScriptName(string(name), instance, offset[0], offset[1], false))
		}
		prefabScriptWriteInt(out, uint32(len(name)))
		out.Write(name)
		prefabScriptCopyWord(in, out)
		prefabScriptCopyWord(in, out)
		prefabScriptLocals(in, out)
		prefabScriptCopyWord(in, out)
		start, _ := out.Seek(0, io.SeekCurrent)
		prefabScriptCopyWord(in, out)
		prefabScriptInstructions(in, out, false)
		end, _ := out.Seek(0, io.SeekCurrent)
		out.Seek(start, io.SeekStart)
		prefabScriptWriteInt(out, uint32(end-start-4))
		out.Seek(end, io.SeekStart)
	}
	if !prefabScriptCopyBytes(in, out, 4) {
		return 0
	}
	input.Close()
	output.Close()
	return uint32(bool2int(ifs.Remove(backup) == nil))
}
func prefabScriptGeneration(dir string) uint32 {
	// Nil/empty directory was undefined in C; avoid constructing host-root paths.
	if dir == "" || len(dir) > 2035 {
		return 0
	}
	*memmap.PtrUint64(0x5D4594, 1549772) = uint64(uint32(PlatformTicks()))
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 35912)), 72))
	for i, v := range []uint32{0, 0, 255, 0, 1, 1, 1, 0, 1} {
		*mapPaintGlobal(i) = v
	}
	*prefabGlobal(prefabSelected) = ^uint32(0)
	*memmap.PtrUint8(0x973F18, 35972) = 2
	prefabResetWaypoint()
	copyPath := func(off uintptr, s string) {
		p := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, off)), 2048)
		copy(p, s)
		p[len(s)] = 0
	}
	copyPath(42152, dir)
	oldPath := dir + alloc.GoString((*byte)(memmap.PtrOff(0x587000, 197556)))
	blendPath := dir + alloc.GoString((*byte)(memmap.PtrOff(0x587000, 197568)))
	copyPath(36008, oldPath)
	copyPath(38056, blendPath)
	ifs.Remove(oldPath)
	ifs.Remove(blendPath)
	GetServer().Nox_xxx_mapReset5028E0()
	library := dir + alloc.GoString((*byte)(memmap.PtrOff(0x587000, 197580)))
	p, free := alloc.CString(library)
	defer free()
	prefabSetPath(mapRoomRaw(unsafe.Pointer(p)), false)
	prefabSetPath(mapRoomRaw(unsafe.Pointer(p)), true)
	ret := prefabLibrary()
	*prefabGlobal(prefabInstance) = 0
	*memmap.PtrUint32(0x973F18, 35880) = 0
	*memmap.PtrUint32(0x5D4594, 1599580) = 0
	return ret
}
func prefabScriptBounds(a *[8]uint32) {
	// Comparisons follow the C expression types, including its mixed signedness.
	x0, y0 := a[0], a[1]
	y := int32(a[1])
	if int32(a[3]) < y {
		x0, y0, y = a[2], a[3], int32(a[3])
	}
	if a[5] < uint32(y) {
		x0, y0, y = a[4], a[5], int32(a[5])
	}
	if a[7] < uint32(y) {
		x0, y0 = a[6], a[7]
	}
	x1, y1 := a[2], a[3]
	if a[0] < x1 {
		x1, y1 = a[0], a[1]
	}
	if a[4] < x1 {
		x1, y1 = a[4], a[5]
	}
	if int32(a[6]) < int32(x1) {
		x1, y1 = a[6], a[7]
	}
	x2, y2 := a[4], a[5]
	if a[0] > x2 {
		x2, y2 = a[0], a[1]
	}
	if a[2] > x2 {
		x2, y2 = a[2], a[3]
	}
	if int32(a[6]) > int32(x2) {
		x2, y2 = a[6], a[7]
	}
	x3, y3 := a[6], a[7]
	if a[1] > y3 {
		x3, y3 = a[0], a[1]
	}
	if int32(a[3]) > int32(y3) {
		x3, y3 = a[2], a[3]
	}
	if a[5] > y3 {
		x3, y3 = a[4], a[5]
	}
	*a = [8]uint32{x0, y0, x1, y1, x2, y2, x3, y3}
}
