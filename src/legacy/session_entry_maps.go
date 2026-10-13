package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"strings"
	"unsafe"
)

func sessionModeIndex(value int16) int32 {
	mask := uint32(uint16(value) & 0x17f0)
	for i := uintptr(0); i < 6; i++ {
		if memmap.Uint32(0x587000, 4704+4*i) == mask {
			return int32(i)
		}
	}
	return 0
}
func sessionMapFilename() *byte { return memmap.PtrUint8(0x5D4594, 2598188) }
func sessionMapName() *byte     { return memmap.PtrUint8(0x85B3FC, 36) }
func sessionSelectedMap() *byte { return memmap.PtrUint8(0x5D4594, 3452) }
func sessionSetSelectedMap(name *byte) uint32 {
	text := alloc.GoString(name)
	dst := unsafe.Slice(sessionSelectedMap(), len(text)+1)
	copy(dst, text)
	dst[len(text)] = 0
	return uint32(len(dst))
}
func sessionSetMapPath(name *byte) *byte {
	text := alloc.GoString(name)
	if mapASCIIEqual(alloc.GoString(sessionMapFilename()), text) {
		return nil
	}
	dst := unsafe.Slice(sessionMapFilename(), 80)
	clear(dst)
	copy(dst, text)
	dst[79] = 0
	base := alloc.GoString(sessionMapFilename())
	if i := strings.LastIndexByte(text, '\\'); i >= 0 {
		base = text[i+1:]
	}
	n := len(base) - 4
	if n < 0 {
		n = 0
	}
	out := unsafe.Slice(sessionMapName(), n+1)
	copy(out[:n], base)
	out[n] = 0
	serverConfigUpdatedSet()
	return sessionMapName()
}
func sessionMapFlags(it *Nox_map_list_item) int32 {
	return int32(Nox_mapToGameFlags(int(int32(it.Field_7))))
}
func sessionMapStateSet(value int32) int32 {
	*memmap.PtrUint32(0x5D4594, 1523072) = uint32(value)
	return value
}
func sessionMapState() int32 { return int32(memmap.Uint32(0x5D4594, 1523072)) }
func sessionLoadStateSet(value int32) int32 {
	*memmap.PtrUint32(0x5D4594, 1563072) = uint32(value)
	return value
}
func sessionScavengerReset()         { *memmap.PtrUint32(0x5D4594, 1548508) = 0 }
func sessionScavengerMaximum() int32 { return int32(memmap.Uint32(0x5D4594, 1548528)) }
func sessionScavengerMaximumReset()  { *memmap.PtrUint32(0x5D4594, 1548528) = 0 }
