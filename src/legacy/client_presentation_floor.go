package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func presentationLocal() *client.Drawable { return LocalPlayerDrawable }
func presentationWallY(wall *[8]byte) int32 {
	y := 23 * int32(wall[6])
	local := presentationLocal()
	if local == nil {
		return y + 11
	}
	var dx int32
	x := 23 * int32(wall[5])
	switch wall[0] {
	case 0, 3, 11:
		dx = -23
		x += 22
	case 1, 4, 12:
		dx = 23
	default:
		return y + 11
	}
	side := dx*(int32(local.PosVec.Y)-y) - 23*(int32(local.PosVec.X)-x)
	if dx < 0 {
		side = -side
	}
	if side < 0 {
		y += 22
	}
	return y
}
func presentationDrawableY(dr *client.Drawable) int32 {
	offset := 8 * uintptr(dr.Field_74_4)
	dx, dy := memmap.Int32(0x587000, 196184+offset), memmap.Int32(0x587000, 196188+offset)
	y := int32(dr.PosVec.Y)
	local := presentationLocal()
	if local == nil {
		return y + dy/2
	}
	side := (int32(local.PosVec.Y)-y)*dx - (int32(local.PosVec.X)-int32(dr.PosVec.X))*dy
	if dx < 0 {
		side = -side
	}
	end := y + dy
	if side < 0 {
		return max(y, end)
	}
	return min(y, end)
}
func presentationCopy(dst, src unsafe.Pointer, n uint32) {
	// Production sprite runs always contain whole 16-bit pixels.
	copy(unsafe.Slice((*byte)(dst), int(n)), unsafe.Slice((*byte)(src), int(n)))
}
func presentationBake(_ *noxrender.Viewport, dr *client.Drawable) {
	var handle noxrender.ImageHandle
	if dr.DrawFuncPtr == drawableDrawKey(drawKey_nox_thing_static_draw) {
		if dr.ObjClass&0x40000 != 0 && dr.ObjFlags&0x1000000 == 0 {
			return
		}
		handle = *(*noxrender.ImageHandle)(unsafe.Add(dr.DrawData, 4))
	} else {
		frames := *(*unsafe.Pointer)(unsafe.Add(dr.DrawData, 4))
		handle = *(*noxrender.ImageHandle)(unsafe.Add(frames, 4*uintptr(dr.AnimFrameSlave)))
	}
	data := GetClient().R2().GetBag().AsImage(handle).Pixdata()
	if len(data) == 0 {
		return
	}
	width, height := int32(binary.LittleEndian.Uint32(data)), int32(binary.LittleEndian.Uint32(data[4:]))
	x := int32(binary.LittleEndian.Uint32(data[8:])) + int32(dr.PosVec.X) - int32(byte(dr.Field_0))
	y := int32(binary.LittleEndian.Uint32(data[12:])) + int32(dr.PosVec.Y) - int32(int16(dr.ZVal2)) - int32(int16(dr.ZVal)) - int32(byte(dr.Field_0>>8))
	originX, originY := int32(dword_5d4594_3798820), int32(dword_5d4594_3798824)
	if x < originX || x+width >= originX+int32(dword_5d4594_3798800) || y < originY || y+height >= originY+int32(dword_5d4594_3798808) {
		dr.Field_86 = 0
		return
	}
	counter := int32(nox_xxx_waypointCounterMB_587000_154948)
	if counter <= 0 {
		dr.Field_86 = 0
	}
	if counter-int32(dr.Field_86) <= 1 && counter > 0 {
		dr.Field_86 = uint32(counter)
		return
	}
	base := legacyGlobals.nox_video_tileBuf_ptr_3798796
	size := int(uintptr(legacyGlobals.nox_video_tileBuf_end_3798844) - uintptr(base))
	stride := int(dword_5d4594_3798804)
	start := int(uint32(stride)*(uint32(y)+uint32(dword_5d4594_3798840)-uint32(originY)) + 2*(uint32(x)+uint32(dword_5d4594_3798836)-uint32(originX)))
	if start >= size {
		start -= size
	}
	src := 17
	for row := int32(0); row < height; row++ {
		dst := start
		for remaining := width; remaining > 0; {
			op, count := data[src]&15, int(data[src+1])
			src += 2
			length := 2 * count
			switch op {
			case 1:
				dst += length
				if dst >= size {
					dst -= size
				}
			case 2:
				if dst+length < size {
					tileBufMark(dst, length)
					presentationCopy(unsafe.Add(base, dst), unsafe.Pointer(&data[src]), uint32(length))
					dst += length
				} else {
					head := size - dst
					tileBufMark(dst, head)
					presentationCopy(unsafe.Add(base, dst), unsafe.Pointer(&data[src]), uint32(head))
					if length > head {
						tileBufMark(0, length-head)
						presentationCopy(base, unsafe.Pointer(&data[src+head]), uint32(length-head))
					}
					dst = length - head
				}
				src += length
			}
			remaining -= int32(count)
		}
		start += stride
		if start >= size {
			start -= size
		}
	}
}
