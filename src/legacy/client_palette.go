package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"sort"
	"unsafe"
)

func clientPaletteSort() {
	source := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 3880)), 1024)
	var entries [256]uint32
	for i := range entries {
		entries[i] = uint32(i) | uint32(source[4*i])<<8 | uint32(source[4*i+1])<<16 | uint32(source[4*i+2])<<24
	}
	sort.Slice(entries[:], func(i, j int) bool { return entries[i] < entries[j] })
	target := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 809604)), 1024)
	indices := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 808304)), 256)
	for i, v := range entries {
		target[4*i], target[4*i+1], target[4*i+2], target[4*i+3] = byte(v>>8), byte(v>>16), byte(v>>24), 0
		indices[i] = byte(v)
	}
}
func clientPaletteExpand(dst, src unsafe.Pointer) {
	d, s := unsafe.Slice((*byte)(dst), 1024), unsafe.Slice((*byte)(src), 768)
	for i := 0; i < 256; i++ {
		d[4*i] = s[3*i]
		d[4*i+1] = s[3*i+1]
		v := s[3*i+2]
		d[4*i+3] = 4
		d[4*i+2] = v
	}
}
func clientPaletteCompact(dst, src unsafe.Pointer) {
	d, s := unsafe.Slice((*byte)(dst), 768), unsafe.Slice((*byte)(src), 1024)
	for i := 0; i < 256; i++ {
		d[3*i] = s[4*i]
		d[3*i+1] = s[4*i+1]
		d[3*i+2] = s[4*i+2]
	}
}
