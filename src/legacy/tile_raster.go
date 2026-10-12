package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

func tileRasterWrapSetup() {
	dword_5d4594_1193188 = 1
	stride := uint32(dword_5d4594_3798804)
	*memmap.PtrUint32(0x973CE0, 376) = 45*stride + (uint32(46) << memmap.Uint8(0x973F18, 7696))
}

// These two private fast paths always use the original 16-bit diamond layout.
func tileRasterCopyPacked(dst, source unsafe.Pointer) {
	stride := int(dword_5d4594_3798804)
	offset := 0
	for y := 0; y < 46; y++ {
		radius := min(y, 45-y)
		left, width := 23-radius, 2*radius+1
		copy(unsafe.Slice((*byte)(unsafe.Add(dst, y*stride+2*left)), 2*width), unsafe.Slice((*byte)(unsafe.Add(source, offset)), 2*width))
		offset += 2 * width
	}
}
func tileRasterFillPacked(dst unsafe.Pointer, color uint32) {
	stride := int(dword_5d4594_3798804)
	for y := 0; y < 46; y++ {
		radius := min(y, 45-y)
		left, width := 23-radius, 2*radius+1
		row := unsafe.Add(dst, y*stride+2*left)
		// The unrolled stores start at a relative word boundary after this halfword.
		if left&1 != 0 {
			*(*uint16)(row) = uint16(color)
			row = unsafe.Add(row, 2)
			width--
		}
		uiRenderFill(row, color, int32(2*width))
	}
}
func tileRasterOffset(pos image.Point) (base unsafe.Pointer, offset, size, stride int) {
	base = legacyGlobals.nox_video_tileBuf_ptr_3798796
	size = int(uintptr(legacyGlobals.nox_video_tileBuf_end_3798844) - uintptr(base))
	stride = int(dword_5d4594_3798804)
	x := uint32(dword_5d4594_3798836) + uint32(pos.X) - uint32(dword_5d4594_3798820)
	y := uint32(dword_5d4594_3798840) + uint32(pos.Y) - uint32(dword_5d4594_3798824)
	offset = int(uint32(stride)*y + (x << memmap.Uint8(0x973F18, 7696)))
	if offset >= size {
		offset -= size
	}
	return
}
func tileRasterTexture(pos image.Point, handle noxrender.ImageHandle) {
	base, offset, size, stride := tileRasterOffset(pos)
	data := GetClient().R2().GetBag().AsImage(handle).Pixdata()
	if len(data) == 0 {
		return
	}
	if offset+int(memmap.Uint32(0x973CE0, 376)) < size {
		tileBufMarkTile(offset)
		tileRasterCopyPacked(unsafe.Add(base, offset), unsafe.Pointer(&data[0]))
		return
	}
	shift := memmap.Uint8(0x973F18, 7696)
	for y := 0; y < 46; y++ {
		left := int(memmap.Uint32(0x973CE0, 192+4*uintptr(y)) << shift)
		n := int(memmap.Uint32(0x973CE0, 384+4*uintptr(y)) << shift)
		at := (offset + y*stride + left) % size
		first := min(n, size-at)
		tileBufMark(at, first)
		copy(unsafe.Slice((*byte)(unsafe.Add(base, at)), first), data[:first])
		if first < n {
			tileBufMark(0, n-first)
			copy(unsafe.Slice((*byte)(base), n-first), data[first:n])
		}
		data = data[n:]
	}
}
func tileRasterFill(pos image.Point, tile uint16) {
	def := &tileDefinitionsAll()[tile]
	if !Get_nox_client_texturedFloors_154956() && def.Field58&1 != 0 {
		dword_5d4594_1193156 = 1
	}
	color := def.Color48
	base, offset, size, stride := tileRasterOffset(pos)
	if offset+int(memmap.Uint32(0x973CE0, 376)) < size {
		tileBufMarkTile(offset)
		tileRasterFillPacked(unsafe.Add(base, offset), color)
		return
	}
	shift := memmap.Uint8(0x973F18, 7696)
	for y := 0; y < 46; y++ {
		left := int(memmap.Uint32(0x973CE0, 192+4*uintptr(y)) << shift)
		n := int(memmap.Uint32(0x973CE0, 384+4*uintptr(y)) << shift)
		at := (offset + y*stride + left) % size
		first := min(n, size-at)
		tileBufMark(at, first)
		uiRenderFill(unsafe.Add(base, at), color, int32(first))
		if first < n {
			tileBufMark(0, n-first)
			uiRenderFill(base, color, int32(n-first))
		}
	}
}
