package legacy

import (
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"unsafe"
)

func mapSectionIO(p unsafe.Pointer, n int) uint32 {
	_, err := cryptfile.Global().ReadWrite(unsafe.Slice((*byte)(p), n))
	if err != nil {
		return 0
	}
	return 1
}
func mapSectionByte(v byte) byte      { mapSectionIO(unsafe.Pointer(&v), 1); return v }
func mapSectionWord(v uint16) uint16  { mapSectionIO(unsafe.Pointer(&v), 2); return v }
func mapSectionDword(v uint32) uint32 { mapSectionIO(unsafe.Pointer(&v), 4); return v }

func mapSectionTile(p *[5]uint32) uint32 {
	p[0] = uint32(mapSectionByte(byte(p[0])))
	p[1] = uint32(int32(int16(mapSectionWord(uint16(p[1])))))
	p[2] = uint32(mapSectionByte(byte(p[2])))
	p[3] = uint32(mapSectionByte(byte(p[3])))
	var count byte
	for n := mapPaintNode(p[4]); n != nil; n = mapPaintNode(n[4]) {
		count++
	}
	count = mapSectionByte(count)
	var ret uint32
	if cryptfile.Global().ReadOnly() {
		p[4] = 0
		tail := p
		ret = uint32(count)
		carry := byte(p[3])
		variation := uint16(p[1])
		for i := 0; i < int(count); i++ {
			n := mapPaintSubtileNew(0, 0, 0, 0)
			carry = mapSectionByte(carry)
			n[0] = uint32(carry)
			variation = mapSectionWord(variation)
			n[1] = uint32(int32(int16(variation)))
			carry = mapSectionByte(carry)
			n[2] = uint32(carry)
			ret = mapSectionIO(unsafe.Pointer(&carry), 1)
			n[3] = uint32(carry)
			tail[4] = mapRoomRaw(unsafe.Pointer(n))
			tail = n
		}
	} else {
		for n := mapPaintNode(p[4]); n != nil; n = mapPaintNode(n[4]) {
			mapSectionByte(byte(n[0]))
			mapSectionWord(uint16(n[1]))
			mapSectionByte(byte(n[2]))
			v := byte(n[3])
			ret = mapSectionIO(unsafe.Pointer(&v), 1)
		}
	}
	return uint32(byte(ret))
}
